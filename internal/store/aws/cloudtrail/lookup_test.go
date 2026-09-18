package cloudtrail

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func newLookupTestStore(t *testing.T) *CloudTrailStore {
	t.Helper()
	tmpDir := "./tmp/cloudtrail-lookup-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	s, err := storage.Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open test storage: %v", err)
	}
	t.Cleanup(func() { s.Close() })

	return NewCloudTrailStore(s, "acc123", "us-east-1")
}

func putTimedEvent(t *testing.T, store *CloudTrailStore, id string, at time.Time) {
	t.Helper()
	event := NewEvent("CreateTrail", "cloudtrail.amazonaws.com", nil, false)
	event.EventID = id
	event.EventTime = at.UTC()
	if err := store.PutEvent(event); err != nil {
		t.Fatalf("put event %s: %v", id, err)
	}
}

// LookupEvents returns the most recent event first on every retrieval path
// (fetched AWS contract: "The events list is sorted by time. The most
// recent event is listed first."), and the ordering survives pagination.
func TestLookupEventsReturnsNewestFirstOnEveryPath(t *testing.T) {
	store := newLookupTestStore(t)

	base := time.Date(2026, 9, 17, 2, 0, 0, 0, time.UTC)
	putTimedEvent(t, store, "evt-old-hour", base)                      // bucket 02, outside the span
	putTimedEvent(t, store, "evt-pre-bound", base.Add(65*time.Minute)) // bucket 03, below the start bound
	putTimedEvent(t, store, "evt-mid", base.Add(90*time.Minute))       // bucket 03, in window
	putTimedEvent(t, store, "evt-new-hour", base.Add(2*time.Hour))     // bucket 04, in window

	start := base.Add(75 * time.Minute) // inside hour 03: the pre-bound event shares its bucket
	end := base.Add(3 * time.Hour)
	events, next, err := store.LookupEvents(EventQuery{
		StartTime: &start, EndTime: &end, MaxResults: 2,
	})
	if err != nil {
		t.Fatalf("time lookup: %v", err)
	}
	if len(events) != 2 || events[0].EventID != "evt-new-hour" || events[1].EventID != "evt-mid" {
		t.Fatalf("first page = %v, want evt-new-hour then evt-mid", ids(events))
	}

	// Pagination continues below the last served event; the pre-bound
	// event is served by the walk but dropped by the bound filter, and
	// the token drains to exhaustion.
	pages := 0
	for next != "" && pages < 10 {
		var more []*Event
		more, next, err = store.LookupEvents(EventQuery{
			StartTime: &start, EndTime: &end, MaxResults: 2, NextToken: next,
		})
		if err != nil {
			t.Fatalf("time lookup page %d: %v", pages+2, err)
		}
		for _, e := range more {
			if e.EventID == "evt-pre-bound" || e.EventID == "evt-old-hour" {
				t.Fatalf("out-of-window event %s returned", e.EventID)
			}
		}
		pages++
	}
	if next != "" {
		t.Fatalf("token did not drain after %d pages", pages)
	}

	// The filterless scan path walks the whole events bucket newest first.
	scanEvents, scanNext, err := store.LookupEvents(EventQuery{MaxResults: 3})
	if err != nil {
		t.Fatalf("scan lookup: %v", err)
	}
	if len(scanEvents) != 3 || scanEvents[0].EventID != "evt-new-hour" ||
		scanEvents[1].EventID != "evt-mid" || scanEvents[2].EventID != "evt-pre-bound" {
		t.Fatalf("scan page = %s, want evt-new-hour, evt-mid, evt-pre-bound", ids(scanEvents))
	}
	if scanNext == "" {
		t.Fatal("scan page must paginate when events remain")
	}

	// A single-attribute path orders newest first within its bucket.
	attrEvents, _, err := store.LookupEvents(EventQuery{
		EventNames: []string{"CreateTrail"}, MaxResults: 4,
	})
	if err != nil {
		t.Fatalf("event-name lookup: %v", err)
	}
	if len(attrEvents) != 4 || attrEvents[0].EventID != "evt-new-hour" {
		t.Fatalf("event-name page = %s, want evt-new-hour first", ids(attrEvents))
	}
}

