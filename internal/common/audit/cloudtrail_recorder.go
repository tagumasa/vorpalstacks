package audit

import (
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

var _ Recorder = (*CloudTrailRecorder)(nil)

// CloudTrailRecorder records audit events to CloudTrail.
type CloudTrailRecorder struct {
	store cloudtrailstore.CloudTrailStoreInterface
}

// NewCloudTrailRecorder creates a new CloudTrail recorder backed by the given store.
func NewCloudTrailRecorder(store cloudtrailstore.CloudTrailStoreInterface) *CloudTrailRecorder {
	return &CloudTrailRecorder{store: store}
}

// RecordEvent records a CloudTrail audit event.
func (r *CloudTrailRecorder) RecordEvent(event *AuditEvent) error {
	userIdentity := &cloudtrailstore.UserIdentity{
		Type:      "AssumedRole",
		AccountID: event.AccountID,
		UserName:  event.PrincipalName,
	}
	if event.AccessKeyID != "" {
		if len(event.AccessKeyID) >= 16 {
			userIdentity.PrincipalID = "AIDAI" + event.AccessKeyID[:16]
		} else {
			userIdentity.PrincipalID = "AIDAI" + event.AccessKeyID
		}
		userIdentity.ARN = "arn:aws:sts::" + event.AccountID + ":assumed-role/vorpalstacks/" + event.AccessKeyID
	} else {
		userIdentity.PrincipalID = "vorpalstacks:vorpalstacks"
		userIdentity.ARN = "arn:aws:sts::" + event.AccountID + ":assumed-role/vorpalstacks/vorpalstacks"
	}

	return r.store.RecordServiceEvent(
		event.EventName,
		event.EventSource,
		userIdentity,
		event.SourceIP,
		event.RequestParameters,
		event.ResponseElements,
	)
}

// Record records an audit event if it is an AuditEvent.
func (r *CloudTrailRecorder) Record(event interface{}) error {
	if e, ok := event.(*AuditEvent); ok {
		return r.RecordEvent(e)
	}
	return nil
}
