// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"time"

	"vorpalstacks/internal/common/pagination"
)

// Page-direction bytes carried inside event-page tokens.
const (
	PageForward  byte = 'f'
	PageBackward byte = 'b'
)

// eventPageTokenVersion identifies the event-token vocabulary; decoding a
// token whose version differs rejects it.
const eventPageTokenVersion byte = 4

// TokenExpiry24h is the documented expiry the pagination tokens of both
// planes share: "The token expires after 24 hours" (the read operations'
// nextToken/nextForwardToken/nextBackwardToken members and the describe
// listings' nextToken). The service-plane listing vocabulary references
// this constant.
const TokenExpiry24h = 24 * time.Hour

// eventTokenNowMs is the clock the event tokens mint from; tests backdate
// a mint by overriding it.
var eventTokenNowMs = func() int64 { return time.Now().UnixMilli() }

// eventTokenExpired reports whether a mint timestamp has passed the
// documented 24-hour expiry; a zero mint (a token from a vocabulary
// predating the field) is always expired.
func eventTokenExpired(mintedMs int64) bool {
	if mintedMs <= 0 {
		return true
	}
	return time.Now().UnixMilli()-mintedMs > TokenExpiry24h.Milliseconds()
}

// ErrInvalidPaginationToken is returned when a pagination token cannot be
// decoded or does not belong to the request that presented it. The
// service maps it to InvalidParameterException.
var ErrInvalidPaginationToken = errors.New("invalid pagination token")

// eventCursor is the value cursor of an event read: the sort key of the
// boundary event a page stopped at. Event reads order by (timestamp,
// ingestion time, log stream, message digest, ordinal). The documented
// sort key is (timestamp, ingestion time, ID of the PutLogEvents
// request); every ingested batch commits as exactly one chunk, so the
// ordinal — the chunk identifier plus the event's position inside the
// chunk — stands in for the per-put request ID and keeps byte-identical
// duplicate events distinct at a page boundary.
type eventCursor struct {
	Timestamp     int64  `json:"ts"`
	IngestionTime int64  `json:"it"`
	LogStream     string `json:"ls,omitempty"`
	MessageDigest string `json:"md,omitempty"`
	Ordinal       string `json:"od,omitempty"`
}

// messageDigest derives the cursor's compact tiebreak for a message.
func messageDigest(msg string) string {
	sum := sha256.Sum256([]byte(msg))
	return hex.EncodeToString(sum[:8])
}

// cursorBeforeZero is the cursor that precedes every event: a forward
// page resumed from it starts at the head of the scope.
func cursorBeforeZero() eventCursor {
	return eventCursor{}
}

// cursorAfterAll is the cursor that follows every event: a backward page
// resumed from it starts at the tail of the scope. Real events never
// reach the int64 maximum, so the timestamp component alone orders this
// cursor above the whole scope.
func cursorAfterAll() eventCursor {
	return eventCursor{
		Timestamp:     int64(^uint64(0) >> 1), // math.MaxInt64
		IngestionTime: int64(^uint64(0) >> 1),
	}
}

// eventPageToken scopes an events page to the request that produced it:
// group, stream (single-stream reads) or explicit stream set (filtered
// reads, identified by digest so the token stays compact), time window
// and filter pattern. A token presented against any other request shape
// fails validation rather than silently repositioning the read.
type eventPageToken struct {
	Version     byte        `json:"v"`
	Group       string      `json:"g"`
	Stream      string      `json:"s,omitempty"`
	StreamSetID string      `json:"ss,omitempty"`
	StartTime   int64       `json:"st"`
	EndTime     int64       `json:"et"`
	Pattern     string      `json:"fp,omitempty"`
	Direction   byte        `json:"d"`
	Cursor      eventCursor `json:"c"`
	MintedMs    int64       `json:"mi"`
}

// TokenVersion implements pagination.ScopedToken.
func (t eventPageToken) TokenVersion() byte { return t.Version }

// streamSetID digests the sorted explicit stream list of a filtered read
// into a compact scope identity.
func streamSetID(streams []string) string {
	if len(streams) == 0 {
		return ""
	}
	sorted := make([]string, len(streams))
	copy(sorted, streams)
	sort.Strings(sorted)
	sum := sha256.Sum256([]byte(strings.Join(sorted, "\x00")))
	return hex.EncodeToString(sum[:8])
}

// encodeEventPageToken serialises an event token to its wire form. The
// mint is stamped once per walk: a token that already carries one (the
// scope a presented token decoded into) keeps it, so every token of one
// walk — including the end-of-stream re-offer, which must equal the
// token the caller presented byte for byte — encodes identically, and
// the expiry ages from the walk's first page.
func encodeEventPageToken(t eventPageToken) (string, error) {
	t.Version = eventPageTokenVersion
	if t.MintedMs == 0 {
		t.MintedMs = eventTokenNowMs()
	}
	return pagination.EncodeScopedToken(t)
}

// decodeEventPageToken reverses encodeEventPageToken and rejects foreign
// or stale vocabularies — a token past its documented 24-hour expiry is
// as invalid as one from another request shape.
func decodeEventPageToken(s string) (eventPageToken, error) {
	var t eventPageToken
	if err := pagination.DecodeScopedToken(s, &t); err != nil {
		return t, ErrInvalidPaginationToken
	}
	if t.Version != eventPageTokenVersion || (t.Direction != PageForward && t.Direction != PageBackward) ||
		eventTokenExpired(t.MintedMs) {
		return t, ErrInvalidPaginationToken
	}
	return t, nil
}
