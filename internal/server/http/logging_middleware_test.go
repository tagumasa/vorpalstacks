package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestLoggingMiddlewarePreservesResponseController pins the wrapper's
// Unwrap: the response-controller deadline controls streaming responses
// rely on must reach the underlying connection through the logging
// middleware's writer wrapper — a wrapper without Unwrap strands the
// controller at "feature not supported" and a long-lived stream cannot
// clear the server's absolute write deadline.
func TestLoggingMiddlewarePreservesResponseController(t *testing.T) {
	var ctrlErr error
	srv := httptest.NewUnstartedServer(LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrlErr = http.NewResponseController(w).SetWriteDeadline(time.Time{})
	})))
	srv.Config.WriteTimeout = 300 * time.Millisecond
	srv.Start()
	defer srv.Close()

	if _, err := srv.Client().Get(srv.URL); err != nil {
		t.Fatalf("get: %v", err)
	}
	if ctrlErr != nil {
		t.Fatalf("ResponseController through the logging wrapper: %v", ctrlErr)
	}
}
