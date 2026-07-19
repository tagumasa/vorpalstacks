package rdbengine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Catalog operations for database/table/column metadata stored in Pebble.
// Keys: catalog/{engine}/db/{db_name}               → DatabaseMeta JSON
//       catalog/{engine}/db/{db_name}/table/{tbl}    → TableSchema JSON
//       catalog/{engine}/db/{db_name}/table/{tbl}/idx/{idx} → IndexDef JSON

// Permitted identifier constraints. The 64-byte limit mirrors MySQL's
// NAME_LEN constant; the character restrictions match MySQL's
// my_charset_latin1 + identifier-quote rules and also rule out the
// path separator '/' and control bytes 0x00 / 0xFF that would corrupt
// Pebble key encoding (catalogDBKey, rowKeyPrefix, indexEndKey, ...).
//
// Without these checks a malicious or careless caller could craft a
// database name like "foo/table/bar" whose catalogDBKey collides with
// the catalogTableKey of table "bar" in database "foo", silently
// corrupting both entries (C-4).
const (
	identifierMaxBytes = 64
)

var (
	// ErrInvalidIdentifier is returned when a database, table, or index
	// name fails ValidateIdentifier.
	ErrInvalidIdentifier = fmt.Errorf("rdbengine: invalid identifier")
)

// ValidateIdentifier enforces the rules common to database, table, and
// index names: non-empty, at most 64 bytes of UTF-8, no path separator
// or NUL/0xFF sentinel bytes, and no leading/trailing whitespace.
// Apply this at every catalog-mutating entry point so that Pebble key
// construction remains unambiguous.
func ValidateIdentifier(kind, name string) error {
	if name == "" {
		return fmt.Errorf("%w: %s name is empty", ErrInvalidIdentifier, kind)
	}
	if len(name) > identifierMaxBytes {
		return fmt.Errorf("%w: %s name %q exceeds %d bytes",
			ErrInvalidIdentifier, kind, name, identifierMaxBytes)
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		if c == 0x00 {
			return fmt.Errorf("%w: %s name %q contains NUL byte",
				ErrInvalidIdentifier, kind, name)
		}
		if c == 0xFF {
			return fmt.Errorf("%w: %s name %q contains 0xFF sentinel byte",
				ErrInvalidIdentifier, kind, name)
		}
		if c == '/' {
			return fmt.Errorf("%w: %s name %q contains path separator '/'",
				ErrInvalidIdentifier, kind, name)
		}
	}
	if strings.TrimSpace(name) != name {
		return fmt.Errorf("%w: %s name %q has leading or trailing whitespace",
			ErrInvalidIdentifier, kind, name)
	}
	return nil
}

// CreateDatabase creates a new database entry in the catalog.
func (s *Store) CreateDatabase(ctx context.Context, db string) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}
	if err := ValidateIdentifier("database", db); err != nil {
		return err
	}

	dbMu := s.DatabaseLock(db)
	dbMu.Lock()
	defer dbMu.Unlock()

	key := catalogDBKey(s.engine, db)
	if s.has(key) {
		return ErrAlreadyExists
	}
	meta := DatabaseMeta{
		Name:      db,
		Engine:    s.engine,
		CreatedAt: time.Now().UTC(),
	}
	data, err := json.Marshal(meta)
	if err != nil {
		return fmtErr("create_database marshal", err)
	}
	return s.backend.set(key, data)
}

