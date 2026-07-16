// Package converter provides tools for converting Smithy models.
package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var serviceFilePattern = regexp.MustCompile(`^([a-z0-9-]+)-(\d{4}-\d{2}-\d{2})\.json$`)

// ServiceFile represents a Smithy service JSON file
type ServiceFile struct {
	ServiceName string
	FilePath    string
	Version     string
}

// DirectoryScanner scans directories for Smithy service JSON files
type DirectoryScanner struct {
	baseDir string
}

// NewDirectoryScanner creates a new DirectoryScanner
func NewDirectoryScanner(baseDir string) *DirectoryScanner {
	return &DirectoryScanner{
		baseDir: baseDir,
	}
}

// Scan scans the base directory for Smithy service JSON files
// Expected pattern: models/<service-name>/service/<date>/<service-name>-<date>.json
func (s *DirectoryScanner) Scan() ([]ServiceFile, error) {
	// Map to store service files by service name
	serviceMap := make(map[string]ServiceFile)

	// Walk through the directory
	err := filepath.Walk(s.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file matches the expected pattern
		serviceFile, ok := s.parseServicePath(path)
		if !ok {
			return nil
		}

		// Keep the latest version for each service
		if existing, exists := serviceMap[serviceFile.ServiceName]; !exists {
			serviceMap[serviceFile.ServiceName] = serviceFile
		} else {
			existingTime, err1 := time.Parse("2006-01-02", existing.Version)
			newTime, err2 := time.Parse("2006-01-02", serviceFile.Version)

			switch {
			case err1 != nil && err2 != nil:
				serviceMap[serviceFile.ServiceName] = serviceFile
			case err1 != nil:
				serviceMap[serviceFile.ServiceName] = serviceFile
			case err2 != nil:
			case newTime.After(existingTime):
				serviceMap[serviceFile.ServiceName] = serviceFile
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(serviceMap) == 0 {
		return nil, fmt.Errorf("no service files found in directory: %s", s.baseDir)
	}

	services := make([]ServiceFile, 0, len(serviceMap))
	for _, service := range serviceMap {
		services = append(services, service)
	}
	sort.Slice(services, func(i, j int) bool {
		return services[i].ServiceName < services[j].ServiceName
	})

	return services, nil
}

// parseServicePath parses a file path and extracts service information
// Expected pattern: .../service/<date>/<service-name>-<date>.json
func (s *DirectoryScanner) parseServicePath(path string) (ServiceFile, bool) {
	// Get relative path from base directory
	relPath, err := filepath.Rel(s.baseDir, path)
	if err != nil {
		return ServiceFile{}, false
	}

	// Check if the path contains "service" directory
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	serviceDirIndex := -1
	for i, part := range parts {
		if part == "service" {
			serviceDirIndex = i
			break
		}
	}

	// If "service" directory is not found, check if the base directory itself is "service"
	// In this case, the pattern is: <date>/<service-name>-<date>.json
	if serviceDirIndex == -1 {
		// Check if the base directory ends with "service"
		if filepath.Base(filepath.Clean(s.baseDir)) != "service" {
			return ServiceFile{}, false
		}
		// Pattern: <date>/<service-name>-<date>.json
		if len(parts) < 2 {
			return ServiceFile{}, false
		}
		filename := parts[1]

		// Check if filename matches pattern: <service-name>-<date>.json
		matches := serviceFilePattern.FindStringSubmatch(filename)
		if matches == nil {
			return ServiceFile{}, false
		}

		return ServiceFile{
			ServiceName: matches[1],
			FilePath:    path,
			Version:     matches[2],
		}, true
	}

	// Standard pattern: .../service/<date>/<service-name>-<date>.json
	if serviceDirIndex+2 >= len(parts) {
		return ServiceFile{}, false
	}

	// Extract service name from the directory before "service"
	serviceName := ""
	if serviceDirIndex > 0 {
		serviceName = parts[serviceDirIndex-1]
	}

	// Extract filename
	filename := parts[serviceDirIndex+2]

	// Check if filename matches pattern: <service-name>-<date>.json
	matches := serviceFilePattern.FindStringSubmatch(filename)
	if matches == nil {
		return ServiceFile{}, false
	}

	// Verify that the service name in the filename matches the directory name
	fileServiceName := matches[1]
	if serviceName != "" && serviceName != fileServiceName {
		return ServiceFile{}, false
	}

	return ServiceFile{
		ServiceName: fileServiceName,
		FilePath:    path,
		Version:     matches[2],
	}, true
}
