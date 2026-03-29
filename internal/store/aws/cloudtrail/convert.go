package cloudtrail

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_cloudtrail"
)

// TrailToProto converts a Trail to its protobuf representation.
func TrailToProto(t *Trail) *pb.Trail {
	if t == nil {
		return nil
	}
	var startedAt, stoppedAt int64
	if t.StartedLoggingAt != nil {
		startedAt = t.StartedLoggingAt.UnixMilli()
	}
	if t.StoppedLoggingAt != nil {
		stoppedAt = t.StoppedLoggingAt.UnixMilli()
	}

	return &pb.Trail{
		Name:                       t.Name,
		TrailArn:                   t.TrailARN,
		S3BucketName:               t.S3BucketName,
		S3KeyPrefix:                t.S3KeyPrefix,
		SnsTopicName:               t.SnsTopicName,
		SnsTopicArn:                t.SnsTopicARN,
		IncludeGlobalServiceEvents: t.IncludeGlobalServiceEvents,
		IsMultiRegionTrail:         t.IsMultiRegionTrail,
		HomeRegion:                 t.HomeRegion,
		IsOrganizationTrail:        t.IsOrganizationTrail,
		IsLogging:                  t.IsLogging,
		LogFileValidationEnabled:   t.LogFileValidationEnabled,
		CloudWatchLogsLogGroupArn:  t.CloudWatchLogsLogGroupARN,
		CloudWatchLogsRoleArn:      t.CloudWatchLogsRoleARN,
		KmsKeyId:                   t.KMSKeyID,
		HasCustomEventSelectors:    t.HasCustomEventSelectors,
		HasInsightSelectors:        t.HasInsightSelectors,
		EventSelectors:             eventSelectorsToProto(t.EventSelectors),
		InsightSelectors:           insightSelectorsToProto(t.InsightSelectors),
		CreatedAt:                  t.CreatedAt.UnixMilli(),
		LastUpdated:                t.LastUpdated.UnixMilli(),
		StartedLoggingAt:           startedAt,
		StoppedLoggingAt:           stoppedAt,
		Tags:                       t.Tags,
	}
}

// ProtoToTrail converts a protobuf Trail to its internal representation.
func ProtoToTrail(p *pb.Trail) *Trail {
	if p == nil {
		return nil
	}
	var startedAt, stoppedAt *time.Time
	if p.StartedLoggingAt > 0 {
		t := time.UnixMilli(p.StartedLoggingAt)
		startedAt = &t
	}
	if p.StoppedLoggingAt > 0 {
		t := time.UnixMilli(p.StoppedLoggingAt)
		stoppedAt = &t
	}

	return &Trail{
		Name:                       p.Name,
		TrailARN:                   p.TrailArn,
		S3BucketName:               p.S3BucketName,
		S3KeyPrefix:                p.S3KeyPrefix,
		SnsTopicName:               p.SnsTopicName,
		SnsTopicARN:                p.SnsTopicArn,
		IncludeGlobalServiceEvents: p.IncludeGlobalServiceEvents,
		IsMultiRegionTrail:         p.IsMultiRegionTrail,
		HomeRegion:                 p.HomeRegion,
		IsOrganizationTrail:        p.IsOrganizationTrail,
		IsLogging:                  p.IsLogging,
		LogFileValidationEnabled:   p.LogFileValidationEnabled,
		CloudWatchLogsLogGroupARN:  p.CloudWatchLogsLogGroupArn,
		CloudWatchLogsRoleARN:      p.CloudWatchLogsRoleArn,
		KMSKeyID:                   p.KmsKeyId,
		HasCustomEventSelectors:    p.HasCustomEventSelectors,
		HasInsightSelectors:        p.HasInsightSelectors,
		EventSelectors:             protoToEventSelectors(p.EventSelectors),
		InsightSelectors:           protoToInsightSelectors(p.InsightSelectors),
		CreatedAt:                  time.UnixMilli(p.CreatedAt),
		LastUpdated:                time.UnixMilli(p.LastUpdated),
		StartedLoggingAt:           startedAt,
		StoppedLoggingAt:           stoppedAt,
		Tags:                       p.Tags,
	}
}

func eventSelectorsToProto(selectors []EventSelector) []*pb.EventSelector {
	if selectors == nil {
		return nil
	}
	result := make([]*pb.EventSelector, len(selectors))
	for i, s := range selectors {
		result[i] = &pb.EventSelector{
			ReadWriteType:           s.ReadWriteType,
			IncludeManagementEvents: s.IncludeManagementEvents,
			DataResources:           dataResourcesToProto(s.DataResources),
		}
	}
	return result
}

