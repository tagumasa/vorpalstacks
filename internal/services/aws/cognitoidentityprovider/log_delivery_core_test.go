package cognitoidentityprovider

import (
	"errors"
	"testing"
)

// The log-delivery configuration contract: the log level is exclusive to
// the event source (userNotification is ERROR, userAuthEvents is INFO),
// every configuration names exactly one destination, and the
// userNotification source goes to CloudWatch Logs alone.
func TestSetLogDeliveryConfigurationValidation(t *testing.T) {
	svc, reqCtx, _, pool := newTokensTestEnv(t)

	setConfig := func(lc map[string]interface{}) error {
		_, err := svc.setLogDeliveryConfigurationCore(reqCtx.GetRegion(), pool.ID, []interface{}{lc})
		return err
	}

	// The documented valid combinations.
	if err := setConfig(map[string]interface{}{
		"EventSource": "userAuthEvents", "LogLevel": "INFO",
		"CloudWatchLogsConfiguration": map[string]interface{}{"LogGroupArn": "arn:aws:logs:us-east-1:000000000000:log-group:auth"},
	}); err != nil {
		t.Fatalf("userAuthEvents to CloudWatch rejected: %v", err)
	}
	if err := setConfig(map[string]interface{}{
		"EventSource": "userAuthEvents", "LogLevel": "INFO",
		"S3Configuration": map[string]interface{}{"BucketArn": "arn:aws:s3:::auth-logs"},
	}); err != nil {
		t.Fatalf("userAuthEvents to S3 rejected: %v", err)
	}
	if err := setConfig(map[string]interface{}{
		"EventSource": "userNotification", "LogLevel": "ERROR",
		"CloudWatchLogsConfiguration": map[string]interface{}{"LogGroupArn": "arn:aws:logs:us-east-1:000000000000:log-group:notify"},
	}); err != nil {
		t.Fatalf("userNotification to CloudWatch rejected: %v", err)
	}

	// The log level is exclusive to the event source.
	if err := setConfig(map[string]interface{}{
		"EventSource": "userNotification", "LogLevel": "INFO",
		"CloudWatchLogsConfiguration": map[string]interface{}{"LogGroupArn": "arn:aws:logs:us-east-1:000000000000:log-group:notify"},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("userNotification at INFO returned %v, want InvalidParameter", err)
	}
	if err := setConfig(map[string]interface{}{
		"EventSource": "userAuthEvents", "LogLevel": "ERROR",
		"S3Configuration": map[string]interface{}{"BucketArn": "arn:aws:s3:::auth-logs"},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("userAuthEvents at ERROR returned %v, want InvalidParameter", err)
	}

	// A configuration without a destination, or with two, is rejected —
	// an accepted-but-dead configuration must not exist.
	if err := setConfig(map[string]interface{}{
		"EventSource": "userAuthEvents", "LogLevel": "INFO",
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("destination-less configuration returned %v, want InvalidParameter", err)
	}
	if err := setConfig(map[string]interface{}{
		"EventSource": "userAuthEvents", "LogLevel": "INFO",
		"S3Configuration":             map[string]interface{}{"BucketArn": "arn:aws:s3:::auth-logs"},
		"CloudWatchLogsConfiguration": map[string]interface{}{"LogGroupArn": "arn:aws:logs:us-east-1:000000000000:log-group:auth"},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("two-destination configuration returned %v, want InvalidParameter", err)
	}

	// The userNotification source is CloudWatch Logs alone.
	if err := setConfig(map[string]interface{}{
		"EventSource": "userNotification", "LogLevel": "ERROR",
		"S3Configuration": map[string]interface{}{"BucketArn": "arn:aws:s3:::notify-logs"},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("userNotification to S3 returned %v, want InvalidParameter", err)
	}
}

// The S3 destination's bucket name reads out of the bucket ARN's resource
// segment, and an unparseable ARN names no bucket.
func TestExtractBucketName(t *testing.T) {
	if got := extractBucketName("arn:aws:s3:::auth-logs"); got != "auth-logs" {
		t.Fatalf("bucket name = %q, want auth-logs", got)
	}
	if got := extractBucketName("not-an-arn"); got != "" {
		t.Fatalf("unparseable ARN yielded bucket %q", got)
	}
}
