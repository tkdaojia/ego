package middleware

import (
	"context"
	"ego/src/boot/global"
	"ego/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

func VerifyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Cookie(global.C_CONFIG.System.Cookiename)
		if user == "" {
			gotoLogin(c)
			return
		}
		sz := utils.DecodeJWT(user)
		if sz == nil {
			gotoLogin(c)
			return
		}

		userRoles := sz.Roles
		if len(userRoles) == 0 {
			c.JSON(401, gin.H{"msg": "未登录或无角色信息"})
			c.Abort()
			return
		}

		module := strings.TrimSpace(c.Query("module"))
		act := strings.TrimSpace(c.Query("act"))

		currentPermission := fmt.Sprintf("%s:%s", module, act)

		if len(module) == 0 {
			gotoLogin(c)
			return
		}

		// 3. 遍历用户的所有角色，只要有一个角色满足权限，就放行
		hasPermission := false
		for _, role := range userRoles {
			if CheckRoleHasPermission(c.Request.Context(), role, currentPermission) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(403, gin.H{"msg": "对不起，您没有操作权限"})
			c.Abort()
			return
		}
		c.Set("User_Id", sz.UserID)
		c.Set("User_Account", sz.Account)
		c.Set("User_Name", sz.Username)
		c.Set("User_Role", sz.Roles)
		c.Set("request_start_time", time.Now())
		c.Next()
	}
}

func VerifyAuthApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func gotoLogin(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/open/?module=login&act=htm&t="+cast.ToString(utils.GetTimestamp()))
}

// 判断某个角色是否拥有指定权限
func CheckRoleHasPermission(ctx context.Context, role string, currentPermission string) bool {
	localKey := "app:role:" + role
	var permMap map[string]bool

	// ----------------------------------------------------
	// 步骤一：尝试从 Go 本地内存（一级缓存）获取
	// ----------------------------------------------------
	if res, err := utils.Pack.RedisCache.CacheGet(localKey); err == nil {
		permList := strings.Split(cast.ToString(res), ",")
		permMap = make(map[string]bool, len(permList))
		for _, p := range permList {
			permMap[p] = true
		}
		fmt.Println("✅ 直接返回啦", permMap[currentPermission])
		return permMap[currentPermission]
	}

	// ----------------------------------------------------
	// 步骤二：本地缓存未命中，从 Redis（二级缓存）获取
	// ----------------------------------------------------
	redisKey := localKey
	// 从 Redis Hash 中获取该角色的权限 JSON 字符串
	val, err := utils.Pack.RedisCache.RedisGet(ctx, redisKey)
	// 如果 Redis 中压根没有这个角色的配置（可能角色被删了，或者没有任何权限）
	if err == redis.Nil {
		return false
	} else if err != nil {
		// 生产环境这里建议记录日志：Redis 异常
		return false
	}

	// ----------------------------------------------------
	// 步骤三：解析数据并同步到本地缓存
	// ----------------------------------------------------
	var permList []string
	permList = strings.Split(val, ",")
	// 将本地结构转换为 Map，方便 $O(1)$ 高效查询
	permMap = make(map[string]bool, len(permList))
	for _, p := range permList {
		permMap[p] = true
	}

	// 序列化后存入本地缓存，设置过期时间为 5 分钟
	if err = utils.Pack.RedisCache.CacheSet(localKey, []byte(val), 5*60); err != nil {
		global.C_LOG.Error("中间件写入Redis错误", zap.String("error", err.Error()))
	}

	// 返回比对结果
	return permMap[currentPermission]
}
