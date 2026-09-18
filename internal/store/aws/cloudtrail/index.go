package cloudtrail

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
)

// indexCursorPrefix distinguishes encoded index cursors from any other
// string in EventQuery.NextToken. Every LookupEvents response token is an
// encoded index cursor — the scan path encodes its continuation marker the
// same way — so a non-empty NextToken that does not carry the prefix (or
// fails to decode) is rejected as invalid rather than silently restarting
// the query from the beginning.
const indexCursorPrefix = "idx1:"

// IndexCursor encodes a pagination position for index queries. It is opaque
// to callers — the store layer encodes/decodes it to/from the string
// NextToken. Retrieval walks newest first, so the position is a LOWER
// bound: iteration resumes strictly below Key.
//
// For single-bucket queries (Username, EventSource, and the filterless
// scan path): Segment is empty and Key is the last scanned storage key
// within the bucket.
//
// For multi-bucket queries (Time, EventName): Segment identifies the
// current segment being iterated (an hour string like "2024-02-25:10" for
// Time, or the event name for EventName). Key is the last scanned storage
// key within that segment's bucket.
type IndexCursor struct {
	Segment string `json:"s,omitempty"`
	Key     string `json:"k,omitempty"`
}

// encodeIndexCursor serialises an IndexCursor into an opaque string token.
// Returns an empty string when the cursor is exhausted (both fields blank),
// so callers can treat a non-empty token as "more results may exist".
func encodeIndexCursor(c IndexCursor) string {
	if c.Segment == "" && c.Key == "" {
		return ""
	}
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return indexCursorPrefix + base64.StdEncoding.EncodeToString(b)
}

// decodeIndexCursor parses a string token back into an IndexCursor. The
// empty token yields an empty cursor (start from the beginning); any other
// token that does not carry the index-cursor prefix, or whose payload does
// not decode, is an error — the caller surfaces it as an invalid token
// instead of silently restarting the walk.
func decodeIndexCursor(s string) (IndexCursor, error) {
	if s == "" {
		return IndexCursor{}, nil
	}
	if !strings.HasPrefix(s, indexCursorPrefix) {
		return IndexCursor{}, fmt.Errorf("token %q does not carry the index-cursor prefix", s)
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(s, indexCursorPrefix))
	if err != nil {
		return IndexCursor{}, err
	}
	var c IndexCursor
	if err := json.Unmarshal(data, &c); err != nil {
		return IndexCursor{}, err
	}
	return c, nil
}

// IndexType represents the type of index for CloudTrail events.
type IndexType int

const (
	// IndexByTime indexes events by time.
	IndexByTime IndexType = iota
	// IndexByEventName indexes events by event name.
	IndexByEventName
	// IndexByUsername indexes events by username.
	IndexByUsername
	// IndexByEventSource indexes events by event source.
	IndexByEventSource
)

// EventIndexKey represents a key for indexing CloudTrail events.
type EventIndexKey struct {
	IndexType IndexType
	AccountID string
	Region    string
	Segment1  string
	Segment2  string
	EventID   string
}

// EncodePrefix encodes the index key as a bucket prefix string.
// Segment2 (timestamp) is intentionally excluded from the bucket name so that
// queries using only Segment1 resolve to the same bucket as writes.
func (k *EventIndexKey) EncodePrefix() string {
	var prefix string
	switch k.IndexType {
	case IndexByTime:
		prefix = "ct_idx_time"
	case IndexByEventName:
		prefix = "ct_idx_event"
	case IndexByUsername:
		prefix = "ct_idx_user"
	case IndexByEventSource:
		prefix = "ct_idx_source"
	}
	if k.Segment1 != "" {
		return fmt.Sprintf("%s:%s:%s:%s", prefix, k.AccountID, k.Region, k.Segment1)
	}
	return fmt.Sprintf("%s:%s:%s", prefix, k.AccountID, k.Region)
}

