package audit

var _ Recorder = (*CloudTrailRecorder)(nil)

// CloudTrailRecorder records audit events to CloudTrail.
type CloudTrailRecorder struct {
	store EventStore
}

// NewCloudTrailRecorder creates a new CloudTrail recorder backed by the given store.
func NewCloudTrailRecorder(store EventStore) *CloudTrailRecorder {
	return &CloudTrailRecorder{store: store}
}

// RecordEvent records a CloudTrail audit event. The userIdentity was
// derived from the authenticated principal by the event builder; the
// recorder persists it verbatim.
func (r *CloudTrailRecorder) RecordEvent(event *AuditEvent) error {
	userIdentity := event.UserIdentity
	if userIdentity == nil {
		userIdentity = &UserIdentity{Type: "Unknown", AccountID: event.AccountID}
	}

	return r.store.RecordServiceEvent(
		event.EventName,
		event.EventSource,
		userIdentity,
		event.SourceIP,
		event.AccessKeyID,
		event.UserAgent,
		event.ReadOnly,
		event.ErrorCode,
		event.ErrorMessage,
		event.RequestParameters,
		event.ResponseElements,
		event.Resources,
	)
}

// Record records an audit event if it is an AuditEvent.
func (r *CloudTrailRecorder) Record(event interface{}) error {
	if e, ok := event.(*AuditEvent); ok {
		return r.RecordEvent(e)
	}
	return nil
}
