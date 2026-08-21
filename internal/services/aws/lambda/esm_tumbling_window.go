package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// Tumbling-window machinery.
//
// With TumblingWindowInSeconds configured the poller groups each read
// batch by the epoch-aligned window of the record insertion times,
// delivers per-window invocations whose responses carry an aggregated
// state, and closes a window when a later window's records arrive, when
// the shard ends, or after the window's end time plus an inactivity grace
// period. The durable checkpoint only advances when a window completes:
// "After successful invocation, your function checkpoints the sequence
// number and stream processing continues. If invocation is unsuccessful,
// your Lambda function suspends further processing until a successful
// invocation." Records inside an open window are therefore re-delivered
// from the window start after a restart — the documented at-least-once
// contract.

// maxWindowStateBytes is the documented per-shard state ceiling: "Your
// state can be a maximum of 1 MB for each shard. If it exceeds that size,
// Lambda terminates the window early."
const maxWindowStateBytes = 1 << 20

// windowInactivityClose is the documented quiet period after a window's
// end time before the window is considered over: "When no new records are
// being added to the current window, Lambda waits for up to 2 minutes
// before assuming that the window is over."
const windowInactivityClose = 2 * time.Minute

// streamWindowSpan is the window member of a time-window event; the
// boundaries are rendered as ISO-8601 UTC timestamps, matching the
// documented event examples ("2020-12-09T07:04:00Z").
type streamWindowSpan struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// streamWindowEnvelope is the documented time-window event shape: the
// batch records plus the window, the aggregated state, and the shard and
// source identification.
type streamWindowEnvelope struct {
	Records                 []interface{}    `json:"Records"`
	Window                  streamWindowSpan `json:"window"`
	State                   json.RawMessage  `json:"state"`
	ShardID                 string           `json:"shardId,omitempty"`
	EventSourceARN          string           `json:"eventSourceARN"`
	IsFinalInvokeForWindow  bool             `json:"isFinalInvokeForWindow"`
	IsWindowTerminatedEarly bool             `json:"isWindowTerminatedEarly"`
}

// windowedStreamItem pairs a batch item with the epoch-aligned start of
// the tumbling window its insertion time falls into.
type windowedStreamItem struct {
	item        streamBatchItem
	windowStart int64
}

// shardWindow is the in-memory state of one open tumbling window. It is
// looked up under the poller's windows mutex and then confined to the
// worker driving its mapping.
type shardWindow struct {
	windowStart  int64           // epoch-aligned start, unix seconds
	windowEnd    int64           // nominal end = windowStart + window seconds
	state        json.RawMessage // aggregated state from the last successful response
	readSeq      string          // sequence to read after; empty reads from the checkpoint
	lastActivity time.Time       // last cycle that consumed records into the window
}

// windowDelivery describes one time-window invocation.
type windowDelivery struct {
	span      streamWindowSpan
	shardID   string
	sourceARN string
	final     bool
	early     bool
}

// windowCycleResult reports what a windowed cycle did, for status
// reporting by the caller.
type windowCycleResult struct {
	processedAny bool
	discarded    bool
}

// windowStartOf returns the epoch-aligned tumbling window containing t.
// "Lambda determines tumbling window boundaries based on the time when
// records were inserted into the stream."
func windowStartOf(t, windowSeconds int64) int64 {
	if windowSeconds <= 0 {
		return 0
	}
	return (t / windowSeconds) * windowSeconds
}

// splitByWindow partitions windowed items into consecutive per-window
// groups. Insertion times are monotonic within a shard, so windows appear
// in order.
func splitByWindow(items []windowedStreamItem) [][]windowedStreamItem {
	var groups [][]windowedStreamItem
	for _, it := range items {
		if n := len(groups); n > 0 && groups[n-1][0].windowStart == it.windowStart {
			groups[n-1] = append(groups[n-1], it)
			continue
		}
		groups = append(groups, []windowedStreamItem{it})
	}
	return groups
}

// ddbArrivalUnix returns the record insertion time a DynamoDB Streams
// record uses for tumbling window boundary determination.
func ddbArrivalUnix(rec *eventbus.DynamoDBStreamRecord) int64 {
	if ts, ok := rec.Dynamodb["ApproximateCreationDateTime"].(float64); ok {
		return int64(ts)
	}
	return 0
}

