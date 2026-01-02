package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"github.com/hamidgh01/cache-proxy/internal/cache"
	"github.com/hamidgh01/cache-proxy/internal/conf"
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

	fmt.Println("DefaultCacheTTL:", config.DefaultCacheTTL)
	fmt.Println("Origin:", config.Origin)
	fmt.Println("Port:", config.Port)
	fmt.Println("LoggingLevel:", config.LoggingLevel)

	// step_3: initialize redis connection
	cache.InitRedis(config)
}
