// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

import "strings"

// SQSBuilder provides methods for constructing SQS (Simple Queue Service) ARNs.
type SQSBuilder struct{ *ARNBuilder }

// SQS returns an SQSBuilder for constructing SQS ARNs.
func (b *ARNBuilder) SQS() *SQSBuilder { return &SQSBuilder{b} }

// Queue constructs an ARN for an SQS queue.
func (b *SQSBuilder) Queue(name string) string { return b.Build("sqs", name) }

// ParseQueueName extracts the queue name from an SQS queue ARN.
func (b *SQSBuilder) ParseQueueName(arn string) string { return ExtractQueueNameFromARN(arn) }

// SNSBuilder provides methods for constructing SNS (Simple Notification Service) ARNs.
type SNSBuilder struct{ *ARNBuilder }

// SNS returns an SNSBuilder for constructing SNS ARNs.
func (b *ARNBuilder) SNS() *SNSBuilder { return &SNSBuilder{b} }

// Topic constructs an ARN for an SNS topic.
func (b *SNSBuilder) Topic(name string) string { return b.Build("sns", name) }

// PlatformApplication constructs an ARN for an SNS platform application.
func (b *SNSBuilder) PlatformApplication(platform, name string) string {
	return b.Build("sns", "app/"+platform+"/"+name)
}

// PlatformEndpoint constructs an ARN for an SNS platform endpoint.
func (b *SNSBuilder) PlatformEndpoint(platform, name, id string) string {
	return b.Build("sns", "endpoint/"+platform+"/"+name+"/"+id)
}

// Subscription constructs an ARN for an SNS subscription.
func (b *SNSBuilder) Subscription(topicArn, id string) string { return b.Build("sns", id) }

// ParseTopicName extracts the topic name from an SNS topic ARN.
func (b *SNSBuilder) ParseTopicName(arn string) string { return ExtractTopicNameFromARN(arn) }

// EventsBuilder provides methods for constructing EventBridge (CloudWatch Events) ARNs.
type EventsBuilder struct{ *ARNBuilder }

// Events returns an EventsBuilder for constructing EventBridge ARNs.
func (b *ARNBuilder) Events() *EventsBuilder { return &EventsBuilder{b} }

// EventBus constructs an ARN for an EventBridge event bus.
func (b *EventsBuilder) EventBus(name string) string { return b.Build("events", "event-bus/"+name) }

// Rule constructs an ARN for an EventBridge rule.
func (b *EventsBuilder) Rule(name string) string { return b.Build("events", "rule/"+name) }

// RuleOnBus constructs an ARN for an EventBridge rule on a specific event bus.
func (b *EventsBuilder) RuleOnBus(bus, rule string) string {
	return b.Build("events", "rule/"+bus+"/"+rule)
}

// Archive constructs an ARN for an EventBridge archive.
func (b *EventsBuilder) Archive(name string) string { return b.Build("events", "archive/"+name) }

// Connection constructs an ARN for an EventBridge connection.
func (b *EventsBuilder) Connection(name string) string { return b.Build("events", "connection/"+name) }

// ApiDestination constructs an ARN for an EventBridge API destination.
func (b *EventsBuilder) ApiDestination(name string) string {
	return b.Build("events", "api-destination/"+name)
}

// Replay constructs an ARN for an EventBridge replay.
func (b *EventsBuilder) Replay(name string) string { return b.Build("events", "replay/"+name) }

// ParseEventBusName extracts the event bus name from an EventBridge ARN.
func (b *EventsBuilder) ParseEventBusName(arn string) string {
	return ExtractEventBusNameFromARN(arn)
}

// ParseRuleName extracts the rule name from an EventBridge rule ARN.
func (b *EventsBuilder) ParseRuleName(arn string) string { return ExtractRuleNameFromARN(arn) }

// KinesisBuilder provides methods for constructing Kinesis ARNs.
type KinesisBuilder struct{ *ARNBuilder }

// Kinesis returns a KinesisBuilder for constructing Kinesis ARNs.
func (b *ARNBuilder) Kinesis() *KinesisBuilder { return &KinesisBuilder{b} }

// Stream constructs an ARN for a Kinesis stream.
func (b *KinesisBuilder) Stream(name string) string { return b.Build("kinesis", "stream/"+name) }

// Consumer constructs an ARN for a Kinesis stream consumer.
func (b *KinesisBuilder) Consumer(stream, name string) string {
	return b.Build("kinesis", "stream/"+stream+"/consumer/"+name)
}

// ParseStreamName extracts the stream name from a Kinesis stream ARN.
func (b *KinesisBuilder) ParseStreamName(arn string) string { return ExtractStreamNameFromARN(arn) }

// SchedulerBuilder provides methods for constructing EventBridge Scheduler ARNs.
type SchedulerBuilder struct{ *ARNBuilder }

// Scheduler returns a SchedulerBuilder for constructing EventBridge Scheduler ARNs.
func (b *ARNBuilder) Scheduler() *SchedulerBuilder { return &SchedulerBuilder{b} }

// ScheduleGroup constructs an ARN for an EventBridge Scheduler schedule group.
func (b *SchedulerBuilder) ScheduleGroup(name string) string {
	return b.Build("scheduler", "schedule-group/"+name)
}

// Schedule constructs an ARN for an EventBridge Scheduler schedule.
func (b *SchedulerBuilder) Schedule(group, name string) string {
	return b.Build("scheduler", "schedule/"+group+"/"+name)
}

// ParseScheduleGroupName extracts the schedule group name from an EventBridge Scheduler ARN.
func (b *SchedulerBuilder) ParseScheduleGroupName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "schedule-group/") {
		return strings.TrimPrefix(resource, "schedule-group/")
	}
	if strings.HasPrefix(resource, "schedule/") {
		parts := strings.Split(strings.TrimPrefix(resource, "schedule/"), "/")
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

// ParseScheduleName extracts the schedule name from an EventBridge Scheduler ARN.
func (b *SchedulerBuilder) ParseScheduleName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "schedule/") {
		parts := strings.Split(strings.TrimPrefix(resource, "schedule/"), "/")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return ""
}
