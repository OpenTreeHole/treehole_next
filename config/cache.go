package config

import (
	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/go-redis/redis/v8"
	"time"
)

var Cache *cache.Cache[[]byte]

func initCache() {
	if Config.RedisURL != "" {
		var redisStore = store.NewRedis(redis.NewClient(&redis.Options{
			Addr: Config.RedisURL,
		}))
		Cache = cache.New[[]byte](redisStore)
	} else {
		var bigcacheClient, _ = bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
		var bigcacheStore = store.NewBigcache(bigcacheClient)
		Cache = cache.New[[]byte](bigcacheStore)
	}
}
