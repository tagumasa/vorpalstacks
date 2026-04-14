package pebbledb

import (
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
	d.mu.Lock()
	if d.closed {
		d.mu.Unlock()
		return
	}

	keysToDelete := [][]byte{}

	iter, err := d.db.NewIter(nil)
	if err != nil {
		d.mu.Unlock()
		return
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()

		val, err := iter.ValueAndErr()
		if err != nil {
			continue
		}

		_, expired, err := decryptAndUnwrapTTL(d.encryptor, val, d.opts.TTL.Enabled)
		if err != nil || !expired {
			continue
		}

		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)
		keysToDelete = append(keysToDelete, keyCopy)
	}
	d.mu.Unlock()

	if len(keysToDelete) > 0 {
		batch := d.db.NewBatch()
		defer batch.Close()

		for _, key := range keysToDelete {
			_ = batch.Delete(key, nil) // intentional: batch errors are caught by Commit
		}

		if err := batch.Commit(pebble.Sync); err != nil {
			slog.Warn("failed to commit batch deletion of expired keys", "error", err, "count", len(keysToDelete))
		}
	}
}
