package cloudwatchlogs

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/utils/aws/arn"
)

// handleBusDelivery handles CloudWatch Logs subscription filter matches,
// delivering matched log events to Lambda or Kinesis destinations.
func (s *LogsService) handleBusDelivery(ctx context.Context, evt *eventbus.CloudWatchLogDeliveryEvent) eventbus.HandlerResult {
	if arn.IsLambdaARN(evt.DestinationArn) {
		if s.bus == nil {
			return eventbus.HandlerResult{}
		}

		encodedData := base64.StdEncoding.EncodeToString(evt.Payload)

		payload := map[string]interface{}{
			"awslogs": map[string]interface{}{
				"data": encodedData,
			},
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return eventbus.HandlerResult{}
		}

		_, _, invokeErr := s.bus.LambdaInvoker().InvokeForGateway(context.Background(), evt.DestinationArn, payloadBytes)
		if invokeErr != nil {
			logs.Warn("failed to invoke Lambda from subscription filter", logs.Err(invokeErr), logs.String("destinationArn", evt.DestinationArn))
		}
	} else if arn.IsKinesisARN(evt.DestinationArn) {
		if s.bus == nil {
			return eventbus.HandlerResult{}
		}

		streamName := arn.ExtractStreamNameFromARN(evt.DestinationArn)

		encodedData := base64.StdEncoding.EncodeToString(evt.Payload)

		shards, err := s.bus.KinesisInvoker().ListShards(context.Background(), streamName)
		if err != nil || len(shards) == 0 {
			return eventbus.HandlerResult{}
		}

		var activeShardID string
		for _, shard := range shards {
			if shard.SequenceNumberRangeEnd == "" {
				activeShardID = shard.ShardID
				break
			}
		}

		if activeShardID != "" {
			envelope := map[string]interface{}{
				"awslogs": map[string]interface{}{
					"data": encodedData,
				},
			}
			envelopeBytes, marshalErr := json.Marshal(envelope)
			if marshalErr != nil {
				return eventbus.HandlerResult{}
			}
			b64Envelope := base64.StdEncoding.EncodeToString(envelopeBytes)
			if _, putErr := s.bus.KinesisInvoker().PutRecord(context.Background(), streamName, activeShardID, []byte(b64Envelope)); putErr != nil {
				logs.Warn("failed to deliver subscription filter log events to Kinesis", logs.Err(putErr), logs.String("streamName", streamName))
			}
		}
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
