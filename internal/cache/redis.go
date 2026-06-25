package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/hamidgh01/cache-proxy/config"

	"github.com/redis/go-redis/v9"
)

func initRedis(cfg config.RedisConf) (*redis.Client, error) {
	options, err := redis.ParseURL(cfg.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse '%s' (provided REDIS_URL) to '*redis.Options'. origin: %w", cfg.Url, err)
	}

	redisClient := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// test the connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed. origin: %w", err)
	}

	return redisClient, nil
}
