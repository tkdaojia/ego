package service

import (
	model "ego/src/model/basic"
	"ego/src/model/msg"
	"ego/src/model/response"
	"ego/src/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

// 列表页面
func RunSysmoduleList(c *gin.Context) {
	c.HTML(200, "sysmodule/list.htm", nil)
}

// 列表接口
func RunSysmoduleGetdata(c *gin.Context) {
	db := utils.GetDB(c)

	var allModules []model.SysModule
	if err := db.Order("id ASC").Find(&allModules).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindListErr, err)
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	// 总数量直接取切片的长度
	count := int64(len(allModules))

	var firstLevelList []map[string]any            // 存放最终的一级菜单（带children）
	subModuleMap := make(map[int][]map[string]any) // 存放二级菜单，Key是Pid

	// 第一次遍历：把所有菜单转化为前端需要的 map 格式，并将二级菜单按 Pid 归类
	for _, item := range allModules {
		m := map[string]any{
			"id":      item.ID,
			"pid":     item.Pid,
			"mname":   item.Mname,
			"name":    item.Mname, // 对应原代码：sortRes[x]["name"] = sortRes[x]["mname"]
			"remarks": item.Remarks,
		}

		if item.Pid == 0 {
			// 先把一级菜单暂存，因为现在二级菜单还没归类完
			firstLevelList = append(firstLevelList, m)
		} else {
			// 二级菜单填上父级ID
			m["parentId"] = item.Pid
			// 按父级 Pid 归类到 Map 里
			subModuleMap[item.Pid] = append(subModuleMap[item.Pid], m)
		}
	}

	//  第二次遍历：把分好类的二级菜单，塞进对应一级菜单的 "children" 属性中
	for i := range firstLevelList {
		parentId := cast.ToInt(firstLevelList[i]["id"])

		// 从 Map 里直接把属于它的子菜单拿出来（纯内存操作，$O(1)$ 复杂度）
		if children, exists := subModuleMap[parentId]; exists {
			firstLevelList[i]["children"] = children
		} else {
			// 如果没有子菜单，主流规范是给一个空切片 `[]` 而不是 `nil`，方便前端解析
			firstLevelList[i]["children"] = []map[string]any{}
		}
	}

	response.OkTableList(firstLevelList, count, c)
}

func RunSysmoduleAdd(c *gin.Context) {
	db := utils.GetDB(c)
	id := cast.ToInt(c.Query("id"))

	var info model.SysModule

	var pid int = 0

	if id > 0 {
		if err := db.Where("id = ?", id).First(&info).Error; err != nil {
			utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
			response.FailWithMessage(msg.SqlFindErr, c)
			return
		}
		pid = info.Pid
	}

	var cate []model.SysModule
	if err := db.Where("pid = 0").Order("id ASC").Find(&cate).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err)
		response.FailWithMessage(msg.SqlFindListErr, c)
		return
	}

	c.HTML(200, "sysmodule/add.htm", gin.H{
		"info": info,
		"cate": cate,
		"pid":  pid,
	})
}

// 新增接口
func RunSysmoduleCrate(c *gin.Context) {
	var req struct {
		Pid     string `json:"pid"`
		Mname   string `json:"mname" binding:"required"` // 放心使用必填校验
		Remarks string `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	db := utils.GetDB(c)
	item := model.SysModule{
		Mname:   strings.TrimSpace(req.Mname),
		Remarks: strings.TrimSpace(req.Remarks),
		Pid:     cast.ToInt(req.Pid),
	}

	if err := db.Create(&item).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCreateErr, err)
		response.OnFailure(c, msg.SqlCreateErr)
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(c, fmt.Sprintf("新增了Module [%s]", item.Mname), req, nil, item)
	response.OnSuccess(c)
}

// 修改接口
func RunSysmoduleUpdate(c *gin.Context) {
	var req struct {
		Id      int    `json:"id" binding:"required,gt=0"`
		Pid     string `json:"pid"`
		Mname   string `json:"mname" binding:"required"`
		Remarks string `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	db := utils.GetDB(c)
	id := cast.ToInt(req.Id)

	var old model.SysModule
	if err := db.Where("id = ?", id).First(&old).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlFindErr)
		return
	}

	updateData := map[string]any{
		"mname":   strings.TrimSpace(req.Mname),
		"remarks": strings.TrimSpace(req.Remarks),
		"pid":     cast.ToInt(req.Pid),
	}
	if err := db.Model(&old).Updates(updateData).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlUpdateErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlUpdateErr)
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(c, fmt.Sprintf("修改了Module [%s]", updateData["mname"]), req, old, updateData)
	response.OnSuccess(c)
}

// 删除接口
func RunSysmoduleDel(c *gin.Context) {
	req := model.SystemReqDel
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}
	db := utils.GetDB(c)
	id := req.Id
	if id <= 0 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}
	var data model.SysModule
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
			return err
		}
		// 如果是顶级(Pid == 0)，先删除所有子件
		if data.Pid == 0 {
			if err := tx.Delete(&model.SysModule{}, "pid = ?", data.ID).Error; err != nil {
				return err
			}
		}
		// 删除当前
		if err := tx.Delete(&data).Error; err != nil {
			return err
		}
		return nil
	})

	// 统一错误处理
	if err != nil {
		// 如果是记录不存在的错误，返回特定提示
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.OnFailure(c, msg.SqlFindErr)
			return
		}
		// 记录系统级日志
		utils.LogSqlErr(c, msg.SqlDeleteErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlDeleteErr)
		return
	}

	//日志审计
	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("删除组件Module [%s]", data.Mname),
		map[string]any{"id": id}, // Param,
		data,
		nil,
	)

	response.OnSuccess(c)
}
