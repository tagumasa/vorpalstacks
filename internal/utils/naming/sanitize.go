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

// ValidatePathWithinDir checks that the path stays within baseDir via
// path traversal (e.g. "../").
//
// baseDir is resolved to its absolute form first, so the containment
// comparison is independent of the process working directory. relPath
// may itself be absolute — the form path-persisting writers record —
// and is then validated in place rather than re-rooted under baseDir:
// a join would absorb the leading separator and silently name a
// different file. Returns the cleaned absolute path, or an error if
// the result escapes baseDir.
func ValidatePathWithinDir(baseDir, relPath string) (string, error) {
	base := filepath.Clean(baseDir)
	if abs, err := filepath.Abs(base); err == nil {
		base = abs
	}
	target := filepath.Clean(relPath)
	if !filepath.IsAbs(target) {
		target = filepath.Join(base, target)
	}

	if target != base && !strings.HasPrefix(target, base+string(filepath.Separator)) {
		return "", &PathTraversalError{Base: base, Path: relPath}
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
