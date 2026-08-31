package ssm

import (
	"errors"
	"testing"
)

// TestDeleteParametersCore_EmptyNamesRejected pins that DeleteParameters
// shares the ParameterNameList required contract: an empty Names list is
// rejected with ValidationException before any store access. A missing
// member is unreachable through the AWS SDK (client-side required
// validation), so the omission shape is pinned here as the nil case.
func TestDeleteParametersCore_EmptyNamesRejected(t *testing.T) {
	s := &SSMService{}
	_, err := s.deleteParametersCore(nil, DeleteParametersInput{Names: nil})
	if !errors.Is(err, ErrValidationException) {
		t.Fatalf("nil err = %v, want ErrValidationException", err)
	}
	_, err = s.deleteParametersCore(nil, DeleteParametersInput{Names: []string{}})
	if !errors.Is(err, ErrValidationException) {
		t.Fatalf("empty err = %v, want ErrValidationException", err)
	}
}
