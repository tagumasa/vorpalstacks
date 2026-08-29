package iot

import (
	"testing"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// TestBehaviorMetricDimensionOperatorRoundTrip pins the parse, proto
// carrier and response projection of the metricDimension operator: the
// {dimensionName, operator} structure survives create to describe with the
// IN/NOT_IN enum validated on the way in.
func TestBehaviorMetricDimensionOperatorRoundTrip(t *testing.T) {
	behaviors, err := rawToBehaviors([]map[string]interface{}{{
		"name":            "dim-behavior",
		"metricDimension": map[string]interface{}{"dimensionName": "d", "operator": "NOT_IN"},
	}})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(behaviors) != 1 {
		t.Fatalf("expected one behavior, got %d", len(behaviors))
	}
	b := behaviors[0]
	if b.MetricDimension != "d" || b.MetricDimensionOperator != "NOT_IN" {
		t.Fatalf("dimension = %q/%q, want d/NOT_IN", b.MetricDimension, b.MetricDimensionOperator)
	}

	// The storage proto round trip carries both parts through the
	// JSON-encoded carrier.
	pb, err := iotstore.BehaviorToProto(b)
	if err != nil {
		t.Fatalf("proto encode failed: %v", err)
	}
	decoded := iotstore.ProtoToBehavior(pb)
	if decoded.MetricDimension != "d" || decoded.MetricDimensionOperator != "NOT_IN" {
		t.Fatalf("decoded dimension = %q/%q, want d/NOT_IN", decoded.MetricDimension, decoded.MetricDimensionOperator)
	}

	// The response projection emits the structured dimension with the
	// operator.
	resp := behaviorResponse(decoded)
	dim, ok := resp["metricDimension"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected a structured metricDimension, got %#v", resp["metricDimension"])
	}
	if dim["dimensionName"] != "d" || dim["operator"] != "NOT_IN" {
		t.Errorf("projected dimension = %#v", dim)
	}

	// An off-enum operator is rejected.
	if _, err := rawToBehaviors([]map[string]interface{}{{
		"name":            "bad",
		"metricDimension": map[string]interface{}{"dimensionName": "d", "operator": "AROUND"},
	}}); err == nil {
		t.Error("expected an off-enum operator to be rejected")
	}
}
