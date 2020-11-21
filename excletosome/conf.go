/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"ExcleToSome/widuu/goini"
)

var conf *goini.Config

var InPath string
var Intype string
var OutType string
var csvoutdir string
var gostrcutdir string
var sqloutdir string
var packagename string
var outType string
var CsvPath string
var jsoninpath string
var jsonoutpath string
var jsonPackageName string

func ReadConf() {
	conf = goini.SetConfig( "conf.ini")
	conf.ReadList()
	InPath = conf.GetValue("transfrom", "inpath")
	Intype = conf.GetValue("transfrom", "intype")
	OutType = conf.GetValue("transfrom", "outtype")
	csvoutdir = conf.GetValue("csv", "outpath")
	gostrcutdir = conf.GetValue("gostrcut", "outpath")
	packagename =conf.GetValue("gostrcut", "packagename")
	CsvPath = conf.GetValue("gostrcut", "csvpath")
	sqloutdir = conf.GetValue("sql", "outpath")
	jsonoutpath = conf.GetValue("json","outpath")
	jsonPackageName = conf.GetValue("json","packagename")
}

