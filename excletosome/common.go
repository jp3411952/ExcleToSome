/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"fmt"
	"github.com/showgo/xutil"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"strings"
	"sync"
)

// 定义处理行数
type HandleFunc func(exclefileName string, wg *sync.WaitGroup)

// 处理函数
var HandlerMap map[string]HandleFunc

func init() {
	HandlerMap = make(map[string]HandleFunc)
	HandlerMap["csv"] = ToCsv
	HandlerMap["go"] = ToGoFile
	HandlerMap["all"] = ToAllFile
}

func GetHandlerFunc(handleName string) HandleFunc {
	return HandlerMap[handleName]
}


func ChechAndMakeDir(dir string) bool{
	return xutil.StringIsNil(dir) || !xutil.MakeDirAll(dir)
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
		if xutil.StringIsNil(cell.String()) {
			noDataColIndex[cellj] = true
		}
	}
	return noDataColIndex
}

// 处理说有文件
func ToAllFile(exclefileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	for k, hanler := range HandlerMap {
		if strings.Compare(k, "all") == 0 {
			continue
		}
		wg.Add(1)
		go hanler(exclefileName, wg)
	}
	
}
