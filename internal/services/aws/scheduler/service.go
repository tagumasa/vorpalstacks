// Package scheduler provides AWS EventBridge Scheduler service operations for vorpalstacks.
package scheduler

import (
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common"
	"vorpalstacks/internal/services/aws/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// SchedulerService provides EventBridge Scheduler operations.
type SchedulerService struct {
	storageManager *storage.RegionStorageManager
	sqsStore       sqsstore.SQSStoreInterface
	snsStore       snsstore.SNSStoreInterface
	lambdaInvoker  common.LambdaInvoker
	accountID      string
	engine         *Engine
	stores         sync.Map
}

// NewSchedulerService creates a new Scheduler service instance.
// Cross-service dependencies and the engine must be set up before starting.
func NewSchedulerService(storageManager *storage.RegionStorageManager, accountID string) *SchedulerService {
	return &SchedulerService{
		storageManager: storageManager,
		accountID:      accountID,
	}
}

// SetSQSStore injects an SQS store for target resolution.
func (s *SchedulerService) SetSQSStore(store sqsstore.SQSStoreInterface) {
	s.sqsStore = store
}

// SetSNSStore injects an SNS store for target resolution.
func (s *SchedulerService) SetSNSStore(store snsstore.SNSStoreInterface) {
	s.snsStore = store
}

// SetLambdaInvoker injects a Lambda invoker for target resolution.
func (s *SchedulerService) SetLambdaInvoker(invoker common.LambdaInvoker) {
	s.lambdaInvoker = invoker
}

// BuildEngine constructs the scheduling engine from the currently injected dependencies.
// Must be called after all setter methods and before StartEngine.
func (s *SchedulerService) BuildEngine() {
	s.engine = NewEngine(s.storageManager, s.sqsStore, s.snsStore, s.lambdaInvoker, s.accountID)
}

func (s *SchedulerService) store(ctx *request.RequestContext) (*schedulerstore.SchedulerStore, error) {
	if store := ctx.GetSchedulerStore(); store != nil {
		return store.Raw(), nil
	}
	return storecommon.GetOrCreateStoreE(&s.stores, ctx.GetRegion(), func() (*schedulerstore.SchedulerStore, error) {
		st, err := s.storageManager.GetStorage(ctx.GetRegion())
		if err != nil {
			return nil, err
		}
		return schedulerstore.NewSchedulerStore(st, s.accountID, ctx.GetRegion()), nil
	})
}

// StartEngine starts the scheduler engine.
func (s *SchedulerService) StartEngine() error {
	if s.engine != nil {
		return s.engine.Start()
	}
	return nil
}

// StopEngine stops the scheduler engine.
func (s *SchedulerService) StopEngine() error {
	if s.engine != nil {
		return s.engine.Stop()
	}
	return nil
}

// RegisterHandlers registers the Scheduler service handlers with the dispatcher.
func (s *SchedulerService) RegisterHandlers(d *dispatcher.Dispatcher) {
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
