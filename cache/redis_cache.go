package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisCache struct {
	redisCli *redis.Client
}

func (r *redisCache) Set(key string, value interface{}, ttl int) error {
	return r.redisCli.Set(context.Background(), key, value, time.Second*time.Duration(ttl)).Err()
}

func (r *redisCache) Get(key string) string {
	return r.redisCli.Get(context.Background(), key).Val()
}

func (r *redisCache) Del(key string) error {
	return r.redisCli.Del(context.Background(), key).Err()
}

func NewRedisCache(addr, password string, db ...int) Cache {
	selectDb := 0
	if len(db) > 0 {
		selectDb = db[0]
	}
	return &redisCache{redisCli: redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       selectDb,
	})}
}
