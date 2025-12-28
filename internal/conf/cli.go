package conf

import (
	"errors"
	"flag"
	"strconv"
)

func parseCLI(c *Configurations) error {
	var port int
	flag.IntVar(&port, "port", 3000, "Port to run the proxy server")
	flag.StringVar(&c.Origin, "origin", "", "Origin server URL")
	flag.Parse()

	if c.Origin == "" {
		return errors.New("--origin flag is required")
	}
	if port <= 0 || port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	c.Port = strconv.Itoa(port)

	return nil
}
