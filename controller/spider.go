package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
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
type Urls struct {
	Name string
	Url  string
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
		jsonconfig := helper.UrlConfig{}
		jsonconfig.Name = names[k]
		jsonconfig.Url = arg
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
				url_config, err := spider_service.GetUrlconfig(name)
				if err == nil {
					url := Urls{name, url_config.Url}
					urls = append(urls, url)
				}
				//fmt.Println(path, info.Size())
				return nil
			})
		t := helper.Table(urls)
		fmt.Println(t)
		if err != nil {
			fmt.Println(err.Error())
		}
		//fmt.Println("列出所有目录")
	}

}
