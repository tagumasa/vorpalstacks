package iot

import (
	"errors"
	"testing"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// validateCertProviderOperations enforces the fixed single-element
// CreateCertificateFromCsr list before any store access.
func TestValidateCertProviderOperations(t *testing.T) {
	valid := []interface{}{"CreateCertificateFromCsr"}
	if err := validateCertProviderOperations(valid); err != nil {
		t.Fatalf("valid operations list rejected: %v", err)
	}

	invalid := []struct {
		name string
		ops  interface{}
	}{
		{"omitted", nil},
		{"empty", []interface{}{}},
		{"two entries", []interface{}{"CreateCertificateFromCsr", "CreateCertificateFromCsr"}},
		{"wrong value", []interface{}{"SomethingElse"}},
		{"non-string entry", []interface{}{42}},
		{"not a list", "CreateCertificateFromCsr"},
	}
	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateCertProviderOperations(tt.ops); !errors.Is(err, iotstore.ErrInvalidRequest) {
				t.Fatalf("expected ErrInvalidRequest, got %v", err)
			}
		})
	}
}
