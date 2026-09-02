// Package main implements the WAF geolocation and ASN table generator.
// This file parses RIR delegated-extended statistics files: each data
// record is pipe-separated as registry|cc|type|start|value|date|status
// with an optional trailing opaque-id. IPv4 records carry a start
// address and a host count; IPv6 records carry a start prefix and a
// prefix length. Summary records (type "*" markers) and records
// without a country code are ignored.
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/netip"
	"strconv"
	"strings"
)

// delegatedRecord is one address allocation: a prefix and the country
// code the registry delegated it to.
type delegatedRecord struct {
	Prefix      netip.Prefix
	CountryCode string
}

// parseDelegated reads one delegated-extended file and appends its
// IPv4 and IPv6 allocation records to the output slice.
func parseDelegated(r io.Reader, out []delegatedRecord) ([]delegatedRecord, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		fields := strings.Split(line, "|")
		if len(fields) < 7 {
			continue
		}
		cc, resType, start, value, status := fields[1], fields[2], fields[3], fields[4], fields[6]
		if cc == "" || cc == "*" || !(status == "allocated" || status == "assigned") {
			continue
		}
		switch resType {
		case "ipv4":
			addr, err := netip.ParseAddr(start)
			if err != nil || !addr.Is4() {
				continue
			}
			count, err := strconv.ParseUint(value, 10, 32)
			if err != nil || count == 0 {
				continue
			}
			for _, p := range alignedPrefixes(addr, count) {
				out = append(out, delegatedRecord{Prefix: p, CountryCode: strings.ToUpper(cc)})
			}
		case "ipv6":
			addr, err := netip.ParseAddr(start)
			if err != nil || !addr.Is6() {
				continue
			}
			length, err := strconv.Atoi(value)
			if err != nil || length <= 0 || length > 128 {
				continue
			}
			// 4-in-6 encoded allocations parse as Is4In6; normalise to
			// the 128-bit space the table uses.
			prefix, err := addr.Prefix(length)
			if err != nil {
				continue
			}
			out = append(out, delegatedRecord{Prefix: prefix, CountryCode: strings.ToUpper(cc)})
		}
	}
	if err := scanner.Err(); err != nil {
		return out, fmt.Errorf("scan delegated file: %w", err)
	}
	return out, nil
}

// alignedPrefixes decomposes an IPv4 start address plus address count
// into the aligned power-of-two prefixes that exactly cover the range
// (registry counts such as 2560 are not single powers of two, and the
// start address is not guaranteed to be aligned to any power of two).
// Block sizes are bounded by the alignment of the absolute address —
// the largest power of two dividing it — so no prefix ever extends
// beyond its boundary and the decomposition covers the range exactly,
// without over- or under-coverage.
func alignedPrefixes(start netip.Addr, count uint64) []netip.Prefix {
	base := start.As4()
	baseAbs := uint32(base[0])<<24 | uint32(base[1])<<16 | uint32(base[2])<<8 | uint32(base[3])
	var prefixes []netip.Prefix
	for consumed := uint64(0); consumed < count; {
		abs := baseAbs + uint32(consumed)
		if abs < baseAbs {
			// The allocation would run past the end of the IPv4 space;
			// registry data never does, so stop rather than wrap.
			return prefixes
		}
		// The largest block aligned to the current absolute address is
		// its lowest set bit (unbounded at address zero).
		align := uint64(1) << 32
		if abs != 0 {
			align = uint64(abs & -abs)
		}
		size := align
		if remaining := count - consumed; size > remaining {
			size = remaining
			for size&(size-1) != 0 {
				size &= size - 1
			}
		}
		bits := 32 - log2(size)
		addr := netip.AddrFrom4([4]byte{byte(abs >> 24), byte(abs >> 16), byte(abs >> 8), byte(abs)})
		if p, err := addr.Prefix(bits); err == nil {
			prefixes = append(prefixes, p)
		}
		consumed += size
	}
	return prefixes
}

func log2(v uint64) int {
	n := 0
	for v > 1 {
		v >>= 1
		n++
	}
	return n
}
