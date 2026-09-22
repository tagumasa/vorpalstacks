package cloudwatchlogs

import (
	"testing"
	"time"
)

// The documented token expiries: GetQueryResults' nextToken "expires
// after 1 hour", the describe listings' tokens "expire after 24 hours",
// and a scoped token repositions no foreign request.

func TestPageTokenExpiry(t *testing.T) {
	// Backdate the mint clock by a day: every listing-scoped vocabulary
	// rejects a token minted past its 24-hour expiry.
	restore := pageTokenNowMs
	pageTokenNowMs = func() int64 { return time.Now().Add(-25 * time.Hour).UnixMilli() }
	defer func() { pageTokenNowMs = restore }()

	stale := time.Now().Add(-25 * time.Hour).UnixMilli()
	if _, err := decodeQueryPageToken(mustEncodeQueryPageToken(t, stale)); err == nil {
		t.Fatal("stale DescribeQueries token accepted")
	}
	if _, err := decodeListingToken(mustEncodeListingToken(t, stale), "scope"); err == nil {
		t.Fatal("stale listing token accepted")
	}
	if _, err := decodeStreamPageToken(mustEncodeStreamPageToken(t, stale)); err == nil {
		t.Fatal("stale stream token accepted")
	}
	if _, err := decodeQueryResultsPageToken(mustEncodeResultsToken(t, stale)); err == nil {
		t.Fatal("stale results token accepted")
	}

	// The one-hour results TTL: a mint two hours old is expired while a
	// fresh mint pages.
	pageTokenNowMs = func() int64 { return time.Now().Add(-2 * time.Hour).UnixMilli() }
	if _, err := decodeQueryResultsPageToken(mustEncodeResultsToken(t, time.Now().Add(-2*time.Hour).UnixMilli())); err == nil {
		t.Fatal("two-hour-old results token accepted")
	}
	if _, err := decodeListingToken(mustEncodeListingToken(t, time.Now().Add(-2*time.Hour).UnixMilli()), "scope"); err != nil {
		t.Fatalf("two-hour-old listing token rejected: %v", err)
	}

	pageTokenNowMs = restore
	if _, err := decodeQueryResultsPageToken(mustEncodeResultsToken(t, time.Now().UnixMilli())); err != nil {
		t.Fatalf("fresh results token rejected: %v", err)
	}
}

// The results token is scoped to the minting query: a token replayed
// against another query rejects instead of repositioning that query's
// walk.
func TestQueryResultsTokenScoping(t *testing.T) {
	tok, err := encodeQueryResultsPageToken(queryResultsPageToken{QueryId: "query-a", Offset: 10})
	if err != nil {
		t.Fatal(err)
	}
	cursor, err := decodeQueryResultsPageToken(tok)
	if err != nil {
		t.Fatal(err)
	}
	if cursor.QueryId != "query-a" || cursor.Offset != 10 {
		t.Fatalf("cursor = %+v", cursor)
	}
}

// The version bytes are a namespace: every vocabulary's token rejects
// at every other vocabulary's decode. The JSON decode ignores unknown
// keys, so a foreign token structurally decodes into any vocabulary —
// without distinct versions it would silently reposition (or restart)
// that walk instead of rejecting.
func TestPageTokenVocabularyNamespace(t *testing.T) {
	streamTok, err := encodeStreamPageToken(streamPageToken{Group: "g"})
	if err != nil {
		t.Fatal(err)
	}
	queryTok, err := encodeQueryPageToken(queryPageToken{CursorCreateNs: 1})
	if err != nil {
		t.Fatal(err)
	}
	resultsTok, err := encodeQueryResultsPageToken(queryResultsPageToken{QueryId: "q", Offset: 1})
	if err != nil {
		t.Fatal(err)
	}
	listingTok, err := encodeListingToken("scope", "cursor")
	if err != nil {
		t.Fatal(err)
	}

	decoders := []struct {
		name   string
		decode func(string) error
	}{
		{"stream", func(s string) error { _, err := decodeStreamPageToken(s); return err }},
		{"query", func(s string) error { _, err := decodeQueryPageToken(s); return err }},
		{"results", func(s string) error { _, err := decodeQueryResultsPageToken(s); return err }},
		{"listing", func(s string) error { _, err := decodeListingToken(s, "scope"); return err }},
	}
	owned := map[string]string{
		"stream":  streamTok,
		"query":   queryTok,
		"results": resultsTok,
		"listing": listingTok,
	}
	for _, dec := range decoders {
		if err := dec.decode(owned[dec.name]); err != nil {
			t.Fatalf("%s vocabulary rejected its own token: %v", dec.name, err)
		}
		for vocab, tok := range owned {
			if vocab == dec.name {
				continue
			}
			if err := dec.decode(tok); err == nil {
				t.Fatalf("%s decode accepted a %s-vocabulary token", dec.name, vocab)
			}
		}
	}
}

func mustEncodeQueryPageToken(t *testing.T, mintedMs int64) string {
	t.Helper()
	tok, err := encodeQueryPageToken(queryPageToken{CursorCreateNs: 1, MintedMs: mintedMs})
	if err != nil {
		t.Fatal(err)
	}
	return tok
}

func mustEncodeListingToken(t *testing.T, mintedMs int64) string {
	t.Helper()
	tok, err := encodeListingToken("scope", "cursor")
	if err != nil {
		t.Fatal(err)
	}
	return tok
}

func mustEncodeStreamPageToken(t *testing.T, mintedMs int64) string {
	t.Helper()
	tok, err := encodeStreamPageToken(streamPageToken{Group: "g", MintedMs: mintedMs})
	if err != nil {
		t.Fatal(err)
	}
	return tok
}

func mustEncodeResultsToken(t *testing.T, mintedMs int64) string {
	t.Helper()
	tok, err := encodeQueryResultsPageToken(queryResultsPageToken{QueryId: "q", Offset: 1, MintedMs: mintedMs})
	if err != nil {
		t.Fatal(err)
	}
	return tok
}
