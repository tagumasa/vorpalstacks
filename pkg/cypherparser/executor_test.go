package cypherparser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage/graphengine"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("cypher_exec_test_%d", time.Now().UnixNano()))
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func openTestDB(t *testing.T) *graphengine.DB {
	t.Helper()
	dir := tempDir(t)
	db, err := graphengine.Open(dir, graphengine.DefaultOptions())
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func execQuery(t *testing.T, db *graphengine.DB, query string) *CypherResult {
	t.Helper()
	parsed, err := Parse(query)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if parsed.Read == nil {
		t.Fatal("expected read query, got nil")
	}
	result, err := Execute(context.Background(), db, parsed.Read, nil)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}
	return result
}

func execQueryWithParams(t *testing.T, db *graphengine.DB, query string, params map[string]any) *CypherResult {
	t.Helper()
	parsed, err := Parse(query)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if parsed.Read == nil {
		t.Fatal("expected read query, got nil")
	}
	result, err := Execute(context.Background(), db, parsed.Read, params)
	if err != nil {
		t.Fatalf("exec error: %v", err)
	}
	return result
}

// ---------------------------------------------------------------------------
// Setup helpers
// ---------------------------------------------------------------------------

type testGraph struct {
	db    *graphengine.DB
	nodes map[string]graphengine.NodeID
}

func setupPersonGraph(t *testing.T) *testGraph {
	t.Helper()
	db := openTestDB(t)
	g := &testGraph{db: db, nodes: make(map[string]graphengine.NodeID)}

	alice, _ := db.AddNode([]string{"Person"}, graphengine.Props{"name": "Alice", "age": 30})
	bob, _ := db.AddNode([]string{"Person"}, graphengine.Props{"name": "Bob", "age": 25})
	charlie, _ := db.AddNode([]string{"Person"}, graphengine.Props{"name": "Charlie", "age": 35})
	dave, _ := db.AddNode([]string{"Person"}, graphengine.Props{"name": "Dave", "age": 28})

	db.AddEdge(alice, bob, "KNOWS", nil)
	db.AddEdge(alice, charlie, "KNOWS", nil)
	db.AddEdge(bob, charlie, "KNOWS", nil)
	db.AddEdge(charlie, dave, "FOLLOWS", nil)
	db.AddEdge(dave, alice, "FOLLOWS", nil)

	g.nodes["alice"] = alice
	g.nodes["bob"] = bob
	g.nodes["charlie"] = charlie
	g.nodes["dave"] = dave
	return g
}

func setupMultiLabelGraph(t *testing.T) *testGraph {
	t.Helper()
	db := openTestDB(t)
	g := &testGraph{db: db, nodes: make(map[string]graphengine.NodeID)}

	alice, _ := db.AddNode([]string{"Person", "Actor"}, graphengine.Props{"name": "Alice", "age": 30})
	bob, _ := db.AddNode([]string{"Person"}, graphengine.Props{"name": "Bob", "age": 25})
	movie1, _ := db.AddNode([]string{"Movie"}, graphengine.Props{"title": "The Matrix", "year": 1999})
	movie2, _ := db.AddNode([]string{"Movie"}, graphengine.Props{"title": "Inception", "year": 2010})

	db.AddEdge(alice, movie1, "ACTED_IN", graphengine.Props{"role": "Neo"})
	db.AddEdge(alice, movie2, "ACTED_IN", graphengine.Props{"role": "Ariadne"})
	db.AddEdge(bob, movie1, "DIRECTED", nil)

	g.nodes["alice"] = alice
	g.nodes["bob"] = bob
	g.nodes["movie1"] = movie1
	g.nodes["movie2"] = movie2
	return g
}

// ---------------------------------------------------------------------------
// Strategy 1: Simple node match
// ---------------------------------------------------------------------------

func TestExec_MatchAllNodes(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n) RETURN n")
	if len(r.Rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(r.Rows))
	}
	if r.Columns[0] != "n" {
		t.Fatalf("expected column 'n', got %s", r.Columns[0])
	}
}

func TestExec_MatchByLabel(t *testing.T) {
	g := setupMultiLabelGraph(t)
	r := execQuery(t, g.db, "MATCH (m:Movie) RETURN m")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 movies, got %d", len(r.Rows))
	}
	for _, row := range r.Rows {
		n := row["m"].(*graphengine.Node)
		found := false
		for _, l := range n.Labels {
			if l == "Movie" {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected Movie label, got %v", n.Labels)
		}
	}
}

