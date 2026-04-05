// Package stepfunction provides Step Functions service operations for vorpalstacks.
package sfn

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/server/eventbus"
	"vorpalstacks/internal/services/aws/common"
	"vorpalstacks/internal/services/aws/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
	arncommon "vorpalstacks/internal/utils/aws/arn"
)

// ExecutorInterface defines the interface for executing state machines.
type ExecutorInterface interface {
	ExecuteStateMachine(ctx context.Context, execution *sfnstore.Execution) error
}

// StepFunctionService provides AWS Step Functions operations.
type StepFunctionService struct {
	executor       ExecutorInterface
	lambdaInvoker  common.LambdaInvoker
	sqsStore       sqsstore.SQSStoreInterface
	snsStore       snsstore.SNSStoreInterface
	eventsStore    *eventsstore.EventsStore
	accountID      string
	storageManager *storage.RegionStorageManager
	bus            *eventbus.EventBus
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
	s.executor = NewExecutor(nil, s.lambdaInvoker)
	return s
}

// SetLambdaInvoker injects a Lambda invoker for AWS SDK and Lambda integration states.
func (s *StepFunctionService) SetLambdaInvoker(invoker common.LambdaInvoker) {
	s.lambdaInvoker = invoker
}

// SetSQSStore injects an SQS store for cross-service SQS integration states.
func (s *StepFunctionService) SetSQSStore(store sqsstore.SQSStoreInterface) {
	s.sqsStore = store
}

// SetSNSStore injects an SNS store for cross-service SNS integration states.
func (s *StepFunctionService) SetSNSStore(store snsstore.SNSStoreInterface) {
	s.snsStore = store
}

// SetEventsStore injects an EventBridge store for cross-service EventBridge integration states.
func (s *StepFunctionService) SetEventsStore(store *eventsstore.EventsStore) {
	s.eventsStore = store
}

// SetEventBus injects the event bus and subscribes to cross-service start
// execution events from EventBridge, Scheduler, and CloudWatch Alarms.
func (s *StepFunctionService) SetEventBus(bus *eventbus.EventBus) {
	s.bus = bus
	_, _ = eventbus.SubscribeTyped[*eventbus.StepFunctionsStartExecutionEvent](bus, s.handleStartExecutionEvent, eventbus.WithAsync())
}

func (s *StepFunctionService) handleStartExecutionEvent(ctx context.Context, evt *eventbus.StepFunctionsStartExecutionEvent) eventbus.HandlerResult {
	region := evt.Region
	if region == "" {
		region = request.DefaultRegion
	}

	store, err := s.getStoreForRegion(region)
	if err != nil {
		logs.Error("sfn: failed to get store for start execution event",
			logs.String("region", region),
			logs.String("stateMachineArn", evt.StateMachineArn),
			logs.Err(err))
		return eventbus.HandlerResult{}
	}

	sm, err := store.GetStateMachine(ctx, evt.StateMachineArn)
	if err != nil {
		logs.Error("sfn: failed to get state machine for start execution event",
			logs.String("arn", evt.StateMachineArn),
			logs.Err(err))
		return eventbus.HandlerResult{}
	}

	name := fmt.Sprintf("bus-%s", uuid.New().String())
	executionArn := fmt.Sprintf("arn:aws:states:%s:%s:execution:%s:%s",
		region, s.accountID,
		arncommon.ExtractStateMachineNameFromARN(sm.StateMachineArn), name)

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, evt.Input, "")
	exec.ExecutionArn = executionArn

	if err := store.CreateExecution(ctx, exec); err != nil {
		logs.Error("sfn: failed to create execution from bus event",
			logs.String("executionArn", executionArn),
			logs.Err(err))
		return eventbus.HandlerResult{}
	}

	executor := NewExecutorWithStores(store, s.lambdaInvoker, s.sqsStore, s.snsStore, s.eventsStore, s.accountID, region)
	executor.SetEventBus(s.bus)
	execCtx, cancel := context.WithCancel(context.Background())
	store.RegisterExecution(executionArn, cancel)
	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer store.UnregisterExecution(executionArn)
		defer func() {
			if r := recover(); r != nil {
				logs.Error("sfn: panic in bus-triggered execution", logs.String("arn", executionArn), logs.Any("panic", r))
			}
		}()
		_ = executor.ExecuteStateMachine(execCtx, exec)
	}()

	logs.Debug("sfn: started execution from bus event",
		logs.String("executionArn", executionArn),
		logs.String("stateMachineArn", evt.StateMachineArn))

	return eventbus.HandlerResult{}
}

func (s *StepFunctionService) store(reqCtx *request.RequestContext) (*sfnstore.StepFunctionStore, error) {
	if sfnStore := reqCtx.GetSFNStore(); sfnStore != nil {
		if concrete, ok := sfnStore.(*sfnstore.StepFunctionStore); ok {
			return concrete, nil
		}
	}

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
func (s *StepFunctionService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("states", "CreateStateMachine", s.CreateStateMachine)
	d.RegisterHandlerForService("states", "DeleteStateMachine", s.DeleteStateMachine)
	d.RegisterHandlerForService("states", "DescribeStateMachine", s.DescribeStateMachine)
	d.RegisterHandlerForService("states", "GetStateMachine", s.GetStateMachine)
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
	d.RegisterHandlerForService("states", "GetActivity", s.GetActivity)
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
