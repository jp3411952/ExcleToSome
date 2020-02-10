//excle生成文件请勿修改
package csvdata

import (
	"github.com/showgo/csvparse"
	"github.com/showgo/xutil"
)

var testCsv map[int]*Test

type  Test struct {
	Id int //#id 字段名称  id
	Intarr []int //int数组 字段名称  intarr
	Strarr []string //字符串数组 字段名称  strarr
	Floatt float64 //浮点类型 字段名称  floatt
}

func AsynSetTestMapData() {
	 go SetTestMapData()
}

func SetTestMapData() {
    if testCsv == nil {
		testCsv = make(map[int]*Test)
	}
	tem := getTestUsedData("./csv/")
	testCsv  = tem
}

func getTestUsedData(csvpath  string ) map[int]*Test{
    csvmapdata := csvparse.GetCsvMapData(csvpath + "test.csv")
	tem := make(map[int]*Test)
	for _, filedData := range csvmapdata {
		one := new(Test)
		for filedName, filedval := range filedData {
			isok := csvparse.SetFieldReflect(one, filedName, filedval)
			xutil.IsError(isok)
		}
		tem[one.Id] = one
	}
	return tem
}

func GetTestPtr(id int) *Test{
	return testCsv[id]
}
