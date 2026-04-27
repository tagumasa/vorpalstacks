// Package audit provides AWS CloudTrail audit logging functionality for vorpalstacks.
package audit

// AuditEvent represents an audit event to be recorded.
type AuditEvent struct {
	EventName         string
	EventSource       string
	SourceIP          string
	RequestParameters map[string]interface{}
	ResponseElements  map[string]interface{}
	ErrorCode         string
	ErrorMessage      string
	ReadOnly          bool
	AccountID         string
	AccessKeyID       string
	PrincipalName     string
	ResourceTypes     map[string]string
}

// UserIdentity represents the identity information for an audit event.
type UserIdentity struct {
	Type        string // identity type, e.g. "AssumedRole"
	PrincipalID string // unique identifier for the principal
	ARN         string // ARN of the identity
	AccountID   string // AWS account ID
	UserName    string // user or role name
}

// EventStore defines the interface for recording events to a backend store,
// decoupling the audit package from concrete store implementations.
type EventStore interface {
	RecordServiceEvent(eventName, eventSource string, userIdentity *UserIdentity, sourceIP string, requestParams, responseElements map[string]interface{}, resources []ResourceEntry) error
}

// ResourceEntry carries resource information for an audit event.
type ResourceEntry struct {
	ResourceType string
	ResourceName string
}

// Recorder defines the interface for recording audit events.
type Recorder interface {
	RecordEvent(event *AuditEvent) error
}
