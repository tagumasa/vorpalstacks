// Package logs provides AWS CloudWatch Logs service operations for vorpalstacks.
package cloudwatchlogs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"crypto/rand"
	"encoding/hex"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/worker"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// LogsService provides CloudWatch Logs service operations.
type LogsService struct {
	storageManager  *storage.RegionStorageManager
	accountID       string
	dataPath        string
	logsStores      sync.Map // region → *logsstore.Store
	cwMetricInvoker invokers.CloudWatchMetricInvoker
	bus             eventbus.ServiceBus
	kms             invokers.KMSInvoker
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	queries         sync.Map // queryId → *queryState (per-service, not global)
	filterJobs      chan filterFanOutJob
	// metricDefaultMinutes tracks, per (log group, metric filter), the
	// current minute's match/default-emission state so a matchless
	// evaluation publishes the DefaultValue at most once per minute.
	metricDefaultMu      sync.Mutex
	metricDefaultMinutes map[string]metricMinuteState
	// liveTailMu guards liveTailSessions, the count of open Live Tail
	// sessions against the documented concurrent-session quota.
	liveTailMu       sync.Mutex
	liveTailSessions int
	// taskAdmissionMu serialises the quota admission windows of the
	// export and import task creations: the active-task census and the
	// RUNNING record's write run as one critical section, so two
	// concurrent creations cannot both observe the quota open and both
	// admit.
	taskAdmissionMu sync.Mutex
	// queryAdmissionMu serialises StartQuery's concurrency admission: the
	// running-query census and the RUNNING state's registration run as
	// one critical section, so two concurrent starts cannot both observe
	// the quota open and both admit.
	queryAdmissionMu sync.Mutex
	// depsMu guards the injected cross-service dependencies (bus, kms,
	// metric invoker): their setters can run while the service's
	// background workers are already reading them.
	depsMu sync.RWMutex
	// delivering guards this instance's delivery pass against overlap: a
	// pass slower than the cadence skips the next tick instead of racing
	// its own cursor write.
	delivering uint32
	// lookupTableAdmissionMu serialises lookup-table creation: the
	// exists-check, the account quota census and the record's write run
	// as one critical section, so two concurrent creations cannot both
	// pass an open quota or both admit the same name.
	lookupTableAdmissionMu sync.Mutex
}

// NewLogsService creates a new CloudWatch Logs service.
// Optional cross-service dependencies (Lambda invoker, Kinesis store) should be
// injected via setter methods before registering handlers.
func NewLogsService(storageMgr *storage.RegionStorageManager, accountID, dataPath string) *LogsService {
	ctx, cancel := context.WithCancel(context.Background())
	s := &LogsService{
		storageManager: storageMgr,
		accountID:      accountID,
		dataPath:       dataPath,
		ctx:            ctx,
		cancel:         cancel,
	}
	s.startRetentionPurger()
	s.startSubscriptionDeliveryRetryLoop()
	s.startScheduledQueryWorker()
	s.startFilterFanOutWorkers()
	s.startDeliveryWorker()
	s.startMetricDefaultMinuteCloser()
	s.reconcileInterruptedTasks()
	s.loadPersistedQueries()
	return s
}

// newRandomHexID mints a random hexadecimal identifier of n bytes — the
// shared minter behind the delivery and live-tail identifiers (both
// 32-character hex strings).
func newRandomHexID(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// randomHexWithFallback mints n random bytes as hex. A rand failure is
// beyond unlikely on this platform, but a minter's fallback keeps a
// working identifier (the clock, a digest) instead of failing the
// operation — the degradation shape every local minter shares; the
// fallback arrives as a closure so the healthy path never pays for it.
func randomHexWithFallback(n int, fallback func() string) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return fallback()
	}
	return hex.EncodeToString(b)
}

// spawnTask launches one asynchronous task worker (export or import)
// under the service's WaitGroup and context: Stop cancels the context and
// waits for in-flight tasks, so a shutdown can no longer strand a writer
// against closed stores, and a worker's S3 I/O aborts promptly.
func (s *LogsService) spawnTask(fn func(ctx context.Context)) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		fn(s.ctx)
	}()
}

