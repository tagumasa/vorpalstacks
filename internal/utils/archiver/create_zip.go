// Package archiver provides archive utility functions for vorpalstacks.
package archiver

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CreateZip creates a zip file from a directory.
//
// Parameters:
//   - sourceDir: The source directory to zip
//
// Returns:
//   - []byte: The zip file data
//   - error: An error if zipping fails
func CreateZip(sourceDir string) ([]byte, error) {
	info, err := os.Stat(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to stat source directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("source is not a directory: %s", sourceDir)
	}

	buf := new(bytes.Buffer)

	zipWriter := zip.NewWriter(buf)
	defer func() { _ = zipWriter.Close() }()

	absSourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filePath == sourceDir {
			return nil
		}

		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			return err
		}
		if !hasPrefix(absFilePath, absSourceDir) {
			return fmt.Errorf("path escapes source directory: %s", filePath)
		}

		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// #nosec G304
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer func() { _ = file.Close() }()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return buf.Bytes(), nil
}

func hasPrefix(s, prefix string) bool {
	if len(s) < len(prefix) || s[:len(prefix)] != prefix {
		return false
	}
	return len(s) == len(prefix) || s[len(prefix)] == '/'
}
