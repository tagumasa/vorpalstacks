package cypherparser

import (
	"testing"

	"vorpalstacks/internal/core/storage/graphengine"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func parseMust(t *testing.T, input string) *ParsedCypher {
	t.Helper()
	p, err := Parse(input)
	if err != nil {
		t.Fatalf("parse(%q): %v", input, err)
	}
	return p
}

func parseMustFail(t *testing.T, input string) {
	t.Helper()
	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected parse error for %q", input)
	}
}

func parseRead(t *testing.T, input string) *CypherQuery {
	t.Helper()
	p := parseMust(t, input)
	if p.Read == nil {
		t.Fatalf("expected read query for %q", input)
	}
	return p.Read
}

func parseWrite(t *testing.T, input string) *CypherWrite {
	t.Helper()
	p := parseMust(t, input)
	if p.Write == nil {
		t.Fatalf("expected write query for %q", input)
	}
	return p.Write
}

func parseMerge(t *testing.T, input string) *CypherMerge {
	t.Helper()
	p := parseMust(t, input)
	if p.Merge == nil {
		t.Fatalf("expected merge query for %q", input)
	}
	return p.Merge
}

// ---------------------------------------------------------------------------
// MATCH — basic patterns
// ---------------------------------------------------------------------------

func TestParser_SimpleNode(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN n")
	if len(q.Match.Pattern.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(q.Match.Pattern.Nodes))
	}
	if q.Match.Pattern.Nodes[0].Variable != "n" {
		t.Fatalf("expected var 'n', got %q", q.Match.Pattern.Nodes[0].Variable)
	}
	if len(q.Return.Items) != 1 || q.Return.Items[0].Expr.Kind != ExprVarRef || q.Return.Items[0].Expr.Variable != "n" {
		t.Fatal("unexpected RETURN item")
	}
}

func TestParser_NodeWithLabel(t *testing.T) {
	q := parseRead(t, "MATCH (n:Person) RETURN n")
	node := q.Match.Pattern.Nodes[0]
	if len(node.Labels) != 1 || node.Labels[0] != "Person" {
		t.Fatalf("expected label Person, got %v", node.Labels)
	}
}

func TestParser_NodeWithMultipleLabels(t *testing.T) {
	q := parseRead(t, "MATCH (n:Person:Actor) RETURN n")
	node := q.Match.Pattern.Nodes[0]
	if len(node.Labels) != 2 || node.Labels[0] != "Person" || node.Labels[1] != "Actor" {
		t.Fatalf("expected [Person, Actor], got %v", node.Labels)
	}
}

func TestParser_NodeWithProps(t *testing.T) {
	q := parseRead(t, `MATCH (n {name: "Alice", age: 30}) RETURN n`)
	node := q.Match.Pattern.Nodes[0]
	if node.Props["name"] != "Alice" {
		t.Fatalf("expected name=Alice, got %v", node.Props["name"])
	}
	if node.Props["age"] != int64(30) {
		t.Fatalf("expected age=30, got %v (type %T)", node.Props["age"], node.Props["age"])
	}
}

func TestParser_NodeWithLabelAndProps(t *testing.T) {
	q := parseRead(t, `MATCH (n:Person {name: "Bob"}) RETURN n`)
	node := q.Match.Pattern.Nodes[0]
	if len(node.Labels) != 1 || node.Labels[0] != "Person" {
		t.Fatalf("expected [Person], got %v", node.Labels)
	}
	if node.Props["name"] != "Bob" {
		t.Fatalf("expected name=Bob, got %v", node.Props["name"])
	}
}

// ---------------------------------------------------------------------------
// MATCH — relationship patterns
// ---------------------------------------------------------------------------

func TestParser_RelOutgoing(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[:FOLLOWS]->(b) RETURN a, b")
	if len(q.Match.Pattern.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(q.Match.Pattern.Nodes))
	}
	if len(q.Match.Pattern.Rels) != 1 {
		t.Fatalf("expected 1 rel, got %d", len(q.Match.Pattern.Rels))
	}
	rel := q.Match.Pattern.Rels[0]
	if rel.Label != "FOLLOWS" {
		t.Fatalf("expected label FOLLOWS, got %q", rel.Label)
	}
	if rel.Dir != graphengine.Outgoing {
		t.Fatalf("expected Outgoing, got %d", rel.Dir)
	}
}

func TestParser_RelIncoming(t *testing.T) {
	q := parseRead(t, "MATCH (a)<-[:KNOWS]-(b) RETURN a, b")
	rel := q.Match.Pattern.Rels[0]
	if rel.Dir != graphengine.Incoming {
		t.Fatalf("expected Incoming, got %d", rel.Dir)
	}
	if rel.Label != "KNOWS" {
		t.Fatalf("expected label KNOWS, got %q", rel.Label)
	}
}

func TestParser_RelUndirected(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[:KNOWS]-(b) RETURN a, b")
	rel := q.Match.Pattern.Rels[0]
	if rel.Dir != graphengine.Both {
		t.Fatalf("expected Both, got %d", rel.Dir)
	}
}

func TestParser_RelWithVariable(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[r:FOLLOWS]->(b) RETURN r")
	rel := q.Match.Pattern.Rels[0]
	if rel.Variable != "r" {
		t.Fatalf("expected variable 'r', got %q", rel.Variable)
	}
}

