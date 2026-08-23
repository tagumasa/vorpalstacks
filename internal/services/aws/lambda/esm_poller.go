package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

const (
	// defaultESMPollInterval is the default interval between polling cycles
	// for SQS event source mappings. Matches the AWS minimum poll interval
	// for Lambda ESM.
	defaultESMPollInterval = 1 * time.Second

	// defaultESMMaxWorkers is the maximum number of concurrent polling
	// goroutines, one per unique event source mapping.
	defaultESMMaxWorkers = 32

	// maxSQSBatchSize is the maximum number of SQS messages that can be
	// delivered to a Lambda function in a single invocation batch.
	maxSQSBatchSize = int32(10000)

	// sqsReceiveMessageMax is the maximum number of messages that can be
	// returned by a single SQS ReceiveMessage call.
	sqsReceiveMessageMax = int32(10)

	// esmCheckpointBucket is the Pebble bucket name used to persist
	// Kinesis ESM checkpoint data across server restarts.
	esmCheckpointBucket = "lambda-esm-checkpoints"
)

// esmSQSRecord represents a single SQS message formatted as an ESM event
// record in the Lambda event payload. The structure matches the AWS Lambda
// SQS event format documented at:
// https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html
type esmSQSRecord struct {
	MessageID                        string                 `json:"messageId"`
	ReceiptHandle                    string                 `json:"receiptHandle"`
	Body                             string                 `json:"body"`
	MD5OfBody                        string                 `json:"md5OfBody"`
	MD5OfMessageAttributes           string                 `json:"md5OfMessageAttributes"`
	MessageAttributes                map[string]interface{} `json:"messageAttributes,omitempty"`
	EventSourceARN                   string                 `json:"eventSourceArn"`
	EventSource                      string                 `json:"eventSource"`
	AWSRegion                        string                 `json:"awsRegion"`
	ApproximateReceiveCount          string                 `json:"approximateReceiveCount"`
	ApproximateFirstReceiveTimestamp string                 `json:"approximateFirstReceiveTimestamp"`
	SentTimestamp                    string                 `json:"sentTimestamp"`
	SenderID                         string                 `json:"senderId"`
}

// esmSQSEvent is the full Lambda event payload for an SQS event source
// mapping. It contains an array of SQS message records.
type esmSQSEvent struct {
	Records []esmSQSRecord `json:"Records"`
}

// esmPoller manages background polling of SQS event source mappings.
// For each enabled mapping whose EventSourceArn begins with "arn:aws:sqs:",
// it polls the source queue for messages, batches them according to the
// mapping's BatchSize and MaximumBatchingWindowInSeconds, and invokes the
// mapped Lambda function with the batch payload.
//
// Messages are deleted from the queue only after successful invocation.
// If invocation fails, messages remain in the queue and become visible
// again after the visibility timeout expires, providing at-least-once
// delivery semantics.
type esmPoller struct {
	mu             sync.Mutex
	running        bool
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	interval       time.Duration
	workers        int
	logger         logs.Logger
	bus            eventbus.Bus
	esmStore       *lambdastore.EventSourceStore
	lambdaSvc      *LambdaService
	region         string
	storageManager *storage.RegionStorageManager
	kinesisCP      map[string]string // "mappingUUID:streamName:shardID" -> lastSeqNum
	kinesisCPMu    sync.RWMutex
	// windows holds the open tumbling-window state per mapping+shard, keyed
	// like the checkpoint map ("mappingUUID:streamName:shardID" for Kinesis,
	// "ddb:mappingUUID" for DynamoDB streams).
	windows   map[string]*shardWindow
	windowsMu sync.Mutex
	// buffers holds the batching-window gathering state per mapping+shard,
	// keyed like the windows map.
	buffers  map[string]*streamBuffer
	bufferMu sync.Mutex
	// invoke overrides the Lambda invoke path; production leaves it nil so
	// invokeLambda falls through to the LambdaService event-source path.
	// Tests inject a fake to exercise retry and bisection logic.
	invoke func(ctx context.Context, functionRef string, payload []byte) (*lambdastore.InvocationResult, error)
}

// newESMPoller creates a new ESM poller with the given poll interval and
// maximum worker count. If interval is zero, defaultESMPollInterval is used;
// if workers is zero, defaultESMMaxWorkers is used.
func newESMPoller(interval time.Duration, workers int, logger logs.Logger) *esmPoller {
	if interval <= 0 {
		interval = defaultESMPollInterval
	}
	if workers <= 0 {
		workers = defaultESMMaxWorkers
	}
	return &esmPoller{
		interval:  interval,
		workers:   workers,
		logger:    logger,
		kinesisCP: make(map[string]string),
		windows:   make(map[string]*shardWindow),
		buffers:   make(map[string]*streamBuffer),
	}
}

// Start launches the background ESM polling loop. It is safe to call
// Start multiple times; subsequent calls are no-ops until Stop has been
// called.
func (p *esmPoller) Start(ctx context.Context) {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return
	}
	p.running = true
	ctx, p.cancel = context.WithCancel(ctx)
	p.mu.Unlock()

	p.loadKinesisCheckpoints()

	p.wg.Add(1)
	go p.pollLoop(ctx)
}

// Stop gracefully shuts down the ESM polling loop, waiting for any
// in-flight invocations to complete.
func (p *esmPoller) Stop() {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return
	}
	p.running = false
	if p.cancel != nil {
		p.cancel()
	}
	p.mu.Unlock()
	p.wg.Wait()
}

// pollLoop ticks at the configured interval, lists all event source mappings,
// and dispatches each SQS mapping to a worker goroutine for polling.
func (p *esmPoller) pollLoop(ctx context.Context) {
	defer p.wg.Done()
	defer func() { resilience.RecoverAndRestart("ESM pollLoop", &p.wg, func() { p.pollLoop(ctx) }) }()
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.pollAll(ctx)
		}
	}
}

// pollAll fetches all event source mappings from the store and dispatches
// each enabled SQS mapping to the worker pool. Mappings with non-SQS event
// source ARNs are skipped. Errors during individual mapping processing are
// logged but do not halt the loop.
func (p *esmPoller) pollAll(ctx context.Context) {
	regions := []string{p.region}
	if p.storageManager != nil {
		if activeRegions := p.storageManager.GetActiveRegions(); len(activeRegions) > 0 {
			regions = activeRegions
		}
	}

	// The checkpoint and window maps are shared across regions, so the
	// stale-state purges must run once per tick over the union of every
	// region's active stream mappings. Purging per region would let any
	// region without Kinesis mappings wipe the checkpoints other regions
	// just wrote, forcing their shards to restart from TRIM_HORIZON on
	// every cycle.
	activeStream := make(map[string]struct{})
	listingsComplete := true
	for _, region := range regions {
		active, ok := p.pollRegion(ctx, region)
		if !ok {
			listingsComplete = false
			continue
		}
		for uuid := range active {
			activeStream[uuid] = struct{}{}
		}
	}
	if listingsComplete {
		p.purgeStaleKinesisCheckpoints(activeStream)
		p.purgeStaleWindows(activeStream)
		p.purgeStaleBuffers(activeStream)
	}
}

