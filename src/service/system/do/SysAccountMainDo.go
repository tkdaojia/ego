package service

import (
	"ego/src/cls"
	model "ego/src/model/basic"
	"ego/src/model/msg"
	"ego/src/model/response"
	"ego/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"slices"
	"strings"
)

func RunSysAccountList(c *gin.Context) {
	c.HTML(200, "sysaccount/list.htm", "")
}

// 返回table列表数据
func RunSysAccountGetdata(c *gin.Context) {
	db := utils.GetDB(c)

	tx := db.Model(&model.SysAccount{})

	// 接收筛选参数
	account := strings.TrimSpace(c.Query("account"))
	truename := strings.TrimSpace(c.Query("truename"))

	if account != "" {
		tx = tx.Where("account = ?", account)
	}
	if truename != "" {
		tx = tx.Where("truename LIKE ?", "%"+truename+"%")
	}

	var count int64
	if err := tx.Count(&count).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCountErr, err)
		response.OnFailure(c, msg.SqlCountErr)
		return
	}

	page := cast.ToInt(c.Query("page"))
	limit := cast.ToInt(c.Query("limit"))
	if limit < 1 {
		limit = 10
	}
	start := 0
	if page > 1 {
		start = (page - 1) * limit
	}

	var accounts []model.SysAccount
	if err := tx.Offset(start).Limit(limit).Find(&accounts).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlFindListErr, err)
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	// 系统组基础字典数据
	sysgroupMap := cls.MapSetIDSysGroup(c)

	roleIDMap := make(map[string]bool)
	for _, u := range accounts {
		if u.Role_id != "" {
			for _, rid := range strings.Split(u.Role_id, ",") {
				if r := strings.TrimSpace(rid); r != "" {
					roleIDMap[r] = true
				}
			}
		}
	}

	// 收集成切片一次性查库
	var allRoleIDs []string
	for rid := range roleIDMap {
		allRoleIDs = append(allRoleIDs, rid)
	}

	// 一次性查出全页所需的全部角色
	var roles []model.SysRole
	roleCache := make(map[string]model.SysRole)
	if len(allRoleIDs) > 0 {
		db.Model(&model.SysRole{}).Where("id IN ?", allRoleIDs).Find(&roles)
		// 转换成内存 Map，方便 $O(1)$ 复杂度读取
		for _, r := range roles {
			roleCache[cast.ToString(r.ID)] = r
		}
	}

	var results []map[string]interface{}
	for _, u := range accounts {
		item := map[string]interface{}{
			"id":       u.ID,
			"account":  u.Account,
			"truename": u.Truename,
			"role_id":  u.Role_id,
			"sysgroup": u.Sysgroup,
			"status":   u.Status,
		}

		var roleNames []string
		groupIDMap := make(map[string]bool) // 用于系统组去重

		if u.Role_id != "" {
			for _, rid := range strings.Split(u.Role_id, ",") {
				if role, ext := roleCache[rid]; ext {
					roleNames = append(roleNames, role.Rname)
					// 提取角色里的系统组
					if role.Sysgroup != "" {
						for _, sgid := range strings.Split(role.Sysgroup, ",") {
							groupIDMap[sgid] = true
						}
					}
				}
			}
		}

		if len(roleNames) > 0 {
			item["role_name"] = strings.Join(roleNames, ",")
		} else {
			item["role_name"] = "无"
		}
		var sysGroupNames []string
		for sgid := range groupIDMap {
			vid := cast.ToInt(sgid)
			if gName, ok := sysgroupMap[vid]; ok {
				sysGroupNames = append(sysGroupNames, cast.ToString(gName))
			}
		}
		if len(sysGroupNames) > 0 {
			item["sys_group"] = strings.Join(sysGroupNames, ",")
		} else {
			item["sys_group"] = "无"
		}

		results = append(results, item)
	}

	response.OkTableList(results, count, c)
}

func RunSysAccountAdd(c *gin.Context) {
	db := utils.GetDB(c)
	id := cast.ToUint(c.Query("id"))
	info := make(map[string]interface{})

	var user_role_arr []string
	if id > 0 {
		if err := db.Model(&model.SysAccount{}).Where("id = ?", id).Find(&info).Error; err != nil {
			response.OnFailure(c, msg.SqlFindErr)
			return
		}
		user_role_arr = strings.Split(cast.ToString(info["role_id"]), ",")
	}

	//读取权限组
	var role []map[string]interface{}
	db.Model(&model.SysRole{}).Select("id,rname").Find(&role)
	for key, value := range role {
		tmp_id := cast.ToString(value["id"])
		if slices.Contains(user_role_arr, tmp_id) {
			role[key]["select"] = true
		}
	}

	c.HTML(200, "sysaccount/add.htm", gin.H{
		"info": info,
		"role": role,
		"id":   id,
	})
}

