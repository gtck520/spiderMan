package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gtck520/spiderMan/helper"
	"github.com/spf13/viper"
)

type Spider struct {
	UrlList map[string]string
}

//url
type Urls struct {
	Name string
	Url  string
}

//定义每个url的配置结构
type UrlConfig struct {
	Name  string //配置名称
	Url   string //url
	Rules []Rule //规则
}

//规则结构
type Rule struct {
	Name       string //规则名 例如 标题
	StartLable string //开始位置的标签 例如 <title>
	EndLable   string //结束位置标签</title>
}

//扫描全部规则文件
func (s *Spider) ScanAll() ([]Urls, error) {
	configpath := viper.Get("SiteconfigDir")
	urls := []Urls{}
	err := filepath.Walk(configpath.(string),
		func(files string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			_, fileName := filepath.Split(files)
			ext := path.Ext(files)
			name := strings.Replace(fileName, ext, "", -1)
			url_config, err := s.GetUrlconfig(name)
			if err == nil {
				url := Urls{name, url_config.Url}
				urls = append(urls, url)
			}
			//fmt.Println(path, info.Size())
			return nil
		})
	return urls, err

}
func (s *Spider) GetUrlconfig(name string) (UrlConfig, error) {
	configpath := viper.Get("SiteconfigDir")
	filename := configpath.(string) + name + ".json"
	url_config := UrlConfig{}
	_, ok := helper.IsFile(filename)
	if ok {
		config_str := helper.JsonRead(filename)

		err := json.Unmarshal(config_str, &url_config)
		return url_config, err
	} else {
		errs := errors.New("没有找到该名称的规则文件")
		return url_config, errs
	}

}
func (s *Spider) AddUrlrule(name string, rule Rule) {

}

//根据名称参数启动爬虫
func (s *Spider) SpiderRun(name string) {
	if name == "" {
		//不指定名称 则全部爬取
		urls, err := s.ScanAll()
		if err == nil {
			for _, arg := range urls {
				s.NormalRun(arg.Name)
			}
		}
	} else {
		s.NormalRun(name)
	}
}

//通用爬虫启动
func (s *Spider) NormalRun(name string) {
	url_config, err := s.GetUrlconfig(name)
	if err != nil {
		fmt.Println("规则文件读取错误")
		return
	}
	httpurl := url_config.Url
	//域名自动补充 http
	if !strings.HasPrefix(url_config.Url, "http://") && !strings.HasPrefix(url_config.Url, "https://") {
		httpurl = "http://" + url_config.Url
	}
	url1 := strings.Split(httpurl, "//")[1]
	url2 := strings.Split(url1, "/")[0]
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains(url2),
	)
	//拷贝一份实例 用于访问分页
	pageLink := c.Clone()
	//拷贝一份实例 用于访问详情链接
	detailLink := c.Clone()
	//设置客户端，模拟浏览器访问
	//c.UserAgent = "xy"
	//允许重复访问
	//c.AllowURLRevisit = true

	// OnRequest 请求执行之前调用
	// OnResponse 响应返回之后调用
	// OnHTML 监听执行 selector
	// OnXML 监听执行 selector
	// OnHTMLDetach，取消监听，参数为 selector 字符串
	// OnXMLDetach，取消监听，参数为 selector 字符串
	// OnScraped，完成抓取后执行，完成所有工作后执行
	// OnError，错误回调
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})
	//提取分页
	pageLink.OnHTML("#kesfxqxq_A01_03_01", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		//content := e.ChildText("a")

		//fmt.Printf("detial link : %s \t", link)
		//fmt.Printf("detial content : %s \t", coverGBKToUTF8(content))
		//fmt.Println()

		detailLink.Visit(link)
	})
	//提取详情
	detailLink.OnHTML("#kesfxqxq_A01_03_01", func(e *colly.HTMLElement) {
		//content := e.ChildText("a")

		//fmt.Printf("detial link : %s \t", link)
		//fmt.Printf("detial content : %s \t", coverGBKToUTF8(content))
		//fmt.Println()

	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(httpurl)
}
