package service

import (
	service "ego/src/service/system/act"
	"github.com/gin-gonic/gin"
)

func RunIndex(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysIndexMainGet(c)
	}
}

func RunIndexPost(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysIndexMainPost(c)
	}
}
