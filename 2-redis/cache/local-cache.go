package cache

type localCache struct {
	data map[string]string
}

// NewLocalCache Create new local cache
func NewLocalCache() StringCache {
	return &localCache{
		data: map[string]string{},
	}
}

func (lc *localCache) Set(key, val string) {
	lc.data[key] = val
}

func (lc *localCache) Get(key string) string {
	return lc.data[key]
}
