package memorycache

import "sync"

type MemoryCache struct {
	cache map[string]string
	mutex sync.Mutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: make(map[string]string),
	}
}

func (mc *MemoryCache) Get(key string) (string, bool) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	value, ok := mc.cache[key]
	return value, ok
}

func (mc *MemoryCache) Set(key string, value string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.cache[key] = value
}
