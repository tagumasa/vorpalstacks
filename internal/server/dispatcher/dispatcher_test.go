package dispatcher

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/store/api"
)

func newTestStorage(t *testing.T) storage.Storage {
	tmpDir := filepath.Join(os.TempDir(), "dispatcher-test-"+t.Name())
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	t.Cleanup(func() { s.Close() })
	return s
}

func newTestDispatcher(t *testing.T) *Dispatcher {
	s := newTestStorage(t)

	serviceStore := api.NewServiceStore(s)
	operationStore := api.NewOperationStore(s)
	shapeStore := api.NewShapeStore(s)
	memberStore := api.NewMemberStore(s)
	configStore := api.NewConfigStore(s)

	storageMgr, err := storage.NewRegionStorageManager(&storage.Config{Path: filepath.Join(os.TempDir(), "dispatcher-regions-"+t.Name())})
	require.NoError(t, err)
	t.Cleanup(func() { storageMgr.Close() })

	return NewDispatcher(
		serviceStore,
		operationStore,
		shapeStore,
		memberStore,
		configStore,
		&resilience.ServiceResilienceConfig{},
		storageMgr,
		nil,
		nil,
		"000000000000",
	)
}

func TestRegisterHandler(t *testing.T) {
	d := newTestDispatcher(t)

	handler := func(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
		return map[string]string{"status": "ok"}, nil
	}

	d.RegisterHandler("TestOperation", handler)

	h, exists := d.getHandler("", "TestOperation")
	if !exists {
		t.Fatal("handler not registered")
	}
	if h == nil {
		t.Fatal("handler is nil")
	}
}

func TestRegisterHandlerForService(t *testing.T) {
	d := newTestDispatcher(t)

	handler := func(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
		return map[string]string{"service": "test"}, nil
	}

	d.RegisterHandlerForService("TestService", "TestOp", handler)

	h, exists := d.getHandler("TestService", "TestOp")
	if !exists {
		t.Fatal("handler not registered for service")
	}
	if h == nil {
		t.Fatal("handler is nil")
	}
}

func TestHasHandlerForService(t *testing.T) {
	d := newTestDispatcher(t)

	handler := func(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
		return nil, nil
	}

	if d.hasHandlerForService("TestService", "TestOp") {
		t.Error("should not have handler before registration")
	}

	d.RegisterHandlerForService("TestService", "TestOp", handler)

	if !d.hasHandlerForService("TestService", "TestOp") {
		t.Error("should have handler after registration")
	}

	if d.hasHandlerForService("OtherService", "TestOp") {
		t.Error("should not have handler for different service")
	}
}

