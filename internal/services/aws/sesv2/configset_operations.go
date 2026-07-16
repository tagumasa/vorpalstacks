package sesv2

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// updateConfigSet is a common helper for PutConfigurationSet* operations.
// It retrieves the configuration set, applies the modifier, and persists.
func (s *SESv2Service) updateConfigSet(reqCtx *request.RequestContext, req *request.ParsedRequest, modify func(*sesv2store.ConfigurationSet, map[string]interface{})) (interface{}, error) {
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

	modify(configSet, req.Parameters)

	if err := store.UpdateConfigurationSet(configSet); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

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
			TlsPolicy:          request.GetStringParam(deliveryOpts, "TlsPolicy"),
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
		if err := store.TagFromSlice(arn, parsedTags); err != nil {
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

	resp := map[string]interface{}{
		"ConfigurationSetName": configSet.Name,
	}

	if configSet.SendingOptions != nil {
		resp["SendingOptions"] = map[string]interface{}{
			"SendingEnabled": configSet.SendingOptions.SendingEnabled,
		}
	}

	if configSet.ReputationOptions != nil {
		resp["ReputationOptions"] = map[string]interface{}{
			"ReputationMetricsEnabled": configSet.ReputationOptions.ReputationMetricsEnabled,
		}
	}

	if configSet.DeliveryOptions != nil {
		resp["DeliveryOptions"] = map[string]interface{}{
			"SendingPoolName":    configSet.DeliveryOptions.SendingPoolName,
			"MaxDeliverySeconds": configSet.DeliveryOptions.MaxDeliverySeconds,
			"TlsPolicy":          configSet.DeliveryOptions.TlsPolicy,
		}
	}

	if configSet.TrackingOptions != nil {
		resp["TrackingOptions"] = map[string]interface{}{
			"CustomRedirectDomain": configSet.TrackingOptions.CustomRedirectDomain,
			"HttpsPolicy":          configSet.TrackingOptions.HttpsPolicy,
		}
	}

	return resp, nil
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

	resp := map[string]interface{}{
		"ConfigurationSets": configSets,
	}

	if result.IsTruncated {
		resp["NextToken"] = result.NextMarker
	}

	return resp, nil
}

// parseEventDestinationDefinition extracts the destination type details (SNS,
// Kinesis Firehose, CloudWatch, Pinpoint, EventBridge) from the request map.
func parseEventDestinationDefinition(params map[string]interface{}) *sesv2store.EventDestinationDefinition {
	def := &sesv2store.EventDestinationDefinition{}
	hasDest := false

	if snsMap := request.GetMapParam(params, "SnsDestination"); snsMap != nil {
		if arn := request.GetStringParam(snsMap, "TopicArn"); arn != "" {
			def.SnsDestination = &sesv2store.SnsDestination{TopicARN: arn}
			hasDest = true
		}
	}
	if kfMap := request.GetMapParam(params, "KinesisFirehoseDestination"); kfMap != nil {
		dsArn := request.GetStringParam(kfMap, "DeliveryStreamArn")
		roleArn := request.GetStringParam(kfMap, "IamRoleArn")
		if dsArn != "" || roleArn != "" {
			def.KinesisFirehoseDestination = &sesv2store.KinesisFirehoseDestination{
				DeliveryStreamARN: dsArn,
				IAMRoleARN:        roleArn,
			}
			hasDest = true
		}
	}
	if cwMap := request.GetMapParam(params, "CloudWatchDestination"); cwMap != nil {
		dims := parseCloudWatchDimensions(cwMap)
		if len(dims) > 0 {
			def.CloudWatchDestination = &sesv2store.CloudWatchDestination{DimensionConfigurations: dims}
			hasDest = true
		}
	}
	if ppMap := request.GetMapParam(params, "PinpointDestination"); ppMap != nil {
		if arn := request.GetStringParam(ppMap, "ApplicationArn"); arn != "" {
			def.PinpointDestination = &sesv2store.PinpointDestination{ApplicationARN: arn}
			hasDest = true
		}
	}
	if ebMap := request.GetMapParam(params, "EventBridgeDestination"); ebMap != nil {
		if arn := request.GetStringParam(ebMap, "EventBusArn"); arn != "" {
			def.EventBridgeDestination = &sesv2store.EventBridgeDestination{EventBusARN: arn}
			hasDest = true
		}
	}

	if !hasDest {
		return nil
	}
	return def
}

// parseCloudWatchDimensions extracts dimension configurations from a CloudWatch
// destination map.
func parseCloudWatchDimensions(params map[string]interface{}) []sesv2store.CloudWatchDimensionConfiguration {
	var result []sesv2store.CloudWatchDimensionConfiguration
	if dimConfigs := request.GetMapParam(params, "DimensionConfigurations"); dimConfigs != nil {
		for i := 1; ; i++ {
			prefix := fmt.Sprintf("member.%d", i)
			dimName := request.GetStringParam(dimConfigs, prefix+".DimensionName")
			if dimName == "" {
				break
			}
			result = append(result, sesv2store.CloudWatchDimensionConfiguration{
				DimensionName:         dimName,
				DimensionValueSource:  request.GetStringParam(dimConfigs, prefix+".DimensionValueSource"),
				DefaultDimensionValue: request.GetStringParam(dimConfigs, prefix+".DefaultDimensionValue"),
			})
		}
	}
	return result
}

// formatEventDestination builds the response map for an event destination,
// including all destination type details.
func formatEventDestination(ed *sesv2store.EventDestination) map[string]interface{} {
	edMap := map[string]interface{}{
		"Name":    ed.Name,
		"Enabled": ed.Enabled,
	}
	if len(ed.MatchingEventTypes) > 0 {
		edMap["MatchingEventTypes"] = ed.MatchingEventTypes
	}
	if ed.EventDestinationDefinition != nil {
		def := ed.EventDestinationDefinition
		if def.SnsDestination != nil {
			edMap["SnsDestination"] = map[string]interface{}{"TopicArn": def.SnsDestination.TopicARN}
		}
		if def.KinesisFirehoseDestination != nil {
			edMap["KinesisFirehoseDestination"] = map[string]interface{}{
				"DeliveryStreamArn": def.KinesisFirehoseDestination.DeliveryStreamARN,
				"IamRoleArn":        def.KinesisFirehoseDestination.IAMRoleARN,
			}
		}
		if def.CloudWatchDestination != nil {
			dims := make([]interface{}, 0, len(def.CloudWatchDestination.DimensionConfigurations))
			for _, dc := range def.CloudWatchDestination.DimensionConfigurations {
				dims = append(dims, map[string]interface{}{
					"DimensionName":         dc.DimensionName,
					"DimensionValueSource":  dc.DimensionValueSource,
					"DefaultDimensionValue": dc.DefaultDimensionValue,
				})
			}
			edMap["CloudWatchDestination"] = map[string]interface{}{"DimensionConfigurations": dims}
		}
		if def.PinpointDestination != nil {
			edMap["PinpointDestination"] = map[string]interface{}{"ApplicationArn": def.PinpointDestination.ApplicationARN}
		}
		if def.EventBridgeDestination != nil {
			edMap["EventBridgeDestination"] = map[string]interface{}{"EventBusArn": def.EventBridgeDestination.EventBusARN}
		}
	}
	return edMap
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

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	for _, ed := range configSet.EventDestinations {
		if ed.Name == eventDestinationName {
			return nil, ErrAlreadyExists
		}
	}

	eventDest := &sesv2store.EventDestination{
		Name:    eventDestinationName,
		Enabled: true,
	}

	if defMap := request.GetMapParam(req.Parameters, "EventDestination"); defMap != nil {
		eventDest.Enabled = request.GetBoolParam(defMap, "Enabled")
		eventDest.MatchingEventTypes = request.GetStringList(defMap, "MatchingEventTypes")
		eventDest.EventDestinationDefinition = parseEventDestinationDefinition(defMap)
	}

	configSet.EventDestinations = append(configSet.EventDestinations, eventDest)
	if err := store.UpdateConfigurationSet(configSet); err != nil {
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

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	eventDests := make([]interface{}, 0, len(configSet.EventDestinations))
	for _, ed := range configSet.EventDestinations {
		edMap := formatEventDestination(ed)
		eventDests = append(eventDests, edMap)
	}

	return map[string]interface{}{
		"EventDestinations": eventDests,
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

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	found := false
	for i, ed := range configSet.EventDestinations {
		if ed.Name == eventDestinationName {
			if defMap := request.GetMapParam(req.Parameters, "EventDestination"); defMap != nil {
				configSet.EventDestinations[i].Enabled = request.GetBoolParam(defMap, "Enabled")
				if types := request.GetStringList(defMap, "MatchingEventTypes"); len(types) > 0 {
					configSet.EventDestinations[i].MatchingEventTypes = types
				}
				if newDef := parseEventDestinationDefinition(defMap); newDef != nil {
					configSet.EventDestinations[i].EventDestinationDefinition = newDef
				}
			}
			found = true
			break
		}
	}
	if !found {
		return nil, ErrNotFound
	}

	if err := store.UpdateConfigurationSet(configSet); err != nil {
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

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return nil, err
	}

	found := false
	filtered := make([]*sesv2store.EventDestination, 0, len(configSet.EventDestinations))
	for _, ed := range configSet.EventDestinations {
		if ed.Name == eventDestinationName {
			found = true
			continue
		}
		filtered = append(filtered, ed)
	}
	if !found {
		return nil, ErrNotFound
	}

	configSet.EventDestinations = filtered
	if err := store.UpdateConfigurationSet(configSet); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutConfigurationSetDeliveryOptions updates the delivery options for a configuration set.
func (s *SESv2Service) PutConfigurationSetDeliveryOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateConfigSet(reqCtx, req, func(cs *sesv2store.ConfigurationSet, params map[string]interface{}) {
		if cs.DeliveryOptions == nil {
			cs.DeliveryOptions = &sesv2store.DeliveryOptions{}
		}
		if v := request.GetStringParam(params, "SendingPoolName"); v != "" {
			cs.DeliveryOptions.SendingPoolName = v
		}
		if v := request.GetIntParam(params, "MaxDeliverySeconds"); v > 0 {
			cs.DeliveryOptions.MaxDeliverySeconds = int32(v)
		}
		if v := request.GetStringParam(params, "TlsPolicy"); v != "" {
			cs.DeliveryOptions.TlsPolicy = v
		}
	})
}

// PutConfigurationSetReputationOptions updates the reputation options for a configuration set.
func (s *SESv2Service) PutConfigurationSetReputationOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateConfigSet(reqCtx, req, func(cs *sesv2store.ConfigurationSet, params map[string]interface{}) {
		if cs.ReputationOptions == nil {
			cs.ReputationOptions = &sesv2store.ReputationOptions{}
		}
		cs.ReputationOptions.ReputationMetricsEnabled = request.GetBoolParam(params, "ReputationMetricsEnabled")
	})
}

// PutConfigurationSetSendingOptions updates the sending options for a configuration set.
func (s *SESv2Service) PutConfigurationSetSendingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateConfigSet(reqCtx, req, func(cs *sesv2store.ConfigurationSet, params map[string]interface{}) {
		if cs.SendingOptions == nil {
			cs.SendingOptions = &sesv2store.SendingOptions{}
		}
		cs.SendingOptions.SendingEnabled = request.GetBoolParam(params, "SendingEnabled")
	})
}

// PutConfigurationSetSuppressionOptions updates the suppression options for a configuration set.
func (s *SESv2Service) PutConfigurationSetSuppressionOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateConfigSet(reqCtx, req, func(cs *sesv2store.ConfigurationSet, params map[string]interface{}) {
		if cs.SuppressionOptions == nil {
			cs.SuppressionOptions = &sesv2store.SuppressionOptions{}
		}
		cs.SuppressionOptions.SuppressedReasons = request.GetStringList(params, "SuppressedReasons")
	})
}

