// Package iotevents provides AWS IoT Events service operations for vorpalstacks.
package iotevents

import (
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// IoTEventsService provides IoT Events service operations.
type IoTEventsService struct {
	accountID      string
	region         string
	stores         sync.Map // region → iotstore.IotStoreInterface
	storageManager *storage.RegionStorageManager
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

func (s *IoTEventsService) store(reqCtx *request.RequestContext) (iotstore.IotStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (iotstore.IotStoreInterface, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return iotstore.NewIotStore(storage, s.accountID, reqCtx.GetRegion()), nil
	})
}

// GetStoreForRegion returns the cached IoT store for the given region.
// Uses LoadOrStore to prevent duplicate store creation under concurrent access.
func (s *IoTEventsService) GetStoreForRegion(region string) (iotstore.IotStoreInterface, error) {
	store, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	st := iotstore.NewIotStore(store, s.accountID, region)
	if actual, loaded := s.stores.LoadOrStore(region, st); loaded {
		return actual.(iotstore.IotStoreInterface), nil
	}
	return st, nil
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
}
