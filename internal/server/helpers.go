package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func isCacheable(r *http.Request) bool {
	// Only responses of GET requests are cached
	if r.Method != http.MethodGet {
		return false
	}
	// Authorization header -> do not cache
	if r.Header.Get("Authorization") != "" {
		return false
	}

	return true
}

func isResponseCacheable(resp *http.Response) (isCacheable bool, cacheTTl time.Duration) {
	// must: have 'status: 200 OK' , not have 'Set-Cookie: ...' header , not have 'Vary: *' header
	// 2. if Set-Cookie present -> user-specific -> do not cache
	// 3. if 'Vary: *' -> Implies that the response is uncacheable. (https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Vary)
	if resp.StatusCode != http.StatusOK || resp.Header.Get("Set-Cookie") != "" || resp.Header.Get("Vary") == "*" {
		return false, 0
	}

	// analyze Cache-Control header (if present)
	if shouldCache, cacheTTl := analyzeCacheControlHeader(resp); !shouldCache {
		return false, cacheTTl // if shouldCache'd be false -> cacheTTl'll be 0
	} else {
		return true, cacheTTl
	}
}

// Analyzes Cache-Control header and determines the response should cache or not,
// and if should cache: tries to extract caching-ttl (from 's-maxage' or 'max-age').
func analyzeCacheControlHeader(resp *http.Response) (shouldCache bool, cacheTTl time.Duration) {
	ccHeader := resp.Header.Get("Cache-Control")
	if ccHeader == "" {
		return true, 0 // cacheTTl=0 --> means use DefaultCacheTTL for expiration
	}

	// parse Cache-Control header
	parsedCCHeader := make(map[string]string)
	parts := strings.SplitSeq(ccHeader, ",")
	for part := range parts {
		part = strings.TrimSpace(part)

		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			parsedCCHeader[strings.ToLower(kv[0])] = strings.Trim(kv[1], `"`)
		} else {
			parsedCCHeader[strings.ToLower(part)] = ""
		}
	}

	// no-store -> must not cache
	if _, ok := parsedCCHeader["no-store"]; ok {
		return false, 0
	}
	// no-cache -> requires revalidation (skip for now)
	if _, ok := parsedCCHeader["no-cache"]; ok {
		return false, 0
	}

	// s-maxage (proxy-specific)
	if v, ok := parsedCCHeader["s-maxage"]; ok {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			return true, time.Duration(secs) * time.Second
		}
	}
	// max-age
	if v, ok := parsedCCHeader["max-age"]; ok {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			return true, time.Duration(secs) * time.Second
		}
	}

	return true, 0
}
