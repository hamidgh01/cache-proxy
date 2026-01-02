package server

import (
	"log"
	"net/http"
	"time"

	"github.com/hamidgh01/cache-proxy/internal/cache"
	"github.com/hamidgh01/cache-proxy/internal/conf"
)

type ProxyServer struct {
	OriginDomain string
	httpClient   *http.Client
}

func NewProxyServer(c *conf.Configurations) *ProxyServer {
	return &ProxyServer{
		OriginDomain: c.Origin,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// construct the target URL
	targetURL := p.OriginDomain + r.URL.String() // String(): Path + Query
	log.Printf("Received %s request for '%s'", r.Method, targetURL)

	// if not cacheable -> serve through origin
	if !isCacheable(r) {
		p.forwardToOrigin(w, r, targetURL)
		return
	}
	// if cacheable:
	// try to get and serve from cache (Cache Lookup)
	entry, err := cache.Redis.Fetch(targetURL)
	if err == nil {
		p.serveFromCache(w, &entry)
		log.Printf("'%s %s' is served from cache (CACHE HIT)", r.Method, targetURL) // log.info
		return
	}
	// if not cached before: get from origin, then cache, and then serve
	p.fetchAndCache(w, r, targetURL)
}
