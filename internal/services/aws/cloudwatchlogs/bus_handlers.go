package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/utils/aws/arn"
)

// handleBusDelivery handles CloudWatch Logs subscription filter matches,
// delivering matched log events to Lambda or Kinesis destinations.
func (s *LogsService) handleBusDelivery(ctx context.Context, evt *eventbus.CloudWatchLogDeliveryEvent) eventbus.HandlerResult {
	if arn.IsLambdaARN(evt.DestinationArn) {
		s.invokeLambda(evt.DestinationArn, evt.Payload)
	} else if arn.IsKinesisARN(evt.DestinationArn) {
		s.putToKinesis(evt.DestinationArn, evt.Payload)
	}
	return eventbus.HandlerResult{}
}

// handleLambdaLogWrite ingests Lambda execution logs into CloudWatch Logs,
// then applies metric filters and subscription filters.
func (s *LogsService) handleLambdaLogWrite(ctx context.Context, evt *eventbus.LambdaLogWriteEvent) eventbus.HandlerResult {
	logsStore := s.ensureLogGroupAndStream(evt.Region, evt.LogGroup, evt.LogStream, s.accountID)
	if logsStore == nil {
		return eventbus.HandlerResult{}
	}

	storeEvents := convertBusLogEntries(evt.LogEvents)
	if !s.writeLogEvents(logsStore, evt.LogGroup, evt.LogStream, storeEvents) {
		return eventbus.HandlerResult{}
	}

	s.applyMetricFiltersByRegion(evt.Region, evt.LogGroup, storeEvents)
	s.applySubscriptionFiltersByRegion(evt.Region, evt.LogGroup, evt.LogStream, storeEvents)

	return eventbus.HandlerResult{}
}

// handleAPIGatewayAccessLog writes a single formatted access log entry from
// API Gateway to the specified CloudWatch Logs log group/stream.
func (s *LogsService) handleAPIGatewayAccessLog(ctx context.Context, evt *eventbus.APIGatewayAccessLogEvent) eventbus.HandlerResult {
	s.writeSingleLogMessage(evt.Region, evt.LogGroup, evt.LogStream, evt.AccountID, evt.FormattedLog)
	return eventbus.HandlerResult{}
}

// handleDirectPutLogEvents writes log events from EventBridge/Scheduler/SFN
// targets directly to a CloudWatch Logs log group/stream.
func (s *LogsService) handleDirectPutLogEvents(ctx context.Context, evt *eventbus.CloudWatchLogsPutEvent) eventbus.HandlerResult {
	logsStore := s.ensureLogGroupAndStream(evt.Region, evt.LogGroup, evt.LogStream, evt.AccountID)
	if logsStore == nil {
		return eventbus.HandlerResult{}
	}

	storeEvents := convertBusLogEntries(evt.LogEvents)
	s.writeLogEvents(logsStore, evt.LogGroup, evt.LogStream, storeEvents)

	return eventbus.HandlerResult{}
}
