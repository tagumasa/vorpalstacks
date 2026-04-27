// Package audit provides AWS CloudTrail audit logging functionality for vorpalstacks.
package audit

import (
	"strings"

	"vorpalstacks/internal/common/request"
)

var serviceToEventSource = map[string]string{
	"s3":               "s3.amazonaws.com",
	"lambda":           "lambda.amazonaws.com",
	"dynamodb":         "dynamodb.amazonaws.com",
	"sqs":              "sqs.amazonaws.com",
	"sns":              "sns.amazonaws.com",
	"kms":              "kms.amazonaws.com",
	"cloudtrail":       "cloudtrail.amazonaws.com",
	"cloudwatch":       "monitoring.amazonaws.com",
	"logs":             "logs.amazonaws.com",
	"iam":              "iam.amazonaws.com",
	"sts":              "sts.amazonaws.com",
	"events":           "events.amazonaws.com",
	"scheduler":        "scheduler.amazonaws.com",
	"stepfunctions":    "states.amazonaws.com",
	"cognito-idp":      "cognito-idp.amazonaws.com",
	"apigateway":       "apigateway.amazonaws.com",
	"cloudfront":       "cloudfront.amazonaws.com",
	"route53":          "route53.amazonaws.com",
	"acm":              "acm.amazonaws.com",
	"secretsmanager":   "secretsmanager.amazonaws.com",
	"athena":           "athena.amazonaws.com",
	"kinesis":          "kinesis.amazonaws.com",
	"ssm":              "ssm.amazonaws.com",
	"ses":              "email.amazonaws.com",
	"sesv2":            "email.amazonaws.com",
	"timestream-query": "timestream.amazonaws.com",
	"neptune":          "rds.amazonaws.com",
	"neptunedata":      "rds.amazonaws.com",
	"neptunegraph":     "neptune-graph.amazonaws.com",
	"appsync":          "appsync.amazonaws.com",
	"wafv2":            "wafv2.amazonaws.com",
	"ec2":              "ec2.amazonaws.com",
}

// GetEventSource returns the event source for a given service name.
func GetEventSource(serviceName string) string {
	if source, ok := serviceToEventSource[serviceName]; ok {
		return source
	}
	return serviceName + ".amazonaws.com"
}

// EventBuilder builds audit events for CloudTrail logging.
type EventBuilder struct {
	serviceName string
	operation   string
	reqCtx      *request.RequestContext
	req         *request.ParsedRequest
}

// NewEventBuilder creates a new event builder.
func NewEventBuilder(serviceName, operation string, reqCtx *request.RequestContext, req *request.ParsedRequest) *EventBuilder {
	return &EventBuilder{
		serviceName: serviceName,
		operation:   operation,
		reqCtx:      reqCtx,
		req:         req,
	}
}

// Build builds an audit event from the request and response.
func (b *EventBuilder) Build(response interface{}, err error) *AuditEvent {
	event := &AuditEvent{
		EventName:   b.operation,
		EventSource: GetEventSource(b.serviceName),
	}

	if b.reqCtx != nil {
		event.SourceIP = b.reqCtx.SourceIP
	}

	if b.req != nil {
		event.AccessKeyID = b.req.AccessKeyID
	}

	if b.reqCtx != nil {
		event.AccountID = b.reqCtx.GetAccountID()
		event.PrincipalName = b.reqCtx.Principal
	}

	if b.req != nil && b.req.Parameters != nil {
		event.RequestParameters = b.sanitizeParameters(b.req.Parameters)
	}

	if response != nil {
		event.ResponseElements = b.buildResponseElements(response)
	}

	if err != nil {
		event.ErrorCode = b.extractErrorCode(err)
		event.ErrorMessage = err.Error()
	}

	event.ReadOnly = b.isReadOnlyOperation()

	return event
}

func (b *EventBuilder) sanitizeParameters(params map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})
	sensitiveKeys := map[string]bool{
		"password":    true,
		"secret":      true,
		"token":       true,
		"apikey":      true,
		"api_key":     true,
		"accesskey":   true,
		"access_key":  true,
		"secretkey":   true,
		"secret_key":  true,
		"privatekey":  true,
		"private_key": true,
		"credential":  true,
		"credentials": true,
	}

	for k, v := range params {
		lowerKey := strings.ToLower(k)
		if sensitiveKeys[lowerKey] || strings.Contains(lowerKey, "password") ||
			strings.Contains(lowerKey, "secret") || strings.Contains(lowerKey, "token") {
			sanitized[k] = "****"
		} else {
			sanitized[k] = v
		}
	}
	return sanitized
}

func (b *EventBuilder) buildResponseElements(response interface{}) map[string]interface{} {
	if response == nil {
		return nil
	}

	if m, ok := response.(map[string]interface{}); ok {
		return m
	}
	return nil
}

func (b *EventBuilder) extractErrorCode(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()
	if strings.Contains(errStr, "NotFound") || strings.Contains(errStr, "not found") {
		return "NoSuchResource"
	}
	if strings.Contains(errStr, "AlreadyExists") || strings.Contains(errStr, "already exists") {
		return "ResourceAlreadyExists"
	}
	if strings.Contains(errStr, "AccessDenied") || strings.Contains(errStr, "access denied") {
		return "AccessDenied"
	}
	if strings.Contains(errStr, "InvalidParameter") || strings.Contains(errStr, "invalid parameter") {
		return "InvalidParameter"
	}
	if strings.Contains(errStr, "MissingParameter") || strings.Contains(errStr, "missing parameter") {
		return "MissingParameter"
	}

	return "InternalError"
}

func (b *EventBuilder) isReadOnlyOperation() bool {
	readOnlyPrefixes := []string{
		"Get", "List", "Describe", "Lookup", "Search",
		"BatchGet", "Query", "Scan", "Head",
		"Simulate", "Generate",
	}

	writeExceptions := []string{
		"GetFederationToken",
		"GetSessionToken",
	}

	for _, exc := range writeExceptions {
		if b.operation == exc {
			return false
		}
	}

	for _, prefix := range readOnlyPrefixes {
		if strings.HasPrefix(b.operation, prefix) {
			return true
		}
	}
	return false
}
