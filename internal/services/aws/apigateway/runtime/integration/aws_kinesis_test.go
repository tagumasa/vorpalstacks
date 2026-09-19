package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
)

// shardReportingKinesisInvoker answers a PutRecord receipt whose shard no
// hardcoded response member could guess.
type shardReportingKinesisInvoker struct {
	shardID string
}

func (f *shardReportingKinesisInvoker) ListShards(context.Context, string, string) ([]invokers.ShardInfo, error) {
	return []invokers.ShardInfo{{ShardID: f.shardID}}, nil
}

func (f *shardReportingKinesisInvoker) PutRecord(context.Context, string, string, string, []byte) (string, string, error) {
	return "496277365171631905666581120543900666567083405037136334941402517362", f.shardID, nil
}

func (f *shardReportingKinesisInvoker) CreateShardIterator(context.Context, string, string, string, string, string, *time.Time) (string, error) {
	return "", nil
}

func (f *shardReportingKinesisInvoker) GetRecords(context.Context, string, string, string, string, int32, bool) ([]invokers.KinesisRecord, string, error) {
	return nil, "", nil
}

func (f *shardReportingKinesisInvoker) StreamExists(context.Context, string, string) (bool, error) {
	return true, nil
}

// TestKinesisPutRecordResponseReportsReceivingShard pins the PutRecord
// action's response vocabulary: the Kinesis model's PutRecordOutput carries
// SequenceNumber and ShardId, and ShardId is the shard the record actually
// landed on — hash-based placement can move any record off shard zero, so a
// fixed shard-zero member is wrong by construction.
func TestKinesisPutRecordResponseReportsReceivingShard(t *testing.T) {
	invoker := &shardReportingKinesisInvoker{shardID: "shardId-000000000042"}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	executor := NewAWSExecutor(bus, "000000000000", "us-east-1")

	resp, err := executor.Execute(context.Background(), &IntegrationRequest{
		URI:     "arn:aws:apigateway:us-east-1:kinesis:action/PutRecord/pinned-stream",
		Headers: map[string]string{"PartitionKey": "pinned-key"},
		Body:    []byte(`{"message":"pinned"}`),
	})
	if err != nil {
		t.Fatalf("execute Kinesis PutRecord action: %v", err)
	}

	var parsed struct {
		PutRecordResponse struct {
			PutRecordResult struct {
				SequenceNumber string `json:"SequenceNumber"`
				ShardId        string `json:"ShardId"`
			} `json:"PutRecordResult"`
		} `json:"PutRecordResponse"`
	}
	if err := json.Unmarshal(resp.Body, &parsed); err != nil {
		t.Fatalf("parse PutRecord response %s: %v", resp.Body, err)
	}
	result := parsed.PutRecordResponse.PutRecordResult
	if result.ShardId != "shardId-000000000042" {
		t.Fatalf("PutRecord response ShardId = %q, want the shard the invoker reported (shardId-000000000042)", result.ShardId)
	}
	if result.SequenceNumber == "" {
		t.Fatalf("PutRecord response SequenceNumber empty, want the invoker's receipt")
	}
}
