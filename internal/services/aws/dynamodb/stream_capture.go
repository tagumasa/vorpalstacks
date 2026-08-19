// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// captureStreamChange generates a DynamoDB Streams record for the given
// item change if the table has streaming enabled. This is the legacy
// non-transactional version used only where a transaction is not available.
// Prefer captureStreamChangeTxn for new call sites.
func (s *DynamoDBService) captureStreamChange(ctx context.Context, reqCtx *request.RequestContext, table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue) {
	if table == nil || table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled {
		return
	}

	streamViewType := table.StreamSpecification.StreamViewType
	streamArn := table.StreamArn

	keysResp, newImageResp, oldImageResp := s.buildStreamImages(streamViewType, keys, newImage, oldImage)

	store, err := s.GetCachedStoreForRegion(reqCtx.GetRegion())
	if err != nil {
		logs.Warn("failed to get store for stream capture",
			logs.String("table", table.Name), logs.Err(err))
		return
	}

	_, err = store.Streams().AddRecord(
		table.Name,
		streamArn,
		string(streamViewType),
		eventName,
		keysResp,
		newImageResp,
		oldImageResp,
		nil,
	)
	if err != nil {
		logs.Warn("failed to capture stream record",
			logs.String("table", table.Name), logs.Err(err))
	}
}

// captureStreamChangeTxn writes a stream record within the same storage
// transaction as the item mutation, ensuring atomicity. If the transaction
// rolls back, the stream record is also discarded.
func (s *DynamoDBService) captureStreamChangeTxn(txn *dbstore.DynamoDBTxn, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue) {
	if table == nil || table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled {
		return
	}

	streamViewType := table.StreamSpecification.StreamViewType
	streamArn := table.StreamArn

	keysResp, newImageResp, oldImageResp := s.buildStreamImages(streamViewType, keys, newImage, oldImage)

	_, err := store.Streams().AddRecordTxn(
		txn.RawTxn(),
		table.Name,
		streamArn,
		string(streamViewType),
		eventName,
		keysResp,
		newImageResp,
		oldImageResp,
		nil,
	)
	if err != nil {
		logs.Error("failed to capture stream record in transaction",
			logs.String("table", table.Name), logs.Err(err))
	}
}

// buildStreamImages converts the raw attribute values into the response
// format expected by stream consumers, filtered by the StreamViewType.
func (s *DynamoDBService) buildStreamImages(streamViewType dbstore.StreamViewType, keys, newImage, oldImage map[string]*dbstore.AttributeValue) (keysResp, newImageResp, oldImageResp map[string]interface{}) {
	keysResp = buildItemResponse(keys)

	switch streamViewType {
	case dbstore.StreamViewTypeNewImage:
		if newImage != nil {
			newImageResp = buildItemResponse(newImage)
		}
	case dbstore.StreamViewTypeOldImage:
		if oldImage != nil {
			oldImageResp = buildItemResponse(oldImage)
		}
	case dbstore.StreamViewTypeNewAndOldImages:
		if newImage != nil {
			newImageResp = buildItemResponse(newImage)
		}
		if oldImage != nil {
			oldImageResp = buildItemResponse(oldImage)
		}
	case dbstore.StreamViewTypeKeysOnly:
		// Only keys are included.
	}

	return keysResp, newImageResp, oldImageResp
}

// kinesisDestinationRecord is the JSON payload sent to Kinesis Data Streams
// when a table has a Kinesis streaming destination configured.
type kinesisDestinationRecord struct {
	Keys                        map[string]interface{} `json:"Keys,omitempty"`
	NewImage                    map[string]interface{} `json:"NewImage,omitempty"`
	OldImage                    map[string]interface{} `json:"OldImage,omitempty"`
	EventName                   string                 `json:"eventName"`
	ApproximateCreationDateTime int64                  `json:"ApproximateCreationDateTime"`
}

// sendToKinesisDestinations dispatches item change records to all active
// Kinesis Data Stream destinations configured on the table. This runs
// asynchronously after the storage transaction has committed.
func (s *DynamoDBService) sendToKinesisDestinations(table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue) {
	if table == nil || len(table.KinesisDataStreamDestinations) == 0 {
		return
	}

	if s.bus == nil {
		return
	}
	kinesisInvoker := s.bus.KinesisInvoker()
	if kinesisInvoker == nil {
		return
	}

	keysResp := buildItemResponse(keys)
	// ApproximateCreationDateTime carries millisecond timestamps by default;
	// a destination configured for MICROSECOND precision receives
	// microsecond timestamps instead.
	now := time.Now()
	precision := kinesisPrecisionMillisecond
	for _, dest := range table.KinesisDataStreamDestinations {
		if dest.DestinationStatus == kinesisDestinationActive && dest.ApproximateCreationDateTimePrecision != "" {
			precision = dest.ApproximateCreationDateTimePrecision
			break
		}
	}
	var createdAt int64
	if precision == kinesisPrecisionMicrosecond {
		createdAt = now.UnixMicro()
	} else {
		createdAt = now.UnixMilli()
	}
	record := kinesisDestinationRecord{
		Keys:                        keysResp,
		EventName:                   string(eventName),
		ApproximateCreationDateTime: createdAt,
	}
	if newImage != nil {
		record.NewImage = buildItemResponse(newImage)
	}
	if oldImage != nil {
		record.OldImage = buildItemResponse(oldImage)
	}

	data, err := json.Marshal(record)
	if err != nil {
		logs.Warn("failed to marshal Kinesis destination record",
			logs.String("table", table.Name), logs.Err(err))
		return
	}
	// The Kinesis store keeps record payloads in their wire representation
	// (base64 text), matching records written through the Kinesis API, so
	// GetRecords consumers decode every record the same way.
	wireData := base64.StdEncoding.EncodeToString(data)

	partitionKey := extractPartitionKeyForKinesis(keys)

	for _, dest := range table.KinesisDataStreamDestinations {
		if dest.DestinationStatus != kinesisDestinationActive {
			continue
		}
		streamName := parseStreamNameFromARN(dest.StreamArn)
		if streamName == "" {
			continue
		}
		go func(sn, pk, payload string) {
			defer func() { resilience.RecoverPanic("dynamodb Kinesis destination emit") }()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if _, err := kinesisInvoker.PutRecord(ctx, sn, pk, []byte(payload)); err != nil {
				logs.Warn("failed to send record to Kinesis destination",
					logs.String("stream", sn), logs.Err(err))
			}
		}(streamName, partitionKey, wireData)
	}
}

