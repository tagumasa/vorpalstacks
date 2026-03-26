package sqs

import (
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	q := NewQueue("test-queue", "us-east-1", "123456789012")

	if q.Name != "test-queue" {
		t.Errorf("Name = %v, want test-queue", q.Name)
	}
	if q.Region != "us-east-1" {
		t.Errorf("Region = %v, want us-east-1", q.Region)
	}
	if q.AccountID != "123456789012" {
		t.Errorf("AccountID = %v, want 123456789012", q.AccountID)
	}
	if q.VisibilityTimeout != 30 {
		t.Errorf("VisibilityTimeout = %v, want 30", q.VisibilityTimeout)
	}
	if q.MaximumMessageSize != 262144 {
		t.Errorf("MaximumMessageSize = %v, want 262144", q.MaximumMessageSize)
	}
	if q.MessageRetentionPeriod != 345600 {
		t.Errorf("MessageRetentionPeriod = %v, want 345600", q.MessageRetentionPeriod)
	}
	if q.DelaySeconds != 0 {
		t.Errorf("DelaySeconds = %v, want 0", q.DelaySeconds)
	}
	if q.ReceiveMessageWaitTimeSeconds != 0 {
		t.Errorf("ReceiveMessageWaitTimeSeconds = %v, want 0", q.ReceiveMessageWaitTimeSeconds)
	}
	if q.Attributes == nil {
		t.Error("Attributes = nil, want non-nil")
	}
	if q.Tags == nil {
		t.Error("Tags = nil, want non-nil")
	}
	if q.Permissions == nil {
		t.Error("Permissions = nil, want non-nil")
	}
}

func TestNewQueueDefaultValues(t *testing.T) {
	before := time.Now().UTC()
	q := NewQueue("test-queue", "us-east-1", "123456789012")
	after := time.Now().UTC()

	if q.CreatedTimestamp.Before(before) || q.CreatedTimestamp.After(after) {
		t.Errorf("CreatedTimestamp not in expected range")
	}
	if q.LastModifiedTimestamp.Before(before) || q.LastModifiedTimestamp.After(after) {
		t.Errorf("LastModifiedTimestamp not in expected range")
	}
}

func TestNewMessage(t *testing.T) {
	m := NewMessage("test body")

	if m.Body != "test body" {
		t.Errorf("Body = %v, want 'test body'", m.Body)
	}
	if m.MD5OfBody == "" {
		t.Error("MD5OfBody = empty, want non-empty")
	}
	if m.SentTimestamp.IsZero() {
		t.Error("SentTimestamp = zero, want non-zero")
	}
	if m.MessageAttributes == nil {
		t.Error("MessageAttributes = nil, want non-nil")
	}
	if m.Attributes == nil {
		t.Error("Attributes = nil, want non-nil")
	}
}

func TestNewMessageEmptyBody(t *testing.T) {
	m := NewMessage("")

	if m.Body != "" {
		t.Errorf("Body = %v, want ''", m.Body)
	}
	if m.MD5OfBody != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("MD5OfBody = %v, want 'd41d8cd98f00b204e9800998ecf8427e'", m.MD5OfBody)
	}
}

func TestMessageAttributeValueDataTypes(t *testing.T) {
	strValue := "test"
	mav := &MessageAttributeValue{
		StringValue: &strValue,
		DataType:    "String",
	}

	if mav.StringValue == nil || *mav.StringValue != "test" {
		t.Errorf("StringValue = %v, want 'test'", mav.StringValue)
	}
	if mav.DataType != "String" {
		t.Errorf("DataType = %v, want 'String'", mav.DataType)
	}
}