func (p *esmPoller) pollRegion(ctx context.Context, region string) (map[string]struct{}, bool) {
	var esmStore *lambdastore.EventSourceStore
	if region == p.region && p.esmStore != nil {
		esmStore = p.esmStore
	} else if p.storageManager != nil {
		st, err := p.storageManager.GetStorage(region)
		if err != nil {
			return nil, true
		}
		if cached, ok := p.lambdaSvc.storeCache.Load(region); ok {
			if typed, ok := cached.(*lambdaStore); ok {
				esmStore = typed.EventSources
			}
		}
		if esmStore == nil {
			newStore := &lambdaStore{
				Functions:    lambdastore.NewFunctionStore(st, p.lambdaSvc.accountID, region),
				Layers:       lambdastore.NewLayerStore(st, p.lambdaSvc.accountID, region),
				EventSources: lambdastore.NewEventSourceStore(st, p.lambdaSvc.accountID, region),
			}
			if actual, loaded := p.lambdaSvc.storeCache.LoadOrStore(region, newStore); loaded {
				esmStore = actual.(*lambdaStore).EventSources
			} else {
				esmStore = newStore.EventSources
			}
		}
	}
	if esmStore == nil {
		return nil, true
	}

	result, err := esmStore.ListAllMappings()
	if err != nil {
		// The listing failed: report incompleteness so the caller skips the
		// checkpoint purge rather than wiping checkpoints for mappings it
		// could not see.
		p.log("failed to list event source mappings", "error", err)
		return nil, false
	}

	type pollJob struct {
		mapping *lambdastore.EventSourceMapping
	}
	jobs := make(chan pollJob, len(result))
	activeStreamUUIDs := make(map[string]struct{})
	for _, m := range result {
		if m.State != "Enabled" {
			continue
		}
		_, esmService, _, _, _ := arnutil.SplitARN(m.EventSourceArn)
		if esmService != "sqs" && esmService != "kinesis" && esmService != "dynamodb" {
			continue
		}
		if esmService == "kinesis" || esmService == "dynamodb" {
			activeStreamUUIDs[m.UUID] = struct{}{}
		}
		jobs <- pollJob{mapping: m}
	}
	close(jobs)

	jobCount := len(result)
	workerCount := p.workers
	if workerCount > jobCount {
		workerCount = jobCount
	}
	if workerCount == 0 {
		return activeStreamUUIDs, true
	}

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { resilience.RecoverPanic("lambda esm worker") }()
			for job := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					p.processMapping(ctx, job.mapping)
				}
			}
		}()
	}
	wg.Wait()
	return activeStreamUUIDs, true
}

// processMapping polls the SQS source queue for a single event source
// mapping. It receives messages, builds the Lambda event payload, invokes
// the mapped function, and deletes successfully processed messages from
// the queue.
func (p *esmPoller) processMapping(ctx context.Context, mapping *lambdastore.EventSourceMapping) {
	esmService := arnutil.GetServiceFromARN(mapping.EventSourceArn)
	if esmService == "kinesis" {
		p.processKinesisMapping(ctx, mapping)
		return
	}
	if esmService == "dynamodb" {
		p.processDynamoDBStreamsMapping(ctx, mapping)
		return
	}
	p.processSQSMapping(ctx, mapping)
}

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

