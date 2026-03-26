// Package events provides AWS EventBridge (CloudWatch Events) service operations for vorpalstacks.
package eventbridge

import (
	"context"
	"sync"

	appconfig "vorpalstacks/internal/config"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common"
	"vorpalstacks/internal/services/aws/common/request"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// SNSPublisher defines the interface for publishing to SNS topics.
type SNSPublisher interface {
	PublishToTopic(ctx context.Context, accountID, region, topicArn, message string) error
}

// EventsService provides AWS EventBridge operations.
type EventsService struct {
	storageManager  *storage.RegionStorageManager
	eventsStore     *eventsstore.EventsStore
	sqsStores       sync.Map // region → *sqsstore.SQSStore
	snsStores       sync.Map // region → *snsstore.SNSStore
	snsPublisher    SNSPublisher
	lambdaInvoker   common.LambdaInvoker
	accountID       string
	region          string
	targetSemaphore chan struct{}
	replayCancels   sync.Map // replayName → context.CancelFunc
}

const targetConcurrencyLimit = 100

// NewEventsService creates a new Events service instance.
func NewEventsService(storageMgr *storage.RegionStorageManager, accountID, region string) *EventsService {
	return &EventsService{
		storageManager:  storageMgr,
		accountID:       accountID,
		region:          region,
		targetSemaphore: make(chan struct{}, targetConcurrencyLimit),
	}
}

// NewEventsServiceWithClients creates a new Events service instance with custom clients.
func NewEventsServiceWithClients(storageMgr *storage.RegionStorageManager, accountID, region string, sqsStore sqsstore.SQSStoreInterface, snsStore snsstore.SNSStoreInterface, lambdaInvoker common.LambdaInvoker) *EventsService {
	s := &EventsService{
		storageManager:  storageMgr,
		lambdaInvoker:   lambdaInvoker,
		accountID:       accountID,
		region:          region,
		targetSemaphore: make(chan struct{}, targetConcurrencyLimit),
	}
	if sqsStore != nil {
		s.sqsStores.Store(region, sqsStore)
	}
	if snsStore != nil {
		s.snsStores.Store(region, snsStore)
	}
	return s
}

// NewEventsServiceWithStore creates a new Events service instance with a custom store.
func NewEventsServiceWithStore(storageMgr *storage.RegionStorageManager, accountID, region string, eventsStore *eventsstore.EventsStore) *EventsService {
	return &EventsService{
		storageManager:  storageMgr,
		eventsStore:     eventsStore,
		accountID:       accountID,
		region:          region,
		targetSemaphore: make(chan struct{}, targetConcurrencyLimit),
	}
}

// NewEventsServiceWithPublisher creates a new Events service instance with a custom SNS publisher.
func NewEventsServiceWithPublisher(storageMgr *storage.RegionStorageManager, accountID, region string, sqsStore *sqsstore.SQSStore, snsPublisher SNSPublisher, lambdaInvoker common.LambdaInvoker) *EventsService {
	s := &EventsService{
		storageManager:  storageMgr,
		snsPublisher:    snsPublisher,
		lambdaInvoker:   lambdaInvoker,
		accountID:       accountID,
		region:          region,
		targetSemaphore: make(chan struct{}, targetConcurrencyLimit),
	}
	if sqsStore != nil {
		s.sqsStores.Store(region, sqsStore)
	}
	return s
}

func (s *EventsService) store(ctx *request.RequestContext) (*eventsstore.EventsStore, error) {
	if s.eventsStore != nil {
		return s.eventsStore, nil
	}
	if store := ctx.GetEventBridgeStore(); store != nil {
		return store.(*eventsstore.EventsStore), nil
	}
	storage, err := ctx.GetStorage()
	if err != nil {
		return nil, err
	}
	return eventsstore.NewEventsStore(storage, s.accountID, ctx.GetRegion()), nil
}

func (s *EventsService) getSQSStore(ctx *request.RequestContext) (sqsstore.SQSStoreInterface, error) {
	if store := ctx.GetSQSStore(); store != nil {
		return store, nil
	}
	region := ctx.GetRegion()
	return s.getSQSStoreForRegion(region)
}

func (s *EventsService) getSQSStoreForRegion(region string) (sqsstore.SQSStoreInterface, error) {
	if cached, ok := s.sqsStores.Load(region); ok {
		if typed, ok := cached.(sqsstore.SQSStoreInterface); ok {
			return typed, nil
		}
	}
	storage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	store := sqsstore.NewSQSStore(storage, s.accountID, region, appconfig.BaseURL())
	if actual, loaded := s.sqsStores.LoadOrStore(region, store); loaded {
		if typed, ok := actual.(sqsstore.SQSStoreInterface); ok {
			return typed, nil
		}
	}
	return store, nil
}

func (s *EventsService) getSNSStore(ctx *request.RequestContext) (snsstore.SNSStoreInterface, error) {
	if store := ctx.GetSNSStore(); store != nil {
		return store, nil
	}
	region := ctx.GetRegion()
	return s.getSNSStoreForRegion(region)
}

func (s *EventsService) getSNSStoreForRegion(region string) (snsstore.SNSStoreInterface, error) {
	if cached, ok := s.snsStores.Load(region); ok {
		if typed, ok := cached.(snsstore.SNSStoreInterface); ok {
			return typed, nil
		}
	}
	storage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	store := snsstore.NewSNSStore(storage, s.accountID, region)
	if actual, loaded := s.snsStores.LoadOrStore(region, store); loaded {
		if typed, ok := actual.(snsstore.SNSStoreInterface); ok {
			return typed, nil
		}
	}
	return store, nil
}

// RegisterHandlers registers the Events service handlers with the dispatcher.
func (s *EventsService) RegisterHandlers(d *dispatcher.Dispatcher) {
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
