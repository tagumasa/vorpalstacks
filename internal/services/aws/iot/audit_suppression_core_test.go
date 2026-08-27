package iot

import (
	"errors"
	"testing"

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
