package sesv2

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	"vorpalstacks/internal/services/aws/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// CreateConfigurationSet creates a new configuration set for SESv2.
func (s *SESv2Service) CreateConfigurationSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	parsedTags := tags.ParseTags(req.Parameters, "Tags")

	configSet := sesv2store.NewConfigurationSet(configSetName)

	if deliveryOpts := request.GetMapParam(req.Parameters, "DeliveryOptions"); deliveryOpts != nil {
		configSet.DeliveryOptions = &sesv2store.DeliveryOptions{
			SendingPoolName:    request.GetStringParam(deliveryOpts, "SendingPoolName"),
			MaxDeliverySeconds: int32(request.GetIntParam(deliveryOpts, "MaxDeliverySeconds")),
		}
	}

	if reputationOpts := request.GetMapParam(req.Parameters, "ReputationOptions"); reputationOpts != nil {
		configSet.ReputationOptions = &sesv2store.ReputationOptions{
			ReputationMetricsEnabled: request.GetBoolParam(reputationOpts, "ReputationMetricsEnabled"),
		}
	}

	if sendingOpts := request.GetMapParam(req.Parameters, "SendingOptions"); sendingOpts != nil {
		configSet.SendingOptions = &sesv2store.SendingOptions{
			SendingEnabled: request.GetBoolParam(sendingOpts, "SendingEnabled"),
		}
	}

	if trackingOpts := request.GetMapParam(req.Parameters, "TrackingOptions"); trackingOpts != nil {
		configSet.TrackingOptions = &sesv2store.TrackingOptions{
			CustomRedirectDomain: request.GetStringParam(trackingOpts, "CustomRedirectDomain"),
			HttpsPolicy:          request.GetStringParam(trackingOpts, "HttpsPolicy"),
		}
	}

	_, err = store.CreateConfigurationSet(configSet)
	if err != nil {
		return nil, err
	}

	if len(parsedTags) > 0 {
		arn := store.BuildConfigSetArn(configSetName)
		if err := store.TagResourceFromSlice(arn, parsedTags); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"ConfigurationSetName": configSetName,
	}, nil
}

// GetConfigurationSet retrieves the configuration set details.
func (s *SESv2Service) GetConfigurationSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"ConfigurationSetName": configSet.Name,
	}

	if configSet.SendingOptions != nil {
		response["SendingOptions"] = map[string]interface{}{
			"SendingEnabled": configSet.SendingOptions.SendingEnabled,
		}
	}

	if configSet.ReputationOptions != nil {
		response["ReputationOptions"] = map[string]interface{}{
			"ReputationMetricsEnabled": configSet.ReputationOptions.ReputationMetricsEnabled,
		}
	}

	if configSet.DeliveryOptions != nil {
		response["DeliveryOptions"] = map[string]interface{}{
			"SendingPoolName":    configSet.DeliveryOptions.SendingPoolName,
			"MaxDeliverySeconds": configSet.DeliveryOptions.MaxDeliverySeconds,
		}
	}

	if configSet.TrackingOptions != nil {
		response["TrackingOptions"] = map[string]interface{}{
			"CustomRedirectDomain": configSet.TrackingOptions.CustomRedirectDomain,
			"HttpsPolicy":          configSet.TrackingOptions.HttpsPolicy,
		}
	}

	return response, nil
}

// DeleteConfigurationSet deletes the specified configuration set.
func (s *SESv2Service) DeleteConfigurationSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	if err := store.DeleteConfigurationSet(configSetName); err != nil {
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

	pageSize := request.GetIntParam(req.Parameters, "PageSize")
	if pageSize == 0 {
		pageSize = 100
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	opts := common.ListOptions{
		MaxItems: pageSize,
		Marker:   nextToken,
	}

	result, err := store.ListConfigurationSets(opts)
	if err != nil {
		return nil, err
	}

	configSets := make([]string, 0, len(result.Items))
	for _, cs := range result.Items {
		configSets = append(configSets, cs.Name)
	}

	response := map[string]interface{}{
		"ConfigurationSets": configSets,
	}

	if result.IsTruncated {
		response["NextToken"] = result.NextMarker
	}

	return response, nil
}

// CreateConfigurationSetEventDestination creates an event destination for a configuration set.
func (s *SESv2Service) CreateConfigurationSetEventDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	eventDestinationName := request.GetStringParam(req.Parameters, "EventDestinationName")

	if configSetName == "" || eventDestinationName == "" {
		return nil, ErrMissingParameter
	}

	_, err = store.GetConfigurationSet(configSetName)
	if err != nil {
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

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	_, err = store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"EventDestinations": []interface{}{},
	}, nil
}

