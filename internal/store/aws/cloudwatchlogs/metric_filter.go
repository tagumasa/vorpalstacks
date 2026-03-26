// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
)

const metricFilterPrefix = "metric-filter:"

func (s *Store) metricFilterKey(logGroupName, filterName string) string {
	return metricFilterPrefix + logGroupName + "#" + filterName
}

func (s *Store) metricFilterPrefixForGroup(logGroupName string) string {
	return metricFilterPrefix + logGroupName + "#"
}

// PutMetricFilter creates or updates a metric filter.
func (s *Store) PutMetricFilter(filter *MetricFilter) error {
	key := s.metricFilterKey(filter.LogGroupName, filter.Name)
	return s.PutProto(key, MetricFilterToProto(filter))
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

// DeleteMetricFilter deletes a metric filter.
func (s *Store) DeleteMetricFilter(logGroupName, filterName string) error {
	key := s.metricFilterKey(logGroupName, filterName)
	if !s.Exists(key) {
		return ErrMetricFilterNotFound
	}
	return s.Delete(key)
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

	var filters []*MetricFilter
	count := 0
	var lastKey string

	err := s.BaseStore.ForEach(func(key string, value []byte) error {
		if !strings.HasPrefix(key, prefix) {
			return nil
		}

		if nextToken != "" && key <= nextToken {
			return nil
		}

		if count >= maxItems {
			return fmt.Errorf("stop")
		}

		var p pb.MetricFilter
		if err := proto.Unmarshal(value, &p); err != nil {
			return nil
		}
		filters = append(filters, ProtoToMetricFilter(&p))
		count++
		lastKey = key
		return nil
	})

	if err != nil && err.Error() != "stop" {
		return nil, "", err
	}

	if count >= maxItems {
		return filters, lastKey, nil
	}
	return filters, "", nil
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
	return fmt.Sprintf("arn:aws:logs:%s:%s:log-group:%s:metric-filter:%s",
		region, accountID, f.LogGroupName, f.Name)
}
