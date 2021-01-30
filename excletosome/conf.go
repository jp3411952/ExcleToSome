/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"encoding/json"
	"os"
)

var Conf *ConfJson


type ConfJson struct {
	InPath string `json:InPath`
	Intype string `json:Intype`
	OutType string `json:OutType`
	ServerOutCsv string   `json:ServerOutCsv`//服务器cvs
	GostrcutOutpath string `json:GostrcutOutpath`
	GopakName string  `json:GopakName`
	ClientOutCsv string  `json:ClientOutCsv`//客户端csv
	CsharpOut string  `json:CsharpOut`//c#输出目录
	LuaTOutPath string `json:LuaTOutPath`
}



func ReadConfJson() {
	Conf = new(ConfJson)
	filePtr, err := os.Open("conf.json")
	if err != nil {
		return
	}
	defer filePtr.Close()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(Conf)
	if err != nil {
		return
	}

}
