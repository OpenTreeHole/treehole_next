package utils

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/goccy/go-json"
	gocache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"treehole_next/config"
)

var Cache *cache.Cache[any]

func InitCache() {
	if config.Config.RedisURL == "" {
		useGoCache()
		return
	}
	redisClient, err := newRedisClient(config.Config.RedisURL)
	if err != nil {
		log.Warn().Err(err).Str("redis_url", config.Config.RedisURL).Msg("redis init failed, fallback to go-cache")
		useGoCache()
		return
	}
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		_ = redisClient.Close()
		log.Warn().Err(err).Str("redis_url", config.Config.RedisURL).Msg("redis ping failed, fallback to go-cache")
		useGoCache()
		return
	}
	Cache = cache.New[any](redis_store.NewRedis(redisClient))
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	if strings.Contains(redisURL, "://") {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			return nil, err
		}
		return redis.NewClient(opt), nil
	}
	return redis.NewClient(&redis.Options{Addr: redisURL}), nil
}

func useGoCache() {
	gocacheStore := gocache_store.NewGoCache(gocache.New(5*time.Minute, 10*time.Minute))
	Cache = cache.New[any](gocacheStore)
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

// GetCache gets a value from cache by key and unmarshals it into value (must be a pointer).
// It supports both Redis store (returns string) and go-cache store (returns []byte).
// Returns true if the key exists and JSON unmarshal succeeds, false otherwise.
func GetCache(key string, value any) bool {
	raw, err := Cache.Get(context.Background(), key)
	if err != nil {
		return false
	}
	// Normalize to []byte: Redis store returns string, go-cache returns []byte.
	var data []byte
	switch v := raw.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		log.Warn().Str("key", key).Msg("cache value type not []byte or string")
		return false
	}
	// Strip leading null bytes to avoid "invalid character '\\u0000'" from legacy or cross-app data.
	data = bytes.TrimLeft(data, "\x00")
	if len(data) == 0 {
		return false
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		log.Warn().Err(err).Str("key", key).Str("data", string(data)).Msg("unmarshal cache failed")
	}
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