func TestParser_RelWithProps(t *testing.T) {
	q := parseRead(t, `MATCH (a)-[r:FOLLOWS {since: 2020}]->(b) RETURN r`)
	rel := q.Match.Pattern.Rels[0]
	if rel.Props == nil {
		t.Fatal("expected edge props")
	}
	if rel.Props["since"] != int64(2020) {
		t.Fatalf("expected since=2020, got %v", rel.Props["since"])
	}
}

func TestParser_RelVarLength(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[:FOLLOWS*1..3]->(b) RETURN b")
	rel := q.Match.Pattern.Rels[0]
	if !rel.VarLength {
		t.Fatal("expected VarLength=true")
	}
	if rel.MinHops != 1 || rel.MaxHops != 3 {
		t.Fatalf("expected 1..3, got %d..%d", rel.MinHops, rel.MaxHops)
	}
}

func TestParser_RelVarLengthUnbounded(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[:FOLLOWS*2..]->(b) RETURN b")
	rel := q.Match.Pattern.Rels[0]
	if rel.MinHops != 2 || rel.MaxHops != -1 {
		t.Fatalf("expected 2..-1, got %d..%d", rel.MinHops, rel.MaxHops)
	}
}

func TestParser_RelVarLengthFixed(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[:FOLLOWS*3]->(b) RETURN b")
	rel := q.Match.Pattern.Rels[0]
	if rel.MinHops != 3 || rel.MaxHops != 3 {
		t.Fatalf("expected 3..3, got %d..%d", rel.MinHops, rel.MaxHops)
	}
}

func TestParser_ThreeNodePath(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[:FOLLOWS]->(b)-[:LIKES]->(c) RETURN a, b, c")
	if len(q.Match.Pattern.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(q.Match.Pattern.Nodes))
	}
	if len(q.Match.Pattern.Rels) != 2 {
		t.Fatalf("expected 2 rels, got %d", len(q.Match.Pattern.Rels))
	}
}

func TestParser_AnyEdgeType(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[r]->(b) RETURN r")
	rel := q.Match.Pattern.Rels[0]
	if rel.Label != "" {
		t.Fatalf("expected empty label, got %q", rel.Label)
	}
}

// ---------------------------------------------------------------------------
// WHERE clause
// ---------------------------------------------------------------------------

func TestParser_WhereSimple(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE n.age > 25 RETURN n")
	if q.Where == nil {
		t.Fatal("expected WHERE")
	}
	if q.Where.Kind != ExprComparison {
		t.Fatalf("expected Comparison, got %d", q.Where.Kind)
	}
	if q.Where.Op != OpGt {
		t.Fatalf("expected >, got %d", q.Where.Op)
	}
	if q.Where.Left.Kind != ExprPropAccess || q.Where.Left.Property != "age" {
		t.Fatal("expected n.age on left")
	}
	if q.Where.Right.Kind != ExprLiteral || q.Where.Right.LitValue != int64(25) {
		t.Fatal("expected 25 on right")
	}
}

func TestParser_WhereAnd(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE n.age > 25 AND n.age < 32 RETURN n")
	if q.Where.Kind != ExprAnd {
		t.Fatalf("expected And, got %d", q.Where.Kind)
	}
	if len(q.Where.Operands) != 2 {
		t.Fatalf("expected 2 operands, got %d", len(q.Where.Operands))
	}
}

func TestParser_WhereOr(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE n.age > 30 OR n.name = 'Alice' RETURN n")
	if q.Where.Kind != ExprOr {
		t.Fatalf("expected Or, got %d", q.Where.Kind)
	}
}

func TestParser_WhereNot(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE NOT n.active RETURN n")
	if q.Where.Kind != ExprNot {
		t.Fatalf("expected Not, got %d", q.Where.Kind)
	}
}

func TestParser_WhereNotCombined(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE NOT n.active AND n.age > 20 RETURN n")
	if q.Where.Kind != ExprAnd {
		t.Fatalf("expected And, got %d", q.Where.Kind)
	}
	if q.Where.Operands[0].Kind != ExprNot {
		t.Fatalf("expected first operand to be Not, got %d", q.Where.Operands[0].Kind)
	}
}

func TestParser_WhereIn(t *testing.T) {
	q := parseRead(t, `MATCH (n) WHERE n.name IN ["Alice", "Bob"] RETURN n`)
	if q.Where.Kind != ExprIn {
		t.Fatalf("expected In, got %d", q.Where.Kind)
	}
}

func TestParser_WhereStartsWith(t *testing.T) {
	q := parseRead(t, `MATCH (n) WHERE n.name STARTS WITH "A" RETURN n`)
	if q.Where.Kind != ExprStartsWith {
		t.Fatalf("expected StartsWith, got %d", q.Where.Kind)
	}
}

func TestParser_WhereEndsWith(t *testing.T) {
	q := parseRead(t, `MATCH (n) WHERE n.name ENDS WITH "e" RETURN n`)
	if q.Where.Kind != ExprEndsWith {
		t.Fatalf("expected EndsWith, got %d", q.Where.Kind)
	}
}

func TestParser_WhereContains(t *testing.T) {
	q := parseRead(t, `MATCH (n) WHERE n.name CONTAINS "li" RETURN n`)
	if q.Where.Kind != ExprContains {
		t.Fatalf("expected Contains, got %d", q.Where.Kind)
	}
}

func TestParser_WhereIsNull(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE n.email IS NULL RETURN n")
	if q.Where.Kind != ExprIsNull {
		t.Fatalf("expected IsNull, got %d", q.Where.Kind)
	}
}

