package lambda

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/eventbus"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// scriptedDynamoDBStream mirrors the store semantics the poller relies on:
// sequence numbers are assigned from one upwards, GetLatestSequence answers
// the sequence of the newest record (zero for an empty stream), and
// GetRecords reads strictly after the requested position.
type scriptedDynamoDBStream struct {
	mu      sync.Mutex
	records []eventbus.DynamoDBStreamRecord
}

func (s *scriptedDynamoDBStream) publish(seq int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append(s.records, eventbus.DynamoDBStreamRecord{
		EventID:      fmt.Sprintf("id-%d", seq),
		EventName:    "INSERT",
		EventVersion: "1.1",
		EventSource:  "aws:dynamodb",
		AWSRegion:    "us-east-1",
		Dynamodb: map[string]interface{}{
			"SequenceNumber":              strconv.FormatInt(seq, 10),
			"ApproximateCreationDateTime": float64(time.Now().Unix()),
			"Keys":                        map[string]interface{}{"id": map[string]interface{}{"S": fmt.Sprintf("key-%d", seq)}},
			"NewImage":                    map[string]interface{}{"id": map[string]interface{}{"S": fmt.Sprintf("key-%d", seq)}},
		},
		EventSourceARN: "arn:aws:dynamodb:us-east-1:123456789012:table/ddb-anchor-table/stream/label",
	})
}

func (s *scriptedDynamoDBStream) recordSeq(r eventbus.DynamoDBStreamRecord) int64 {
	seq, _ := r.Dynamodb["SequenceNumber"].(string)
	v, _ := strconv.ParseInt(seq, 10, 64)
	return v
}

func (s *scriptedDynamoDBStream) GetRecords(_ context.Context, _, _ string, fromSeq int64, limit int) ([]eventbus.DynamoDBStreamRecord, int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []eventbus.DynamoDBStreamRecord
	for _, r := range s.records {
		if s.recordSeq(r) <= fromSeq {
			continue
		}
		out = append(out, r)
		if len(out) >= limit {
			break
		}
	}
	next := fromSeq
	if len(out) > 0 {
		next = s.recordSeq(out[len(out)-1])
	}
	return out, next, nil
}

func (s *scriptedDynamoDBStream) GetLatestSequence(context.Context, string, string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.records) == 0 {
		return 0, nil
	}
	return s.recordSeq(s.records[len(s.records)-1]), nil
}

func (s *scriptedDynamoDBStream) ShardIDForStream(string) string { return "shard-0001" }

// fakeDDBStreamsBus exposes only the DynamoDB Streams invoker; the poller
// calls nothing else on the bus for a mapping without failure destinations.
type fakeDDBStreamsBus struct {
	eventbus.Bus
	invoker *scriptedDynamoDBStream
}

func (b *fakeDDBStreamsBus) DynamoDBStreamsInvoker() eventbus.DynamoDBStreamsInvoker {
	return b.invoker
}

func newDDBAnchorPoller(t *testing.T, startingPosition string) (*esmPoller, *scriptedDynamoDBStream, *lambdastore.EventSourceMapping, *capturePayloads) {
	t.Helper()
	const tableARN = "arn:aws:dynamodb:us-east-1:123456789012:table/ddb-anchor-table/stream/label"
	const funcARN = "arn:aws:lambda:us-east-1:123456789012:function:ddb-anchor-fn"

	esmStore := lambdastore.NewEventSourceStore(&memStorage{bucket: newMemBucket()}, "123456789012", "us-east-1")
	created, err := esmStore.Create(&lambdastore.EventSourceMapping{
		EventSourceArn:   tableARN,
		FunctionArn:      funcARN,
		StartingPosition: startingPosition,
		BatchSize:        10,
		State:            "Enabled",
	})
	if err != nil {
		t.Fatalf("create mapping: %v", err)
	}

	stream := &scriptedDynamoDBStream{}
	invoke := &capturePayloads{}
	p := &esmPoller{
		bus:       &fakeDDBStreamsBus{invoker: stream},
		invoke:    invoke.invoke,
		kinesisCP: make(map[string]string),
		windows:   make(map[string]*shardWindow),
		buffers:   make(map[string]*streamBuffer),
		esmStore:  esmStore,
		lambdaSvc: &LambdaService{accountID: "123456789012"},
	}
	return p, stream, created, invoke
}

