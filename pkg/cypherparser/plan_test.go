package cypherparser

import (
	"context"
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage/graphengine"
)

func TestExplain_NodeScan(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	parsed, err := Parse("EXPLAIN MATCH (n:Person) RETURN n")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err := Execute(context.Background(), db, parsed.Read, nil)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("EXPLAIN should return 1 row with plan, got %d", len(result.Rows))
	}
	if len(result.Columns) != 1 || result.Columns[0] != "plan" {
		t.Fatalf("EXPLAIN should have column 'plan', got %v", result.Columns)
	}
}

func TestExplain_SingleHop(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	parsed, err := Parse("EXPLAIN MATCH (a)-[r:KNOWS]->(b) RETURN a, b")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err := Execute(context.Background(), db, parsed.Read, nil)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("EXPLAIN should return 1 row with plan, got %d", len(result.Rows))
	}
}

func TestExplain_VarLength(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	parsed, err := Parse("EXPLAIN MATCH (a)-[:KNOWS*1..3]->(b) RETURN b")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	_, err = Execute(context.Background(), db, parsed.Read, nil)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}
}

func TestExplain_WithUnwind(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	parsed, err := Parse("EXPLAIN UNWIND [1,2,3] AS x MATCH (n) WHERE n.age > x RETURN n")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	_, err = Execute(context.Background(), db, parsed.Read, nil)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}
}

func TestExplain_WithOrderAndLimit(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	parsed, err := Parse("EXPLAIN MATCH (n:Person) RETURN n ORDER BY n.name LIMIT 10")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	_, err = Execute(context.Background(), db, parsed.Read, nil)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}
}

func TestExplain_WithSkip(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	parsed, err := Parse("EXPLAIN MATCH (n:Person) RETURN n SKIP 5")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	_, err = Execute(context.Background(), db, parsed.Read, nil)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}
}

func TestExplain_Distinct(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	parsed, err := Parse("EXPLAIN MATCH (n:Person) RETURN DISTINCT n.age")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	_, err = Execute(context.Background(), db, parsed.Read, nil)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}
}

func TestProfile_ReturnsResult(t *testing.T) {
	g := setupPersonGraph(t)
	defer g.db.Close()

	parsed, err := Parse("PROFILE MATCH (n:Person) RETURN n.name")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err := Execute(context.Background(), g.db, parsed.Read, nil)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("PROFILE should return 1 row with plan+results, got %d rows", len(result.Rows))
	}
	if len(result.Columns) != 2 || result.Columns[0] != "plan" || result.Columns[1] != "results" {
		t.Fatalf("PROFILE should have columns ['plan','results'], got %v", result.Columns)
	}
	cypherResult, ok := result.Rows[0]["results"].(*CypherResult)
	if !ok {
		t.Fatalf("PROFILE results should be *CypherResult, got %T", result.Rows[0]["results"])
	}
	if len(cypherResult.Rows) != 4 {
		t.Fatalf("PROFILE embedded result should have 4 rows, got %d", len(cypherResult.Rows))
	}
}

// ---------------------------------------------------------------------------
// QueryPlan tree building tests
// ---------------------------------------------------------------------------

func TestBuildPlan_AllNodesScan(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
			{Variable: "n"},
		}}},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("n")},
		}},
	}

	plan := buildPlan(q, db)
	if plan.Operator != OpProduceResult {
		t.Fatalf("root should be ProduceResults, got %s", plan.Operator)
	}
	if len(plan.Children) != 1 {
		t.Fatalf("root should have 1 child, got %d", len(plan.Children))
	}

	scan := findOperator(plan, OpAllNodesScan)
	if scan == nil {
		t.Fatal("expected AllNodesScan operator in plan")
	}
}

func TestBuildPlan_LabelScan(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
			{Variable: "n", Labels: []string{"Person"}},
		}}},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("n")},
		}},
	}

	plan := buildPlan(q, db)
	scan := findOperator(plan, OpLabelScan)
	if scan == nil {
		t.Fatal("expected LabelScan operator in plan")
	}
	if !strings.Contains(scan.Details, "Person") {
		t.Fatalf("LabelScan details should contain 'Person', got %s", scan.Details)
	}
}

func TestBuildPlan_SingleHop(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{
			Nodes: []NodePattern{
				{Variable: "a"},
				{Variable: "b"},
			},
			Rels: []RelPattern{
				{Label: "KNOWS", Dir: graphengine.Outgoing},
			},
		}},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("a")},
		}},
	}

	plan := buildPlan(q, db)
	expand := findOperator(plan, OpExpand)
	if expand == nil {
		t.Fatal("expected Expand operator in plan")
	}
	if !strings.Contains(expand.Details, "KNOWS") {
		t.Fatalf("Expand details should contain 'KNOWS', got %s", expand.Details)
	}
}

func TestBuildPlan_VarLength(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{
			Nodes: []NodePattern{
				{Variable: "a"},
				{Variable: "b"},
			},
			Rels: []RelPattern{
				{Label: "KNOWS", VarLength: true, MinHops: 1, MaxHops: 3},
			},
		}},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("b")},
		}},
	}

	plan := buildPlan(q, db)
	vl := findOperator(plan, OpExpandVarLen)
	if vl == nil {
		t.Fatal("expected VarLengthExpand operator in plan")
	}
}

func TestBuildPlan_Filter(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
			{Variable: "n", Props: map[string]any{"name": "Alice"}},
		}}},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("n")},
		}},
	}

	plan := buildPlan(q, db)
	filter := findOperator(plan, OpFilter)
	if filter == nil {
		t.Fatal("expected Filter operator in plan")
	}
}

