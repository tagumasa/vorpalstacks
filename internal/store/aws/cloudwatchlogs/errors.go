package cloudwatchlogs

import (
	"encoding/json"
	"errors"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrLogGroupNotFound is returned when the specified CloudWatch Logs
	// log group does not exist.
	ErrLogGroupNotFound = errors.New("log group not found")

	// ErrLogGroupAlreadyExists is returned when attempting to create a log group
	// that already exists.
	ErrLogGroupAlreadyExists = errors.New("log group already exists")

	// ErrLogStreamNotFound is returned when the specified log stream
	// does not exist.
	ErrLogStreamNotFound = errors.New("log stream not found")

	// ErrLogStreamAlreadyExists is returned when attempting to create a log stream
	// that already exists.
	ErrLogStreamAlreadyExists = errors.New("log stream already exists")

	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = errors.New("resource not found")

	// ErrLimitExceeded is returned when a service limit has been exceeded.
	ErrLimitExceeded = errors.New("limit exceeded")

	// ErrMetricFilterNotFound is returned when the specified metric filter
	// does not exist.
	ErrMetricFilterNotFound = errors.New("metric filter not found")

	// ErrSubscriptionFilterNotFound is returned when the specified subscription filter
	// does not exist.
	ErrSubscriptionFilterNotFound = errors.New("subscription filter not found")

	// ErrDestinationNotFound is returned when the specified destination
	// does not exist. A destination-already-exists sentinel does not
	// belong here: PutDestination is the documented create-or-update
	// (an existing destination is an update target, never a conflict).
	ErrDestinationNotFound = errors.New("destination not found")

	// ErrLogGroupDeletionProtected is returned when the guarded group
	// delete meets a record with deletion protection enabled.
	ErrLogGroupDeletionProtected = errors.New("log group deletion protection enabled")
)

// getJSONRecord loads one JSON-regime record: a missing record answers
// the given sentinel, any other storage failure propagates — one body
// behind the JSON-regime getters (the proto-regime getters carry the
// same discrimination through common.IsNotFound at their own sites).
func getJSONRecord[T any](s *Store, key string, sentinel error) (*T, error) {
	var record T
	if err := s.Get(key, &record); err != nil {
		if common.IsNotFound(err) {
			return nil, sentinel
		}
		return nil, err
	}
	return &record, nil
}

// listJSONRecords scans one prefix in the JSON regime and returns the
// records filter accepts (a nil filter accepts every record). A record
// that fails to decode is invisible to the caller; the corruption log
// keeps it discoverable instead of silently shrinking the list — one
// body behind the JSON-regime listings, so that convention cannot drift
// per family. The label names the record family in the log line ("Corrupt
// " + label). The proto-regime listings carry the same discipline at
// their own sites (the regime split getJSONRecord documents).
func listJSONRecords[T any](s *Store, prefix, label string, filter func(*T) bool) ([]*T, error) {
	var records []*T
	if err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var record T
		if err := json.Unmarshal(value, &record); err != nil {
			logs.Warn("Corrupt "+label,
				logs.String("key", key), logs.Err(err))
			return nil
		}
		if filter == nil || filter(&record) {
			records = append(records, &record)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return records, nil
}
