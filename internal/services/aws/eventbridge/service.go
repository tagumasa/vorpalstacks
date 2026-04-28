// Package events provides AWS EventBridge (CloudWatch Events) service operations for vorpalstacks.
package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// EventsService provides AWS EventBridge operations.
type EventsService struct {
	storageManager  *storage.RegionStorageManager
	eventsStores    sync.Map // region → *eventsstore.EventsStore
	accountID       string
	bus             eventbus.Bus
	targetSemaphore chan struct{}
	replayWg        sync.WaitGroup
	replayCancels   sync.Map // replayName → context.CancelFunc
}

const targetConcurrencyLimit = 100

// NewEventsService creates a new Events service instance.
// Optional cross-service dependencies should be injected via setter methods
// before registering handlers.
func NewEventsService(storageMgr *storage.RegionStorageManager, accountID string) *EventsService {
	return &EventsService{
		storageManager:  storageMgr,
		accountID:       accountID,
		targetSemaphore: make(chan struct{}, targetConcurrencyLimit),
	}
}

// Close waits for all in-flight replay goroutines to finish, then cancels
// any remaining replay contexts.
func (s *EventsService) Close() {
	s.replayWg.Wait()
	s.replayCancels.Range(func(key, value any) bool {
		if cancel, ok := value.(context.CancelFunc); ok {
			cancel()
		}
		s.replayCancels.Delete(key)
		return true
	})
}

// SetEventsStore sets a pre-built events store for the given region,
// bypassing per-request store creation.
func (s *EventsService) SetEventsStore(region string, store *eventsstore.EventsStore) {
	if store != nil {
		s.eventsStores.Store(region, store)
	}
}

// SetEventBus injects the event bus and registers the EventBridge delivery handler.
func (s *EventsService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
	_, _ = eventbus.SubscribeTyped[*eventbus.EventBridgeDeliveryEvent](bus, s.handleBusDelivery, eventbus.WithAsync())
	_, _ = eventbus.SubscribeTyped[*eventbus.EventBridgePutEventsEvent](bus, s.handlePutEventsEvent, eventbus.WithAsync())
}

func (s *EventsService) handleBusDelivery(ctx context.Context, evt *eventbus.EventBridgeDeliveryEvent) eventbus.HandlerResult {
	var store *eventsstore.EventsStore
	if v, ok := s.eventsStores.Load(evt.Region); ok {
		store = v.(*eventsstore.EventsStore)
	}
	if store == nil && strings.Contains(evt.TargetARN, ":event-bus/") {
		return s.handleEventBusDelivery(ctx, evt)
	}

	targetType := s.parseTargetType(evt.TargetARN)
	target := eventsstore.Target{
		ARN: evt.TargetARN,
		ID:  evt.TargetID,
	}

	event := &eventsstore.Event{
		ID: evt.EventID(),
	}

	s.dispatchToTarget(ctx, evt.Region, event, target, targetType, evt.Input)
	return eventbus.HandlerResult{}
}

func (s *EventsService) handleEventBusDelivery(ctx context.Context, evt *eventbus.EventBridgeDeliveryEvent) eventbus.HandlerResult {
	var event eventsstore.Event
	if err := json.Unmarshal(evt.Input, &event); err != nil {
		logs.Warn("eventbridge: failed to unmarshal event for rule matching", logs.String("error", err.Error()))
		return eventbus.HandlerResult{}
	}

	eventBusName := "default"
	if idx := strings.Index(evt.TargetARN, ":event-bus/"); idx != -1 {
		eventBusName = evt.TargetARN[idx+len(":event-bus/"):]
	}

	var es *eventsstore.EventsStore
	if v, ok := s.eventsStores.Load(evt.Region); ok {
		es = v.(*eventsstore.EventsStore)
	}
	if es == nil {
		return eventbus.HandlerResult{}
	}

	if err := s.deliverEventWithStore(ctx, evt.Region, &event, eventBusName, es); err != nil {
		logs.Warn("eventbridge: failed to deliver event via rule matching", logs.String("error", err.Error()))
	}

	return eventbus.HandlerResult{}
}

func (s *EventsService) handlePutEventsEvent(ctx context.Context, evt *eventbus.EventBridgePutEventsEvent) eventbus.HandlerResult {
	region := evt.Region
	if region == "" {
		return eventbus.HandlerResult{}
	}

	var es *eventsstore.EventsStore
	if v, ok := s.eventsStores.Load(region); ok {
		es = v.(*eventsstore.EventsStore)
	}
	if es == nil {
		st, err := s.storageManager.GetStorage(region)
		if err != nil {
			logs.Warn("eventbridge: failed to get storage for putEvents bus event",
				logs.String("region", region),
				logs.Err(err))
			return eventbus.HandlerResult{}
		}
		es = eventsstore.NewEventsStore(st, s.accountID, region)
		s.eventsStores.Store(region, es)
	}

	var inputMap map[string]interface{}
	if err := json.Unmarshal([]byte(evt.Input), &inputMap); err != nil {
		logs.Warn("eventbridge: failed to unmarshal putEvents input",
			logs.String("region", region),
			logs.Err(err))
		return eventbus.HandlerResult{}
	}

	source, _ := inputMap["Source"].(string)
	detailType, _ := inputMap["DetailType"].(string)
	if source == "" || detailType == "" {
		logs.Warn("eventbridge: putEvents input missing Source or DetailType",
			logs.String("region", region))
		return eventbus.HandlerResult{}
	}

	var detail map[string]interface{}
	if detailRaw, ok := inputMap["Detail"]; ok {
		switch d := detailRaw.(type) {
		case map[string]interface{}:
			detail = d
		case string:
			_ = json.Unmarshal([]byte(d), &detail)
		}
	}

	eventBusName := evt.EventBusName
	if eventBusName == "" {
		eventBusName = "default"
	}

	event := &eventsstore.Event{
		ID:           generateEventID(),
		Version:      "0",
		DetailType:   detailType,
		Source:       source,
		Account:      s.accountID,
		Time:         time.Now().UTC(),
		Region:       region,
		Detail:       detail,
		EventBusName: eventBusName,
	}

	if resources, ok := inputMap["Resources"].([]interface{}); ok {
		for _, r := range resources {
			if rStr, ok := r.(string); ok {
				event.Resources = append(event.Resources, rStr)
			}
		}
	}

	if err := s.deliverEventWithStore(ctx, region, event, eventBusName, es); err != nil {
		logs.Warn("eventbridge: failed to deliver putEvents bus event",
			logs.String("region", region),
			logs.String("error", err.Error()))
	}

	return eventbus.HandlerResult{}
}

