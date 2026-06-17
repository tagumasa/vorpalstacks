// Package iotevents provides AWS IoT Events service operations for vorpalstacks.
//
// This package imports the IoT store (internal/store/aws/iot) as a declared
// exception per AGENTS.md rule #19, because IoT Events is an IoT sub-service
// that shares the same data model (detector models, inputs, state machines).
package iotevents

import (
	"log/slog"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// IoTEventsServiceDeps holds dependencies for IoT Events initialisation.
type IoTEventsServiceDeps struct {
	EventBus       eventbus.Bus
	StorageManager *storage.RegionStorageManager
}

// IoTEventsService provides IoT Events service operations.
type IoTEventsService struct {
	accountID      string
	region         string
	stores         sync.Map // region → iotstore.IotStoreInterface
	storageManager *storage.RegionStorageManager
	actionAdapter  *DetectorActionAdapter
	once           sync.Once
	initialised    bool
}

// NewIoTEventsService creates a new IoT Events service.
func NewIoTEventsService(accountID, region string) *IoTEventsService {
	return &IoTEventsService{
		accountID: accountID,
		region:    region,
	}
}

// SetStorageManager injects the region storage manager for lazy store creation.
func (s *IoTEventsService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// Init sets all dependencies atomically, guarded by sync.Once.
// Must be called before RegisterHandlers or any operation.
func (s *IoTEventsService) Init(deps IoTEventsServiceDeps) {
	s.once.Do(func() {
		adapter := NewDetectorActionAdapter(deps.EventBus, slog.Default())
		s.actionAdapter = adapter
		s.hydrateDetectorModels()
		s.wireBatchEvaluate(adapter)
		s.initialised = true
	})
}

func (s *IoTEventsService) hydrateDetectorModels() {
	if s.actionAdapter == nil || s.storageManager == nil {
		return
	}
	regions := s.storageManager.ListRegions()
	if len(regions) == 0 {
		return
	}
	total := 0
	for _, region := range regions {
		store, err := s.GetStoreForRegion(region)
		if err != nil {
			slog.Warn("failed to get store for IoT Events detector hydration", "region", region, "error", err)
			continue
		}

		// LoadDetectorModel/BatchEvaluate are on *IotStore, not the interface.
		// Type-assert to the concrete type for state machine hydration.
		concrete, ok := store.(*iotstore.IotStore)
		if !ok {
			slog.Warn("store type assertion failed for detector model hydration", "region", region)
			continue
		}

		var opts storecommon.ListOptions
		for {
			result, err := store.ListDetectorModels(opts)
			if err != nil {
				slog.Warn("failed to list detector models for hydration", "region", region, "error", err)
				break
			}
			for _, dm := range result.Items {
				concrete.LoadDetectorModel(dm)
				total++
			}
			if result.NextMarker == "" {
				break
			}
			opts.Marker = result.NextMarker
		}
	}
	slog.Info("hydrated IoT Events detector models", "total", total)
}

func (s *IoTEventsService) wireBatchEvaluate(adapter *DetectorActionAdapter) {
	if s.storageManager == nil {
		return
	}
	regions := s.storageManager.ListRegions()
	for _, region := range regions {
		store, err := s.GetStoreForRegion(region)
		if err != nil {
			continue
		}
		concrete, ok := store.(*iotstore.IotStore)
		if !ok {
			continue
		}
		adapter.SetBatchEvaluate(concrete.BatchEvaluate)
		break // all regions share the same adapter; one is sufficient
	}
}

func (s *IoTEventsService) store(reqCtx *request.RequestContext) (iotstore.IotStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (iotstore.IotStoreInterface, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		var onAction func(string, string, string, map[string]interface{})
		if s.actionAdapter != nil {
			onAction = s.actionAdapter.OnAction()
		}
		return iotstore.NewIotStore(storage, s.accountID, reqCtx.GetRegion(), onAction), nil
	})
}

// GetStoreForRegion returns the cached IoT store for the given region.
// Uses LoadOrStore to prevent duplicate store creation under concurrent access.
func (s *IoTEventsService) GetStoreForRegion(region string) (iotstore.IotStoreInterface, error) {
	store, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	var onAction func(string, string, string, map[string]interface{})
	if s.actionAdapter != nil {
		onAction = s.actionAdapter.OnAction()
	}
	st := iotstore.NewIotStore(store, s.accountID, region, onAction)
	if actual, loaded := s.stores.LoadOrStore(region, st); loaded {
		return actual.(iotstore.IotStoreInterface), nil
	}
	return st, nil
}

func (s *IoTEventsService) SetRepublishFn(fn func(topic string, payload []byte) error) {
	if s.actionAdapter != nil {
		s.actionAdapter.SetRepublishFn(fn)
	}
}

// RegisterHandlers registers the IoT Events handlers with the dispatcher.
func (s *IoTEventsService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("iotevents", "CreateDetectorModel", s.CreateDetectorModel)
	d.RegisterHandlerForService("iotevents", "DescribeDetectorModel", s.DescribeDetectorModel)
	d.RegisterHandlerForService("iotevents", "UpdateDetectorModel", s.UpdateDetectorModel)
	d.RegisterHandlerForService("iotevents", "DeleteDetectorModel", s.DeleteDetectorModel)
	d.RegisterHandlerForService("iotevents", "ListDetectorModels", s.ListDetectorModels)
	d.RegisterHandlerForService("iotevents", "CreateInput", s.CreateInput)
	d.RegisterHandlerForService("iotevents", "DescribeInput", s.DescribeInput)
	d.RegisterHandlerForService("iotevents", "UpdateInput", s.UpdateInput)
	d.RegisterHandlerForService("iotevents", "DeleteInput", s.DeleteInput)
	d.RegisterHandlerForService("iotevents", "ListInputs", s.ListInputs)
	d.RegisterHandlerForService("iotevents", "TagResource", s.TagResource)
	d.RegisterHandlerForService("iotevents", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("iotevents", "ListTagsForResource", s.ListTagsForResource)

	d.RegisterHandlerForService("iotevents", "BatchPutMessage", s.BatchPutMessage)
	d.RegisterHandler("BatchPutMessage", s.BatchPutMessage)
}
