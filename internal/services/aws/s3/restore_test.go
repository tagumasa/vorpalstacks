package s3

import (
	"strings"
	"testing"
	"time"

	s3store "vorpalstacks/internal/store/aws/s3"
)

// nextRestoreExpiry adds the requested days to the restore time and rounds
// up to the following midnight UTC, as the restore API documents (a copy
// restored on Oct 15 2012 at 10:30 UTC for 3 days expires on Oct 19 2012
// at 00:00 UTC).
func TestNextRestoreExpiry(t *testing.T) {
	restored := time.Date(2012, 10, 15, 10, 30, 0, 0, time.UTC)
	got := nextRestoreExpiry(restored, 3)
	want := time.Date(2012, 10, 19, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("expiry = %s, want %s", got, want)
	}

	// A restore time whose end lands exactly on midnight keeps that
	// midnight instead of rounding into the next day.
	midnight := time.Date(2012, 10, 18, 0, 0, 0, 0, time.UTC)
	got = nextRestoreExpiry(midnight, 3)
	want = time.Date(2012, 10, 21, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("midnight expiry = %s, want %s", got, want)
	}
}

// objectRestored treats the expiry timestamp as the sole restore state and
// restoreHeaderValue renders the x-amz-restore response for active copies.
func TestObjectRestoredAndHeader(t *testing.T) {
	now := time.Now()
	obj := &s3store.Object{StorageClass: s3store.StorageClassGlacier}
	if objectRestored(obj, now) {
		t.Fatal("an object without restore expiry must not count as restored")
	}
	if h := restoreHeaderValue(obj, now); h != "" {
		t.Fatalf("unrestored object must not report a restore header, got %q", h)
	}

	expiry := now.Add(time.Hour)
	obj.RestoreExpiry = &expiry
	if !objectRestored(obj, now) {
		t.Fatal("an unexpired restore must count as restored")
	}
	h := restoreHeaderValue(obj, now)
	if !strings.Contains(h, `ongoing-request="false"`) || !strings.Contains(h, "expiry-date=") {
		t.Fatalf("restore header must carry ongoing-request and expiry-date, got %q", h)
	}

	if objectRestored(obj, now.Add(2*time.Hour)) {
		t.Fatal("an expired restore must not count as restored")
	}
	if h := restoreHeaderValue(obj, now.Add(2*time.Hour)); h != "" {
		t.Fatalf("expired restore must not report a header, got %q", h)
	}
}
