package cypherparser

import (
	"context"
	"testing"

	"vorpalstacks/pkg/graphengine"
)

func setupSocialGraph(t *testing.T) *testGraph {
	t.Helper()
	db := openTestDB(t)
	g := &testGraph{db: db, nodes: make(map[string]graphengine.NodeID)}

	alice, _ := db.AddNode([]string{"Person"}, graphengine.Props{"name": "Alice", "age": 30})
	bob, _ := db.AddNode([]string{"Person"}, graphengine.Props{"name": "Bob", "age": 25})
	charlie, _ := db.AddNode([]string{"Person"}, graphengine.Props{"name": "Charlie", "age": 35})
	dave, _ := db.AddNode([]string{"Person"}, graphengine.Props{"name": "Dave", "age": 28})
	eve, _ := db.AddNode([]string{"Person"}, graphengine.Props{"name": "Eve", "age": 22})

	db.AddEdge(alice, bob, "KNOWS", nil)
	db.AddEdge(alice, charlie, "KNOWS", nil)
	db.AddEdge(bob, dave, "KNOWS", nil)
	db.AddEdge(charlie, eve, "KNOWS", nil)
	db.AddEdge(dave, eve, "FOLLOWS", nil)

	g.nodes["alice"] = alice
	g.nodes["bob"] = bob
	g.nodes["charlie"] = charlie
	g.nodes["dave"] = dave
	g.nodes["eve"] = eve
	return g
}

// ---------------------------------------------------------------------------
// UNWIND tests
// ---------------------------------------------------------------------------

func TestUnwind_BasicList(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	for _, name := range []string{"Alice", "Bob", "Charlie"} {
		db.AddNode([]string{"Person"}, graphengine.Props{"name": name})
	}

	result := execQuery(t, db, `UNWIND [1, 2, 3] AS x RETURN x`)
	if len(result.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(result.Rows))
	}

	expected := []int64{1, 2, 3}
	for i, row := range result.Rows {
		val, ok := row["x"].(int64)
		if !ok {
			t.Fatalf("row %d: expected int64, got %T", i, row["x"])
		}
		if val != expected[i] {
			t.Fatalf("row %d: expected %d, got %d", i, expected[i], val)
		}
	}
}

func TestUnwind_EmptyList(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	result := execQuery(t, db, `UNWIND [] AS x RETURN x`)
	if len(result.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(result.Rows))
	}
}

func TestUnwind_WithMatch(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	db.AddNode([]string{"Person"}, graphengine.Props{"name": "Alice"})
	db.AddNode([]string{"Person"}, graphengine.Props{"name": "Bob"})
	db.AddNode([]string{"Person"}, graphengine.Props{"name": "Charlie"})

	result := execQuery(t, db, `UNWIND ["Alice", "Bob"] AS name MATCH (n:Person {name: name}) RETURN n.name`)
	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result.Rows))
	}

	names := map[string]bool{}
	for _, row := range result.Rows {
		name := row["n.name"].(string)
		names[name] = true
	}
	if !names["Alice"] || !names["Bob"] {
		t.Fatalf("expected Alice and Bob, got %v", names)
	}
}

func TestUnwind_WithMatchNoResults(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	db.AddNode([]string{"Person"}, graphengine.Props{"name": "Alice"})

	result := execQuery(t, db, `UNWIND ["NonExistent"] AS name MATCH (n:Person {name: name}) RETURN n.name`)
	if len(result.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(result.Rows))
	}
}

func TestUnwind_WithMatchAndEdges(t *testing.T) {
	g := setupSocialGraph(t)
	defer g.db.Close()

	result := execQuery(t, g.db, `UNWIND ["Alice", "Bob"] AS name MATCH (n:Person {name: name})-[:KNOWS]->(friend) RETURN n.name, friend.name`)
	if len(result.Rows) < 1 {
		t.Fatalf("expected at least 1 row, got %d", len(result.Rows))
	}
}

