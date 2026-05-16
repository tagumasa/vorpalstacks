package rdbengine

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

// Catalog operations for database/table/column metadata stored in Pebble.
// Keys: catalog/{engine}/db/{db_name}               → DatabaseMeta JSON
//       catalog/{engine}/db/{db_name}/table/{tbl}    → TableSchema JSON
//       catalog/{engine}/db/{db_name}/table/{tbl}/idx/{idx} → IndexDef JSON

// CreateDatabase creates a new database entry in the catalog.
func (s *Store) CreateDatabase(ctx context.Context, db string) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	key := catalogDBKey(s.engine, db)
	if s.backend.has(key) {
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
func (s *Store) DropDatabase(ctx context.Context, db string) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	key := catalogDBKey(s.engine, db)
	if !s.backend.has(key) {
		return ErrNotFound
	}

	tables, err := s.ListTables(ctx, db)
	if err != nil {
		return err
	}
	for _, tbl := range tables {
		if err := s.DropTable(ctx, db, tbl); err != nil {
			return err
		}
	}

	return s.backend.delete(key)
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

	it := s.backend.newIter(prefix, end)
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
	dbKey := catalogDBKey(s.engine, db)
	if !s.backend.has(dbKey) {
		return ErrNotFound
	}

	tblKey := catalogTableKey(s.engine, db, schema.Name)
	if s.backend.has(tblKey) {
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
func (s *Store) DropTable(ctx context.Context, db, table string) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	tblKey := catalogTableKey(s.engine, db, table)
	if !s.backend.has(tblKey) {
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
	data, err := s.backend.get(key)
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
	if !s.backend.has(dbKey) {
		return nil, ErrNotFound
	}

	prefix := catalogTablePrefix(s.engine, db)
	end := make([]byte, len(prefix)+1)
	copy(end, prefix)
	end[len(prefix)] = 0xFF

	it := s.backend.newIter(prefix, end)
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

// getDatabaseMeta returns the metadata for a database.
func (s *Store) getDatabaseMeta(db string) (*DatabaseMeta, error) {
	key := catalogDBKey(s.engine, db)
	data, err := s.backend.get(key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrNotFound
	}
	var meta DatabaseMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// listIndexDefs returns all index definitions for a table.
func (s *Store) listIndexDefs(db, table string) ([]IndexDef, error) {
	prefix := catalogIndexPrefix(s.engine, db, table)
	end := make([]byte, len(prefix)+1)
	copy(end, prefix)
	end[len(prefix)] = 0xFF

	it := s.backend.newIter(prefix, end)
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
