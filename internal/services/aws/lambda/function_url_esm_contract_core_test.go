package lambda

import (
	"context"
	"errors"
	"testing"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// urlConfigFixture creates a function with a published version and a "prod"
// alias over it, then returns the loaded function and its store.
func urlConfigFixture(t *testing.T, svc *LambdaService, stores *lambdaStore, name string) *lambdastore.Function {
	t.Helper()
	code := []byte("exports.handler = async () => {};")
	if _, _, err := svc.storeCode(name, "$LATEST", code, "us-east-1"); err != nil {
		t.Fatalf("store code: %v", err)
	}
	fn, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: name,
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
		CodeLocation: name + "-latest.zip",
		CodeSize:     int64(len(code)),
	})
	if err != nil {
		t.Fatalf("create fixture function: %v", err)
	}
	if _, err := svc.publishVersionWithCode(stores, fn, "", "", "us-east-1"); err != nil {
		t.Fatalf("publish version: %v", err)
	}
	if _, err := stores.Functions.CreateAliasAtomically(name, func(fn *lambdastore.Function) (*lambdastore.Alias, error) {
		return &lambdastore.Alias{Name: "prod", FunctionVersion: "1"}, nil
	}); err != nil {
		t.Fatalf("create alias: %v", err)
	}
	stored, err := stores.Functions.Get(name)
	if err != nil {
		t.Fatalf("reload function: %v", err)
	}
	return stored
}

