package lambda

import (
	"context"
	"time"

	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// Shared stream-poll driver.
//
// Every stream record source (Kinesis, DynamoDB Streams) runs the same
// pipeline once its source-specific fetch completes: the checkpoint
// read-ahead resolution, the initial-anchor persistence, the
// windowed/buffered/plain three-way batch dispatch, and the cycle-level
// status reporting whose precedence the cycleReport owns. The source flows
// fetch their batches, prepare each into a streamFetchedBatch, and hand
// them to a streamCycle; the delivery machinery (retry, bisection,
// windows, buffering) lives in the shared batch files.

// cycleReport folds one poll cycle's reporting precedence: a failure
// outranks a discard, which outranks a partial batch response, which
// outranks success. The end-of-cycle success report is deferred behind
// this precedence so a batch that failed or was discarded earlier in the
// same cycle cannot be overwritten by a later successful batch.
type cycleReport struct {
	failure      bool
	discarded    bool
	processedAny bool
	// partial counts the records the function reported in
	// batchItemFailures this cycle; a partial report outranks the
	// success report but stays behind a failure or discard.
	partial int
}

func (c *cycleReport) recordFailure()   { c.failure = true }
func (c *cycleReport) recordDiscard()   { c.discarded = true }
func (c *cycleReport) recordProcessed() { c.processedAny = true }

func (c *cycleReport) recordPartial(reported int) { c.partial += reported }

// shouldReportSuccess reports whether the cycle may write the success
// result: only a cycle without any failure, discard or partial batch
// response that actually delivered records.
func (c *cycleReport) shouldReportSuccess() bool {
	return !c.failure && !c.discarded && c.partial == 0 && c.processedAny
}

// streamFetchedBatch is one fetched batch prepared for delivery: the
// surviving records as plain delivery items and as tumbling-window
// pairings, the raw batch's ordering keys, and the sequence of the last
// record the read advanced through.
type streamFetchedBatch struct {
	items   []streamBatchItem
	witems  []windowedStreamItem
	keys    map[string]struct{}
	readSeq string
}

// streamCycle carries the per-mapping state of one poll cycle through the
// shared dispatch. Kinesis drives several shards through one cycle, so the
// per-shard checkpoint key and source identity are dispatch parameters
// while the report aggregates across every dispatch of the cycle.
type streamCycle struct {
	p         *esmPoller
	mapping   *lambdastore.EventSourceMapping
	label     string // source label for status logs, e.g. "Kinesis ESM"
	batchSize int32
	windowed  bool
	buffered  bool
	report    cycleReport
}

// newStreamCycle derives the dispatch configuration of one mapping: a
// positive tumbling window aggregates records into window-bound batches, a
// positive batching window gathers records across cycles, and neither
// invokes every read immediately.
func (p *esmPoller) newStreamCycle(mapping *lambdastore.EventSourceMapping, label string) *streamCycle {
	return &streamCycle{
		p:         p,
		mapping:   mapping,
		label:     label,
		batchSize: clampESMBatchSize(mapping.BatchSize, mapping.EventSourceArn),
		windowed:  mapping.TumblingWindowInSeconds > 0,
		buffered:  mapping.TumblingWindowInSeconds <= 0 && mapping.MaximumBatchingWindowInSeconds > 0,
		report:    cycleReport{},
	}
}

// streamReadPosition resolves where a stream mapping's next read starts:
// the durable checkpoint, unless an open tumbling window or a gathering
// batching-window buffer has read ahead of it. While a tumbling window is
// open the read position advances per delivered chunk; the durable
// checkpoint only moves when the window completes. A batching-window
// buffer likewise reads ahead of the checkpoint while records are being
// gathered.
func (p *esmPoller) streamReadPosition(cpKey string, windowed, buffered bool) string {
	p.kinesisCPMu.RLock()
	readFrom := p.kinesisCP[cpKey]
	p.kinesisCPMu.RUnlock()
	if windowed {
		p.windowsMu.Lock()
		if w := p.windows[cpKey]; w != nil && w.readSeq != "" {
			readFrom = w.readSeq
		}
		p.windowsMu.Unlock()
	} else if buffered {
		if seq := p.streamBufferReadSeq(cpKey); seq != "" {
			readFrom = seq
		}
	}
	return readFrom
}

// dispatchBatches runs the three-way batch dispatch for one shard's
// fetched batches. anchor is the initial read position to persist when the
// first read came back empty, so the next cycle reads strictly after it
// instead of re-anchoring at a newer position, which would skip everything
// that arrived in between. Windowed and buffering mappings need this as
// well: until their first window or buffer opens they hold no read
// position of their own, so without a durable anchor every empty poll
// would re-anchor at the newest record and lose all but the newest record
// of an inter-poll burst. A cycle that returned records advances the
// checkpoint through the normal delivery flow, which never moves it
// backwards.
func (c *streamCycle) dispatchBatches(ctx context.Context, cpKey string, src streamSource, batches []streamFetchedBatch, anchor string) {
	if len(batches) == 0 && anchor != "" {
		c.p.advanceCheckpoint(cpKey, anchor)
	}

	switch {
	case c.windowed:
		// The tumbling window state is order-sensitive, so the fetched
		// batches flow through the window machinery sequentially.
		for _, b := range batches {
			c.windowCycle(ctx, cpKey, src, b.witems, b.readSeq)
		}
		if len(batches) == 0 {
			// No records this cycle, but the inactivity grace period may
			// have expired for an open window.
			c.windowCycle(ctx, cpKey, src, nil, "")
		}
	case c.buffered:
		c.runBuffered(ctx, cpKey, src, batches)
	default:
		c.runPlain(ctx, cpKey, src, batches)
	}
}

// windowCycle advances the tumbling window of one dispatch and folds the
// cycle result into the report.
func (c *streamCycle) windowCycle(ctx context.Context, cpKey string, src streamSource, items []windowedStreamItem, readThrough string) {
	res, werr := c.p.processStreamWindow(ctx, c.mapping, cpKey, src, items, readThrough)
	if werr != nil {
		c.recordFailure("tumbling window", werr)
		return
	}
	if res.discarded {
		c.recordDiscard("")
		return
	}
	if res.processedAny {
		c.report.recordProcessed()
	}
}

// runBuffered gathers the surviving records of every fetched batch into
// the batching-window buffer; the read position moves with the buffer, the
// checkpoint only at flush.
func (c *streamCycle) runBuffered(ctx context.Context, cpKey string, src streamSource, batches []streamFetchedBatch) {
	buf := c.p.getStreamBuffer(cpKey)
	for _, b := range batches {
		if len(b.items) == 0 {
			continue
		}
		if len(buf.items) == 0 {
			buf.firstAt = time.Now()
		}
		buf.items = append(buf.items, b.items...)
		buf.readThrough = b.readSeq
	}
	full, expired := bufferReady(len(buf.items), int(c.batchSize), buf.firstAt, time.Now(), c.mapping.MaximumBatchingWindowInSeconds)
	if len(buf.items) > 0 && !full && !expired {
		// The window is still gathering; hold the records.
		return
	}
	items := c.p.dropStreamBuffer(cpKey)
	readThrough := buf.readThrough
	if len(items) == 0 {
		// Nothing survived, but the read position still moves.
		if readThrough != "" {
			c.p.advanceCheckpoint(cpKey, readThrough)
		}
		return
	}
	outcome := c.p.flushStreamBuffer(ctx, c.mapping, src, cpKey, readThrough, items)
	switch {
	case outcome.err != nil:
		c.recordFailure("batching window", outcome.err)
	case outcome.discarded:
		c.recordDiscard("discarded failed stream records after exhausting retries")
	case outcome.reported > 0:
		c.report.recordPartial(outcome.reported)
	case outcome.delivered:
		c.report.recordProcessed()
	}
}

// runPlain delivers the fetched batches through the ordered parallel
// runner. "The number of batches to process from each shard concurrently":
// batches run in parallel except where they share a partition key, and the
// checkpoint takes the contiguous consumed prefix.
func (c *streamCycle) runPlain(ctx context.Context, cpKey string, src streamSource, batches []streamFetchedBatch) {
	if len(batches) == 0 {
		return
	}

	keys := make([]map[string]struct{}, len(batches))
	for i, b := range batches {
		keys[i] = b.keys
	}
	outcomes := runOrderedBatches(ctx, len(batches), keys, func(ctx context.Context, idx int) batchOutcome {
		return c.plainBatch(ctx, src, batches[idx])
	})
	lastConsumed, delivered, discarded, reported, failure := prefixOutcome(outcomes)
	if failure != nil {
		c.recordFailure("", failure)
	} else if discarded {
		c.recordDiscard("discarded failed stream records after exhausting retries")
	}
	if reported > 0 {
		c.report.recordPartial(reported)
	}
	if delivered {
		c.report.recordProcessed()
	}

	if lastConsumed != "" {
		c.p.advanceCheckpoint(cpKey, lastConsumed)
	}
}

// plainBatch delivers one fetched batch; a batch whose records all expired
// or were all filtered out still consumes its position.
func (c *streamCycle) plainBatch(ctx context.Context, src streamSource, b streamFetchedBatch) batchOutcome {
	if len(b.items) == 0 {
		return batchOutcome{lastConsumed: b.readSeq}
	}
	return c.p.processStreamBatch(ctx, c.mapping, src, b.items)
}

// recordFailure folds a failure into the cycle report, the poller
// diagnostic and the mapping status. kind names the flow branch in the
// diagnostic, e.g. "tumbling window"; empty names the plain batch flow.
func (c *streamCycle) recordFailure(kind string, err error) {
	c.report.recordFailure()
	detail := c.label
	if kind != "" {
		detail += " " + kind
	}
	c.p.log("lambda invocation failed for "+detail, "function", c.mapping.FunctionArn, "error", err)
	if rerr := c.p.esmStore.SetProcessingResult(c.mapping.UUID, err.Error()); rerr != nil {
		logs.Warn("esm: failed to set state after "+detail+" error",
			logs.String("mapping", c.mapping.UUID), logs.Err(rerr))
	}
}

// recordDiscard folds a discard into the cycle report and the mapping
// status. A non-empty diagnostic is also emitted to the poller log.
func (c *streamCycle) recordDiscard(diagnostic string) {
	c.report.recordDiscard()
	if diagnostic != "" {
		c.p.log(diagnostic, "function", c.mapping.FunctionArn)
	}
	if rerr := c.p.esmStore.SetProcessingResult(c.mapping.UUID, "Records discarded after exhausting retries"); rerr != nil {
		logs.Warn("esm: failed to set discard result", logs.String("mapping", c.mapping.UUID), logs.Err(rerr))
	}
}

// finish writes the end-of-cycle status under the report's precedence:
// only a cycle that actually delivered records may report success, and
// never over a failure, discard or partial batch response the same cycle
// recorded: idle cycles leave the previous result untouched, a discard
// stays visible even when a later batch of the same cycle succeeded, and a
// partial response reports its records instead.
func (c *streamCycle) finish() {
	if c.report.shouldReportSuccess() {
		if err := c.p.esmStore.SetProcessingResult(c.mapping.UUID, "No errors."); err != nil {
			logs.Error("esm: failed to set state", logs.String("mapping", c.mapping.UUID), logs.String("error", err.Error()))
		}
	} else if !c.report.failure && !c.report.discarded && c.report.partial > 0 {
		if err := c.p.esmStore.SetProcessingResult(c.mapping.UUID, streamPartialResult(c.report.partial)); err != nil {
			logs.Error("esm: failed to set state", logs.String("mapping", c.mapping.UUID), logs.Err(err))
		}
	}
}
