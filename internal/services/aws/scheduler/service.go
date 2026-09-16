// Package scheduler provides AWS EventBridge Scheduler service operations for vorpalstacks.
package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// SchedulerService provides EventBridge Scheduler operations.
type SchedulerService struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	engine         *Engine
	stores         sync.Map
	roleProvider   iam.RolePolicyProvider
}

// NewSchedulerService creates a new Scheduler service instance.
// Cross-service dependencies and the engine must be set up before starting.
func NewSchedulerService(storageManager *storage.RegionStorageManager, accountID string) *SchedulerService {
	return &SchedulerService{
		storageManager: storageManager,
		accountID:      accountID,
	}
}

// BuildEngine constructs the scheduling engine from the currently injected
// dependencies and hands it the service's store cache as its single store
// source — one construction path and one closer (StopEngine) serve both the
// HTTP/admin planes and the engine. Must be called after all setter methods
// and before StartEngine.
func (s *SchedulerService) BuildEngine() {
	s.engine = NewEngine(s.storageManager, s.accountID)
	s.engine.SetStoreProvider(s.GetStoreForRegion)
}

// SetEventBus injects the event bus into the scheduler engine and registers
// the ScheduleFiredEvent handler. When the bus is set, schedule execution
// routes through the bus instead of direct store/invoker calls.
func (s *SchedulerService) SetEventBus(bus eventbus.ServiceBus) error {
	if s.engine != nil {
		s.engine.SetEventBus(bus)
		if _, err := eventbus.SubscribeTyped[*eventbus.ScheduleFiredEvent](bus, s.handleBusDelivery, eventbus.WithAsync()); err != nil {
			return fmt.Errorf("scheduler: subscribe ScheduleFiredEvent: %w", err)
		}
	}
	return nil
}

func (s *SchedulerService) handleBusDelivery(ctx context.Context, evt *eventbus.ScheduleFiredEvent) eventbus.HandlerResult {
	if s.engine == nil {
		return eventbus.HandlerResult{}
	}

	// Reconstruct the full target from the serialised payload so that all
	// sub-parameters (SqsParameters, KinesisParameters, etc.) are available
	// for correct delivery. Falls back to a minimal target if the payload
	// is absent (e.g. events from older engine instances).
	var target *schedulerstore.Target
	if evt.TargetPayload != "" {
		var t schedulerstore.Target
		if err := json.Unmarshal([]byte(evt.TargetPayload), &t); err == nil {
			target = &t
		}
	}
	if target == nil {
		target = &schedulerstore.Target{
			Arn:   evt.TargetArn,
			Input: evt.Input,
		}
	}

	schedule := &schedulerstore.Schedule{
		Name:                  evt.ScheduleName,
		GroupName:             evt.GroupName,
		Region:                evt.Region,
		ScheduleExpression:    evt.ScheduleExpression,
		Target:                target,
		ActionAfterCompletion: schedulerstore.ActionAfterCompletion(evt.ActionAfterCompletion),
	}

	// Unified delivery with retry and DLQ routing. deliverWithRetry handles
	// immediate retry, persisted background retries, and DLQ routing on
	// permanent failure.
	s.engine.deliverWithRetry(ctx, schedule, target)

	return eventbus.HandlerResult{}
}

// GetStoreForRegion returns the cached SchedulerStore for the given region,
// creating a new store instance if not already cached. It is the single
// construction path, shared by the admin plane and (as the engine's store
// provider) the scheduling engine.
func (s *SchedulerService) GetStoreForRegion(region string) (*schedulerstore.SchedulerStore, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (*schedulerstore.SchedulerStore, error) {
		if s.storageManager == nil {
			return nil, fmt.Errorf("scheduler storage manager not initialised")
		}
		st, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}
		return schedulerstore.NewSchedulerStore(st, s.accountID, region), nil
	})
}

// store resolves the regional store for a request-scoped region through
// the shared construction path (the engine's store provider binds to the
// same cache).
func (s *SchedulerService) store(ctx *request.RequestContext) (*schedulerstore.SchedulerStore, error) {
	return s.GetStoreForRegion(ctx.GetRegion())
}

// SetRoleProvider injects the IAM role policy provider so that the admin
// console handler can validate RoleArn trust policies.
func (s *SchedulerService) SetRoleProvider(rp iam.RolePolicyProvider) {
	s.roleProvider = rp
}

// RoleProvider returns the injected IAM role policy provider, or nil.
func (s *SchedulerService) RoleProvider() iam.RolePolicyProvider {
	return s.roleProvider
}

// AccountID returns the account ID for this service.
func (s *SchedulerService) AccountID() string {
	return s.accountID
}

// StartEngine starts the scheduler engine.
func (s *SchedulerService) StartEngine() error {
	if s.engine != nil {
		return s.engine.Start()
	}
	return nil
}

// StopEngine stops the scheduler engine and closes the shared per-region
// store cache, releasing each store's ClientTokenStore background cleanup
// goroutine. The engine owns no stores, so this single close covers
// everything the service created.
func (s *SchedulerService) StopEngine() error {
	var firstErr error
	if s.engine != nil {
		firstErr = s.engine.Stop()
	}
	// Close all cached per-region stores to release their ClientTokenStore
	// background cleanup goroutines.
	s.stores.Range(func(_, v interface{}) bool {
		if store, ok := v.(*schedulerstore.SchedulerStore); ok {
			store.Close()
		}
		return true
	})
	return firstErr
}

// RegisterHandlers registers the Scheduler service handlers with the dispatcher.
func (s *SchedulerService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("scheduler", "CreateScheduleGroup", s.CreateScheduleGroup)
	d.RegisterHandlerForService("scheduler", "DeleteScheduleGroup", s.DeleteScheduleGroup)
	d.RegisterHandlerForService("scheduler", "GetScheduleGroup", s.GetScheduleGroup)
	d.RegisterHandlerForService("scheduler", "ListScheduleGroups", s.ListScheduleGroups)

	d.RegisterHandlerForService("scheduler", "CreateSchedule", s.CreateSchedule)
	d.RegisterHandlerForService("scheduler", "DeleteSchedule", s.DeleteSchedule)
	d.RegisterHandlerForService("scheduler", "GetSchedule", s.GetSchedule)
	d.RegisterHandlerForService("scheduler", "UpdateSchedule", s.UpdateSchedule)
	d.RegisterHandlerForService("scheduler", "ListSchedules", s.ListSchedules)

	d.RegisterHandlerForService("scheduler", "TagResource", s.TagResource)
	d.RegisterHandlerForService("scheduler", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("scheduler", "ListTagsForResource", s.ListTagsForResource)
}
