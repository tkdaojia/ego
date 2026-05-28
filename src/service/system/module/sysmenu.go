package service

import (
	service "ego/src/service/system/act"
	"github.com/gin-gonic/gin"
)

func RunSysmenu(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysMenuMainGet(c)
	}
}

func RunSysMenuPost(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysMenuMainPost(c)
	}
}