func (p *esmPoller) processKinesisMapping(ctx context.Context, mapping *lambdastore.EventSourceMapping) {
	if p.bus == nil {
		return
	}

	_, _, streamRegion, _, resource := arnutil.SplitARN(mapping.EventSourceArn)
	if streamRegion == "" {
		p.log("failed to parse Kinesis event source ARN", "arn", mapping.EventSourceArn)
		return
	}

	streamName := resource
	if idx := strings.Index(resource, "stream/"); idx != -1 {
		streamName = resource[idx+len("stream/"):]
	}

	shards, err := p.bus.KinesisInvoker().ListShards(ctx, streamName)
	if err != nil {
		p.log("failed to list shards for Kinesis ESM", "stream", streamName, "error", err)
		return
	}

	batchSize := int32(mapping.BatchSize)
	if batchSize <= 0 {
		batchSize = 100
	}
	if batchSize > 10000 {
		batchSize = 10000
	}

	windowed := mapping.TumblingWindowInSeconds > 0
	windowSeconds := int64(mapping.TumblingWindowInSeconds)
	pf := parallelizationFactorOf(mapping)
	// A positive batching window gathers records across cycles instead of
	// invoking every read immediately; tumbling windows bypass it because
	// their aggregation batches are already window-bound.
	buffered := !windowed && mapping.MaximumBatchingWindowInSeconds > 0
	batchingWindow := mapping.MaximumBatchingWindowInSeconds

	// report keeps the cycle-level precedence: failure outranks discard,
	// which outranks the deferred end-of-cycle success report.
	report := cycleReport{}

	// reportWindow folds one tumbling-window cycle result into the mapping
	// status and the per-cycle flags.
	reportWindow := func(res windowCycleResult, werr error) {
		if werr != nil {
			report.recordFailure()
			p.log("lambda invocation failed for Kinesis ESM tumbling window",
				"function", mapping.FunctionArn, "stream", streamName, "error", werr)
			if rerr := p.esmStore.SetProcessingResult(mapping.UUID, werr.Error()); rerr != nil {
				logs.Warn("esm: failed to set state after Kinesis window error",
					logs.String("mapping", mapping.UUID), logs.Err(rerr))
			}
			return
		}
		if res.discarded {
			report.recordDiscard()
			if rerr := p.esmStore.SetProcessingResult(mapping.UUID, "Records discarded after exhausting retries"); rerr != nil {
				logs.Warn("esm: failed to set discard result", logs.String("mapping", mapping.UUID), logs.Err(rerr))
			}
			return
		}
		if res.processedAny {
			report.recordProcessed()
		}
	}
	windowCycle := func(cpKey, shardID string, items []windowedStreamItem, readThrough string) {
		src := streamSource{kind: streamSourceKinesis, streamArn: mapping.EventSourceArn, shardID: shardID}
		res, werr := p.processStreamWindow(ctx, mapping, cpKey, src, items, readThrough)
		reportWindow(res, werr)
	}

	// renderRecords builds the wire-format maps for one batch, recording
	// each record's arrival time for tumbling window boundaries.
	renderRecords := func(shardID string, records []eventbus.KinesisRecord) ([]map[string]interface{}, map[string]int64) {
		arrivals := make(map[string]int64, len(records))
		rendered := make([]map[string]interface{}, 0, len(records))
		for _, rec := range records {
			arrivals[rec.SequenceNumber] = rec.ApproximateArrivalTimestamp.Unix()
			record := map[string]interface{}{
				"kinesis": map[string]interface{}{
					"kinesisSchemaVersion":        "1.0",
					"partitionKey":                rec.PartitionKey,
					"sequenceNumber":              rec.SequenceNumber,
					"data":                        string(rec.Data),
					"approximateArrivalTimestamp": rec.ApproximateArrivalTimestamp.Format(timeutils.ISO8601UTCFormat),
				},
				"eventSource":       "aws:kinesis",
				"eventVersion":      "1.0",
				"eventID":           fmt.Sprintf("%s:%s:%s", shardID, rec.SequenceNumber, mapping.UUID),
				"awsRegion":         streamRegion,
				"eventName":         "aws:kinesis:record",
				"invokeIdentityArn": arnutil.NewARNBuilder(p.lambdaSvc.accountID, "").IAM().Role("vorpalstacks-lambda"),
			}
			if windowed {
				// The documented KinesisTimeWindowEvent carries the source
				// ARN on every record as well as at the envelope level.
				record["eventSourceARN"] = mapping.EventSourceArn
			}
			rendered = append(rendered, record)
		}
		return rendered, arrivals
	}

	// dropExpiredKinesis removes records older than
	// MaximumRecordAgeInSeconds (-1, the default, keeps every record) and
	// reports the expired batch to the on-failure destination: "Lambda
	// retries until the records expire, exceed the maximum age ... If the
	// error handling measures fail, Lambda discards the records".
	dropExpiredKinesis := func(shardID string, records []eventbus.KinesisRecord) []eventbus.KinesisRecord {
		if mapping.MaximumRecordAgeInSeconds <= 0 {
			return records
		}
		src := streamSource{kind: streamSourceKinesis, streamArn: mapping.EventSourceArn, shardID: shardID}
		cutoff := time.Now().Add(-time.Duration(mapping.MaximumRecordAgeInSeconds) * time.Second)
		fresh := records[:0]
		var expired []eventbus.KinesisRecord
		for _, rec := range records {
			if rec.ApproximateArrivalTimestamp.After(cutoff) {
				fresh = append(fresh, rec)
				continue
			}
			expired = append(expired, rec)
		}
		if len(expired) > 0 {
			rendered, _ := renderRecords(shardID, expired)
			items := make([]streamBatchItem, len(rendered))
			for i, kr := range rendered {
				items[i] = streamBatchItem{record: kr, seq: kinesisRecordSeq(kr)}
			}
			p.deliverDiscardedBatch(ctx, mapping, src, streamFailureBatchInfoOf(src, items),
				marshalStreamBatch(items), 0, uninvokedBatchResponse())
		}
		return fresh
	}

	// windowedItemsOf renders one batch and pairs the surviving records
	// with their tumbling window boundaries.
	windowedItemsOf := func(shardID string, records []eventbus.KinesisRecord) []windowedStreamItem {
		records = dropExpiredKinesis(shardID, records)
		rendered, arrivals := renderRecords(shardID, records)
		rendered = filterKinesisRecords(rendered, mapping.FilterCriteria)
		items := make([]streamBatchItem, len(rendered))
		for i, kr := range rendered {
			items[i] = streamBatchItem{record: kr, seq: kinesisRecordSeq(kr)}
		}
		witems := make([]windowedStreamItem, len(items))
		for i, it := range items {
			witems[i] = windowedStreamItem{
				item:        it,
				windowStart: windowStartOf(arrivals[it.seq], windowSeconds),
			}
		}
		return witems
	}

	// prepareItems filters and renders one fetched batch into delivery
	// items, after the age expiry pass, with event filter criteria
	// dropping non-matching records.
	prepareItems := func(shardID string, records []eventbus.KinesisRecord) []streamBatchItem {
		records = dropExpiredKinesis(shardID, records)
		if len(records) == 0 {
			return nil
		}
		rendered, _ := renderRecords(shardID, records)
		rendered = filterKinesisRecords(rendered, mapping.FilterCriteria)
		items := make([]streamBatchItem, len(rendered))
		for i, kr := range rendered {
			items[i] = streamBatchItem{record: kr, seq: kinesisRecordSeq(kr)}
		}
		return items
	}

	// processItems delivers prepared items with the retry and bisection
	// policy.
	processItems := func(src streamSource, items []streamBatchItem) batchOutcome {
		return p.processStreamBatch(ctx, mapping, src, items)
	}

	// processBatch runs the non-windowed, unbuffered delivery pipeline for
	// one batch. A batch whose records all expired or were all filtered out
	// still consumes its position.
	processBatch := func(src streamSource, records []eventbus.KinesisRecord) batchOutcome {
		latestSeqAll := records[len(records)-1].SequenceNumber
		items := prepareItems(src.shardID, records)
		if len(items) == 0 {
			// Every record expired or was filtered out; advance past the
			// whole batch to avoid re-reading the same records.
			return batchOutcome{lastConsumed: latestSeqAll}
		}
		return processItems(src, items)
	}

	fnName := arnutil.ExtractFunctionNameFromARN(mapping.FunctionArn)
	if fnName == "" {
		p.log("failed to extract function name from ARN", "arn", mapping.FunctionArn)
		return
	}

	for _, shard := range shards {
		cpKey := fmt.Sprintf("%s:%s:%s", mapping.UUID, streamName, shard.ShardID)
		src := streamSource{kind: streamSourceKinesis, streamArn: mapping.EventSourceArn, shardID: shard.ShardID}

		if shard.SequenceNumberRangeEnd != "" {
			// A closed shard never yields new records; deliver the final
			// invocation for any window still open on it and flush any
			// records its batching window was still gathering.
			if windowed {
				p.closeEndedShardWindow(ctx, mapping, cpKey, shard.ShardID)
			}
			if buffered {
				if items := p.dropStreamBuffer(cpKey); len(items) > 0 {
					if out := p.flushStreamBuffer(ctx, mapping, src, cpKey, "", items); out.err != nil {
						p.log("failed to flush batching window of ended shard",
							"mapping", mapping.UUID, "shard", shard.ShardID, "error", out.err)
					}
				}
			}
			continue
		}

		p.kinesisCPMu.RLock()
		lastSeq := p.kinesisCP[cpKey]
		p.kinesisCPMu.RUnlock()

		// While a tumbling window is open the read position advances per
		// delivered chunk; the durable checkpoint only moves when the window
		// completes. A batching-window buffer likewise reads ahead of the
		// checkpoint while records are being gathered.
		readFrom := lastSeq
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

		var iteratorType, iteratorSeqNum string
		var iteratorTimestamp *time.Time
		if readFrom != "" {
			iteratorType = "AFTER_SEQUENCE_NUMBER"
			iteratorSeqNum = readFrom
		} else {
			// Honour the user-configured StartingPosition for the initial read.
			// Default to TRIM_HORIZON for backward compatibility.
			switch mapping.StartingPosition {
			case "LATEST":
				iteratorType = "LATEST"
			case "AT_TIMESTAMP":
				iteratorType = "AT_TIMESTAMP"
				if !mapping.StartingPositionTimestamp.IsZero() {
					ts := mapping.StartingPositionTimestamp.UTC()
					iteratorTimestamp = &ts
				}
			default:
				iteratorType = "TRIM_HORIZON"
			}
		}

		iteratorSeq, err := p.bus.KinesisInvoker().CreateShardIterator(ctx, streamName, shard.ShardID, iteratorType, iteratorSeqNum, iteratorTimestamp)
		if err != nil {
			p.log("failed to create shard iterator", "stream", streamName, "shard", shard.ShardID, "error", err)
			continue
		}

		// Fetch up to ParallelizationFactor batches; a chained read picks
		// up strictly after the previous batch's last record. The first
		// read of a LATEST anchor includes the anchor record itself: the
		// mapping activates immediately on this platform, without the
		// activation lag that on AWS lets records published between
		// CreateEventSourceMapping and the iterator anchoring still be
		// delivered, so the anchor record — the stream's latest at first
		// poll — must be read inclusively.
		anchorInitialLATEST := readFrom == "" && iteratorType == "LATEST"
		var batches [][]eventbus.KinesisRecord
		pos := iteratorSeq
		for i := 0; i < pf; i++ {
			records, next, gerr := p.bus.KinesisInvoker().GetRecords(ctx, streamName, shard.ShardID, pos, batchSize, i == 0 && anchorInitialLATEST)
			if gerr != nil {
				p.log("failed to get records from Kinesis", "stream", streamName, "shard", shard.ShardID, "error", gerr)
				break
			}
			if len(records) == 0 {
				break
			}
			batches = append(batches, records)
			if int32(len(records)) < batchSize {
				break
			}
			pos = next
		}
		if len(batches) == 0 && readFrom == "" {
			// Persist the initial anchor so the next cycle reads strictly
			// after it instead of re-anchoring LATEST at a newer latest
			// record, which would skip everything that arrived in between.
			// Windowed and buffering mappings need this as well: until their
			// first window or buffer opens they hold no read position of
			// their own, so without a durable anchor every empty poll would
			// re-anchor LATEST and lose all but the newest record of an
			// inter-poll burst. An open window or gathering buffer always
			// carries a read position, so the condition cannot fire then.
			p.advanceCheckpoint(cpKey, iteratorSeq)
		}

		if windowed {
			// The tumbling window state is order-sensitive, so the fetched
			// batches flow through the window machinery sequentially.
			for _, records := range batches {
				windowCycle(cpKey, shard.ShardID, windowedItemsOf(shard.ShardID, records), records[len(records)-1].SequenceNumber)
			}
			if len(batches) == 0 {
				// No records this cycle, but the inactivity grace period may
				// have expired for an open window.
				windowCycle(cpKey, shard.ShardID, nil, "")
			}
			continue
		}

		if buffered {
			// Gather the surviving records of every fetched batch; the read
			// position moves with the buffer, the checkpoint only at flush.
			buf := p.getStreamBuffer(cpKey)
			for _, records := range batches {
				items := prepareItems(shard.ShardID, records)
				if len(items) == 0 {
					continue
				}
				if len(buf.items) == 0 {
					buf.firstAt = time.Now()
				}
				buf.items = append(buf.items, items...)
				buf.readThrough = records[len(records)-1].SequenceNumber
			}
			full, expired := bufferReady(len(buf.items), int(batchSize), buf.firstAt, time.Now(), batchingWindow)
			if len(buf.items) > 0 && !full && !expired {
				// The window is still gathering; hold the records.
				continue
			}
			items := p.dropStreamBuffer(cpKey)
			readThrough := buf.readThrough
			if len(items) == 0 {
				// Nothing survived, but the read position still moves.
				if readThrough != "" {
					p.advanceCheckpoint(cpKey, readThrough)
				}
				continue
			}
			outcome := p.flushStreamBuffer(ctx, mapping, src, cpKey, readThrough, items)
			if outcome.err != nil {
				report.recordFailure()
				p.log("lambda invocation failed for Kinesis ESM batching window",
					"function", fnName, "stream", streamName, "error", outcome.err)
				if rerr := p.esmStore.SetProcessingResult(mapping.UUID, outcome.err.Error()); rerr != nil {
					logs.Warn("esm: failed to set state after Kinesis window error",
						logs.String("mapping", mapping.UUID), logs.Err(rerr))
				}
				continue
			}
			if outcome.discarded {
				report.recordDiscard()
				p.log("discarded failed stream records after exhausting retries", "function", fnName, "stream", streamName)
				if rerr := p.esmStore.SetProcessingResult(mapping.UUID, "Records discarded after exhausting retries"); rerr != nil {
					logs.Warn("esm: failed to set discard result", logs.String("mapping", mapping.UUID), logs.Err(rerr))
				}
			} else if outcome.reported > 0 {
				report.recordPartial(outcome.reported)
			} else if outcome.delivered {
				report.recordProcessed()
			}
			continue
		}

		if len(batches) == 0 {
			continue
		}

		// "The number of batches to process from each shard concurrently":
		// batches run in parallel except where they share a partition key,
		// and the checkpoint takes the contiguous consumed prefix.
		keySets := make([]map[string]struct{}, len(batches))
		for i, records := range batches {
			keySets[i] = kinesisRecordKeys(records)
		}
		outcomes := runOrderedBatches(ctx, len(batches), keySets, func(ctx context.Context, idx int) batchOutcome {
			return processBatch(src, batches[idx])
		})
		lastConsumed, delivered, discarded, reported, failure := prefixOutcome(outcomes)
		if failure != nil {
			report.recordFailure()
			p.log("lambda invocation failed for Kinesis ESM", "function", fnName, "stream", streamName, "error", failure)
			if rerr := p.esmStore.SetProcessingResult(mapping.UUID, failure.Error()); rerr != nil {
				logs.Warn("esm: failed to set state after Kinesis invocation error",
					logs.String("mapping", mapping.UUID), logs.Err(rerr))
			}
		} else if discarded {
			report.recordDiscard()
			p.log("discarded failed stream records after exhausting retries", "function", fnName, "stream", streamName)
			if rerr := p.esmStore.SetProcessingResult(mapping.UUID, "Records discarded after exhausting retries"); rerr != nil {
				logs.Warn("esm: failed to set discard result", logs.String("mapping", mapping.UUID), logs.Err(rerr))
			}
		}
		if reported > 0 {
			report.recordPartial(reported)
		}
		if delivered {
			report.recordProcessed()
		}

		if lastConsumed == "" {
			continue
		}
		p.kinesisCPMu.Lock()
		p.kinesisCP[cpKey] = lastConsumed
		p.kinesisCPMu.Unlock()
		if err := p.persistKinesisCheckpoint(cpKey, lastConsumed); err != nil {
			logs.Warn("esm: failed to persist Kinesis checkpoint, in-memory state may diverge on restart",
				logs.String("key", cpKey), logs.Err(err))
		}
	}

	// Only a cycle that actually delivered records may report success, and
	// never over a failure, discard or partial batch response the same
	// cycle recorded: idle cycles leave the previous result untouched, a
	// discard stays visible even when a later batch of the same cycle
	// succeeded, and a partial response reports its records instead.
	if report.shouldReportSuccess() {
		if err := p.esmStore.SetProcessingResult(mapping.UUID, "No errors."); err != nil {
			logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.String("error", err.Error()))
		}
	} else if !report.failure && !report.discarded && report.partial > 0 {
		if err := p.esmStore.SetProcessingResult(mapping.UUID, streamPartialResult(report.partial)); err != nil {
			logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.Err(err))
		}
	}
}

