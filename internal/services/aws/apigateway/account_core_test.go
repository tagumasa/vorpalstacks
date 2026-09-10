package apigateway

import (
	"testing"

	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// TestAccountCore pins the account singleton contract: the default view
// carries the UsagePlans feature and the documented default throttle
// limits, the documented patch surface (replace /cloudwatchRoleArn, add
// and remove /features) applies, and everything else — including removing
// UsagePlans and any undocumented path — is a BadRequest.
func TestAccountCore(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	stores := newImportStores(t)

	account, err := svc.getAccountCore(stores)
	if err != nil {
		t.Fatalf("getAccountCore: %v", err)
	}
	if account.ThrottleSettings == nil {
		t.Fatal("default throttleSettings missing")
	}
	if account.ThrottleSettings.RateLimit != apigatewaystore.AccountDefaultRateLimit ||
		account.ThrottleSettings.BurstLimit != apigatewaystore.AccountDefaultBurstLimit {
		t.Fatalf("default throttleSettings = %+v, want the documented account defaults",
			account.ThrottleSettings)
	}
	if !accountHasFeature(account, "UsagePlans") {
		t.Fatalf("default features = %v, want UsagePlans", account.Features)
	}

	updated, err := svc.updateAccountCore(stores, []PatchOperation{
		{Op: "replace", Path: "/cloudwatchRoleArn", Value: "arn:aws:iam::123456789012:role/probe"},
		{Op: "add", Path: "/features", Value: "Probe"},
	})
	if err != nil {
		t.Fatalf("updateAccountCore: %v", err)
	}
	if updated.CloudwatchRoleArn != "arn:aws:iam::123456789012:role/probe" || !accountHasFeature(updated, "Probe") {
		t.Fatalf("update not applied: %+v", updated)
	}

	// The stored record round-trips through a fresh core call.
	reread, err := svc.getAccountCore(stores)
	if err != nil {
		t.Fatalf("reread: %v", err)
	}
	if reread.CloudwatchRoleArn != "arn:aws:iam::123456789012:role/probe" {
		t.Fatalf("cloudwatchRoleArn not persisted: %+v", reread)
	}

	removed, err := svc.updateAccountCore(stores, []PatchOperation{
		{Op: "remove", Path: "/features", Value: "Probe"},
	})
	if err != nil {
		t.Fatalf("feature removal: %v", err)
	}
	if accountHasFeature(removed, "Probe") {
		t.Fatalf("feature not removed: %v", removed.Features)
	}

	for name, ops := range map[string][]PatchOperation{
		"remove UsagePlans":   {{Op: "remove", Path: "/features", Value: "UsagePlans"}},
		"add roleArn":         {{Op: "add", Path: "/cloudwatchRoleArn", Value: "x"}},
		"undocumented path":   {{Op: "replace", Path: "/apiKeyVersion", Value: "2"}},
		"throttle not listed": {{Op: "replace", Path: "/throttleSettings/burstLimit", Value: "10"}},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := svc.updateAccountCore(stores, ops)
			apiErr, ok := err.(*ApiGatewayError)
			if !ok {
				t.Fatalf("expected *ApiGatewayError, got %T: %v", err, err)
			}
			if apiErr.Code != "BadRequestException" {
				t.Errorf("code = %s, want BadRequestException", apiErr.Code)
			}
		})
	}
}

// TestAccountStoreRoundTrip pins the store default: a missing record reads
// as an empty account, and an update persists.
func TestAccountStoreRoundTrip(t *testing.T) {
	stores := newImportStores(t)
	account, err := stores.account.Get()
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if account.CloudwatchRoleArn != "" || len(account.Features) != 0 {
		t.Fatalf("default account not empty: %+v", account)
	}

	stored := &apigatewaystore.Account{
		CloudwatchRoleArn: "arn:probe",
		Features:          []string{"UsagePlans"},
		ThrottleSettings:  &apigatewaystore.ThrottleSettings{BurstLimit: 100, RateLimit: 5.5},
	}
	if err := stores.account.Update(stored); err != nil {
		t.Fatalf("update: %v", err)
	}
	reread, err := stores.account.Get()
	if err != nil {
		t.Fatalf("reread: %v", err)
	}
	if reread.CloudwatchRoleArn != "arn:probe" || reread.ThrottleSettings.BurstLimit != 100 || reread.ThrottleSettings.RateLimit != 5.5 {
		t.Fatalf("round-trip mismatch: %+v", reread)
	}
}
