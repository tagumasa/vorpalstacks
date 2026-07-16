package pebbledb

import (
	"github.com/cockroachdb/pebble/v2"
)

// LazyIterator provides lazy, on-demand iteration over a key range.
// Unlike ScanRange which materialises all results upfront, this iterator
// wraps a pebble.Iterator and decrypts values one at a time.
type LazyIterator struct {
	iter      *pebble.Iterator
	encryptor Encryptor
	ttlOpts   *TTLOptions
	key       []byte
	value     []byte
	err       error
	closed    bool
	first     bool
}

// NewLazyIterator creates a lazy iterator over the given range.
// The caller must call Close when finished.
func (d *DB) NewLazyIterator(start, end []byte) *LazyIterator {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.closed {
		return &LazyIterator{err: ErrClosed}
	}

	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	if err != nil {
		return &LazyIterator{err: err}
	}

	var ttlOpts *TTLOptions
	if d.opts.TTL.Enabled {
		ttlOpts = &d.opts.TTL
	}

	return &LazyIterator{
		iter:      iter,
		encryptor: d.encryptor,
		ttlOpts:   ttlOpts,
		first:     true,
	}
}

func (li *LazyIterator) ttlEnabled() bool {
	return li.ttlOpts != nil
}

func (li *LazyIterator) advance() {
	for li.iter.Valid() {
		val, err := li.iter.ValueAndErr()
		if err != nil {
			li.err = err
			return
		}

		data, expired, err := decryptAndUnwrapTTL(li.encryptor, val, li.ttlEnabled())
		if err != nil {
			li.err = err
			return
		}
		if expired {
			li.iter.Next()
			continue
		}

		key := li.iter.Key()
		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)
		li.key = keyCopy
		li.value = data
		return
	}
	li.key = nil
	li.value = nil
}

// Next advances the iterator to the next key-value pair.
func (li *LazyIterator) Next() bool {
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

// Key returns the current key.
func (li *LazyIterator) Key() []byte {
	return li.key
}

// Value returns the current decrypted value.
func (li *LazyIterator) Value() []byte {
	return li.value
}

// Error returns any error encountered during iteration.
func (li *LazyIterator) Error() error {
	return li.err
}

// Close releases the underlying iterator resources.
func (li *LazyIterator) Close() {
	if li.closed {
		return
	}
	li.closed = true
	if li.iter != nil {
		li.iter.Close()
	}
}
