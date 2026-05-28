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
	case "out":
		RunLogOut(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{"message": "action not found"})
	}
}

func RunLoginHtm(c *gin.Context) {
	webname := global.C_CONFIG.System.Webname
	cookiename := global.C_CONFIG.System.Cookiename
	c.HTML(200, "open/login.htm", gin.H{
		"webname":    webname,
		"cookiename": cookiename,
	})
}

func RunLogOut(c *gin.Context) {
	cookieName := "user_cookie"
	if global.C_CONFIG.System.Cookiename != "" {
		cookieName = global.C_CONFIG.System.Cookiename
	}
	host := c.Request.Host
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	c.SetCookie(cookieName, "", -1, "/", host, false, true)
	c.Redirect(http.StatusFound, "/open/?module=login&act=htm")
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

	response.Result(200, signedToken, "登录成功", c)
}
