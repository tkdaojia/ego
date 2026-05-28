package service

import (
	service "ego/src/service/system/act"
	"github.com/gin-gonic/gin"
)

func RunSysLog(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "audit":
		service.RunSysLogAuditGet(c)
	}
}
