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
}

// Recorder defines the interface for recording audit events.
type Recorder interface {
	RecordEvent(event *AuditEvent) error
}
