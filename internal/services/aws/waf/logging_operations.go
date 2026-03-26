package waf

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// PutLoggingConfiguration creates or updates a logging configuration for a web ACL.
func (s *WAFService) PutLoggingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	loggingConfigMap := request.GetMapParam(req.Parameters, "LoggingConfiguration")
	if loggingConfigMap == nil {
		return nil, NewAPIError("InvalidParameterValue", "LoggingConfiguration is required", 400)
	}

	resourceArn := request.GetStringParam(loggingConfigMap, "ResourceArn")
	if resourceArn == "" {
		return nil, NewAPIError("InvalidParameterValue", "ResourceArn is required", 400)
	}

	logDestinationConfigs := request.GetStringList(loggingConfigMap, "LogDestinationConfigs")
	if len(logDestinationConfigs) == 0 {
		return nil, NewAPIError("InvalidParameterValue", "LogDestinationConfigs is required", 400)
	}

	logScope := request.GetStringParam(loggingConfigMap, "LogScope")
	logType := request.GetStringParam(loggingConfigMap, "LogType")
	managedByFirewallManager := request.GetBoolParam(loggingConfigMap, "ManagedByFirewallManager")

	loggingFilter := convertPBToStoreLoggingFilter(request.GetMapParam(loggingConfigMap, "LoggingFilter"))
	redactedFields := convertPBToStoreRedactedFields(request.GetArrayParam(loggingConfigMap, "RedactedFields"))

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	webACL, err := stores.webACLs.GetByARN(resourceArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "WebACL not found", 400)
		}
		return nil, err
	}
	_ = webACL

	existingConfig, err := stores.loggingConfigs.GetByResourceArn(resourceArn)
	if err == nil && existingConfig != nil {
		config, err := stores.loggingConfigs.Update(resourceArn, logDestinationConfigs, logScope, logType, loggingFilter, managedByFirewallManager, redactedFields)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"LoggingConfiguration": convertStoreLoggingConfigToPB(config),
		}, nil
	}

	config, err := stores.loggingConfigs.Create(resourceArn, logDestinationConfigs, logScope, logType, loggingFilter, managedByFirewallManager, redactedFields)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"LoggingConfiguration": convertStoreLoggingConfigToPB(config),
	}, nil
}

// GetLoggingConfiguration returns the logging configuration for a web ACL.
func (s *WAFService) GetLoggingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, NewAPIError("InvalidParameterValue", "ResourceArn is required", 400)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := stores.loggingConfigs.GetByResourceArn(resourceArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "The logging configuration does not exist", 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"LoggingConfiguration": convertStoreLoggingConfigToPB(config),
	}, nil
}

// DeleteLoggingConfiguration deletes the logging configuration for a web ACL.
func (s *WAFService) DeleteLoggingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, NewAPIError("InvalidParameterValue", "ResourceArn is required", 400)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := stores.loggingConfigs.Delete(resourceArn); err != nil {
		if wafstore.IsNotFound(err) {
			return map[string]interface{}{}, nil
		}
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListLoggingConfigurations returns a list of logging configurations.
func (s *WAFService) ListLoggingConfigurations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	scope := request.GetStringParam(req.Parameters, "Scope")
	if scope == "" {
		scope = "REGIONAL"
	}

	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit == 0 {
		limit = 100
	}
	marker := request.GetStringParam(req.Parameters, "NextMarker")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := stores.loggingConfigs.List(scope, marker, limit)
	if err != nil {
		return nil, err
	}

	configs := make([]interface{}, 0, len(result.LoggingConfigurations))
	for _, config := range result.LoggingConfigurations {
		configs = append(configs, convertStoreLoggingConfigToPB(config))
	}

	response := map[string]interface{}{
		"LoggingConfigurations": configs,
	}

	if result.IsTruncated && result.NextMarker != "" {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}

func convertPBToStoreLoggingFilter(pf map[string]interface{}) *wafstore.LoggingFilter {
	if pf == nil {
		return nil
	}

	defaultBehavior, _ := pf["DefaultBehavior"].(string)
	filtersSlice, _ := pf["Filters"].([]interface{})

	filters := make([]wafstore.Filter, 0, len(filtersSlice))
	for _, f := range filtersSlice {
		if fm, ok := f.(map[string]interface{}); ok {
			filter := wafstore.Filter{
				Behavior:    request.GetStringParam(fm, "Behavior"),
				Requirement: request.GetStringParam(fm, "Requirement"),
			}

			conditionsSlice, _ := fm["Conditions"].([]interface{})
			for _, c := range conditionsSlice {
				if cm, ok := c.(map[string]interface{}); ok {
					cond := wafstore.FilterCondition{}
					if ac, ok := cm["ActionCondition"].(map[string]interface{}); ok {
						cond.ActionCondition = &wafstore.ActionCondition{
							Action: request.GetStringParam(ac, "Action"),
						}
					}
					if lc, ok := cm["LabelNameCondition"].(map[string]interface{}); ok {
						cond.LabelNameCondition = &wafstore.LabelNameCondition{
							LabelName: request.GetStringParam(lc, "LabelName"),
						}
					}
					filter.Conditions = append(filter.Conditions, cond)
				}
			}
			filters = append(filters, filter)
		}
	}

	return &wafstore.LoggingFilter{
		DefaultBehavior: defaultBehavior,
		Filters:         filters,
	}
}

func convertPBToStoreRedactedFields(fields []interface{}) []*wafstore.FieldToMatch {
	if len(fields) == 0 {
		return nil
	}

	result := make([]*wafstore.FieldToMatch, 0, len(fields))
	for _, f := range fields {
		if fm, ok := f.(map[string]interface{}); ok {
			result = append(result, convertPBToStoreFieldToMatch(fm))
		}
	}
	return result
}

func convertStoreLoggingConfigToPB(config *wafstore.LoggingConfiguration) map[string]interface{} {
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
		result["LoggingFilter"] = convertStoreLoggingFilterToPB(config.LoggingFilter)
	}
	if len(config.RedactedFields) > 0 {
		result["RedactedFields"] = convertStoreRedactedFieldsToPB(config.RedactedFields)
	}

	return result
}

