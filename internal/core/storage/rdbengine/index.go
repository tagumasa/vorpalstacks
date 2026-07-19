package rdbengine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// Index operations for secondary index management stored in Pebble.
// Index data keys:  idx/{engine}/{db}/{table}/{idx_name}/{col_val}/{pk} → empty
// Unique keys:      uniq/{engine}/{db}/{table}/{idx_name}/{col_val}     → {pk}
// Catalog keys:     catalog/{engine}/db/{db}/table/{tbl}/idx/{idx_name} → IndexDef JSON

// CreateIndex creates a secondary index definition and builds index entries
// from existing row data.
func (s *Store) CreateIndex(ctx context.Context, db, table, indexName string, columns []string, unique bool) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}
	if err := ValidateIdentifier("database", db); err != nil {
		return err
	}
	if err := ValidateIdentifier("table", table); err != nil {
		return err
	}
	if err := ValidateIdentifier("index", indexName); err != nil {
		return err
	}
	if len(columns) == 0 {
		return fmt.Errorf("rdbengine: create_index: columns list is empty")
	}

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	tblKey := catalogTableKey(s.engine, db, table)
	if !s.has(tblKey) {
		return ErrNotFound
	}

	idxCatalogKey := catalogIndexKey(s.engine, db, table, indexName)
	if s.has(idxCatalogKey) {
		return ErrAlreadyExists
	}

	schema, err := s.GetTableSchema(ctx, db, table)
	if err != nil {
		return err
	}
	for _, col := range columns {
		if schema.ColumnByName(col) == nil {
			return fmtErr("create_index", ErrNotFound)
		}
	}

	def := IndexDef{
		Name:    indexName,
		Table:   table,
		Columns: columns,
		Unique:  unique,
	}
	data, err := json.Marshal(def)
	if err != nil {
		return fmtErr("create_index marshal", err)
	}

	batch := s.backend.newBatch()
	defer batch.close()

	batch.put(idxCatalogKey, data)

	prefix := rowKeyPrefix(s.engine, db, table)
	end := rowEndKey(s.engine, db, table)
	iter := s.newIter(prefix, end)
	defer iter.close()

	for iter.first(); iter.valid(); iter.next() {
		pk := iter.key()[len(prefix):]
		row, decodeErr := decodeRow(iter.value())
		if decodeErr != nil {
			return fmtErr("create_index backfill decode", decodeErr)
		}

		colVal, hasNull, encErr := s.encodeIndexColumns(def, row)
		if encErr != nil {
			return encErr
		}
		// SQL standard / MySQL behaviour: NULL values are distinct for
		// uniqueness; never reject a row for colliding on an empty
		// unique key, and never insert a uniq/<empty> entry for a row
		// whose indexed columns are NULL.
		iKey := indexKey(s.engine, db, table, indexName, colVal, pk)
		batch.put(iKey, pk)

		if unique && !hasNull {
			uKey := uniqueKey(s.engine, db, table, indexName, colVal)
			existing, getErr := s.get(uKey)
			if getErr != nil {
				return fmtErr("create_index backfill unique check", getErr)
			}
			if existing != nil {
				return ErrConstraintViolation
			}
			batch.put(uKey, pk)
		}
	}
	if iter.err() != nil {
		return fmtErr("create_index backfill scan", iter.err())
	}

	return batch.commit()
}

// DropIndex removes an index definition and all its data entries.
func (s *Store) DropIndex(ctx context.Context, db, table, indexName string) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	idxCatalogKey := catalogIndexKey(s.engine, db, table, indexName)
	if !s.has(idxCatalogKey) {
		return ErrNotFound
	}

	batch := s.backend.newBatch()
	defer batch.close()

	batch.del(idxCatalogKey)

	idxStart := indexKeyPrefix(s.engine, db, table, indexName)
	idxEnd := indexEndKey(s.engine, db, table, indexName)
	batch.deleteRange(idxStart, idxEnd)

	uniqStart := uniqueKeyPrefix(s.engine, db, table, indexName)
	uniqEnd := uniqueEndKey(s.engine, db, table, indexName)
	batch.deleteRange(uniqStart, uniqEnd)

	return batch.commit()
}

// checkUniqueConstraints validates that no unique constraint would be violated
// by the given row. The pk parameter identifies the current row — when updating,
// existing unique entries belonging to the same pk are skipped.
//
// Per the SQL standard and MySQL behaviour, NULL values in unique-indexed
// columns are distinct: rows with NULL in any indexed column bypass the
// unique check entirely. This matches MySQL which allows multiple NULL
// values in a UNIQUE column.
func (s *Store) checkUniqueConstraints(db, table string, pk []byte, row Row) error {
	indexes, err := s.listIndexDefs(db, table)
	if err != nil {
		return err
	}
	for _, idx := range indexes {
		if !idx.Unique {
			continue
		}
		colVal, hasNull, encErr := s.encodeIndexColumns(idx, row)
		if encErr != nil {
			return encErr
		}
		if hasNull {
			// NULLs are distinct for uniqueness; skip the check entirely.
			continue
		}
		uKey := uniqueKey(s.engine, db, table, idx.Name, colVal)
		existing, err := s.get(uKey)
		if err != nil {
			return err
		}
		if existing != nil && !bytes.Equal(existing, pk) {
			return ErrConstraintViolation
		}
	}
	return nil
}

