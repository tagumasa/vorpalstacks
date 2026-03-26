// Package conditions provides AWS IAM condition evaluation utilities for vorpalstacks.
package conditions

import (
	"net/http"
	"strings"
	"time"
)

// EvaluationContext holds the context for evaluating IAM policy conditions.
type EvaluationContext struct {
	CurrentTime    time.Time
	SourceIP       string
	UserAgent      string
	Referer        string
	RequestHeaders map[string]string
	Environment    map[string]string
	PrincipalARN   string
	PrincipalID    string
	ResourceARN    string
	Action         string
}

// BuildEvaluationContext creates an EvaluationContext from an HTTP request.
//
// Parameters:
//   - r: The HTTP request
//
// Returns:
//   - *EvaluationContext: The evaluation context
func BuildEvaluationContext(r *http.Request) *EvaluationContext {
	ctx := &EvaluationContext{
		CurrentTime:    time.Now(),
		RequestHeaders: make(map[string]string),
		Environment:    make(map[string]string),
	}

	if r == nil {
		return ctx
	}

	ctx.SourceIP = extractSourceIP(r)
	ctx.UserAgent = r.UserAgent()
	ctx.Referer = r.Referer()
	ctx.RequestHeaders = extractHeaders(r.Header)

	return ctx
}

// extractSourceIP extracts the source IP address from an HTTP request.
func extractSourceIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	addr := r.RemoteAddr
	if colonIdx := strings.LastIndex(addr, ":"); colonIdx != -1 {
		addr = addr[:colonIdx]
	}
	if strings.HasPrefix(addr, "[") && strings.HasSuffix(addr, "]") {
		addr = addr[1 : len(addr)-1]
	}

	return addr
}

// extractHeaders extracts headers from an HTTP request.
func extractHeaders(header http.Header) map[string]string {
	result := make(map[string]string)
	for k, v := range header {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}
