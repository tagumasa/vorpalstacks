package storage

import (
	"path/filepath"
	"strings"
	"testing"
)

// sanitizeKey must be traversal-safe and injective: no mapping contains a
// parent-directory step or an absolute escape, and distinct keys map to
// distinct paths — legal AWS keys such as "a/../b" and "b", or "a//b" and
// "a/b", must never share a blob file.
func TestSanitizeKeyTraversalSafeAndInjective(t *testing.T) {
	s := &HybridBlobStore{}
	keys := []string{
		"docs/year/file.txt",
		"file..backup.txt",
		"a..b",
		"..",
		"../escape",
		"a/../b",
		"b",
		"a//b",
		"a",
		"a/b",
		"a/b/c",
		// Keys that spell marker names themselves must not alias with
		// the marker-prefixed layout of other keys.
		"f:a",
		"f:a/b",
		"d:a",
		"d:a/b",
		"f:",
		"d:",
		"a/",
		"/leading",
		".",
		"a/./b",
		"100%sure",
		"%2E%2E",
		"..%2F..",
	}

	files := map[string]string{} // mapped file path -> key
	dirs := map[string]string{}  // implied directory path -> first key needing it
	for _, key := range keys {
		mapped := s.sanitizeKey(key)

		if strings.Contains(mapped, ".."+string(filepath.Separator)) || mapped == ".." {
			t.Fatalf("key %q maps to %q — parent traversal survives", key, mapped)
		}
		if filepath.IsAbs(mapped) {
			t.Fatalf("key %q maps to absolute path %q", key, mapped)
		}
		for _, seg := range strings.Split(mapped, string(filepath.Separator)) {
			if seg == "." || seg == ".." || seg == "" {
				t.Fatalf("key %q maps to %q — segment %q is path-meaningful", key, mapped, seg)
			}
		}

		if prev, dup := files[mapped]; dup {
			t.Fatalf("keys %q and %q alias to %q", prev, key, mapped)
		}
		files[mapped] = key
		for dir := filepath.Dir(mapped); dir != "."; dir = filepath.Dir(dir) {
			if _, dup := dirs[dir]; !dup {
				dirs[dir] = key
			}
		}
	}

	// The S3 key namespace is flat: a key and its prefix extension are
	// independent objects, so no key's blob file may sit on a path where
	// another key's directories must live.
	for path, fileKey := range files {
		if dirKey, dup := dirs[path]; dup {
			t.Fatalf("path %q is %q's blob file but also a directory %q's keys require", path, fileKey, dirKey)
		}
	}

	// The role markers place files and directories on disjoint names.
	if got, want := s.sanitizeKey("a"), "f:a"; got != want {
		t.Fatalf("single-segment key layout changed: got %q want %q", got, want)
	}
	if got, want := s.sanitizeKey("a/b"), filepath.Join("d:a", "f:b"); got != want {
		t.Fatalf("multi-segment key layout changed: got %q want %q", got, want)
	}
}
