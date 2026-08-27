package iot

import (
	"testing"
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// activeViolationResponse must project the ActiveViolation member set: the
// ViolationEvent-only members (violationEventType, violationEventTime,
// metricValue) must not appear, and the metric value is emitted under
// lastViolationValue. No violation engine generates records on this
// platform, so the wire shape is pinned here rather than through the SDK.
func TestActiveViolationResponseShape(t *testing.T) {
	ts := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	resp := activeViolationResponse(&iotstore.ViolationEvent{
		ViolationID:         "v-1",
		ThingName:           "thing-1",
		SecurityProfileName: "profile-1",
		ViolationEventType:  "in-violation",
		MetricValue:         &iotstore.MetricValue{Count: 3},
		ViolationEventTime:  ts,
	})
	for _, forbidden := range []string{"violationEventType", "violationEventTime", "metricValue"} {
		if _, ok := resp[forbidden]; ok {
			t.Errorf("active violation response carries %s, which is not an ActiveViolation member", forbidden)
		}
	}
	if resp["violationId"] != "v-1" || resp["thingName"] != "thing-1" {
		t.Errorf("identity members missing: %v", resp)
	}
	if last, ok := resp["lastViolationValue"].(map[string]interface{}); !ok || last["count"] != int64(3) {
		t.Errorf("expected lastViolationValue.count=3, got %v", resp["lastViolationValue"])
	}
	if resp["lastViolationTime"] != ts.Unix() || resp["violationStartTime"] != ts.Unix() {
		t.Errorf("expected violation time members as epoch seconds, got %v/%v", resp["lastViolationTime"], resp["violationStartTime"])
	}
}