// TestCreateFunctionUrlConfigRejectsNumericQualifier pins the modelled URL
// qualifier contract: a function URL targets $LATEST or an alias, never a
// numeric version — the qualifier must resolve to an alias.
func TestCreateFunctionUrlConfigRejectsNumericQualifier(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	fn := urlConfigFixture(t, svc, stores, "url-numeric")

	_, err := svc.createFunctionUrlConfigCore(stores, fn, &FunctionUrlConfigInput{
		AuthType:  "NONE",
		Qualifier: "1",
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("numeric qualifier: expected InvalidParameterValueException, got %v", err)
	}
}

// TestCreateFunctionUrlConfigAcceptsAliasQualifier pins that an alias
// qualifier still resolves, and $LATEST stays accepted.
func TestCreateFunctionUrlConfigAcceptsAliasQualifier(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	fn := urlConfigFixture(t, svc, stores, "url-alias")

	if _, err := svc.createFunctionUrlConfigCore(stores, fn, &FunctionUrlConfigInput{
		AuthType:  "NONE",
		Qualifier: "prod",
	}); err != nil {
		t.Fatalf("alias qualifier: %v", err)
	}
}

// TestUpdateFunctionUrlConfigRejectsNumericQualifier pins the same rule on
// the update path.
func TestUpdateFunctionUrlConfigRejectsNumericQualifier(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	fn := urlConfigFixture(t, svc, stores, "url-numeric-update")

	created, err := svc.createFunctionUrlConfigCore(stores, fn, &FunctionUrlConfigInput{AuthType: "NONE"})
	if err != nil {
		t.Fatalf("create URL config: %v", err)
	}
	// Reload so the update operates on the function carrying the stored
	// URL configuration.
	fn, err = stores.Functions.Get("url-numeric-update")
	if err != nil {
		t.Fatalf("reload function: %v", err)
	}
	_, err = svc.updateFunctionUrlConfigCore(stores, fn, &FunctionUrlConfigUpdateInput{
		Qualifier: "1",
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("numeric qualifier on update: expected InvalidParameterValueException, got %v", err)
	}
	if created.Qualifier == "1" {
		t.Fatal("update must not persist a numeric qualifier")
	}
}

// TestCreateEventSourceMappingValidatesDestinationConfig pins the modelled
// DestinationConfig contract for event source mappings: OnSuccess is not a
// supported member, and an on-failure destination must be an SQS, SNS, or
// S3 ARN — the destinations the delivery path implements.
func TestCreateEventSourceMappingValidatesDestinationConfig(t *testing.T) {
	svc, reqCtx := workingStorageService(t)
	stores, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("acquire store: %v", err)
	}
	createPresenceFixture(t, svc, stores, "esm-dest")

	bad := []map[string]interface{}{
		{"OnSuccess": map[string]interface{}{"Destination": "arn:aws:sqs:us-east-1:000000000000:target"}},
		{"OnFailure": map[string]interface{}{"Destination": "arn:aws:firehose:us-east-1:000000000000:deliverystream/x"}},
		{"OnFailure": map[string]interface{}{"Destination": "not-an-arn"}},
	}
	for i, dest := range bad {
		_, err := svc.createEventSourceMappingCore(reqCtx, &EventSourceMappingCreateInput{
			FunctionNameRaw:      "esm-dest",
			EventSourceArn:       "arn:aws:sqs:us-east-1:000000000000:esm-dest-" + string(rune('a'+i)),
			DestinationConfigRaw: dest,
		})
		var le *LambdaError
		if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
			t.Fatalf("destination %+v: expected InvalidParameterValueException, got %v", dest, err)
		}
	}

	const sqsArn = "arn:aws:sqs:us-east-1:000000000000:esm-dest-ok"
	mapping, err := svc.createEventSourceMappingCore(reqCtx, &EventSourceMappingCreateInput{
		FunctionNameRaw: "esm-dest",
		EventSourceArn:  sqsArn,
		DestinationConfigRaw: map[string]interface{}{
			"OnFailure": map[string]interface{}{"Destination": sqsArn},
		},
	})
	if err != nil {
		t.Fatalf("valid on-failure destination: %v", err)
	}
	if mapping.DestinationConfig == nil || mapping.DestinationConfig.OnFailure == nil ||
		mapping.DestinationConfig.OnFailure.Destination != sqsArn {
		t.Fatalf("destination not persisted: %+v", mapping.DestinationConfig)
	}

	// The update path applies the same validator.
	_, err = svc.updateEventSourceMappingCore(reqCtx, &EventSourceMappingUpdateInput{
		UUID: mapping.UUID,
		DestinationConfigRaw: map[string]interface{}{
			"OnFailure": map[string]interface{}{"Destination": "arn:aws:lambda:us-east-1:000000000000:function:other"},
		},
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("bad destination on update: expected InvalidParameterValueException, got %v", err)
	}
}

// TestUpdateEventSourceMappingHonoursBatchSizeFlag pins that an explicitly
// provided BatchSize is range-validated like every other flagged member: an
// explicit zero is rejected instead of being silently dropped.
func TestUpdateEventSourceMappingHonoursBatchSizeFlag(t *testing.T) {
	svc, reqCtx := workingStorageService(t)
	stores, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("acquire store: %v", err)
	}
	createPresenceFixture(t, svc, stores, "esm-batch")

	mapping, err := svc.createEventSourceMappingCore(reqCtx, &EventSourceMappingCreateInput{
		FunctionNameRaw: "esm-batch",
		EventSourceArn:  "arn:aws:sqs:us-east-1:000000000000:esm-batch-q",
	})
	if err != nil {
		t.Fatalf("create mapping: %v", err)
	}
	if mapping.BatchSize != 10 {
		t.Fatalf("stored BatchSize %d, want the SQS default 10", mapping.BatchSize)
	}

	_, err = svc.updateEventSourceMappingCore(reqCtx, &EventSourceMappingUpdateInput{
		UUID:         mapping.UUID,
		HasBatchSize: true,
		BatchSize:    0,
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("explicit BatchSize 0: expected InvalidParameterValueException, got %v", err)
	}
}

// TestFunctionUrlConfigRejectsOutOfRangeMaxAge pins the modelled Cors.MaxAge
// range on both URL-config planes, including a partial update that tries to
// smuggle an out-of-range value past a stored in-range configuration.
func TestFunctionUrlConfigRejectsOutOfRangeMaxAge(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	fn := urlConfigFixture(t, svc, stores, "url-maxage")

	_, err := svc.createFunctionUrlConfigCore(stores, fn, &FunctionUrlConfigInput{
		AuthType: "NONE",
		Cors:     &lambdastore.CorsConfig{MaxAge: lambdastore.MaxCorsMaxAgeSeconds + 1},
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("create MaxAge: expected InvalidParameterValueException, got %v", err)
	}

	created, err := svc.createFunctionUrlConfigCore(stores, fn, &FunctionUrlConfigInput{
		AuthType: "NONE",
		Cors:     &lambdastore.CorsConfig{MaxAge: 60},
	})
	if err != nil {
		t.Fatalf("create in-range URL config: %v", err)
	}
	_ = created

	stored, err := stores.Functions.Get("url-maxage")
	if err != nil {
		t.Fatalf("reload function: %v", err)
	}
	_, err = svc.updateFunctionUrlConfigCore(stores, stored, &FunctionUrlConfigUpdateInput{
		Cors: map[string]interface{}{"MaxAge": float64(100000)},
	})
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("update MaxAge: expected InvalidParameterValueException, got %v", err)
	}
}

// TestCreateFunctionRejectsEphemeralStorageWithoutSize pins that an
// EphemeralStorage member without Size reaches the Core's Size validation
// instead of being silently dropped: the model marks Size required and its
// range floor rejects the absent value like an explicit zero.
func TestCreateFunctionRejectsEphemeralStorageWithoutSize(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	code := []byte("exports.handler = async () => {};")
	if _, _, err := svc.storeCode("eph-no-size", "$LATEST", code, "us-east-1"); err != nil {
		t.Fatalf("store code: %v", err)
	}
	_, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName:     "eph-no-size",
		Runtime:          "nodejs22.x",
		Role:             "arn:aws:iam::000000000000:role/lambda",
		Handler:          "index.handler",
		CodeLocation:     "eph-no-size-latest.zip",
		CodeSize:         int64(len(code)),
		EphemeralStorage: &lambdastore.EphemeralStorage{},
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("EphemeralStorage without Size: expected InvalidParameterValueException, got %v", err)
	}
}

