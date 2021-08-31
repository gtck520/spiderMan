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
	"github.com/gtck520/spiderMan/controller"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var spider = controller.Spider{}
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加一个目标url地址",
	Long:  `添加一个url地址，将此地址加入到工作队列中，并同时生成这个地址爬虫规则文件。`,
	Run:   spider.Add,
}

func init() {
	rootCmd.AddCommand(addCmd)

	// 下面定义了一个Flag foo, foo后面接的值会被赋值给Foo
	//Foo = addCmd.PersistentFlags().String("foo", "", "A help for foo")
	// 下面定义了一个Flag print ,print后面的值会被赋值给Print变量
	//addCmd.PersistentFlags().StringVar(&Print, "print", "", "print")
	// 下面定义了一个Flag show,show默认为false, 有两种调用方式--show\-s，命令后面接了show则上面定义的show变量就会变成true
	//addCmd.PersistentFlags().BoolVarP(&show, "show", "s", false, "show")
	// 下面定义了一个Flag show,show默认为false, 有两种调用方式--show\-s，命令后面接了show则上面定义的show变量就会变成true
	//showL = *addCmd.Flags().BoolP("showL", "S", false, "show")
	// 下面定义了一个Flag print ,print后面的值会被赋值给Print变量
	//addCmd.Flags().StringVar(&PrintL, "printL", "", "print")
	// 下面定义了一个Flag fooL, foo后面接的值会被赋值给FooL
	//FooL = addCmd.Flags().String("fooL", "", "A help for foo")
	//show = *testCmd.Flags().BoolP("show", "s", false, "show")
	// 设置使用test的时候后面必须接show
	//_ = testCmd.MarkFlagRequired("show")

}