func convertStoreLoggingFilterToPB(lf *wafstore.LoggingFilter) map[string]interface{} {
	if lf == nil {
		return nil
	}

	filters := make([]interface{}, 0, len(lf.Filters))
	for _, f := range lf.Filters {
		filter := map[string]interface{}{
			"Behavior":    f.Behavior,
			"Requirement": f.Requirement,
		}
		conditions := make([]interface{}, 0, len(f.Conditions))
		for _, c := range f.Conditions {
			cond := map[string]interface{}{}
			if c.ActionCondition != nil {
				cond["ActionCondition"] = map[string]interface{}{
					"Action": c.ActionCondition.Action,
				}
			}
			if c.LabelNameCondition != nil {
				cond["LabelNameCondition"] = map[string]interface{}{
					"LabelName": c.LabelNameCondition.LabelName,
				}
			}
			conditions = append(conditions, cond)
		}
		filter["Conditions"] = conditions
		filters = append(filters, filter)
	}

	return map[string]interface{}{
		"DefaultBehavior": lf.DefaultBehavior,
		"Filters":         filters,
	}
}

func convertStoreRedactedFieldsToPB(fields []*wafstore.FieldToMatch) []interface{} {
	result := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		result = append(result, convertStoreFieldToMatchToPB(f))
	}
	return result
}

func convertStoreFieldToMatchToPB(ftm *wafstore.FieldToMatch) map[string]interface{} {
	if ftm == nil {
		return nil
	}

	result := map[string]interface{}{}
	if ftm.SingleQueryArg != "" {
		result["SingleQueryArgument"] = map[string]interface{}{"Name": ftm.SingleQueryArg}
	}
	if ftm.UriPath != "" {
		result["UriPath"] = map[string]interface{}{}
	}
	if ftm.QueryString != "" {
		result["QueryString"] = map[string]interface{}{}
	}
	if ftm.Body != "" {
		result["Body"] = map[string]interface{}{}
	}
	if ftm.Method != "" {
		result["Method"] = map[string]interface{}{}
	}
	if ftm.Header != "" {
		result["SingleHeader"] = map[string]interface{}{"Name": ftm.Header}
	}
	if ftm.Cookie != "" {
		result["Cookies"] = map[string]interface{}{}
	}

	return result
}

func convertPBToStoreFieldToMatch(fm map[string]interface{}) *wafstore.FieldToMatch {
	if fm == nil {
		return nil
	}

	ftm := &wafstore.FieldToMatch{}

	if sqarg, ok := fm["SingleQueryArgument"].(map[string]interface{}); ok {
		ftm.SingleQueryArg = request.GetStringParam(sqarg, "Name")
	}
	if _, ok := fm["UriPath"].(map[string]interface{}); ok {
		ftm.UriPath = "UriPath"
	}
	if _, ok := fm["QueryString"].(map[string]interface{}); ok {
		ftm.QueryString = "QueryString"
	}
	if _, ok := fm["Body"].(map[string]interface{}); ok {
		ftm.Body = "Body"
	}
	if _, ok := fm["Method"].(map[string]interface{}); ok {
		ftm.Method = "Method"
	}
	if sh, ok := fm["SingleHeader"].(map[string]interface{}); ok {
		ftm.Header = request.GetStringParam(sh, "Name")
	}

	return ftm
}
