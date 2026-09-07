package dynamodb

import (
	"reflect"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// allTypesAttributeValue carries one attribute value of every wire type.
func allTypesAttributeValue() *AttributeValue {
	s, n := "str", "42.5"
	b := true
	null := true
	inner := "inner"
	return &AttributeValue{
		S:    &s,
		N:    &n,
		B:    []byte{0x00, 0x01, 0xff},
		BOOL: &b,
		NULL: &null,
		SS:   []string{"a", "b"},
		NS:   []string{"1", "2"},
		BS:   [][]byte{{0x01}, {0x02}},
		M: map[string]*AttributeValue{
			"nested": {S: &inner},
		},
		L: []*AttributeValue{{S: &inner}, {N: &n}},
	}
}

// The wire codec must round-trip every attribute-value type exactly; it is
// the single definition of the wire shape for both the service responses
// and the stored stream-record images.
func TestAttributeValueWireCodecRoundTrip(t *testing.T) {
	av := allTypesAttributeValue()
	wire := BuildAttributeValueWire(av)
	parsed := ParseAttributeValueWire(wire)
	if !reflect.DeepEqual(av, parsed) {
		t.Fatalf("round-trip mismatch:\n got %#v\nwant %#v", parsed, av)
	}

	if got := BuildAttributeValueWire(nil); got != nil {
		t.Fatalf("nil value must render nil, got %v", got)
	}
	if got := ParseItemWire(nil); got != nil {
		t.Fatalf("nil map must parse nil, got %v", got)
	}
	if got := BuildItemWire(nil); !reflect.DeepEqual(got, map[string]interface{}{}) {
		t.Fatalf("nil item must render an empty map, got %v", got)
	}
}

// Stream records persist as protobuf and read back with every field —
// including the wire-shaped images and the service identity — intact.
func TestStreamRecordProtoRoundTrip(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	store := NewStreamStore(st, "123456789012", "us-east-1")

	keys := BuildItemWire(map[string]*AttributeValue{"id": {S: ptr("k1")}})
	newImage := BuildItemWire(map[string]*AttributeValue{"id": {S: ptr("k1")}, "n": {N: ptr("7")}})
	original, err := store.AddRecord("Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
		StreamEventInsert, keys, newImage, nil, TTLServiceIdentity)
	if err != nil {
		t.Fatalf("add record: %v", err)
	}

	loaded, next, err := store.GetRecords("Tbl", 0, 10)
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	if len(loaded) != 1 || next != 1 {
		t.Fatalf("loaded = %v, next = %d", loaded, next)
	}
	rec := loaded[0]
	if rec.EventID != original.EventID || rec.EventName != StreamEventInsert ||
		rec.EventVersion != original.EventVersion || rec.EventSourceARN != testStreamArn {
		t.Fatalf("record header mismatch: %+v", rec)
	}
	if rec.Dynamodb.SequenceNumber != original.Dynamodb.SequenceNumber ||
		rec.Dynamodb.StreamViewType != "NEW_AND_OLD_IMAGES" ||
		rec.Dynamodb.SizeBytes != original.Dynamodb.SizeBytes {
		t.Fatalf("record data mismatch: %+v", rec.Dynamodb)
	}
	if rec.Dynamodb.Keys["id"].(map[string]interface{})["S"] != "k1" {
		t.Fatalf("keys mismatch: %v", rec.Dynamodb.Keys)
	}
	if rec.Dynamodb.NewImage["n"].(map[string]interface{})["N"] != "7" {
		t.Fatalf("new image mismatch: %v", rec.Dynamodb.NewImage)
	}
	if rec.Dynamodb.OldImage != nil {
		t.Fatalf("absent old image must stay nil, got %v", rec.Dynamodb.OldImage)
	}
	if rec.UserIdentity == nil || rec.UserIdentity.Type != "Service" ||
		rec.UserIdentity.PrincipalID != "dynamodb.amazonaws.com" {
		t.Fatalf("user identity mismatch: %+v", rec.UserIdentity)
	}

	latest, err := store.GetLatestSequence("Tbl")
	if err != nil || latest != 1 {
		t.Fatalf("latest sequence = %d, %v", latest, err)
	}
}

func ptr(s string) *string { return &s }
