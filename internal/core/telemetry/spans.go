// Package telemetry provides OpenTelemetry tracing and metrics for vorpalstacks.
package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ServiceSpan starts a new span for a service operation.
func ServiceSpan(ctx context.Context, service, operation string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return StartSpan(ctx, service+"."+operation,
		trace.WithAttributes(append([]attribute.KeyValue{
			attribute.String("service", service),
			attribute.String("operation", operation),
		}, attrs...)...),
	)
}

// S3Span starts a new span for an S3 operation.
func S3Span(ctx context.Context, operation, bucket, key string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		attribute.String("aws.service", "s3"),
		attribute.String("aws.operation", operation),
	}
	if bucket != "" {
		attrs = append(attrs, attribute.String("aws.s3.bucket", bucket))
	}
	if key != "" {
		attrs = append(attrs, attribute.String("aws.s3.key", key))
	}
	return ServiceSpan(ctx, "s3", operation, attrs...)
}

// LambdaSpan starts a new span for a Lambda operation.
func LambdaSpan(ctx context.Context, operation, functionName string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		attribute.String("aws.service", "lambda"),
		attribute.String("aws.operation", operation),
	}
	if functionName != "" {
		attrs = append(attrs, attribute.String("aws.lambda.function_name", functionName))
	}
	return ServiceSpan(ctx, "lambda", operation, attrs...)
}

// DynamoDBSpan starts a new span for a DynamoDB operation.
func DynamoDBSpan(ctx context.Context, operation, tableName string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		attribute.String("aws.service", "dynamodb"),
		attribute.String("aws.operation", operation),
	}
	if tableName != "" {
		attrs = append(attrs, attribute.String("aws.dynamodb.table_name", tableName))
	}
	return ServiceSpan(ctx, "dynamodb", operation, attrs...)
}

// SNSSpan starts a new span for an SNS operation.
func SNSSpan(ctx context.Context, operation, topicArn string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		attribute.String("aws.service", "sns"),
		attribute.String("aws.operation", operation),
	}
	if topicArn != "" {
		attrs = append(attrs, attribute.String("aws.sns.topic_arn", topicArn))
	}
	return ServiceSpan(ctx, "sns", operation, attrs...)
}

// SQSSpan starts a new span for an SQS operation.
func SQSSpan(ctx context.Context, operation, queueUrl string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		attribute.String("aws.service", "sqs"),
		attribute.String("aws.operation", operation),
	}
	if queueUrl != "" {
		attrs = append(attrs, attribute.String("aws.sqs.queue_url", queueUrl))
	}
	return ServiceSpan(ctx, "sqs", operation, attrs...)
}

// SetSpanError marks a span as having an error.
func SetSpanError(span trace.Span, err error) {
	if span.IsRecording() && err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
}

// SetSpanSuccess marks a span as successful.
func SetSpanSuccess(span trace.Span) {
	if span.IsRecording() {
		span.SetStatus(codes.Ok, "")
	}
}
