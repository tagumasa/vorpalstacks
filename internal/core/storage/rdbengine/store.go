package rdbengine

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"vorpalstacks/internal/core/storage"
)

// Store implements RowStore backed by Pebble key-value storage via
// storage.BatchBucket obtained from RegionStorageManager.
// A per-table mutex serialises write operations (InsertRow, UpdateRow,
// DeleteRow) so that the read-check-write cycle is free from TOCTOU
// races under concurrent access.
type Store struct {
	backend *kvBackend
	engine  string
	closed  atomic.Bool
	mu      sync.Map // map[string]*sync.Mutex — keyed by "db/table"
	snap    storage.Snapshot
}

var errClosed = errors.New("rdbengine: store is closed")

// New creates a new relational engine backed by a shared storage bucket.
// The bucket is obtained from RegionStorageManager via BatchBucket.
// The caller is responsible for bucket lifecycle; Close() is a no-op.
func New(bucket storage.BatchBucket, opts Options) (*Store, error) {
	engine := opts.Engine
	if engine == "" {
		engine = "mysql"
	}
	return &Store{
		backend: &kvBackend{bucket: bucket},
		engine:  engine,
	}, nil
}

// errReadOnly is returned by every write method on a snapshot Store
// (one created by NewSnapshotReader). Snapshot stores are read-only;
// attempting to write through them would mix snapshot reads with live
// writes, producing an inconsistent state.
var errReadOnly = fmt.Errorf("rdbengine: store is read-only (snapshot)")

// NewSnapshotReader returns a new Store whose read operations (get, has,
// newIter) go through a point-in-time Pebble snapshot for consistent
// cross-table reads. All reads through the returned Store observe a
// consistent view as of the moment of creation, regardless of concurrent
// writes. The caller must call CloseSnapshot (or Close) when finished to
// release the snapshot.
//
// If the backing bucket does not implement storage.Snapshotter (e.g. an
// in-memory test fake), an error is returned.
func (s *Store) NewSnapshotReader() (*Store, error) {
	snap, ok := s.backend.bucket.(storage.Snapshotter)
	if !ok {
		return nil, fmt.Errorf("rdbengine: backing bucket %T does not support snapshots", s.backend.bucket)
	}
	return &Store{
		backend: s.backend,
		engine:  s.engine,
		snap:    snap.NewSnapshot(),
	}, nil
}

// CloseSnapshot releases the snapshot if this store was created by
// NewSnapshotReader. Safe to call on a regular (non-snapshot) store.
func (s *Store) CloseSnapshot() {
	if s.snap != nil {
		s.snap.Close()
		s.snap = nil
	}
}

// get retrieves a single key. When a snapshot is active, reads go
// through the snapshot for point-in-time consistency.
func (s *Store) get(key []byte) ([]byte, error) {
	if s.snap != nil {
		return s.snap.Get(key)
	}
	return s.backend.get(key)
}

// has checks whether a key exists. Snapshot-aware.
func (s *Store) has(key []byte) bool {
	if s.snap != nil {
		v, err := s.snap.Get(key)
		return err == nil && v != nil
	}
	return s.backend.has(key)
}

// newIter creates a key-range iterator. When a snapshot is active, the
// iterator observes the snapshot view. The returned kvIterator must be
// closed by the caller.
func (s *Store) newIter(lower, upper []byte) kvIterator {
	if s.snap != nil {
		return &snapshotIterWrap{iter: s.snap.ScanRange(lower, upper)}
	}
	return s.backend.newIter(lower, upper)
}

// Close marks the store as closed. All subsequent operations return
// errClosed. If the store was created by NewSnapshotReader, the
// underlying Pebble snapshot is also released.
func (s *Store) Close() error {
	s.closed.Store(true)
	if s.snap != nil {
		s.snap.Close()
		s.snap = nil
	}
	return nil
}

func (s *Store) checkOpen() error {
	if s.closed.Load() {
		return errClosed
	}
	return nil
}

func (s *Store) TableLock(db, table string) *sync.Mutex {
	key := db + "/" + table
	v, _ := s.mu.LoadOrStore(key, &sync.Mutex{})
	return v.(*sync.Mutex)
}