// DropDatabase removes a database and all its tables/indexes from the catalog
// and deletes all row data and index entries for that database.
//
// Atomicity: every mutation (database meta, every table's schema, every
// index definition, every row range, every index range, every unique
// range) is appended to a single Pebble batch that is committed exactly
// once. A crash between commits therefore cannot leave the catalog
// half-deleted with orphaned tables pointing at a missing database.
//
// Concurrency: the database-level mutex is acquired first, then every
// table-level mutex. This blocks concurrent InsertRow / UpdateRow /
// DeleteRow on any table from committing after the deleteRange batch
// fires, which would leave orphaned row data. CreateTable also
// acquires DatabaseLock, so no new tables can appear during the
// iteration. TableLock-only callers (InsertRow, UpdateRow, DeleteRow,
// DropTable) never acquire DatabaseLock, so this ordering cannot
// deadlock.
func (s *Store) DropDatabase(ctx context.Context, db string) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}
	if err := ValidateIdentifier("database", db); err != nil {
		return err
	}

	dbMu := s.DatabaseLock(db)
	dbMu.Lock()
	defer dbMu.Unlock()

	dbKey := catalogDBKey(s.engine, db)
	if !s.has(dbKey) {
		return ErrNotFound
	}

	tables, err := s.ListTables(ctx, db)
	if err != nil {
		return err
	}

	// Acquire every table-level lock so that a concurrent InsertRow /
	// UpdateRow / DeleteRow / DropTable on any table in this database
	// cannot commit after the deleteRange batch fires. Without these
	// locks, a row committed by a concurrent writer after DropDatabase's
	// batch commit would survive as unreachable orphan data (catalog
	// says the table does not exist, but the row bytes remain in Pebble).
	//
	// Lock ordering: DatabaseLock first (already held), then all
	// TableLocks. TableLock-only callers never acquire DatabaseLock,
	// so this cannot deadlock.
	tableMus := make([]*sync.Mutex, len(tables))
	for i, tbl := range tables {
		tableMus[i] = s.TableLock(db, tbl)
		tableMus[i].Lock()
	}
	defer func() {
		for _, mu := range tableMus {
			mu.Unlock()
		}
	}()

	batch := s.backend.newBatch()
	defer batch.close()

	for _, tbl := range tables {
		tblKey := catalogTableKey(s.engine, db, tbl)
		batch.del(tblKey)

		indexes, idxErr := s.listIndexDefs(db, tbl)
		if idxErr != nil {
			return idxErr
		}
		for _, idx := range indexes {
			idxKey := catalogIndexKey(s.engine, db, tbl, idx.Name)
			batch.del(idxKey)
			idxDataStart := indexKeyPrefix(s.engine, db, tbl, idx.Name)
			idxDataEnd := indexEndKey(s.engine, db, tbl, idx.Name)
			batch.deleteRange(idxDataStart, idxDataEnd)
			if idx.Unique {
				uniqStart := uniqueKeyPrefix(s.engine, db, tbl, idx.Name)
				uniqEnd := uniqueEndKey(s.engine, db, tbl, idx.Name)
				batch.deleteRange(uniqStart, uniqEnd)
			}
		}

		rowStart := rowKeyPrefix(s.engine, db, tbl)
		rowEnd := rowEndKey(s.engine, db, tbl)
		batch.deleteRange(rowStart, rowEnd)
	}

	batch.del(dbKey)
	if err := batch.commit(); err != nil {
		return fmtErr("drop_database commit", err)
	}
	return nil
}

// ListDatabases returns all database names for the current engine.
func (s *Store) ListDatabases(ctx context.Context) ([]string, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	prefix := catalogDBPrefix(s.engine)
	end := make([]byte, len(prefix)+1)
	copy(end, prefix)
	end[len(prefix)] = 0xFF

	suffix := "/"

	it := s.newIter(prefix, end)
	defer it.close()

	var names []string
	for it.first(); it.valid(); it.next() {
		k := it.key()
		rest := string(k[len(prefix):])
		if idx := strings.Index(rest, suffix); idx >= 0 {
			rest = rest[:idx]
		}
		if rest != "" && !strings.Contains(rest, "/") {
			names = append(names, rest)
		}
	}
	if err := it.err(); err != nil {
		return nil, fmtErr("list_databases", err)
	}
	return names, nil
}

// CreateTable creates a new table with the given schema in the catalog.
func (s *Store) CreateTable(ctx context.Context, db string, schema *TableSchema) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if s.snap != nil {
		return errReadOnly
	}
	if err := ValidateIdentifier("database", db); err != nil {
		return err
	}
	if schema == nil {
		return fmt.Errorf("rdbengine: create_table: nil schema")
	}
	if err := ValidateIdentifier("table", schema.Name); err != nil {
		return err
	}
	for _, col := range schema.Columns {
		if err := ValidateIdentifier("column", col.Name); err != nil {
			return err
		}
	}

	// Acquire DatabaseLock first so that a concurrent DropDatabase
	// cannot race: if DropDatabase has already committed, the dbKey
	// check below will fail; if DropDatabase has not yet listed tables,
	// it will see this table in its iteration. Without DatabaseLock,
	// CreateTable could create a new table after DropDatabase's table
	// list is frozen, leaving an orphan catalog entry pointing at a
	// deleted database.
	dbMu := s.DatabaseLock(db)
	dbMu.Lock()
	defer dbMu.Unlock()

	dbKey := catalogDBKey(s.engine, db)
	if !s.has(dbKey) {
		return ErrNotFound
	}

	tblMu := s.TableLock(db, schema.Name)
	tblMu.Lock()
	defer tblMu.Unlock()

	tblKey := catalogTableKey(s.engine, db, schema.Name)
	if s.has(tblKey) {
		return ErrAlreadyExists
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return fmtErr("create_table marshal", err)
	}

	batch := s.backend.newBatch()
	defer batch.close()
	batch.put(tblKey, data)
	return batch.commit()
}

