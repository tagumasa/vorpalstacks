package pebbledb

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/cockroachdb/pebble/v2"
)

var (
	// ErrTxnClosed is returned when attempting to operate on a transaction that has been closed.
	ErrTxnClosed = errors.New("transaction is closed")
)

// Txn provides a transaction interface for the database.
// Transactions can be read-only (using snapshots) or read-write (using batches).
type Txn struct {
	db       *DB
	snapshot *pebble.Snapshot
	batch    *Batch
	readonly bool
	closed   bool
}

// NewSnapshot creates a new snapshot of the database.
// Snapshots provide a consistent view of the database at a point in time.
func (d *DB) NewSnapshot() *pebble.Snapshot {
	return d.db.NewSnapshot()
}

// View executes a read-only transaction against the database.
// The provided function receives a read-only transaction.
// Changes made within the function are not persisted.
func (d *DB) View(ctx context.Context, fn func(*Txn) error) error {
	snapshot := d.db.NewSnapshot()
	defer snapshot.Close()
	txn := &Txn{
		db:       d,
		snapshot: snapshot,
		readonly: true,
	}
	return fn(txn)
}

// Update executes a read-write transaction against the database.
// The provided function receives a read-write transaction.
// Changes are only persisted after the function returns without error.
func (d *DB) Update(ctx context.Context, fn func(*Txn) error) error {
	batch := d.NewBatch()
	defer batch.Close()
	txn := &Txn{
		db:       d,
		batch:    batch,
		readonly: false,
	}
	err := fn(txn)
	if err != nil {
		return err
	}
	return batch.Commit(pebble.Sync)
}

// Get retrieves a value from the transaction.
// For read-only transactions, it reads from the snapshot.
// For read-write transactions, it reads from the database.
func (t *Txn) Get(key []byte) ([]byte, error) {
	if t.closed {
		return nil, ErrTxnClosed
	}
	if t.readonly && t.snapshot != nil {
		return t.getAndDecrypt(t.snapshot, key)
	}
	if t.batch != nil {
		return t.getAndDecrypt(t.db.db, key)
	}
	return nil, ErrTxnClosed
}

// getAndDecrypt retrieves a value from a pebble.Getter, decrypts it, and
// unwraps the TTL envelope if TTL is enabled.
func (t *Txn) getAndDecrypt(getter interface {
	Get([]byte) ([]byte, io.Closer, error)
}, key []byte) ([]byte, error) {
	val, closer, err := getter.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	defer closer.Close()
	result := make([]byte, len(val))
	copy(result, val)
	data, expired, err := decryptAndUnwrapTTL(t.db.encryptor, result, t.db.opts.TTL.Enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt value: %w", err)
	}
	if expired {
		return nil, ErrKeyNotFound
	}
	return data, nil
}

// Set stores a key-value pair in the transaction.
// This is only valid for read-write transactions.
func (t *Txn) Set(key, value []byte) error {
	if t.closed || t.readonly {
		return ErrTxnClosed
	}
	if t.batch == nil {
		return ErrTxnClosed
	}
	return t.batch.Set(key, value)
}

// Delete removes a key from the transaction.
// This is only valid for read-write transactions.
func (t *Txn) Delete(key []byte) error {
	if t.closed || t.readonly {
		return ErrTxnClosed
	}
	if t.batch == nil {
		return ErrTxnClosed
	}
	return t.batch.Delete(key)
}

// DeleteRange removes all keys in the range [start, end) from the transaction.
// This is only valid for read-write transactions.
func (t *Txn) DeleteRange(start, end []byte) error {
	if t.closed || t.readonly {
		return ErrTxnClosed
	}
	if t.batch == nil {
		return ErrTxnClosed
	}
	return t.batch.DeleteRange(start, end)
}

// NewIter returns an iterator over the transaction.
// For read-only transactions, it iterates over the snapshot.
// For read-write transactions, it iterates over the database state at the
// time the transaction began. Uncommitted batch writes are not visible
// through the iterator; this matches Pebble's transactional semantics.
func (t *Txn) NewIter(opts *pebble.IterOptions) (*pebble.Iterator, error) {
	if t.closed {
		return nil, ErrTxnClosed
	}
	if t.snapshot != nil {
		return t.snapshot.NewIter(opts)
	}
	if t.batch != nil {
		return t.db.db.NewIter(opts)
	}
	return nil, ErrTxnClosed
}

// Close closes the transaction and releases associated resources.
// For read-write transactions, this does not commit changes.
func (t *Txn) Close() error {
	if t.closed {
		return nil
	}
	t.closed = true
	if t.snapshot != nil {
		t.snapshot.Close()
	}
	if t.batch != nil {
		return t.batch.Close()
	}
	return nil
}

// NewLazyIterator creates a lazy iterator for the transaction.
// The returned iterator uses the same LazyIterator type as DB.NewLazyIterator.
func (t *Txn) NewLazyIterator(start, end []byte) *LazyIterator {
	if t.closed {
		return &LazyIterator{err: ErrTxnClosed}
	}

	iter, err := t.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	if err != nil {
		return &LazyIterator{err: err}
	}

	var ttlOpts *TTLOptions
	if t.db.opts.TTL.Enabled {
		ttlOpts = &t.db.opts.TTL
	}

	return &LazyIterator{
		iter:      iter,
		encryptor: t.db.encryptor,
		ttlOpts:   ttlOpts,
		first:     true,
	}
}
