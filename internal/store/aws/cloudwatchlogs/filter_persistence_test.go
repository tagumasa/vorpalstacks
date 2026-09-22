// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"reflect"
	"testing"
)

// TestFilterSelectionMembersSurvivePersistence pins the six transform-era
// members through the storage round-trip: the proto messages carry them as
// append-only fields, so a Put followed by a store read returns exactly the
// members that were accepted — persistence no longer silently drops client
// state.
func TestFilterSelectionMembersSurvivePersistence(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "persist"}); err != nil {
		t.Fatal(err)
	}

	if err := s.PutMetricFilter(&MetricFilter{
		Name:                      "mf",
		LogGroupName:              "persist",
		FilterPattern:             "{ $.latency > 0 }",
		MetricTransformations:     []MetricTransformation{{MetricName: "Lat", MetricNamespace: "N", MetricValue: "$.latency"}},
		ApplyOnTransformedLogs:    true,
		FieldSelectionCriteria:    "$.user,$.latency",
		EmitSystemFieldDimensions: []string{"$.host", "@logStream"},
	}); err != nil {
		t.Fatal(err)
	}
	mf, err := s.GetMetricFilter("persist", "mf")
	if err != nil {
		t.Fatal(err)
	}
	if !mf.ApplyOnTransformedLogs {
		t.Fatal("applyOnTransformedLogs must survive persistence")
	}
	if mf.FieldSelectionCriteria != "$.user,$.latency" {
		t.Fatalf("fieldSelectionCriteria must survive persistence, got %q", mf.FieldSelectionCriteria)
	}
	if want := []string{"$.host", "@logStream"}; !reflect.DeepEqual(mf.EmitSystemFieldDimensions, want) {
		t.Fatalf("emitSystemFieldDimensions must survive persistence, got %v", mf.EmitSystemFieldDimensions)
	}

	if err := s.PutSubscriptionFilter(&SubscriptionFilter{
		LogGroupName:           "persist",
		FilterName:             "sf",
		FilterPattern:          "ERROR",
		DestinationArn:         s.arnBuilder.CloudWatch().LogGroup("persist"),
		ApplyOnTransformedLogs: true,
		FieldSelectionCriteria: "$.user",
		EmitSystemFields:       []string{"@timestamp", "@message"},
	}); err != nil {
		t.Fatal(err)
	}
	filters, err := s.ListSubscriptionFilters("persist", "sf")
	if err != nil {
		t.Fatal(err)
	}
	if len(filters) != 1 {
		t.Fatalf("stored filter count = %d, want 1", len(filters))
	}
	sf := filters[0]
	if !sf.ApplyOnTransformedLogs {
		t.Fatal("applyOnTransformedLogs must survive persistence")
	}
	if sf.FieldSelectionCriteria != "$.user" {
		t.Fatalf("fieldSelectionCriteria must survive persistence, got %q", sf.FieldSelectionCriteria)
	}
	if want := []string{"@timestamp", "@message"}; !reflect.DeepEqual(sf.EmitSystemFields, want) {
		t.Fatalf("emitSystemFields must survive persistence, got %v", sf.EmitSystemFields)
	}
}
