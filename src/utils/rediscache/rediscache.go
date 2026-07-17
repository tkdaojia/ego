package rediscache

import (
	"context"
	"ego/src/boot/global"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct{}

// sec 是分钟
func (redisCache *RedisCache) CacheSet(key string, value []byte, sec int) error {
	key = global.C_CONFIG.Redis.Pre + key
	return global.C_REDIS.CACHE.Set([]byte(key), value, sec*60)
}

func (redisCache *RedisCache) CacheGet(key string) ([]byte, error) {
	key = global.C_CONFIG.Redis.Pre + key
	return global.C_REDIS.CACHE.Get([]byte(key))
}

func (redisCache *RedisCache) RedisSet(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	key = global.C_CONFIG.Redis.Pre + key
	err = global.C_REDIS.REDIS.Set(ctx, key, value, expiration).Err()
	return err
}

func (redisCache *RedisCache) RedisGet(ctx context.Context, key string) (string, error) {
	key = global.C_CONFIG.Redis.Pre + key
	val, err := global.C_REDIS.REDIS.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return val, nil
}
