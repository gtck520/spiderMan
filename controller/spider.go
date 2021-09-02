package controller

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Spider struct {
	UrlList map[string]string
}

// 添加一个目标url地址
func (s *Spider) Add(cmd *cobra.Command, args []string, globals ...string) {
	for _, arg := range args {
		fmt.Println("Arg:", arg)
	}
	for _, glo := range globals {
		fmt.Println("glo:", glo)
	}
}