func TestParser_WhereIsNotNull(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE n.email IS NOT NULL RETURN n")
	if q.Where.Kind != ExprIsNotNull {
		t.Fatalf("expected IsNotNull, got %d", q.Where.Kind)
	}
}

func TestParser_WhereRegexMatch(t *testing.T) {
	q := parseRead(t, `MATCH (n) WHERE n.name =~ "A.*" RETURN n`)
	if q.Where.Kind != ExprRegexMatch {
		t.Fatalf("expected RegexMatch, got %d", q.Where.Kind)
	}
}

// ---------------------------------------------------------------------------
// RETURN clause
// ---------------------------------------------------------------------------

func TestParser_ReturnMultiple(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[:FOLLOWS]->(b) RETURN a, b")
	if len(q.Return.Items) != 2 {
		t.Fatalf("expected 2 return items, got %d", len(q.Return.Items))
	}
}

func TestParser_ReturnAlias(t *testing.T) {
	q := parseRead(t, `MATCH (n {name: "Alice"}) RETURN n.name AS person_name`)
	if len(q.Return.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(q.Return.Items))
	}
	item := q.Return.Items[0]
	if item.Alias != "person_name" {
		t.Fatalf("expected alias 'person_name', got %q", item.Alias)
	}
	if item.Expr.Kind != ExprPropAccess {
		t.Fatalf("expected PropAccess, got %d", item.Expr.Kind)
	}
}

func TestParser_ReturnDistinct(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN DISTINCT n.name")
	if !q.Return.Distinct {
		t.Fatal("expected Distinct=true")
	}
}

// ---------------------------------------------------------------------------
// ORDER BY / SKIP / LIMIT
// ---------------------------------------------------------------------------

func TestParser_OrderBy(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN n ORDER BY n.age DESC")
	if len(q.OrderBy) != 1 {
		t.Fatalf("expected 1 order item, got %d", len(q.OrderBy))
	}
	if !q.OrderBy[0].Desc {
		t.Fatal("expected Desc=true")
	}
}

func TestParser_OrderByAsc(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN n ORDER BY n.age ASC")
	if q.OrderBy[0].Desc {
		t.Fatal("expected Desc=false")
	}
}

func TestParser_OrderByMultiple(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN n ORDER BY n.age, n.name DESC")
	if len(q.OrderBy) != 2 {
		t.Fatalf("expected 2 order items, got %d", len(q.OrderBy))
	}
	if q.OrderBy[0].Desc {
		t.Fatal("first item should be ASC")
	}
	if !q.OrderBy[1].Desc {
		t.Fatal("second item should be DESC")
	}
}

func TestParser_SkipLimit(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN n SKIP 10 LIMIT 5")
	if q.Skip == nil || *q.Skip != 10 {
		t.Fatalf("expected Skip=10, got %d", derefInt(q.Skip))
	}
	if q.Limit == nil || *q.Limit != 5 {
		t.Fatalf("expected Limit=5, got %d", derefInt(q.Limit))
	}
}

func TestParser_LimitOnly(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN n LIMIT 3")
	if q.Limit == nil || *q.Limit != 3 {
		t.Fatalf("expected Limit=3, got %d", derefInt(q.Limit))
	}
	if q.Skip != nil {
		t.Fatalf("expected Skip=nil, got %d", *q.Skip)
	}
}

// ---------------------------------------------------------------------------
// EXPLAIN / PROFILE
// ---------------------------------------------------------------------------

func TestParser_Explain(t *testing.T) {
	q := parseRead(t, "EXPLAIN MATCH (n) RETURN n")
	if q.Explain != ExplainOnly {
		t.Fatalf("expected ExplainOnly, got %d", q.Explain)
	}
}

func TestParser_Profile(t *testing.T) {
	q := parseRead(t, "PROFILE MATCH (n) RETURN n")
	if q.Explain != ExplainProfile {
		t.Fatalf("expected ExplainProfile, got %d", q.Explain)
	}
}

// ---------------------------------------------------------------------------
// WITH clause
// ---------------------------------------------------------------------------

func TestParser_With(t *testing.T) {
	q := parseRead(t, "WITH n AS node MATCH (node)-[:FOLLOWS]->(b) RETURN b")
	if q.With == nil {
		t.Fatal("expected WITH clause")
	}
	if len(q.With.Items) != 1 {
		t.Fatalf("expected 1 WITH item, got %d", len(q.With.Items))
	}
	if q.With.Items[0].Alias != "node" {
		t.Fatalf("expected alias 'node', got %q", q.With.Items[0].Alias)
	}
}

func TestParser_WithDistinct(t *testing.T) {
	q := parseRead(t, "WITH DISTINCT n.name AS name RETURN name")
	if q.With == nil {
		t.Fatal("expected WITH clause")
	}
	if !q.With.Distinct {
		t.Fatal("expected Distinct=true")
	}
}

func TestParser_WithMultiple(t *testing.T) {
	q := parseRead(t, "WITH n.name AS name, n.age AS age RETURN name, age")
	if len(q.With.Items) != 2 {
		t.Fatalf("expected 2 WITH items, got %d", len(q.With.Items))
	}
}

// ---------------------------------------------------------------------------
// UNWIND clause
// ---------------------------------------------------------------------------

