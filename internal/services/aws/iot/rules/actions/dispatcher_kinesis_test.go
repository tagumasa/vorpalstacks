package actions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
)

// capturingKinesisInvoker records the bytes one Kinesis action wrote.
type capturingKinesisInvoker struct {
	data []byte
	key  string
}

func (c *capturingKinesisInvoker) ListShards(context.Context, string, string) ([]invokers.ShardInfo, error) {
	return nil, nil
}

func (c *capturingKinesisInvoker) PutRecord(_ context.Context, _, _, partitionKey string, data []byte) (string, string, error) {
	c.data = data
	c.key = partitionKey
	return "seq-1", "shardId-000000000000", nil
}

func (c *capturingKinesisInvoker) CreateShardIterator(context.Context, string, string, string, string, string, *time.Time) (string, error) {
	return "", nil
}

func (c *capturingKinesisInvoker) GetRecords(context.Context, string, string, string, string, int32, bool) ([]invokers.KinesisRecord, string, error) {
	return nil, "", nil
}

func (c *capturingKinesisInvoker) StreamExists(context.Context, string, string) (bool, error) {
	return true, nil
}

// TestDispatchKinesisEncodesPayloadLikeEveryProducer pins the record body's
// encoding: the platform's Kinesis Data member carries the base64 form and
// GetRecords returns it as-is, so the action pre-encodes its JSON payload
// exactly as every other cross-service producer does — a raw JSON body on a
// shared stream would decode differently from the other producers' records.
func TestDispatchKinesisEncodesPayloadLikeEveryProducer(t *testing.T) {
	invoker := &capturingKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	d := NewDispatcher(bus, nil)

	payload := map[string]interface{}{"temp": float64(25)}
	config := &ActionConfig{
		Type:       "kinesis",
		StreamName: "pinned-stream",
		Extra:      map[string]interface{}{"partitionKey": "device-1"},
	}
	actionPayload := &ActionPayload{Raw: payload, JSONBytes: mustJSON(payload)}

	if err := d.dispatchKinesis(context.Background(), config, actionPayload); err != nil {
		t.Fatalf("dispatch kinesis action: %v", err)
	}

	if invoker.key != "device-1" {
		t.Fatalf("partition key = %q, want the action's configured key", invoker.key)
	}
	if string(invoker.data) != base64.StdEncoding.EncodeToString(actionPayload.JSONBytes) {
		t.Fatalf("record data = %q, want the base64 form of the JSON payload %q",
			invoker.data, actionPayload.JSONBytes)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(invoker.data))
	if err != nil {
		t.Fatalf("record data is not valid base64: %v", err)
	}
	var round map[string]interface{}
	if err := json.Unmarshal(decoded, &round); err != nil {
		t.Fatalf("decoded record data is not the JSON payload: %v", err)
	}
	if round["temp"] != float64(25) {
		t.Fatalf("decoded payload = %v, want the action's payload", round)
	}
}
