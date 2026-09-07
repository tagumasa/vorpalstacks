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
	// seekFirst positions the iterator at the first entry to yield (First
	// for a forward walk; Last or SeekLT(before) for a reverse walk), and
	// move advances the underlying pebble iterator one entry in the
	// iterator's direction (Next forward, Prev reverse), including when the
	// TTL sweep skips expired entries. move's result is discarded: advance
	// re-checks validity itself.
	seekFirst func() bool
	move      func() bool
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
		seekFirst: iter.First,
		move:      iter.Next,
	}
}

// NewReverseLazyIterator creates a lazy iterator that walks the given range
// backwards. Iteration starts at the largest key strictly less than before
// (or the largest key in range when before is nil). The caller must call
// Close when finished.
func (d *DB) NewReverseLazyIterator(start, end, before []byte) *LazyIterator {
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

	seek := iter.Last
	if before != nil {
		seek = func() bool { return iter.SeekLT(before) }
	}

	return &LazyIterator{
		iter:      iter,
		encryptor: d.encryptor,
		ttlOpts:   ttlOpts,
		first:     true,
		seekFirst: seek,
		move:      iter.Prev,
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
			li.move()
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
		li.seekFirst()
	} else {
		li.move()
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
