package sfn

import (
	"context"
	"testing"
)

// TestMapRunSeqFromARN pins the shared run-identifier parser: the sequence
// is the number between "mapRun-" and the timestamp of the run id, and an
// ARN without a run id reports false rather than a fake zero sequence.
func TestMapRunSeqFromARN(t *testing.T) {
	cases := []struct {
		arn  string
		want int64
		ok   bool
	}{
		{"arn:aws:states:us-east-1:000000000000:mapRun:sm/M/mapRun-12-20260913120000", 12, true},
		{"arn:aws:states:us-east-1:000000000000:mapRun:sm/M/mapRun-9-20260913120001", 9, true},
		{"arn:aws:states:us-east-1:000000000000:mapRun:sm/M/mapRun-1048576-20260913120002", 1048576, true},
		{"arn:aws:states:us-east-1:000000000000:mapRun:recover-sm/existing", 0, false},
		{"arn:aws:states:us-east-1:000000000000:execution:sm/e1", 0, false},
	}
	for _, tc := range cases {
		got, ok := MapRunSeqFromARN(tc.arn)
		if ok != tc.ok || got != tc.want {
			t.Errorf("MapRunSeqFromARN(%s) = %d,%v want %d,%v", tc.arn, got, ok, tc.want, tc.ok)
		}
	}
}

// TestNextMapRunSeqRecoversHighestFromStorage pins the boot recovery of
// the run sequence: a store opened over existing runs continues past the
// highest stored sequence instead of restarting at one.
func TestNextMapRunSeqRecoversHighestFromStorage(t *testing.T) {
	store := newHistoryTestStore(t)
	for _, arn := range []string{
		"arn:aws:states:us-east-1:000000000000:mapRun:sm/M/mapRun-9-20260913120000",
		"arn:aws:states:us-east-1:000000000000:mapRun:sm/M/mapRun-42-20260913120100",
	} {
		if err := store.CreateMapRun(context.Background(), &MapRun{MapRunArn: arn, ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:sm/e", Name: "M"}); err != nil {
			t.Fatalf("seed map run: %v", err)
		}
	}
	if got := store.NextMapRunSeq(); got != 43 {
		t.Errorf("NextMapRunSeq after recovery = %d, want 43 (highest stored 42 + 1)", got)
	}
}

// TestListAllMapRunsExcludesInlineBeforePagination pins that inline-map
// checkpoint records never occupy a page slot: with inline records
// interleaved among describable runs, every page is filled entirely with
// visible Map Runs.
func TestListAllMapRunsExcludesInlineBeforePagination(t *testing.T) {
	store := newHistoryTestStore(t)
	execArn := "arn:aws:states:us-east-1:000000000000:execution:sm/e"
	base := "arn:aws:states:us-east-1:000000000000:mapRun:sm/M/mapRun-"
	// Keys interleave: seq 2 and 4 are inline checkpoints, 1/3/5/6 are
	// describable runs, so a post-pagination filter would return empty or
	// short pages while the inline records consumed the page size.
	seed := []struct {
		arn    string
		inline bool
	}{
		{base + "1-20260913120000", false},
		{base + "2-20260913120001", true},
		{base + "3-20260913120002", false},
		{base + "4-20260913120003", true},
		{base + "5-20260913120004", false},
		{base + "6-20260913120005", false},
	}
	for _, s := range seed {
		if err := store.CreateMapRun(context.Background(), &MapRun{MapRunArn: s.arn, ExecutionArn: execArn, Name: "M", Inline: s.inline}); err != nil {
			t.Fatalf("seed map run: %v", err)
		}
	}

	seen := []string{}
	token := ""
	pages := 0
	for {
		page, err := store.ListAllMapRuns(context.Background(), execArn, 2, token)
		if err != nil {
			t.Fatalf("list page: %v", err)
		}
		for _, mr := range page.MapRuns {
			if mr.Inline {
				t.Fatalf("inline record %s surfaced through ListAllMapRuns", mr.MapRunArn)
			}
			seen = append(seen, mr.MapRunArn)
		}
		if len(page.MapRuns) > 2 {
			t.Fatalf("page returned %d records, want at most the 2-record page size", len(page.MapRuns))
		}
		pages++
		if page.NextToken == "" {
			break
		}
		token = page.NextToken
		if pages > 10 {
			t.Fatal("pagination did not terminate")
		}
	}
	if len(seen) != 4 {
		t.Errorf("visible runs = %d (%v), want the 4 non-inline records", len(seen), seen)
	}
}
