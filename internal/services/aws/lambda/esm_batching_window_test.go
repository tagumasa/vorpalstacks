package lambda

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// TestValidateESMBatchWindowPair pins the documented pairing: a batch size
// above 10 requires a batching window of at least one second.
func TestValidateESMBatchWindowPair(t *testing.T) {
	cases := []struct {
		batchSize, window int32
		wantErr           bool
	}{
		{10, 0, false},
		{10, 5, false},
		{11, 0, true},
		{100, 0, true},
		{11, 1, false},
		{100, 1, false},
		{100, 5, false},
	}
	for _, tc := range cases {
		err := validateESMBatchWindowPair(tc.batchSize, tc.window)
		if tc.wantErr && err == nil {
			t.Fatalf("BatchSize %d with window %d must be rejected", tc.batchSize, tc.window)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("BatchSize %d with window %d must be accepted: %v", tc.batchSize, tc.window, err)
		}
	}
}

// TestBufferReady pins the gathering decision: the batch flushes when full,
// or when the window elapsed since the first buffered record.
func TestBufferReady(t *testing.T) {
	now := time.Now()
	if full, expired := bufferReady(3, 5, now, now.Add(1*time.Second), 5); full || expired {
		t.Fatalf("partial buffer inside the window must hold: full=%v expired=%v", full, expired)
	}
	if full, expired := bufferReady(5, 5, now, now, 5); !full || expired {
		t.Fatalf("full buffer must flush: full=%v expired=%v", full, expired)
	}
	if full, expired := bufferReady(3, 5, now, now.Add(6*time.Second), 5); full || !expired {
		t.Fatalf("expired window must flush: full=%v expired=%v", full, expired)
	}
	if full, expired := bufferReady(0, 5, time.Time{}, now, 5); full || expired {
		t.Fatalf("empty buffer must not flush: full=%v expired=%v", full, expired)
	}
}

// TestStreamBuffer_HoldAndReadPosition pins the buffer bookkeeping: the
// read position moves with the buffer while the durable checkpoint stays
// put, and dropping returns the held items.
func TestStreamBuffer_HoldAndReadPosition(t *testing.T) {
	p := &esmPoller{kinesisCP: make(map[string]string), buffers: make(map[string]*streamBuffer)}
	buf := p.getStreamBuffer("k")
	if p.streamBufferReadSeq("k") != "" {
		t.Fatal("fresh buffer must read from the checkpoint")
	}
	buf.items = append(buf.items, kinesisTestItem("1", "a"))
	buf.readThrough = "1"
	if got := p.streamBufferReadSeq("k"); got != "1" {
		t.Fatalf("read position must follow the buffer, got %q", got)
	}
	held := p.dropStreamBuffer("k")
	if len(held) != 1 || held[0].seq != "1" {
		t.Fatalf("drop must return the held items: %+v", held)
	}
	if p.streamBufferReadSeq("k") != "" {
		t.Fatal("dropped buffer must not hold a read position")
	}
}

// TestFlushStreamBuffer_ChunksByBatchSize pins the flush behaviour: the
// gathered records are delivered in batch-size chunks in order and the
// checkpoint advances to the buffer's read position.
func TestFlushStreamBuffer_ChunksByBatchSize(t *testing.T) {
	inv := &capturedInvoke{}
	p := newWindowTestPoller(inv)
	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:                    "arn:aws:lambda:us-east-1:123456789012:function:buf",
		BatchSize:                      5,
		MaximumBatchingWindowInSeconds: 10,
	}
	items := make([]streamBatchItem, 7)
	for i := range items {
		items[i] = kinesisTestItem(string(rune('1'+i)), "d")
	}
	outcome := p.flushStreamBuffer(context.Background(), mapping, testStreamSource("kinesis", "s1"), "k", "7", items)
	if outcome.err != nil {
		t.Fatalf("flush: %v", outcome.err)
	}
	if !outcome.delivered || outcome.discarded {
		t.Fatalf("outcome wrong: %+v", outcome)
	}
	if len(inv.payloads) != 2 {
		t.Fatalf("7 items with batch size 5 must deliver two chunks, got %d", len(inv.payloads))
	}
	first := decodeBatchSeqs(t, inv.payloads[0])
	second := decodeBatchSeqs(t, inv.payloads[1])
	if len(first) != 5 || len(second) != 2 || first[0] != "1" || second[0] != "6" || second[1] != "7" {
		t.Fatalf("chunking wrong: %v then %v", first, second)
	}
	p.kinesisCPMu.RLock()
	cp := p.kinesisCP["k"]
	p.kinesisCPMu.RUnlock()
	if cp != "7" {
		t.Fatalf("checkpoint must advance to the read position, got %q", cp)
	}
}

// decodeBatchSeqs reads the sequence identifiers out of a plain stream
// batch payload.
func decodeBatchSeqs(t *testing.T, payload []byte) []string {
	t.Helper()
	var decoded struct {
		Records []struct {
			Kinesis struct {
				SequenceNumber string `json:"sequenceNumber"`
			} `json:"kinesis"`
		} `json:"Records"`
	}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	seqs := make([]string, len(decoded.Records))
	for i, r := range decoded.Records {
		seqs[i] = r.Kinesis.SequenceNumber
	}
	return seqs
}
