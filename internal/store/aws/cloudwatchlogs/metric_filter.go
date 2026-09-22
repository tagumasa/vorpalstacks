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

	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
	"vorpalstacks/internal/store/aws/common"
)

func (s *Store) metricFilterKey(logGroupName, filterName string) string {
	return keyPrefixMetricFilter + escapePath(logGroupName) + ":" + filterName
}

func (s *Store) metricFilterPrefixForGroup(logGroupName string) string {
	return keyPrefixMetricFilter + escapePath(logGroupName) + ":"
}

// PutMetricFilter creates or updates a metric filter. An update preserves
// the stored creation time (DescribeMetricFilters reports the filter's
// original creation, not the last update). When the filter is newly
// created, the per-log-group metric filter quota is enforced and
// MetricFilterCount on the parent LogGroup is incremented atomically
// within a storage transaction so concurrent PutLogEvents operations
// observe a consistent count. The whole mutation holds the group's lock
// so the count update cannot lose a concurrent configuration change to
// the same record (and vice versa).
func (s *Store) PutMetricFilter(filter *MetricFilter) error {
	key := s.metricFilterKey(filter.LogGroupName, filter.Name)

	gLock := s.groupLock(filter.LogGroupName)
	gLock.Lock()
	defer gLock.Unlock()

	if s.ts != nil {
		ctx := context.Background()
		return s.ts.Update(ctx, func(txn storage.Transaction) error {
			// Check existence within the transaction for atomic check-then-act.
			filterBytes, err := txn.Bucket(s.bucketName).Get([]byte(key))
			if err != nil {
				return err
			}
			isNew := filterBytes == nil
			if !isNew {
				var old pb.MetricFilter
				if err := proto.Unmarshal(filterBytes, &old); err != nil {
					return err
				}
				if old.CreatedAt > 0 {
					filter.CreatedAt = time.UnixMilli(old.CreatedAt).UTC()
				}
			}

			var lgProto *pb.LogGroup
			if isNew {
				lgBytes, err := txn.Bucket(s.bucketName).Get([]byte(s.logGroupKey(filter.LogGroupName)))
				if err != nil {
					return err
				}
				if lgBytes == nil {
					return ErrLogGroupNotFound
				}
				lgProto = new(pb.LogGroup)
				if err := proto.Unmarshal(lgBytes, lgProto); err != nil {
					return err
				}
				if int(lgProto.MetricFilterCount) >= MaxMetricFiltersPerGroup {
					return ErrLimitExceeded
				}
			}

			filterData, err := proto.Marshal(MetricFilterToProto(filter))
			if err != nil {
				return err
			}
			if err := txn.Bucket(s.bucketName).Put([]byte(key), filterData); err != nil {
				return err
			}

			if isNew {
				lgProto.MetricFilterCount++
				lgData, err := proto.Marshal(lgProto)
				if err != nil {
					return err
				}
				return txn.Bucket(s.bucketName).Put([]byte(s.logGroupKey(filter.LogGroupName)), lgData)
			}
			return nil
		})
	}

	// Fallback: sequential writes with no cross-path atomicity guarantee.
	oldFilter, err := s.GetMetricFilter(filter.LogGroupName, filter.Name)
	if err != nil && !errors.Is(err, ErrMetricFilterNotFound) {
		return err
	}
	isNew := oldFilter == nil
	if isNew {
		lg, err := s.GetLogGroup(filter.LogGroupName)
		if err != nil {
			return err
		}
		if int(lg.MetricFilterCount) >= MaxMetricFiltersPerGroup {
			return ErrLimitExceeded
		}
	} else if !oldFilter.CreatedAt.IsZero() {
		filter.CreatedAt = oldFilter.CreatedAt
	}
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
		if err := s.putLogGroupLocked(lg); err != nil {
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
// A missing record reports the not-found sentinel; any other storage
// failure propagates as a storage error.
func (s *Store) GetMetricFilter(logGroupName, filterName string) (*MetricFilter, error) {
	key := s.metricFilterKey(logGroupName, filterName)
	var p pb.MetricFilter
	if err := s.GetProto(key, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrMetricFilterNotFound
		}
		return nil, err
	}
	return ProtoToMetricFilter(&p), nil
}

// DeleteMetricFilter deletes a metric filter and decrements
// MetricFilterCount on the parent LogGroup atomically within a storage
// transaction so concurrent readers always observe a consistent count.
// Like PutMetricFilter it holds the group's lock for the whole mutation.
func (s *Store) DeleteMetricFilter(logGroupName, filterName string) error {
	key := s.metricFilterKey(logGroupName, filterName)

	gLock := s.groupLock(logGroupName)
	gLock.Lock()
	defer gLock.Unlock()

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

	// Fallback path used when no transactional store is configured. The
	// steps run in the order whose partial failure leaves the least
	// drift: the parent LogGroup is read BEFORE the filter record is
	// deleted (a read failure now changes nothing — retrying is clean),
	// and the count decrement follows the successful delete. A failure
	// of the final LogGroup write still leaves the count one high; that
	// residual tear is unavoidable without transactions and is logged.
	lg, err := s.GetLogGroup(logGroupName)
	if err != nil && !errors.Is(err, ErrLogGroupNotFound) {
		return err
	}
	if !s.Exists(key) {
		return ErrMetricFilterNotFound
	}
	if err := s.Delete(key); err != nil {
		return err
	}
	if lg == nil {
		// The group is gone; its count went with it.
		return nil
	}
	if lg.MetricFilterCount > 0 {
		lg.MetricFilterCount--
		if err := s.putLogGroupLocked(lg); err != nil {
			logs.Error("Metric filter count left high: LogGroup write failed after filter deletion",
				logs.String("logGroupName", logGroupName), logs.Err(err))
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

// MetricFilter represents a CloudWatch Logs metric filter.
type MetricFilter struct {
	Name                      string                 `json:"name"`
	LogGroupName              string                 `json:"logGroupName"`
	FilterPattern             string                 `json:"filterPattern"`
	MetricTransformations     []MetricTransformation `json:"metricTransformations"`
	ApplyOnTransformedLogs    bool                   `json:"applyOnTransformedLogs,omitempty"`
	FieldSelectionCriteria    string                 `json:"fieldSelectionCriteria,omitempty"`
	EmitSystemFieldDimensions []string               `json:"emitSystemFieldDimensions,omitempty"`
	CreatedAt                 time.Time              `json:"createdAt"`
}

// MetricTransformation represents a metric transformation for a metric filter.
type MetricTransformation struct {
	MetricName      string  `json:"metricName"`
	MetricNamespace string  `json:"metricNamespace"`
	MetricValue     string  `json:"metricValue"`
	DefaultValue    float64 `json:"defaultValue,omitempty"`
	DefaultValueSet bool    `json:"defaultValueSet,omitempty"`
	// Dimensions is the transformation's dimensions member — "The fields
	// to use as dimensions for the metric" — whose values are value
	// references into the matched event (the operation documentation's
	// own worked example runs "dimensions": {"Request": "$request"}).
	Dimensions map[string]string `json:"dimensions,omitempty"`
	// Unit is the transformation's unit member ("The unit to assign to
	// the metric. If you omit this, the unit is set as None"; the empty
	// string is the omitted, None form).
	Unit string `json:"unit,omitempty"`

	// Wire-presence flags for the required members whose shapes allow the
	// empty string (minimum length zero): the members are required on
	// MetricTransformation, so their absence rejects in the Core while a
	// present empty value stays legal. Presence is a parse-time fact and
	// never serialises.
	MetricNameSet      bool `json:"-"`
	MetricNamespaceSet bool `json:"-"`
	MetricValueSet     bool `json:"-"`
}
