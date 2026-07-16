// Package mobyclient provides Docker/Moby API client implementation.
package mobyclient

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/logs"

	"github.com/moby/moby/client"
)

// ContainerSockDir is the directory to mount sockets in containers.
const (
	ContainerSockDir = "/var/opt/socket"
)

// ContainerLifecycle abstracts the Docker container operations required by
// services that manage their own containers (Lambda, RDS, etc.).  Defining an
// interface decouples consumers from the concrete mobyclient.Client and allows
// testing without a running Docker daemon.
type ContainerLifecycle interface {
	GetContainerStatus(ctx context.Context, containerID string) (ContainerStatus, error)
	CreateContainerFromConfig(ctx context.Context, cfg AdvancedContainerConfig) (*CreateContainerResult, error)
	StartContainer(ctx context.Context, containerID string) error
	RemoveContainer(ctx context.Context, containerID string, force bool) error
	CreateFileInContainer(ctx context.Context, containerID string, path string, content []byte) error
	Exec(ctx context.Context, containerID string, cfg ExecConfig) (*ExecResult, error)
	Close() error
}

// Client provides Docker/Moby API operations.
type Client struct {
	cli    *client.Client
	logger logs.Logger
	config Config
}

// NewClient creates a new mobyclient.
func NewClient(cfg Config, logger logs.Logger) (*Client, error) {
	opts := []client.Opt{
		client.WithHost(cfg.Host),
	}

	if cfg.Version != "" {
		opts = append(opts, client.WithAPIVersion(cfg.Version))
	}

	if cfg.TLSVerify {
		opts = append(opts, client.WithTLSClientConfig(
			cfg.CertPath+"/ca.pem",
			cfg.CertPath+"/cert.pem",
			cfg.CertPath+"/key.pem",
		))
	}

	cli, err := client.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	logger.Info("MobyClient initialized",
		logs.String("host", cfg.Host),
		logs.String("version", cfg.Version),
	)

	return &Client{
		cli:    cli,
		logger: logger,
		config: cfg,
	}, nil
}

// NewClientFromEnv creates a new mobyclient from environment variables.
func NewClientFromEnv(logger logs.Logger) (*Client, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client from env: %w", err)
	}

	logger.Info("MobyClient initialized from environment")

	return &Client{
		cli:    cli,
		logger: logger,
		config: DefaultConfig(),
	}, nil
}

// Close closes the client.
func (c *Client) Close() error {
	if c.cli != nil {
		return c.cli.Close()
	}
	return nil
}

// Ping checks connectivity to Docker daemon.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx, client.PingOptions{})
	return err
}

// Version returns Docker daemon version info.
func (c *Client) Version(ctx context.Context) (string, error) {
	v, err := c.cli.ServerVersion(ctx, client.ServerVersionOptions{})
	if err != nil {
		return "", err
	}
	return v.Version, nil
}