func TestBuildPlan_Sort(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
			{Variable: "n"},
		}}},
		OrderBy: []OrderItem{
			{Expr: propExpr("n", "name")},
		},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("n")},
		}},
	}

	plan := buildPlan(q, db)
	sort := findOperator(plan, OpSort)
	if sort == nil {
		t.Fatal("expected Sort operator in plan")
	}
}

func TestBuildPlan_Limit(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
			{Variable: "n"},
		}}},
		Limit: intPtr(10),
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("n")},
		}},
	}

	plan := buildPlan(q, db)
	limit := findOperator(plan, OpLimit)
	if limit == nil {
		t.Fatal("expected Limit operator in plan")
	}
}

func TestBuildPlan_Unwind(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Unwind: &UnwindClause{
			Expr: listExpr(litExpr(1), litExpr(2)),
			Var:  "x",
		},
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
			{Variable: "n"},
		}}},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("n")},
		}},
	}

	plan := buildPlan(q, db)
	unwind := findOperator(plan, OpUnwind)
	if unwind == nil {
		t.Fatal("expected Unwind operator in plan")
	}
}

func TestBuildPlan_With(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
			{Variable: "n"},
		}}},
		With: &WithClause{
			Items: []WithItem{
				{Expr: varRefExpr("n"), Alias: "node"},
			},
		},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("node")},
		}},
	}

	plan := buildPlan(q, db)
	with := findOperator(plan, OpWith)
	if with == nil {
		t.Fatal("expected With operator in plan")
	}
}

func TestQueryPlan_String(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
			{Variable: "n", Labels: []string{"Person"}},
		}}},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("n")},
		}},
	}

	plan := BuildExplainPlan(q, db)
	s := plan.String()
	if !strings.Contains(s, "EXPLAIN:") {
		t.Fatalf("plan string should contain 'EXPLAIN:', got:\n%s", s)
	}
	if !strings.Contains(s, "ProduceResults") {
		t.Fatalf("plan string should contain 'ProduceResults', got:\n%s", s)
	}
	if !strings.Contains(s, "LabelScan") {
		t.Fatalf("plan string should contain 'LabelScan', got:\n%s", s)
	}
}

func TestBuildExplainPlan_ProfileFlag(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
			{Variable: "n"},
		}}},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("n")},
		}},
	}

	plan := BuildExplainPlan(q, db)
	if plan.Profile {
		t.Fatal("EXPLAIN plan should not have Profile=true")
	}
	if plan.Result != nil {
		t.Fatal("EXPLAIN plan should not have Result")
	}
}

func TestBuildProfilePlan(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	q := &CypherQuery{
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{
			{Variable: "n"},
		}}},
		Return: ReturnClause{Items: []ReturnItem{
			{Expr: varRefExpr("n")},
		}},
	}

	result := &CypherResult{Columns: []string{"n"}, Rows: []map[string]any{{"n": "test"}}}
	plan := BuildProfilePlan(q, db, result)
	if !plan.Profile {
		t.Fatal("PROFILE plan should have Profile=true")
	}
	if plan.Result == nil {
		t.Fatal("PROFILE plan should have Result")
	}

	s := plan.String()
	if !strings.Contains(s, "PROFILE:") {
		t.Fatalf("plan string should contain 'PROFILE:', got:\n%s", s)
	}
}

// ---------------------------------------------------------------------------
// Plan formatting tests
// ---------------------------------------------------------------------------

func TestFormatExprBrief(t *testing.T) {
	tests := []struct {
		e    Expression
		want string
	}{
		{varRefExpr("x"), "x"},
		{propExpr("n", "name"), "n.name"},
		{funcCallExpr("count", varRefExpr("x")), "count(...)"},
		{litExpr(42), "42"},
		{aggExpr(AggCount, nil), "count(*)"},
	}

	for _, tc := range tests {
		got := formatExprBrief(tc.e)
		if got != tc.want {
			t.Errorf("formatExprBrief(%v) = %q, want %q", tc.e, got, tc.want)
		}
	}
}

func TestFormatOrderBy(t *testing.T) {
	items := []OrderItem{
		{Expr: propExpr("n", "name")},
		{Expr: propExpr("n", "age"), Desc: true},
	}
	got := formatOrderBy(items)
	if !strings.Contains(got, "n.name") || !strings.Contains(got, "DESC") {
		t.Fatalf("formatOrderBy = %q", got)
	}
}

func TestFormatReturn(t *testing.T) {
	rc := ReturnClause{
		Distinct: true,
		Items: []ReturnItem{
			{Expr: varRefExpr("n"), Alias: "node"},
		},
	}
	got := formatReturn(rc)
	if !strings.Contains(got, "DISTINCT") || !strings.Contains(got, "AS node") {
		t.Fatalf("formatReturn = %q", got)
	}
}

func TestFormatWithItems(t *testing.T) {
	wc := &WithClause{
		Distinct: true,
		Items: []WithItem{
			{Expr: propExpr("n", "name"), Alias: "name"},
		},
	}
	got := formatWithItems(wc)
	if !strings.Contains(got, "DISTINCT") || !strings.Contains(got, "AS name") {
		t.Fatalf("formatWithItems = %q", got)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func findOperator(node *PlanNode, op PlanOperator) *PlanNode {
	if node.Operator == op {
		return node
	}
	for _, child := range node.Children {
		if found := findOperator(child, op); found != nil {
			return found
		}
	}
	return nil
}
