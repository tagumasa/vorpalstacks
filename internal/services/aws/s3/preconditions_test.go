package s3

import (
	"testing"
	"time"

	s3store "vorpalstacks/internal/store/aws/s3"
)

func precondTime(offset time.Duration) *time.Time {
	t := time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC).Add(offset)
	return &t
}

// The evaluation order is RFC 7232 section 6's — the order the S3 API
// reference points to — and yields its two documented pairwise outcomes.
func TestCheckObjectPreconditions(t *testing.T) {
	obj := &s3store.Object{ETag: `"abc123"`, LastModified: time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC)}
	past := precondTime(-time.Hour)  // strictly before LastModified
	future := precondTime(time.Hour) // strictly after LastModified

	cases := []struct {
		name              string
		ifMatch           string
		ifNoneMatch       string
		ifModifiedSince   *time.Time
		ifUnmodifiedSince *time.Time
		want              error
	}{
		{"documented pair: true If-Match with false If-Unmodified-Since passes", `"abc123"`, "", nil, future, nil},
		{"false If-Match fails regardless of If-Unmodified-Since", `"other"`, "", nil, future, ErrPreconditionFailed},
		{"If-Match wildcard passes", "*", "", nil, future, nil},
		{"If-Unmodified-Since alone fails when modified after", "", "", nil, past, ErrPreconditionFailed},
		{"If-Unmodified-Since alone passes when unmodified", "", "", nil, future, nil},
		{"mixed pair: false If-Unmodified-Since wins over matching If-None-Match", "", `"abc123"`, nil, past, ErrPreconditionFailed},
		{"documented pair: false If-None-Match with true If-Modified-Since is 304", "", `"abc123"`, past, nil, ErrNotModified},
		{"matching If-None-Match is 304 even when modified since IMS", "", `"abc123"`, past, future, ErrNotModified},
		{"If-None-Match wildcard is 304", "", "*", nil, nil, ErrNotModified},
		{"differing If-None-Match ignores If-Modified-Since", "", `"other"`, future, nil, nil},
		{"no conditions passes", "", "", nil, nil, nil},
	}
	for _, c := range cases {
		got := checkObjectPreconditions(obj, c.ifMatch, c.ifNoneMatch, c.ifModifiedSince, c.ifUnmodifiedSince)
		if got != c.want {
			t.Fatalf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

// The copy-source preconditions share the RFC 7232 section 6 order with
// every failure mapping to 412 — a copy cannot answer 304.
func TestCheckCopySourcePreconditions(t *testing.T) {
	obj := &s3store.Object{ETag: `"abc123"`, LastModified: time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC)}
	past := precondTime(-time.Hour)
	future := precondTime(time.Hour)

	cases := []struct {
		name              string
		ifMatch           string
		ifNoneMatch       string
		ifModifiedSince   *time.Time
		ifUnmodifiedSince *time.Time
		want              error
	}{
		{"true If-Match suppresses a false If-Unmodified-Since", `"abc123"`, "", nil, future, nil},
		{"false If-Unmodified-Since wins over matching If-None-Match", "", `"abc123"`, nil, past, ErrPreconditionFailed},
		{"matching If-None-Match fails the copy", "", `"abc123"`, nil, nil, ErrPreconditionFailed},
		{"If-None-Match wildcard fails the copy", "", "*", nil, nil, ErrPreconditionFailed},
		{"differing If-None-Match ignores If-Modified-Since", "", `"other"`, future, nil, nil},
		{"unmodified If-Modified-Since alone fails the copy", "", "", future, nil, ErrPreconditionFailed},
	}
	for _, c := range cases {
		got := checkCopySourcePreconditions(obj, c.ifMatch, c.ifNoneMatch, c.ifModifiedSince, c.ifUnmodifiedSince)
		if got != c.want {
			t.Fatalf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
