package lambda

import (
	"context"
	"errors"
	"testing"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// TestCreateFunctionCoreRejectsBadLoggingConfig pins that the modelled
// LoggingConfig enums are enforced on the create path, not only by the
// standalone validator.
func TestCreateFunctionCoreRejectsBadLoggingConfig(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)

	_, _, err := svc.createFunctionCore(context.Background(), stores, &CreateFunctionInput{
		FunctionName: "logging-bad",
		Runtime:      "nodejs22.x",
		Role:         "arn:aws:iam::000000000000:role/lambda",
		Handler:      "index.handler",
		LoggingConfig: &lambdastore.LoggingConfig{
			LogFormat: "Banana",
		},
	})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("bad LogFormat on create: expected InvalidParameterValueException, got %v", err)
	}
}

// TestUpdateFunctionConfigurationCoreRejectsBadImageConfig pins the same
// wiring for ImageConfig on the update path.
func TestUpdateFunctionConfigurationCoreRejectsBadImageConfig(t *testing.T) {
	svc, stores := canonicalRuntimeService(t)
	createPresenceFixture(t, svc, stores, "image-bad")

	longDir := make([]byte, lambdastore.MaxImageConfigWorkingDirectory+1)
	for i := range longDir {
		longDir[i] = 'd'
	}
	_, err := svc.updateFunctionConfigurationCore(context.Background(), stores,
		&UpdateFunctionConfigurationInput{
			FunctionName: "image-bad",
			ImageConfig:  &lambdastore.ImageConfig{WorkingDirectory: string(longDir)},
		})
	var le *LambdaError
	if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
		t.Fatalf("over-long WorkingDirectory on update: expected InvalidParameterValueException, got %v", err)
	}
}
