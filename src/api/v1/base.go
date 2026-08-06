package v1

import (
	"ego/src/model/response"
	service "ego/src/service/system/module"
	"github.com/gin-gonic/gin"
)

func ApiBaseGet(c *gin.Context) {
	module := c.Query("module")
	switch module {
	case "index":
		service.RunIndex(c)
	case "sysdict":
		service.RunSysparameter(c)
	case "sysmodule":
		service.RunSysmodule(c)
	case "sysmenu":
		service.RunSysmenu(c)
	case "sysrole":
		service.RunSysRole(c)
	case "sysaccount":
		service.RunSysAccount(c)
	case "sysgroup":
		service.RunSysGroup(c)
	case "syslog":
		service.RunSysLog(c)
	default:
		response.OkWithMessage("Hi", c)
	}
}

func ApiBasePost(c *gin.Context) {

	module := c.Query("module")
	switch module {
	case "index":
		service.RunIndexPost(c)
	case "sysdict":
		service.RunSysparameterPost(c)
	case "sysmodule":
		service.RunSysmodulePost(c)
	case "sysmenu":
		service.RunSysMenuPost(c)
	case "sysrole":
		service.RunSysRolePost(c)
	case "sysaccount":
		service.RunSysAccountPost(c)
	case "sysgroup":
		service.RunSysGroupPost(c)
	default:
		response.OkWithMessage("创建成功1", c)
	}
}

func ApiOpen(c *gin.Context) {
	module := c.Query("module")
	switch module {
	case "login":
		service.RunLogin(c)
	case "out":
		c.HTML(200, "open/out.htm", gin.H{})
	case "autologin":
		service.RunAutoLogin(c)
	default:
		response.OkWithMessage("Hi1", c)
	}
}
func ApiOpenPost(c *gin.Context) {
	module := c.Query("module")
	switch module {
	case "login":
		service.RunLoginPost(c)
	default:
		response.OkWithMessage("Hi", c)
	}
}
