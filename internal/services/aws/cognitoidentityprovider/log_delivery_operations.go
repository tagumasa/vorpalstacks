package cognitoidentityprovider

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// SetLogDeliveryConfiguration configures log delivery for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SetLogDeliveryConfiguration.html
func (s *CognitoService) SetLogDeliveryConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	cfg := &cognitostore.LogDeliveryConfiguration{
		UserPoolID: userPoolID,
	}

	if rawConfigs, ok := req.Parameters["LogConfigurations"]; ok {
		if slice, ok := rawConfigs.([]interface{}); ok {
			for _, c := range slice {
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
		}
	}

	if err := store.SaveLogDeliveryConfiguration(cfg); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"LogDeliveryConfiguration": formatLogDeliveryConfiguration(cfg),
	}, nil
}

// GetLogDeliveryConfiguration retrieves the log delivery configuration for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetLogDeliveryConfiguration.html
func (s *CognitoService) GetLogDeliveryConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cfg, err := store.GetLogDeliveryConfiguration(userPoolID)
	if err != nil {
		return map[string]interface{}{
			"LogDeliveryConfiguration": map[string]interface{}{
				"UserPoolId":        userPoolID,
				"LogConfigurations": []interface{}{},
			},
		}, nil
	}

	return map[string]interface{}{
		"LogDeliveryConfiguration": formatLogDeliveryConfiguration(cfg),
	}, nil
}

// publishAuthEventLog checks whether userAuthEvents logging is configured and
// publishes the event to the appropriate delivery targets via EventBus.
func (s *CognitoService) publishAuthEventLog(reqCtx *request.RequestContext, userPoolID string, event *cognitostore.AuthEvent) {
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

func (s *CognitoService) publishToCloudWatchLogs(logGroupArn, userPoolID, message string) {
	if s.bus == nil {
		return
	}

	logGroup := extractLogGroupName(logGroupArn)
	if logGroup == "" {
		return
	}

	evt := &eventbus.CloudWatchLogsPutEvent{
		LogGroup:  logGroup,
		LogStream: "cognito-" + userPoolID,
		LogEvents: []eventbus.LogEntry{
			{
				Timestamp: time.Now().UnixMilli(),
				Message:   message,
			},
		},
	}
	if err := s.bus.Publish(context.Background(), evt); err != nil {
		log.Printf("warning: failed to deliver auth event to CloudWatch Logs for pool %s: %v", userPoolID, err)
	}
}

func extractLogGroupName(arn string) string {
	parts := strings.SplitN(arn, ":", 6)
	if len(parts) < 6 {
		return ""
	}
	resource := parts[5]
	if strings.HasPrefix(resource, "log-group:") {
		return strings.TrimPrefix(resource, "log-group:")
	}
	return ""
}

// formatAuthEventLogMessage renders an authentication event in the JSON
// envelope AWS exports for userAuthEvents log delivery: an outer record with
// the event timestamp, source, and level, and a message object whose members
// mirror the AdminListUserAuthEvents event (fields this platform does not
// populate are logged empty, matching the documented "some fields might be
// logged with null values" contract).
// authEventCreationDateLayout renders the creationDate member the way AWS
// renders it in the exported user activity log (the Java Date.toString
// style, e.g. "Wed Jul 17 17:25:55 UTC 2024").
const authEventCreationDateLayout = "Mon Jan 2 15:04:05 MST 2006"

func formatAuthEventLogMessage(event *cognitostore.AuthEvent) string {
	record := map[string]interface{}{
		"eventTimestamp": strconv.FormatInt(event.CreationDate.UnixMilli(), 10),
		"eventSource":    "USER_ACTIVITY",
		"logLevel":       "INFO",
		"message": map[string]interface{}{
			"version":                       "1",
			"eventId":                       event.EventID,
			"eventType":                     event.EventType,
			"userSub":                       event.UserID,
			"userName":                      event.UserName,
			"userPoolId":                    event.UserPoolID,
			"clientId":                      event.ClientID,
			"creationDate":                  event.CreationDate.Format(authEventCreationDateLayout),
			"eventResponse":                 event.EventResponse,
			"riskLevel":                     event.RiskLevel,
			"riskDecision":                  event.RiskDecision,
			"challenges":                    formatAuthEventChallenges(event),
			"deviceName":                    event.ContextDeviceName,
			"ipAddress":                     event.ContextIPAddress,
			"requestId":                     "",
			"idpName":                       "",
			"compromisedCredentialDetected": strconv.FormatBool(event.CompromisedFlag),
			"city":                          "",
			"country":                       "",
			"eventFeedbackValue":            "",
			"eventFeedbackDate":             "",
			"eventFeedbackProvider":         "",
			"hasContextData":                strconv.FormatBool(event.ContextIPAddress != "" || event.ContextDeviceName != ""),
		},
		"logSourceId": map[string]interface{}{
			"userPoolId": event.UserPoolID,
		},
	}
	buf, err := json.Marshal(record)
	if err != nil {
		return ""
	}
	return string(buf)
}

// formatAuthEventChallenges renders the recorded challenge responses in the
// exported-log shape (name/response pairs).
func formatAuthEventChallenges(event *cognitostore.AuthEvent) []map[string]string {
	challenges := make([]map[string]string, 0, len(event.ChallengeResponses))
	for _, cr := range event.ChallengeResponses {
		challenges = append(challenges, map[string]string{
			"challengeName":     cr.ChallengeName,
			"challengeResponse": cr.ChallengeResponse,
		})
	}
	return challenges
}

func formatLogDeliveryConfiguration(cfg *cognitostore.LogDeliveryConfiguration) map[string]interface{} {
	configs := make([]map[string]interface{}, 0, len(cfg.LogConfigurations))
	for _, lc := range cfg.LogConfigurations {
		entry := map[string]interface{}{
			"LogLevel":    lc.LogLevel,
			"EventSource": lc.EventSource,
		}
		if lc.CloudWatchLogsConfiguration != nil {
			entry["CloudWatchLogsConfiguration"] = map[string]interface{}{
				"LogGroupArn": lc.CloudWatchLogsConfiguration.LogGroupArn,
			}
		}
		if lc.S3Configuration != nil {
			entry["S3Configuration"] = map[string]interface{}{
				"BucketArn": lc.S3Configuration.BucketArn,
			}
		}
		if lc.FirehoseConfiguration != nil {
			entry["FirehoseConfiguration"] = map[string]interface{}{
				"StreamArn": lc.FirehoseConfiguration.StreamArn,
			}
		}
		configs = append(configs, entry)
	}
	return map[string]interface{}{
		"UserPoolId":        cfg.UserPoolID,
		"LogConfigurations": configs,
	}
}

func getStringParam(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
