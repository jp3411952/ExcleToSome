/*
创建时间: 2020/2/6
作者: zjy
功能介绍:

*/

package main

import (
	"fmt"
	"github.com/zjytra/ExcleToSome/excletosome"
	"github.com/zjytra/devlop/xutil/strutil"
	"path/filepath"
	"time"
)

func Test()  {
	//tem := csvdata.SetTestMapData("csv/")
	//fmt.Printf("%v",tem)
}


func main() {
	//Test()
	DoWriteFile()
}

func DoWriteFile() {

	excletosome.ReadConfJson() // 读取配置文件

	if !excletosome.SetOutTypeFun() {
		fmt.Println("outType must be csv||go||sql||all")
		time.Sleep(time.Second * 5)
		return
	}
	if strutil.StringIsNil(excletosome.Conf.InPath) {
		return
	}
	ReadDirFiles()
}




func ReadDirFiles() {
	fmt.Println("读取需要解析的文件")
	// //遍历xlsx目录遍历指定目录下所有文件
	filepath.Walk(excletosome.GetInPath(),excletosome.WalkOnefile)
	excletosome.Wg.Wait()
	excletosome.WriteCfgMgr()
	fmt.Println("执行完毕!")
}




