package kinesis

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// SubscribeToShardEventStream represents an event stream for SubscribeToShard operations.
type SubscribeToShardEventStream struct {
	reader io.Reader
}

// GetStream returns the reader for the stream.
func (e *SubscribeToShardEventStream) GetStream() io.Reader {
	return e.reader
}

// GetStreamHeaders returns the HTTP headers for the stream.
func (e *SubscribeToShardEventStream) GetStreamHeaders() http.Header {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/vnd.amazon.eventstream")
	return headers
}

// RegisterStreamConsumer registers a consumer for a Kinesis stream.
func (s *KinesisService) RegisterStreamConsumer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	consumerName := request.GetParamLowerFirst(req.Parameters, "ConsumerName")

	if streamARN == "" || consumerName == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	consumer, err := store.RegisterStreamConsumer(streamARN, consumerName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"Consumer": map[string]interface{}{
			"ConsumerName":              consumer.ConsumerName,
			"ConsumerARN":               consumer.ConsumerARN,
			"StreamARN":                 consumer.StreamARN,
			"ConsumerStatus":            consumer.ConsumerStatus,
			"ConsumerCreationTimestamp": float64(consumer.ConsumerCreationTimestamp.Unix()),
		},
	}, nil
}

// DeregisterStreamConsumer deregisters a consumer from a Kinesis stream.
func (s *KinesisService) DeregisterStreamConsumer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	consumerARN := request.GetParamLowerFirst(req.Parameters, "ConsumerARN")
	consumerName := request.GetParamLowerFirst(req.Parameters, "ConsumerName")

	if consumerARN == "" && consumerName == "" {
		return nil, ErrInvalidArgument
	}

	if consumerName != "" && streamARN != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		consumer, err := store.GetStreamConsumerByName(streamARN, consumerName)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		consumerARN = consumer.ConsumerARN
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeregisterStreamConsumer(consumerARN); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DescribeStreamConsumer returns details about a Kinesis stream consumer.
func (s *KinesisService) DescribeStreamConsumer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	consumerARN := request.GetParamLowerFirst(req.Parameters, "ConsumerARN")
	consumerName := request.GetParamLowerFirst(req.Parameters, "ConsumerName")

	if consumerARN == "" && consumerName == "" {
		return nil, ErrInvalidArgument
	}

	var consumer *kinesisstore.Consumer

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if consumerARN != "" {
		consumer, err = store.GetStreamConsumer(consumerARN)
	} else if consumerName != "" && streamARN != "" {
		consumer, err = store.GetStreamConsumerByName(streamARN, consumerName)
	} else {
		return nil, ErrInvalidArgument
	}

	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"ConsumerDescription": map[string]interface{}{
			"ConsumerName":              consumer.ConsumerName,
			"ConsumerARN":               consumer.ConsumerARN,
			"StreamARN":                 consumer.StreamARN,
			"ConsumerStatus":            consumer.ConsumerStatus,
			"ConsumerCreationTimestamp": float64(consumer.ConsumerCreationTimestamp.Unix()),
		},
	}, nil
}

// ListStreamConsumers lists consumers of a Kinesis stream.
func (s *KinesisService) ListStreamConsumers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

	if streamARN == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, err := store.GetStreamByARN(streamARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	consumers, err := store.ListStreamConsumers(stream.StreamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var formattedConsumers []map[string]interface{}
	for _, c := range consumers {
		formattedConsumers = append(formattedConsumers, map[string]interface{}{
			"ConsumerName":              c.ConsumerName,
			"ConsumerARN":               c.ConsumerARN,
			"StreamARN":                 c.StreamARN,
			"ConsumerStatus":            c.ConsumerStatus,
			"ConsumerCreationTimestamp": float64(c.ConsumerCreationTimestamp.Unix()),
		})
	}

	return map[string]interface{}{
		"Consumers":        formattedConsumers,
		"NextToken":        nil,
		"HasMoreConsumers": false,
	}, nil
}

// SubscribeToShard subscribes a consumer to receive records from a Kinesis shard.
func (s *KinesisService) SubscribeToShard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	consumerARN := request.GetParamLowerFirst(req.Parameters, "ConsumerARN")
	shardID := request.GetParamLowerFirst(req.Parameters, "ShardId")

	if consumerARN == "" || shardID == "" {
		return nil, ErrInvalidArgument
	}

	startingPosition := ""
	seqNum := ""
	tsStr := ""
	if sp := request.GetMapParam(req.Parameters, "StartingPosition"); sp != nil {
		startingPosition = request.GetStringParam(sp, "Type")
		seqNum = request.GetStringParam(sp, "SequenceNumber")
		tsStr = request.GetStringParam(sp, "Timestamp")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	consumer, err := store.GetStreamConsumer(consumerARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	stream, err := store.GetStreamByARN(consumer.StreamARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var timestamp *time.Time
	if tsStr != "" {
		if unixTs, err := strconv.ParseInt(tsStr, 10, 64); err == nil {
			t := time.Unix(unixTs, 0).UTC()
			timestamp = &t
		}
	}

	iterator, err := store.CreateShardIterator(stream.StreamName, shardID, startingPosition, seqNum, timestamp)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		writer := NewSubscribeToShardEventStreamWriter(pw)

		if err := writer.WriteInitialResponse(); err != nil {
			logs.Error("Failed to write initial response", logs.Err(err))
			return
		}

		includeStart := startingPosition == "AT_SEQUENCE_NUMBER"
		records, lastSeqNum, err := store.GetRecords(stream.StreamName, shardID, iterator.SequenceNumber, 1000, includeStart)
		if err != nil {
			if writeErr := writer.WriteResourceNotFoundException(err.Error()); writeErr != nil {
				logs.Error("Failed to write error event", logs.Err(writeErr))
			}
			return
		}

		if err := writer.WriteSubscribeToShardEvent(records, lastSeqNum, 0, nil); err != nil {
			return
		}

		if err := writer.WriteEndEvent(); err != nil {
			logs.Error("Failed to write end event", logs.Err(err))
		}
	}()

	return &SubscribeToShardEventStream{reader: pr}, nil
}
