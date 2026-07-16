// Package telemetry provides OpenTelemetry tracing and metrics for vorpalstacks.
package telemetry

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// HTTPMiddleware creates an HTTP middleware that adds OpenTelemetry tracing.
func HTTPMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "http.request",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}

// SpanMiddleware creates a middleware that adds custom span attributes.
func SpanMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := StartSpan(r.Context(), "http.request",
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.host", r.Host),
				attribute.String("http.scheme", r.URL.Scheme),
				attribute.String("http.remote_addr", r.RemoteAddr),
				attribute.String("http.user_agent", r.UserAgent()),
			),
		)
		defer span.End()

		span.SetAttributes(attribute.String("http.route", r.URL.Path))

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r.WithContext(ctx))

		span.SetAttributes(attribute.Int("http.status_code", rw.statusCode))
		if rw.statusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(rw.statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	})
}

// responseWriter is a wrapper around http.ResponseWriter that captures status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before delegating to the underlying ResponseWriter.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
