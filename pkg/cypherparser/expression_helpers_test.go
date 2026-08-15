package cypherparser

import (
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage/graphengine"
)

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