// extractPartitionKeyForKinesis extracts a string partition key from the item's
// key attributes for use as the Kinesis partition key.
func extractPartitionKeyForKinesis(keys map[string]*dbstore.AttributeValue) string {
	for _, v := range keys {
		if v != nil {
			if v.S != nil {
				return *v.S
			}
			if v.N != nil {
				return *v.N
			}
		}
	}
	return "default"
}

// parseStreamNameFromARN extracts the stream name from a Kinesis ARN.
// Format: arn:aws:kinesis:<region>:<account>:stream/<name>
func parseStreamNameFromARN(arn string) string {
	idx := strings.LastIndex(arn, "/")
	if idx < 0 || idx+1 >= len(arn) {
		return ""
	}
	return arn[idx+1:]
}

// replicateToGlobalTableReplicas propagates item changes to all other
// replica regions if the table is part of a global table. This implements
// the multi-active replication behaviour of DynamoDB Global Tables.
//
// The replication runs asynchronously in a goroutine. The callback receives
// a context with a 30-second timeout and the destination store; it should
// use destStore.Update(ctx, ...) so that index entries, item count, and
// table size are updated atomically.
func (s *DynamoDBService) replicateToGlobalTableReplicas(sourceStore dbstore.DynamoDBStoreInterface, sourceRegion, tableName string, op func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error) {
	if s.busStoreFactory == nil {
		return
	}

	globalTable, err := sourceStore.GlobalTables().Get(tableName)
	if err != nil || globalTable == nil || len(globalTable.ReplicationGroup) <= 1 {
		return
	}

	go func() {
		defer func() { resilience.RecoverPanic("dynamodb global table replication") }()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		for _, replica := range globalTable.ReplicationGroup {
			if replica.RegionName == sourceRegion {
				continue
			}
			destStore, err := s.GetStoreForRegion(replica.RegionName)
			if err != nil {
				logs.Warn("DynamoDB Global Tables: failed to get store for replica region",
					logs.String("table", tableName),
					logs.String("region", replica.RegionName),
					logs.Err(err))
				continue
			}
			if err := op(ctx, destStore); err != nil {
				logs.Warn("DynamoDB Global Tables: failed to replicate to replica region",
					logs.String("table", tableName),
					logs.String("region", replica.RegionName),
					logs.Err(err))
			}
		}
	}()
}

// replicaPutOp returns a replica-region operation that stores a committed
// item together with its index entries and table counters, mirroring the
// write path of putItemCore.
func replicaPutOp(table *dbstore.Table, key, attrs map[string]*dbstore.AttributeValue) func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
	itemSize := calculateItemSize(attrs)
	return func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
		return destStore.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			existing, getErr := txn.GetItem(table.Name, key)
			isNewRep := false
			if getErr != nil {
				if dbstore.IsItemNotFound(getErr) {
					isNewRep = true
				} else {
					return getErr
				}
			}
			if !isNewRep && existing != nil {
				if delErr := txn.DeleteIndexEntries(table.Name, existing); delErr != nil {
					return delErr
				}
			}
			if putErr := txn.PutItem(table.Name, key, attrs); putErr != nil {
				return putErr
			}
			repItem := &dbstore.Item{
				TableName:  table.Name,
				Key:        key,
				Attributes: attrs,
			}
			if putIdxErr := txn.PutIndexEntries(table.Name, repItem); putIdxErr != nil {
				return putIdxErr
			}
			if isNewRep {
				if upErr := txn.UpdateItemCount(table.Name, 1); upErr != nil {
					return upErr
				}
				if upErr := txn.UpdateTableSize(table.Name, itemSize); upErr != nil {
					return upErr
				}
			} else {
				oldSize := calculateItemSize(existing.Attributes)
				if upErr := txn.UpdateTableSize(table.Name, itemSize-oldSize); upErr != nil {
					return upErr
				}
			}
			return nil
		})
	}
}

// replicaDeleteOp returns a replica-region operation that removes a
// committed item together with its index entries and table counters,
// mirroring the write path of deleteItemCore.
func replicaDeleteOp(table *dbstore.Table, key map[string]*dbstore.AttributeValue) func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
	return func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
		return destStore.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			existing, getErr := txn.GetItem(table.Name, key)
			if getErr != nil {
				if dbstore.IsItemNotFound(getErr) {
					return nil
				}
				return getErr
			}
			if delIdxErr := txn.DeleteIndexEntries(table.Name, existing); delIdxErr != nil {
				return delIdxErr
			}
			if delErr := txn.DeleteItem(table.Name, key); delErr != nil {
				return delErr
			}
			if upErr := txn.UpdateItemCount(table.Name, -1); upErr != nil {
				return upErr
			}
			existingSize := calculateItemSize(existing.Attributes)
			if upErr := txn.UpdateTableSize(table.Name, -existingSize); upErr != nil {
				return upErr
			}
			return nil
		})
	}
}
