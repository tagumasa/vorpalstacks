package appsync

import (
	"testing"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

// selectionDepth is the query-body nesting the documented QueryDepthLimit
// bounds, so the counting is pinned against a document table: plain fields,
// inline fragments, and fragment spreads (which are transparent and
// cycle-safe).
func TestSelectionDepthCountsQueryBodyNesting(t *testing.T) {
	schema := gqlparser.MustLoadSchema(&ast.Source{
		Name:  "schema.graphql",
		Input: "type Query { s: String a: A } type A { s: String b: B } type B { c: String }",
	})

	depthOf := func(t *testing.T, query string) int {
		t.Helper()
		doc, err := gqlparser.LoadQuery(schema, query)
		if err != nil {
			t.Fatalf("load %q: %v", query, err)
		}
		return selectionDepth(doc.Operations[0].SelectionSet, doc.Fragments)
	}

	cases := map[string]int{
		"{ s }":                        1,
		"{ a { s } }":                  2,
		"{ a { b { c } } }":            3,
		"{ ... on Query { a { s } } }": 2,
		"query { a { b { ...F } } } fragment F on B { c }": 3,
	}
	for query, want := range cases {
		if got := depthOf(t, query); got != want {
			t.Errorf("%q: selectionDepth = %d, want %d", query, got, want)
		}
	}
}
