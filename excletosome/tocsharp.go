package excletosome

import (
	"fmt"
	"github.com/zjytra/devlop/xutil"
	"github.com/zjytra/devlop/xutil/osutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func WriteToCsharp(exclefileName string,exelContent [][]string) {
	defer Wg.Done()

	//转换成自己关心的数据
	excelContent := GetPlatfCol(exelContent,"C")
	if excelContent == nil {
		return
	}
	if len(excelContent) == 0 {
		return
	}
	if !osutil.MakeDirAll(Conf.CsharpOut) {
		return
	}

	structName := xutil.Capitalize(exclefileName) //
	filename := filepath.Join(Conf.CsharpOut, structName+".cs")
	filebytes, _ := ioutil.ReadFile(filename)
	oldStr := string(filebytes)  //添加文件是否变化的验证

	templpltFile,err :=  os.Open("./templt/csharpcsvtemplt")
	if err != nil {
		fmt.Printf("打开csharp模板文件出错 %v",err)
		return
	}

	//unity可寻址系统名称
	addrName := "Csv/" + exclefileName  //6
	var buf  = make([]byte,4096)
	blen,err2 := templpltFile.Read(buf)
	if buf == nil ||  err2 != nil {
		fmt.Printf("读取模板文件出错 %v",err2)
		return
	}
	for i := 0;i < blen;i++ {
		if buf[i] == 0 {
			blen = i
			break
		}
	}
	buf = buf[0:blen - 1]
	///处理null串
	//赋值字段
	nameRow  := excelContent[0]
	typeRow  := excelContent[1]
	sb := new(strings.Builder)
	for i := 0;i < len(nameRow); i++ {
		filedName := xutil.Capitalize(nameRow[i])
		rowStr := fmt.Sprintf(`
				typeIsMatch = CheckTypeIsMatch(oneData.%s.GetType().Name, rowType[col], rowName[col]);
				if (!typeIsMatch)
				{
					return;
				}
				oneData.%s %s;
		`,filedName,filedName,GetCsharpConvertFun(typeRow[i]))
		sb.WriteString(rowStr)
	}

	keyType := CsvTypeToCSharpType(excelContent[1][0])  //3
	keyName := xutil.Capitalize(excelContent[0][0]) //key的名称
	allstr := fmt.Sprintf(string(buf),structName,keyType,structName,
		getCsharpClassFieldStr(excelContent),
		structName,addrName,
		structName,structName,
	sb.String(),keyName,keyName,keyName)
	//没有变化就不写
	if strings.Compare(allstr,oldStr) == 0 {
		return
	}
	// 创建c#文件
	fs, err := os.Create(filename)
	defer fs.Close()
	if xutil.IsError(err) {
		return
	}
	 // 1
	fs.WriteString(allstr)
}

// 获取csharp字段字符串
func getCsharpClassFieldStr(excelContent [][]string) string {
	var sb = new(strings.Builder)
	colCount := len(excelContent[0])
	for i := 0; i < colCount; i++ {
		// public int id = 0;
		csharptypename := CsvTypeToCSharpType(excelContent[1][i])
		defaltVal := GetDefalutValByCsvType(excelContent[1][i])
		oneFieldstr := fmt.Sprintf("\tpublic %s %s %s  %s \n",csharptypename,xutil.Capitalize(excelContent[0][i]),defaltVal,
		fmt.Sprint("//",excelContent[2][i], " 字段名称  ", excelContent[0][i]))
		sb.WriteString(oneFieldstr)
	}
	return  sb.String()
}

//获取c# 默认值
func GetDefalutValByCsvType(csvtype string) string {
	switch csvtype {
	case "int","int8","uint8","int16","uint16","int32","uint32":
		return " = 0;"
	case "string":
		return ` = "";`
	case "bool":
		return " = false;"
	case "float64":
		return " = 0.0;"
	case "[]int":
		return " = null;"
	case "[]string":
		return " = null;"
	}
	return "null;"
}

//Csv类型转换未C# 类型
func CsvTypeToCSharpType(csvType string) string {
	switch csvType {
	case "int":
		return "int"
	case "int8":
		return "byte"
	case "uint8":
		return "byte"
	case "int16":
		return "short"
	case "uint16":
		return "ushort"
	case "int32":
		return "int"
	case "uint32":
		return "uint"
	case "int64":
		return "long"
	case "uint64":
		return "ulong"
	case "float64":
		return "double"
	case "string", "bool":
		return csvType
	case "[]int":
		return "int[]"
	case "[]string":
		return "string[]"
	}
	return ""
}


func GetCsharpConvertFun(fliedtype string) string {
	parsestr := ".Parse(oneRow[col++]);"
	switch fliedtype {
	case "int":
		return "= int" + parsestr
	case "int8":
		return "= byte" + parsestr
	case "uint8":
		return "= byte"+ parsestr
	case "int16":
		return "= short"+ parsestr
	case "uint16":
		return "= ushort"+ parsestr
	case "int32":
		return "= int"+ parsestr
	case "uint32":
		return "= uint"+ parsestr
	case "float64":
		return "= double"+ parsestr
	case "string":
		return "= oneRow[col++];"
	case "bool":
		return "= bool"+ parsestr
	case "[]int","int[]":
		return "= StringUtil.StringToIntArr(oneRow[col++])"
	case "[]string","string[]":
		return "= oneRow[col++].Split(',')"
	default:
		fmt.Println(fliedtype, "is an unknown type.")
		return ""
	}
	return ""
}