func protoToEventSelectors(selectors []*pb.EventSelector) []EventSelector {
	if selectors == nil {
		return nil
	}
	result := make([]EventSelector, len(selectors))
	for i, s := range selectors {
		result[i] = EventSelector{
			ReadWriteType:           s.ReadWriteType,
			IncludeManagementEvents: s.IncludeManagementEvents,
			DataResources:           protoToDataResources(s.DataResources),
		}
	}
	return result
}

func dataResourcesToProto(resources []DataResource) []*pb.DataResource {
	if resources == nil {
		return nil
	}
	result := make([]*pb.DataResource, len(resources))
	for i, r := range resources {
		result[i] = &pb.DataResource{
			Type:   r.Type,
			Values: r.Values,
		}
	}
	return result
}

func protoToDataResources(resources []*pb.DataResource) []DataResource {
	if resources == nil {
		return nil
	}
	result := make([]DataResource, len(resources))
	for i, r := range resources {
		result[i] = DataResource{
			Type:   r.Type,
			Values: r.Values,
		}
	}
	return result
}

func insightSelectorsToProto(selectors []InsightSelector) []*pb.InsightSelector {
	if selectors == nil {
		return nil
	}
	result := make([]*pb.InsightSelector, len(selectors))
	for i, s := range selectors {
		result[i] = &pb.InsightSelector{
			InsightType: s.InsightType,
		}
	}
	return result
}

func protoToInsightSelectors(selectors []*pb.InsightSelector) []InsightSelector {
	if selectors == nil {
		return nil
	}
	result := make([]InsightSelector, len(selectors))
	for i, s := range selectors {
		result[i] = InsightSelector{
			InsightType: s.InsightType,
		}
	}
	return result
}

// EventToProto converts an Event to its protobuf representation.
func EventToProto(e *Event) *pb.Event {
	if e == nil {
		return nil
	}
	var reqParams, respElements string
	if e.RequestParameters != nil {
		if b, err := json.Marshal(e.RequestParameters); err == nil {
			reqParams = string(b)
		}
	}
	if e.ResponseElements != nil {
		if b, err := json.Marshal(e.ResponseElements); err == nil {
			respElements = string(b)
		}
	}

	return &pb.Event{
		EventId:               e.EventID,
		EventName:             e.EventName,
		ReadOnly:              e.ReadOnly,
		AccessKeyId:           e.AccessKeyId,
		EventSource:           e.EventSource,
		EventTime:             e.EventTime.UnixMilli(),
		EventType:             e.EventType,
		EventVersion:          e.EventVersion,
		UserIdentity:          userIdentityToProto(e.UserIdentity),
		Resources:             resourcesToProto(e.Resources),
		CloudTrailEvent:       e.CloudTrailEvent,
		RequestParametersJson: reqParams,
		ResponseElementsJson:  respElements,
		RequestId:             e.RequestID,
		SourceIpAddress:       e.SourceIPAddress,
		UserAgent:             e.UserAgent,
		ErrorCode:             e.ErrorCode,
		ErrorMessage:          e.ErrorMessage,
		Tags:                  e.Tags,
	}
}

// ProtoToEvent converts a protobuf Event to its internal representation.
func ProtoToEvent(p *pb.Event) *Event {
	if p == nil {
		return nil
	}
	var reqParams, respElements map[string]interface{}
	if p.RequestParametersJson != "" {
		if err := json.Unmarshal([]byte(p.RequestParametersJson), &reqParams); err != nil {
			logs.Error("Failed to unmarshal RequestParametersJson", logs.Err(err))
		}
	}
	if p.ResponseElementsJson != "" {
		if err := json.Unmarshal([]byte(p.ResponseElementsJson), &respElements); err != nil {
			logs.Error("Failed to unmarshal ResponseElementsJson", logs.Err(err))
		}
	}

	return &Event{
		EventID:           p.EventId,
		EventName:         p.EventName,
		ReadOnly:          p.ReadOnly,
		AccessKeyId:       p.AccessKeyId,
		EventSource:       p.EventSource,
		EventTime:         time.UnixMilli(p.EventTime),
		EventType:         p.EventType,
		EventVersion:      p.EventVersion,
		UserIdentity:      protoToUserIdentity(p.UserIdentity),
		Resources:         protoToResources(p.Resources),
		CloudTrailEvent:   p.CloudTrailEvent,
		RequestParameters: reqParams,
		ResponseElements:  respElements,
		RequestID:         p.RequestId,
		SourceIPAddress:   p.SourceIpAddress,
		UserAgent:         p.UserAgent,
		ErrorCode:         p.ErrorCode,
		ErrorMessage:      p.ErrorMessage,
		Tags:              p.Tags,
	}
}

