package cypherparser

import (
	"math"
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage/graphengine"
)

func ep(e Expression) *Expression { return &e }

// ---------------------------------------------------------------------------
// evalExpr tests
// ---------------------------------------------------------------------------

func TestEvalExpr_Literal(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	tests := []struct {
		expr  Expression
		want  any
		isErr bool
	}{
		{litExpr("hello"), "hello", false},
		{litExpr(42), 42, false},
		{litExpr(3.14), 3.14, false},
		{litExpr(true), true, false},
		{litExpr(nil), nil, false},
	}

	for _, tt := range tests {
		got, err := evalExpr(ctx, &tt.expr)
		if tt.isErr {
			if err == nil {
				t.Errorf("expected error for %v", tt.expr)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		if got != tt.want {
			t.Errorf("got %v (%T), want %v (%T)", got, got, tt.want, tt.want)
		}
	}
}

func TestEvalExpr_VarRef(t *testing.T) {
	node := &graphengine.Node{ID: 1, Labels: []string{"Person"}, Props: graphengine.Props{"name": "Alice"}}
	ctx := &EvalContext{Bindings: map[string]any{"n": node}}

	got, err := evalExpr(ctx, ep(varRefExpr("n")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != node {
		t.Errorf("got %v, want %v", got, node)
	}

	_, err = evalExpr(ctx, ep(varRefExpr("x")))
	if err == nil {
		t.Error("expected error for unbound variable")
	}
}

func TestEvalExpr_PropAccess(t *testing.T) {
	node := &graphengine.Node{ID: 1, Props: graphengine.Props{"name": "Alice", "age": 30}}
	ctx := &EvalContext{Bindings: map[string]any{"n": node}}

	got, err := evalExpr(ctx, ep(propExpr("n", "name")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Alice" {
		t.Errorf("got %v, want Alice", got)
	}

	got, err = evalExpr(ctx, ep(propExpr("n", "age")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 30 {
		t.Errorf("got %v, want 30", got)
	}

	got, err = evalExpr(ctx, ep(propExpr("n", "missing")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}

	nilCtx := &EvalContext{Bindings: map[string]any{"n": nil}}
	got, err = evalExpr(nilCtx, ep(propExpr("n", "name")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("got %v, want nil for nil binding", got)
	}
}

func TestEvalExpr_Comparison(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	tests := []struct {
		name string
		expr Expression
		want bool
	}{
		{"eq true", compExpr(litExpr(1), OpEq, litExpr(1)), true},
		{"eq false", compExpr(litExpr(1), OpEq, litExpr(2)), false},
		{"neq true", compExpr(litExpr(1), OpNeq, litExpr(2)), true},
		{"neq false", compExpr(litExpr(1), OpNeq, litExpr(1)), false},
		{"lt true", compExpr(litExpr(1), OpLt, litExpr(2)), true},
		{"lt false", compExpr(litExpr(2), OpLt, litExpr(1)), false},
		{"gt true", compExpr(litExpr(2), OpGt, litExpr(1)), true},
		{"lte true", compExpr(litExpr(1), OpLte, litExpr(1)), true},
		{"gte true", compExpr(litExpr(1), OpGte, litExpr(1)), true},
		{"string eq", compExpr(litExpr("abc"), OpEq, litExpr("abc")), true},
		{"string lt", compExpr(litExpr("a"), OpLt, litExpr("b")), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evalExpr(ctx, &tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvalExpr_And(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(andExpr(litExpr(true), litExpr(true))))
	if got != true {
		t.Errorf("true AND true = %v", got)
	}

	got, _ = evalExpr(ctx, ep(andExpr(litExpr(true), litExpr(false))))
	if got != false {
		t.Errorf("true AND false = %v", got)
	}

	got, _ = evalExpr(ctx, ep(andExpr(litExpr(false), litExpr(true))))
	if got != false {
		t.Errorf("false AND true = %v", got)
	}

	got, _ = evalExpr(ctx, ep(andExpr(litExpr(false), litExpr(false))))
	if got != false {
		t.Errorf("false AND false = %v", got)
	}
}

func TestEvalExpr_Or(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(orExpr(litExpr(true), litExpr(false))))
	if got != true {
		t.Errorf("true OR false = %v", got)
	}

	got, _ = evalExpr(ctx, ep(orExpr(litExpr(false), litExpr(false))))
	if got != false {
		t.Errorf("false OR false = %v", got)
	}
}

func TestEvalExpr_Not(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(notExpr(litExpr(true))))
	if got != false {
		t.Errorf("NOT true = %v", got)
	}

	got, _ = evalExpr(ctx, ep(notExpr(litExpr(false))))
	if got != true {
		t.Errorf("NOT false = %v", got)
	}
}

func TestEvalExpr_Param(t *testing.T) {
	ctx := &EvalContext{
		Bindings: map[string]any{},
		Params:   map[string]any{"name": "Alice", "age": 30},
	}

	expr := Expression{Kind: ExprParam, ParamName: "name"}
	got, err := evalExpr(ctx, &expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Alice" {
		t.Errorf("got %v, want Alice", got)
	}

	expr2 := Expression{Kind: ExprParam, ParamName: "missing"}
	_, err = evalExpr(ctx, &expr2)
	if err == nil {
		t.Error("expected error for missing param")
	}

	noParamsCtx := &EvalContext{Bindings: map[string]any{}}
	_, err = evalExpr(noParamsCtx, &expr)
	if err == nil {
		t.Error("expected error when no params provided")
	}
}

// ---------------------------------------------------------------------------
// String predicates
// ---------------------------------------------------------------------------

func TestEvalExpr_In(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	expr := inExpr(litExpr("a"), listExpr(litExpr("a"), litExpr("b"), litExpr("c")))
	got, err := evalExpr(ctx, &expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Errorf("got %v, want true", got)
	}

	expr2 := inExpr(litExpr("x"), listExpr(litExpr("a"), litExpr("b")))
	got, err = evalExpr(ctx, &expr2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != false {
		t.Errorf("got %v, want false", got)
	}

	expr3 := inExpr(litExpr(1), litExpr("not a list"))
	got, err = evalExpr(ctx, &expr3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestEvalExpr_StartsWith(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	expr := startsWithExpr(litExpr("hello world"), litExpr("hello"))
	got, _ := evalExpr(ctx, &expr)
	if got != true {
		t.Errorf("got %v, want true", got)
	}

	expr2 := startsWithExpr(litExpr("hello"), litExpr("world"))
	got, _ = evalExpr(ctx, &expr2)
	if got != false {
		t.Errorf("got %v, want false", got)
	}

	expr3 := startsWithExpr(litExpr(42), litExpr("4"))
	got, _ = evalExpr(ctx, &expr3)
	if got != false {
		t.Errorf("non-string starts with should return false, got %v", got)
	}
}

func TestEvalExpr_EndsWith(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	expr := endsWithExpr(litExpr("hello world"), litExpr("world"))
	got, _ := evalExpr(ctx, &expr)
	if got != true {
		t.Errorf("got %v, want true", got)
	}

	expr2 := endsWithExpr(litExpr("hello"), litExpr("world"))
	got, _ = evalExpr(ctx, &expr2)
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}

func TestEvalExpr_Contains(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	expr := containsExpr(litExpr("hello world"), litExpr("lo wo"))
	got, _ := evalExpr(ctx, &expr)
	if got != true {
		t.Errorf("got %v, want true", got)
	}

	expr2 := containsExpr(litExpr("hello"), litExpr("xyz"))
	got, _ = evalExpr(ctx, &expr2)
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}

func TestEvalExpr_IsNull(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(isNullExpr(litExpr(nil))))
	if got != true {
		t.Errorf("nil IS NULL = %v", got)
	}

	got, _ = evalExpr(ctx, ep(isNullExpr(litExpr("hello"))))
	if got != false {
		t.Errorf("'hello' IS NULL = %v", got)
	}

	got, _ = evalExpr(ctx, ep(isNotNullExpr(litExpr(nil))))
	if got != false {
		t.Errorf("nil IS NOT NULL = %v", got)
	}

	got, _ = evalExpr(ctx, ep(isNotNullExpr(litExpr("hello"))))
	if got != true {
		t.Errorf("'hello' IS NOT NULL = %v", got)
	}
}

func TestEvalExpr_RegexMatch(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	expr := regexMatchExpr(litExpr("hello123"), litExpr(`^h\w+3$`))
	got, err := evalExpr(ctx, &expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Errorf("got %v, want true", got)
	}

	expr2 := regexMatchExpr(litExpr("hello"), litExpr(`^\d+$`))
	got, _ = evalExpr(ctx, &expr2)
	if got != false {
		t.Errorf("got %v, want false", got)
	}

	expr3 := regexMatchExpr(litExpr("test"), litExpr(`[invalid`))
	_, err = evalExpr(ctx, &expr3)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

// ---------------------------------------------------------------------------
// CASE expression
// ---------------------------------------------------------------------------

func TestEvalExpr_CaseSimple(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	expr := caseExpr(nil, []CaseWhen{
		{Condition: ep(compExpr(litExpr(1), OpLt, litExpr(2))), Result: ep(litExpr("small"))},
		{Condition: ep(compExpr(litExpr(1), OpLt, litExpr(5))), Result: ep(litExpr("medium"))},
	}, ep(litExpr("big")))

	got, _ := evalExpr(ctx, &expr)
	if got != "small" {
		t.Errorf("got %v, want small", got)
	}
}

func TestEvalExpr_CaseGeneral(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	expr := caseExpr(
		ep(litExpr("b")),
		[]CaseWhen{
			{Condition: ep(litExpr("a")), Result: ep(litExpr(1))},
			{Condition: ep(litExpr("b")), Result: ep(litExpr(2))},
			{Condition: ep(litExpr("c")), Result: ep(litExpr(3))},
		},
		ep(litExpr(0)),
	)

	got, _ := evalExpr(ctx, &expr)
	if got != 2 {
		t.Errorf("got %v, want 2", got)
	}

	expr2 := caseExpr(
		ep(litExpr("x")),
		[]CaseWhen{
			{Condition: ep(litExpr("a")), Result: ep(litExpr(1))},
		},
		ep(litExpr(42)),
	)

	got, _ = evalExpr(ctx, &expr2)
	if got != 42 {
		t.Errorf("got %v, want 42 (ELSE)", got)
	}

	expr3 := caseExpr(
		ep(litExpr("x")),
		[]CaseWhen{
			{Condition: ep(litExpr("a")), Result: ep(litExpr(1))},
		},
		nil,
	)

	got, _ = evalExpr(ctx, &expr3)
	if got != nil {
		t.Errorf("got %v, want nil (no ELSE, no match)", got)
	}
}

// ---------------------------------------------------------------------------
// Arithmetic
// ---------------------------------------------------------------------------

func TestEvalExpr_Arithmetic(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	tests := []struct {
		name string
		expr Expression
		want any
	}{
		{"int add", addExpr(litExpr(3), litExpr(4)), int64(7)},
		{"int sub", subExpr(litExpr(10), litExpr(3)), int64(7)},
		{"int mul", mulExpr(litExpr(3), litExpr(4)), int64(12)},
		{"int div", divExpr(litExpr(10), litExpr(3)), int64(3)},
		{"int mod", modExpr(litExpr(10), litExpr(3)), int64(1)},
		{"float add", addExpr(litExpr(3.5), litExpr(2.5)), 3.5 + 2.5},
		{"float sub", subExpr(litExpr(10.0), litExpr(3.5)), 10.0 - 3.5},
		{"float mul", mulExpr(litExpr(3.0), litExpr(4.0)), 12.0},
		{"float div", divExpr(litExpr(10.0), litExpr(3.0)), 10.0 / 3.0},
		{"float mod", modExpr(litExpr(10.5), litExpr(3.0)), math.Mod(10.5, 3.0)},
		{"string concat", addExpr(litExpr("hello"), litExpr(" world")), "hello world"},
		{"precedence", addExpr(mulExpr(litExpr(2), litExpr(3)), litExpr(4)), int64(10)},
		{"unary minus", subExpr(litExpr(0), litExpr(5)), int64(-5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evalExpr(ctx, &tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if compareValues(got, tt.want) != 0 || got != tt.want {
				switch gv := got.(type) {
				case float64:
					if math.Abs(gv-tt.want.(float64)) > 1e-9 {
						t.Errorf("got %v, want %v", got, tt.want)
					}
				default:
					if got != tt.want {
						t.Errorf("got %v (%T), want %v (%T)", got, got, tt.want, tt.want)
					}
				}
			}
		})
	}
}

func TestEvalExpr_DivByZero(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	v, err := evalExpr(ctx, ep(divExpr(litExpr(1), litExpr(0))))
	if err != nil {
		t.Errorf("division by zero: expected nil, got error: %v", err)
	}
	if v != nil {
		t.Errorf("division by zero: expected nil, got %v", v)
	}

	v, err = evalExpr(ctx, ep(modExpr(litExpr(1), litExpr(0))))
	if err != nil {
		t.Errorf("modulo by zero: expected nil, got error: %v", err)
	}
	if v != nil {
		t.Errorf("modulo by zero: expected nil, got %v", v)
	}
}

// ---------------------------------------------------------------------------
// List / Map literals
// ---------------------------------------------------------------------------

func TestEvalExpr_ListLiteral(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	expr := listExpr(litExpr(1), litExpr(2), litExpr(3))
	got, err := evalExpr(ctx, &expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, ok := got.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", got)
	}
	if len(list) != 3 || list[0] != 1 || list[1] != 2 || list[2] != 3 {
		t.Errorf("got %v, want [1, 2, 3]", list)
	}
}

func TestEvalExpr_MapLiteral(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	expr := mapExpr(MapPair{Key: "a", Value: litExpr(1)}, MapPair{Key: "b", Value: litExpr("hello")})
	got, err := evalExpr(ctx, &expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", got)
	}
	if m["a"] != 1 || m["b"] != "hello" {
		t.Errorf("got %v", m)
	}
}

// ---------------------------------------------------------------------------
// Built-in functions
// ---------------------------------------------------------------------------

func TestEvalExpr_FuncType(t *testing.T) {
	edge := &graphengine.Edge{ID: 1, Label: "KNOWS", From: 1, To: 2}
	ctx := &EvalContext{Bindings: map[string]any{"r": edge}}

	got, err := evalExpr(ctx, ep(funcCallExpr("type", varRefExpr("r"))))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "KNOWS" {
		t.Errorf("got %v, want KNOWS", got)
	}

	node := &graphengine.Node{ID: 1}
	ctx2 := &EvalContext{Bindings: map[string]any{"n": node}}
	got, _ = evalExpr(ctx2, ep(funcCallExpr("type", varRefExpr("n"))))
	if got != nil {
		t.Errorf("type() on node should return nil, got %v", got)
	}
}

func TestEvalExpr_FuncID(t *testing.T) {
	node := &graphengine.Node{ID: 42}
	ctx := &EvalContext{Bindings: map[string]any{"n": node}}

	got, err := evalExpr(ctx, ep(funcCallExpr("id", varRefExpr("n"))))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != int64(42) {
		t.Errorf("got %v, want 42", got)
	}

	edge := &graphengine.Edge{ID: 7}
	ctx2 := &EvalContext{Bindings: map[string]any{"r": edge}}
	got, _ = evalExpr(ctx2, ep(funcCallExpr("id", varRefExpr("r"))))
	if got != int64(7) {
		t.Errorf("got %v, want 7", got)
	}
}

func TestEvalExpr_FuncCoalesce(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(funcCallExpr("coalesce", litExpr(nil), litExpr("fallback"))))
	if got != "fallback" {
		t.Errorf("got %v, want fallback", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("coalesce", litExpr("first"), litExpr("second"))))
	if got != "first" {
		t.Errorf("got %v, want first", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("coalesce", litExpr(nil), litExpr(nil))))
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestEvalExpr_FuncSize(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(funcCallExpr("size", litExpr("hello"))))
	if got != int64(5) {
		t.Errorf("got %v, want 5", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("size", listExpr(litExpr(1), litExpr(2), litExpr(3)))))
	if got != int64(3) {
		t.Errorf("got %v, want 3", got)
	}
}

func TestEvalExpr_FuncToString(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(funcCallExpr("tostring", litExpr(42))))
	if got != "42" {
		t.Errorf("got %v, want '42'", got)
	}
}

func TestEvalExpr_FuncToUpper(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(funcCallExpr("upper", litExpr("hello"))))
	if got != "HELLO" {
		t.Errorf("got %v, want HELLO", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("lower", litExpr("HELLO"))))
	if got != "hello" {
		t.Errorf("got %v, want hello", got)
	}
}

func TestEvalExpr_FuncReplace(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(funcCallExpr("replace", litExpr("hello world"), litExpr("world"), litExpr("there"))))
	if got != "hello there" {
		t.Errorf("got %v, want 'hello there'", got)
	}
}

func TestEvalExpr_FuncSplit(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(funcCallExpr("split", litExpr("a,b,c"), litExpr(","))))
	list, ok := got.([]any)
	if !ok || len(list) != 3 || list[0] != "a" || list[1] != "b" || list[2] != "c" {
		t.Errorf("got %v, want [a, b, c]", got)
	}
}

func TestEvalExpr_FuncRange(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(funcCallExpr("range", litExpr(1), litExpr(3))))
	list, ok := got.([]any)
	if !ok || len(list) != 3 {
		t.Fatalf("got %v, want [1, 2, 3]", got)
	}
	if list[0] != int64(1) || list[1] != int64(2) || list[2] != int64(3) {
		t.Errorf("got %v", list)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("range", litExpr(0), litExpr(10), litExpr(3))))
	list, ok = got.([]any)
	if !ok || len(list) != 4 {
		t.Fatalf("range(0, 10, 3): got %v, want [0, 3, 6, 9]", got)
	}
}

func TestEvalExpr_FuncAbs(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(funcCallExpr("abs", litExpr(-5))))
	if got != 5.0 {
		t.Errorf("got %v, want 5.0", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("ceil", litExpr(3.2))))
	if got != 4.0 {
		t.Errorf("got %v, want 4.0", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("floor", litExpr(3.8))))
	if got != 3.0 {
		t.Errorf("got %v, want 3.0", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("round", litExpr(3.5))))
	if got != 4.0 {
		t.Errorf("got %v, want 4.0", got)
	}
}

func TestEvalExpr_UnknownFunc(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	_, err := evalExpr(ctx, ep(funcCallExpr("nonexistent", litExpr(1))))
	if err == nil {
		t.Error("expected error for unknown function")
	}
}

// ---------------------------------------------------------------------------
// LABELS / PROPERTIES
// ---------------------------------------------------------------------------

func TestEvalExpr_Labels(t *testing.T) {
	node := &graphengine.Node{ID: 1, Labels: []string{"Person", "Actor"}}
	ctx := &EvalContext{Bindings: map[string]any{"n": node}}

	got, _ := evalExpr(ctx, ep(labelsExpr(varRefExpr("n"))))
	labels, ok := got.([]string)
	if !ok || len(labels) != 2 || labels[0] != "Person" || labels[1] != "Actor" {
		t.Errorf("got %v, want [Person, Actor]", got)
	}

	edge := &graphengine.Edge{ID: 1, Label: "KNOWS"}
	ctx2 := &EvalContext{Bindings: map[string]any{"r": edge}}
	got, _ = evalExpr(ctx2, ep(labelsExpr(varRefExpr("r"))))
	if got != nil {
		t.Errorf("labels() on edge should return nil, got %v", got)
	}
}

func TestEvalExpr_Properties(t *testing.T) {
	node := &graphengine.Node{ID: 1, Props: graphengine.Props{"name": "Alice", "age": 30}}
	ctx := &EvalContext{Bindings: map[string]any{"n": node}}

	got, _ := evalExpr(ctx, ep(propertiesExpr(varRefExpr("n"))))
	m, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", got)
	}
	if m["name"] != "Alice" || m["age"] != 30 {
		t.Errorf("got %v", m)
	}

	nilNode := &graphengine.Node{ID: 2}
	ctx2 := &EvalContext{Bindings: map[string]any{"n": nilNode}}
	got, _ = evalExpr(ctx2, ep(propertiesExpr(varRefExpr("n"))))
	m, ok = got.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", got)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

// ---------------------------------------------------------------------------
// Value helpers
// ---------------------------------------------------------------------------

func TestCompareValues(t *testing.T) {
	tests := []struct {
		a, b any
		want int
	}{
		{nil, nil, 0},
		{nil, 1, -1},
		{1, nil, 1},
		{1, 1, 0},
		{1, 2, -1},
		{2, 1, 1},
		{1.5, 1.5, 0},
		{1.5, 2.5, -1},
		{"a", "a", 0},
		{"a", "b", -1},
		{"b", "a", 1},
		{true, false, 1},
		{false, true, -1},
		{true, true, 0},
	}

	for _, tt := range tests {
		got := compareValues(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("compareValues(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		v    any
		want float64
		ok   bool
	}{
		{float64(3.14), 3.14, true},
		{int(42), 42.0, true},
		{int64(42), 42.0, true},
		{uint(42), 42.0, true},
		{"hello", 0, false},
		{nil, 0, false},
	}

	for _, tt := range tests {
		got, ok := toFloat64(tt.v)
		if ok != tt.ok {
			t.Errorf("toFloat64(%v) ok = %v, want %v", tt.v, ok, tt.ok)
			continue
		}
		if ok && got != tt.want {
			t.Errorf("toFloat64(%v) = %v, want %v", tt.v, got, tt.want)
		}
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		v    any
		want bool
	}{
		{nil, false},
		{true, true},
		{false, false},
		{int64(1), true},
		{int64(0), false},
		{float64(1.0), true},
		{float64(0.0), false},
		{"hello", true},
		{"", false},
	}

	for _, tt := range tests {
		got := toBool(tt.v)
		if got != tt.want {
			t.Errorf("toBool(%v) = %v, want %v", tt.v, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Property / label matching
// ---------------------------------------------------------------------------

func TestGetProperty(t *testing.T) {
	node := &graphengine.Node{Props: graphengine.Props{"name": "Alice", "age": 30}}
	if getProperty(node, "name") != "Alice" {
		t.Error("expected Alice")
	}
	if getProperty(node, "missing") != nil {
		t.Error("expected nil")
	}

	edge := &graphengine.Edge{Label: "KNOWS", Props: graphengine.Props{"since": 2020}}
	if getProperty(edge, "label") != "KNOWS" {
		t.Error("expected KNOWS")
	}
	if getProperty(edge, "type") != "KNOWS" {
		t.Error("expected KNOWS (type alias)")
	}
	if getProperty(edge, "since") != 2020 {
		t.Error("expected 2020")
	}

	m := map[string]any{"key": "value"}
	if getProperty(m, "key") != "value" {
		t.Error("expected value")
	}
}

func TestMatchProps(t *testing.T) {
	tests := []struct {
		name        string
		actual      graphengine.Props
		constraints map[string]any
		want        bool
	}{
		{"nil constraints", graphengine.Props{"a": 1}, nil, true},
		{"match all", graphengine.Props{"a": 1, "b": "x"}, map[string]any{"a": 1, "b": "x"}, true},
		{"missing key", graphengine.Props{"a": 1}, map[string]any{"b": 1}, false},
		{"value mismatch", graphengine.Props{"a": 1}, map[string]any{"a": 2}, false},
		{"partial match", graphengine.Props{"a": 1, "b": 2}, map[string]any{"a": 1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchProps(tt.actual, tt.constraints); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchLabels(t *testing.T) {
	tests := []struct {
		nodeLabels []string
		required   []string
		want       bool
	}{
		{[]string{"Person"}, []string{"Person"}, true},
		{[]string{"Person", "Actor"}, []string{"Person"}, true},
		{[]string{"Person", "Actor"}, []string{"Person", "Actor"}, true},
		{[]string{"Person"}, []string{"Actor"}, false},
		{[]string{"Person"}, []string{}, true},
		{[]string{}, []string{}, true},
		{[]string{}, []string{"Person"}, false},
	}

	for _, tt := range tests {
		if got := matchLabels(tt.nodeLabels, tt.required); got != tt.want {
			t.Errorf("matchLabels(%v, %v) = %v, want %v", tt.nodeLabels, tt.required, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Aggregation
// ---------------------------------------------------------------------------

func TestHasAggregation(t *testing.T) {
	tests := []struct {
		name string
		expr Expression
		want bool
	}{
		{"literal", litExpr(1), false},
		{"var ref", varRefExpr("n"), false},
		{"agg count", aggExpr(AggCount, nil), true},
		{"agg sum", aggExpr(AggSum, ep(varRefExpr("n.age"))), true},
		{"comparison no agg", compExpr(varRefExpr("n.age"), OpGt, litExpr(10)), false},
		{"and with agg", andExpr(litExpr(true), aggExpr(AggCount, nil)), true},
		{"arithmetic with agg", addExpr(litExpr(1), aggExpr(AggSum, ep(varRefExpr("n.age")))), true},
		{"nested in func", funcCallExpr("size", aggExpr(AggCollect, ep(varRefExpr("n.name")))), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAggregation(&tt.expr); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeAggregation(t *testing.T) {
	rows := []map[string]any{
		{"n": &graphengine.Node{Props: graphengine.Props{"age": float64(20)}}},
		{"n": &graphengine.Node{Props: graphengine.Props{"age": float64(30)}}},
		{"n": &graphengine.Node{Props: graphengine.Props{"age": float64(40)}}},
	}

	t.Run("COUNT(*)", func(t *testing.T) {
		got, err := ComputeAggregation(AggCount, nil, false, rows)
		if err != nil || got != int64(3) {
			t.Errorf("got %v, err %v, want 3", got, err)
		}
	})

	t.Run("COUNT(expr)", func(t *testing.T) {
		got, err := ComputeAggregation(AggCount, ep(propExpr("n", "age")), false, rows)
		if err != nil || got != int64(3) {
			t.Errorf("got %v, err %v, want 3", got, err)
		}
	})

	t.Run("SUM", func(t *testing.T) {
		got, err := ComputeAggregation(AggSum, ep(propExpr("n", "age")), false, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 90.0 {
			t.Errorf("got %v, want 90.0", got)
		}
	})

	t.Run("AVG", func(t *testing.T) {
		got, err := ComputeAggregation(AggAvg, ep(propExpr("n", "age")), false, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 30.0 {
			t.Errorf("got %v, want 30.0", got)
		}
	})

	t.Run("MIN", func(t *testing.T) {
		got, err := ComputeAggregation(AggMin, ep(propExpr("n", "age")), false, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != float64(20) {
			t.Errorf("got %v, want 20.0", got)
		}
	})

	t.Run("MAX", func(t *testing.T) {
		got, err := ComputeAggregation(AggMax, ep(propExpr("n", "age")), false, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != float64(40) {
			t.Errorf("got %v, want 40.0", got)
		}
	})

	t.Run("COLLECT", func(t *testing.T) {
		got, err := ComputeAggregation(AggCollect, ep(propExpr("n", "age")), false, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		list, ok := got.([]any)
		if !ok || len(list) != 3 {
			t.Fatalf("got %v, want list of 3", got)
		}
	})

	t.Run("COUNT DISTINCT", func(t *testing.T) {
		distinctRows := []map[string]any{
			{"n": &graphengine.Node{Props: graphengine.Props{"city": "London"}}},
			{"n": &graphengine.Node{Props: graphengine.Props{"city": "Paris"}}},
			{"n": &graphengine.Node{Props: graphengine.Props{"city": "London"}}},
		}
		got, err := ComputeAggregation(AggCount, ep(propExpr("n", "city")), true, distinctRows)
		if err != nil || got != int64(2) {
			t.Errorf("got %v, err %v, want 2", got, err)
		}
	})
}

// ---------------------------------------------------------------------------
// Result projection helpers
// ---------------------------------------------------------------------------

func TestReturnItemName(t *testing.T) {
	tests := []struct {
		item ReturnItem
		want string
	}{
		{ReturnItem{Expr: varRefExpr("n")}, "n"},
		{ReturnItem{Expr: propExpr("n", "name")}, "n.name"},
		{ReturnItem{Expr: varRefExpr("n"), Alias: "person"}, "person"},
		{ReturnItem{Expr: aggExpr(AggCount, nil)}, "count(*)"},
		{ReturnItem{Expr: aggExpr(AggSum, ep(propExpr("n", "age")))}, "sum(n.age)"},
	}

	for _, tt := range tests {
		got := returnItemName(tt.item)
		if got != tt.want {
			t.Errorf("returnItemName(%+v) = %q, want %q", tt.item, got, tt.want)
		}
	}
}

func TestExprName(t *testing.T) {
	tests := []struct {
		expr Expression
		want string
	}{
		{varRefExpr("n"), "n"},
		{propExpr("n", "age"), "n.age"},
		{funcCallExpr("type", varRefExpr("r")), "type(r)"},
		{addExpr(varRefExpr("a"), varRefExpr("b")), "a + b"},
		{labelsExpr(varRefExpr("n")), "labels(n)"},
	}

	for _, tt := range tests {
		got := exprName(tt.expr)
		if got != tt.want {
			t.Errorf("exprName() = %q, want %q", got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Sort rows
// ---------------------------------------------------------------------------

func TestSortRows(t *testing.T) {
	rows := []map[string]any{
		{"n.age": float64(30)},
		{"n.age": float64(10)},
		{"n.age": float64(20)},
	}

	sortRows(rows, []OrderItem{
		{Expr: propExpr("n", "age"), Desc: false},
	})

	if rows[0]["n.age"] != float64(10) || rows[1]["n.age"] != float64(20) || rows[2]["n.age"] != float64(30) {
		t.Errorf("ASC sort failed: %v", rows)
	}

	sortRows(rows, []OrderItem{
		{Expr: propExpr("n", "age"), Desc: true},
	})

	if rows[0]["n.age"] != float64(30) || rows[1]["n.age"] != float64(20) || rows[2]["n.age"] != float64(10) {
		t.Errorf("DESC sort failed: %v", rows)
	}
}

// ---------------------------------------------------------------------------
// Top-K Heap
// ---------------------------------------------------------------------------

func TestTopKHeap(t *testing.T) {
	orderBy := []OrderItem{{Expr: varRefExpr("v"), Desc: false}}
	h := newTopKHeap(orderBy, 3)

	for i := 5; i >= 1; i-- {
		h.offer(topKItem{
			sortKey: []any{int64(i)},
			source:  i,
		})
	}

	sorted := h.sorted()
	if len(sorted) != 3 {
		t.Fatalf("expected 3 items, got %d", len(sorted))
	}
	if sorted[0].source != 1 || sorted[1].source != 2 || sorted[2].source != 3 {
		t.Errorf("expected [1, 2, 3], got %v", sorted)
	}
}

func TestTopKHeap_Desc(t *testing.T) {
	orderBy := []OrderItem{{Expr: varRefExpr("v"), Desc: true}}
	h := newTopKHeap(orderBy, 3)

	for i := 1; i <= 5; i++ {
		h.offer(topKItem{
			sortKey: []any{int64(i)},
			source:  i,
		})
	}

	sorted := h.sorted()
	if sorted[0].source != 5 || sorted[1].source != 4 || sorted[2].source != 3 {
		t.Errorf("expected [5, 4, 3], got %v", sorted)
	}
}

// ---------------------------------------------------------------------------
// Parameter resolution
// ---------------------------------------------------------------------------

func TestResolveParams(t *testing.T) {
	query := &CypherQuery{
		Match: MatchClause{
			Pattern: Pattern{
				Nodes: []NodePattern{
					{Variable: "n", Props: map[string]any{"name": paramRef("name")}},
				},
			},
		},
		Where: &Expression{
			Kind:  ExprComparison,
			Left:  ep(propExpr("n", "age")),
			Op:    OpGt,
			Right: &Expression{Kind: ExprParam, ParamName: "minAge"},
		},
		Return: ReturnClause{
			Items: []ReturnItem{
				{Expr: varRefExpr("n")},
			},
		},
	}

	params := map[string]any{"name": "Alice", "minAge": 25}

	err := ResolveParams(query, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if query.Match.Pattern.Nodes[0].Props["name"] != "Alice" {
		t.Errorf("param not resolved in node props: %v", query.Match.Pattern.Nodes[0].Props)
	}

	if query.Where.Kind != ExprComparison {
		t.Errorf("WHERE param not resolved, kind = %d", query.Where.Kind)
	}
	if query.Where.Right.Kind != ExprLiteral || query.Where.Right.LitValue != 25 {
		t.Errorf("WHERE param value not correct: %+v", query.Where.Right)
	}
}

func TestResolveParams_Missing(t *testing.T) {
	query := &CypherQuery{
		Where: &Expression{Kind: ExprParam, ParamName: "missing"},
	}

	err := ResolveParams(query, map[string]any{})
	if err == nil {
		t.Error("expected error for missing param")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Errorf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Aggregation with DISTINCT flag (parser integration)
// ---------------------------------------------------------------------------

func TestAggDistinctFlag(t *testing.T) {
	parsed, err := Parse("RETURN COUNT(DISTINCT n.city)")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	aggExpr := parsed.Read.Return.Items[0].Expr
	if aggExpr.Kind != ExprAggregation {
		t.Fatalf("expected ExprAggregation, got %d", aggExpr.Kind)
	}
	if !aggExpr.AggDistinct {
		t.Error("expected AggDistinct = true")
	}
	if aggExpr.AggFunc != AggCount {
		t.Errorf("expected AggCount, got %d", aggExpr.AggFunc)
	}
}

func TestAggExprName(t *testing.T) {
	tests := []struct {
		expr Expression
		want string
	}{
		{aggExpr(AggCount, nil), "count(*)"},
		{aggExpr(AggSum, ep(propExpr("n", "age"))), "sum(n.age)"},
		{
			Expression{Kind: ExprAggregation, AggFunc: AggCount, AggArg: ep(propExpr("n", "city")), AggDistinct: true},
			"count(DISTINCT n.city)",
		},
	}

	for _, tt := range tests {
		got := exprName(tt.expr)
		if got != tt.want {
			t.Errorf("got %q, want %q", got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Edge property access
// ---------------------------------------------------------------------------

func TestGetProperty_Edge(t *testing.T) {
	edge := &graphengine.Edge{
		ID: 1, Label: "KNOWS", From: 1, To: 2,
		Props: graphengine.Props{"since": 2020, "weight": 1.5},
	}

	if getProperty(edge, "label") != "KNOWS" {
		t.Error("expected label = KNOWS")
	}
	if getProperty(edge, "type") != "KNOWS" {
		t.Error("expected type = KNOWS")
	}
	if getProperty(edge, "since") != 2020 {
		t.Error("expected since = 2020")
	}
	if getProperty(edge, "weight") != 1.5 {
		t.Error("expected weight = 1.5")
	}
	if getProperty(edge, "missing") != nil {
		t.Error("expected nil for missing property")
	}
}

// ---------------------------------------------------------------------------
// Exists callback
// ---------------------------------------------------------------------------

func TestEvalExpr_Exists(t *testing.T) {
	ctx := &EvalContext{
		Bindings: map[string]any{"n": &graphengine.Node{ID: 1}},
		ExistsFn: func(p *Pattern, bindings map[string]any) (bool, error) {
			return true, nil
		},
	}

	pattern := &Pattern{
		Nodes: []NodePattern{{Variable: "m"}},
	}
	expr := existsExpr(pattern)
	got, err := evalExpr(ctx, &expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Errorf("got %v, want true", got)
	}

	noCallbackCtx := &EvalContext{Bindings: map[string]any{}}
	got, err = evalExpr(noCallbackCtx, &expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}

// ---------------------------------------------------------------------------
// Integration: parse + evaluate
// ---------------------------------------------------------------------------

func TestParseAndEval_SimpleWhere(t *testing.T) {
	parsed, err := Parse("MATCH (n) WHERE n.age > 25 AND n.name = 'Alice' RETURN n")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	node := &graphengine.Node{Props: graphengine.Props{"age": float64(30), "name": "Alice"}}
	ctx := &EvalContext{Bindings: map[string]any{"n": node}}

	ok, err := evalBool(ctx, parsed.Read.Where)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	if !ok {
		t.Error("expected WHERE to match")
	}

	node2 := &graphengine.Node{Props: graphengine.Props{"age": float64(20), "name": "Alice"}}
	ctx2 := &EvalContext{Bindings: map[string]any{"n": node2}}
	ok, _ = evalBool(ctx2, parsed.Read.Where)
	if ok {
		t.Error("expected WHERE to not match (age 20 < 25)")
	}
}

func TestParseAndEval_WithParams(t *testing.T) {
	parsed, err := Parse("MATCH (n {name: $name}) WHERE n.age > $minAge RETURN n")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	params := map[string]any{"name": "Alice", "minAge": 25}
	if err := ResolveParams(parsed.Read, params); err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	if parsed.Read.Match.Pattern.Nodes[0].Props["name"] != "Alice" {
		t.Error("param not resolved in node props")
	}
}

func TestParseAndEval_ComplexExpr(t *testing.T) {
	parsed, err := Parse("MATCH (n) WHERE n.name STARTS WITH 'A' AND n.age >= 18 RETURN n.name, n.age * 2 AS double_age")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	node := &graphengine.Node{Props: graphengine.Props{"name": "Alice", "age": float64(25)}}
	ctx := &EvalContext{Bindings: map[string]any{"n": node}}

	ok, err := evalBool(ctx, parsed.Read.Where)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	if !ok {
		t.Error("expected WHERE to match")
	}

	got, err := evalExpr(ctx, &parsed.Read.Return.Items[1].Expr)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	if got != float64(50) {
		t.Errorf("got %v (%T), want 50", got, got)
	}
}

func TestParseAndEval_CaseInReturn(t *testing.T) {
	parsed, err := Parse("MATCH (n) RETURN CASE WHEN n.age < 18 THEN 'minor' WHEN n.age < 65 THEN 'adult' ELSE 'senior' END AS category")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	tests := []struct {
		age  float64
		want string
	}{
		{10, "minor"},
		{25, "adult"},
		{70, "senior"},
	}

	for _, tt := range tests {
		node := &graphengine.Node{Props: graphengine.Props{"age": tt.age}}
		ctx := &EvalContext{Bindings: map[string]any{"n": node}}

		got, err := evalExpr(ctx, &parsed.Read.Return.Items[0].Expr)
		if err != nil {
			t.Fatalf("eval error: %v", err)
		}
		if got != tt.want {
			t.Errorf("age %v: got %v, want %v", tt.age, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Edge property eval in expression
// ---------------------------------------------------------------------------

func TestParseAndEval_EdgeProps(t *testing.T) {
	parsed, err := Parse("MATCH (a)-[r:KNOWS {since: 2020}]->(b) WHERE r.weight > 0.5 RETURN r.weight")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	edge := &graphengine.Edge{
		ID: 1, Label: "KNOWS", From: 1, To: 2,
		Props: graphengine.Props{"since": 2020, "weight": 0.8},
	}
	ctx := &EvalContext{Bindings: map[string]any{"r": edge}}

	ok, err := evalBool(ctx, parsed.Read.Where)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	if !ok {
		t.Error("expected WHERE to match")
	}

	got, err := evalExpr(ctx, &parsed.Read.Return.Items[0].Expr)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	if got != 0.8 {
		t.Errorf("got %v, want 0.8", got)
	}
}

// ---------------------------------------------------------------------------
// List comparison in expressions
// ---------------------------------------------------------------------------

func TestCompareValues_Lists(t *testing.T) {
	tests := []struct {
		a, b any
		want int
	}{
		{[]any{1, 2}, []any{1, 2}, 0},
		{[]any{1, 2}, []any{1, 3}, -1},
		{[]any{1, 3}, []any{1, 2}, 1},
		{[]any{1}, []any{1, 2}, -1},
		{[]any{1, 2}, []any{1}, 1},
	}

	for _, tt := range tests {
		got := compareValues(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("compareValues(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// String concat with non-string
// ---------------------------------------------------------------------------

func TestEvalArith_StringConcat(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	got, _ := evalExpr(ctx, ep(addExpr(litExpr("count: "), litExpr(5))))
	if got != "count: 5" {
		t.Errorf("got %v, want 'count: 5'", got)
	}

	got, _ = evalExpr(ctx, ep(addExpr(litExpr(5), litExpr(" items"))))
	if got != "5 items" {
		t.Errorf("got %v, want '5 items'", got)
	}
}

// ---------------------------------------------------------------------------
// Exists function (scalar)
// ---------------------------------------------------------------------------

func TestEvalExpr_FuncExists(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{
		"n": &graphengine.Node{Props: graphengine.Props{"name": "Alice"}},
	}}

	got, _ := evalExpr(ctx, ep(funcCallExpr("exists", propExpr("n", "name"))))
	if got != true {
		t.Errorf("exists(n.name) should be true, got %v", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("exists", propExpr("n", "missing"))))
	if got != false {
		t.Errorf("exists(n.missing) should be false, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// head / last / reverse
// ---------------------------------------------------------------------------

func TestEvalExpr_ListFuncs(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	list := listExpr(litExpr(1), litExpr(2), litExpr(3))

	got, _ := evalExpr(ctx, ep(funcCallExpr("head", list)))
	if got != 1 {
		t.Errorf("head() = %v, want 1", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("last", list)))
	if got != 3 {
		t.Errorf("last() = %v, want 3", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("reverse", list)))
	rev, ok := got.([]any)
	if !ok || len(rev) != 3 || rev[0] != 3 || rev[1] != 2 || rev[2] != 1 {
		t.Errorf("reverse() = %v, want [3, 2, 1]", got)
	}

	got, _ = evalExpr(ctx, ep(funcCallExpr("head", listExpr())))
	if got != nil {
		t.Errorf("head() on empty list should be nil, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// evalRowExpr
// ---------------------------------------------------------------------------

func TestEvalRowExpr(t *testing.T) {
	row := map[string]any{
		"n.age": float64(25),
		"n":     &graphengine.Node{Props: graphengine.Props{"age": float64(25)}},
	}

	got := evalRowExpr(ep(propExpr("n", "age")), row)
	if got != float64(25) {
		t.Errorf("got %v, want 25", got)
	}

	got = evalRowExpr(ep(varRefExpr("n")), row)
	if got == nil {
		t.Error("got nil for variable lookup")
	}
}

// ---------------------------------------------------------------------------
// Benchmark
// ---------------------------------------------------------------------------

func BenchmarkEvalExpr_SimpleComparison(b *testing.B) {
	ctx := &EvalContext{Bindings: map[string]any{
		"n": &graphengine.Node{Props: graphengine.Props{"age": float64(25), "name": "Alice"}},
	}}
	expr := compExpr(propExpr("n", "age"), OpGt, litExpr(20))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evalExpr(ctx, &expr)
	}
}

func BenchmarkEvalExpr_ComplexWhere(b *testing.B) {
	ctx := &EvalContext{Bindings: map[string]any{
		"n": &graphengine.Node{Props: graphengine.Props{"age": float64(25), "name": "Alice", "city": "London"}},
	}}
	expr := andExpr(
		compExpr(propExpr("n", "age"), OpGt, litExpr(18)),
		compExpr(propExpr("n", "name"), OpEq, litExpr("Alice")),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evalExpr(ctx, &expr)
	}
}

func BenchmarkMatchProps(b *testing.B) {
	nodeProps := graphengine.Props{"name": "Alice", "age": float64(25), "city": "London"}
	constraints := map[string]any{"name": "Alice", "city": "London"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matchProps(nodeProps, constraints)
	}
}

func BenchmarkCompareValues(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compareValues(float64(25), float64(30))
	}
}

// Helper to test Parse + eval integration without needing graphengine.DB
func TestParseIntegrations(t *testing.T) {
	tests := []struct {
		name  string
		query string
		valid bool
	}{
		{"simple match", "MATCH (n) RETURN n", true},
		{"with where", "MATCH (n) WHERE n.age > 18 RETURN n.name", true},
		{"with params", "MATCH (n {name: $name}) RETURN n", true},
		{"with case", "MATCH (n) RETURN CASE WHEN n.age > 18 THEN 'adult' END", true},
		{"with aggregation", "MATCH (n) RETURN COUNT(*) AS cnt", true},
		{"with distinct agg", "MATCH (n) RETURN COUNT(DISTINCT n.city) AS cities", true},
		{"with in", "MATCH (n) WHERE n.city IN ['London', 'Paris'] RETURN n", true},
		{"with exists", "MATCH (n) WHERE EXISTS { (n)-[:KNOWS]->() } RETURN n", true},
		{"with regex", "MATCH (n) WHERE n.name =~ '^A.*' RETURN n", true},
		{"with is null", "MATCH (n) WHERE n.email IS NULL RETURN n", true},
		{"with labels", "MATCH (n) RETURN labels(n)", true},
		{"with properties", "MATCH (n) RETURN properties(n)", true},
		{"with arithmetic", "MATCH (n) RETURN n.age * 2 + 1 AS adjusted", true},
		{"with list literal", "MATCH (n) WHERE n.city IN ['London', 'Paris'] RETURN n", true},
		{"with map literal", "MATCH (n) RETURN {name: n.name, age: n.age} AS info", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.query)
			if tt.valid && err != nil {
				t.Errorf("expected valid parse, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("expected parse error, got nil")
			}
		})
	}
}

// Ensure we test error formatting
func TestEvalExpr_ErrorFormatting(t *testing.T) {
	ctx := &EvalContext{Bindings: map[string]any{}}

	_, err := evalExpr(ctx, &Expression{Kind: 999})
	if err == nil {
		t.Error("expected error for unsupported expression kind")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("error should contain 'unsupported', got: %v", err)
	}

	_, err = evalExpr(ctx, ep(divExpr(litExpr(1), litExpr(0))))
	if err != nil {
		t.Errorf("division by zero: expected nil result, got error: %v", err)
	}

	_, err = evalExpr(ctx, ep(modExpr(litExpr(1), litExpr(0))))
	if err != nil {
		t.Errorf("modulo by zero: expected nil result, got error: %v", err)
	}

	_, err = evalExpr(ctx, &Expression{
		Kind:       ExprArithMul,
		ArithLeft:  ep(litExpr(1)),
		ArithRight: ep(litExpr("not a number")),
	})
	if err == nil {
		t.Error("expected error for non-numeric arithmetic")
	}
}

// Test that the expression evaluator properly handles the full query AST
func TestResolveParams_AllClauses(t *testing.T) {
	q := &CypherQuery{
		Match: MatchClause{
			Pattern: Pattern{
				Nodes: []NodePattern{{Props: map[string]any{"k": paramRef("p1")}}},
			},
		},
		Where: &Expression{Kind: ExprParam, ParamName: "p2"},
		Return: ReturnClause{
			Items: []ReturnItem{{Expr: Expression{Kind: ExprParam, ParamName: "p3"}}},
		},
		Set: []SetItem{
			{Variable: "n", Property: "x", Value: Expression{Kind: ExprParam, ParamName: "p4"}},
		},
		OrderBy: []OrderItem{
			{Expr: Expression{Kind: ExprParam, ParamName: "p5"}},
		},
		With: &WithClause{
			Items: []WithItem{
				{Expr: Expression{Kind: ExprParam, ParamName: "p6"}},
			},
		},
		Unwind: &UnwindClause{
			Expr: Expression{Kind: ExprParam, ParamName: "p7"},
			Var:  "x",
		},
	}

	params := map[string]any{
		"p1": "v1", "p2": "v2", "p3": "v3",
		"p4": "v4", "p5": "v5", "p6": "v6", "p7": "v7",
	}

	err := ResolveParams(q, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q.Match.Pattern.Nodes[0].Props["k"] != "v1" {
		t.Error("node param not resolved")
	}
	if q.Where.Kind != ExprLiteral {
		t.Error("WHERE param not resolved")
	}
	if q.Return.Items[0].Expr.LitValue != "v3" {
		t.Error("RETURN param not resolved")
	}
	if q.Set[0].Value.LitValue != "v4" {
		t.Error("SET param not resolved")
	}
	if q.OrderBy[0].Expr.LitValue != "v5" {
		t.Error("ORDER BY param not resolved")
	}
	if q.With.Items[0].Expr.LitValue != "v6" {
		t.Error("WITH param not resolved")
	}
	if q.Unwind.Expr.LitValue != "v7" {
		t.Error("UNWIND param not resolved")
	}
}
