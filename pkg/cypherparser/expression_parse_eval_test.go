package cypherparser

import (
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage/graphengine"
)

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
// evalOrderValue (was evalRowExpr)
// ---------------------------------------------------------------------------

func TestEvalOrderValue(t *testing.T) {
	row := map[string]any{
		"n.age": float64(25),
		"n":     &graphengine.Node{Props: graphengine.Props{"age": float64(25)}},
	}

	got := evalOrderValue(ep(propExpr("n", "age")), row)
	if got != float64(25) {
		t.Errorf("got %v, want 25", got)
	}

	got = evalOrderValue(ep(varRefExpr("n")), row)
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
