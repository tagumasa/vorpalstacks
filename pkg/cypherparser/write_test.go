package cypherparser

import (
	"context"
	"fmt"
	"testing"

	"vorpalstacks/pkg/graphengine"
)

func setupTestDB(t *testing.T) (graphengine.GraphReader, graphengine.GraphWriter, func()) {
	t.Helper()
	dir := t.TempDir()
	db, err := graphengine.Open(dir, graphengine.Options{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return db, db, func() { db.Close() }
}

func assertIntEqual(t *testing.T, name string, expected int64, got any) {
	t.Helper()
	f, ok := toFloat64(got)
	if !ok || int64(f) != expected {
		t.Fatalf("%s: expected %d, got %v (%T)", name, expected, got, got)
	}
}

func assertStrEqual(t *testing.T, name string, expected string, got any) {
	t.Helper()
	if got != expected {
		t.Fatalf("%s: expected %q, got %v (%T)", name, expected, got, got)
	}
}

// ---------------------------------------------------------------------------
// CREATE
// ---------------------------------------------------------------------------

func TestWrite_CreateSingleNode(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	parsed, err := Parse(`CREATE (n:Person {name: "Alice", age: 30}) RETURN n.name`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteWrite(context.Background(), r, w, parsed.Write, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
	assertStrEqual(t, "n.name", "Alice", result.Rows[0]["n.name"])

	if r.Stats().NodeCount != 1 {
		t.Fatalf("expected 1 node, got %d", r.Stats().NodeCount)
	}
}

func TestWrite_CreateEdge(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	parsed, err := Parse(`CREATE (a:Person {name: "Alice"})-[r:KNOWS {since: 2020}]->(b:Person {name: "Bob"}) RETURN a.name, b.name, r.since`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteWrite(context.Background(), r, w, parsed.Write, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
	assertStrEqual(t, "a.name", "Alice", result.Rows[0]["a.name"])
	assertStrEqual(t, "b.name", "Bob", result.Rows[0]["b.name"])
	assertIntEqual(t, "r.since", 2020, result.Rows[0]["r.since"])

	if r.Stats().NodeCount != 2 {
		t.Fatalf("expected 2 nodes, got %d", r.Stats().NodeCount)
	}
	if r.Stats().EdgeCount != 1 {
		t.Fatalf("expected 1 edge, got %d", r.Stats().EdgeCount)
	}
}

func TestWrite_CreateMultiplePatterns(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	parsed, err := Parse(`CREATE (a:Person {name: "Alice"}), (b:Person {name: "Bob"}), (c:Person {name: "Carol"}) RETURN a.name, b.name, c.name`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteWrite(context.Background(), r, w, parsed.Write, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
	assertStrEqual(t, "a.name", "Alice", result.Rows[0]["a.name"])
	assertStrEqual(t, "b.name", "Bob", result.Rows[0]["b.name"])
	assertStrEqual(t, "c.name", "Carol", result.Rows[0]["c.name"])

	if r.Stats().NodeCount != 3 {
		t.Fatalf("expected 3 nodes, got %d", r.Stats().NodeCount)
	}
}

func TestWrite_CreateIncomingEdge(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	parsed, err := Parse(`CREATE (a:Person {name: "Alice"})<-[:KNOWS]-(b:Person {name: "Bob"}) RETURN a.name, b.name`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	_, err = ExecuteWrite(context.Background(), r, w, parsed.Write, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if r.Stats().NodeCount != 2 {
		t.Fatalf("expected 2 nodes, got %d", r.Stats().NodeCount)
	}
	if r.Stats().EdgeCount != 1 {
		t.Fatalf("expected 1 edge, got %d", r.Stats().EdgeCount)
	}
}

func TestWrite_CreateNoReturn(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	parsed, err := Parse(`CREATE (n:Person {name: "Dave"})`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteWrite(context.Background(), r, w, parsed.Write, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(result.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(result.Rows))
	}

	if r.Stats().NodeCount != 1 {
		t.Fatalf("expected 1 node, got %d", r.Stats().NodeCount)
	}
}

func TestWrite_CreateExpressionProps(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	parsed, err := Parse(`CREATE (n:Calc {val: 10 + 5}) RETURN n.val`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteWrite(context.Background(), r, w, parsed.Write, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	assertIntEqual(t, "n.val", 15, result.Rows[0]["n.val"])
}

// ---------------------------------------------------------------------------
// MERGE
// ---------------------------------------------------------------------------

func TestWrite_MergeCreate(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	parsed, err := Parse(`MERGE (n:Person {name: "Alice"}) RETURN n.name`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteMerge(context.Background(), r, w, parsed.Merge, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	assertStrEqual(t, "n.name", "Alice", result.Rows[0]["n.name"])

	if r.Stats().NodeCount != 1 {
		t.Fatalf("expected 1 node, got %d", r.Stats().NodeCount)
	}

	_, err = ExecuteMerge(context.Background(), r, w, parsed.Merge, nil)
	if err != nil {
		t.Fatalf("execute 2: %v", err)
	}

	if r.Stats().NodeCount != 1 {
		t.Fatalf("MERGE should not create duplicate, got %d", r.Stats().NodeCount)
	}
}

func TestWrite_MergeOnCreateSet(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	parsed, err := Parse(`MERGE (n:Person {name: "Alice"}) ON CREATE SET n.age = 30 RETURN n.name, n.age`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteMerge(context.Background(), r, w, parsed.Merge, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	assertIntEqual(t, "n.age", 30, result.Rows[0]["n.age"])
}

func TestWrite_MergeOnMatchSet(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	seedParsed, err := Parse(`MERGE (n:Person {name: "Alice"})`)
	if err != nil {
		t.Fatalf("parse seed: %v", err)
	}
	_, err = ExecuteMerge(context.Background(), r, w, seedParsed.Merge, nil)
	if err != nil {
		t.Fatalf("execute seed: %v", err)
	}

	parsed, err := Parse(`MERGE (n:Person {name: "Alice"}) ON MATCH SET n.age = 40 RETURN n.name, n.age`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteMerge(context.Background(), r, w, parsed.Merge, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	assertIntEqual(t, "n.age", 40, result.Rows[0]["n.age"])

	if r.Stats().NodeCount != 1 {
		t.Fatalf("expected still 1 node, got %d", r.Stats().NodeCount)
	}
}

// ---------------------------------------------------------------------------
// SET (on MATCH results)
// ---------------------------------------------------------------------------

func TestWrite_MatchSet(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	seedParsed, err := Parse(`CREATE (n:Person {name: "Alice"})`)
	if err != nil {
		t.Fatalf("parse seed: %v", err)
	}
	_, err = ExecuteWrite(context.Background(), r, w, seedParsed.Write, nil)
	if err != nil {
		t.Fatalf("execute seed: %v", err)
	}

	parsed, err := Parse(`MATCH (n:Person {name: "Alice"}) SET n.age = 30 RETURN n.name, n.age`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteQueryWrite(context.Background(), r, w, parsed.Read, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
	assertIntEqual(t, "n.age", 30, result.Rows[0]["n.age"])
}

func TestWrite_MatchSetMultipleRows(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	seedParsed, err := Parse(`CREATE (a:Person {name: "Alice"}), (b:Person {name: "Bob"})`)
	if err != nil {
		t.Fatalf("parse seed: %v", err)
	}
	_, err = ExecuteWrite(context.Background(), r, w, seedParsed.Write, nil)
	if err != nil {
		t.Fatalf("execute seed: %v", err)
	}

	parsed, err := Parse(`MATCH (n:Person) SET n.active = true RETURN n.name, n.active`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteQueryWrite(context.Background(), r, w, parsed.Read, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result.Rows))
	}
	for i, row := range result.Rows {
		if row["n.active"] != true {
			t.Fatalf("row %d: expected active=true, got %v", i, row["n.active"])
		}
	}
}

// ---------------------------------------------------------------------------
// DELETE
// ---------------------------------------------------------------------------

func TestWrite_MatchDelete(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	seedParsed, err := Parse(`CREATE (n:Person {name: "ToDelete"})`)
	if err != nil {
		t.Fatalf("parse seed: %v", err)
	}
	_, err = ExecuteWrite(context.Background(), r, w, seedParsed.Write, nil)
	if err != nil {
		t.Fatalf("execute seed: %v", err)
	}

	parsed, err := Parse(`MATCH (n:Person {name: "ToDelete"}) DELETE n`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	_, err = ExecuteQueryWrite(context.Background(), r, w, parsed.Read, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if r.Stats().NodeCount != 0 {
		t.Fatalf("expected 0 nodes after delete, got %d", r.Stats().NodeCount)
	}
}

func TestWrite_MatchDeleteWithEdges(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	seedParsed, err := Parse(`CREATE (a:Person {name: "Alice"})-[:KNOWS]->(b:Person {name: "Bob"})`)
	if err != nil {
		t.Fatalf("parse seed: %v", err)
	}
	_, err = ExecuteWrite(context.Background(), r, w, seedParsed.Write, nil)
	if err != nil {
		t.Fatalf("execute seed: %v", err)
	}

	parsed, err := Parse(`MATCH (n:Person {name: "Alice"}) DELETE n`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	_, err = ExecuteQueryWrite(context.Background(), r, w, parsed.Read, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if r.Stats().EdgeCount != 0 {
		t.Fatalf("expected 0 edges, got %d", r.Stats().EdgeCount)
	}
	if r.Stats().NodeCount != 1 {
		t.Fatalf("expected 1 node (Bob) remaining, got %d", r.Stats().NodeCount)
	}
}

func TestWrite_DeleteEdge(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	seedParsed, err := Parse(`CREATE (a:Person {name: "Alice"})-[r:KNOWS]->(b:Person {name: "Bob"})`)
	if err != nil {
		t.Fatalf("parse seed: %v", err)
	}
	_, err = ExecuteWrite(context.Background(), r, w, seedParsed.Write, nil)
	if err != nil {
		t.Fatalf("execute seed: %v", err)
	}

	parsed, err := Parse(`MATCH (a:Person {name: "Alice"})-[r:KNOWS]->(b:Person {name: "Bob"}) DELETE r`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	_, err = ExecuteQueryWrite(context.Background(), r, w, parsed.Read, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if r.Stats().EdgeCount != 0 {
		t.Fatalf("expected 0 edges, got %d", r.Stats().EdgeCount)
	}
	if r.Stats().NodeCount != 2 {
		t.Fatalf("expected 2 nodes, got %d", r.Stats().NodeCount)
	}
}

// ---------------------------------------------------------------------------
// REMOVE
// ---------------------------------------------------------------------------

func TestWrite_RemoveLabel(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	seedParsed, err := Parse(`CREATE (n:Person:Employee {name: "Alice"})`)
	if err != nil {
		t.Fatalf("parse seed: %v", err)
	}
	_, err = ExecuteWrite(context.Background(), r, w, seedParsed.Write, nil)
	if err != nil {
		t.Fatalf("execute seed: %v", err)
	}

	parsed, err := Parse(`MATCH (n:Person {name: "Alice"}) REMOVE n:Employee RETURN labels(n)`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteQueryWrite(context.Background(), r, w, parsed.Read, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
	labels, ok := result.Rows[0]["labels(n)"].([]string)
	if !ok {
		t.Fatalf("expected []string, got %T", result.Rows[0]["labels(n)"])
	}
	for _, l := range labels {
		if l == "Employee" {
			t.Fatalf("Employee label should have been removed, got %v", labels)
		}
	}
	hasPerson := false
	for _, l := range labels {
		if l == "Person" {
			hasPerson = true
		}
	}
	if !hasPerson {
		t.Fatalf("Person label should remain, got %v", labels)
	}
}

func TestWrite_RemoveProperty(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	seedParsed, err := Parse(`CREATE (n:Person {name: "Alice", temp: "yes"})`)
	if err != nil {
		t.Fatalf("parse seed: %v", err)
	}
	_, err = ExecuteWrite(context.Background(), r, w, seedParsed.Write, nil)
	if err != nil {
		t.Fatalf("execute seed: %v", err)
	}

	parsed, err := Parse(`MATCH (n:Person {name: "Alice"}) REMOVE n.temp RETURN n.name`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result, err := ExecuteQueryWrite(context.Background(), r, w, parsed.Read, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
	assertStrEqual(t, "n.name", "Alice", result.Rows[0]["n.name"])
}

// ---------------------------------------------------------------------------
// Combined: CREATE then MATCH
// ---------------------------------------------------------------------------

func TestWrite_CreateThenMatch(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	createParsed, err := Parse(`CREATE (a:Person {name: "Alice"})-[r:KNOWS {since: 2020}]->(b:Person {name: "Bob"})`)
	if err != nil {
		t.Fatalf("parse create: %v", err)
	}
	_, err = ExecuteWrite(context.Background(), r, w, createParsed.Write, nil)
	if err != nil {
		t.Fatalf("execute create: %v", err)
	}

	matchParsed, err := Parse(`MATCH (a:Person {name: "Alice"})-[r:KNOWS]->(b:Person) RETURN a.name, r.since, b.name`)
	if err != nil {
		t.Fatalf("parse match: %v", err)
	}

	result, err := Execute(context.Background(), r, matchParsed.Read, nil)
	if err != nil {
		t.Fatalf("execute match: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
	assertStrEqual(t, "a.name", "Alice", result.Rows[0]["a.name"])
	assertStrEqual(t, "b.name", "Bob", result.Rows[0]["b.name"])
	assertIntEqual(t, "r.since", 2020, result.Rows[0]["r.since"])
}

// ---------------------------------------------------------------------------
// Context cancellation
// ---------------------------------------------------------------------------

func TestWrite_ContextCancellation(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	parsed, err := Parse(`CREATE (n:Person {name: "Alice"})`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	_, err = ExecuteWrite(ctx, r, w, parsed.Write, nil)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestWrite_MergeMissingLabel(t *testing.T) {
	_, err := Parse(`MERGE (n {name: "Alice"})`)
	if err == nil {
		t.Fatal("expected error for MERGE without label")
	}
}

func TestWrite_CreateEdgeWithoutLabel(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	parsed, err := Parse(`CREATE (a:Person)-[]->(b:Person)`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	_, err = ExecuteWrite(context.Background(), r, w, parsed.Write, nil)
	if err == nil {
		t.Fatal("expected error for edge without label")
	}
}

func TestWrite_MissingParam(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	parsed, err := Parse(`CREATE (n:Person {name: $name})`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	_, err = ExecuteWrite(context.Background(), r, w, parsed.Write, nil)
	if err == nil {
		t.Fatal("expected error for missing parameter")
	}
}

// ---------------------------------------------------------------------------
// Integration: full write-read cycle
// ---------------------------------------------------------------------------

func TestWrite_FullCycle(t *testing.T) {
	r, w, close := setupTestDB(t)
	defer close()

	createQ := `CREATE (a:User {name: "Alice", score: 100})-[r:PLAYS {level: 5}]->(g:Game {name: "Chess"})`
	parsed, err := Parse(createQ)
	if err != nil {
		t.Fatalf("parse create: %v", err)
	}
	if _, err := ExecuteWrite(context.Background(), r, w, parsed.Write, nil); err != nil {
		t.Fatalf("execute create: %v", err)
	}

	setQ := `MATCH (u:User {name: "Alice"}) SET u.score = 150`
	parsed, err = Parse(setQ)
	if err != nil {
		t.Fatalf("parse set: %v", err)
	}
	if _, err := ExecuteQueryWrite(context.Background(), r, w, parsed.Read, nil); err != nil {
		t.Fatalf("execute set: %v", err)
	}

	matchQ := `MATCH (u:User {name: "Alice"})-[r:PLAYS]->(g:Game) RETURN u.name, u.score, r.level, g.name`
	parsed, err = Parse(matchQ)
	if err != nil {
		t.Fatalf("parse match: %v", err)
	}
	result, err := Execute(context.Background(), r, parsed.Read, nil)
	if err != nil {
		t.Fatalf("execute match: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
	assertStrEqual(t, "u.name", "Alice", result.Rows[0]["u.name"])
	assertIntEqual(t, "u.score", 150, result.Rows[0]["u.score"])
	assertIntEqual(t, "r.level", 5, result.Rows[0]["r.level"])
	assertStrEqual(t, "g.name", "Chess", result.Rows[0]["g.name"])

	deleteQ := `MATCH (u:User {name: "Alice"})-[r:PLAYS]->(g:Game) DELETE r`
	parsed, err = Parse(deleteQ)
	if err != nil {
		t.Fatalf("parse delete: %v", err)
	}
	if _, err := ExecuteQueryWrite(context.Background(), r, w, parsed.Read, nil); err != nil {
		t.Fatalf("execute delete: %v", err)
	}

	if r.Stats().EdgeCount != 0 {
		t.Fatalf("expected 0 edges after delete, got %d", r.Stats().EdgeCount)
	}
	if r.Stats().NodeCount != 2 {
		t.Fatalf("expected 2 nodes remaining, got %d", r.Stats().NodeCount)
	}

	_ = fmt.Sprintf
}
