package sqs

import (
	"errors"
	"testing"
)

func TestParseRedrivePolicyDefaultMaxReceiveCount(t *testing.T) {
	rdp, err := ParseRedrivePolicy(`{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:dlq"}`)
	if err != nil {
		t.Fatalf("ParseRedrivePolicy: %v", err)
	}
	if rdp.MaxReceiveCount != DefaultMaxReceiveCount {
		t.Errorf("absent maxReceiveCount = %d, want default %d", rdp.MaxReceiveCount, DefaultMaxReceiveCount)
	}
	if rdp.DeadLetterTargetARN != "arn:aws:sqs:us-east-1:123456789012:dlq" {
		t.Errorf("DeadLetterTargetARN = %q", rdp.DeadLetterTargetARN)
	}
}

func TestParseRedrivePolicyExplicitMaxReceiveCount(t *testing.T) {
	rdp, err := ParseRedrivePolicy(`{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:dlq","maxReceiveCount":"5"}`)
	if err != nil {
		t.Fatalf("ParseRedrivePolicy: %v", err)
	}
	if rdp.MaxReceiveCount != 5 {
		t.Errorf("maxReceiveCount = %d, want 5", rdp.MaxReceiveCount)
	}
	rdp, err = ParseRedrivePolicy(`{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:dlq","maxReceiveCount":7}`)
	if err != nil {
		t.Fatalf("ParseRedrivePolicy: %v", err)
	}
	if rdp.MaxReceiveCount != 7 {
		t.Errorf("integer maxReceiveCount = %d, want 7", rdp.MaxReceiveCount)
	}
}

func TestValidateSSEExclusion(t *testing.T) {
	cases := []struct {
		name  string
		attrs map[string]string
		want  error
	}{
		{"both set", map[string]string{"KmsMasterKeyId": "alias/key", "SqsManagedSseEnabled": "true"}, ErrInvalidParameterValue},
		{"kms only", map[string]string{"KmsMasterKeyId": "alias/key"}, nil},
		{"sse only", map[string]string{"SqsManagedSseEnabled": "true"}, nil},
		{"neither", map[string]string{}, nil},
		{"empty kms is unset", map[string]string{"KmsMasterKeyId": "", "SqsManagedSseEnabled": "true"}, nil},
		{"sse false", map[string]string{"KmsMasterKeyId": "alias/key", "SqsManagedSseEnabled": "false"}, nil},
	}
	for _, tc := range cases {
		got := validateSSEExclusion(tc.attrs)
		if !errors.Is(got, tc.want) {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestValidateHighThroughputFifo(t *testing.T) {
	cases := []struct {
		name  string
		attrs map[string]string
		want  error
	}{
		{"per message group with messageGroup scope", map[string]string{"DeduplicationScope": "messageGroup", "FifoThroughputLimit": "perMessageGroupId"}, nil},
		{"per message group with queue scope", map[string]string{"DeduplicationScope": "queue", "FifoThroughputLimit": "perMessageGroupId"}, ErrInvalidParameterValue},
		{"per message group with absent scope", map[string]string{"FifoThroughputLimit": "perMessageGroupId"}, ErrInvalidParameterValue},
		{"per queue with any scope", map[string]string{"DeduplicationScope": "queue", "FifoThroughputLimit": "perQueue"}, nil},
	}
	for _, tc := range cases {
		got := validateHighThroughputFifo(tc.attrs)
		if !errors.Is(got, tc.want) {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestValidateDeduplicationScopeAcceptsDocumentedValues(t *testing.T) {
	for _, v := range []string{"messageGroup", "queue"} {
		if err := validateDeduplicationScope(v); err != nil {
			t.Errorf("documented value %q rejected: %v", v, err)
		}
	}
	if err := validateDeduplicationScope("queueMessageGroup"); err == nil {
		t.Errorf("non-existent value queueMessageGroup accepted")
	}
}

func TestValidateSQSActionListAcceptsAnyActionName(t *testing.T) {
	// The AWS API Reference documents Actions as "the name of any action or
	// `*`", so action names outside any fixed list must be accepted.
	if err := validateSQSActionList([]string{"CreateQueue", "DeleteQueue", "ListQueues"}); err != nil {
		t.Errorf("valid action names rejected: %v", err)
	}
	if err := validateSQSActionList([]string{"*"}); err != nil {
		t.Errorf("* rejected: %v", err)
	}
	if err := validateSQSActionList([]string{"not an action!"}); err == nil {
		t.Errorf("name with invalid characters accepted")
	}
}