func TestParser_Unwind(t *testing.T) {
	q := parseRead(t, `UNWIND [1, 2, 3] AS x MATCH (n {id: x}) RETURN n`)
	if q.Unwind == nil {
		t.Fatal("expected UNWIND clause")
	}
	if q.Unwind.Var != "x" {
		t.Fatalf("expected var 'x', got %q", q.Unwind.Var)
	}
	if q.Unwind.Expr.Kind != ExprListLiteral {
		t.Fatalf("expected ListLiteral, got %d", q.Unwind.Expr.Kind)
	}
	if len(q.Unwind.Expr.ListItems) != 3 {
		t.Fatalf("expected 3 list items, got %d", len(q.Unwind.Expr.ListItems))
	}
}

// ---------------------------------------------------------------------------
// SET clause
// ---------------------------------------------------------------------------

func TestParser_Set(t *testing.T) {
	q := parseRead(t, `MATCH (n) SET n.name = "Alice" RETURN n`)
	if len(q.Set) != 1 {
		t.Fatalf("expected 1 SET item, got %d", len(q.Set))
	}
	item := q.Set[0]
	if item.Variable != "n" || item.Property != "name" {
		t.Fatalf("expected n.name, got %s.%s", item.Variable, item.Property)
	}
	if item.Value.Kind != ExprLiteral || item.Value.LitValue != "Alice" {
		t.Fatal("expected value 'Alice'")
	}
}

func TestParser_SetMultiple(t *testing.T) {
	q := parseRead(t, `MATCH (n) SET n.name = "Alice", n.age = 30 RETURN n`)
	if len(q.Set) != 2 {
		t.Fatalf("expected 2 SET items, got %d", len(q.Set))
	}
}

// ---------------------------------------------------------------------------
// DELETE clause
// ---------------------------------------------------------------------------

func TestParser_Delete(t *testing.T) {
	q := parseRead(t, "MATCH (n:Old) DELETE n")
	if len(q.Delete) != 1 || q.Delete[0] != "n" {
		t.Fatalf("expected [n], got %v", q.Delete)
	}
}

func TestParser_DeleteMultiple(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[r]->(b) DELETE a, r, b")
	if len(q.Delete) != 3 {
		t.Fatalf("expected 3 items, got %d", len(q.Delete))
	}
}

// ---------------------------------------------------------------------------
// REMOVE clause
// ---------------------------------------------------------------------------

func TestParser_RemoveLabel(t *testing.T) {
	q := parseRead(t, "MATCH (n) REMOVE n:Temp RETURN n")
	if len(q.Remove) != 1 {
		t.Fatalf("expected 1 REMOVE item, got %d", len(q.Remove))
	}
	item := q.Remove[0]
	if item.Kind != RemoveLabel || item.Label != "Temp" {
		t.Fatalf("expected RemoveLabel Temp, got Kind=%d Label=%q", item.Kind, item.Label)
	}
}

func TestParser_RemoveProperty(t *testing.T) {
	q := parseRead(t, "MATCH (n) REMOVE n.temp RETURN n")
	if len(q.Remove) != 1 {
		t.Fatalf("expected 1 REMOVE item, got %d", len(q.Remove))
	}
	item := q.Remove[0]
	if item.Kind != RemoveProperty || item.Property != "temp" {
		t.Fatalf("expected RemoveProperty temp, got Kind=%d Property=%q", item.Kind, item.Property)
	}
}

// ---------------------------------------------------------------------------
// OPTIONAL MATCH
// ---------------------------------------------------------------------------

func TestParser_OptionalMatch(t *testing.T) {
	q := parseRead(t, "MATCH (a)-[:FOLLOWS]->(b) OPTIONAL MATCH (b)-[:LIKES]->(c) RETURN a, b, c")
	if q.OptionalMatch == nil {
		t.Fatal("expected OPTIONAL MATCH")
	}
	if len(q.OptionalMatch.Pattern.Nodes) != 2 {
		t.Fatalf("expected 2 nodes in optional pattern, got %d", len(q.OptionalMatch.Pattern.Nodes))
	}
}

func TestParser_OptionalMatchWithWhere(t *testing.T) {
	q := parseRead(t, "MATCH (a) OPTIONAL MATCH (b) WHERE b.active = true RETURN a, b")
	if q.OptionalWhere == nil {
		t.Fatal("expected OPTIONAL WHERE")
	}
	if q.OptionalWhere.Kind != ExprComparison {
		t.Fatalf("expected Comparison, got %d", q.OptionalWhere.Kind)
	}
}

// ---------------------------------------------------------------------------
// Multiple MATCH clauses
// ---------------------------------------------------------------------------

func TestParser_MultipleMatch(t *testing.T) {
	q := parseRead(t, "MATCH (a:Person) MATCH (b:Movie) WHERE a.name = b.title RETURN a, b")
	if len(q.Matches) != 2 {
		t.Fatalf("expected 2 MATCH clauses, got %d", len(q.Matches))
	}
	if q.Matches[0].Pattern.Nodes[0].Labels[0] != "Person" {
		t.Fatal("first MATCH should be Person")
	}
	if q.Matches[1].Pattern.Nodes[0].Labels[0] != "Movie" {
		t.Fatal("second MATCH should be Movie")
	}
}

// ---------------------------------------------------------------------------
// Expression tests
// ---------------------------------------------------------------------------

func TestParser_ExprLiteralString(t *testing.T) {
	q := parseRead(t, `RETURN "hello"`)
	if q.Return.Items[0].Expr.Kind != ExprLiteral || q.Return.Items[0].Expr.LitValue != "hello" {
		t.Fatal("expected string literal")
	}
}

