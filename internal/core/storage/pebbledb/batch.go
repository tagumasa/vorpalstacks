package pebbledb

import (
	"errors"

	"github.com/cockroachdb/pebble/v2"
)

var (
	// ErrBatchClosed is returned when attempting to operate on a closed batch.
	ErrBatchClosed = errors.New("batch is closed")
)

// Batch provides a batched write operation for PebbleDB.
// It allows multiple writes to be grouped and committed atomically.
type Batch struct {
	db        *DB
	batch     *pebble.Batch
	encryptor Encryptor
	closed    bool
}

// NewBatch creates a new Batch for the database.
// The batch inherits the database's encryption settings.
func (d *DB) NewBatch() *Batch {
	return &Batch{
		db:        d,
		batch:     d.db.NewBatch(),
		encryptor: d.encryptor,
	}
}

// Set adds a key-value pair to the batch.
// The value is encrypted if encryption is enabled on the database.
func (b *Batch) Set(key, value []byte) error {
	if b.closed {
		return ErrBatchClosed
	}
	encrypted, err := b.encryptor.Encrypt(value)
	if err != nil {
		return err
	}
	return b.batch.Set(key, encrypted, nil)
}

// Delete adds a key deletion to the batch.
func (b *Batch) Delete(key []byte) error {
	if b.closed {
		return ErrBatchClosed
	}
	return b.batch.Delete(key, nil)
}

// DeleteRange adds a range deletion to the batch.
// All keys in the range [start, end) are marked for deletion.
func (b *Batch) DeleteRange(start, end []byte) error {
	if b.closed {
		return ErrBatchClosed
	}
	return b.batch.DeleteRange(start, end, nil)
}

// Commit writes all pending operations in the batch to the database.
// After committing, the batch is closed and cannot be used further.
func (b *Batch) Commit(opts *pebble.WriteOptions) error {
	if b.closed {
		return ErrBatchClosed
	}
	b.closed = true
	if opts == nil {
		opts = &pebble.WriteOptions{Sync: false}
	}
	err := b.batch.Commit(opts)
	b.batch.Close()
	return err
}

// Close closes the batch without committing.
// Any pending operations are discarded.
func (b *Batch) Close() error {
	if b.closed {
		return nil
	}
	b.closed = true
	return b.batch.Close()
}
