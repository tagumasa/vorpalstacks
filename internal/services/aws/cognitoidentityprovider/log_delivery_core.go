package cognitoidentityprovider

import (
	"encoding/json"
	"strconv"
	"time"

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
		if lc.LogLevel != "ERROR" && lc.LogLevel != "INFO" {
			return nil, ErrInvalidParameter
		}
		if lc.EventSource != "userNotification" && lc.EventSource != "userAuthEvents" {
			return nil, ErrInvalidParameter
		}
		// The developer guide binds the log level to the event source:
		// message-delivery logs are userNotification at ERROR alone, and
		// user-activity logs are userAuthEvents at INFO alone.
		switch lc.EventSource {
		case "userNotification":
			if lc.LogLevel != "ERROR" {
				return nil, ErrInvalidParameter
			}
		case "userAuthEvents":
			if lc.LogLevel != "INFO" {
				return nil, ErrInvalidParameter
			}
		}
		// A configuration names exactly one destination, and the
		// userNotification source goes to CloudWatch Logs alone — "You
		// can't send user-notification logs to any destination other than
		// CloudWatch Logs".
		cwRaw, cwPresent := m["CloudWatchLogsConfiguration"]
		s3Raw, s3Present := m["S3Configuration"]
		fhRaw, fhPresent := m["FirehoseConfiguration"]
		destinations := 0
		for _, present := range []bool{cwPresent, s3Present, fhPresent} {
			if present {
				destinations++
			}
		}
		if destinations != 1 {
			return nil, ErrInvalidParameter
		}
		if lc.EventSource == "userNotification" && !cwPresent {
			return nil, ErrInvalidParameter
		}
		if cwPresent {
			cw, isMap := cwRaw.(map[string]interface{})
			if !isMap {
				return nil, ErrInvalidParameter
			}
			lc.CloudWatchLogsConfiguration = &cognitostore.CloudWatchLogsConfig{
				LogGroupArn: getStringParam(cw, "LogGroupArn"),
			}
		}
		if s3Present {
			s3, isMap := s3Raw.(map[string]interface{})
			if !isMap {
				return nil, ErrInvalidParameter
			}
			lc.S3Configuration = &cognitostore.S3Config{
				BucketArn: getStringParam(s3, "BucketArn"),
			}
		}
		if fhPresent {
			fh, isMap := fhRaw.(map[string]interface{})
			if !isMap {
				return nil, ErrInvalidParameter
			}
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
		return nil, ErrInternalError
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
			s.publishToCloudWatchLogs(lc.CloudWatchLogsConfiguration.LogGroupArn, reqCtx.GetRegion(), userPoolID, message)
		}
		if lc.S3Configuration != nil && lc.S3Configuration.BucketArn != "" {
			s.publishToS3(lc.S3Configuration.BucketArn, reqCtx.GetRegion(), userPoolID, event.EventID, message)
		}
		// The Firehose destination is accepted and stored but delivers
		// nothing until the platform Firehose service exists (release
		// blocker; category-(a) missing substrate).
	}
}

// publishNotificationLogCore delivers a message-delivery notification
// record to the userNotification log destination. The userNotification
// source is CloudWatch Logs alone, and delivery failures are swallowed
// like every log destination — logging must never fail the notification
// flow.
func (s *CognitoService) publishNotificationLogCore(reqCtx *request.RequestContext, userPoolID, details string) {
	store, err := s.store(reqCtx)
	if err != nil {
		return
	}

	cfg, err := store.GetLogDeliveryConfiguration(userPoolID)
	if err != nil || cfg == nil {
		return
	}

	for _, lc := range cfg.LogConfigurations {
		if lc.EventSource != "userNotification" {
			continue
		}
		if lc.CloudWatchLogsConfiguration == nil || lc.CloudWatchLogsConfiguration.LogGroupArn == "" {
			continue
		}
		record, merr := json.Marshal(map[string]interface{}{
			"eventTimestamp": strconv.FormatInt(time.Now().UnixMilli(), 10),
			"eventSource":    "USER_NOTIFICATION",
			"logLevel":       "ERROR",
			"message":        map[string]interface{}{"details": details},
			"logSourceId":    map[string]interface{}{"userPoolId": userPoolID},
		})
		if merr != nil {
			return
		}
		s.publishToCloudWatchLogs(lc.CloudWatchLogsConfiguration.LogGroupArn, reqCtx.GetRegion(), userPoolID, string(record))
	}
}
