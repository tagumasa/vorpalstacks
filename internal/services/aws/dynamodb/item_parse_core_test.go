package dynamodb

import (
	"reflect"
	"strings"
	"testing"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// The AttributeValue union carries exactly one member; the decoder enforces
// that structurally instead of resolving multi-member shapes by member order
// (and rejects a NULL member whose value is not true).

func TestParseAttributeValueRejectsMultiMemberShapes(t *testing.T) {
	cases := []struct {
		name  string
		shape map[string]interface{}
	}{
		{"string and number", map[string]interface{}{"S": "x", "N": "1"}},
		{"bool and null", map[string]interface{}{"BOOL": true, "NULL": true}},
		{"string and map", map[string]interface{}{"S": "x", "M": map[string]interface{}{}}},
		{"list and string set", map[string]interface{}{"L": []interface{}{}, "SS": []interface{}{"a"}}},
	}
	for _, tc := range cases {
		_, err := parseAttributeValue(tc.shape)
		if err == nil {
			t.Fatalf("%s: expected union violation error, got value", tc.name)
		}
		if !strings.Contains(err.Error(), "more than one datatypes set") {
			t.Fatalf("%s: expected multi-member message, got %q", tc.name, err.Error())
		}
	}
}

func TestParseAttributeValueRejectsNullFalse(t *testing.T) {
	_, err := parseAttributeValue(map[string]interface{}{"NULL": false})
	if err == nil {
		t.Fatal("NULL:false: expected rejection, got value")
	}
	if !strings.Contains(err.Error(), "Null value must be true") {
		t.Fatalf("NULL:false: expected null-true message, got %q", err.Error())
	}

	// A non-boolean NULL payload is the same degenerate shape.
	if _, err := parseAttributeValue(map[string]interface{}{"NULL": "yes"}); err == nil {
		t.Fatal("NULL:string: expected rejection, got value")
	}

	// NULL:true remains the one valid NULL form, and a bare JSON null still
	// decodes as a NULL value.
	av, err := parseAttributeValue(map[string]interface{}{"NULL": true})
	if err != nil || av == nil || av.NULL == nil || !*av.NULL {
		t.Fatalf("NULL:true: expected null value, got (%v, %v)", av, err)
	}
	if av, err := parseAttributeValue(nil); err != nil || av == nil || av.NULL == nil || !*av.NULL {
		t.Fatalf("JSON null: expected null value, got (%v, %v)", av, err)
	}
}

func TestParseAttributeValueSingleMembersStillParse(t *testing.T) {
	av, err := parseAttributeValue(map[string]interface{}{"S": "x"})
	if err != nil || av == nil || av.S == nil || *av.S != "x" {
		t.Fatalf("S member: got (%v, %v)", av, err)
	}
	av, err = parseAttributeValue(map[string]interface{}{"BOOL": false})
	if err != nil || av == nil || av.BOOL == nil || *av.BOOL {
		t.Fatalf("BOOL:false member: got (%v, %v)", av, err)
	}

	// A structural violation nested inside a container value rejects the
	// whole value, not just the offending element.
	_, err = parseAttributeValue(map[string]interface{}{
		"M": map[string]interface{}{
			"inner": map[string]interface{}{"S": "x", "N": "1"},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "more than one datatypes set") {
		t.Fatalf("nested violation: expected propagated union error, got %v", err)
	}
}

func TestParseItemPropagatesValueErrors(t *testing.T) {
	_, err := parseItem(map[string]interface{}{
		"pk":  map[string]interface{}{"S": "a"},
		"bad": map[string]interface{}{"S": "x", "N": "1"},
	})
	if err == nil || !strings.Contains(err.Error(), "more than one datatypes set") {
		t.Fatalf("parseItem: expected propagated union error, got %v", err)
	}

	item, err := parseItem(map[string]interface{}{
		"pk": map[string]interface{}{"S": "a"},
	})
	if err != nil || item == nil || item["pk"] == nil {
		t.Fatalf("parseItem valid: got (%v, %v)", item, err)
	}
}

// TestRequestParserAcceptsTheWireCodecOutput pins the agreement between
// the two wire-shape halves: every value the response codec renders, the
// request-side parser accepts back as the same typed value — one decoder
// renders, the other validates, and neither may drift from the shared
// wire shape.
func TestRequestParserAcceptsTheWireCodecOutput(t *testing.T) {
	s, n := "str", "42.5"
	b := true
	inner := "inner"
	values := []*dbstore.AttributeValue{
		{S: &s},
		{N: &n},
		{B: []byte{0x00, 0x01, 0xff}},
		{BOOL: &b},
		dbstore.NullValue(),
		{SS: []string{"a", "b"}},
		{NS: []string{"1", "2"}},
		{BS: [][]byte{{0x01}, {0x02}}},
		{M: map[string]*dbstore.AttributeValue{"nested": {S: &inner}}},
		{L: []*dbstore.AttributeValue{{S: &inner}, {N: &n}}},
	}
	for _, av := range values {
		wire := dbstore.BuildAttributeValueWire(av)
		parsed, err := parseAttributeValue(wire)
		if err != nil {
			t.Fatalf("request parser rejected the codec output of %#v: %v", av, err)
		}
		if !reflect.DeepEqual(parsed, av) {
			t.Fatalf("codec output re-parsed differently:\n got %#v\nwant %#v", parsed, av)
		}
	}
}
