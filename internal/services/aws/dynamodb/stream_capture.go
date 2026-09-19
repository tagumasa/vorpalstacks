// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/internal/utils/aws/arn"
)

// captureStreamChangeTxn writes a stream record within the same storage
// transaction as the item mutation, ensuring atomicity. If the transaction
// rolls back, the stream record is also discarded.
func (s *DynamoDBService) captureStreamChangeTxn(txn *dbstore.DynamoDBTxn, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue) {
	s.captureStreamChangeTxnAs(txn, store, table, eventName, keys, newImage, oldImage, nil)
}

// captureStreamChangeTxnAs is captureStreamChangeTxn with the actor identity
// AWS attaches to service-initiated records; global-table replication marks
// its records with the Service principal so consumers can tell replicated
// writes from direct client writes.
func (s *DynamoDBService) captureStreamChangeTxnAs(txn *dbstore.DynamoDBTxn, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue, userIdentity *dbstore.StreamUserIdentity) {
	if table == nil || table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled {
		return
	}

	streamViewType := table.StreamSpecification.StreamViewType
	streamArn := table.StreamArn

	keysResp, newImageResp, oldImageResp := dbstore.BuildStreamImages(streamViewType, keys, newImage, oldImage)

	_, err := store.Streams().AddRecordTxn(
		txn.RawTxn(),
		table.Name,
		streamArn,
		string(streamViewType),
		eventName,
		keysResp,
		newImageResp,
		oldImageResp,
		userIdentity,
	)
	if err != nil {
		logs.Error("failed to capture stream record in transaction",
			logs.String("table", table.Name), logs.Err(err))
	}
}

// streamEventForWrite maps a committed write's shape to its DynamoDB
// Streams event name: REMOVE for deletions, INSERT for creations, MODIFY for
// every other write.
func streamEventForWrite(isDelete, wasNew bool) dbstore.StreamEventName {
	if isDelete {
		return dbstore.StreamEventRemove
	}
	if wasNew {
		return dbstore.StreamEventInsert
	}
	return dbstore.StreamEventModify
}

// emitChangePropagation fires the asynchronous side effects every committed
// item change owes: delivery to the table's Kinesis data stream destinations
// and replication to the global table's replica regions. It runs after the
// storage transaction has committed; replicaOp builds the destination-region
// write for the change's post-image (or removal), and a nil table or replicaOp
// leaves only the parts the caller actually has.
func (s *DynamoDBService) emitChangePropagation(store dbstore.DynamoDBStoreInterface, region string, table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue, replicaOp func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error) {
	s.sendToKinesisDestinations(table, eventName, keys, newImage, oldImage)
	if table == nil || replicaOp == nil {
		return
	}
	s.replicateToGlobalTableReplicas(store, region, table.Name, replicaOp)
}

// kinesisDestinationEnvelope is the JSON document DynamoDB writes into a
// Kinesis data stream destination for one item change: the DynamoDB Streams
// record's dynamodb envelope carrying the key and images, with the timestamp
// precision the destination is configured for.
type kinesisDestinationEnvelope struct {
	EventName string                       `json:"eventName"`
	DynamoDB  kinesisDestinationRecordData `json:"dynamodb"`
}

type kinesisDestinationRecordData struct {
	Keys                                 map[string]interface{} `json:"Keys,omitempty"`
	NewImage                             map[string]interface{} `json:"NewImage,omitempty"`
	OldImage                             map[string]interface{} `json:"OldImage,omitempty"`
	ApproximateCreationDateTime          int64                  `json:"ApproximateCreationDateTime"`
	ApproximateCreationDateTimePrecision string                 `json:"ApproximateCreationDateTimePrecision"`
}

