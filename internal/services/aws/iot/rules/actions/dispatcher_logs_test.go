package actions

import (
	"testing"
)

// parseBatchLogEntries must accept the documented batchMode array format
// ({"timestamp": <ms>, "message": "..."} records) and reject payloads that
// do not conform, because AWS requires every record to carry a timestamp
// and a message.
func TestParseBatchLogEntries(t *testing.T) {
	entries, err := parseBatchLogEntries([]byte(`[
		{"timestamp": 1673520691093, "message": "Test message 1"},
		{"timestamp": 1673520692879, "message": "Test message 2"},
		{"timestamp": 1673520693442, "message": "Test message 3"}
	]`))
	if err != nil {
		t.Fatalf("valid array rejected: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Timestamp != 1673520691093 || entries[0].Message != "Test message 1" {
		t.Fatalf("first entry mismatch: %+v", entries[0])
	}
	if entries[2].Timestamp != 1673520693442 || entries[2].Message != "Test message 3" {
		t.Fatalf("last entry mismatch: %+v", entries[2])
	}
}

func TestParseBatchLogEntriesEmptyArrayIsNoOp(t *testing.T) {
	entries, err := parseBatchLogEntries([]byte(`[]`))
	if err != nil {
		t.Fatalf("empty array must be a successful no-op: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestParseBatchLogEntriesRejectsNonArray(t *testing.T) {
	if _, err := parseBatchLogEntries([]byte(`{"message":"not an array"}`)); err == nil {
		t.Fatal("non-array payload must be rejected")
	}
	if _, err := parseBatchLogEntries([]byte(`"string"`)); err == nil {
		t.Fatal("string payload must be rejected")
	}
}

func TestParseBatchLogEntriesRejectsMissingTimestamp(t *testing.T) {
	if _, err := parseBatchLogEntries([]byte(`[{"message":"no timestamp"}]`)); err == nil {
		t.Fatal("record without timestamp must be rejected")
	}
	// A missing message is an empty message, not a malformed record.
	entries, err := parseBatchLogEntries([]byte(`[{"timestamp":1673520691093}]`))
	if err != nil {
		t.Fatalf("record without message must be accepted: %v", err)
	}
	if len(entries) != 1 || entries[0].Message != "" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

// cloudWatchLogsEntries must route by batchMode: the whole payload as one
// entry when off, per-record entries when on.
func TestCloudWatchLogsEntriesBatchModeSwitch(t *testing.T) {
	single := &ActionConfig{Extra: map[string]interface{}{}}
	entries, err := cloudWatchLogsEntries(single, &ActionPayload{JSONString: `{"value":7}`, JSONBytes: []byte(`{"value":7}`)})
	if err != nil {
		t.Fatalf("batchMode off failed: %v", err)
	}
	if len(entries) != 1 || entries[0].Message != `{"value":7}` {
		t.Fatalf("batchMode off must emit the whole payload once: %+v", entries)
	}

	batch := &ActionConfig{Extra: map[string]interface{}{"batchMode": true}}
	entries, err = cloudWatchLogsEntries(batch, &ActionPayload{JSONBytes: []byte(`[{"timestamp":5,"message":"a"},{"timestamp":6,"message":"b"}]`)})
	if err != nil {
		t.Fatalf("batchMode on failed: %v", err)
	}
	if len(entries) != 2 || entries[0].Message != "a" || entries[1].Timestamp != 6 {
		t.Fatalf("batchMode on must expand records: %+v", entries)
	}
}
