package memory

import (
	"github.com/patrickmn/go-cache"
)

type MemoryCache struct {
	cache *cache.Cache
}

func NewMemoryCache(cache *cache.Cache) *MemoryCache {
	return &MemoryCache{cache: cache}
}

func (m *MemoryCache) Get(key string) (interface{}, bool) {
	return m.cache.Get(key)
}

func (m *MemoryCache) Set(key string, data interface{}) {
	m.cache.Set(key, data, 0)
}

func (m *MemoryCache) Delete(key string) {
	m.cache.Delete(key)
}
