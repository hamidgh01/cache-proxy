package server

import (
	"io"
	"log"
	"maps"
	"net/http"

	"github.com/hamidgh01/cache-proxy/internal/cache"
)

func (p *ProxyServer) serveFromCache(w http.ResponseWriter, e *cache.CacheEntry) {
	maps.Copy(w.Header(), e.Headers)
	w.Header().Set("X-Cache", "HIT")
	w.WriteHeader(http.StatusOK)
	w.Write(e.Body)
}

func (p *ProxyServer) fetchFromOrigin(
	r *http.Request, w http.ResponseWriter, targetURL string,
) *http.Response {
	// Prepare outgoing request (Origin Server Connection Manager)
	outReq, _ := http.NewRequest(r.Method, targetURL, r.Body)
	maps.Copy(outReq.Header, r.Header)

	response, err := p.httpClient.Do(outReq)
	if err != nil {
		http.Error(w, "Origin unreachable", http.StatusBadGateway)
		log.Printf("Origin unreachable for '%s %s'. Error message: %s\n", r.Method, targetURL, err) // log.info
		return nil
	}
	return response
}

func (p *ProxyServer) forwardToOrigin(
	w http.ResponseWriter, r *http.Request, targetURL string,
) {
	// simple passthrough for non-cacheable requests
	resp := p.fetchFromOrigin(r, w, targetURL)
	if resp == nil {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	maps.Copy(w.Header(), resp.Header)
	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Printf("'%s %s' (non-cacheable) is served through origin.\n", r.Method, targetURL) // log.info
}

func (p *ProxyServer) fetchAndCache(
	w http.ResponseWriter, r *http.Request, targetURL string,
) {
	resp := p.fetchFromOrigin(r, w, targetURL)
	if resp == nil {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// first cache the response (if cacheable)
	if isResponseCacheable(resp) {
		if err := cache.Redis.Save(resp, body, targetURL); err != nil {
			log.Printf(
				"Failed to cache response for '%s %s'. Error message: %s\n", r.Method, targetURL, err,
			) // log.error
		} else {
			log.Printf("response for '%s %s' cached successfully!\n", r.Method, targetURL) // log.info
		}
	}
	// then send response to client
	maps.Copy(w.Header(), resp.Header)
	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Printf("'%s %s' is served through origin.\n", r.Method, targetURL) // log.info
}