// A single-bound lookup covers the full range from its bound to the newest
// (or from the oldest to its bound) recorded event — never a single hour.
func TestLookupEventsSingleBoundCoversFullRange(t *testing.T) {
	store := newLookupTestStore(t)

	base := time.Date(2026, 9, 17, 2, 0, 0, 0, time.UTC)
	putTimedEvent(t, store, "evt-h2", base.Add(2*time.Hour))
	putTimedEvent(t, store, "evt-h0", base)

	// StartTime only: everything from the bound up, across hours.
	start := base.Add(time.Minute)
	events, _, err := store.LookupEvents(EventQuery{StartTime: &start, MaxResults: 10})
	if err != nil {
		t.Fatalf("start-only lookup: %v", err)
	}
	if len(events) != 1 || events[0].EventID != "evt-h2" {
		t.Fatalf("start-only = %s, want only evt-h2", ids(events))
	}

	// EndTime only: everything up to the bound, across hours.
	end := base.Add(2*time.Hour + time.Minute)
	events, _, err = store.LookupEvents(EventQuery{EndTime: &end, MaxResults: 10})
	if err != nil {
		t.Fatalf("end-only lookup: %v", err)
	}
	if len(events) != 2 || events[0].EventID != "evt-h2" || events[1].EventID != "evt-h0" {
		t.Fatalf("end-only = %s, want evt-h2 then evt-h0", ids(events))
	}
}

// A bound expressed with a timezone offset addresses the UTC buckets its
// instant defines — 11:30+09:00 is 02:30 UTC, not the 11:00 bucket.
func TestLookupEventsOffsetBoundAddressesUTCBuckets(t *testing.T) {
	store := newLookupTestStore(t)

	utc := time.Date(2026, 9, 17, 2, 45, 0, 0, time.UTC)
	putTimedEvent(t, store, "evt-utc", utc)

	tokyo := time.FixedZone("tokyo", 9*60*60)
	start := utc.Add(-time.Hour).In(tokyo) // 10:45+09:00 — the 02:45 UTC instant
	end := utc.Add(time.Hour).In(tokyo)
	events, _, err := store.LookupEvents(EventQuery{StartTime: &start, EndTime: &end, MaxResults: 10})
	if err != nil {
		t.Fatalf("offset-bound lookup: %v", err)
	}
	if len(events) != 1 || events[0].EventID != "evt-utc" {
		t.Fatalf("offset-bound = %s, want evt-utc", ids(events))
	}
}

// A NextToken this store never issued is refused with ErrInvalidNextToken;
// an EventId that matches nothing is a valid empty result, and a
// well-formed continuation still works.
func TestLookupEventsTokenAndEventIDPaths(t *testing.T) {
	store := newLookupTestStore(t)

	base := time.Date(2026, 9, 17, 2, 0, 0, 0, time.UTC)
	for i := 0; i < 3; i++ {
		putTimedEvent(t, store, fmt.Sprintf("evt-%d", i), base.Add(time.Duration(i)*time.Minute))
	}

	if _, _, err := store.LookupEvents(EventQuery{MaxResults: 10, NextToken: "not-a-token"}); !errors.Is(err, ErrInvalidNextToken) {
		t.Fatalf("garbage token: error = %v, want ErrInvalidNextToken", err)
	}
	if _, _, err := store.LookupEvents(EventQuery{MaxResults: 10, NextToken: "idx1:%%%"}); !errors.Is(err, ErrInvalidNextToken) {
		t.Fatalf("corrupt cursor payload: error = %v, want ErrInvalidNextToken", err)
	}

	events, _, err := store.LookupEvents(EventQuery{EventID: "evt-1"})
	if err != nil || len(events) != 1 || events[0].EventID != "evt-1" {
		t.Fatalf("event-id lookup = %s (err %v), want evt-1", ids(events), err)
	}
	events, _, err = store.LookupEvents(EventQuery{EventID: "evt-absent"})
	if err != nil {
		t.Fatalf("absent event-id lookup must be a valid empty result, got %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("absent event-id lookup = %s, want empty", ids(events))
	}
}

func ids(events []*Event) []string {
	var out []string
	for _, e := range events {
		out = append(out, e.EventID)
	}
	return out
}
