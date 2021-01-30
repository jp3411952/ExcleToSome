package excletosome

import (
	"fmt"
	"github.com/zjytra/devlop/xutil"
	"os"
	"strings"
)

//生成的文件只包含对应平台关心的字段


func SetOutTypeFun() bool {
	Hadler = GetHandlerFunc()
	if Hadler == nil {
		fmt.Println("outType=",Conf.OutType)
		return false
	}
	return true
}

func WalkOnefile(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil // 如果是文件就不管
	}
	//输入的是xlsx 的处理
	if strings.Compare(Conf.Intype,"xlsx") == 0 && xutil.IsXlsx(info.Name()) {
		Wg.Add(1)
		go ParseXlsx(info.Name())
		return nil
	}

	////输入的是json 的处理
	//if strings.Compare(Conf.Intype,"json") == 0 {
	//	go ParseJson(info.Name())
	//	return nil
	//}
	return nil
}


func ParseXlsx(fileName string) {
	defer Wg.Done()
	confName,excelContent := readxlsx(fileName)
	Wg.Add(1)
	go Hadler(confName,excelContent)
}


func ParseJson(fileName string) {
	Wg.Add(1)
}