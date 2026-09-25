package dynamodb

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func ptrStr(s string) *string { return &s }

func ptrBool(b bool) *bool { return &b }

// evalCondition compiles and evaluates one expression against the item —
// the per-case form of the two-step contract every plane applies: parse
// and validate once per request, evaluate per item.
func evalCondition(t *testing.T, item *dbstore.Item, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) bool {
	t.Helper()
	compiled, err := compileConditionExpression(expr, names, values)
	require.NoError(t, err)
	return compiled.matches(item)
}

func TestConditionNOTCaseInsensitive(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"status": {S: ptrStr("active")},
		},
	}
	values := map[string]*dbstore.AttributeValue{
		":v": {S: ptrStr("active")},
	}

	cases := []struct {
		name string
		expr string
		want bool
	}{
		{"NOT uppercase", "NOT #s = :v", false},
		{"not lowercase", "not #s = :v", false},
		{"Not mixed case", "Not #s = :v", false},
		{"NoT mixed case", "NoT #s = :v", false},
		{"double NOT", "NOT NOT #s = :v", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			names := map[string]string{"#s": "status"}
			got := evalCondition(t, item, tc.expr, names, values)
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestConditionValidationRejectsUndefinedSubstitutions pins the
// parse-stage contract: an expression that references a substitution
// neither map defines is a ValidationException before any item is read —
// the read planes answered HTTP 200 with a silently filtered result for
// the same requests.
func TestConditionValidationRejectsUndefinedSubstitutions(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"name": {S: ptrStr("test")},
		},
	}
	values := map[string]*dbstore.AttributeValue{
		":v": {S: ptrStr("active")},
	}

	for _, expr := range []string{
		"#nonexistent = :v", // alias absent from the names map
		"name = :missing",   // placeholder absent from the values map
		"name BETWEEN :missing AND :v",
		"name IN (:v, :missing)",
	} {
		_, err := compileConditionExpression(expr, nil, values)
		require.Error(t, err, "%q: expected a validation error", expr)
		var apiErr *APIError
		require.True(t, errors.As(err, &apiErr), "%q: expected an APIError, got %T", expr, err)
		assert.Equal(t, "com.amazon.coral.validate#ValidationException", apiErr.Code, "%q: wrong error code", expr)
	}

	// A spec with no expression and no tree carries no condition at all.
	compiled, err := conditionSpec{}.compile()
	require.NoError(t, err)
	assert.Nil(t, compiled)

	// A well-formed expression still evaluates.
	assert.True(t, evalCondition(t, item, "#n = :t", map[string]string{"#n": "name"},
		map[string]*dbstore.AttributeValue{":t": {S: ptrStr("test")}}))
}

// TestConditionParseRejectsMalformedGrammar pins the syntax gate: an
// operator outside the expression grammar, a literal where the grammar
// has no literal, or a function call with the wrong shape is a
// ValidationException — never a condition that silently evaluates false.
func TestConditionParseRejectsMalformedGrammar(t *testing.T) {
	values := map[string]*dbstore.AttributeValue{
		":p": {S: ptrStr("prefix")},
	}
	for _, expr := range []string{
		"a NE :v",            // legacy ComparisonOperator name
		"a != :v",            // SQL-style inequality
		"s = 'x'",            // string literal: not an operand of this grammar
		"a = ",               // missing right-hand side
		"(a = :v",            // unbalanced parenthesis
		"attribute_exists",   // bare function name
		"attribute_exists()", // empty argument list
		"begins_with(a)",     // one-operand call
		"contains(a, :p, :p)",
		"a EQ :v",
	} {
		_, err := compileConditionExpression(expr, nil, values)
		require.Error(t, err, "%q: expected a parse error", expr)
	}
}

