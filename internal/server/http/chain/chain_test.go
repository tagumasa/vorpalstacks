package chain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequestContext(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	rc := NewRequestContext(ctx, w, r)

	if rc.Context != ctx {
		t.Error("context should be set")
	}
	if rc.Request != r {
		t.Error("request should be set")
	}
	if rc.Response != w {
		t.Error("response should be set")
	}
	if rc.Params == nil {
		t.Error("params should be initialized")
	}
}

func TestRequestContext_SetHandled(t *testing.T) {
	rc := &RequestContext{}
	rc.SetHandled()

	if !rc.IsHandled() {
		t.Error("should be handled after SetHandled")
	}
}

func TestRequestContext_IsHandled_Default(t *testing.T) {
	rc := &RequestContext{}

	if rc.IsHandled() {
		t.Error("should not be handled by default")
	}
}

func TestRequestContext_SetError(t *testing.T) {
	rc := &RequestContext{}
	testErr := &testError{msg: "test error"}

	rc.SetError(testErr)

	if rc.Error != testErr {
		t.Error("error should be set")
	}
}

func TestRequestContext_SetResponse(t *testing.T) {
	rc := &RequestContext{}

	rc.SetResponse(http.StatusOK, "test body")

	if rc.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rc.StatusCode)
	}
	if rc.ResponseBody != "test body" {
		t.Error("response body should be set")
	}
}

func TestHandlerChain_NewHandlerChain(t *testing.T) {
	hc := NewHandlerChain()

	if hc == nil {
		t.Error("handler chain should not be nil")
	}
	if hc.requestHandlers == nil {
		t.Error("request handlers should be initialized")
	}
	if hc.exceptionHandlers == nil {
		t.Error("exception handlers should be initialized")
	}
	if hc.responseHandlers == nil {
		t.Error("response handlers should be initialized")
	}
	if hc.finalizers == nil {
		t.Error("finalizers should be initialized")
	}
}

func TestHandlerChain_AddRequestHandler(t *testing.T) {
	hc := NewHandlerChain()

	hc.AddRequestHandler(func(ctx *RequestContext) error {
		return nil
	})

	if len(hc.requestHandlers) != 1 {
		t.Errorf("expected 1 request handler, got %d", len(hc.requestHandlers))
	}
}

func TestHandlerChain_AddExceptionHandler(t *testing.T) {
	hc := NewHandlerChain()

	hc.AddExceptionHandler(func(ctx *RequestContext, err error) error {
		return nil
	})

	if len(hc.exceptionHandlers) != 1 {
		t.Errorf("expected 1 exception handler, got %d", len(hc.exceptionHandlers))
	}
}

func TestHandlerChain_AddResponseHandler(t *testing.T) {
	hc := NewHandlerChain()

	hc.AddResponseHandler(func(ctx *RequestContext) error {
		return nil
	})

	if len(hc.responseHandlers) != 1 {
		t.Errorf("expected 1 response handler, got %d", len(hc.responseHandlers))
	}
}

func TestHandlerChain_AddFinalizer(t *testing.T) {
	hc := NewHandlerChain()

	hc.AddFinalizer(func(ctx *RequestContext) {
	})

	if len(hc.finalizers) != 1 {
		t.Errorf("expected 1 finalizer, got %d", len(hc.finalizers))
	}
}

func TestHandlerChain_Execute_Success(t *testing.T) {
	hc := NewHandlerChain()

	hc.AddRequestHandler(func(ctx *RequestContext) error {
		return nil
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	hc.Execute(w, r)
}

func TestHandlerChain_Execute_WithError(t *testing.T) {
	hc := NewHandlerChain()
	handlerCalled := false
	exceptionCalled := false

	hc.AddRequestHandler(func(ctx *RequestContext) error {
		handlerCalled = true
		return &testError{msg: "handler error"}
	})

	hc.AddExceptionHandler(func(ctx *RequestContext, err error) error {
		exceptionCalled = true
		return nil
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	hc.Execute(w, r)

	if !handlerCalled {
		t.Error("request handler should be called")
	}
	if !exceptionCalled {
		t.Error("exception handler should be called")
	}
}

func TestHandlerChain_Execute_Handled(t *testing.T) {
	hc := NewHandlerChain()
	handler1Called := false
	handler2Called := false

	hc.AddRequestHandler(func(ctx *RequestContext) error {
		handler1Called = true
		ctx.SetHandled()
		return nil
	})

	hc.AddRequestHandler(func(ctx *RequestContext) error {
		handler2Called = true
		return nil
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	hc.Execute(w, r)

	if !handler1Called {
		t.Error("first handler should be called")
	}
	if handler2Called {
		t.Error("second handler should not be called after SetHandled")
	}
}

func TestHandlerChain_Execute_FinalizerCalled(t *testing.T) {
	hc := NewHandlerChain()

	hc.AddFinalizer(func(ctx *RequestContext) {
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	hc.Execute(w, r)
}

func TestHandlerChain_Clone(t *testing.T) {
	hc := NewHandlerChain()

	hc.AddRequestHandler(func(ctx *RequestContext) error {
		return nil
	})

	clone := hc.Clone()

	if len(clone.requestHandlers) != 1 {
		t.Errorf("clone should have 1 request handler, got %d", len(clone.requestHandlers))
	}

	clone.AddRequestHandler(func(ctx *RequestContext) error {
		return nil
	})

	if len(hc.requestHandlers) != 1 {
		t.Error("original should not be affected by clone modification")
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
