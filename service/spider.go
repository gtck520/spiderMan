package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gtck520/spiderMan/helper"
	"github.com/spf13/viper"
)

type Spider struct {
	Co Colly
}

//url
type Urls struct {
	Name string
	Url  string
}

//定义每个url的配置结构
type UrlConfig struct {
	Name  string `json:"name"`  //配置名称
	Url   string `json:"url"`   //url
	Rules []Rule `json:"rules"` //规则
	Out   int    `json:"out"`   //选择输出方式 1 csv 2 mysql
}

//规则结构
type Rule struct {
	Type      int       `json:"type"`      //采集规则  1 html 2接口数据
	Name      string    `json:"name"`      //规则名称
	Field     string    `json:"field"`     //规则映射的字段名称
	Match     string    `json:"match"`     //type=2  则为正则表达式前缀，type=1 html 则为jquery selector规则
	PageMatch string    `json:"pagematch"` //type=2  链接中的分页字段，type=1 分页的jquery selector规则
	SubRule   []SubRule `json:"subrules"`  //下一级规则
}
type SubRule struct {
	Type  int    `json:"type"`  //采集规则  1 html 2接口数据
	Name  string `json:"name"`  //规则名称
	Field string `json:"field"` //规则映射的字段名称
	Match string `json:"match"` //type=2  则为正则表达式前缀，type=1 html 则为jquery selector规则
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

		fmt.Println("文件内容:" + string(config_str))
		err := json.Unmarshal(config_str, &url_config)
		return url_config, err
	} else {
		errs := errors.New("没有找到该名称的规则文件")
		return url_config, errs
	}

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
		fmt.Println("规则文件读取错误" + err.Error())
		return
	}
	s.Co.BuildC(url_config.Url)
	for _, rule := range url_config.Rules {
		s.Co.GetContent(rule)
		for _, subrule := range rule.SubRule {
			s.Co.GetDContent(subrule)
		}
	}
	s.Co.C.Visit(url_config.Url)
}
