package http

import (
	"net/http"
	"strings"
	"testing"
)

// TestIsUnauthenticatedOperation pins the request shapes that qualify as
// unauthenticated STS operations: the Action parameter is read from the
// form body (fully percent-decoded) or the query string, and anything
// else fails closed.
func TestIsUnauthenticatedOperation(t *testing.T) {
	tests := []struct {
		name   string
		target string
		body   string
		want   bool
	}{
		{"plain body action", "/", "Action=AssumeRoleWithSAML&RoleArn=x", true},
		{"percent-encoded body action", "/", "Action=AssumeRoleWith%53AML", true},
		{"query string action", "/?Action=AssumeRoleWithWebIdentity", "", true},
		{"percent-encoded query action", "/?Action=AssumeRoleWith%57ebIdentity", "", true},
		{"authenticated action", "/", "Action=AssumeRole", false},
		{"other body action", "/", "Action=GetCallerIdentity", false},
		{"no action", "/", "Version=2011-06-15", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "http://localhost"+tt.target, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if got := isUnauthenticatedOperation(req); got != tt.want {
				t.Errorf("isUnauthenticatedOperation(%q, %q) = %v, want %v", tt.target, tt.body, got, tt.want)
			}
		})
	}
}
