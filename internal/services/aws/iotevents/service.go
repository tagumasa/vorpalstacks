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
		s.hydrateAlarmModels()
		s.wireBatchEvaluate(adapter)
		s.wireActionCallback(adapter)
		s.initialised = true
	})
}

// wireActionCallback installs the detector action callback onto every
// regional IotStore singleton so that BatchEvaluate dispatches transitions
// through this service's action adapter. Because GetOrCreateStore returns
// the same instance shared with the iot service and the MQTT auth
// provider, wiring the callback here makes every transition observe the
// adapter regardless of which caller drives BatchEvaluate. Init runs under
// sync.Once, so each (account, region) singleton receives the callback
// exactly once per process; a single OnAction closure is reused across
// regions because it only captures the shared adapter pointer.
func (s *IoTEventsService) wireActionCallback(adapter *DetectorActionAdapter) {
	if s.storageManager == nil {
		return
	}
	onAction := adapter.OnAction()
	for _, region := range s.storageManager.ListRegions() {
		store, err := s.storageManager.GetStorage(region)
		if err != nil {
			continue
		}
		iotstore.GetOrCreateStore(store, s.accountID, region).SetActionCallback(onAction)
	}
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

// hydrateAlarmModels loads all persisted alarm models into each regional
// AlarmStateMachine so that BatchPutMessage evaluations work immediately
// after server restart (H-SM2).
func (s *IoTEventsService) hydrateAlarmModels() {
	if s.storageManager == nil {
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
			slog.Warn("failed to get store for IoT Events alarm hydration", "region", region, "error", err)
			continue
		}
		concrete, ok := store.(*iotstore.IotStore)
		if !ok {
			continue
		}
		var opts storecommon.ListOptions
		for {
			result, err := store.ListAlarmModels(opts)
			if err != nil {
				slog.Warn("failed to list alarm models for hydration", "region", region, "error", err)
				break
			}
			for _, am := range result.Items {
				concrete.LoadAlarmModel(am)
				total++
			}
			if result.NextMarker == "" {
				break
			}
			opts.Marker = result.NextMarker
		}
	}
	slog.Info("hydrated IoT Events alarm models", "total", total)
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
	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	return iotstore.GetOrCreateStore(storage, s.accountID, reqCtx.GetRegion()), nil
}

// GetStoreForRegion returns the shared IotStore singleton for the given
// region. The same instance is used by the iot service, the MQTT auth
// provider and this service, so detector model hydration, BatchEvaluate
// and rule-action dispatch all observe one DetectorStateMachine.
func (s *IoTEventsService) GetStoreForRegion(region string) (iotstore.IotStoreInterface, error) {
	store, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return iotstore.GetOrCreateStore(store, s.accountID, region), nil
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
	d.RegisterHandlerForService("iotevents", "CreateAlarmModel", s.CreateAlarmModel)
	d.RegisterHandlerForService("iotevents", "DescribeAlarmModel", s.DescribeAlarmModel)
	d.RegisterHandlerForService("iotevents", "UpdateAlarmModel", s.UpdateAlarmModel)
	d.RegisterHandlerForService("iotevents", "DeleteAlarmModel", s.DeleteAlarmModel)
	d.RegisterHandlerForService("iotevents", "ListAlarmModels", s.ListAlarmModels)
	d.RegisterHandlerForService("iotevents", "ListAlarmModelVersions", s.ListAlarmModelVersions)
	d.RegisterHandlerForService("iotevents", "StartDetectorModelAnalysis", s.StartDetectorModelAnalysis)
	d.RegisterHandlerForService("iotevents", "DescribeDetectorModelAnalysis", s.DescribeDetectorModelAnalysis)
	d.RegisterHandlerForService("iotevents", "GetDetectorModelAnalysisResults", s.GetDetectorModelAnalysisResults)
	d.RegisterHandlerForService("iotevents", "PutLoggingOptions", s.PutLoggingOptions)
	d.RegisterHandlerForService("iotevents", "DescribeLoggingOptions", s.DescribeLoggingOptions)
	d.RegisterHandlerForService("iotevents", "ListInputRoutings", s.ListInputRoutings)
	d.RegisterHandlerForService("iotevents", "ListDetectorModelVersions", s.ListDetectorModelVersions)
	d.RegisterHandlerForService("iotevents", "TagResource", s.TagResource)
	d.RegisterHandlerForService("iotevents", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("iotevents", "ListTagsForResource", s.ListTagsForResource)

	d.RegisterHandlerForService("iotevents-data", "BatchPutMessage", s.BatchPutMessage)
	d.RegisterHandler("BatchPutMessage", s.BatchPutMessage)
	d.RegisterHandlerForService("iotevents-data", "BatchAcknowledgeAlarm", s.BatchAcknowledgeAlarm)
	d.RegisterHandlerForService("iotevents-data", "BatchDisableAlarm", s.BatchDisableAlarm)
	d.RegisterHandlerForService("iotevents-data", "BatchEnableAlarm", s.BatchEnableAlarm)
	d.RegisterHandlerForService("iotevents-data", "BatchResetAlarm", s.BatchResetAlarm)
	d.RegisterHandlerForService("iotevents-data", "BatchSnoozeAlarm", s.BatchSnoozeAlarm)
	d.RegisterHandlerForService("iotevents-data", "BatchDeleteDetector", s.BatchDeleteDetector)
	d.RegisterHandlerForService("iotevents-data", "BatchUpdateDetector", s.BatchUpdateDetector)
	d.RegisterHandlerForService("iotevents-data", "DescribeAlarm", s.DescribeAlarm)
	d.RegisterHandlerForService("iotevents-data", "DescribeDetector", s.DescribeDetector)
	d.RegisterHandlerForService("iotevents-data", "ListAlarms", s.ListAlarms)
	d.RegisterHandlerForService("iotevents-data", "ListDetectors", s.ListDetectors)
}
