/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"fmt"
	"github.com/zjytra/ExcleToSome/generatepgl"
	"github.com/zjytra/devlop/xutil"
	"github.com/zjytra/devlop/xutil/osutil"
	"github.com/zjytra/devlop/xutil/strutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)


func WriteToGoFile(exclefileName string, celContent [][]string) {
	defer Wg.Done()

	//拥有平台对应列的的数量
	excelContent := GetPlatfCol(celContent,"S")
	if excelContent == nil {
		return
	}
	if len(excelContent) == 0 {
		return
	}
	if !osutil.MakeDirAll(Conf.GostrcutOutpath) || strutil.StringIsNil(Conf.GopakName) {
		fmt.Printf("GostrcutOutpath = ", Conf.GostrcutOutpath, "packagename = ", Conf.GopakName)
		return
	}


	filename := filepath.Join(Conf.GostrcutOutpath, exclefileName+".go")
	filebytes, _ := ioutil.ReadFile(filename)
	oldStr := string(filebytes)  //添加文件是否变化的验证
	structName := xutil.Capitalize(exclefileName)
	// 一次写入多行
	genGoLang := generatepgl.NewGenerGoLang()
	typeRow := excelContent[1]
	var isConvert,HasStrConv bool
	//查看是否需要加入转换包
	for _,typestr := range typeRow {
		if strings.Contains(typestr,"bool") || strings.Contains(typestr,"float")  {
			isConvert = true
			if HasStrConv { //另外有个条件也达成就跳出
				break
			}
		}
		if typestr != "string" {  //有字符串转换
			HasStrConv = true
			if isConvert {
				break
			}
		}
	}
	//写文件头
	writeHeder(genGoLang,isConvert,HasStrConv)
	// 定义变量
	varName := fmt.Sprintf("%sMap", exclefileName)
	csvtypeName := fmt.Sprintf("map[%s]*%s", excelContent[1][0], structName)
	//atomictype := fmt.Sprintf("atomic.Value")
	genGoLang.WriteVar(varName, csvtypeName)
	genGoLang.WriteNextLine()
	// 写结构体内容
	writeStruct(genGoLang, structName, excelContent)
	setfunName := fmt.Sprintf("Set%sMapData", structName)
	getfunName := fmt.Sprintf("load%sCsv", structName)
	// 写方法
	setcsvData(setfunName, varName, getfunName, genGoLang)
	getCsvData(getfunName, csvtypeName, excelContent, exclefileName, structName, genGoLang)
	getAllfunName := fmt.Sprintf("GetAll%s", structName)
	getptrFunc(varName,structName, excelContent, genGoLang)
	getAllfun(getAllfunName, varName, csvtypeName, genGoLang)

	newStr := genGoLang.String()
	if strings.Compare(oldStr,newStr)  ==  0 {
		return
	}

	// 创建go文件 要验证不一样才写
	fs, err := os.Create(filename)
	defer fs.Close()
	if xutil.IsError(err) {
		return
	}
	// 一次性写入所有数据
	fs.WriteString(newStr)
}

func writeHeder(genGoLang *generatepgl.GenerGoLang,isConvert bool,HasStrConv bool) {
	genGoLang.WriteHead("//excle生成文件请勿修改\n", Conf.GopakName)
	genGoLang.WriteMoreNextLine(1)

	var  pakageArr  []string
	pakageArr = append(pakageArr,xutil.GetPackageStr("fmt"))
	pakageArr = append(pakageArr,xutil.GetPackageStr("github.com/zjytra/wengo/csvsys/csvparse"))
	if HasStrConv {
		pakageArr = append(pakageArr,xutil.GetPackageStr("github.com/zjytra/devlop/xutil/strutil"))
	}
	if isConvert {
		pakageArr = append(pakageArr,xutil.GetPackageStr("strconv"))
	}
	pakageArr = append(pakageArr,xutil.GetPackageStr("strings"))
	genGoLang.WriteImport(pakageArr)
		//xutil.GetPackageStr("github.com/zjytra/wengo/engine_core/xlog")
}

