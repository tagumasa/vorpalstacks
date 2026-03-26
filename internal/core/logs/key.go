// Package logs provides logging functionality for vorpalstacks.
package logs

import (
	"fmt"
	"strings"
)

// Key represents a unique identifier for a log entry.
type Key struct {
	Namespace string
	TenantID  string
	Region    string
	ID        string
}

// NewKey creates a new Key with the given parameters.
func NewKey(namespace, tenantID, region, id string) *Key {
	return &Key{
		Namespace: namespace,
		TenantID:  tenantID,
		Region:    region,
		ID:        id,
	}
}

// NewLogKey creates a new Key for a log entry.
func NewLogKey(tenantID, region, id string) *Key {
	return NewKey("logs", tenantID, region, id)
}

// Encode converts the Key to a string representation.
func (k *Key) Encode() string {
	return fmt.Sprintf("%s:%s:%s:%s", k.Namespace, k.TenantID, k.Region, k.ID)
}

// EncodePrefix returns the prefix portion of the Key for prefix-based queries.
func (k *Key) EncodePrefix() string {
	return fmt.Sprintf("%s:%s:%s", k.Namespace, k.TenantID, k.Region)
}

// DecodeKey parses a Key from its string representation.
func DecodeKey(s string) *Key {
	parts := strings.SplitN(s, ":", 4)
	if len(parts) < 4 {
		return nil
	}
	return &Key{
		Namespace: parts[0],
		TenantID:  parts[1],
		Region:    parts[2],
		ID:        parts[3],
	}
}

// IndexKey represents a key for indexing log entries.
type IndexKey struct {
	IndexType IndexType
	TenantID  string
	Region    string
	Segment1  string
	Segment2  string
	EntryID   string
}

// IndexType represents the type of index.
type IndexType int

// IndexType constants define the different ways log entries can be indexed.
const (
	IndexByTime IndexType = iota
	IndexByLevel
	IndexBySource
)

// Encode converts the IndexKey to a string representation.
func (k *IndexKey) Encode() string {
	var prefix string
	switch k.IndexType {
	case IndexByTime:
		prefix = "idx_time"
	case IndexByLevel:
		prefix = "idx_level"
	case IndexBySource:
		prefix = "idx_source"
	}
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", prefix, k.TenantID, k.Region, k.Segment1, k.Segment2, k.EntryID)
}

// EncodePrefix returns the prefix portion of the IndexKey for prefix-based queries.
func (k *IndexKey) EncodePrefix() string {
	var prefix string
	switch k.IndexType {
	case IndexByTime:
		prefix = "idx_time"
	case IndexByLevel:
		prefix = "idx_level"
	case IndexBySource:
		prefix = "idx_source"
	}
	if k.Segment2 != "" {
		return fmt.Sprintf("%s:%s:%s:%s:%s", prefix, k.TenantID, k.Region, k.Segment1, k.Segment2)
	}
	if k.Segment1 != "" {
		return fmt.Sprintf("%s:%s:%s:%s", prefix, k.TenantID, k.Region, k.Segment1)
	}
	return fmt.Sprintf("%s:%s:%s", prefix, k.TenantID, k.Region)
}

// NewTimeIndexKey creates an IndexKey for time-based indexing.
func NewTimeIndexKey(tenantID, region, dateHour, entryID string) *IndexKey {
	return &IndexKey{
		IndexType: IndexByTime,
		TenantID:  tenantID,
		Region:    region,
		Segment1:  dateHour,
		EntryID:   entryID,
	}
}

// NewLevelIndexKey creates an IndexKey for level-based indexing.
func NewLevelIndexKey(tenantID, region, level string, timestamp int64, entryID string) *IndexKey {
	return &IndexKey{
		IndexType: IndexByLevel,
		TenantID:  tenantID,
		Region:    region,
		Segment1:  level,
		Segment2:  fmt.Sprintf("%d", timestamp),
		EntryID:   entryID,
	}
}

// NewSourceIndexKey creates an IndexKey for source-based indexing.
func NewSourceIndexKey(tenantID, region, source string, timestamp int64, entryID string) *IndexKey {
	return &IndexKey{
		IndexType: IndexBySource,
		TenantID:  tenantID,
		Region:    region,
		Segment1:  source,
		Segment2:  fmt.Sprintf("%d", timestamp),
		EntryID:   entryID,
	}
}