// spanOf renders the window boundaries in the documented ISO-8601 UTC
// form.
func spanOf(start, windowSeconds int64) streamWindowSpan {
	return streamWindowSpan{
		Start: time.Unix(start, 0).UTC().Format(time.RFC3339),
		End:   time.Unix(start+windowSeconds, 0).UTC().Format(time.RFC3339),
	}
}

// marshalWindowedBatch renders one time-window invocation payload.
func marshalWindowedBatch(d windowDelivery, state json.RawMessage, items []streamBatchItem) []byte {
	records := make([]interface{}, len(items))
	for i, item := range items {
		records[i] = item.record
	}
	if len(state) == 0 {
		state = json.RawMessage("{}")
	}
	payload, err := json.Marshal(streamWindowEnvelope{
		Records:                 records,
		Window:                  d.span,
		State:                   state,
		ShardID:                 d.shardID,
		EventSourceARN:          d.sourceARN,
		IsFinalInvokeForWindow:  d.final,
		IsWindowTerminatedEarly: d.early,
	})
	if err != nil {
		return nil
	}
	return payload
}

// extractWindowState returns the state member of a time-window invocation
// response. "when using tumbling windows, your Lambda function response
// must contain a state property. If the response does not contain a state
// property, Lambda considers this a failed invocation."
//
// The container runtime appends the handler's return value to its stdout
// after any console output, so when the payload carries log lines the
// state object sits on the final line.
func extractWindowState(payload []byte) (json.RawMessage, error) {
	doc, ok := finalJSONDocument(payload)
	if !ok {
		return nil, fmt.Errorf("esm: time-window response carries no state member")
	}
	state, err := stateFromJSONObject(doc)
	if err != nil {
		return nil, fmt.Errorf("esm: time-window response carries no state member")
	}
	return state, nil
}

// stateFromJSONObject parses a time-window response and returns its state
// member verbatim.
func stateFromJSONObject(payload []byte) (json.RawMessage, error) {
	var resp map[string]json.RawMessage
	if err := json.Unmarshal(payload, &resp); err != nil {
		return nil, err
	}
	state, ok := resp["state"]
	if !ok || len(state) == 0 {
		return nil, fmt.Errorf("no state member")
	}
	return state, nil
}

// invokeWindowedWithRetry drives one time-window invocation with the
// mapping's retry policy; a response without the state member counts as a
// failed invocation and consumes a retry attempt.
func (p *esmPoller) invokeWindowedWithRetry(ctx context.Context, mapping *lambdastore.EventSourceMapping, d windowDelivery, state json.RawMessage, items []streamBatchItem) (json.RawMessage, error) {
	payload := marshalWindowedBatch(d, state, items)
	if payload == nil {
		return nil, fmt.Errorf("esm: failed to marshal time-window payload")
	}
	var responseState json.RawMessage
	err := p.invokeWithRetry(ctx, mapping, payload, func(result *lambdastore.InvocationResult) error {
		st, serr := extractWindowState(result.Payload)
		if serr != nil {
			return serr
		}
		responseState = st
		return nil
	})
	if err != nil {
		return nil, err
	}
	return responseState, nil
}

// processWindowedBatch delivers one chunk of an open window, bisecting on
// function error exactly like the non-windowed path and threading the
// aggregated state through the halves in order. A chunk that exhausts a
// finite retry budget is discarded: the window stays open and keeps
// everything aggregated so far.
func (p *esmPoller) processWindowedBatch(ctx context.Context, mapping *lambdastore.EventSourceMapping, d windowDelivery, src streamSource, state json.RawMessage, items []streamBatchItem) (newState json.RawMessage, lastConsumed string, discarded, oversize bool, err error) {
	if len(items) == 0 {
		return state, "", false, false, nil
	}
	newState, err = p.invokeWindowedWithRetry(ctx, mapping, d, state, items)
	if err == nil {
		return newState, items[len(items)-1].seq, false, len(newState) > maxWindowStateBytes, nil
	}
	if mapping.BisectBatchOnFunctionError && len(items) > 1 {
		mid := len(items) / 2
		firstState, firstLast, firstDiscarded, firstOversize, firstErr := p.processWindowedBatch(ctx, mapping, d, src, state, items[:mid])
		if firstErr != nil {
			// The first half is still failing with an infinite retry
			// budget; preserve whatever it managed to consume.
			return firstState, firstLast, false, false, firstErr
		}
		if firstOversize {
			return firstState, firstLast, false, true, nil
		}
		secondState, secondLast, secondDiscarded, secondOversize, secondErr := p.processWindowedBatch(ctx, mapping, d, src, firstState, items[mid:])
		if secondErr != nil {
			return firstState, firstLast, firstDiscarded || secondDiscarded, false, secondErr
		}
		if secondLast == "" {
			return firstState, firstLast, firstDiscarded || secondDiscarded, secondOversize, nil
		}
		return secondState, secondLast, firstDiscarded || secondDiscarded, secondOversize, nil
	}
	if mapping.MaximumRetryAttempts >= 0 {
		p.deliverDiscardedBatch(ctx, mapping, src, streamFailureBatchInfoOf(src, items),
			marshalWindowedBatch(d, state, items), retryAttemptsOf(mapping), discardedBatchResponse(err))
		return state, items[len(items)-1].seq, true, false, nil
	}
	return nil, "", false, false, err
}

