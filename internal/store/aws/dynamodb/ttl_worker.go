package dynamodb

import (
	"context"
	"strconv"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
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

func (w *ttlWorker) run() {
	defer w.wg.Done()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.doTTLCleanup()
		}
	}
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

// cleanupTableTTL scans a single table and deletes items whose TTL attribute
// (a Number holding a Unix epoch timestamp in seconds) is less than now.
func (w *ttlWorker) cleanupTableTTL(table *Table, ttlAttrName string, now int64) {
	var deletedCount int64

	err := w.store.items.Scan(table.Name, func(item *Item) error {
		attr, exists := item.Attributes[ttlAttrName]
		if !exists || attr == nil || attr.N == nil {
			return nil
		}
		expiry, err := strconv.ParseFloat(*attr.N, 64)
		if err != nil {
			return nil
		}
		if int64(expiry) >= now {
			return nil
		}

		// Delete the expired item within a transaction so that index
		// entries, item count, and table size are updated atomically.
		delErr := w.store.Update(w.ctx, func(txn *DynamoDBTxn) error {
			existingItem, err := txn.GetItem(table.Name, item.Key)
			if err != nil {
				if IsItemNotFound(err) {
					return nil
				}
				return err
			}
			if err := txn.DeleteItemWrite(table.Name, item.Key, existingItem, true, CalculateItemSize(existingItem.Attributes)); err != nil {
				return err
			}

			// Emit stream record if streaming is enabled.
			if table.StreamSpecification != nil && table.StreamSpecification.StreamEnabled && table.StreamArn != "" {
				keysResp, _, oldImageResp := BuildStreamImages(
					table.StreamSpecification.StreamViewType,
					item.Key, nil, existingItem.Attributes,
				)
				_, streamErr := w.store.streams.AddRecordTxn(
					txn.RawTxn(),
					table.Name,
					table.StreamArn,
					string(table.StreamSpecification.StreamViewType),
					StreamEventRemove,
					keysResp,
					nil,
					oldImageResp,
					TTLServiceIdentity,
				)
				if streamErr != nil {
					logs.Warn("DynamoDB TTL: failed to emit stream record",
						logs.String("table", table.Name),
						logs.Err(streamErr))
				}
			}
			return nil
		})
		if delErr != nil {
			logs.Warn("DynamoDB TTL: failed to delete expired item",
				logs.String("table", table.Name),
				logs.Err(delErr))
		} else {
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
