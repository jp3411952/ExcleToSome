/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"encoding/csv"
	"fmt"
	"github.com/zjytra/devlop/xutil"
	"github.com/zjytra/devlop/xutil/osutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)


func WriteC_Csv(exclefileName string,celContent [][]string)  {
	doWriteCsv(exclefileName,celContent,"C")
}

func WriteS_Csv(exclefileName string,celContent [][]string)  {
	doWriteCsv(exclefileName,celContent,"S")
}

func doWriteCsv(exclefileName string,celContent [][]string,pltName string)  {
	defer Wg.Done()
	//拥有平台对应列的的数量
	excelContent := GetPlatfCol(celContent,pltName)
	if excelContent == nil {
		return
	}
	if len(excelContent) == 0 {
		return
	}
	csvoutdir := Conf.ServerOutCsv
	if pltName == "C" {
		csvoutdir = Conf.ClientOutCsv
	}
	if !osutil.MakeDirAll(csvoutdir) {
		return
	}

	filename := filepath.Join(csvoutdir,exclefileName+ ".csv")
	filebytes, _ := ioutil.ReadFile(filename)
	oldStr := string(filebytes)  //添加文件是否变化的验证
	var sb  = new(strings.Builder)
	//一次写入多行
	csvfileWt := csv.NewWriter(sb)
	if csvfileWt == nil {
		return
	}
	var newContent [][]string
	for _,row := range excelContent{
		var oneRow []string
		for _,val := range row{
			oneRow = append(oneRow,val)
		}
		newContent = append(newContent,oneRow)
	}
	csvfileWt.WriteAll(newContent)
	allstr := sb.String()
	//没有变化就不写
	if strings.Compare(allstr,oldStr) == 0 {
		return
	}
	//创建csv文件
	fs, err := os.Create(filename)
	defer fs.Close()
	if xutil.IsError(err) {
		return
	}
	fs.WriteString(allstr)
	fmt.Println( "csv 写完毕")
}

