package cognitoidentityprovider

import (
	"testing"

	"vorpalstacks/internal/common/request"
)

// An explicitly supplied MinimumLength of zero is an out-of-range value,
// not the "unset" marker the stored-policy check tolerates.
func TestParsePasswordPolicyExplicitMinimumLengthZeroRejected(t *testing.T) {
	flat := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies.PasswordPolicy.MinimumLength": "0",
	}}
	if _, err := parsePasswordPolicyWithBase(flat, nil); err == nil {
		t.Fatal("expected an explicit MinimumLength of 0 to be rejected on the flat parameter path")
	}

	jsonPath := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies": map[string]interface{}{
			"PasswordPolicy": map[string]interface{}{
				"MinimumLength": 0,
			},
		},
	}}
	if _, err := parsePasswordPolicyWithBase(jsonPath, nil); err == nil {
		t.Fatal("expected an explicit MinimumLength of 0 to be rejected on the structured parameter path")
	}
}

// The Smithy range boundaries {6, 99} apply to explicitly supplied values;
// an absent member keeps the stored or default value untouched.
func TestParsePasswordPolicyMinimumLengthBoundaries(t *testing.T) {
	for _, tc := range []struct {
		value   string
		wantErr bool
	}{
		{"6", false},
		{"8", false},
		{"99", false},
		{"100", true},
	} {
		req := &request.ParsedRequest{Parameters: map[string]interface{}{
			"Policies.PasswordPolicy.MinimumLength": tc.value,
		}}
		policy, err := parsePasswordPolicyWithBase(req, nil)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("expected MinimumLength %s to be rejected", tc.value)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for MinimumLength %s: %v", tc.value, err)
		}
		if policy.MinimumLength != 8 && tc.value == "8" {
			t.Fatalf("expected the parsed MinimumLength to be 8, got %d", policy.MinimumLength)
		}
	}
}

// An absent MinimumLength member leaves the base policy value in place.
func TestParsePasswordPolicyAbsentMinimumLengthKeepsBase(t *testing.T) {
	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies.PasswordPolicy.RequireSymbols": "false",
	}}
	policy, err := parsePasswordPolicyWithBase(req, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy.MinimumLength != 8 {
		t.Fatalf("expected the default MinimumLength of 8, got %d", policy.MinimumLength)
	}
}
