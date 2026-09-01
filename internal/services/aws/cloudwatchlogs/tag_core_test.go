package cloudwatchlogs

import (
	"context"
	"errors"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// newTagConfigTestStore opens a throwaway CloudWatch Logs store seeded with
// one log group.
func newTagConfigTestStore(t *testing.T) *logsstore.Store {
	t.Helper()
	raw, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { raw.Close() })
	store, err := logsstore.NewStore(raw, raw.Bucket("logs-us-east-1"), "000000000000", "us-east-1", t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	lg := logsstore.NewLogGroup("app-logs", "us-east-1", "000000000000")
	if err := store.CreateLogGroup(lg); err != nil {
		t.Fatalf("seed log group: %v", err)
	}
	return store
}

// TestTagHandlerConfigValidateResource pins the tag-target resolution behind
// the CloudWatch Logs tag trio: a real log-group ARN resolves, a nonexistent
// or non-log-group ARN fails with the modelled ResourceNotFoundException.
func TestTagHandlerConfigValidateResource(t *testing.T) {
	svc := &LogsService{}
	store := newTagConfigTestStore(t)
	cfg := svc.tagHandlerConfig(store)
	ctx := context.Background()

	t.Run("existing log group resolves", func(t *testing.T) {
		for _, arn := range []string{
			// DescribeLogGroups' arn field carries the log-stream
			// namespace suffix; both forms address the group itself.
			"arn:aws:logs:us-east-1:000000000000:log-group:app-logs",
			"arn:aws:logs:us-east-1:000000000000:log-group:app-logs:*",
		} {
			if err := cfg.ValidateResource(ctx, arn); err != nil {
				t.Errorf("ValidateResource(%q) = %v, want nil", arn, err)
			}
		}
	})

	t.Run("nonexistent log group is rejected", func(t *testing.T) {
		for _, arn := range []string{
			"arn:aws:logs:us-east-1:000000000000:log-group:no-such-group",
			"arn:aws:logs:us-east-1:000000000000:destination:whatever",
		} {
			err := cfg.ValidateResource(ctx, arn)
			if err == nil {
				t.Fatalf("ValidateResource(%q) = nil, want not-found", arn)
			}
			var awsErr *awserrors.AWSError
			if !errors.As(err, &awsErr) {
				t.Fatalf("ValidateResource(%q) = %T, want *awserrors.AWSError", arn, err)
			}
			if awsErr.Code != "ResourceNotFoundException" {
				t.Errorf("%s: code = %q, want ResourceNotFoundException", arn, awsErr.Code)
			}
		}
	})
}
