package initialize

import (
	"context" // 必须加这个
	"ego/src/boot/global"
	model "ego/src/model/basic"
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
	"time" // 可选
)

func Redis() {
	redisCfg := global.C_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic("Redis 连接失败:" + err.Error())
	}
	// 初始化本地缓存
	cacheSize := 10 * 1024 * 1024 // 10MB
	var a model.C_RedisCache
	a.REDIS = client
	a.CACHE = freecache.NewCache(cacheSize)
	global.C_REDIS = a
}
