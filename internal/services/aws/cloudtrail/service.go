// Package cloudtrail provides AWS CloudTrail service operations for vorpalstacks.
package cloudtrail

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// CloudTrailService provides AWS CloudTrail operations.
type CloudTrailService struct {
	accountID string
	region    string
	stores    sync.Map // region → cloudtrailstore.CloudTrailStoreInterface
}

// NewCloudTrailService creates a new CloudTrail service instance.
func NewCloudTrailService(accountID, region string) *CloudTrailService {
	return &CloudTrailService{
		accountID: accountID,
		region:    region,
	}
}

func (s *CloudTrailService) store(reqCtx *request.RequestContext) (cloudtrailstore.CloudTrailStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (cloudtrailstore.CloudTrailStoreInterface, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		return cloudtrailstore.NewCloudTrailStore(storage, s.accountID, reqCtx.GetRegion()), nil
	})
}

// RegisterHandlers registers the CloudTrail service handlers with the dispatcher.
func (s *CloudTrailService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("cloudtrail", "CreateTrail", s.CreateTrail)
	d.RegisterHandlerForService("cloudtrail", "DeleteTrail", s.DeleteTrail)
	d.RegisterHandlerForService("cloudtrail", "UpdateTrail", s.UpdateTrail)
	d.RegisterHandlerForService("cloudtrail", "DescribeTrails", s.DescribeTrails)
	d.RegisterHandlerForService("cloudtrail", "GetTrail", s.GetTrail)
	d.RegisterHandlerForService("cloudtrail", "GetTrailStatus", s.GetTrailStatus)
	d.RegisterHandlerForService("cloudtrail", "ListTrails", s.ListTrails)
	d.RegisterHandlerForService("cloudtrail", "StartLogging", s.StartLogging)
	d.RegisterHandlerForService("cloudtrail", "StopLogging", s.StopLogging)

	d.RegisterHandlerForService("cloudtrail", "LookupEvents", s.LookupEvents)
	d.RegisterHandlerForService("cloudtrail", "ListPublicKeys", s.ListPublicKeys)

	d.RegisterHandlerForService("cloudtrail", "GetEventSelectors", s.GetEventSelectors)
	d.RegisterHandlerForService("cloudtrail", "PutEventSelectors", s.PutEventSelectors)
	d.RegisterHandlerForService("cloudtrail", "GetInsightSelectors", s.GetInsightSelectors)
	d.RegisterHandlerForService("cloudtrail", "PutInsightSelectors", s.PutInsightSelectors)

	d.RegisterHandlerForService("cloudtrail", "AddTags", s.AddTags)
	d.RegisterHandlerForService("cloudtrail", "TagResource", s.AddTags)
	d.RegisterHandlerForService("cloudtrail", "RemoveTags", s.RemoveTags)
	d.RegisterHandlerForService("cloudtrail", "UntagResource", s.RemoveTags)
	d.RegisterHandlerForService("cloudtrail", "ListTags", s.ListTags)
	d.RegisterHandlerForService("cloudtrail", "ListTagsForResource", s.ListTags)

	d.RegisterHandlerForService("cloudtrail", "GetResourcePolicy", s.GetResourcePolicy)
	d.RegisterHandlerForService("cloudtrail", "PutResourcePolicy", s.PutResourcePolicy)
	d.RegisterHandlerForService("cloudtrail", "DeleteResourcePolicy", s.DeleteResourcePolicy)
}
