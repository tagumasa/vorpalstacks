package s3

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"vorpalstacks/internal/common/auth"
	"vorpalstacks/internal/utils/crypto"
)

// A presigned request signed with different credentials must be rejected
// when signature verification is enabled, and pass with only the
// freshness check when it is disabled.
func TestPresignedURLVerificationGating(t *testing.T) {
	svc := NewS3ServiceWithStorage(nil, "000000000000")
	svc.SetCredentialsProvider(auth.NewStaticCredentialsProvider(
		"AKIA-SERVER", "server-secret", "us-east-1", ""))

	other := auth.NewStaticCredentialsProvider("AKIA-OTHER", "other-secret", "us-east-1", "")
	otherCreds, err := other.GetCredentials()
	if err != nil {
		t.Fatalf("credentials: %v", err)
	}
	raw := crypto.PresignS3URL(http.MethodGet, "http", "localhost:50080", "some-bucket",
		"some-key", "us-east-1", 15*time.Minute, otherCreds.AccessKeyID, otherCreds.SecretAccessKey)

	req, err := http.NewRequest(http.MethodGet, raw, nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	h := &S3Handler{svc: svc, region: "us-east-1", signatureVerification: true}
	if err := h.verifyPresignedURL(req, "some-bucket"); err == nil {
		t.Fatal("want signature mismatch rejection with verification enabled")
	}

	h.signatureVerification = false
	if err := h.verifyPresignedURL(req, "some-bucket"); err != nil {
		t.Fatalf("freshness-only path rejected a fresh URL: %v", err)
	}

	// An expired URL is rejected on both paths.
	stale := u.Query()
	stale.Set("X-Amz-Date", time.Now().UTC().Add(-time.Hour).Format("20060102T150405Z"))
	staleReq := req.Clone(req.Context())
	staleReq.URL.RawQuery = stale.Encode()
	h.signatureVerification = false
	if err := h.verifyPresignedURL(staleReq, "some-bucket"); err == nil {
		t.Fatal("want expiry rejection with verification disabled")
	}
}