// TestDynamoDBLatestAnchorSkipsHistory pins the LATEST start position: the
// anchor is the stream's latest record at the first poll, history before it
// is not delivered, and the anchor record itself is (the activation-lag
// semantics the Kinesis path applies).
func TestDynamoDBLatestAnchorSkipsHistory(t *testing.T) {
	p, stream, mapping, invoke := newDDBAnchorPoller(t, "LATEST")

	// Records that predate the mapping must be skipped by LATEST.
	stream.publish(1)
	stream.publish(2)
	// The anchor itself: the newest record at the first poll.
	stream.publish(3)

	p.processDynamoDBStreamsMapping(context.Background(), mapping)

	allPayloads := strings.Join(invoke.snapshot(), "\n")
	if !strings.Contains(allPayloads, `"SequenceNumber":"3"`) {
		t.Fatalf("the anchor record must be delivered, got payloads: %s", allPayloads)
	}
	for _, seq := range []string{`"SequenceNumber":"1"`, `"SequenceNumber":"2"`} {
		if strings.Contains(allPayloads, seq) {
			t.Fatalf("records predating the mapping must not be delivered under LATEST, got payloads: %s", allPayloads)
		}
	}
	if got := p.kinesisCP["ddb:"+mapping.UUID]; got != "3" {
		t.Fatalf("checkpoint must advance past the anchor record, got %q", got)
	}

	// Records arriving after the anchor are delivered on the next cycle.
	stream.publish(4)
	p.processDynamoDBStreamsMapping(context.Background(), mapping)
	allPayloads = strings.Join(invoke.snapshot(), "\n")
	if !strings.Contains(allPayloads, `"SequenceNumber":"4"`) {
		t.Fatalf("records after the anchor must be delivered, got payloads: %s", allPayloads)
	}
}

// TestDynamoDBLatestAnchorPersistedOnEmptyStream pins the empty-stream
// anchoring: the first poll durably records the anchor so a burst arriving
// between polls is read in full instead of only its newest record.
func TestDynamoDBLatestAnchorPersistedOnEmptyStream(t *testing.T) {
	p, stream, mapping, invoke := newDDBAnchorPoller(t, "LATEST")

	p.processDynamoDBStreamsMapping(context.Background(), mapping)
	if got := p.kinesisCP["ddb:"+mapping.UUID]; got != "0" {
		t.Fatalf("after the first empty poll the anchor must be persisted at zero, got checkpoint %q", got)
	}
	if invoke.callCount() != 0 {
		t.Fatalf("an empty stream must not invoke the function, got %d invokes", invoke.callCount())
	}

	// A burst arrives between polls.
	stream.publish(1)
	stream.publish(2)
	stream.publish(3)

	p.processDynamoDBStreamsMapping(context.Background(), mapping)
	allPayloads := strings.Join(invoke.snapshot(), "\n")
	for _, seq := range []string{"1", "2", "3"} {
		if !strings.Contains(allPayloads, fmt.Sprintf(`"SequenceNumber":%q`, seq)) {
			t.Fatalf("invocation payloads must carry record %q of the burst, got %s", seq, allPayloads)
		}
	}
	if got := p.kinesisCP["ddb:"+mapping.UUID]; got != "3" {
		t.Fatalf("checkpoint must advance past the burst, got %q", got)
	}
}

// TestDynamoDBTrimHorizonReadsHistory pins the default start position: with
// no durable checkpoint, TRIM_HORIZON reads the stream from the beginning.
func TestDynamoDBTrimHorizonReadsHistory(t *testing.T) {
	p, stream, mapping, invoke := newDDBAnchorPoller(t, "TRIM_HORIZON")

	stream.publish(1)
	stream.publish(2)

	p.processDynamoDBStreamsMapping(context.Background(), mapping)

	allPayloads := strings.Join(invoke.snapshot(), "\n")
	for _, seq := range []string{"1", "2"} {
		if !strings.Contains(allPayloads, fmt.Sprintf(`"SequenceNumber":%q`, seq)) {
			t.Fatalf("TRIM_HORIZON must deliver the history, got payloads: %s", allPayloads)
		}
	}
}
