// Package stepfunction provides Step Functions service operations for vorpalstacks.
package sfn

import (
	"context"
	"sync"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// ExecutorInterface defines the interface for executing state machines.
type ExecutorInterface interface {
	ExecuteStateMachine(ctx context.Context, execution *sfnstore.Execution) error
}

// StepFunctionService provides AWS Step Functions operations.
type StepFunctionService struct {
	executor       ExecutorInterface
	accountID      string
	storageManager *storage.RegionStorageManager
	bus            eventbus.Bus
	stores         sync.Map
	asyncWg        sync.WaitGroup
}

// NewStepFunctionService creates a new Step Functions service instance.
// Optional cross-service dependencies should be injected via setter methods
// before registering handlers.
func NewStepFunctionService(storageMgr *storage.RegionStorageManager, accountID string) *StepFunctionService {
	s := &StepFunctionService{
		accountID:      accountID,
		storageManager: storageMgr,
	}
	s.executor = NewExecutor(nil, nil)
	return s
}

// SetEventBus injects the event bus and subscribes to cross-service start
// execution events from EventBridge, Scheduler, and CloudWatch Alarms.
func (s *StepFunctionService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
	_, _ = eventbus.SubscribeTyped[*eventbus.StepFunctionsStartExecutionEvent](bus, s.handleStartExecutionEvent, eventbus.WithAsync())
}

func (s *StepFunctionService) handleStartExecutionEvent(ctx context.Context, evt *eventbus.StepFunctionsStartExecutionEvent) eventbus.HandlerResult {
	region := evt.Region
	if region == "" {
		region = defaults.DefaultRegion
	}

	store, err := s.getStoreForRegion(region)
	if err != nil {
		logs.Error("sfn: failed to get store for start execution event",
			logs.String("region", region),
			logs.String("stateMachineArn", evt.StateMachineArn),
			logs.Err(err))
		return eventbus.HandlerResult{}
	}

	// The bus start path runs the same validation, ARN resolution and
	// launch logic as the HTTP StartExecution handler; IoT rule actions
	// send stateMachineName (per Smithy), which the Core resolves.
	if err := s.startExecutionForBusCore(ctx, store, evt.StateMachineArn, evt.StateMachineName, evt.Input); err != nil {
		logs.Error("sfn: failed to start execution from bus event",
			logs.String("arn", evt.StateMachineArn),
			logs.String("name", evt.StateMachineName),
			logs.Err(err))
	}
	return eventbus.HandlerResult{}
}

func (s *StepFunctionService) store(reqCtx *request.RequestContext) (*sfnstore.StepFunctionStore, error) {
	region := reqCtx.GetRegion()
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (*sfnstore.StepFunctionStore, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return sfnstore.NewStepFunctionStore(storage, s.accountID, region), nil
	})
}

func (s *StepFunctionService) getStoreForRegion(region string) (*sfnstore.StepFunctionStore, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (*sfnstore.StepFunctionStore, error) {
		storage, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}
		return sfnstore.NewStepFunctionStore(storage, s.accountID, region), nil
	})
}

// Shutdown gracefully stops all Step Function stores and waits for
// any pending asynchronous operations to complete.
func (s *StepFunctionService) Shutdown() {
	s.stores.Range(func(key, value interface{}) bool {
		store := value.(*sfnstore.StepFunctionStore)
		store.CancelAllExecutions()
		return true
	})
	s.asyncWg.Wait()
}

// RegisterHandlers registers the Step Functions service handlers with the dispatcher.
func (s *StepFunctionService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("states", "CreateStateMachine", s.CreateStateMachine)
	d.RegisterHandlerForService("states", "DeleteStateMachine", s.DeleteStateMachine)
	d.RegisterHandlerForService("states", "DescribeStateMachine", s.DescribeStateMachine)
	d.RegisterHandlerForService("states", "ListStateMachines", s.ListStateMachines)
	d.RegisterHandlerForService("states", "UpdateStateMachine", s.UpdateStateMachine)

	d.RegisterHandlerForService("states", "StartExecution", s.StartExecution)
	d.RegisterHandlerForService("states", "StartSyncExecution", s.StartSyncExecution)
	d.RegisterHandlerForService("states", "StopExecution", s.StopExecution)
	d.RegisterHandlerForService("states", "DescribeExecution", s.DescribeExecution)
	d.RegisterHandlerForService("states", "DescribeStateMachineForExecution", s.DescribeStateMachineForExecution)
	d.RegisterHandlerForService("states", "ListExecutions", s.ListExecutions)
	d.RegisterHandlerForService("states", "GetExecutionHistory", s.GetExecutionHistory)

	d.RegisterHandlerForService("states", "DescribeMapRun", s.DescribeMapRun)
	d.RegisterHandlerForService("states", "ListMapRuns", s.ListMapRuns)

	d.RegisterHandlerForService("states", "CreateActivity", s.CreateActivity)
	d.RegisterHandlerForService("states", "DeleteActivity", s.DeleteActivity)
	d.RegisterHandlerForService("states", "DescribeActivity", s.DescribeActivity)
	d.RegisterHandlerForService("states", "ListActivities", s.ListActivities)

	d.RegisterHandlerForService("states", "GetActivityTask", s.GetActivityTask)
	d.RegisterHandlerForService("states", "SendTaskSuccess", s.SendTaskSuccess)
	d.RegisterHandlerForService("states", "SendTaskFailure", s.SendTaskFailure)
	d.RegisterHandlerForService("states", "SendTaskHeartbeat", s.SendTaskHeartbeat)

	d.RegisterHandlerForService("states", "ValidateStateMachineDefinition", s.ValidateStateMachineDefinition)

	d.RegisterHandlerForService("states", "TagResource", s.TagResource)
	d.RegisterHandlerForService("states", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("states", "ListTagsForResource", s.ListTagsForResource)

	d.RegisterHandlerForService("states", "RedriveExecution", s.RedriveExecution)
	d.RegisterHandlerForService("states", "TestState", s.TestState)

	d.RegisterHandlerForService("states", "PublishStateMachineVersion", s.PublishStateMachineVersion)
	d.RegisterHandlerForService("states", "DeleteStateMachineVersion", s.DeleteStateMachineVersion)
	d.RegisterHandlerForService("states", "ListStateMachineVersions", s.ListStateMachineVersions)

	d.RegisterHandlerForService("states", "CreateStateMachineAlias", s.CreateStateMachineAlias)
	d.RegisterHandlerForService("states", "DescribeStateMachineAlias", s.DescribeStateMachineAlias)
	d.RegisterHandlerForService("states", "DeleteStateMachineAlias", s.DeleteStateMachineAlias)
	d.RegisterHandlerForService("states", "UpdateStateMachineAlias", s.UpdateStateMachineAlias)
	d.RegisterHandlerForService("states", "ListStateMachineAliases", s.ListStateMachineAliases)

	d.RegisterHandlerForService("states", "UpdateMapRun", s.UpdateMapRun)
}
