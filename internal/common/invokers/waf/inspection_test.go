package waf

import (
	"strings"
	"testing"
)

// TestBuildWebACLInspectionRequestHeaderOrder pins the header emission
// order: a captured wire order drives the emission (with the request's
// own casing), headers the capture missed are appended, and without a
// capture the map's unspecified order stands.
func TestBuildWebACLInspectionRequestHeaderOrder(t *testing.T) {
	headers := map[string][]string{
		"Host":       {"example.test"},
		"X-Tenant":   {"a"},
		"User-Agent": {"curl/8"},
		"X-Extra":    {"late"},
	}

	captured := []string{"x-tenant", "host", "user-agent"}
	req := BuildWebACLInspectionRequest("GET", "/", "", "127.0.0.1", "HTTP/1.1", headers, captured, nil, false)
	var names []string
	for _, h := range req.Headers {
		names = append(names, h.Name)
	}
	got := strings.Join(names, ",")
	want := "x-tenant,host,user-agent,X-Extra"
	if len(names) != 4 || got != want {
		t.Fatalf("header emission = %q, want %q", got, want)
	}

	fallback := BuildWebACLInspectionRequest("GET", "/", "", "127.0.0.1", "HTTP/1.1", headers, nil, nil, false)
	if len(fallback.Headers) != 4 {
		t.Fatalf("fallback emission = %d headers, want 4", len(fallback.Headers))
	}

	cookies := map[string][]string{"Cookie": {"session=abc; theme=dark"}}
	req = BuildWebACLInspectionRequest("GET", "/", "", "127.0.0.1", "HTTP/1.1", cookies, []string{"cookie"}, nil, false)
	if len(req.Cookies) != 2 || req.Cookies[0].Name != "session" || req.Cookies[1].Name != "theme" {
		t.Fatalf("parsed cookies = %+v", req.Cookies)
	}
}
