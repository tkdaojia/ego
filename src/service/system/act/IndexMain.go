package service

import (
	"ego/src/boot/global"
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

func RunSysIndexMainGet(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "welcome":
		RunWelcome(c)
	case "wait":
		RunWait(c)
	case "other":
		RunIndexOther(c)
	case "editpwd":
		RunEditPwd(c)
	case "choose":
		RunChoose(c)
	default:
		RunIndexDefault(c)

	}
}

func RunSysIndexMainPost(c *gin.Context) {
	do := c.Query("do")
	switch do {
	case "editpwd":
		RunEditPwdPost(c)
	}
}

func RunEditPwdPost(c *gin.Context) {
	db := utils.GetDB(c)
	uid := utils.LibGetUid(c)

	// 1. 只查询需要的密码字段，提高查询效率
	var info model.SysAccount
	if err := db.Model(&model.SysAccount{}).Select("id", "password").Where("id = ?", uid).First(&info).Error; err != nil {
		response.BackMsgErr("账号不存在或已停用", c)
		return
	}

	oldPwd := c.PostForm("oldPassword")
	password := c.PostForm("password")
	rePwd := c.PostForm("repassword")

	// 2. 基础输入校验
	if len(password) < 6 { // 顺便限制最小密码长度，提升安全性
		response.BackMsgErr("新密码长度不能少于 6 位", c)
		return
	}
	if password != rePwd {
		response.BackMsgErr("两次输入的新密码不一致", c)
		return
	}
	if oldPwd == password {
		response.BackMsgErr("新密码不能与原密码相同", c)
		return
	}

	// 3. 验证原密码 (Bcrypt 验证)
	if err := bcrypt.CompareHashAndPassword([]byte(info.Password), []byte(oldPwd)); err != nil {
		response.BackMsgErr("原密码错误", c)
		return
	}

	// 4. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.BackMsgErr("密码加密失败，请稍后重试", c)
		return
	}

	// 5. 优雅更新（直接使用模型更新，避免硬编码表名）
	if err := db.Model(&info).Update("password", string(hashedPassword)).Error; err != nil {
		response.BackMsgErr("密码更新失败", c)
		return
	}

	response.OnSuccess(c)
}

func RunWait(c *gin.Context) {
	c.HTML(200, "index/wait.htm", gin.H{})
}

func RunEditPwd(c *gin.Context) {
	c.HTML(200, "index/editpwd.htm", gin.H{})
}

func RunIndexOther(c *gin.Context) {
	c.HTML(200, "index/other.htm", gin.H{})
}

func RunIndexDefault(c *gin.Context) {
	db := utils.GetDB(c)
	uid := utils.LibGetUid(c)

	var user model.SysAccount
	if err := db.Where("id = ?", uid).First(&user).Error; err != nil {
		response.OnFailure(c, msg.UserDataLost)
		return
	}

	var roleArr []model.SysRole
	if err := db.Where("id IN ?", strings.Split(user.Role_id, ",")).Find(&roleArr).Error; err != nil {
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	var gropuIdArr = []string{}

	for _, role := range roleArr {
		if len(role.Sysgroup) > 0 {
			tmp := strings.Split(role.Sysgroup, ",")
			for _, rid := range tmp {
				if slices.Contains(gropuIdArr, rid) == false {
					gropuIdArr = append(gropuIdArr, rid)
				}
			}
		}
	}
	userGroup := gropuIdArr
	var sysgroup []model.SysGroup
	if err := db.Where("state = 1 AND id IN ?", userGroup).Order("ordnum ASC").Find(&sysgroup).Error; err != nil {
		response.OnFailure(c, msg.SqlFindListErr)
		return
	}

	t := utils.GetTimestamp()
	webname := global.C_CONFIG.System.Webname
	c.HTML(200, "index/first.htm", gin.H{
		"t":         t,
		"sysgroup":  sysgroup,
		"userGroup": userGroup,
		"webname":   webname,
	})
}

func RunWelcome(c *gin.Context) {
	c.HTML(200, "index/welcome.htm", gin.H{})
}

func RunChoose(c *gin.Context) {
	db := utils.GetDB(c)
	sysID := cast.ToInt(c.Query("system"))
	uid := utils.LibGetUid(c)

	if sysID < 1 {
		response.OnFailure(c, msg.IdInvalidErr)
		return
	}

	var sysgroup model.SysGroup
	if err := db.Where("id = ? AND state = 1", sysID).First(&sysgroup).Error; err != nil {
		response.OnFailure(c, "本系统未开启或不存在")
		return
	}

	roles := utils.LibGetUrole(c)
	var menucateIDs []string
	var menulistIDs []string

	if len(roles) > 0 {
		var sysRoles []model.SysRole
		db.Where("name IN ?", roles).Find(&sysRoles)

		// 使用 map 进行去重并过滤空字符串，防止脏数据干扰 SQL 索引
		cateMap := make(map[string]struct{})
		listMap := make(map[string]struct{})

		for _, role := range sysRoles {
			for _, id := range strings.Split(role.Menucate, ",") {
				if id = strings.TrimSpace(id); id != "" {
					cateMap[id] = struct{}{}
				}
			}
			for _, id := range strings.Split(role.Menulist, ",") {
				if id = strings.TrimSpace(id); id != "" {
					listMap[id] = struct{}{}
				}
			}
		}

		// 转回切片
		for id := range cateMap {
			menucateIDs = append(menucateIDs, id)
		}
		for id := range listMap {
			menulistIDs = append(menulistIDs, id)
		}
	}

	if len(menucateIDs) == 0 && len(menulistIDs) == 0 {
		response.OnFailure(c, "您暂无该系统的访问权限")
		return
	}

	var allMenus []map[string]interface{}
	allIDs := append(menucateIDs, menulistIDs...) // 合并一二级菜单 ID 统一查询

	db.Model(&model.SysMenu{}).
		Where("typeid = ? AND status = 1 AND id IN ?", sysID, allIDs).
		Order("ordnum ASC, id ASC").
		Find(&allMenus)

	var menuTree []map[string]interface{}
	timestamp := utils.GetTimestamp()

	for _, m := range allMenus {
		pid := cast.ToInt(m["pid"])
		idStr := fmt.Sprintf("%v", m["id"])

		// 判断是否属于用户拥有的“一级菜单权限集”且自身为顶级菜单
		if pid == 0 && slices.Contains(menucateIDs, idStr) {
			var subMenus []map[string]interface{}

			// 寻找当前一级菜单下的子菜单
			for _, sub := range allMenus {
				subPid := cast.ToInt(sub["pid"])
				subIDStr := fmt.Sprintf("%v", sub["id"])

				if subPid == m["id"] && slices.Contains(menulistIDs, subIDStr) {
					sub["time"] = timestamp
					subMenus = append(subMenus, sub)
				}
			}
			m["sub"] = subMenus
			menuTree = append(menuTree, m)
		}
	}

	var user map[string]interface{}
	db.Model(&model.SysAccount{}).Where("id = ? AND status = 1", uid).Select("truename, status").First(&user)

	//  统一渲染输出
	c.HTML(200, sysgroup.Indextpl, gin.H{
		"menu":       menuTree,
		"user":       user,
		"webname":    global.C_CONFIG.System.Webname,
		"cookiename": global.C_CONFIG.System.Cookiename,
	})
}
