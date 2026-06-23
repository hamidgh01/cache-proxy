package main

import (
	"fmt"
	"os"

	"github.com/hamidgh01/cache-proxy/config"
	"github.com/hamidgh01/cache-proxy/internal/cache"
	"github.com/hamidgh01/cache-proxy/internal/server"
	"github.com/hamidgh01/cache-proxy/pkg/logger"
)

func main() {
	// init configurations
	config, err := config.InitConfig()
	if err != nil {
		fmt.Println("failed to init configurations. reason: ", err)
		os.Exit(1)
	}

	// setup logger
	logger := logger.NewLogger(config.LoggerCfg)

	// establish redis connection
	cache.InitRedis(config.RedisCfg)

	// init and run proxy server
	proxyServer := server.NewProxyServer(config.ServerCfg, logger)
	if err := proxyServer.Run(); err != nil {
		fmt.Println("failed to run cache proxy server. reason: ", err)
		// close redis
		os.Exit(1)
	}

	// ToDo: add graceful shutdown
}