// processDynamoDBStreamsMapping polls a DynamoDB stream for new records,
// batches them, and invokes the mapped Lambda function with the records
// in DynamoDB Streams event format. Checkpoints are persisted per mapping
// to survive restarts.
func (p *esmPoller) processDynamoDBStreamsMapping(ctx context.Context, mapping *lambdastore.EventSourceMapping) {
	if p.bus == nil {
		return
	}

	invoker := p.bus.DynamoDBStreamsInvoker()
	if invoker == nil {
		return
	}

	// Extract table name from the stream ARN:
	// arn:aws:dynamodb:region:account:table/tableName/stream/label
	_, _, streamRegion, _, resource := arnutil.SplitARN(mapping.EventSourceArn)
	if streamRegion == "" {
		p.log("failed to parse DynamoDB stream ARN", "arn", mapping.EventSourceArn)
		return
	}

	tableName := resource
	if idx := strings.Index(resource, "table/"); idx != -1 {
		rest := resource[idx+len("table/"):]
		if slashIdx := strings.Index(rest, "/"); slashIdx != -1 {
			tableName = rest[:slashIdx]
		} else {
			tableName = rest
		}
	}

	// Load or initialise the checkpoint for this mapping.
	checkpointKey := fmt.Sprintf("ddb:%s", mapping.UUID)
	p.kinesisCPMu.RLock()
	lastSeqStr := p.kinesisCP[checkpointKey]
	p.kinesisCPMu.RUnlock()

	windowed := mapping.TumblingWindowInSeconds > 0
	windowSeconds := int64(mapping.TumblingWindowInSeconds)

	// While a tumbling window is open the read position advances per
	// delivered chunk; the durable checkpoint only moves when the window
	// completes. A batching-window buffer likewise reads ahead of the
	// checkpoint while records are being gathered.
	shardID := invoker.ShardIDForStream(mapping.EventSourceArn)
	src := streamSource{kind: streamSourceDynamoDB, streamArn: mapping.EventSourceArn, shardID: shardID}
	buffered := !windowed && mapping.MaximumBatchingWindowInSeconds > 0
	readFrom := lastSeqStr
	if windowed {
		p.windowsMu.Lock()
		if w := p.windows[checkpointKey]; w != nil && w.readSeq != "" {
			readFrom = w.readSeq
		}
		p.windowsMu.Unlock()
	} else if buffered {
		if seq := p.streamBufferReadSeq(checkpointKey); seq != "" {
			readFrom = seq
		}
	}

	// Reads are exclusive of the read position, so the numeric sequence
	// of a record is a valid "read everything after this" cursor.
	var fromSeq int64
	anchorLATEST := int64(0)
	anchored := false
	if readFrom != "" {
		if v, err := strconv.ParseInt(readFrom, 10, 64); err == nil {
			fromSeq = v
		}
	} else {
		// No durable read position yet: honour the user-configured
		// StartingPosition for the initial read, mirroring the Kinesis
		// path. AT_TIMESTAMP never reaches the poller for DynamoDB —
		// mapping creation rejects it.
		switch mapping.StartingPosition {
		case "LATEST":
			latest, aerr := invoker.GetLatestSequence(ctx, streamRegion, tableName)
			if aerr != nil {
				p.log("failed to read latest DynamoDB stream sequence",
					"table", tableName, "error", aerr.Error())
				return
			}
			anchorLATEST = latest
			anchored = true
			// The anchor is read inclusively: the mapping activates
			// immediately on this platform, without the activation lag
			// that on AWS lets records published between
			// CreateEventSourceMapping and the iterator anchoring still
			// be delivered, so the anchor record — the stream's latest
			// at first poll — must be delivered too.
			fromSeq = latest - 1
			if fromSeq < 0 {
				fromSeq = 0
			}
		default:
			// TRIM_HORIZON (or unset): read from the beginning of the
			// stream. Sequence numbers are assigned from one upwards, so
			// zero reads every record.
			fromSeq = 0
		}
	}

	batchSize := int32(mapping.BatchSize)
	if batchSize <= 0 {
		batchSize = 100
	}

	// reportWindow mirrors the Kinesis path's status reporting for the
	// windowed flow; the non-windowed flow keeps its own inline reporting.
	// The success write is deferred to the end of the poll behind the
	// cycle precedence, so a batch discarded earlier in the same cycle
	// cannot be overwritten by a later successful one.
	report := cycleReport{}
	reportWindow := func(res windowCycleResult, werr error) {
		if werr != nil {
			report.recordFailure()
			p.log("failed to invoke function for DynamoDB ESM tumbling window",
				"function", mapping.FunctionArn, "error", werr.Error())
			if rerr := p.esmStore.SetProcessingResult(mapping.UUID, werr.Error()); rerr != nil {
				logs.Warn("esm: failed to set result after DynamoDB window error",
					logs.String("mapping", mapping.UUID), logs.Err(rerr))
			}
			return
		}
		if res.discarded {
			report.recordDiscard()
			if rerr := p.esmStore.SetProcessingResult(mapping.UUID, "Records discarded after exhausting retries"); rerr != nil {
				logs.Warn("esm: failed to set discard result", logs.String("mapping", mapping.UUID), logs.Err(rerr))
			}
			return
		}
		if res.processedAny {
			report.recordProcessed()
		}
	}
	windowCycle := func(items []windowedStreamItem, readThrough string) {
		res, werr := p.processStreamWindow(ctx, mapping, checkpointKey, src, items, readThrough)
		reportWindow(res, werr)
	}

	pf := parallelizationFactorOf(mapping)

	// Fetch up to ParallelizationFactor batches; a chained read picks up
	// strictly after the previous batch's last record. The invoker answers
	// with the last record's sequence, which is a valid read position
	// because reads are exclusive of it.
	type ddbBatch struct {
		records []eventbus.DynamoDBStreamRecord
		nextSeq int64
	}
	var batches []ddbBatch
	from := fromSeq
	for i := 0; i < pf; i++ {
		records, nextSeq, gerr := invoker.GetRecords(ctx, streamRegion, tableName, from, int(batchSize))
		if gerr != nil {
			p.log("failed to get DynamoDB stream records",
				"table", tableName, "error", gerr.Error())
			if len(batches) == 0 {
				return
			}
			break
		}
		if len(records) == 0 {
			break
		}
		batches = append(batches, ddbBatch{records: records, nextSeq: nextSeq})
		if int32(len(records)) < batchSize {
			break
		}
		from = nextSeq
	}

	if len(batches) == 0 && anchored {
		// Persist the initial anchor so the next cycle reads strictly
		// after it instead of re-anchoring LATEST at a newer latest
		// record, which would skip everything that arrived in between —
		// the same inter-poll burst protection the Kinesis path applies.
		// A cycle that returned records advances the checkpoint through
		// the normal delivery flow, which never moves it backwards.
		p.advanceCheckpoint(checkpointKey, strconv.FormatInt(anchorLATEST, 10))
	}

	if windowed {
		// The tumbling window state is order-sensitive, so the fetched
		// batches flow through the window machinery sequentially.
		for _, b := range batches {
			readThrough := strconv.FormatInt(b.nextSeq, 10)

			// Apply age expiry and event filtering before invocation.
			records := p.discardExpiredDynamoDBRecords(ctx, mapping, src, b.records)
			records = filterDynamoDBRecords(records, mapping.FilterCriteria)
			if len(records) == 0 {
				// All records were filtered out; advance the read position
				// to avoid re-reading them on the next poll.
				windowCycle(nil, readThrough)
				continue
			}

			witems := make([]windowedStreamItem, len(records))
			for i := range records {
				witems[i] = windowedStreamItem{
					item:        streamBatchItem{record: &records[i], seq: dynamoDBRecordSeq(&records[i])},
					windowStart: windowStartOf(ddbArrivalUnix(&records[i]), windowSeconds),
				}
			}
			windowCycle(witems, readThrough)
		}
		if len(batches) == 0 {
			// No records this cycle, but the inactivity grace period may
			// have expired for an open window.
			windowCycle(nil, "")
		}
		// Only a cycle without a failure or discard that delivered records
		// may report success; the deferred write keeps a discard recorded
		// earlier in the cycle visible.
		if report.shouldReportSuccess() {
			if err := p.esmStore.SetProcessingResult(mapping.UUID, "No errors."); err != nil {
				logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.Err(err))
			}
		}
		return
	}

	if buffered {
		// Gather the surviving records of every fetched batch; the read
		// position moves with the buffer, the checkpoint only at flush.
		buf := p.getStreamBuffer(checkpointKey)
		for _, b := range batches {
			records := p.discardExpiredDynamoDBRecords(ctx, mapping, src, b.records)
			records = filterDynamoDBRecords(records, mapping.FilterCriteria)
			if len(records) == 0 {
				continue
			}
			items := make([]streamBatchItem, len(records))
			for i := range records {
				items[i] = streamBatchItem{record: &records[i], seq: dynamoDBRecordSeq(&records[i])}
			}
			if len(buf.items) == 0 {
				buf.firstAt = time.Now()
			}
			buf.items = append(buf.items, items...)
			buf.readThrough = strconv.FormatInt(b.nextSeq, 10)
		}
		full, expired := bufferReady(len(buf.items), int(batchSize), buf.firstAt, time.Now(), mapping.MaximumBatchingWindowInSeconds)
		if len(buf.items) > 0 && !full && !expired {
			// The window is still gathering; hold the records.
			return
		}
		items := p.dropStreamBuffer(checkpointKey)
		readThrough := buf.readThrough
		if len(items) == 0 {
			// Nothing survived, but the read position still moves.
			if readThrough != "" {
				p.advanceCheckpoint(checkpointKey, readThrough)
			}
			return
		}
		outcome := p.flushStreamBuffer(ctx, mapping, src, checkpointKey, readThrough, items)
		if outcome.err != nil {
			p.log("failed to invoke function for DynamoDB ESM batching window",
				"function", mapping.FunctionArn, "error", outcome.err.Error())
			if rerr := p.esmStore.SetProcessingResult(mapping.UUID, outcome.err.Error()); rerr != nil {
				logs.Warn("esm: failed to set result after DynamoDB window error",
					logs.String("mapping", mapping.UUID), logs.Err(rerr))
			}
			return
		}
		if outcome.discarded {
			if rerr := p.esmStore.SetProcessingResult(mapping.UUID, "Records discarded after exhausting retries"); rerr != nil {
				logs.Warn("esm: failed to set discard result", logs.String("mapping", mapping.UUID), logs.Err(rerr))
			}
			return
		}
		if outcome.reported > 0 {
			if err := p.esmStore.SetProcessingResult(mapping.UUID, streamPartialResult(outcome.reported)); err != nil {
				logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.Err(err))
			}
			return
		}
		if outcome.delivered {
			if err := p.esmStore.SetProcessingResult(mapping.UUID, "No errors."); err != nil {
				logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.Err(err))
			}
		}
		return
	}

	if len(batches) == 0 {
		return
	}

	// processBatch runs the non-windowed delivery pipeline for one batch;
	// a batch whose records were all filtered out still consumes its
	// position.
	processBatch := func(b ddbBatch) batchOutcome {
		records := p.discardExpiredDynamoDBRecords(ctx, mapping, src, b.records)
		records = filterDynamoDBRecords(records, mapping.FilterCriteria)
		if len(records) == 0 {
			return batchOutcome{lastConsumed: strconv.FormatInt(b.nextSeq, 10)}
		}

		items := make([]streamBatchItem, len(records))
		for i := range records {
			items[i] = streamBatchItem{record: &records[i], seq: dynamoDBRecordSeq(&records[i])}
		}
		return p.processStreamBatch(ctx, mapping, src, items)
	}

	// Batches run concurrently except where they share an item key; the
	// checkpoint takes the contiguous consumed prefix.
	keySets := make([]map[string]struct{}, len(batches))
	for i, b := range batches {
		keySets[i] = dynamoDBRecordKeys(b.records)
	}
	outcomes := runOrderedBatches(ctx, len(batches), keySets, func(ctx context.Context, idx int) batchOutcome {
		return processBatch(batches[idx])
	})
	lastConsumed, delivered, discarded, reported, failure := prefixOutcome(outcomes)
	if failure != nil {
		p.log("failed to invoke function for DynamoDB ESM",
			"function", mapping.FunctionArn, "error", failure.Error())
		if rerr := p.esmStore.SetProcessingResult(mapping.UUID, failure.Error()); rerr != nil {
			logs.Warn("esm: failed to set result after DynamoDB invocation error",
				logs.String("mapping", mapping.UUID), logs.Err(rerr))
		}
	} else if discarded {
		p.log("discarded failed stream records after exhausting retries",
			"function", mapping.FunctionArn)
		if rerr := p.esmStore.SetProcessingResult(mapping.UUID, "Records discarded after exhausting retries"); rerr != nil {
			logs.Warn("esm: failed to set discard result", logs.String("mapping", mapping.UUID), logs.Err(rerr))
		}
	}

	// Advance the checkpoint to the contiguous consumed prefix (the full
	// batch on success, the prefix before a failure or a partial batch
	// response).
	if lastConsumed != "" {
		p.kinesisCPMu.Lock()
		p.kinesisCP[checkpointKey] = lastConsumed
		p.kinesisCPMu.Unlock()
		if err := p.persistKinesisCheckpoint(checkpointKey, lastConsumed); err != nil {
			logs.Warn("esm: failed to persist DynamoDB checkpoint", logs.Err(err))
		}
	}

	if failure == nil && !discarded && delivered {
		result := "No errors."
		if reported > 0 {
			result = streamPartialResult(reported)
		}
		if err := p.esmStore.SetProcessingResult(mapping.UUID, result); err != nil {
			logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.Err(err))
		}
	}
}

