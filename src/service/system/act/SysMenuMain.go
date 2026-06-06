package service

import (
	web "ego/src/service/system/do"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunSysMenuMainGet(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "list":
		web.RunSysMenuList(c)
	case "add":
		web.RunSysMenuAdd(c)
	case "getdata":
		web.RunSysMenuGetdata(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}

func RunSysMenuMainPost(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "create":
		web.RunSysMenuCrate(c)
	case "update":
		web.RunSysMenuUpdate(c)
	case "del":
		web.RunSysMenuDel(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}
