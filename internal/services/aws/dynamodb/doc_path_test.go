package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// TestParseDocPath pins the shared document-path grammar: name and
// bracket-index segmentation, bracket validation, and the requirement that
// a path begin with an attribute name.
func TestParseDocPath(t *testing.T) {
	type want struct {
		names   []string
		indexes []int
	}
	cases := map[string]want{
		"a":          {names: []string{"a"}},
		"a.b":        {names: []string{"a", "b"}},
		"a[0].b":     {names: []string{"a", "b"}, indexes: []int{0}},
		"a[0][1]":    {names: []string{"a"}, indexes: []int{0, 1}},
		"a[10].b[2]": {names: []string{"a", "b"}, indexes: []int{10, 2}},
		"a..b":       {names: []string{"a", "b"}},
	}
	for path, w := range cases {
		parts, err := parseDocPath(path)
		require.NoError(t, err, path)
		var gotNames []string
		var gotIndexes []int
		for _, p := range parts {
			if p.isIndex {
				gotIndexes = append(gotIndexes, p.index)
			} else {
				gotNames = append(gotNames, p.name)
			}
		}
		assert.Equal(t, w.names, gotNames, "names for %q", path)
		assert.Equal(t, w.indexes, gotIndexes, "indexes for %q", path)
	}

	for _, path := range []string{"a[x]", "a[-1]", "a[]", "a[0.5]", "[0].a", "[1]", "a[0", "a[0][5"} {
		if _, err := parseDocPath(path); err == nil {
			t.Errorf("parseDocPath(%q): expected an error", path)
		}
	}
}

// TestGetDocPathValue pins the read walk over Map and List containers.
func TestGetDocPathValue(t *testing.T) {
	attrs := map[string]*dbstore.AttributeValue{
		"top": dbstore.StringValue("x"),
		"doc": dbstore.MapValue(map[string]*dbstore.AttributeValue{
			"list": dbstore.ListValue([]*dbstore.AttributeValue{
				dbstore.StringValue("zero"),
				dbstore.MapValue(map[string]*dbstore.AttributeValue{
					"name": dbstore.StringValue("nested"),
				}),
			}),
		}),
	}

	parts := func(path string) []docPathPart {
		parsed, err := parseDocPath(path)
		require.NoError(t, err)
		return parsed
	}

	if v := getDocPathValue(attrs, parts("top")); v == nil || v.S == nil || *v.S != "x" {
		t.Errorf("top-level lookup failed: %+v", v)
	}
	if v := getDocPathValue(attrs, parts("doc.list[0]")); v == nil || v.S == nil || *v.S != "zero" {
		t.Errorf("list index lookup failed: %+v", v)
	}
	if v := getDocPathValue(attrs, parts("doc.list[1].name")); v == nil || v.S == nil || *v.S != "nested" {
		t.Errorf("nested map lookup failed: %+v", v)
	}
	for _, path := range []string{"missing", "doc.missing", "doc.list[2]", "doc.top"} {
		// "doc.top" addresses the top-level string through a map container;
		// the walk stops where the container type does not match.
		if v := getDocPathValue(attrs, parts(path)); v != nil {
			t.Errorf("getDocPathValue(%q): expected nil, got %+v", path, v)
		}
	}
	if v := getDocPathValue(attrs, nil); v != nil {
		t.Errorf("empty parts: expected nil, got %+v", v)
	}
}

// TestNestedValueMutation pins the write walk: single-part assignment,
// intermediate Map creation, index bounds and type-mismatch rejection.
func TestNestedValueMutation(t *testing.T) {
	attrs := map[string]*dbstore.AttributeValue{}
	require.NoError(t, setNestedValue(attrs, "flat", dbstore.StringValue("v")))
	assert.Equal(t, "v", *attrs["flat"].S)

	require.NoError(t, setNestedValue(attrs, "a.b.c", dbstore.StringValue("deep")))
	assert.Equal(t, "deep", *attrs["a"].M["b"].M["c"].S)

	list := dbstore.ListValue([]*dbstore.AttributeValue{dbstore.StringValue("x")})
	attrs["arr"] = list
	require.NoError(t, setNestedValue(attrs, "arr[0]", dbstore.StringValue("y")))
	assert.Equal(t, "y", *list.L[0].S)

	if err := setNestedValue(attrs, "arr[5]", dbstore.StringValue("oob")); err == nil {
		t.Error("out-of-range index: expected an error")
	}
	if err := setNestedValue(attrs, "flat[0]", dbstore.StringValue("idx-on-scalar")); err == nil {
		t.Error("index into a non-List top level: expected an error")
	}
	if err := setNestedValue(attrs, "flat.nested[0]", dbstore.StringValue("idx-on-scalar")); err == nil {
		t.Error("index into a created Map: expected an error")
	}

	require.NoError(t, removeNestedValue(attrs, "a.b"))
	if _, ok := attrs["a"].M["b"]; ok {
		t.Error("removeNestedValue left the removed key behind")
	}
	require.NoError(t, removeNestedValue(attrs, "arr[0]"))
	assert.Empty(t, list.L)
	require.NoError(t, removeNestedValue(attrs, "missing.path"))
	require.NoError(t, removeNestedValue(attrs, "flat.deep.missing"))

	// A leading bracket index has no attribute to address.
	if err := setNestedValue(attrs, "[0]", dbstore.StringValue("no")); err == nil {
		t.Error("leading bracket index: expected an error")
	}
}

// TestConditionRejectsMalformedDocPath pins the condition-side strictness:
// a malformed path is a request error, not a missing-attribute answer —
// writes surface it through the condition checker instead of evaluating
// the path's prefix.
func TestConditionRejectsMalformedDocPath(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"a": dbstore.MapValue(map[string]*dbstore.AttributeValue{
				"b": dbstore.StringValue("v"),
			}),
		},
	}
	for _, expr := range []string{"a[x] = :v", "a[-1] = :v", "a[0 = :v", "attribute_exists(a[x])", "attribute_not_exists([0].a)"} {
		if _, err := evaluateConditionExpr(item, expr, nil, map[string]*dbstore.AttributeValue{
			":v": dbstore.StringValue("v"),
		}); err == nil {
			t.Errorf("evaluateConditionExpr(%q): expected an error", expr)
		}
	}

	// Well-formed nested paths keep evaluating against the item.
	met, err := evaluateConditionExpr(item, "a.b = :v", nil, map[string]*dbstore.AttributeValue{
		":v": dbstore.StringValue("v"),
	})
	require.NoError(t, err)
	assert.True(t, met)
}
