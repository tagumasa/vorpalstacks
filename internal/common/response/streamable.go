package response

import (
	"io"
	"net/http"
)

// StreamableResponse is an interface for responses that should be streamed
// as raw bytes (e.g. Lambda Invoke payload).
type StreamableResponse interface {
	GetStream() io.Reader
	GetStreamHeaders() http.Header
}

// StatusCodeResponse is an interface for responses that need a custom HTTP
// status code (e.g. Lambda Invoke error responses).
type StatusCodeResponse interface {
	GetStreamStatusCode() int
}
