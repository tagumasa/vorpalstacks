package cognitoidentityprovider

import (
	"context"
	"errors"
	"testing"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
)

// stubTriggerInvoker is a LambdaInvoker stand-in that replays a canned
// invocation outcome so the trigger error-classification paths can be
// exercised without a real Lambda runtime.
type stubTriggerInvoker struct {
	functionError string
	payload       []byte
	invokeErr     error
}

func (s *stubTriggerInvoker) InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error) {
	return 200, s.payload, s.invokeErr
}

func (s *stubTriggerInvoker) InvokeForTrigger(ctx context.Context, functionName string, payload []byte) (invokers.LambdaInvocation, error) {
	return invokers.LambdaInvocation{
		StatusCode:    200,
		Payload:       s.payload,
		FunctionError: s.functionError,
	}, s.invokeErr
}

func (s *stubTriggerInvoker) GetFunctionARN(ctx context.Context, functionName string) (string, error) {
	return functionName, nil
}

func newTriggerTestService(t *testing.T, invoker invokers.LambdaInvoker) *CognitoService {
	t.Helper()
	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("failed to start event bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	bus.SetLambdaInvoker(invoker)
	svc := NewCognitoService("123456789012", "us-east-1")
	svc.bus = bus
	if _, err := eventbus.SubscribeTyped[*eventbus.CognitoTriggerEvent](bus, svc.handleCognitoTrigger); err != nil {
		t.Fatalf("failed to subscribe trigger handler: %v", err)
	}
	return svc
}

const triggerTestARN = "arn:aws:lambda:us-east-1:123456789012:function:migration"

// A Lambda function that raises an error must fail the migration with
// NotAuthorizedException semantics (an incorrect-password outcome), not an
// internal error.
func TestInvokeTriggerFunctionErrorClassifiedAsIncorrectPassword(t *testing.T) {
	svc := newTriggerTestService(t, &stubTriggerInvoker{functionError: "Unhandled"})

	_, err := svc.invokeTrigger(
		context.Background(),
		UserMigrationAuthentication, "us-east-1_POOL", "user", "client",
		triggerTestARN,
		map[string]interface{}{"password": "wrong"},
		map[string]interface{}{"finalUserStatus": "CONFIRMED"},
		true,
	)
	if err == nil {
		t.Fatal("expected an error from a failed trigger function")
	}
	var fnErr *lambdaFunctionError
	if !errors.As(err, &fnErr) {
		t.Fatalf("expected a lambdaFunctionError, got %T: %v", err, err)
	}
	if got := classifyMigrationFailure(err); got != ErrIncorrectPassword {
		t.Fatalf("expected ErrIncorrectPassword, got %v", got)
	}
}

// An invocation-transport failure is an infrastructure error, not a
// wrong-password outcome.
func TestInvokeTriggerTransportFailureClassifiedAsInternalError(t *testing.T) {
	svc := newTriggerTestService(t, &stubTriggerInvoker{invokeErr: errors.New("function not found")})

	_, err := svc.invokeTrigger(
		context.Background(),
		UserMigrationAuthentication, "us-east-1_POOL", "user", "client",
		triggerTestARN,
		map[string]interface{}{"password": "secret"},
		map[string]interface{}{"finalUserStatus": "CONFIRMED"},
		true,
	)
	if err == nil {
		t.Fatal("expected an error from a failed invocation")
	}
	if got := classifyMigrationFailure(err); got != ErrInternalError {
		t.Fatalf("expected ErrInternalError, got %v", got)
	}
}

// A successful invocation returns the unmarshalled Lambda response.
func TestInvokeTriggerSuccessReturnsResponse(t *testing.T) {
	svc := newTriggerTestService(t, &stubTriggerInvoker{payload: []byte(`{"finalUserStatus":"RESET_REQUIRED"}`)})

	resp, err := svc.invokeTrigger(
		context.Background(),
		UserMigrationAuthentication, "us-east-1_POOL", "user", "client",
		triggerTestARN,
		map[string]interface{}{"password": "secret"},
		map[string]interface{}{"finalUserStatus": "CONFIRMED"},
		true,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp["finalUserStatus"] != "RESET_REQUIRED" {
		t.Fatalf("expected the Lambda-provided finalUserStatus, got %v", resp["finalUserStatus"])
	}
}
