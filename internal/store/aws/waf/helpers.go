package waf

import (
	"github.com/google/uuid"
)

// GenerateLockToken generates a unique lock token for WAF resource
// optimistic concurrency control. The token is a UUID (RFC 4122) to
// match the Smithy LockToken pattern:
// ^[0-9a-f]{8}-(?:[0-9a-f]{4}-){3}[0-9a-f]{12}$
func GenerateLockToken() string {
	return uuid.New().String()
}
