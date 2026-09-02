package main

import (
	"bytes"
	"encoding/binary"
	"net/netip"
	"strings"
	"testing"
)

func TestParseDelegatedRecords(t *testing.T) {
	input := `# header comment
ripencc|*|ipv4|*|100645|summary
ripencc|PS|ipv4|1.178.112.0|4096|20071126|allocated|opaque
apnic|JP|ipv4|126.0.0.0|1048576|20100101|allocated
apnic|AU|ipv4|1.0.0.0|2560|20100801|allocated
ripencc|DE|ipv6|2001:db8::|32|19990801|allocated
ripencc|ZZ|ipv4|9.9.9.9|256|20100101|reserved
ripencc|FR|ipv4|8.8.8.0|24|20100101|available
`
	var records []delegatedRecord
	records, err := parseDelegated(strings.NewReader(input), records)
	if err != nil {
		t.Fatal(err)
	}
	byCountry := map[string][]netip.Prefix{}
	for _, rec := range records {
		byCountry[rec.CountryCode] = append(byCountry[rec.CountryCode], rec.Prefix)
	}
	if got := byCountry["PS"]; len(got) != 1 || got[0].String() != "1.178.112.0/20" {
		t.Fatalf("PS prefixes = %v", got)
	}
	if got := byCountry["JP"]; len(got) != 1 || got[0].String() != "126.0.0.0/12" {
		t.Fatalf("JP prefixes = %v", got)
	}
	// 2560 addresses decompose into aligned prefixes /21+/23+/23... via
	// the greedy split: 2560 = 2048+512.
	got := byCountry["AU"]
	if len(got) == 0 {
		t.Fatalf("AU prefixes missing: %v", records)
	}
	var covered uint64
	for _, p := range got {
		covered += 1 << uint(32-p.Bits())
	}
	if covered != 2560 {
		t.Fatalf("AU coverage = %d, want 2560 (%v)", covered, got)
	}
	if got := byCountry["DE"]; len(got) != 1 || got[0].String() != "2001:db8::/32" {
		t.Fatalf("DE prefixes = %v", got)
	}
	if _, ok := byCountry["ZZ"]; ok {
		t.Fatal("reserved records must be skipped")
	}
	if _, ok := byCountry["FR"]; ok {
		t.Fatal("available records must be skipped")
	}
}

// A delegated start address is not guaranteed to be aligned to any
// power of two; the decomposition must cover exactly [start,
// start+count) without masking the start address down to a coarser
// boundary.
func TestAlignedPrefixesUnalignedStart(t *testing.T) {
	prefixes := alignedPrefixes(netip.MustParseAddr("1.2.3.64"), 192)
	// The start's alignment caps the first block at /26; the remaining
	// 128 addresses are one aligned /25.
	want := []string{"1.2.3.64/26", "1.2.3.128/25"}
	if len(prefixes) != len(want) {
		t.Fatalf("prefixes = %v, want %v", prefixes, want)
	}
	for i, p := range prefixes {
		if p.String() != want[i] {
			t.Fatalf("prefix %d = %s, want %s", i, p.String(), want[i])
		}
	}
}

// assertExactCoverage verifies the decomposition invariants: the
// prefixes are contiguous, disjoint, start at start and together cover
// exactly count addresses.
func assertExactCoverage(t *testing.T, start netip.Addr, count uint64, prefixes []netip.Prefix) {
	t.Helper()
	base := start.As4()
	next := uint32(base[0])<<24 | uint32(base[1])<<16 | uint32(base[2])<<8 | uint32(base[3])
	var covered uint64
	for _, p := range prefixes {
		addr := p.Masked().Addr().As4()
		got := uint32(addr[0])<<24 | uint32(addr[1])<<16 | uint32(addr[2])<<8 | uint32(addr[3])
		if got != next {
			t.Fatalf("prefix %s starts at %08x, want contiguous %08x", p.String(), got, next)
		}
		size := uint64(1) << uint(32-p.Bits())
		covered += size
		next += uint32(size)
	}
	if covered != count {
		t.Fatalf("coverage = %d, want %d", covered, count)
	}
}

func TestAlignedPrefixesExactCoverage(t *testing.T) {
	for _, tc := range []struct {
		start string
		count uint64
	}{
		{"1.0.0.0", 2560},
		{"1.2.3.64", 192},
		{"192.0.2.13", 77},
		{"203.0.113.128", 1024},
		{"9.9.9.7", 9},
	} {
		start := netip.MustParseAddr(tc.start)
		prefixes := alignedPrefixes(start, tc.count)
		if len(prefixes) == 0 {
			t.Fatalf("%s/%d decomposed to nothing", tc.start, tc.count)
		}
		assertExactCoverage(t, start, tc.count, prefixes)
	}
}

// A malformed ipv6 start address must be skipped like any other
// unparseable record, not abort the generator.
func TestParseDelegatedSkipsMalformedIPv6(t *testing.T) {
	input := `ripencc|DE|ipv6|not-an-addr|32|19990801|allocated
ripencc|DE|ipv6|2001:db8::|32|19990801|allocated
ripencc|DE|ipv4|1.2.3.64|192|20260902|allocated
`
	var records []delegatedRecord
	records, err := parseDelegated(strings.NewReader(input), records)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 {
		t.Fatalf("records = %v, want the valid ipv6 allocation and the two ipv4 blocks", records)
	}
	if records[0].Prefix.String() != "2001:db8::/32" {
		t.Fatalf("ipv6 record = %v", records[0])
	}
	var ipv4Prefixes []netip.Prefix
	for _, rec := range records[1:] {
		ipv4Prefixes = append(ipv4Prefixes, rec.Prefix)
	}
	assertExactCoverage(t, netip.MustParseAddr("1.2.3.64"), 192, ipv4Prefixes)
}

