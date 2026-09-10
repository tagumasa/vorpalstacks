package auth

import (
	"context"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// TestAuthenticateRecordsUsageDespiteCancelledContext pins the detachment
// of usage recording from the request context: a context already cancelled
// before authentication still leaves the day's usage record written — the
// recording runs on its own bounded context, not the request's.
func TestAuthenticateRecordsUsageDespiteCancelledContext(t *testing.T) {
	ps, err := storage.NewPebbleStorage(&storage.Config{
		Path:       t.TempDir(),
		TTLEnabled: false,
	})
	if err != nil {
		t.Fatalf("NewPebbleStorage: %v", err)
	}
	t.Cleanup(func() { ps.Close() })

	store := apigatewaystore.NewUsageStore(ps, "123456789012", "us-east-1")
	if _, err := store.CreateApiKey(&apigatewaystore.ApiKey{
		Id:        "key1",
		Value:     "value-aaaaaaaaaaaaaaaaaaaa",
		Name:      "usage-detachment",
		Enabled:   true,
		StageKeys: []string{"api1/prod"},
	}); err != nil {
		t.Fatalf("CreateApiKey: %v", err)
	}
	if _, err := store.CreateUsagePlan(&apigatewaystore.UsagePlan{Id: "plan1", Name: "detachment"}); err != nil {
		t.Fatalf("CreateUsagePlan: %v", err)
	}
	if _, err := store.CreateUsagePlanKey("plan1", &apigatewaystore.UsagePlanKey{Id: "key1"}); err != nil {
		t.Fatalf("CreateUsagePlanKey: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	authenticator := NewAPIKeyAuthenticator(store)
	if err := authenticator.Authenticate(ctx, "value-aaaaaaaaaaaaaaaaaaaa",
		&apigatewaystore.Method{ApiKeyRequired: true}, "api1", "prod"); err != nil {
		t.Fatalf("Authenticate: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for {
		record, err := store.GetUsage("plan1", "key1", time.Now().Format("2006-01-02"))
		if err == nil && record.RequestCount >= 1 {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("usage record not written for a cancelled request context")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
