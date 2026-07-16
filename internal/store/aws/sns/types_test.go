package sns

import (
	"testing"
)

func TestNewTopic(t *testing.T) {
	topic := NewTopic("my-topic", "arn:aws:sns:us-east-1:123456789012:my-topic")

	if topic.Name != "my-topic" {
		t.Errorf("Name = %v, want 'my-topic'", topic.Name)
	}
	if topic.Arn != "arn:aws:sns:us-east-1:123456789012:my-topic" {
		t.Errorf("Arn = %v, want 'arn:aws:sns:us-east-1:123456789012:my-topic'", topic.Arn)
	}
	if topic.Attributes == nil {
		t.Error("Attributes = nil, want non-nil")
	}
}

func TestNewSubscription(t *testing.T) {
	sub := NewSubscription(
		"arn:aws:sns:us-east-1:123456789012:my-topic:subscription-id",
		"arn:aws:sns:us-east-1:123456789012:my-topic",
		"https",
		"https://example.com/hook",
	)

	if sub.SubscriptionArn != "arn:aws:sns:us-east-1:123456789012:my-topic:subscription-id" {
		t.Errorf("SubscriptionArn = %v, want 'arn:aws:sns:us-east-1:123456789012:my-topic:subscription-id'", sub.SubscriptionArn)
	}
	if sub.TopicArn != "arn:aws:sns:us-east-1:123456789012:my-topic" {
		t.Errorf("TopicArn = %v, want 'arn:aws:sns:us-east-1:123456789012:my-topic'", sub.TopicArn)
	}
	if sub.Protocol != "https" {
		t.Errorf("Protocol = %v, want 'https'", sub.Protocol)
	}
	if sub.Endpoint != "https://example.com/hook" {
		t.Errorf("Endpoint = %v, want 'https://example.com/hook'", sub.Endpoint)
	}
	if sub.Attributes == nil {
		t.Error("Attributes = nil, want non-nil")
	}
	if !sub.PendingConfirmation {
		t.Error("PendingConfirmation = false, want true")
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage("msg-id-123", "arn:aws:sns:us-east-1:123456789012:my-topic", "Hello World")

	if msg.MessageId != "msg-id-123" {
		t.Errorf("MessageId = %v, want 'msg-id-123'", msg.MessageId)
	}
	if msg.TopicArn != "arn:aws:sns:us-east-1:123456789012:my-topic" {
		t.Errorf("TopicArn = %v, want 'arn:aws:sns:us-east-1:123456789012:my-topic'", msg.TopicArn)
	}
	if msg.Message != "Hello World" {
		t.Errorf("Message = %v, want 'Hello World'", msg.Message)
	}
	if msg.MessageAttributes == nil {
		t.Error("MessageAttributes = nil, want non-nil")
	}
}