// kinesisRecordForDestination builds one destination's payload and Kinesis
// partition key. ApproximateCreationDateTimePrecision is a per-destination
// setting (UpdateKinesisStreamingDestination), so the timestamp is formatted
// here for this destination only — MILLISECOND unless the destination asks
// for MICROSECOND. The partition key comes from the table's HASH attribute:
// AWS does not document the derivation, and the HASH value keeps every
// record of one item on the same Kinesis shard.
func kinesisRecordForDestination(dest *dbstore.KinesisDataStreamDestination, table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue, now time.Time) ([]byte, string) {
	precision := dest.ApproximateCreationDateTimePrecision
	var createdAt int64
	if precision == kinesisPrecisionMicrosecond {
		createdAt = now.UnixMicro()
	} else {
		precision = kinesisPrecisionMillisecond
		createdAt = now.UnixMilli()
	}

	data := kinesisDestinationRecordData{
		Keys:                                 buildItemResponse(keys),
		ApproximateCreationDateTime:          createdAt,
		ApproximateCreationDateTimePrecision: precision,
	}
	if newImage != nil {
		data.NewImage = buildItemResponse(newImage)
	}
	if oldImage != nil {
		data.OldImage = buildItemResponse(oldImage)
	}

	// The payload is pure JSON data, so marshalling cannot fail.
	payload, _ := json.Marshal(kinesisDestinationEnvelope{
		EventName: string(eventName),
		DynamoDB:  data,
	})
	return payload, extractPartitionKeyForKinesis(keys, getHashKeyName(table))
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

	now := time.Now()
	for _, dest := range table.KinesisDataStreamDestinations {
		if dest.DestinationStatus != kinesisDestinationActive {
			continue
		}
		streamName := arn.ExtractStreamNameFromARN(dest.StreamArn)
		if streamName == "" {
			continue
		}
		payload, partitionKey := kinesisRecordForDestination(dest, table, eventName, keys, newImage, oldImage, now)
		// The Kinesis store keeps record payloads in their wire representation
		// (base64 text), matching records written through the Kinesis API, so
		// GetRecords consumers decode every record the same way.
		wireData := base64.StdEncoding.EncodeToString(payload)
		_, _, destRegion, _, _ := arn.SplitARN(dest.StreamArn)

		go func(sn, region, pk, payload string) {
			defer func() {
				if r := recover(); r != nil {
					resilience.LogPanic("dynamodb Kinesis destination emit", r)
				}
			}()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if _, _, err := kinesisInvoker.PutRecord(ctx, region, sn, pk, []byte(payload)); err != nil {
				logs.Warn("failed to send record to Kinesis destination",
					logs.String("stream", sn), logs.Err(err))
			}
		}(streamName, destRegion, partitionKey, wireData)
	}
}

// extractPartitionKeyForKinesis renders the item's HASH key attribute as the
// Kinesis partition key. AWS does not document the derivation; using the
// HASH attribute's value is deterministic per item (map iteration order
// never selects it), and anything the HASH attribute cannot render falls
// back to a constant key.
func extractPartitionKeyForKinesis(keys map[string]*dbstore.AttributeValue, pkName string) string {
	if pkName != "" {
		if v := keys[pkName]; v != nil {
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

// replicateToGlobalTableReplicas propagates item changes to all other
// replica regions if the table is part of a global table. This implements
// the multi-active replication behaviour of DynamoDB Global Tables.
//
// The replication runs asynchronously in a goroutine. The callback receives
// a context with a 30-second timeout and the destination store; it should
// use destStore.Update(ctx, ...) so that index entries, item count, and
// table size are updated atomically.
func (s *DynamoDBService) replicateToGlobalTableReplicas(sourceStore dbstore.DynamoDBStoreInterface, sourceRegion, tableName string, op func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error) {
	if s.storageManager == nil {
		return
	}

	globalTable, err := sourceStore.GlobalTables().Get(tableName)
	if err != nil || globalTable == nil || len(globalTable.ReplicationGroup) <= 1 {
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				resilience.LogPanic("dynamodb global table replication", r)
			}
		}()
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
// write path of putItemCore. When the replica table itself streams, the
// replicated write captures a stream record in the replica's own stream,
// marked with the replication service identity; a replicated write never
// re-replicates or delivers to the replica's Kinesis destinations.
func (s *DynamoDBService) replicaPutOp(table *dbstore.Table, key, attrs map[string]*dbstore.AttributeValue) func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
	itemSize := dbstore.CalculateItemSize(attrs)
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
				oldSize := dbstore.CalculateItemSize(existing.Attributes)
				if upErr := txn.UpdateTableSize(table.Name, itemSize-oldSize); upErr != nil {
					return upErr
				}
			}

			if repTable, tblErr := txn.GetTable(table.Name); tblErr == nil {
				var oldAttrs map[string]*dbstore.AttributeValue
				if !isNewRep && existing != nil {
					oldAttrs = existing.Attributes
				}
				s.captureStreamChangeTxnAs(txn, destStore, repTable, streamEventForWrite(false, isNewRep), key, attrs, oldAttrs, dbstore.ReplicationServiceIdentity)
			}
			return nil
		})
	}
}

// replicaDeleteOp returns a replica-region operation that removes a
// committed item together with its index entries and table counters,
// mirroring the write path of deleteItemCore, and captures the removal in
// the replica table's own stream when it streams, marked with the
// replication service identity.
func (s *DynamoDBService) replicaDeleteOp(table *dbstore.Table, key map[string]*dbstore.AttributeValue) func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
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
			existingSize := dbstore.CalculateItemSize(existing.Attributes)
			if upErr := txn.UpdateTableSize(table.Name, -existingSize); upErr != nil {
				return upErr
			}

			if repTable, tblErr := txn.GetTable(table.Name); tblErr == nil {
				s.captureStreamChangeTxnAs(txn, destStore, repTable, dbstore.StreamEventRemove, key, nil, existing.Attributes, dbstore.ReplicationServiceIdentity)
			}
			return nil
		})
	}
}
