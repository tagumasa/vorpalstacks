package dynamodb

import (
	"fmt"
	"math/big"
	"strings"
)

// key_encoding.go is the single owner of the storage-key layout for item and
// index keys. Every component that carries an attribute value is encoded
// here; joining raw values with KeySep elsewhere is prohibited: binary key
// values may contain any byte (including NUL), which would shift component
// boundaries, and numbers need an exact order-preserving rendering (38
// significant digits — an intermediate float64 conversion loses precision).
//
// Component encoding properties (unit-pinned in key_encoding_test.go):
//   - injective: distinct values encode distinctly, so distinct items get
//     distinct storage keys;
//   - prefix-free: no component encoding is a prefix of another, so a
//     partition scan prefix cannot cross into another partition;
//   - order-preserving: lexicographic key order equals DynamoDB value order
//     (S by UTF-8 bytes, B as unsigned bytes, N numerically), which the
//     ascending lookahead cap and the reverse-page read rely on.

const (
	// keyTerm terminates ascending components (S, B, non-negative numbers).
	// It can never occur inside one because every 0x00/0x01 byte in the
	// payload is escaped.
	keyTerm = "\x00"
	// keyTermNeg terminates negative-number components. Negative encodings
	// are 9's-complemented, so a numerically smaller negative has MORE
	// digits and must sort BEFORE; that needs a terminator above every digit
	// byte, which keyTerm (below all digits) cannot provide.
	keyTermNeg = "\xff"
)

// numberSortFractionalDigits is the fractional-digit width of the sort
// rendering. Validated DynamoDB numbers carry at most 38 significant digits
// within the documented range 1E-130 to just under 1E+126, so the lowest
// significant digit of the smallest magnitudes (1.0…01E-130) sits at decimal
// place 130+37=167; rendering 168 fractional digits represents every
// validated value exactly, which keeps the numeric component injective.
const numberSortFractionalDigits = 168

// EncodeKeyValue renders one attribute value as a storage key component.
// Key attributes are S, N, or B; an empty S or B payload encodes as "" (the
// "missing key attribute" signal used by the item-key builders), while the
// number zero is a legal key value and encodes non-empty.
func EncodeKeyValue(av *AttributeValue) string {
	if av == nil {
		return ""
	}
	if av.S != nil {
		if *av.S == "" {
			return ""
		}
		return "S" + escapeKeyPayload(*av.S) + keyTerm
	}
	if av.N != nil {
		body := encodeNumberForSort(*av.N)
		if len(body) > 0 && body[0] == '0' {
			return "N" + body + keyTermNeg
		}
		return "N" + body + keyTerm
	}
	if av.B != nil {
		if len(av.B) == 0 {
			return ""
		}
		return "B" + escapeKeyPayload(string(av.B)) + keyTerm
	}
	return ""
}

// escapeKeyPayload escapes bytes so the terminator (0x00) can never appear
// inside a component: 0x00 becomes 0x01 0x01 and 0x01 becomes 0x01 0x02,
// every other byte passes through. The mapping preserves byte order (an
// escaped 0x00 sorts below an escaped 0x01, which sorts below 0x02) and is
// injective (every 0x01 in the output is the lead of an escape pair).
func escapeKeyPayload(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 0x00:
			b.WriteByte(0x01)
			b.WriteByte(0x01)
		case 0x01:
			b.WriteByte(0x01)
			b.WriteByte(0x02)
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// encodeNumberForSort renders a number as a fixed-shape decimal that sorts
// numerically: a sign half ("1" = non-negative, "0" = negative), a 3-digit
// integer-digit count, then the exact significant digits. Negative encodings
// are 9's-complemented so descending numeric order becomes ascending byte
// order. The rendering is derived from big.Rat via FloatString (exact
// integer scaling, no float mantissa) at numberSortFractionalDigits
// fractional digits, so all 38 significant digits of validated DynamoDB
// numbers survive at every exponent of the documented range.
func encodeNumberForSort(numStr string) string {
	rat := new(big.Rat)
	if _, ok := rat.SetString(numStr); numStr == "" || !ok {
		// Unparseable input never reaches storage keys through validated
		// paths; keep it outside the numeric ordering rather than panicking.
		return "X" + numStr
	}
	neg := rat.Sign() < 0
	plain := new(big.Rat).Abs(rat).FloatString(numberSortFractionalDigits)
	dot := strings.IndexByte(plain, '.')
	intPart, fracPart := plain[:dot], plain[dot+1:]
	fracPart = strings.TrimRight(fracPart, "0")
	head := fmt.Sprintf("%03d", len(intPart)) + intPart + fracPart
	if neg {
		return "0" + complementDigits(head)
	}
	return "1" + head
}

// complementDigits 9's-complements a decimal digit string, flipping its
// sort order.
func complementDigits(s string) string {
	b := []byte(s)
	for i := range b {
		b[i] = '9' - b[i] + '0'
	}
	return string(b)
}

// EncodeItemKey builds the storage key for an item: tableName, the encoded
// partition value, then the encoded sort value when the table has one. It
// returns "" when a schema key attribute is missing or empty, which callers
// treat as ErrInvalidKey.
func EncodeItemKey(tableName string, key map[string]*AttributeValue, table *Table) string {
	if table == nil {
		return ""
	}
	pkName, skName := schemaKeyNames(table.KeySchema)
	pkValue := EncodeKeyValue(key[pkName])
	if pkValue == "" {
		return ""
	}
	if skName != "" {
		skValue := EncodeKeyValue(key[skName])
		if skValue == "" {
			return ""
		}
		return tableName + KeySep + pkValue + KeySep + skValue
	}
	return tableName + KeySep + pkValue
}

// schemaKeyNames returns the hash and range key attribute names of a key
// schema ("" when absent); it serves table and index schemas alike.
func schemaKeyNames(schema []*KeySchemaElement) (string, string) {
	hashName, rangeName := "", ""
	for _, ks := range schema {
		switch ks.KeyType {
		case KeyTypeHash:
			hashName = ks.AttributeName
		case KeyTypeRange:
			rangeName = ks.AttributeName
		}
	}
	return hashName, rangeName
}
