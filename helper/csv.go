package helper

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"time"

	"github.com/spf13/viper"
)

type Csv struct {
	Writer *csv.Writer
	file   *os.File
}

func (c *Csv) CsvWrite(Name string, content map[string]string, writenum int) bool {
	ouputpath := viper.Get("OutputDir")
	now := time.Now()
	filename := ouputpath.(string) + Name + now.Format("2006_01_02_15_04_05") + ".csv"
	Checkdir(filename)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	defer f.Close()
	w := csv.NewWriter(f) //创建一个新的写入文件流
	defer w.Flush()
	title := make([]string, 0)
	data := [][]string{}
	values := make([]string, 0)

	for ka := range content {
		title = append(title, ka)
	}
	//排序
	// sort.Strings(title) 升序
	sort.Sort(sort.Reverse(sort.StringSlice(title)))
	for _, v := range title {
		values = append(values, content[v])
	}
	//第一次写入文件才添加标题栏
	if writenum == 0 {
		data = append(data, title)
	}

	data = append(data, values)

	err = w.WriteAll(data) //写入数据
	if err != nil {
		return false
	} else {
		return true
	}
}