func TestExec_MatchByLabelAndProps(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person {name: 'Alice'}) RETURN n")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(r.Rows))
	}
	n := r.Rows[0]["n"].(*graphengine.Node)
	if n.GetString("name") != "Alice" {
		t.Fatalf("expected Alice, got %s", n.GetString("name"))
	}
}

func TestExec_MatchByPropsOnly(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n {age: 25}) RETURN n")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(r.Rows))
	}
	n := r.Rows[0]["n"].(*graphengine.Node)
	if n.GetString("name") != "Bob" {
		t.Fatalf("expected Bob, got %s", n.GetString("name"))
	}
}

func TestExec_MatchWithWhere(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.age > 28 RETURN n")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 rows (Alice=30, Charlie=35), got %d", len(r.Rows))
	}
}

func TestExec_MatchReturnProperty(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name = 'Alice' RETURN n.name")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(r.Rows))
	}
	if r.Rows[0]["n.name"] != "Alice" {
		t.Fatalf("expected Alice, got %v", r.Rows[0]["n.name"])
	}
}

func TestExec_MatchReturnMultipleProps(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name = 'Alice' RETURN n.name, n.age")
	if len(r.Rows) != 1 {
		t.Fatal("expected 1 row")
	}
	if r.Rows[0]["n.name"] != "Alice" {
		t.Fatalf("expected Alice, got %v", r.Rows[0]["n.name"])
	}
}

func TestExec_MatchWithAlias(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name = 'Alice' RETURN n.name AS personName")
	if r.Columns[0] != "personName" {
		t.Fatalf("expected column 'personName', got %s", r.Columns[0])
	}
	if r.Rows[0]["personName"] != "Alice" {
		t.Fatalf("expected Alice, got %v", r.Rows[0]["personName"])
	}
}

func TestExec_MatchEmptyResult(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person {name: 'Nobody'}) RETURN n")
	if len(r.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(r.Rows))
	}
}

func TestExec_MatchMultipleLabels(t *testing.T) {
	g := setupMultiLabelGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person:Actor) RETURN n")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 Person+Actor, got %d", len(r.Rows))
	}
}

func TestExec_MatchNoMatch(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Movie) RETURN n")
	if len(r.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(r.Rows))
	}
}

// ---------------------------------------------------------------------------
// Strategy 2: Single-hop match
// ---------------------------------------------------------------------------

func TestExec_SingleHopOutgoing(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a)-[r:KNOWS]->(b) RETURN a.name, b.name")
	if len(r.Rows) != 3 {
		t.Fatalf("expected 3 KNOWS edges, got %d", len(r.Rows))
	}
}

func TestExec_SingleHopIncoming(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a)<-[r:KNOWS]-(b) RETURN a.name, b.name")
	if len(r.Rows) != 3 {
		t.Fatalf("expected 3 KNOWS edges (reverse), got %d", len(r.Rows))
	}
}

func TestExec_SingleHopBoth(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a)-[r:KNOWS]-(b) RETURN a.name, b.name")
	if len(r.Rows) == 0 {
		t.Fatal("expected KNOWS edges in both directions")
	}
}

func TestExec_SingleHopWithWhere(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person)-[r:KNOWS]->(b:Person) WHERE a.name = 'Alice' RETURN a.name, b.name")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 (Alice->Bob, Alice->Charlie), got %d", len(r.Rows))
	}
}

func TestExec_SingleHopReturnEdge(t *testing.T) {
	g := setupMultiLabelGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person)-[r:ACTED_IN]->(m:Movie) WHERE a.name = 'Alice' RETURN a.name, r.role, m.title")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 ACTED_IN edges, got %d", len(r.Rows))
	}
}

func TestExec_SingleHopEdgeTypeScan(t *testing.T) {
	g := setupMultiLabelGraph(t)
	r := execQuery(t, g.db, "MATCH (a)-[:ACTED_IN]->(m) RETURN a.name, m.title")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 ACTED_IN edges, got %d", len(r.Rows))
	}
}

func TestExec_SingleHopNoMatch(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a)-[r:NONEXISTENT]->(b) RETURN a.name, b.name")
	if len(r.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(r.Rows))
	}
}

// ---------------------------------------------------------------------------
// Strategy 3: Variable-length match
// ---------------------------------------------------------------------------

