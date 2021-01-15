package cache

import (
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)


type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

// NewRedisCache Create new Redis cache
func NewRedisCache(host string, db int, exp time.Duration) StringCache {
	return &redisCache{
		host: host,
		db: db,
		expires: exp,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cache.host,
		Password: "",
		DB: cache.db,
	})
}

func (cache *redisCache) Set(key string, value string) {
	client := cache.getClient()
	err := client.Set(key, value, cache.expires*time.Second).Err()
	if err != nil {
		log.Fatal("Redis Set", err)
	}
}

func (cache *redisCache) Get(key string) string {
	client := cache.getClient()
	val, err := client.Get(key).Result()
	if err != nil && err != redis.Nil {
		log.Fatal("Redis Get ", err)
		return ""
	}
	return val
}