// DatabaseLock returns a mutex serialising database-wide schema mutations
// (DropDatabase, CreateDatabase, CreateTable). DropDatabase acquires
// this lock, then every table-level lock, before issuing the atomic
// deleteRange batch; CreateDatabase and CreateTable acquire only the
// database lock (CreateTable additionally acquires the per-table lock).
// Data operations (InsertRow / UpdateRow / DeleteRow) acquire only
// TableLock and never DatabaseLock, so they cannot deadlock against
// schema operations.
func (s *Store) DatabaseLock(db string) *sync.Mutex {
	key := "db:" + db
	v, _ := s.mu.LoadOrStore(key, &sync.Mutex{})
	return v.(*sync.Mutex)
}

// TxnBatch is an atomic write batch for transaction support. Writes are
// accumulated and only persisted when Commit is called. Rollback discards all.
type TxnBatch struct {
	store  *Store
	batch  kvBatch
	closed bool
}

func (tb *TxnBatch) Commit() error {
	if tb.closed {
		return fmtErr("txn_batch commit", errors.New("batch already closed"))
	}
	tb.closed = true
	return tb.batch.commit()
}

func (tb *TxnBatch) Rollback() {
	if tb.closed {
		return
	}
	tb.closed = true
	tb.batch.close()
}

func (s *Store) NewTxnBatch() *TxnBatch {
	if s.snap != nil {
		return nil
	}
	return &TxnBatch{store: s, batch: s.backend.newBatch()}
}

func (s *Store) TxnInsertRow(tb *TxnBatch, db, table string, pk []byte, row Row) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}
	key := rowKey(s.engine, db, table, pk)
	data, err := encodeRow(row)
	if err != nil {
		return fmtErr("txn_insert encode", err)
	}
	tb.batch.put(key, data)
	if err := s.appendIndexEntries(tb.batch, db, table, pk, row); err != nil {
		return err
	}
	return nil
}

func (s *Store) TxnUpdateRow(tb *TxnBatch, db, table string, pk []byte, row Row) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}
	key := rowKey(s.engine, db, table, pk)

	// Read the existing row from committed storage so we can remove its
	// old index entries. Rows inserted earlier in the same uncommitted
	// transaction are not visible here; that edge case leaves orphaned
	// index entries for never-committed rows, which is acceptable.
	if existing, err := s.get(key); err != nil {
		return fmtErr("txn_update read_old", err)
	} else if existing != nil {
		if oldRow, err := decodeRow(existing); err != nil {
			return fmtErr("txn_update decode_old", err)
		} else {
			if err := s.appendRemoveIndexEntries(tb.batch, db, table, pk, oldRow); err != nil {
				return err
			}
		}
	}

	data, err := encodeRow(row)
	if err != nil {
		return fmtErr("txn_update encode", err)
	}
	tb.batch.put(key, data)
	if err := s.appendIndexEntries(tb.batch, db, table, pk, row); err != nil {
		return err
	}
	return nil
}

func (s *Store) TxnDeleteRow(tb *TxnBatch, db, table string, pk []byte) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}
	key := rowKey(s.engine, db, table, pk)

	// Read the existing row from committed storage so we can remove its
	// index entries atomically with the row deletion.
	if existing, err := s.get(key); err != nil {
		return fmtErr("txn_delete read_old", err)
	} else if existing != nil {
		if oldRow, err := decodeRow(existing); err != nil {
			return fmtErr("txn_delete decode_old", err)
		} else {
			if err := s.appendRemoveIndexEntries(tb.batch, db, table, pk, oldRow); err != nil {
				return err
			}
		}
	}

	tb.batch.del(key)
	return nil
}

// GetRow retrieves a single row by primary key.
func (s *Store) GetRow(ctx context.Context, db, table string, pk []byte) (Row, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	key := rowKey(s.engine, db, table, pk)
	data, err := s.get(key)
	if err != nil {
		return nil, fmtErr("get_row", err)
	}
	if data == nil {
		return nil, ErrNotFound
	}
	row, err := decodeRow(data)
	if err != nil {
		return nil, fmtErr("get_row decode", err)
	}
	return row, nil
}

