package service

import (
	service "ego/src/service/system/act"
	"github.com/gin-gonic/gin"
)

func RunSysAccount(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysAccountMainGet(c)

	}
}

func RunSysAccountPost(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysAccountMainPost(c)
	}
}