func (p *esmPoller) processSQSMapping(ctx context.Context, mapping *lambdastore.EventSourceMapping) {
	if p.bus == nil {
		return
	}

	_, _, region, accountID, _ := arnutil.SplitARN(mapping.EventSourceArn)
	if region == "" || accountID == "" {
		p.log("failed to parse event source ARN", "arn", mapping.EventSourceArn)
		return
	}

	queueName := arnutil.ExtractQueueNameFromARN(mapping.EventSourceArn)
	if queueName == "" {
		p.log("failed to extract queue name from ARN", "arn", mapping.EventSourceArn)
		return
	}

	queueURL, err := p.bus.SQSInvoker().GetQueueByName(ctx, region, queueName)
	if err != nil {
		p.log("sqs queue not found by name", "queue", queueName, "mapping", mapping.UUID, "error", err)
		return
	}

	batchSize := mapping.BatchSize
	if batchSize <= 0 {
		batchSize = 100
	}
	if batchSize > maxSQSBatchSize {
		batchSize = maxSQSBatchSize
	}

	perCallMax := sqsReceiveMessageMax
	if perCallMax > batchSize {
		perCallMax = batchSize
	}

	// The receive wait is bounded to one poll interval: the poll cycle runs
	// synchronously (pollAll waits for every mapping to finish), so a
	// receive that long-polls longer than the interval would stall the
	// whole cycle and delay every other mapping's poll. Long polling lives
	// in the SQS store; the poller only needs each cycle to stay
	// responsive, and the batching window still governs record gathering.
	sqsReceiveWaitSeconds := int32(p.interval / time.Second)
	if sqsReceiveWaitSeconds < 1 {
		sqsReceiveWaitSeconds = 1
	}

	var allMessages []eventbus.ReceivedSQSMessage
	remaining := batchSize
	for remaining > 0 {
		fetchCount := perCallMax
		if fetchCount > remaining {
			fetchCount = remaining
		}

		msgs, err := p.bus.SQSInvoker().ReceiveMessage(ctx, region, queueURL, fetchCount, nil, sqsReceiveWaitSeconds)
		if err != nil {
			p.log("sqs receive failed", "queue", queueName, "mapping", mapping.UUID, "error", err)
			break
		}

		if len(msgs) == 0 {
			break
		}

		allMessages = append(allMessages, msgs...)
		remaining -= int32(len(msgs))

		if int32(len(msgs)) < fetchCount {
			break
		}
	}

	if len(allMessages) == 0 {
		return
	}

	records := make([]esmSQSRecord, 0, len(allMessages))
	receiptHandles := make([]string, 0, len(allMessages))
	messageIDs := make([]string, 0, len(allMessages))
	for _, msg := range allMessages {
		records = append(records, receivedSQSMessageToRecord(msg, mapping.EventSourceArn, region))
		receiptHandles = append(receiptHandles, msg.ReceiptHandle)
		messageIDs = append(messageIDs, msg.MessageID)
	}

	// Apply event filtering before invocation.
	records = filterSQSRecords(records, mapping.FilterCriteria)
	if len(records) == 0 {
		// All messages were filtered out.  AWS treats filtered-out
		// messages as successfully processed — delete them from the
		// queue to prevent infinite re-polling after visibility
		// timeout expires.
		for _, handle := range receiptHandles {
			if err := p.bus.SQSInvoker().DeleteMessage(ctx, region, queueURL, handle); err != nil {
				p.log("failed to delete filtered-out message", "queue", queueName, "error", err)
			}
		}
		if err := p.esmStore.SetProcessingResult(mapping.UUID, "No errors."); err != nil {
			logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.String("error", err.Error()))
		}
		return
	}

	eventPayload := esmSQSEvent{Records: records}
	payload, err := json.Marshal(eventPayload)
	if err != nil {
		p.log("failed to marshal ESM event payload", "queue", queueName, "error", err)
		return
	}

	fnName := arnutil.ExtractFunctionNameFromARN(mapping.FunctionArn)
	if fnName == "" {
		p.log("failed to extract function name from ARN", "arn", mapping.FunctionArn)
		return
	}

	var report batchResponseReport
	invokeErr := p.invokeWithRetry(ctx, mapping, payload, batchResponseSink(mapping, &report))

	if invokeErr != nil {
		p.log("lambda invocation failed", "function", fnName, "queue", queueName, "error", invokeErr)
		if err := p.esmStore.SetProcessingResult(mapping.UUID, invokeErr.Error()); err != nil {
			logs.Warn("esm: failed to set state after SQS invocation error",
				logs.String("mapping", mapping.UUID), logs.Err(err))
		}
		return
	}

	// A partial batch response deletes only the messages the function did
	// not report: "To make messages id2 and id4 visible again in your
	// queue, your function should return" their identifiers — the reported
	// messages return with the queue's visibility timeout.
	deleteFailures := 0
	for i, handle := range receiptHandles {
		if _, failed := report.failedIDs[messageIDs[i]]; failed {
			continue
		}
		if err := p.bus.SQSInvoker().DeleteMessage(ctx, region, queueURL, handle); err != nil {
			p.log("failed to delete message", "queue", queueName, "error", err)
			deleteFailures++
		}
	}

	lastResult := "No errors."
	if deleteFailures > 0 {
		lastResult = fmt.Sprintf("%d message(s) failed to delete", deleteFailures)
	} else if reported := reportedSQSFailureCount(messageIDs, report); reported > 0 {
		lastResult = sqsPartialResult(reported)
	}

	if err := p.esmStore.SetProcessingResult(mapping.UUID, lastResult); err != nil {
		logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.String("error", err.Error()))
	}
}

