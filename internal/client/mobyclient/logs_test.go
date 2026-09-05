package mobyclient

import (
	"encoding/binary"
	"strings"
	"testing"
	"time"
)

func TestParseLogWindowStripsTimestamps(t *testing.T) {
	since := time.Date(2026, 9, 4, 1, 2, 3, 0, time.UTC)
	raw := "2026-09-04T01:02:03.100000000Z INIT start\n" +
		"2026-09-04T01:02:04.200000001Z handler done\n"

	text, cursor, err := parseLogWindow(raw, since)
	if err != nil {
		t.Fatalf("parseLogWindow: %v", err)
	}
	wantText := "INIT start\nhandler done\n"
	if text != wantText {
		t.Fatalf("text = %q, want %q", text, wantText)
	}
	wantCursor := time.Date(2026, 9, 4, 1, 2, 4, 200000001, time.UTC)
	if !cursor.Equal(wantCursor) {
		t.Fatalf("cursor = %v, want the newest line timestamp %v", cursor, wantCursor)
	}
}

func TestParseLogWindowEmptyKeepsCursor(t *testing.T) {
	since := time.Date(2026, 9, 4, 1, 2, 3, 5, time.UTC)
	text, cursor, err := parseLogWindow("", since)
	if err != nil {
		t.Fatalf("parseLogWindow: %v", err)
	}
	if text != "" {
		t.Fatalf("text = %q, want empty", text)
	}
	if !cursor.Equal(since) {
		t.Fatalf("cursor = %v, want %v unchanged", cursor, since)
	}
}

func TestParseLogWindowRejectsUnstampedLine(t *testing.T) {
	_, _, err := parseLogWindow("not a timestamped line\n", time.Time{})
	if err == nil {
		t.Fatalf("a line without a docker timestamp must be an error, not silent corruption")
	}
}

func TestParseLogWindowNewestTimestampWins(t *testing.T) {
	// Log order is not guaranteed to be monotonic across the demultiplexed
	// stdout/stderr streams; the cursor must still carry the newest one.
	raw := "2026-09-04T01:02:05.000000000Z second stream\n" +
		"2026-09-04T01:02:04.000000000Z first stream\n"
	_, cursor, err := parseLogWindow(raw, time.Time{})
	if err != nil {
		t.Fatalf("parseLogWindow: %v", err)
	}
	want := time.Date(2026, 9, 4, 1, 2, 5, 0, time.UTC)
	if !cursor.Equal(want) {
		t.Fatalf("cursor = %v, want the newest timestamp %v", cursor, want)
	}
}

// TestParseLogWindowExcludesTheReDeliveredBoundary pins the compensation
// for docker's inclusive since filter: the record at exactly `since` is
// the previous window's newest line re-delivered and must not appear,
// while everything after it must.
func TestParseLogWindowExcludesTheReDeliveredBoundary(t *testing.T) {
	since := time.Date(2026, 9, 4, 1, 2, 4, 0, time.UTC)
	raw := "2026-09-04T01:02:03.000000000Z round one output\n" +
		"2026-09-04T01:02:04.000000000Z boundary output\n" +
		"2026-09-04T01:02:05.000000000Z round two output\n"

	text, cursor, err := parseLogWindow(raw, since)
	if err != nil {
		t.Fatalf("parseLogWindow: %v", err)
	}
	if strings.Contains(text, "boundary output") || strings.Contains(text, "round one output") {
		t.Fatalf("records at or before since were re-delivered: %q", text)
	}
	if !strings.Contains(text, "round two output") {
		t.Fatalf("the record after the boundary was dropped: %q", text)
	}
	wantCursor := time.Date(2026, 9, 4, 1, 2, 5, 0, time.UTC)
	if !cursor.Equal(wantCursor) {
		t.Fatalf("cursor = %v, want %v", cursor, wantCursor)
	}
}

// TestParseLogWindowDropsStragglersAtOrBelowTheCursor pins the filter
// against the non-monotonic record order docker can produce across the
// demultiplexed streams: a record stamped at or below the cursor was
// already inside the previous window's read no matter where it appears in
// the file, so it must be dropped wherever it surfaces.
func TestParseLogWindowDropsStragglersAtOrBelowTheCursor(t *testing.T) {
	since := time.Date(2026, 9, 4, 1, 2, 5, 0, time.UTC)
	raw := "2026-09-04T01:02:05.000000000Z boundary output\n" +
		"2026-09-04T01:02:03.000000000Z early straggler\n" +
		"2026-09-04T01:02:06.000000000Z new record\n" +
		"2026-09-04T01:02:04.000000000Z late straggler\n"

	text, cursor, err := parseLogWindow(raw, since)
	if err != nil {
		t.Fatalf("parseLogWindow: %v", err)
	}
	if strings.Contains(text, "boundary output") ||
		strings.Contains(text, "early straggler") ||
		strings.Contains(text, "late straggler") {
		t.Fatalf("records at or below the cursor were re-delivered: %q", text)
	}
	if !strings.Contains(text, "new record") {
		t.Fatalf("the record after the boundary was dropped: %q", text)
	}
	wantCursor := time.Date(2026, 9, 4, 1, 2, 6, 0, time.UTC)
	if !cursor.Equal(wantCursor) {
		t.Fatalf("cursor = %v, want %v", cursor, wantCursor)
	}
}

// TestParseLogWindowBoundaryOnlyWindowKeepsCursor pins the empty window: a
// read that only sees the re-delivered boundary returns no text and keeps
// the cursor exactly where it was, so the next window retries the same
// boundary cleanly.
func TestParseLogWindowBoundaryOnlyWindowKeepsCursor(t *testing.T) {
	since := time.Date(2026, 9, 4, 1, 2, 4, 0, time.UTC)
	raw := "2026-09-04T01:02:03.000000000Z old output\n" +
		"2026-09-04T01:02:04.000000000Z boundary output\n"

	text, cursor, err := parseLogWindow(raw, since)
	if err != nil {
		t.Fatalf("parseLogWindow: %v", err)
	}
	if text != "" {
		t.Fatalf("text = %q, want empty", text)
	}
	if !cursor.Equal(since) {
		t.Fatalf("cursor = %v, want %v unchanged", cursor, since)
	}
}

// frame builds one docker non-TTY stream frame of the given stream byte.
func frame(stream byte, payload string) []byte {
	buf := make([]byte, 8+len(payload))
	buf[0] = stream
	binary.BigEndian.PutUint32(buf[4:8], uint32(len(payload)))
	copy(buf[8:], payload)
	return buf
}

func TestDemuxLogStreamMergesStdoutStderrInWriteOrder(t *testing.T) {
	raw := append(frame(1, "out line\n"), frame(2, "err line\n")...)
	raw = append(raw, frame(1, "out again\n")...)

	got := demuxLogStream(raw)
	want := "out line\nerr line\nout again\n"
	if got != want {
		t.Fatalf("demux = %q, want %q", got, want)
	}
}

func TestDemuxLogStreamEmptyAndTruncated(t *testing.T) {
	if got := demuxLogStream(nil); got != "" {
		t.Fatalf("empty stream demuxed to %q", got)
	}
	// A trailing truncated frame contributes nothing decodable.
	got := demuxLogStream(append(frame(1, "whole\n"), 0x01, 0x00, 0x00))
	if got != "whole\n" {
		t.Fatalf("truncated-tail demux = %q, want the whole frame only", got)
	}
}