func TestParser_ExprLiteralInt(t *testing.T) {
	q := parseRead(t, "RETURN 42")
	if q.Return.Items[0].Expr.Kind != ExprLiteral || q.Return.Items[0].Expr.LitValue != int64(42) {
		t.Fatalf("expected int 42, got %v", q.Return.Items[0].Expr.LitValue)
	}
}

func TestParser_ExprLiteralFloat(t *testing.T) {
	q := parseRead(t, "RETURN 3.14")
	if q.Return.Items[0].Expr.Kind != ExprLiteral || q.Return.Items[0].Expr.LitValue != 3.14 {
		t.Fatalf("expected float 3.14, got %v", q.Return.Items[0].Expr.LitValue)
	}
}

func TestParser_ExprLiteralBool(t *testing.T) {
	q := parseRead(t, "RETURN true")
	if q.Return.Items[0].Expr.Kind != ExprLiteral || q.Return.Items[0].Expr.LitValue != true {
		t.Fatal("expected bool true")
	}
}

func TestParser_ExprLiteralNull(t *testing.T) {
	q := parseRead(t, "RETURN null")
	if q.Return.Items[0].Expr.Kind != ExprLiteral || q.Return.Items[0].Expr.LitValue != nil {
		t.Fatal("expected null")
	}
}

func TestParser_ExprVarRef(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN n")
	if q.Return.Items[0].Expr.Kind != ExprVarRef || q.Return.Items[0].Expr.Variable != "n" {
		t.Fatal("expected var ref 'n'")
	}
}

func TestParser_ExprPropAccess(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN n.name")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprPropAccess {
		t.Fatalf("expected PropAccess, got %d", item.Expr.Kind)
	}
	if item.Expr.Object != "n" || item.Expr.Property != "name" {
		t.Fatalf("expected n.name, got %s.%s", item.Expr.Object, item.Expr.Property)
	}
}

func TestParser_ExprFuncCall(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN coalesce(n.name, n.email, 'unknown')")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprFuncCall {
		t.Fatalf("expected FuncCall, got %d", item.Expr.Kind)
	}
	if item.Expr.FuncName != "coalesce" {
		t.Fatalf("expected coalesce, got %q", item.Expr.FuncName)
	}
	if len(item.Expr.Args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(item.Expr.Args))
	}
}

func TestParser_ExprFuncCallMultipleArgs(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN coalesce(n.name, n.email, 'unknown')")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprFuncCall {
		t.Fatalf("expected FuncCall, got %d", item.Expr.Kind)
	}
	if len(item.Expr.Args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(item.Expr.Args))
	}
}

func TestParser_ExprParenthesised(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN (n.age + 1) * 2")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprArithMul {
		t.Fatalf("expected ArithMul, got %d", item.Expr.Kind)
	}
	if item.Expr.ArithLeft.Kind != ExprArithAdd {
		t.Fatalf("expected left to be ArithAdd, got %d", item.Expr.ArithLeft.Kind)
	}
}

func TestParser_ExprParamRef(t *testing.T) {
	q := parseRead(t, "MATCH (n {name: $name}) RETURN n")
	node := q.Match.Pattern.Nodes[0]
	pr, ok := node.Props["name"].(paramRef)
	if !ok {
		t.Fatalf("expected paramRef, got %T", node.Props["name"])
	}
	if string(pr) != "name" {
		t.Fatalf("expected param 'name', got %q", string(pr))
	}
}

// ---------------------------------------------------------------------------
// Arithmetic expressions
// ---------------------------------------------------------------------------

func TestParser_ExprAddition(t *testing.T) {
	q := parseRead(t, "RETURN 1 + 2")
	if q.Return.Items[0].Expr.Kind != ExprArithAdd {
		t.Fatalf("expected ArithAdd, got %d", q.Return.Items[0].Expr.Kind)
	}
}

func TestParser_ExprSubtraction(t *testing.T) {
	q := parseRead(t, "RETURN 10 - 3")
	if q.Return.Items[0].Expr.Kind != ExprArithSub {
		t.Fatalf("expected ArithSub, got %d", q.Return.Items[0].Expr.Kind)
	}
}

func TestParser_ExprMultiplication(t *testing.T) {
	q := parseRead(t, "RETURN 3 * 4")
	if q.Return.Items[0].Expr.Kind != ExprArithMul {
		t.Fatalf("expected ArithMul, got %d", q.Return.Items[0].Expr.Kind)
	}
}

func TestParser_ExprDivision(t *testing.T) {
	q := parseRead(t, "RETURN 10 / 3")
	if q.Return.Items[0].Expr.Kind != ExprArithDiv {
		t.Fatalf("expected ArithDiv, got %d", q.Return.Items[0].Expr.Kind)
	}
}

func TestParser_ExprModulo(t *testing.T) {
	q := parseRead(t, "RETURN 10 % 3")
	if q.Return.Items[0].Expr.Kind != ExprArithMod {
		t.Fatalf("expected ArithMod, got %d", q.Return.Items[0].Expr.Kind)
	}
}

func TestParser_ExprArithmeticPrecedence(t *testing.T) {
	q := parseRead(t, "RETURN 2 + 3 * 4")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprArithAdd {
		t.Fatalf("expected ArithAdd (2 + (3*4)), got %d", item.Expr.Kind)
	}
	if item.Expr.ArithRight.Kind != ExprArithMul {
		t.Fatal("expected right side to be multiplication")
	}
}

