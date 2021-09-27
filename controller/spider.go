package controller

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gtck520/spiderMan/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

	filepath := viper.Get("SiteconfigDir")
	for k, arg := range args {
		jsonconfig := helper.UrlConfig{}
		jsonconfig.Name = names[k]
		jsonconfig.Url = arg
		filename := names[k] + ".json"
		content, _ := json.Marshal(jsonconfig)
		helper.JsonWrite(content, filepath.(string)+filename)

	}
}
