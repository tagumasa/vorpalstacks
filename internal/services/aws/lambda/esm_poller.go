package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

const (
	// defaultESMPollInterval is the default interval between polling cycles
	// for SQS event source mappings. Matches the AWS minimum poll interval
	// for Lambda ESM.
	defaultESMPollInterval = 1 * time.Second

	// defaultESMMaxWorkers is the maximum number of concurrent polling
	// goroutines, one per unique event source mapping.
	defaultESMMaxWorkers = 32

	// defaultSQSWaitTimeSeconds is the SQS ReceiveMessage WaitTimeSeconds
	// parameter. Using long polling (20s) reduces empty responses and cost.
	defaultSQSWaitTimeSeconds = int32(20)

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
	accountID      string
	region         string
	storageManager *storage.RegionStorageManager
	kinesisCP      map[string]string // "mappingUUID:streamName:shardID" -> lastSeqNum
	kinesisCPMu    sync.RWMutex
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

	for _, region := range regions {
		p.pollRegion(ctx, region)
	}
}

func (p *esmPoller) pollRegion(ctx context.Context, region string) {
	var esmStore *lambdastore.EventSourceStore
	if region == p.region && p.esmStore != nil {
		esmStore = p.esmStore
	} else if p.storageManager != nil {
		st, err := p.storageManager.GetStorage(region)
		if err != nil {
			return
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
				EventSources: lambdastore.NewEventSourceStore(st, p.accountID, region),
			}
			if actual, loaded := p.lambdaSvc.storeCache.LoadOrStore(region, newStore); loaded {
				esmStore = actual.(*lambdaStore).EventSources
			} else {
				esmStore = newStore.EventSources
			}
		}
	}
	if esmStore == nil {
		return
	}

	result, err := esmStore.List(storecommon.ListOptions{MaxItems: 1000})
	if err != nil {
		p.log("failed to list event source mappings", "error", err)
		return
	}

	type pollJob struct {
		mapping *lambdastore.EventSourceMapping
	}
	jobs := make(chan pollJob, len(result.Items))
	activeKinesisUUIDs := make(map[string]struct{})
	for _, m := range result.Items {
		if m.State != "Enabled" {
			continue
		}
		_, esmService, _, _, _ := arnutil.SplitARN(m.EventSourceArn)
		if esmService != "sqs" && esmService != "kinesis" && esmService != "dynamodb" {
			continue
		}
		if esmService == "kinesis" {
			activeKinesisUUIDs[m.UUID] = struct{}{}
		}
		jobs <- pollJob{mapping: m}
	}
	close(jobs)

	p.purgeStaleKinesisCheckpoints(activeKinesisUUIDs)

	jobCount := len(result.Items)
	workerCount := p.workers
	if workerCount > jobCount {
		workerCount = jobCount
	}
	if workerCount == 0 {
		return
	}

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
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
		p.log("DynamoDB Streams ESM not yet implemented, skipping mapping",
			"eventSourceArn", mapping.EventSourceArn,
			"functionArn", mapping.FunctionArn)
		return
	}
	p.processSQSMapping(ctx, mapping)
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

	for _, shard := range shards {
		if shard.SequenceNumberRangeEnd != "" {
			continue
		}

		cpKey := fmt.Sprintf("%s:%s:%s", mapping.UUID, streamName, shard.ShardID)

		p.kinesisCPMu.RLock()
		lastSeq := p.kinesisCP[cpKey]
		p.kinesisCPMu.RUnlock()

		var iteratorType, iteratorSeqNum string
		if lastSeq != "" {
			iteratorType = "AFTER_SEQUENCE_NUMBER"
			iteratorSeqNum = lastSeq
		} else {
			iteratorType = "TRIM_HORIZON"
		}

		iteratorSeq, err := p.bus.KinesisInvoker().CreateShardIterator(ctx, streamName, shard.ShardID, iteratorType, iteratorSeqNum)
		if err != nil {
			p.log("failed to create shard iterator", "stream", streamName, "shard", shard.ShardID, "error", err)
			continue
		}

		records, _, err := p.bus.KinesisInvoker().GetRecords(ctx, streamName, shard.ShardID, iteratorSeq, batchSize)
		if err != nil {
			p.log("failed to get records from Kinesis", "stream", streamName, "shard", shard.ShardID, "error", err)
			continue
		}

		if len(records) == 0 {
			continue
		}

		kinesisRecords := make([]map[string]interface{}, 0, len(records))
		for _, rec := range records {
			kinesisRecords = append(kinesisRecords, map[string]interface{}{
				"kinesis": map[string]interface{}{
					"kinesisSchemaVersion":        "1.0",
					"partitionKey":                rec.PartitionKey,
					"sequenceNumber":              rec.SequenceNumber,
					"data":                        string(rec.Data),
					"approximateArrivalTimestamp": rec.ApproximateArrivalTimestamp.Format("2006-01-02T15:04:05.000Z"),
				},
				"eventSource":       "aws:kinesis",
				"eventVersion":      "1.0",
				"eventID":           fmt.Sprintf("%s:%s:%s", shard.ShardID, rec.SequenceNumber, mapping.UUID),
				"awsRegion":         streamRegion,
				"eventName":         "aws:kinesis:record",
				"invokeIdentityArn": arnutil.NewARNBuilder(p.accountID, "").IAM().Role("vorpalstacks-lambda"),
			})
		}

		eventPayload := map[string]interface{}{
			"Records": kinesisRecords,
		}
		payload, err := json.Marshal(eventPayload)
		if err != nil {
			p.log("failed to marshal Kinesis ESM payload", "stream", streamName, "error", err)
			continue
		}

		fnName := arnutil.ExtractFunctionNameFromARN(mapping.FunctionArn)
		if fnName == "" {
			p.log("failed to extract function name from ARN", "arn", mapping.FunctionArn)
			continue
		}

		_, _, invokeErr := p.invokeLambda(ctx, mapping.FunctionArn, payload)
		if invokeErr != nil {
			p.log("lambda invocation failed for Kinesis ESM", "function", fnName, "stream", streamName, "error", invokeErr)
			if err := p.esmStore.SetState(mapping.UUID, "Enabled", fmt.Sprintf("Last processing result: %s", invokeErr.Error())); err != nil {
				logs.Warn("esm: failed to set state after Kinesis invocation error",
					logs.String("mapping", mapping.UUID), logs.Err(err))
			}
			continue
		}

		latestSeq := records[len(records)-1].SequenceNumber
		p.kinesisCPMu.Lock()
		p.kinesisCP[cpKey] = latestSeq
		p.kinesisCPMu.Unlock()
		if err := p.persistKinesisCheckpoint(cpKey, latestSeq); err != nil {
			logs.Warn("esm: failed to persist Kinesis checkpoint, in-memory state may diverge on restart",
				logs.String("key", cpKey), logs.Err(err))
		}
	}

	if err := p.esmStore.SetState(mapping.UUID, "Enabled", "Last processing result: No errors."); err != nil {
		logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.String("error", err.Error()))
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

	queueURL, err := p.bus.SQSInvoker().GetQueueByName(ctx, queueName)
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

	maxMessages := batchSize
	if maxMessages > sqsReceiveMessageMax {
		maxMessages = sqsReceiveMessageMax
	}

	waitTime := defaultSQSWaitTimeSeconds
	if mapping.MaximumBatchingWindowInSeconds > 0 && mapping.MaximumBatchingWindowInSeconds < waitTime {
		waitTime = mapping.MaximumBatchingWindowInSeconds
	}

	messages, err := p.bus.SQSInvoker().ReceiveMessage(ctx, queueURL, maxMessages, nil, waitTime)
	if err != nil {
		p.log("sqs receive failed", "queue", queueName, "mapping", mapping.UUID, "error", err)
		return
	}

	if len(messages) == 0 {
		return
	}

	records := make([]esmSQSRecord, 0, len(messages))
	receiptHandles := make([]string, 0, len(messages))
	for _, msg := range messages {
		records = append(records, receivedSQSMessageToRecord(msg, mapping.EventSourceArn, region))
		receiptHandles = append(receiptHandles, msg.ReceiptHandle)
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

	_, _, invokeErr := p.invokeLambda(ctx, mapping.FunctionArn, payload)

	if invokeErr != nil {
		p.log("lambda invocation failed", "function", fnName, "queue", queueName, "error", invokeErr)
		if err := p.esmStore.SetState(mapping.UUID, "Enabled", fmt.Sprintf("Last processing result: %s", invokeErr.Error())); err != nil {
			logs.Warn("esm: failed to set state after SQS invocation error",
				logs.String("mapping", mapping.UUID), logs.Err(err))
		}
		return
	}

	for _, handle := range receiptHandles {
		if err := p.bus.SQSInvoker().DeleteMessage(ctx, queueURL, handle); err != nil {
			p.log("failed to delete message", "queue", queueName, "error", err)
		}
	}

	if err := p.esmStore.SetState(mapping.UUID, "Enabled", "Last processing result: No errors."); err != nil {
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
// It constructs the necessary store context and delegates to the internal
// invokeFunction method on the LambdaService.
func (p *esmPoller) invokeLambda(ctx context.Context, functionRef string, payload []byte) (int64, []byte, error) {
	if p.lambdaSvc == nil {
		return 0, nil, fmt.Errorf("esm: lambda service not available")
	}

	region := p.region
	functionName := functionRef
	if strings.HasPrefix(functionRef, "arn:") {
		if _, _, arnRegion, _, _ := arnutil.SplitARN(functionRef); arnRegion != "" {
			region = arnRegion
		}
		functionName = arnutil.ExtractFunctionNameFromARN(functionRef)
	}

	fnStore := p.lambdaSvc.getOrCreateFunctionStore(region)
	function, ver, _, err := p.lambdaSvc.resolveQualifier(fnStore, functionName, "")
	if err != nil {
		return 0, nil, fmt.Errorf("esm: failed to resolve function %s: %w", functionName, err)
	}

	result, err := p.lambdaSvc.invokeFunction(function, ver, fnStore, region, payload)
	if err != nil {
		return 0, nil, fmt.Errorf("esm: invocation failed for %s: %w", functionName, err)
	}
	if result == nil {
		return 0, nil, fmt.Errorf("esm: invocation returned nil result for %s", functionName)
	}

	return result.StatusCode, result.Payload, nil
}

// purgeStaleKinesisCheckpoints removes checkpoint entries for Kinesis ESM
// mappings that no longer exist or are not in the enabled state.
func (p *esmPoller) purgeStaleKinesisCheckpoints(activeKinesisUUIDs map[string]struct{}) {
	p.kinesisCPMu.Lock()
	for key := range p.kinesisCP {
		uuidEnd := strings.IndexByte(key, ':')
		if uuidEnd < 0 {
			continue
		}
		uuid := key[:uuidEnd]
		if _, active := activeKinesisUUIDs[uuid]; !active {
			delete(p.kinesisCP, key)
			if err := p.deleteKinesisCheckpoint(key); err != nil {
				logs.Warn("esm: failed to delete stale Kinesis checkpoint from persistence",
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
	if p.logger == nil {
		return
	}
	fields := make([]logs.Field, 0, len(keyvals)/2)
	for i := 0; i+1 < len(keyvals); i += 2 {
		fields = append(fields, logs.Field{Key: fmt.Sprint(keyvals[i]), Value: keyvals[i+1]})
	}
	p.logger.Info(msg, fields...)
}
