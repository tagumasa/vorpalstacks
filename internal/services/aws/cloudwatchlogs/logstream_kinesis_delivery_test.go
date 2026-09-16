package cloudwatchlogs

import (
	"context"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
)

// regionRecordingKinesisInvoker captures the region every Kinesis call of a
// subscription-filter delivery addressed.
type regionRecordingKinesisInvoker struct {
	mu          sync.Mutex
	listRegions []string
	putRegions  []string
}

func (r *regionRecordingKinesisInvoker) ListShards(_ context.Context, region, _ string) ([]invokers.ShardInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listRegions = append(r.listRegions, region)
	return []invokers.ShardInfo{{ShardID: "shardId-000000000000"}}, nil
}

func (r *regionRecordingKinesisInvoker) PutRecord(_ context.Context, region, _ string, _ string, _ []byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.putRegions = append(r.putRegions, region)
	return "seq-1", nil
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

	svc.putToKinesis("arn:aws:kinesis:eu-west-1:000000000000:stream/logs-dest", []byte("compressed-payload"))

	if len(invoker.listRegions) != 1 || invoker.listRegions[0] != "eu-west-1" {
		t.Errorf("ListShards regions = %v, want exactly [eu-west-1] (the ARN's region)", invoker.listRegions)
	}
	if len(invoker.putRegions) != 1 || invoker.putRegions[0] != "eu-west-1" {
		t.Errorf("PutRecord regions = %v, want exactly [eu-west-1] (the ARN's region)", invoker.putRegions)
	}
}
