package cloudwatchlogs

import (
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
)

// PutSubscriptionFilter creates or updates a subscription filter under
// the family mutex: the creation-time preservation read and the limit
// check's admission count serialise with every other family write.
func (s *Store) PutSubscriptionFilter(filter *SubscriptionFilter) error {
	s.subFilterMu.Lock()
	defer s.subFilterMu.Unlock()
	return s.putSubscriptionFilterLocked(filter)
}

// putSubscriptionFilterLocked is the lock-free write body for callers
// already holding subFilterMu (the limit-checked admission).
func (s *Store) putSubscriptionFilterLocked(filter *SubscriptionFilter) error {
	key := s.subscriptionFilterKey(filter.LogGroupName, filter.FilterName)
	var existing pb.SubscriptionFilter
	if s.GetProto(key, &existing) == nil {
		filter.CreationTime = time.UnixMilli(existing.CreatedAt)
	} else {
		filter.CreationTime = time.Now().UTC()
	}
	return s.PutProto(key, SubscriptionFilterToProto(filter))
}

// PutSubscriptionFilterWithLimitCheck atomically checks the per-log-group
// subscription filter limit (the documented per-group quota, carried by
// MaxSubscriptionFiltersPerLogGroup) and creates or updates the filter. This prevents race conditions where concurrent
// PutSubscriptionFilter calls could exceed the limit. The admission also
// holds the group's write lock — the DeleteLogGroup teardown lock — and
// re-checks the group inside it: the service layer's earlier group check
// runs outside any critical section, so without the re-check a teardown
// interleaving this put would leave the filter permanently outliving its
// deleted group (the metric-filter family's pattern).
func (s *Store) PutSubscriptionFilterWithLimitCheck(filter *SubscriptionFilter, maxPerGroup int) error {
	gLock := s.groupLock(filter.LogGroupName)
	gLock.Lock()
	defer gLock.Unlock()

	if _, err := s.GetLogGroup(filter.LogGroupName); err != nil {
		return err
	}

	s.subFilterMu.Lock()
	defer s.subFilterMu.Unlock()

	filters, err := s.ListSubscriptionFilters(filter.LogGroupName, "")
	if err != nil {
		return err
	}

	existingCount := 0
	isUpdate := false
	for _, f := range filters {
		if f.FilterName == filter.FilterName {
			isUpdate = true
			continue
		}
		existingCount++
	}

	if !isUpdate && existingCount >= maxPerGroup {
		return ErrLimitExceeded
	}

	return s.putSubscriptionFilterLocked(filter)
}

// DeleteSubscriptionFilter deletes a subscription filter under the same
// mutex as the limit check, so a concurrent admission counts the live
// filter set this delete is part of instead of a stale snapshot.
func (s *Store) DeleteSubscriptionFilter(logGroupName, filterName string) error {
	s.subFilterMu.Lock()
	defer s.subFilterMu.Unlock()

	key := s.subscriptionFilterKey(logGroupName, filterName)
	if !s.Exists(key) {
		return ErrSubscriptionFilterNotFound
	}
	return s.Delete(key)
}

// ListSubscriptionFilters lists subscription filters for a log group. A
// scan failure propagates: the per-group limit check counts from this
// list, and a partial list could admit a filter beyond the limit.
func (s *Store) ListSubscriptionFilters(logGroupName, filterNamePrefix string) ([]*SubscriptionFilter, error) {
	prefix := s.subscriptionFilterKey(logGroupName, "")
	var filters []*SubscriptionFilter

	if err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var p pb.SubscriptionFilter
		if err := proto.Unmarshal(value, &p); err != nil {
			// A record that fails to decode is invisible to the caller;
			// log it so corruption is discoverable instead of silently
			// shrinking the list.
			logs.Warn("Corrupt subscription filter record",
				logs.String("key", key), logs.Err(err))
			return nil
		}
		filter := ProtoToSubscriptionFilter(&p)
		if filterNamePrefix == "" || strings.HasPrefix(filter.FilterName, filterNamePrefix) {
			filters = append(filters, filter)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return filters, nil
}

// SubscriptionFilter represents a CloudWatch Logs subscription filter.
type SubscriptionFilter struct {
	LogGroupName           string    `json:"logGroupName"`
	FilterName             string    `json:"filterName"`
	FilterPattern          string    `json:"filterPattern"`
	DestinationArn         string    `json:"destinationArn"`
	RoleArn                string    `json:"roleArn"`
	Distribution           string    `json:"distribution"`
	ApplyOnTransformedLogs bool      `json:"applyOnTransformedLogs,omitempty"`
	FieldSelectionCriteria string    `json:"fieldSelectionCriteria,omitempty"`
	EmitSystemFields       []string  `json:"emitSystemFields,omitempty"`
	CreationTime           time.Time `json:"creationTime"`
}
