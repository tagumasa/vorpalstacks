package chunk

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

// TimestreamEntry represents a single Timestream record as a chunk entry.
type TimestreamEntry struct {
	DatabaseName     string         `json:"databaseName"`
	TableName        string         `json:"tableName"`
	Dimensions       []Dimension    `json:"dimensions,omitempty"`
	MeasureName      string         `json:"measureName"`
	MeasureValue     string         `json:"measureValue,omitempty"`
	MeasureValueType string         `json:"measureValueType,omitempty"`
	MeasureValues    []MeasureValue `json:"measureValues,omitempty"`
	Time             time.Time      `json:"timestamp"`
	Version          int64          `json:"version,omitempty"`
}

// Dimension represents a Timestream dimension.
type Dimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MeasureValue represents a Timestream measure value.
type MeasureValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

// Timestamp returns the entry's timestamp as int64.
func (e *TimestreamEntry) Timestamp() int64 {
	return e.Time.UnixNano()
}

// Message returns the entry's serialized data.
func (e *TimestreamEntry) Message() []byte {
	data, _ := json.Marshal(e)
	return data
}

// BuildDimensionsKey creates a sorted key from dimensions.
func (e *TimestreamEntry) BuildDimensionsKey() string {
	if len(e.Dimensions) == 0 {
		return "_empty"
	}
	parts := make([]string, len(e.Dimensions))
	for i, d := range e.Dimensions {
		parts[i] = d.Name + "=" + d.Value
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}
