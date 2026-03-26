// Package logs provides logging functionality for vorpalstacks.
package logs

import (
	"time"
)

// QueryFilter defines criteria for querying log entries.
type QueryFilter struct {
	TenantID  *string
	Region    *string
	Level     *Level
	Source    *string
	StartTime *time.Time
	EndTime   *time.Time
	Limit     int
}

// NewQueryFilter creates a new empty QueryFilter.
func NewQueryFilter() *QueryFilter {
	return &QueryFilter{}
}

// WithTenantID sets the tenant ID filter.
func (f *QueryFilter) WithTenantID(id string) *QueryFilter {
	f.TenantID = &id
	return f
}

// WithRegion sets the region filter.
func (f *QueryFilter) WithRegion(region string) *QueryFilter {
	f.Region = &region
	return f
}

// WithLevel sets the log level filter.
func (f *QueryFilter) WithLevel(level Level) *QueryFilter {
	f.Level = &level
	return f
}

// WithSource sets the source filter.
func (f *QueryFilter) WithSource(source string) *QueryFilter {
	f.Source = &source
	return f
}

// WithStartTime sets the start time filter.
func (f *QueryFilter) WithStartTime(t time.Time) *QueryFilter {
	f.StartTime = &t
	return f
}

// WithEndTime sets the end time filter.
func (f *QueryFilter) WithEndTime(t time.Time) *QueryFilter {
	f.EndTime = &t
	return f
}

// WithLimit sets the maximum number of results.
func (f *QueryFilter) WithLimit(limit int) *QueryFilter {
	f.Limit = limit
	return f
}

// WithTimeRange sets both start and end time filters.
func (f *QueryFilter) WithTimeRange(start, end time.Time) *QueryFilter {
	f.StartTime = &start
	f.EndTime = &end
	return f
}
