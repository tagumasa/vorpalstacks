package cloudwatchlogs

import (
	"testing"

	"vorpalstacks/internal/core/storage"
)

// PutLogEvents acknowledges a write only after it is durable: a store that
// reopens over the same storage (as a restarted process does) must see every
// acknowledged event without depending on any in-process buffer.
func TestPutLogEventsDurableAcrossStoreReopen(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	dataPath := t.TempDir()

	writer, err := NewStore(st, st.Bucket("logs-us-east-1"), "000000000000", "us-east-1", dataPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}
	if err := writer.CreateLogStream(&LogStream{LogGroupName: "group", Name: "stream"}); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.PutLogEvents("group", "stream", []LogEntry{
		{Timestamp: 1000, Message: "first"},
		{Timestamp: 2000, Message: "second"},
	}); err != nil {
		t.Fatal(err)
	}

	reader, err := NewStore(st, st.Bucket("logs-us-east-1"), "000000000000", "us-east-1", dataPath)
	if err != nil {
		t.Fatal(err)
	}
	events, _, _, err := reader.GetLogEvents("group", "stream", 0, 0, 0, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("reopened store must see both acknowledged events, got %d", len(events))
	}
	if events[0].Message != "first" || events[1].Message != "second" {
		t.Fatalf("unexpected event order: %q then %q", events[0].Message, events[1].Message)
	}
}
