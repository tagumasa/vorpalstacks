package cloudwatchlogs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
)

// regionRecordingKinesisInvoker captures the region and partition key of
// every Kinesis call a subscription-filter delivery made.
type regionRecordingKinesisInvoker struct {
	mu          sync.Mutex
	listRegions []string
	putRegions  []string
	putKeys     []string
}

func (r *regionRecordingKinesisInvoker) ListShards(_ context.Context, region, _ string) ([]invokers.ShardInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listRegions = append(r.listRegions, region)
	return []invokers.ShardInfo{{ShardID: "shardId-000000000000"}}, nil
}

func (r *regionRecordingKinesisInvoker) PutRecord(_ context.Context, region, _ string, partitionKey string, _ []byte) (string, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.putRegions = append(r.putRegions, region)
	r.putKeys = append(r.putKeys, partitionKey)
	return "seq-1", "", nil
}

func (r *regionRecordingKinesisInvoker) CreateShardIterator(context.Context, string, string, string, string, string, *time.Time) (string, error) {
	return "", nil
}

func (r *regionRecordingKinesisInvoker) GetRecords(context.Context, string, string, string, string, int32, bool) ([]invokers.KinesisRecord, string, error) {
	return nil, "", nil
}

func (r *regionRecordingKinesisInvoker) StreamExists(context.Context, string, string) (bool, error) {
	return true, nil
}

// TestPutToKinesisAddressesDestinationARNRegion pins that the whole
// subscription-filter delivery pipeline resolves the destination stream's
// region from the ARN: the shard probe and the record write must address the
// same regional store, or a cross-region destination fails the probe against
// the wrong region and the delivery vanishes.
func TestPutToKinesisAddressesDestinationARNRegion(t *testing.T) {
	invoker := &regionRecordingKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	svc := &LogsService{accountID: "000000000000", bus: bus}

	svc.putToKinesis("arn:aws:kinesis:eu-west-1:000000000000:stream/logs-dest", "CloudTrail/logs", "delivery-stream", []byte("compressed-payload"))

	if len(invoker.listRegions) != 1 || invoker.listRegions[0] != "eu-west-1" {
		t.Errorf("ListShards regions = %v, want exactly [eu-west-1] (the ARN's region)", invoker.listRegions)
	}
	if len(invoker.putRegions) != 1 || invoker.putRegions[0] != "eu-west-1" {
		t.Errorf("PutRecord regions = %v, want exactly [eu-west-1] (the ARN's region)", invoker.putRegions)
	}
}

// TestPutToKinesisPartitionKeyIsDerivedFromLogIdentity pins the delivery's
// partition key: a digest of the delivery's log-group/log-stream identity —
// never the probed shard ID, which placement ignores and which carries no
// identity a consumer could use — so retries keep their shard and distinct
// log streams spread across the destination as the documented by-log-stream
// distribution describes.
func TestPutToKinesisPartitionKeyIsDerivedFromLogIdentity(t *testing.T) {
	invoker := &regionRecordingKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	svc := &LogsService{accountID: "000000000000", bus: bus}

	svc.putToKinesis("arn:aws:kinesis:eu-west-1:000000000000:stream/logs-dest", "CloudTrail/logs", "delivery-stream", []byte("compressed-payload"))

	digest := sha256.Sum256([]byte("CloudTrail/logs\x00delivery-stream"))
	wantKey := hex.EncodeToString(digest[:])
	if len(invoker.putKeys) != 1 || invoker.putKeys[0] != wantKey {
		t.Fatalf("PutRecord partition keys = %v, want exactly the log-identity digest %q", invoker.putKeys, wantKey)
	}
	if invoker.putKeys[0] == "shardId-000000000000" {
		t.Fatalf("PutRecord partition key is the probed shard ID — placement ignores it and it carries no identity")
	}
}
