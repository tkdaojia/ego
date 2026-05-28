package service

import (
	service "ego/src/service/system/act"
	"github.com/gin-gonic/gin"
)

func RunSysRole(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysRoleMainGet(c)
	}
}

func RunSysRolePost(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysRoleMainPost(c)

	}
}
