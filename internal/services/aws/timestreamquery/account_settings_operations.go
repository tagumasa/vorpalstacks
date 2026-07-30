package timestreamquery

import (
	"context"

	"vorpalstacks/internal/common/request"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// validQueryPricingModels enumerates the Smithy QueryPricingModel enum.
var validQueryPricingModels = map[string]bool{
	"BYTES_SCANNED": true,
	"COMPUTE_UNITS": true,
}

// validComputeModes enumerates the Smithy ComputeMode enum.
var validComputeModes = map[string]bool{
	"ON_DEMAND":   true,
	"PROVISIONED": true,
}

// DescribeAccountSettings returns the account settings for Timestream Query.
func (s *TimestreamQueryService) DescribeAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	settings, err := store.accountSettingsStore.GetAccountSettings()
	if err != nil {
		return nil, err
	}

	return formatAccountSettingsResponse(settings), nil
}

// UpdateAccountSettings updates the account settings for Timestream Query.
func (s *TimestreamQueryService) UpdateAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// M2: Use existence check instead of value check. MaxQueryTCU has no
	// Smithy range trait — 0 is a valid value meaning "use default".
	var maxQueryTCU *int64
	if _, ok := req.Parameters["MaxQueryTCU"]; ok {
		val := request.GetIntParam(req.Parameters, "MaxQueryTCU")
		tcu := int64(val)
		maxQueryTCU = &tcu
	}

	queryPricingModel := request.GetParamCaseInsensitive(req.Parameters, "QueryPricingModel")

	// M13: Validate QueryPricingModel enum (Smithy: BYTES_SCANNED, COMPUTE_UNITS).
	if queryPricingModel != "" && !validQueryPricingModels[queryPricingModel] {
		return nil, ErrValidationException
	}

	queryComputeType := ""
	var provisionedCapacity *tsstore.ProvisionedCapacitySettings

	if qcMap := request.GetMapParamCaseInsensitive(req.Parameters, "QueryCompute"); qcMap != nil {
		if mode, ok := qcMap["ComputeMode"].(string); ok {
			// M13: Validate ComputeMode enum (Smithy: ON_DEMAND, PROVISIONED).
			if mode != "" && !validComputeModes[mode] {
				return nil, ErrValidationException
			}
			queryComputeType = mode
		}

		// M2/M3: Parse ProvisionedCapacityRequest (Smithy: TargetQueryTCU +
		// NotificationConfiguration).
		if pcMap, ok := qcMap["ProvisionedCapacity"].(map[string]interface{}); ok {
			provisionedCapacity = parseProvisionedCapacityRequest(pcMap)
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	settings, err := store.accountSettingsStore.UpdateAccountSettings(maxQueryTCU, queryPricingModel, queryComputeType, "", provisionedCapacity)
	if err != nil {
		return nil, err
	}

	return formatAccountSettingsResponse(settings), nil
}

// parseProvisionedCapacityRequest parses a ProvisionedCapacityRequest from
// the raw input map (Smithy: TargetQueryTCU + NotificationConfiguration).
func parseProvisionedCapacityRequest(raw map[string]interface{}) *tsstore.ProvisionedCapacitySettings {
	pc := &tsstore.ProvisionedCapacitySettings{}

	if tcu, ok := raw["TargetQueryTCU"].(float64); ok {
		pc.TargetQueryTCU = int64(tcu)
	}

	if notifConfig, ok := raw["NotificationConfiguration"].(map[string]interface{}); ok {
		if snsConfig, ok := notifConfig["SnsConfiguration"].(map[string]interface{}); ok {
			pc.SNSTopicARN, _ = snsConfig["TopicArn"].(string)
		}
		pc.RoleARN, _ = notifConfig["RoleArn"].(string)
	}

	pc.LastUpdateStatus = "PENDING"

	return pc
}

// formatAccountSettingsResponse builds the API response map for account
// settings, including the full QueryCompute structure with ProvisionedCapacity.
func formatAccountSettingsResponse(settings *tsstore.AccountSettings) map[string]interface{} {
	queryCompute := map[string]interface{}{
		"ComputeMode": settings.QueryComputeType,
	}

	// M3: Include ProvisionedCapacity in output when configured.
	if settings.ProvisionedCapacity != nil {
		pcMap := map[string]interface{}{
			"ActiveQueryTCU": settings.ProvisionedCapacity.ActiveQueryTCU,
		}
		if settings.ProvisionedCapacity.TargetQueryTCU > 0 {
			pcMap["LastUpdate"] = map[string]interface{}{
				"TargetQueryTCU": settings.ProvisionedCapacity.TargetQueryTCU,
				"Status":         settings.ProvisionedCapacity.LastUpdateStatus,
			}
			if settings.ProvisionedCapacity.LastUpdateStatusMessage != "" {
				pcMap["LastUpdate"].(map[string]interface{})["StatusMessage"] = settings.ProvisionedCapacity.LastUpdateStatusMessage
			}
		}
		if settings.ProvisionedCapacity.SNSTopicARN != "" || settings.ProvisionedCapacity.RoleARN != "" {
			notifMap := map[string]interface{}{}
			if settings.ProvisionedCapacity.SNSTopicARN != "" {
				notifMap["SnsConfiguration"] = map[string]interface{}{
					"TopicArn": settings.ProvisionedCapacity.SNSTopicARN,
				}
			}
			if settings.ProvisionedCapacity.RoleARN != "" {
				notifMap["RoleArn"] = settings.ProvisionedCapacity.RoleARN
			}
			pcMap["NotificationConfiguration"] = notifMap
		}
		queryCompute["ProvisionedCapacity"] = pcMap
	}

	return map[string]interface{}{
		"MaxQueryTCU":       settings.MaxQueryTCU,
		"QueryPricingModel": settings.QueryPricingMode,
		"QueryCompute":      queryCompute,
	}
}
