package service

import (
	service "ego/src/service/system/act"
	"github.com/gin-gonic/gin"
)

func RunSysparameter(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysDictMainGet(c)

	}
}

func RunSysparameterPost(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "main":
		service.RunSysDictMainPost(c)

	}
}
