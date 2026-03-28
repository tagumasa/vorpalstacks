package pebbledb

import (
	"context"
	"errors"
	"fmt"
	"time"

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
		val, closer, err := t.snapshot.Get(key)
		if err != nil {
			if err == pebble.ErrNotFound {
				return nil, ErrKeyNotFound
			}
			return nil, err
		}
		defer closer.Close()
		result := make([]byte, len(val))
		copy(result, val)
		decrypted, err := t.db.encryptor.Decrypt(result)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt value: %w", err)
		}
		return decrypted, nil
	}
	if t.batch != nil {
		val, closer, err := t.db.db.Get(key)
		if err != nil {
			if err == pebble.ErrNotFound {
				return nil, ErrKeyNotFound
			}
			return nil, err
		}
		defer closer.Close()
		result := make([]byte, len(val))
		copy(result, val)
		decrypted, err := t.db.encryptor.Decrypt(result)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt value: %w", err)
		}
		return decrypted, nil
	}
	return nil, ErrTxnClosed
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
// For read-write transactions, it iterates over the batch.
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

// TxnLazyIterator provides lazy iteration over a transaction range.
type TxnLazyIterator struct {
	iter      *pebble.Iterator
	encryptor Encryptor
	ttlOpts   *TTLOptions
	key       []byte
	value     []byte
	err       error
	closed    bool
	first     bool
}

// NewTxnLazyIterator creates a lazy iterator for the transaction.
func (t *Txn) NewTxnLazyIterator(start, end []byte) *TxnLazyIterator {
	if t.closed {
		return &TxnLazyIterator{err: ErrTxnClosed}
	}

	iter, err := t.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	if err != nil {
		return &TxnLazyIterator{err: err}
	}

	var ttlOpts *TTLOptions
	if t.db.opts.TTL.Enabled {
		ttlOpts = &t.db.opts.TTL
	}

	li := &TxnLazyIterator{
		iter:      iter,
		encryptor: t.db.encryptor,
		ttlOpts:   ttlOpts,
	}
	li.first = true
	return li
}

func (li *TxnLazyIterator) advance() {
	for li.iter.Valid() {
		key := li.iter.Key()
		val, err := li.iter.ValueAndErr()
		if err != nil {
			li.err = err
			return
		}

		decrypted, err := li.encryptor.Decrypt(val)
		if err != nil {
			li.err = err
			return
		}

		if li.ttlOpts != nil {
			ttlVal, err := UnmarshalTTLValue(decrypted)
			if err == nil && ttlVal.ExpiresAt > 0 {
				if time.Now().Unix() > ttlVal.ExpiresAt {
					li.iter.Next()
					continue
				}
				decrypted = ttlVal.Data
			}
		}

		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)
		li.key = keyCopy
		li.value = decrypted
		return
	}
	li.key = nil
	li.value = nil
}

func (li *TxnLazyIterator) Next() bool {
	if li.err != nil || li.closed {
		return false
	}
	if li.first {
		li.first = false
		li.iter.First()
	} else {
		li.iter.Next()
	}
	li.advance()
	return li.key != nil
}

func (li *TxnLazyIterator) Key() []byte {
	return li.key
}

func (li *TxnLazyIterator) Value() []byte {
	return li.value
}

func (li *TxnLazyIterator) Error() error {
	return li.err
}

func (li *TxnLazyIterator) Close() {
	if li.closed {
		return
	}
	li.closed = true
	if li.iter != nil {
		li.iter.Close()
	}
}
