package dynamodb

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	crypto "vorpalstacks/internal/utils/crypto"
)

// iteratorTestKey is a fixed signing key for the pure encode/decode tests;
// the store-backed key lifecycle is pinned by the stream store tests.
var iteratorTestKey = []byte("unit-test signing key")

func TestShardIteratorRoundTripAndExpiry(t *testing.T) {
	iterator := encodeShardIterator(iteratorTestKey, "IteratorTable", 7)
	tableName, seq, issuedAt, err := decodeShardIterator(iteratorTestKey, iterator)
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
}

// The iterator is opaque and server-issued: a token a client could have
// constructed by hand must be rejected, and so must an issued token whose
// payload or signature was altered, under the signing key of the issuing
// server or any other key.
func TestShardIteratorRejectsForgedAndTamperedTokens(t *testing.T) {
	forged := base64.RawURLEncoding.EncodeToString([]byte("VictimTable|0|" + time.Now().Format("150405")))
	if _, _, _, err := decodeShardIterator(iteratorTestKey, forged); err == nil {
		t.Fatalf("hand-crafted plaintext iterator must be rejected")
	}

	issued := encodeShardIterator(iteratorTestKey, "IteratorTable", 7)
	if _, _, _, err := decodeShardIterator(iteratorTestKey, issued); err != nil {
		t.Fatalf("issued iterator must decode: %v", err)
	}

	raw, err := base64.RawURLEncoding.DecodeString(issued)
	if err != nil {
		t.Fatalf("decode issued iterator: %v", err)
	}
	raw[len(raw)-1] ^= 0xFF
	tampered := base64.RawURLEncoding.EncodeToString(raw)
	if _, _, _, err := decodeShardIterator(iteratorTestKey, tampered); err == nil {
		t.Fatalf("tampered signature must be rejected")
	}

	// A payload rewritten to another table keeps no valid signature.
	rewritten := base64.RawURLEncoding.EncodeToString(
		append([]byte("VictimTable|0|"+strings.Repeat("9", 10)), crypto.HMACSHA256(iteratorTestKey, []byte("IteratorTable|7|1"))...))
	if _, _, _, err := decodeShardIterator(iteratorTestKey, rewritten); err == nil {
		t.Fatalf("rewritten payload must be rejected")
	}

	otherKey := []byte("another server key")
	if _, _, _, err := decodeShardIterator(otherKey, issued); err == nil {
		t.Fatalf("token signed by another key must be rejected")
	}

	if _, _, _, err := decodeShardIterator(iteratorTestKey, "IteratorTable|7"); err == nil {
		t.Fatalf("non-base64 plaintext iterator must be rejected")
	}
}