// openStreamWindow creates the state for a window starting at start.
func (p *esmPoller) openStreamWindow(key string, start, windowSeconds int64) *shardWindow {
	win := &shardWindow{
		windowStart:  start,
		windowEnd:    start + windowSeconds,
		state:        json.RawMessage("{}"),
		lastActivity: time.Now(),
	}
	p.windowsMu.Lock()
	if p.windows == nil {
		p.windows = make(map[string]*shardWindow)
	}
	p.windows[key] = win
	p.windowsMu.Unlock()
	return win
}

// advanceCheckpoint records a consumed sequence both in memory and in the
// durable checkpoint store.
func (p *esmPoller) advanceCheckpoint(cpKey, seq string) {
	p.kinesisCPMu.Lock()
	p.kinesisCP[cpKey] = seq
	p.kinesisCPMu.Unlock()
	if err := p.persistKinesisCheckpoint(cpKey, seq); err != nil {
		logs.Warn("esm: failed to persist stream checkpoint, in-memory state may diverge on restart",
			logs.String("key", cpKey), logs.Err(err))
	}
}

// completeStreamWindow checkpoints the window's last consumed sequence and
// drops the window state: "After processing, the window completes and your
// final invocation completes, and then the state is dropped."
func (p *esmPoller) completeStreamWindow(key string, win *shardWindow, cpSeq string) {
	if cpSeq != "" {
		p.advanceCheckpoint(key, cpSeq)
	}
	p.windowsMu.Lock()
	if cur, ok := p.windows[key]; ok && cur == win {
		delete(p.windows, key)
	}
	p.windowsMu.Unlock()
}

// closeStreamWindow delivers the final invocation of an open window and
// checkpoints the window's last consumed sequence on success. The final
// invocation is never discarded: "If invocation is unsuccessful, your
// Lambda function suspends further processing until a successful
// invocation." early marks a termination caused by an oversized state.
func (p *esmPoller) closeStreamWindow(ctx context.Context, mapping *lambdastore.EventSourceMapping, key, shardID, sourceARN string, win *shardWindow, early bool) error {
	d := windowDelivery{
		span:      spanOf(win.windowStart, win.windowEnd-win.windowStart),
		shardID:   shardID,
		sourceARN: sourceARN,
		final:     true,
		early:     early,
	}
	if _, err := p.invokeWindowedWithRetry(ctx, mapping, d, win.state, nil); err != nil {
		return err
	}
	p.completeStreamWindow(key, win, win.readSeq)
	return nil
}