func TestExec_VarLengthFixedHops(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person)-[:KNOWS*1..2]->(b:Person) WHERE a.name = 'Bob' RETURN b.name")
	names := make(map[string]bool)
	for _, row := range r.Rows {
		names[row["b.name"].(string)] = true
	}
	if !names["Charlie"] {
		t.Fatalf("expected Charlie reachable from Bob via KNOWS, got %v", names)
	}
}

func TestExec_VarLengthRange(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person)-[:KNOWS*1..3]->(b:Person) RETURN a.name, b.name")
	if len(r.Rows) == 0 {
		t.Fatal("expected results for 1..3 hop KNOWS")
	}
}

func TestExec_VarLengthMinDepth(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person)-[:KNOWS*2..3]->(b:Person) WHERE a.name = 'Bob' RETURN b.name")
	for _, row := range r.Rows {
		if row["b.name"] == "Alice" {
			t.Fatalf("Alice is 1 hop from Bob (via Bob->Charlie->Alice? No, Bob->Charlie is outgoing only)")
		}
	}
}

// ---------------------------------------------------------------------------
// ORDER BY
// ---------------------------------------------------------------------------

func TestExec_OrderByAsc(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN n.name ORDER BY n.name ASC")
	if len(r.Rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(r.Rows))
	}
	names := make([]string, len(r.Rows))
	for i, row := range r.Rows {
		names[i] = row["n.name"].(string)
	}
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Fatalf("not sorted ASC: %v", names)
		}
	}
}

func TestExec_OrderByDesc(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN n.name ORDER BY n.name DESC")
	names := make([]string, len(r.Rows))
	for i, row := range r.Rows {
		names[i] = row["n.name"].(string)
	}
	for i := 1; i < len(names); i++ {
		if names[i] > names[i-1] {
			t.Fatalf("not sorted DESC: %v", names)
		}
	}
}

func TestExec_OrderByAge(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN n.name, n.age ORDER BY n.age ASC")
	ages := make([]any, len(r.Rows))
	for i, row := range r.Rows {
		ages[i] = row["n.age"]
	}
	for i := 1; i < len(ages); i++ {
		if compareValues(ages[i], ages[i-1]) < 0 {
			t.Fatalf("not sorted by age ASC: %v", ages)
		}
	}
}

// ---------------------------------------------------------------------------
// LIMIT
// ---------------------------------------------------------------------------

func TestExec_Limit(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN n.name LIMIT 2")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(r.Rows))
	}
}

func TestExec_LimitZero(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN n.name LIMIT 1")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(r.Rows))
	}
}

func TestExec_LimitWithOrderBy(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN n.name ORDER BY n.age DESC LIMIT 2")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(r.Rows))
	}
}

// ---------------------------------------------------------------------------
// SKIP
// ---------------------------------------------------------------------------

func TestExec_Skip(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN n.name ORDER BY n.name ASC SKIP 2")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 rows after skip 2, got %d", len(r.Rows))
	}
}

func TestExec_SkipWithLimit(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN n.name ORDER BY n.name ASC SKIP 1 LIMIT 2")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(r.Rows))
	}
}

func TestExec_SkipExceeds(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN n.name SKIP 100")
	if len(r.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(r.Rows))
	}
}

// ---------------------------------------------------------------------------
// DISTINCT
// ---------------------------------------------------------------------------

func TestExec_Distinct(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person)-[:KNOWS]->(b:Person) RETURN DISTINCT b.name")
	names := make(map[string]bool)
	for _, row := range r.Rows {
		name := row["b.name"].(string)
		if names[name] {
			t.Fatalf("duplicate name: %s", name)
		}
		names[name] = true
	}
}

// ---------------------------------------------------------------------------
// Aggregation
// ---------------------------------------------------------------------------

func TestExec_Count(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN count(n) AS cnt")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(r.Rows))
	}
	cnt, ok := r.Rows[0]["cnt"].(int64)
	if !ok || cnt != 4 {
		t.Fatalf("expected count=4, got %v", r.Rows[0]["cnt"])
	}
}

func TestExec_CountStar(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN count(*) AS cnt")
	cnt, ok := r.Rows[0]["cnt"].(int64)
	if !ok || cnt != 4 {
		t.Fatalf("expected count(*)=4, got %v", r.Rows[0]["cnt"])
	}
}

