package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func ptrStr(s string) *string { return &s }

func TestEvaluateConditionExpr_NOT_CaseInsensitive(t *testing.T) {
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
		{"Not mixed case (Bug D-9)", "Not #s = :v", false},
		{"NoT mixed case", "NoT #s = :v", false},
		{"double NOT", "NOT NOT #s = :v", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			names := map[string]string{"#s": "status"}
			got, err := evaluateConditionExpr(item, tc.expr, names, values)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEvaluateFilterExpression_MalformedReturnsFalse(t *testing.T) {
	item := &dbstore.Item{
		Attributes: map[string]*dbstore.AttributeValue{
			"name": {S: ptrStr("test")},
		},
	}

	t.Run("missing attribute name mapping returns false (Bug D-3)", func(t *testing.T) {
		values := map[string]*dbstore.AttributeValue{
			":v": {S: ptrStr("active")},
		}
		got := evaluateFilterExpression(item, "#nonexistent = :v", nil, values)
		assert.False(t, got, "expression referencing unmapped name should return false")
	})

	t.Run("empty expression returns true", func(t *testing.T) {
		got := evaluateFilterExpression(item, "", nil, nil)
		assert.True(t, got)
	})

	t.Run("valid expression evaluates correctly", func(t *testing.T) {
		values := map[string]*dbstore.AttributeValue{
			":v": {S: ptrStr("test")},
		}
		got := evaluateFilterExpression(item, "#n = :v", map[string]string{"#n": "name"}, values)
		assert.True(t, got)
	})
}

func TestEvaluateConditionExpr_AND_OR(t *testing.T) {
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
		got, err := evaluateConditionExpr(item, "#a = :age AND #s = :status", names, values)
		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("AND one false", func(t *testing.T) {
		got, err := evaluateConditionExpr(item, "#a = :age AND #s = :old", names, values)
		assert.NoError(t, err)
		assert.False(t, got)
	})

	t.Run("OR one true", func(t *testing.T) {
		got, err := evaluateConditionExpr(item, "#a = :old OR #s = :status", names, values)
		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("OR both false", func(t *testing.T) {
		got, err := evaluateConditionExpr(item, "#a = :old OR #s = :age", names, values)
		assert.NoError(t, err)
		assert.False(t, got)
	})
}
