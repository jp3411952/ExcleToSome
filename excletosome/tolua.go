/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

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

var LuaCfg []string

//写Lua文件
func writeLuaTable(exclefileName string, celContent [][]string) {
	defer Wg.Done()

	//拥有平台对应列的的数量
	excelContent := GetPlatfCol(celContent, "L")
	if len(excelContent) == 0 {
		return
	}
	if !osutil.MakeDirAll(Conf.LuaTOutPath) {
		fmt.Printf("LuaTOutPath = ", Conf.LuaTOutPath)
		return
	}
	if excelContent == nil {
		fmt.Printf("excle 数据为空")
		return
	}
	filename := filepath.Join(Conf.LuaTOutPath, exclefileName+".lua.txt")
	filebytes, _ := ioutil.ReadFile(filename)
	oldStr := string(filebytes) //添加文件是否变化的验证
	structName := xutil.Capitalize(exclefileName)
	//写总管理文件
	LuaCfg = append(LuaCfg, exclefileName)
	newStr := new(strings.Builder)
	writeLuaHead(newStr, excelContent, structName)
	writeData(newStr, excelContent)
	writeFun(newStr, structName)
	newStr.WriteString("return " + structName)
	var writeStr = newStr.String()
	if strings.Compare(oldStr, writeStr) == 0 {
		return
	}
	fs, err := os.Create(filename)
	defer fs.Close()
	if xutil.IsError(err) {
		return
	}
	// 一次性写入所有数据
	fs.WriteString(writeStr)
}

//写头部
func writeLuaHead(sb *strings.Builder, excelContent [][]string, structName string) {
	nameRow := excelContent[0]
	typeRow := excelContent[1]
	noteRow := excelContent[2]
	sb.WriteString("local Debug = Debug\n")
	sb.WriteString("---go工具生成的配置文件不要修改")
	sb.WriteString("\n---@class " + structName + "\n")
	for i, s := range nameRow {
		sb.WriteString("---@" + s + " " + typeToLuaType(typeRow[i]) + " @ " + noteRow[i] + "\n")
	}
	sb.WriteString("local  " + structName + " =")
}

func writeData(sb *strings.Builder, excelContent [][]string) {
	sb.WriteString(" {\n")
	nameRow := excelContent[0]
	typeRow := excelContent[1]
	var rowLen = len(excelContent)
	for i, row := range excelContent {
		if i < 3 { //排除前三行
			continue
		}
		var colLen = len(row)
		for col, s := range row {
			if col == 0 { // 写key
				sb.WriteString("\t[" + s + "] = { ")
			}
			writeLuaDataByType(sb, nameRow[col], s, typeRow[col])
			if col != colLen-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(" }")
		if i != rowLen-1 { //不是最后一行写逗号
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
}

func writeLuaDataByType(sb *strings.Builder, fieldName string, data string, ctype string) {
	sb.WriteString(fieldName + " = ")
	switch ctype {
	case "int", "int8", "uint8", "int16", "uint16", "int32", "uint32":
		sb.WriteString(data)
		break
	case "string":
		sb.WriteString(fmt.Sprintf(`"%v"`, data))
		break
	case "table":
		sb.WriteString(data)
		break
	case "[]int":
		sb.WriteString(fmt.Sprintf(`{%v}`, data))
		break
	case "[]string":
		strArr := strings.Split(data, ",")
		var strs string
		var strArrLn = len(strArr)
		for i, s := range strArr {
			strs += fmt.Sprintf(`"%v"`, s)
			if i != strArrLn-1 {
				strs += ","
			}
		}
		sb.WriteString(fmt.Sprintf(`{%v}`, strs))
		break
	}

}

func writeFun(sb *strings.Builder, structName string) {
	sb.WriteString(fmt.Sprintf(`
---只读表
Local_text.__newindex = function(table, key, value)
	
end

---@return %s
function %s.GetCfg(_id)
	local data = %s[_id]
	if data == nil then
		Debug.Error("没有找到对应的配置项=", _id)
		return nil
	end
	return data
end

`, structName, structName, structName))
}

func typeToLuaType(ctype string) string {
	switch ctype {
	case "int", "int8", "uint8", "int16", "uint16", "int32", "uint32":
		return "number"
	case "string":
		return ctype
	case "[]int":
		return "number[]"
	case "[]string":
		return "string[]"
	case "table":
		return "table"
	}
	return "nil"
}

func WriteCfgMgr() {
	if !osutil.MakeDirAll(Conf.LuaTOutPath) {
		fmt.Printf("LuaTOutPath = ", Conf.LuaTOutPath)
		return
	}
	//没有可写的文件
	if LuaCfg == nil {
		return
	}

	var writeStr = new(strings.Builder)
	writeStr.WriteString("---go工具生成的配置文件不要修改\n")
	writeStr.WriteString("---@class CfgMgr: _Class\n")

	for _, s := range LuaCfg {
		writeStr.WriteString("---@field public " + s + " " + s + "\n")
	}

	writeStr.WriteString("local CfgMgr = Class(\"CfgMgr\")\n\n")
	writeStr.WriteString("function CfgMgr:Ctor(...)\n")
	for _, s := range LuaCfg {
		writeStr.WriteString("\tself." + s + " = require(\"LuaCfg/" + s + "\")\n")
	}
	writeStr.WriteString("end")

	writeStr.WriteString("\n\nreturn CfgMgr")

	filename := filepath.Join(Conf.LuaTOutPath, "CfgMgr.lua.txt")
	fs, err := os.Create(filename)
	defer fs.Close()
	if xutil.IsError(err) {
		return
	}
	// 一次性写入所有数据
	fs.WriteString(writeStr.String())
}
