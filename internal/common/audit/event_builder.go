// Package audit provides AWS CloudTrail audit logging functionality for vorpalstacks.
package audit

import (
	"errors"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	arnutil "vorpalstacks/internal/utils/aws/arn"
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
		event.UserAgent = b.reqCtx.UserAgent
	}

	if b.req != nil {
		event.AccessKeyID = b.req.AccessKeyID
	}

	if b.reqCtx != nil {
		event.AccountID = b.reqCtx.GetAccountID()
		event.PrincipalName = b.reqCtx.Principal
	}

	event.UserIdentity = b.buildUserIdentity()

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

	event.Resources = b.extractResources(response)

	return event
}

// buildUserIdentity derives the record's userIdentity from the
// authenticated principal on the request context. The CloudTrail record
// vocabulary is matched to what the platform actually resolved: the account
// root principal records as Root, a resolved IAM user as IAMUser with the
// user's real unique ID, a role session as AssumedRole with the session
// issuer context, a federated session as FederatedUser, and an unresolvable
// caller as Unknown — no value is invented.
func (b *EventBuilder) buildUserIdentity() *UserIdentity {
	id := &UserIdentity{}
	if b.reqCtx != nil {
		id.AccountID = b.reqCtx.GetAccountID()
	}
	if b.req != nil {
		id.AccessKeyID = b.req.AccessKeyID
	}
	if b.reqCtx == nil {
		id.Type = "Unknown"
		return id
	}

	switch {
	case b.reqCtx.PrincipalType == request.PrincipalTypeRoot,
		b.reqCtx.Principal == iam.RootUserName:
		id.Type = "Root"
		id.PrincipalID = id.AccountID
		id.ARN = arnutil.NewARNBuilder(id.AccountID, "").IAM().Root()
		return id
	case b.reqCtx.PrincipalType == request.PrincipalTypeRole:
		id.Type = "AssumedRole"
		id.ARN = b.reqCtx.Principal
		id.PrincipalID = b.reqCtx.PrincipalID
		if b.reqCtx.Session != nil {
			id.UserName = b.reqCtx.Session.SessionName
			if b.reqCtx.Session.IssuerPrincipalID != "" && b.reqCtx.Session.SessionName != "" {
				id.PrincipalID = b.reqCtx.Session.IssuerPrincipalID + ":" + b.reqCtx.Session.SessionName
			}
			id.SessionContext = b.buildSessionContext(b.reqCtx.Session)
		}
		return id
	case b.reqCtx.PrincipalType == request.PrincipalTypeUser:
		if b.reqCtx.Session != nil && b.reqCtx.Session.CredentialPrincipalType == "FederatedUser" {
			id.Type = "FederatedUser"
			id.UserName = b.reqCtx.Principal
			id.PrincipalID = b.reqCtx.PrincipalID
			if b.reqCtx.Session.IssuerPrincipalID != "" && b.reqCtx.Session.SessionName != "" {
				id.PrincipalID = b.reqCtx.Session.IssuerPrincipalID + ":" + b.reqCtx.Session.SessionName
			}
			id.ARN = arnutil.NewARNBuilder(id.AccountID, "").STS().FederatedUser(b.reqCtx.Principal)
			id.SessionContext = b.buildSessionContext(b.reqCtx.Session)
			return id
		}
		id.Type = "IAMUser"
		id.PrincipalID = b.reqCtx.PrincipalID
		id.UserName = b.reqCtx.Principal
		if b.reqCtx.Principal != "" {
			id.ARN = arnutil.NewARNBuilder(id.AccountID, "").IAM().User(b.reqCtx.Principal)
		}
		return id
	default:
		// No principal type resolved (permissive plane with an
		// unresolvable access key, or an anonymous caller): CloudTrail
		// records Unknown rather than an invented identity.
		id.Type = "Unknown"
		return id
	}
}

