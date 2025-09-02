package redis

import "github.com/redis/go-redis/v9"

func LoadCache(addr, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return rdb
}
