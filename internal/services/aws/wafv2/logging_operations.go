package wafv2

import (
	"context"

	"vorpalstacks/internal/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// PutLoggingConfiguration creates or updates the logging configuration for the specified web ACL.
func (s *WAFv2Service) PutLoggingConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	loggingConfigMap := request.GetMapParam(req.Parameters, "LoggingConfiguration")
	if loggingConfigMap == nil {
		return nil, validationError("LoggingConfiguration is required")
	}

	resourceArn := request.GetStringParam(loggingConfigMap, "ResourceArn")
	if resourceArn == "" {
		return nil, validationError("ResourceArn is required")
	}

	logDestinationConfigs := request.GetStringList(loggingConfigMap, "LogDestinationConfigs")
	if len(logDestinationConfigs) == 0 {
		return nil, validationError("LogDestinationConfigs is required")
	}

	logScope := request.GetStringParam(loggingConfigMap, "LogScope")
	logType := request.GetStringParam(loggingConfigMap, "LogType")
	managedByFirewallManager := request.GetBoolParam(loggingConfigMap, "ManagedByFirewallManager")

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
		config, err := stores.loggingConfigs.Update(resourceArn, logDestinationConfigs, logScope, logType, nil, managedByFirewallManager, nil)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"LoggingConfiguration": convertLoggingConfigToResponse(config),
		}, nil
	}

	config, err := stores.loggingConfigs.Create(resourceArn, logDestinationConfigs, logScope, logType, nil, managedByFirewallManager, nil)
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
		return nil, validationError("ResourceArn is required")
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
		return nil, validationError("ResourceArn is required")
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
		configs = append(configs, convertLoggingConfigToResponse(config))
	}

	resp := map[string]interface{}{
		"LoggingConfigurations": configs,
	}

	if result.IsTruncated && result.NextMarker != "" {
		resp["NextMarker"] = result.NextMarker
	}

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

	return result
}
