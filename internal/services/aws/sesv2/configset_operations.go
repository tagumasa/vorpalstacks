package sesv2

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
)

// CreateConfigurationSet creates a new configuration set for SESv2.
func (s *SESv2Service) CreateConfigurationSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.createConfigurationSetCore(store, CreateConfigurationSetInput{
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		Tags:                 tags.ParseTags(req.Parameters, "Tags"),
		DeliveryOptions:      request.GetMapParam(req.Parameters, "DeliveryOptions"),
		ReputationOptions:    request.GetMapParam(req.Parameters, "ReputationOptions"),
		SendingOptions:       request.GetMapParam(req.Parameters, "SendingOptions"),
		TrackingOptions:      request.GetMapParam(req.Parameters, "TrackingOptions"),
		SuppressionOptions:   request.GetMapParam(req.Parameters, "SuppressionOptions"),
		VdmOptions:           request.GetMapParam(req.Parameters, "VdmOptions"),
		ArchivingOptions:     request.GetMapParam(req.Parameters, "ArchivingOptions"),
	})
}

// GetConfigurationSet retrieves the configuration set details.
func (s *SESv2Service) GetConfigurationSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getConfigurationSetCore(store, request.GetStringParam(req.Parameters, "ConfigurationSetName"))
}

// DeleteConfigurationSet deletes the specified configuration set.
func (s *SESv2Service) DeleteConfigurationSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteConfigurationSetCore(store, request.GetStringParam(req.Parameters, "ConfigurationSetName")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListConfigurationSets returns a list of configuration sets for the account.
func (s *SESv2Service) ListConfigurationSets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listConfigurationSetsCore(store,
		request.GetIntParam(req.Parameters, "PageSize"),
		request.GetStringParam(req.Parameters, "NextToken"))
}

// CreateConfigurationSetEventDestination creates an event destination for a configuration set.
func (s *SESv2Service) CreateConfigurationSetEventDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	defMap := request.GetMapParam(req.Parameters, "EventDestination")
	if err := s.createConfigurationSetEventDestinationCore(store, CreateConfigurationSetEventDestinationInput{
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		EventDestinationName: request.GetStringParam(req.Parameters, "EventDestinationName"),
		EventDestination:     EventDestinationInput{Map: defMap, Provided: defMap != nil},
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// GetConfigurationSetEventDestinations retrieves the event destinations for a configuration set.
func (s *SESv2Service) GetConfigurationSetEventDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getConfigurationSetEventDestinationsCore(store, request.GetStringParam(req.Parameters, "ConfigurationSetName"))
}

// UpdateConfigurationSetEventDestination updates an event destination for a configuration set.
func (s *SESv2Service) UpdateConfigurationSetEventDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	defMap := request.GetMapParam(req.Parameters, "EventDestination")
	if err := s.updateConfigurationSetEventDestinationCore(store, UpdateConfigurationSetEventDestinationInput{
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		EventDestinationName: request.GetStringParam(req.Parameters, "EventDestinationName"),
		EventDestination:     EventDestinationInput{Map: defMap, Provided: defMap != nil},
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DeleteConfigurationSetEventDestination deletes an event destination from a configuration set.
func (s *SESv2Service) DeleteConfigurationSetEventDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteConfigurationSetEventDestinationCore(store,
		request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		request.GetStringParam(req.Parameters, "EventDestinationName")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// PutConfigurationSetDeliveryOptions updates the delivery options for a configuration set.
func (s *SESv2Service) PutConfigurationSetDeliveryOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putConfigurationSetDeliveryOptionsCore(store, PutConfigurationSetDeliveryOptionsInput{
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		SendingPoolName:      request.GetStringParam(req.Parameters, "SendingPoolName"),
		MaxDeliverySeconds:   request.GetIntParam(req.Parameters, "MaxDeliverySeconds"),
		TlsPolicy:            request.GetStringParam(req.Parameters, "TlsPolicy"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// PutConfigurationSetReputationOptions updates the reputation options for a configuration set.
func (s *SESv2Service) PutConfigurationSetReputationOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putConfigurationSetReputationOptionsCore(store, PutConfigurationSetReputationOptionsInput{
		ConfigurationSetName:     request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		ReputationMetricsEnabled: request.GetBoolParam(req.Parameters, "ReputationMetricsEnabled"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// PutConfigurationSetSendingOptions updates the sending options for a configuration set.
func (s *SESv2Service) PutConfigurationSetSendingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putConfigurationSetSendingOptionsCore(store, PutConfigurationSetSendingOptionsInput{
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		SendingEnabled:       request.GetBoolParam(req.Parameters, "SendingEnabled"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// PutConfigurationSetSuppressionOptions updates the suppression options for a configuration set.
func (s *SESv2Service) PutConfigurationSetSuppressionOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putConfigurationSetSuppressionOptionsCore(store, PutConfigurationSetSuppressionOptionsInput{
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		SuppressedReasons:    request.GetStringList(req.Parameters, "SuppressedReasons"),
		SuppressionScope:     request.GetStringParam(req.Parameters, "SuppressionScope"),
		ValidationOptions:    request.GetMapParam(req.Parameters, "ValidationOptions"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// PutConfigurationSetTrackingOptions updates the tracking options for a configuration set.
func (s *SESv2Service) PutConfigurationSetTrackingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putConfigurationSetTrackingOptionsCore(store, PutConfigurationSetTrackingOptionsInput{
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		CustomRedirectDomain: request.GetStringParam(req.Parameters, "CustomRedirectDomain"),
		HttpsPolicy:          request.GetStringParam(req.Parameters, "HttpsPolicy"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// PutConfigurationSetVdmOptions updates the VDM options for a configuration set.
func (s *SESv2Service) PutConfigurationSetVdmOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putConfigurationSetVdmOptionsCore(store, PutConfigurationSetVdmOptionsInput{
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		VdmOptions:           request.GetMapParam(req.Parameters, "VdmOptions"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// PutConfigurationSetArchivingOptions updates the archiving options for a configuration set.
func (s *SESv2Service) PutConfigurationSetArchivingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putConfigurationSetArchivingOptionsCore(store, PutConfigurationSetArchivingOptionsInput{
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		ArchiveArn:           request.GetStringParam(req.Parameters, "ArchiveArn"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}
