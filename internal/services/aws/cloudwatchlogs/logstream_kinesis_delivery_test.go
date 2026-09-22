package cloudwatchlogs

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
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
	putStreams  []string
	putKeys     []string
	putData     [][]byte
}

func (r *regionRecordingKinesisInvoker) ListShards(_ context.Context, region, _ string) ([]invokers.ShardInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listRegions = append(r.listRegions, region)
	return []invokers.ShardInfo{{ShardID: "shardId-000000000000"}}, nil
}

func (r *regionRecordingKinesisInvoker) PutRecord(_ context.Context, region, streamName, partitionKey string, data []byte) (string, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.putRegions = append(r.putRegions, region)
	r.putStreams = append(r.putStreams, streamName)
	r.putKeys = append(r.putKeys, partitionKey)
	r.putData = append(r.putData, data)
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

// putCount reports the number of captured PutRecord calls under the
// recorder's mutex, for tests that poll delivery progress.
func (r *regionRecordingKinesisInvoker) putCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.putRegions)
}

// snapshot returns a copy of the captured PutRecord calls under the
// recorder's mutex.
func (r *regionRecordingKinesisInvoker) snapshot() (regions, streams, keys []string, data [][]byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string{}, r.putRegions...), append([]string{}, r.putStreams...),
		append([]string{}, r.putKeys...), append([][]byte{}, r.putData...)
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

	svc.putToKinesis("arn:aws:kinesis:eu-west-1:000000000000:stream/logs-dest", "CloudTrail/logs", "delivery-stream", "ByLogStream", []byte("compressed-payload"))

	regions, _, _, _ := invoker.snapshot()
	if len(regions) != 1 || len(invoker.listRegions) != 1 || invoker.listRegions[0] != "eu-west-1" {
		t.Errorf("regions = list:%v put:%v, want exactly [eu-west-1] (the ARN's region)", invoker.listRegions, regions)
	}
}

// TestPutToKinesisRecordDataIsRawGzipPayload pins the delivered blob
// against the documented consumer pipeline: the CloudWatch Logs User
// Guide's Kinesis example examines a record with `base64 -d | zcat` over
// the Kinesis API's base64 presentation — one gzip layer of the
// DATA_MESSAGE JSON, no JSON envelope and no additional base64 wrap (those
// belong to the Lambda invocation format alone). The invoker receives the
// wire form (the base64 text the platform stores for every SDK-written
// record); decoding it once must yield the gzip, and one gunzip the
// payload.
func TestPutToKinesisRecordDataIsRawGzipPayload(t *testing.T) {
	invoker := &regionRecordingKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	svc := &LogsService{accountID: "000000000000", bus: bus}

	payload := map[string]interface{}{
		"owner":       "000000000000",
		"logGroup":    "CloudTrail/logs",
		"logStream":   "delivery-stream",
		"messageType": "DATA_MESSAGE",
		"logEvents":   []map[string]interface{}{{"id": "1", "timestamp": 1432826855000, "message": "{}"}},
	}
	compressed, err := compressJSON(payload)
	if err != nil {
		t.Fatalf("compressJSON: %v", err)
	}

	svc.putToKinesis("arn:aws:kinesis:eu-west-1:000000000000:stream/logs-dest", "CloudTrail/logs", "delivery-stream", "ByLogStream", compressed)

	_, _, _, data := invoker.snapshot()
	if len(data) != 1 {
		t.Fatalf("expected exactly one PutRecord, got %d", len(data))
	}
	blob, err := base64.StdEncoding.DecodeString(string(data[0]))
	if err != nil {
		t.Fatalf("the delivered wire form must be the record blob's base64 presentation: %v", err)
	}
	if !bytes.Equal(blob, compressed) {
		t.Fatalf("the decoded blob must be the gzipped payload as-is, got %d bytes", len(blob))
	}
	var roundTrip map[string]interface{}
	if err := gunzipJSON(blob, &roundTrip); err != nil {
		t.Fatalf("one gunzip must yield the payload JSON (no extra base64/JSON layer): %v", err)
	}
	if roundTrip["messageType"] != "DATA_MESSAGE" {
		t.Fatalf("gunzipped record must be the DATA_MESSAGE payload, got messageType %v", roundTrip["messageType"])
	}
}

// TestPutToKinesisPartitionKeyFollowsDistribution pins the partition-key
// semantics of the filter's distribution: ByLogStream (the documented
// default, and the empty value a stored filter predating the member may
// carry) derives a stable digest of the delivery's log-group/log-stream
// identity, while Random takes a fresh key per delivery so hash placement
// spreads the records — the guide's recommendation against single-shard
// throttling.
func TestPutToKinesisPartitionKeyFollowsDistribution(t *testing.T) {
	invoker := &regionRecordingKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	svc := &LogsService{accountID: "000000000000", bus: bus}
	const dest = "arn:aws:kinesis:eu-west-1:000000000000:stream/logs-dest"

	svc.putToKinesis(dest, "CloudTrail/logs", "delivery-stream", "ByLogStream", []byte("p"))
	svc.putToKinesis(dest, "CloudTrail/logs", "delivery-stream", "", []byte("p"))
	digest := sha256.Sum256([]byte("CloudTrail/logs\x00delivery-stream"))
	wantKey := hex.EncodeToString(digest[:])

	svc.putToKinesis(dest, "CloudTrail/logs", "delivery-stream", "Random", []byte("p"))
	svc.putToKinesis(dest, "CloudTrail/logs", "delivery-stream", "Random", []byte("p"))

	_, _, keys, _ := invoker.snapshot()
	if len(keys) != 4 {
		t.Fatalf("expected four PutRecord calls, got %d", len(keys))
	}
	if keys[0] != wantKey || keys[1] != wantKey {
		t.Fatalf("ByLogStream and the empty default must both derive the log-identity digest %q, got %q and %q", wantKey, keys[0], keys[1])
	}
	if keys[2] == wantKey || keys[3] == wantKey {
		t.Fatalf("Random must not fall back to the log-identity digest, got %q", keys[2])
	}
	if keys[2] == keys[3] {
		t.Fatalf("Random must take a fresh key per delivery, got the same key twice: %q", keys[2])
	}
}