func (s *Store) appendIndexEntries(batch kvBatch, db, table string, pk []byte, row Row) error {
	indexes, err := s.listIndexDefs(db, table)
	if err != nil {
		return err
	}
	if len(indexes) == 0 {
		return nil
	}
	for _, idx := range indexes {
		colVal, hasNull, encErr := s.encodeIndexColumns(idx, row)
		if encErr != nil {
			return encErr
		}
		iKey := indexKey(s.engine, db, table, idx.Name, colVal, pk)
		batch.put(iKey, pk)
		if idx.Unique && !hasNull {
			uKey := uniqueKey(s.engine, db, table, idx.Name, colVal)
			batch.put(uKey, pk)
		}
	}
	return nil
}

func (s *Store) appendRemoveIndexEntries(batch kvBatch, db, table string, pk []byte, row Row) error {
	indexes, err := s.listIndexDefs(db, table)
	if err != nil {
		return err
	}
	if len(indexes) == 0 {
		return nil
	}
	for _, idx := range indexes {
		colVal, hasNull, encErr := s.encodeIndexColumns(idx, row)
		if encErr != nil {
			return encErr
		}
		iKey := indexKey(s.engine, db, table, idx.Name, colVal, pk)
		batch.del(iKey)
		// When the row had NULL indexed columns, no uniq entry was
		// written at insert time, so there is nothing to remove.
		if idx.Unique && !hasNull {
			uKey := uniqueKey(s.engine, db, table, idx.Name, colVal)
			batch.del(uKey)
		}
	}
	return nil
}

// ScanIndex returns rows matching an index scan via prefix iteration.
func (s *Store) ScanIndex(ctx context.Context, db, table, indexName string, opts IndexScanOptions) (RowIterator, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	schema, err := s.GetTableSchema(ctx, db, table)
	if err != nil {
		return nil, err
	}

	var lower, upper []byte
	if len(opts.Start) > 0 {
		lower = indexKey(s.engine, db, table, indexName, opts.Start, nil)
	} else {
		lower = indexKeyPrefix(s.engine, db, table, indexName)
	}
	if len(opts.End) > 0 {
		upper = indexKey(s.engine, db, table, indexName, opts.End, nil)
	} else {
		upper = indexEndKey(s.engine, db, table, indexName)
	}

	return &indexRowIter{
		store:    s,
		schema:   schema,
		db:       db,
		table:    table,
		iter:     s.newIter(lower, upper),
		limit:    opts.Limit,
		offset:   opts.Offset,
		consumed: 0,
	}, nil
}

// encodeIndexColumns encodes the indexed columns from a row into a single
// key component. The returned `hasNull` flag is true when any indexed
// column's value was NULL; per the SQL standard and MySQL semantics,
// NULL values are distinct for uniqueness purposes and must NOT collide
// on a common empty key. Callers (checkUniqueConstraints, CreateIndex
// backfill) skip the unique-check entirely when hasNull is true.
func (s *Store) encodeIndexColumns(idx IndexDef, row Row) ([]byte, bool, error) {
	var buf []byte
	hasNull := false
	for i, col := range idx.Columns {
		if i > 0 {
			buf = append(buf, 0x00)
		}
		cv, ok := row[col]
		if !ok || cv.Value == nil {
			hasNull = true
			continue
		}
		enc, err := encodeValue(cv)
		if err != nil {
			return nil, false, fmtErr("encode_index_columns", err)
		}
		buf = append(buf, enc...)
	}
	return buf, hasNull, nil
}

// indexRowIter iterates over index entries and fetches the corresponding rows.
type indexRowIter struct {
	store    *Store
	schema   *TableSchema
	db       string
	table    string
	iter     kvIterator
	limit    int
	offset   int
	consumed int
	current  Row
	err      error
}

func (it *indexRowIter) Next() bool {
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

		pk := it.iter.value()
		if len(pk) == 0 {
			continue
		}

		row, err := it.store.GetRow(context.Background(), it.db, it.table, pk)
		if err != nil {
			it.err = err
			return false
		}
		it.current = row
		return true
	}
}

func (it *indexRowIter) Row() Row     { return it.current }
func (it *indexRowIter) Error() error { return it.err }
func (it *indexRowIter) Close()       { it.iter.close() }
