package pebbledb

import (
	"bytes"
	"log/slog"
	"time"

	"github.com/cockroachdb/pebble/v2"
)

func (d *DB) ttlGC() {
	ticker := time.NewTicker(d.opts.TTL.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-d.ttlStop:
			return
		case <-ticker.C:
			d.cleanExpiredKeys()
		}
	}
}

func (d *DB) cleanExpiredKeys() {
	d.mu.RLock()
	if d.closed {
		d.mu.RUnlock()
		return
	}

	now := time.Now().Unix()
	lower := []byte(ttlIndexPrefix)

	snap := d.db.NewSnapshot()

	var keysToDelete [][]byte
	var indexKeysToDelete [][]byte

	iter, err := snap.NewIter(&pebble.IterOptions{
		LowerBound: lower,
	})
	if err != nil {
		snap.Close()
		d.mu.RUnlock()
		return
	}

	for iter.First(); iter.Valid(); iter.Next() {
		k := iter.Key()
		if !bytes.HasPrefix(k, lower) {
			break
		}
		expiresAt, dataKey, ok := parseTTLIndexKey(k)
		if !ok {
			continue
		}
		if expiresAt >= now {
			break
		}
		idxCopy := make([]byte, len(k))
		copy(idxCopy, k)
		indexKeysToDelete = append(indexKeysToDelete, idxCopy)
		dataCopy := make([]byte, len(dataKey))
		copy(dataCopy, dataKey)
		keysToDelete = append(keysToDelete, dataCopy)
	}
	iterErr := iter.Error()
	iter.Close()
	snap.Close()
	d.mu.RUnlock()

	if iterErr != nil {
		slog.Warn("TTL GC index scan failed", "error", iterErr)
		return
	}

	if len(indexKeysToDelete) == 0 {
		return
	}

	d.mu.Lock()
	if d.closed {
		d.mu.Unlock()
		return
	}
	defer d.mu.Unlock()

	batch := d.db.NewBatch()
	defer batch.Close()

	for _, key := range keysToDelete {
		_ = batch.Delete(key, nil)
	}
	for _, key := range indexKeysToDelete {
		_ = batch.Delete(key, nil)
	}

	if err := batch.Commit(d.syncFlag()); err != nil {
		slog.Warn("failed to commit batch deletion of expired keys", "error", err, "count", len(keysToDelete))
	}
}
