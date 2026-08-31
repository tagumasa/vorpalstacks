package timestreamquery

import (
	"context"

	"vorpalstacks/internal/common/request"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// DescribeAccountSettings returns the account settings for Timestream Query.
func (s *TimestreamQueryService) DescribeAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	settings, err := s.describeAccountSettingsCore(store)
	if err != nil {
		return nil, err
	}

	return formatAccountSettingsResponse(settings), nil
}

// UpdateAccountSettings updates the account settings for Timestream Query.
func (s *TimestreamQueryService) UpdateAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Use existence check instead of value check. MaxQueryTCU has no
	// Smithy range trait — 0 is a valid value meaning "use default".
	var maxQueryTCU *int64
	if _, ok := req.Parameters["MaxQueryTCU"]; ok {
		val := request.GetIntParam(req.Parameters, "MaxQueryTCU")
		tcu := int64(val)
		maxQueryTCU = &tcu
	}

	input := UpdateAccountSettingsInput{
		MaxQueryTCU:       maxQueryTCU,
		QueryPricingModel: request.GetParamCaseInsensitive(req.Parameters, "QueryPricingModel"),
	}

	if qcMap := request.GetMapParamCaseInsensitive(req.Parameters, "QueryCompute"); qcMap != nil {
		if mode, ok := qcMap["ComputeMode"].(string); ok {
			input.HasComputeMode = true
			input.ComputeMode = mode
		}

		// Parse ProvisionedCapacityRequest (Smithy: TargetQueryTCU +
		// NotificationConfiguration).
		if pcMap, ok := qcMap["ProvisionedCapacity"].(map[string]interface{}); ok {
			input.ProvisionedCapacity = parseProvisionedCapacityRequest(pcMap)
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	settings, err := s.updateAccountSettingsCore(store, input)
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

	// Include ProvisionedCapacity in output when configured.
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
