// Package logs provides logging functionality for vorpalstacks.
package logs

import (
	"time"

	"github.com/google/uuid"
	"vorpalstacks/internal/core/storage"
)

// LogStore manages log entries and provides query functionality.
type LogStore struct {
	storage storage.MaintenanceStorage
	indexer *IndexManager
}

// NewLogStore creates a new LogStore with the given storage backend.
func NewLogStore(s storage.MaintenanceStorage) *LogStore {
	return &LogStore{
		storage: s,
		indexer: NewIndexManager(s),
	}
}

// Write stores a log entry.
func (s *LogStore) Write(entry *LogEntry) error {
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	tenantID := entry.TenantID
	if tenantID == "" {
		tenantID = "default"
	}
	region := entry.Region
	if region == "" {
		region = "local"
	}

	key := NewLogKey(tenantID, region, entry.ID)
	data, err := entry.Encode()
	if err != nil {
		return err
	}

	bucket := s.storage.Bucket(key.EncodePrefix())
	if err := bucket.Put([]byte(entry.ID), data); err != nil {
		return err
	}

	if s.indexer != nil {
		if err := s.indexer.AddIndex(entry); err != nil {
			return err
		}
	}

	return nil
}

// Get retrieves a log entry by its ID.
func (s *LogStore) Get(tenantID, region, id string) (*LogEntry, error) {
	key := NewLogKey(tenantID, region, id)
	bucket := s.storage.Bucket(key.EncodePrefix())
	data, err := bucket.Get([]byte(id))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return DecodeLogEntry(data)
}

// Query retrieves log entries matching the given filter.
func (s *LogStore) Query(filter QueryFilter) ([]*LogEntry, error) {
	var entryIDs []string
	var err error

	if filter.StartTime != nil || filter.EndTime != nil {
		entryIDs, err = s.indexer.QueryByTime(filter)
	} else if filter.Level != nil {
		entryIDs, err = s.indexer.QueryByLevel(filter)
	} else if filter.Source != nil {
		entryIDs, err = s.indexer.QueryBySource(filter)
	} else {
		entryIDs, err = s.listAllEntries(filter)
	}

	if err != nil {
		return nil, err
	}

	if filter.Limit > 0 && len(entryIDs) > filter.Limit {
		entryIDs = entryIDs[:filter.Limit]
	}

	entries := make([]*LogEntry, 0, len(entryIDs))
	tenantID := "default"
	if filter.TenantID != nil {
		tenantID = *filter.TenantID
	}
	region := "local"
	if filter.Region != nil {
		region = *filter.Region
	}
	for _, id := range entryIDs {
		entry, err := s.Get(tenantID, region, id)
		if err != nil {
			continue
		}
		if entry != nil {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func (s *LogStore) listAllEntries(filter QueryFilter) ([]string, error) {
	tenantID := "default"
	if filter.TenantID != nil {
		tenantID = *filter.TenantID
	}
	region := "local"
	if filter.Region != nil {
		region = *filter.Region
	}

	key := NewLogKey(tenantID, region, "")
	bucket := s.storage.Bucket(key.EncodePrefix())

	var ids []string
	iter := bucket.ScanPrefix(nil)
	defer iter.Close()

	for iter.Next() {
		ids = append(ids, string(iter.Key()))
	}

	return ids, nil
}

// DeleteBefore removes all log entries before the specified time.
func (s *LogStore) DeleteBefore(t time.Time) error {
	tenantID := "default"
	region := "local"

	key := NewLogKey(tenantID, region, "")
	bucket := s.storage.Bucket(key.EncodePrefix())

	var toDelete []*LogEntry
	iter := bucket.ScanPrefix(nil)
	for iter.Next() {
		entry, err := DecodeLogEntry(iter.Value())
		if err != nil {
			continue
		}
		if entry.Timestamp.Before(t) {
			toDelete = append(toDelete, entry)
		}
	}
	iter.Close()

	for _, entry := range toDelete {
		if err := bucket.Delete([]byte(entry.ID)); err != nil {
			return err
		}
		if s.indexer != nil {
			if err := s.indexer.RemoveIndex(entry); err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteByID removes a log entry by its ID.
func (s *LogStore) DeleteByID(tenantID, region, id string) error {
	entry, err := s.Get(tenantID, region, id)
	if err != nil {
		return err
	}
	if entry == nil {
		return nil
	}

	key := NewLogKey(tenantID, region, id)
	bucket := s.storage.Bucket(key.EncodePrefix())
	if err := bucket.Delete([]byte(id)); err != nil {
		return err
	}

	if s.indexer != nil {
		return s.indexer.RemoveIndex(entry)
	}

	return nil
}

// Count returns the number of log entries matching the filter.
func (s *LogStore) Count(filter QueryFilter) (int, error) {
	entries, err := s.Query(filter)
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}

// Compact compacts the underlying storage.
func (s *LogStore) Compact() error {
	return s.storage.Compact()
}

// Close closes the LogStore and releases resources.
func (s *LogStore) Close() error {
	return s.storage.Close()
}
