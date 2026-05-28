package service

import (
	model "ego/src/model/basic"
	"ego/src/model/msg"
	"ego/src/model/response"
	"ego/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func RunSysGroupMainList(c *gin.Context) {
	c.HTML(200, "sysgroup/sysgroupList.htm", nil)
}

func RunSysGroupMainAdd(c *gin.Context) {
	db := utils.GetDB(c)
	id := cast.ToInt(c.Query("id"))
	var info model.SysGroup
	if id > 0 {
		if err := db.Where("id = ?", id).First(&info).Error; err != nil {
			utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
			response.ErrHtml(c, msg.SqlFindErr)
			return
		}
	}
	c.HTML(http.StatusOK, "sysgroup/sysgroupAdd.htm", gin.H{
		"info": info,
	})
}

// 创建系统组
func RunSysGroupCreate(c *gin.Context) {

	var req struct {
		Gname    string `json:"gname" binding:"required"`
		Icon     string `json:"icon"`
		Indextpl string `json:"indextpl"`
		State    string `json:"state"`
		Ordnum   string `json:"ordnum"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}
	db := utils.GetDB(c)
	item := model.SysGroup{
		Gname:    strings.TrimSpace(req.Gname),
		Icon:     strings.TrimSpace(req.Icon),
		Indextpl: strings.TrimSpace(req.Indextpl),
		State:    cast.ToInt(req.State),
		Ordnum:   cast.ToInt(req.Ordnum),
	}

	if err := db.Create(&item).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCreateErr, err)
		response.OnFailure(c, msg.SqlCreateErr)
		return
	}
	utils.Pack.AuditLog.SaveAuditLog(c, fmt.Sprintf("新增系统组 [%s]", item.Gname), req, nil, item)
	response.OnSuccess(c)
}

// 修改系统组
func RunSysGroupUpdate(c *gin.Context) {
	var req struct {
		Id       int    `json:"id" binding:"required,gt=0"`
		Gname    string `json:"gname" binding:"required"`
		Icon     string `json:"icon"`
		Indextpl string `json:"indextpl"`
		State    string `json:"state"`
		Ordnum   string `json:"ordnum"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	db := utils.GetDB(c)
	id := cast.ToInt(req.Id)
	var old model.SysGroup
	if err := db.Where("id = ?", id).First(&old).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlFindErr)
		return
	}

	updateData := map[string]any{
		"gname":    strings.TrimSpace(req.Gname),
		"icon":     strings.TrimSpace(req.Icon),
		"indextpl": strings.TrimSpace(req.Indextpl),
		"state":    cast.ToInt(req.State),
		"ordnum":   cast.ToInt(req.Ordnum),
	}
	if err := db.Model(&old).Updates(updateData).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlUpdateErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlUpdateErr)
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(c, fmt.Sprintf("修改了系统组Group [%s]", updateData["gname"]), req, old, updateData)
	response.OnSuccess(c)
}

func RunSysGroupsMainGetdata(c *gin.Context) {
	db := utils.GetDB(c)
	var conds []any
	var exprs []string
	gname := strings.TrimSpace(c.Query("gname"))
	if len(gname) > 0 {
		exprs = append(exprs, "gname = ?")
		conds = append(conds, gname)
	}
	tx := db
	if len(exprs) > 0 {
		tx = tx.Where(strings.Join(exprs, " AND "), conds...)
	}

	page := cast.ToInt(c.Query("page"))
	if page < 1 {
		page = 1
	}
	limit := cast.ToInt(c.Query("limit"))
	if limit < 1 || limit > 200 {
		limit = 10
	}
	start := (page - 1) * limit

	var count int64
	if err := tx.Model(&model.SysGroup{}).Count(&count).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCountErr, err)
		response.OnFailure(c, msg.SqlCountErr)
		return
	}

	var results []map[string]any
	if err := tx.Model(&model.SysGroup{}).
		Offset(start).
		Limit(limit).
		Order("ordnum ASC").
		Find(&results).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindListErr, err)
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	for i := range results {
		results[i]["created_at"] = utils.Pack.DateTime.RFC3339ToYmdHis(results[i]["created_at"])
		results[i]["updated_at"] = utils.Pack.DateTime.RFC3339ToYmdHis(results[i]["updated_at"])
	}

	response.OkTableList(results, count, c)
}

// 删除系统组
func RunSysGroupMainDel(c *gin.Context) {
	db := utils.GetDB(c)
	id := cast.ToInt(c.Query("id"))
	if id <= 0 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}
	var data model.SysGroup
	if err := db.Where("id = ?", id).First(&data).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlFindErr)
		return
	}
	err := db.Delete(&data).Error
	if err != nil {
		utils.LogSqlErr(c, msg.SqlDeleteErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlDeleteErr)
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("删除系统组 [%s]", data.Gname),
		map[string]any{"id": id}, // Param,
		data,
		nil,
	)
	response.OnSuccess(c)
}
