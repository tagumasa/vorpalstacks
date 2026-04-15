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

// cleanExpiredKeys scans for expired TTL entries using a snapshot and
// deletes them in a single batch. The write lock is held only during the
// batch commit, not during the full scan.
func (d *DB) cleanExpiredKeys() {
	d.mu.RLock()
	if d.closed {
		d.mu.RUnlock()
		return
	}

	snap := d.db.NewSnapshot()
	d.mu.RUnlock()

	keysToDelete := d.collectExpiredKeys(snap)
	snap.Close()

	if len(keysToDelete) == 0 {
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

	if err := batch.Commit(d.syncFlag()); err != nil {
		slog.Warn("failed to commit batch deletion of expired keys", "error", err, "count", len(keysToDelete))
	}
}

// collectExpiredKeys iterates a snapshot and returns keys whose TTL
// values have expired. The caller is responsible for closing the snapshot.
func (d *DB) collectExpiredKeys(snap *pebble.Snapshot) [][]byte {
	iter, err := snap.NewIter(nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var keys [][]byte
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
		keys = append(keys, keyCopy)
	}

	return keys
}
