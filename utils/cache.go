package utils

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache *cache.Cache

func InitCache() {
	Cache = cache.New(5*time.Minute, 10*time.Minute)
}

func SetCache(key string, value interface{}, duration time.Duration) {
	Cache.Set(key, value, duration)
}

func GetCache(key string) (interface{}, bool) {
	return Cache.Get(key)
}