func TestDispatch_WithHandler(t *testing.T) {
	d := newTestDispatcher(t)

	handler := func(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
		return map[string]string{"result": "success"}, nil
	}
	d.RegisterHandler("TestOperation", handler)

	body := `{"test": "data"}`
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Amz-Target", "TestService.TestOperation")

	w := httptest.NewRecorder()
	d.Dispatch(w, req, "TestService", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestDispatch_NoService(t *testing.T) {
	d := newTestDispatcher(t)

	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	d.Dispatch(w, req, "", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestWriteResponse_JSON(t *testing.T) {
	d := newTestDispatcher(t)

	op := &api.Operation{
		Name:           "TestOp",
		SmithyProtocol: strPtr("aws.protocols#json"),
		ContentType:    strPtr("application/json"),
	}

	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	d.writeResponse(w, req, op, "", map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	respBody, _ := io.ReadAll(w.Body)
	if !bytes.Contains(respBody, []byte("key")) {
		t.Errorf("expected JSON response, got: %s", respBody)
	}
}

func TestWriteResponse_XML(t *testing.T) {
	d := newTestDispatcher(t)

	op := &api.Operation{
		Name:           "TestOp",
		SmithyProtocol: strPtr("aws.protocols#restXml"),
		ContentType:    strPtr("application/xml"),
	}

	req := httptest.NewRequest("POST", "/2013-04-01/test", nil)
	req.Header.Set("Content-Type", "application/xml")

	w := httptest.NewRecorder()
	d.writeResponse(w, req, op, "", map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestWriteResponse_QueryXML(t *testing.T) {
	d := newTestDispatcher(t)

	op := &api.Operation{
		Name:           "TestOp",
		SmithyProtocol: strPtr("aws.protocols#awsQuery"),
		ContentType:    strPtr("text/xml"),
	}

	req := httptest.NewRequest("POST", "/", bytes.NewBufferString("Action=TestOp&Version=2012-10-17"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	d.writeResponse(w, req, op, "TestOp", map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHandleError_AWSError(t *testing.T) {
	d := newTestDispatcher(t)

	op := &api.Operation{
		Name:           "TestOp",
		SmithyProtocol: strPtr("aws.protocols#json"),
		ContentType:    strPtr("application/json"),
	}

	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	testErr := errors.New("test error")
	d.handleError(w, req, op, testErr)

	if w.Code < 400 {
		t.Errorf("expected error status, got %d", w.Code)
	}
}

func TestGetErrorContentType(t *testing.T) {
	d := newTestDispatcher(t)

	t.Run("with operation content type", func(t *testing.T) {
		op := &api.Operation{
			ContentType: strPtr("application/json"),
		}
		req := httptest.NewRequest("POST", "/", nil)
		ct := d.getErrorContentType(req, op)
		if ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
	})

	t.Run("with query protocol", func(t *testing.T) {
		op := &api.Operation{
			SmithyProtocol: strPtr("aws.protocols#awsQuery"),
		}
		req := httptest.NewRequest("POST", "/", nil)
		ct := d.getErrorContentType(req, op)
		if ct != "text/xml" {
			t.Errorf("expected text/xml, got %s", ct)
		}
	})

	t.Run("from request", func(t *testing.T) {
		op := &api.Operation{}
		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		ct := d.getErrorContentType(req, op)
		if ct != "text/xml" {
			t.Errorf("expected text/xml, got %s", ct)
		}
	})
}

func TestGenerateFallbackResponse(t *testing.T) {
	d := newTestDispatcher(t)

	t.Run("nil operation", func(t *testing.T) {
		resp := d.generateFallbackResponse(nil)
		if len(resp) != 0 {
			t.Errorf("expected empty response for nil operation, got %v", resp)
		}
	})

	t.Run("operation without output shape", func(t *testing.T) {
		op := &api.Operation{
			Name: "TestOp",
		}
		resp := d.generateFallbackResponse(op)
		if len(resp) != 0 {
			t.Errorf("expected empty response, got %v", resp)
		}
	})
}

func TestDetectOperationFromParsed(t *testing.T) {
	d := newTestDispatcher(t)

	t.Run("nil request", func(t *testing.T) {
		op := d.detectOperationFromParsed(nil)
		if op != nil {
			t.Error("expected nil for nil request")
		}
	})

	t.Run("empty operation", func(t *testing.T) {
		req := &request.ParsedRequest{Operation: ""}
		op := d.detectOperationFromParsed(req)
		if op != nil {
			t.Error("expected nil for empty operation")
		}
	})
}

func TestListServices(t *testing.T) {
	d := newTestDispatcher(t)

	services, err := d.ListServices()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Empty store returns nil slice, which is valid
	_ = services
}

func TestListOperations(t *testing.T) {
	d := newTestDispatcher(t)

	ops, err := d.ListOperations("TestService")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Empty store returns nil slice, which is valid
	_ = ops
}

func TestBuildParsedRequest(t *testing.T) {
	d := newTestDispatcher(t)

	httpReq := httptest.NewRequest("POST", "/test?key=value", bytes.NewBufferString(`{"data": "test"}`))
	httpReq.Header.Set("Content-Type", "application/json")

	dispatchCtx := &DispatchContext{
		ServiceName: "TestService",
		Operation:   "TestOperation",
		Params: map[string]interface{}{
			"param1": "value1",
		},
	}

	parsedReq := d.buildParsedRequest(httpReq, dispatchCtx)

	if parsedReq.Operation != "TestOperation" {
		t.Errorf("expected TestOperation, got %s", parsedReq.Operation)
	}

	if parsedReq.Parameters["param1"] != "value1" {
		t.Errorf("expected param1=value1, got %v", parsedReq.Parameters["param1"])
	}
}

func TestCloudFrontPayloadOperations(t *testing.T) {
	tests := []struct {
		opName   string
		expected bool
		root     string
	}{
		{"CreateDistribution", true, "Distribution"},
		{"GetDistribution", true, "Distribution"},
		{"ListDistributions", true, "DistributionList"},
		{"UnknownOperation", false, "UnknownOperation"},
	}

	for _, tt := range tests {
		t.Run(tt.opName, func(t *testing.T) {
			isPayload := isCloudFrontPayloadOperation(tt.opName)
			if isPayload != tt.expected {
				t.Errorf("isCloudFrontPayloadOperation(%s) = %v, want %v", tt.opName, isPayload, tt.expected)
			}

			root := getCloudFrontPayloadRoot(tt.opName)
			if root != tt.root {
				t.Errorf("getCloudFrontPayloadRoot(%s) = %s, want %s", tt.opName, root, tt.root)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
