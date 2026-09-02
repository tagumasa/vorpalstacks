package inspection

import (
	"net/netip"
	"testing"
)

// The embedded table is a real registry snapshot; these assertions use
// allocations whose delegation is stable (1.0.0.0/8 APNIC AU,
// 8.8.8.0/24 ARIN US with Google's origin ASN 15169).
func TestEmbeddedGeoTableLookups(t *testing.T) {
	if cc, ok := CountryForIP(netip.MustParseAddr("1.0.0.1")); !ok || cc != "AU" {
		t.Fatalf("1.0.0.1 country = %q (ok %v), want AU", cc, ok)
	}
	if cc, ok := CountryForIP(netip.MustParseAddr("8.8.8.8")); !ok || cc != "US" {
		t.Fatalf("8.8.8.8 country = %q (ok %v), want US", cc, ok)
	}
	if asn, ok := ASNForIP(netip.MustParseAddr("8.8.8.8")); !ok || asn != 15169 {
		t.Fatalf("8.8.8.8 asn = %d (ok %v), want 15169", asn, ok)
	}
	// The documentation example prefix is reserved and must not resolve.
	if _, ok := CountryForIP(netip.MustParseAddr("192.0.2.1")); ok {
		t.Fatal("192.0.2.1 is documentation space and must not resolve to a country")
	}
	// IPv6 form of a v4 address resolves identically.
	if cc, ok := CountryForIP(netip.MustParseAddr("::ffff:8.8.8.8")); !ok || cc != "US" {
		t.Fatalf("::ffff:8.8.8.8 country = %q (ok %v), want US", cc, ok)
	}
}
