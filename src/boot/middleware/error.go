package middleware

import (
	"ego/src/boot/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
)

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//  安全地获取错误描述，防止非 error 类型的 panic 导致二次崩溃
				var errMsg string
				if e, ok := err.(error); ok {
					errMsg = e.Error()
				} else {
					errMsg = fmt.Sprintf("%v", err)
				}

				//  检查是否为 broken pipe (客户端断开连接)
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取原始 HTTP 请求内容
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				//  构建基础日志字段
				logFields := []zap.Field{
					zap.String("error", errMsg),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
					zap.String("request", string(httpRequest)),
				}

				if brokenPipe {
					// 对方连接已经断开，无法写入状态码，直接记录日志并中止
					global.C_LOG.Error("[Recovery from broken pipe]", logFields...)

					if errAsError, ok := err.(error); ok {
						_ = c.Error(errAsError)
					}
					c.Abort()
					return
				}

				//  普通 Panic 默认获取堆栈信息，方便肉眼排查 Bug
				stackInfo := make([]byte, 2048)
				n := runtime.Stack(stackInfo, false)
				logFields = append(logFields, zap.String("stack", string(stackInfo[:n])))

				//  记录 Error 级别日志，并返回 500
				global.C_LOG.Error("[Recovery from panic]", logFields...)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
