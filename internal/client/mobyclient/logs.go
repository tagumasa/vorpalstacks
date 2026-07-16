// Package mobyclient provides Docker/Moby API client operations for container logs.
package mobyclient

import (
	"context"
	"fmt"
	"io"

	"vorpalstacks/internal/core/logs"

	"github.com/moby/moby/client"
)

// GetLogs gets container logs.
func (c *Client) GetLogs(ctx context.Context, containerID string, stdout, stderr bool) (string, error) {
	opts := client.ContainerLogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
		Timestamps: false,
		Follow:     false,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(logs), nil
}

// GetLogsWithOptions gets container logs with options.
func (c *Client) GetLogsWithOptions(ctx context.Context, containerID string, opts LogsOptions) (string, error) {
	dockerOpts := client.ContainerLogsOptions{
		ShowStdout: opts.ShowStdout,
		ShowStderr: opts.ShowStderr,
		Timestamps: opts.Timestamps,
		Follow:     opts.Follow,
		Tail:       opts.Tail,
		Since:      opts.Since,
		Until:      opts.Until,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, dockerOpts)
	if err != nil {
		return "", fmt.Errorf("failed to get logs with options: %w", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(logs), nil
}

// StreamLogs streams container logs.
func (c *Client) StreamLogs(ctx context.Context, containerID string, stdout, stderr bool) (io.ReadCloser, error) {
	opts := client.ContainerLogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
		Timestamps: true,
		Follow:     true,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to stream logs: %w", err)
	}

	return reader, nil
}

// StreamLogsWithOptions streams container logs with options.
func (c *Client) StreamLogsWithOptions(ctx context.Context, containerID string, opts LogsOptions) (io.ReadCloser, error) {
	dockerOpts := client.ContainerLogsOptions{
		ShowStdout: opts.ShowStdout,
		ShowStderr: opts.ShowStderr,
		Timestamps: opts.Timestamps,
		Follow:     opts.Follow,
		Tail:       opts.Tail,
		Since:      opts.Since,
		Until:      opts.Until,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, dockerOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to stream logs with options: %w", err)
	}

	return reader, nil
}

// FollowLogs follows container logs and writes to the provided writer.
func (c *Client) FollowLogs(ctx context.Context, containerID string, stdout, stderr bool, writer io.Writer) error {
	c.logger.Debug("Following logs", logs.String("id", containerID))

	opts := client.ContainerLogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
		Timestamps: true,
		Follow:     true,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		return fmt.Errorf("failed to follow logs: %w", err)
	}
	defer reader.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		return fmt.Errorf("failed to copy logs: %w", err)
	}

	return nil
}

// TailLogs gets the last N lines of container logs.
func (c *Client) TailLogs(ctx context.Context, containerID string, lines string, stdout, stderr bool) (string, error) {
	opts := client.ContainerLogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
		Timestamps: false,
		Follow:     false,
		Tail:       lines,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		return "", fmt.Errorf("failed to tail logs: %w", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(logs), nil
}

// LogsSince gets logs since a specific time.
func (c *Client) LogsSince(ctx context.Context, containerID string, since string, stdout, stderr bool) (string, error) {
	opts := client.ContainerLogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
		Timestamps: true,
		Follow:     false,
		Since:      since,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		return "", fmt.Errorf("failed to get logs since: %w", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(logs), nil
}

// LogsUntil gets logs until a specific time.
func (c *Client) LogsUntil(ctx context.Context, containerID string, until string, stdout, stderr bool) (string, error) {
	opts := client.ContainerLogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
		Timestamps: true,
		Follow:     false,
		Until:      until,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		return "", fmt.Errorf("failed to get logs until: %w", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(logs), nil
}
