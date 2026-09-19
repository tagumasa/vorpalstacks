package kinesis

import (
	"fmt"
	"testing"
	"time"
)

// seedRecord writes a record at a hand-chosen sequence number and arrival
// time so the retention pins can age data without waiting out a real
// retention window.
func seedRecord(t *testing.T, store *KinesisStore, streamName, shardID, seqNum string, arrival time.Time) {
	t.Helper()
	record := &Record{
		SequenceNumber:              seqNum,
		ApproximateArrivalTimestamp: arrival,
		Data:                        "data-" + seqNum,
		PartitionKey:                "pk",
	}
	if err := store.recordsStore.PutProto(fmt.Sprintf("%s#%s#%s", streamName, shardID, seqNum), RecordToProto(record)); err != nil {
		t.Fatalf("seed record: %v", err)
	}
}

// recordKeys censuses the stream's record keys in the store.
func recordKeys(t *testing.T, store *KinesisStore, streamName string) []string {
	t.Helper()
	var keys []string
	if err := store.recordsStore.ScanPrefix(streamName+"#", func(key string, _ []byte) error {
		keys = append(keys, key)
		return nil
	}); err != nil {
		t.Fatalf("census records: %v", err)
	}
	return keys
}

// TestGetRecordsFiltersBeforeTheLimitCut pins the age filter's position:
// the scan skips aged-out records first, then the limit cut bounds the page,
// so a reader behind a wall of expired data still receives its full page of
// in-retention records.
func TestGetRecordsFiltersBeforeTheLimitCut(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("ret_filter", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	shardID := "shardId-000000000000"

	aged := time.Now().UTC().Add(-48 * time.Hour)
	for i := 1; i <= 3; i++ {
		seedRecord(t, store, "ret_filter", shardID, formatSequenceNumber(aged.UnixNano(), int64(i), 0), aged)
	}
	var freshSeqs [3]string
	fresh := time.Now().UTC()
	for i := 1; i <= 3; i++ {
		freshSeqs[i-1] = formatSequenceNumber(fresh.UnixNano(), int64(i), 0)
		seedRecord(t, store, "ret_filter", shardID, freshSeqs[i-1], fresh)
	}

	page, advance, err := store.GetRecords("ret_filter", shardID, "", 2, true, time.Now().UTC().Add(-time.Hour))
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	if len(page) != 2 || page[0].SequenceNumber != freshSeqs[0] || page[1].SequenceNumber != freshSeqs[1] {
		t.Fatalf("page: got %d records starting at %v, want the first two in-retention records", len(page), pageSeqs(page))
	}
	if advance != freshSeqs[1] {
		t.Fatalf("advance: got %s, want the page's last consumed record %s", advance, freshSeqs[1])
	}
}

// TestGetRecordsAgedPageAdvancesToTip pins the advance position on a fully
// expired page: aged-out records do not enter the page but do advance the
// reader, so the follow-up iterator sits after the last aged record instead
// of re-reading the same expired window forever.
func TestGetRecordsAgedPageAdvancesToTip(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("ret_advance", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	shardID := "shardId-000000000000"

	aged := time.Now().UTC().Add(-48 * time.Hour)
	var lastAged string
	for i := 1; i <= 3; i++ {
		lastAged = formatSequenceNumber(aged.UnixNano(), int64(i), 0)
		seedRecord(t, store, "ret_advance", shardID, lastAged, aged)
	}

	page, advance, err := store.GetRecords("ret_advance", shardID, "", 10, true, time.Now().UTC().Add(-time.Hour))
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	if len(page) != 0 {
		t.Fatalf("page: got %d records, want none — every record is aged out", len(page))
	}
	if advance != lastAged {
		t.Fatalf("advance on a fully aged page: got %q, want the last aged record %q — the reader must move to the tip", advance, lastAged)
	}
}

// TestGetRecordsZeroCutoffReturnsRawHistory pins the AT_TIMESTAMP contract:
// a zero cutoff disables the age filter, so the position scan sees the
// shard's raw history regardless of retention.
func TestGetRecordsZeroCutoffReturnsRawHistory(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("ret_raw", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	shardID := "shardId-000000000000"

	aged := time.Now().UTC().Add(-48 * time.Hour)
	seq := formatSequenceNumber(aged.UnixNano(), 1, 0)
	seedRecord(t, store, "ret_raw", shardID, seq, aged)

	page, _, err := store.GetRecords("ret_raw", shardID, "", 10, true, time.Time{})
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	if len(page) != 1 || page[0].SequenceNumber != seq {
		t.Fatalf("raw history: got %d records, want the aged record visible under a zero cutoff", len(page))
	}
}

// TestTrimExpiredRecordsReclaimsAgedRecords pins the physical trim: records
// past the stream's retention window leave the store, in-window records
// stay, and the read path no longer sees the aged ones.
func TestTrimExpiredRecordsReclaimsAgedRecords(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("ret_trim", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	shardID := "shardId-000000000000"

	aged := time.Now().UTC().Add(-48 * time.Hour)
	var agedSeqs []string
	for i := 1; i <= 2; i++ {
		seq := formatSequenceNumber(aged.UnixNano(), int64(i), 0)
		agedSeqs = append(agedSeqs, seq)
		seedRecord(t, store, "ret_trim", shardID, seq, aged)
	}
	fresh := time.Now().UTC()
	var freshSeqs []string
	for i := 1; i <= 2; i++ {
		seq := formatSequenceNumber(fresh.UnixNano(), int64(i), 0)
		freshSeqs = append(freshSeqs, seq)
		seedRecord(t, store, "ret_trim", shardID, seq, fresh)
	}

	if err := store.TrimExpiredRecords("ret_trim"); err != nil {
		t.Fatalf("trim: %v", err)
	}

	keys := recordKeys(t, store, "ret_trim")
	if len(keys) != len(freshSeqs) {
		t.Fatalf("record census after trim: %d keys remain (%v), want the %d in-window records", len(keys), keys, len(freshSeqs))
	}
	for _, agedSeq := range agedSeqs {
		for _, key := range keys {
			if key == fmt.Sprintf("ret_trim#%s#%s", shardID, agedSeq) {
				t.Fatalf("aged record %s survived the trim", agedSeq)
			}
		}
	}

	page, _, err := store.GetRecords("ret_trim", shardID, "", 10, true, time.Time{})
	if err != nil {
		t.Fatalf("get records after trim: %v", err)
	}
	if len(page) != len(freshSeqs) {
		t.Fatalf("read after trim: got %d records, want the %d in-window records", len(page), len(freshSeqs))
	}
}

// pageSeqs renders a page's sequence numbers for failure messages.
func pageSeqs(page []*Record) []string {
	seqs := make([]string, len(page))
	for i, r := range page {
		seqs[i] = r.SequenceNumber
	}
	return seqs
}