func TestUnwind_WithParam(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	db.AddNode([]string{"Person"}, graphengine.Props{"name": "Alice"})
	db.AddNode([]string{"Person"}, graphengine.Props{"name": "Bob"})

	result := execQueryWithParams(t, db, `UNWIND $names AS name MATCH (n:Person {name: name}) RETURN n.name`, map[string]any{
		"names": []any{"Alice", "Bob"},
	})
	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result.Rows))
	}
}

func TestUnwind_WithWhere(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	db.AddNode([]string{"Person"}, graphengine.Props{"name": "Alice", "age": 30})
	db.AddNode([]string{"Person"}, graphengine.Props{"name": "Bob", "age": 25})
	db.AddNode([]string{"Person"}, graphengine.Props{"name": "Charlie", "age": 35})

	result := execQuery(t, db, `UNWIND [10, 20, 30, 40] AS x MATCH (n:Person) WHERE n.age >= x RETURN n.name, x ORDER BY x`)
	if len(result.Rows) == 0 {
		t.Fatal("expected rows, got 0")
	}
}

func TestUnwind_WithOrderBy(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	result := execQuery(t, db, `UNWIND [3, 1, 2] AS x RETURN x ORDER BY x`)
	if len(result.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(result.Rows))
	}
	if result.Rows[0]["x"].(int64) != 1 {
		t.Fatalf("first row should be 1, got %v", result.Rows[0]["x"])
	}
	if result.Rows[2]["x"].(int64) != 3 {
		t.Fatalf("last row should be 3, got %v", result.Rows[2]["x"])
	}
}

func TestUnwind_WithLimit(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	result := execQuery(t, db, `UNWIND [1, 2, 3, 4, 5] AS x RETURN x LIMIT 3`)
	if len(result.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(result.Rows))
	}
}

func TestUnwind_StringList(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	result := execQuery(t, db, `UNWIND ["hello", "world"] AS s RETURN s`)
	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result.Rows))
	}
	if result.Rows[0]["s"].(string) != "hello" {
		t.Fatalf("expected 'hello', got %v", result.Rows[0]["s"])
	}
}

// ---------------------------------------------------------------------------
// WITH tests
// ---------------------------------------------------------------------------

func TestWith_BasicProjection(t *testing.T) {
	g := setupPersonGraph(t)
	defer g.db.Close()

	result := execQuery(t, g.db, `MATCH (n:Person) WITH n.name AS name RETURN name ORDER BY name`)
	if len(result.Rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(result.Rows))
	}
	if result.Rows[0]["name"].(string) != "Alice" {
		t.Fatalf("first should be Alice, got %v", result.Rows[0]["name"])
	}
	if result.Rows[3]["name"].(string) != "Dave" {
		t.Fatalf("last should be Dave, got %v", result.Rows[3]["name"])
	}
}

func TestWith_FilterAfterWith(t *testing.T) {
	g := setupPersonGraph(t)
	defer g.db.Close()

	result := execQuery(t, g.db, `MATCH (n:Person) WITH n.name AS name WHERE name STARTS WITH "A" OR name STARTS WITH "B" RETURN name ORDER BY name`)
	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result.Rows))
	}
}

func TestWith_PropAccessProjection(t *testing.T) {
	g := setupPersonGraph(t)
	defer g.db.Close()

	result := execQuery(t, g.db, `MATCH (n:Person) WITH n.age AS age RETURN age ORDER BY age`)
	if len(result.Rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(result.Rows))
	}
	f, ok := toFloat64(result.Rows[0]["age"])
	if !ok || f != 25 {
		t.Fatalf("youngest should be 25, got %v", result.Rows[0]["age"])
	}
}

func TestWith_Aggregation(t *testing.T) {
	g := setupPersonGraph(t)
	defer g.db.Close()

	result := execQuery(t, g.db, `MATCH (n:Person) WITH count(n) AS cnt RETURN cnt`)
	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
	if result.Rows[0]["cnt"].(int64) != 4 {
		t.Fatalf("expected count 4, got %v", result.Rows[0]["cnt"])
	}
}

func TestWith_Distinct(t *testing.T) {
	g := setupPersonGraph(t)
	defer g.db.Close()

	result := execQuery(t, g.db, `MATCH (n:Person) WITH DISTINCT n.age AS age RETURN age ORDER BY age`)
	if len(result.Rows) != 4 {
		t.Fatalf("expected 4 distinct ages, got %d", len(result.Rows))
	}
}

