package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
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
	mu        sync.Mutex
	running   bool
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	interval  time.Duration
	workers   int
	logger    logs.Logger
	sqsStore  *sqsstore.SQSStore
	esmStore  *lambdastore.EventSourceStore
	lambdaSvc *LambdaService
	accountID string
	region    string
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
		interval: interval,
		workers:  workers,
		logger:   logger,
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
	if p.esmStore == nil {
		return
	}

	result, err := p.esmStore.List(lambdastore.ListOptions{MaxItems: 1000})
	if err != nil {
		p.log("failed to list event source mappings", "error", err)
		return
	}

	type pollJob struct {
		mapping *lambdastore.EventSourceMapping
	}
	jobs := make(chan pollJob, len(result.Items))
	for _, m := range result.Items {
		if m.State != "Enabled" {
			continue
		}
		if !strings.HasPrefix(m.EventSourceArn, "arn:aws:sqs:") {
			continue
		}
		jobs <- pollJob{mapping: m}
	}
	close(jobs)

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
	if p.sqsStore == nil {
		return
	}

	_, _, region, accountID, _ := svcarn.SplitARN(mapping.EventSourceArn)
	if region == "" || accountID == "" {
		p.log("failed to parse event source ARN", "arn", mapping.EventSourceArn)
		return
	}

	queueName := svcarn.ExtractQueueNameFromARN(mapping.EventSourceArn)
	if queueName == "" {
		p.log("failed to extract queue name from ARN", "arn", mapping.EventSourceArn)
		return
	}

	queueURL := fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", region, accountID, queueName)

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

	messages, err := p.sqsStore.ReceiveMessage(queueURL, maxMessages, nil, waitTime)
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
		records = append(records, sqsMessageToRecord(msg, mapping.EventSourceArn, region))
		receiptHandles = append(receiptHandles, msg.ReceiptHandle)
	}

	eventPayload := esmSQSEvent{Records: records}
	payload, err := json.Marshal(eventPayload)
	if err != nil {
		p.log("failed to marshal ESM event payload", "queue", queueName, "error", err)
		return
	}

	fnName := svcarn.ExtractFunctionNameFromARN(mapping.FunctionArn)
	if fnName == "" {
		p.log("failed to extract function name from ARN", "arn", mapping.FunctionArn)
		return
	}

	_, _, invokeErr := p.invokeLambda(ctx, fnName, payload)

	if invokeErr != nil {
		p.log("lambda invocation failed", "function", fnName, "queue", queueName, "error", invokeErr)
		_ = p.esmStore.SetState(mapping.UUID, "Enabled", fmt.Sprintf("Last processing result: %s", invokeErr.Error()))
		return
	}

	for _, handle := range receiptHandles {
		if err := p.sqsStore.DeleteMessage(queueURL, handle); err != nil {
			p.log("failed to delete message", "queue", queueName, "error", err)
		}
	}

	_ = p.esmStore.SetState(mapping.UUID, "Enabled", "Last processing result: No errors.")
}

// sqsMessageToRecord converts an SQS store Message into an ESM SQS record
// matching the Lambda event format.
func sqsMessageToRecord(msg *sqsstore.Message, eventSourceArn, region string) esmSQSRecord {
	record := esmSQSRecord{
		MessageID:               msg.ID,
		ReceiptHandle:           msg.ReceiptHandle,
		Body:                    msg.Body,
		MD5OfBody:               msg.MD5OfBody,
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
		record.MessageAttributes = map[string]interface{}{
			"SequenceNumber": map[string]string{
				"stringValue": msg.SequenceNumber,
				"dataType":    "String",
			},
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
func (p *esmPoller) invokeLambda(ctx context.Context, functionName string, payload []byte) (int64, []byte, error) {
	if p.lambdaSvc == nil {
		return 0, nil, fmt.Errorf("esm: lambda service not available")
	}

	fnStore := lambdastore.NewFunctionStore(p.lambdaSvc.getRegionalStorage(p.region), p.accountID, p.region)
	function, ver, _, err := p.lambdaSvc.resolveQualifier(fnStore, functionName, "")
	if err != nil {
		return 0, nil, fmt.Errorf("esm: failed to resolve function %s: %w", functionName, err)
	}

	result, err := p.lambdaSvc.invokeFunction(function, ver, fnStore, p.region, payload)
	if err != nil {
		return 0, nil, fmt.Errorf("esm: invocation failed for %s: %w", functionName, err)
	}
	if result == nil {
		return 0, nil, fmt.Errorf("esm: invocation returned nil result for %s", functionName)
	}

	return result.StatusCode, result.Payload, nil
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
