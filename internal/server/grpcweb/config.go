// Package grpcweb provides gRPC-Web server functionality for vorpalstacks.
package grpcweb

import (
	"fmt"

	"vorpalstacks/internal/common/serviceports"
)

// Config holds configuration for the gRPC-Web server.
type Config struct {
	Port     int
	BindAddr string
}

// DefaultPort returns the default port for the server, or the configured port if set.
func (c *Config) DefaultPort() int {
	if c.Port == 0 {
		return serviceports.GRPCWeb
	}
	return c.Port
}

// DefaultPortString returns the default port as a string.
func (c *Config) DefaultPortString() string {
	return fmt.Sprintf("%d", c.DefaultPort())
}

func (c *Config) getBindAddr() string {
	if c.BindAddr != "" {
		return c.BindAddr
	}
	return "127.0.0.1"
}
