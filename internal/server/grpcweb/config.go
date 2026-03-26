// Package grpcweb provides gRPC-Web server functionality for vorpalstacks.
package grpcweb

// Config holds configuration for the gRPC-Web server.
type Config struct {
	Port     string
	BindAddr string
}

// DefaultPort returns the default port for the server, or the configured port if set.
func (c *Config) DefaultPort() string {
	if c.Port == "" {
		return "9090"
	}
	return c.Port
}

func (c *Config) getBindAddr() string {
	if c.BindAddr != "" {
		return c.BindAddr
	}
	return "127.0.0.1"
}
