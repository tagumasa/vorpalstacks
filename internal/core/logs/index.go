// Package logs provides logging functionality for vorpalstacks.
package logs

import (
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
)

// IndexManager manages indexes for efficient log entry queries.
type IndexManager struct {
	storage storage.BasicStorage
}

// NewIndexManager creates a new IndexManager with the given storage.
func NewIndexManager(s storage.BasicStorage) *IndexManager {
	return &IndexManager{storage: s}
}

// AddIndex adds index entries for a log entry.
func (m *IndexManager) AddIndex(entry *LogEntry) error {
	tenantID := entry.TenantID
	if tenantID == "" {
		tenantID = "default"
	}
	region := entry.Region
	if region == "" {
		region = "local"
	}

	dateHour := entry.Timestamp.Format("2006-01-02:15")
	ts := entry.Timestamp.UnixNano()

	timeKey := NewTimeIndexKey(tenantID, region, dateHour, entry.ID)
	if err := m.putIndex(timeKey); err != nil {
		return err
	}

	levelKey := NewLevelIndexKey(tenantID, region, entry.Level.String(), ts, entry.ID)
	if err := m.putIndex(levelKey); err != nil {
		return err
	}

	if entry.Source != "" {
		sourceKey := NewSourceIndexKey(tenantID, region, entry.Source, ts, entry.ID)
		if err := m.putIndex(sourceKey); err != nil {
			return err
		}
	}

	return nil
}

func (m *IndexManager) putIndex(key *IndexKey) error {
	bucket := m.storage.Bucket(key.EncodePrefix())
	return bucket.Put([]byte(key.EntryID), []byte{1})
}

// RemoveIndex removes index entries for a log entry.
func (m *IndexManager) RemoveIndex(entry *LogEntry) error {
	tenantID := entry.TenantID
	if tenantID == "" {
		tenantID = "default"
	}
	region := entry.Region
	if region == "" {
		region = "local"
	}

	dateHour := entry.Timestamp.Format("2006-01-02:15")
	ts := entry.Timestamp.UnixNano()

	timeKey := NewTimeIndexKey(tenantID, region, dateHour, entry.ID)
	if err := m.deleteIndex(timeKey); err != nil {
		return err
	}

	levelKey := NewLevelIndexKey(tenantID, region, entry.Level.String(), ts, entry.ID)
	if err := m.deleteIndex(levelKey); err != nil {
		return err
	}

	if entry.Source != "" {
		sourceKey := NewSourceIndexKey(tenantID, region, entry.Source, ts, entry.ID)
		if err := m.deleteIndex(sourceKey); err != nil {
			return err
		}
	}

	return nil
}

func (m *IndexManager) deleteIndex(key *IndexKey) error {
	bucket := m.storage.Bucket(key.EncodePrefix())
	return bucket.Delete([]byte(key.EntryID))
}

// QueryByTime retrieves entry IDs matching the time filter.
func (m *IndexManager) QueryByTime(filter QueryFilter) ([]string, error) {
	tenantID := "default"
	if filter.TenantID != nil {
		tenantID = *filter.TenantID
	}
	region := "local"
	if filter.Region != nil {
		region = *filter.Region
	}

	var ids []string

	if filter.StartTime != nil && filter.EndTime != nil {
		startHour := filter.StartTime.Truncate(time.Hour)
		endHour := filter.EndTime.Truncate(time.Hour)

		for t := startHour; !t.After(endHour); t = t.Add(time.Hour) {
			dateHour := t.Format("2006-01-02:15")
			idxKey := &IndexKey{
				IndexType: IndexByTime,
				TenantID:  tenantID,
				Region:    region,
				Segment1:  dateHour,
			}
			hourIDs := m.scanIndex(idxKey)
			ids = append(ids, hourIDs...)
		}
	} else if filter.StartTime != nil {
		dateHour := filter.StartTime.Truncate(time.Hour).Format("2006-01-02:15")
		idxKey := &IndexKey{
			IndexType: IndexByTime,
			TenantID:  tenantID,
			Region:    region,
			Segment1:  dateHour,
		}
		ids = m.scanIndex(idxKey)
	} else if filter.EndTime != nil {
		dateHour := filter.EndTime.Truncate(time.Hour).Format("2006-01-02:15")
		idxKey := &IndexKey{
			IndexType: IndexByTime,
			TenantID:  tenantID,
			Region:    region,
			Segment1:  dateHour,
		}
		ids = m.scanIndex(idxKey)
	}

	return ids, nil
}

// QueryByLevel retrieves entry IDs matching the level filter.
func (m *IndexManager) QueryByLevel(filter QueryFilter) ([]string, error) {
	tenantID := "default"
	if filter.TenantID != nil {
		tenantID = *filter.TenantID
	}
	region := "local"
	if filter.Region != nil {
		region = *filter.Region
	}

	idxKey := &IndexKey{
		IndexType: IndexByLevel,
		TenantID:  tenantID,
		Region:    region,
		Segment1:  filter.Level.String(),
	}

	return m.scanIndex(idxKey), nil
}

// QueryBySource retrieves entry IDs matching the source filter.
func (m *IndexManager) QueryBySource(filter QueryFilter) ([]string, error) {
	tenantID := "default"
	if filter.TenantID != nil {
		tenantID = *filter.TenantID
	}
	region := "local"
	if filter.Region != nil {
		region = *filter.Region
	}

	if filter.Source == nil {
		return nil, nil
	}

	idxKey := &IndexKey{
		IndexType: IndexBySource,
		TenantID:  tenantID,
		Region:    region,
		Segment1:  *filter.Source,
	}

	return m.scanIndex(idxKey), nil
}

func (m *IndexManager) scanIndex(key *IndexKey) []string {
	prefix := key.EncodePrefix()
	bucket := m.storage.Bucket(prefix)

	var ids []string
	iter := bucket.ScanPrefix(nil)
	defer iter.Close()

	for iter.Next() {
		k := string(iter.Key())
		if idx := strings.LastIndex(k, ":"); idx >= 0 {
			ids = append(ids, k[idx+1:])
		} else {
			ids = append(ids, k)
		}
	}

	return ids
}

// ClearIndexes removes all indexes for the given tenant and region.
func (m *IndexManager) ClearIndexes(tenantID, region string) error {
	prefixes := []string{
		"idx_time:" + tenantID + ":" + region,
		"idx_level:" + tenantID + ":" + region,
		"idx_source:" + tenantID + ":" + region,
	}

	for _, prefix := range prefixes {
		bucket := m.storage.Bucket(prefix)
		iter := bucket.ScanPrefix(nil)
		var keys [][]byte
		for iter.Next() {
			keys = append(keys, iter.Key())
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
