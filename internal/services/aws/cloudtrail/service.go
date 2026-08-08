// Package cloudtrail provides AWS CloudTrail service operations for vorpalstacks.
package cloudtrail

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// StoreInterface is a type alias for the CloudTrail store interface. It
// allows admin_handler.go to reference the store type without importing the
// store package directly (store-import prohibition).
type StoreInterface = cloudtrailstore.CloudTrailStoreInterface

// CloudTrailService provides AWS CloudTrail operations.
type CloudTrailService struct {
	accountID      string
	region         string
	stores         sync.Map // region → cloudtrailstore.CloudTrailStoreInterface
	storageManager *storage.RegionStorageManager
}

// NewCloudTrailService creates a new CloudTrail service instance.
func NewCloudTrailService(accountID, region string) *CloudTrailService {
	return &CloudTrailService{
		accountID: accountID,
		region:    region,
	}
}

// SetStorageManager injects the region storage manager for lazy store creation.
func (s *CloudTrailService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// GetEventStore returns the shared CloudTrail event store for the given region.
// This ensures a single store instance per region across the service, the audit
// recorder factory, and the S3 audit recorder.
func (s *CloudTrailService) GetEventStore(store storage.BasicStorage, region string) cloudtrailstore.CloudTrailStoreInterface {
	st, _ := storecommon.GetOrCreateStoreE(&s.stores, region, func() (cloudtrailstore.CloudTrailStoreInterface, error) {
		return cloudtrailstore.NewCloudTrailStore(store, s.accountID, region), nil
	})
	return st
}

// GetStoreForRegion returns the cached CloudTrail store for the given region,
// creating a new store instance if not already cached.
func (s *CloudTrailService) GetStoreForRegion(region string) (cloudtrailstore.CloudTrailStoreInterface, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(cloudtrailstore.CloudTrailStoreInterface), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("cloudtrail storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	store := cloudtrailstore.NewCloudTrailStore(st, s.accountID, region)
	actual, _ := s.stores.LoadOrStore(region, store)
	return actual.(cloudtrailstore.CloudTrailStoreInterface), nil
}

func (s *CloudTrailService) store(reqCtx *request.RequestContext) (cloudtrailstore.CloudTrailStoreInterface, error) {
	st, err := storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (cloudtrailstore.CloudTrailStoreInterface, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		return cloudtrailstore.NewCloudTrailStore(storage, s.accountID, reqCtx.GetRegion()), nil
	})
	return st, err
}

// RegisterHandlers registers CloudTrail operation handlers with the dispatcher.
func (s *CloudTrailService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("cloudtrail", "CreateTrail", s.CreateTrail)
	d.RegisterHandlerForService("cloudtrail", "GetTrail", s.GetTrail)
	d.RegisterHandlerForService("cloudtrail", "UpdateTrail", s.UpdateTrail)
	d.RegisterHandlerForService("cloudtrail", "DeleteTrail", s.DeleteTrail)
	d.RegisterHandlerForService("cloudtrail", "StartLogging", s.StartLogging)
	d.RegisterHandlerForService("cloudtrail", "StopLogging", s.StopLogging)
	d.RegisterHandlerForService("cloudtrail", "GetTrailStatus", s.GetTrailStatus)
	d.RegisterHandlerForService("cloudtrail", "ListTrails", s.ListTrails)
	d.RegisterHandlerForService("cloudtrail", "AddTags", s.AddTags)
	d.RegisterHandlerForService("cloudtrail", "RemoveTags", s.RemoveTags)
	d.RegisterHandlerForService("cloudtrail", "ListTags", s.ListTags)
	d.RegisterHandlerForService("cloudtrail", "LookupEvents", s.LookupEvents)
	d.RegisterHandlerForService("cloudtrail", "GetEventSelectors", s.GetEventSelectors)
	d.RegisterHandlerForService("cloudtrail", "PutEventSelectors", s.PutEventSelectors)
	d.RegisterHandlerForService("cloudtrail", "DescribeTrails", s.DescribeTrails)
	d.RegisterHandlerForService("cloudtrail", "ListPublicKeys", s.ListPublicKeys)
	d.RegisterHandlerForService("cloudtrail", "GetInsightSelectors", s.GetInsightSelectors)
	d.RegisterHandlerForService("cloudtrail", "PutInsightSelectors", s.PutInsightSelectors)
	d.RegisterHandlerForService("cloudtrail", "GetResourcePolicy", s.GetResourcePolicy)
	d.RegisterHandlerForService("cloudtrail", "PutResourcePolicy", s.PutResourcePolicy)
	d.RegisterHandlerForService("cloudtrail", "DeleteResourcePolicy", s.DeleteResourcePolicy)
	d.RegisterHandlerForService("cloudtrail", "CreateEventDataStore", s.CreateEventDataStore)
	d.RegisterHandlerForService("cloudtrail", "GetEventDataStore", s.GetEventDataStore)
	d.RegisterHandlerForService("cloudtrail", "ListEventDataStores", s.ListEventDataStores)
	d.RegisterHandlerForService("cloudtrail", "UpdateEventDataStore", s.UpdateEventDataStore)
	d.RegisterHandlerForService("cloudtrail", "DeleteEventDataStore", s.DeleteEventDataStore)
	d.RegisterHandlerForService("cloudtrail", "StartEventDataStoreIngestion", s.StartEventDataStoreIngestion)
	d.RegisterHandlerForService("cloudtrail", "StopEventDataStoreIngestion", s.StopEventDataStoreIngestion)
	d.RegisterHandlerForService("cloudtrail", "RestoreEventDataStore", s.RestoreEventDataStore)
	d.RegisterHandlerForService("cloudtrail", "EnableFederation", s.EnableFederation)
	d.RegisterHandlerForService("cloudtrail", "DisableFederation", s.DisableFederation)
	d.RegisterHandlerForService("cloudtrail", "StartQuery", s.StartQuery)
	d.RegisterHandlerForService("cloudtrail", "GetQueryResults", s.GetQueryResults)
	d.RegisterHandlerForService("cloudtrail", "DescribeQuery", s.DescribeQuery)
	d.RegisterHandlerForService("cloudtrail", "CancelQuery", s.CancelQuery)
	d.RegisterHandlerForService("cloudtrail", "ListQueries", s.ListQueries)
	d.RegisterHandlerForService("cloudtrail", "CreateChannel", s.CreateChannel)
	d.RegisterHandlerForService("cloudtrail", "DeleteChannel", s.DeleteChannel)
	d.RegisterHandlerForService("cloudtrail", "GetChannel", s.GetChannel)
	d.RegisterHandlerForService("cloudtrail", "ListChannels", s.ListChannels)
	d.RegisterHandlerForService("cloudtrail", "UpdateChannel", s.UpdateChannel)
	d.RegisterHandlerForService("cloudtrail", "GetEventConfiguration", s.GetEventConfiguration)
	d.RegisterHandlerForService("cloudtrail", "PutEventConfiguration", s.PutEventConfiguration)
	d.RegisterHandlerForService("cloudtrail", "RegisterOrganizationDelegatedAdmin", s.RegisterOrganizationDelegatedAdmin)
	d.RegisterHandlerForService("cloudtrail", "DeregisterOrganizationDelegatedAdmin", s.DeregisterOrganizationDelegatedAdmin)
	d.RegisterHandlerForService("cloudtrail", "StartImport", s.StartImport)
	d.RegisterHandlerForService("cloudtrail", "StopImport", s.StopImport)
	d.RegisterHandlerForService("cloudtrail", "GetImport", s.GetImport)
	d.RegisterHandlerForService("cloudtrail", "ListImports", s.ListImports)
	d.RegisterHandlerForService("cloudtrail", "ListImportFailures", s.ListImportFailures)
	d.RegisterHandlerForService("cloudtrail", "GenerateQuery", s.GenerateQuery)
	d.RegisterHandlerForService("cloudtrail", "SearchSampleQueries", s.SearchSampleQueries)
}