// PutConfigurationSetTrackingOptions updates the tracking options for a configuration set.
func (s *SESv2Service) PutConfigurationSetTrackingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateConfigSet(reqCtx, req, func(cs *sesv2store.ConfigurationSet, params map[string]interface{}) {
		if cs.TrackingOptions == nil {
			cs.TrackingOptions = &sesv2store.TrackingOptions{}
		}
		if v := request.GetStringParam(params, "CustomRedirectDomain"); v != "" {
			cs.TrackingOptions.CustomRedirectDomain = v
		}
		if v := request.GetStringParam(params, "HttpsPolicy"); v != "" {
			cs.TrackingOptions.HttpsPolicy = v
		}
	})
}

// PutConfigurationSetVdmOptions updates the VDM options for a configuration set.
func (s *SESv2Service) PutConfigurationSetVdmOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateConfigSet(reqCtx, req, func(cs *sesv2store.ConfigurationSet, params map[string]interface{}) {
		if cs.VdmOptions == nil {
			cs.VdmOptions = &sesv2store.VdmOptions{}
		}

		if dashboardOpts := request.GetMapParam(params, "DashboardOptions"); dashboardOpts != nil {
			if cs.VdmOptions.DashboardOptions == nil {
				cs.VdmOptions.DashboardOptions = &sesv2store.VDMDashboardOptions{}
			}
			cs.VdmOptions.DashboardOptions.EngagementMetrics = request.GetStringParam(dashboardOpts, "EngagementMetrics")
		}

		if guardianOpts := request.GetMapParam(params, "GuardianOptions"); guardianOpts != nil {
			if cs.VdmOptions.GuardianOptions == nil {
				cs.VdmOptions.GuardianOptions = &sesv2store.VDMGuardianOptions{}
			}
			cs.VdmOptions.GuardianOptions.OptimizedSharedDelivery = request.GetStringParam(guardianOpts, "OptimizedSharedDelivery")
		}
	})
}

// PutConfigurationSetArchivingOptions updates the archiving options for a configuration set.
func (s *SESv2Service) PutConfigurationSetArchivingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateConfigSet(reqCtx, req, func(cs *sesv2store.ConfigurationSet, params map[string]interface{}) {
		if cs.ArchivingOptions == nil {
			cs.ArchivingOptions = &sesv2store.ArchivingOptions{}
		}
		cs.ArchivingOptions.Enabled = request.GetBoolParam(params, "Enabled")
		if v := request.GetStringParam(params, "TargetArn"); v != "" {
			cs.ArchivingOptions.TargetArn = v
		}
		if v := int32(request.GetIntParam(params, "RetentionPeriod")); v > 0 {
			cs.ArchivingOptions.RetentionPeriod = v
		}
	})
}