func TestWith_MultipleItems(t *testing.T) {
	g := setupPersonGraph(t)
	defer g.db.Close()

	result := execQuery(t, g.db, `MATCH (n:Person) WITH n.name AS name, n.age AS age RETURN name, age ORDER BY name`)
	if len(result.Rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(result.Rows))
	}
	if _, ok := result.Rows[0]["name"]; !ok {
		t.Fatal("missing 'name' column")
	}
	if _, ok := result.Rows[0]["age"]; !ok {
		t.Fatal("missing 'age' column")
	}
}

func TestWith_SkipAndLimit(t *testing.T) {
	g := setupPersonGraph(t)
	defer g.db.Close()

	result := execQuery(t, g.db, `MATCH (n:Person) WITH n.name AS name SKIP 1 LIMIT 2 RETURN name`)
	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result.Rows))
	}
}

// ---------------------------------------------------------------------------
// Pipeline tests
// ---------------------------------------------------------------------------

func TestPipeline_TwoSegments(t *testing.T) {
	g := setupSocialGraph(t)
	defer g.db.Close()

	pipe := &CypherPipeline{
		Segments: []PipelineSegment{
			{
				Query: &CypherQuery{
					Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
						{Variable: "n", Labels: []string{"Person"}},
					}}},
					With: &WithClause{
						Items: []WithItem{
							{Expr: varRefExpr("n"), Alias: "node"},
						},
					},
				},
			},
			{
				With: &WithClause{
					Items: []WithItem{
						{Expr: varRefExpr("node"), Alias: "n"},
					},
				},
				Query: &CypherQuery{
					Return: ReturnClause{Items: []ReturnItem{
						{Expr: propExpr("n", "name")},
					}},
				},
			},
		},
	}

	result, err := ExecutePipeline(context.Background(), g.db, g.db, pipe, nil)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(result.Rows) != 5 {
		t.Fatalf("expected 5 rows, got %d", len(result.Rows))
	}
}

func TestPipeline_WithUnwind(t *testing.T) {
	g := setupSocialGraph(t)
	defer g.db.Close()

	pipe := &CypherPipeline{
		Segments: []PipelineSegment{
			{
				Query: &CypherQuery{
					Unwind: &UnwindClause{
						Expr: listExpr(litExpr("Alice"), litExpr("Bob")),
						Var:  "searchName",
					},
					Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
						{Variable: "n", Labels: []string{"Person"}},
					}}},
					Where: ep(compExpr(
						propExpr("n", "name"),
						OpEq,
						varRefExpr("searchName"),
					)),
					With: &WithClause{
						Items: []WithItem{
							{Expr: varRefExpr("n"), Alias: "node"},
						},
					},
				},
			},
			{
				With: &WithClause{
					Items: []WithItem{
						{Expr: varRefExpr("node"), Alias: "n"},
					},
				},
				Query: &CypherQuery{
					Return: ReturnClause{Items: []ReturnItem{
						{Expr: propExpr("n", "name")},
					}},
				},
			},
		},
	}

	result, err := ExecutePipeline(context.Background(), g.db, g.db, pipe, nil)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d (rows: %v)", len(result.Rows), result.Rows)
	}
}

func TestPipeline_Empty(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	pipe := &CypherPipeline{Segments: []PipelineSegment{}}
	result, err := ExecutePipeline(context.Background(), db, db, pipe, nil)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(result.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(result.Rows))
	}
}

func TestPipeline_SingleSegment(t *testing.T) {
	g := setupSocialGraph(t)
	defer g.db.Close()

	pipe := &CypherPipeline{
		Segments: []PipelineSegment{
			{
				Query: &CypherQuery{
					Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
						{Variable: "n", Labels: []string{"Person"}},
					}}},
					Return: ReturnClause{Items: []ReturnItem{
						{Expr: propExpr("n", "name")},
					}},
				},
			},
		},
	}

	result, err := ExecutePipeline(context.Background(), g.db, g.db, pipe, nil)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(result.Rows) != 5 {
		t.Fatalf("expected 5 rows, got %d", len(result.Rows))
	}
}

