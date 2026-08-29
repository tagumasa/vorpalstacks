package cloudtrail

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	tags "vorpalstacks/internal/common/tags"
)

// parseBool handles bool values from request parameters which may be bool or
// json.Number.
func parseBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return strings.ToLower(val) == "true"
	default:
		return false
	}
}

// boolParam reads a presence-checked boolean member into a pointer so the
// Core can distinguish an explicit value from an absent member.
func boolParam(params map[string]interface{}, key string) *bool {
	v, ok := params[key]
	if !ok {
		return nil
	}
	b := parseBool(v)
	return &b
}

// CreateEventDataStore creates a new CloudTrail event data store.
func (s *CloudTrailService) CreateEventDataStore(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := CreateEventDataStoreInput{
		Name:                         request.GetStringParam(req.Parameters, "Name"),
		TerminationProtectionEnabled: boolParam(req.Parameters, "TerminationProtectionEnabled"),
		MultiRegionEnabled:           boolParam(req.Parameters, "MultiRegionEnabled"),
		OrganizationEnabled:          boolParam(req.Parameters, "OrganizationEnabled"),
		IngestionEnabled:             boolParam(req.Parameters, "IngestionEnabled"),
		StartIngestion:               boolParam(req.Parameters, "StartIngestion"),
		KmsKeyId:                     request.GetStringParam(req.Parameters, "KmsKeyId"),
		BillingMode:                  request.GetStringParam(req.Parameters, "BillingMode"),
	}
	in.RetentionPeriodRaw, in.RetentionPeriodSet = req.Parameters["RetentionPeriod"]
	if aesRaw, ok := req.Parameters["AdvancedEventSelectors"]; ok {
		in.AdvancedEventSelectorsRaw = aesRaw
		in.AdvancedEventSelectorsSet = true
	}
	if tagsRaw, ok := req.Parameters["TagsList"]; ok {
		in.TagsRaw = tagsRaw
		in.TagsSet = true
		in.TagList = tags.ParseTags(req.Parameters, "TagsList")
	}

	return s.createEventDataStoreCore(store, in)
}

// GetEventDataStore retrieves details about the specified event data store.
func (s *CloudTrailService) GetEventDataStore(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.getEventDataStoreCore(store, EventDataStoreIDInput{
		EventDataStore: request.GetStringParam(req.Parameters, "EventDataStore"),
	})
}

// ListEventDataStores returns event data stores with pagination.
func (s *CloudTrailService) ListEventDataStores(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.listEventDataStoresCore(store, ListEventDataStoresInput{
		NextToken:  req.GetParam("NextToken"),
		MaxResults: request.GetIntParam(req.Parameters, "MaxResults"),
	})
}

// UpdateEventDataStore updates the specified event data store.
func (s *CloudTrailService) UpdateEventDataStore(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := UpdateEventDataStoreInput{
		EventDataStore:               request.GetStringParam(req.Parameters, "EventDataStore"),
		Name:                         request.GetStringParam(req.Parameters, "Name"),
		TerminationProtectionEnabled: boolParam(req.Parameters, "TerminationProtectionEnabled"),
		MultiRegionEnabled:           boolParam(req.Parameters, "MultiRegionEnabled"),
		OrganizationEnabled:          boolParam(req.Parameters, "OrganizationEnabled"),
		IngestionEnabled:             boolParam(req.Parameters, "IngestionEnabled"),
		KmsKeyId:                     request.GetStringParam(req.Parameters, "KmsKeyId"),
		BillingMode:                  request.GetStringParam(req.Parameters, "BillingMode"),
	}
	in.RetentionPeriodRaw, in.RetentionPeriodSet = req.Parameters["RetentionPeriod"]
	if aesRaw, ok := req.Parameters["AdvancedEventSelectors"]; ok {
		in.AdvancedEventSelectorsRaw = aesRaw
		in.AdvancedEventSelectorsSet = true
	}

	return s.updateEventDataStoreCore(store, in)
}

// DeleteEventDataStore deletes the specified event data store.
func (s *CloudTrailService) DeleteEventDataStore(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.deleteEventDataStoreCore(store, EventDataStoreIDInput{
		EventDataStore: request.GetStringParam(req.Parameters, "EventDataStore"),
	})
}

// StartEventDataStoreIngestion enables ingestion on the specified event data
// store.
func (s *CloudTrailService) StartEventDataStoreIngestion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.startEventDataStoreIngestionCore(store, EventDataStoreIDInput{
		EventDataStore: request.GetStringParam(req.Parameters, "EventDataStore"),
	})
}

// StopEventDataStoreIngestion disables ingestion on the specified event data
// store.
func (s *CloudTrailService) StopEventDataStoreIngestion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.stopEventDataStoreIngestionCore(store, EventDataStoreIDInput{
		EventDataStore: request.GetStringParam(req.Parameters, "EventDataStore"),
	})
}

// RestoreEventDataStore restores a PENDING_DELETION event data store.
func (s *CloudTrailService) RestoreEventDataStore(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.restoreEventDataStoreCore(store, EventDataStoreIDInput{
		EventDataStore: request.GetStringParam(req.Parameters, "EventDataStore"),
	})
}

// EnableFederation enables Lake query federation on the specified event data
// store.
func (s *CloudTrailService) EnableFederation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.enableFederationCore(ctx, store, EnableFederationInput{
		EventDataStore:    request.GetStringParam(req.Parameters, "EventDataStore"),
		FederationRoleArn: request.GetStringParam(req.Parameters, "FederationRoleArn"),
		IAMValidator:      reqCtx.GetIAMValidator(),
	})
}

// DisableFederation disables Lake query federation on the specified event data
// store.
func (s *CloudTrailService) DisableFederation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.disableFederationCore(store, DisableFederationInput{
		EventDataStore: request.GetStringParam(req.Parameters, "EventDataStore"),
	})
}
