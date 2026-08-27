package iot

import "testing"

// TestValidateRegistrationStatus pins the CertificateStatus enum enforced
// on the certificate registration operations.
func TestValidateRegistrationStatus(t *testing.T) {
	for _, status := range []string{"ACTIVE", "INACTIVE", "REVOKED", "PENDING_TRANSFER", "REGISTER_INACTIVE", "PENDING_ACTIVATION"} {
		if err := validateRegistrationStatus(status); err != nil {
			t.Fatalf("validateRegistrationStatus rejected documented status %q", status)
		}
	}
	for _, status := range []string{"", "active", "GARBAGE", "PENDING"} {
		if err := validateRegistrationStatus(status); err == nil {
			t.Fatalf("validateRegistrationStatus accepted invalid status %q", status)
		}
	}
}
