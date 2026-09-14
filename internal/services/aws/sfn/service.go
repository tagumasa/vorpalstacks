// Package stepfunction provides Step Functions service operations for vorpalstacks.
package sfn

import (
	"context"
	"fmt"
	"sync"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// StepFunctionService provides AWS Step Functions operations.
type StepFunctionService struct {
	accountID            string
	storageManager       *storage.RegionStorageManager
	bus                  eventbus.ServiceBus
	stores               sync.Map
	asyncWg              sync.WaitGroup
	roleProvider         iam.RolePolicyProvider
	taskCredentialsAuthz TaskCredentialsAuthoriser
}

// NewStepFunctionService creates a new Step Functions service instance.
// Optional cross-service dependencies should be injected via setter methods
// before registering handlers.
func NewStepFunctionService(storageMgr *storage.RegionStorageManager, accountID string) *StepFunctionService {
	return &StepFunctionService{
		accountID:      accountID,
		storageManager: storageMgr,
	}
}

// SetRoleProvider injects the IAM role policy provider so the state
// machine Cores can validate role ARNs on both planes without a request
// context. A nil provider (not injected) leaves role validation skipped.
func (s *StepFunctionService) SetRoleProvider(rp iam.RolePolicyProvider) {
	s.roleProvider = rp
}

// SetTaskCredentialsAuthoriser injects the seam the executors use to run
// the Task Credentials assume-role authorisation at dispatch. A nil
// authoriser (not injected) leaves task invocations on the machine-role
// behaviour.
func (s *StepFunctionService) SetTaskCredentialsAuthoriser(a TaskCredentialsAuthoriser) {
	s.taskCredentialsAuthz = a
}

// iamValidator builds the IAM validator the Cores use. A nil result (no
// role provider injected) leaves role validation skipped.
func (s *StepFunctionService) iamValidator() *iam.IAMValidator {
	if s.roleProvider == nil {
		return nil
	}
	return iam.NewIAMValidator(s.roleProvider, s.accountID)
}

// SetEventBus injects the event bus and subscribes to cross-service start
// execution events from EventBridge, Scheduler, and CloudWatch Alarms.
func (s *StepFunctionService) SetEventBus(bus eventbus.ServiceBus) error {
	s.bus = bus
	if _, err := eventbus.SubscribeTyped[*eventbus.StepFunctionsStartExecutionEvent](bus, s.handleStartExecutionEvent, eventbus.WithAsync()); err != nil {
		return fmt.Errorf("sfn: subscribe StepFunctionsStartExecutionEvent: %w", err)
	}
	return nil
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
		// The same propagation contract as the core-error path below: a
		// store-acquisition fault is a failed delivery, never a success.
		return eventbus.HandlerResult{Error: err}
	}

	// The bus start path runs the same validation, ARN resolution and
	// launch logic as the HTTP StartExecution handler; IoT rule actions
	// send stateMachineName (per Smithy), which the Core resolves.
	if err := s.startExecutionForBusCore(ctx, store, evt.StateMachineArn, evt.StateMachineName, evt.Input); err != nil {
		logs.Error("sfn: failed to start execution from bus event",
			logs.String("arn", evt.StateMachineArn),
			logs.String("name", evt.StateMachineName),
			logs.Err(err))
		// Synchronous publishers (the Scheduler engine publishes via
		// PublishSync) rely on the handler error to drive their retry and
		// dead-letter policy, so the failure must be propagated.
		return eventbus.HandlerResult{Error: err}
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
