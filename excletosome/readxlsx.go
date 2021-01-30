package excletosome

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/zjytra/devlop/xutil"
	"github.com/zjytra/devlop/xutil/strutil"
	"path/filepath"
	"strings"
)

// 读取excle文件
func readxlsx(exclefileName string) (string,[][]string) {
	filename := filepath.Join(Conf.InPath, exclefileName)
	xlFile, err := xlsx.OpenFile(filename)
	var confName string // 配置表名称
	if xutil.IsError(err) {
		return confName,nil
	}
	var sheet1 *xlsx.Sheet
	for sheetName, sheet := range xlFile.Sheet {
		if strings.Contains(sheetName, "Sheet") || strings.Contains(sheetName, "注释") { //排除未取名的表
			continue
		}
		sheet1 = sheet
		confName = sheetName
		break
	}

	if sheet1 == nil {
		fmt.Printf(exclefileName, "没有找到表")
		return  confName,nil
	}
	rownum := sheet1.MaxRow
	if rownum == 0 {
		return  confName,nil
	}
	// 构建表数据,二维数组  先不make,避免产生无用的数据
	var newContent [][]string
	noDataColIndex := GetNoFiledColIndex(sheet1.Rows[0]) // 记录没有字段的下标

	for rownum, row := range sheet1.Rows {
		firstColcell := row.Cells[0].String()   // 第一列的数据
		if !ValidCsvRow(firstColcell, rownum) { // 如果无效就记录为无效数据
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
	return confName,newContent
}


// 验证csv行数据是否有效
// 只能第三行保留注释,  有注释 str首字符 != #  ASCII表 35
// 第3行是代表那个平台  rownum = 2
// 第4行是注释  rownum = 3
// 并且id不为nil
func ValidCsvRow(str string, rownum int) bool {
	if strutil.StringIsNil(str) {
		return false
	}
	if rownum != 3 && strings.Index(str,"#") == 0 {
		return false
	}
	return true
}

