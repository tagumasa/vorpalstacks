// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

import (
	"strings"
)

// CloudWatchBuilder provides methods for constructing CloudWatch ARNs.
type CloudWatchBuilder struct{ *ARNBuilder }

// CloudWatch returns a CloudWatchBuilder for constructing CloudWatch ARNs.
func (b *ARNBuilder) CloudWatch() *CloudWatchBuilder { return &CloudWatchBuilder{b} }

// LogGroup constructs an ARN for a CloudWatch Logs log group.
func (b *CloudWatchBuilder) LogGroup(name string) string { return b.Build("logs", "log-group:"+name) }

// LogStream constructs an ARN for a CloudWatch Logs log stream.
func (b *CloudWatchBuilder) LogStream(group, stream string) string {
	return b.Build("logs", "log-group:"+group+":log-stream:"+stream)
}

// Destination constructs an ARN for a CloudWatch Logs destination.
func (b *CloudWatchBuilder) Destination(name string) string {
	return b.Build("logs", "destination:"+name)
}

// ScheduledQuery constructs an ARN for a CloudWatch Logs scheduled query.
func (b *CloudWatchBuilder) ScheduledQuery(id string) string {
	return b.Build("logs", "scheduled-query:"+id)
}

// LookupTable constructs an ARN for a CloudWatch Logs lookup table.
func (b *CloudWatchBuilder) LookupTable(name string) string {
	return b.Build("logs", "lookup-table:"+name)
}

// DeliverySource constructs an ARN for a CloudWatch Logs vended-logs
// delivery source (delivery-source:<name>, the resource form the
// PutDeliveryDestinationPolicy documentation's example policy shows).
func (b *CloudWatchBuilder) DeliverySource(name string) string {
	return b.Build("logs", "delivery-source:"+name)
}

// DeliveryDestination constructs an ARN for a CloudWatch Logs vended-logs
// delivery destination (delivery-destination:<name>, the resource form the
// PutDeliveryDestinationPolicy documentation's example policy shows).
func (b *CloudWatchBuilder) DeliveryDestination(name string) string {
	return b.Build("logs", "delivery-destination:"+name)
}

// Delivery constructs an ARN for a CloudWatch Logs vended-logs delivery
// (delivery:<id>, the resource form the PutDeliveryDestinationPolicy
// documentation's example policy shows).
func (b *CloudWatchBuilder) Delivery(id string) string {
	return b.Build("logs", "delivery:"+id)
}

// Alarm constructs an ARN for a CloudWatch alarm.
func (b *CloudWatchBuilder) Alarm(name string) string { return b.Build("cloudwatch", "alarm:"+name) }

// Dashboard constructs an ARN for a CloudWatch dashboard.
func (b *CloudWatchBuilder) Dashboard(name string) string {
	return b.Build("cloudwatch", "dashboard/"+name)
}

// ParseLogGroupName extracts the log group name from a CloudWatch Logs ARN.
func (b *CloudWatchBuilder) ParseLogGroupName(arn string) string {
	return ExtractLogGroupNameFromARN(arn)
}

// CloudTrailBuilder provides methods for constructing CloudTrail ARNs.
type CloudTrailBuilder struct{ *ARNBuilder }

// CloudTrail returns a CloudTrailBuilder for constructing CloudTrail ARNs.
func (b *ARNBuilder) CloudTrail() *CloudTrailBuilder { return &CloudTrailBuilder{b} }

// Trail constructs an ARN for a CloudTrail trail.
func (b *CloudTrailBuilder) Trail(name string) string { return b.Build("cloudtrail", "trail/"+name) }

// EventDataStore constructs an ARN for a CloudTrail event data store. The
// resource word is "eventdatastore/" per the AWS event data store ARN
// format: arn:aws:cloudtrail:region:account:eventdatastore/uuid.
func (b *CloudTrailBuilder) EventDataStore(id string) string {
	return b.Build("cloudtrail", "eventdatastore/"+id)
}

// Channel constructs an ARN for a CloudTrail channel.
func (b *CloudTrailBuilder) Channel(id string) string {
	return b.Build("cloudtrail", "channel/"+id)
}

// ParseTrailName extracts the trail name from a CloudTrail ARN.
func (b *CloudTrailBuilder) ParseTrailName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "trail/") {
		return strings.TrimPrefix(resource, "trail/")
	}
	return ""
}
