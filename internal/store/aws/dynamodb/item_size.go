package dynamodb

import (
	"strings"
)

// This file is the single owner of the DynamoDB item-size algorithm. The
// Developer Guide (DynamoDB item sizes and formats) defines it: an item's
// size is the sum of its attribute name lengths and value sizes, where a
// Number costs 1 byte per two significant digits plus 1 byte (leading and
// trailing zeros trimmed), a Boolean or Null costs 1 byte, a set costs the
// sum of its element sizes (number elements carry the same +1 byte, without
// a name), and a List or Map costs 3 bytes plus 1 byte per element on top of
// its nested contents. Both the item write paths (TableSizeBytes accounting)
// and the TTL worker's expiry accounting must derive the same figure, which
// is why the algorithm lives here once instead of per consumer.

// MaxItemSizeBytes is the documented maximum size of a single item.
const MaxItemSizeBytes = 400 * 1024

// ReadCapacityUnitBytes is the documented read-capacity granularity: one
// read capacity unit covers one strongly consistent read of an item up to
// this many bytes, with larger items rounding up to multiples of it.
const ReadCapacityUnitBytes int64 = 4096

// WriteCapacityUnitBytes is the documented write-capacity granularity: one
// write capacity unit covers one write of an item up to this many bytes,
// with larger items rounding up to multiples of it.
const WriteCapacityUnitBytes int64 = 1024

// CalculateItemSize returns the size in bytes of an item's attribute map
// under the documented formula: the sum of every attribute name length plus
// its value size.
func CalculateItemSize(attrs map[string]*AttributeValue) int64 {
	var size int64
	for name, av := range attrs {
		size += int64(len(name))
		size += attributeValueSize(av)
	}
	return size
}

// attributeValueSize returns the size in bytes of one attribute value,
// excluding the attribute name (the caller adds it).
func attributeValueSize(av *AttributeValue) int64 {
	if av == nil {
		return 0
	}
	if av.S != nil {
		return int64(len(*av.S))
	}
	if av.N != nil {
		return numberSize(*av.N)
	}
	if av.B != nil {
		return int64(len(av.B))
	}
	if av.BOOL != nil {
		return 1
	}
	if av.NULL != nil {
		return 1
	}
	if av.SS != nil {
		var size int64
		for _, s := range av.SS {
			size += int64(len(s))
		}
		return size
	}
	if av.NS != nil {
		var size int64
		for _, n := range av.NS {
			size += numberSize(n)
		}
		return size
	}
	if av.BS != nil {
		var size int64
		for _, b := range av.BS {
			size += int64(len(b))
		}
		return size
	}
	if av.M != nil {
		size := int64(3)
		for k, v := range av.M {
			size += 1 + int64(len(k)) + attributeValueSize(v)
		}
		return size
	}
	if av.L != nil {
		size := int64(3)
		for _, v := range av.L {
			size += 1 + attributeValueSize(v)
		}
		return size
	}
	return 0
}

// numberSize returns the size in bytes of a Number value: 1 byte per two
// significant digits plus 1 byte, with leading and trailing zeros trimmed
// from the mantissa before counting.
func numberSize(n string) int64 {
	return (int64(CountSignificantDigits(n))+1)/2 + 1
}

// CountSignificantDigits counts the mantissa's significant digits: leading
// and trailing zeros are trimmed, matching the DynamoDB number
// normalisation the 38-digit precision limit and the item-size formula both
// use. Zero alone carries no significant digits.
func CountSignificantDigits(n string) int {
	mantissa := n
	if idx := strings.IndexAny(mantissa, "eE"); idx != -1 {
		mantissa = mantissa[:idx]
	}
	mantissa = strings.TrimLeft(mantissa, "+-")
	mantissa = strings.ReplaceAll(mantissa, ".", "")
	mantissa = strings.TrimLeft(mantissa, "0")
	mantissa = strings.TrimRight(mantissa, "0")
	return len(mantissa)
}
