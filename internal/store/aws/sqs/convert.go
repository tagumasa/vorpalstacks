package sqs

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
)

// Queue conversion

// QueueToProto converts a Queue to its protobuf representation.
func QueueToProto(q *Queue) *pb.Queue {
	if q == nil {
		return nil
	}
	return &pb.Queue{
		Name:                          q.Name,
		Url:                           q.URL,
		Arn:                           q.ARN,
		Region:                        q.Region,
		AccountId:                     q.AccountID,
		CreatedTimestamp:              timestamppb.New(q.CreatedTimestamp),
		LastModifiedTimestamp:         timestamppb.New(q.LastModifiedTimestamp),
		VisibilityTimeout:             q.VisibilityTimeout,
		MaximumMessageSize:            q.MaximumMessageSize,
		MessageRetentionPeriod:        q.MessageRetentionPeriod,
		DelaySeconds:                  q.DelaySeconds,
		ReceiveMessageWaitTimeSeconds: q.ReceiveMessageWaitTimeSeconds,
		Attributes:                    q.Attributes,
		Tags:                          q.Tags,
		Permissions:                   permissionsToProto(q.Permissions),
		Policy:                        q.Policy,
		RedrivePolicy:                 redrivePolicyToProto(q.RedrivePolicy),
		FifoQueue:                     q.FifoQueue,
		ContentBasedDeduplication:     q.ContentBasedDeduplication,
	}
}

// ProtoToQueue converts a protobuf Queue to its internal representation.
func ProtoToQueue(p *pb.Queue) *Queue {
	if p == nil {
		return nil
	}
	return &Queue{
		Name:                          p.Name,
		URL:                           p.Url,
		ARN:                           p.Arn,
		Region:                        p.Region,
		AccountID:                     p.AccountId,
		CreatedTimestamp:              p.CreatedTimestamp.AsTime(),
		LastModifiedTimestamp:         p.LastModifiedTimestamp.AsTime(),
		VisibilityTimeout:             p.VisibilityTimeout,
		MaximumMessageSize:            p.MaximumMessageSize,
		MessageRetentionPeriod:        p.MessageRetentionPeriod,
		DelaySeconds:                  p.DelaySeconds,
		ReceiveMessageWaitTimeSeconds: p.ReceiveMessageWaitTimeSeconds,
		Attributes:                    p.Attributes,
		Tags:                          p.Tags,
		Permissions:                   protoToPermissions(p.Permissions),
		Policy:                        p.Policy,
		RedrivePolicy:                 protoToRedrivePolicy(p.RedrivePolicy),
		FifoQueue:                     p.FifoQueue,
		ContentBasedDeduplication:     p.ContentBasedDeduplication,
	}
}

// Permission conversion

// PermissionToProto converts a Permission to its protobuf representation.
func PermissionToProto(p *Permission) *pb.Permission {
	if p == nil {
		return nil
	}
	return &pb.Permission{
		Label:         p.Label,
		AwsAccountIds: p.AWSAccountIDs,
		Actions:       p.Actions,
	}
}

// ProtoToPermission converts a protobuf Permission to its internal representation.
func ProtoToPermission(p *pb.Permission) *Permission {
	if p == nil {
		return nil
	}
	return &Permission{
		Label:         p.Label,
		AWSAccountIDs: p.AwsAccountIds,
		Actions:       p.Actions,
	}
}

func permissionsToProto(m map[string]*Permission) map[string]*pb.Permission {
	if m == nil {
		return nil
	}
	result := make(map[string]*pb.Permission, len(m))
	for k, v := range m {
		result[k] = PermissionToProto(v)
	}
	return result
}

func protoToPermissions(m map[string]*pb.Permission) map[string]*Permission {
	if m == nil {
		return nil
	}
	result := make(map[string]*Permission, len(m))
	for k, v := range m {
		result[k] = ProtoToPermission(v)
	}
	return result
}

// RedrivePolicy conversion

func redrivePolicyToProto(r *RedrivePolicy) *pb.RedrivePolicy {
	if r == nil {
		return nil
	}
	return &pb.RedrivePolicy{
		DeadLetterTargetArn: r.DeadLetterTargetARN,
		MaxReceiveCount:     r.MaxReceiveCount,
	}
}

func protoToRedrivePolicy(p *pb.RedrivePolicy) *RedrivePolicy {
	if p == nil {
		return nil
	}
	return &RedrivePolicy{
		DeadLetterTargetARN: p.DeadLetterTargetArn,
		MaxReceiveCount:     p.MaxReceiveCount,
	}
}

// Message conversion

// MessageToProto converts a Message to its protobuf representation.
func MessageToProto(m *Message) *pb.Message {
	if m == nil {
		return nil
	}
	return &pb.Message{
		Id:                               m.ID,
		Body:                             m.Body,
		Md5OfBody:                        m.MD5OfBody,
		Md5OfMessageAttributes:           m.MD5OfMessageAttributes,
		MessageAttributes:                messageAttributesToProto(m.MessageAttributes),
		ReceiptHandle:                    m.ReceiptHandle,
		QueueUrl:                         m.QueueURL,
		QueueArn:                         m.QueueARN,
		SentTimestamp:                    timestamppb.New(m.SentTimestamp),
		DelaySeconds:                     m.DelaySeconds,
		VisibilityTimeout:                m.VisibilityTimeout,
		ReceivedAt:                       timeToProto(m.ReceivedAt),
		VisibleAfter:                     timeToProto(m.VisibleAfter),
		ApproximateReceiveCount:          m.ApproximateReceiveCount,
		ApproximateFirstReceiveTimestamp: timeToProto(m.ApproximateFirstReceiveTimestamp),
		SequenceNumber:                   m.SequenceNumber,
		MessageDeduplicationId:           m.MessageDeduplicationID,
		MessageGroupId:                   m.MessageGroupID,
		Attributes:                       m.Attributes,
	}
}

