package auditLog

import (
	"ego/src/boot/global" // 换成你实际的 global 路径
	model "ego/src/model/basic"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"time"
)

type AuditLog struct{}

// SaveAuditLog 异步保存审计日志
// description: 操作描述 (如 "修改了用户资料")
// param: 请求参数 (支持任何结构体/Map/常规类型)
// oldData: 变更前的数据 (没有传 nil)
// newData: 变更后的数据 (没有传 nil)
func (auditLog *AuditLog) SaveAuditLog(c *gin.Context, description string, param, oldData, newData any) {
	endTime := time.Now()
	// 在主线程立即复制 Context，防止闭包逃逸后 Context 被 Gin 释放
	ctxCopy := c.Copy()

	var latency int64 = 0
	if startTimeVal, exists := ctxCopy.Get("request_start_time"); exists {
		if startTime, ok := startTimeVal.(time.Time); ok {
			// 结束时间 减去 开始时间，并转为毫秒 (ms)
			latency = endTime.Sub(startTime).Milliseconds()
		}
	}

	// 从 Context 中获取当前登录用户 (这里根据你实际的 Auth 中间件调整)
	username := "guest"
	nickname := "guest"
	if user, exists := ctxCopy.Get("User_Account"); exists {
		username = user.(string)
	}
	if user, exists := ctxCopy.Get("User_Name"); exists {
		nickname = user.(string)
	}

	// 3. 开启 Goroutine 异步执行序列化和入库，不阻塞业务
	go func() {
		defer func() {
			if r := recover(); r != nil {
				global.C_LOG.Error("审计日志异步写入崩溃捕获", zap.Any("recover", r))
			}
		}()

		module := strings.TrimSpace(ctxCopy.Query("module"))
		action := strings.TrimSpace(ctxCopy.Query("act"))
		do := strings.TrimSpace(ctxCopy.Query("do"))

		// 统一将参数/新旧数据转为 JSON 字符串
		paramStr := convertToJSONStr(param)
		oldDataStr := convertToJSONStr(oldData)
		newDataStr := convertToJSONStr(newData)

		logRecord := model.SysOperationLog{
			Username:    username,
			Nickname:    nickname,
			Module:      module,
			Action:      action,
			Do:          do,
			Description: description,
			Url:         ctxCopy.Request.RequestURI,
			Method:      ctxCopy.Request.Method,
			Ip:          ctxCopy.ClientIP(),
			UserAgent:   ctxCopy.Request.UserAgent(),
			Param:       paramStr,
			DataOld:     oldDataStr,
			DataNew:     newDataStr,
			Status:      1,
			Latency:     latency,
		}

		// 使用全局无 Context 的 DB 写入数据库
		if err := global.C_DB.Create(&logRecord).Error; err != nil {
			global.C_LOG.Error("写入审计日志失败", zap.Error(err))
		}
	}()
}

func (auditLog *AuditLog) SaveAuditLogError(c *gin.Context, errText string) {
	endTime := time.Now()
	// 在主线程立即复制 Context，防止闭包逃逸后 Context 被 Gin 释放
	ctxCopy := c.Copy()

	var latency int64 = 0
	if startTimeVal, exists := ctxCopy.Get("request_start_time"); exists {
		if startTime, ok := startTimeVal.(time.Time); ok {
			// 结束时间 减去 开始时间，并转为毫秒 (ms)
			latency = endTime.Sub(startTime).Milliseconds()
		}
	}

	// 从 Context 中获取当前登录用户 (这里根据你实际的 Auth 中间件调整)
	username := "guest"
	nickname := "guest"
	if user, exists := ctxCopy.Get("User_Account"); exists {
		username = user.(string)
	}
	if user, exists := ctxCopy.Get("User_Name"); exists {
		nickname = user.(string)
	}

	// 3. 开启 Goroutine 异步执行序列化和入库，不阻塞业务
	go func() {
		defer func() {
			if r := recover(); r != nil {
				global.C_LOG.Error("审计日志异步写入崩溃捕获", zap.Any("recover", r))
			}
		}()

		module := strings.TrimSpace(ctxCopy.Query("module"))
		action := strings.TrimSpace(ctxCopy.Query("act"))
		do := strings.TrimSpace(ctxCopy.Query("do"))

		logRecord := model.SysOperationLog{
			Username:  username,
			Nickname:  nickname,
			Module:    module,
			Action:    action,
			Do:        do,
			Url:       ctxCopy.Request.RequestURI,
			Method:    ctxCopy.Request.Method,
			Ip:        ctxCopy.ClientIP(),
			UserAgent: ctxCopy.Request.UserAgent(),
			Status:    2,
			ErrorMsg:  errText,
			Latency:   latency,
		}

		// 使用全局无 Context 的 DB 写入数据库
		if err := global.C_DB.Create(&logRecord).Error; err != nil {
			global.C_LOG.Error("写入审计日志失败", zap.Error(err))
		}
	}()
}

// 辅助函数：将任意对象转为 JSON 字符串，若为空或失败则返回空字符串
func convertToJSONStr(v any) string {
	if v == nil {
		return ""
	}
	// 如果本身就是字符串，直接返回，避免二次转义
	if str, ok := v.(string); ok {
		return str
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}
