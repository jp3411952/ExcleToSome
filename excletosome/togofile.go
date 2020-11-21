/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"ExcleToSome/generatepgl"
	"fmt"
	"github.com/wengo/xutil"
	"github.com/wengo/xutil/osutil"
	"github.com/wengo/xutil/strutil"
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
	if !osutil.MakeDirAll(gostrcutdir) || strutil.StringIsNil(packagename) {
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
	//写文件头
	writeHeder(genGoLang)
	// 定义变量
	varName := fmt.Sprintf("%sAtomic", exclefileName)
	csvtypeName := fmt.Sprintf("map[%s]*%s", excelContent[1][0], structName)
	atomictype := fmt.Sprintf("atomic.Value")
	genGoLang.WriteVar(varName, atomictype)
	genGoLang.WriteNextLine()
	// 写结构体内容
	writeStruct(genGoLang, structName, excelContent, exclefileName)
	setfunName := fmt.Sprintf("Set%sMapData", structName)
	getfunName := fmt.Sprintf("load%sUsedData", structName)
	// 写方法
	setcsvData(setfunName, varName, getfunName, genGoLang)
	getCsvData(getfunName, csvtypeName, excelContent, exclefileName, structName, genGoLang)
	getAllfunName := fmt.Sprintf("GetAll%s", structName)
	getptrFunc(structName, getAllfunName, excelContent, genGoLang)
	getAllfun(getAllfunName, varName, csvtypeName, genGoLang)
	
	// 一次性写入所有数据
	fs.WriteString(genGoLang.String())
}

func writeHeder(genGoLang *generatepgl.GenerGoLang) {
	genGoLang.WriteHead("//excle生成文件请勿修改\n", packagename)
	genGoLang.WriteMoreNextLine(1)
	genGoLang.WriteImport([]string{
		xutil.GetPackageStr("fmt"),
		xutil.GetPackageStr("github.com/wengo/csvparse"),
		xutil.GetPackageStr("github.com/wengo/xutil"),
		xutil.GetPackageStr("github.com/wengo/xlog"),
		xutil.GetPackageStr("sync/atomic")})
}

func writeStruct(genGoLang *generatepgl.GenerGoLang, structName string, excelContent [][]string, exclefileName string) {
	genGoLang.WriteStruct(structName)
	colCount := len(excelContent[0])
	for i := 0; i < colCount; i++ {
		if len(excelContent) < 3 {
			fmt.Println(exclefileName, "小于三行")
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
}

func setcsvData(setfunName string, varName string, getfunName string, genGoLang *generatepgl.GenerGoLang) {
	csvPath := "csvpath"
	funcInfo := generatepgl.NewFuncInfo(setfunName)
	funcInfo.FuncParam[csvPath] = generatepgl.GoString
	funcInfo.FuncContent = fmt.Sprintf(
		`  	defer xlog.RecoverToStd()
	%s.Store(%s(csvpath))`,
		varName, getfunName)
	genGoLang.WriteFunc(funcInfo)
}

func getCsvData(getfunName string, csvtypeName string, excelContent [][]string, exclefileName string, structName string, genGoLang *generatepgl.GenerGoLang) {
	// 写方法3
	funcInfo := generatepgl.NewFuncInfo(getfunName)
	funcInfo.FuncParam["csvpath"] = generatepgl.GoString
	funcInfo.FuncReturn["tem"] = csvtypeName
	idName := xutil.Capitalize(excelContent[0][0])
	funcInfo.FuncContent = fmt.Sprintf(
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
	}`,exclefileName, csvtypeName, structName, idName, idName, idName)
	genGoLang.WriteFunc(funcInfo)
}

func getptrFunc(structName string, getAllfunName string, excelContent [][]string, genGoLang *generatepgl.GenerGoLang) {
	// 写方法4
	funcInfo4 := generatepgl.NewFuncInfo(fmt.Sprintf("Get%sPtr", structName))
	funcInfo4.FuncContent = fmt.Sprintf(
		`    alldata := %s()
	if alldata == nil {
		return nil
	}
	if data, ok := alldata[%s]; ok {
		return data
	}`, getAllfunName, excelContent[0][0])
	funcInfo4.FuncParam[excelContent[0][0]] = excelContent[1][0]
	funcInfo4.FuncReturn["nil"] = fmt.Sprintf("*%s", structName)
	genGoLang.WriteFunc(funcInfo4)
}

func getAllfun(getAllfunName string, varName string, csvtypeName string, genGoLang *generatepgl.GenerGoLang) {
	// getAll
	funcInfo := generatepgl.NewFuncInfo(getAllfunName)
	funcInfo.FuncContent = fmt.Sprintf(
		`    val := %s.Load()
	if data, ok := val.(%s); ok {
		return data
	}`, varName, csvtypeName)
	funcInfo.FuncReturn["nil"] = fmt.Sprintf("%s", csvtypeName)
	genGoLang.WriteFunc(funcInfo)
}
