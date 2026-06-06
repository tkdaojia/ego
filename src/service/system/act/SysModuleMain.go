package service

import (
	web "ego/src/service/system/do"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunSysModuleMainGet(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "list":
		web.RunSysmoduleList(c)
	case "add":
		web.RunSysmoduleAdd(c)
	case "getdata":
		web.RunSysmoduleGetdata(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}

func RunSysModuleMainPost(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "create":
		web.RunSysmoduleCrate(c)
	case "update":
		web.RunSysmoduleUpdate(c)
	case "del":
		web.RunSysmoduleDel(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}
