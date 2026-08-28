package cognitoidentityprovider

import (
	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// Core functions for the user-pool log delivery family. The handlers
// extract the wire members; validation and store access live here.

// setLogDeliveryConfigurationCore parses and validates the log-configuration
// list and persists the pool's log delivery configuration. LogLevel and
// EventSource are enum-validated per the API model.
func (s *CognitoService) setLogDeliveryConfigurationCore(region, userPoolID string, logConfigurations []interface{}) (*cognitostore.LogDeliveryConfiguration, error) {
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	cfg := &cognitostore.LogDeliveryConfiguration{
		UserPoolID: userPoolID,
	}

	for _, c := range logConfigurations {
		m, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		lc := cognitostore.LogConfiguration{
			LogLevel:    getStringParam(m, "LogLevel"),
			EventSource: getStringParam(m, "EventSource"),
		}
		if lc.LogLevel != "" && lc.LogLevel != "ERROR" && lc.LogLevel != "INFO" {
			return nil, ErrInvalidParameter
		}
		if lc.EventSource != "" && lc.EventSource != "userNotification" && lc.EventSource != "userAuthEvents" {
			return nil, ErrInvalidParameter
		}
		if cw, ok := m["CloudWatchLogsConfiguration"].(map[string]interface{}); ok {
			lc.CloudWatchLogsConfiguration = &cognitostore.CloudWatchLogsConfig{
				LogGroupArn: getStringParam(cw, "LogGroupArn"),
			}
		}
		if s3, ok := m["S3Configuration"].(map[string]interface{}); ok {
			lc.S3Configuration = &cognitostore.S3Config{
				BucketArn: getStringParam(s3, "BucketArn"),
			}
		}
		if fh, ok := m["FirehoseConfiguration"].(map[string]interface{}); ok {
			lc.FirehoseConfiguration = &cognitostore.FirehoseConfig{
				StreamArn: getStringParam(fh, "StreamArn"),
			}
		}
		cfg.LogConfigurations = append(cfg.LogConfigurations, lc)
	}

	if err := store.SaveLogDeliveryConfiguration(cfg); err != nil {
		return nil, ErrInternalError
	}

	return cfg, nil
}

// getLogDeliveryConfigurationCore loads the pool's log delivery
// configuration. A nil configuration means none is stored; the transport
// layer renders the documented empty default in that case.
func (s *CognitoService) getLogDeliveryConfigurationCore(region, userPoolID string) (*cognitostore.LogDeliveryConfiguration, error) {
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	cfg, err := store.GetLogDeliveryConfiguration(userPoolID)
	if err != nil {
		return nil, nil
	}
	return cfg, nil
}

// publishAuthEventLogCore delivers an authentication event to the log
// destinations configured for the userAuthEvents event source. Delivery
// failures are swallowed: logging must never fail the auth flow.
func (s *CognitoService) publishAuthEventLogCore(reqCtx *request.RequestContext, userPoolID string, event *cognitostore.AuthEvent) {
	store, err := s.store(reqCtx)
	if err != nil {
		return
	}

	cfg, err := store.GetLogDeliveryConfiguration(userPoolID)
	if err != nil || cfg == nil {
		return
	}

	for _, lc := range cfg.LogConfigurations {
		if lc.EventSource != "userAuthEvents" {
			continue
		}
		message := formatAuthEventLogMessage(event)
		if lc.CloudWatchLogsConfiguration != nil && lc.CloudWatchLogsConfiguration.LogGroupArn != "" {
			s.publishToCloudWatchLogs(lc.CloudWatchLogsConfiguration.LogGroupArn, userPoolID, message)
		}
	}
}
