package controller

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gtck520/spiderMan/helper"
	"github.com/gtck520/spiderMan/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var spider_service = service.Spider{}

type Spider struct {
	UrlList map[string]string
}

// 添加一个目标url地址
func (s *Spider) Add(cmd *cobra.Command, args []string, globals ...string) {

	names := strings.Split(globals[0], ",")
	if len(args) != len(names) {
		fmt.Println("url参数个数与名称个数不一致！")
		return
	}
	configpath := viper.Get("SiteconfigDir")
	for k, arg := range args {
		jsonconfig := service.UrlConfig{}
		jsonconfig.Name = names[k]
		jsonconfig.Url = arg
		//子页面规则
		subrule := service.SubRule{}
		subrule.Field = "title"
		subrule.Match = ".main-title"
		subrule.Type = 1
		subrule.Name = "标题"
		//列表页面规则
		rule := service.Rule{}
		rule.Field = "list"
		rule.Match = `"url":"(.*?)"`
		rule.Name = "列表"
		rule.Type = 2
		rule.SubRule = append(rule.SubRule, subrule)

		jsonconfig.Rules = append(jsonconfig.Rules, rule)
		filename := names[k] + ".json"
		content, _ := json.Marshal(jsonconfig)
		helper.JsonWrite(content, configpath.(string)+filename)

	}
}

//
func (s *Spider) List(cmd *cobra.Command, args []string, globals ...string) {
	if len(globals) > 0 && globals[0] != "" {
		url_config, err := spider_service.GetUrlconfig(globals[0])
		if err == nil {
			fmt.Println("列出规则：", url_config)
		} else {
			fmt.Println(err.Error())
		}
	} else {
		urls, err := spider_service.ScanAll()
		if err != nil {
			fmt.Println(err.Error())
		}
		//fmt.Println("列出所有目录")
		t := helper.Table(urls)
		fmt.Println(t)
	}

}
func (s *Spider) SpiderRun(cmd *cobra.Command, args []string, globals ...string) {
	var name = ""
	if len(globals) > 0 && globals[0] != "" {
		name = globals[0]
	}
	spider_service.SpiderRun(name)
}
