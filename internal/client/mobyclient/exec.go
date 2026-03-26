// Package mobyclient provides Docker/Moby API client operations for exec.
package mobyclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"vorpalstacks/internal/core/logs"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

// execReadCloser wraps ExecAttachResult to implement io.ReadCloser.
type execReadCloser struct {
	*client.ExecAttachResult
}

// Read implements the io.Reader interface.
func (e *execReadCloser) Read(p []byte) (n int, err error) {
	return e.ExecAttachResult.Reader.Read(p)
}

// Close closes the exec result.
func (e *execReadCloser) Close() error {
	if e.ExecAttachResult != nil {
		e.ExecAttachResult.Close()
	}
	return nil
}

// Exec executes a command in a container.
func (c *Client) Exec(ctx context.Context, containerID string, cfg ExecConfig) (*ExecResult, error) {
	c.logger.Debug("Executing command in container",
		logs.String("id", containerID),
		logs.Any("cmd", cfg.Cmd),
	)

	options := client.ExecCreateOptions{
		AttachStdout: cfg.AttachStdout,
		AttachStderr: cfg.AttachStderr,
		Cmd:          cfg.Cmd,
		Env:          cfg.Env,
		WorkingDir:   cfg.WorkingDir,
		Privileged:   cfg.Privileged,
	}

	result, err := c.cli.ExecCreate(ctx, containerID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create exec: %w", err)
	}

	attachResult, err := c.cli.ExecAttach(ctx, result.ID, client.ExecAttachOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to attach exec: %w", err)
	}
	defer attachResult.Close()

	var stdout, stderr bytes.Buffer
	if _, err = stdcopy.StdCopy(&stdout, &stderr, attachResult.Reader); err != nil {
		return nil, fmt.Errorf("failed to read exec output: %w", err)
	}

	inspectResult, err := c.cli.ExecInspect(ctx, result.ID, client.ExecInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect exec: %w", err)
	}

	return &ExecResult{
		ExitCode: inspectResult.ExitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}, nil
}

// ExecWithStdin executes a command in a container, piping stdin to the process.
func (c *Client) ExecWithStdin(ctx context.Context, containerID string, cfg ExecConfig, stdin io.Reader) (*ExecResult, error) {
	c.logger.Debug("Executing command in container with stdin",
		logs.String("id", containerID),
		logs.Any("cmd", cfg.Cmd),
	)

	options := client.ExecCreateOptions{
		AttachStdin:  true,
		AttachStdout: cfg.AttachStdout,
		AttachStderr: cfg.AttachStderr,
		Cmd:          cfg.Cmd,
		Env:          cfg.Env,
		WorkingDir:   cfg.WorkingDir,
		Privileged:   cfg.Privileged,
	}

	result, err := c.cli.ExecCreate(ctx, containerID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create exec with stdin: %w", err)
	}

	attachResult, err := c.cli.ExecAttach(ctx, result.ID, client.ExecAttachOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to attach exec with stdin: %w", err)
	}
	defer attachResult.Close()

	stdinDone := make(chan struct{})
	go func() {
		defer close(stdinDone)
		if _, err := io.Copy(attachResult.Conn, stdin); err != nil {
			log.Printf("Failed to copy stdin: %v", err)
		}
		if err := attachResult.CloseWrite(); err != nil {
			log.Printf("Failed to close write: %v", err)
		}
	}()

	var stdout, stderr bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdout, &stderr, attachResult.Reader); err != nil {
		<-stdinDone
		return nil, fmt.Errorf("failed to read exec output: %w", err)
	}
	<-stdinDone

	inspectResult, err := c.cli.ExecInspect(ctx, result.ID, client.ExecInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect exec: %w", err)
	}

	return &ExecResult{
		ExitCode: inspectResult.ExitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}, nil
}

