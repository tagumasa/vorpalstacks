package iot

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// TestDimensionTokenCrossNameUniqueness pins the documented
// clientRequestToken uniqueness: reusing an existing dimension's token for
// a different dimension is rejected, the same-name replay and
// duplicate-conflict rules are unchanged, and deleting the owning
// dimension releases the token for later creates.
func TestDimensionTokenCrossNameUniqueness(t *testing.T) {
	rawStore, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { rawStore.Close() })
	store := iotstore.NewIotStore(rawStore, "000000000000", "us-east-1", nil)
	svc := &IoTService{}
	mkInput := func(name, token string) CreateDimensionInput {
		return CreateDimensionInput{
			Name:               name,
			Type:               "TOPIC_FILTER",
			StringValues:       []string{"sdk/test/#"},
			ClientRequestToken: token,
		}
	}
	if _, err := svc.createDimensionCore(store, mkInput("dim-alpha", "token-shared")); err != nil {
		t.Fatalf("initial create: %v", err)
	}
	if _, err := svc.createDimensionCore(store, mkInput("dim-alpha", "token-shared")); err != nil {
		t.Fatalf("same-name same-token replay must succeed: %v", err)
	}
	if _, err := svc.createDimensionCore(store, mkInput("dim-alpha", "token-other")); !errors.Is(err, iotstore.ErrResourceAlreadyExists) {
		t.Fatalf("same-name different-token: got %v, want ErrResourceAlreadyExists", err)
	}
	if _, err := svc.createDimensionCore(store, mkInput("dim-beta", "token-shared")); !errors.Is(err, iotstore.ErrResourceAlreadyExists) {
		t.Fatalf("cross-name token reuse: got %v, want ErrResourceAlreadyExists", err)
	}
	// Deleting the owning dimension releases the token.
	if err := svc.deleteDimensionCore(store, "dim-alpha"); err != nil {
		t.Fatalf("delete owner: %v", err)
	}
	if _, err := svc.createDimensionCore(store, mkInput("dim-beta", "token-shared")); err != nil {
		t.Fatalf("token reuse after owner deletion: %v", err)
	}
}
