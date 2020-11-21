/*
创建时间: 2020/2/8
作者: zjy
功能介绍:
生成编程对应的编程语言文件
*/

package generatepgl

import (
	"fmt"
	"github.com/wengo/xutil/strutil"
	"strings"
)

const(
	GoPkg = "package "
	GoImport = "import "
	GoType = "type "
	Gofunc = "func "
	Gostruct = "struct "
	Gointerface = "interface "
	GoConst = "const "
	GoString =" string "
	GoReturn = "return "
	GoVar = "var "
	Goassig = " := "
	NextLine = "\n"
	Tab ="\t"
	Space = " "
	Leftbrace="{\n"
	Rightbrace="}\n"
	Leftcurves="("
	Rightcurves=")"
	CommentLine="//"
	Euqual = " = "
)

type GenerGoLang struct {
	strbd  strings.Builder
}

type FiledInfo struct {
	FiledName string  //字段名称
	FileType string   //字段类型
	FiledComment string //字段注释
}

type FuncInfo struct {
	FuncName string  //方法名称
	FuncParam map[string]string  //方法参数名称及参数类型
	FuncReturn map[string]string  //返回值名称及类型
	FuncComment string  //方法注释
	FuncContent string  //方法注释
}

func NewFuncInfo(funcName string) *FuncInfo  {
	return &FuncInfo{
		funcName,
		make(map[string]string),
		make(map[string]string),
		"",
		"",
	}
}

func NewGenerGoLang() *GenerGoLang  {
	return new(GenerGoLang)
}

func (ggl *GenerGoLang)WriteHead(comment string,pkgname string)  {
	ggl.StringBuilderWrite(comment)
	ggl.StringBuilderWrite(GoPkg)
	ggl.StringBuilderWrite(pkgname)
	ggl.WriteNextLine()
}

func (ggl *GenerGoLang)WriteImport(pkgnames []string)  {
	pkgnum := len(pkgnames)
	if pkgnum <= 0 {
		return
	}
	ggl.StringBuilderWrite(GoImport)
	if pkgnum > 1 {
		ggl.StringBuilderWrite(Leftcurves)
		ggl.WriteNextLine()
	}
	for _,pkgname := range pkgnames {
		ggl.StringBuilderWrite(Tab)
		ggl.StringBuilderWrite(pkgname)
		ggl.WriteNextLine()
	}
	if pkgnum > 1 {
		ggl.StringBuilderWrite(Rightcurves)
	}
	ggl.WriteNextLine()
	ggl.WriteNextLine()
}

func (ggl *GenerGoLang)WriteConst(constName string,val string)  {
	ggl.StringBuilderWrite(GoConst)
	ggl.StringBuilderWrite(constName)
	ggl.StringBuilderWrite(Euqual)
	ggl.StringBuilderWrite(fmt.Sprintf("\"%s\"",val))
	ggl.WriteNextLine()
}

func (ggl *GenerGoLang)WriteVar(varName string, typename string)  {
	ggl.StringBuilderWrite(GoVar)
	ggl.StringBuilderWrite(varName)
	ggl.writeSapce()
	ggl.StringBuilderWrite(typename)
	ggl.WriteNextLine()
}

func (ggl *GenerGoLang)WriteInterface(typeName string)  {
	ggl.writeType(typeName,Gointerface)
}

func (ggl *GenerGoLang)WriteStruct(typeName string)  {
	ggl.writeType(typeName,Gostruct)
}

func (ggl *GenerGoLang)WriteField(filed *FiledInfo)    {
	ggl.StringBuilderWrite(Tab)
	ggl.StringBuilderWrite(filed.FiledName)
	ggl.writeSapce()
	ggl.StringBuilderWrite(filed.FileType)
	ggl.writeSapce()
	ggl.WriteComment(filed.FiledComment)
}

func (ggl *GenerGoLang)WriteFunc(funcInfo *FuncInfo)    {
	ggl.WriteNextLine()
	//方法注释
	ggl.WriteComment(funcInfo.FuncComment)
	// 方法定义
	ggl.StringBuilderWrite(Gofunc)
	ggl.StringBuilderWrite(funcInfo.FuncName)
	ggl.StringBuilderWrite(Leftcurves)
	//方法参数
	var parmstr string
	for paramnmae, paramtype := range funcInfo.FuncParam {
		parmstr += fmt.Sprint(paramnmae, Space, paramtype, ",")
	}
	// 去掉最后一个逗号
	if len(funcInfo.FuncParam) > 0 {
		ggl.StringBuilderWrite(parmstr[:len(parmstr)-1])
	}
	ggl.StringBuilderWrite(Rightcurves)
	ggl.writeSapce()
	returnTypeCount := len(funcInfo.FuncReturn)
	if returnTypeCount > 1 { // 返回值数量大于一个的时候猜写括号
		ggl.StringBuilderWrite(Leftcurves)
	}
	returnNames := make([]string, returnTypeCount)
	returnTypes := make([]string, returnTypeCount)
	i := 0
	for returnName, returnType := range funcInfo.FuncReturn {
		returnNames[i] = returnName
		returnTypes[i] = returnType
		i++
	}
	ggl.StringBuilderWrite(strings.Join(returnTypes, ","))
	if returnTypeCount > 1 {
		ggl.StringBuilderWrite(Rightcurves)
	}
	ggl.WriteStartBrace()
	
	// 方法体内容
	if !strutil.StringIsNil(funcInfo.FuncContent) {
		ggl.StringBuilderWrite(funcInfo.FuncContent)
		ggl.WriteNextLine()
	}
	// 方法返回值
	ggl.WriteRuturn(returnNames)

	ggl.WriteEndBrace()
	//ggl.WriteNextLine()
}

//写注释
func (ggl *GenerGoLang) WriteComment(str string) {
	if !strutil.StringIsNil(str) {
		ggl.StringBuilderWrite(CommentLine)
		ggl.StringBuilderWrite(str)
		ggl.WriteNextLine()
	}
}

//写返回值
func (ggl *GenerGoLang)WriteRuturn(rets []string) {
	if len(rets) <= 0 {
		return
	}
	ggl.StringBuilderWrite(Tab)
	ggl.StringBuilderWrite(GoReturn)
	ggl.StringBuilderWrite(strings.Join(rets,","))
	ggl.WriteNextLine()
}


func (ggl *GenerGoLang)StringBuilderWrite(str string)  {
	ggl.strbd.WriteString(str)
}

func (ggl *GenerGoLang)writeType(typeName string,gotype string)  {
	ggl.StringBuilderWrite(GoType)
	ggl.writeSapce()
	ggl.StringBuilderWrite(typeName)
	ggl.writeSapce()
	ggl.StringBuilderWrite(gotype)
	ggl.WriteStartBrace()
}

func (ggl *GenerGoLang)WriteStartBrace() {
	ggl.StringBuilderWrite(Leftbrace)
}

func (ggl *GenerGoLang)WriteEndBrace() {
	ggl.StringBuilderWrite(Rightbrace)
}


func (ggl *GenerGoLang)writeSapce()  {
	ggl.StringBuilderWrite(Space)
}

func (ggl *GenerGoLang) WriteNextLine()  {
	ggl.StringBuilderWrite(NextLine)
}

func (ggl *GenerGoLang) WriteMoreNextLine(n int)  {
	for ;n > 0 ; n--  {
		ggl.WriteNextLine()
	}
}

func (ggl *GenerGoLang)String() string  {
	return ggl.strbd.String()
}