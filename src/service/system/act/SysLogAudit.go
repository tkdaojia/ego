package service

import (
	model "ego/src/model/basic"
	"ego/src/model/msg"
	"ego/src/model/response"
	"ego/src/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func RunSysLogAuditGet(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "list":
		RunSysLogAuditList(c)
	case "getdata":
		RunSysLogAuditGetdata(c)
	case "detail":
		RunSysLogAuditDetail(c)
	}
}

func RunSysLogAuditList(c *gin.Context) {
	c.HTML(200, "syslog/syslogList.htm", nil)
}

func RunSysLogAuditGetdata(c *gin.Context) {
	db := utils.GetDB(c)

	page := cast.ToInt(c.Query("page"))
	limit := cast.ToInt(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20 // 默认每页 20 条
	}

	// 2. 获取前端搜索栏传过来的过滤参数
	username := strings.TrimSpace(c.Query("username"))
	//module := strings.TrimSpace(c.Query("module"))
	status := cast.ToInt(c.Query("status"))

	// 3. 构建动态查询 Scope（作用域代理）
	query := db.Model(&model.SysOperationLog{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCountErr, err)
		response.OnFailure(c, msg.SqlCountErr)
		return
	}

	if count == 0 {
		response.OkTableList([]map[string]any{}, 0, c)
		return
	}

	var list []model.SysOperationLog
	offset := (page - 1) * limit

	err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		utils.LogSqlErr(c, msg.SqlFindListErr, err)
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	results := make([]map[string]any, len(list))
	for i, v := range list {
		results[i] = map[string]any{
			"id":          v.ID,
			"username":    v.Username,
			"nickname":    v.Nickname,
			"module":      v.Module,
			"action":      v.Action,
			"do":          v.Do,
			"description": v.Description,
			"method":      v.Method,
			"url":         v.Url,
			"ip":          v.Ip,
			"status":      v.Status,
			"latency":     v.Latency,
			"createdAt":   v.CreatedAt.Format("2006-01-02 15:04:05"), // 人类友好时间
		}
	}

	response.OkTableList(results, count, c)
}

func RunSysLogAuditDetail(c *gin.Context) {
	db := utils.GetDB(c)

	// 1. 获取并转换 URL 参数中的日志 ID
	id := cast.ToInt(c.Query("id"))
	if id < 1 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}

	// 2. 声明模型变量，精准查询单条日志记录（包含 text 大文本字段）
	var info model.SysOperationLog
	if err := db.Where("id = ?", id).First(&info).Error; err != nil {
		// 拦截：记录不存在或已被自动化脚本清理
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.OnFailure(c, "该操作日志不存在或已被清理")
			return
		}
		// 拦截：其他底层 SQL 报错
		utils.LogSqlErr(c, msg.SqlFindErr, err)
		response.OnFailure(c, msg.SqlFindErr)
		return
	}

	// 3. 渲染详情模板，并注入查出来的全量 info 对象
	c.HTML(http.StatusOK, "syslog/syslogDetail.htm", gin.H{
		"info": info,
	})
}
