/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"fmt"
	"github.com/wengo/xutil/osutil"
	"os"
	"path"
)

//写Lua文件
func writeLuaTable(paths string, fileName string, dataDict interface{}) {
	if !osutil.MakeDirAll(paths) {
		return
	}
	realpath := path.Join(paths,fileName+".lua")
	file, err := os.Create(realpath) //不存在创建清空内容覆写
	if err != nil {
		fmt.Println("open file failed.", err.Error())
		return
	}
	defer file.Close()
	writeLuaTableContent(file, dataDict, 0)
	file.WriteString("return ")
}

//写Lua表内容
func writeLuaTableContent(fileHandle *os.File, data interface{}, idx int) {
	switch t := data.(type) {
	case int:
		fileHandle.WriteString(fmt.Sprintf("%v",data)) //对于interface{}, %v会打印实际类型的值
	case float64:
		fileHandle.WriteString(fmt.Sprintf("%v",data)) //对于interface{}, %v会打印实际类型的值
	case string:
		fileHandle.WriteString(fmt.Sprintf(`"%s"`, data)) //对于interface{}, %v会打印实际类型的值
	case []interface{}:
		fileHandle.WriteString("{\n")
		a := data.([]interface{})
		for _, v := range a {
			addTabs(fileHandle, idx)
			writeLuaTableContent(fileHandle, v, idx+1)
			fileHandle.WriteString(",\n")
		}
		addTabs(fileHandle, idx-1)
		fileHandle.WriteString("}")
	case []string:
		fileHandle.WriteString("{\n")
		a := data.([]string)
		for _, v := range a {
			addTabs(fileHandle, idx)
			writeLuaTableContent(fileHandle, v, idx+1)
			fileHandle.WriteString(",\n")
		}
		addTabs(fileHandle, idx-1)
		fileHandle.WriteString("}")
	
	case map[string]interface{}:
		m := data.(map[string]interface{})
		fileHandle.WriteString("{\n")
		for k, v := range m {
			addTabs(fileHandle, idx)
			fileHandle.WriteString("[")
			writeLuaTableContent(fileHandle, k, idx+1)
			fileHandle.WriteString("] = ")
			writeLuaTableContent(fileHandle, v, idx+1)
			fileHandle.WriteString(",\n")
		}
		addTabs(fileHandle, idx-1)
		fileHandle.WriteString("}")
	default:
		fileHandle.WriteString(fmt.Sprintf("%t", data))
		_ = t
	}
}

//在文件中添加制表符
func addTabs(fileHandle *os.File, idx int) {
	for i := 0; i < idx; i++ {
		fileHandle.WriteString("\t")
	}
}
