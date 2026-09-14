package sfn

import (
	"testing"
)

// TestResolveUpdateExpressionParsesDelete pins that the DELETE clause of a
// DynamoDB UpdateExpression is recorded instead of being silently consumed.
func TestResolveUpdateExpressionParsesDelete(t *testing.T) {
	parsed, err := resolveUpdateExpression(
		"SET a = :one ADD c :two DELETE tags :tags REMOVE d",
		map[string]interface{}{":one": float64(1), ":two": float64(2), ":tags": []interface{}{"x"}},
		nil)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if got := parsed.Sets["a"]; got != float64(1) {
		t.Errorf("SET a = %v, want 1", got)
	}
	if got := parsed.Adds["c"]; got != float64(2) {
		t.Errorf("ADD c = %v, want 2", got)
	}
	if got, ok := parsed.Deletes["tags"]; !ok || got == nil {
		t.Errorf("DELETE tags = %v, want the :tags set — the clause must take effect", got)
	}
	if len(parsed.Removes) != 1 || parsed.Removes[0] != "d" {
		t.Errorf("REMOVE = %v, want [d]", parsed.Removes)
	}
}

// TestResolveUpdateExpressionRejectsMalformed pins that a malformed
// UpdateExpression is a validation error, never a silent no-op update —
// DynamoDB rejects the update request itself, so an SFN-mediated update
// with the same expression must fail the task.
func TestResolveUpdateExpressionRejectsMalformed(t *testing.T) {
	values := map[string]interface{}{":v": float64(1), ":tags": []interface{}{"x"}}
	for _, expr := range []string{
		"SET",
		"SET a b c",
		"SET a = :missing",
		"ADD c",
		"DELETE tags",
		"REMOVE",
		"a = :v",
		"SET a = :v REMOVE",
	} {
		if _, err := resolveUpdateExpression(expr, values, nil); err == nil {
			t.Errorf("expression %q must be rejected, not silently committed as no operations", expr)
		}
	}
}

// TestApplyDynamoDBAdd pins the ADD semantics: numbers add, sets union, an
// absent attribute takes the value, and incompatible types error instead
// of silently replacing the attribute.
func TestApplyDynamoDBAdd(t *testing.T) {
	if got, err := applyDynamoDBAdd(float64(1), float64(2)); err != nil || got != float64(3) {
		t.Errorf("numeric ADD = %v, %v; want 3", got, err)
	}
	got, err := applyDynamoDBAdd([]interface{}{"a"}, []interface{}{"b"})
	if err != nil {
		t.Fatalf("set ADD errored: %v", err)
	}
	if len(got.([]interface{})) != 2 {
		t.Errorf("set ADD = %v, want the union", got)
	}
	if got, err := applyDynamoDBAdd(nil, float64(5)); err != nil || got != float64(5) {
		t.Errorf("ADD onto absent attribute = %v, %v; want 5", got, err)
	}
	if _, err := applyDynamoDBAdd("existing", float64(5)); err == nil {
		t.Error("ADD of a number onto a string attribute must error, not replace the value")
	}
}

// TestApplyDynamoDBAddSetUnionDeduplicates pins that a set ADD is a set
// union: a DynamoDB set holds no duplicates, so elements already present
// are not appended a second time and the existing order is preserved.
func TestApplyDynamoDBAddSetUnionDeduplicates(t *testing.T) {
	got, err := applyDynamoDBAdd([]interface{}{"b"}, []interface{}{"a", "b"})
	if err != nil {
		t.Fatalf("set ADD errored: %v", err)
	}
	want := []interface{}{"b", "a"}
	if len(got.([]interface{})) != len(want) {
		t.Fatalf("set ADD = %v, want %v", got, want)
	}
	for i, v := range want {
		if got.([]interface{})[i] != v {
			t.Errorf("set ADD = %v, want %v — duplicates must not accumulate", got, want)
		}
	}

	// Re-adding the exact same set changes nothing.
	again, err := applyDynamoDBAdd(want, []interface{}{"b", "a"})
	if err != nil {
		t.Fatalf("idempotent set ADD errored: %v", err)
	}
	if len(again.([]interface{})) != 2 {
		t.Errorf("idempotent set ADD = %v, want %v", again, want)
	}
}

// TestApplyDynamoDBDelete pins the DELETE semantics: set elements are
// removed, an emptied set drops the attribute, an absent attribute is a
// no-op, and a non-set attribute errors.
func TestApplyDynamoDBDelete(t *testing.T) {
	got, err := applyDynamoDBDelete([]interface{}{"a", "b", "c"}, []interface{}{"b"})
	if err != nil {
		t.Fatalf("set DELETE errored: %v", err)
	}
	remaining := got.([]interface{})
	if len(remaining) != 2 || remaining[0] != "a" || remaining[1] != "c" {
		t.Errorf("set DELETE = %v, want [a c]", remaining)
	}
	if got, err := applyDynamoDBDelete([]interface{}{"a"}, []interface{}{"a"}); err != nil || got != nil {
		t.Errorf("emptied set DELETE = %v, %v; want the attribute dropped", got, err)
	}
	if got, err := applyDynamoDBDelete(nil, []interface{}{"a"}); err != nil || got != nil {
		t.Errorf("DELETE on absent attribute = %v, %v; want a no-op", got, err)
	}
	if _, err := applyDynamoDBDelete("scalar", []interface{}{"a"}); err == nil {
		t.Error("DELETE on a non-set attribute must error")
	}
}

// TestAttributeNumberRoundTrip pins the attribute type round-trip: a
// numeric attribute (N) read back after a putItem stays numeric (N) rather
// than drifting to a string (S), and numeric sets keep NS.
func TestAttributeNumberRoundTrip(t *testing.T) {
	plain := awsAttrValueToPlain(map[string]interface{}{"N": "42"})
	if got := plain.(float64); got != 42 {
		t.Fatalf("plain value = %v (%T), want 42", plain, plain)
	}
	back := plainToAWSAttrValue(plain)
	n, ok := back.(map[string]interface{})["N"].(string)
	if !ok || n != "42" {
		t.Fatalf("round-tripped attribute = %v, want N=42", back)
	}

	setPlain := awsAttrValueToPlain(map[string]interface{}{"NS": []interface{}{"1", "2"}})
	backSet := plainToAWSAttrValue(setPlain)
	if _, ok := backSet.(map[string]interface{})["NS"]; !ok {
		t.Errorf("numeric set round-tripped as %v, want NS", backSet)
	}

	str := plainToAWSAttrValue(awsAttrValueToPlain(map[string]interface{}{"S": "text"}))
	if _, ok := str.(map[string]interface{})["S"]; !ok {
		t.Errorf("string attribute round-tripped as %v, want S", str)
	}
}
