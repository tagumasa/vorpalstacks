package eventbridge

import (
	"encoding/json"
	"testing"
)

// vocabTestEventJSON renders the documentation's EC2 state-change event as
// the TestEventPattern Event payload (all seven mandatory members present).
func vocabTestEventJSON() string {
	return `{
		"version": "0",
		"id": "7bf73129-1428-4cd3-a780-95db273d1602",
		"detail-type": "EC2 Instance State-change Notification",
		"source": "aws.ec2",
		"account": "123456789012",
		"time": "2015-11-11T21:29:54Z",
		"region": "us-east-1",
		"resources": ["arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111"],
		"detail": {"instance-id": "i-abcd1111", "state": "pending"}
	}`
}

// eventJSONWithoutMember re-renders the test event with one member removed.
func eventJSONWithoutMember(t *testing.T, member string) string {
	t.Helper()
	var eventMap map[string]interface{}
	if err := json.Unmarshal([]byte(vocabTestEventJSON()), &eventMap); err != nil {
		t.Fatal(err)
	}
	delete(eventMap, member)
	without, err := json.Marshal(eventMap)
	if err != nil {
		t.Fatal(err)
	}
	return string(without)
}

// TestEventPatternCoreContract pins the TestEventPattern plane: the same
// structural validation as every acceptance site, the seven mandatory event
// members from the model's Event documentation, and evaluation through the
// shared level matcher ($or included).
func TestEventPatternCoreContract(t *testing.T) {
	svc := &EventsService{}

	result, err := svc.testEventPatternCore(
		`{"detail": {"state": ["pending"]}}`, vocabTestEventJSON())
	if err != nil || result != true {
		t.Fatalf("a matching pattern must answer true, got result=%v err=%v", result, err)
	}

	result, err = svc.testEventPatternCore(
		`{"detail": {"state": ["terminated"]}}`, vocabTestEventJSON())
	if err != nil || result != false {
		t.Fatalf("a non-matching pattern must answer false, got result=%v err=%v", result, err)
	}

	// $or through the TestEventPattern plane shares the delivery matcher's
	// semantics.
	result, err = svc.testEventPatternCore(
		`{"$or": [{"detail": {"state": ["terminated"]}}, {"source": ["aws.ec2"]}]}`,
		vocabTestEventJSON())
	if err != nil || result != true {
		t.Fatalf("$or second disjunct must match, got result=%v err=%v", result, err)
	}

	// An invalid pattern is rejected with the same error shape as PutRule,
	// CreateArchive and UpdateArchive.
	if _, err := svc.testEventPatternCore(`{"detail": {"state": [{"regexp": "x"}]}}`, vocabTestEventJSON()); err == nil {
		t.Fatal("an unknown operator must be rejected")
	}

	// Every mandatory member is required; one omission is enough.
	for _, member := range mandatoryEventMembers {
		if _, err := svc.testEventPatternCore(`{"source": ["aws.ec2"]}`, eventJSONWithoutMember(t, member)); err == nil {
			t.Fatalf("an event missing mandatory member %q must be rejected", member)
		}
	}
}
