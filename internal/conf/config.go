package conf

import (
	"time"
)

type Configurations struct {
	Port            string
	Origin          string
	LoggingLevel    string
	DefaultCacheTTL time.Duration
}

var Config Configurations

func InitConfig() error {
	// configuration provided by '.env' file
	if err := extractFromDotEnv(&Config); err != nil {
		return err
	}
	// configuration presented in CLI
	if err := parseCLI(&Config); err != nil {
		return err
	}

	return nil
}