// ScanRows returns an iterator over rows in a table, optionally bounded by
// primary key range, with limit/offset support.
func (s *Store) ScanRows(ctx context.Context, db, table string, opts ScanOptions) (RowIterator, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	schema, err := s.GetTableSchema(ctx, db, table)
	if err != nil {
		return nil, err
	}

	var lower, upper []byte
	if len(opts.StartPK) > 0 {
		lower = rowKey(s.engine, db, table, opts.StartPK)
	} else {
		lower = rowKeyPrefix(s.engine, db, table)
	}
	if len(opts.EndPK) > 0 {
		upper = rowKey(s.engine, db, table, opts.EndPK)
	} else {
		upper = rowEndKey(s.engine, db, table)
	}

	return &scanRowIter{
		schema:   schema,
		iter:     s.newIter(lower, upper),
		limit:    opts.Limit,
		offset:   opts.Offset,
		consumed: 0,
	}, nil
}

// InsertRow writes a new row. Returns an error if the primary key already exists.
// The row data and index entries are written atomically in a single batch.
func (s *Store) InsertRow(ctx context.Context, db, table string, pk []byte, row Row) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	key := rowKey(s.engine, db, table, pk)
	existing, err := s.get(key)
	if err != nil {
		return fmtErr("insert_row check", err)
	}
	if existing != nil {
		return ErrAlreadyExists
	}

	data, err := encodeRow(row)
	if err != nil {
		return fmtErr("insert_row encode", err)
	}

	if err := s.checkUniqueConstraints(db, table, pk, row); err != nil {
		return err
	}

	batch := s.backend.newBatch()
	defer batch.close()

	batch.put(key, data)
	if err := s.appendIndexEntries(batch, db, table, pk, row); err != nil {
		return err
	}

	if err := batch.commit(); err != nil {
		return fmtErr("insert_row commit", err)
	}
	return nil
}

// UpdateRow replaces an existing row's data and updates index entries.
// The old index removal, new data, and new index entries are written atomically.
func (s *Store) UpdateRow(ctx context.Context, db, table string, pk []byte, row Row) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	key := rowKey(s.engine, db, table, pk)
	existing, err := s.get(key)
	if err != nil {
		return fmtErr("update_row check", err)
	}
	if existing == nil {
		return ErrNotFound
	}

	oldRow, err := decodeRow(existing)
	if err != nil {
		return fmtErr("update_row decode_old", err)
	}

	data, err := encodeRow(row)
	if err != nil {
		return fmtErr("update_row encode", err)
	}

	if err := s.checkUniqueConstraints(db, table, pk, row); err != nil {
		return err
	}

	batch := s.backend.newBatch()
	defer batch.close()

	batch.put(key, data)
	if err := s.appendRemoveIndexEntries(batch, db, table, pk, oldRow); err != nil {
		return err
	}
	if err := s.appendIndexEntries(batch, db, table, pk, row); err != nil {
		return err
	}

	if err := batch.commit(); err != nil {
		return fmtErr("update_row commit", err)
	}
	return nil
}

// DeleteRow removes a row and its index entries atomically.
func (s *Store) DeleteRow(ctx context.Context, db, table string, pk []byte) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	key := rowKey(s.engine, db, table, pk)
	data, err := s.get(key)
	if err != nil {
		return fmtErr("delete_row check", err)
	}
	if data == nil {
		return ErrNotFound
	}

	oldRow, err := decodeRow(data)
	if err != nil {
		return fmtErr("delete_row decode", err)
	}

	batch := s.backend.newBatch()
	defer batch.close()

	batch.del(key)
	if err := s.appendRemoveIndexEntries(batch, db, table, pk, oldRow); err != nil {
		return err
	}

	if err := batch.commit(); err != nil {
		return fmtErr("delete_row commit", err)
	}
	return nil
}

