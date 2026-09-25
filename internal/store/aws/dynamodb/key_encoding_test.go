package dynamodb

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func avS(s string) *AttributeValue { return &AttributeValue{S: &s} }
func avN(n string) *AttributeValue { return &AttributeValue{N: &n} }
func avB(b []byte) *AttributeValue { return &AttributeValue{B: b} }

func keySchemaTable(pkName, skName string) *Table {
	schema := []*KeySchemaElement{{AttributeName: pkName, KeyType: KeyTypeHash}}
	if skName != "" {
		schema = append(schema, &KeySchemaElement{AttributeName: skName, KeyType: KeyTypeRange})
	}
	return &Table{Name: "t", KeySchema: schema}
}

// A NUL byte inside a binary partition value must not shift the component
// boundary: the two item layouts below produce the same raw-joined key under
// separator joining and must remain distinct.
func TestEncodeItemKeyNulInjectionDistinct(t *testing.T) {
	table := keySchemaTable("pk", "sk")
	k1 := EncodeItemKey("t", map[string]*AttributeValue{
		"pk": avS("a\x00b"), "sk": avS("c")}, table)
	k2 := EncodeItemKey("t", map[string]*AttributeValue{
		"pk": avS("a"), "sk": avS("b\x00c")}, table)
	if k1 == "" || k2 == "" {
		t.Fatalf("expected non-degenerate keys, got %q and %q", k1, k2)
	}
	if k1 == k2 {
		t.Fatalf("NUL injection: distinct key maps encoded identically: %q", k1)
	}
}

// DynamoDB numbers carry 38 significant digits; the encoding must keep
// adjacent 38-digit values distinct (a float64 mantissa collapses them).
func TestEncodeKeyValueNumberPrecision38(t *testing.T) {
	a := avN("12345678901234567890123456789012345678")
	b := avN("12345678901234567890123456789012345679")
	ea, eb := EncodeKeyValue(a), EncodeKeyValue(b)
	if ea == eb {
		t.Fatalf("adjacent 38-digit numbers encoded identically: %q", ea)
	}
	if !(ea < eb) {
		t.Fatalf("expected enc(a) < enc(b), got %q vs %q", ea, eb)
	}
}

// The deep-exponent band of the documented number range carries 38
// significant digits at exponents down to -130, placing the lowest digit as
// far as decimal place 167; the rendering must keep such values distinct at
// every exponent, in both signs.
func TestEncodeKeyValueDeepExponentPrecision(t *testing.T) {
	pairs := [][2]string{
		{"1.0000000000000000000000000000000000001E-100", "1.0000000000000000000000000000000000002E-100"},
		{"1.0000000000000000000000000000000000001E-130", "1.0000000000000000000000000000000000002E-130"},
		{"-1.0000000000000000000000000000000000002E-100", "-1.0000000000000000000000000000000000001E-100"},
		{"1E-100", "1.0000000000000000000000000000000000001E-100"},
	}
	for _, p := range pairs {
		ea, eb := EncodeKeyValue(avN(p[0])), EncodeKeyValue(avN(p[1]))
		if ea == eb {
			t.Fatalf("distinct values %q and %q encoded identically: %q", p[0], p[1], ea)
		}
		if !(ea < eb) {
			t.Fatalf("expected enc(%q) < enc(%q), got %q vs %q", p[0], p[1], ea, eb)
		}
	}
}