// NewTimeIndexKey creates a new index key by time. The timestamp keeps the
// entries within an hour bucket in time order, which the descending
// retrieval walk (most recent event first, per the LookupEvents contract)
// depends on.
func NewTimeIndexKey(accountID, region, dateHour string, timestamp int64, eventID string) *EventIndexKey {
	return &EventIndexKey{
		IndexType: IndexByTime,
		AccountID: accountID,
		Region:    region,
		Segment1:  dateHour,
		Segment2:  fmt.Sprintf("%d", timestamp),
		EventID:   eventID,
	}
}

// NewEventNameIndexKey creates a new index key by event name.
func NewEventNameIndexKey(accountID, region, eventName string, timestamp int64, eventID string) *EventIndexKey {
	return &EventIndexKey{
		IndexType: IndexByEventName,
		AccountID: accountID,
		Region:    region,
		Segment1:  eventName,
		Segment2:  fmt.Sprintf("%d", timestamp),
		EventID:   eventID,
	}
}

// NewUsernameIndexKey creates a new index key by username.
func NewUsernameIndexKey(accountID, region, username string, timestamp int64, eventID string) *EventIndexKey {
	return &EventIndexKey{
		IndexType: IndexByUsername,
		AccountID: accountID,
		Region:    region,
		Segment1:  username,
		Segment2:  fmt.Sprintf("%d", timestamp),
		EventID:   eventID,
	}
}

// NewEventSourceIndexKey creates a new index key by event source.
func NewEventSourceIndexKey(accountID, region, eventSource string, timestamp int64, eventID string) *EventIndexKey {
	return &EventIndexKey{
		IndexType: IndexByEventSource,
		AccountID: accountID,
		Region:    region,
		Segment1:  eventSource,
		Segment2:  fmt.Sprintf("%d", timestamp),
		EventID:   eventID,
	}
}

// EventIndexManager manages indexes for CloudTrail events.
type EventIndexManager struct {
	storage   storage.BasicStorage
	accountID string
	region    string
}

// NewEventIndexManager creates a new EventIndexManager instance.
func NewEventIndexManager(s storage.BasicStorage, accountID, region string) *EventIndexManager {
	return &EventIndexManager{
		storage:   s,
		accountID: accountID,
		region:    region,
	}
}

// buildIndexKeys generates all index keys for the given event.
func buildIndexKeys(accountID, region string, event *Event) []*EventIndexKey {
	// The bucket hour is derived from the UTC-normalised event time: the
	// read side addresses buckets from UTC bounds, so a located event time
	// would address the wrong wall-hour bucket.
	dateHour := event.EventTime.UTC().Format("2006-01-02:15")
	ts := event.EventTime.UnixNano()

	var keys []*EventIndexKey
	keys = append(keys, NewTimeIndexKey(accountID, region, dateHour, ts, event.EventID))

	if event.EventName != "" {
		keys = append(keys, NewEventNameIndexKey(accountID, region, event.EventName, ts, event.EventID))
	}
	if event.UserIdentity != nil && event.UserIdentity.UserName != "" {
		keys = append(keys, NewUsernameIndexKey(accountID, region, event.UserIdentity.UserName, ts, event.EventID))
	}
	if event.EventSource != "" {
		keys = append(keys, NewEventSourceIndexKey(accountID, region, event.EventSource, ts, event.EventID))
	}
	return keys
}

// AddIndex adds an event to the index.
func (m *EventIndexManager) AddIndex(event *Event) error {
	return m.applyIndexKeys(buildIndexKeys(m.accountID, m.region, event), m.putIndex)
}

func (m *EventIndexManager) applyIndexKeys(keys []*EventIndexKey, fn func(*EventIndexKey) error) error {
	for _, k := range keys {
		if err := fn(k); err != nil {
			return err
		}
	}
	return nil
}

func (m *EventIndexManager) putIndex(key *EventIndexKey) error {
	bucket := m.storage.Bucket(key.EncodePrefix())
	return bucket.Put(indexStorageKey(key), []byte{1})
}

