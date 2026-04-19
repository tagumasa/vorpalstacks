// Package naming provides naming and string manipulation utilities for vorpalstacks.
package naming

import (
	"path/filepath"
	"strings"
)

// SanitizePathComponent replaces characters that are unsafe in filesystem paths.
//
// Only alphanumeric characters, hyphens, underscores, and dots are preserved.
// All other characters are replaced with underscores. If the result is empty,
// "unnamed" is returned.
func SanitizePathComponent(name string) string {
	safe := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, name)
	if safe == "" {
		safe = "unnamed"
	}
	return safe
}

// ValidatePathWithinDir checks that joining baseDir with relPath does not
// escape baseDir via path traversal (e.g. "../").
//
// Returns the cleaned absolute path, or an error if the result escapes baseDir.
func ValidatePathWithinDir(baseDir, relPath string) (string, error) {
	cleanBase := filepath.Clean(baseDir)
	target := filepath.Clean(filepath.Join(cleanBase, relPath))

	if !strings.HasPrefix(target, cleanBase+string(filepath.Separator)) && target != cleanBase {
		return "", &PathTraversalError{Base: cleanBase, Path: relPath}
	}
	return target, nil
}

// PathTraversalError indicates a path traversal attempt was detected.
type PathTraversalError struct {
	Base string
	Path string
}

// Error returns a human-readable description of the path traversal error.
func (e *PathTraversalError) Error() string {
	return "path traversal: " + e.Path + " escapes " + e.Base
}
