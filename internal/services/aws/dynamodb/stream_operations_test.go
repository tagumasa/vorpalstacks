package dynamodb

import (
	"testing"
	"time"
)

func TestShardIteratorRoundTripAndExpiry(t *testing.T) {
	iterator := encodeShardIterator("IteratorTable", 7)
	tableName, seq, issuedAt, err := decodeShardIterator(iterator)
	if err != nil {
		t.Fatalf("decode iterator: %v", err)
	}
	if tableName != "IteratorTable" || seq != 7 {
		t.Fatalf("expected IteratorTable/7, got %s/%d", tableName, seq)
	}

	issued := time.Unix(issuedAt, 0)
	if shardIteratorExpired(issuedAt, issued) {
		t.Fatalf("iterator must be valid at issue time")
	}
	if shardIteratorExpired(issuedAt, issued.Add(899*time.Second)) {
		t.Fatalf("iterator must still be valid one second before expiry")
	}
	if !shardIteratorExpired(issuedAt, issued.Add(900*time.Second)) {
		t.Fatalf("iterator must expire after fifteen minutes")
	}

	// Iterators without an issue timestamp are not decodable.
	if _, _, _, err := decodeShardIterator("IteratorTable|7"); err == nil {
		t.Fatalf("two-part iterator must be rejected")
	}
}
