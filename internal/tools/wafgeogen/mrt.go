// MRT TABLE_DUMP_V2 RIB parsing (RFC 6396). The generator extracts one
// prefix-to-origin-ASN mapping per announced prefix: the origin is the
// rightmost autonomous system of the RIB entry's AS_PATH.
package main

import (
	"compress/bzip2"
	"encoding/binary"
	"fmt"
	"io"
	"net/netip"
)

const (
	mrtTypeTableDumpV2 = 13
	// TABLE_DUMP_V2 subtypes per RFC 6396: 1 is the peer index table,
	// 2-5 are the RIB record flavours.
	subtypePeerIndex   = 1
	subtypeRIBIPv4     = 2
	subtypeRIBIPv4Path = 3
	subtypeRIBIPv6     = 4
	subtypeRIBIPv6Path = 5
)

// parseMRT reads a bz2-compressed MRT RIB dump and appends one record
// per prefix (the first RIB entry seen for the prefix wins; the
// collector's own view is authoritative enough for a derived table).
func parseMRT(r io.Reader, out map[netip.Prefix]uint32) error {
	br := bzip2.NewReader(r)
	header := make([]byte, 12)
	for {
		if _, err := io.ReadFull(br, header); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("read mrt header: %w", err)
		}
		mrtType := binary.BigEndian.Uint16(header[4:6])
		subtype := binary.BigEndian.Uint16(header[6:8])
		length := int(binary.BigEndian.Uint32(header[8:12]))
		body := make([]byte, length)
		if _, err := io.ReadFull(br, body); err != nil {
			return fmt.Errorf("read mrt body: %w", err)
		}
		if mrtType != mrtTypeTableDumpV2 {
			continue
		}
		if err := parseRIBEntry(subtype, body, out); err != nil {
			// A malformed record must not abort the whole dump; the
			// remaining records still form a usable table.
			continue
		}
	}
}

// parseRIBEntry handles the RIB_IPV4/6_UNICAST(_ADDPATH) subtypes.
// Per RFC 6396 these records carry no AFI/SAFI (only the GENERIC
// subtypes do): the body is sequence(4), prefix length(1), prefix,
// entry count(2), then the RIB entries.
func parseRIBEntry(subtype uint16, body []byte, out map[netip.Prefix]uint32) error {
	if subtype < subtypeRIBIPv4 || subtype > subtypeRIBIPv6Path {
		return nil
	}
	addPath := subtype == subtypeRIBIPv4Path || subtype == subtypeRIBIPv6Path
	if len(body) < 7 {
		return fmt.Errorf("short rib header")
	}
	prefixLen := int(body[4])
	var addrLen int
	var isIPv6 bool
	switch subtype {
	case subtypeRIBIPv4, subtypeRIBIPv4Path:
		if prefixLen > 32 {
			return fmt.Errorf("ipv4 prefix length %d", prefixLen)
		}
		addrLen = (prefixLen + 7) / 8
		isIPv6 = false
	default:
		addrLen = 16
		isIPv6 = true
	}
	if prefixLen > 128 {
		return fmt.Errorf("prefix length %d", prefixLen)
	}
	offset := 5
	if len(body) < offset+addrLen {
		return fmt.Errorf("short prefix")
	}
	var prefix netip.Prefix
	if isIPv6 {
		var raw [16]byte
		copy(raw[:], body[offset:offset+addrLen])
		p, err := netip.AddrFrom16(raw).Prefix(prefixLen)
		if err != nil {
			return err
		}
		prefix = p
	} else {
		if prefixLen > 32 {
			return fmt.Errorf("ipv4 prefix length %d", prefixLen)
		}
		var raw [4]byte
		copy(raw[:], body[offset:offset+addrLen])
		p, err := netip.AddrFrom4(raw).Prefix(prefixLen)
		if err != nil {
			return err
		}
		prefix = p
	}
	offset += addrLen
	if len(body) < offset+2 {
		return fmt.Errorf("short entry count")
	}
	entryCount := int(binary.BigEndian.Uint16(body[offset : offset+2]))
	offset += 2
	for i := 0; i < entryCount; i++ {
		if addPath {
			offset += 4
		}
		if len(body) < offset+8 {
			return fmt.Errorf("short rib entry")
		}
		attrLen := int(binary.BigEndian.Uint16(body[offset+6 : offset+8]))
		offset += 8
		if len(body) < offset+attrLen {
			return fmt.Errorf("short attributes")
		}
		attrs := body[offset : offset+attrLen]
		offset += attrLen
		if _, seen := out[prefix]; seen {
			continue
		}
		if asn, ok := originFromASPath(attrs); ok {
			out[prefix] = asn
		}
	}
	return nil
}

// originFromASPath scans BGP path attributes for AS_PATH (type 2) and
// returns its rightmost ASN, handling both 2-byte and 4-byte encodings.
func originFromASPath(attrs []byte) (uint32, bool) {
	for len(attrs) >= 3 {
		flags := attrs[0]
		attrType := attrs[1]
		length := 0
		payload := 0
		if flags&0x10 != 0 {
			if len(attrs) < 4 {
				return 0, false
			}
			length = int(binary.BigEndian.Uint16(attrs[2:4]))
			payload = 4
		} else {
			length = int(attrs[2])
			payload = 3
		}
		if payload+length > len(attrs) {
			return 0, false
		}
		if attrType == 2 {
			return lastASFromPath(attrs[payload : payload+length])
		}
		attrs = attrs[payload+length:]
	}
	return 0, false
}

// lastASFromPath walks the AS_PATH segments and returns the final ASN.
// Segments tile the path exactly, so the remaining length distinguishes
// the 4-byte encoding (2+count*4) from the legacy 2-byte one.
func lastASFromPath(path []byte) (uint32, bool) {
	var last uint32
	found := false
	for len(path) >= 2 {
		count := int(path[1])
		if count == 0 {
			path = path[2:]
			continue
		}
		switch {
		case len(path) == 2+count*4:
			last = binary.BigEndian.Uint32(path[2+count*4-4 : 2+count*4])
			path = path[2+count*4:]
		case len(path) == 2+count*2:
			last = uint32(binary.BigEndian.Uint16(path[2+count*2-2 : 2+count*2]))
			path = path[2+count*2:]
		default:
			return last, found
		}
		found = true
	}
	return last, found
}
