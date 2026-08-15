package cypherparser

import (
	"math"
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
