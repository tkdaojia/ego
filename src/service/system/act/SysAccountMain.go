package service

import (
	web "ego/src/service/system/do"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunSysAccountMainGet(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "list":
		web.RunSysAccountList(c)
	case "add":
		web.RunSysAccountAdd(c)

	case "getdata":
		web.RunSysAccountGetdata(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}

func RunSysAccountMainPost(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case `create`:
		web.RunSysAccountCreate(c)
	case `update`:
		web.RunSysAccountUpdate(c)
	case "del":
		web.RunSysAccountDel(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}
