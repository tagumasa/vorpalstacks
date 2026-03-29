// Package http provides HTTP server functionality for vorpalstacks.
package http

import (
	"bufio"
	"net"
	"net/http"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
)

func sanitizeForLog(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// LoggingMiddleware logs HTTP requests and responses with timing information.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		logs.Info("HTTP request",
			logs.String("method", r.Method),
			logs.String("path", sanitizeForLog(r.URL.Path)),
			logs.String("query", sanitizeForLog(r.URL.RawQuery)),
			logs.Int("status", ww.statusCode),
			logs.Any("duration", duration),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader writes the header and captures the status code for logging.
//
// Parameters:
//   - statusCode: The HTTP status code to write
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Flush implements http.Flusher for HTTP/2 support.
func (w *responseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack implements http.Hijacker for WebSocket support.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// Push implements http.Pusher for HTTP/2 server push support.
func (w *responseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}
