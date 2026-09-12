package cognitoidentityprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// SetLogDeliveryConfiguration configures log delivery for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SetLogDeliveryConfiguration.html
func (s *CognitoService) SetLogDeliveryConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var logConfigs []interface{}
	if raw, ok := req.Parameters["LogConfigurations"]; ok {
		if slice, ok := raw.([]interface{}); ok {
			logConfigs = slice
		}
	}

	cfg, err := s.setLogDeliveryConfigurationCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), logConfigs)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"LogDeliveryConfiguration": formatLogDeliveryConfiguration(cfg),
	}, nil
}

// GetLogDeliveryConfiguration retrieves the log delivery configuration for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetLogDeliveryConfiguration.html
func (s *CognitoService) GetLogDeliveryConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")

	cfg, err := s.getLogDeliveryConfigurationCore(reqCtx.GetRegion(), userPoolID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
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

// publishToCloudWatchLogs publishes a log-delivery record to a CloudWatch
// Logs destination through the service bus. The event carries the region
// and account: the logs handler routes to its regional store on those
// fields, so an event without them lands outside every region the API
// serves.
func (s *CognitoService) publishToCloudWatchLogs(logGroupArn, region, userPoolID, message string) {
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
	evt.Region = region
	evt.AccountID = s.accountID
	if err := s.bus.Publish(context.Background(), evt); err != nil {
		logs.Warn("cognito auth event delivery to cloudwatch logs failed", logs.String("pool", userPoolID), logs.Err(err))
	}
}

// publishToS3 delivers a log record to the destination bucket as one JSON
// object under the AWSLogs prefix — the folder structure AWS's vended logs
// create in the destination bucket.
func (s *CognitoService) publishToS3(bucketArn, region, userPoolID, eventID, message string) {
	if s.bus == nil {
		return
	}
	bucket := extractBucketName(bucketArn)
	if bucket == "" {
		return
	}
	invoker := s.bus.S3Invoker()
	if invoker == nil {
		return
	}
	key := fmt.Sprintf("AWSLogs/%s/cognito-idp/%s/%s/%d-%s.json",
		s.accountID, region, userPoolID, time.Now().UnixMilli(), eventID)
	if err := invoker.PutObject(context.Background(), region, bucket, key, []byte(message), "application/json"); err != nil {
		logs.Warn("cognito auth event delivery to s3 failed", logs.String("pool", userPoolID), logs.Err(err))
	}
}

// extractBucketName reads the bucket name out of an S3 bucket ARN, whose
// resource segment is the bucket name alone.
func extractBucketName(arn string) string {
	parsed, err := svcarn.ParseARN(arn)
	if err != nil {
		return ""
	}
	return parsed.Resource
}

// extractLogGroupName reads the log-group name out of a CloudWatch Logs log
// group ARN. The resource segment of the ARN carries it under the
// log-group: prefix; an ARN of any other resource shape names no group.
func extractLogGroupName(arn string) string {
	parsed, err := svcarn.ParseARN(arn)
	if err != nil {
		return ""
	}
	if strings.HasPrefix(parsed.Resource, "log-group:") {
		return strings.TrimPrefix(parsed.Resource, "log-group:")
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