// Lexicographic order of encoded values must equal DynamoDB value order per
// type family: numbers numerically (negatives descending into ascending byte
// order), strings and binaries by unsigned byte order.
func TestEncodeKeyValueOrderPreserving(t *testing.T) {
	numbers := []string{
		"-9.9999999999999999999999999999999999999E+125",
		"-1000000", "-3.5", "-0.55", "-0.5", "-1E-130",
		"0", "1E-130", "0.05", "0.5", "0.55", "1", "9", "10", "50",
		"1E+125", "9.9999999999999999999999999999999999999E+125",
	}
	prev := ""
	for i, n := range numbers {
		cur := EncodeKeyValue(avN(n))
		if cur == "" {
			t.Fatalf("number %q encoded empty", n)
		}
		if i > 0 && !(prev < cur) {
			t.Fatalf("number order broken at %q: %q !< %q", n, prev, cur)
		}
		prev = cur
	}

	strs := []string{"", "a", "a\x00", "a\x00b", "ab", "b", "z", "\xff"}
	prev = ""
	for i, s := range strs {
		cur := EncodeKeyValue(avS(s))
		if i > 0 && !(prev < cur) {
			t.Fatalf("string order broken at %q: %q !< %q", s, prev, cur)
		}
		prev = cur
	}

	bins := [][]byte{{}, {0x00}, {0x00, 0x01}, {0x00, 0x02}, {0x01}, {0x01, 0x00}, {0xff}}
	prev = ""
	for i, b := range bins {
		cur := EncodeKeyValue(avB(b))
		if i > 0 && !(prev < cur) {
			t.Fatalf("binary order broken at %x: %q !< %q", b, prev, cur)
		}
		prev = cur
	}
}

// No partition component encoding may be a prefix of another, so a partition
// scan prefix cannot cross into another partition.
func TestEncodeKeyValuePrefixFree(t *testing.T) {
	values := []*AttributeValue{
		avS("a"), avS("ab"), avS("abc"), avS("a\x00b"), avS("b"),
		avB([]byte{0x00}), avB([]byte{0x00, 0x00}), avB([]byte{0x00, 0x01}),
		avN("1"), avN("1.5"), avN("15"), avN("-1"), avN("0"), avN("0.5"),
	}
	encs := make([]string, 0, len(values))
	seen := map[string]string{}
	for _, v := range values {
		e := EncodeKeyValue(v)
		if e == "" {
			t.Fatalf("value %+v encoded empty", v)
		}
		if dup, ok := seen[e]; ok {
			t.Fatalf("distinct values %+v and %s encoded identically: %q", v, dup, e)
		}
		seen[e] = fmt.Sprintf("%+v", v)
		encs = append(encs, e)
	}
	sort.Strings(encs)
	for i := 0; i+1 < len(encs); i++ {
		if len(encs[i]) <= len(encs[i+1]) && encs[i+1][:len(encs[i])] == encs[i] {
			t.Fatalf("encoding %q is a prefix of %q", encs[i], encs[i+1])
		}
	}
}

// The unparseable-number fallback is a component payload like any other:
// every escape-relevant byte in it is escaped, so the terminator — and the
// separator it doubles as — occurs exactly once, as the component's last
// byte. The invariant the header asserts holds on this path too, not only
// on the validated ones.
func TestEncodeKeyValueUnparseableNumberFallbackEscapes(t *testing.T) {
	for _, n := range []string{"1\x002", "1\x012", "1\x01\x002"} {
		enc := EncodeKeyValue(avN(n))
		if strings.Count(enc, "\x00") != 1 || !strings.HasSuffix(enc, "\x00") {
			t.Fatalf("fallback %q left its payload unescaped: %q", n, enc)
		}
	}
}

// Empty S and B payloads signal a missing key attribute; the number zero is
// a legal key value and must encode non-empty.
func TestEncodeKeyValueEmptyPayloads(t *testing.T) {
	if got := EncodeKeyValue(avS("")); got != "" {
		t.Fatalf("empty S encoded as %q, want \"\"", got)
	}
	if got := EncodeKeyValue(avB(nil)); got != "" {
		t.Fatalf("empty B encoded as %q, want \"\"", got)
	}
	if got := EncodeKeyValue(avN("0")); got == "" {
		t.Fatal("number zero encoded empty, want non-empty")
	}
}

// Scientific notation, plus signs, and leading/trailing zeros all normalise
// to the same numeric value and must encode identically.
func TestEncodeKeyValueNumberNormalisation(t *testing.T) {
	canonical := EncodeKeyValue(avN("1.5"))
	for _, form := range []string{"1.5", "1.50", "+1.5", "0.15E1", "15E-1"} {
		if got := EncodeKeyValue(avN(form)); got != canonical {
			t.Fatalf("form %q encoded %q, want canonical %q", form, got, canonical)
		}
	}
}
