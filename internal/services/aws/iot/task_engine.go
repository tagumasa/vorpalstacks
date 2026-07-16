package iot

import (
	"context"
	"os"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

const (
	defaultTaskInterval = 30 * time.Second
	testTaskInterval    = 2 * time.Second
)

// taskEngine periodically scans for IN_PROGRESS IoT tasks (Detect
// Mitigation and On-Demand Audit) and transitions them to COMPLETED.
// Without this engine, tasks remain IN_PROGRESS indefinitely.
//
// In a full implementation, the engine would execute the actual task
// logic (violation scanning, audit checks) before transitioning. In our
// edge/on-prem deployment, the simplified behaviour is to transition
// tasks to COMPLETED on the next tick, since the underlying security
// engines (Device Defender, audit check runners) are not implemented.
type taskEngine struct {
	mu             sync.Mutex
	running        bool
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	interval       time.Duration
	storageManager *storage.RegionStorageManager
	accountID      string
}

func newTaskEngine(sm *storage.RegionStorageManager, accountID string) *taskEngine {
	interval := defaultTaskInterval
	if os.Getenv("TEST_MODE") == "true" {
		interval = testTaskInterval
	}
	return &taskEngine{
		interval:       interval,
		storageManager: sm,
		accountID:      accountID,
	}
}

// Start launches the background task processing loop.
func (e *taskEngine) Start(ctx context.Context) {
	e.mu.Lock()
	if e.running {
		e.mu.Unlock()
		return
	}
	e.running = true
	ctx, e.cancel = context.WithCancel(ctx)
	e.mu.Unlock()

	e.wg.Add(1)
	go e.taskLoop(ctx)
}

// Stop gracefully shuts down the task engine.
func (e *taskEngine) Stop() {
	e.mu.Lock()
	if !e.running {
		e.mu.Unlock()
		return
	}
	e.running = false
	if e.cancel != nil {
		e.cancel()
	}
	e.mu.Unlock()
	e.wg.Wait()
}

func (e *taskEngine) taskLoop(ctx context.Context) {
	defer e.wg.Done()
	defer func() {
		resilience.RecoverAndRestart("iot taskEngine", &e.wg, func() { e.taskLoop(ctx) })
	}()

	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.scanAndProcessTasks(ctx)
		}
	}
}

// scanAndProcessTasks iterates all configured regions, finds IN_PROGRESS
// tasks, and transitions them to COMPLETED.
func (e *taskEngine) scanAndProcessTasks(ctx context.Context) {
	if e.storageManager == nil {
		return
	}

	for _, region := range e.storageManager.ListRegions() {
		st, err := e.storageManager.GetStorage(region)
		if err != nil {
			continue
		}
		store := iotstore.GetOrCreateStore(st, e.accountID, region)

		e.processDetectMitigationTasks(store)
		e.processAuditTasks(store)
	}
}

// processDetectMitigationTasks scans for IN_PROGRESS detect mitigation
// tasks and transitions them to COMPLETED.
func (e *taskEngine) processDetectMitigationTasks(store iotstore.IotStoreInterface) {
	items, err := store.ListGeneric("detectMitigationTask/")
	if err != nil {
		logs.Warn("iot task engine: failed to list detect mitigation tasks", logs.Err(err))
		return
	}

	for _, rec := range items {
		status, _ := rec["status"].(string)
		if status != "IN_PROGRESS" {
			continue
		}

		taskId, _ := rec["taskId"].(string)
		if taskId == "" {
			continue
		}

		rec["status"] = "COMPLETED"
		rec["endTime"] = time.Now().UTC().Unix()
		if err := store.PutGeneric("detectMitigationTask/"+taskId, rec); err != nil {
			logs.Warn("iot task engine: failed to update detect mitigation task",
				logs.String("taskId", taskId), logs.Err(err))
		} else {
			logs.Debug("iot task engine: detect mitigation task completed",
				logs.String("taskId", taskId))
		}
	}
}

// processAuditTasks scans for IN_PROGRESS audit tasks and transitions
// them to COMPLETED, generating placeholder findings where applicable.
func (e *taskEngine) processAuditTasks(store iotstore.IotStoreInterface) {
	items, err := store.ListGeneric("auditTask/")
	if err != nil {
		logs.Warn("iot task engine: failed to list audit tasks", logs.Err(err))
		return
	}

	for _, rec := range items {
		status, _ := rec["status"].(string)
		if status != "IN_PROGRESS" {
			continue
		}

		taskId, _ := rec["taskId"].(string)
		if taskId == "" {
			continue
		}

		rec["status"] = "COMPLETED"
		rec["endTime"] = time.Now().UTC().Unix()
		if err := store.PutGeneric("auditTask/"+taskId, rec); err != nil {
			logs.Warn("iot task engine: failed to update audit task",
				logs.String("taskId", taskId), logs.Err(err))
		} else {
			logs.Debug("iot task engine: audit task completed",
				logs.String("taskId", taskId))
		}
	}
}
