package lambda

import (
	"strings"
	"testing"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// versionListFixture seeds a function store for the entry-pagination tests.
func versionListFixture(t *testing.T, functions ...*lambdastore.Function) *lambdaStore {
	t.Helper()
	store := &lambdaStore{
		Functions: lambdastore.NewFunctionStore(&memStorage{bucket: newMemBucket()}, "123456789012", "us-east-1"),
	}
	for _, fn := range functions {
		if _, err := store.Functions.Create(fn); err != nil {
			t.Fatalf("create function %s: %v", fn.FunctionName, err)
		}
	}
	return store
}

func versionedFunction(name string, versions ...string) *lambdastore.Function {
	fn := &lambdastore.Function{
		FunctionName: name,
		Runtime:      "nodejs22.x",
		Handler:      "index.handler",
		Role:         "arn:aws:iam::123456789012:role/lambda",
	}
	for _, v := range versions {
		fn.Versions = append(fn.Versions, lambdastore.Version{Version: v})
	}
	return fn
}

// entryLabel renders one listing entry as name@version for assertions.
func entryLabel(entries []versionListEntry, i int) string {
	if entries[i].versionIndex < 0 {
		return entries[i].fn.FunctionName + "@$LATEST"
	}
	return entries[i].fn.FunctionName + "@" + entries[i].fn.Versions[entries[i].versionIndex].Version
}

// TestListFunctionVersionEntries_CapsPagesPerEntry pins the documented
// contract: "ListFunctions returns a maximum of 50 items in each response"
// and FunctionVersion=ALL makes every published version its own entry, so
// the page cap counts entries, not functions, and the marker resumes inside
// a partially listed function.
func TestListFunctionVersionEntries_CapsPagesPerEntry(t *testing.T) {
	svc := &LambdaService{}
	store := versionListFixture(t,
		versionedFunction("alpha", "1", "2", "3"),
		versionedFunction("beta"),
	)

	type page struct {
		labels []string
		next   string
	}
	var pages []page
	marker := ""
	for {
		entries, next, err := svc.listFunctionVersionEntries(store, marker, 2)
		if err != nil {
			t.Fatalf("page with marker %q: %v", marker, err)
		}
		if len(entries) > 2 {
			t.Fatalf("page must carry at most 2 entries, got %d", len(entries))
		}
		p := page{next: next}
		for i := range entries {
			p.labels = append(p.labels, entryLabel(entries, i))
		}
		pages = append(pages, p)
		if next == "" {
			break
		}
		marker = next
		if len(pages) > 6 {
			t.Fatalf("pagination did not terminate, last marker %q", marker)
		}
	}

	want := [][]string{
		{"alpha@$LATEST", "alpha@1"},
		{"alpha@2", "alpha@3"},
		{"beta@$LATEST"},
	}
	if len(pages) != len(want) {
		t.Fatalf("expected %d pages, got %d: %+v", len(want), len(pages), pages)
	}
	for i, w := range want {
		if strings.Join(pages[i].labels, ",") != strings.Join(w, ",") {
			t.Fatalf("page %d = %v, want %v", i, pages[i].labels, w)
		}
	}
}

// TestListFunctionVersionEntries_ResumeAtFunctionStart pins the marker that
// stops before a function's own entry: "name|0" lists the function's
// configuration entry and versions from the beginning.
func TestListFunctionVersionEntries_ResumeAtFunctionStart(t *testing.T) {
	svc := &LambdaService{}
	store := versionListFixture(t,
		versionedFunction("alpha", "1", "2"),
		versionedFunction("beta"),
	)

	entries, next, err := svc.listFunctionVersionEntries(store, "alpha|0", 10)
	if err != nil {
		t.Fatalf("resume at function start: %v", err)
	}
	if len(entries) != 4 || entryLabel(entries, 0) != "alpha@$LATEST" || entryLabel(entries, 3) != "beta@$LATEST" {
		labels := make([]string, 0, len(entries))
		for i := range entries {
			labels = append(labels, entryLabel(entries, i))
		}
		t.Fatalf("unexpected entries %v", labels)
	}
	if next != "" {
		t.Fatalf("exhausted listing must not carry a marker, got %q", next)
	}
}

// TestListFunctionVersionEntries_RejectsMalformedMarkers pins the error
// path for tokens that are not name|count pairs.
func TestListFunctionVersionEntries_RejectsMalformedMarkers(t *testing.T) {
	svc := &LambdaService{}
	store := versionListFixture(t, versionedFunction("alpha", "1"))
	for _, marker := range []string{"alpha|x", "alpha|-1", "|2"} {
		if _, _, err := svc.listFunctionVersionEntries(store, marker, 10); err == nil {
			t.Fatalf("marker %q must be rejected", marker)
		}
	}
}