// ProtoToMessage converts a protobuf Message to its internal representation.
func ProtoToMessage(p *pb.Message) *Message {
	if p == nil {
		return nil
	}
	return &Message{
		ID:                               p.Id,
		Body:                             p.Body,
		MD5OfBody:                        p.Md5OfBody,
		MD5OfMessageAttributes:           p.Md5OfMessageAttributes,
		MessageAttributes:                protoToMessageAttributes(p.MessageAttributes),
		ReceiptHandle:                    p.ReceiptHandle,
		QueueURL:                         p.QueueUrl,
		QueueARN:                         p.QueueArn,
		SentTimestamp:                    p.SentTimestamp.AsTime(),
		DelaySeconds:                     p.DelaySeconds,
		VisibilityTimeout:                p.VisibilityTimeout,
		ReceivedAt:                       protoToTime(p.ReceivedAt),
		VisibleAfter:                     protoToTime(p.VisibleAfter),
		ApproximateReceiveCount:          p.ApproximateReceiveCount,
		ApproximateFirstReceiveTimestamp: protoToTime(p.ApproximateFirstReceiveTimestamp),
		SequenceNumber:                   p.SequenceNumber,
		MessageDeduplicationID:           p.MessageDeduplicationId,
		MessageGroupID:                   p.MessageGroupId,
		Attributes:                       p.Attributes,
	}
}

// MessageAttributeValue conversion

func messageAttributeValueToProto(v *MessageAttributeValue) *pb.MessageAttributeValue {
	if v == nil {
		return nil
	}
	result := &pb.MessageAttributeValue{
		BinaryValue:      v.BinaryValue,
		StringListValues: v.StringListValues,
		BinaryListValues: v.BinaryListValues,
		DataType:         v.DataType,
	}
	if v.StringValue != nil {
		result.StringValue = *v.StringValue
	}
	return result
}

func protoToMessageAttributeValue(p *pb.MessageAttributeValue) *MessageAttributeValue {
	if p == nil {
		return nil
	}
	result := &MessageAttributeValue{
		BinaryValue:      p.BinaryValue,
		StringListValues: p.StringListValues,
		BinaryListValues: p.BinaryListValues,
		DataType:         p.DataType,
	}
	if p.StringValue != "" {
		result.StringValue = &p.StringValue
	}
	return result
}

func messageAttributesToProto(m map[string]*MessageAttributeValue) map[string]*pb.MessageAttributeValue {
	if m == nil {
		return nil
	}
	result := make(map[string]*pb.MessageAttributeValue, len(m))
	for k, v := range m {
		result[k] = messageAttributeValueToProto(v)
	}
	return result
}

func protoToMessageAttributes(m map[string]*pb.MessageAttributeValue) map[string]*MessageAttributeValue {
	if m == nil {
		return nil
	}
	result := make(map[string]*MessageAttributeValue, len(m))
	for k, v := range m {
		result[k] = protoToMessageAttributeValue(v)
	}
	return result
}

// MessageMoveTask conversion

// MessageMoveTaskToProto converts a MessageMoveTask to its protobuf representation.
func MessageMoveTaskToProto(t *MessageMoveTask) *pb.MessageMoveTask {
	if t == nil {
		return nil
	}
	return &pb.MessageMoveTask{
		TaskId:              t.TaskId,
		SourceQueueArn:      t.SourceQueueARN,
		DestinationQueueArn: t.DestinationQueueARN,
		Status:              t.Status,
		MaxNumberOfMessages: t.MaxNumberOfMessages,
		MovedMessages:       t.MovedMessages,
		FailureMessages:     t.FailureMessages,
		StartTime:           timestamppb.New(t.StartTime),
		EndTime:             timeToProto(t.EndTime),
	}
}

// ProtoToMessageMoveTask converts a protobuf MessageMoveTask to its internal representation.
func ProtoToMessageMoveTask(p *pb.MessageMoveTask) *MessageMoveTask {
	if p == nil {
		return nil
	}
	return &MessageMoveTask{
		TaskId:              p.TaskId,
		SourceQueueARN:      p.SourceQueueArn,
		DestinationQueueARN: p.DestinationQueueArn,
		Status:              p.Status,
		MaxNumberOfMessages: p.MaxNumberOfMessages,
		MovedMessages:       p.MovedMessages,
		FailureMessages:     p.FailureMessages,
		StartTime:           p.StartTime.AsTime(),
		EndTime:             protoToTime(p.EndTime),
	}
}

// Helper functions

func timeToProto(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func protoToTime(p *timestamppb.Timestamp) time.Time {
	if p == nil {
		return time.Time{}
	}
	return p.AsTime()
}