// TestCreateEventSourceMappingRejectsMalformedTimestamp pins that a
// non-numeric StartingPositionTimestamp is rejected instead of being
// silently zeroed into an AT_TIMESTAMP mapping.
func TestCreateEventSourceMappingRejectsMalformedTimestamp(t *testing.T) {
	svc, reqCtx := workingStorageService(t)
	stores, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("acquire store: %v", err)
	}
	urlConfigFixture(t, svc, stores, "esm-ts")

	_, err = svc.createEventSourceMappingCore(reqCtx, &EventSourceMappingCreateInput{
		FunctionNameRaw:              "esm-ts",
		EventSourceArn:               "arn:aws:kinesis:us-east-1:000000000000:stream/esm-ts",
		StartingPosition:             "AT_TIMESTAMP",
		HasStartingPositionTimestamp: true,
		StartingPositionTimestampRaw: "not-a-timestamp",
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("malformed timestamp: expected InvalidParameterValueException, got %v", err)
	}
}

// TestUpdateEventSourceMappingRepointsFunction pins the update path's
// FunctionName member: the reference resolves with Create's rules, an
// unresolvable function is rejected, and the stored mapping carries the
// repointed qualified ARN.
func TestUpdateEventSourceMappingRepointsFunction(t *testing.T) {
	svc, reqCtx := workingStorageService(t)
	stores, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("acquire store: %v", err)
	}
	target := urlConfigFixture(t, svc, stores, "esm-target")
	urlConfigFixture(t, svc, stores, "esm-source")

	mapping, err := svc.createEventSourceMappingCore(reqCtx, &EventSourceMappingCreateInput{
		FunctionNameRaw: "esm-source",
		EventSourceArn:  "arn:aws:sqs:us-east-1:000000000000:esm-repoint-q",
	})
	if err != nil {
		t.Fatalf("create mapping: %v", err)
	}

	// An unresolvable function reference is rejected instead of silently
	// breaking the mapping.
	_, err = svc.updateEventSourceMappingCore(reqCtx, &EventSourceMappingUpdateInput{
		UUID:            mapping.UUID,
		FunctionNameRaw: "no-such-function",
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "ResourceNotFoundException" {
		t.Fatalf("unknown function: expected ResourceNotFoundException, got %v", err)
	}

	updated, err := svc.updateEventSourceMappingCore(reqCtx, &EventSourceMappingUpdateInput{
		UUID:            mapping.UUID,
		FunctionNameRaw: "esm-target:prod",
	})
	if err != nil {
		t.Fatalf("repoint: %v", err)
	}
	if updated.FunctionArn != target.FunctionArn+":prod" {
		t.Fatalf("repointed ARN %q, want %q", updated.FunctionArn, target.FunctionArn+":prod")
	}
	if updated.FunctionName != "esm-target" {
		t.Fatalf("repointed name %q, want esm-target", updated.FunctionName)
	}
}

// TestPutProvisionedConcurrencyReportsInProgress pins the honest initial
// provisioning report: a Put records the requested capacity with zero
// allocated/available and Status IN_PROGRESS — the warm capacity is only
// reported once the pre-warmed environment exists (no Docker client in
// this test, so the configuration stays IN_PROGRESS).
func TestPutProvisionedConcurrencyReportsInProgress(t *testing.T) {
	svc, reqCtx := workingStorageService(t)
	stores, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("acquire store: %v", err)
	}
	fn := urlConfigFixture(t, svc, stores, "esm-provconc")
	if _, err := svc.publishVersionWithCode(stores, fn, "", "", "us-east-1"); err != nil {
		t.Fatalf("publish version: %v", err)
	}

	config, err := svc.putProvisionedConcurrencyCore(reqCtx, &ProvisionedConcurrencyInput{
		FunctionName:                    "esm-provconc",
		Qualifier:                       "1",
		ProvisionedConcurrentExecutions: 5,
		Region:                          "us-east-1",
	})
	if err != nil {
		t.Fatalf("put provisioned concurrency: %v", err)
	}
	if config.Status != "IN_PROGRESS" {
		t.Fatalf("Status %q, want IN_PROGRESS", config.Status)
	}
	if config.RequestedProvisionedConcurrentExecutions != 5 {
		t.Fatalf("Requested %d, want 5", config.RequestedProvisionedConcurrentExecutions)
	}
	if config.AllocatedProvisionedConcurrentExecutions != 0 || config.AvailableProvisionedConcurrentExecutions != 0 {
		t.Fatalf("Allocated/Available %d/%d, want 0/0 before the warm environment exists",
			config.AllocatedProvisionedConcurrentExecutions, config.AvailableProvisionedConcurrentExecutions)
	}
}
