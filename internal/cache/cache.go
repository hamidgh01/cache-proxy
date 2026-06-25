package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/hamidgh01/cache-proxy/config"
	"github.com/hamidgh01/cache-proxy/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type cacheEntry struct {
	Status  int
	Headers http.Header
	Body    []byte
}

type CacheService struct {
	redis           *redis.Client
	defaultCacheTTL time.Duration
	logger          *logger.Logger
}

func NewCacheService(cfg config.RedisConf, l *logger.Logger) (*CacheService, func() error, error) {
	c := &CacheService{
		defaultCacheTTL: time.Minute * time.Duration(cfg.DefaultCacheTTL),
		logger:          l,
	}

	// establish redis connection
	redisClient, err := initRedis(cfg)
	if err != nil {
		return nil, nil, err
	}

	c.logger.Info("redis connection established successfully.")
	c.redis = redisClient

	return c, c.redis.Close, nil
}

func generateCacheKey(targetURL string) string {
	return fmt.Sprintf("CacheProxy:%s", targetURL)
}

func (c *CacheService) Save(
	ctx context.Context, resp *http.Response, respBody []byte, targetURL string, ttl time.Duration,
) error {
	// filter hop-by-hop headers
	cleanedHeaders := filterResponseHeaders(resp.Header)

	// create a CacheEntry from response
	entry := cacheEntry{
		Status:  resp.StatusCode,
		Headers: cleanedHeaders,
		Body:    respBody,
	}

	// encode to buffer
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode CacheEntry: %w", err)
	}

	// cache buffered data
	key := generateCacheKey(targetURL)
	if ttl == 0 {
		ttl = c.defaultCacheTTL
	}
	if err := c.redis.Set(ctx, key, buf.Bytes(), ttl).Err(); err != nil {
		return fmt.Errorf("failed to save encoded CacheEntry in Redis: %v", err)
	}

	return nil
}

func (c *CacheService) Fetch(ctx context.Context, targetURL string) (*cacheEntry, error) {
	// get bytes from redis
	key := generateCacheKey(targetURL)
	data, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to fetch cached data (ket='%s') from Redis: %w", key, err)
	}

	// decode fetched bytes to CacheEntry
	entry := &cacheEntry{}
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(entry); err != nil {
		return nil, fmt.Errorf("failed to decode fetched data to CacheEntry: %w", err)
	}

	return entry, nil
}
