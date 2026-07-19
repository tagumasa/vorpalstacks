package pebbledb

import (
	"github.com/cockroachdb/pebble/v2"
)

// SnapshotScanner provides point-in-time consistent reads over the
// database via a Pebble snapshot. All Get and iteration operations
// observe the database state as of the moment NewSnapshotScanner was
// called, regardless of concurrent writes. Decryption and TTL
// unwrapping are applied transparently, matching the behaviour of
// non-snapshot reads through DB.Get and DB.NewLazyIterator.
//
// The caller must call Close when finished to release the snapshot.
type SnapshotScanner struct {
	db       *DB
	snapshot *pebble.Snapshot
	closed   bool
}

// NewSnapshotScanner creates a new point-in-time snapshot scanner.
func (d *DB) NewSnapshotScanner() *SnapshotScanner {
	return &SnapshotScanner{
		db:       d,
		snapshot: d.db.NewSnapshot(),
	}
}

// Get retrieves a value by key from the snapshot, with decryption and
// TTL unwrapping applied. Returns ErrKeyNotFound if the key does not
// exist in the snapshot view.
func (s *SnapshotScanner) Get(key []byte) ([]byte, error) {
	val, closer, err := s.snapshot.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	defer closer.Close()

	result := make([]byte, len(val))
	copy(result, val)

	data, expired, err := decryptAndUnwrapTTL(s.db.encryptor, result, s.db.opts.TTL.Enabled)
	if err != nil {
		return nil, err
	}
	if expired {
		return nil, ErrKeyNotFound
	}
	return data, nil
}

// Has returns true if the key exists in the snapshot view.
func (s *SnapshotScanner) Has(key []byte) bool {
	_, err := s.Get(key)
	return err == nil
}

// NewLazyIterator creates a lazy iterator over the given key range
// backed by the snapshot. Values are decrypted and TTL-filtered as
// they are read.
func (s *SnapshotScanner) NewLazyIterator(start, end []byte) *LazyIterator {
	iter, err := s.snapshot.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	if err != nil {
		return &LazyIterator{err: err}
	}

	var ttlOpts *TTLOptions
	if s.db.opts.TTL.Enabled {
		ttlOpts = &s.db.opts.TTL
	}

	return &LazyIterator{
		iter:      iter,
		encryptor: s.db.encryptor,
		ttlOpts:   ttlOpts,
		first:     true,
	}
}

// Close releases the snapshot.
func (s *SnapshotScanner) Close() {
	if s.closed {
		return
	}
	s.closed = true
	s.snapshot.Close()
}
