// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"context"
	"errors"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
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
// atomically within a storage transaction so concurrent PutLogEvents
// operations observe a consistent count.
func (s *Store) PutMetricFilter(filter *MetricFilter) error {
	key := s.metricFilterKey(filter.LogGroupName, filter.Name)

	if s.ts != nil {
		ctx := context.Background()
		return s.ts.Update(ctx, func(txn storage.Transaction) error {
			// Check existence within the transaction for atomic check-then-act.
			filterBytes, err := txn.Bucket(s.bucketName).Get([]byte(key))
			if err != nil {
				return err
			}
			isNew := filterBytes == nil

			filterData, err := proto.Marshal(MetricFilterToProto(filter))
			if err != nil {
				return err
			}
			if err := txn.Bucket(s.bucketName).Put([]byte(key), filterData); err != nil {
				return err
			}

			if isNew {
				lgBytes, err := txn.Bucket(s.bucketName).Get([]byte(s.logGroupKey(filter.LogGroupName)))
				if err != nil || lgBytes == nil {
					return ErrLogGroupNotFound
				}
				var lgProto pb.LogGroup
				if err := proto.Unmarshal(lgBytes, &lgProto); err != nil {
					return err
				}
				lgProto.MetricFilterCount++
				lgData, err := proto.Marshal(&lgProto)
				if err != nil {
					return err
				}
				return txn.Bucket(s.bucketName).Put([]byte(s.logGroupKey(filter.LogGroupName)), lgData)
			}
			return nil
		})
	}

	// Fallback: sequential writes with no cross-path atomicity guarantee.
	isNew := !s.Exists(key)
	if err := s.PutProto(key, MetricFilterToProto(filter)); err != nil {
		return err
	}
	if isNew {
		lg, err := s.GetLogGroup(filter.LogGroupName)
		if err != nil {
			// Rollback: remove the filter we just wrote so the count
			// does not drift if the parent LogGroup lookup fails.
			// The rollback error is logged so silent failures remain
			// visible, but the original error is preserved for the caller.
			if rmErr := s.Delete(key); rmErr != nil {
				logs.Warn("Failed to rollback metric filter after LogGroup lookup failure",
					logs.String("key", key), logs.Err(rmErr))
			}
			return err
		}
		lg.MetricFilterCount++
		if err := s.PutLogGroup(lg); err != nil {
			if rmErr := s.Delete(key); rmErr != nil {
				logs.Warn("Failed to rollback metric filter after LogGroup update failure",
					logs.String("key", key), logs.Err(rmErr))
			}
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
// MetricFilterCount on the parent LogGroup atomically within a storage
// transaction so concurrent readers always observe a consistent count.
func (s *Store) DeleteMetricFilter(logGroupName, filterName string) error {
	key := s.metricFilterKey(logGroupName, filterName)

	if s.ts != nil {
		ctx := context.Background()
		return s.ts.Update(ctx, func(txn storage.Transaction) error {
			filterBytes, err := txn.Bucket(s.bucketName).Get([]byte(key))
			if err != nil {
				return err
			}
			if filterBytes == nil {
				return ErrMetricFilterNotFound
			}
			if err := txn.Bucket(s.bucketName).Delete([]byte(key)); err != nil {
				return err
			}

			lgBytes, err := txn.Bucket(s.bucketName).Get([]byte(s.logGroupKey(logGroupName)))
			if err != nil {
				return err
			}
			if lgBytes == nil {
				// LogGroup already deleted; the count is implicitly zero
				// so there is nothing to decrement.
				return nil
			}
			var lgProto pb.LogGroup
			if err := proto.Unmarshal(lgBytes, &lgProto); err != nil {
				return err
			}
			if lgProto.MetricFilterCount > 0 {
				lgProto.MetricFilterCount--
				lgData, err := proto.Marshal(&lgProto)
				if err != nil {
					return err
				}
				return txn.Bucket(s.bucketName).Put([]byte(s.logGroupKey(logGroupName)), lgData)
			}
			return nil
		})
	}

	// Fallback path used when no transactional store is configured:
	// sequential writes with explicit error-type distinction so the caller
	// still receives ErrMetricFilterNotFound vs. a generic storage error.
	if !s.Exists(key) {
		return ErrMetricFilterNotFound
	}
	if err := s.Delete(key); err != nil {
		return err
	}
	lg, err := s.GetLogGroup(logGroupName)
	if err != nil {
		if errors.Is(err, ErrLogGroupNotFound) {
			return nil
		}
		return err
	}
	if lg.MetricFilterCount > 0 {
		lg.MetricFilterCount--
		return s.PutLogGroup(lg)
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