func TestConditionANDOR(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"age":    {N: ptrStr("25")},
			"status": {S: ptrStr("active")},
		},
	}
	names := map[string]string{
		"#a": "age",
		"#s": "status",
	}
	values := map[string]*dbstore.AttributeValue{
		":age":    {N: ptrStr("25")},
		":status": {S: ptrStr("active")},
		":old":    {N: ptrStr("30")},
	}

	t.Run("AND both true", func(t *testing.T) {
		assert.True(t, evalCondition(t, item, "#a = :age AND #s = :status", names, values))
	})
	t.Run("AND one false", func(t *testing.T) {
		assert.False(t, evalCondition(t, item, "#a = :age AND #s = :old", names, values))
	})
	t.Run("OR one true", func(t *testing.T) {
		assert.True(t, evalCondition(t, item, "#a = :old OR #s = :status", names, values))
	})
	t.Run("OR both false", func(t *testing.T) {
		assert.False(t, evalCondition(t, item, "#a = :old OR #s = :age", names, values))
	})
	t.Run("keyword case folding", func(t *testing.T) {
		assert.True(t, evalCondition(t, item, "#a = :age and #s = :status", names, values))
		assert.True(t, evalCondition(t, item, "#a = :old or #s = :status", names, values))
	})
	t.Run("precedence AND before OR", func(t *testing.T) {
		assert.True(t, evalCondition(t, item, "#a = :age AND #s = :old OR #s = :status", names, values))
	})
	t.Run("parenthesised grouping", func(t *testing.T) {
		assert.False(t, evalCondition(t, item, "(#a = :age OR #s = :status) AND #a = :old", names, values))
		assert.True(t, evalCondition(t, item, "(#a = :old OR #s = :old) OR #s = :status", names, values))
	})
}

// TestConditionUnicodeAroundKeywords pins the lexer's byte-wise path
// accumulation: a multi-byte rune whose uppercase form changes byte
// length next to a keyword still splits on the keyword, never inside the
// rune.
func TestConditionUnicodeAroundKeywords(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"s": {S: ptrStr("ı")},
			"n": {N: ptrStr("5")},
		},
	}
	values := map[string]*dbstore.AttributeValue{
		":s":    {S: ptrStr("ı")},
		":n":    {N: ptrStr("7")},
		":five": {N: ptrStr("5")},
	}

	assert.True(t, evalCondition(t, item, "s = :s OR n = :n", nil, values))
	assert.True(t, evalCondition(t, item, "s = :s AND n = :five", nil, values))

	// A string literal stays a parse error even beside multi-byte runes.
	_, err := compileConditionExpression("s = 'x' OR n = :n", nil, values)
	require.Error(t, err)
}

// TestConditionFunctionArityAndOperandGrammar pins the function grammar:
// the value operand follows the operand grammar — a bare token is a
// document path, so begins_with(a, b) compares against the value AT b,
// not the literal string "b"; a path that addresses nothing makes the
// condition false.
func TestConditionFunctionArityAndOperandGrammar(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"name": {S: ptrStr("prefix-middle")},
			"pfx":  {S: ptrStr("prefix")},
		},
	}
	values := map[string]*dbstore.AttributeValue{
		":p": {S: ptrStr("prefix")},
	}

	assert.True(t, evalCondition(t, item, "begins_with(name, pfx)", nil, nil))
	assert.False(t, evalCondition(t, item, "begins_with(name, absent)", nil, nil))
	assert.True(t, evalCondition(t, item, "begins_with(name, :p)", nil, values))

	// The path and the operand must be distinct — the same resolved path
	// on both sides of contains/begins_with is rejected, including through
	// aliases that resolve to the same path; distinct paths keep parsing.
	_, err := compileConditionExpression("contains(name, name)", nil, nil)
	require.Error(t, err, "contains(name, name): expected a distinctness rejection")
	_, err = compileConditionExpression("begins_with(#n, #n)", map[string]string{"#n": "name"}, nil)
	require.Error(t, err, "begins_with(#n, #n): expected a distinctness rejection")
	_, err = compileConditionExpression("begins_with(name, pfx)", nil, nil)
	require.NoError(t, err, "begins_with(name, pfx): distinct paths must parse")
}

