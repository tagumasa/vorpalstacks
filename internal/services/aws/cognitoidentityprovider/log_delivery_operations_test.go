package cognitoidentityprovider

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// The exported record must mirror the AWS user activity log sample:
// eventTimestamp as a string of epoch milliseconds, a message version,
// the user name and client ID, and the Java-style creationDate.
func TestFormatAuthEventLogMessageEnvelope(t *testing.T) {
	event := &cognitostore.AuthEvent{
		EventID:       "evt-1",
		UserName:      "alice",
		ClientID:      "client-7",
		UserPoolID:    "us-east-1_POOL",
		UserID:        "sub-1",
		EventType:     "SignIn",
		CreationDate:  time.Date(2024, 7, 17, 17, 25, 55, 0, time.UTC),
		EventResponse: "Pass",
		RiskDecision:  "NoRisk",
		RiskLevel:     "Low",
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(formatAuthEventLogMessage(event)), &record); err != nil {
		t.Fatalf("failed to parse the rendered log record: %v", err)
	}

	if ts, ok := record["eventTimestamp"].(string); !ok || ts != "1721237155000" {
		t.Fatalf("expected eventTimestamp to be the epoch-millis string 1721237155000, got %v", record["eventTimestamp"])
	}
	message, ok := record["message"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected a message object, got %T", record["message"])
	}
	if message["version"] != "1" {
		t.Fatalf("expected message version 1, got %v", message["version"])
	}
	if message["userName"] != "alice" {
		t.Fatalf("expected the recorded user name, got %v", message["userName"])
	}
	if message["clientId"] != "client-7" {
		t.Fatalf("expected the recorded client ID, got %v", message["clientId"])
	}
	creationDate, _ := message["creationDate"].(string)
	if !strings.HasSuffix(creationDate, " 2024") || !strings.Contains(creationDate, "UTC") {
		t.Fatalf("expected a Java-style UTC creationDate, got %q", creationDate)
	}
}
