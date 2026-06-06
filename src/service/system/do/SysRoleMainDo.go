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
	"slices"
	"strings"
	"time"
)

func RunSysRoleList(c *gin.Context) {
	c.HTML(200, "sysrole/list.htm", nil)
}

func RunSysRoleAddMenu(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))
	if id <= 0 {
		response.ErrHtml(c, msg.IdInvalidErr)
		return
	}
	c.HTML(200, "sysrole/addmenu.htm", gin.H{
		"id": id,
	})
}

func RunSysRoleAddModule(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))
	if id <= 0 {
		response.ErrHtml(c, msg.IdInvalidErr)
		return
	}
	c.HTML(200, "sysrole/addmodule.htm", gin.H{
		"id": id,
	})
}

func RunSysRoleAdd(c *gin.Context) {
	db := utils.GetDB(c)
	id := cast.ToInt(c.Query("id"))

	var info model.SysRole

	if id > 0 {
		if err := db.Where("id = ?", id).First(&info).Error; err != nil {
			response.ErrHtml(c, msg.SqlFindErr)
			utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
			return
		}
	}

	c.HTML(200, "sysrole/add.htm", gin.H{
		"info": info,
	})
}

// 删除菜单接口
func RunSysRoleDel(c *gin.Context) {
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

	var old model.SysRole
	if err := db.Where("id = ?", id).First(&old).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlFindErr)
		return
	}
	if err := db.Delete(&old).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlDeleteErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlDeleteErr)
		return
	}
	//日志审计
	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("删除角色菜单 [%s]", old.Name),
		map[string]any{"id": id}, // Param,
		old,
		nil,
	)
	response.OnSuccess(c)
}

// 编辑菜单组权限
func RunSysRoleUpdateMenu(c *gin.Context) {
	db := utils.GetDB(c)
	var req struct {
		Id  string `json:"id" binding:"required"`
		Mid string `json:"mid"` // 接收到的菜单ID字符串，如 "1,2,3,4"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	id := cast.ToInt(req.Id)
	if id <= 0 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}

	var old model.SysRole
	if err := db.Where("id = ?", id).First(&old).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlFindErr)
		return
	}

	midStr := strings.TrimSpace(req.Mid)

	cidMap := make(map[string]struct{})
	sidMap := make(map[string]struct{})

	if midStr != "" {
		mArr := strings.Split(midStr, ",")

		var menus []model.SysMenu
		if err := db.Select("id, pid").Where("id IN ?", mArr).Find(&menus).Error; err != nil {
			utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
			response.OnFailure(c, msg.SqlFindListErr)
			return
		}

		for _, menu := range menus {
			sid := cast.ToString(menu.ID)
			if menu.Pid == 0 {
				cidMap[sid] = struct{}{} // 顶级菜单存入父级
			} else {
				pid := cast.ToString(menu.Pid)
				cidMap[pid] = struct{}{} // 子菜单的父ID也必须存入父级，保证菜单链条完整
				sidMap[sid] = struct{}{} // 属于子菜单
			}
		}
	}

	var cidArr, sidArr []string
	for k := range cidMap {
		cidArr = append(cidArr, k)
	}
	for k := range sidMap {
		sidArr = append(sidArr, k)
	}

	updateData := map[string]any{
		"menucate": strings.Join(cidArr, ","),
		"menulist": strings.Join(sidArr, ","),
	}

	if err := db.Model(&old).Updates(updateData).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlUpdateErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlUpdateErr)
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("修改角色菜单权限 [%s]", old.Name),
		req,
		old,
		updateData,
	)

	response.OnSuccess(c)
}