// buildSessionContext reports the issuer and attributes of a resolved
// temporary-credential session.
func (b *EventBuilder) buildSessionContext(info *request.SessionInfo) *SessionContext {
	ctx := &SessionContext{
		SessionIssuer: &SessionIssuer{
			Type:        info.IssuerType,
			PrincipalID: info.IssuerPrincipalID,
			ARN:         info.IssuerARN,
			AccountID:   info.IssuerAccountID,
			UserName:    info.IssuerUserName,
		},
		Attributes: &SessionAttributes{
			MFAAuthenticated: "false",
			CreationDate:     info.CreationDate.UTC().Truncate(time.Second),
		},
	}
	if info.MFAAuthenticated {
		ctx.Attributes.MFAAuthenticated = "true"
	}
	return ctx
}

// extractResources returns the operation's resources as ordered entries.
// Any ARN the request or response names is recorded, with the CloudTrail
// type identifier derived from the ARN itself (AWS::service::type); the S3
// bucket response names the bucket by location rather than ARN, so that
// one form keeps its service-specific reading.
func (b *EventBuilder) extractResources(response interface{}) []ResourceEntry {
	var params map[string]interface{}
	if b.req != nil {
		params = b.req.Parameters
	}
	var resp map[string]interface{}
	if m, ok := response.(map[string]interface{}); ok {
		resp = m
	}

	var entries []ResourceEntry
	if loc, ok := resp["Location"].(string); ok && loc != "" {
		entries = append(entries, ResourceEntry{
			ResourceType: "AWS::S3::Bucket",
			ResourceName: strings.TrimPrefix(loc, "/"),
		})
	}
	for _, arn := range collectResourceARNs(params, resp) {
		entries = append(entries, ResourceEntry{
			ResourceType: resourceTypeFromARN(arn),
			ResourceName: arn,
		})
	}
	return entries
}

func (b *EventBuilder) sanitizeParameters(params map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{}, len(params))
	for k, v := range params {
		if isSensitiveKey(k) {
			sanitized[k] = "****"
			continue
		}
		sanitized[k] = redactSensitiveValues(v)
	}
	return sanitized
}

// redactSensitiveValues applies the same key-based redaction to nested
// structures: a sensitive key buries its value however deeply the wire
// shape nests it.
func redactSensitiveValues(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		nested := make(map[string]interface{}, len(v))
		for k, item := range v {
			if isSensitiveKey(k) {
				nested[k] = "****"
				continue
			}
			nested[k] = redactSensitiveValues(item)
		}
		return nested
	case []interface{}:
		redacted := make([]interface{}, len(v))
		for i, item := range v {
			redacted[i] = redactSensitiveValues(item)
		}
		return redacted
	default:
		return value
	}
}

func isSensitiveKey(k string) bool {
	lowerKey := strings.ToLower(k)
	return sensitiveKeys[lowerKey] ||
		strings.Contains(lowerKey, "password") ||
		strings.Contains(lowerKey, "secret") ||
		strings.Contains(lowerKey, "token")
}

var sensitiveKeys = map[string]bool{
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

// buildResponseElements returns the response of a state-changing action.
// Per the record contents reference, responseElements carries the elements
// "for actions that make changes (create, update, or delete actions). For
// readOnly APIs, this field is null."
func (b *EventBuilder) buildResponseElements(response interface{}) map[string]interface{} {
	if response == nil || b.isReadOnlyOperation() {
		return nil
	}

	if m, ok := response.(map[string]interface{}); ok {
		return m
	}
	return nil
}

// extractErrorCode returns the wire error code of a failed operation. The
// operation errors carry their model-declared code on the AWSError type, so
// the recorded code is the one the caller received; errors without a wire
// identity are internal faults and record AWS's generic InternalError code.
func (b *EventBuilder) extractErrorCode(err error) string {
	if err == nil {
		return ""
	}

	var awsErr *awserrors.AWSError
	if errors.As(err, &awsErr) {
		return awsErr.GetCode()
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
