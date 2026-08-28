package neptunegraph

import (
	"strings"
	"testing"
)

// TestResetDeleteRequireSkipSnapshot pins the server-side enforcement of the
// model-required skipSnapshot member on ResetGraph and DeleteGraph. The typed
// AWS SDK rejects an omitted skipSnapshot client-side, so the server path is
// only reachable through raw HTTP and is pinned here instead.
func TestResetDeleteRequireSkipSnapshot(t *testing.T) {
	s := NewNeptuneGraphService("123456789012", "us-east-1", t.TempDir())

	_, err := s.resetGraphCore(nil, &ResetGraphInput{GraphIdentifier: "g-test"})
	if err == nil || !strings.Contains(err.Error(), "skipSnapshot") {
		t.Fatalf("expected ResetGraph without skipSnapshot to be rejected, got %v", err)
	}

	_, err = s.deleteGraphCore(nil, &DeleteGraphInput{GraphIdentifier: "g-test", Region: "us-east-1"})
	if err == nil || !strings.Contains(err.Error(), "skipSnapshot") {
		t.Fatalf("expected DeleteGraph without skipSnapshot to be rejected, got %v", err)
	}
}
