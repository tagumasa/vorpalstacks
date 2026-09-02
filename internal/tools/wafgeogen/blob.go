// Blob layout shared with the runtime reader in the WAF inspection
// package:
//
//	magic "WGEOASN1"
//	uint32 generation date (YYYYMMDD)
//	uint32 country-code count; then count 2-byte codes
//	uint32 country-table byte length; then the encoded interval table
//	uint32 asn-table byte length; then the encoded interval table
//
// Interval tables are a uvarint interval count followed by per-interval
// uvarint records: start-delta hi, start-delta lo (the 128-bit
// distance from the previous interval's end), bits (size 1<<bits) and
// value (country-code index or ASN).
package main

import (
	"encoding/binary"
)

var blobMagic = [8]byte{'W', 'G', 'E', 'O', 'A', 'S', 'N', '1'}

func buildBlob(date string, countries []string, countryTable, asnTable []byte) []byte {
	var dateValue uint32
	if len(date) == 8 {
		var v uint64
		for i := 0; i < 8; i++ {
			v = v*10 + uint64(date[i]-'0')
		}
		dateValue = uint32(v)
	}
	out := make([]byte, 0, 32+len(countryTable)+len(asnTable))
	out = append(out, blobMagic[:]...)
	out = binary.BigEndian.AppendUint32(out, dateValue)
	out = binary.BigEndian.AppendUint32(out, uint32(len(countries)))
	for _, cc := range countries {
		out = append(out, cc[0], cc[1])
	}
	out = binary.BigEndian.AppendUint32(out, uint32(len(countryTable)))
	out = append(out, countryTable...)
	out = binary.BigEndian.AppendUint32(out, uint32(len(asnTable)))
	out = append(out, asnTable...)
	return out
}
