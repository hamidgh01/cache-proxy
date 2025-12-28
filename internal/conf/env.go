package conf

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func extractFromDotEnv(c *Configurations) error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	c.LoggingLevel = strings.ToLower(os.Getenv("LOGGING_LEVEL"))
	switch c.LoggingLevel {
	case "debug", "info", "warning", "error", "critical":
	case "":
		return errors.New("LOGGING_LEVEL is not provided in '.env' file or Environment Variables!")
	default:
		return errors.New(`LOGGING_LEVEL options: "debug", "info", "warning", "error", "critical"`)
	}

	DEFAULT_CACHE_TTL := os.Getenv("DEFAULT_CACHE_TTL")
	if DEFAULT_CACHE_TTL == "" {
		return errors.New("DEFAULT_CACHE_TTL is not provided in '.env' file or Environment Variables!")
	}
	DEFAULT_CACHE_TTL_INT, err := strconv.Atoi(DEFAULT_CACHE_TTL)
	if err != nil {
		return fmt.Errorf("DEFAULT_CACHE_TTL EnvVar must be an integer number! current value: %s", DEFAULT_CACHE_TTL)
	}
	c.DefaultCacheTTL = time.Duration(DEFAULT_CACHE_TTL_INT) * time.Minute

	return nil
}
