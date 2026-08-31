package route53

import (
	"strings"
	"testing"
)

// TestCreateHostedZoneCoreRequiresCallerReference pins that the required
// CallerReference member of CreateHostedZoneRequest is rejected rather than
// silently replaced by a synthesised idempotency token: a server-side token
// changes on every retry and defeats the execute-once semantics the member
// provides. Both checks run before any store access, so a nil store proves
// the rejection precedes persistence.
func TestCreateHostedZoneCoreRequiresCallerReference(t *testing.T) {
	s := &Route53Service{}

	_, err := s.createHostedZoneCore(nil, CreateHostedZoneInput{
		Name:            "emptyref.example.com.",
		CallerReference: "",
	})
	if err == nil {
		t.Fatal("empty CallerReference must be rejected")
	}
	if !strings.Contains(err.Error(), "InvalidInput") {
		t.Fatalf("expected InvalidInput, got %v", err)
	}

	_, err = s.createHostedZoneCore(nil, CreateHostedZoneInput{
		Name:            "longref.example.com.",
		CallerReference: strings.Repeat("a", 129),
	})
	if err == nil {
		t.Fatal("129-character CallerReference must be rejected")
	}
	if !strings.Contains(err.Error(), "InvalidInput") {
		t.Fatalf("expected InvalidInput, got %v", err)
	}
}
