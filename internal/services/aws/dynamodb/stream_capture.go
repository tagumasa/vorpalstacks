// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/internal/utils/aws/arn"
)

// captureStreamChangeTxn writes a stream record within the same storage
// transaction as the item mutation, ensuring atomicity. If the transaction
// rolls back, the stream record is also discarded; a record write that
// fails carries the failure back, so the item write can never commit
// without the stream record its table's StreamSpecification promises.
func (s *DynamoDBService) captureStreamChangeTxn(txn *dbstore.DynamoDBTxn, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue) error {
	return s.captureStreamChangeTxnAs(txn, store, table, eventName, keys, newImage, oldImage, nil)
}

// captureStreamChangeTxnAs is captureStreamChangeTxn with the actor identity
// AWS attaches to service-initiated records; global-table replication marks
// its records with the Service principal so consumers can tell replicated
// writes from direct client writes.
func (s *DynamoDBService) captureStreamChangeTxnAs(txn *dbstore.DynamoDBTxn, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue, userIdentity *dbstore.StreamUserIdentity) error {
	if table == nil || table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled {
		return nil
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
	return err
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
	if precision == dbstore.ACDTPrecisionMicrosecond {
		createdAt = now.UnixMicro()
	} else {
		precision = dbstore.ACDTPrecisionMillisecond
		createdAt = now.UnixMilli()
	}

	data := kinesisDestinationRecordData{
		Keys:                                 buildItemResponse(keys),
		ApproximateCreationDateTime:          createdAt,
		ApproximateCreationDateTimePrecision: string(precision),
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
// asynchronously after the storage transaction has committed. Each emit
// registers on the service's background WaitGroup: Close waits for the
// in-flight emits to finish, and the emit's own bounded timeout keeps that
// wait finite — a shutdown completes a half-delivered emit rather than
// killing it.
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

		s.bgWg.Add(1)
		go func(sn, region, pk, payload string) {
			defer s.bgWg.Done()
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
// The replication runs asynchronously in a goroutine registered on the
// service's background WaitGroup: Close waits for the in-flight delivery to
// finish, and the delivery's own bounded timeout keeps that wait finite — a
// shutdown completes a half-delivered replication rather than killing it and
// leaving the replicas diverged. A delivery that exhausts its bounded
// retries escalates to the convergence worker, whose in-memory queue dies
// with the process: escalations pending at shutdown are lost, and the
// divergence they carried stands until a later write to the same item
// re-replicates. The callback receives a context with a
// 30-second timeout and the destination store; it should use
// destStore.Update(ctx, ...) so that index entries, item count, and
// table size are updated atomically.
func (s *DynamoDBService) replicateToGlobalTableReplicas(sourceStore dbstore.DynamoDBStoreInterface, sourceRegion, tableName string, op func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error) {
	if s.storageManager == nil {
		return
	}

	globalTable, err := sourceStore.GlobalTables().Get(tableName)
	if err != nil {
		// A missing record is the ordinary case: the table simply does not
		// belong to a global table. Any other read failure would silently
		// skip replication and let the replicas diverge, so it is logged
		// rather than swallowed.
		if !commonstore.IsNotFound(err) {
			logs.Warn("DynamoDB Global Tables: failed to read the global-table record; change not replicated",
				logs.String("table", tableName),
				logs.String("sourceRegion", sourceRegion),
				logs.Err(err))
		}
		return
	}
	if globalTable == nil || len(globalTable.ReplicationGroup) <= 1 {
		return
	}

	s.bgWg.Add(1)
	go func() {
		defer s.bgWg.Done()
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
			// One delivery per replica with bounded retries: a transient
			// fault — store resolution or the write itself — heals inside
			// the same delivery, and a persistent one escalates to the
			// convergence worker rather than ending as a logged
			// divergence. The writes are idempotent (last-write-wins
			// puts, absent-key deletes), so a retried delivery cannot
			// double-apply.
			deliverErr := withReplicationRetry(ctx, func() error {
				destStore, err := s.GetStoreForRegion(replica.RegionName)
				if err != nil {
					return err
				}
				return op(ctx, destStore)
			})
			if deliverErr != nil {
				s.escalateReplicationDelivery(tableName, replica.RegionName, op, deliverErr)
			}
		}
	}()
}

// replicationEscalation carries one exhausted delivery to the convergence
// worker: the destination region and the op itself, which closes over the
// change (table record snapshot, key, attributes) and re-runs unchanged.
type replicationEscalation struct {
	tableName string
	region    string
	op        func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error
}

// replicationEscalationQueueCap bounds the pending escalations: each entry
// holds one change's data, so the convergence path's memory stays bounded
// — an overflowing queue drops the escalation with an Error log, the
// disclosed limit of the stand-in.
const replicationEscalationQueueCap = 1024

// escalateReplicationDelivery hands an exhausted per-replica delivery to
// the convergence worker, creating the worker and its queue lazily under
// the escalation mutex on first use (the replication goroutines run on
// services that do not all pass through NewDynamoDBService, so the worker
// cannot start there alone). The non-blocking send keeps the replication
// goroutine delivery-paced; a full queue is the drop the cap discloses.
func (s *DynamoDBService) escalateReplicationDelivery(tableName, region string, op func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error, cause error) {
	logs.Warn("DynamoDB Global Tables: delivery exhausted its retries; escalating for background convergence",
		logs.String("table", tableName),
		logs.String("region", region),
		logs.Err(cause))
	s.escalationMu.Lock()
	queue := s.replicationEscalations
	if queue == nil {
		// The worker's lifetime is the service's background context — the
		// same lifecycle every other background goroutine derives from —
		// so Close's single cancellation reaches it however late the
		// delivery creates it; the delivery goroutine's own wait-group
		// entry holds the counter open until the worker registers.
		ctx, _ := s.ensureBackgroundLifecycle()
		queue = make(chan replicationEscalation, replicationEscalationQueueCap)
		s.replicationEscalations = queue
		s.bgWg.Add(1)
		go s.runReplicationEscalationWorker(ctx, queue)
	}
	s.escalationMu.Unlock()
	s.escalationsEnqueued.Add(1)
	select {
	case queue <- replicationEscalation{tableName: tableName, region: region, op: op}:
	default:
		logs.Error("DynamoDB Global Tables: escalation queue full; divergence stands",
			logs.String("table", tableName),
			logs.String("region", region))
	}
}

// runReplicationEscalationWorker drains the escalation queue until the
// cancellation Close triggers: every failed delivery returns to the
// queue's tail, so convergence keeps trying while the fault stands.
func (s *DynamoDBService) runReplicationEscalationWorker(ctx context.Context, queue chan replicationEscalation) {
	defer s.bgWg.Done()
	defer func() {
		if r := recover(); r != nil {
			resilience.LogPanic("dynamodb replication escalation worker", r)
		}
	}()
	for {
		select {
		case e := <-queue:
			if s.runEscalatedDelivery(e) {
				// The fault outlived this delivery's bounded retries:
				// the escalation returns to the queue's tail and
				// convergence keeps trying.
				s.escalationsEnqueued.Add(1)
				select {
				case queue <- e:
				default:
					logs.Error("DynamoDB Global Tables: escalation queue full; divergence stands",
						logs.String("table", e.tableName),
						logs.String("region", e.region))
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// pendingEscalationCount reports how many exhausted deliveries the
// convergence path has enqueued, counting the worker's own re-queues — the
// monotonic observation point a test waits on to know a delivery has
// escalated (the queue's instantaneous length oscillates as the worker
// drains it).
func (s *DynamoDBService) pendingEscalationCount() int64 {
	return s.escalationsEnqueued.Load()
}

// runEscalatedDelivery re-runs one escalated delivery under a fresh
// bounded timeout with the same retry shape the in-flight delivery used,
// and reports whether the delivery failed again (the caller re-queues it).
// A panic inside the op drops the escalation: the recovery log discloses
// the loss, and a panicking op would fail every future attempt alike.
func (s *DynamoDBService) runEscalatedDelivery(e replicationEscalation) (failedAgain bool) {
	defer func() {
		if r := recover(); r != nil {
			resilience.LogPanic("dynamodb replication escalation delivery", r)
			failedAgain = false
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	deliverErr := withReplicationRetry(ctx, func() error {
		destStore, err := s.GetStoreForRegion(e.region)
		if err != nil {
			return err
		}
		return e.op(ctx, destStore)
	})
	if deliverErr == nil {
		logs.Info("DynamoDB Global Tables: escalated delivery converged",
			logs.String("table", e.tableName),
			logs.String("region", e.region))
		return false
	}
	logs.Warn("DynamoDB Global Tables: escalated delivery failed again; re-queued",
		logs.String("table", e.tableName),
		logs.String("region", e.region),
		logs.Err(deliverErr))
	return true
}

// replicateRetryAttempts bounds the per-replica delivery attempts. AWS's
// global tables converge through continuous internal retry; the platform's
// stand-in bounds the in-flight attempts to keep a delivery finite, and a
// fault that survives them escalates to the convergence worker rather
// than ending as a logged divergence.
const replicateRetryAttempts = 3

// withReplicationRetry runs one per-replica delivery with bounded retries:
// attempts separated by a short backoff that yields to the replication
// context's deadline.
func withReplicationRetry(ctx context.Context, run func() error) error {
	var err error
	for attempt := 0; attempt < replicateRetryAttempts; attempt++ {
		if attempt > 0 {
			timer := time.NewTimer(time.Duration(attempt) * 100 * time.Millisecond)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
		if err = run(); err == nil {
			return nil
		}
	}
	return err
}

// replicaPutOp returns a replica-region operation that stores a committed
// item together with its index entries and table counters, mirroring the
// write path of putItemCore. When the replica table itself streams, the
// replicated write captures a stream record in the replica's own stream,
// marked with the replication service identity; a replicated write never
// re-replicates or delivers to the replica's Kinesis destinations.
func (s *DynamoDBService) replicaPutOp(table *dbstore.Table, key, attrs map[string]*dbstore.AttributeValue) func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
	return func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
		return destStore.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			existing, getErr := txn.GetItem(table.Name, key)
			if getErr != nil && !dbstore.IsItemNotFound(getErr) {
				return getErr
			}
			var oldItem *dbstore.Item
			var oldSize int64
			if getErr == nil && existing != nil {
				oldItem = existing
				oldSize = dbstore.CalculateItemSize(existing.Attributes)
			}
			if err := txn.StoreItemWrite(table.Name, key, attrs, oldItem, oldItem != nil, oldSize); err != nil {
				return err
			}

			// The replica table's absence skips the capture (a replica
			// without the table has no stream to record into), but a
			// read fault aborts the transaction: the item write and its
			// stream record are one replica update, so committing the
			// write while silently dropping the record would diverge
			// the replica's stream from its own state — the enclosing
			// retry re-runs both halves.
			repTable, tblErr := txn.GetTable(table.Name)
			if tblErr != nil && !dbstore.IsTableNotFound(tblErr) {
				return tblErr
			}
			if tblErr == nil {
				var oldAttrs map[string]*dbstore.AttributeValue
				if oldItem != nil {
					oldAttrs = oldItem.Attributes
				}
				if err := s.captureStreamChangeTxnAs(txn, destStore, repTable, streamEventForWrite(false, oldItem == nil), key, attrs, oldAttrs, dbstore.ReplicationServiceIdentity); err != nil {
					return err
				}
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
			existingSize := dbstore.CalculateItemSize(existing.Attributes)
			if err := txn.DeleteItemWrite(table.Name, key, existing, true, existingSize); err != nil {
				return err
			}

			// The same fault rule the put half applies: absence skips the
			// capture, a read fault aborts so the removal and its
			// stream record commit or roll back together.
			repTable, tblErr := txn.GetTable(table.Name)
			if tblErr != nil && !dbstore.IsTableNotFound(tblErr) {
				return tblErr
			}
			if tblErr == nil {
				if err := s.captureStreamChangeTxnAs(txn, destStore, repTable, dbstore.StreamEventRemove, key, nil, existing.Attributes, dbstore.ReplicationServiceIdentity); err != nil {
					return err
				}
			}
			return nil
		})
	}
}
