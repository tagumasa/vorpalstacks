package wafv2

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

func validateLogScope(scope string) error {
	switch scope {
	case "", "CUSTOMER", "SECURITY_LAKE", "CLOUDWATCH_TELEMETRY_RULE_MANAGED":
		return nil
	default:
		return invalidParamError("LogScope must be one of: CUSTOMER, SECURITY_LAKE, CLOUDWATCH_TELEMETRY_RULE_MANAGED")
	}
}

func validateLogType(logType string) error {
	switch logType {
	case "", "WAF_LOGS":
		return nil
	default:
		return invalidParamError("LogType must be WAF_LOGS")
	}
}

func parseLoggingFilter(raw interface{}) (*wafstore.LoggingFilter, error) {
	if raw == nil {
		return nil, nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, invalidParamError("LoggingFilter must be an object")
	}
	lf := &wafstore.LoggingFilter{
		DefaultBehavior: request.GetStringParam(m, "DefaultBehavior"),
	}
	if lf.DefaultBehavior != "" {
		if err := validateFilterBehavior(lf.DefaultBehavior); err != nil {
			return nil, err
		}
	}
	if filtersRaw, ok := m["Filters"]; ok {
		if arr, ok := filtersRaw.([]interface{}); ok {
			for _, fRaw := range arr {
				fMap, ok := fRaw.(map[string]interface{})
				if !ok {
					return nil, invalidParamError("LoggingFilter Filters entries must be objects")
				}
				f := wafstore.Filter{
					Behavior:    request.GetStringParam(fMap, "Behavior"),
					Requirement: request.GetStringParam(fMap, "Requirement"),
				}
				if f.Behavior != "" {
					if err := validateFilterBehavior(f.Behavior); err != nil {
						return nil, err
					}
				}
				if f.Requirement != "" {
					if err := validateFilterRequirement(f.Requirement); err != nil {
						return nil, err
					}
				}
				if condsRaw, ok := fMap["Conditions"]; ok {
					if condArr, ok := condsRaw.([]interface{}); ok {
						for _, cRaw := range condArr {
							cMap, ok := cRaw.(map[string]interface{})
							if !ok {
								return nil, invalidParamError("LoggingFilter Conditions entries must be objects")
							}
							fc := wafstore.FilterCondition{}
							if acRaw, ok := cMap["ActionCondition"]; ok {
								if acMap, ok := acRaw.(map[string]interface{}); ok {
									fc.ActionCondition = &wafstore.ActionCondition{
										Action: request.GetStringParam(acMap, "Action"),
									}
								}
							}
							if lncRaw, ok := cMap["LabelNameCondition"]; ok {
								if lncMap, ok := lncRaw.(map[string]interface{}); ok {
									fc.LabelNameCondition = &wafstore.LabelNameCondition{
										LabelName: request.GetStringParam(lncMap, "LabelName"),
									}
								}
							}
							f.Conditions = append(f.Conditions, fc)
						}
					}
				}
				lf.Filters = append(lf.Filters, f)
			}
		}
	}
	if lf.DefaultBehavior == "" && len(lf.Filters) == 0 {
		return nil, nil
	}
	return lf, nil
}

// PutLoggingConfiguration creates or updates the logging configuration for the specified web ACL.
func (s *WAFv2Service) PutLoggingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	loggingConfigMap := request.GetMapParam(req.Parameters, "LoggingConfiguration")
	if loggingConfigMap == nil {
		return nil, invalidParamError("LoggingConfiguration is required")
	}

	var redactedFields []interface{}
	if rfRaw := loggingConfigMap["RedactedFields"]; rfRaw != nil {
		if arr, ok := rfRaw.([]interface{}); ok {
			redactedFields = arr
		}
	}

	config, err := s.putLoggingConfigurationCore(reqCtx, LoggingConfigInput{
		ResourceArn:              request.GetStringParam(loggingConfigMap, "ResourceArn"),
		LogDestinationConfigs:    request.GetStringList(loggingConfigMap, "LogDestinationConfigs"),
		LogScope:                 request.GetStringParam(loggingConfigMap, "LogScope"),
		LogType:                  request.GetStringParam(loggingConfigMap, "LogType"),
		ManagedByFirewallManager: request.GetBoolParam(loggingConfigMap, "ManagedByFirewallManager"),
		RedactedFields:           redactedFields,
		LoggingFilterRaw:         loggingConfigMap["LoggingFilter"],
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"LoggingConfiguration": convertLoggingConfigToResponse(config),
	}, nil
}