func TestTrieFlattenLongestPrefixWins(t *testing.T) {
	trie := &trieNode{}
	insert4 := func(cidr string, value uint32) {
		p := netip.MustParsePrefix(cidr)
		raw, bits := prefix128(p)
		trie.insert(raw, bits, value)
	}
	insert4("10.0.0.0/8", 1)
	insert4("10.1.0.0/16", 2)
	insert4("10.1.2.0/24", 3)
	insert4("192.168.0.0/16", 4)

	intervals := trie.flatten()
	lookup := func(addr string) (uint32, bool) {
		raw, _ := prefix128(netip.MustParsePrefix(addr + "/32"))
		for _, iv := range intervals {
			if containsAddr(iv, raw) {
				return iv.value, true
			}
		}
		return 0, false
	}
	for addr, want := range map[string]uint32{
		"10.0.0.1": 1, "10.1.0.1": 2, "10.1.2.3": 3, "10.2.0.1": 1,
		"192.168.1.1": 4, "11.0.0.1": 0,
	} {
		got, ok := lookup(addr)
		if want == 0 {
			if ok {
				t.Fatalf("%s unexpectedly in interval (value %d)", addr, got)
			}
			continue
		}
		if !ok || got != want {
			t.Fatalf("%s value = %d (ok %v), want %d", addr, got, ok, want)
		}
	}
}

// containsAddr reports whether the 128-bit address falls in the
// interval (test-side check mirroring the runtime lookup).
func containsAddr(iv interval, addr [16]byte) bool {
	if addrLess(addr, iv.start) {
		return false
	}
	end := addPow2(iv.start, iv.bits)
	return addrLess(addr, end)
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	trie := &trieNode{}
	for _, spec := range []struct {
		cidr  string
		value uint32
	}{
		{"1.0.0.0/24", 10}, {"1.0.1.0/24", 10}, {"2001:db8::/32", 64512},
	} {
		p := netip.MustParsePrefix(spec.cidr)
		raw, bits := prefix128(p)
		trie.insert(raw, bits, spec.value)
	}
	intervals := trie.flatten()
	encoded := encodeTable(intervals)

	// Decode with the same layout the runtime reader uses.
	reader := bytes.NewReader(encoded)
	count, err := binary.ReadUvarint(reader)
	if err != nil {
		t.Fatal(err)
	}
	if int(count) != len(intervals) {
		t.Fatalf("decoded count = %d, want %d", count, len(intervals))
	}
	var prevEnd [16]byte
	for i := 0; i < int(count); i++ {
		hi, _ := binary.ReadUvarint(reader)
		lo, _ := binary.ReadUvarint(reader)
		bits, _ := binary.ReadUvarint(reader)
		value, _ := binary.ReadUvarint(reader)
		start := add128(prevEnd, hi, lo)
		want := intervals[i]
		if start != want.start || int(bits) != want.bits || uint32(value) != want.value {
			t.Fatalf("interval %d decoded = {%v %d %d}, want {%v %d %d}",
				i, start, bits, value, want.start, want.bits, want.value)
		}
		prevEnd = addPow2(start, int(bits))
	}
}

// add128 adds a 128-bit value given as hi/lo halves (test-side mirror
// of the runtime decoder).
func add128(base [16]byte, hi, lo uint64) [16]byte {
	out := base
	addLo := lo
	addHi := hi
	for i := 15; i >= 0; i-- {
		var add byte
		if i >= 8 {
			add = byte(addHi >> (8 * uint(i-8)))
		} else {
			add = byte(addLo >> (8 * uint(i)))
		}
		v := int(out[i]) + int(add)
		out[i] = byte(v)
		if v > 255 && i > 0 {
			out[i-1]++
		}
	}
	return out
}

func TestMRTASPathOrigin(t *testing.T) {
	// 4-byte AS_PATH: one AS_SEQUENCE segment with two ASNs
	// (type 1, count 2, 64500, 64512) — origin is the last ASN.
	path4 := []byte{1, 2, 0x00, 0x00, 0xfc, 0x04, 0x00, 0x00, 0xfc, 0x00}
	if asn, ok := lastASFromPath(path4); !ok || asn != 64512 {
		t.Fatalf("4-byte origin = %d (ok %v), want 64512", asn, ok)
	}
	// Legacy 2-byte encoding: type 2 (AS_SET), count 3.
	path2 := []byte{2, 3, 0xfc, 0x04, 0xfc, 0x05, 0xfc, 0xff}
	if asn, ok := lastASFromPath(path2); !ok || asn != 0xfcff {
		t.Fatalf("2-byte origin = %d (ok %v), want %d", asn, ok, 0xfcff)
	}
}

func TestMRTASPathAttributeScan(t *testing.T) {
	// ORIGIN (type 1, 1 byte), then AS_PATH (type 2, flags 0x40,
	// extended length) with one 4-byte ASN.
	attrs := []byte{
		0x40, 0x01, 0x01, 0x01,
		0x50, 0x02, 0x00, 0x06, 0x01, 0x01, 0x00, 0x00, 0xfd, 0xe8,
	}
	if asn, ok := originFromASPath(attrs); !ok || asn != 65000 {
		t.Fatalf("origin = %d (ok %v), want 65000", asn, ok)
	}
}
