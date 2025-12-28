package main

import (
	"fmt"

	"github.com/hamidgh01/cache-proxy/internal/conf"
)

func main() {

	// step_1 : parse CLI args & init configurations
	if err := conf.InitConfig(); err != nil {
		panic(fmt.Sprintf("Failure at 'conf.InitConfig'. Error Message: %s", err))
	}

	fmt.Println("DefaultCacheTTL:", conf.Config.DefaultCacheTTL)
	fmt.Println("Origin:", conf.Config.Origin)
	fmt.Println("Port:", conf.Config.Port)
	fmt.Println("LoggingLevel:", conf.Config.LoggingLevel)
}
