package dynamodb

import (
	"strings"
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
		// A multi-byte name is one segment: the walk takes whole runes, so
		// the segment keeps the name it was given — never a byte-at-a-time
		// re-encoding of it.
		"m.年":    {names: []string{"m", "年"}},
		"年[0].c": {names: []string{"年", "c"}, indexes: []int{0}},
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

// resolvePath parses a literal document path (no expression attribute
// names) for the mutation tests.
func resolvePath(t *testing.T, path string) []docPathPart {
	t.Helper()
	parts, err := resolveDocPathParts(path, nil)
	require.NoError(t, err, path)
	return parts
}

// TestResolveDocPathParts pins the per-segment expression-attribute-name
// resolution: each "#alias" segment of a compound path substitutes on the
// parsed parts, an undefined alias is a strict-resolution error while the
// lenient read variant keeps the literal segment, and a whole-token alias
// standing for a name that itself contains '.' stays a single segment.
func TestResolveDocPathParts(t *testing.T) {
	names := map[string]string{"#pr": "m", "#k": "k", "#d": "dotted.name"}

	parts, err := resolveDocPathParts("#pr.#k", names)
	require.NoError(t, err)
	require.Len(t, parts, 2)
	assert.Equal(t, "m", parts[0].name)
	assert.Equal(t, "k", parts[1].name)

	parts, err = resolveDocPathParts("#pr.#k[1]", names)
	require.NoError(t, err)
	require.Len(t, parts, 3)
	assert.Equal(t, "m", parts[0].name)
	assert.Equal(t, "k", parts[1].name)
	assert.True(t, parts[2].isIndex)
	assert.Equal(t, 1, parts[2].index)

	// A whole-token alias resolves to one segment even when the target
	// name carries a separator: resolution on parsed parts, never on a
	// re-serialised string.
	parts, err = resolveDocPathParts("#d", names)
	require.NoError(t, err)
	require.Len(t, parts, 1)
	assert.Equal(t, "dotted.name", parts[0].name)

	// Strict resolution rejects an alias the names map does not define;
	// the lenient read variant keeps the literal segment so the read
	// reports the attribute as absent.
	if _, err := resolveDocPathParts("#pr.#missing", names); err == nil {
		t.Error("strict resolution: expected an undefined-alias error")
	}
	parts, err = resolveDocPathPartsLenient("#pr.#missing", names)
	require.NoError(t, err)
	require.Len(t, parts, 2)
	assert.Equal(t, "m", parts[0].name)
	assert.Equal(t, "#missing", parts[1].name)
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
// intermediate map traversal, final-segment list append and path-validity
// rejection: a nested SET descends through existing containers only — a
// missing or wrong-typed parent at any depth is the documented invalid
// document path, while a final list segment naming a nonexistent element
// appends at the end of the list.
func TestNestedValueMutation(t *testing.T) {
	attrs := map[string]*dbstore.AttributeValue{}
	require.NoError(t, setNestedValue(attrs, resolvePath(t, "flat"), dbstore.StringValue("v")))
	assert.Equal(t, "v", *attrs["flat"].S)

	// A nested SET requires its parents to exist: no implicit map creation.
	attrs["a"] = dbstore.MapValue(map[string]*dbstore.AttributeValue{
		"b": dbstore.MapValue(map[string]*dbstore.AttributeValue{}),
	})
	require.NoError(t, setNestedValue(attrs, resolvePath(t, "a.b.c"), dbstore.StringValue("deep")))
	assert.Equal(t, "deep", *attrs["a"].M["b"].M["c"].S)
	if err := setNestedValue(attrs, resolvePath(t, "missing.nested"), dbstore.StringValue("v")); err == nil || !strings.Contains(err.Error(), "invalid for update") {
		t.Errorf("missing parent: expected the invalid-document-path rejection, got %v", err)
	}
	if err := setNestedValue(attrs, resolvePath(t, "a.missing.deeper"), dbstore.StringValue("v")); err == nil || !strings.Contains(err.Error(), "invalid for update") {
		t.Errorf("missing deeper parent: expected the invalid-document-path rejection, got %v", err)
	}
	if err := setNestedValue(attrs, resolvePath(t, "flat.nested"), dbstore.StringValue("v")); err == nil || !strings.Contains(err.Error(), "invalid for update") {
		t.Errorf("scalar parent: expected the invalid-document-path rejection, got %v", err)
	}

	list := dbstore.ListValue([]*dbstore.AttributeValue{dbstore.StringValue("x")})
	attrs["arr"] = list
	require.NoError(t, setNestedValue(attrs, resolvePath(t, "arr[0]"), dbstore.StringValue("y")))
	assert.Equal(t, "y", *list.L[0].S)

	// A final segment naming a nonexistent element appends at the list end,
	// whatever out-of-range index the path carries.
	require.NoError(t, setNestedValue(attrs, resolvePath(t, "arr[5]"), dbstore.StringValue("oob")))
	require.Len(t, list.L, 2)
	assert.Equal(t, "oob", *list.L[1].S)
	// Traversal through a missing element keeps the invalid-document-path
	// rejection: only the final segment may append.
	if err := setNestedValue(attrs, resolvePath(t, "arr[5].name"), dbstore.StringValue("v")); err == nil {
		t.Error("traversal through a missing element: expected an error")
	}
	if err := setNestedValue(attrs, resolvePath(t, "flat[0]"), dbstore.StringValue("idx-on-scalar")); err == nil {
		t.Error("index into a non-List top level: expected an error")
	}
	if err := setNestedValue(attrs, resolvePath(t, "flat.nested[0]"), dbstore.StringValue("idx-on-scalar")); err == nil {
		t.Error("index into a created Map: expected an error")
	}

	require.NoError(t, removeNestedValue(attrs, resolvePath(t, "a.b")))
	if _, ok := attrs["a"].M["b"]; ok {
		t.Error("removeNestedValue left the removed key behind")
	}
	require.NoError(t, removeNestedValue(attrs, resolvePath(t, "arr[0]")))
	require.Len(t, list.L, 1)
	assert.Equal(t, "oob", *list.L[0].S)
	require.NoError(t, removeNestedValue(attrs, resolvePath(t, "missing.path")))
	require.NoError(t, removeNestedValue(attrs, resolvePath(t, "flat.deep.missing")))

	// A leading bracket index has no attribute to address: the rejection
	// now surfaces at path resolution, before the mutation walk runs.
	if _, err := resolveDocPathParts("[0]", nil); err == nil {
		t.Error("leading bracket index: expected an error")
	}
}

// TestConditionRejectsMalformedDocPath pins the condition-side strictness:
// a malformed path is a request error at parse time, not a
// missing-attribute answer — every plane that carries a condition
// expression rejects it before any evaluation runs.
func TestConditionRejectsMalformedDocPath(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"a": dbstore.MapValue(map[string]*dbstore.AttributeValue{
				"b": dbstore.StringValue("v"),
			}),
		},
	}
	values := map[string]*dbstore.AttributeValue{
		":v": dbstore.StringValue("v"),
	}
	for _, expr := range []string{"a[x] = :v", "a[-1] = :v", "a[0 = :v", "attribute_exists(a[x])", "attribute_not_exists([0].a)"} {
		if _, err := compileConditionExpression(expr, nil, values); err == nil {
			t.Errorf("compileConditionExpression(%q): expected an error", expr)
		}
	}

	// Well-formed nested paths keep evaluating against the item.
	compiled, err := compileConditionExpression("a.b = :v", nil, values)
	require.NoError(t, err)
	assert.True(t, compiled.matches(item))
}
