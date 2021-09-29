package helper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func JsonRead(filename string) []byte {
	fp, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 1000)
	n, err := fp.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(data[:n]))
	return data[:n]
}

func JsonWrite(data []byte, filename string) {
	Checkdir(filename)
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func Checkdir(files string) bool {
	// check
	paths, _ := filepath.Split(files)
	_, err := IsExists(paths)
	if err != true {
		fmt.Println("path not exists ", paths)
		err := os.MkdirAll(paths, 0711)
		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
			return false
		}
		//return true
	}
	_, ferr := IsFile(files)
	if ferr != true {
		//新建文件
		os.Create(files)
	}
	return true
}

// 判断路径是否存在
func IsExists(path string) (os.FileInfo, bool) {
	f, err := os.Stat(path)
	return f, err == nil || os.IsExist(err)
}

// 判断所给路径是否为文件夹
func IsDir(path string) (os.FileInfo, bool) {
	f, flag := IsExists(path)
	return f, flag && f.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) (os.FileInfo, bool) {
	f, flag := IsExists(path)
	return f, flag && !f.IsDir()
}