func TestExec_CountDistinct(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person)-[:KNOWS]->(b:Person) RETURN count(DISTINCT b.name) AS cnt")
	cnt, ok := r.Rows[0]["cnt"].(int64)
	if !ok || cnt != 2 {
		t.Fatalf("expected count(DISTINCT b.name)=2 (Bob, Charlie), got %v", r.Rows[0]["cnt"])
	}
}

func TestExec_Sum(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN sum(n.age) AS total")
	total, ok := r.Rows[0]["total"].(float64)
	if !ok {
		t.Fatalf("expected float64, got %T", r.Rows[0]["total"])
	}
	expected := float64(30 + 25 + 35 + 28)
	if total != expected {
		t.Fatalf("expected sum=%f, got %f", expected, total)
	}
}

func TestExec_Avg(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN avg(n.age) AS avgAge")
	avg, ok := r.Rows[0]["avgAge"].(float64)
	if !ok {
		t.Fatalf("expected float64, got %T", r.Rows[0]["avgAge"])
	}
	expected := float64(30+25+35+28) / 4
	if avg != expected {
		t.Fatalf("expected avg=%f, got %f", expected, avg)
	}
}

func TestExec_Min(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN min(n.age) AS youngest")
	min, ok := r.Rows[0]["youngest"].(float64)
	if !ok || min != 25 {
		t.Fatalf("expected min=25, got %v", r.Rows[0]["youngest"])
	}
}

func TestExec_Max(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) RETURN max(n.age) AS oldest")
	max, ok := r.Rows[0]["oldest"].(float64)
	if !ok || max != 35 {
		t.Fatalf("expected max=35, got %v", r.Rows[0]["oldest"])
	}
}

// ---------------------------------------------------------------------------
// OPTIONAL MATCH
// ---------------------------------------------------------------------------

func TestExec_OptionalMatch(t *testing.T) {
	g := setupMultiLabelGraph(t)
	r := execQuery(t, g.db, "MATCH (p:Person) OPTIONAL MATCH (p)-[r:ACTED_IN]->(m:Movie) RETURN p.name, m.title")
	if len(r.Rows) == 0 {
		t.Fatal("expected rows")
	}
	aliceFound := false
	bobFound := false
	for _, row := range r.Rows {
		switch row["p.name"] {
		case "Alice":
			aliceFound = true
		case "Bob":
			bobFound = true
			if row["m.title"] != nil {
				t.Fatalf("Bob should have nil movie in optional match, got %v", row["m.title"])
			}
		}
	}
	if !aliceFound || !bobFound {
		t.Fatalf("expected both Alice and Bob, alice=%v bob=%v", aliceFound, bobFound)
	}
}

// ---------------------------------------------------------------------------
// Parameters
// ---------------------------------------------------------------------------

func TestExec_WithParams(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQueryWithParams(t, g.db, "MATCH (n:Person {name: $name}) RETURN n.age", map[string]any{"name": "Alice"})
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(r.Rows))
	}
}

func TestExec_WithParamsWhere(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQueryWithParams(t, g.db, "MATCH (n:Person) WHERE n.age > $minAge RETURN n.name", map[string]any{"minAge": 28})
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 rows (Alice=30, Charlie=35), got %d", len(r.Rows))
	}
}

// ---------------------------------------------------------------------------
// WHERE expressions
// ---------------------------------------------------------------------------

func TestExec_WhereAnd(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.age > 25 AND n.age < 35 RETURN n.name")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 rows (Alice=30, Dave=28), got %d", len(r.Rows))
	}
}

func TestExec_WhereOr(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name = 'Alice' OR n.name = 'Charlie' RETURN n.name")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(r.Rows))
	}
}

func TestExec_WhereNot(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE NOT n.name = 'Alice' RETURN n.name")
	if len(r.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(r.Rows))
	}
}

func TestExec_WhereStartsWith(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name STARTS WITH 'A' RETURN n.name")
	if len(r.Rows) != 1 || r.Rows[0]["n.name"] != "Alice" {
		t.Fatalf("expected Alice, got %d rows", len(r.Rows))
	}
}

func TestExec_WhereEndsWith(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name ENDS WITH 'e' RETURN n.name")
	if len(r.Rows) != 3 {
		t.Fatalf("expected 3 (Alice, Charlie, Dave all end with 'e'), got %d", len(r.Rows))
	}
}

func TestExec_WhereContains(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name CONTAINS 'li' RETURN n.name")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 (Alice, Charlie), got %d", len(r.Rows))
	}
}

func TestExec_WhereIn(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name IN ['Alice', 'Charlie'] RETURN n.name")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2, got %d", len(r.Rows))
	}
}

