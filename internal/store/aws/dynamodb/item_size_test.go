package dynamodb

import (
	"context"
	"testing"
)

func numAttr(n string) *AttributeValue { return &AttributeValue{N: &n} }

func boolAttr(b bool) *AttributeValue { return &AttributeValue{BOOL: &b} }

func nullAttr() *AttributeValue {
	t := true
	return &AttributeValue{NULL: &t}
}

// The documented item-size formula (Developer Guide): Numbers cost 1 byte
// per two significant digits plus 1 byte, Booleans/Nulls 1 byte, sets the
// sum of their element sizes (number elements carry the +1 without a name),
// Lists/Maps 3 bytes plus 1 byte per element.
func TestAttributeValueSizeDocumentedFormula(t *testing.T) {
	cases := []struct {
		name string
		av   *AttributeValue
		want int64
	}{
		{"string", strAttr("abc"), 3},
		{"binary", &AttributeValue{B: []byte{1, 2, 3, 4}}, 4},
		{"bool", boolAttr(true), 1},
		{"null", nullAttr(), 1},
		{"number three digits", numAttr("123"), 3},
		{"number zero", numAttr("0"), 1},
		{"number trailing zeros trimmed", numAttr("1000"), 2},
		{"number leading fractional zeros trimmed", numAttr("0.0001"), 2},
		{"number negative exponent form", numAttr("-2.5e10"), 2},
		{"string set", &AttributeValue{SS: []string{"a", "bc"}}, 3},
		{"number set", &AttributeValue{NS: []string{"1", "22"}}, 4},
		{"binary set", &AttributeValue{BS: [][]byte{{1}, {2, 3}}}, 3},
		{"list", &AttributeValue{L: []*AttributeValue{strAttr("a"), numAttr("1")}}, 3 + (1 + 1) + (1 + 2)},
		{"map", &AttributeValue{M: map[string]*AttributeValue{"k": strAttr("v")}}, 3 + (1 + 1 + 1)},
		{"nested document", &AttributeValue{M: map[string]*AttributeValue{
			"doc": {L: []*AttributeValue{strAttr("xy")}},
		}}, 3 + (1 + 3 + (3 + 1 + 2))},
	}
	for _, tc := range cases {
		if got := attributeValueSize(tc.av); got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestCalculateItemSizeSumsNameAndValue(t *testing.T) {
	item := map[string]*AttributeValue{"name": strAttr("abc")}
	if got := CalculateItemSize(item); got != 4+3 {
		t.Errorf("got %d, want 7", got)
	}
}

func TestCountSignificantDigitsTrimsZeros(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"0", 0},
		{"0.000", 0},
		{"1000", 1},
		{"1.23", 3},
		{"0.0001", 1},
		{"-2.5e10", 2},
	}
	for _, tc := range cases {
		if got := CountSignificantDigits(tc.in); got != tc.want {
			t.Errorf("%q: got %d, want %d", tc.in, got, tc.want)
		}
	}
}

// A TTL expiry must subtract exactly what the write path added, so a table
// whose every item expires reports zero size and zero item count. The write
// side below mirrors the service write cores (PutItem plus the counter
// updates); the worker derives its subtraction from the same single owner.
func TestTTLExpiryKeepsTableSizeSymmetric(t *testing.T) {
	store := newTablePersistStore(t)
	if err := store.Tables().SetTimeToLive("PersistTbl", &TimeToLiveSpecification{
		Enabled:       true,
		AttributeName: "exp",
	}); err != nil {
		t.Fatalf("set ttl: %v", err)
	}

	key := map[string]*AttributeValue{"id": strAttr("doc1")}
	attrs := map[string]*AttributeValue{
		"id":  strAttr("doc1"),
		"n":   numAttr("123.45"),
		"exp": numAttr("1"),
		"doc": {M: map[string]*AttributeValue{"l": {L: []*AttributeValue{strAttr("payload")}}}},
	}
	err := store.Update(context.Background(), func(txn *DynamoDBTxn) error {
		if err := txn.PutItem("PersistTbl", key, attrs); err != nil {
			return err
		}
		if err := txn.PutIndexEntries("PersistTbl", &Item{TableName: "PersistTbl", Key: key, Attributes: attrs}); err != nil {
			return err
		}
		if err := txn.UpdateItemCount("PersistTbl", 1); err != nil {
			return err
		}
		return txn.UpdateTableSize("PersistTbl", CalculateItemSize(attrs))
	})
	if err != nil {
		t.Fatalf("put: %v", err)
	}
	table, err := store.Tables().Get("PersistTbl")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if table.TableSizeBytes != CalculateItemSize(attrs) || table.ItemCount != 1 {
		t.Fatalf("after put: size=%d count=%d, want %d/1", table.TableSizeBytes, table.ItemCount, CalculateItemSize(attrs))
	}

	w := &ttlWorker{store: store, ctx: context.Background()}
	w.doTTLCleanup()

	table, err = store.Tables().Get("PersistTbl")
	if err != nil {
		t.Fatalf("get after cleanup: %v", err)
	}
	if table.TableSizeBytes != 0 || table.ItemCount != 0 {
		t.Fatalf("after TTL expiry: size=%d count=%d, want 0/0", table.TableSizeBytes, table.ItemCount)
	}
}