// QueryByTime queries events by time range, iterating the hour buckets of
// the time index NEWEST HOUR FIRST, and within each hour newest event
// first — the LookupEvents contract ("The events list is sorted by time.
// The most recent event is listed first."). The hour string format is
// "2006-01-02:15"; bounds are normalised to UTC before the hours are
// derived, because the buckets are addressed by UTC wall hours.
//
// Both bounds must be non-nil: an unbounded side is resolved against the
// recorded event span by the caller (store.LookupEvents) before this
// method runs, so a single-bound lookup covers every hour from its bound
// to the newest (or from the oldest to its bound) recorded event instead
// of a single hour.
//
// Pagination across multiple hour buckets is supported via the cursor:
//   - cursor.Segment is the hour being partially scanned.
//   - cursor.Key is the last scanned key within that hour's bucket.
//
// On resume, hours newer than cursor.Segment are skipped (already served),
// the matching hour is scanned from cursor.Key downwards, and older hours
// are scanned from their newest key. When maxResults is reached, the
// cursor records the current hour and last scanned key so the next call
// can resume exactly where this one stopped.
func (m *EventIndexManager) QueryByTime(startTime, endTime *time.Time, maxResults int32, cursor IndexCursor) ([]string, IndexCursor, error) {
	if startTime == nil || endTime == nil {
		return nil, IndexCursor{}, errors.New("cloudtrail index: QueryByTime requires resolved bounds")
	}
	startHour := startTime.UTC().Truncate(time.Hour)
	endHour := endTime.UTC().Truncate(time.Hour)

	var hours []string
	for t := endHour; !t.Before(startHour); t = t.Add(-time.Hour) {
		hours = append(hours, t.Format("2006-01-02:15"))
	}

	return m.querySegments(IndexByTime, hours, maxResults, cursor)
}

// QueryByEventName queries events by one or more event names, iterating
// through each event name's index bucket. Pagination across multiple
// buckets follows the same cursor pattern as QueryByTime:
//   - cursor.Segment is the event name being partially scanned.
//   - cursor.Key is the last scanned key within that bucket.
func (m *EventIndexManager) QueryByEventName(eventNames []string, maxResults int32, cursor IndexCursor) ([]string, IndexCursor, error) {
	return m.querySegments(IndexByEventName, eventNames, maxResults, cursor)
}

// QueryByUsername queries events by username. The username index is a
// single bucket, so the cursor only tracks the last scanned storage key.
// When the number of returned IDs reaches maxResults, a non-empty nextCursor
// is returned so the caller can fetch the next page.
func (m *EventIndexManager) QueryByUsername(username string, maxResults int32, cursor IndexCursor) ([]string, IndexCursor, error) {
	return m.querySingleBucket(IndexByUsername, username, maxResults, cursor)
}

// QueryByEventSource queries events by event source. Same single-bucket
// pagination pattern as QueryByUsername.
func (m *EventIndexManager) QueryByEventSource(eventSource string, maxResults int32, cursor IndexCursor) ([]string, IndexCursor, error) {
	return m.querySingleBucket(IndexByEventSource, eventSource, maxResults, cursor)
}

// querySegments walks the multi-bucket cursor pagination shared by the Time
// and EventName indexes, NEWEST FIRST within every segment bucket: segments
// before the cursor's segment (in the caller's segment order — newest hour
// first for Time) are skipped as already served, the cursor's segment
// resumes strictly below cursor.Key, and each segment's bucket is scanned
// downwards up to the remaining result budget. When the budget is reached
// mid-walk the returned cursor carries the current segment and the last
// scanned key so the next call resumes exactly where this one stopped.
func (m *EventIndexManager) querySegments(indexType IndexType, segments []string, maxResults int32, cursor IndexCursor) ([]string, IndexCursor, error) {
	var ids []string
	var nextCursor IndexCursor
	cursorActive := cursor.Segment != ""

	for _, segment := range segments {
		if cursorActive && segment != cursor.Segment {
			continue
		}

		remaining := maxResults - int32(len(ids))
		if remaining <= 0 {
			nextCursor = IndexCursor{Segment: segment}
			break
		}

		idxKey := &EventIndexKey{
			IndexType: indexType,
			AccountID: m.accountID,
			Region:    m.region,
			Segment1:  segment,
		}

		var beforeKey string
		if cursorActive {
			beforeKey = cursor.Key
			cursorActive = false
		}

		segmentIDs, lastKey, err := m.scanIndexReverse(idxKey, remaining, beforeKey)
		if err != nil {
			return nil, IndexCursor{}, err
		}
		ids = append(ids, segmentIDs...)

		if int32(len(segmentIDs)) >= remaining && lastKey != "" {
			nextCursor = IndexCursor{Segment: segment, Key: lastKey}
			break
		}
	}

	return ids, nextCursor, nil
}

