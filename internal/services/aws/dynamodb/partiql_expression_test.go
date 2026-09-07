package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// whereEval parses a WHERE clause and evaluates it against attrs, so the
// tests exercise the same parse→evaluate path the engines use.
func whereEval(t *testing.T, attrs map[string]*dbstore.AttributeValue, where string, params *partiQLParams) bool {
	t.Helper()
	_, expr := parseSelectStatement("SELECT * FROM \"t\" WHERE " + where)
	require.NotNil(t, expr, "WHERE clause must parse: %s", where)
	return evaluateExpr(attrs, expr, params)
}

func typedAttrs() map[string]*dbstore.AttributeValue {
	ten := "10"
	one := "1"
	return map[string]*dbstore.AttributeValue{
		"s1":       {S: &one},
		"sfoo":     {S: ptrStr("foo")},
		"n1":       {N: &one},
		"n10":      {N: &ten},
		"b":        {B: []byte{0x00, 0x01}},
		"flag":     {BOOL: proto.Bool(true)},
		"l":        {L: []*dbstore.AttributeValue{{S: ptrStr("x")}}},
		"m":        {M: map[string]*dbstore.AttributeValue{"k": {S: ptrStr("v")}}},
		"nullattr": {NULL: proto.Bool(true)},
	}
}

// TestPartiQLWhereEqualityIsTypeStrict pins that equality requires the same
// DynamoDB type with the same payload: a string "1" never equals the number
// 1, and a binary or document attribute never equals an empty string.
func TestPartiQLWhereEqualityIsTypeStrict(t *testing.T) {
	attrs := typedAttrs()

	cases := []struct {
		where string
		want  bool
	}{
		{"s1 = '1'", true},
		{"s1 = 1", false},
		{"n1 = 1", true},
		{"n1 = '1'", false},
		{"n10 = 10", true},
		{"b = ''", false},
		{"l = ''", false},
		{"m = ''", false},
		{"flag = true", true},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, whereEval(t, attrs, tc.where, nil), tc.where)
	}

	assert.True(t, whereEval(t, attrs, "s1 = ?",
		wireParams(map[string]interface{}{"S": "1"})))
	assert.False(t, whereEval(t, attrs, "s1 = ?",
		wireParams(map[string]interface{}{"N": "1"})))
	assert.True(t, whereEval(t, attrs, "b = ?",
		wireParams(map[string]interface{}{"B": "AAE="})))
}

// TestPartiQLWhereInequality pins <> semantics: a cross-type pair is not
// equal, so <> holds; same-type pairs compare their payloads.
func TestPartiQLWhereInequality(t *testing.T) {
	attrs := typedAttrs()

	assert.True(t, whereEval(t, attrs, "s1 <> 1", nil))
	assert.False(t, whereEval(t, attrs, "n1 <> 1", nil))
	assert.True(t, whereEval(t, attrs, "n1 <> 2", nil))
}

// TestPartiQLWhereOrdering pins ordering semantics: only same-type S/N/B
// operands are ordered, numbers compare numerically, and incomparable
// pairs never satisfy an ordering predicate.
func TestPartiQLWhereOrdering(t *testing.T) {
	attrs := typedAttrs()

	cases := []struct {
		where string
		want  bool
	}{
		{"n1 < 2", true},
		{"n1 < 1", false},
		// numeric ordering, not lexicographic: 10 > 9 even though "10" < "9"
		{"n10 > 9", true},
		{"n10 < 9", false},
		// cross-type ordering never matches
		{"n1 < '2'", false},
		{"s1 < 2", false},
		{"s1 <= 1", false},
		{"s1 >= 1", false},
		{"sfoo > 'bar'", true},
		// BOOL and documents are not orderable
		{"flag > false", false},
		{"l > ''", false},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, whereEval(t, attrs, tc.where, nil), tc.where)
	}
}

