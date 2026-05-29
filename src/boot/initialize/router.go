package initialize

import (
	"ego/src/boot/middleware"
	"ego/src/boot/router"
	"ego/src/utils"
	"ego/src/utils/backrun"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

// 初始化总路由
func Routers() *gin.Engine {
	var Router = gin.New()

	Router.Use(middleware.GinRecovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}

	Router.SetFuncMap(template.FuncMap{
		"inarray": utils.BoolInarray,
	})
	Router.Delims("<?", "?>")

	Router.LoadHTMLGlob("web/**/*")
	Router.StaticFS("/static", http.Dir("./static"))

	Router.Use(middleware.Cors()) //跨域
	PublicGroup := Router.Group("")
	router.InitApiRouter(PublicGroup)

	Background()
	return Router
}

func Background() {
	backrun.BackRunQueue()
}
