package rdbengine

import (
	"bytes"
	"context"
	"encoding/json"
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

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	tblKey := catalogTableKey(s.engine, db, table)
	if !s.backend.has(tblKey) {
		return ErrNotFound
	}

	idxCatalogKey := catalogIndexKey(s.engine, db, table, indexName)
	if s.backend.has(idxCatalogKey) {
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
	iter := s.backend.newIter(prefix, end)
	defer iter.close()

	for iter.first(); iter.valid(); iter.next() {
		pk := iter.key()[len(prefix):]
		row, decodeErr := decodeRow(iter.value())
		if decodeErr != nil {
			return fmtErr("create_index backfill decode", decodeErr)
		}

		colVal := s.encodeIndexColumns(def, row)
		iKey := indexKey(s.engine, db, table, indexName, colVal, pk)
		batch.put(iKey, pk)

		if unique {
			uKey := uniqueKey(s.engine, db, table, indexName, colVal)
			existing, getErr := s.backend.get(uKey)
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

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	idxCatalogKey := catalogIndexKey(s.engine, db, table, indexName)
	if !s.backend.has(idxCatalogKey) {
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
func (s *Store) checkUniqueConstraints(db, table string, pk []byte, row Row) error {
	indexes, err := s.listIndexDefs(db, table)
	if err != nil {
		return err
	}
	for _, idx := range indexes {
		if !idx.Unique {
			continue
		}
		colVal := s.encodeIndexColumns(idx, row)
		uKey := uniqueKey(s.engine, db, table, idx.Name, colVal)
		existing, err := s.backend.get(uKey)
		if err != nil {
			return err
		}
		if existing != nil && !bytes.Equal(existing, pk) {
			return ErrConstraintViolation
		}
	}
	return nil
}

func (s *Store) appendIndexEntries(batch kvBatch, db, table string, pk []byte, row Row) {
	indexes, err := s.listIndexDefs(db, table)
	if err != nil || len(indexes) == 0 {
		return
	}
	for _, idx := range indexes {
		colVal := s.encodeIndexColumns(idx, row)
		iKey := indexKey(s.engine, db, table, idx.Name, colVal, pk)
		batch.put(iKey, pk)
		if idx.Unique {
			uKey := uniqueKey(s.engine, db, table, idx.Name, colVal)
			batch.put(uKey, pk)
		}
	}
}

func (s *Store) appendRemoveIndexEntries(batch kvBatch, db, table string, pk []byte, row Row) {
	indexes, err := s.listIndexDefs(db, table)
	if err != nil || len(indexes) == 0 {
		return
	}
	for _, idx := range indexes {
		colVal := s.encodeIndexColumns(idx, row)
		iKey := indexKey(s.engine, db, table, idx.Name, colVal, pk)
		batch.del(iKey)
		if idx.Unique {
			uKey := uniqueKey(s.engine, db, table, idx.Name, colVal)
			batch.del(uKey)
		}
	}
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
		iter:     s.backend.newIter(lower, upper),
		limit:    opts.Limit,
		offset:   opts.Offset,
		consumed: 0,
	}, nil
}

// encodeIndexColumns encodes the indexed columns from a row into a single key component.
func (s *Store) encodeIndexColumns(idx IndexDef, row Row) []byte {
	var buf []byte
	first := true
	for _, col := range idx.Columns {
		if !first {
			buf = append(buf, 0x00)
		}
		first = false
		cv, ok := row[col]
		if !ok || cv.Value == nil {
			continue
		}
		buf = append(buf, encodeValue(cv)...)
	}
	return buf
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
