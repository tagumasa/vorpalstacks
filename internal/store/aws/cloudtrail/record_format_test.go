package cloudtrail

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestRecordFormatJSON pins the stored CloudTrailEvent JSON against the AWS
// CloudTrail record format: eventID/requestID spellings, boolean readOnly,
// awsRegion, record-shaped resources entries, and no storage-only keys at
// the top level.
func TestRecordFormatJSON(t *testing.T) {
	store := newLookupTestStore(t)

	identity := &UserIdentity{
		Type:        "IAMUser",
		PrincipalID: "AIDAJDPLRKLG7UEXAMPLE",
		ARN:         "arn:aws:iam::123456789012:user/Alice",
		AccountID:   "123456789012",
		UserName:    "Alice",
	}
	resources := []Resource{{
		ResourceType: "AWS::CloudTrail::Trail",
		ResourceName: "arn:aws:cloudtrail:us-east-1:123456789012:trail/example-trail",
	}}
	err := store.RecordServiceEvent(
		"CreateTrail", "cloudtrail.amazonaws.com", identity,
		"192.0.2.0", "AKIAIOSFODNN7EXAMPLE", "aws-cli/2.0",
		false, "", "",
		map[string]interface{}{"name": "example-trail"},
		map[string]interface{}{"trailARN": resources[0].ResourceName},
		resources,
	)
	if err != nil {
		t.Fatalf("record service event: %v", err)
	}

	events, _, err := store.LookupEvents(EventQuery{
		EventNames: []string{"CreateTrail"},
		ReadOnly:   "false",
		MaxResults: 10,
	})
	if err != nil {
		t.Fatalf("lookup events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("lookup returned %d events, want 1", len(events))
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(events[0].CloudTrailEvent), &record); err != nil {
		t.Fatalf("unmarshal CloudTrailEvent record: %v", err)
	}

	if got := record["eventID"]; got != events[0].EventID {
		t.Errorf("record eventID = %v, want the event's ID %s", got, events[0].EventID)
	}
	if got, ok := record["readOnly"].(bool); !ok || got {
		t.Errorf("record readOnly = %#v, want boolean false", record["readOnly"])
	}
	if got := record["eventVersion"]; got != "1.08" {
		t.Errorf("record eventVersion = %v, want 1.08", got)
	}
	if got := record["eventType"]; got != "AwsApiCall" {
		t.Errorf("record eventType = %v, want AwsApiCall", got)
	}
	if got := record["eventCategory"]; got != "Management" {
		t.Errorf("record eventCategory = %v, want Management", got)
	}
	if got, ok := record["managementEvent"].(bool); !ok || !got {
		t.Errorf("record managementEvent = %#v, want boolean true", record["managementEvent"])
	}
	if got := record["awsRegion"]; got != "us-east-1" {
		t.Errorf("record awsRegion = %v, want us-east-1", got)
	}
	if got := record["eventTime"].(string); !strings.HasSuffix(got, "Z") || strings.Contains(got, ".") {
		t.Errorf("record eventTime = %v, want whole-second UTC RFC3339", got)
	}
	if _, ok := record["eventId"]; ok {
		t.Error("record carries the storage spelling eventId, want eventID only")
	}
	if _, ok := record["requestId"]; ok {
		t.Error("record carries the storage spelling requestId, want requestID only")
	}
	for _, storageKey := range []string{"accessKeyId", "tags", "cloudTrailEvent"} {
		if _, ok := record[storageKey]; ok {
			t.Errorf("record carries storage-only key %q at the top level", storageKey)
		}
	}

	identityJSON, ok := record["userIdentity"].(map[string]interface{})
	if !ok {
		t.Fatalf("record userIdentity = %#v, want an object", record["userIdentity"])
	}
	if got := identityJSON["type"]; got != "IAMUser" {
		t.Errorf("record userIdentity.type = %v, want IAMUser", got)
	}
	if got := identityJSON["principalId"]; got != "AIDAJDPLRKLG7UEXAMPLE" {
		t.Errorf("record userIdentity.principalId = %v, want the user's real ID", got)
	}

	resourcesJSON, ok := record["resources"].([]interface{})
	if !ok || len(resourcesJSON) != 1 {
		t.Fatalf("record resources = %#v, want one entry", record["resources"])
	}
	entry, ok := resourcesJSON[0].(map[string]interface{})
	if !ok {
		t.Fatalf("record resources entry = %#v, want an object", resourcesJSON[0])
	}
	if got := entry["ARN"]; got != resources[0].ResourceName {
		t.Errorf("record resources entry ARN = %v, want %s", got, resources[0].ResourceName)
	}
	if got := entry["accountId"]; got != "123456789012" {
		t.Errorf("record resources entry accountId = %v, want 123456789012", got)
	}
	if got := entry["type"]; got != "AWS::CloudTrail::Trail" {
		t.Errorf("record resources entry type = %v, want AWS::CloudTrail::Trail", got)
	}
}

// TestRecordFormatReadOnlyNullResponse pins the record for a read-only
// call: boolean readOnly true and a null responseElements, per the record
// contents reference ("For readOnly APIs, this field is null").
func TestRecordFormatReadOnlyNullResponse(t *testing.T) {
	store := newLookupTestStore(t)

	err := store.RecordServiceEvent(
		"ListTrails", "cloudtrail.amazonaws.com", &UserIdentity{Type: "IAMUser", AccountID: "123456789012"},
		"192.0.2.0", "AKIAIOSFODNN7EXAMPLE", "aws-cli/2.0",
		true, "", "",
		map[string]interface{}{}, nil, nil,
	)
	if err != nil {
		t.Fatalf("record service event: %v", err)
	}

	events, _, err := store.LookupEvents(EventQuery{
		EventNames: []string{"ListTrails"},
		ReadOnly:   "true",
		MaxResults: 10,
	})
	if err != nil {
		t.Fatalf("lookup events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("lookup returned %d events, want 1", len(events))
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(events[0].CloudTrailEvent), &record); err != nil {
		t.Fatalf("unmarshal CloudTrailEvent record: %v", err)
	}
	if got, ok := record["readOnly"].(bool); !ok || !got {
		t.Errorf("record readOnly = %#v, want boolean true", record["readOnly"])
	}
	if got, present := record["responseElements"]; !present || got != nil {
		t.Errorf("record responseElements = %#v (present=%v), want present null", got, present)
	}
}

// TestSessionAttributesPersistenceRoundTrip pins the proto persistence
// boundary's zero-time guard for session attributes: a zero CreationDate
// crosses EventToProto/ProtoToEvent as zero (never as the negative
// UnixMilli encoding that decodes to year 1), and a populated CreationDate
// survives the round trip unchanged.
func TestSessionAttributesPersistenceRoundTrip(t *testing.T) {
	created := time.Date(2026, 9, 17, 12, 30, 5, 0, time.UTC)

	for _, tc := range []struct {
		name string
		in   time.Time
	}{
		{"populated creation date", created},
		{"zero creation date", time.Time{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			event := &Event{
				EventID:   "evt-session-" + tc.name,
				EventName: "AssumeRole",
				UserIdentity: &UserIdentity{
					Type: "AssumedRole",
					SessionContext: &SessionContext{
						Attributes: &SessionAttributes{
							CreationDate:     tc.in,
							MFAAuthenticated: "false",
						},
					},
				},
			}

			roundTripped := ProtoToEvent(EventToProto(event))
			got := roundTripped.UserIdentity.SessionContext.Attributes
			if got.MFAAuthenticated != "false" {
				t.Fatalf("MFAAuthenticated = %q, want %q", got.MFAAuthenticated, "false")
			}
			if !got.CreationDate.Equal(tc.in) {
				t.Fatalf("CreationDate = %v, want %v", got.CreationDate, tc.in)
			}
		})
	}
}
