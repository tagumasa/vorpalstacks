package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestServiceSpan(t *testing.T) {
	ctx := context.Background()
	_, span := ServiceSpan(ctx, "test-service", "test-operation")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestServiceSpan_WithAttributes(t *testing.T) {
	ctx := context.Background()
	_, span := ServiceSpan(ctx, "test-service", "test-operation",
		attribute.String("key", "value"))
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestS3Span(t *testing.T) {
	ctx := context.Background()
	_, span := S3Span(ctx, "GetObject", "test-bucket", "test-key")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestS3Span_EmptyBucket(t *testing.T) {
	ctx := context.Background()
	_, span := S3Span(ctx, "ListBuckets", "", "")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestLambdaSpan(t *testing.T) {
	ctx := context.Background()
	_, span := LambdaSpan(ctx, "Invoke", "test-function")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestLambdaSpan_EmptyFunction(t *testing.T) {
	ctx := context.Background()
	_, span := LambdaSpan(ctx, "ListFunctions", "")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestDynamoDBSpan(t *testing.T) {
	ctx := context.Background()
	_, span := DynamoDBSpan(ctx, "GetItem", "test-table")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestDynamoDBSpan_EmptyTable(t *testing.T) {
	ctx := context.Background()
	_, span := DynamoDBSpan(ctx, "ListTables", "")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestSNSSpan(t *testing.T) {
	ctx := context.Background()
	_, span := SNSSpan(ctx, "Publish", "arn:aws:sns:us-east-1:123456789012:test-topic")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestSNSSpan_EmptyTopic(t *testing.T) {
	ctx := context.Background()
	_, span := SNSSpan(ctx, "ListTopics", "")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestSQSSpan(t *testing.T) {
	ctx := context.Background()
	_, span := SQSSpan(ctx, "SendMessage", "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestSQSSpan_EmptyQueue(t *testing.T) {
	ctx := context.Background()
	_, span := SQSSpan(ctx, "ListQueues", "")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestSetSpanError(t *testing.T) {
	ctx := context.Background()
	_, span := StartSpan(ctx, "test-span")
	defer span.End()

	SetSpanError(span, nil)

	SetSpanError(span, &testError{msg: "test error"})
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestSetSpanSuccess(t *testing.T) {
	ctx := context.Background()
	_, span := StartSpan(ctx, "test-span")
	defer span.End()

	SetSpanSuccess(span)
}