// reconcileInterruptedTasks fails export and import records that a
// previous process left in an active state. A restart kills every worker
// without a chance to finalise, so any RUNNING/PENDING/IN_PROGRESS
// record at construction time has no worker behind it; FAILED with the
// interruption reason is the honest terminal state. The region set comes
// from the storage manager, not the store cache (the purgeAllRegions
// rule: active-state enforcement must not depend on foreground traffic
// having touched the region).
func (s *LogsService) reconcileInterruptedTasks() {
	for _, region := range worker.ActiveRegions(s.storageManager) {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			continue
		}
		exports, err := store.ListExportTasks("")
		if err != nil {
			continue
		}
		for _, t := range exports {
			if t.Status != logsstore.ExportStatusRunning && t.Status != logsstore.ExportStatusPending && t.Status != logsstore.ExportStatusPendingCancel {
				continue
			}
			// A PENDING_CANCEL record's cancellation was already
			// requested: the restart killed the worker that would have
			// honoured it, so the terminal state the request sought is
			// written here rather than the interruption failure.
			if t.Status == logsstore.ExportStatusPendingCancel {
				t.Status = logsstore.ExportStatusCancelled
				t.StatusMessage = "Cancelled by user"
			} else {
				t.Status = logsstore.ExportStatusFailed
				t.StatusMessage = "interrupted by service restart"
			}
			if err := store.PutExportTask(t); err != nil {
				logs.Error("Failed to reconcile interrupted export task",
					logs.String("taskId", t.TaskId), logs.Err(err))
			}
		}
		imports, err := store.ListImportTasks("", "", "")
		if err != nil {
			continue
		}
		for _, t := range imports {
			if t.ImportStatus != logsstore.ImportStatusInProgress {
				continue
			}
			t.ImportStatus = logsstore.ImportStatusFailed
			t.ErrorMessage = "interrupted by service restart"
			if err := store.PutImportTask(t); err != nil {
				logs.Error("Failed to reconcile interrupted import task",
					logs.String("importId", t.ImportId), logs.Err(err))
			}
		}
		// A RUNNING scheduled-query execution is equally orphaned, and the
		// ticker's in-flight guard skips a query whose latest execution is
		// stuck RUNNING — failing the orphan is what unblocks the schedule.
		scheduled, err := store.ListScheduledQueries("")
		if err != nil {
			continue
		}
		for _, sq := range scheduled {
			execs, err := store.ListScheduledQueryExecutions(sq.Id, 0, 0)
			if err != nil {
				continue
			}
			for _, exec := range execs {
				if exec.Status != logsstore.ScheduledExecutionStatusRunning {
					continue
				}
				exec.Status = logsstore.ScheduledExecutionStatusFailed
				exec.ErrorMessage = "interrupted by service restart"
				if err := store.PutScheduledQueryExecution(exec); err != nil {
					logs.Error("Failed to reconcile interrupted scheduled query execution",
						logs.String("queryId", exec.QueryId), logs.Err(err))
				}
			}
		}
	}
}

// SetEventBus injects the event bus and registers handlers for CloudWatch Logs
// delivery, Lambda log writes, API Gateway access logs, and direct log event
// ingestion from EventBridge/Scheduler targets.
func (s *LogsService) SetEventBus(bus eventbus.ServiceBus) error {
	s.depsMu.Lock()
	s.bus = bus
	if bus != nil {
		s.kms = bus.KMSInvoker()
	}
	s.depsMu.Unlock()
	if _, err := eventbus.SubscribeTyped[*eventbus.CloudWatchLogDeliveryEvent](bus, s.handleBusDelivery, eventbus.WithAsync()); err != nil {
		return fmt.Errorf("logs: subscribe CloudWatchLogDeliveryEvent: %w", err)
	}
	if _, err := eventbus.SubscribeTyped[*eventbus.LambdaLogWriteEvent](bus, s.handleLambdaLogWrite, eventbus.WithAsync()); err != nil {
		return fmt.Errorf("logs: subscribe LambdaLogWriteEvent: %w", err)
	}
	if _, err := eventbus.SubscribeTyped[*eventbus.APIGatewayAccessLogEvent](bus, s.handleAPIGatewayAccessLog, eventbus.WithAsync()); err != nil {
		return fmt.Errorf("logs: subscribe APIGatewayAccessLogEvent: %w", err)
	}
	if _, err := eventbus.SubscribeTyped[*eventbus.CloudWatchLogsPutEvent](bus, s.handleDirectPutLogEvents, eventbus.WithAsync()); err != nil {
		return fmt.Errorf("logs: subscribe CloudWatchLogsPutEvent: %w", err)
	}
	return nil
}