// TestPartiQLWhereInBetweenLikeNull pins the remaining predicates with
// typed operands.
func TestPartiQLWhereInBetweenLikeNull(t *testing.T) {
	attrs := typedAttrs()

	assert.True(t, whereEval(t, attrs, "n1 IN (1, 2)", nil))
	assert.False(t, whereEval(t, attrs, "n1 IN ('1')", nil))
	assert.True(t, whereEval(t, attrs, "n1 BETWEEN 0 AND 2", nil))
	assert.False(t, whereEval(t, attrs, "n1 BETWEEN '0' AND '2'", nil))
	assert.True(t, whereEval(t, attrs, "n1 NOT BETWEEN 2 AND 3", nil))
	assert.False(t, whereEval(t, attrs, "n1 NOT BETWEEN '2' AND '3'", nil))
	assert.True(t, whereEval(t, attrs, "sfoo LIKE 'f%'", nil))
	assert.False(t, whereEval(t, attrs, "sfoo LIKE 'b%'", nil))
	assert.False(t, whereEval(t, attrs, "n1 LIKE '1%'", nil))
	assert.True(t, whereEval(t, attrs, "nullattr IS NULL", nil))
	assert.True(t, whereEval(t, attrs, "s1 IS NOT NULL", nil))
	assert.False(t, whereEval(t, attrs, "s1 IS NULL", nil))
}

// TestPartiQLOrderByUsesTypedOrdering pins ORDER BY ordering: numbers sort
// numerically and incomparable operands keep the stable order.
func TestPartiQLOrderByUsesTypedOrdering(t *testing.T) {
	num := func(v string) *dbstore.Item {
		return &dbstore.Item{Attributes: map[string]*dbstore.AttributeValue{
			"n": {N: &v},
		}}
	}
	items := []*dbstore.Item{num("10"), num("9"), num("2")}

	sorted := sortItemsByOrderBy(items, &orderByClause{column: "n", direction: "ASC"})
	require.Len(t, sorted, 3)
	assert.Equal(t, "2", *sorted[0].Attributes["n"].N)
	assert.Equal(t, "9", *sorted[1].Attributes["n"].N)
	assert.Equal(t, "10", *sorted[2].Attributes["n"].N)

	desc := sortItemsByOrderBy(items, &orderByClause{column: "n", direction: "DESC"})
	require.Len(t, desc, 3)
	assert.Equal(t, "10", *desc[0].Attributes["n"].N)
	assert.Equal(t, "2", *desc[2].Attributes["n"].N)

	// An incomparable pair (S vs N on the same column) keeps stable order.
	mixed := []*dbstore.Item{
		{Attributes: map[string]*dbstore.AttributeValue{"a": {S: ptrStr("x")}}},
		{Attributes: map[string]*dbstore.AttributeValue{"a": {N: ptrStr("1")}}},
	}
	mixedSorted := sortItemsByOrderBy(mixed, &orderByClause{column: "a", direction: "ASC"})
	assert.Equal(t, "x", *mixedSorted[0].Attributes["a"].S)
	assert.Equal(t, "1", *mixedSorted[1].Attributes["a"].N)
}

// TestFilterExpressionOrderingIsTypeStrict pins the unified filter-engine
// ordering: cross-type scan/query filter comparisons never match, instead
// of comparing flattened strings.
func TestFilterExpressionOrderingIsTypeStrict(t *testing.T) {
	str := &dbstore.AttributeValue{S: ptrStr("9")}
	num := &dbstore.AttributeValue{N: ptrStr("1")}

	assert.False(t, compareAttributeValues(str, "<=", num))
	assert.False(t, compareAttributeValues(str, "<", num))
	assert.False(t, compareAttributeValues(num, ">=", str))
	assert.True(t, compareBetween(num, &dbstore.AttributeValue{N: ptrStr("0")}, &dbstore.AttributeValue{N: ptrStr("2")}))
	assert.False(t, compareBetween(num, &dbstore.AttributeValue{S: ptrStr("0")}, &dbstore.AttributeValue{N: ptrStr("2")}))
}
