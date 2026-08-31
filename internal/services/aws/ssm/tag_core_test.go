package ssm

import (
	"errors"
	"testing"

	"vorpalstacks/internal/common/request"
)

// TestValidateResourceType_Required pins the Smithy required contract of
// ResourceType on the three tag operations: a missing member is rejected
// with ValidationException (AWS rejects null required members) and an
// unimplemented value with InvalidResourceType. The omission path is
// unreachable through the AWS SDK (client-side required validation), so it
// is pinned here.
func TestValidateResourceType_Required(t *testing.T) {
	req := &request.ParsedRequest{Parameters: map[string]interface{}{}}
	if err := validateResourceType(req); !errors.Is(err, ErrValidationException) {
		t.Fatalf("missing err = %v, want ErrValidationException", err)
	}
	req.Parameters["ResourceType"] = "Document"
	if err := validateResourceType(req); !errors.Is(err, ErrInvalidResourceType) {
		t.Fatalf("Document err = %v, want ErrInvalidResourceType", err)
	}
	req.Parameters["ResourceType"] = "Parameter"
	if err := validateResourceType(req); err != nil {
		t.Fatalf("Parameter err = %v, want nil", err)
	}
}
