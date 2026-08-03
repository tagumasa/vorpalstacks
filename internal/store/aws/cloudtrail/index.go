package cloudtrail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
)

// indexCursorPrefix distinguishes encoded index cursors from plain scan-path
// markers in EventQuery.NextToken. The LookupEvents switch chooses exactly
// one path per call (determined by the query fields), so the token format
// never mixes within a single pagination sequence; the prefix is a safety
// net in case a caller accidentally passes the wrong token type.
const indexCursorPrefix = "idx1:"

// IndexCursor encodes a pagination position for index queries. It is opaque
// to callers — the store layer encodes/decodes it to/from the string
// NextToken.
//
// For single-bucket queries (Username, EventSource): Segment is empty and
// Key is the last scanned storage key within the bucket.
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

// decodeIndexCursor parses a string token back into an IndexCursor. Tokens
// that do not carry the index-cursor prefix (e.g. plain scan-path markers)
// yield an empty cursor, causing the query to start from the beginning.
func decodeIndexCursor(s string) (IndexCursor, error) {
	if s == "" || !strings.HasPrefix(s, indexCursorPrefix) {
		return IndexCursor{}, nil
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

// NewTimeIndexKey creates a new index key by time.
func NewTimeIndexKey(accountID, region, dateHour, eventID string) *EventIndexKey {
	return &EventIndexKey{
		IndexType: IndexByTime,
		AccountID: accountID,
		Region:    region,
		Segment1:  dateHour,
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
	dateHour := event.EventTime.Format("2006-01-02:15")
	ts := event.EventTime.UnixNano()

	var keys []*EventIndexKey
	keys = append(keys, NewTimeIndexKey(accountID, region, dateHour, event.EventID))

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

// RemoveIndex removes an event from the index.
func (m *EventIndexManager) RemoveIndex(event *Event) error {
	return m.applyIndexKeys(buildIndexKeys(m.accountID, m.region, event), m.deleteIndex)
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

func (m *EventIndexManager) deleteIndex(key *EventIndexKey) error {
	bucket := m.storage.Bucket(key.EncodePrefix())
	return bucket.Delete(indexStorageKey(key))
}

// QueryByTime queries events by time range, iterating hour-by-hour through
// the time index. The hour string format is "2006-01-02:15".
//
// Pagination across multiple hour buckets is supported via the cursor:
//   - cursor.Segment is the hour being partially scanned.
//   - cursor.Key is the last scanned key within that hour's bucket.
//
// On resume, hours before cursor.Segment are skipped, the matching hour is
// scanned from cursor.Key onwards, and subsequent hours are scanned from
// the beginning. When maxResults is reached, the cursor records the current
// hour and last key so the next call can resume exactly where this one
// stopped.
func (m *EventIndexManager) QueryByTime(startTime, endTime *time.Time, maxResults int32, cursor IndexCursor) ([]string, IndexCursor, error) {
	var hours []string
	if startTime != nil && endTime != nil {
		startHour := startTime.Truncate(time.Hour)
		endHour := endTime.Truncate(time.Hour)
		for t := startHour; !t.After(endHour); t = t.Add(time.Hour) {
			hours = append(hours, t.Format("2006-01-02:15"))
		}
	} else if startTime != nil {
		hours = append(hours, startTime.Truncate(time.Hour).Format("2006-01-02:15"))
	} else if endTime != nil {
		hours = append(hours, endTime.Truncate(time.Hour).Format("2006-01-02:15"))
	}

	var ids []string
	var nextCursor IndexCursor
	cursorActive := cursor.Segment != ""

	for _, dateHour := range hours {
		if cursorActive && dateHour != cursor.Segment {
			continue
		}

		remaining := maxResults - int32(len(ids))
		if remaining <= 0 {
			nextCursor = IndexCursor{Segment: dateHour}
			break
		}

		idxKey := &EventIndexKey{
			IndexType: IndexByTime,
			AccountID: m.accountID,
			Region:    m.region,
			Segment1:  dateHour,
		}

		var startAfterKey string
		if cursorActive {
			startAfterKey = cursor.Key
			cursorActive = false
		}

		hourIDs, lastKey := m.scanIndex(idxKey, remaining, startAfterKey)
		ids = append(ids, hourIDs...)

		if int32(len(hourIDs)) >= remaining && lastKey != "" {
			nextCursor = IndexCursor{Segment: dateHour, Key: lastKey}
			break
		}
	}

	return ids, nextCursor, nil
}

// QueryByEventName queries events by one or more event names, iterating
// through each event name's index bucket. Pagination across multiple
// buckets follows the same cursor pattern as QueryByTime:
//   - cursor.Segment is the event name being partially scanned.
//   - cursor.Key is the last scanned key within that bucket.
func (m *EventIndexManager) QueryByEventName(eventNames []string, maxResults int32, cursor IndexCursor) ([]string, IndexCursor, error) {
	var ids []string
	var nextCursor IndexCursor
	cursorActive := cursor.Segment != ""

	for _, eventName := range eventNames {
		if cursorActive && eventName != cursor.Segment {
			continue
		}

		remaining := maxResults - int32(len(ids))
		if remaining <= 0 {
			nextCursor = IndexCursor{Segment: eventName}
			break
		}

		idxKey := &EventIndexKey{
			IndexType: IndexByEventName,
			AccountID: m.accountID,
			Region:    m.region,
			Segment1:  eventName,
		}

		var startAfterKey string
		if cursorActive {
			startAfterKey = cursor.Key
			cursorActive = false
		}

		eventIDs, lastKey := m.scanIndex(idxKey, remaining, startAfterKey)
		ids = append(ids, eventIDs...)

		if int32(len(eventIDs)) >= remaining && lastKey != "" {
			nextCursor = IndexCursor{Segment: eventName, Key: lastKey}
			break
		}
	}

	return ids, nextCursor, nil
}

// QueryByUsername queries events by username. The username index is a
// single bucket, so the cursor only tracks the last scanned storage key.
// When the number of returned IDs reaches maxResults, a non-empty nextCursor
// is returned so the caller can fetch the next page.
func (m *EventIndexManager) QueryByUsername(username string, maxResults int32, cursor IndexCursor) ([]string, IndexCursor, error) {
	idxKey := &EventIndexKey{
		IndexType: IndexByUsername,
		AccountID: m.accountID,
		Region:    m.region,
		Segment1:  username,
	}
	ids, lastKey := m.scanIndex(idxKey, maxResults, cursor.Key)
	var nextCursor IndexCursor
	if int32(len(ids)) >= maxResults && lastKey != "" {
		nextCursor = IndexCursor{Key: lastKey}
	}
	return ids, nextCursor, nil
}

// QueryByEventSource queries events by event source. Same single-bucket
// pagination pattern as QueryByUsername.
func (m *EventIndexManager) QueryByEventSource(eventSource string, maxResults int32, cursor IndexCursor) ([]string, IndexCursor, error) {
	idxKey := &EventIndexKey{
		IndexType: IndexByEventSource,
		AccountID: m.accountID,
		Region:    m.region,
		Segment1:  eventSource,
	}
	ids, lastKey := m.scanIndex(idxKey, maxResults, cursor.Key)
	var nextCursor IndexCursor
	if int32(len(ids)) >= maxResults && lastKey != "" {
		nextCursor = IndexCursor{Key: lastKey}
	}
	return ids, nextCursor, nil
}

// scanIndex reads up to maxResults event IDs from the bucket addressed by
// key.EncodePrefix(). When startAfterKey is non-empty, iteration resumes
// strictly after that key (by appending \x00, the smallest byte value, to
// obtain the smallest key greater than startAfterKey in lexicographic order).
//
// The function returns the extracted event IDs and the raw storage key of
// the last item read. Callers use lastKey to build the next IndexCursor for
// subsequent pages.
func (m *EventIndexManager) scanIndex(key *EventIndexKey, maxResults int32, startAfterKey string) (ids []string, lastKey string) {
	if maxResults <= 0 {
		return nil, ""
	}
	prefix := key.EncodePrefix()
	bucket := m.storage.Bucket(prefix)

	var iter storage.Iterator
	if startAfterKey != "" {
		start := append([]byte(startAfterKey), 0x00)
		iter = bucket.ScanRange(start, nil)
	} else {
		iter = bucket.ScanPrefix(nil)
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
		return ids, ""
	}
	return ids, lastKey
}

// ClearIndexes clears all indexes for a given account and region.
func (m *EventIndexManager) ClearIndexes(accountID, region string) error {
	prefixes := []string{
		"ct_idx_time:" + accountID + ":" + region,
		"ct_idx_event:" + accountID + ":" + region,
		"ct_idx_user:" + accountID + ":" + region,
		"ct_idx_source:" + accountID + ":" + region,
	}

	for _, prefix := range prefixes {
		bucket := m.storage.Bucket(prefix)
		iter := bucket.ScanPrefix(nil)
		var keys [][]byte
		for iter.Next() {
			keys = append(keys, iter.Key())
		}
		if err := iter.Error(); err != nil {
			iter.Close()
			return err
		}
		iter.Close()

		for _, key := range keys {
			if err := bucket.Delete(key); err != nil {
				return err
			}
		}
	}

	return nil
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
