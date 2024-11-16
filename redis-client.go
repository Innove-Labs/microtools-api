package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

func initRedis() {
	log.Println("Connecting to Redis...")
	config := LoadConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisURI,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}

func IncrementAPIHit(endpoint string) {
	key := "api_count:" + endpoint
	rdb.Incr(ctx, key)
}
