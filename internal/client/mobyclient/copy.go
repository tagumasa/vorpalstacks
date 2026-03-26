// Package mobyclient provides Docker/Moby API client operations for copying files.
package mobyclient

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"vorpalstacks/internal/core/logs"

	"github.com/moby/moby/client"
)

// CopyIntoContainer copies files from local path into a container.
func (c *Client) CopyIntoContainer(ctx context.Context, containerID string, srcPath string, destPath string) error {
	c.logger.Debug("Copying files into container",
		logs.String("container", containerID),
		logs.String("src", srcPath),
		logs.String("dest", destPath),
	)

	stat, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to stat source path: %w", err)
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	if stat.IsDir() {
		if err := c.tarDirectory(tw, srcPath, ""); err != nil {
			return fmt.Errorf("failed to tar directory: %w", err)
		}
	} else {
		if err := c.tarFile(tw, srcPath, filepath.Base(srcPath)); err != nil {
			return fmt.Errorf("failed to tar file: %w", err)
		}
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}

	options := client.CopyToContainerOptions{
		DestinationPath:           destPath,
		Content:                   &buf,
		AllowOverwriteDirWithFile: true,
	}

	_, err = c.cli.CopyToContainer(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("failed to copy to container: %w", err)
	}

	c.logger.Info("Files copied into container",
		logs.String("container", containerID),
		logs.String("dest", destPath),
	)
	return nil
}

// CopyFromContainer copies files from a container to local path.
func (c *Client) CopyFromContainer(ctx context.Context, containerID string, srcPath string, destPath string) error {
	c.logger.Debug("Copying files from container",
		logs.String("container", containerID),
		logs.String("src", srcPath),
		logs.String("dest", destPath),
	)

	result, err := c.cli.CopyFromContainer(ctx, containerID, client.CopyFromContainerOptions{
		SourcePath: srcPath,
	})
	if err != nil {
		return fmt.Errorf("failed to copy from container: %w", err)
	}
	defer result.Content.Close()

	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	tr := tar.NewReader(result.Content)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		target := filepath.Join(destPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}
			outFile.Close()
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("failed to create symlink: %w", err)
			}
		}
	}

	c.logger.Info("Files copied from container",
		logs.String("container", containerID),
		logs.String("src", srcPath),
		logs.Any("size", result.Stat.Size),
	)
	return nil
}

// CreateFileInContainer creates a file with content in a container.
func (c *Client) CreateFileInContainer(ctx context.Context, containerID string, filePath string, content []byte) error {
	c.logger.Debug("Creating file in container",
		logs.String("container", containerID),
		logs.String("path", filePath),
	)

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	hdr := &tar.Header{
		Name: filepath.Base(filePath),
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("failed to write tar header: %w", err)
	}
	if _, err := tw.Write(content); err != nil {
		return fmt.Errorf("failed to write tar content: %w", err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}

	dirPath := filepath.Dir(filePath)
	if dirPath == "" {
		dirPath = "/"
	}

	options := client.CopyToContainerOptions{
		DestinationPath:           dirPath,
		Content:                   &buf,
		AllowOverwriteDirWithFile: true,
	}

	_, err := c.cli.CopyToContainer(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("failed to copy file to container: %w", err)
	}

	return nil
}

// tarDirectory recursively tars a directory.
func (c *Client) tarDirectory(tw *tar.Writer, srcDir string, base string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		if base != "" {
			relPath = filepath.Join(base, relPath)
		}

		if relPath == "." {
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		if strings.HasPrefix(relPath, ".") && strings.Contains(relPath, string(filepath.Separator)) {
			header.Name = strings.TrimPrefix(relPath, ".")
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tw, file)
		return err
	})
}

// tarFile tars a single file.
func (c *Client) tarFile(tw *tar.Writer, srcFile string, name string) error {
	info, err := os.Stat(srcFile)
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = name

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	file, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(tw, file)
	return err
}
