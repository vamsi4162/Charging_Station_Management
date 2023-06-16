package utils

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	once  sync.Once
	Cache *cache.Cache
)

const (
	CacheTTL = 5 * time.Minute // Cache TTL is set to 5 minutes
)

// InitCache initializes the cache with a default TTL
func InitCache() {
	once.Do(func() {
		Cache = cache.New(CacheTTL, CacheTTL)
	})
}

// GetCache returns the instance of the cache
func GetCache() *cache.Cache {
	return Cache
}
