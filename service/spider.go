package service

import (
	"encoding/json"
	"errors"

	"github.com/gtck520/spiderMan/helper"
	"github.com/spf13/viper"
)

type Spider struct {
	UrlList map[string]string
}

func (s *Spider) GetUrlconfig(name string) (helper.UrlConfig, error) {
	configpath := viper.Get("SiteconfigDir")
	filename := configpath.(string) + name + ".json"
	url_config := helper.UrlConfig{}
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