// ExecWithOptions executes a command in a container with options.
func (c *Client) ExecWithOptions(ctx context.Context, containerID string, cfg ExecConfig, opts ExecOptions) (*ExecResult, error) {
	c.logger.Debug("Executing command in container with options",
		logs.String("id", containerID),
		logs.Any("cmd", cfg.Cmd),
	)

	createOptions := client.ExecCreateOptions{
		AttachStdout: opts.AttachStdout,
		AttachStderr: opts.AttachStderr,
		AttachStdin:  opts.AttachStdin,
		TTY:          opts.TTY,
		Cmd:          cfg.Cmd,
		Env:          opts.Env,
		WorkingDir:   opts.WorkingDir,
		Privileged:   opts.Privileged,
		User:         opts.User,
	}

	result, err := c.cli.ExecCreate(ctx, containerID, createOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create exec with options: %w", err)
	}

	attachResult, err := c.cli.ExecAttach(ctx, result.ID, client.ExecAttachOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to attach exec with options: %w", err)
	}
	defer attachResult.Close()

	var stdout, stderr bytes.Buffer
	if _, err = stdcopy.StdCopy(&stdout, &stderr, attachResult.Reader); err != nil {
		return nil, fmt.Errorf("failed to read exec output with options: %w", err)
	}

	inspectResult, err := c.cli.ExecInspect(ctx, result.ID, client.ExecInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect exec with options: %w", err)
	}

	return &ExecResult{
		ExitCode: inspectResult.ExitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}, nil
}

// ExecStream executes a command in a container and streams output.
func (c *Client) ExecStream(ctx context.Context, containerID string, cfg ExecConfig) (io.ReadCloser, error) {
	c.logger.Debug("Executing command in container (stream)",
		logs.String("id", containerID),
		logs.Any("cmd", cfg.Cmd),
	)

	options := client.ExecCreateOptions{
		AttachStdout: cfg.AttachStdout,
		AttachStderr: cfg.AttachStderr,
		Cmd:          cfg.Cmd,
		Env:          cfg.Env,
		WorkingDir:   cfg.WorkingDir,
		Privileged:   cfg.Privileged,
	}

	result, err := c.cli.ExecCreate(ctx, containerID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create exec stream: %w", err)
	}

	attachResult, err := c.cli.ExecAttach(ctx, result.ID, client.ExecAttachOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to attach exec stream: %w", err)
	}

	return &execReadCloser{ExecAttachResult: &attachResult}, nil
}

// ExecInteractive executes a command in a container with interactive I/O.
func (c *Client) ExecInteractive(ctx context.Context, containerID string, cfg ExecConfig, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	c.logger.Debug("Executing command in container (interactive)",
		logs.String("id", containerID),
		logs.Any("cmd", cfg.Cmd),
	)

	options := client.ExecCreateOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		TTY:          true,
		Cmd:          cfg.Cmd,
		Env:          cfg.Env,
		WorkingDir:   cfg.WorkingDir,
		Privileged:   cfg.Privileged,
	}

	result, err := c.cli.ExecCreate(ctx, containerID, options)
	if err != nil {
		return -1, fmt.Errorf("failed to create interactive exec: %w", err)
	}

	attachResult, err := c.cli.ExecAttach(ctx, result.ID, client.ExecAttachOptions{})
	if err != nil {
		return -1, fmt.Errorf("failed to attach interactive exec: %w", err)
	}
	defer attachResult.Close()

	stdinDone := make(chan struct{})
	go func() {
		defer close(stdinDone)
		if _, err := io.Copy(attachResult.Conn, stdin); err != nil {
			log.Printf("Failed to copy stdin in interactive exec: %v", err)
		}
	}()

	_, err = io.Copy(stdout, attachResult.Reader)
	if err != nil {
		<-stdinDone
		return -1, fmt.Errorf("failed to read interactive exec output: %w", err)
	}

	<-stdinDone

	inspectResult, err := c.cli.ExecInspect(ctx, result.ID, client.ExecInspectOptions{})
	if err != nil {
		return -1, fmt.Errorf("failed to inspect interactive exec: %w", err)
	}

	return inspectResult.ExitCode, nil
}

// ExecResize resizes the exec terminal.
func (c *Client) ExecResize(ctx context.Context, execID string, height, width uint) error {
	options := client.ExecResizeOptions{
		Height: height,
		Width:  width,
	}

	_, err := c.cli.ExecResize(ctx, execID, options)
	if err != nil {
		return fmt.Errorf("failed to resize exec: %w", err)
	}
	return nil
}
