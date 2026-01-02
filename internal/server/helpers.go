package server

import (
	"net/http"
	"strings"
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

func isResponseCacheable(resp *http.Response) bool {
	// must have: 'status: 200 OK'
	if resp.StatusCode != http.StatusOK {
		return false
	}
	// Set-Cookie present -> user-specific -> do not cache
	if resp.Header.Get("Set-Cookie") != "" {
		return false
	}
	// Don't cache if -> Cache-Control: no-store
	if cc := resp.Header.Get("Cache-Control"); cc != "" {
		if strings.Contains(strings.ToLower(cc), "no-store") {
			return false
		}
	}
	// Vary: *
	if resp.Header.Get("Vary") == "*" {
		return false
	}

	return true
}
