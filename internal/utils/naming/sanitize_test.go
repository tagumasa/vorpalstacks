package naming

import (
	"path/filepath"
	"runtime"
	"testing"
)

// TestValidatePathWithinDirAbsoluteForms pins the two legitimate input
// forms: a relative path joins under the directory, and an absolute path
// — the form path-persisting writers record — validates in place instead
// of being silently re-rooted by the join (which would name a different
// file). The containment comparison must be independent of a relative
// base directory, because a server started with a relative data path
// validates absolute stored paths against it.
func TestValidatePathWithinDirAbsoluteForms(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("path forms are POSIX-shaped")
	}

	base := "rel-dir/nested"
	absBase, err := filepath.Abs(base)
	if err != nil {
		t.Fatal(err)
	}

	// A relative input resolves under the absolute base.
	p, err := ValidatePathWithinDir(base, "file.vlog")
	if err != nil {
		t.Fatalf("relative input: %v", err)
	}
	if want := absBase + "/file.vlog"; p != want {
		t.Fatalf("relative input resolved to %q, want %q", p, want)
	}

	// An absolute input inside the base is returned as named, not
	// re-rooted (the join would produce rel-dir/nested/abs-base/...).
	inside := absBase + "/chunk.vlog"
	p, err = ValidatePathWithinDir(base, inside)
	if err != nil {
		t.Fatalf("absolute in-dir input: %v", err)
	}
	if p != inside {
		t.Fatalf("absolute in-dir input resolved to %q, want %q", p, inside)
	}

	// An absolute input outside the base is rejected.
	if _, err := ValidatePathWithinDir(base, "/etc/passwd"); err == nil {
		t.Fatal("absolute outside-dir input must be rejected")
	}

	// A relative escape is rejected on the absolute base.
	if _, err := ValidatePathWithinDir(base, "../escape"); err == nil {
		t.Fatal("relative escape must be rejected")
	}

	// The base itself is admissible.
	if p, err = ValidatePathWithinDir(base, "."); err != nil || p != absBase {
		t.Fatalf("base itself: p=%q err=%v, want %q", p, err, absBase)
	}
}
