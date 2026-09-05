package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/invokers"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ddbStreamsInvoker adapts the DynamoDB service to the
// invokers.DynamoDBStreamsInvoker interface. It delegates to the
// per-region StreamStore.
type ddbStreamsInvoker struct {
	svc *DynamoDBService
}

// NewDynamoDBStreamsInvoker creates an invokers.DynamoDBStreamsInvoker
// backed by the given DynamoDB service.
func NewDynamoDBStreamsInvoker(svc *DynamoDBService) invokers.DynamoDBStreamsInvoker {
	return &ddbStreamsInvoker{svc: svc}
}

func (i *ddbStreamsInvoker) GetRecords(ctx context.Context, region, tableName string, fromSeq int64, limit int) ([]invokers.DynamoDBStreamRecord, int64, error) {
	store, err := i.svc.GetCachedStoreForRegion(region)
	if err != nil {
		return nil, 0, err
	}

	records, nextSeq, err := store.Streams().GetRecords(tableName, fromSeq, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]invokers.DynamoDBStreamRecord, len(records))
	for idx, rec := range records {
		dynamodbMap := map[string]interface{}{
			"Keys":                        rec.Dynamodb.Keys,
			"SequenceNumber":              rec.Dynamodb.SequenceNumber,
			"SizeBytes":                   rec.Dynamodb.SizeBytes,
			"StreamViewType":              rec.Dynamodb.StreamViewType,
			"ApproximateCreationDateTime": rec.Dynamodb.ApproximateCreationDateTime,
		}
		if rec.Dynamodb.NewImage != nil {
			dynamodbMap["NewImage"] = rec.Dynamodb.NewImage
		}
		if rec.Dynamodb.OldImage != nil {
			dynamodbMap["OldImage"] = rec.Dynamodb.OldImage
		}

		result[idx] = invokers.DynamoDBStreamRecord{
			EventID:        rec.EventID,
			EventName:      string(rec.EventName),
			EventVersion:   rec.EventVersion,
			EventSource:    rec.EventSource,
			AWSRegion:      rec.AWSRegion,
			Dynamodb:       dynamodbMap,
			EventSourceARN: rec.EventSourceARN,
		}
	}
	return result, nextSeq, nil
}

func (i *ddbStreamsInvoker) GetLatestSequence(ctx context.Context, region, tableName string) (int64, error) {
	store, err := i.svc.GetCachedStoreForRegion(region)
	if err != nil {
		return 0, err
	}
	return store.Streams().GetLatestSequence(tableName)
}

// ShardIDForStream returns the deterministic shard identifier of a DynamoDB
// stream ARN, the same value the DynamoDB Streams API reports for the
// stream. Lambda's tumbling-window events carry it as the envelope
// shardId.
func (i *ddbStreamsInvoker) ShardIDForStream(streamARN string) string {
	return dbstore.ShardIDForStream(streamARN)
}
