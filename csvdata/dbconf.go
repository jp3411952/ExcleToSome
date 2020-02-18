//excle生成文件请勿修改
package csvdata

import (
	"fmt"
	"github.com/showgo/csvparse"
	"github.com/showgo/xutil"
)

var dbconfCsv map[string]*Dbconf

type  Dbconf struct {
	Dbname string //#数据库名称 字段名称  dbname
	Ip string //ip地址 字段名称  ip
	Dbport string //端口号 字段名称  dbport
	Dbusername string //用户名 字段名称  dbusername
	Dbpwd string //密码 字段名称  dbpwd
	Maxopenconns int //最大链接数 字段名称  maxopenconns
	Maxidleconns int //闲置连接数 字段名称  maxidleconns
}

func SetDbconfMapData(csvpath  string ) {
    if dbconfCsv == nil {
		dbconfCsv = make(map[string]*Dbconf)
	}
	tem := getDbconfUsedData(csvpath)
	dbconfCsv  = tem
}

func getDbconfUsedData(csvpath  string ) map[string]*Dbconf{
    csvmapdata := csvparse.GetCsvMapData(csvpath + "dbconf.csv")
	tem := make(map[string]*Dbconf)
	for _, filedData := range csvmapdata {
		one := new(Dbconf)
		for filedName, filedval := range filedData {
			isok := csvparse.SetFieldReflect(one, filedName, filedval)
			xutil.IsError(isok)
			if _,ok := tem[one.Dbname]; ok {
				fmt.Println(one.Dbname,"重复")
			}
		}
		tem[one.Dbname] = one
	}
	return tem
}

func GetDbconfPtr(dbname string) *Dbconf{
    if _,ok := dbconfCsv[dbname]; !ok  {
		return nil
	}
	return dbconfCsv[dbname]
}
