package graphengine

import (
	"fmt"
	"os"

	"github.com/cockroachdb/pebble/v2"
	"vorpalstacks/internal/core/storage"
)

// kvBackend abstracts raw key-value operations for the graph engine.
// Two implementations exist: pebbleBackend (legacy independent DB) and
// bucketBackend (shared storage via RegionStorageManager).
type kvBackend interface {
	get(key []byte) ([]byte, error)
	set(key, value []byte) error
	delete(key []byte) error
	newBatch() kvBatch
	newIter(lower, upper []byte) (kvIterator, error)
	deleteRange(start, end []byte) error
	diskSize() int64
	close() error
}

type kvBatch interface {
	put(key, value []byte)
	del(key []byte)
	deleteRange(start, end []byte)
	commit() error
	commitSync() error
	close()
}

type kvIterator interface {
	first()
	valid() bool
	next()
	key() []byte
	value() []byte
	err() error
	close()
}

// --- pebbleBackend: wraps *pebble.DB (legacy) ---

type pebbleBackend struct {
	db    *pebble.DB
	cache *pebble.Cache
}

func (p *pebbleBackend) get(key []byte) ([]byte, error) {
	val, closer, err := p.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	defer closer.Close()
	result := make([]byte, len(val))
	copy(result, val)
	return result, nil
}

func (p *pebbleBackend) set(key, value []byte) error {
	return p.db.Set(key, value, pebble.NoSync)
}

func (p *pebbleBackend) delete(key []byte) error {
	return p.db.Delete(key, pebble.NoSync)
}

func (p *pebbleBackend) newBatch() kvBatch {
	return &pebbleBatchWrap{batch: p.db.NewBatch()}
}

func (p *pebbleBackend) newIter(lower, upper []byte) (kvIterator, error) {
	iter, err := p.db.NewIter(&pebble.IterOptions{LowerBound: lower, UpperBound: upper})
	if err != nil {
		return nil, err
	}
	return &pebbleIterWrap{iter: iter}, nil
}

func (p *pebbleBackend) deleteRange(start, end []byte) error {
	return p.db.DeleteRange(start, end, pebble.NoSync)
}

func (p *pebbleBackend) diskSize() int64 {
	return int64(p.db.Metrics().DiskSpaceUsage())
}

func (p *pebbleBackend) close() error {
	return p.db.Close()
}

type pebbleBatchWrap struct {
	batch *pebble.Batch
}

func (w *pebbleBatchWrap) put(key, value []byte)      { _ = w.batch.Set(key, value, nil) }
func (w *pebbleBatchWrap) del(key []byte)              { _ = w.batch.Delete(key, nil) }
func (w *pebbleBatchWrap) deleteRange(s, e []byte)     { _ = w.batch.DeleteRange(s, e, nil) }
func (w *pebbleBatchWrap) commit() error               { return w.batch.Commit(pebble.NoSync) }
func (w *pebbleBatchWrap) commitSync() error           { return w.batch.Commit(pebble.Sync) }
func (w *pebbleBatchWrap) close()                      { w.batch.Close() }

type pebbleIterWrap struct {
	iter *pebble.Iterator
}

func (w *pebbleIterWrap) first()          { w.iter.First() }
func (w *pebbleIterWrap) valid() bool     { return w.iter.Valid() }
func (w *pebbleIterWrap) next()           { w.iter.Next() }
func (w *pebbleIterWrap) key() []byte     { return w.iter.Key() }
func (w *pebbleIterWrap) value() []byte   { return w.iter.Value() }
func (w *pebbleIterWrap) err() error      { return w.iter.Error() }
func (w *pebbleIterWrap) close()          { w.iter.Close() }

// --- bucketBackend: wraps storage.BatchBucket (shared storage) ---

type bucketBackend struct {
	bucket storage.BatchBucket
}

func (b *bucketBackend) get(key []byte) ([]byte, error) {
	return b.bucket.Get(key)
}

func (b *bucketBackend) set(key, value []byte) error {
	return b.bucket.Put(key, value)
}

func (b *bucketBackend) delete(key []byte) error {
	return b.bucket.Delete(key)
}

func (b *bucketBackend) newBatch() kvBatch {
	return &bucketBatchWrap{batch: b.bucket.NewBatch()}
}

func (b *bucketBackend) newIter(lower, upper []byte) (kvIterator, error) {
	return &bucketIterWrap{iter: b.bucket.ScanRange(lower, upper)}, nil
}

func (b *bucketBackend) deleteRange(start, end []byte) error {
	return b.bucket.DeleteRange(start, end)
}

func (b *bucketBackend) diskSize() int64 {
	return 0
}

func (b *bucketBackend) close() error {
	return nil
}

type bucketBatchWrap struct {
	batch storage.Batch
}

func (w *bucketBatchWrap) put(key, value []byte)  { _ = w.batch.Put(key, value) }
func (w *bucketBatchWrap) del(key []byte)         { _ = w.batch.Delete(key) }
func (w *bucketBatchWrap) deleteRange(s, e []byte) { _ = w.batch.DeleteRange(s, e) }
func (w *bucketBatchWrap) commit() error          { return w.batch.Commit() }
func (w *bucketBatchWrap) commitSync() error      { return w.batch.Commit() }
func (w *bucketBatchWrap) close()                 { w.batch.Close() }

type bucketIterWrap struct {
	iter  storage.Iterator
	ok    bool
	start bool
}

func (w *bucketIterWrap) first() {
	w.start = true
	w.ok = w.iter.Next()
}

func (w *bucketIterWrap) valid() bool   { return w.ok }
func (w *bucketIterWrap) next()         { w.ok = w.iter.Next() }
func (w *bucketIterWrap) key() []byte   { return w.iter.Key() }
func (w *bucketIterWrap) value() []byte { return w.iter.Value() }
func (w *bucketIterWrap) err() error    { return w.iter.Error() }
func (w *bucketIterWrap) close()        { w.iter.Close() }

// New creates a graph engine backed by a shared storage bucket.
// The bucket is obtained from RegionStorageManager via BatchBucket.
// The caller is responsible for bucket lifecycle; Close() is a no-op.
func New(bucket storage.BatchBucket, opts Options) (*DB, error) {
	backend := &bucketBackend{bucket: bucket}

	d := &DB{
		backend: backend,
		dir:     "",
	}

	if err := d.loadCounters(); err != nil {
		return nil, err
	}

	d.Embeddings = newEmbeddingStore(d)
	return d, nil
}

func newEmbeddingStore(d *DB) *EmbeddingStore {
	es := &EmbeddingStore{
		db:    &pebbleDBAdapter{db: d},
		cache: make(map[NodeID][]float64),
	}
	es.loadCache()
	return es
}

// reconstructFromBackend rebuilds counters from existing data after migration.
func (d *DB) reconstructFromBackend() {
	d.selfHealCounters()
	if err := d.persistCounters(); err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: failed to persist counters after reconstruction: %v\n", err)
	}
}

// RebuildCountersIfNeeded checks if counters are zero and rebuilds from data.
// Call after New() when migrating from independent pebble to shared storage.
func (d *DB) RebuildCountersIfNeeded() {
	if d.nodeCount.Load() == 0 && d.edgeCount.Load() == 0 {
		d.reconstructFromBackend()
	}
}