func (s *LogsService) store(reqCtx *request.RequestContext) (*logsstore.Store, error) {
	return storecommon.GetOrCreateStoreE(&s.logsStores, reqCtx.GetRegion(), func() (*logsstore.Store, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		store, err := logsstore.NewStore(storage, storage.Bucket("logs-"+reqCtx.GetRegion()), s.accountID, reqCtx.GetRegion(), s.dataPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create CloudWatch Logs store: %w", err)
		}
		return store, nil
	})
}

// getLogsStoreByRegion resolves the CloudWatch Logs store for the given region.
// Used by bus handlers that operate outside of an HTTP request context.
func (s *LogsService) getLogsStoreByRegion(region string) (*logsstore.Store, error) {
	return storecommon.GetOrCreateStoreE(&s.logsStores, region, func() (*logsstore.Store, error) {
		regionStorage, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, fmt.Errorf("failed to get storage for region %q: %w", region, err)
		}
		store, err := logsstore.NewStore(regionStorage, regionStorage.Bucket("logs-"+region), s.accountID, region, s.dataPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create CloudWatch Logs store: %w", err)
		}
		return store, nil
	})
}

// GetStoreForRegion resolves the CloudWatch Logs store for the given region,
// creating it on first use. Cross-service consumers (the eventbus logs
// invoker) resolve stores through this method so that every writer and the
// API read plane share one store instance per region.
func (s *LogsService) GetStoreForRegion(region string) (*logsstore.Store, error) {
	return s.getLogsStoreByRegion(region)
}

// SetCloudWatchMetricInvoker injects the CloudWatch metric invoker for emitting
// metric data when metric filters match log events.
func (s *LogsService) SetCloudWatchMetricInvoker(invoker invokers.CloudWatchMetricInvoker) {
	s.depsMu.Lock()
	s.cwMetricInvoker = invoker
	s.depsMu.Unlock()
}

// The injected dependencies are published while the service's workers
// already run (the delivery worker and the metric fan-out read them on
// every pass), so every read and write goes through the deps mutex.
func (s *LogsService) eventBus() eventbus.ServiceBus {
	s.depsMu.RLock()
	defer s.depsMu.RUnlock()
	return s.bus
}

func (s *LogsService) kmsInvoker() invokers.KMSInvoker {
	s.depsMu.RLock()
	defer s.depsMu.RUnlock()
	return s.kms
}

func (s *LogsService) metricInvoker() invokers.CloudWatchMetricInvoker {
	s.depsMu.RLock()
	defer s.depsMu.RUnlock()
	return s.cwMetricInvoker
}

