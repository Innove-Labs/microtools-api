package database

import (
	"github.com/go-redis/redis/v8"
	"github.com/innovelabs/microtools-go/internal/config"
)

// InitRedis initializes Redis client
func InitRedis() *redis.Client {
	cfg := config.LoadConfig()
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisURI,
	})
}
