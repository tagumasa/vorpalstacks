package lambda

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
)

// failingStorageReqCtx builds a request context whose region storage can
// never be created: the manager's base path is a regular file, so every
// store acquisition fails with an infrastructure error. Operations that
// validate their parameters before acquiring the store must therefore
// surface the parameter error, not the storage failure.
func failingStorageReqCtx(t *testing.T) *request.RequestContext {
	t.Helper()
	base := filepath.Join(t.TempDir(), "occupied")
	if err := os.WriteFile(base, []byte("occupied"), 0o600); err != nil {
		t.Fatalf("write blocker file: %v", err)
	}
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: base})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	return request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
}

// workingStorageReqCtx builds a request context backed by a freshly created
// empty region store so core operations can reach their store calls.
func workingStorageReqCtx(t *testing.T) *request.RequestContext {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	return request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
}

func parsedRequest(params map[string]interface{}) *request.ParsedRequest {
	if params == nil {
		params = map[string]interface{}{}
	}
	return &request.ParsedRequest{Operation: "Test", Parameters: params}
}

// TestListLayerVersionsEmptyNameRejected pins the required-member contract
// of ListLayerVersions: an empty LayerName is InvalidParameterValueException
// (the model marks the httpLabel member required with a minimum length of 1)
// and never the ResourceNotFoundException the unguarded store lookup would
// produce — on both a healthy store and one that cannot be acquired.
func TestListLayerVersionsEmptyNameRejected(t *testing.T) {
	// Each case gets its own service: the store cache is keyed by region, so
	// a shared service would hand the second case the first case's store.
	cases := []struct {
		name   string
		reqCtx *request.RequestContext
	}{
		{"working store", workingStorageReqCtx(t)},
		{"failing storage", failingStorageReqCtx(t)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewLambdaService(nil, "000000000000", "us-east-1", t.TempDir())
			_, err := svc.ListLayerVersions(context.Background(), tc.reqCtx,
				parsedRequest(map[string]interface{}{"LayerName": ""}))
			var le *LambdaError
			if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
				t.Fatalf("expected InvalidParameterValueException, got: %v", err)
			}
		})
	}
}

// TestValidationPrecedesStoreAcquisition pins the failure precedence of the
// operations that validate their parameters ahead of the store acquisition:
// a request that is simultaneously invalid and undeliverable to storage must
// fail with InvalidParameterValueException, never with the storage error.
// Every case is shaped to pass the handler-side wire resolution and hit the
// member validation that lives on the Core layer.
func TestValidationPrecedesStoreAcquisition(t *testing.T) {
	reqCtx := failingStorageReqCtx(t)
	ctx := context.Background()

	cases := []struct {
		name   string
		invoke func(svc *LambdaService) error
	}{
		{"PutFunctionConcurrency", func(svc *LambdaService) error {
			_, err := svc.PutFunctionConcurrency(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"FunctionName":                 "fn",
				"ReservedConcurrentExecutions": float64(-1),
			}))
			return err
		}},
		{"PutFunctionEventInvokeConfig", func(svc *LambdaService) error {
			_, err := svc.PutFunctionEventInvokeConfig(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"FunctionName":         "fn",
				"MaximumRetryAttempts": float64(99),
			}))
			return err
		}},
		{"CreateEventSourceMapping", func(svc *LambdaService) error {
			_, err := svc.CreateEventSourceMapping(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"DeleteEventSourceMapping", func(svc *LambdaService) error {
			_, err := svc.DeleteEventSourceMapping(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"GetEventSourceMapping", func(svc *LambdaService) error {
			_, err := svc.GetEventSourceMapping(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"UpdateEventSourceMapping", func(svc *LambdaService) error {
			_, err := svc.UpdateEventSourceMapping(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"PublishLayerVersion", func(svc *LambdaService) error {
			_, err := svc.PublishLayerVersion(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"DeleteLayerVersion", func(svc *LambdaService) error {
			_, err := svc.DeleteLayerVersion(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"GetLayerVersion", func(svc *LambdaService) error {
			_, err := svc.GetLayerVersion(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"AddLayerVersionPermission", func(svc *LambdaService) error {
			_, err := svc.AddLayerVersionPermission(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"RemoveLayerVersionPermission", func(svc *LambdaService) error {
			_, err := svc.RemoveLayerVersionPermission(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"GetLayerVersionPolicy", func(svc *LambdaService) error {
			_, err := svc.GetLayerVersionPolicy(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"PutProvisionedConcurrencyConfig", func(svc *LambdaService) error {
			_, err := svc.PutProvisionedConcurrencyConfig(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"FunctionName":                    "fn",
				"Qualifier":                       "prod",
				"ProvisionedConcurrentExecutions": float64(0),
			}))
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewLambdaService(nil, "000000000000", "us-east-1", t.TempDir())
			err := tc.invoke(svc)
			var le *LambdaError
			if !errors.As(err, &le) || le.GetCode() != "InvalidParameterValueException" {
				t.Fatalf("expected InvalidParameterValueException, got: %v", err)
			}
		})
	}
}
