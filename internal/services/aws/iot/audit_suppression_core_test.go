package iot

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// createAuditSuppressionCore enforces the model's required
// clientRequestToken member before any store access; the AWS SDK fills the
// idempotency token automatically, so the missing-member negative is only
// reachable through the raw wire and is pinned here (a nil store stands in
// for the real one; the guard runs first).
func TestCreateAuditSuppressionRequiresToken(t *testing.T) {
	svc := &IoTService{}
	err := svc.createAuditSuppressionCore(nil, CreateAuditSuppressionInput{
		CheckName:                    "DEVICE_CERTIFICATE_EXPIRING_CHECK",
		ResourceIdentifier:           map[string]interface{}{"deviceCertificateId": "cert-1"},
		SuppressIndefinitelyProvided: true,
	})
	if !errors.Is(err, iotstore.ErrMissingParam) {
		t.Fatalf("expected ErrMissingParam for an empty clientRequestToken, got %v", err)
	}
}

// TestAuditSuppressionTokenCrossTupleUniqueness pins the documented
// clientRequestToken uniqueness: reusing an existing suppression's token
// for a different (checkName, resourceIdentifier) tuple is rejected, and
// deleting the owning suppression releases the token for later creates.
func TestAuditSuppressionTokenCrossTupleUniqueness(t *testing.T) {
	rawStore, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { rawStore.Close() })
	store := iotstore.NewIotStore(rawStore, "000000000000", "us-east-1", nil)
	svc := &IoTService{}
	mkInput := func(checkName, token string) CreateAuditSuppressionInput {
		return CreateAuditSuppressionInput{
			CheckName:                    checkName,
			ResourceIdentifier:           map[string]interface{}{"deviceCertificateId": "cert-1"},
			SuppressIndefinitely:         true,
			SuppressIndefinitelyProvided: true,
			ClientRequestToken:           token,
		}
	}
	if err := svc.createAuditSuppressionCore(store, mkInput("CHECK_A", "token-shared")); err != nil {
		t.Fatalf("initial create: %v", err)
	}
	if err := svc.createAuditSuppressionCore(store, mkInput("CHECK_A", "token-shared")); err != nil {
		t.Fatalf("same-tuple same-token replay must succeed: %v", err)
	}
	if err := svc.createAuditSuppressionCore(store, mkInput("CHECK_B", "token-shared")); !errors.Is(err, iotstore.ErrResourceAlreadyExists) {
		t.Fatalf("cross-tuple token reuse: got %v, want ErrResourceAlreadyExists", err)
	}
	// Deleting the owning suppression releases the token.
	if err := svc.deleteAuditSuppressionCore(store, "CHECK_A", map[string]interface{}{"deviceCertificateId": "cert-1"}); err != nil {
		t.Fatalf("delete owner: %v", err)
	}
	if err := svc.createAuditSuppressionCore(store, mkInput("CHECK_B", "token-shared")); err != nil {
		t.Fatalf("token reuse after owner deletion: %v", err)
	}
}

// TestAuditSuppressionExpirationPairings pins the documented expiration
// pairings: an explicit false alongside an expiration date is the pairing
// the CLI example sends and is accepted on create and update, an
// indefinite suppression together with a date is rejected, and a
// suppression without any expiration is rejected on create.
func TestAuditSuppressionExpirationPairings(t *testing.T) {
	rawStore, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { rawStore.Close() })
	store := iotstore.NewIotStore(rawStore, "000000000000", "us-east-1", nil)
	svc := &IoTService{}
	rid := map[string]interface{}{"deviceCertificateId": "cert-pair"}

	// No expiration at all is rejected on create.
	err = svc.createAuditSuppressionCore(store, CreateAuditSuppressionInput{
		CheckName:          "CHECK_NONE",
		ResourceIdentifier: rid,
		ClientRequestToken: "token-none",
	})
	if !errors.Is(err, iotstore.ErrInvalidRequest) {
		t.Fatalf("expiration-less create: got %v, want ErrInvalidRequest", err)
	}

	// false + expirationDate is accepted on create.
	err = svc.createAuditSuppressionCore(store, CreateAuditSuppressionInput{
		CheckName:                    "CHECK_PAIR",
		ResourceIdentifier:           rid,
		ExpirationDate:               1700000000,
		ExpirationProvided:           true,
		SuppressIndefinitely:         false,
		SuppressIndefinitelyProvided: true,
		ClientRequestToken:           "token-pair",
	})
	if err != nil {
		t.Fatalf("false + expirationDate create: %v", err)
	}

	// ... and on update.
	err = svc.updateAuditSuppressionCore(store, UpdateAuditSuppressionInput{
		CheckName:                    "CHECK_PAIR",
		ResourceIdentifier:           rid,
		ExpirationDate:               1800000000,
		ExpirationProvided:           true,
		SuppressIndefinitely:         false,
		SuppressIndefinitelyProvided: true,
	})
	if err != nil {
		t.Fatalf("false + expirationDate update: %v", err)
	}

	// An indefinite suppression together with a date is rejected on update.
	err = svc.updateAuditSuppressionCore(store, UpdateAuditSuppressionInput{
		CheckName:                    "CHECK_PAIR",
		ResourceIdentifier:           rid,
		ExpirationDate:               1900000000,
		ExpirationProvided:           true,
		SuppressIndefinitely:         true,
		SuppressIndefinitelyProvided: true,
	})
	if !errors.Is(err, iotstore.ErrInvalidRequest) {
		t.Fatalf("indefinite + date update: got %v, want ErrInvalidRequest", err)
	}

	// An indefinite-only update replaces the stored date.
	err = svc.updateAuditSuppressionCore(store, UpdateAuditSuppressionInput{
		CheckName:                    "CHECK_PAIR",
		ResourceIdentifier:           rid,
		SuppressIndefinitely:         true,
		SuppressIndefinitelyProvided: true,
	})
	if err != nil {
		t.Fatalf("indefinite-only update: %v", err)
	}
	rec, err := svc.describeAuditSuppressionCore(store, "CHECK_PAIR", rid)
	if err != nil {
		t.Fatalf("describe after indefinite update: %v", err)
	}
	switch v := rec.ExpirationDate.(type) {
	case int64:
		if v != 0 {
			t.Fatalf("indefinite update must clear the stored date, got %v", v)
		}
	case float64:
		if v != 0 {
			t.Fatalf("indefinite update must clear the stored date, got %v", v)
		}
	}
}
