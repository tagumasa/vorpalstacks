package dynamodb

import (
	"context"
	"strconv"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
)

// ttlWorker periodically scans tables for expired TTL items and deletes them.
// AWS typically deletes expired items within 48 hours; this worker runs on a
// shorter interval for edge/on-prem deployments where timely cleanup matters.
type ttlWorker struct {
	store  *DynamoDBStore
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// newTTLWorker creates and starts a TTL cleanup goroutine for the given store.
func newTTLWorker(store *DynamoDBStore) *ttlWorker {
	ctx, cancel := context.WithCancel(context.Background())
	w := &ttlWorker{
		store:  store,
		ctx:    ctx,
		cancel: cancel,
	}
	w.wg.Add(1)
	go w.run()
	return w
}

// Close stops the TTL worker and waits for the goroutine to exit.
func (w *ttlWorker) Close() {
	w.cancel()
	w.wg.Wait()
}

// ttlSweepInterval is the expiry sweep's cadence: expired items are
// removed at most this often per worker, trading deletion latency for
// scan cost.
const ttlSweepInterval = 5 * time.Minute

func (w *ttlWorker) run() {
	defer w.wg.Done()
	ticker := time.NewTicker(ttlSweepInterval)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.tick()
		}
	}
}

// tick runs one cleanup pass with crash containment. The worker holds no
// lock across a tick — every store call it makes completes its own lock,
// iterator and transaction lifetimes inside the call — so a panic inside
// the pass leaves no half-held state the next tick would deadlock on: the
// panic is logged and the ticker keeps ticking. Without the containment a
// single panicking pass would take the whole process down.
func (w *ttlWorker) tick() {
	defer func() {
		if r := recover(); r != nil {
			resilience.LogPanic("dynamodb ttl worker tick", r)
		}
	}()
	w.doTTLCleanup()
}

// doTTLCleanup scans all tables, finds those with TTL enabled, and deletes
// items whose TTL attribute value is in the past.
func (w *ttlWorker) doTTLCleanup() {
	now := time.Now().Unix()

	var marker string
	for {
		if w.ctx.Err() != nil {
			return
		}
		tables, nextMarker, err := w.store.tables.List(marker, 100)
		if err != nil {
			logs.Warn("DynamoDB TTL worker: failed to list tables",
				logs.Err(err))
			return
		}
		for _, table := range tables {
			if w.ctx.Err() != nil {
				return
			}
			ttl, err := w.store.tables.GetTimeToLive(table.Name)
			if err != nil || ttl == nil || !ttl.Enabled || ttl.AttributeName == "" {
				continue
			}
			w.cleanupTableTTL(table, ttl.AttributeName, now)
		}
		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}
}

// ttlExpired reports whether the TTL attribute holds a Number whose Unix
// epoch-seconds value is strictly before now. An absent, null or
// non-numeric attribute never expires. The scan filter and the
// in-transaction re-check share this predicate so both ask the same
// question of the item they are looking at. The comparison stays in the
// float domain: converting a Number beyond the int64 range (a valid
// DynamoDB Number far in the future) to an integer type is
// implementation-dependent — amd64 answers MinInt64, which would read as
// long expired — so the clock converts to float64 instead.
func ttlExpired(attributes map[string]*AttributeValue, ttlAttrName string, now int64) bool {
	attr, exists := attributes[ttlAttrName]
	if !exists || attr == nil || attr.N == nil {
		return false
	}
	expiry, err := strconv.ParseFloat(*attr.N, 64)
	if err != nil {
		return false
	}
	return expiry < float64(now)
}

// cleanupTableTTL scans a single table and deletes items whose TTL attribute
// (a Number holding a Unix epoch timestamp in seconds) is less than now.
func (w *ttlWorker) cleanupTableTTL(table *Table, ttlAttrName string, now int64) {
	var deletedCount int64

	err := w.store.items.Scan(table.Name, func(item *Item) error {
		if !ttlExpired(item.Attributes, ttlAttrName, now) {
			return nil
		}

		deleted, delErr := w.deleteIfStillExpired(table, ttlAttrName, now, item.Key)
		if delErr != nil {
			logs.Warn("DynamoDB TTL: failed to delete expired item",
				logs.String("table", table.Name),
				logs.Err(delErr))
		} else if deleted {
			deletedCount++
		}
		return nil
	})

	if err != nil {
		logs.Warn("DynamoDB TTL: scan error during cleanup",
			logs.String("table", table.Name),
			logs.Err(err))
	}
	if deletedCount > 0 {
		logs.Info("DynamoDB TTL: deleted expired items",
			logs.String("table", table.Name),
			logs.Int64("count", deletedCount))
	}
}

// deleteIfStillExpired deletes one item the scan judged expired and reports
// whether it did. The scan iterates a point-in-time snapshot while client
// writes keep committing, so the expiry decision is re-evaluated inside the
// transaction against the live record: an item whose TTL attribute was
// extended into the future or removed between the snapshot and this
// transaction is no longer expired and must survive the cleanup pass (the
// documented TTL contract — a write that lands before the delete keeps the
// item alive). The re-check uses the pass-scoped clock the scan used, so it
// asks exactly whether the item still belongs to this pass's expired set.
func (w *ttlWorker) deleteIfStillExpired(table *Table, ttlAttrName string, now int64, key map[string]*AttributeValue) (bool, error) {
	deleted := false
	err := w.store.Update(w.ctx, func(txn *DynamoDBTxn) error {
		existingItem, err := txn.GetItem(table.Name, key)
		if err != nil {
			if IsItemNotFound(err) {
				// A client deleted the item between the scan
				// snapshot and this transaction.
				return nil
			}
			return err
		}
		if !ttlExpired(existingItem.Attributes, ttlAttrName, now) {
			// The live record no longer expires under this pass's
			// decision clock: a concurrent write extended or
			// removed its TTL attribute, and that write survives.
			return nil
		}

		// Delete the expired item within a transaction so that index
		// entries, item count, and table size are updated atomically.
		// Both the delete and the stream record key off the re-read
		// record, so the whole removal describes one consistent item.
		if err := txn.DeleteItemWrite(table.Name, existingItem.Key, existingItem, true, CalculateItemSize(existingItem.Attributes)); err != nil {
			return err
		}

		// Emit stream record if streaming is enabled. A capture failure
		// aborts the transaction — the delete and its stream record are
		// one removal ("the whole removal describes one consistent
		// item"), so committing the delete without the record would
		// publish an item removal no stream consumer ever sees; the next
		// pass re-runs the whole removal, the same propagation the main
		// write paths apply to capture errors.
		if table.StreamSpecification != nil && table.StreamSpecification.StreamEnabled && table.StreamArn != "" {
			keysResp, _, oldImageResp := BuildStreamImages(
				table.StreamSpecification.StreamViewType,
				existingItem.Key, nil, existingItem.Attributes,
			)
			if _, streamErr := w.store.streams.AddRecordTxn(
				txn.RawTxn(),
				table.Name,
				table.StreamArn,
				string(table.StreamSpecification.StreamViewType),
				StreamEventRemove,
				keysResp,
				nil,
				oldImageResp,
				TTLServiceIdentity,
			); streamErr != nil {
				logs.Error("DynamoDB TTL: stream record capture failed, delete aborted for retry",
					logs.String("table", table.Name),
					logs.Err(streamErr))
				return streamErr
			}
		}
		deleted = true
		return nil
	})
	return deleted, err
}
