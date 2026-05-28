package service

import (
	web "ego/src/service/system/do"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunSysGroupMainGet(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "list":
		web.RunSysGroupMainList(c)
	case `add`:
		web.RunSysGroupMainAdd(c)
	case `getdata`:
		web.RunSysGroupsMainGetdata(c)
	case `del`:
		web.RunSysGroupMainDel(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}

func RunSysGroupMainPost(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case `create`:
		web.RunSysGroupCreate(c)
	case `update`:
		web.RunSysGroupUpdate(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}