// querySingleBucket scans one index bucket addressed by the segment, the
// cursor pagination shared by the Username and EventSource indexes: a
// single bucket means the cursor only tracks the last scanned storage key.
// The scan runs newest first, matching the multi-bucket walks.
func (m *EventIndexManager) querySingleBucket(indexType IndexType, segment string, maxResults int32, cursor IndexCursor) ([]string, IndexCursor, error) {
	idxKey := &EventIndexKey{
		IndexType: indexType,
		AccountID: m.accountID,
		Region:    m.region,
		Segment1:  segment,
	}
	ids, lastKey, err := m.scanIndexReverse(idxKey, maxResults, cursor.Key)
	if err != nil {
		return nil, IndexCursor{}, err
	}
	var nextCursor IndexCursor
	if int32(len(ids)) >= maxResults && lastKey != "" {
		nextCursor = IndexCursor{Key: lastKey}
	}
	return ids, nextCursor, nil
}

// scanIndexReverse reads up to maxResults event IDs from the bucket
// addressed by key.EncodePrefix(), walking the storage keys DOWNWARDS
// (newest first — every index bucket keys its entries
// "<timestamp>:<eventID>"). When before is non-empty, iteration starts at
// the largest key strictly less than before, resuming a previous page.
//
// The function returns the extracted event IDs and the raw storage key of
// the last item read. Callers use lastKey to build the next IndexCursor for
// subsequent pages. An iteration failure propagates — a partial page must
// never masquerade as a complete result.
func (m *EventIndexManager) scanIndexReverse(key *EventIndexKey, maxResults int32, before string) (ids []string, lastKey string, err error) {
	if maxResults <= 0 {
		return nil, "", nil
	}
	prefix := key.EncodePrefix()
	bucket := m.storage.Bucket(prefix)

	var iter storage.Iterator
	if before != "" {
		iter = bucket.ScanPrefixReverse(nil, []byte(before))
	} else {
		iter = bucket.ScanPrefixReverse(nil, nil)
	}
	defer iter.Close()

	for iter.Next() {
		if int32(len(ids)) >= maxResults {
			break
		}
		k := string(iter.Key())
		lastKey = k
		if idx := strings.LastIndex(k, ":"); idx >= 0 {
			ids = append(ids, k[idx+1:])
		} else {
			ids = append(ids, k)
		}
	}
	if err := iter.Error(); err != nil {
		return nil, "", err
	}
	return ids, lastKey, nil
}

// AddIndexInTxn adds an event to the index within a transaction.
func (m *EventIndexManager) AddIndexInTxn(txn storage.Transaction, event *Event) error {
	return m.applyIndexKeys(buildIndexKeys(m.accountID, m.region, event), func(key *EventIndexKey) error {
		return m.putIndexInTxn(txn, key)
	})
}

func (m *EventIndexManager) putIndexInTxn(txn storage.Transaction, key *EventIndexKey) error {
	bucket := txn.Bucket(key.EncodePrefix())
	return bucket.Put(indexStorageKey(key), []byte{1})
}

// indexStorageKey returns the storage key for an index entry.
// When Segment2 (timestamp) is present, it is prepended to the EventID
// to maintain time-ordered iteration within the bucket.
func indexStorageKey(key *EventIndexKey) []byte {
	if key.Segment2 != "" {
		return []byte(key.Segment2 + ":" + key.EventID)
	}
	return []byte(key.EventID)
}
