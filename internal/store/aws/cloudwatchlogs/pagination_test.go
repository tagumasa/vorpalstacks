package cloudwatchlogs

import (
	"testing"
	"time"

	"vorpalstacks/internal/common/pagination"
)

// The event-token expiry: "The token expires after 24 hours" (the read
// operations' nextToken members). A token minted past the expiry is as
// invalid as one from a foreign vocabulary, and the wire form that
// predates the mint field rejects with it.
func TestEventPageTokenExpiry(t *testing.T) {
	restore := eventTokenNowMs
	defer func() { eventTokenNowMs = restore }()

	eventTokenNowMs = func() int64 { return time.Now().Add(-25 * time.Hour).UnixMilli() }
	stale, err := encodeEventPageToken(eventPageToken{Group: "g", Direction: PageForward})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := decodeEventPageToken(stale); err == nil {
		t.Fatal("token minted 25 hours ago accepted")
	}

	// A mint just inside the window still decodes.
	eventTokenNowMs = func() int64 { return time.Now().Add(-23 * time.Hour).UnixMilli() }
	aged, err := encodeEventPageToken(eventPageToken{Group: "g", Direction: PageForward})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := decodeEventPageToken(aged); err != nil {
		t.Fatalf("token minted 23 hours ago rejected: %v", err)
	}

	eventTokenNowMs = func() int64 { return time.Now().UnixMilli() }
	fresh, err := encodeEventPageToken(eventPageToken{Group: "g", Direction: PageForward})
	if err != nil {
		t.Fatal(err)
	}
	tok, err := decodeEventPageToken(fresh)
	if err != nil {
		t.Fatalf("fresh token rejected: %v", err)
	}
	if tok.MintedMs <= 0 {
		t.Fatalf("mint = %d, want the stamped clock", tok.MintedMs)
	}

	// The wire form without the mint field (the pre-expiry vocabulary)
	// rejects: the version moved with the field.
	legacy, err := pagination.EncodeScopedToken(eventPageToken{
		Version: 3, Group: "g", Direction: PageForward,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := decodeEventPageToken(legacy); err == nil {
		t.Fatal("mint-less vocabulary accepted")
	}
}

// The end-of-stream re-offer contract: a walk's tokens carry the mint of
// the walk's first page, so re-encoding an unchanged scope reproduces
// the presented token byte for byte and the client's equality check —
// the documented pagination-finished signal — holds.
func TestEventPageTokenReofferStability(t *testing.T) {
	restore := eventTokenNowMs
	defer func() { eventTokenNowMs = restore }()

	eventTokenNowMs = func() int64 { return 1111111111111 }
	// Realistic clocks: the synthetic stamps below sit near the real
	// now so the decode-side expiry check evaluates the mint, not the
	// synthetic epoch.
	nowMs := time.Now().UnixMilli()
	eventTokenNowMs = func() int64 { return nowMs - 60000 }
	first, err := encodeEventPageToken(eventPageToken{Group: "g", Stream: "s", Direction: PageForward})
	if err != nil {
		t.Fatal(err)
	}

	// The clock moves; a scope carrying the walk's mint encodes stably.
	eventTokenNowMs = func() int64 { return nowMs }
	tok, err := decodeEventPageToken(first)
	if err != nil {
		t.Fatal(err)
	}
	again, err := encodeEventPageToken(tok)
	if err != nil {
		t.Fatal(err)
	}
	if again != first {
		t.Fatal("re-offered token differs from the presented token")
	}
	if tok.MintedMs != nowMs-60000 {
		t.Fatalf("walk mint = %d, want the first page's stamp", tok.MintedMs)
	}

	// A fresh walk (no carried mint) mints at the current clock.
	fresh, err := encodeEventPageToken(eventPageToken{Group: "g", Stream: "s", Direction: PageForward})
	if err != nil {
		t.Fatal(err)
	}
	freshTok, err := decodeEventPageToken(fresh)
	if err != nil {
		t.Fatal(err)
	}
	if freshTok.MintedMs != nowMs {
		t.Fatalf("fresh walk mint = %d, want the current clock", freshTok.MintedMs)
	}
}
