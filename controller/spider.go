package controller

import "github.com/spf13/cobra"

type Spider struct {
	UrlList map[string]string
}

// 添加一个目标url地址
func (s *Spider) Add(cmd *cobra.Command, args []string) {
	println("被调用")
}
