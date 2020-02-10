/*
创建时间: 2020/2/7
作者: zjy
功能介绍:

*/

package excletosome

import (
	"github.com/widuu/goini"
)

var conf *goini.Config

var InPath string
var Intype string
var csvoutdir string
var gostrcutdir string
var sqloutdir string
var packagename string
var outType string
var CsvPath string

func ReadConf() {
	conf = goini.SetConfig( "conf.ini")
	conf.ReadList()
	InPath = conf.GetValue("infile", "inpath")
	Intype = conf.GetValue("infile", "intype")
	csvoutdir = conf.GetValue("csv", "outpath")
	gostrcutdir = conf.GetValue("gostrcut", "outpath")
	packagename =conf.GetValue("gostrcut", "packagename")
	CsvPath = conf.GetValue("gostrcut", "csvpath")
	sqloutdir = conf.GetValue("sql", "outpath")
}
