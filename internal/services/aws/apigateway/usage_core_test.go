package apigateway

import (
	"strings"
	"testing"
)

// TestCreateApiKeyCoreValidationOrder pins the data-plane failure
// precedence: the stageKey format check runs before the
// generateDistinctId/value pairing check, and an empty name is accepted
// by the validation layer (name is optional in the API Gateway model).
func TestCreateApiKeyCoreValidationOrder(t *testing.T) {
	svc := &APIGatewayService{}
	generateFalse := false

	cases := []struct {
		name    string
		in      ApiKeyInput
		wantMsg string
	}{
		{
			name: "invalid stageKey wins over the value pairing check",
			in: ApiKeyInput{
				Name:               "",
				StageKeys:          []string{"bad-stage-key"},
				GenerateDistinctId: &generateFalse,
				Value:              "",
			},
			wantMsg: "invalid stageKey format, expected restApiId/stageName: bad-stage-key",
		},
		{
			name: "value pairing check fires without stageKeys",
			in: ApiKeyInput{
				GenerateDistinctId: &generateFalse,
				Value:              "",
			},
			wantMsg: "value is required when generateDistinctId is false",
		},
		{
			name: "supplied value below the length floor is rejected",
			in: ApiKeyInput{
				GenerateDistinctId: &generateFalse,
				Value:              "short-value",
			},
			wantMsg: "value must be between 20 and 128 characters",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := tc.in
			// The validations under test return before any store access,
			// so a nil stores bundle is safe for these error paths.
			_, err := svc.createApiKeyCore(nil, &in)
			if err == nil {
				t.Fatal("expected a validation error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantMsg) {
				t.Fatalf("error %q does not contain %q", err.Error(), tc.wantMsg)
			}
		})
	}
}
