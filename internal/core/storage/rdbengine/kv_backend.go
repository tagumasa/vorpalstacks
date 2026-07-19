package rdbengine

import "vorpalstacks/internal/core/storage"

// kvBackend abstracts raw key-value operations for the relational engine.
// Wraps storage.BatchBucket obtained from RegionStorageManager.
type kvBackend struct {
	bucket storage.BatchBucket
}

func (b *kvBackend) get(key []byte) ([]byte, error) {
	return b.bucket.Get(key)
}

func (b *kvBackend) set(key, value []byte) error {
	return b.bucket.Put(key, value)
}

func (b *kvBackend) delete(key []byte) error {
	return b.bucket.Delete(key)
}

func (b *kvBackend) newBatch() kvBatch {
	return &batchWrap{batch: b.bucket.NewBatch()}
}

func (b *kvBackend) newIter(lower, upper []byte) kvIterator {
	return &iterWrap{iter: b.bucket.ScanRange(lower, upper)}
}

func (b *kvBackend) has(key []byte) bool {
	return b.bucket.Has(key)
}

type kvBatch interface {
	put(key, value []byte)
	del(key []byte)
	deleteRange(start, end []byte)
	commit() error
	commitSync() error
	close()
}

type batchWrap struct {
	batch storage.Batch
}

func (w *batchWrap) put(key, value []byte)   { _ = w.batch.Put(key, value) }
func (w *batchWrap) del(key []byte)          { _ = w.batch.Delete(key) }
func (w *batchWrap) deleteRange(s, e []byte) { _ = w.batch.DeleteRange(s, e) }
func (w *batchWrap) commit() error           { return w.batch.Commit() }
func (w *batchWrap) commitSync() error       { return w.batch.CommitSync() }
func (w *batchWrap) close()                  { w.batch.Close() }

type kvIterator interface {
	first()
	valid() bool
	next()
	key() []byte
	value() []byte
	err() error
	close()
}

type iterWrap struct {
	iter  storage.Iterator
	ok    bool
	start bool
}

func (w *iterWrap) first() {
	w.start = true
	w.ok = w.iter.Next()
}

func (w *iterWrap) valid() bool   { return w.ok }
func (w *iterWrap) next()         { w.ok = w.iter.Next() }
func (w *iterWrap) key() []byte   { return w.iter.Key() }
func (w *iterWrap) value() []byte { return w.iter.Value() }
func (w *iterWrap) err() error    { return w.iter.Error() }
func (w *iterWrap) close()        { w.iter.Close() }

// snapshotIterWrap adapts a storage.Iterator (returned by Snapshot.ScanRange)
// to the internal kvIterator interface.
type snapshotIterWrap struct {
	iter storage.Iterator
	ok   bool
}

func (w *snapshotIterWrap) first()        { w.ok = w.iter.Next() }
func (w *snapshotIterWrap) valid() bool   { return w.ok }
func (w *snapshotIterWrap) next()         { w.ok = w.iter.Next() }
func (w *snapshotIterWrap) key() []byte   { return w.iter.Key() }
func (w *snapshotIterWrap) value() []byte { return w.iter.Value() }
func (w *snapshotIterWrap) err() error    { return w.iter.Error() }
func (w *snapshotIterWrap) close()        { w.iter.Close() }
