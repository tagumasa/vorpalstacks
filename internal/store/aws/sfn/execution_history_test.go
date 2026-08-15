package sfn

import (
	"context"
	"fmt"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func newHistoryTestStore(t *testing.T) *StepFunctionStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return NewStepFunctionStore(st, "000000000000", "us-east-1")
}

func seedHistory(t *testing.T, store *StepFunctionStore, arn string, count int) {
	t.Helper()
	for i := int64(1); i <= int64(count); i++ {
		event := &ExecutionHistoryEvent{
			ExecutionArn: arn,
			EventId:      i,
			Type:         fmt.Sprintf("Event%d", i),
			Timestamp:    time.Now().UTC(),
		}
		if err := store.AddExecutionHistoryEvent(context.Background(), event); err != nil {
			t.Fatal(err)
		}
	}
}

func eventIDs(events []*ExecutionHistoryEvent) []int64 {
	ids := make([]int64, len(events))
	for i, e := range events {
		ids[i] = e.EventId
	}
	return ids
}

func equalIDs(got []int64, want ...int64) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// History keys zero-pad the event ID, so paginating an execution with more
// than nine events must return them in ascending numeric order; an
// unpadded decimal key would sort event 10 between 1 and 2.
func TestExecutionHistoryForwardOrderAcrossPages(t *testing.T) {
	store := newHistoryTestStore(t)
	arn := "arn:aws:states:us-east-1:000000000000:execution:sm:exec"
	seedHistory(t, store, arn, 12)

	page1, next1, err := store.GetExecutionHistory(context.Background(), arn, 5, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if !equalIDs(eventIDs(page1), 1, 2, 3, 4, 5) {
		t.Fatalf("page 1 order wrong: %v", eventIDs(page1))
	}
	if next1 == "" {
		t.Fatal("expected continuation marker")
	}

	page2, next2, err := store.GetExecutionHistory(context.Background(), arn, 5, next1, false)
	if err != nil {
		t.Fatal(err)
	}
	if !equalIDs(eventIDs(page2), 6, 7, 8, 9, 10) {
		t.Fatalf("page 2 order wrong: %v", eventIDs(page2))
	}

	page3, next3, err := store.GetExecutionHistory(context.Background(), arn, 5, next2, false)
	if err != nil {
		t.Fatal(err)
	}
	if !equalIDs(eventIDs(page3), 11, 12) {
		t.Fatalf("page 3 order wrong: %v", eventIDs(page3))
	}
	if next3 != "" {
		t.Fatalf("expected exhausted marker, got %q", next3)
	}
}

// reverseOrder must reverse the whole history, not just each page: pages
// walk from the newest event backwards with a direction-consistent marker.
func TestExecutionHistoryReverseOrderAcrossPages(t *testing.T) {
	store := newHistoryTestStore(t)
	arn := "arn:aws:states:us-east-1:000000000000:execution:sm:exec"
	seedHistory(t, store, arn, 12)

	page1, next1, err := store.GetExecutionHistory(context.Background(), arn, 5, "", true)
	if err != nil {
		t.Fatal(err)
	}
	if !equalIDs(eventIDs(page1), 12, 11, 10, 9, 8) {
		t.Fatalf("reverse page 1 order wrong: %v", eventIDs(page1))
	}

	page2, next2, err := store.GetExecutionHistory(context.Background(), arn, 5, next1, true)
	if err != nil {
		t.Fatal(err)
	}
	if !equalIDs(eventIDs(page2), 7, 6, 5, 4, 3) {
		t.Fatalf("reverse page 2 order wrong: %v", eventIDs(page2))
	}

	page3, next3, err := store.GetExecutionHistory(context.Background(), arn, 5, next2, true)
	if err != nil {
		t.Fatal(err)
	}
	if !equalIDs(eventIDs(page3), 2, 1) {
		t.Fatalf("reverse page 3 order wrong: %v", eventIDs(page3))
	}
	if next3 != "" {
		t.Fatalf("expected exhausted marker, got %q", next3)
	}
}

func TestExecutionHistoryReverseRejectsInvalidToken(t *testing.T) {
	store := newHistoryTestStore(t)
	arn := "arn:aws:states:us-east-1:000000000000:execution:sm:exec"
	seedHistory(t, store, arn, 3)

	if _, _, err := store.GetExecutionHistory(context.Background(), arn, 5, "not-a-number", true); err == nil {
		t.Fatal("non-numeric reverse marker accepted")
	}
	if _, _, err := store.GetExecutionHistory(context.Background(), arn, 5, "-1", true); err == nil {
		t.Fatal("negative reverse marker accepted")
	}
}

// The reverse marker anchors to a fixed event ID, so appending events
// between pages must neither duplicate the previous page's tail nor skip
// events: a count-from-the-newest-end marker would shift under growth and
// re-emit the boundary event.
func TestExecutionHistoryReverseStableUnderAppends(t *testing.T) {
	store := newHistoryTestStore(t)
	arn := "arn:aws:states:us-east-1:000000000000:execution:sm:exec"
	seedHistory(t, store, arn, 12)

	page1, next1, err := store.GetExecutionHistory(context.Background(), arn, 5, "", true)
	if err != nil {
		t.Fatal(err)
	}
	if !equalIDs(eventIDs(page1), 12, 11, 10, 9, 8) {
		t.Fatalf("reverse page 1 order wrong: %v", eventIDs(page1))
	}

	// The execution grows between pages.
	seedHistoryFrom(t, store, arn, 13, 14)

	page2, next2, err := store.GetExecutionHistory(context.Background(), arn, 5, next1, true)
	if err != nil {
		t.Fatal(err)
	}
	if !equalIDs(eventIDs(page2), 7, 6, 5, 4, 3) {
		t.Fatalf("reverse page 2 wrong after appends: %v", eventIDs(page2))
	}

	page3, next3, err := store.GetExecutionHistory(context.Background(), arn, 5, next2, true)
	if err != nil {
		t.Fatal(err)
	}
	if !equalIDs(eventIDs(page3), 2, 1) {
		t.Fatalf("reverse page 3 wrong: %v", eventIDs(page3))
	}
	if next3 != "" {
		t.Fatalf("expected exhausted marker, got %q", next3)
	}
}

func seedHistoryFrom(t *testing.T, store *StepFunctionStore, arn string, from, to int64) {
	t.Helper()
	for i := from; i <= to; i++ {
		event := &ExecutionHistoryEvent{
			ExecutionArn: arn,
			EventId:      i,
			Type:         fmt.Sprintf("Event%d", i),
			Timestamp:    time.Now().UTC(),
		}
		if err := store.AddExecutionHistoryEvent(context.Background(), event); err != nil {
			t.Fatal(err)
		}
	}
}
