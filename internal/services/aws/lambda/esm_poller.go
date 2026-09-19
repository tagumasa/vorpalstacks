package lambda

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
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

	// esmCheckpointBucket is the Pebble bucket name used to persist
	// Kinesis ESM checkpoint data across server restarts.
	esmCheckpointBucket = "lambda-esm-checkpoints"
)

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
	bus            eventbus.ServiceBus
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
	defer func() {
		if r := recover(); r != nil {
			resilience.RestartAfterPanic("ESM pollLoop", r, &p.wg, func() { p.pollLoop(ctx) })
		}
	}()
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
	} else {
		stores, err := p.lambdaSvc.getOrCreateLambdaStoreE(region)
		if err != nil {
			// Report incompleteness so the caller skips the checkpoint
			// purge rather than wiping checkpoints for mappings it could
			// not see.
			p.log("failed to build lambda store", "region", region, "error", err)
			return nil, false
		}
		esmStore = stores.EventSources
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
			defer func() {
				if r := recover(); r != nil {
					resilience.LogPanic("lambda esm worker", r)
				}
			}()
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
