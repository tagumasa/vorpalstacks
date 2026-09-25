package dynamodb

import (
	"encoding/base64"
	"testing"
	"time"
)

// iteratorTestKey is a fixed signing key for the pure encode/decode tests;
// the store-backed key lifecycle is pinned by the stream store tests.
var iteratorTestKey = []byte("unit-test signing key")

// testIterStreamArn names the issuing stream generation in the pure
// encode/decode tests.
const testIterStreamArn = "arn:aws:dynamodb:us-east-1:123456789012:table/IteratorTable/stream/2026-01-01T00:00:00.000"

func TestShardIteratorRoundTripAndExpiry(t *testing.T) {
	iterator := encodeShardIterator(iteratorTestKey, testIterStreamArn, "IteratorTable", 7, "TRIM_HORIZON")
	streamArn, tableName, seq, issuedAt, iteratorType, err := decodeShardIterator(iteratorTestKey, iterator)
	if err != nil {
		t.Fatalf("decode iterator: %v", err)
	}
	if streamArn != testIterStreamArn || tableName != "IteratorTable" || seq != 7 || iteratorType != "TRIM_HORIZON" {
		t.Fatalf("expected %s/IteratorTable/7, got %s/%s/%d", testIterStreamArn, streamArn, tableName, seq)
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
	// A hand-crafted plaintext — whatever the forger guesses the internal
	// layout to be — carries no signature and must be rejected. The test
	// stays decoupled from the layout: the forged material is arbitrary.
	forged := base64.RawURLEncoding.EncodeToString([]byte("whatever a client might hand-assemble"))
	if _, _, _, _, _, err := decodeShardIterator(iteratorTestKey, forged); err == nil {
		t.Fatalf("hand-crafted plaintext iterator must be rejected")
	}

	issued := encodeShardIterator(iteratorTestKey, testIterStreamArn, "IteratorTable", 7, "LATEST")
	if _, _, _, _, _, err := decodeShardIterator(iteratorTestKey, issued); err != nil {
		t.Fatalf("issued iterator must decode: %v", err)
	}

	raw, err := base64.RawURLEncoding.DecodeString(issued)
	if err != nil {
		t.Fatalf("decode issued iterator: %v", err)
	}
	raw[len(raw)-1] ^= 0xFF
	tampered := base64.RawURLEncoding.EncodeToString(raw)
	if _, _, _, _, _, err := decodeShardIterator(iteratorTestKey, tampered); err == nil {
		t.Fatalf("tampered signature must be rejected")
	}

	// A real encoding with one payload byte altered keeps no valid
	// signature: the alteration is positional, not layout-shaped, so the
	// test never forges the internal payload form.
	altered := append([]byte(nil), raw...)
	altered[0] ^= 0xFF
	if _, _, _, _, _, err := decodeShardIterator(iteratorTestKey, base64.RawURLEncoding.EncodeToString(altered)); err == nil {
		t.Fatalf("altered payload must be rejected")
	}

	otherKey := []byte("another server key")
	if _, _, _, _, _, err := decodeShardIterator(otherKey, issued); err == nil {
		t.Fatalf("token signed by another key must be rejected")
	}

	if _, _, _, _, _, err := decodeShardIterator(iteratorTestKey, "IteratorTable|7|1|2"); err == nil {
		t.Fatalf("non-base64 plaintext iterator must be rejected")
	}
}
