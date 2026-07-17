// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"time"

	arnutil "vorpalstacks/internal/utils/aws/arn"

	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
	"vorpalstacks/internal/store/aws/common"
)

const metricFilterPrefix = "metric-filter:"

func (s *Store) metricFilterKey(logGroupName, filterName string) string {
	return metricFilterPrefix + logGroupName + "#" + filterName
}

func (s *Store) metricFilterPrefixForGroup(logGroupName string) string {
	return metricFilterPrefix + logGroupName + "#"
}

// PutMetricFilter creates or updates a metric filter. When the filter is
// newly created, MetricFilterCount on the parent LogGroup is incremented
// atomically under filterMutex.
func (s *Store) PutMetricFilter(filter *MetricFilter) error {
	s.filterMutex.Lock()
	defer s.filterMutex.Unlock()

	key := s.metricFilterKey(filter.LogGroupName, filter.Name)
	isNew := !s.Exists(key)

	if err := s.PutProto(key, MetricFilterToProto(filter)); err != nil {
		return err
	}

	if isNew {
		lg, err := s.GetLogGroup(filter.LogGroupName)
		if err != nil {
			return err
		}
		lg.MetricFilterCount++
		if err := s.PutLogGroup(lg); err != nil {
			return err
		}
	}

	return nil
}

// GetMetricFilter retrieves a metric filter by log group and filter name.
func (s *Store) GetMetricFilter(logGroupName, filterName string) (*MetricFilter, error) {
	key := s.metricFilterKey(logGroupName, filterName)
	var p pb.MetricFilter
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrMetricFilterNotFound
	}
	return ProtoToMetricFilter(&p), nil
}

// DeleteMetricFilter deletes a metric filter and decrements
// MetricFilterCount on the parent LogGroup atomically under filterMutex.
func (s *Store) DeleteMetricFilter(logGroupName, filterName string) error {
	s.filterMutex.Lock()
	defer s.filterMutex.Unlock()

	key := s.metricFilterKey(logGroupName, filterName)
	if !s.Exists(key) {
		return ErrMetricFilterNotFound
	}
	if err := s.Delete(key); err != nil {
		return err
	}

	lg, err := s.GetLogGroup(logGroupName)
	if err != nil {
		return nil
	}
	if lg.MetricFilterCount > 0 {
		lg.MetricFilterCount--
		if err := s.PutLogGroup(lg); err != nil {
			return err
		}
	}

	return nil
}

// ListMetricFilters returns metric filters for a log group.
func (s *Store) ListMetricFilters(logGroupName, filterNamePrefix string, nextToken string, maxItems int) ([]*MetricFilter, string, error) {
	if maxItems <= 0 {
		maxItems = 50
	}

	prefix := s.metricFilterPrefixForGroup(logGroupName)
	if filterNamePrefix != "" {
		prefix = prefix + filterNamePrefix
	}

	result, err := common.ListProto(s.BaseStore, common.ListOptions{
		Prefix:   prefix,
		Marker:   nextToken,
		MaxItems: maxItems,
	}, func() *pb.MetricFilter { return new(pb.MetricFilter) }, nil)
	if err != nil {
		return nil, "", err
	}

	filters := make([]*MetricFilter, len(result.Items))
	for i := range result.Items {
		filters[i] = ProtoToMetricFilter(result.Items[i])
	}

	var token string
	if result.IsTruncated {
		token = result.NextMarker
	}

	return filters, token, nil
}

// NewMetricFilter creates a new metric filter.
func NewMetricFilter(logGroupName, filterName, filterPattern string, transformations []MetricTransformation) *MetricFilter {
	return &MetricFilter{
		Name:                  filterName,
		LogGroupName:          logGroupName,
		FilterPattern:         filterPattern,
		MetricTransformations: transformations,
		CreatedAt:             time.Now().UTC(),
	}
}

// ARN returns the ARN of the metric filter.
func (f *MetricFilter) ARN(accountID, region string) string {
	return arnutil.NewARNBuilder(accountID, region).CloudWatch().MetricFilter(f.LogGroupName, f.Name)
}