func TestParser_ExprUnaryMinus(t *testing.T) {
	q := parseRead(t, "RETURN -5")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprArithSub {
		t.Fatalf("expected ArithSub (0 - 5), got %d", item.Expr.Kind)
	}
}

// ---------------------------------------------------------------------------
// Comparison operators
// ---------------------------------------------------------------------------

func TestParser_ExprEq(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE n.name = 'Alice' RETURN n")
	if q.Where.Op != OpEq {
		t.Fatalf("expected Eq, got %d", q.Where.Op)
	}
}

func TestParser_ExprNeq(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE n.name != 'Bob' RETURN n")
	if q.Where.Op != OpNeq {
		t.Fatalf("expected Neq, got %d", q.Where.Op)
	}
}

func TestParser_ExprLt(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE n.age < 30 RETURN n")
	if q.Where.Op != OpLt {
		t.Fatalf("expected Lt, got %d", q.Where.Op)
	}
}

func TestParser_ExprGte(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE n.age >= 30 RETURN n")
	if q.Where.Op != OpGte {
		t.Fatalf("expected Gte, got %d", q.Where.Op)
	}
}

func TestParser_ExprLte(t *testing.T) {
	q := parseRead(t, "MATCH (n) WHERE n.age <= 30 RETURN n")
	if q.Where.Op != OpLte {
		t.Fatalf("expected Lte, got %d", q.Where.Op)
	}
}

// ---------------------------------------------------------------------------
// List and map literals
// ---------------------------------------------------------------------------

func TestParser_ExprListLiteral(t *testing.T) {
	q := parseRead(t, `RETURN [1, 2, 3]`)
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprListLiteral {
		t.Fatalf("expected ListLiteral, got %d", item.Expr.Kind)
	}
	if len(item.Expr.ListItems) != 3 {
		t.Fatalf("expected 3 items, got %d", len(item.Expr.ListItems))
	}
}

func TestParser_ExprListLiteralEmpty(t *testing.T) {
	q := parseRead(t, "RETURN []")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprListLiteral {
		t.Fatalf("expected ListLiteral, got %d", item.Expr.Kind)
	}
	if len(item.Expr.ListItems) != 0 {
		t.Fatalf("expected 0 items, got %d", len(item.Expr.ListItems))
	}
}

func TestParser_ExprMapLiteral(t *testing.T) {
	q := parseRead(t, `RETURN {name: "Alice", age: 30}`)
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprMapLiteral {
		t.Fatalf("expected MapLiteral, got %d", item.Expr.Kind)
	}
	if len(item.Expr.MapPairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(item.Expr.MapPairs))
	}
}

// ---------------------------------------------------------------------------
// Aggregation functions
// ---------------------------------------------------------------------------

func TestParser_ExprCountStar(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN count(*)")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprAggregation {
		t.Fatalf("expected Aggregation, got %d", item.Expr.Kind)
	}
	if item.Expr.AggFunc != AggCount {
		t.Fatalf("expected AggCount, got %d", item.Expr.AggFunc)
	}
	if item.Expr.AggArg != nil {
		t.Fatal("expected nil arg for COUNT(*)")
	}
}

func TestParser_ExprCountDistinct(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN count(DISTINCT n.name)")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprAggregation {
		t.Fatalf("expected Aggregation, got %d", item.Expr.Kind)
	}
	if item.Expr.AggArg == nil {
		t.Fatal("expected non-nil arg")
	}
}

func TestParser_ExprSum(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN sum(n.age)")
	item := q.Return.Items[0]
	if item.Expr.AggFunc != AggSum {
		t.Fatalf("expected AggSum, got %d", item.Expr.AggFunc)
	}
}

func TestParser_ExprAvg(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN avg(n.age)")
	item := q.Return.Items[0]
	if item.Expr.AggFunc != AggAvg {
		t.Fatalf("expected AggAvg, got %d", item.Expr.AggFunc)
	}
}

func TestParser_ExprMin(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN min(n.age)")
	item := q.Return.Items[0]
	if item.Expr.AggFunc != AggMin {
		t.Fatalf("expected AggMin, got %d", item.Expr.AggFunc)
	}
}

func TestParser_ExprMax(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN max(n.age)")
	item := q.Return.Items[0]
	if item.Expr.AggFunc != AggMax {
		t.Fatalf("expected AggMax, got %d", item.Expr.AggFunc)
	}
}

func TestParser_ExprCollect(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN collect(n.name)")
	item := q.Return.Items[0]
	if item.Expr.AggFunc != AggCollect {
		t.Fatalf("expected AggCollect, got %d", item.Expr.AggFunc)
	}
}

// ---------------------------------------------------------------------------
// CASE expression
// ---------------------------------------------------------------------------

func TestParser_ExprCaseSimple(t *testing.T) {
	q := parseRead(t, `MATCH (n) RETURN CASE WHEN n.age > 30 THEN "senior" ELSE "junior" END`)
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprCase {
		t.Fatalf("expected Case, got %d", item.Expr.Kind)
	}
	if item.Expr.CaseOperand != nil {
		t.Fatal("expected nil operand for simple CASE")
	}
	if len(item.Expr.CaseWhens) != 1 {
		t.Fatalf("expected 1 when, got %d", len(item.Expr.CaseWhens))
	}
	if item.Expr.CaseElse == nil {
		t.Fatal("expected ELSE clause")
	}
}