func userIdentityToProto(ui *UserIdentity) *pb.UserIdentity {
	if ui == nil {
		return nil
	}
	return &pb.UserIdentity{
		Type:           ui.Type,
		PrincipalId:    ui.PrincipalID,
		Arn:            ui.ARN,
		AccountId:      ui.AccountID,
		AccessKeyId:    ui.AccessKeyID,
		UserName:       ui.UserName,
		SessionContext: sessionContextToProto(ui.SessionContext),
	}
}

func protoToUserIdentity(p *pb.UserIdentity) *UserIdentity {
	if p == nil {
		return nil
	}
	return &UserIdentity{
		Type:           p.Type,
		PrincipalID:    p.PrincipalId,
		ARN:            p.Arn,
		AccountID:      p.AccountId,
		AccessKeyID:    p.AccessKeyId,
		UserName:       p.UserName,
		SessionContext: protoToSessionContext(p.SessionContext),
	}
}

func sessionContextToProto(sc *SessionContext) *pb.SessionContext {
	if sc == nil {
		return nil
	}
	return &pb.SessionContext{
		SessionIssuer: sessionIssuerToProto(sc.SessionIssuer),
		Attributes:    sessionAttributesToProto(sc.Attributes),
	}
}

func protoToSessionContext(p *pb.SessionContext) *SessionContext {
	if p == nil {
		return nil
	}
	return &SessionContext{
		SessionIssuer: protoToSessionIssuer(p.SessionIssuer),
		Attributes:    protoToSessionAttributes(p.Attributes),
	}
}

func sessionIssuerToProto(si *SessionIssuer) *pb.SessionIssuer {
	if si == nil {
		return nil
	}
	return &pb.SessionIssuer{
		Type:        si.Type,
		PrincipalId: si.PrincipalID,
		Arn:         si.ARN,
		UserName:    si.UserName,
	}
}

func protoToSessionIssuer(p *pb.SessionIssuer) *SessionIssuer {
	if p == nil {
		return nil
	}
	return &SessionIssuer{
		Type:        p.Type,
		PrincipalID: p.PrincipalId,
		ARN:         p.Arn,
		UserName:    p.UserName,
	}
}

func sessionAttributesToProto(sa *SessionAttributes) *pb.SessionAttributes {
	if sa == nil {
		return nil
	}
	return &pb.SessionAttributes{
		CreationDate:     sa.CreationDate.UnixMilli(),
		MfaAuthenticated: sa.MFAAuthenticated,
	}
}

func protoToSessionAttributes(p *pb.SessionAttributes) *SessionAttributes {
	if p == nil {
		return nil
	}
	return &SessionAttributes{
		CreationDate:     time.UnixMilli(p.CreationDate),
		MFAAuthenticated: p.MfaAuthenticated,
	}
}

func resourcesToProto(resources []Resource) []*pb.Resource {
	if resources == nil {
		return nil
	}
	result := make([]*pb.Resource, len(resources))
	for i, r := range resources {
		result[i] = &pb.Resource{
			ResourceType: r.ResourceType,
			ResourceName: r.ResourceName,
			Arn:          r.ARN,
		}
	}
	return result
}

func protoToResources(resources []*pb.Resource) []Resource {
	if resources == nil {
		return nil
	}
	result := make([]Resource, len(resources))
	for i, r := range resources {
		result[i] = Resource{
			ResourceType: r.ResourceType,
			ResourceName: r.ResourceName,
			ARN:          r.Arn,
		}
	}
	return result
}

// ResourcePolicyToProto converts a ResourcePolicy to its protobuf representation.
func ResourcePolicyToProto(rp *ResourcePolicy) *pb.ResourcePolicy {
	if rp == nil {
		return nil
	}
	return &pb.ResourcePolicy{
		ResourceArn: rp.ResourceARN,
		Policy:      rp.Policy,
	}
}

// ProtoToResourcePolicy converts a protobuf ResourcePolicy to its internal representation.
func ProtoToResourcePolicy(p *pb.ResourcePolicy) *ResourcePolicy {
	if p == nil {
		return nil
	}
	return &ResourcePolicy{
		ResourceARN: p.ResourceArn,
		Policy:      p.Policy,
	}
}