// DropTable removes a table schema, all its indexes, and all row data.
// TableLock is acquired so that a concurrent InsertRow / UpdateRow /
// DeleteRow cannot commit after the deleteRange batch fires, which
// would leave orphaned row data in Pebble.
func (s *Store) DropTable(ctx context.Context, db, table string) error {
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

	mu := s.TableLock(db, table)
	mu.Lock()
	defer mu.Unlock()

	tblKey := catalogTableKey(s.engine, db, table)
	if !s.has(tblKey) {
		return ErrNotFound
	}

	indexes, err := s.listIndexDefs(db, table)
	if err != nil {
		return err
	}

	batch := s.backend.newBatch()
	defer batch.close()

	batch.del(tblKey)

	for _, idx := range indexes {
		idxKey := catalogIndexKey(s.engine, db, table, idx.Name)
		batch.del(idxKey)
		idxDataStart := indexKeyPrefix(s.engine, db, table, idx.Name)
		idxDataEnd := indexEndKey(s.engine, db, table, idx.Name)
		batch.deleteRange(idxDataStart, idxDataEnd)
		if idx.Unique {
			uniqStart := uniqueKeyPrefix(s.engine, db, table, idx.Name)
			uniqEnd := uniqueEndKey(s.engine, db, table, idx.Name)
			batch.deleteRange(uniqStart, uniqEnd)
		}
	}

	rowStart := rowKeyPrefix(s.engine, db, table)
	rowEnd := rowEndKey(s.engine, db, table)
	batch.deleteRange(rowStart, rowEnd)

	return batch.commit()
}

// GetTableSchema returns the schema for a specific table.
func (s *Store) GetTableSchema(ctx context.Context, db, table string) (*TableSchema, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	key := catalogTableKey(s.engine, db, table)
	data, err := s.get(key)
	if err != nil {
		return nil, fmtErr("get_table_schema", err)
	}
	if data == nil {
		return nil, ErrNotFound
	}
	var schema TableSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmtErr("get_table_schema unmarshal", err)
	}
	return &schema, nil
}

// ListTables returns all table names in a database.
func (s *Store) ListTables(ctx context.Context, db string) ([]string, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	dbKey := catalogDBKey(s.engine, db)
	if !s.has(dbKey) {
		return nil, ErrNotFound
	}

	prefix := catalogTablePrefix(s.engine, db)
	end := make([]byte, len(prefix)+1)
	copy(end, prefix)
	end[len(prefix)] = 0xFF

	it := s.newIter(prefix, end)
	defer it.close()

	var names []string
	for it.first(); it.valid(); it.next() {
		k := it.key()
		rest := string(k[len(prefix):])
		if idx := strings.Index(rest, "/"); idx >= 0 {
			rest = rest[:idx]
		}
		if rest != "" {
			names = append(names, rest)
		}
	}
	if err := it.err(); err != nil {
		return nil, fmtErr("list_tables", err)
	}
	return names, nil
}

// listIndexDefs returns all index definitions for a table.
func (s *Store) listIndexDefs(db, table string) ([]IndexDef, error) {
	prefix := catalogIndexPrefix(s.engine, db, table)
	end := make([]byte, len(prefix)+1)
	copy(end, prefix)
	end[len(prefix)] = 0xFF

	it := s.newIter(prefix, end)
	defer it.close()

	var defs []IndexDef
	for it.first(); it.valid(); it.next() {
		var def IndexDef
		if err := json.Unmarshal(it.value(), &def); err != nil {
			continue
		}
		defs = append(defs, def)
	}
	if err := it.err(); err != nil {
		return nil, err
	}
	return defs, nil
}

// ListIndexes returns all index definitions for a table.
func (s *Store) ListIndexes(ctx context.Context, db, table string) ([]IndexDef, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	return s.listIndexDefs(db, table)
}