func writeStruct(genGoLang *generatepgl.GenerGoLang, structName string, excelContent [][]string) {
	genGoLang.WriteStruct(structName)
	colCount := len(excelContent[0])
	for i := 0; i < colCount; i++ {
		filed := generatepgl.FiledInfo{
			xutil.Capitalize(excelContent[0][i]),                          // 第一行名称保证首字母大写
			excelContent[1][i],                                            // 第二行类型
			fmt.Sprint(excelContent[2][i], " 字段名称  ", excelContent[0][i]), // 第四行注释 平台行被剔除了
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
		`   %s = %s(csvpath)`,varName, getfunName)
	genGoLang.WriteFunc(funcInfo)
}

func getCsvData(getfunName string, csvtypeName string, excelContent [][]string, exclefileName string, structName string, genGoLang *generatepgl.GenerGoLang) {
	// 写方法3
	funcInfo := generatepgl.NewFuncInfo(getfunName)
	funcInfo.FuncParam["csvpath"] = generatepgl.GoString
	funcInfo.FuncReturn["tem"] = csvtypeName
	idName := xutil.Capitalize(excelContent[0][0])
	funcInfo.FuncContent = fmt.Sprintf(
`
	csvName := "/%s.csv"
	csvmapdata := csvparse.GetCsvSliceData(csvpath + csvName)
	if csvmapdata == nil {
		fmt.Printf("获取csv字符串错误%v",csvName)
		return nil
	}`,exclefileName,"%v")
	//赋值字段
	nameRow  := excelContent[0]
	typeRow  := excelContent[1]
	sb := new(strings.Builder)
	for i := 0;i < len(nameRow); i++ {
		filedName := xutil.Capitalize(nameRow[i])
		rowStr := fmt.Sprintf(`
		done = csvparse.CheckType(one.%s, typeRow[col], nameRow[col],csvName)
		if !done {
			return nil
		}
		one.%s %s
		col++
		`,filedName,filedName,GetTypeConvertFun(typeRow[i]))
		sb.WriteString(rowStr);
	}
	//字段名称

	funcInfo.FuncContent = funcInfo.FuncContent + fmt.Sprintf(
	`
	tem := make(%s)
	nameRow  := csvmapdata[0]
	typeRow  := csvmapdata[1]
	var col int
    var done bool
	for rowNum, oneRow := range csvmapdata {
		if rowNum < csvparse.Invalid_Row { // 排除前三行
			continue
		}
		col = 0 //重置变量
		//第一个是#的字符行忽略掉
		if strings.Index(oneRow[col],"#") == 0 {
			continue
		}
		one := new(%s)
		%s
		if _,ok := tem[one.%s]; ok {
			fmt.Println(one.%s,"重复")
		}
		tem[one.%s] = one
	}`,csvtypeName, structName,sb.String(), idName, idName, idName)
	genGoLang.WriteFunc(funcInfo)
}

//变量名称
func getptrFunc(varmapName string, structName string, excelContent [][]string, genGoLang *generatepgl.GenerGoLang) {
	// 写方法4
	funcInfo4 := generatepgl.NewFuncInfo(fmt.Sprintf("Get%sPtr", structName))
	funcInfo4.FuncContent = fmt.Sprintf(
		`   if data, ok := %s[%s]; ok {
		return data
	}`, varmapName, excelContent[0][0])
	funcInfo4.FuncParam[excelContent[0][0]] = excelContent[1][0]
	funcInfo4.FuncReturn["nil"] = fmt.Sprintf("*%s", structName)
	genGoLang.WriteFunc(funcInfo4)
}

func getAllfun(getAllfunName string, varName string, csvtypeName string, genGoLang *generatepgl.GenerGoLang) {
	// getAll
	funcInfo := generatepgl.NewFuncInfo(getAllfunName)
	funcInfo.FuncReturn[varName] = fmt.Sprintf("%s", csvtypeName)
	genGoLang.WriteFunc(funcInfo)
}


func GetTypeConvertFun(fliedtype string) string {
	switch fliedtype {
	case "int":
		return "= strutil.StrToInt(oneRow[col])"
	case "int8":
		return "= strutil.StrToInt8(oneRow[col])"
	case "uint8":
		return "= strutil.StrToUint8(oneRow[col])"
	case "int16":
		return "= strutil.StrToInt16(oneRow[col])"
	case "uint16":
		return "= strutil.StrToUint16(oneRow[col])"
	case "int32":
		return "= strutil.StrToInt32(oneRow[col])"
	case "uint32":
		return "= strutil.StrToUint32(oneRow[col])"
	case "float64":
		return ",_ = strconv.ParseFloat(oneRow[col], 64)"
	case "string":
		return "= oneRow[col]"
	case "bool":
		return ",_ = strconv.ParseBool(oneRow[col])"
	case "[]int","int[]":
		return "= StringsToIntArr(RepleaceBrackets(oneRow[col]))"
	case "[]string","string[]":
		return "= strings.Split(RepleaceBrackets(oneRow[col]))"
	default:
		fmt.Println(fliedtype, "is an unknown type.")
		return ""
	}
	return ""
}