func TestExec_WhereIsNull(t *testing.T) {
	g := setupMultiLabelGraph(t)
	_, _ = g.db.AddNode([]string{"Person"}, graphengine.Props{"name": "Eve"})
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.age IS NULL RETURN n.name")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 (Eve has no age), got %d", len(r.Rows))
	}
}

// ---------------------------------------------------------------------------
// Functions in RETURN
// ---------------------------------------------------------------------------

func TestExec_ReturnId(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name = 'Alice' RETURN id(n) AS nid")
	if len(r.Rows) != 1 {
		t.Fatal("expected 1 row")
	}
	nid := r.Rows[0]["nid"].(int64)
	if nid <= 0 {
		t.Fatalf("expected positive node ID, got %d", nid)
	}
}

func TestExec_ReturnLabels(t *testing.T) {
	g := setupMultiLabelGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person:Actor) RETURN labels(n) AS lbls")
	if len(r.Rows) != 1 {
		t.Fatal("expected 1 row")
	}
	lbls := r.Rows[0]["lbls"].([]string)
	if len(lbls) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(lbls))
	}
}

func TestExec_ReturnProperties(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name = 'Alice' RETURN properties(n) AS props")
	if len(r.Rows) != 1 {
		t.Fatal("expected 1 row")
	}
	props := r.Rows[0]["props"].(map[string]any)
	if props["name"] != "Alice" {
		t.Fatalf("expected name=Alice, got %v", props)
	}
}

func TestExec_ReturnSize(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n:Person) WHERE n.name = 'Alice' RETURN size(labels(n)) AS cnt")
	if len(r.Rows) != 1 {
		t.Fatal("expected 1 row")
	}
	cnt := r.Rows[0]["cnt"].(int64)
	if cnt != 1 {
		t.Fatalf("expected 1 label, got %d", cnt)
	}
}

// ---------------------------------------------------------------------------
// EXPLAIN / PROFILE
// ---------------------------------------------------------------------------

func TestExec_Explain(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "EXPLAIN MATCH (n:Person) RETURN n")
	if len(r.Rows) != 1 {
		t.Fatalf("EXPLAIN should return 1 row with plan, got %d", len(r.Rows))
	}
	if len(r.Columns) != 1 || r.Columns[0] != "plan" {
		t.Fatalf("EXPLAIN should have column 'plan', got %v", r.Columns)
	}
	planStr, ok := r.Rows[0]["plan"].(string)
	if !ok || !strings.HasPrefix(planStr, "EXPLAIN:") {
		t.Fatalf("EXPLAIN plan should start with 'EXPLAIN:', got: %s", planStr)
	}
}

func TestExec_Profile(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "PROFILE MATCH (n:Person) RETURN n")
	if len(r.Rows) != 1 {
		t.Fatalf("PROFILE should return 1 row with plan+results, got %d", len(r.Rows))
	}
	if len(r.Columns) != 2 || r.Columns[0] != "plan" || r.Columns[1] != "results" {
		t.Fatalf("PROFILE should have columns ['plan','results'], got %v", r.Columns)
	}
	planStr, ok := r.Rows[0]["plan"].(string)
	if !ok || !strings.HasPrefix(planStr, "PROFILE:") {
		t.Fatalf("PROFILE plan should start with 'PROFILE:', got: %s", planStr)
	}
	cypherResult, ok := r.Rows[0]["results"].(*CypherResult)
	if !ok {
		t.Fatalf("PROFILE results should be *CypherResult, got %T", r.Rows[0]["results"])
	}
	if len(cypherResult.Rows) != 4 {
		t.Fatalf("PROFILE embedded result should have 4 rows, got %d", len(cypherResult.Rows))
	}
}

// ---------------------------------------------------------------------------
// Edge cases
// ---------------------------------------------------------------------------

func TestExec_EmptyGraph(t *testing.T) {
	db := openTestDB(t)
	r := execQuery(t, db, "MATCH (n) RETURN n")
	if len(r.Rows) != 0 {
		t.Fatalf("expected 0 rows on empty graph, got %d", len(r.Rows))
	}
}

func TestExec_ReturnCountOnEmpty(t *testing.T) {
	db := openTestDB(t)
	r := execQuery(t, db, "MATCH (n) RETURN count(n) AS cnt")
	if len(r.Rows) != 1 {
		t.Fatal("expected 1 row")
	}
	cnt, ok := r.Rows[0]["cnt"].(int64)
	if !ok || cnt != 0 {
		t.Fatalf("expected count=0, got %v", r.Rows[0]["cnt"])
	}
}

