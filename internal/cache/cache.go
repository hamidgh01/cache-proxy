package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
)

type CacheEntry struct {
	Status  int
	Headers http.Header
	Body    []byte
}

func generateCacheKey(targetURL string) string {
	return fmt.Sprintf("CacheProxy:%s", targetURL)
}

func (r *RedisIntegration) Save(resp *http.Response, respBody []byte, targetURL string) error {
	// 1. create a CacheEntry from response
	entry := CacheEntry{
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Body:    respBody,
	}
	// 2. encode to buffer
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("Failed to encode CacheEntry: %v", err)
	}
	// 3. cache buffered data
	key := generateCacheKey(targetURL)
	if err := r.Set(r.ctx, key, buf.Bytes(), r.DefaultCacheTTL).Err(); err != nil {
		return fmt.Errorf("Failed to save CacheEntry in Redis: %v", err) // log.warning
	}

	return nil
}

func (r *RedisIntegration) Fetch(targetURL string) (CacheEntry, error) {
	// 1. get bytes from redis
	key := generateCacheKey(targetURL)
	data, err := r.Get(r.ctx, key).Bytes()
	if err != nil {
		return CacheEntry{}, err
	}
	// 2. try to decode fetched bytes to CacheEntry
	var entry CacheEntry
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&entry); err != nil {
		// log.warning: "failure at decoding cached data to CacheEntry:", err
		return CacheEntry{}, err
	}

	return entry, nil
}
