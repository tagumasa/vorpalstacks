package sns

// Pins for the signing-key initialiser's retry contract: a failed attempt
// leaves the fields unset so the next call retries — a transient storage
// fault at first use must not strip signatures from every later envelope
// for the process lifetime.

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// TestSigningKeyInitRetriesAfterFailure pins the retry: the initialiser
// runs while the region storage cannot open (its directory path is
// occupied by a regular file), leaves both fields unset, and the SAME
// service object completes the initialisation once the obstruction clears
// — the envelope then carries a signature and a certificate URL.
func TestSigningKeyInitRetriesAfterFailure(t *testing.T) {
	base := t.TempDir()
	blocker := filepath.Join(base, "us-east-1")
	if err := os.WriteFile(blocker, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("write blocker: %v", err)
	}

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: base})
	if err != nil {
		t.Fatalf("region storage manager: %v", err)
	}
	t.Cleanup(func() { mgr.Close() })
	svc := NewSNSService(mgr, "123456789012", "us-east-1")

	svc.initSigningKey()
	if svc.signingKey != nil || svc.signingCertPEM != nil {
		t.Fatal("failed initialisation set the signing fields — the retry path is gone")
	}

	if err := os.Remove(blocker); err != nil {
		t.Fatalf("remove blocker: %v", err)
	}
	svc.initSigningKey()
	if svc.signingKey == nil || svc.signingCertPEM == nil {
		t.Fatal("initialisation did not complete after the obstruction cleared")
	}

	envelope := map[string]interface{}{
		"Type":      "Notification",
		"MessageId": "68f0dbd4-09a5-4d5f-99a0-21a6ba41e8c3",
		"TopicArn":  "arn:aws:sns:us-east-1:123456789012:signing-pin-topic",
		"Message":   "the signing key recovered",
		"Timestamp": envelopeTimestamp(time.Now()),
	}
	signature, certURL := svc.signEnvelope(envelope, "us-east-1", "1")
	if signature == "" || certURL == "" {
		t.Fatalf("recovered key produced signature=%q certURL=%q, want both", signature, certURL)
	}
}

// TestSigningKeyInitStableOnSuccess pins that a completed initialisation
// keeps its key: the second call must not regenerate (a regenerated key
// invalidates the signatures every already-delivered envelope carries).
func TestSigningKeyInitStableOnSuccess(t *testing.T) {
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("region storage manager: %v", err)
	}
	t.Cleanup(func() { mgr.Close() })
	svc := NewSNSService(mgr, "123456789012", "us-east-1")

	svc.initSigningKey()
	first := svc.signingKey
	if first == nil {
		t.Fatal("initialisation did not complete")
	}
	svc.initSigningKey()
	if svc.signingKey != first {
		t.Fatal("second initialisation regenerated the key — completed state must be stable")
	}
}
