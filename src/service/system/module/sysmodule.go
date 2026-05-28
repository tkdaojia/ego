package service

import (
	service "ego/src/service/system/act"
	"github.com/gin-gonic/gin"
)

func RunSysmodule(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysModuleMainGet(c)
	}
}

func RunSysmodulePost(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysModuleMainPost(c)
	}
}
