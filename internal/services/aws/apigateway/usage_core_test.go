package apigateway

import (
	"strings"
	"testing"

	"vorpalstacks/internal/store/aws/apigateway"
)

// TestCreateApiKeyCoreValidationOrder pins the data-plane failure
// precedence: the stageKey format check runs before the name bound and
// the generateDistinctId/value pairing check, and an empty name is
// accepted by the validation layer (name is optional in the API Gateway
// model, bounded at the documented 1024 characters when supplied).
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
		{
			name: "name above the length bound is rejected",
			in: ApiKeyInput{
				Name: strings.Repeat("n", apigateway.ApiKeyNameMaxLength+1),
			},
			wantMsg: "name must not exceed 1024 characters",
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

// TestUpdateApiKeyNameBound pins the /name patch surface against the
// documented name bound: a replacement at exactly the bound is stored and
// one character more is a BadRequest.
func TestUpdateApiKeyNameBound(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	stores := newImportStores(t)
	created, err := stores.usage.CreateApiKey(&apigateway.ApiKey{Name: "original"})
	if err != nil {
		t.Fatalf("seed key: %v", err)
	}

	atBound := strings.Repeat("n", apigateway.ApiKeyNameMaxLength)
	updated, err := svc.updateApiKeyCore(stores, created.Id, []PatchOperation{
		{Op: "replace", Path: "/name", Value: atBound},
	})
	if err != nil {
		t.Fatalf("replace at the bound: %v", err)
	}
	if updated.Name != atBound {
		t.Fatalf("stored name length = %d, want %d", len(updated.Name), len(atBound))
	}

	_, err = svc.updateApiKeyCore(stores, created.Id, []PatchOperation{
		{Op: "replace", Path: "/name", Value: atBound + "x"},
	})
	apiErr, ok := err.(*ApiGatewayError)
	if !ok || apiErr.Code != "BadRequestException" {
		t.Fatalf("expected BadRequestException above the bound, got %v", err)
	}
}
