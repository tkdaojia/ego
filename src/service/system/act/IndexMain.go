package service

import (
	web "ego/src/service/system/do"
	"github.com/gin-gonic/gin"
)

func RunSysIndexMainGet(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "welcome":
		web.RunWelcome(c)
	case "wait":
		web.RunWait(c)
	case "other":
		web.RunIndexOther(c)
	case "editpwd":
		web.RunEditPwd(c)
	case "choose":
		web.RunChoose(c)
	case "logout":
		web.RunSystemLogOut(c)
	default:
		web.RunIndexDefault(c)

	}
}

func RunSysIndexMainPost(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "editpwd":
		web.RunEditPwdPost(c)
	}
}