// GetLoggingConfiguration retrieves the logging configuration for the specified web ACL.
func (s *WAFv2Service) GetLoggingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	config, err := s.getLoggingConfigurationCore(reqCtx, request.GetStringParam(req.Parameters, "ResourceArn"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"LoggingConfiguration": convertLoggingConfigToResponse(config),
	}, nil
}

// DeleteLoggingConfiguration removes the logging configuration for the specified web ACL.
func (s *WAFv2Service) DeleteLoggingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteLoggingConfigurationCore(reqCtx, request.GetStringParam(req.Parameters, "ResourceArn")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListLoggingConfigurations returns a paginated list of all logging configurations.
func (s *WAFv2Service) ListLoggingConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.listLoggingConfigurationsCore(reqCtx, LoggingListInput{
		Scope:      request.GetStringParam(req.Parameters, "Scope"),
		Limit:      pagination.GetMaxItems(req.Parameters, 100, "Limit"),
		NextMarker: pagination.GetMarker(req.Parameters, "NextMarker"),
	})
	if err != nil {
		return nil, err
	}

	configs := make([]interface{}, 0, len(result.LoggingConfigurations))
	for _, config := range result.LoggingConfigurations {
		configs = append(configs, convertLoggingConfigToResponse(config))
	}

	resp := map[string]interface{}{
		"LoggingConfigurations": configs,
	}
	pagination.SetNextToken(resp, "NextMarker", result.NextMarker)
	return resp, nil
}

func convertLoggingConfigToResponse(config *wafstore.LoggingConfiguration) map[string]interface{} {
	if config == nil {
		return nil
	}

	result := map[string]interface{}{
		"ResourceArn":              config.ResourceArn,
		"LogDestinationConfigs":    config.LogDestinationConfigs,
		"ManagedByFirewallManager": config.ManagedByFirewallManager,
	}

	if config.LogScope != "" {
		result["LogScope"] = config.LogScope
	}
	if config.LogType != "" {
		result["LogType"] = config.LogType
	}
	if config.LoggingFilter != nil {
		result["LoggingFilter"] = convertLoggingFilterToResponse(config.LoggingFilter)
	}
	if len(config.RedactedFields) > 0 {
		result["RedactedFields"] = config.RedactedFields
	}

	return result
}

func convertLoggingFilterToResponse(lf *wafstore.LoggingFilter) map[string]interface{} {
	if lf == nil {
		return nil
	}
	result := map[string]interface{}{
		"DefaultBehavior": lf.DefaultBehavior,
	}
	if len(lf.Filters) > 0 {
		filters := make([]map[string]interface{}, 0, len(lf.Filters))
		for _, f := range lf.Filters {
			fm := map[string]interface{}{
				"Behavior":    f.Behavior,
				"Requirement": f.Requirement,
			}
			if len(f.Conditions) > 0 {
				conds := make([]map[string]interface{}, 0, len(f.Conditions))
				for _, c := range f.Conditions {
					cm := map[string]interface{}{}
					if c.ActionCondition != nil {
						cm["ActionCondition"] = map[string]interface{}{
							"Action": c.ActionCondition.Action,
						}
					}
					if c.LabelNameCondition != nil {
						cm["LabelNameCondition"] = map[string]interface{}{
							"LabelName": c.LabelNameCondition.LabelName,
						}
					}
					conds = append(conds, cm)
				}
				fm["Conditions"] = conds
			}
			filters = append(filters, fm)
		}
		result["Filters"] = filters
	}
	return result
}
