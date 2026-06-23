package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hamidgh01/cache-proxy/config"

	"github.com/redis/go-redis/v9"
)

type RedisIntegration struct {
	*redis.Client
	ctx             context.Context
	DefaultCacheTTL time.Duration
}

var Redis *RedisIntegration

func InitRedis(cfg config.RedisConf) {
	options, err := redis.ParseURL(cfg.Url)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse Redis options from REDIS_URL: %v", err))
	}
	Redis = &RedisIntegration{
		redis.NewClient(options),
		context.Background(),
		time.Duration(cfg.DefaultCacheTTL) * time.Minute,
	}

	// Test the connection
	_, err = Redis.Ping(Redis.ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	log.Println("Connected to Redis") // log.info
}

func CloseRedis() {
	// to implement
	// close redis at shutdown to prevent memory leaks
}
