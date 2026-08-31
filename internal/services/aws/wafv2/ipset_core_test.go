package wafv2

import (
	"strings"
	"testing"
)

// TestParseAddressListRequiredMember pins the Addresses contract of
// CreateIPSet and UpdateIPSet: the member is required, so an omitted
// (nil) value and a non-list value are rejected, while the documented
// empty array "Addresses": [] stays valid.
func TestParseAddressListRequiredMember(t *testing.T) {
	if _, err := parseAddressList(nil, "IPV4"); err == nil || !strings.Contains(err.Error(), "WAFInvalidParameterException") {
		t.Fatalf("nil Addresses must be rejected with WAFInvalidParameterException, got %v", err)
	}
	if _, err := parseAddressList("192.0.2.0/24", "IPV4"); err == nil || !strings.Contains(err.Error(), "WAFInvalidParameterException") {
		t.Fatalf("non-list Addresses must be rejected with WAFInvalidParameterException, got %v", err)
	}
	addresses, err := parseAddressList([]interface{}{}, "IPV4")
	if err != nil {
		t.Fatalf("empty Addresses array must be valid: %v", err)
	}
	if len(addresses) != 0 {
		t.Fatalf("empty Addresses array must parse to zero entries, got %d", len(addresses))
	}
}

// TestParseAddressListEntries pins entry-level parsing: string CIDR
// entries pass validation, non-string entries are rejected.
func TestParseAddressListEntries(t *testing.T) {
	addresses, err := parseAddressList([]interface{}{"192.0.2.0/24", "198.51.100.7/32"}, "IPV4")
	if err != nil {
		t.Fatalf("valid CIDR entries must parse: %v", err)
	}
	if len(addresses) != 2 {
		t.Fatalf("expected 2 addresses, got %d", len(addresses))
	}
	if _, err := parseAddressList([]interface{}{42}, "IPV4"); err == nil || !strings.Contains(err.Error(), "WAFInvalidParameterException") {
		t.Fatalf("non-string Addresses entry must be rejected, got %v", err)
	}
}
