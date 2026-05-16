package rdbengine

import (
	"context"
	"encoding/binary"
	"errors"
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

// Close marks the store as closed. All subsequent operations return errClosed.
func (s *Store) Close() error {
	s.closed.Store(true)
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

// GetRow retrieves a single row by primary key.
func (s *Store) GetRow(ctx context.Context, db, table string, pk []byte) (Row, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	key := rowKey(s.engine, db, table, pk)
	data, err := s.backend.get(key)
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
		iter:     s.backend.newIter(lower, upper),
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

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	key := rowKey(s.engine, db, table, pk)
	existing, err := s.backend.get(key)
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
	s.appendIndexEntries(batch, db, table, pk, row)

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

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	key := rowKey(s.engine, db, table, pk)
	existing, err := s.backend.get(key)
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
	s.appendRemoveIndexEntries(batch, db, table, pk, oldRow)
	s.appendIndexEntries(batch, db, table, pk, row)

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

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	key := rowKey(s.engine, db, table, pk)
	data, err := s.backend.get(key)
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
	s.appendRemoveIndexEntries(batch, db, table, pk, oldRow)

	if err := batch.commit(); err != nil {
		return fmtErr("delete_row commit", err)
	}
	return nil
}

// TruncateTable removes all rows from a table and returns the count deleted.
func (s *Store) TruncateTable(ctx context.Context, db, table string) (int, error) {
	if err := s.checkOpen(); err != nil {
		return 0, err
	}

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	tblKey := catalogTableKey(s.engine, db, table)
	if !s.backend.has(tblKey) {
		return 0, ErrNotFound
	}

	prefix := rowKeyPrefix(s.engine, db, table)
	end := rowEndKey(s.engine, db, table)

	iter := s.backend.newIter(prefix, end)
	defer iter.close()

	count := 0
	for iter.first(); iter.valid(); iter.next() {
		count++
	}
	if iter.err() != nil {
		return 0, fmtErr("truncate_table scan", iter.err())
	}

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
	return count, nil
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
	data, err := s.backend.get(key)
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
	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	key := autoincKey(s.engine, db, table)
	data, err := s.backend.get(key)
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