// processStreamWindow advances the tumbling window of one mapping+shard.
// It delivers each per-window chunk of the read batch — the last chunk of
// a window is marked final when a later window's records are already
// visible in the same batch — and closes an open window whose end time and
// inactivity grace period have both expired. readThrough is the sequence
// of the last record read this cycle; every record in the batch has been
// delivered, discarded or dropped by filtering, so an open window's read
// position advances to it.
func (p *esmPoller) processStreamWindow(ctx context.Context, mapping *lambdastore.EventSourceMapping, key string, src streamSource, items []windowedStreamItem, readThrough string) (windowCycleResult, error) {
	windowSeconds := int64(mapping.TumblingWindowInSeconds)
	shardID, sourceARN := src.shardID, src.streamArn

	p.windowsMu.Lock()
	win := p.windows[key]
	p.windowsMu.Unlock()

	var res windowCycleResult
	groups := splitByWindow(items)

	for gi, group := range groups {
		if win == nil {
			win = p.openStreamWindow(key, group[0].windowStart, windowSeconds)
		} else if group[0].windowStart > win.windowStart {
			// A later window's records arrived: the open window is complete.
			// Its own records were delivered in earlier cycles, so the final
			// invocation carries only the aggregated state.
			if err := p.closeStreamWindow(ctx, mapping, key, shardID, sourceARN, win, false); err != nil {
				return res, err
			}
			win = p.openStreamWindow(key, group[0].windowStart, windowSeconds)
		}

		finalChunk := gi+1 < len(groups)
		chunk := make([]streamBatchItem, len(group))
		for i, wi := range group {
			chunk[i] = wi.item
		}
		d := windowDelivery{
			span:      spanOf(win.windowStart, windowSeconds),
			shardID:   shardID,
			sourceARN: sourceARN,
			final:     finalChunk,
		}
		state, lastConsumed, discarded, oversize, err := p.processWindowedBatch(ctx, mapping, d, src, win.state, chunk)
		if err != nil {
			// The chunk still fails with an infinite retry budget. Keep any
			// partial bisection progress; nothing after it may be consumed.
			if state != nil && lastConsumed != "" {
				win.state = state
				win.readSeq = lastConsumed
				win.lastActivity = time.Now()
			}
			return res, err
		}
		win.state = state
		win.readSeq = lastConsumed
		win.lastActivity = time.Now()
		if discarded {
			res.discarded = true
		} else {
			res.processedAny = true
		}

		if oversize {
			// "Your state can be a maximum of 1 MB for each shard. If it
			// exceeds that size, Lambda terminates the window early." The
			// remaining records re-read after the checkpoint advances land
			// in a fresh window with a fresh state.
			if err := p.closeStreamWindow(ctx, mapping, key, shardID, sourceARN, win, true); err != nil {
				return res, err
			}
			return res, nil
		}
		if finalChunk {
			if discarded {
				// The records are gone but the window still needs a
				// successful final invocation before it may complete.
				if err := p.closeStreamWindow(ctx, mapping, key, shardID, sourceARN, win, false); err != nil {
					return res, err
				}
			} else {
				p.completeStreamWindow(key, win, lastConsumed)
			}
			win = nil
		}
	}

	if win == nil {
		// Nothing survived this cycle and no window is open, so the durable
		// checkpoint may take over the read position.
		if readThrough != "" {
			p.advanceCheckpoint(key, readThrough)
		}
		return res, nil
	}
	if readThrough != "" {
		win.readSeq = readThrough
	}

	// "When no new records are being added to the current window, Lambda
	// waits for up to 2 minutes before assuming that the window is over."
	// The grace period only applies once the window's end time has passed;
	// earlier silence still lets further records join the window.
	if time.Now().Unix() >= win.windowEnd && time.Since(win.lastActivity) >= windowInactivityClose {
		if err := p.closeStreamWindow(ctx, mapping, key, shardID, sourceARN, win, false); err != nil {
			return res, err
		}
	}
	return res, nil
}

// closeEndedShardWindow delivers the final invocation for a window that is
// still open on a shard that has ended: "When a shard ends, Lambda
// considers the current window to be closed, and any child shards start
// their own window in a fresh state." The closed shard yields no further
// records, so the final invocation carries only the aggregated state.
func (p *esmPoller) closeEndedShardWindow(ctx context.Context, mapping *lambdastore.EventSourceMapping, cpKey, shardID string) {
	p.windowsMu.Lock()
	win := p.windows[cpKey]
	p.windowsMu.Unlock()
	if win == nil {
		return
	}
	if err := p.closeStreamWindow(ctx, mapping, cpKey, shardID, mapping.EventSourceArn, win, false); err != nil {
		p.log("failed to deliver final tumbling window invocation for ended shard",
			"mapping", mapping.UUID, "shard", shardID, "error", err)
		if rerr := p.esmStore.SetProcessingResult(mapping.UUID, err.Error()); rerr != nil {
			logs.Warn("esm: failed to set result after ended-shard window close",
				logs.String("mapping", mapping.UUID), logs.Err(rerr))
		}
	}
}

// purgeStaleWindows drops tumbling-window state for mappings that no
// longer have an active stream mapping, mirroring the checkpoint purge.
func (p *esmPoller) purgeStaleWindows(activeUUIDs map[string]struct{}) {
	p.windowsMu.Lock()
	for key := range p.windows {
		uuid := strings.TrimPrefix(key, "ddb:")
		if idx := strings.IndexByte(uuid, ':'); idx >= 0 {
			uuid = uuid[:idx]
		}
		if _, active := activeUUIDs[uuid]; !active {
			delete(p.windows, key)
		}
	}
	p.windowsMu.Unlock()
}
