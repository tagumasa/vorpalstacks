package dynamodb

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vorpalstacks/pkg/sqlparser"
)

func wireParams(vals ...map[string]interface{}) *partiQLParams {
	raw := make([]interface{}, 0, len(vals))
	for _, v := range vals {
		raw = append(raw, v)
	}
	return &partiQLParams{Parameters: raw}
}

// TestParseInsertStatementWithParams_BindsAllTypes pins that every
// AttributeValue union member survives parameter binding with its type
// intact: the parser delegates to the item-path wire parser instead of
// stringifying unknown members.
func TestParseInsertStatementWithParams_BindsAllTypes(t *testing.T) {
	stmt := `INSERT INTO "T" VALUE {'a': ?, 'b': ?, 'c': ?, 'd': ?, 'e': ?, 'f': ?, 'g': ?, 'h': ?, 'i': ?, 'j': ?}`
	params := wireParams(
		map[string]interface{}{"S": "str"},
		map[string]interface{}{"N": "42"},
		map[string]interface{}{"B": base64.StdEncoding.EncodeToString([]byte{0x01, 0x02})},
		map[string]interface{}{"BOOL": true},
		map[string]interface{}{"NULL": true},
		map[string]interface{}{"SS": []interface{}{"x", "y"}},
		map[string]interface{}{"NS": []interface{}{"1", "2"}},
		map[string]interface{}{"BS": []interface{}{base64.StdEncoding.EncodeToString([]byte{0x03})}},
		map[string]interface{}{"L": []interface{}{map[string]interface{}{"N": "7"}}},
		map[string]interface{}{"M": map[string]interface{}{"inner": map[string]interface{}{"S": "v"}}},
	)

	tableName, item, err := parseInsertStatementWithParams(stmt, params)
	require.NoError(t, err)
	assert.Equal(t, "T", tableName)

	assert.Equal(t, "str", *item["a"].S)
	assert.Equal(t, "42", *item["b"].N)
	assert.Equal(t, []byte{0x01, 0x02}, item["c"].B)
	assert.True(t, *item["d"].BOOL)
	assert.True(t, *item["e"].NULL)
	assert.Equal(t, []string{"x", "y"}, item["f"].SS)
	assert.Equal(t, []string{"1", "2"}, item["g"].NS)
	assert.Equal(t, [][]byte{{0x03}}, item["h"].BS)
	require.Len(t, item["i"].L, 1)
	assert.Equal(t, "7", *item["i"].L[0].N)
	assert.Equal(t, "v", *item["j"].M["inner"].S)
}

// TestParseInsertStatementWithParams_NestedPlaceholders pins that
// placeholders inside nested lists and objects are resolved, not stored as
// placeholder text.
func TestParseInsertStatementWithParams_NestedPlaceholders(t *testing.T) {
	stmt := `INSERT INTO "T" VALUE {'m': {'x': ?}, 'l': [?], 'll': [[?]]}`
	params := wireParams(
		map[string]interface{}{"N": "1"},
		map[string]interface{}{"S": "two"},
		map[string]interface{}{"BOOL": false},
	)

	_, item, err := parseInsertStatementWithParams(stmt, params)
	require.NoError(t, err)

	assert.Equal(t, "1", *item["m"].M["x"].N)
	require.Len(t, item["l"].L, 1)
	assert.Equal(t, "two", *item["l"].L[0].S)
	require.Len(t, item["ll"].L, 1)
	require.Len(t, item["ll"].L[0].L, 1)
	assert.False(t, *item["ll"].L[0].L[0].BOOL)
}

// TestResolvePlaceholder pins the parameter-binding error contract: one
// bound parameter per placeholder, and the parameter must be a valid
// AttributeValue — anything else is a ValidationException.
func TestResolvePlaceholder(t *testing.T) {
	valid := wireParams(map[string]interface{}{"S": "v"})

	v, err := resolvePlaceholder([]byte(":v1"), valid)
	require.NoError(t, err)
	assert.Equal(t, "v", *v.S)

	_, err = resolvePlaceholder([]byte(":v1"), nil)
	assert.ErrorIs(t, err, ErrInvalidParameter)

	_, err = resolvePlaceholder([]byte(":v2"), valid)
	assert.ErrorIs(t, err, ErrInvalidParameter, "out-of-range placeholder index")

	_, err = resolvePlaceholder([]byte(":vn"), valid)
	assert.ErrorIs(t, err, ErrInvalidParameter, "non-numeric placeholder index")

	_, err = resolvePlaceholder([]byte(":v1"), wireParams(map[string]interface{}{"S": 1}))
	assert.ErrorIs(t, err, ErrInvalidParameter, "malformed parameter value")
}

