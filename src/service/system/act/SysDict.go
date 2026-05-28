package service

import (
	web "ego/src/service/system/do"
	"github.com/gin-gonic/gin"
)

func RunSysDictMainGet(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "list":
		web.RunSysDictList(c)
	case "add":
		web.RunSysDictAdd(c)
	case "del":
		web.RunSysDictDel(c)
	case "getdata":
		web.RunSysDictGetdata(c)
	case "set":
		web.RunSysDictSet(c)
	case "get_setdata":
		web.RunSysDictGetSetdata(c)
	case "get_setdata_table":
		web.RunSysDictGetSetdataTable(c)

	case "svngetdata":
		web.RunSysDictSvnGetdata(c)
	}
}
func RunSysDictMainPost(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "create":
		web.RunSysDictCreate(c)
	case "update":
		web.RunSysDictUpdate(c)
	case "set":
		web.RunSysDictSetPost(c)
	}
}
