package service

import (
	service "ego/src/service/system/act"
	"github.com/gin-gonic/gin"
)

func RunSysGroup(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysGroupMainGet(c)
	}
}

func RunSysGroupPost(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysGroupMainPost(c)

	}
}
