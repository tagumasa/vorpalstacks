package sqs

import (
	"strings"
	"testing"
)

// TestReceiveMessageCoreExplicitMaxRejected pins the MaxNumberOfMessages
// contract "Valid values: 1 to 10. Default: 1.": a member that is present
// on the wire with a value outside the range — including an explicit 0 —
// is rejected, and only an absent member falls back to the default of 1.
// The Go SDK cannot reach the explicit-0 path (its serialiser drops a
// zero-valued MaxNumberOfMessages), so the rejection is pinned here for
// both wire forms: the query-protocol string and the JSON number. The
// validation runs before any store access, so a nil store is safe.
func TestReceiveMessageCoreExplicitMaxRejected(t *testing.T) {
	s := &SQSService{}
	cases := map[string]interface{}{
		"explicit zero (query string form)": "0",
		"explicit zero (JSON number form)":  float64(0),
		"above the maximum":                 float64(11),
	}
	for name, value := range cases {
		_, err := s.receiveMessageCore(nil, ReceiveMessageInput{
			QueueURL: "http://localhost:50080/000000000000/queue-a",
			Parameters: map[string]interface{}{
				"MaxNumberOfMessages": value,
			},
		})
		if err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
			t.Fatalf("%s: expected InvalidParameterValue, got %v", name, err)
		}
	}
}

// TestReceiveMessageCoreNonIntegerMembersRejected pins the wire-type
// contract for the Integer members of ReceiveMessage: a member that is
// present but not an integer is a shape violation rejected with
// SerializationException (the awsJson1_0 protocol SQS serves), never
// silently treated as an omitted member falling back to a default. The
// rejection runs before any store access, so a nil store is safe.
func TestReceiveMessageCoreNonIntegerMembersRejected(t *testing.T) {
	s := &SQSService{}
	queueURL := "http://localhost:50080/000000000000/queue-a"
	for _, member := range []string{"MaxNumberOfMessages", "WaitTimeSeconds", "VisibilityTimeout"} {
		_, err := s.receiveMessageCore(nil, ReceiveMessageInput{
			QueueURL: queueURL,
			Parameters: map[string]interface{}{
				member: "abc",
			},
		})
		if err == nil || !strings.Contains(err.Error(), "SerializationException") {
			t.Fatalf("%s: expected SerializationException, got %v", member, err)
		}
	}

}