// receivedSQSMessageToRecord converts an eventbus.ReceivedSQSMessage into an
// ESM SQS record matching the Lambda event format.
func receivedSQSMessageToRecord(msg eventbus.ReceivedSQSMessage, eventSourceArn, region string) esmSQSRecord {
	record := esmSQSRecord{
		MessageID:               msg.MessageID,
		ReceiptHandle:           msg.ReceiptHandle,
		Body:                    msg.Body,
		MD5OfBody:               msg.MD5OfBody,
		MD5OfMessageAttributes:  msg.MD5OfMessageAttributes,
		EventSourceARN:          eventSourceArn,
		EventSource:             "aws:sqs",
		AWSRegion:               region,
		ApproximateReceiveCount: fmt.Sprintf("%d", msg.ApproximateReceiveCount),
		SentTimestamp:           fmt.Sprintf("%d", msg.SentTimestamp.UnixMilli()),
	}

	if msg.ApproximateFirstReceiveTimestamp.IsZero() {
		record.ApproximateFirstReceiveTimestamp = record.SentTimestamp
	} else {
		record.ApproximateFirstReceiveTimestamp = fmt.Sprintf("%d", msg.ApproximateFirstReceiveTimestamp.UnixMilli())
	}

	if msg.SequenceNumber != "" {
		if record.MessageAttributes == nil {
			record.MessageAttributes = make(map[string]interface{})
		}
		record.MessageAttributes["SequenceNumber"] = map[string]string{
			"stringValue": msg.SequenceNumber,
			"dataType":    "String",
		}
	}

	if msg.MessageDeduplicationID != "" {
		if record.MessageAttributes == nil {
			record.MessageAttributes = make(map[string]interface{})
		}
		record.MessageAttributes["MessageDeduplicationId"] = map[string]string{
			"stringValue": msg.MessageDeduplicationID,
			"dataType":    "String",
		}
	}

	if msg.MessageGroupID != "" {
		if record.MessageAttributes == nil {
			record.MessageAttributes = make(map[string]interface{})
		}
		record.MessageAttributes["MessageGroupId"] = map[string]string{
			"stringValue": msg.MessageGroupID,
			"dataType":    "String",
		}
	}

	return record
}

