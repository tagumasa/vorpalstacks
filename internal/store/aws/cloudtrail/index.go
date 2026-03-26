package cloudtrail

import (
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
)

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

// EncodePrefix encodes the index key as a prefix string.
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
	if k.Segment2 != "" {
		return fmt.Sprintf("%s:%s:%s:%s:%s", prefix, k.AccountID, k.Region, k.Segment1, k.Segment2)
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

// AddIndex adds an event to the index.
func (m *EventIndexManager) AddIndex(event *Event) error {
	accountID := m.accountID
	region := m.region
	dateHour := event.EventTime.Format("2006-01-02:15")
	ts := event.EventTime.UnixNano()

	timeKey := NewTimeIndexKey(accountID, region, dateHour, event.EventID)
	if err := m.putIndex(timeKey); err != nil {
		return err
	}

	if event.EventName != "" {
		eventNameKey := NewEventNameIndexKey(accountID, region, event.EventName, ts, event.EventID)
		if err := m.putIndex(eventNameKey); err != nil {
			return err
		}
	}

	if event.UserIdentity != nil && event.UserIdentity.UserName != "" {
		usernameKey := NewUsernameIndexKey(accountID, region, event.UserIdentity.UserName, ts, event.EventID)
		if err := m.putIndex(usernameKey); err != nil {
			return err
		}
	}

	if event.EventSource != "" {
		eventSourceKey := NewEventSourceIndexKey(accountID, region, event.EventSource, ts, event.EventID)
		if err := m.putIndex(eventSourceKey); err != nil {
			return err
		}
	}

	return nil
}

// RemoveIndex removes an event from the index.
func (m *EventIndexManager) RemoveIndex(event *Event) error {
	accountID := m.accountID
	region := m.region
	dateHour := event.EventTime.Format("2006-01-02:15")
	ts := event.EventTime.UnixNano()

	timeKey := NewTimeIndexKey(accountID, region, dateHour, event.EventID)
	if err := m.deleteIndex(timeKey); err != nil {
		return err
	}

	if event.EventName != "" {
		eventNameKey := NewEventNameIndexKey(accountID, region, event.EventName, ts, event.EventID)
		if err := m.deleteIndex(eventNameKey); err != nil {
			return err
		}
	}

	if event.UserIdentity != nil && event.UserIdentity.UserName != "" {
		usernameKey := NewUsernameIndexKey(accountID, region, event.UserIdentity.UserName, ts, event.EventID)
		if err := m.deleteIndex(usernameKey); err != nil {
			return err
		}
	}

	if event.EventSource != "" {
		eventSourceKey := NewEventSourceIndexKey(accountID, region, event.EventSource, ts, event.EventID)
		if err := m.deleteIndex(eventSourceKey); err != nil {
			return err
		}
	}

	return nil
}

func (m *EventIndexManager) putIndex(key *EventIndexKey) error {
	bucket := m.storage.Bucket(key.EncodePrefix())
	return bucket.Put([]byte(key.EventID), []byte{1})
}

func (m *EventIndexManager) deleteIndex(key *EventIndexKey) error {
	bucket := m.storage.Bucket(key.EncodePrefix())
	return bucket.Delete([]byte(key.EventID))
}

// QueryByTime queries events by time range.
func (m *EventIndexManager) QueryByTime(startTime, endTime *time.Time, maxResults int32) ([]string, error) {
	accountID := m.accountID
	region := m.region

	var ids []string

	if startTime != nil && endTime != nil {
		startHour := startTime.Truncate(time.Hour)
		endHour := endTime.Truncate(time.Hour)

		for t := startHour; !t.After(endHour); t = t.Add(time.Hour) {
			dateHour := t.Format("2006-01-02:15")
			idxKey := &EventIndexKey{
				IndexType: IndexByTime,
				AccountID: accountID,
				Region:    region,
				Segment1:  dateHour,
			}
			hourIDs := m.scanIndex(idxKey, maxResults-int32(len(ids)))
			ids = append(ids, hourIDs...)
			if int32(len(ids)) >= maxResults {
				break
			}
		}
	} else if startTime != nil {
		dateHour := startTime.Truncate(time.Hour).Format("2006-01-02:15")
		idxKey := &EventIndexKey{
			IndexType: IndexByTime,
			AccountID: accountID,
			Region:    region,
			Segment1:  dateHour,
		}
		ids = m.scanIndex(idxKey, maxResults)
	} else if endTime != nil {
		dateHour := endTime.Truncate(time.Hour).Format("2006-01-02:15")
		idxKey := &EventIndexKey{
			IndexType: IndexByTime,
			AccountID: accountID,
			Region:    region,
			Segment1:  dateHour,
		}
		ids = m.scanIndex(idxKey, maxResults)
	}

	return ids, nil
}

// QueryByEventName queries events by event name.
func (m *EventIndexManager) QueryByEventName(eventNames []string, maxResults int32) ([]string, error) {
	accountID := m.accountID
	region := m.region

	var ids []string
	for _, eventName := range eventNames {
		idxKey := &EventIndexKey{
			IndexType: IndexByEventName,
			AccountID: accountID,
			Region:    region,
			Segment1:  eventName,
		}
		eventIDs := m.scanIndex(idxKey, maxResults-int32(len(ids)))
		ids = append(ids, eventIDs...)
		if int32(len(ids)) >= maxResults {
			break
		}
	}

	return ids, nil
}

// QueryByUsername queries events by username.
func (m *EventIndexManager) QueryByUsername(username string, maxResults int32) ([]string, error) {
	accountID := m.accountID
	region := m.region

	idxKey := &EventIndexKey{
		IndexType: IndexByUsername,
		AccountID: accountID,
		Region:    region,
		Segment1:  username,
	}

	return m.scanIndex(idxKey, maxResults), nil
}

// QueryByEventSource queries events by event source.
func (m *EventIndexManager) QueryByEventSource(eventSource string, maxResults int32) ([]string, error) {
	accountID := m.accountID
	region := m.region

	idxKey := &EventIndexKey{
		IndexType: IndexByEventSource,
		AccountID: accountID,
		Region:    region,
		Segment1:  eventSource,
	}

	return m.scanIndex(idxKey, maxResults), nil
}

func (m *EventIndexManager) scanIndex(key *EventIndexKey, maxResults int32) []string {
	if maxResults <= 0 {
		return nil
	}
	prefix := key.EncodePrefix()
	bucket := m.storage.Bucket(prefix)

	var ids []string
	iter := bucket.ScanPrefix(nil)
	defer iter.Close()

	for iter.Next() {
		if int32(len(ids)) >= maxResults {
			break
		}
		k := string(iter.Key())
		if idx := strings.LastIndex(k, ":"); idx >= 0 {
			ids = append(ids, k[idx+1:])
		} else {
			ids = append(ids, k)
		}
	}
	if err := iter.Error(); err != nil {
		return ids
	}

	return ids
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
	accountID := m.accountID
	region := m.region
	dateHour := event.EventTime.Format("2006-01-02:15")
	ts := event.EventTime.UnixNano()

	timeKey := NewTimeIndexKey(accountID, region, dateHour, event.EventID)
	if err := m.putIndexInTxn(txn, timeKey); err != nil {
		return err
	}

	if event.EventName != "" {
		eventNameKey := NewEventNameIndexKey(accountID, region, event.EventName, ts, event.EventID)
		if err := m.putIndexInTxn(txn, eventNameKey); err != nil {
			return err
		}
	}

	if event.UserIdentity != nil && event.UserIdentity.UserName != "" {
		usernameKey := NewUsernameIndexKey(accountID, region, event.UserIdentity.UserName, ts, event.EventID)
		if err := m.putIndexInTxn(txn, usernameKey); err != nil {
			return err
		}
	}

	if event.EventSource != "" {
		eventSourceKey := NewEventSourceIndexKey(accountID, region, event.EventSource, ts, event.EventID)
		if err := m.putIndexInTxn(txn, eventSourceKey); err != nil {
			return err
		}
	}

	return nil
}

func (m *EventIndexManager) putIndexInTxn(txn storage.Transaction, key *EventIndexKey) error {
	bucket := txn.Bucket(key.EncodePrefix())
	return bucket.Put([]byte(key.EventID), []byte{1})
}
