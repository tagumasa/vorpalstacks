package s3

import (
	"strings"
	"testing"
)

// Restore-index keys must stay disjoint from object record keys (which
// start with a bucket name) and must round-trip through the parser.
func TestRestoreIndexKey(t *testing.T) {
	key := restoreIndexKey("bucket", "obj.txt", "")
	if key != restoreIndexKey("bucket", "obj.txt", "null") {
		t.Fatalf("empty versionId must normalise to null: %q", key)
	}
	if !strings.HasPrefix(key, "\x01restore\x00") {
		t.Fatalf("index key must carry the control-char marker prefix: %q", key)
	}
	if strings.HasPrefix(key, "bucket\x00") {
		t.Fatalf("index key must not collide with object record keys: %q", key)
	}

	entry, ok := parseRestoreIndexKey(key)
	if !ok || entry.Bucket != "bucket" || entry.Key != "obj.txt" || entry.VersionID != "null" {
		t.Fatalf("round-trip failed: %+v ok=%v", entry, ok)
	}

	entry, ok = parseRestoreIndexKey(restoreIndexKey("b", "k", "v-1"))
	if !ok || entry.VersionID != "v-1" {
		t.Fatalf("versioned round-trip failed: %+v ok=%v", entry, ok)
	}

	if _, ok := parseRestoreIndexKey("bucket\x00obj.txt\x00null"); ok {
		t.Fatal("an object record key must not parse as an index entry")
	}
	if _, ok := parseRestoreIndexKey(restoreIndexPrefix + "onlyonefield"); ok {
		t.Fatal("a malformed index key must not parse")
	}
}
