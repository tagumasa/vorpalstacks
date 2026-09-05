// Package mobyclient provides Docker/Moby API client operations for container logs.
package mobyclient

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"

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

// ReadLogWindow returns the stdout and stderr written after `since` with
// the per-line docker timestamps stripped, plus the newest line's
// timestamp — the exact cursor for the next window. A zero `since` reads
// the whole log; an empty window leaves `since` unchanged. Docker's since
// filter is inclusive — a record whose timestamp equals `since` is
// delivered again — so the window re-filters client-side and the net
// contract is strictly-after. Consecutive windows therefore partition the
// log by record timestamp; the only loss a timestamp cursor can suffer is
// a record stamped at or below the cursor by a wall-clock step backwards
// between two reads, which the docker logs API offers no handle to
// distinguish from re-delivered output.
func (c *Client) ReadLogWindow(ctx context.Context, containerID string, since time.Time) (string, time.Time, error) {
	sinceStr := ""
	if !since.IsZero() {
		sinceStr = since.Format(time.RFC3339Nano)
	}
	opts := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     false,
		Since:      sinceStr,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		return "", since, fmt.Errorf("failed to read log window: %w", err)
	}
	defer reader.Close()

	raw, err := io.ReadAll(reader)
	if err != nil {
		return "", since, fmt.Errorf("failed to read log window body: %w", err)
	}
	return parseLogWindow(demuxLogStream(raw), since)
}

// demuxLogStream strips docker's non-TTY stream framing: each frame is an
// 8-byte header (stream byte, 3 padding bytes, big-endian payload length)
// followed by its payload, and stdout with stderr frames interleave in
// write order — the log window consumes them as one merged text, the
// shape CloudWatch interleaves console output in.
func demuxLogStream(raw []byte) string {
	var b strings.Builder
	for len(raw) >= 8 {
		length := binary.BigEndian.Uint32(raw[4:8])
		if uint64(length) > uint64(len(raw)-8) {
			// A truncated trailing frame means the read ended mid-write;
			// nothing further can be decoded.
			break
		}
		b.Write(raw[8 : 8+length])
		raw = raw[8+length:]
	}
	return b.String()
}

// parseLogWindow strips the RFC3339Nano timestamp prefix docker prepends
// to every log line when timestamps are requested and tracks the newest
// one as the next window's cursor. Every record at or before `since` is
// dropped, wherever it appears: the cursor is the newest timestamp of the
// previous window's read, a record delivered between two reads is stamped
// at its flush time and so lands above the cursor, and docker's inclusive
// filter re-delivers the boundary record — only records stamped strictly
// after the cursor can be new output.
func parseLogWindow(raw string, since time.Time) (string, time.Time, error) {
	if raw == "" {
		return "", since, nil
	}
	var out strings.Builder
	cursor := since
	for len(raw) > 0 {
		var line string
		if idx := strings.IndexByte(raw, '\n'); idx >= 0 {
			line, raw = raw[:idx], raw[idx+1:]
		} else {
			line, raw = raw, ""
		}
		line = strings.TrimSuffix(line, "\r")
		if line == "" {
			continue
		}
		ts, text, ok := strings.Cut(line, " ")
		t, err := time.Parse(time.RFC3339Nano, ts)
		if !ok || err != nil {
			return "", since, fmt.Errorf("unexpected log line without a docker timestamp: %q", line)
		}
		if !t.After(since) {
			continue
		}
		if t.After(cursor) {
			cursor = t
		}
		out.WriteString(text)
		out.WriteByte('\n')
	}
	return out.String(), cursor, nil
}
