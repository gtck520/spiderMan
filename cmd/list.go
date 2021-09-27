/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var list_name string
var ListCmd = &cobra.Command{
	Use: "list",
	// 定义arguments数量最少为1个
	//Args:  cobra.MinimumNArgs(1),
	Short: "list",
	Long:  `列出已添加的站点列表`,
	Run: func(cmd *cobra.Command, args []string) {
		spider.List(cmd, args, list_name)
	},
}

func init() {
	rootCmd.AddCommand(ListCmd)
	ListCmd.Flags().StringVarP(&list_name, "name", "n", "", "取出对应名称的配置信息")
	// _ = ListCmd.MarkFlagRequired("name")
}