func TestParser_ExprCaseWithOperand(t *testing.T) {
	q := parseRead(t, `MATCH (n) RETURN CASE n.status WHEN "active" THEN 1 WHEN "inactive" THEN 0 ELSE -1 END`)
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprCase {
		t.Fatalf("expected Case, got %d", item.Expr.Kind)
	}
	if item.Expr.CaseOperand == nil {
		t.Fatal("expected non-nil operand")
	}
	if len(item.Expr.CaseWhens) != 2 {
		t.Fatalf("expected 2 whens, got %d", len(item.Expr.CaseWhens))
	}
}

// ---------------------------------------------------------------------------
// EXISTS { pattern }
// ---------------------------------------------------------------------------

func TestParser_ExprExists(t *testing.T) {
	q := parseRead(t, "MATCH (a) RETURN EXISTS { (a)-[:FOLLOWS]->(b) }")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprExists {
		t.Fatalf("expected Exists, got %d", item.Expr.Kind)
	}
	if item.Expr.ExistsPattern == nil {
		t.Fatal("expected pattern")
	}
}

// ---------------------------------------------------------------------------
// LABELS() and PROPERTIES()
// ---------------------------------------------------------------------------

func TestParser_ExprLabels(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN labels(n)")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprLabels {
		t.Fatalf("expected Labels, got %d", item.Expr.Kind)
	}
}

func TestParser_ExprProperties(t *testing.T) {
	q := parseRead(t, "MATCH (n) RETURN properties(n)")
	item := q.Return.Items[0]
	if item.Expr.Kind != ExprProperties {
		t.Fatalf("expected Properties, got %d", item.Expr.Kind)
	}
}

// ---------------------------------------------------------------------------
// CREATE queries
// ---------------------------------------------------------------------------

func TestParser_CreateSingleNode(t *testing.T) {
	w := parseWrite(t, `CREATE (n:Person {name: "Alice", age: 30}) RETURN n`)
	if len(w.Creates) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(w.Creates))
	}
	node := w.Creates[0].Nodes[0]
	if node.Variable != "n" {
		t.Fatalf("expected var 'n', got %q", node.Variable)
	}
	if len(node.Labels) != 1 || node.Labels[0] != "Person" {
		t.Fatalf("expected [Person], got %v", node.Labels)
	}
}

func TestParser_CreateNodeWithEdge(t *testing.T) {
	w := parseWrite(t, `CREATE (a:Person {name: "Alice"})-[:FOLLOWS]->(b:Person {name: "Bob"}) RETURN a, b`)
	if len(w.Creates) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(w.Creates))
	}
	if len(w.Creates[0].Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(w.Creates[0].Nodes))
	}
	if len(w.Creates[0].Rels) != 1 {
		t.Fatalf("expected 1 rel, got %d", len(w.Creates[0].Rels))
	}
}

func TestParser_CreateMultiplePatterns(t *testing.T) {
	w := parseWrite(t, `CREATE (a:Person {name: "Alice"}), (b:Person {name: "Bob"})`)
	if len(w.Creates) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(w.Creates))
	}
}

func TestParser_CreateNoReturn(t *testing.T) {
	w := parseWrite(t, `CREATE (n:Movie {title: "The Matrix"})`)
	if w.Return != nil {
		t.Fatal("expected nil RETURN")
	}
}

func TestParser_CreateReturnDistinct(t *testing.T) {
	w := parseWrite(t, `CREATE (n:Person) RETURN DISTINCT n`)
	if w.Return == nil {
		t.Fatal("expected RETURN")
	}
	if !w.Return.Distinct {
		t.Fatal("expected Distinct=true")
	}
}

// ---------------------------------------------------------------------------
// MERGE queries
// ---------------------------------------------------------------------------

func TestParser_MergeBasic(t *testing.T) {
	m := parseMerge(t, `MERGE (n:Person {name: "Alice"})`)
	if m.Pattern.Variable != "n" {
		t.Fatalf("expected var 'n', got %q", m.Pattern.Variable)
	}
	if len(m.Pattern.Labels) != 1 || m.Pattern.Labels[0] != "Person" {
		t.Fatalf("expected [Person], got %v", m.Pattern.Labels)
	}
}

func TestParser_MergeOnCreateSet(t *testing.T) {
	m := parseMerge(t, `MERGE (n:Person {name: "Alice"}) ON CREATE SET n.createdAt = 2024`)
	if len(m.OnCreateSet) != 1 {
		t.Fatalf("expected 1 OnCreateSet item, got %d", len(m.OnCreateSet))
	}
	item := m.OnCreateSet[0]
	if item.Variable != "n" || item.Property != "createdAt" {
		t.Fatalf("expected n.createdAt, got %s.%s", item.Variable, item.Property)
	}
}

func TestParser_MergeOnMatchSet(t *testing.T) {
	m := parseMerge(t, `MERGE (n:Person {name: "Alice"}) ON MATCH SET n.lastSeen = 2024`)
	if len(m.OnMatchSet) != 1 {
		t.Fatalf("expected 1 OnMatchSet item, got %d", len(m.OnMatchSet))
	}
}

