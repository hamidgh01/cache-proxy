package conf

import (
	"time"
)

type Configurations struct {
	Port            string
	Origin          string
	LoggingLevel    string
	RedisURL        string
	DefaultCacheTTL time.Duration
}

func InitConfig() (*Configurations, error) {
	var config Configurations

	// configuration provided by '.env' file
	if err := extractFromDotEnv(&config); err != nil {
		return nil, err
	}
	// configuration presented in CLI
	if err := parseCLI(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
