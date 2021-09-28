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

func (s *Spider) SpiderRun(name string) {
	if name == "" {

	}
}

//通用爬虫启动
func (s *Spider) NormalRun(name string) {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://hackerspaces.org/")
}
