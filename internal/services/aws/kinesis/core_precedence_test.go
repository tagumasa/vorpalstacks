package kinesis

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
)

// failingStorageReqCtx builds a request context whose region storage can
// never be created: the manager's base path is a regular file, so every
// store acquisition fails with an infrastructure error. Operations that
// validate their parameters before acquiring the store must therefore
// surface the parameter error, not the storage failure.
func failingStorageReqCtx(t *testing.T) *request.RequestContext {
	t.Helper()
	base := filepath.Join(t.TempDir(), "occupied")
	if err := os.WriteFile(base, []byte("occupied"), 0o600); err != nil {
		t.Fatalf("write blocker file: %v", err)
	}
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: base})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	return request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
}

func parsedRequest(params map[string]interface{}) *request.ParsedRequest {
	if params == nil {
		params = map[string]interface{}{}
	}
	return &request.ParsedRequest{Operation: "Test", Parameters: params}
}

// TestValidationPrecedesStoreAcquisition pins the failure precedence of the
// operations that validate their parameters ahead of the store acquisition:
// a request that is simultaneously invalid and undeliverable to storage
// must fail with InvalidArgumentException, never with the storage error.
func TestValidationPrecedesStoreAcquisition(t *testing.T) {
	svc := NewKinesisService("000000000000")
	reqCtx := failingStorageReqCtx(t)
	ctx := context.Background()

	cases := []struct {
		name   string
		invoke func() error
	}{
		{"UpdateStreamMode", func() error {
			_, err := svc.UpdateStreamMode(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"GetResourcePolicy", func() error {
			_, err := svc.GetResourcePolicy(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"ResourceARN": "not-an-arn",
			}))
			return err
		}},
		{"PutResourcePolicy", func() error {
			_, err := svc.PutResourcePolicy(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"ResourceARN": "not-an-arn",
			}))
			return err
		}},
		{"DeleteResourcePolicy", func() error {
			_, err := svc.DeleteResourcePolicy(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"ResourceARN": "not-an-arn",
			}))
			return err
		}},
		{"UpdateMaxRecordSize", func() error {
			_, err := svc.UpdateMaxRecordSize(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"UpdateStreamWarmThroughput", func() error {
			_, err := svc.UpdateStreamWarmThroughput(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"WarmThroughputMiBps": float64(-1),
			}))
			return err
		}},
		{"RegisterStreamConsumer", func() error {
			_, err := svc.RegisterStreamConsumer(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"DeregisterStreamConsumer", func() error {
			_, err := svc.DeregisterStreamConsumer(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"DescribeStreamConsumer", func() error {
			_, err := svc.DescribeStreamConsumer(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"ListStreamConsumers", func() error {
			_, err := svc.ListStreamConsumers(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"SubscribeToShard", func() error {
			_, err := svc.SubscribeToShard(ctx, reqCtx, parsedRequest(nil))
			return err
		}},
		{"GetRecords", func() error {
			_, err := svc.GetRecords(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"ShardIterator": "iterator",
				"Limit":         float64(0),
			}))
			return err
		}},
		// The migrated store-parameter cores: each validates its members
		// before the store acquisition, so the same precedence guarantee
		// now holds for the whole operation set.
		{"CreateStream", func() error {
			_, err := svc.CreateStream(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "invalid stream name!",
			}))
			return err
		}},
		{"ListStreams", func() error {
			_, err := svc.ListStreams(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"Limit": float64(0),
			}))
			return err
		}},
		{"ListShards", func() error {
			_, err := svc.ListShards(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"MaxResults": float64(0),
			}))
			return err
		}},
		{"SplitShard", func() error {
			_, err := svc.SplitShard(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "split",
			}))
			return err
		}},
		{"MergeShards", func() error {
			_, err := svc.MergeShards(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName":   "merge",
				"ShardToMerge": "shardId-000000000000",
			}))
			return err
		}},
		{"UpdateShardCount", func() error {
			_, err := svc.UpdateShardCount(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "reshape",
			}))
			return err
		}},
		{"PutRecord", func() error {
			_, err := svc.PutRecord(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName":   "put",
				"PartitionKey": "",
			}))
			return err
		}},
		{"PutRecords", func() error {
			_, err := svc.PutRecords(ctx, reqCtx, parsedRequest(map[string]interface{}{}))
			return err
		}},
		{"GetShardIterator", func() error {
			_, err := svc.GetShardIterator(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName":        "read",
				"ShardIteratorType": "LATEST",
			}))
			return err
		}},
		{"IncreaseStreamRetentionPeriod", func() error {
			_, err := svc.IncreaseStreamRetentionPeriod(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"RetentionPeriodHours": float64(0),
			}))
			return err
		}},
		{"DecreaseStreamRetentionPeriod", func() error {
			_, err := svc.DecreaseStreamRetentionPeriod(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"RetentionPeriodHours": float64(0),
			}))
			return err
		}},
		{"StartStreamEncryption", func() error {
			_, err := svc.StartStreamEncryption(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"EncryptionType": "NONE",
			}))
			return err
		}},
		{"StopStreamEncryption", func() error {
			_, err := svc.StopStreamEncryption(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "stop-enc",
			}))
			return err
		}},
		{"UpdateAccountSettings", func() error {
			_, err := svc.UpdateAccountSettings(ctx, reqCtx, parsedRequest(map[string]interface{}{}))
			return err
		}},
		{"EnableEnhancedMonitoring", func() error {
			_, err := svc.EnableEnhancedMonitoring(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName":        "monitor",
				"ShardLevelMetrics": []interface{}{"NotAMetric"},
			}))
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.invoke(); !errors.Is(err, ErrInvalidArgument) {
				t.Fatalf("expected InvalidArgumentException, got: %v", err)
			}
		})
	}
}