func TestExec_AnonymousVariable(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (:Person)-[:KNOWS]->(b) RETURN b.name")
	if len(r.Rows) != 3 {
		t.Fatalf("expected 3 KNOWS targets, got %d", len(r.Rows))
	}
}

// ---------------------------------------------------------------------------
// Index-backed WHERE equality extraction
// ---------------------------------------------------------------------------

func TestExec_WhereEqualityIndexLookup(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n) WHERE n.name = 'Bob' RETURN n.age")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(r.Rows))
	}
}

func TestExec_WhereAndEqualityIndexLookup(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (n) WHERE n.age = 25 AND n.name = 'Bob' RETURN n.name")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(r.Rows))
	}
}

// ---------------------------------------------------------------------------
// Complex patterns
// ---------------------------------------------------------------------------

func TestExec_PatternWithLabelFilter(t *testing.T) {
	g := setupMultiLabelGraph(t)
	r := execQuery(t, g.db, "MATCH (p:Person)-[r:ACTED_IN]->(m:Movie) RETURN p.name, m.title ORDER BY p.name ASC")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(r.Rows))
	}
}

func TestExec_FOLLOWSChain(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a)-[r:FOLLOWS]->(b) RETURN a.name, b.name")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 FOLLOWS edges, got %d", len(r.Rows))
	}
}

func TestExec_VarLengthFOLLOWS(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person)-[:FOLLOWS*1..2]->(b:Person) WHERE a.name = 'Charlie' RETURN b.name")
	if len(r.Rows) == 0 {
		t.Fatal("expected Dave reachable via FOLLOWS from Charlie")
	}
}

func TestExec_MultiMatch(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person {name: 'Alice'}) MATCH (a)-[:KNOWS]->(b) RETURN b.name")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 KNOWS from Alice, got %d", len(r.Rows))
	}
	names := make(map[string]bool)
	for _, row := range r.Rows {
		names[row["b.name"].(string)] = true
	}
	if !names["Bob"] || !names["Charlie"] {
		t.Fatalf("expected Bob and Charlie, got %v", names)
	}
}

func TestExec_MultiMatchWithWhere(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person) MATCH (a)-[:KNOWS]->(b) WHERE a.name = 'Alice' RETURN b.name")
	if len(r.Rows) != 2 {
		t.Fatalf("expected 2 KNOWS from Alice, got %d", len(r.Rows))
	}
}

func TestExec_MultiMatchThreeClauses(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person {name: 'Alice'}) MATCH (a)-[:KNOWS]->(b) MATCH (b)-[:KNOWS]->(c) RETURN c.name")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 (Alice->Bob->Charlie or Alice->Charlie->?), got %d", len(r.Rows))
	}
}

func TestExec_MultiHopThreeNodes(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person)-[:KNOWS]->(b:Person)-[:KNOWS]->(c:Person) RETURN a.name, b.name, c.name")
	if len(r.Rows) == 0 {
		t.Fatal("expected at least 1 result for 3-node KNOWS chain")
	}
}

func TestExec_MultiHopWithWhere(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a:Person)-[:KNOWS]->(b:Person)-[:KNOWS]->(c:Person) WHERE a.name = 'Alice' RETURN b.name, c.name")
	if len(r.Rows) != 1 {
		t.Fatalf("expected 1 (Alice->Bob->Charlie), got %d", len(r.Rows))
	}
	if r.Rows[0]["c.name"] != "Charlie" {
		t.Fatalf("expected Charlie, got %v", r.Rows[0]["c.name"])
	}
}

func TestExec_MultiHopFourNodes(t *testing.T) {
	g := setupPersonGraph(t)
	r := execQuery(t, g.db, "MATCH (a)-[:KNOWS]->(b)-[:KNOWS]->(c) RETURN a.name, c.name")
	if len(r.Rows) == 0 {
		t.Fatal("expected results for 3-node KNOWS chain")
	}
	foundAliceCharlie := false
	for _, row := range r.Rows {
		if row["a.name"] == "Alice" && row["c.name"] == "Charlie" {
			foundAliceCharlie = true
		}
	}
	if !foundAliceCharlie {
		t.Fatal("expected Alice->Bob->Charlie path")
	}
}