func TestParser_MergeOnCreateAndMatchSet(t *testing.T) {
	m := parseMerge(t, `MERGE (n:Person {name: "Alice"}) ON CREATE SET n.createdAt = 2024 ON MATCH SET n.lastSeen = 2024`)
	if len(m.OnCreateSet) != 1 || len(m.OnMatchSet) != 1 {
		t.Fatalf("expected 1 OnCreateSet and 1 OnMatchSet, got %d and %d", len(m.OnCreateSet), len(m.OnMatchSet))
	}
}

func TestParser_MergeWithReturn(t *testing.T) {
	m := parseMerge(t, `MERGE (n:Person {name: "Alice"}) RETURN n`)
	if m.Return == nil {
		t.Fatal("expected RETURN")
	}
}

func TestParser_MergeRequiresLabel(t *testing.T) {
	parseMustFail(t, `MERGE (n {name: "Alice"})`)
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestParser_NoReturnIsValid(t *testing.T) {
	q := parseRead(t, "MATCH (n:Old) DELETE n")
	if len(q.Delete) != 1 {
		t.Fatalf("expected 1 DELETE item, got %d", len(q.Delete))
	}
}

func TestParser_ErrorUnexpectedToken(t *testing.T) {
	parseMustFail(t, "RETURN n MATCH (x)")
}

func TestParser_ErrorBadRelPattern(t *testing.T) {
	parseMustFail(t, "MATCH (a)FOO(b) RETURN a, b")
}

func TestParser_ErrorUnterminatedString(t *testing.T) {
	parseMustFail(t, `MATCH (n {name: "Alice}) RETURN n`)
}

func TestParser_ErrorInvalidQuery(t *testing.T) {
	parseMustFail(t, "INVALID QUERY")
}

func TestParser_ErrorMergeNoLabel(t *testing.T) {
	parseMustFail(t, "MERGE (n) RETURN n")
}

func TestParser_ErrorBadAfterMerge(t *testing.T) {
	parseMustFail(t, "MERGE (n:Person) ON DELETE SET n.x = 1")
}

// ---------------------------------------------------------------------------
// Expression value in property maps (our extension)
// ---------------------------------------------------------------------------

func TestParser_PropMapExprValue(t *testing.T) {
	q := parseRead(t, "MATCH (n {age: 20 + 10}) RETURN n")
	node := q.Match.Pattern.Nodes[0]
	val, ok := node.Props["age"].(Expression)
	if !ok {
		t.Fatalf("expected Expression for age, got %T", node.Props["age"])
	}
	if val.Kind != ExprArithAdd {
		t.Fatalf("expected ArithAdd, got %d", val.Kind)
	}
}

// ---------------------------------------------------------------------------
// Complex real-world queries
// ---------------------------------------------------------------------------

func TestParser_ComplexQuery(t *testing.T) {
	input := `MATCH (a:Person {name: "Alice"})-[:FOLLOWS*1..3]->(b:Person) WHERE b.age > 25 AND (b.name STARTS WITH "A" OR b.name STARTS WITH "B") RETURN DISTINCT b.name, b.age ORDER BY b.age DESC SKIP 5 LIMIT 10`
	q := parseRead(t, input)

	if len(q.Match.Pattern.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(q.Match.Pattern.Nodes))
	}
	rel := q.Match.Pattern.Rels[0]
	if !rel.VarLength || rel.MinHops != 1 || rel.MaxHops != 3 {
		t.Fatalf("expected *1..3, got %v..%v (var=%v)", rel.MinHops, rel.MaxHops, rel.VarLength)
	}
	if q.Where == nil {
		t.Fatal("expected WHERE")
	}
	if !q.Return.Distinct {
		t.Fatal("expected DISTINCT")
	}
	if len(q.Return.Items) != 2 {
		t.Fatalf("expected 2 return items, got %d", len(q.Return.Items))
	}
	if q.Skip == nil || *q.Skip != 5 || q.Limit == nil || *q.Limit != 10 {
		t.Fatalf("expected SKIP 5 LIMIT 10, got %d/%d", derefInt(q.Skip), derefInt(q.Limit))
	}
}

func TestParser_ComplexCreateQuery(t *testing.T) {
	input := `CREATE (a:Person {name: "Alice", age: 30})-[:KNOWS {since: 2020}]->(b:Person {name: "Bob", age: 25}) RETURN a.name, b.name AS friend`
	w := parseWrite(t, input)

	if len(w.Creates) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(w.Creates))
	}
	rel := w.Creates[0].Rels[0]
	if rel.Label != "KNOWS" {
		t.Fatalf("expected KNOWS, got %q", rel.Label)
	}
	if rel.Props["since"] != int64(2020) {
		t.Fatalf("expected since=2020, got %v", rel.Props["since"])
	}
	if len(w.Return.Items) != 2 {
		t.Fatalf("expected 2 return items, got %d", len(w.Return.Items))
	}
	if w.Return.Items[1].Alias != "friend" {
		t.Fatalf("expected alias 'friend', got %q", w.Return.Items[1].Alias)
	}
}

func TestParser_WithMatchReturn(t *testing.T) {
	input := `MATCH (n:Person) WITH n.name AS name WHERE name STARTS WITH "A" RETURN name`
	q := parseRead(t, input)

	if q.With == nil {
		t.Fatal("expected WITH")
	}
	if q.With.Items[0].Alias != "name" {
		t.Fatalf("expected alias 'name', got %q", q.With.Items[0].Alias)
	}
	if len(q.Match.Pattern.Nodes) != 1 {
		t.Fatalf("expected 1 node in MATCH, got %d", len(q.Match.Pattern.Nodes))
	}
}
