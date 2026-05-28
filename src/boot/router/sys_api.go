package router

import (
	"ego/src/api/app"
	"ego/src/api/ego"
	v1 "ego/src/api/v1"
	"ego/src/boot/global"
	"ego/src/boot/middleware"
	"github.com/gin-gonic/gin"
)

func InitApiRouter(Router *gin.RouterGroup) {
	m := global.C_CONFIG.System

	ZrRouter := Router.Group("/" + m.Webroute).Use(middleware.VerifyAuth())
	{
		ZrRouter.GET("/", ego.ApiZrcodeGet)
		ZrRouter.POST("/", ego.ApiZrcodePost)
	}

	OpenRouter := Router.Group("/open")
	{
		OpenRouter.GET("/", v1.ApiOpen)
		OpenRouter.POST("/", v1.ApiOpenPost)
	}

	AdminRouter := Router.Group("/admin").Use(middleware.VerifyAuth())
	{
		AdminRouter.GET("/", v1.ApiBaseGet)
		AdminRouter.POST("/", v1.ApiBasePost)
	}

	AppRouter := Router.Group("/app").Use(middleware.VerifyAuthApp())
	{
		AppRouter.GET("/", app.ApiAppGet)
		AppRouter.POST("/", app.ApiAppPost)
	}

}
