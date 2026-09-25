package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// The shared write-response builder is the single response face of the
// standalone UPDATE and DELETE planes: a result without a RETURNING clause
// serves the empty Items list, each clause value picks its image, and the
// MODIFIED forms restrict to the written attributes — a target with no old
// value simply does not appear.
func TestBuildPartiqlWriteResponseServesReturningImages(t *testing.T) {
	sAttr := func(v string) *dbstore.AttributeValue { return &dbstore.AttributeValue{S: &v} }
	oldImage := map[string]*dbstore.AttributeValue{"a": sAttr("old-a"), "b": sAttr("kept")}
	newImage := map[string]*dbstore.AttributeValue{"a": sAttr("new-a"), "b": sAttr("kept"), "c": sAttr("created")}

	wireString := func(item map[string]interface{}, name string) (string, bool) {
		attr, ok := item[name].(map[string]interface{})
		if !ok {
			return "", false
		}
		v, ok := attr["S"].(string)
		return v, ok
	}
	itemsOf := func(resp map[string]interface{}) []map[string]interface{} {
		t.Helper()
		items, ok := resp["Items"].([]map[string]interface{})
		if !ok {
			t.Fatalf("response carries no Items list: %v", resp["Items"])
		}
		return items
	}

	for _, noClause := range []*partiqlWriteResult{nil, {oldImage: oldImage, newImage: newImage}} {
		if items := itemsOf(buildPartiqlWriteResponse(noClause)); len(items) != 0 {
			t.Fatalf("result without a clause must serve an empty Items list, got %v", items)
		}
	}

	cases := []struct {
		returning string
		wantKeys  []string
		wantA     string
	}{
		{"ALL OLD *", []string{"a", "b"}, "old-a"},
		{"ALL NEW *", []string{"a", "b", "c"}, "new-a"},
		{"MODIFIED OLD *", []string{"a"}, "old-a"},
		{"MODIFIED NEW *", []string{"a", "c"}, "new-a"},
	}
	for _, tc := range cases {
		result := &partiqlWriteResult{oldImage: oldImage, newImage: newImage, returning: tc.returning, modified: []string{"a", "c"}}
		items := itemsOf(buildPartiqlWriteResponse(result))
		if len(items) != 1 {
			t.Fatalf("%s: expected one item, got %d", tc.returning, len(items))
		}
		if got, ok := wireString(items[0], "a"); !ok || got != tc.wantA {
			t.Fatalf("%s: attribute a = %q (ok=%v), want %q", tc.returning, got, ok, tc.wantA)
		}
		if len(items[0]) != len(tc.wantKeys) {
			t.Fatalf("%s: image keys %v, want exactly %v", tc.returning, items[0], tc.wantKeys)
		}
	}
}

// A malformed ExecuteStatement NextToken is a rejected request, never an
// absent cursor: a token that is not base64, not the offset JSON, or a
// non-positive offset would otherwise restart the page at the first item
// and re-serve what the caller already saw. A well-formed token whose
// offset is past a shrunken result set stays the documented last page.
func TestExecuteStatementRejectsMalformedNextToken(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()

	cases := []struct {
		name  string
		token string
	}{
		{"not base64", `!!!not-base64!!!`},
		{"base64 of non-JSON", `aGVsbG8=`},
		{"JSON of a non-positive offset", `MA==`},
	}
	for _, tc := range cases {
		_, err := svc.executePartiQLSelectEnhanced(ctx, reqCtx, `SELECT * FROM LegacyTable`, parsePartiQLParams(nil), false, 0, tc.token)
		if !errors.Is(err, ErrInvalidParameter) {
			t.Fatalf("%s: expected the parameter rejection, got %v", tc.name, err)
		}
	}

	// A well-formed offset past the (empty) result set is the last page,
	// not an error — the stale-token contract the strict decode preserves.
	resp, err := svc.executePartiQLSelectEnhanced(ctx, reqCtx, `SELECT * FROM LegacyTable`, parsePartiQLParams(nil), false, 0, `MTA=`)
	if err != nil {
		t.Fatalf("well-formed offset token rejected: %v", err)
	}
	result, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("select response shape: %T", resp)
	}
	if items, ok := result["Items"].([]map[string]interface{}); !ok || len(items) != 0 {
		t.Fatalf("past-the-end offset must answer an empty page, got %v", result["Items"])
	}
}

// The ORDER BY sort is stable: equal values tie, and only a stable sort
// keeps tied elements in their input order — the deterministic-item
// contract paginated responses rely on. The pin uses equal-valued ties:
// with type-mismatched values the comparator's ties are not a valid weak
// ordering (a mismatch ties with both sides of a strict pair, cycling),
// so "stable order" there means the sort's own determinism, not a
// position contract that could be pinned.
func TestOrderBySortKeepsTiesInInputOrder(t *testing.T) {
	num := func(v string) *dbstore.AttributeValue { return &dbstore.AttributeValue{N: &v} }
	item := func(tag string, col string) *dbstore.Item {
		return &dbstore.Item{Attributes: map[string]*dbstore.AttributeValue{"tag": {S: &tag}, "col": num(col)}}
	}

	// Eighteen items exercise the sort's partition paths (past the
	// insertion-sort cutover): sixteen equal-valued 3s, with the strictly
	// ordered 1 and 5 anchors appended at the end of the input.
	var items []*dbstore.Item
	var wantTies []string
	for i := 0; i < 16; i++ {
		tag := fmt.Sprintf("tie-%02d", i)
		items = append(items, item(tag, "3"))
		wantTies = append(wantTies, tag)
	}
	items = append(items, item("one", "1"), item("five", "5"))

	sorted := sortItemsByOrderBy(items, &orderByClause{column: "col", direction: "ASC"})

	var tags []string
	for _, it := range sorted {
		tags = append(tags, *it.Attributes["tag"].S)
	}

	if len(tags) != 18 {
		t.Fatalf("sort lost items: %v", tags)
	}
	if tags[0] != "one" || tags[17] != "five" {
		t.Fatalf("strictly ordered anchors must bound the result: %v", tags)
	}
	if strings.Join(tags[1:17], ",") != strings.Join(wantTies, ",") {
		t.Fatalf("tied equal-valued items left their input order: %v", tags[1:17])
	}
}
