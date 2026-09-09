package s3

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// The read-error mapping surfaces a delete-marker hit the documented way:
// an explicitly requested marker version is the 405 surface, a marker that
// is the latest version stays a NoSuchKey 404, and both carry the marker
// record for their response headers; every other miss keeps the addressing
// rules of mapVersionLookupError.
func TestMapVersionReadError(t *testing.T) {
	marker := &s3store.Object{
		Key:            "k",
		BucketName:     "b",
		VersionID:      "mv-1",
		IsDeleteMarker: true,
		LastModified:   time.Date(2026, 9, 8, 10, 0, 0, 0, time.UTC),
	}

	err := mapVersionReadError(s3store.ErrObjectNotFound, marker, "k", "mv-1")
	var dmErr *deleteMarkerReadError
	if !errors.As(err, &dmErr) {
		t.Fatalf("explicit marker version: got %T, want deleteMarkerReadError", err)
	}
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.Code != "MethodNotAllowed" || awsErr.HTTPStatus != http.StatusMethodNotAllowed {
		t.Fatalf("explicit marker version error = %v, want MethodNotAllowed 405", err)
	}
	if dmErr.marker != marker {
		t.Fatal("the wrapper must carry the marker record")
	}

	err = mapVersionReadError(s3store.ErrObjectNotFound, marker, "k", "")
	if !errors.As(err, &dmErr) {
		t.Fatalf("latest marker: got %T, want deleteMarkerReadError", err)
	}
	if !errors.As(err, &awsErr) || awsErr.Code != "NoSuchKey" || awsErr.HTTPStatus != http.StatusNotFound {
		t.Fatalf("latest marker error = %v, want NoSuchKey 404", err)
	}

	if err := mapVersionReadError(s3store.ErrObjectNotFound, nil, "k", "some-vid"); !errors.Is(err, ErrNoSuchVersion) {
		t.Fatalf("plain version miss = %v, want NoSuchVersion", err)
	}
	other := errors.New("store fault")
	if got := mapVersionReadError(other, nil, "k", "some-vid"); got != other {
		t.Fatalf("non-not-found errors propagate: got %v", got)
	}
	if got := mapVersionReadError(s3store.ErrObjectNotFound, nil, "k", ""); !errors.Is(got, s3store.ErrObjectNotFound) {
		t.Fatalf("plain key miss = %v, want the store sentinel to reach the writer", got)
	}
}

// writeError renders the documented delete-marker read contracts: the 405
// carries Last-Modified plus the marker identification headers, the
// latest-marker 404 carries the identification headers only, and both
// bodies name their error code.
func TestWriteErrorDeleteMarkerHeaders(t *testing.T) {
	stamp := time.Date(2026, 9, 8, 10, 0, 0, 0, time.UTC)
	marker := &s3store.Object{
		Key:            "k",
		BucketName:     "b",
		VersionID:      "mv-1",
		IsDeleteMarker: true,
		LastModified:   stamp,
	}

	rec := httptest.NewRecorder()
	(&S3Handler{}).writeError(rec, mapVersionReadError(s3store.ErrObjectNotFound, marker, "k", "mv-1"), "b", "k", "req-1")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
	if got := rec.Header().Get("Last-Modified"); got != stamp.UTC().Format(http.TimeFormat) {
		t.Fatalf("Last-Modified = %q, want the marker's timestamp", got)
	}
	if got := rec.Header().Get("x-amz-delete-marker"); got != "true" {
		t.Fatalf("x-amz-delete-marker = %q, want true", got)
	}
	if got := rec.Header().Get("x-amz-version-id"); got != "mv-1" {
		t.Fatalf("x-amz-version-id = %q, want mv-1", got)
	}
	if !strings.Contains(rec.Body.String(), "MethodNotAllowed") {
		t.Fatalf("body = %q, want MethodNotAllowed", rec.Body.String())
	}

	rec = httptest.NewRecorder()
	(&S3Handler{}).writeError(rec, mapVersionReadError(s3store.ErrObjectNotFound, marker, "k", ""), "b", "k", "req-2")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("latest-marker status = %d, want 404", rec.Code)
	}
	if got := rec.Header().Get("Last-Modified"); got != "" {
		t.Fatalf("latest-marker 404 must not carry Last-Modified, got %q", got)
	}
	if got := rec.Header().Get("x-amz-delete-marker"); got != "true" {
		t.Fatalf("latest-marker x-amz-delete-marker = %q, want true", got)
	}
	if got := rec.Header().Get("x-amz-version-id"); got != "mv-1" {
		t.Fatalf("latest-marker x-amz-version-id = %q, want mv-1", got)
	}
	if !strings.Contains(rec.Body.String(), "NoSuchKey") {
		t.Fatalf("body = %q, want NoSuchKey", rec.Body.String())
	}
}