// UpdateConfigurationSetEventDestination updates an event destination for a configuration set.
func (s *SESv2Service) UpdateConfigurationSetEventDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	eventDestinationName := request.GetStringParam(req.Parameters, "EventDestinationName")

	if configSetName == "" || eventDestinationName == "" {
		return nil, ErrMissingParameter
	}

	_, err = store.GetConfigurationSet(configSetName)
	if err != nil {
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

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	eventDestinationName := request.GetStringParam(req.Parameters, "EventDestinationName")

	if configSetName == "" || eventDestinationName == "" {
		return nil, ErrMissingParameter
	}

	_, err = store.GetConfigurationSet(configSetName)
	if err != nil {
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

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	sendingPoolName := request.GetStringParam(req.Parameters, "SendingPoolName")
	maxDeliverySeconds := request.GetIntParam(req.Parameters, "MaxDeliverySeconds")

	if configSet.DeliveryOptions == nil {
		configSet.DeliveryOptions = &sesv2store.DeliveryOptions{}
	}

	if sendingPoolName != "" {
		configSet.DeliveryOptions.SendingPoolName = sendingPoolName
	}
	if maxDeliverySeconds > 0 {
		configSet.DeliveryOptions.MaxDeliverySeconds = int32(maxDeliverySeconds)
	}

	if err := store.UpdateConfigurationSet(configSet); err != nil {
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

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	reputationMetricsEnabled := request.GetBoolParam(req.Parameters, "ReputationMetricsEnabled")

	if configSet.ReputationOptions == nil {
		configSet.ReputationOptions = &sesv2store.ReputationOptions{}
	}
	configSet.ReputationOptions.ReputationMetricsEnabled = reputationMetricsEnabled

	if err := store.UpdateConfigurationSet(configSet); err != nil {
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

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	sendingEnabled := request.GetBoolParam(req.Parameters, "SendingEnabled")

	if configSet.SendingOptions == nil {
		configSet.SendingOptions = &sesv2store.SendingOptions{}
	}
	configSet.SendingOptions.SendingEnabled = sendingEnabled

	if err := store.UpdateConfigurationSet(configSet); err != nil {
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

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	suppressedReasons := request.GetStringList(req.Parameters, "SuppressedReasons")

	if configSet.SuppressionOptions == nil {
		configSet.SuppressionOptions = &sesv2store.SuppressionOptions{}
	}
	configSet.SuppressionOptions.SuppressedReasons = suppressedReasons

	if err := store.UpdateConfigurationSet(configSet); err != nil {
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

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	customRedirectDomain := request.GetStringParam(req.Parameters, "CustomRedirectDomain")
	httpsPolicy := request.GetStringParam(req.Parameters, "HttpsPolicy")

	if configSet.TrackingOptions == nil {
		configSet.TrackingOptions = &sesv2store.TrackingOptions{}
	}
	if customRedirectDomain != "" {
		configSet.TrackingOptions.CustomRedirectDomain = customRedirectDomain
	}
	if httpsPolicy != "" {
		configSet.TrackingOptions.HttpsPolicy = httpsPolicy
	}

	if err := store.UpdateConfigurationSet(configSet); err != nil {
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

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	if configSet.VdmOptions == nil {
		configSet.VdmOptions = &sesv2store.VdmOptions{}
	}

	if dashboardOpts := request.GetMapParam(req.Parameters, "DashboardOptions"); dashboardOpts != nil {
		if configSet.VdmOptions.DashboardOptions == nil {
			configSet.VdmOptions.DashboardOptions = &sesv2store.VDMDashboardOptions{}
		}
		configSet.VdmOptions.DashboardOptions.EngagementMetrics = request.GetStringParam(dashboardOpts, "EngagementMetrics")
	}

	if guardianOpts := request.GetMapParam(req.Parameters, "GuardianOptions"); guardianOpts != nil {
		if configSet.VdmOptions.GuardianOptions == nil {
			configSet.VdmOptions.GuardianOptions = &sesv2store.VDMGuardianOptions{}
		}
		configSet.VdmOptions.GuardianOptions.OptimizedSharedDelivery = request.GetStringParam(guardianOpts, "OptimizedSharedDelivery")
	}

	if err := store.UpdateConfigurationSet(configSet); err != nil {
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

	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName == "" {
		return nil, ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	if configSet.ArchivingOptions == nil {
		configSet.ArchivingOptions = &sesv2store.ArchivingOptions{}
	}

	configSet.ArchivingOptions.Enabled = request.GetBoolParam(req.Parameters, "Enabled")
	if targetArn := request.GetStringParam(req.Parameters, "TargetArn"); targetArn != "" {
		configSet.ArchivingOptions.TargetArn = targetArn
	}
	if retentionPeriod := int32(request.GetIntParam(req.Parameters, "RetentionPeriod")); retentionPeriod > 0 {
		configSet.ArchivingOptions.RetentionPeriod = retentionPeriod
	}

	if err := store.UpdateConfigurationSet(configSet); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
