package iot

import (
	"testing"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// A supplied maxResults or nextToken that is malformed — non-numeric, or
// numeric but outside the documented SearchIndex range (1..100) — must be
// rejected with InvalidRequestException: silently paging at the default or
// accepting an unbounded page would return results the caller never asked
// for.
func TestApplyPaginationRejectsNonNumericValues(t *testing.T) {
	for _, params := range []map[string]interface{}{
		{"maxResults": "abc"},
		{"maxResults": true},
		{"nextToken": "not-a-token"},
	} {
		_, _, err := applyPagination(10, params)
		if err != iotstore.ErrInvalidRequest {
			t.Fatalf("malformed %v must surface as InvalidRequestException, got %v", params, err)
		}
	}
}

// The Smithy model documents SearchIndex maxResults with range(min=1) and
// "This maximum number cannot exceed 100"; both bounds are enforced.
func TestApplyPaginationEnforcesMaxResultsRange(t *testing.T) {
	outOfRange := []interface{}{
		0, -3, float64(0), float64(-1),
		101, 500, float64(101), float64(1000),
		"0", "-1", "101", "500",
	}
	for _, v := range outOfRange {
		if _, _, err := applyPagination(500, map[string]interface{}{"maxResults": v}); err != iotstore.ErrInvalidRequest {
			t.Fatalf("out-of-range maxResults %#v must be rejected with InvalidRequestException, got %v", v, err)
		}
	}

	// Both bounds are inclusive.
	start, end, err := applyPagination(500, map[string]interface{}{"maxResults": iotstore.SearchIndexMaxResultsMin})
	if err != nil || start != 0 || end != iotstore.SearchIndexMaxResultsMin {
		t.Fatalf("maxResults=min not honoured: start=%d end=%d err=%v", start, end, err)
	}
	start, end, err = applyPagination(500, map[string]interface{}{"maxResults": iotstore.SearchIndexMaxResultsCap})
	if err != nil || start != 0 || end != iotstore.SearchIndexMaxResultsCap {
		t.Fatalf("maxResults=cap not honoured: start=%d end=%d err=%v", start, end, err)
	}
	start, end, err = applyPagination(500, map[string]interface{}{"maxResults": "100"})
	if err != nil || start != 0 || end != 100 {
		t.Fatalf("numeric-string maxResults not honoured: start=%d end=%d err=%v", start, end, err)
	}
}

func TestApplyPaginationHonoursNumericValues(t *testing.T) {
	start, end, err := applyPagination(10, map[string]interface{}{"maxResults": "3"})
	if err != nil {
		t.Fatal(err)
	}
	if start != 0 || end != 3 {
		t.Fatalf("numeric-string maxResults not honoured: start=%d end=%d", start, end)
	}

	start, end, err = applyPagination(10, map[string]interface{}{"maxResults": 4, "nextToken": "3"})
	if err != nil {
		t.Fatal(err)
	}
	if start != 3 || end != 7 {
		t.Fatalf("numeric nextToken not honoured: start=%d end=%d", start, end)
	}

	start, end, err = applyPagination(10, nil)
	if err != nil {
		t.Fatal(err)
	}
	if start != 0 || end != 10 {
		t.Fatalf("absent maxResults must default to the page size: start=%d end=%d", start, end)
	}
}
