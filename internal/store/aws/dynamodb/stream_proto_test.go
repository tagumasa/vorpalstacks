package dynamodb

import (
	"context"
	"encoding/base64"
	"reflect"
	"testing"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
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

// The NULL member's pointer presence — not the flag's value — is the type
// discriminator on every output path: the wire parser normalises any
// {"NULL": ...} value to true and the persisted NullValue enum cannot carry
// false, so a hand-built false flag still renders, converts and answers
// IsNull as a null.
func TestNullMemberConvertsOnPresence(t *testing.T) {
	flag := false
	av := &AttributeValue{NULL: &flag}

	wire := BuildAttributeValueWire(av)
	if got, ok := wire["NULL"].(bool); !ok || !got {
		t.Fatalf("wire render of a NULL-typed value = %v, want NULL: true", wire)
	}
	parsed := ParseAttributeValueWire(wire)
	if parsed.NULL == nil || !*parsed.NULL {
		t.Fatalf("parse of the rendered wire form = %+v, want NULL true", parsed)
	}

	pbAV := attributeValueToProto(av)
	if _, isNull := pbAV.Value.(*pb.AttributeValue_Null); !isNull {
		t.Fatalf("proto conversion of a NULL-typed value = %#v, want the Null member", pbAV)
	}
	back := protoToAttributeValue(pbAV)
	if back.NULL == nil || !*back.NULL {
		t.Fatalf("proto round-trip of a NULL-typed value = %+v, want NULL true", back)
	}

	if !av.IsNull() {
		t.Fatal("IsNull must answer true on pointer presence alone")
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
	var original *StreamRecord
	if err := st.Update(context.Background(), func(txn storage.Transaction) error {
		r, err := store.AddRecordTxn(txn, "Tbl", testStreamArn, "NEW_AND_OLD_IMAGES",
			StreamEventInsert, keys, newImage, nil, TTLServiceIdentity)
		original = r
		return err
	}); err != nil {
		t.Fatalf("add record: %v", err)
	}

	loaded, next, err := store.GetRecords("Tbl", testStreamArn, 0, 10)
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

	latest, err := store.GetLatestSequenceForStream("Tbl", "")
	if err != nil || latest != 1 {
		t.Fatalf("latest sequence = %d, %v", latest, err)
	}
}

func ptr(s string) *string { return &s }

// TestStringSetsDecodeFromBothWireShapes pins the set normaliser's
// contract: SS and NS decode from the JSON-decoded []interface{} shape as
// they do from the builder's []string shape, and BS keeps both — a
// document decoded straight from JSON must not lose its string or number
// sets while its binary set survives.
func TestStringSetsDecodeFromBothWireShapes(t *testing.T) {
	fromJSON := ParseAttributeValueWire(map[string]interface{}{
		"SS": []interface{}{"a", "b"},
		"NS": []interface{}{"1", "2"},
		"BS": []interface{}{base64.StdEncoding.EncodeToString([]byte{0x01})},
	})
	if len(fromJSON.SS) != 2 || fromJSON.SS[0] != "a" || fromJSON.SS[1] != "b" {
		t.Fatalf("SS from []interface{} = %v", fromJSON.SS)
	}
	if len(fromJSON.NS) != 2 || fromJSON.NS[0] != "1" || fromJSON.NS[1] != "2" {
		t.Fatalf("NS from []interface{} = %v", fromJSON.NS)
	}
	if len(fromJSON.BS) != 1 || string(fromJSON.BS[0]) != "\x01" {
		t.Fatalf("BS from []interface{} = %v", fromJSON.BS)
	}

	roundTrip := ParseAttributeValueWire(BuildAttributeValueWire(&AttributeValue{
		SS: []string{"a", "b"},
		NS: []string{"1", "2"},
		BS: [][]byte{{0x01}},
	}))
	if len(roundTrip.SS) != 2 || len(roundTrip.NS) != 2 || len(roundTrip.BS) != 1 {
		t.Fatalf("builder round-trip lost a set: %+v", roundTrip)
	}
}
