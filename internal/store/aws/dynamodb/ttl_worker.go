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
			if err := txn.DeleteIndexEntries(table.Name, existingItem); err != nil {
				return err
			}
			if err := txn.DeleteItem(table.Name, item.Key); err != nil {
				return err
			}
			if err := txn.UpdateItemCount(table.Name, -1); err != nil {
				return err
			}
			oldItemSize := calculateItemSizeForStore(existingItem.Attributes)
			if err := txn.UpdateTableSize(table.Name, -oldItemSize); err != nil {
				return err
			}

			// Emit stream record if streaming is enabled.
			if table.StreamSpecification != nil && table.StreamSpecification.StreamEnabled && table.StreamArn != "" {
				keysResp, _, oldImageResp := buildStreamImagesForStore(
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

// calculateItemSizeForStore estimates the storage size of an item's attributes.
func calculateItemSizeForStore(attrs map[string]*AttributeValue) int64 {
	var size int64
	for k, v := range attrs {
		size += int64(len(k))
		size += attributeValueSizeForStore(v)
	}
	return size
}

func attributeValueSizeForStore(v *AttributeValue) int64 {
	if v == nil {
		return 0
	}
	if v.S != nil {
		return int64(len(*v.S))
	}
	if v.N != nil {
		return int64(len(*v.N))
	}
	if v.B != nil {
		return int64(len(v.B))
	}
	if v.BOOL != nil {
		return 1
	}
	if v.NULL != nil {
		return 1
	}
	if v.SS != nil {
		var s int64
		for _, str := range v.SS {
			s += int64(len(str))
		}
		return s
	}
	if v.NS != nil {
		var s int64
		for _, n := range v.NS {
			s += int64(len(n))
		}
		return s
	}
	if v.BS != nil {
		var s int64
		for _, b := range v.BS {
			s += int64(len(b))
		}
		return s
	}
	if v.M != nil {
		var s int64
		for k, val := range v.M {
			s += int64(len(k)) + attributeValueSizeForStore(val)
		}
		return s
	}
	if v.L != nil {
		var s int64
		for _, val := range v.L {
			s += attributeValueSizeForStore(val)
		}
		return s
	}
	return 0
}

// buildStreamImagesForStore converts attribute values into the response format
// expected by stream consumers, filtered by the StreamViewType. This is a
// store-layer helper used by the TTL worker.
func buildStreamImagesForStore(streamViewType StreamViewType, keys, newImage, oldImage map[string]*AttributeValue) (keysResp, newImageResp, oldImageResp map[string]interface{}) {
	keysResp = make(map[string]interface{})
	for k, v := range keys {
		keysResp[k] = buildStreamAttributeValueResponse(v)
	}

	switch streamViewType {
	case StreamViewTypeNewImage:
		if newImage != nil {
			newImageResp = make(map[string]interface{})
			for k, v := range newImage {
				newImageResp[k] = buildStreamAttributeValueResponse(v)
			}
		}
	case StreamViewTypeOldImage:
		if oldImage != nil {
			oldImageResp = make(map[string]interface{})
			for k, v := range oldImage {
				oldImageResp[k] = buildStreamAttributeValueResponse(v)
			}
		}
	case StreamViewTypeNewAndOldImages:
		if newImage != nil {
			newImageResp = make(map[string]interface{})
			for k, v := range newImage {
				newImageResp[k] = buildStreamAttributeValueResponse(v)
			}
		}
		if oldImage != nil {
			oldImageResp = make(map[string]interface{})
			for k, v := range oldImage {
				oldImageResp[k] = buildStreamAttributeValueResponse(v)
			}
		}
	case StreamViewTypeKeysOnly:
		// Only keys are included.
	}

	return keysResp, newImageResp, oldImageResp
}

// buildStreamAttributeValueResponse converts a single AttributeValue to the
// DynamoDB JSON wire format used in stream records.
func buildStreamAttributeValueResponse(v *AttributeValue) map[string]interface{} {
	if v == nil {
		return nil
	}
	if v.S != nil {
		return map[string]interface{}{"S": *v.S}
	}
	if v.N != nil {
		return map[string]interface{}{"N": *v.N}
	}
	if v.B != nil {
		return map[string]interface{}{"B": v.B}
	}
	if v.BOOL != nil {
		return map[string]interface{}{"BOOL": *v.BOOL}
	}
	if v.NULL != nil {
		return map[string]interface{}{"NULL": *v.NULL}
	}
	if v.SS != nil {
		items := make([]interface{}, len(v.SS))
		for i, s := range v.SS {
			items[i] = s
		}
		return map[string]interface{}{"SS": items}
	}
	if v.NS != nil {
		items := make([]interface{}, len(v.NS))
		for i, n := range v.NS {
			items[i] = n
		}
		return map[string]interface{}{"NS": items}
	}
	if v.BS != nil {
		items := make([]interface{}, len(v.BS))
		for i, b := range v.BS {
			items[i] = b
		}
		return map[string]interface{}{"BS": items}
	}
	if v.M != nil {
		m := make(map[string]interface{})
		for k, val := range v.M {
			m[k] = buildStreamAttributeValueResponse(val)
		}
		return map[string]interface{}{"M": m}
	}
	if v.L != nil {
		items := make([]interface{}, len(v.L))
		for i, val := range v.L {
			items[i] = buildStreamAttributeValueResponse(val)
		}
		return map[string]interface{}{"L": items}
	}
	return nil
}
