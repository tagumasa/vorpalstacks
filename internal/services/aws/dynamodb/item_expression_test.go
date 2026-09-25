package dynamodb

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestFindMatchingCloseParen(t *testing.T) {
	t.Run("simple match", func(t *testing.T) {
		assert.Equal(t, 4, findMatchingCloseParen("(abc)", 0))
	})

	// findMatchingCloseParen must skip over inner groups so that the
	// outermost close position is returned for nested function calls.
	t.Run("nested parentheses", func(t *testing.T) {
		assert.Equal(t, 10, findMatchingCloseParen("(a,(b,c),d)", 0))
	})

	t.Run("deeply nested", func(t *testing.T) {
		assert.Equal(t, 8, findMatchingCloseParen("(a(b(c)))", 0))
	})

	// Nested list_append calls must each be balanced correctly so that the
	// parser identifies the outer expression boundaries rather than the
	// first inner close.
	t.Run("list_append nested", func(t *testing.T) {
		s := "list_append(:a, list_append(:b, :c))"
		openPos := strings.Index(s, "(")
		closePos := findMatchingCloseParen(s, openPos)
		assert.Equal(t, len(s)-1, closePos,
			"should match outermost close paren")
		assert.Equal(t, s[openPos:closePos+1], "(:a, list_append(:b, :c))")
	})

	t.Run("no matching close returns -1", func(t *testing.T) {
		assert.Equal(t, -1, findMatchingCloseParen("(abc", 0))
	})

	t.Run("empty parens", func(t *testing.T) {
		assert.Equal(t, 1, findMatchingCloseParen("()", 0))
	})

	t.Run("multiple groups", func(t *testing.T) {
		assert.Equal(t, 2, findMatchingCloseParen("(a)(b)", 0))
		assert.Equal(t, 5, findMatchingCloseParen("(a)(b)", 3))
	})
}

// evalUpdateOperand compiles and evaluates one operand string of the
// UpdateExpression plane against the shared operand model, the entry the
// SET action itself goes through.
func evalUpdateOperand(t *testing.T, expr string, values map[string]*dbstore.AttributeValue, attrs map[string]*dbstore.AttributeValue) (*dbstore.AttributeValue, error) {
	t.Helper()
	operand, err := parseUpdateOperand(expr, values, nil)
	if err != nil {
		return nil, err
	}
	return operand.eval(attrs)
}

// list_append demands both operands be lists: a path operand naming an
// attribute the item does not hold is an evaluation error on either side,
// never a silent single-operand result wearing the other operand's value.
func TestListAppendMissingOperandIsAnError(t *testing.T) {
	str := func(v string) *dbstore.AttributeValue { return &dbstore.AttributeValue{S: &v} }
	values := map[string]*dbstore.AttributeValue{
		":new": dbstore.ListValue([]*dbstore.AttributeValue{str("b")}),
	}
	attrs := map[string]*dbstore.AttributeValue{"id": str("i1")}

	if _, err := evalUpdateOperand(t, "list_append(noSuchAttr, :new)", values, attrs); !errors.Is(err, ErrTypeMismatch) {
		t.Fatalf("missing first operand: err = %v, want ErrTypeMismatch", err)
	}
	if _, err := evalUpdateOperand(t, "list_append(id, :new)", values, attrs); !errors.Is(err, ErrTypeMismatch) {
		t.Fatalf("non-list first operand: err = %v, want ErrTypeMismatch", err)
	}
	attrs["items"] = dbstore.ListValue([]*dbstore.AttributeValue{str("a")})
	if _, err := evalUpdateOperand(t, "list_append(items, noSuchAttr)", values, attrs); !errors.Is(err, ErrTypeMismatch) {
		t.Fatalf("missing second operand: err = %v, want ErrTypeMismatch", err)
	}

	// The two-list case still appends in order.
	got, err := evalUpdateOperand(t, "list_append(items, :new)", values, attrs)
	if err != nil {
		t.Fatalf("append: %v", err)
	}
	if len(got.L) != 2 || *got.L[0].S != "a" || *got.L[1].S != "b" {
		t.Fatalf("appended list = %+v, want [a b]", got.L)
	}
}

// The nested-operand case of if_not_exists: the default operand is itself
// a function call, so the close-paren and comma splits must respect the
// nesting — list_append inside if_not_exists is the documented composition.
func TestIfNotExistsDefaultCarryingListAppend(t *testing.T) {
	str := func(v string) *dbstore.AttributeValue { return &dbstore.AttributeValue{S: &v} }
	values := map[string]*dbstore.AttributeValue{
		":l1": dbstore.ListValue([]*dbstore.AttributeValue{str("a")}),
		":l2": dbstore.ListValue([]*dbstore.AttributeValue{str("b")}),
	}

	// The path addresses nothing: the default's list_append concatenates.
	attrs := map[string]*dbstore.AttributeValue{}
	got, err := evalUpdateOperand(t, "if_not_exists(missing, list_append(:l1, :l2))", values, attrs)
	if err != nil {
		t.Fatalf("nested default: %v", err)
	}
	if l := got.L; len(l) != 2 || *l[0].S != "a" || *l[1].S != "b" {
		t.Fatalf("nested default: %+v, want [a b]", l)
	}

	// The path carries a value: the default never evaluates.
	attrs = map[string]*dbstore.AttributeValue{"present": str("kept")}
	got, err = evalUpdateOperand(t, "if_not_exists(present, list_append(:l1, :l2))", values, attrs)
	if err != nil {
		t.Fatalf("present path: %v", err)
	}
	if got == nil || got.S == nil || *got.S != "kept" {
		t.Fatalf("present path: %+v, want kept", got)
	}

	// Deeper nesting still splits on the matching parens.
	values[":l3"] = dbstore.ListValue([]*dbstore.AttributeValue{str("c")})
	got, err = evalUpdateOperand(t, "if_not_exists(missing, list_append(:l1, list_append(:l2, :l3)))", values, attrs)
	if err != nil {
		t.Fatalf("deep nesting: %v", err)
	}
	if l := got.L; len(l) != 3 || *l[0].S != "a" || *l[1].S != "b" || *l[2].S != "c" {
		t.Fatalf("deep nesting: %+v, want [a b c]", l)
	}

	// An unbalanced call is a validation error, not a mis-sliced operand.
	if _, err := evalUpdateOperand(t, "if_not_exists(missing, list_append(:l1, :l2)", values, attrs); err == nil {
		t.Fatal("unbalanced call: expected a validation error")
	}
}