func RunSysAccountDel(c *gin.Context) {
	db := utils.GetDB(c)
	id := cast.ToUint(c.Query("id"))

	result := db.Delete(&model.SysAccount{}, id)
	if result.RowsAffected != 1 {
		response.BackMsgErr("删除异常", c)
		return
	}
	response.BackMsgOk(c)
}

func RunSysAccountCreate(c *gin.Context) {
	type SysAccountReq struct {
		Account  string `json:"account" binding:"required"`  // 账号必填
		Truename string `json:"truename" binding:"required"` // 姓名必填
		RoleID   string `json:"role_id"`                     // 权限组
		Password string `json:"password" binding:"required"` // 密码
		Sysgroup string `json:"sysgroup"`                    // 系统组
	}

	var req SysAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}
	db := utils.GetDB(c)
	// 密码加密处理
	password := strings.TrimSpace(req.Password)
	if len(password) == 0 {
		response.OnFailure(c, msg.UserPasswordEmpty)
		return
	}
	var have int64
	db.Model(&model.SysAccount{}).Where("account = ?", req.Account).Count(&have)
	if have > 0 {
		response.OnFailure(c, msg.NameExistsErr)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.OnFailure(c, msg.PasswordEncryptErr)
		return
	}

	newPassword := string(hashedPassword)
	// 初始化新用户结构体
	item := model.SysAccount{
		Account:  req.Account,
		Truename: req.Truename,
		Password: newPassword,
		Role_id:  req.RoleID,
		Status:   1,
		Sysgroup: req.Sysgroup,
	}

	if err := db.Create(&item).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlCreateErr, err)
		response.OnFailure(c, msg.SqlCreateErr)
	}

	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("新增了用户 [%s]", req.Account),
		req,
		nil,
		item,
	)
	response.OnSuccess(c)
}

func RunSysAccountUpdate(c *gin.Context) {
	type SysAccountReq struct {
		ID       uint   `json:"id" binding:"required"`       // 编辑时 ID 必填
		Account  string `json:"account" binding:"required"`  // 账号必填
		Truename string `json:"truename" binding:"required"` // 姓名必填
		RoleID   string `json:"role_id"`                     // 权限组
		Password string `json:"password"`                    // 编辑时密码选填（留空代表不修改）
		Sysgroup string `json:"sysgroup"`                    // 系统组
	}

	var req SysAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OnFailure(c, msg.ReqParamErr)
		return
	}

	db := utils.GetDB(c)

	// 检查账号是否存在
	var oldItem model.SysAccount
	if err := db.Model(&model.SysAccount{}).Where("id = ?", req.ID).First(&oldItem).Error; err != nil {
		response.OnFailure(c, msg.SqlFindErr) // 或者你定义的未找到数据的错误码
		return
	}

	//  检查修改后的账号是否与其他账号冲突
	var have int64
	db.Model(&model.SysAccount{}).Where("account = ? AND id != ?", req.Account, req.ID).Count(&have)
	if have > 0 {
		response.OnFailure(c, msg.NameExistsErr)
		return
	}

	// 构建需要更新的字段（使用 map 可以精确控制更新，避免空字符串覆盖）
	data := map[string]interface{}{
		"account":  req.Account,
		"truename": req.Truename,
		"role_id":  req.RoleID,
		"sysgroup": req.Sysgroup,
	}

	// 只有当用户输入了新密码时，才进行加密并加入更新队列
	password := strings.TrimSpace(req.Password)
	if len(password) > 0 {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.OnFailure(c, msg.PasswordEncryptErr)
			return
		}
		newPassword := string(hashedPassword)
		data["password"] = newPassword
	}

	if err := db.Model(&model.SysAccount{}).Where("id = ?", req.ID).Updates(data).Error; err != nil {
		utils.LogSqlErr(c, msg.SqlUpdateErr, err)
		response.OnFailure(c, msg.SqlUpdateErr)
		return
	}

	// 5. 获取更新后的最新数据，用于记录审计日志
	var newItem model.SysAccount
	_ = db.Model(&model.SysAccount{}).Where("id = ?", req.ID).First(&newItem)

	// 6. 记录审计日志
	utils.Pack.AuditLog.SaveAuditLog(
		c,
		fmt.Sprintf("更新了用户 [%s]", req.Account),
		req,
		oldItem, // 传入旧数据以便对比
		newItem, // 传入新数据
	)

	response.OnSuccess(c)
}
