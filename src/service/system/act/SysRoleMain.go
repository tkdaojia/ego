package service

import (
	web "ego/src/service/system/do"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunSysRoleMainGet(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "list":
		web.RunSysRoleList(c)
	case "add":
		web.RunSysRoleAdd(c)
	case "del":
		web.RunSysRoleDel(c)
	case "getdata":
		web.RunSysRoleGetdata(c)
	case "menutree":
		web.RunSysRoleMenuTree(c)
	case "moduletree":
		web.RunSysRoleModuleTree(c)
	case "addmenu":
		web.RunSysRoleAddMenu(c)
	case "addmodule":
		web.RunSysRoleAddModule(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}

func RunSysRoleMainPost(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "create":
		web.RunSysRoleCreate(c)
	case "update":
		web.RunSysRoleUpdate(c)
	case "updatemenu":
		web.RunSysRoleUpdateMenu(c)
	case "updatemodule":
		web.RunSysRoleUpdateModule(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}
