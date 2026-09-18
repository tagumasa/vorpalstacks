package cloudtrail

import (
	"errors"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
)

// The pins below drive lookupEventsCore's request validation against a
// real store: the time-range rule, the one-attribute rule, the
// EventCategory enum, and the invalid-continuation-token refusal — each
// surfacing the model-declared error shape.

func TestLookupEventsCoreValidation(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	code := func(t *testing.T, err error, want string) {
		t.Helper()
		if err == nil {
			t.Fatalf("expected %s, got success", want)
		}
		var awsErr *awserrors.AWSError
		if !errors.As(err, &awsErr) || awsErr.GetCode() != want {
			t.Fatalf("error = %v, want %s", err, want)
		}
	}

	t.Run("EndTime before StartTime is refused", func(t *testing.T) {
		start := time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC)
		end := start.Add(-time.Minute)
		_, err := svc.lookupEventsCore(store, LookupEventsInput{
			StartTimeRaw: float64(start.Unix()), EndTimeRaw: float64(end.Unix()),
		})
		code(t, err, "InvalidTimeRangeException")
	})

	t.Run("EndTime equal to StartTime is accepted", func(t *testing.T) {
		at := time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC)
		if _, err := svc.lookupEventsCore(store, LookupEventsInput{
			StartTimeRaw: float64(at.Unix()), EndTimeRaw: float64(at.Unix()),
		}); err != nil {
			t.Fatalf("equal bounds must be accepted: %v", err)
		}
	})

	t.Run("More than one lookup attribute is refused", func(t *testing.T) {
		_, err := svc.lookupEventsCore(store, LookupEventsInput{
			LookupAttributes: []interface{}{
				map[string]interface{}{"AttributeKey": "EventName", "AttributeValue": "CreateTrail"},
				map[string]interface{}{"AttributeKey": "Username", "AttributeValue": "root"},
			},
		})
		code(t, err, "InvalidLookupAttributesException")
	})

	t.Run("A single lookup attribute is accepted", func(t *testing.T) {
		if _, err := svc.lookupEventsCore(store, LookupEventsInput{
			LookupAttributes: []interface{}{
				map[string]interface{}{"AttributeKey": "EventName", "AttributeValue": "CreateTrail"},
			},
		}); err != nil {
			t.Fatalf("single attribute must be accepted: %v", err)
		}
	})

	t.Run("EventCategory outside the enum is refused", func(t *testing.T) {
		_, err := svc.lookupEventsCore(store, LookupEventsInput{EventCategory: "management"})
		code(t, err, "InvalidEventCategoryException")
	})

	t.Run("EventCategory insight is accepted", func(t *testing.T) {
		if _, err := svc.lookupEventsCore(store, LookupEventsInput{EventCategory: "insight"}); err != nil {
			t.Fatalf("insight must be accepted: %v", err)
		}
	})

	t.Run("A token the store never issued is refused", func(t *testing.T) {
		_, err := svc.lookupEventsCore(store, LookupEventsInput{NextToken: "not-a-token"})
		code(t, err, "InvalidNextTokenException")
	})
}
