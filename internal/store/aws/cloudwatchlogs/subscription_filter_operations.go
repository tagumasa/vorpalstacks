package cloudwatchlogs

import (
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
)

// PutSubscriptionFilter creates or updates a subscription filter.
func (s *Store) PutSubscriptionFilter(filter *SubscriptionFilter) error {
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
// subscription filter limit (max 2 distinct filter names) and creates or
// updates the filter. This prevents race conditions where concurrent
// PutSubscriptionFilter calls could exceed the limit.
func (s *Store) PutSubscriptionFilterWithLimitCheck(filter *SubscriptionFilter, maxPerGroup int) error {
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

	return s.PutSubscriptionFilter(filter)
}

// DeleteSubscriptionFilter deletes a subscription filter.
func (s *Store) DeleteSubscriptionFilter(logGroupName, filterName string) error {
	key := s.subscriptionFilterKey(logGroupName, filterName)
	if !s.Exists(key) {
		return ErrSubscriptionFilterNotFound
	}
	return s.Delete(key)
}

// GetSubscriptionFilter retrieves a subscription filter by name.
func (s *Store) GetSubscriptionFilter(logGroupName, filterName string) (*SubscriptionFilter, error) {
	key := s.subscriptionFilterKey(logGroupName, filterName)
	var p pb.SubscriptionFilter
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrSubscriptionFilterNotFound
	}
	return ProtoToSubscriptionFilter(&p), nil
}

// ListSubscriptionFilters lists subscription filters for a log group.
func (s *Store) ListSubscriptionFilters(logGroupName, filterNamePrefix string) ([]*SubscriptionFilter, error) {
	prefix := s.subscriptionFilterKey(logGroupName, "")
	var filters []*SubscriptionFilter

	if err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var p pb.SubscriptionFilter
		if err := proto.Unmarshal(value, &p); err == nil {
			filter := ProtoToSubscriptionFilter(&p)
			if filterNamePrefix == "" || strings.HasPrefix(filter.FilterName, filterNamePrefix) {
				filters = append(filters, filter)
			}
		}
		return nil
	}); err != nil {
		logs.Error("Failed to scan subscription filters", logs.String("logGroup", logGroupName), logs.Err(err))
	}

	return filters, nil
}
