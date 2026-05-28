package service

import (
	"ego/src/cls"
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

func RunSysMenuList(c *gin.Context) {
	c.HTML(200, "sysmenu/list.htm", nil)
}

func RunSysMenuAdd(c *gin.Context) {
	db := utils.GetDB(c)
	id := cast.ToInt(c.Query("id"))
	var info model.SysMenu

	pid := 0
	mid := 0
	if id > 0 {
		if err := db.Where("id = ?", id).First(&info).Error; err != nil {
			utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
			response.OnFailure(c, msg.SqlFindErr)
			return
		}
		pid = info.Pid
		mid = info.Typeid
	}
	//读取大类
	var cate []model.SysMenu
	if err := db.Where("pid = 0").Find(&cate).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindListErr, err)
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	menuType := cls.MapSetIDSysGroup(c)
	c.HTML(200, "sysmenu/add.htm", gin.H{
		"info":     info,
		"cate":     cate,
		"pid":      pid,
		"menuType": menuType,
		"mid":      mid,
	})
}

// 菜单列表接口
func RunSysMenuGetdata(c *gin.Context) {
	db := utils.GetDB(c)

	var allMenus []model.SysMenu
	if err := db.Order("ordnum ASC, id ASC").Find(&allMenus).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindListErr, err)
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	count := int64(len(allMenus))

	// 获取类型映射 map
	mtype := cls.MapSetIDSysGroup(c)

	var firstLevelList []map[string]any          // 存放最终的一级菜单
	subMenuMap := make(map[int][]map[string]any) // 存放二级菜单，Key 是 Pid

	// 第一次遍历：转化为 map 格式，并将二级菜单按 Pid 归类
	for _, item := range allMenus {
		// 获取当前菜单的类型名称
		tid := item.Typeid
		mtypeName := mtype[tid]

		m := map[string]any{
			"id":         item.ID,
			"pid":        item.Pid,
			"mname":      item.Mname,
			"name":       item.Mname,
			"mlink":      item.Mlink,
			"icon":       item.Icon,
			"status":     item.Status,
			"typeid":     item.Typeid,
			"mtype_name": mtypeName,
		}

		if item.Pid == 0 {
			m["open"] = true // 只有一级菜单需要 open = true
			firstLevelList = append(firstLevelList, m)
		} else {
			m["parentId"] = item.Pid // 二级菜单加上 parentId
			subMenuMap[item.Pid] = append(subMenuMap[item.Pid], m)
		}
	}

	//  第二次遍历：把分好类的二级菜单，塞进对应一级菜单的 "children" 属性中
	for i := range firstLevelList {
		parentId := cast.ToInt(firstLevelList[i]["id"])

		// 从 Map 里直接提取子菜单（纯内存操作，效率极高）
		if children, exists := subMenuMap[parentId]; exists {
			firstLevelList[i]["children"] = children
		} else {
			// 如果没有子菜单，返回空切片 `[]` 避免前端解析报错
			firstLevelList[i]["children"] = []map[string]any{}
		}
	}

	response.OkTableList(firstLevelList, count, c)
}

// 删除菜单接口
func RunSysMenuDel(c *gin.Context) {
	db := utils.GetDB(c)

	id := cast.ToUint(c.Query("id"))
	if id == 0 {
		response.OnFailure(c, "无效的菜单ID")
		return
	}
	var menuData model.SysMenu

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&menuData).Error; err != nil {
			return err
		}
		// 如果是顶级菜单(Pid == 0)，先删除所有子菜单
		if menuData.Pid == 0 {
			if err := tx.Delete(&model.SysMenu{}, "pid = ?", menuData.ID).Error; err != nil {
				return err
			}
		}
		// 删除当前菜单
		if err := tx.Delete(&menuData).Error; err != nil {
			return err
		}
		return nil
	})

	// 统一错误处理
	if err != nil {
		// 如果是记录不存在的错误，返回特定提示
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.OnFailure(c, "菜单不存在或已被删除")
			return
		}
		// 记录系统级日志
		utils.LogSqlErr(c, "Delete_SysMenu", err, zap.Uint("id", id))
		response.OnFailure(c, "删除菜单失败，请稍后再试")
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("删除了菜单 [%s]", menuData.Mname),
		map[string]any{"id": id}, // Param: 可以直接传个 map 或者 query 结构体
		menuData,                 // DataOld: 刚才查出来的旧数据
		nil,                      // DataNew: 删除操作没有新数据，传 nil
	)

	response.OnSuccess(c)
}

// 新增接口
func RunSysMenuCrate(c *gin.Context) {
	var req struct {
		Pid    string `json:"pid"`
		Mname  string `json:"mname" binding:"required"`
		Mlink  string `json:"mlink" binding:"required"`
		Icon   string `json:"icon"`
		Ordnum string `json:"ordnum"`
		Status string `json:"status"`
		Typeid string `json:"typeid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	db := utils.GetDB(c)
	item := model.SysMenu{
		Mname:  strings.TrimSpace(req.Mname),
		Mlink:  strings.TrimSpace(req.Mlink),
		Icon:   strings.TrimSpace(req.Icon),
		Typeid: cast.ToInt(req.Typeid),
		Status: cast.ToInt(req.Status),
		Ordnum: cast.ToInt(req.Ordnum),
		Pid:    cast.ToInt(req.Pid),
	}

	if err := db.Create(&item).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCreateErr, err)
		response.OnFailure(c, msg.SqlCreateErr)
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(c, fmt.Sprintf("新增菜单 [%s]", item.Mname), req, nil, item)
	response.OnSuccess(c)
}

// 修改接口
func RunSysMenuUpdate(c *gin.Context) {
	var req struct {
		Id     int    `json:"id" binding:"required,gt=0"`
		Pid    string `json:"pid"`
		Mname  string `json:"mname" binding:"required"`
		Mlink  string `json:"mlink" binding:"required"`
		Icon   string `json:"icon"`
		Ordnum string `json:"ordnum"`
		Status string `json:"status"`
		Typeid string `json:"typeid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	db := utils.GetDB(c)
	id := cast.ToInt(req.Id)

	var old model.SysMenu
	if err := db.Where("id = ?", id).First(&old).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlFindErr)
		return
	}

	updateData := map[string]any{
		"mname":  strings.TrimSpace(req.Mname),
		"mlink":  strings.TrimSpace(req.Mlink),
		"icon":   strings.TrimSpace(req.Icon),
		"ordnum": cast.ToInt(req.Ordnum),
		"status": cast.ToInt(req.Status),
		"typeid": cast.ToInt(req.Typeid),
		"pid":    cast.ToInt(req.Pid),
	}
	if err := db.Model(&old).Updates(updateData).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlUpdateErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlUpdateErr)
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(c, fmt.Sprintf("修改了菜单 [%s]", updateData["mname"]), req, old, updateData)
	response.OnSuccess(c)
}
