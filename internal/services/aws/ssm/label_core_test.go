package ssm

import (
	"errors"
	"testing"
)

// TestUnlabelParameterVersionCore_RequiresVersion pins the Smithy required
// contract of UnlabelParameterVersion.ParameterVersion ("If it isn't
// present, the call will fail."): an omitted version is rejected before any
// store access instead of silently defaulting to the latest version. The
// omission path is unreachable through the AWS SDK (client-side required
// validation), so it is pinned here. LabelParameterVersion keeps the
// latest-version default its model documentation mandates.
func TestUnlabelParameterVersionCore_RequiresVersion(t *testing.T) {
	s := &SSMService{}
	_, err := s.unlabelParameterVersionCore(nil, UnlabelParameterVersionInput{
		Name:   "/p",
		Labels: []string{"alpha"},
	})
	if !errors.Is(err, ErrValidationException) {
		t.Fatalf("err = %v, want ErrValidationException", err)
	}
}
