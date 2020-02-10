/*
创建时间: 2020/2/6
作者: zjy
功能介绍:

*/

package main

import (
	"./excletosome"
	"flag"
	"fmt"
	"github.com/showgo/xutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 具体的处理函数
var Hadler excletosome.HandleFunc

var wg sync.WaitGroup

func main() {
	excletosome.ReadConf() // 读取配置文件
	if !ReadCommondLine() {
		fmt.Println("outType must be csv||go||sql||all")
		time.Sleep(time.Second * 5)
		return
	}
	if xutil.StringIsNil(excletosome.InPath) {
		return
	}
	ReadExcleDirFile()
}

func ReadCommondLine() bool {
	var outType string
	flag.StringVar(&outType, "outType", "csv", "请输入outType类型 csv||go||sql||all")
	flag.Parse()
	Hadler = excletosome.GetHandlerFunc(outType)
	if Hadler == nil {
		fmt.Println("outType=",outType)
		return false
	}
	return true
}



func ReadExcleDirFile() {
	fmt.Println("读取需要解析的文件")
	// //遍历xlsx目录遍历指定目录下所有文件
	filepath.Walk(excletosome.InPath,walkOnefile)
	wg.Wait()
	fmt.Println("执行完毕!")
	time.Sleep(time.Second * 5)
}

func walkOnefile(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil // 如果是文件就不管
	}
	//输入的是xlsx 的处理
	if strings.Compare(excletosome.Intype,"xlsx") == 0 &&
		xutil.IsXlsx(info.Name()) {
		wg.Add(1)
		go Hadler(info.Name(), &wg)
	}
	return nil
}