// ---------------------------------------------------------------------------
// WITH + UNWIND combined tests
// ---------------------------------------------------------------------------

func TestWithUnwind_CollectAndUnwind(t *testing.T) {
	g := setupPersonGraph(t)
	defer g.db.Close()

	result := execQuery(t, g.db, `MATCH (n:Person) WITH collect(n.name) AS names UNWIND names AS name RETURN name ORDER BY name`)
	if len(result.Rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(result.Rows))
	}
	if result.Rows[0]["name"].(string) != "Alice" {
		t.Fatalf("first should be Alice, got %v", result.Rows[0]["name"])
	}
}

// ---------------------------------------------------------------------------
// applyWith / applyUnwind unit tests
// ---------------------------------------------------------------------------

func TestApplyWith_Basic(t *testing.T) {
	bindings := []map[string]any{
		{"n": map[string]any{"name": "Alice", "age": 30}},
		{"n": map[string]any{"name": "Bob", "age": 25}},
	}

	wc := &WithClause{
		Items: []WithItem{
			{Expr: propExpr("n", "name"), Alias: "name"},
		},
	}

	result, err := applyWith(wc, nil, bindings)
	if err != nil {
		t.Fatalf("applyWith error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 bindings, got %d", len(result))
	}
	if result[0]["name"] != "Alice" {
		t.Fatalf("expected Alice, got %v", result[0]["name"])
	}
	if _, ok := result[0]["n"]; ok {
		t.Fatal("n should not be in projected binding")
	}
}

func TestApplyWith_WithWhere(t *testing.T) {
	bindings := []map[string]any{
		{"n": map[string]any{"name": "Alice", "age": 30}},
		{"n": map[string]any{"name": "Bob", "age": 25}},
	}

	wc := &WithClause{
		Items: []WithItem{
			{Expr: propExpr("n", "name"), Alias: "name"},
		},
	}

	where := compExpr(propExpr("n", "age"), OpGt, litExpr(28))
	result, err := applyWith(wc, &where, bindings)
	if err != nil {
		t.Fatalf("applyWith error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 binding (Alice only), got %d", len(result))
	}
}

func TestApplyWith_Distinct(t *testing.T) {
	bindings := []map[string]any{
		{"n": map[string]any{"age": 30}},
		{"n": map[string]any{"age": 25}},
		{"n": map[string]any{"age": 30}},
	}

	wc := &WithClause{
		Distinct: true,
		Items: []WithItem{
			{Expr: propExpr("n", "age"), Alias: "age"},
		},
	}

	result, err := applyWith(wc, nil, bindings)
	if err != nil {
		t.Fatalf("applyWith error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 distinct bindings, got %d", len(result))
	}
}

func TestApplyUnwind_Basic(t *testing.T) {
	binding := map[string]any{"key": "value"}
	uc := &UnwindClause{
		Expr: listExpr(litExpr(int64(1)), litExpr(int64(2)), litExpr(int64(3))),
		Var:  "x",
	}

	result, err := applyUnwind(binding, uc)
	if err != nil {
		t.Fatalf("applyUnwind error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 bindings, got %d", len(result))
	}
	if result[0]["x"] != int64(1) {
		t.Fatalf("expected 1, got %v", result[0]["x"])
	}
	if result[0]["key"] != "value" {
		t.Fatalf("expected key=value preserved, got %v", result[0]["key"])
	}
}

func TestApplyUnwind_EmptyList(t *testing.T) {
	binding := map[string]any{}
	uc := &UnwindClause{
		Expr: listExpr(),
		Var:  "x",
	}

	result, err := applyUnwind(binding, uc)
	if err != nil {
		t.Fatalf("applyUnwind error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 bindings, got %d", len(result))
	}
}

func TestApplyUnwind_NonList(t *testing.T) {
	binding := map[string]any{}
	uc := &UnwindClause{
		Expr: litExpr("not a list"),
		Var:  "x",
	}

	_, err := applyUnwind(binding, uc)
	if err == nil {
		t.Fatal("expected error for non-list UNWIND")
	}
}
