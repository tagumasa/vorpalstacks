package sesv2

import (
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// ---------------------------------------------------------------------------
// Input DTOs — configuration-set family
// ---------------------------------------------------------------------------

// CreateConfigurationSetInput carries every create-time configuration-set
// member. The nested option structures travel as raw wire maps and are
// interpreted by the Core.
type CreateConfigurationSetInput struct {
	ConfigurationSetName string
	Tags                 []tags.Tag
	DeliveryOptions      map[string]interface{}
	ReputationOptions    map[string]interface{}
	SendingOptions       map[string]interface{}
	TrackingOptions      map[string]interface{}
	SuppressionOptions   map[string]interface{}
	VdmOptions           map[string]interface{}
	ArchivingOptions     map[string]interface{}
}

// EventDestinationInput carries the raw EventDestination wire map of the
// event-destination CRUD operations.
type EventDestinationInput struct {
	Map      map[string]interface{}
	Provided bool
}

// CreateConfigurationSetEventDestinationInput carries the create members.
type CreateConfigurationSetEventDestinationInput struct {
	ConfigurationSetName string
	EventDestinationName string
	EventDestination     EventDestinationInput
}

// UpdateConfigurationSetEventDestinationInput carries the update members.
type UpdateConfigurationSetEventDestinationInput struct {
	ConfigurationSetName string
	EventDestinationName string
	EventDestination     EventDestinationInput
}

// PutConfigurationSetDeliveryOptionsInput carries the flat
// PutConfigurationSetDeliveryOptions members.
type PutConfigurationSetDeliveryOptionsInput struct {
	ConfigurationSetName string
	SendingPoolName      string
	MaxDeliverySeconds   int
	TlsPolicy            string
}

// PutConfigurationSetReputationOptionsInput carries the flat reputation
// members.
type PutConfigurationSetReputationOptionsInput struct {
	ConfigurationSetName     string
	ReputationMetricsEnabled bool
}

// PutConfigurationSetSendingOptionsInput carries the flat sending members.
type PutConfigurationSetSendingOptionsInput struct {
	ConfigurationSetName string
	SendingEnabled       bool
}

// PutConfigurationSetSuppressionOptionsInput carries the flat suppression
// members plus the raw ValidationOptions wire map.
type PutConfigurationSetSuppressionOptionsInput struct {
	ConfigurationSetName string
	SuppressedReasons    []string
	SuppressionScope     string
	ValidationOptions    map[string]interface{}
}

// PutConfigurationSetTrackingOptionsInput carries the flat tracking members.
type PutConfigurationSetTrackingOptionsInput struct {
	ConfigurationSetName string
	CustomRedirectDomain string
	HttpsPolicy          string
}

// PutConfigurationSetVdmOptionsInput carries the raw VdmOptions wire map.
type PutConfigurationSetVdmOptionsInput struct {
	ConfigurationSetName string
	VdmOptions           map[string]interface{}
}

// PutConfigurationSetArchivingOptionsInput carries the flat archiving
// members.
type PutConfigurationSetArchivingOptionsInput struct {
	ConfigurationSetName string
	ArchiveArn           string
}

// ---------------------------------------------------------------------------
// Shared Core helpers
// ---------------------------------------------------------------------------

// updateConfigSetCore is the shared skeleton for the
// PutConfigurationSet* family: it resolves the configuration set, applies
// the modifier, and persists.
func (s *SESv2Service) updateConfigSetCore(store sesv2store.SESv2StoreInterface, configSetName string, modify func(*sesv2store.ConfigurationSet) error) error {
	if configSetName == "" {
		return ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return err
	}

	if err := modify(configSet); err != nil {
		return err
	}

	return store.UpdateConfigurationSet(configSet)
}

// parseEventDestinationDefinition extracts the destination-type details
// (SNS, Kinesis Firehose, CloudWatch, Pinpoint, EventBridge) from the
// request map. It returns the parsed definition, a boolean indicating
// whether at least one destination type was supplied, and any validation
// error.
//
// Enabled and MatchingEventTypes are intentionally NOT parsed here —
// the caller manages those top-level fields directly from the
// EventDestination map, avoiding the redundant double-assignment that
// previously existed.
func parseEventDestinationDefinition(params map[string]interface{}) (*sesv2store.EventDestinationDefinition, bool, error) {
	def := &sesv2store.EventDestinationDefinition{}

	hasSns := false
	hasKinesis := false
	hasCloudWatch := false
	hasPinpoint := false
	hasEventBridge := false

	if snsMap := request.GetMapParam(params, "SnsDestination"); snsMap != nil {
		if arn := request.GetStringParam(snsMap, "TopicArn"); arn != "" {
			def.SnsDestination = &sesv2store.SnsDestination{TopicARN: arn}
			hasSns = true
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
			hasKinesis = true
		}
	}
	if cwMap := request.GetMapParam(params, "CloudWatchDestination"); cwMap != nil {
		dims := parseCloudWatchDimensions(cwMap)
		if len(dims) > 0 {
			def.CloudWatchDestination = &sesv2store.CloudWatchDestination{DimensionConfigurations: dims}
			hasCloudWatch = true
		}
	}
	if ppMap := request.GetMapParam(params, "PinpointDestination"); ppMap != nil {
		if arn := request.GetStringParam(ppMap, "ApplicationArn"); arn != "" {
			def.PinpointDestination = &sesv2store.PinpointDestination{ApplicationARN: arn}
			hasPinpoint = true
		}
	}
	if ebMap := request.GetMapParam(params, "EventBridgeDestination"); ebMap != nil {
		if arn := request.GetStringParam(ebMap, "EventBusArn"); arn != "" {
			def.EventBridgeDestination = &sesv2store.EventBridgeDestination{EventBusARN: arn}
			hasEventBridge = true
		}
	}

	// Per AWS docs, only one destination type may be specified per event
	// destination. Zero destinations is permitted — the EventDestination
	// can be created with only Enabled and MatchingEventTypes.
	destCount := countEventDestinations(hasSns, hasKinesis, hasCloudWatch, hasPinpoint, hasEventBridge)
	if destCount > 1 {
		return nil, false, ErrBadRequest
	}

	return def, destCount > 0, nil
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

// ---------------------------------------------------------------------------
// Core functions — configuration-set family
// ---------------------------------------------------------------------------

// createConfigurationSetCore is the single entry point for
// configuration-set creation, including every nested option structure and
// create-time tags.
func (s *SESv2Service) createConfigurationSetCore(store sesv2store.SESv2StoreInterface, in CreateConfigurationSetInput) (map[string]interface{}, error) {
	configSetName := in.ConfigurationSetName
	if configSetName == "" {
		return nil, ErrMissingParameter
	}
	if !validateConfigurationSetName(configSetName) {
		return nil, ErrBadRequest
	}

	configSet := sesv2store.NewConfigurationSet(configSetName)

	if in.DeliveryOptions != nil {
		tls := request.GetStringParam(in.DeliveryOptions, "TlsPolicy")
		if tls != "" && !validateTlsPolicy(tls) {
			return nil, ErrBadRequest
		}
		mds := int32(request.GetIntParam(in.DeliveryOptions, "MaxDeliverySeconds"))
		if mds > 0 && !validateMaxDeliverySeconds(mds) {
			return nil, ErrBadRequest
		}
		configSet.DeliveryOptions = &sesv2store.DeliveryOptions{
			SendingPoolName:    request.GetStringParam(in.DeliveryOptions, "SendingPoolName"),
			MaxDeliverySeconds: mds,
			TlsPolicy:          tls,
		}
	}

	if in.ReputationOptions != nil {
		configSet.ReputationOptions = &sesv2store.ReputationOptions{
			ReputationMetricsEnabled: request.GetBoolParam(in.ReputationOptions, "ReputationMetricsEnabled"),
		}
	}

	if in.SendingOptions != nil {
		configSet.SendingOptions = &sesv2store.SendingOptions{
			SendingEnabled: request.GetBoolParam(in.SendingOptions, "SendingEnabled"),
		}
	}

	if in.TrackingOptions != nil {
		https := request.GetStringParam(in.TrackingOptions, "HttpsPolicy")
		if https != "" && !validateHttpsPolicy(https) {
			return nil, ErrBadRequest
		}
		configSet.TrackingOptions = &sesv2store.TrackingOptions{
			CustomRedirectDomain: request.GetStringParam(in.TrackingOptions, "CustomRedirectDomain"),
			HttpsPolicy:          https,
		}
	}

	if in.SuppressionOptions != nil {
		for _, r := range request.GetStringList(in.SuppressionOptions, "SuppressedReasons") {
			if !validateSuppressionListReason(r) {
				return nil, ErrBadRequest
			}
		}
		scope := request.GetStringParam(in.SuppressionOptions, "SuppressionScope")
		if scope != "" && !validateSuppressionListScope(scope) {
			return nil, ErrBadRequest
		}
		configSet.SuppressionOptions = &sesv2store.SuppressionOptions{
			SuppressedReasons: request.GetStringList(in.SuppressionOptions, "SuppressedReasons"),
			SuppressionScope:  scope,
		}
	}

	if in.VdmOptions != nil {
		configSet.VdmOptions = &sesv2store.VdmOptions{}
		if dashboardOpts := request.GetMapParam(in.VdmOptions, "DashboardOptions"); dashboardOpts != nil {
			em := request.GetStringParam(dashboardOpts, "EngagementMetrics")
			if em != "" && !validateFeatureStatus(em) {
				return nil, ErrBadRequest
			}
			configSet.VdmOptions.DashboardOptions = &sesv2store.VDMDashboardOptions{
				EngagementMetrics: em,
			}
		}
		if guardianOpts := request.GetMapParam(in.VdmOptions, "GuardianOptions"); guardianOpts != nil {
			osd := request.GetStringParam(guardianOpts, "OptimizedSharedDelivery")
			if osd != "" && !validateFeatureStatus(osd) {
				return nil, ErrBadRequest
			}
			configSet.VdmOptions.GuardianOptions = &sesv2store.VDMGuardianOptions{
				OptimizedSharedDelivery: osd,
			}
		}
	}

	if in.ArchivingOptions != nil {
		if archiveArn := request.GetStringParam(in.ArchivingOptions, "ArchiveArn"); archiveArn != "" {
			configSet.ArchivingOptions = &sesv2store.ArchivingOptions{
				ArchiveArn: archiveArn,
			}
		}
	}

	if _, err := store.CreateConfigurationSet(configSet); err != nil {
		return nil, err
	}

	if len(in.Tags) > 0 {
		arn := store.BuildConfigSetArn(configSetName)
		if err := store.TagFromSlice(arn, in.Tags); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"ConfigurationSetName": configSetName,
	}, nil
}

// getConfigurationSetCore is the single entry point for reading a
// configuration set, including its tags.
func (s *SESv2Service) getConfigurationSetCore(store sesv2store.SESv2StoreInterface, configSetName string) (map[string]interface{}, error) {
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

	if configSet.SuppressionOptions != nil {
		resp["SuppressionOptions"] = map[string]interface{}{
			"SuppressedReasons": configSet.SuppressionOptions.SuppressedReasons,
		}
	}

	if configSet.VdmOptions != nil {
		vdmMap := map[string]interface{}{}
		if configSet.VdmOptions.DashboardOptions != nil {
			vdmMap["DashboardOptions"] = map[string]interface{}{
				"EngagementMetrics": configSet.VdmOptions.DashboardOptions.EngagementMetrics,
			}
		}
		if configSet.VdmOptions.GuardianOptions != nil {
			vdmMap["GuardianOptions"] = map[string]interface{}{
				"OptimizedSharedDelivery": configSet.VdmOptions.GuardianOptions.OptimizedSharedDelivery,
			}
		}
		resp["VdmOptions"] = vdmMap
	}

	if configSet.ArchivingOptions != nil {
		resp["ArchivingOptions"] = map[string]interface{}{
			"ArchiveArn": configSet.ArchivingOptions.ArchiveArn,
		}
	}

	arn := store.BuildConfigSetArn(configSet.Name)
	if tags, err := store.ListAsSlice(arn); err == nil && len(tags) > 0 {
		resp["Tags"] = tags
	}

	return resp, nil
}

// deleteConfigurationSetCore is the single entry point for deleting a
// configuration set.
func (s *SESv2Service) deleteConfigurationSetCore(store sesv2store.SESv2StoreInterface, configSetName string) error {
	if configSetName == "" {
		return ErrMissingParameter
	}
	return store.DeleteConfigurationSet(configSetName)
}

// listConfigurationSetsCore is the single entry point for listing
// configuration sets.
func (s *SESv2Service) listConfigurationSetsCore(store sesv2store.SESv2StoreInterface, pageSize int, nextToken string) (map[string]interface{}, error) {
	if pageSize == 0 {
		pageSize = 100
	}

	result, err := store.ListConfigurationSets(common.ListOptions{
		MaxItems: pageSize,
		Marker:   nextToken,
	})
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

// createConfigurationSetEventDestinationCore is the single entry point for
// creating an event destination. Per AWS, when Enabled=true the request must
// supply a non-empty MatchingEventTypes list. We replicate that validation
// rather than silently persisting a destination that would never fire.
func (s *SESv2Service) createConfigurationSetEventDestinationCore(store sesv2store.SESv2StoreInterface, in CreateConfigurationSetEventDestinationInput) error {
	configSetName := in.ConfigurationSetName
	eventDestinationName := in.EventDestinationName

	if configSetName == "" || eventDestinationName == "" {
		return ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return err
	}

	for _, ed := range configSet.EventDestinations {
		if ed.Name == eventDestinationName {
			return ErrAlreadyExists
		}
	}

	eventDest := &sesv2store.EventDestination{
		Name: eventDestinationName,
	}

	if defMap := in.EventDestination.Map; in.EventDestination.Provided && defMap != nil {
		// GetBoolParam returns false for an absent key, which
		// would clobber the intended default. Use GetBoolParamDefault so
		// the AWS default (true) is preserved when the caller omits
		// Enabled.
		eventDest.Enabled = request.GetBoolParamDefault(defMap, "Enabled", true)

		// Parse MatchingEventTypes directly here instead of
		// inside parseEventDestinationDefinition, eliminating the
		// redundant double-assignment.
		if types := request.GetStringList(defMap, "MatchingEventTypes"); len(types) > 0 {
			if !validateEventTypes(types) {
				return ErrBadRequest
			}
			eventDest.MatchingEventTypes = types
		}

		def, _, err := parseEventDestinationDefinition(defMap)
		if err != nil {
			return err
		}
		eventDest.EventDestinationDefinition = def
	} else {
		eventDest.Enabled = true
	}

	if eventDest.Enabled && len(eventDest.MatchingEventTypes) == 0 {
		return ErrInvalidParameter
	}

	configSet.EventDestinations = append(configSet.EventDestinations, eventDest)
	return store.UpdateConfigurationSet(configSet)
}

// getConfigurationSetEventDestinationsCore is the single entry point for
// listing a configuration set's event destinations.
func (s *SESv2Service) getConfigurationSetEventDestinationsCore(store sesv2store.SESv2StoreInterface, configSetName string) (map[string]interface{}, error) {
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

// updateConfigurationSetEventDestinationCore is the single entry point for
// updating an event destination.
func (s *SESv2Service) updateConfigurationSetEventDestinationCore(store sesv2store.SESv2StoreInterface, in UpdateConfigurationSetEventDestinationInput) error {
	configSetName := in.ConfigurationSetName
	eventDestinationName := in.EventDestinationName

	if configSetName == "" || eventDestinationName == "" {
		return ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return err
	}

	found := false
	for i, ed := range configSet.EventDestinations {
		if ed.Name == eventDestinationName {
			if defMap := in.EventDestination.Map; in.EventDestination.Provided && defMap != nil {
				// Only update Enabled when the key is actually
				// present. Previously GetBoolParam returned false for an
				// absent key, silently disabling an active destination.
				if _, ok := defMap["Enabled"]; ok {
					configSet.EventDestinations[i].Enabled = request.GetBoolParam(defMap, "Enabled")
				}

				// Parse MatchingEventTypes directly.
				if types := request.GetStringList(defMap, "MatchingEventTypes"); len(types) > 0 {
					if !validateEventTypes(types) {
						return ErrBadRequest
					}
					configSet.EventDestinations[i].MatchingEventTypes = types
				}

				def, hasDestType, err := parseEventDestinationDefinition(defMap)
				if err != nil {
					return err
				}
				// Only overwrite the destination definition
				// when the caller supplied a destination type. Without
				// this guard, a partial update (Enabled-only or
				// MatchingEventTypes-only) would replace the stored SNS
				// topic / Kinesis stream / etc. with nil pointers.
				if hasDestType {
					configSet.EventDestinations[i].EventDestinationDefinition = def
				}
			}
			if configSet.EventDestinations[i].Enabled && len(configSet.EventDestinations[i].MatchingEventTypes) == 0 {
				return ErrInvalidParameter
			}
			found = true
			break
		}
	}
	if !found {
		return ErrNotFound
	}

	return store.UpdateConfigurationSet(configSet)
}

// deleteConfigurationSetEventDestinationCore is the single entry point for
// deleting an event destination.
func (s *SESv2Service) deleteConfigurationSetEventDestinationCore(store sesv2store.SESv2StoreInterface, configSetName, eventDestinationName string) error {
	if configSetName == "" || eventDestinationName == "" {
		return ErrMissingParameter
	}

	configSet, err := store.GetConfigurationSet(configSetName)
	if err != nil {
		return err
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
		return ErrNotFound
	}

	configSet.EventDestinations = filtered
	return store.UpdateConfigurationSet(configSet)
}

// ---------------------------------------------------------------------------
// Core functions — PutConfigurationSet* family
// ---------------------------------------------------------------------------

// putConfigurationSetDeliveryOptionsCore is the single entry point for the
// delivery-options update.
func (s *SESv2Service) putConfigurationSetDeliveryOptionsCore(store sesv2store.SESv2StoreInterface, in PutConfigurationSetDeliveryOptionsInput) error {
	return s.updateConfigSetCore(store, in.ConfigurationSetName, func(cs *sesv2store.ConfigurationSet) error {
		if cs.DeliveryOptions == nil {
			cs.DeliveryOptions = &sesv2store.DeliveryOptions{}
		}
		if in.SendingPoolName != "" {
			cs.DeliveryOptions.SendingPoolName = in.SendingPoolName
		}
		if in.MaxDeliverySeconds > 0 {
			if !validateMaxDeliverySeconds(int32(in.MaxDeliverySeconds)) {
				return ErrBadRequest
			}
			cs.DeliveryOptions.MaxDeliverySeconds = int32(in.MaxDeliverySeconds)
		}
		if in.TlsPolicy != "" {
			if !validateTlsPolicy(in.TlsPolicy) {
				return ErrBadRequest
			}
			cs.DeliveryOptions.TlsPolicy = in.TlsPolicy
		}
		return nil
	})
}

// putConfigurationSetReputationOptionsCore is the single entry point for the
// reputation-options update.
func (s *SESv2Service) putConfigurationSetReputationOptionsCore(store sesv2store.SESv2StoreInterface, in PutConfigurationSetReputationOptionsInput) error {
	return s.updateConfigSetCore(store, in.ConfigurationSetName, func(cs *sesv2store.ConfigurationSet) error {
		if cs.ReputationOptions == nil {
			cs.ReputationOptions = &sesv2store.ReputationOptions{}
		}
		cs.ReputationOptions.ReputationMetricsEnabled = in.ReputationMetricsEnabled
		return nil
	})
}

// putConfigurationSetSendingOptionsCore is the single entry point for the
// sending-options update.
func (s *SESv2Service) putConfigurationSetSendingOptionsCore(store sesv2store.SESv2StoreInterface, in PutConfigurationSetSendingOptionsInput) error {
	return s.updateConfigSetCore(store, in.ConfigurationSetName, func(cs *sesv2store.ConfigurationSet) error {
		if cs.SendingOptions == nil {
			cs.SendingOptions = &sesv2store.SendingOptions{}
		}
		cs.SendingOptions.SendingEnabled = in.SendingEnabled
		return nil
	})
}

// putConfigurationSetSuppressionOptionsCore is the single entry point for
// the suppression-options update. Per Smithy
// com.amazonaws.sesv2#PutConfigurationSetSuppressionOptionsRequest the input
// carries SuppressedReasons, SuppressionScope (ACCOUNT/TENANT), and
// ValidationOptions (Auto Validation threshold settings).
func (s *SESv2Service) putConfigurationSetSuppressionOptionsCore(store sesv2store.SESv2StoreInterface, in PutConfigurationSetSuppressionOptionsInput) error {
	return s.updateConfigSetCore(store, in.ConfigurationSetName, func(cs *sesv2store.ConfigurationSet) error {
		if cs.SuppressionOptions == nil {
			cs.SuppressionOptions = &sesv2store.SuppressionOptions{}
		}
		for _, r := range in.SuppressedReasons {
			if !validateSuppressionListReason(r) {
				return ErrBadRequest
			}
		}
		cs.SuppressionOptions.SuppressedReasons = in.SuppressedReasons
		if in.SuppressionScope != "" {
			if !validateSuppressionListScope(in.SuppressionScope) {
				return ErrBadRequest
			}
			cs.SuppressionOptions.SuppressionScope = in.SuppressionScope
		}

		if vo := in.ValidationOptions; vo != nil {
			if ct := request.GetMapParam(vo, "ConditionThreshold"); ct != nil {
				thresholdEnabled := request.GetStringParam(ct, "ConditionThresholdEnabled")
				if thresholdEnabled != "" && !validateFeatureStatus(thresholdEnabled) {
					return ErrBadRequest
				}
				threshold := &sesv2store.SuppressionConditionThreshold{
					ConditionThresholdEnabled: thresholdEnabled,
				}
				if oct := request.GetMapParam(ct, "OverallConfidenceThreshold"); oct != nil {
					cvt := request.GetStringParam(oct, "ConfidenceVerdictThreshold")
					if cvt != "" && !validateSuppressionConfidenceVerdictThreshold(cvt) {
						return ErrBadRequest
					}
					threshold.OverallConfidenceThreshold = &sesv2store.SuppressionConfidenceThreshold{
						ConfidenceVerdictThreshold: cvt,
					}
				}
				cs.SuppressionOptions.ValidationOptions = &sesv2store.SuppressionValidationOptions{
					ConditionThreshold: threshold,
				}
			}
		}
		return nil
	})
}

// putConfigurationSetTrackingOptionsCore is the single entry point for the
// tracking-options update.
func (s *SESv2Service) putConfigurationSetTrackingOptionsCore(store sesv2store.SESv2StoreInterface, in PutConfigurationSetTrackingOptionsInput) error {
	return s.updateConfigSetCore(store, in.ConfigurationSetName, func(cs *sesv2store.ConfigurationSet) error {
		if cs.TrackingOptions == nil {
			cs.TrackingOptions = &sesv2store.TrackingOptions{}
		}
		if in.CustomRedirectDomain != "" {
			cs.TrackingOptions.CustomRedirectDomain = in.CustomRedirectDomain
		}
		if in.HttpsPolicy != "" {
			if !validateHttpsPolicy(in.HttpsPolicy) {
				return ErrBadRequest
			}
			cs.TrackingOptions.HttpsPolicy = in.HttpsPolicy
		}
		return nil
	})
}

// putConfigurationSetVdmOptionsCore is the single entry point for the VDM
// options update. Per Smithy
// com.amazonaws.sesv2#PutConfigurationSetVdmOptionsRequest the VdmOptions
// member carries DashboardOptions and GuardianOptions as nested structs.
func (s *SESv2Service) putConfigurationSetVdmOptionsCore(store sesv2store.SESv2StoreInterface, in PutConfigurationSetVdmOptionsInput) error {
	return s.updateConfigSetCore(store, in.ConfigurationSetName, func(cs *sesv2store.ConfigurationSet) error {
		if cs.VdmOptions == nil {
			cs.VdmOptions = &sesv2store.VdmOptions{}
		}

		vdmOpts := in.VdmOptions
		if vdmOpts == nil {
			return nil
		}

		if dashboardOpts := request.GetMapParam(vdmOpts, "DashboardOptions"); dashboardOpts != nil {
			if cs.VdmOptions.DashboardOptions == nil {
				cs.VdmOptions.DashboardOptions = &sesv2store.VDMDashboardOptions{}
			}
			em := request.GetStringParam(dashboardOpts, "EngagementMetrics")
			if em != "" && !validateFeatureStatus(em) {
				return ErrBadRequest
			}
			cs.VdmOptions.DashboardOptions.EngagementMetrics = em
		}

		if guardianOpts := request.GetMapParam(vdmOpts, "GuardianOptions"); guardianOpts != nil {
			if cs.VdmOptions.GuardianOptions == nil {
				cs.VdmOptions.GuardianOptions = &sesv2store.VDMGuardianOptions{}
			}
			osd := request.GetStringParam(guardianOpts, "OptimizedSharedDelivery")
			if osd != "" && !validateFeatureStatus(osd) {
				return ErrBadRequest
			}
			cs.VdmOptions.GuardianOptions.OptimizedSharedDelivery = osd
		}
		return nil
	})
}

// putConfigurationSetArchivingOptionsCore is the single entry point for the
// archiving-options update. Per Smithy
// `PutConfigurationSetArchivingOptionsRequest`, only `ArchiveArn` is
// accepted.
func (s *SESv2Service) putConfigurationSetArchivingOptionsCore(store sesv2store.SESv2StoreInterface, in PutConfigurationSetArchivingOptionsInput) error {
	return s.updateConfigSetCore(store, in.ConfigurationSetName, func(cs *sesv2store.ConfigurationSet) error {
		if in.ArchiveArn == "" {
			cs.ArchivingOptions = nil
			return nil
		}
		cs.ArchivingOptions = &sesv2store.ArchivingOptions{
			ArchiveArn: in.ArchiveArn,
		}
		return nil
	})
}

// ---------------------------------------------------------------------------
// HTTP handlers — parse → DTO → Core → serialise
// ---------------------------------------------------------------------------