// TestParseInsertStatementWithParams_UnresolvedPlaceholderErrors pins that
// an INSERT carrying a placeholder with no bound parameter fails instead of
// persisting the placeholder text.
func TestParseInsertStatementWithParams_UnresolvedPlaceholderErrors(t *testing.T) {
	_, _, err := parseInsertStatementWithParams(`INSERT INTO "T" VALUE {'id': ?}`, &partiQLParams{})
	assert.ErrorIs(t, err, ErrInvalidParameter)

	_, item, err := parseInsertStatementWithParams(`INSERT INTO "T" VALUE {'id': 'literal'}`, nil)
	require.NoError(t, err)
	assert.Equal(t, "literal", *item["id"].S, "literal statements need no parameters")
}

// TestExprToAttributeValueLiterals pins the literal-only materialiser: every
// literal kind keeps its DynamoDB type.
func TestExprToAttributeValueLiterals(t *testing.T) {
	_, item, err := parseInsertStatementWithParams(
		`INSERT INTO "T" VALUE {'s': 'a', 'n': 5, 'f': 5.5, 'b': true, 'z': null, 'l': [1, 'x'], 'm': {'k': 'v'}}`, nil)
	require.NoError(t, err)

	assert.Equal(t, "a", *item["s"].S)
	assert.Equal(t, "5", *item["n"].N)
	assert.Equal(t, "5.5", *item["f"].N)
	assert.True(t, *item["b"].BOOL)
	assert.True(t, *item["z"].NULL)
	require.Len(t, item["l"].L, 2)
	assert.Equal(t, "1", *item["l"].L[0].N)
	assert.Equal(t, "x", *item["l"].L[1].S)
	assert.Equal(t, "v", *item["m"].M["k"].S)
}

// TestPartiQLSetLiteralMaterialises pins the set literal: a homogeneous
// member list materialises as SS or NS, while mixed kinds and repeated
// members are validation errors — never a silently collapsed or retyped
// set.
func TestPartiQLSetLiteralMaterialises(t *testing.T) {
	_, item, err := parseInsertStatementWithParams(
		`INSERT INTO "T" VALUE {'ss': <<'a', 'b'>>, 'ns': <<1, 2>>}`, nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, item["ss"].SS)
	assert.Equal(t, []string{"1", "2"}, item["ns"].NS)

	for _, stmt := range []string{
		`INSERT INTO "T" VALUE {'s': <<'a', 1>>}`,
		`INSERT INTO "T" VALUE {'s': <<'a', 'a'>>}`,
		`INSERT INTO "T" VALUE {'s': <<1, 1.0>>}`,
		`INSERT INTO "T" VALUE {'s': <<true>>}`,
	} {
		_, _, err := parseInsertStatementWithParams(stmt, nil)
		assert.ErrorIs(t, err, ErrInvalidParameter, stmt)
	}
}

// TestPartiQLNumberLiteralsHeldToNumberContract pins that number literals
// follow the same DynamoDB Number contract as wire AttributeValues: the
// documented magnitude range and the 38-significant-digit limit. The key
// case is a sub-1E-130 magnitude such as 1e-500 — the storage-key number
// encoder renders it identically to zero, so accepting it would let a
// literal overwrite the zero-valued item's key.
func TestPartiQLNumberLiteralsHeldToNumberContract(t *testing.T) {
	_, _, err := parseInsertStatementWithParams(`INSERT INTO "T" VALUE {'id': 1e-500}`, nil)
	assert.ErrorIs(t, err, ErrInvalidParameter)

	_, _, err = parseInsertStatementWithParams(`INSERT INTO "T" VALUE {'id': 1e-131}`, nil)
	assert.ErrorIs(t, err, ErrInvalidParameter)

	_, _, err = parseInsertStatementWithParams(`INSERT INTO "T" VALUE {'id': 1e126}`, nil)
	assert.ErrorIs(t, err, ErrInvalidParameter, "the documented top magnitude is below 1E+126")

	_, _, err = parseInsertStatementWithParams(
		`INSERT INTO "T" VALUE {'id': 1234567890123456789012345678901234567890}`, nil)
	assert.ErrorIs(t, err, ErrInvalidParameter, "39 significant digits exceed the 38-digit precision")

	_, _, err = parseInsertStatementWithParams(`INSERT INTO "T" VALUE {'m': {'n': [1e-500]}}`, nil)
	assert.ErrorIs(t, err, ErrInvalidParameter, "nested literals follow the same contract")

	_, item, err := parseInsertStatementWithParams(
		`INSERT INTO "T" VALUE {'a': 0, 'b': -5, 'c': +7, 'd': 1e-130, 'e': 9.9999999999999999999999999999999999999e125}`, nil)
	require.NoError(t, err, "zero, signed literals, and both range extremes are valid Numbers")
	assert.Equal(t, "0", *item["a"].N)
	assert.Equal(t, "-5", *item["b"].N, "a sign-prefixed literal is a Number, not NULL")
	assert.Equal(t, "7", *item["c"].N)
	assert.Equal(t, "1e-130", *item["d"].N)
	assert.Equal(t, "9.9999999999999999999999999999999999999e125", *item["e"].N)
}

