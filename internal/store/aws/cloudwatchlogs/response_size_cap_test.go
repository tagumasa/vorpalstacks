package cloudwatchlogs

import (
	"strings"
	"testing"
	"time"
)

// TestGetLogEventsResponseSizeCap pins the documented GetLogEvents default
// page budget: "as many log events as can fit in a response size of 1 MB
// (up to 10,000 log events)". A stream holding more than the budget is
// served in pages whose serialized estimate stays within it, pagination
// still reaches every event, and the page keeps at least one event.
func TestGetLogEventsResponseSizeCap(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(NewLogGroup("cap-group", "us-east-1", "000000000000")); err != nil {
		t.Fatalf("group: %v", err)
	}
	if err := s.CreateLogStream(&LogStream{Name: "s", LogGroupName: "cap-group"}); err != nil {
		t.Fatalf("stream: %v", err)
	}

	// Two batches keep each PutLogEvents write under the batch ceiling
	// while the stream itself exceeds the 1 MB read budget.
	now := time.Now().UnixMilli()
	for i := 0; i < 2; i++ {
		entries := make([]LogEntry, 3)
		for j := range entries {
			entries[j] = LogEntry{
				Timestamp:     now + int64(i*3+j),
				Message:       strings.Repeat("x", 200000),
				IngestionTime: now,
			}
		}
		if _, err := s.PutLogEvents("cap-group", "s", entries); err != nil {
			t.Fatalf("seed batch %d: %v", i, err)
		}
	}

	events, nextForward, _, err := s.GetLogEvents("cap-group", "s", 0, 0, 0, true, "")
	if err != nil {
		t.Fatalf("first page: %v", err)
	}
	if len(events) == 0 || len(events) >= 6 {
		t.Fatalf("first page should hold a strict subset of the six events, got %d", len(events))
	}
	estimated := 0
	for _, e := range events {
		estimated += len(e.Message) + EventResponseEnvelopeBytes
	}
	if estimated > MaxReadResponseBytes {
		t.Fatalf("page estimate %d exceeds the %d budget", estimated, MaxReadResponseBytes)
	}

	served := len(events)
	token := nextForward
	for {
		page, forward, _, err := s.GetLogEvents("cap-group", "s", 0, 0, 0, true, token)
		if err != nil {
			t.Fatalf("page: %v", err)
		}
		if forward == token || len(page) == 0 {
			break
		}
		served += len(page)
		token = forward
	}
	if served != 6 {
		t.Fatalf("pagination should serve all six events, served %d", served)
	}

	// An explicit count limit below the byte budget still governs.
	events, _, _, err = s.GetLogEvents("cap-group", "s", 0, 0, 1, true, "")
	if err != nil || len(events) != 1 {
		t.Fatalf("explicit limit of 1: events=%d err=%v", len(events), err)
	}
}

// TestFilterLogEventsResponseSizeCap pins the same documented budget on
// the FilterLogEvents plane: "Each page returned can contain up to 1 MB of
// log events or up to 10,000 log events." The filtered walk pages within
// the byte budget and pagination still reaches every event.
func TestFilterLogEventsResponseSizeCap(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(NewLogGroup("cap-filter-group", "us-east-1", "000000000000")); err != nil {
		t.Fatalf("group: %v", err)
	}
	if err := s.CreateLogStream(&LogStream{Name: "s", LogGroupName: "cap-filter-group"}); err != nil {
		t.Fatalf("stream: %v", err)
	}

	now := time.Now().UnixMilli()
	for i := 0; i < 2; i++ {
		entries := make([]LogEntry, 3)
		for j := range entries {
			entries[j] = LogEntry{
				Timestamp:     now + int64(i*3+j),
				Message:       strings.Repeat("x", 200000),
				IngestionTime: now,
			}
		}
		if _, err := s.PutLogEvents("cap-filter-group", "s", entries); err != nil {
			t.Fatalf("seed batch %d: %v", i, err)
		}
	}

	events, next, err := s.FilterLogEvents("cap-filter-group", nil, 0, 0, "", 0, true, "")
	if err != nil {
		t.Fatalf("first page: %v", err)
	}
	if len(events) == 0 || len(events) >= 6 {
		t.Fatalf("first page should hold a strict subset of the six events, got %d", len(events))
	}
	estimated := 0
	for _, e := range events {
		estimated += len(e.Message) + EventResponseEnvelopeBytes
	}
	if estimated > MaxReadResponseBytes {
		t.Fatalf("page estimate %d exceeds the %d budget", estimated, MaxReadResponseBytes)
	}

	served := len(events)
	// FilterLogEvents mints a token only while more pages remain; the walk
	// ends at the empty token (calling back with it would restart from the
	// head), and the page it arrives with still counts.
	for token := next; token != ""; {
		page, nextToken, err := s.FilterLogEvents("cap-filter-group", nil, 0, 0, "", 0, true, token)
		if err != nil {
			t.Fatalf("page: %v", err)
		}
		served += len(page)
		if served > 6 {
			t.Fatalf("pagination over-served: %d", served)
		}
		token = nextToken
	}
	if served != 6 {
		t.Fatalf("pagination should serve all six events, served %d", served)
	}
}