// invokeLambda invokes the Lambda function with the given name and payload.
// It delegates to the event-source invoke path, which resolves name/ARN
// reference forms (including alias qualifiers), the target region, and —
// unlike the gateway invoke surface — reports the function-level error
// classification so failed batches are not acknowledged as successes.
func (p *esmPoller) invokeLambda(ctx context.Context, functionRef string, payload []byte) (*lambdastore.InvocationResult, error) {
	if p.invoke != nil {
		return p.invoke(ctx, functionRef, payload)
	}
	if p.lambdaSvc == nil {
		return nil, fmt.Errorf("esm: lambda service not available")
	}
	return p.lambdaSvc.InvokeForEventSource(ctx, functionRef, payload)
}

// purgeStaleKinesisCheckpoints removes checkpoint entries for stream ESM
// mappings — Kinesis and DynamoDB Streams alike — that no longer exist or
// are not in the enabled state.
func (p *esmPoller) purgeStaleKinesisCheckpoints(activeUUIDs map[string]struct{}) {
	p.kinesisCPMu.Lock()
	for key := range p.kinesisCP {
		// Both key forms share the map: "<uuid>:<stream>:<shard>" for
		// Kinesis and "ddb:<uuid>" for DynamoDB Streams. The mapping UUID
		// sits before the first colon either way.
		uuid := strings.TrimPrefix(key, "ddb:")
		if idx := strings.IndexByte(uuid, ':'); idx >= 0 {
			uuid = uuid[:idx]
		}
		if uuid == "" {
			continue
		}
		if _, active := activeUUIDs[uuid]; !active {
			delete(p.kinesisCP, key)
			if err := p.deleteKinesisCheckpoint(key); err != nil {
				logs.Warn("esm: failed to delete stale stream checkpoint from persistence",
					logs.String("key", key), logs.Err(err))
			}
		}
	}
	p.kinesisCPMu.Unlock()
}

