// Package audit provides AWS CloudTrail audit logging functionality for vorpalstacks.
package audit

import "time"

// AuditEvent represents an audit event to be recorded.
type AuditEvent struct {
	EventName         string
	EventSource       string
	SourceIP          string
	UserAgent         string
	RequestParameters map[string]interface{}
	ResponseElements  map[string]interface{}
	ErrorCode         string
	ErrorMessage      string
	ReadOnly          bool
	AccountID         string
	AccessKeyID       string
	PrincipalName     string
	Resources         []ResourceEntry
	// UserIdentity is derived from the authenticated principal on the
	// request context; the recorder persists it verbatim.
	UserIdentity *UserIdentity
}

// UserIdentity represents the identity information for an audit event,
// mirroring the CloudTrail record userIdentity element. Type values follow
// the recorded CloudTrail vocabulary: Root, IAMUser, AssumedRole,
// FederatedUser, and Unknown when the platform cannot classify the caller.
type UserIdentity struct {
	Type           string
	PrincipalID    string
	ARN            string
	AccountID      string
	AccessKeyID    string
	UserName       string
	SessionContext *SessionContext
}

// SessionContext reports how temporary credentials were obtained, per the
// CloudTrail record userIdentity.sessionContext element.
type SessionContext struct {
	SessionIssuer *SessionIssuer
	Attributes    *SessionAttributes
}

// SessionIssuer describes the principal that issued a temporary credential.
type SessionIssuer struct {
	Type        string // Root, IAMUser, or Role
	PrincipalID string
	ARN         string
	AccountID   string
	UserName    string
}

// SessionAttributes carries the session attributes of a temporary credential.
type SessionAttributes struct {
	MFAAuthenticated string // "true"/"false"
	CreationDate     time.Time
}

// EventStore defines the interface for recording events to a backend store,
// decoupling the audit package from concrete store implementations. The
// signature carries every field an audit record persists: the operation's
// read-only classification, the typed error outcome of a failed call, and
// the caller's user agent, none of which the store may re-derive or drop.
type EventStore interface {
	RecordServiceEvent(eventName, eventSource string, userIdentity *UserIdentity, sourceIP, accessKeyID, userAgent string, readOnly bool, errorCode, errorMessage string, requestParams, responseElements map[string]interface{}, resources []ResourceEntry) error
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
