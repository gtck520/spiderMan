package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gtck520/spiderMan/helper"
)

type Colly struct {
	C        *colly.Collector
	Cdetail  *colly.Collector
	num      int //采集次数
	WriteNum int //写入次数
}

//初始化建立colly实例
func (co *Colly) BuildC(Url string) {

	//域名自动补充 http
	httpurl := Url
	if !strings.HasPrefix(Url, "http://") && !strings.HasPrefix(Url, "https://") {
		httpurl = "http://" + Url
	}

	//httpurl = "https://feed.sina.com.cn/api/roll/get?pageid=121&lid=1356&num=20&versionNumber=1.2.4&page=2&encode=utf-8&callback=feedCardJsonpCallback&_=1632886493594"
	url1 := strings.Split(httpurl, "//")[1]
	url2 := strings.Split(url1, "/")[0]
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains(url2),
	)
	//设置客户端，模拟浏览器访问
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36"
	//拷贝一份实例 用于访问详情链接
	co.Cdetail = c.Clone()
	co.C = c
	co.num = 0
	co.WriteNum = 0
}

//根据规则抓取数据
func (co *Colly) GetContent(Rule Rule) {
	// OnRequest 请求执行之前调用
	// OnResponse 响应返回之后调用
	// OnHTML 监听执行 selector
	// OnXML 监听执行 selector
	// OnHTMLDetach，取消监听，参数为 selector 字符串
	// OnXMLDetach，取消监听，参数为 selector 字符串
	// OnScraped，完成抓取后执行，完成所有工作后执行
	// OnError，错误回调
	// On every a element which has href attribute call callback
	if Rule.Type == 2 {
		co.C.OnResponse(func(r *colly.Response) {
			str := string(r.Body)
			//解析正则表达式，如果成功返回解释器
			reg1 := regexp.MustCompile(Rule.Match)
			if reg1 == nil { //解释失败，返回nil
				fmt.Println("regexp err")
				return
			}
			//detailLink.Visit(`https:\/\/news.sina.cn\/gn\/2021-09-29\/detail-iktzscyx6975697.d.html`)
			//根据规则提取关键信息
			result1 := reg1.FindAllStringSubmatch(str, -1)
			if len(result1) > 0 {
				if len(Rule.SubRule) > 0 {
					for _, v := range result1 {
						fmt.Printf("Link found: %s\n", v[1])
						childurl := strings.Replace(v[1], "\\", "", -1)

						//子链接加入允许访问
						url1 := strings.Split(childurl, "//")[1]
						url2 := strings.Split(url1, "/")[0]
						co.Cdetail.AllowedDomains = []string{url2}

						err := co.Cdetail.Visit(childurl)
						if err != nil {
							fmt.Println(childurl+" visit err:", err.Error())
						}
					}
				}

			}
		})
	} else {
		co.C.OnHTML(Rule.Match, func(e *colly.HTMLElement) {
			link := e.Attr("href")
			// Print link
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			if len(Rule.SubRule) > 0 {
				co.Cdetail.Visit(e.Request.AbsoluteURL(link))
			}
			// Visit link found on page
			// Only those links are visited which are in AllowedDomains
			//C.Visit(e.Request.AbsoluteURL(link))
		})
	}

}

//根据规则提取分页数据
func (co *Colly) GetPageContent(Rule PageRule) {
	if Rule.Type == 1 {
		co.C.OnHTML(Rule.Match, func(e *colly.HTMLElement) {
			link := e.Attr("href")
			fmt.Printf("pageLink found: %q -> %s\n", e.Text, link)
			page, _ := strconv.Atoi(e.Text)
			if Rule.Page > 0 {
				if Rule.Page == page {
					co.C.Visit(e.Request.AbsoluteURL(link))
				}
			} else {
				if co.num < Rule.Num {
					co.C.Visit(e.Request.AbsoluteURL(link))
					co.num++
				}
			}
		})
	} else if Rule.Type == 2 {
		if Rule.Page > 0 {
			link := strings.Replace(Rule.Match, "{page}", strconv.Itoa(Rule.Page), -1)
			co.C.Visit(link)
		} else {
			for i := 1; i < Rule.Num; i++ {
				link := strings.Replace(Rule.Match, "{page}", strconv.Itoa(i), -1)
				co.C.Visit(link)
			}
		}
	} else {
		return
	}
}

//根据规则抓取数据
func (co *Colly) GetDContent(Rule Rule, SubRule []SubRule, Output int) {
	co.Cdetail.OnHTML(Rule.SubMatch, func(e *colly.HTMLElement) {
		for _, subrule := range SubRule {
			var content string
			if subrule.Type == 2 {
				str := string(e.Text)
				//解析正则表达式，如果成功返回解释器
				reg1 := regexp.MustCompile(subrule.Match)
				if reg1 == nil { //解释失败，返回nil
					fmt.Println("regexp err")
					return
				}
				//detailLink.Visit(`https:\/\/news.sina.cn\/gn\/2021-09-29\/detail-iktzscyx6975697.d.html`)
				//根据规则提取关键信息
				result1 := reg1.FindAllStringSubmatch(str, -1)
				content = result1[0][1]
				fmt.Printf("content found: %s\n", result1)

			} else {
				content = e.ChildText(subrule.Match)
				fmt.Printf("detial %s : %s \t", subrule.Name, content)
			}
			detaildata := make(map[string]string)
			subdata := make(map[string]map[string]string)
			detaildata[subrule.Field] = content
			subdata[e.Request.URL.String()] = detaildata
			if Output == 1 {
				if helper.CsvWrite(subdata, co.WriteNum) {
					co.WriteNum++
				}
			} else {

			}

		}
	})

}
