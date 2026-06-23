package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hamidgh01/cache-proxy/config"
	"github.com/hamidgh01/cache-proxy/internal/cache"
	"github.com/hamidgh01/cache-proxy/pkg/logger"
)

type ProxyServer struct {
	originUrl  string
	proxyPort  int
	httpClient *http.Client
	logger     *logger.Logger
}

func NewProxyServer(cfg config.ServerConf, l *logger.Logger) *ProxyServer {
	return &ProxyServer{
		originUrl:  cfg.Origin,
		proxyPort:  cfg.Port,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		logger:     l,
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

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// construct the target URL
	targetURL := p.originUrl + r.URL.String() // URL.String(): Path + Query
	p.logger.Infof("Received %s request for '%s'", r.Method, targetURL)

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
		p.logger.Infof("'%s %s' is served from cache (CACHE HIT)", r.Method, targetURL)
		return
	}
	// if not cached before: get from origin, then cache, and then serve
	p.fetchAndCache(w, r, targetURL)
}
