package config

import (
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/go-redis/redis/v8"
	gocache "github.com/patrickmn/go-cache"
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
		gocacheClient := gocache.New(5*time.Minute, 10*time.Minute)
		gocacheStore := store.NewGoCache(gocacheClient)
		Cache = cache.New[[]byte](gocacheStore)
	}
}
