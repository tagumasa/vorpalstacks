package vmysql

import (
	"context"
	"io"
	"sync"

	"github.com/dolthub/go-mysql-server/sql"

	"vorpalstacks/internal/core/storage/rdbengine"
)

// pebbleDatabase implements sql.Database backed by rdbengine.
// Satisfies: Database, TableCreator, TableDropper, MutableDatabaseProvider (via provider)
type pebbleDatabase struct {
	name  string
	store *rdbengine.Store
}

func newPebbleDatabase(name string, store *rdbengine.Store) *pebbleDatabase {
	return &pebbleDatabase{name: name, store: store}
}

func (d *pebbleDatabase) Name() string {
	return d.name
}

func (d *pebbleDatabase) GetTableInsensitive(ctx *sql.Context, tblName string) (sql.Table, bool, error) {
	schema, err := d.store.GetTableSchema(context.Background(), d.name, tblName)
	if err != nil {
		if err == rdbengine.ErrNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	tbl, err := newPebbleTable(d.name, schema, d.store)
	if err != nil {
		return nil, false, err
	}
	return tbl, true, nil
}

func (d *pebbleDatabase) GetTableNames(ctx *sql.Context) ([]string, error) {
	return d.store.ListTables(context.Background(), d.name)
}

func (d *pebbleDatabase) CreateTable(ctx *sql.Context, name string, schema sql.PrimaryKeySchema, collation sql.CollationID, comment string) error {
	cols := make([]rdbengine.ColumnDef, len(schema.Schema))
	for i, col := range schema.Schema {
		cols[i] = sqlColumnToDef(col)
	}
	tblSchema := &rdbengine.TableSchema{
		Name:    name,
		Columns: cols,
	}
	return d.store.CreateTable(context.Background(), d.name, tblSchema)
}

func (d *pebbleDatabase) DropTable(ctx *sql.Context, name string) error {
	return d.store.DropTable(context.Background(), d.name, name)
}

// pebbleProvider implements sql.MutableDatabaseProvider backed by rdbengine.
type pebbleProvider struct {
	store *rdbengine.Store
	mu    sync.RWMutex
	dbs   map[string]*pebbleDatabase
}

func newPebbleProvider(store *rdbengine.Store) *pebbleProvider {
	return &pebbleProvider{
		store: store,
		dbs:   make(map[string]*pebbleDatabase),
	}
}

func (p *pebbleProvider) Database(ctx *sql.Context, name string) (sql.Database, error) {
	p.mu.RLock()
	if db, ok := p.dbs[name]; ok {
		p.mu.RUnlock()
		return db, nil
	}
	p.mu.RUnlock()

	names, err := p.store.ListDatabases(context.Background())
	if err != nil {
		return nil, err
	}
	for _, n := range names {
		if n == name {
			db := newPebbleDatabase(name, p.store)
			p.mu.Lock()
			p.dbs[name] = db
			p.mu.Unlock()
			return db, nil
		}
	}
	return nil, sql.ErrDatabaseNotFound.New(name)
}

func (p *pebbleProvider) HasDatabase(ctx *sql.Context, name string) bool {
	_, err := p.Database(ctx, name)
	return err == nil
}

func (p *pebbleProvider) AllDatabases(ctx *sql.Context) []sql.Database {
	names, err := p.store.ListDatabases(context.Background())
	if err != nil {
		return nil
	}
	var result []sql.Database
	for _, name := range names {
		p.mu.RLock()
		db, ok := p.dbs[name]
		p.mu.RUnlock()
		if !ok {
			db = newPebbleDatabase(name, p.store)
			p.mu.Lock()
			p.dbs[name] = db
			p.mu.Unlock()
		}
		result = append(result, db)
	}
	return result
}

func (p *pebbleProvider) CreateDatabase(ctx *sql.Context, name string) error {
	if err := p.store.CreateDatabase(context.Background(), name); err != nil {
		return err
	}
	db := newPebbleDatabase(name, p.store)
	p.mu.Lock()
	p.dbs[name] = db
	p.mu.Unlock()
	return nil
}

func (p *pebbleProvider) DropDatabase(ctx *sql.Context, name string) error {
	if err := p.store.DropDatabase(context.Background(), name); err != nil {
		return err
	}
	p.mu.Lock()
	delete(p.dbs, name)
	p.mu.Unlock()
	return nil
}

func sqlColumnToDef(col *sql.Column) rdbengine.ColumnDef {
	ct := sqlTypeToColumnType(col.Type)
	var defVal *string
	if col.Default != nil {
		s := col.Default.String()
		defVal = &s
	}
	return rdbengine.ColumnDef{
		Name:         col.Name,
		Type:         ct,
		Nullable:     col.Nullable,
		PrimaryKey:   col.PrimaryKey,
		AutoIncr:     col.AutoIncrement,
		DefaultValue: defVal,
	}
}

// rowIterWrapper wraps rdbengine.RowIterator to satisfy sql.RowIter.
type rowIterWrapper struct {
	iter rdbengine.RowIterator
	sch  sql.Schema
}

func (r *rowIterWrapper) Next(ctx *sql.Context) (sql.Row, error) {
	if !r.iter.Next() {
		if err := r.iter.Error(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
	row := r.iter.Row()
	sqlRow := make(sql.Row, len(r.sch))
	for i, col := range r.sch {
		cv, ok := row[col.Name]
		if !ok || cv.Value == nil {
			sqlRow[i] = nil
			continue
		}
		sqlRow[i] = cv.Value
	}
	return sqlRow, nil
}

func (r *rowIterWrapper) Close(ctx *sql.Context) error {
	r.iter.Close()
	return nil
}
