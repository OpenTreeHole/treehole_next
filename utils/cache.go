package utils

import (
	"context"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/goccy/go-json"
	gocache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"time"
	"treehole_next/config"
)

var Cache *cache.Cache[[]byte]

func InitCache() {
	if config.Config.RedisURL != "" {
		redisStore := redis_store.NewRedis(redis.NewClient(&redis.Options{
			Addr: config.Config.RedisURL,
		}))
		Cache = cache.New[[]byte](redisStore)
	} else {
		gocacheStore := gocache_store.NewGoCache(gocache.New(5*time.Minute, 10*time.Minute))
		Cache = cache.New[[]byte](gocacheStore)
	}
}

const maxDuration time.Duration = 1<<63 - 1

func SetCache(key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if expiration == 0 {
		expiration = maxDuration
	}
	return Cache.Set(context.Background(), key, data, store.WithExpiration(expiration))
}

func GetCache(key string, value any) bool {
	data, err := Cache.Get(context.Background(), key)
	if err != nil {
		return false
	}
	err = json.Unmarshal(data, value)
	return err == nil
}

func DeleteCache(key string) error {
	err := Cache.Delete(context.Background(), key)
	if err == nil {
		return nil
	}
	if err.Error() == "Entry not found" {
		return nil
	}
	return err
}
