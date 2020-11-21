/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"wengo/xutil/osutil"
	"wengo/xutil/strutil"
)

// 定义处理行数
type HandleFunc func(exclefileName string, wg *sync.WaitGroup)

// excle处理函数
var ExlceHandlerMap map[string]HandleFunc

// json处理函数
var JsonHandlerMap map[string]HandleFunc

func init() {
	ExlceHandlerMap = make(map[string]HandleFunc)
	JsonHandlerMap = make(map[string]HandleFunc)
	ExlceHandlerMap["csv"] = ToCsv
	ExlceHandlerMap["go"] = ToGoFile
	ExlceHandlerMap["all"] = ToAllFile
	JsonHandlerMap["go"] = JsonToGo
}

func GetHandlerFunc() HandleFunc {
	
	//根据输入的类型及输出类型获得对应的处理函数
	switch Intype {
	case "json": //json 的输出方法
		if fun,ok := JsonHandlerMap[OutType] ; ok  {
			return fun
		}
	case "xlsx": //xlsx 输出的方法
		if fun,ok := ExlceHandlerMap[OutType] ; ok  {
			return fun
		}
	}

	return  nil
}

func GetInPath() string {
	switch Intype {
	case "json","xlsx":
		return InPath
	}
return ""
}


func ChechAndMakeDir(dir string) bool{
	return !osutil.MakeDirAll(dir)
}

// 读取excle文件
func readxlsx(exclefileName string) [][]string {
	filename := filepath.Join(InPath, exclefileName)
	xlFile, err := xlsx.OpenFile(filename)
	if xutil.IsError(err) {
		return nil
	}
	sheet1 := xlFile.Sheet["Sheet1"]
	if sheet1 == nil {
		fmt.Printf(exclefileName, "没有Sheet1 表,只使用Sheet1")
		return nil
	}
	rownum := len(sheet1.Rows)
	if rownum == 0 {
		return nil
	}
	// 构建表数据,二维数组  先不make,避免产生无用的数据
	var newContent [][]string
	noDataColIndex := GetNoFiledColIndex(sheet1.Rows[0]) // 记录没有字段的下标
	for rownum, row := range sheet1.Rows {
		firstColcell := row.Cells[0].String() // 第一列的数据
		if !xutil.ValidCsvRow(firstColcell, rownum) { // 如果无效就记录为无效数据
			continue
		}
		var oneRow []string
		for cellj, cell := range row.Cells {
			if noDataColIndex[cellj] { // 如果没有字段就不写数据
				continue
			}
			oneRow = append(oneRow, cell.String())
		}
		if oneRow != nil {
			newContent = append(newContent, oneRow)
		}
	}
	return newContent
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
func ToAllFile(exclefileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	for k, hanler := range ExlceHandlerMap {
		if strings.Compare(k, "all") == 0 {
			continue
		}
		wg.Add(1)
		go hanler(exclefileName, wg)
	}
	
}

func ReadJson()  {
	
}