package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/hamidgh01/cache-proxy/internal/cache"
	"github.com/hamidgh01/cache-proxy/internal/conf"
	"github.com/hamidgh01/cache-proxy/internal/server"
)

func main() {

	// step_1: load .env file
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("Failed to load '.env' file. Error Message: %s", err))
	}

	// step_2: parse CLI args & init configurations
	config, err := conf.InitConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to initial configurations. Error Message: %s", err))
	}

	// step_3: initialize redis connection
	cache.InitRedis(config)

	// step_4: start the proxy server
	fmt.Printf(
		"Caching Proxy Server running on port '%s', forwarding to '%s'\n\n",
		config.Port,
		config.Origin,
	)
	proxyServer := server.NewProxyServer(config)
	log.Fatal(http.ListenAndServe(":"+config.Port, proxyServer))

	// ToDo: add graceful shutdown
}
