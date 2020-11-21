/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"encoding/json"
	"fmt"
	"github.com/zjytra/wengo/xutil/osutil"
	"os"
	"path"
)

//字典转字符串
func map2JsonStr(dataDict map[string]interface{}) string {
	mjson, err := json.Marshal(dataDict)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(mjson)
}



//写JSON文件
func writeJSON(paths string, fileName string, dataDict map[string]interface{}) {
	if !osutil.MakeDirAll(paths) {
		return
	}
	realpath := path.Join(paths,fileName+".json")
	file, err := os.OpenFile(realpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) //不存在创建清空内容覆写
	if err != nil {
		fmt.Println("open file failed.", err.Error())
		return
	}
	defer file.Close()
	//字典转字符串
	file.WriteString(map2JsonStr(dataDict))
}

