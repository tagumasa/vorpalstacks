package lambda

import (
	"context"
	"errors"
	"testing"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// createPresenceFixture creates a function carrying the optional members the
// presence-flag tests manipulate: a non-empty Description, a non-default
// Timeout, and a DeadLetterConfig.
func createPresenceFixture(t *testing.T, svc *LambdaService, stores *lambdaStore, name string) {
	t.Helper()
	_, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName:     name,
		Runtime:          "nodejs22.x",
		Role:             "arn:aws:iam::000000000000:role/lambda",
		Handler:          "index.handler",
		Description:      "original description",
		Timeout:          30,
		DeadLetterConfig: &lambdastore.DeadLetterConfig{TargetArn: "arn:aws:sqs:us-east-1:000000000000:dlq"},
	})
	if err != nil {
		t.Fatalf("create fixture function: %v", err)
	}
}

// TestUpdateFunctionConfigurationCoreExplicitEmptyDescriptionClears pins the
// clear-on-empty rule: an explicitly provided empty Description clears the
// stored value instead of being ignored as "not provided".
func TestUpdateFunctionConfigurationCoreExplicitEmptyDescriptionClears(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	createPresenceFixture(t, svc, stores, "presence-desc")

	if _, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
		&UpdateFunctionConfigurationInput{FunctionName: "presence-desc", HasDescription: true}); err != nil {
		t.Fatalf("update with explicit empty Description: %v", err)
	}
	stored, err := stores.Functions.Get("presence-desc")
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if stored.Description != "" {
		t.Fatalf("stored Description %q, want cleared", stored.Description)
	}
}

// TestUpdateFunctionConfigurationCoreOmittedDescriptionUnchanged pins the
// other half of the presence distinction: an omitted member leaves the stored
// value untouched even though its zero value is indistinguishable from empty.
func TestUpdateFunctionConfigurationCoreOmittedDescriptionUnchanged(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	createPresenceFixture(t, svc, stores, "presence-desc-keep")

	if _, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
		&UpdateFunctionConfigurationInput{FunctionName: "presence-desc-keep"}); err != nil {
		t.Fatalf("update omitting Description: %v", err)
	}
	stored, err := stores.Functions.Get("presence-desc-keep")
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if stored.Description != "original description" {
		t.Fatalf("stored Description %q, want unchanged", stored.Description)
	}
}

// TestUpdateFunctionConfigurationCoreExplicitZeroTimeoutRejected pins that an
// explicitly provided Timeout of 0 (or negative) is range-rejected by the
// Core instead of being silently treated as unset.
func TestUpdateFunctionConfigurationCoreExplicitZeroTimeoutRejected(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	createPresenceFixture(t, svc, stores, "presence-timeout")

	for _, timeout := range []int32{0, -5} {
		_, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
			&UpdateFunctionConfigurationInput{FunctionName: "presence-timeout", HasTimeout: true, Timeout: timeout})
		var le *LambdaError
		if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
			t.Fatalf("explicit Timeout %d: expected InvalidParameterValueException, got %v", timeout, err)
		}
	}
}

// TestUpdateFunctionConfigurationCoreOmittedTimeoutUnchanged pins that an
// omitted Timeout leaves the stored value untouched.
func TestUpdateFunctionConfigurationCoreOmittedTimeoutUnchanged(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	createPresenceFixture(t, svc, stores, "presence-timeout-keep")

	if _, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
		&UpdateFunctionConfigurationInput{FunctionName: "presence-timeout-keep"}); err != nil {
		t.Fatalf("update omitting Timeout: %v", err)
	}
	stored, err := stores.Functions.Get("presence-timeout-keep")
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if stored.Timeout != 30 {
		t.Fatalf("stored Timeout %d, want unchanged 30", stored.Timeout)
	}
}

// TestUpdateFunctionConfigurationCoreExplicitZeroMemorySizeRejected pins the
// same present-member range rule for MemorySize.
func TestUpdateFunctionConfigurationCoreExplicitZeroMemorySizeRejected(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	createPresenceFixture(t, svc, stores, "presence-memory")

	_, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
		&UpdateFunctionConfigurationInput{FunctionName: "presence-memory", HasMemorySize: true, MemorySize: 0})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("explicit MemorySize 0: expected InvalidParameterValueException, got %v", err)
	}
}

// TestUpdateFunctionConfigurationCoreDeadLetterConfigPresence pins the
// DeadLetterConfig presence semantics: a present nil clears the stored
// configuration, a present config replaces it, and an omitted member leaves
// it untouched.
func TestUpdateFunctionConfigurationCoreDeadLetterConfigPresence(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	createPresenceFixture(t, svc, stores, "presence-dlc")

	if _, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
		&UpdateFunctionConfigurationInput{FunctionName: "presence-dlc", HasDeadLetterConfig: true}); err != nil {
		t.Fatalf("update with present-but-empty DeadLetterConfig: %v", err)
	}
	stored, err := stores.Functions.Get("presence-dlc")
	if err != nil {
		t.Fatalf("read back after clear: %v", err)
	}
	if stored.DeadLetterConfig != nil {
		t.Fatalf("stored DeadLetterConfig %+v, want cleared", stored.DeadLetterConfig)
	}

	replacement := &lambdastore.DeadLetterConfig{TargetArn: "arn:aws:sns:us-east-1:000000000000:topic"}
	if _, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
		&UpdateFunctionConfigurationInput{FunctionName: "presence-dlc", HasDeadLetterConfig: true, DeadLetterConfig: replacement}); err != nil {
		t.Fatalf("update with replacement DeadLetterConfig: %v", err)
	}
	stored, err = stores.Functions.Get("presence-dlc")
	if err != nil {
		t.Fatalf("read back after replace: %v", err)
	}
	if stored.DeadLetterConfig == nil || stored.DeadLetterConfig.TargetArn != replacement.TargetArn {
		t.Fatalf("stored DeadLetterConfig %+v, want %q", stored.DeadLetterConfig, replacement.TargetArn)
	}
}

// TestUpdateFunctionConfigurationHandlerPresenceMapping pins the HTTP
// handler's presence detection: wire members present with empty/zero values
// reach the Core as present (Description "" clears; Timeout 0 is rejected),
// so the two planes cannot drift back to zero-value-as-unset.
func TestUpdateFunctionConfigurationHandlerPresenceMapping(t *testing.T) {
	svc, reqCtx := workingStorageService(t)
	stores, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("acquire store: %v", err)
	}
	createPresenceFixture(t, svc, stores, "presence-handler")

	if _, err := svc.UpdateFunctionConfiguration(context.Background(), reqCtx,
		parsedRequest(map[string]interface{}{"FunctionName": "presence-handler", "Description": ""})); err != nil {
		t.Fatalf("handler update with explicit empty Description: %v", err)
	}
	stored, err := stores.Functions.Get("presence-handler")
	if err != nil {
		t.Fatalf("read back after clear: %v", err)
	}
	if stored.Description != "" {
		t.Fatalf("stored Description %q, want cleared via wire", stored.Description)
	}

	_, err = svc.UpdateFunctionConfiguration(context.Background(), reqCtx,
		parsedRequest(map[string]interface{}{"FunctionName": "presence-handler", "Timeout": float64(0)}))
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("wire Timeout 0: expected InvalidParameterValueException, got %v", err)
	}
}
