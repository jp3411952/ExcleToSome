/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"github.com/tealeg/xlsx"
	"github.com/zjytra/devlop/xutil/osutil"
	"github.com/zjytra/devlop/xutil/strutil"
	"strings"
	"sync"
)

//LangField 对应语言需要的字段
type LangField struct {

}

// 定义处理行数
type HandleFunc func(exclefileName string, excelContent [][]string)

// 具体的处理函数
var Hadler HandleFunc

var Wg sync.WaitGroup

// excle处理函数
var HandlerMap map[string]HandleFunc

func init() {
	HandlerMap = make(map[string]HandleFunc)
	HandlerMap["Ccsv"] = WriteC_Csv
	HandlerMap["Scsv"] = WriteS_Csv
	HandlerMap["S"] = WriteToGoFile
	HandlerMap["C"] = WriteToCsharp
	HandlerMap["L"] = writeLuaTable
	HandlerMap["A"] = ToAllFile
}

func GetHandlerFunc() HandleFunc {

	//根据输入的类型及输出类型获得对应的处理函数
	switch Conf.Intype {
	case "json", "xlsx":
		if fun, ok := HandlerMap[Conf.OutType]; ok {
			return fun
		}
	}

	return nil
}

func GetInPath() string {
	switch Conf.Intype {
	case "json", "xlsx":
		return Conf.InPath
	}
	return ""
}

func ChechAndMakeDir(dir string) bool {
	return !osutil.MakeDirAll(dir)
}

// 获取没有字段的下标
func GetNoFiledColIndex(row *xlsx.Row) map[int]bool {
	if row == nil {
		return nil
	}
	noDataColIndex := make(map[int]bool)
	for cellj, cell := range row.Cells {
		if strutil.StringIsNil(cell.String()) {
			noDataColIndex[cellj] = true
		}
	}
	return noDataColIndex
}

// 处理说有文件
func ToAllFile(exclefileName string, excelContent [][]string) {
	defer Wg.Done()
	for k, hanler := range HandlerMap {
		if strings.Compare(k, Conf.OutType) == 0 {
			continue
		}
		Wg.Add(1)
		go hanler(exclefileName, excelContent)
	}

}

//查看对应表格是否包含对应平台
func GetPlatfCol(excelContent [][]string, plts string) [][]string {
	platrow := excelContent[2] //第3行涉及具体语言字段 下标从0开始算的
	var field = make(map[int]int)
	//重新组装数据 只取自己关心的字段 第一层是行 第二层是列
	var newContent [][]string
	for col, cell := range platrow {
		//包含本平台以及全部平台的字段
		if strings.Contains(cell, plts) || strings.Contains(cell, "A") {
			field[col] = col
		}
	}
	//没有自己关心的数据
	if len(field) == 0 {
		return  nil
	}

	for i,row  := range  excelContent {
		if i == 2  { //把平台的字段干掉
			continue
		}
		var oneRow []string
		for j,col  := range  row {
			if _,ok := field[j];ok { //包含该列
				oneRow = append(oneRow,col)
			}
		}
		newContent = append(newContent,oneRow)
	}


	return newContent
}
