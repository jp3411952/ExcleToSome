/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"../generatepgl"
	"fmt"
	"github.com/showgo/xutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func ToGoFile(exclefileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	// 重命名文件
	excelContent := readxlsx(exclefileName)
	filenme := strings.TrimSuffix(exclefileName, ".xlsx")
	WriteToGoFile(filenme, excelContent)
}

func WriteToGoFile(exclefileName string, excelContent [][]string) {
	if !xutil.MakeDirAll(gostrcutdir) || xutil.StringIsNil(packagename) {
		fmt.Printf("gostrcutdir = ", gostrcutdir, "packagename = ", gostrcutdir)
		return
	}
	if excelContent == nil {
		fmt.Printf("excle 数据为空")
		return
	}
	// 创建csv文件
	fs, err := os.Create(filepath.Join(gostrcutdir, exclefileName+".go"))
	if xutil.IsError(err) {
		return
	}
	defer fs.Close()
	
	structName := xutil.Capitalize(exclefileName)
	// 一次写入多行
	genGoLang := generatepgl.NewGenerGoLang()
	genGoLang.WriteHead("//excle生成文件请勿修改\n", packagename)
	genGoLang.WriteMoreNextLine(1)
	
	genGoLang.WriteImport([]string{"\"fmt\"","\"github.com/showgo/csvparse\"", "\"github.com/showgo/xutil\""})
	// 定义变量
	varName := fmt.Sprintf("%sCsv", structName)
	vartypeName := fmt.Sprintf("map[%s]*%s", excelContent[1][0], structName)
	genGoLang.WriteVar(varName, vartypeName)
	genGoLang.WriteNextLine()
	// 写结构体内容
	genGoLang.WriteStruct(structName)
	colCount := len(excelContent[0])
	for i := 0; i < colCount; i++ {
		if len(excelContent) < 3 {
			fmt.Println(exclefileName,"小于三行")
			continue
		}
		filed := generatepgl.FiledInfo{
			xutil.Capitalize(excelContent[0][i]),                          // 第一行名称保证首字母大写
			excelContent[1][i],                                            // 第二行类型
			fmt.Sprint(excelContent[2][i], " 字段名称  ", excelContent[0][i]), // 第三行注释
		}
		genGoLang.WriteField(&filed)
	}
	// 结束括号
	genGoLang.WriteEndBrace()
	
	setfunName := fmt.Sprintf("Set%sMapData", structName)
	getfunName := fmt.Sprintf("get%sUsedData", structName)
	// // 写方法1
	// AsynSetfunName := fmt.Sprintf("AsynSet%sMapData", structName)
	// funcInfo := generatepgl.NewFuncInfo(AsynSetfunName)
	// funcInfo.FuncContent = fmt.Sprintf("\t go %s()", setfunName)
	// genGoLang.WriteFunc(funcInfo)
	
	// 写方法2
	csvPath := "csvpath"
	funcInfo2 := generatepgl.NewFuncInfo(setfunName)
	funcInfo2.FuncParam[csvPath] = generatepgl.GoString
	funcInfo2.FuncContent = fmt.Sprintf(
		`    if %s == nil {
		%s = make(%s)
	}
	tem := %s(%s)
	%s  = tem`,
		varName, varName, vartypeName, getfunName, csvPath, varName)
	genGoLang.WriteFunc(funcInfo2)
	// 写方法3
	funcInfo3 := generatepgl.NewFuncInfo(getfunName)
	funcInfo3.FuncParam["csvpath"] = generatepgl.GoString
	funcInfo3.FuncReturn["tem"] = vartypeName
	idName := xutil.Capitalize(excelContent[0][0])
	funcInfo3.FuncContent = fmt.Sprintf(
		`    csvmapdata := csvparse.GetCsvMapData(csvpath + "/%s.csv")
	tem := make(%s)
	for _, filedData := range csvmapdata {
		one := new(%s)
		for filedName, filedval := range filedData {
			isok := csvparse.ReflectSetField(one, filedName, filedval)
			xutil.IsError(isok)
			if _,ok := tem[one.%s]; ok {
				fmt.Println(one.%s,"重复")
			}
		}
		tem[one.%s] = one
	}`,
		exclefileName, vartypeName, structName,idName,idName,idName )
	genGoLang.WriteFunc(funcInfo3)
	
	// 写方法4
	funcInfo4 := generatepgl.NewFuncInfo(fmt.Sprintf("Get%sPtr", structName))
	funcInfo4.FuncContent = fmt.Sprintf(
		`    data, ok := %s[%s];
	if  !ok  {
		return nil
	}`,varName,excelContent[0][0])
	funcInfo4.FuncParam[excelContent[0][0]] = excelContent[1][0]
	returnName := fmt.Sprintf("%s","data")
	funcInfo4.FuncReturn[returnName] = fmt.Sprintf("*%s", structName)
	genGoLang.WriteFunc(funcInfo4)
	
	// 一次性写入所有数据
	fs.WriteString(genGoLang.String())
}
