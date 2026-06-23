package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerCfg ServerConf
	LoggerCfg LoggerConf
	RedisCfg  RedisConf
}

type ServerConf struct {
	Port   int
	Origin string
}

type LoggerConf struct {
	Level      string
	OutputFile string
}

type RedisConf struct {
	Url             string
	DefaultCacheTTL int
}

var urlPattern = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*)`)

func InitConfig() (*Config, error) {
	var cfg = &Config{}

	// CLI configurations input
	if err := parseCLI(cfg); err != nil {
		return nil, err
	}

	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("failed to load '.env' file. origin: %s", err)
	}

	cfg.RedisCfg.Url = os.Getenv("REDIS_URL")
	if cfg.RedisCfg.Url == "" {
		return nil, errors.New("REDIS_URL is not provided in '.env' file or environment variables!")
	}

	return cfg, nil
}

func parseCLI(c *Config) error {
	flag.IntVar(
		&c.ServerCfg.Port, "port", 3000, "port to run the proxy server",
	)
	flag.StringVar(
		&c.ServerCfg.Origin, "origin", "", "origin server url (required)",
	)
	flag.StringVar(
		&c.LoggerCfg.Level,
		"log-level",
		"info",
		"logging level. OPTIONS: debug, info, warning, error, fatal",
	)
	flag.StringVar(
		&c.LoggerCfg.OutputFile,
		"log-file",
		"",
		"path to the log file. e.g. `./app.log` (default: os.Stdout)",
	)
	flag.IntVar(
		&c.RedisCfg.DefaultCacheTTL, "cache-ttl", 10, "default cache ttl in minutes. minimum: 1, maximum: 30",
	)
	flag.Parse()

	// validate origin url input
	if c.ServerCfg.Origin == "" {
		return errors.New("-origin flag is required")
	}

	if !urlPattern.MatchString(c.ServerCfg.Origin) {
		return fmt.Errorf("invalid -origin input. '%s' is not a valid url.", c.ServerCfg.Origin)
	}

	// validate port input
	if c.ServerCfg.Port <= 0 || c.ServerCfg.Port > 65535 {
		return fmt.Errorf("invalid -port input: '%d'. \nport must be between 1 and 65535", c.ServerCfg.Port)
	}

	// validate logging level input
	switch strings.ToLower(c.LoggerCfg.Level) {
	case "debug", "info", "warning", "error", "fatal":
	default:
		return fmt.Errorf(
			"invalid -log-level input: '%s'. \nvalid options: 'debug', 'info', 'warning', 'error', 'fatal'",
			c.LoggerCfg.Level,
		)
	}

	if c.LoggerCfg.OutputFile != "" {
		if !(strings.HasSuffix(c.LoggerCfg.OutputFile, ".log") || strings.HasSuffix(c.LoggerCfg.OutputFile, ".logs")) {
			return fmt.Errorf(
				"it's better the logger's output file ends with '.log' or '.logs'. \ncurrent input: '%s'",
				c.LoggerCfg.OutputFile,
			)
		}
	} // skip if OutputFile == "" -> output will be set to os.Stdout

	// check cache ttl input
	if c.RedisCfg.DefaultCacheTTL < 1 {
		c.RedisCfg.DefaultCacheTTL = 1
	} else if c.RedisCfg.DefaultCacheTTL > 30 {
		c.RedisCfg.DefaultCacheTTL = 30
	}

	return nil
}
