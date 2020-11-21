/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"encoding/csv"
	"fmt"
	"wengo/xutil"
	"wengo/xutil/osutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func ToCsv(exclefileName string ,wg *sync.WaitGroup)  {
	defer wg.Done()
	//重命名文件
	excelContent := readxlsx(exclefileName)
	filenme :=	strings.TrimSuffix(exclefileName,".xlsx")
	WriteToCsv(filenme,excelContent)
}





func WriteToCsv(exclefileName string,excelContent [][]string)  {
	if !osutil.MakeDirAll(csvoutdir) {
		fmt.Printf("csvoutdir = ",csvoutdir)
		return
	}
	if excelContent == nil {
		fmt.Printf( "excle 数据为空")
		return
	}
	//创建csv文件
	fs, err := os.Create(filepath.Join(csvoutdir,exclefileName+ ".csv"))
	if xutil.IsError(err) {
		return
	}
	defer fs.Close()
	//一次写入多行
	csvfileWt := csv.NewWriter(fs)
	if csvfileWt == nil {
		return
	}
	csvfileWt.WriteAll(excelContent)
	fmt.Println( "csv 写完毕")
}

