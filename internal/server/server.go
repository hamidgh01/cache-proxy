package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hamidgh01/cache-proxy/config"
	"github.com/hamidgh01/cache-proxy/internal/cache"
	"github.com/hamidgh01/cache-proxy/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type ProxyServer struct {
	ctx          context.Context
	originUrl    string
	proxyPort    int
	httpClient   *http.Client
	cacheService *cache.CacheService
	logger       *logger.Logger
}

func NewProxyServer(cfg config.ServerConf, l *logger.Logger, c *cache.CacheService) *ProxyServer {
	return &ProxyServer{
		ctx:          context.Background(),
		originUrl:    cfg.Origin,
		proxyPort:    cfg.Port,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		cacheService: c,
		logger:       l,
	}
}

func (p *ProxyServer) Run() error {
	address := fmt.Sprintf("localhost:%d", p.proxyPort)

	p.logger.Infof(
		"running cache proxy server on port '%d', forwarding to '%s'\n",
		p.proxyPort,
		p.originUrl,
	)

	return http.ListenAndServe(address, p)
}

// how this proxy works ???
//  1. if incoming requests is not cacheable -> forward and serve through origin
//  2. if is cacheable
//     2.1 Cache Lookup: try to get and serve from cache (if cached before)
//     2.2 if not cached before: get from origin, then cache, and then serve
func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// construct the target URL
	targetURL := p.originUrl + r.URL.String() // URL.String(): Path + Query
	p.logger.Infof("received '%s %s'", r.Method, targetURL)

	// 1. if not cacheable -> serve through origin
	if !isCacheable(r) {
		p.forwardToOrigin(w, r, targetURL)
		return
	}

	// 2. if cacheable:

	// 2.1 try to get and serve from cache (Cache Lookup)
	entry, err := p.cacheService.Fetch(p.ctx, targetURL)
	switch err {
	case nil: // serve directly from cache
		sendFinalResponse(w, entry.Headers, "HIT", entry.Status, entry.Body)
		p.logger.Infof("(CACHE HIT) '%s %s' is served from cache", r.Method, targetURL)
		return
	case redis.Nil: // not cached before
	default: // cached before, but there's CacheService error
		p.logger.Errorf("failed to serve from cache. reason: %s", err.Error())
	}

	// 2.2 if not cached before: get from origin, then cache, and then serve
	p.fetchFromOriginAndCacheThenServe(w, r, targetURL)
}