// TestConditionBinaryOperandFunctions pins the binary branches of the
// operand-typed functions: begins_with prefix-tests bytes, contains
// matches a contiguous byte subsequence, NOT contains negates it, and a
// string operand against a binary attribute never matches.
func TestConditionBinaryOperandFunctions(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"blob": {B: []byte("prefix-middle")},
		},
	}
	values := map[string]*dbstore.AttributeValue{
		":pfx":  {B: []byte("prefix")},
		":mid":  {B: []byte("ix-mid")},
		":miss": {B: []byte("absent")},
		":str":  {S: ptrStr("prefix")},
	}

	assert.True(t, evalCondition(t, item, "begins_with(blob, :pfx)", nil, values))
	assert.False(t, evalCondition(t, item, "begins_with(blob, :mid)", nil, values))
	assert.False(t, evalCondition(t, item, "begins_with(blob, :str)", nil, values))

	assert.True(t, evalCondition(t, item, "contains(blob, :mid)", nil, values))
	assert.False(t, evalCondition(t, item, "contains(blob, :miss)", nil, values))
	assert.False(t, evalCondition(t, item, "contains(blob, :str)", nil, values))

	assert.False(t, evalCondition(t, item, "NOT contains(blob, :mid)", nil, values))
	assert.True(t, evalCondition(t, item, "NOT contains(blob, :miss)", nil, values))
}

// TestConditionAttributeTypeOperand pins attribute_type's operand
// contract: the placeholder carries a String naming a type, a defined
// placeholder of any other type is rejected instead of silently matching
// nothing, and a type name that does not describe the attribute answers
// false.
func TestConditionAttributeTypeOperand(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"name": {S: ptrStr("x")},
		},
	}
	values := map[string]*dbstore.AttributeValue{
		":s": {S: ptrStr("S")},
		":n": {N: ptrStr("7")},
	}

	assert.True(t, evalCondition(t, item, "attribute_type(name, :s)", nil, values))
	assert.False(t, evalCondition(t, item, "attribute_type(name, :no)", nil, map[string]*dbstore.AttributeValue{
		":no": {S: ptrStr("N")},
	}))
	_, err := compileConditionExpression("attribute_type(name, :n)", nil, values)
	require.Error(t, err, "attribute_type(name, :n): expected a non-String operand rejection")
}

// TestConditionSizeAndPathOperands pins the operand model: size() is an
// operand wherever an operand stands — BETWEEN bounds, IN candidates, and
// the right-hand side of a comparator — and a document path resolves on
// either side of a comparison.
func TestConditionSizeAndPathOperands(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"ss":    {SS: []string{"a", "b"}},
			"count": {N: ptrStr("2")},
			"copy":  {SS: []string{"a", "b"}},
		},
	}
	values := map[string]*dbstore.AttributeValue{
		":one":   {N: ptrStr("1")},
		":two":   {N: ptrStr("2")},
		":three": {N: ptrStr("3")},
		":big":   {N: ptrStr("9")},
	}

	assert.True(t, evalCondition(t, item, "size(ss) BETWEEN :one AND :three", nil, values))
	assert.False(t, evalCondition(t, item, "size(ss) BETWEEN :three AND :big", nil, values))
	assert.True(t, evalCondition(t, item, "size(ss) IN (:one, :two)", nil, values))
	assert.False(t, evalCondition(t, item, "size(ss) IN (:one, :three)", nil, values))
	assert.True(t, evalCondition(t, item, "size(ss) <> :big", nil, values))
	assert.True(t, evalCondition(t, item, "size(ss) = size(copy)", nil, nil))
	// count BETWEEN size(ss) AND :three: 2 within [2, 3].
	assert.True(t, evalCondition(t, item, "count BETWEEN size(ss) AND :three", nil, values))
}

// TestMissingAttributeComparisonsAreFalse pins the missing-attribute rule
// with no = :null carve-out: every comparator on an attribute absent from
// the item is false — an absent attribute is not a NULL-valued one, which
// is what the cached operators page separates explicitly — while an
// attribute that exists with the NULL type still equals :null.
func TestMissingAttributeComparisonsAreFalse(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"nul": {NULL: ptrBool(true)},
		},
	}
	values := map[string]*dbstore.AttributeValue{
		":null": {NULL: ptrBool(true)},
		":s":    {S: ptrStr("x")},
	}

	assert.False(t, evalCondition(t, item, "absent = :null", nil, values))
	assert.False(t, evalCondition(t, item, "absent <> :s", nil, values))
	assert.True(t, evalCondition(t, item, "nul = :null", nil, values))
}