// 编辑module权限
func RunSysRoleUpdateModule(c *gin.Context) {
	db := utils.GetDB(c)
	var req struct {
		Id  string `json:"id"`
		Mid string `json:"mid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	id := cast.ToInt(req.Id)
	if id <= 0 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}

	var old model.SysRole
	if err := db.Where("id = ?", id).First(&old).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlFindErr)
		return
	}

	midStr := strings.TrimSpace(req.Mid)
	var roleListArr []string

	if midStr != "" {
		mArr := strings.Split(midStr, ",")

		var subModules []model.SysModule
		if err := db.Select("id, pid, mname").Where("id IN ?", mArr).Find(&subModules).Error; err != nil {
			utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
			response.OnFailure(c, msg.SqlFindListErr)
			return
		}

		parentIdMap := make(map[int]struct{})
		for _, sub := range subModules {
			if sub.Pid > 0 {
				parentIdMap[sub.Pid] = struct{}{}
			}
		}

		var parentIds []int
		for pId := range parentIdMap {
			parentIds = append(parentIds, pId)
		}

		parentNameMap := make(map[int]string)
		if len(parentIds) > 0 {
			var parentModules []model.SysModule
			if err := db.Select("id, mname").Where("id IN ?", parentIds).Find(&parentModules).Error; err == nil {
				for _, p := range parentModules {
					parentNameMap[p.ID] = p.Mname
				}
			}
		}

		// 纯内存拼装字符串
		for _, sub := range subModules {
			if sub.Pid == 0 {
				//  如果自己就是顶级大类，直接存入自己的名字
				roleListArr = append(roleListArr, sub.Mname)
			} else {
				// 如果是子菜单，则拼接成 "父模块:子模块" 格式
				pName := parentNameMap[sub.Pid]
				if pName != "" {
					roleListArr = append(roleListArr, pName+":"+sub.Mname)
				}
			}
		}
	}

	updateData := map[string]any{
		"rolelist":   strings.Join(roleListArr, ","),
		"modulelist": midStr,
	}

	if err := db.Model(&old).Updates(updateData).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlUpdateErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlUpdateErr)
		return
	}

	_ = utils.Pack.RedisCache.RedisSet(c, "app:role:"+old.Name, updateData["rolelist"], 2*time.Hour)

	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("修改角色模块 [%s]", old.Name),
		req,
		old,
		updateData,
	)
	response.OnSuccess(c)
}

// 列表接口
func RunSysRoleGetdata(c *gin.Context) {
	db := utils.GetDB(c) // WithContext 保证了当前实例的请求级绝对安全
	var count int64
	if err := db.Model(&model.SysRole{}).Count(&count).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCountErr, err)
		response.OnFailure(c, msg.SqlCountErr)
		return
	}
	var results []map[string]any
	if err := db.Model(&model.SysRole{}).Find(&results).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindListErr, err)
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}
	groupIdMap := make(map[int]struct{})
	for _, v := range results {
		sgStr := cast.ToString(v["sysgroup"])
		if len(sgStr) > 0 {
			for _, idStr := range strings.Split(sgStr, ",") {
				if id := cast.ToInt(strings.TrimSpace(idStr)); id > 0 {
					groupIdMap[id] = struct{}{}
				}
			}
		}
	}

	var groupIds []int
	for id := range groupIdMap {
		groupIds = append(groupIds, id)
	}

	groupNameMap := make(map[int]string)
	if len(groupIds) > 0 {
		var sysGroups []model.SysGroup
		// 从底层干净的 db 发起批量查询
		if err := db.Select("id, gname").Where("id IN ?", groupIds).Find(&sysGroups).Error; err == nil {
			for _, g := range sysGroups {
				groupNameMap[cast.ToInt(g.ID)] = g.Gname
			}
		} else {
			utils.LogSqlErr(c, "Batch_Find_SysGroup_Err", err)
		}
	}

	for k, v := range results {
		sgStr := cast.ToString(v["sysgroup"])
		results[k]["sysgroup"] = "" // 默认置空，防止前端看到历史脏数据

		if len(sgStr) > 0 {
			var names []string
			for _, idStr := range strings.Split(sgStr, ",") {
				id := cast.ToInt(strings.TrimSpace(idStr))
				if name, exists := groupNameMap[id]; exists {
					names = append(names, name)
				}
			}
			// 用逗号重新拼接成可读的文本返回给前端表格
			results[k]["sysgroup"] = strings.Join(names, ",")
		}
	}

	if results == nil {
		results = []map[string]any{}
	}

	response.OkTableList(results, count, c)
}

// 新增接口
func RunSysRoleCreate(c *gin.Context) {
	db := utils.GetDB(c)

	var req struct {
		Name     string `json:"name" binding:"required"` // 借助标签确保不为空
		Rname    string `json:"rname" binding:"required"`
		Sysgroup string `json:"sysgroup"`
		Remarks  string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	name := strings.TrimSpace(req.Name)

	// 校验唯一性
	var count int64
	if err := db.Model(&model.SysRole{}).Where("name = ?", name).Count(&count).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCountErr, err)
		response.OnFailure(c, msg.SqlCountErr)
		return
	}
	if count > 0 {
		response.OnFailure(c, msg.NameExistsErr+name)
		return
	}

	// 构造新模型
	item := model.SysRole{
		Name:     name,
		Rname:    strings.TrimSpace(req.Rname),
		Sysgroup: strings.TrimSpace(req.Sysgroup),
		Remarks:  strings.TrimSpace(req.Remarks),
	}

	if err := db.Create(&item).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCreateErr, err)
		response.OnFailure(c, msg.SqlCreateErr)
		return
	}

	// 记录审计日志
	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("新增权限角色 [%s]", item.Rname),
		req,
		nil,
		item,
	)

	response.OnSuccess(c)
}

// 修改接口
func RunSysRoleUpdate(c *gin.Context) {
	var req struct {
		Id       int    `json:"id" binding:"required,gt=0"`
		Name     string `json:"name" binding:"required"`
		Rname    string `json:"rname" binding:"required"`
		Sysgroup string `json:"sysgroup"`
		Remarks  string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	db := utils.GetDB(c)
	id := cast.ToInt(req.Id)
	if id <= 0 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}

	var old model.SysRole
	if err := db.Where("id = ?", id).First(&old).Error; err != nil {
		response.OnFailure(c, msg.SqlFindErr)
		utils.LogSqlErr(c, msg.SqlFindErr, err, zap.Int("id", id))
		return
	}

	name := strings.TrimSpace(req.Name)
	var count int64
	err := db.Model(&model.SysRole{}).
		Where("name = ? AND id != ?", name, id).
		Count(&count).Error
	if err != nil {
		utils.LogSqlErr(c, msg.SqlCountErr, err)
		response.OnFailure(c, msg.SqlCountErr)
		return
	}
	if count > 0 {
		response.OnFailure(c, msg.NameExistsErr+name)
		return
	}

	updateData := map[string]any{
		"name":     name,
		"rname":    strings.TrimSpace(req.Rname),
		"sysgroup": strings.TrimSpace(req.Sysgroup),
		"remarks":  strings.TrimSpace(req.Remarks),
	}

	err = db.Model(&old).Updates(updateData).Error

	if err != nil {
		utils.LogSqlErr(c, msg.SqlUpdateErr, err, zap.Int("id", id))
		response.OnFailure(c, msg.SqlUpdateErr)
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("修改权限角色 [%s]", updateData["rname"]),
		req,
		old,
		updateData,
	)

	response.OnSuccess(c)
}

func RunSysRoleMenuTree(c *gin.Context) {
	db := utils.GetDB(c)
	id := cast.ToInt(c.Query("id"))
	if id <= 0 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}
	var role model.SysRole
	if err := db.Where("id = ?", id).First(&role).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlUpdateErr, err, zap.Int("id", id))
		return
	}

	catearr := strings.Split(role.Menucate, ",")
	sortarr := strings.Split(role.Menulist, ",")

	var cate []map[string]interface{}
	db.Model(&model.SysMenu{}).Select("id,mname").Where("pid = 0").Order("ordnum ASC,id ASC").Find(&cate)

	for key, value := range cate {
		var sort []map[string]interface{}
		db.Model(&model.SysMenu{}).Select("id,mname").Where("pid = ?", value["id"]).Order("ordnum ASC,id ASC").Find(&sort)
		cate[key]["name"] = value["mname"]
		cate[key]["value"] = value["id"]

		if slices.Contains(catearr, cast.ToString(value["id"])) == true {
			cate[key]["selected"] = true
		}
		for k, v := range sort {
			sort[k]["name"] = v["mname"]
			sort[k]["value"] = v["id"]

			if slices.Contains(sortarr, cast.ToString(v["id"])) == true {
				sort[k]["selected"] = true
			}
		}
		cate[key]["children"] = sort
	}
	response.OkTableList(cate, 0, c)
}

func RunSysRoleModuleTree(c *gin.Context) {
	db := utils.GetDB(c)
	id := c.Query("id")
	var role model.SysRole
	db.Where("id = ?", id).Find(&role)

	arr := strings.Split(role.Modulelist, ",")

	var cate []map[string]interface{}
	db.Model(&model.SysModule{}).Where("pid = 0").Order("id ASC").Find(&cate)

	for key, value := range cate {
		var sort []map[string]interface{}
		db.Model(&model.SysModule{}).Where("pid = ?", value["id"]).Order("id ASC").Find(&sort)
		cate[key]["name"] = value["mname"]
		cate[key]["value"] = value["id"]

		for k, v := range sort {
			tmp_mname := cast.ToString(v["mname"])
			tmp_remarks := cast.ToString(v["remarks"])
			sort[k]["name"] = v["mname"]
			sort[k]["value"] = v["id"]
			if tmp_remarks != "" {
				sort[k]["name"] = tmp_mname + "(" + cast.ToString(v["remarks"]) + ")"
			}
			if slices.Contains(arr, cast.ToString(v["id"])) == true {
				sort[k]["selected"] = true
			}
		}
		cate[key]["children"] = sort
	}
	response.OkTableList(cate, 0, c)
}
