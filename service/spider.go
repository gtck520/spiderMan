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
	Name   string `json:"name"`   //配置名称
	Tbname string `json:"tbname"` //对应表名
	Url    string `json:"url"`    //url
	Rules  []Rule `json:"rules"`  //规则
	Out    int    `json:"out"`    //选择输出方式 1 csv 2 mysql
}

//规则结构
type Rule struct {
	Type     int       `json:"type"`     //采集规则  1 html 2接口数据
	Name     string    `json:"name"`     //规则名称
	Field    string    `json:"field"`    //规则映射的字段名称
	Match    string    `json:"match"`    //type=2  则为正则表达式前缀，type=1 html 则为jquery selector规则
	PageRule PageRule  `json:"pagerule"` //分页规则
	SubMatch string    `json:"submatch"` //内容规则 jquery selector
	SubRule  []SubRule `json:"subrules"` //内容详细内容
}
type SubRule struct {
	Type  int    `json:"type"`  //采集规则  1 html 2接口数据
	Name  string `json:"name"`  //规则名称
	Field string `json:"field"` //规则映射的字段名称
	Match string `json:"match"` //type=2  则为正则表达式前缀，type=1 html 则为jquery selector规则
}

//分页规则
type PageRule struct {
	Type  int    `json:"type"`  //采集规则 0 不采集分页  1 html 2接口数据
	Page  int    `json:"page"`  //采集第几页 0为采集全部
	Num   int    `json:"num"`   //总共采集几页 0为无限制
	Match string `json:"match"` //分页提取规则 type为1：则输入jquery提取分页按钮  2：输入分页连接其中 分页处用{page}替换
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

		//fmt.Println("文件内容:" + string(config_str))
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
		s.Co.GetDContent(name, rule, rule.SubRule)
		s.Co.GetPageContent(rule.PageRule)
	}
	s.Co.C.Visit(url_config.Url)
	csvwriter := helper.Csv{}
	mysqlconn := helper.GetMysqlInstance()
	if url_config.Out == 2 {
		fields := make(map[string]string)
		fields["link"] = "varchar(100) not null default ''"
		fields["title"] = "varchar(100) not null default ''"
		fields["content"] = "text not null default ''"
		mysqlconn.CreateTable(url_config.Tbname, fields)
	}
	//写入存储
	for {
		select {
		case i, ok := <-s.Co.WriteChannle:
			if !ok {
				//通道关闭 则退出不阻塞
				return
			}
			//fmt.Printf("读取：%s \n", i)
			if url_config.Out == 1 {
				if csvwriter.CsvWrite(name, i, s.Co.WriteNum) {
					s.Co.WriteNum++
					fmt.Printf("写入数：%d \n", s.Co.WriteNum)
				}
			} else if url_config.Out == 2 {
				if mysqlconn.MysqlWrite(url_config.Tbname, i) {
					s.Co.WriteNum++
					fmt.Printf("写入数：%d \n", s.Co.WriteNum)
				}
			}
		default:
			fmt.Println("数据接收完毕")
			//没有数据也退出不阻塞
			return
		}
	}
}
