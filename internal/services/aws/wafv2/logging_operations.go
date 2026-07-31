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

func parseLoggingFilter(raw interface{}) *wafstore.LoggingFilter {
	if raw == nil {
		return nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	lf := &wafstore.LoggingFilter{
		DefaultBehavior: request.GetStringParam(m, "DefaultBehavior"),
	}
	if filtersRaw, ok := m["Filters"]; ok {
		if arr, ok := filtersRaw.([]interface{}); ok {
			for _, fRaw := range arr {
				fMap, ok := fRaw.(map[string]interface{})
				if !ok {
					continue
				}
				f := wafstore.Filter{
					Behavior:    request.GetStringParam(fMap, "Behavior"),
					Requirement: request.GetStringParam(fMap, "Requirement"),
				}
				if condsRaw, ok := fMap["Conditions"]; ok {
					if condArr, ok := condsRaw.([]interface{}); ok {
						for _, cRaw := range condArr {
							cMap, ok := cRaw.(map[string]interface{})
							if !ok {
								continue
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
		return nil
	}
	return lf
}

// PutLoggingConfiguration creates or updates the logging configuration for the specified web ACL.
func (s *WAFv2Service) PutLoggingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	loggingConfigMap := request.GetMapParam(req.Parameters, "LoggingConfiguration")
	if loggingConfigMap == nil {
		return nil, invalidParamError("LoggingConfiguration is required")
	}

	resourceArn := request.GetStringParam(loggingConfigMap, "ResourceArn")
	if resourceArn == "" {
		return nil, invalidParamError("ResourceArn is required")
	}

	logDestinationConfigs := request.GetStringList(loggingConfigMap, "LogDestinationConfigs")
	if len(logDestinationConfigs) == 0 {
		return nil, invalidParamError("LogDestinationConfigs is required")
	}

	logScope := request.GetStringParam(loggingConfigMap, "LogScope")
	if err := validateLogScope(logScope); err != nil {
		return nil, err
	}
	logType := request.GetStringParam(loggingConfigMap, "LogType")
	if err := validateLogType(logType); err != nil {
		return nil, err
	}
	managedByFirewallManager := request.GetBoolParam(loggingConfigMap, "ManagedByFirewallManager")

	var redactedFields []interface{}
	if rfRaw := loggingConfigMap["RedactedFields"]; rfRaw != nil {
		if arr, ok := rfRaw.([]interface{}); ok {
			redactedFields = arr
		}
	}

	loggingFilter := parseLoggingFilter(loggingConfigMap["LoggingFilter"])

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = stores.webACLs.GetByARN(resourceArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}

	existingConfig, err := stores.loggingConfigs.GetByResourceArn(resourceArn)
	if err == nil && existingConfig != nil {
		config, err := stores.loggingConfigs.Update(resourceArn, logDestinationConfigs, logScope, logType, loggingFilter, managedByFirewallManager, redactedFields)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"LoggingConfiguration": convertLoggingConfigToResponse(config),
		}, nil
	}

	config, err := stores.loggingConfigs.Create(resourceArn, logDestinationConfigs, logScope, logType, loggingFilter, managedByFirewallManager, redactedFields)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"LoggingConfiguration": convertLoggingConfigToResponse(config),
	}, nil
}

// GetLoggingConfiguration retrieves the logging configuration for the specified web ACL.
func (s *WAFv2Service) GetLoggingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, invalidParamError("ResourceArn is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := stores.loggingConfigs.GetByResourceArn(resourceArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("LoggingConfiguration")
		}
		return nil, err
	}

	return map[string]interface{}{
		"LoggingConfiguration": convertLoggingConfigToResponse(config),
	}, nil
}

// DeleteLoggingConfiguration removes the logging configuration for the specified web ACL.
func (s *WAFv2Service) DeleteLoggingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, invalidParamError("ResourceArn is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := stores.loggingConfigs.Delete(resourceArn); err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("LoggingConfiguration")
		}
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListLoggingConfigurations returns a paginated list of all logging configurations.
func (s *WAFv2Service) ListLoggingConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	scope := request.GetStringParam(req.Parameters, "Scope")
	if err := validateScope(scope); err != nil {
		return nil, err
	}

	maxItems := pagination.GetMaxItems(req.Parameters, 100, "Limit")
	nextMarker := pagination.GetMarker(req.Parameters, "NextMarker")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := stores.loggingConfigs.List(scope, nextMarker, maxItems)
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
