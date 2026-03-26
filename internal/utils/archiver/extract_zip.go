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

// ExtractZip extracts a zip file to a directory.
//
// Parameters:
//   - zipData: The zip file data as bytes
//   - destDir: The destination directory to extract to
//
// Returns:
//   - error: An error if extraction fails
func ExtractZip(zipData []byte, destDir string) error {
	if zipData == nil {
		return fmt.Errorf("zip data cannot be nil")
	}

	cleanDestDir := filepath.Clean(destDir)
	absDestDir, err := filepath.Abs(cleanDestDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("failed to read zip: %w", err)
	}

	for _, file := range reader.File {
		if file.FileInfo().Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink not allowed: %s", file.Name)
		}

		path := filepath.Join(cleanDestDir, file.Name) // #nosec G305
		cleanPath := filepath.Clean(path)

		absPath, err := filepath.Abs(cleanPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		if !hasPrefix(absPath, absDestDir) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(cleanPath, file.Mode()&0777); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(cleanPath), 0750); err != nil { // #nosec G301
			return fmt.Errorf("failed to create directory: %w", err)
		}

		outFile, err := os.OpenFile(cleanPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode()&0644) // #nosec G306
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		rc, err := file.Open()
		if err != nil {
			_ = outFile.Close()
			return fmt.Errorf("failed to open zip file: %w", err)
		}

		_, err = io.Copy(outFile, rc) // #nosec G110
		if cerr := rc.Close(); cerr != nil && err == nil {
			err = cerr
		}
		if ferr := outFile.Close(); ferr != nil && err == nil {
			err = ferr
		}

		if err != nil {
			return fmt.Errorf("failed to extract file: %w", err)
		}
	}

	return nil
}
