package ego

import (
	"ego/src/model/response"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*gin.Context)

var (
	getRoutes  = make(map[string]HandlerFunc)
	postRoutes = make(map[string]HandlerFunc)
)

func init() {
	//AutoStart

	//AutoEnd12
}

func RegisterGetRoute(module string, handler HandlerFunc) {
	getRoutes[module] = handler
}
func RegisterPostRoute(module string, handler HandlerFunc) {
	postRoutes[module] = handler
}

func ApiZrcodeGet(c *gin.Context) {
	moduleName := c.Query("module")
	handlerFunc, ok := getRoutes[moduleName]
	if !ok {
		response.OkWithMessage("路由未加载", c)
		return
	}
	handlerFunc(c)
}

func ApiZrcodePost(c *gin.Context) {
	moduleName := c.Query("module")
	handlerFunc, ok := postRoutes[moduleName]
	if !ok {
		response.OkWithMessage("路由未加载2", c)
		return
	}
	handlerFunc(c)
}
