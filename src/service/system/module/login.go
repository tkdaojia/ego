package service

import (
	"ego/src/boot/global"
	"ego/src/model/basic"
	"ego/src/model/response"
	"ego/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func RunLogin(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "htm":
		RunLoginHtm(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}

func RunLoginHtm(c *gin.Context) {
	webname := global.C_CONFIG.System.Webname
	c.HTML(200, "open/login.htm", gin.H{
		"webname": webname,
	})
}

func RunLoginPost(c *gin.Context) {
	type LoginParam struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var param LoginParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.OnFailure(c, "用户名或密码不能为空")
		return
	}

	db := utils.GetDB(c)
	username := strings.TrimSpace(param.Username)
	password := param.Password

	if username == "" || password == "" {
		response.OnFailure(c, "用户名或密码不能为空")
		return
	}

	var info model.SysAccount
	err := db.Model(&model.SysAccount{}).
		Select("id", "account", "truename", "password", "status", "login_count", "role_id").
		Where("account = ?", username).
		First(&info).Error

	if err != nil {
		utils.Pack.AuditLog.SaveAuditLog(c, fmt.Sprintf("异常用户登陆 [%s]", param.Username), param, nil, nil)
		response.OnFailure(c, "用户名或密码错误")
		return
	}

	// 3. 校验账号是否被停用
	if info.Status != 1 {
		utils.Pack.AuditLog.SaveAuditLog(c, fmt.Sprintf("停用账号登陆拦截 [%s]", info.Account), param, nil, nil)
		response.OnFailure(c, "该账号已被停用，请联系管理员")
		return
	}

	isPasswordCorrect := false

	if err := bcrypt.CompareHashAndPassword([]byte(info.Password), []byte(password)); err == nil {
		isPasswordCorrect = true
	}

	if !isPasswordCorrect {
		utils.Pack.AuditLog.SaveAuditLog(c, fmt.Sprintf("用户登陆密码错误 [%s]", param.Username), param, nil, nil)
		response.OnFailure(c, "用户名或密码错误")
		return
	}

	utils.Pack.AuditLog.SaveAuditLog(c, fmt.Sprintf("账号登陆 [%s]", info.Account), param, nil, nil)

	updates := map[string]interface{}{
		"lastlogin":   utils.GetTimestamp(),
		"login_count": info.Login_count + 1,
	}

	db.Model(&info).Updates(updates)

	var roles []string
	if info.Role_id != "" {
		roleIDs := strings.Split(info.Role_id, ",")
		var roleArr []model.SysRole
		// 批量捞出角色标识符
		db.Select("name").Where("id IN ?", roleIDs).Find(&roleArr)
		for _, v := range roleArr {
			if v.Name != "" {
				roles = append(roles, v.Name)
			}
		}
	}
	signedToken := utils.CreateJWT(cast.ToInt(info.ID), info.Account, info.Truename, roles)

	online := model.Online{
		Uid:    info.ID,
		Ip:     c.ClientIP(),
		Active: utils.GetTimestamp(),
		Status: 1,
	}
	db.Where("uid = ?", info.ID).Assign(online).FirstOrCreate(&online)

	c.SetCookie(
		global.C_CONFIG.System.Cookiename, // Cookie 键名
		signedToken,                       // Token 值
		3600*24,                           // 有效期 1 天（需与 utils.ExpireTime 保持一致）
		"/",                               // 全局路径有效
		"",                                // 留空代表当前域名
		false,                             // Secure: 如果是 HTTPS 环境，生产环境记得改为 true
		true,                              // HttpOnly: 👈 开启此项，前端 JS 将无法通过 document.cookie 窃取 Token
	)

	response.Result(200, nil, "登录成功", c)
}
