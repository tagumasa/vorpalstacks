package route53

import (
	"strings"
	"testing"
)

// TestCreateHealthCheckCoreRequiresCallerReference pins the required
// CallerReference member of CreateHealthCheckRequest: an omitted member is
// rejected with InvalidInput instead of being silently replaced by a
// synthesised idempotency token. The check runs before any store access,
// so a nil store proves the rejection precedes persistence.
func TestCreateHealthCheckCoreRequiresCallerReference(t *testing.T) {
	s := &Route53Service{}
	_, err := s.createHealthCheckCore(nil, CreateHealthCheckInput{
		CallerReference: "",
	})
	if err == nil {
		t.Fatal("empty CallerReference must be rejected")
	}
	if !strings.Contains(err.Error(), "InvalidInput") {
		t.Fatalf("expected InvalidInput, got %v", err)
	}
}
