package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/hamidgh01/cache-proxy/internal/conf"
)

type RedisIntegration struct {
	*redis.Client
	ctx             context.Context
	DefaultCacheTTL time.Duration
}

var Redis *RedisIntegration

func InitRedis(c *conf.Configurations) {
	options, err := redis.ParseURL(c.RedisURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse Redis options from REDIS_URL: %v", err))
	}
	Redis = &RedisIntegration{
		redis.NewClient(options),
		context.Background(),
		c.DefaultCacheTTL,
	}

	// Test the connection
	_, err = Redis.Ping(Redis.ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	log.Println("Connected to Redis") // log.info
}
