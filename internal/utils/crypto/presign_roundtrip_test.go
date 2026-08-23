package crypto

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/auth"
)

// A URL produced by PresignS3URL must pass the verifier's freshness and
// signature checks when the same credentials verify it.
func TestPresignS3URLRoundTrip(t *testing.T) {
	provider := auth.NewStaticCredentialsProvider("AKIA-TEST", "secret-test", "us-east-1", "")
	creds, err := provider.GetCredentials()
	if err != nil {
		t.Fatalf("credentials: %v", err)
	}

	raw := PresignS3URL(http.MethodPut, "http", "localhost:50080", "cognito-import",
		"us-east-1_POOL/import-job1", "us-east-1", 15*time.Minute, creds.AccessKeyID, creds.SecretAccessKey)

	if !strings.HasPrefix(raw, "http://localhost:50080/cognito-import/") {
		t.Fatalf("url = %q", raw)
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if exp := u.Query().Get("X-Amz-Expires"); exp != "900" {
		t.Fatalf("X-Amz-Expires = %q, want 900", exp)
	}

	// Freshness check passes now and fails after expiry.
	if err := CheckPresignedURLFreshness(u.Query()); err != nil {
		t.Fatalf("freshness: %v", err)
	}
	stale := u.Query()
	stale.Set("X-Amz-Date", time.Now().UTC().Add(-time.Hour).Format("20060102T150405Z"))
	if err := CheckPresignedURLFreshness(stale); err == nil {
		t.Fatal("want expiry error for a stale date")
	}

	// Full verification against the matching credentials and platform
	// path-style request.
	verifier := NewPresignedURLVerifier(provider)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, raw, nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if err := verifier.VerifyPresignedURL(req, "cognito-import", "us-east-1"); err != nil {
		t.Fatalf("verification: %v", err)
	}
}
