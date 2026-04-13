// Package scheduler provides AWS EventBridge Scheduler service operations for vorpalstacks.
package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/common"
	"vorpalstacks/internal/common/request"
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

// SetEventBus injects the event bus into the scheduler engine and registers
// the ScheduleFiredEvent handler. When the bus is set, schedule execution
// routes through the bus instead of direct store/invoker calls.
func (s *SchedulerService) SetEventBus(bus *eventbus.EventBus) {
	if s.engine != nil {
		s.engine.SetEventBus(bus)
		_, _ = eventbus.SubscribeTyped[*eventbus.ScheduleFiredEvent](bus, s.handleBusDelivery, eventbus.WithAsync())
	}
}

func (s *SchedulerService) handleBusDelivery(ctx context.Context, evt *eventbus.ScheduleFiredEvent) eventbus.HandlerResult {
	if s.engine == nil {
		return eventbus.HandlerResult{}
	}

	target := &schedulerstore.Target{
		Arn:   evt.TargetArn,
		Input: evt.Input,
	}
	schedule := &schedulerstore.Schedule{
		Name:   evt.ScheduleName,
		Region: evt.Region,
		Target: target,
	}

	if strings.Contains(evt.TargetArn, ":lambda:") {
		s.engine.invokeLambda(ctx, schedule, target)
	} else if strings.Contains(evt.TargetArn, ":sqs:") {
		s.engine.sendToSQS(ctx, schedule, target)
	} else if strings.Contains(evt.TargetArn, ":sns:") {
		if s.engine.bus != nil {
			message := target.Input
			if message == "" {
				msgPayload := map[string]interface{}{
					"schedule":  schedule.Name,
					"timestamp": time.Now().UTC().Format(time.RFC3339),
				}
				if msgBytes, err := json.Marshal(msgPayload); err == nil {
					message = string(msgBytes)
				}
			}
			snsEvt := &eventbus.SNSDeliveryEvent{
				TopicARN:  evt.TargetArn,
				MessageID: fmt.Sprintf("%d", time.Now().UnixNano()),
				Message:   message,
			}
			snsEvt.Region = evt.Region
			s.engine.bus.Publish(context.Background(), snsEvt)
		} else {
			s.engine.publishToSNS(ctx, schedule, target)
		}
	} else if strings.Contains(evt.TargetArn, ":logs:") {
		s.engine.sendToCloudWatchLogs(ctx, schedule, target)
	} else if strings.Contains(evt.TargetArn, ":states:") {
		s.engine.startStepFunctionExecution(ctx, schedule, target)
	} else if strings.Contains(evt.TargetArn, ":events:") {
		s.engine.sendToEventBridge(ctx, schedule, target)
	}

	return eventbus.HandlerResult{}
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
