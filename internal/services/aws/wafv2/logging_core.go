package wafv2

import (
	"vorpalstacks/internal/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// LoggingConfigInput is the transport-agnostic input for
// PutLoggingConfiguration. The logging filter travels raw because its wire
// parse (which carries the filter-behavior validation) must run after the
// scalar member validations, matching the original failure precedence.
type LoggingConfigInput struct {
	ResourceArn              string
	LogDestinationConfigs    []string
	LogScope                 string
	LogType                  string
	ManagedByFirewallManager bool
	RedactedFields           []interface{}
	LoggingFilterRaw         interface{}
}

// putLoggingConfigurationCore is the single entry point for creating or
// updating the logging configuration of a WebACL. The request context is
// taken directly because the member validations precede the store
// acquisition in the original failure precedence.
func (s *WAFv2Service) putLoggingConfigurationCore(reqCtx *request.RequestContext, in LoggingConfigInput) (*wafstore.LoggingConfiguration, error) {
	if in.ResourceArn == "" {
		return nil, invalidParamError("ResourceArn is required")
	}

	if len(in.LogDestinationConfigs) == 0 {
		return nil, invalidParamError("LogDestinationConfigs is required")
	}

	for _, arn := range in.LogDestinationConfigs {
		if err := validateLogDestinationARN(arn); err != nil {
			return nil, err
		}
	}

	if err := validateLogScope(in.LogScope); err != nil {
		return nil, err
	}
	if err := validateLogType(in.LogType); err != nil {
		return nil, err
	}

	loggingFilter, err := parseLoggingFilter(in.LoggingFilterRaw)
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err = stores.webACLs.GetByARN(in.ResourceArn); err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}

	existingConfig, err := stores.loggingConfigs.GetByResourceArn(in.ResourceArn)
	if err == nil && existingConfig != nil {
		return stores.loggingConfigs.Update(in.ResourceArn, in.LogDestinationConfigs, in.LogScope, in.LogType, loggingFilter, in.ManagedByFirewallManager, in.RedactedFields)
	}

	return stores.loggingConfigs.Create(in.ResourceArn, in.LogDestinationConfigs, in.LogScope, in.LogType, loggingFilter, in.ManagedByFirewallManager, in.RedactedFields)
}

// getLoggingConfigurationCore is the single entry point for retrieving the
// logging configuration of a WebACL.
func (s *WAFv2Service) getLoggingConfigurationCore(reqCtx *request.RequestContext, resourceArn string) (*wafstore.LoggingConfiguration, error) {
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

	return config, nil
}

// deleteLoggingConfigurationCore is the single entry point for removing the
// logging configuration of a WebACL.
func (s *WAFv2Service) deleteLoggingConfigurationCore(reqCtx *request.RequestContext, resourceArn string) error {
	if resourceArn == "" {
		return invalidParamError("ResourceArn is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if err := stores.loggingConfigs.Delete(resourceArn); err != nil {
		if wafstore.IsNotFound(err) {
			return notFoundError("LoggingConfiguration")
		}
		return err
	}

	return nil
}

// LoggingListInput is the transport-agnostic input for listing logging
// configurations.
type LoggingListInput struct {
	Scope      string
	Limit      int
	NextMarker string
}

// listLoggingConfigurationsCore is the single entry point for listing
// logging configurations.
func (s *WAFv2Service) listLoggingConfigurationsCore(reqCtx *request.RequestContext, in LoggingListInput) (*wafstore.LoggingConfigurationListResult, error) {
	if err := validateScope(in.Scope); err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return stores.loggingConfigs.List(in.Scope, in.NextMarker, in.Limit)
}