// RegisterHandlers registers the CloudWatch Logs service handlers with the dispatcher.
func (s *LogsService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("logs", "CreateLogGroup", s.CreateLogGroup)
	d.RegisterHandlerForService("logs", "DeleteLogGroup", s.DeleteLogGroup)
	d.RegisterHandlerForService("logs", "DescribeLogGroups", s.DescribeLogGroups)
	d.RegisterHandlerForService("logs", "ListLogGroups", s.ListLogGroups)
	d.RegisterHandlerForService("logs", "PutRetentionPolicy", s.PutRetentionPolicy)
	d.RegisterHandlerForService("logs", "DeleteRetentionPolicy", s.DeleteRetentionPolicy)
	d.RegisterHandlerForService("logs", "TagResource", s.TagResource)
	d.RegisterHandlerForService("logs", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("logs", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("logs", "TagLogGroup", s.TagLogGroup)
	d.RegisterHandlerForService("logs", "ListTagsLogGroup", s.ListTagsLogGroup)
	d.RegisterHandlerForService("logs", "UntagLogGroup", s.UntagLogGroup)

	d.RegisterHandlerForService("logs", "CreateLogStream", s.CreateLogStream)
	d.RegisterHandlerForService("logs", "DeleteLogStream", s.DeleteLogStream)
	d.RegisterHandlerForService("logs", "DescribeLogStreams", s.DescribeLogStreams)

	d.RegisterHandlerForService("logs", "PutLogEvents", s.PutLogEvents)
	d.RegisterHandlerForService("logs", "GetLogEvents", s.GetLogEvents)
	d.RegisterHandlerForService("logs", "FilterLogEvents", s.FilterLogEvents)

	d.RegisterHandlerForService("logs", "PutMetricFilter", s.PutMetricFilter)
	d.RegisterHandlerForService("logs", "DeleteMetricFilter", s.DeleteMetricFilter)
	d.RegisterHandlerForService("logs", "DescribeMetricFilters", s.DescribeMetricFilters)
	d.RegisterHandlerForService("logs", "TestMetricFilter", s.TestMetricFilter)

	d.RegisterHandlerForService("logs", "PutSubscriptionFilter", s.PutSubscriptionFilter)
	d.RegisterHandlerForService("logs", "DeleteSubscriptionFilter", s.DeleteSubscriptionFilter)
	d.RegisterHandlerForService("logs", "DescribeSubscriptionFilters", s.DescribeSubscriptionFilters)

	d.RegisterHandlerForService("logs", "PutDestination", s.PutDestination)
	d.RegisterHandlerForService("logs", "PutDestinationPolicy", s.PutDestinationPolicy)
	d.RegisterHandlerForService("logs", "DeleteDestination", s.DeleteDestination)
	d.RegisterHandlerForService("logs", "DescribeDestinations", s.DescribeDestinations)

	d.RegisterHandlerForService("logs", "AssociateKmsKey", s.AssociateKmsKey)
	d.RegisterHandlerForService("logs", "DisassociateKmsKey", s.DisassociateKmsKey)
	d.RegisterHandlerForService("logs", "PutLogGroupDeletionProtection", s.PutLogGroupDeletionProtection)

	d.RegisterHandlerForService("logs", "PutResourcePolicy", s.PutResourcePolicy)
	d.RegisterHandlerForService("logs", "DeleteResourcePolicy", s.DeleteResourcePolicy)
	d.RegisterHandlerForService("logs", "DescribeResourcePolicies", s.DescribeResourcePolicies)

	d.RegisterHandlerForService("logs", "PutAccountPolicy", s.PutAccountPolicy)
	d.RegisterHandlerForService("logs", "DeleteAccountPolicy", s.DeleteAccountPolicy)
	d.RegisterHandlerForService("logs", "DescribeAccountPolicies", s.DescribeAccountPolicies)

	d.RegisterHandlerForService("logs", "PutDataProtectionPolicy", s.PutDataProtectionPolicy)
	d.RegisterHandlerForService("logs", "GetDataProtectionPolicy", s.GetDataProtectionPolicy)
	d.RegisterHandlerForService("logs", "DeleteDataProtectionPolicy", s.DeleteDataProtectionPolicy)

	d.RegisterHandlerForService("logs", "PutQueryDefinition", s.PutQueryDefinition)
	d.RegisterHandlerForService("logs", "DeleteQueryDefinition", s.DeleteQueryDefinition)
	d.RegisterHandlerForService("logs", "DescribeQueryDefinitions", s.DescribeQueryDefinitions)

	d.RegisterHandlerForService("logs", "GetLogGroupFields", s.GetLogGroupFields)
	d.RegisterHandlerForService("logs", "GetLogRecord", s.GetLogRecord)
	d.RegisterHandlerForService("logs", "GetLogObject", s.GetLogObject)
	d.RegisterHandlerForService("logs", "GetLogFields", s.GetLogFields)

	d.RegisterHandlerForService("logs", "CreateExportTask", s.CreateExportTask)
	d.RegisterHandlerForService("logs", "DescribeExportTasks", s.DescribeExportTasks)
	d.RegisterHandlerForService("logs", "CancelExportTask", s.CancelExportTask)

	d.RegisterHandlerForService("logs", "CreateImportTask", s.CreateImportTask)
	d.RegisterHandlerForService("logs", "DescribeImportTasks", s.DescribeImportTasks)
	d.RegisterHandlerForService("logs", "CancelImportTask", s.CancelImportTask)
	d.RegisterHandlerForService("logs", "DescribeImportTaskBatches", s.DescribeImportTaskBatches)

	d.RegisterHandlerForService("logs", "StartQuery", s.StartQuery)
	d.RegisterHandlerForService("logs", "StopQuery", s.StopQuery)
	d.RegisterHandlerForService("logs", "DescribeQueries", s.DescribeQueries)
	d.RegisterHandlerForService("logs", "GetQueryResults", s.GetQueryResults)

	d.RegisterHandlerForService("logs", "CreateScheduledQuery", s.CreateScheduledQuery)
	d.RegisterHandlerForService("logs", "DeleteScheduledQuery", s.DeleteScheduledQuery)
	d.RegisterHandlerForService("logs", "UpdateScheduledQuery", s.UpdateScheduledQuery)
	d.RegisterHandlerForService("logs", "GetScheduledQuery", s.GetScheduledQuery)
	d.RegisterHandlerForService("logs", "GetScheduledQueryHistory", s.GetScheduledQueryHistory)
	d.RegisterHandlerForService("logs", "ListScheduledQueries", s.ListScheduledQueries)

	d.RegisterHandlerForService("logs", "CreateLookupTable", s.CreateLookupTable)
	d.RegisterHandlerForService("logs", "DeleteLookupTable", s.DeleteLookupTable)
	d.RegisterHandlerForService("logs", "GetLookupTable", s.GetLookupTable)
	d.RegisterHandlerForService("logs", "UpdateLookupTable", s.UpdateLookupTable)
	d.RegisterHandlerForService("logs", "DescribeLookupTables", s.DescribeLookupTables)

	d.RegisterHandlerForService("logs", "PutDeliverySource", s.PutDeliverySource)
	d.RegisterHandlerForService("logs", "DescribeDeliverySources", s.DescribeDeliverySources)
	d.RegisterHandlerForService("logs", "GetDeliverySource", s.GetDeliverySource)
	d.RegisterHandlerForService("logs", "DeleteDeliverySource", s.DeleteDeliverySource)
	d.RegisterHandlerForService("logs", "PutDeliveryDestination", s.PutDeliveryDestination)
	d.RegisterHandlerForService("logs", "DescribeDeliveryDestinations", s.DescribeDeliveryDestinations)
	d.RegisterHandlerForService("logs", "GetDeliveryDestination", s.GetDeliveryDestination)
	d.RegisterHandlerForService("logs", "DeleteDeliveryDestination", s.DeleteDeliveryDestination)
	d.RegisterHandlerForService("logs", "PutDeliveryDestinationPolicy", s.PutDeliveryDestinationPolicy)
	d.RegisterHandlerForService("logs", "GetDeliveryDestinationPolicy", s.GetDeliveryDestinationPolicy)
	d.RegisterHandlerForService("logs", "DeleteDeliveryDestinationPolicy", s.DeleteDeliveryDestinationPolicy)
	d.RegisterHandlerForService("logs", "CreateDelivery", s.CreateDelivery)
	d.RegisterHandlerForService("logs", "DescribeDeliveries", s.DescribeDeliveries)
	d.RegisterHandlerForService("logs", "GetDelivery", s.GetDelivery)
	d.RegisterHandlerForService("logs", "DeleteDelivery", s.DeleteDelivery)
	d.RegisterHandlerForService("logs", "UpdateDeliveryConfiguration", s.UpdateDeliveryConfiguration)
	d.RegisterHandlerForService("logs", "DescribeConfigurationTemplates", s.DescribeConfigurationTemplates)

	d.RegisterHandlerForService("logs", "PutTransformer", s.PutTransformer)
	d.RegisterHandlerForService("logs", "GetTransformer", s.GetTransformer)
	d.RegisterHandlerForService("logs", "DeleteTransformer", s.DeleteTransformer)
	d.RegisterHandlerForService("logs", "TestTransformer", s.TestTransformer)

	d.RegisterHandlerForService("logs", "PutIndexPolicy", s.PutIndexPolicy)
	d.RegisterHandlerForService("logs", "DescribeIndexPolicies", s.DescribeIndexPolicies)
	d.RegisterHandlerForService("logs", "DeleteIndexPolicy", s.DeleteIndexPolicy)
	d.RegisterHandlerForService("logs", "DescribeFieldIndexes", s.DescribeFieldIndexes)

	d.RegisterHandlerForService("logs", "PutStorageTierPolicy", s.PutStorageTierPolicy)
	d.RegisterHandlerForService("logs", "GetStorageTierPolicy", s.GetStorageTierPolicy)
	d.RegisterHandlerForService("logs", "ListLogGroupsForQuery", s.ListLogGroupsForQuery)
	d.RegisterHandlerForService("logs", "ListAggregateLogGroupSummaries", s.ListAggregateLogGroupSummaries)
	d.RegisterHandlerForService("logs", "StartLiveTail", s.StartLiveTail)
	d.RegisterHandlerForService("logs", "PutBearerTokenAuthentication", s.PutBearerTokenAuthentication)
}

// startRespawningWorker runs fn under the service's WaitGroup with the
// panic-respawn policy — a panic is logged and the worker restarts in
// the same goroutine, and a normal return or context cancellation ends
// it. The respawn stays inside the one WaitGroup member the goroutine
// already holds: a respawn that re-Add()ed from a fresh goroutine had no
// ordering edge against the panicking member's Done, so a shutdown Wait
// could observe zero between them. This is the one lifecycle rule behind
// the delivery worker, the scheduled-query ticker and the retention
// purger. The supervision itself is the shared worker shell.
func (s *LogsService) startRespawningWorker(name string, fn func()) {
	worker.RunSupervised(s.ctx, &s.wg, name, fn)
}

func (s *LogsService) startRetentionPurger() {
	s.startRespawningWorker("retention purger",
		worker.TickerLoop(s.ctx.Done(), time.Hour, func() {
			s.purgeAllRegions()
			s.evictExpiredQueries()
		}))
}

// evictExpiredQueries drops query state older than the documented
// query-result retention window. Without eviction the query map grows
// with every query for the process lifetime; the sweep bounds it the way
// AWS bounds result availability. The persisted records follow the same
// cutoff across the configured regions, so the retention window bounds
// durability too — a restart cannot resurrect an expired result.
func (s *LogsService) evictExpiredQueries() {
	cutoff := time.Now().Add(-logsstore.QueryResultRetentionPeriod)
	s.queries.Range(func(key, value interface{}) bool {
		qs := value.(*queryState)
		qs.mu.RLock()
		expired := qs.createdAt.Before(cutoff)
		qs.mu.RUnlock()
		if expired {
			s.queries.Delete(key)
		}
		return true
	})
	for _, region := range worker.ActiveRegions(s.storageManager) {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			continue
		}
		records, err := store.ListQueryRecords()
		if err != nil {
			continue
		}
		for _, rec := range records {
			if time.Unix(0, rec.CreatedAtUnixNano).Before(cutoff) {
				if err := store.DeleteQueryRecord(rec.QueryId); err != nil {
					logs.Error("Failed to evict expired query record",
						logs.String("queryId", rec.QueryId), logs.Err(err))
				}
			}
		}
	}
}

// purgeAllRegions purges expired chunks in every configured region. The
// region set comes from the storage manager, not the store instance
// cache: the cache fills only when API traffic constructs a region's
// store, and retention enforcement must not depend on unrelated
// foreground traffic having touched the region since process start.
func (s *LogsService) purgeAllRegions() {
	for _, region := range worker.ActiveRegions(s.storageManager) {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			logs.Error("Failed to resolve logs store for retention purge",
				logs.String("region", region), logs.Err(err))
			continue
		}
		if err := store.PurgeAllExpiredChunks(); err != nil {
			logs.Error("Failed to purge expired chunks",
				logs.String("region", region), logs.Err(err))
		}
	}
}

// AccountID returns the AWS account ID for this service.
func (s *LogsService) AccountID() string {
	return s.accountID
}

// Stop stops the CloudWatch Logs service by canceling the context and waiting for goroutines to complete.
func (s *LogsService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}