// TestPartiQLValueMaterialiserRejectsNonValues pins the corruption-class
// removal: an expression that is not a value this grammar defines — a
// column reference, a function call, an operator combination — is a
// validation error, never a materialised NULL stored as the attribute.
func TestPartiQLValueMaterialiserRejectsNonValues(t *testing.T) {
	for _, stmt := range []string{
		`INSERT INTO "T" VALUE {'a': other}`,        // column reference
		`INSERT INTO "T" VALUE {'a': upper('x')}`,   // function call
		`INSERT INTO "T" VALUE {'a': 1 + 2}`,        // operator combination
		`INSERT INTO "T" VALUE {'a': [other]}`,      // nested inside a list
		`INSERT INTO "T" VALUE {'a': {'b': upper}}`, // nested inside an object
	} {
		_, _, err := parseInsertStatementWithParams(stmt, nil)
		assert.ErrorIs(t, err, ErrInvalidParameter, stmt)
	}
}

// TestPreparePartiQLStatementNumbersInOrder pins the placeholder rewrite
// that gives every `?` its whole-statement index: numbering follows text
// order across clauses, literals and quoted identifiers never contribute,
// and a second pass over the output is a no-op.
func TestPreparePartiQLStatementNumbersInOrder(t *testing.T) {
	rewritten, count := preparePartiQLStatement(`INSERT INTO "T" VALUE {'a': ?, 'b': ?}`)
	assert.Equal(t, 2, count)
	assert.Equal(t, `INSERT INTO "T" VALUE {'a': :v1, 'b': :v2}`, rewritten)

	rewritten, count = preparePartiQLStatement(`SELECT * FROM "T" WHERE s = 'why?' AND n = ?`)
	assert.Equal(t, 1, count)
	assert.Equal(t, `SELECT * FROM "T" WHERE s = 'why?' AND n = :v1`, rewritten)

	again, count := preparePartiQLStatement(rewritten)
	assert.Equal(t, 0, count)
	assert.Equal(t, rewritten, again)
}

func TestValidateParameterCount(t *testing.T) {
	assert.NoError(t, validateParameterCount(2, wireParams(
		map[string]interface{}{"S": "a"},
		map[string]interface{}{"S": "b"},
	)))
	assert.Error(t, validateParameterCount(3, wireParams(
		map[string]interface{}{"S": "a"},
		map[string]interface{}{"S": "b"},
	)))
	assert.Error(t, validateParameterCount(1, wireParams(
		map[string]interface{}{"S": "a"},
		map[string]interface{}{"S": "b"},
	)))
}

// placeholderValues walks an expression and returns the bound placeholder
// references it carries, in walk order.
func placeholderValues(t *testing.T, expr interface{}) []string {
	t.Helper()
	var found []string
	err := sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		if val, ok := node.(*sqlparser.SQLVal); ok && val.Type == sqlparser.ValArg {
			found = append(found, string(val.Val))
		}
		return true, nil
	}, expr.(sqlparser.SQLNode))
	require.NoError(t, err)
	return found
}

// TestManualUpdateNumberingSpansSegments pins the whole-statement
// numbering across the manual UPDATE path: the ADD-clause value and the
// WHERE equality carry distinct, statement-ordered indices even though the
// path re-parses each segment through its own fake statement.
func TestManualUpdateNumberingSpansSegments(t *testing.T) {
	stmt, placeholders := preparePartiQLStatement(`UPDATE "T" ADD n ? WHERE id = ?`)
	require.Equal(t, 2, placeholders)

	_, clauses, whereExpr, parseErr := parseUpdateStatement(stmt)
	require.NoError(t, parseErr)
	require.Len(t, clauses.addAssignments, 1)
	require.NotNil(t, whereExpr)

	assert.Equal(t, []string{":v1"}, placeholderValues(t, clauses.addAssignments[0].value))
	assert.Equal(t, []string{":v2"}, placeholderValues(t, whereExpr))
}
