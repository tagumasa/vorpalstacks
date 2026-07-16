package cognitoidentityprovider

import (
	"context"
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
	_ = s.bus.Publish(context.Background(), evt)
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

func formatAuthEventLogMessage(event *cognitostore.AuthEvent) string {
	return event.CreationDate.Format(time.RFC3339) + " " +
		event.EventType + " " + event.EventResponse +
		" eventId=" + event.EventID +
		" userId=" + event.UserID
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
