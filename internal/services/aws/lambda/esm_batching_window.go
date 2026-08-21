package lambda

import (
	"context"
	"strings"
	"time"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// Batching-window buffering for stream event sources.
//
// "The maximum amount of time, in seconds, that Lambda spends gathering
// records before invoking the function" (CreateEventSourceMapping model,
// MaximumBatchingWindowInSeconds). With a positive window the poller holds
// a shard's surviving records in a per-shard buffer until the batch is
// full or the window has elapsed since the first buffered record, instead
// of invoking every cycle's read immediately. While records are held the
// durable checkpoint stays put and the buffer's read position advances, so
// a restart re-reads the held records — the at-least-once contract.
//
// Tumbling-window mappings bypass this buffering: their aggregation
// batches are already window-bound, and holding a window chunk past its
// boundary would delay the documented window close.

// streamBuffer accumulates one mapping+shard's surviving records between
// flushes. It is looked up under the poller's buffer mutex and then
// confined to the worker driving its mapping.
type streamBuffer struct {
	items       []streamBatchItem
	firstAt     time.Time // when the oldest buffered record arrived
	readThrough string    // sequence of the last record read into the buffer
}

// bufferReady reports whether a buffer may flush: the batch is full, or
// the window elapsed since the first buffered record.
func bufferReady(itemCount, batchSize int, firstAt, now time.Time, windowSeconds int32) (full, expired bool) {
	if batchSize > 0 && itemCount >= batchSize {
		return true, false
	}
	if !firstAt.IsZero() && now.Sub(firstAt) >= time.Duration(windowSeconds)*time.Second {
		return false, true
	}
	return false, false
}

// getStreamBuffer returns (and creates on first use) the buffer for a key.
func (p *esmPoller) getStreamBuffer(key string) *streamBuffer {
	p.bufferMu.Lock()
	defer p.bufferMu.Unlock()
	if p.buffers == nil {
		p.buffers = make(map[string]*streamBuffer)
	}
	buf, ok := p.buffers[key]
	if !ok {
		buf = &streamBuffer{}
		p.buffers[key] = buf
	}
	return buf
}

// dropStreamBuffer removes a key's buffer, returning the held items for a
// final flush.
func (p *esmPoller) dropStreamBuffer(key string) []streamBatchItem {
	p.bufferMu.Lock()
	defer p.bufferMu.Unlock()
	buf, ok := p.buffers[key]
	if !ok {
		return nil
	}
	delete(p.buffers, key)
	return buf.items
}

// streamBufferReadSeq returns the buffered read position for a key, if any.
func (p *esmPoller) streamBufferReadSeq(key string) string {
	p.bufferMu.Lock()
	defer p.bufferMu.Unlock()
	if buf, ok := p.buffers[key]; ok {
		return buf.readThrough
	}
	return ""
}

// purgeStaleBuffers drops batching-window gathering state for mappings that
// no longer have an active stream mapping, mirroring the other purges.
func (p *esmPoller) purgeStaleBuffers(activeUUIDs map[string]struct{}) {
	p.bufferMu.Lock()
	for key := range p.buffers {
		uuid := strings.TrimPrefix(key, "ddb:")
		if idx := strings.IndexByte(uuid, ':'); idx >= 0 {
			uuid = uuid[:idx]
		}
		if _, active := activeUUIDs[uuid]; !active {
			delete(p.buffers, key)
		}
	}
	p.bufferMu.Unlock()
}

// flushStreamBuffer delivers the buffered items in batch-size chunks, in
// order, and checkpoints the buffer's read position on success. A failed
// chunk with an infinite retry budget suspends the flush: the items after
// it stay unacknowledged and the caller must not advance anything. A
// chunk that ends on a partial batch response suspends the flush the same
// way — its cursor ends the consumed prefix, so the remaining chunks are
// not invoked; they re-read from the partial checkpoint on the next cycle.
func (p *esmPoller) flushStreamBuffer(ctx context.Context, mapping *lambdastore.EventSourceMapping, src streamSource, key, cpAdvanceTo string, items []streamBatchItem) batchOutcome {
	batchSize := int(mapping.BatchSize)
	if batchSize <= 0 {
		batchSize = 100
	}
	var agg batchOutcome
	for start := 0; start < len(items); start += batchSize {
		end := start + batchSize
		if end > len(items) {
			end = len(items)
		}
		chunk := items[start:end]
		outcome := p.processStreamBatch(ctx, mapping, src, chunk)
		if outcome.err != nil {
			// A chunk still failing with an infinite retry budget blocks
			// everything after it; nothing may be checkpointed.
			agg.err = outcome.err
			return agg
		}
		if outcome.lastConsumed != "" {
			agg.lastConsumed = outcome.lastConsumed
		}
		agg.discarded = agg.discarded || outcome.discarded
		agg.delivered = agg.delivered || outcome.delivered
		agg.reported += outcome.reported
		if outcome.short {
			agg.short = true
			break
		}
	}
	switch {
	case agg.short && agg.lastConsumed == "":
		// Nothing was consumable; keep the existing checkpoint.
	case agg.short:
		p.advanceCheckpoint(key, agg.lastConsumed)
	default:
		p.advanceCheckpoint(key, cpAdvanceTo)
	}
	return agg
}
