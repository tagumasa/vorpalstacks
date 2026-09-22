// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"fmt"
	"testing"
	"time"
)

// TestPutMetricFilterPerGroupQuota pins the documented per-log-group quota
// (100): the hundred-and-first distinct filter name rejects with the
// limit sentinel while an update of an existing name at the cap still
// succeeds.
func TestPutMetricFilterPerGroupQuota(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "quota-group"}); err != nil {
		t.Fatal(err)
	}

	transform := []MetricTransformation{{
		MetricName:      "Count",
		MetricNamespace: "Quota",
		MetricValue:     "1",
	}}
	for i := 0; i < MaxMetricFiltersPerGroup; i++ {
		if err := s.PutMetricFilter(NewMetricFilter("quota-group", fmt.Sprintf("f-%03d", i), "ERROR", transform)); err != nil {
			t.Fatalf("filter %d within the quota must succeed: %v", i, err)
		}
	}

	if err := s.PutMetricFilter(NewMetricFilter("quota-group", "f-over", "ERROR", transform)); err == nil {
		t.Fatal("the filter beyond the quota must reject")
	} else if err != ErrLimitExceeded {
		t.Fatalf("the filter beyond the quota must reject with the limit sentinel, got %v", err)
	}

	// Updating an existing name at the cap is not a new filter.
	updated := NewMetricFilter("quota-group", "f-000", "WARN", transform)
	if err := s.PutMetricFilter(updated); err != nil {
		t.Fatalf("updating an existing filter at the quota cap must succeed: %v", err)
	}
	got, err := s.GetMetricFilter("quota-group", "f-000")
	if err != nil {
		t.Fatal(err)
	}
	if got.FilterPattern != "WARN" {
		t.Fatalf("the update must take effect, pattern still %q", got.FilterPattern)
	}
}

// TestPutMetricFilterPreservesCreationTime pins that an update keeps the
// stored creation time: DescribeMetricFilters reports the filter's
// original creation, not the last update.
func TestPutMetricFilterPreservesCreationTime(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "ts-group"}); err != nil {
		t.Fatal(err)
	}

	transform := []MetricTransformation{{
		MetricName:      "Count",
		MetricNamespace: "Ts",
		MetricValue:     "1",
	}}
	if err := s.PutMetricFilter(NewMetricFilter("ts-group", "f", "ERROR", transform)); err != nil {
		t.Fatal(err)
	}
	first, err := s.GetMetricFilter("ts-group", "f")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(20 * time.Millisecond)
	if err := s.PutMetricFilter(NewMetricFilter("ts-group", "f", "WARN", transform)); err != nil {
		t.Fatal(err)
	}
	second, err := s.GetMetricFilter("ts-group", "f")
	if err != nil {
		t.Fatal(err)
	}

	if !second.CreatedAt.Equal(first.CreatedAt) {
		t.Fatalf("an update must preserve the creation time: first %v, second %v", first.CreatedAt, second.CreatedAt)
	}
	if second.FilterPattern != "WARN" {
		t.Fatalf("the update must still change the pattern, got %q", second.FilterPattern)
	}
}
