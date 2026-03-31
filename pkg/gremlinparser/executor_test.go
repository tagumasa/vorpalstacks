package gremlinparser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"vorpalstacks/pkg/graphengine"
)

func openTestDB(t *testing.T) *graphengine.DB {
	t.Helper()
	dir := filepath.Join("tmp", "gremlin_test_"+t.Name())
	os.RemoveAll(dir)
	db, err := graphengine.Open(dir, graphengine.Options{})
	if err != nil {
		t.Fatalf("failed to open graph DB: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
		os.RemoveAll(dir)
	})
	return db
}

func setupSocialGraph(t *testing.T) *graphengine.DB {
	t.Helper()
	db := openTestDB(t)

	alice, _ := db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice", "age": 30})
	bob, _ := db.AddNode([]string{"person"}, graphengine.Props{"name": "Bob", "age": 25})
	charlie, _ := db.AddNode([]string{"person"}, graphengine.Props{"name": "Charlie", "age": 35})
	dave, _ := db.AddNode([]string{"person"}, graphengine.Props{"name": "Dave", "age": 28})
	eve, _ := db.AddNode([]string{"person"}, graphengine.Props{"name": "Eve", "age": 22})

	db.AddEdge(alice, bob, "knows", nil)
	db.AddEdge(alice, charlie, "knows", nil)
	db.AddEdge(bob, dave, "knows", nil)
	db.AddEdge(charlie, eve, "knows", nil)
	db.AddEdge(dave, eve, "follows", nil)

	return db
}

func execQuery(t *testing.T, db *graphengine.DB, query string) any {
	t.Helper()
	ctx := context.Background()
	parsed, err := Parse(query)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	result, err := ExecuteQuery(ctx, db, db, parsed, nil)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}
	return result
}

func execQueryWrite(t *testing.T, db *graphengine.DB, query string) any {
	t.Helper()
	return execQuery(t, db, query)
}

func toStringSlice(result any) []string {
	switch v := result.(type) {
	case []any:
		s := make([]string, len(v))
		for i, item := range v {
			s[i] = fmt.Sprintf("%v", item)
		}
		return s
	default:
		return []string{fmt.Sprintf("%v", v)}
	}
}

func toInt64Slice(result any) []int64 {
	switch v := result.(type) {
	case []any:
		s := make([]int64, len(v))
		for i, item := range v {
			s[i] = toInt64(item)
		}
		return s
	case int64:
		return []int64{v}
	case float64:
		return []int64{int64(v)}
	case int:
		return []int64{int64(v)}
	default:
		return nil
	}
}

func toInt64(val any) int64 {
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}

func toStringVal(t *testing.T, result any) string {
	t.Helper()
	switch v := result.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func containsAll(haystack []string, needles ...string) bool {
	haySet := make(map[string]bool)
	for _, h := range haystack {
		haySet[h] = true
	}
	for _, n := range needles {
		if !haySet[n] {
			return false
		}
	}
	return true
}
