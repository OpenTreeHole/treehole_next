package utils

import (
	"context"
	"github.com/eko/gocache/v3/store"
	"github.com/goccy/go-json"
	"time"
	"treehole_next/config"
)

const maxDuration time.Duration = 1<<63 - 1

func SetCache(key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if expiration == 0 {
		expiration = maxDuration
	}
	return config.Cache.Set(context.Background(), key, data, store.WithExpiration(expiration))
}

func GetCache(key string, value any) bool {
	data, err := config.Cache.Get(context.Background(), key)
	if err != nil {
		return false
	}
	err = json.Unmarshal(data, value)
	return err == nil
}

func DeleteCache(key string) error {
	err := config.Cache.Delete(context.Background(), key)
	if err == nil {
		return nil
	}
	if err.Error() == "Entry not found" {
		return nil
	}
	return err
}
