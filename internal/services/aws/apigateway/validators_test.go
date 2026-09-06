package apigateway

import "testing"

// TestValidCacheClusterSizesMessage pins the derived cacheClusterSize
// message byte-for-byte: the three deployment validation sites share this
// string, and a reordering or reformatting would change the wire error.
func TestValidCacheClusterSizesMessage(t *testing.T) {
	want := "Invalid cacheClusterSize: must be one of 0.5, 1.6, 6.1, 13.5, 28.4, 58.2, 118, 237"
	if got := validCacheClusterSizesMessage(); got != want {
		t.Errorf("validCacheClusterSizesMessage() = %q, want %q", got, want)
	}
}
