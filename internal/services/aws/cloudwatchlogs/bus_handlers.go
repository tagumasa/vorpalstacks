package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
)

// handleBusDelivery handles CloudWatch Logs subscription filter matches,
// delivering matched payloads through the shared destination dispatch. A
// failed dispatch rides the documented retry window instead of dropping.
func (s *LogsService) handleBusDelivery(ctx context.Context, evt *eventbus.CloudWatchLogDeliveryEvent) eventbus.HandlerResult {
	if err := s.dispatchSubscriptionDelivery(evt.Region, evt.DestinationArn, evt.LogGroup, evt.LogStream, evt.Distribution, evt.Payload); err != nil {
		if store, storeErr := s.getLogsStoreByRegion(evt.Region); storeErr == nil {
			s.retryFailedDelivery(store, evt.Region, evt.DestinationArn, evt.LogGroup, evt.LogStream, evt.Distribution, evt.Payload, err)
		} else {
			logs.Warn("Failed to open store for subscription delivery retry",
				logs.String("region", evt.Region), logs.Err(storeErr))
		}
	}
	return eventbus.HandlerResult{}
}

// handleLambdaLogWrite ingests Lambda execution logs through the shared
// ingestion seam (validation + write + metric/subscription filter
// fan-out), auto-creating the log group and stream.
func (s *LogsService) handleLambdaLogWrite(ctx context.Context, evt *eventbus.LambdaLogWriteEvent) eventbus.HandlerResult {
	s.ingestBusEvents("Lambda log write", evt.Region, evt.LogGroup, evt.LogStream, s.accountID, evt.LogEvents)
	return eventbus.HandlerResult{}
}

// handleAPIGatewayAccessLog ingests a single formatted access log entry
// from API Gateway through the shared ingestion seam.
func (s *LogsService) handleAPIGatewayAccessLog(ctx context.Context, evt *eventbus.APIGatewayAccessLogEvent) eventbus.HandlerResult {
	s.ingestBusEvents("API Gateway access log", evt.Region, evt.LogGroup, evt.LogStream, evt.AccountID,
		[]eventbus.LogEntry{busAccessLogEntry(evt.FormattedLog)})
	return eventbus.HandlerResult{}
}

// handleDirectPutLogEvents ingests log events from EventBridge/Scheduler/
// SFN targets through the shared ingestion seam.
func (s *LogsService) handleDirectPutLogEvents(ctx context.Context, evt *eventbus.CloudWatchLogsPutEvent) eventbus.HandlerResult {
	s.ingestBusEvents("direct log events", evt.Region, evt.LogGroup, evt.LogStream, evt.AccountID, evt.LogEvents)
	return eventbus.HandlerResult{}
}