func (s *EventsService) store(ctx *request.RequestContext) (*eventsstore.EventsStore, error) {
	return storecommon.GetOrCreateStoreE(&s.eventsStores, ctx.GetRegion(), func() (*eventsstore.EventsStore, error) {
		storage, err := ctx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		es := eventsstore.NewEventsStore(storage, s.accountID, ctx.GetRegion())
		defaultBus := &eventsstore.EventBus{Name: "default"}
		if err := es.CreateEventBus(context.Background(), defaultBus); err != nil {
			if err != eventsstore.ErrEventBusAlreadyExists {
				logs.Warn("eventbridge: failed to auto-create default event bus", logs.Err(err))
			}
		}
		return es, nil
	})
}

// RegisterHandlers registers the Events service handlers with the dispatcher.
func (s *EventsService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("eventbridge", "CreateEventBus", s.CreateEventBus)
	d.RegisterHandlerForService("eventbridge", "DeleteEventBus", s.DeleteEventBus)
	d.RegisterHandlerForService("eventbridge", "DescribeEventBus", s.DescribeEventBus)
	d.RegisterHandlerForService("eventbridge", "ListEventBuses", s.ListEventBuses)
	d.RegisterHandlerForService("eventbridge", "UpdateEventBus", s.UpdateEventBus)

	d.RegisterHandlerForService("eventbridge", "PutRule", s.PutRule)
	d.RegisterHandlerForService("eventbridge", "DeleteRule", s.DeleteRule)
	d.RegisterHandlerForService("eventbridge", "DescribeRule", s.DescribeRule)
	d.RegisterHandlerForService("eventbridge", "ListRules", s.ListRules)
	d.RegisterHandlerForService("eventbridge", "EnableRule", s.EnableRule)
	d.RegisterHandlerForService("eventbridge", "DisableRule", s.DisableRule)

	d.RegisterHandlerForService("eventbridge", "PutTargets", s.PutTargets)
	d.RegisterHandlerForService("eventbridge", "RemoveTargets", s.RemoveTargets)
	d.RegisterHandlerForService("eventbridge", "ListTargetsByRule", s.ListTargetsByRule)
	d.RegisterHandlerForService("eventbridge", "ListRuleNamesByTarget", s.ListRuleNamesByTarget)

	d.RegisterHandlerForService("eventbridge", "PutEvents", s.PutEvents)

	d.RegisterHandlerForService("eventbridge", "TagResource", s.TagResource)
	d.RegisterHandlerForService("eventbridge", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("eventbridge", "ListTagsForResource", s.ListTagsForResource)

	d.RegisterHandlerForService("eventbridge", "CreateArchive", s.CreateArchive)
	d.RegisterHandlerForService("eventbridge", "DeleteArchive", s.DeleteArchive)
	d.RegisterHandlerForService("eventbridge", "DescribeArchive", s.DescribeArchive)
	d.RegisterHandlerForService("eventbridge", "UpdateArchive", s.UpdateArchive)
	d.RegisterHandlerForService("eventbridge", "ListArchives", s.ListArchives)

	d.RegisterHandlerForService("eventbridge", "StartReplay", s.StartReplay)
	d.RegisterHandlerForService("eventbridge", "DescribeReplay", s.DescribeReplay)
	d.RegisterHandlerForService("eventbridge", "ListReplays", s.ListReplays)
	d.RegisterHandlerForService("eventbridge", "CancelReplay", s.CancelReplay)

	d.RegisterHandlerForService("eventbridge", "CreateConnection", s.CreateConnection)
	d.RegisterHandlerForService("eventbridge", "DeleteConnection", s.DeleteConnection)
	d.RegisterHandlerForService("eventbridge", "DescribeConnection", s.DescribeConnection)
	d.RegisterHandlerForService("eventbridge", "UpdateConnection", s.UpdateConnection)
	d.RegisterHandlerForService("eventbridge", "DeauthorizeConnection", s.DeauthorizeConnection)
	d.RegisterHandlerForService("eventbridge", "ListConnections", s.ListConnections)

	d.RegisterHandlerForService("eventbridge", "CreateApiDestination", s.CreateApiDestination)
	d.RegisterHandlerForService("eventbridge", "DeleteApiDestination", s.DeleteApiDestination)
	d.RegisterHandlerForService("eventbridge", "DescribeApiDestination", s.DescribeApiDestination)
	d.RegisterHandlerForService("eventbridge", "UpdateApiDestination", s.UpdateApiDestination)
	d.RegisterHandlerForService("eventbridge", "ListApiDestinations", s.ListApiDestinations)

	d.RegisterHandlerForService("eventbridge", "TestEventPattern", s.TestEventPattern)
}
