/*
创建时间: 2020/2/8
作者: zjy
功能介绍:

*/

package main

import (
	"fmt"
	"github.com/showgo/csvparse"
	"github.com/showgo/xutil"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"./gofile"
)

// 根据配置表生成, 这里最好与结构体一起生成

func main() {
	
  csvdata.SetServerconfCsvMapData()
  fmt.Println("%t",csvdata.ServerconfCsv)
}

func SliceParse() {
	// 需要把csv的数据,映射到map
	csvdt := csvparse.GetCsvSliceData("./csv/test.csv")
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 2048)
			l := runtime.Stack(buf, false)
			fmt.Printf("%v %s", r, buf[:l])
		}
	}()
	fileds := csvdt[0]
	// 字段名称 为了验证是否是这个歌数据
	types := csvdt[1]
	// 类型
	for i := 3; i < len(csvdt); i++ {
		one := new(csvdata.Test)
		valueof := reflect.ValueOf(one)
		for j, csvval := range csvdt[i] {
			temval := valueof.Elem()
			filedName := xutil.Capitalize(fileds[j])
			fliedval := temval.FieldByName(filedName)
			if fliedval.IsValid() {
				setval := CsvStrToInterfaceIType(fliedval.Interface(), csvval)
				if setval == nil {
					fmt.Println("类型=", types[j], "数据=", csvval, "不能转换")
					continue
				}
				isok := csvparse.SetFieldReflect(one, filedName, setval)
				xutil.IsError(isok)
			}
		}
		fmt.Println("one", one)
	}
}


// IType版
func CsvStrToInterfaceIType(fliedval interface{}, strval string) interface{} {
	switch fliedval.(type) {
	case int:
		inval, erro := strconv.Atoi(strval)
		if !xutil.IsError(erro) {
			return inval
		}
	case float64:
		flt64, erro := strconv.ParseFloat(strval, 64)
		if !xutil.IsError(erro) {
			return flt64
		}
	case string:
		return strval
	case []int:
		intArr := csvparse.StringsToIntArr(strval)
		return intArr
	case []string:
		return strings.Split(strval, ",")
	default:
		fmt.Println(fliedval, "is an unknown type.")
		return nil
	}
	
	return nil
}