// TruncateTable removes all rows from a table. The returned count is
// always 0 because MySQL's TRUNCATE TABLE is a DDL operation that, per
// the MySQL reference manual, returns "0 rows affected" to the client
// regardless of how many rows were removed.
//
// Implementation note: the previous implementation iterated every row
// to compute the count before deleting, which made TRUNCATE O(N) and
// held the per-table mutex for the entire scan. On a large table this
// blocked writes for seconds. The current implementation issues a
// single Pebble DeleteRange covering the row key range plus an
// analogous DeleteRange per secondary index, making the operation
// effectively O(number-of-indexes) regardless of table size.
func (s *Store) TruncateTable(ctx context.Context, db, table string) (int, error) {
	if err := s.checkOpen(); err != nil {
		return 0, err
	}
	if s.snap != nil {
		return 0, errReadOnly
	}

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	tblKey := catalogTableKey(s.engine, db, table)
	if !s.has(tblKey) {
		return 0, ErrNotFound
	}

	prefix := rowKeyPrefix(s.engine, db, table)
	end := rowEndKey(s.engine, db, table)

	batch := s.backend.newBatch()
	defer batch.close()

	batch.deleteRange(prefix, end)

	indexes, err := s.listIndexDefs(db, table)
	if err != nil {
		return 0, err
	}
	for _, idx := range indexes {
		idxStart := indexKeyPrefix(s.engine, db, table, idx.Name)
		idxEnd := indexEndKey(s.engine, db, table, idx.Name)
		batch.deleteRange(idxStart, idxEnd)
		if idx.Unique {
			uniqStart := uniqueKeyPrefix(s.engine, db, table, idx.Name)
			uniqEnd := uniqueEndKey(s.engine, db, table, idx.Name)
			batch.deleteRange(uniqStart, uniqEnd)
		}
	}

	if err := batch.commit(); err != nil {
		return 0, fmtErr("truncate_table commit", err)
	}
	return 0, nil
}

// scanRowIter iterates over row key-value pairs in Pebble key order.
type scanRowIter struct {
	schema   *TableSchema
	iter     kvIterator
	limit    int
	offset   int
	consumed int
	current  Row
	err      error
}

func (it *scanRowIter) Next() bool {
	started := false
	for {
		if !started {
			it.iter.first()
			started = true
		} else {
			it.iter.next()
		}
		if !it.iter.valid() {
			return false
		}

		it.consumed++
		if it.consumed <= it.offset {
			continue
		}
		if it.limit > 0 && it.consumed-it.offset > it.limit {
			return false
		}

		row, err := decodeRow(it.iter.value())
		if err != nil {
			it.err = err
			return false
		}
		it.current = row
		return true
	}
}

func (it *scanRowIter) Row() Row     { return it.current }
func (it *scanRowIter) Error() error { return it.err }
func (it *scanRowIter) Close()       { it.iter.close() }

// GetAutoIncrement returns the current auto-increment counter for a table.
func (s *Store) GetAutoIncrement(db, table string) (uint64, error) {
	if err := s.checkOpen(); err != nil {
		return 0, err
	}
	key := autoincKey(s.engine, db, table)
	data, err := s.get(key)
	if err != nil {
		return 0, fmtErr("get_autoinc", err)
	}
	if data == nil {
		return 1, nil
	}
	return binary.BigEndian.Uint64(data), nil
}

// SetAutoIncrement atomically sets the auto-increment counter for a table.
func (s *Store) SetAutoIncrement(db, table string, val uint64) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}
	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], val)
	key := autoincKey(s.engine, db, table)
	return s.backend.set(key, buf[:])
}

// NextAutoIncrement atomically increments and returns the next value.
func (s *Store) NextAutoIncrement(db, table string) (uint64, error) {
	if err := s.checkOpen(); err != nil {
		return 0, err
	}
	if s.snap != nil {
		return 0, errReadOnly
	}
	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	key := autoincKey(s.engine, db, table)
	data, err := s.get(key)
	if err != nil {
		return 0, fmtErr("next_autoinc", err)
	}

	var current uint64
	if data != nil {
		current = binary.BigEndian.Uint64(data)
	}
	if current == 0 {
		current = 1
	}

	next := current + 1
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], next)
	if err := s.backend.set(key, buf[:]); err != nil {
		return 0, fmtErr("next_autoinc set", err)
	}
	return current, nil
}
