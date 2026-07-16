package storage

import (
	"vorpalstacks/internal/core/storage/pebbledb"

	"github.com/cockroachdb/pebble/v2"
)

// PebbleBatch implements Batch for a namespaced bucket backed by Pebble.
type PebbleBatch struct {
	batch  *pebbledb.Batch
	prefix []byte
}

func newPebbleBatch(db *pebbledb.DB, prefix []byte) *PebbleBatch {
	return &PebbleBatch{
		batch:  db.NewBatch(),
		prefix: prefix,
	}
}

func (b *PebbleBatch) makeKey(key []byte) []byte {
	return makePrefixedKey(b.prefix, key)
}

// Put writes a key-value pair to the batch.
func (b *PebbleBatch) Put(key, value []byte) error {
	return b.batch.Set(b.makeKey(key), value)
}

// Delete removes a key from the batch.
func (b *PebbleBatch) Delete(key []byte) error {
	return b.batch.Delete(b.makeKey(key))
}

// DeleteRange removes all keys in [start, end) from the batch.
func (b *PebbleBatch) DeleteRange(start, end []byte) error {
	return b.batch.DeleteRange(b.makeKey(start), b.makeKey(end))
}

// Commit writes the batch to storage without fsync.
func (b *PebbleBatch) Commit() error {
	return b.batch.Commit(pebble.NoSync)
}

// CommitSync writes the batch to storage and fsyncs the WAL.
func (b *PebbleBatch) CommitSync() error {
	return b.batch.Commit(pebble.Sync)
}

// Close releases batch resources without committing.
func (b *PebbleBatch) Close() error {
	return b.batch.Close()
}

// NewBatch creates an atomic batch for this bucket's key namespace.
func (b *PebbleBucket) NewBatch() Batch {
	return newPebbleBatch(b.db, b.prefix)
}

// DeleteRange removes all keys in [start, end) from the bucket.
func (b *PebbleBucket) DeleteRange(start, end []byte) error {
	return b.db.DeleteRange(b.makeKey(start), b.makeKey(end))
}
