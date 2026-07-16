// Package fqdnrouter provides a Host-header multiplexer that routes
// requests to the correct handler based on the Host header suffix.
// It is used to support FQDN mode where multiple services share a
// single listener port.
package fqdnrouter

import (
	"context"
	"net/http"
	"strings"
	"sync"
)

type contextKey string

const resourceIDKey contextKey = "fqdn_resource_id"

// ResourceIDFromContext extracts the resource ID injected by the Router.
// Returns empty string if not present.
func ResourceIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(resourceIDKey).(string); ok {
		return v
	}
	return ""
}

// route binds a host suffix to a handler.
type route struct {
	hostSuffix string
	handler    http.Handler
}

// Router is an http.Handler that dispatches based on the Host header.
// It matches the Host against registered suffixes in registration order.
// The Host header is stripped of port before matching.
// When no suffix matches and a fallback handler is set, the fallback
// is invoked instead of returning 404.
type Router struct {
	mu       sync.RWMutex
	routes   []route
	fallback http.Handler
}

// New creates a new Router.
func New() *Router {
	return &Router{}
}

// Register adds a handler for the given host suffix.
// The suffix is matched against the end of the Host header (after
// stripping the port). For example, Register("mybucket.s3-website.us-east-1.amazonaws.com", h)
// matches requests with Host "mybucket.s3-website.us-east-1.amazonaws.com:50080".
func (r *Router) Register(hostSuffix string, handler http.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes = append(r.routes, route{hostSuffix: strings.ToLower(hostSuffix), handler: handler})
}

// SetFallback sets the handler invoked when no host suffix matches.
// This is used when a listener has a direct handler that should serve
// as the default for unmatched Host headers.
func (r *Router) SetFallback(h http.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fallback = h
}

// ServeHTTP matches the request Host against registered suffixes and
// dispatches to the matching handler. The resource ID (Host prefix before
// the matched suffix) is injected into the request context.
// If no suffix matches, the fallback handler is used if set, otherwise 404.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := strings.ToLower(req.Host)
	// Strip port
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	r.mu.RLock()
	routes := r.routes
	r.mu.RUnlock()

	for _, rt := range routes {
		if strings.HasSuffix(host, rt.hostSuffix) {
			// Extract resource ID: everything before the matched suffix.
			// E.g. host="mybucket.s3-website.us-east-1.amazonaws.com",
			//      suffix=".s3-website.us-east-1.amazonaws.com"
			//      resourceID="mybucket"
			resourceID := strings.TrimSuffix(host, rt.hostSuffix)
			resourceID = strings.TrimSuffix(resourceID, ".")

			ctx := context.WithValue(req.Context(), resourceIDKey, resourceID)
			rt.handler.ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}

	r.mu.RLock()
	fb := r.fallback
	r.mu.RUnlock()

	if fb != nil {
		fb.ServeHTTP(w, req)
		return
	}
	http.NotFound(w, req)
}

// HasRoutes reports whether any routes are registered.
func (r *Router) HasRoutes() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.routes) > 0
}
