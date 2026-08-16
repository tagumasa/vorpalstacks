package cloudfront

import (
	"sort"
	"testing"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
)

// The routing guard keeps the CloudFront REST surface whole: every operation
// the service registers must be routable by the request parser, and every
// route must belong to a registered operation. Both are pure code-to-code
// comparisons with no external file dependencies; wire-level conformance to
// the API model is pinned by the SDK test suite.

type recordingRegistrar struct {
	ops map[string]bool
}

func (r *recordingRegistrar) RegisterHandler(operationName string, h handler.Handler) {}

func (r *recordingRegistrar) RegisterHandlerForService(serviceName, operationName string, h handler.Handler) {
	if serviceName == "cloudfront" {
		r.ops[operationName] = true
	}
}

func TestCloudFrontRoutingMatchesRegistry(t *testing.T) {
	registrar := &recordingRegistrar{ops: map[string]bool{}}
	NewCloudFrontService("000000000000").RegisterHandlers(registrar)

	routed := map[string]bool{}
	for _, route := range request.CloudFrontRouteTable() {
		routed[route.Op] = true
	}

	var unrouted, unregistered []string
	for op := range registrar.ops {
		if !routed[op] {
			unrouted = append(unrouted, op)
		}
	}
	for op := range routed {
		if !registrar.ops[op] {
			unregistered = append(unregistered, op)
		}
	}
	sort.Strings(unrouted)
	sort.Strings(unregistered)
	if len(unrouted) > 0 {
		t.Errorf("operations registered in the handler registry but absent from the routing table (unroutable over HTTP): %v", unrouted)
	}
	if len(unregistered) > 0 {
		t.Errorf("operations in the routing table but not registered in the handler registry: %v", unregistered)
	}
}