func (p *esmPoller) loadKinesisCheckpoints() {
	bucket := p.checkpointBucket()
	if bucket == nil {
		return
	}
	p.kinesisCPMu.Lock()
	if err := bucket.ForEach(func(k, v []byte) error {
		p.kinesisCP[string(k)] = string(v)
		return nil
	}); err != nil {
		logs.Error("esm: failed to load Kinesis checkpoints from persistence", logs.Err(err))
	}
	p.kinesisCPMu.Unlock()
}

func (p *esmPoller) persistKinesisCheckpoint(cpKey, seqNum string) error {
	bucket := p.checkpointBucket()
	if bucket == nil {
		return fmt.Errorf("checkpoint bucket unavailable")
	}
	return bucket.Put([]byte(cpKey), []byte(seqNum))
}

func (p *esmPoller) deleteKinesisCheckpoint(cpKey string) error {
	bucket := p.checkpointBucket()
	if bucket == nil {
		return fmt.Errorf("checkpoint bucket unavailable")
	}
	return bucket.Delete([]byte(cpKey))
}

func (p *esmPoller) checkpointBucket() storage.Bucket {
	if p.storageManager == nil {
		return nil
	}
	st, err := p.storageManager.GetStorage(p.region)
	if err != nil {
		return nil
	}
	return st.Bucket(esmCheckpointBucket)
}

// log emits a structured log message if a logger is configured on the
// poller.
func (p *esmPoller) log(msg string, keyvals ...interface{}) {
	// Fall back to the package logger: a nil logger would otherwise drop
	// every poller diagnostic, leaving event source delivery failures
	// invisible in the server log.
	target := p.logger
	if target == nil {
		target = logs.Default()
	}
	fields := make([]logs.Field, 0, len(keyvals)/2)
	for i := 0; i+1 < len(keyvals); i += 2 {
		fields = append(fields, logs.Field{Key: fmt.Sprint(keyvals[i]), Value: keyvals[i+1]})
	}
	target.Info(msg, fields...)
}
