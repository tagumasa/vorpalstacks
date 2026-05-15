package rdbengine

import "context"

// RowReader provides read-only access to relational data, decoupled from
// the concrete Pebble-backed storage. SQL engine adapters depend on this
// interface rather than *Store directly.
type RowReader interface {
	GetRow(ctx context.Context, db, table string, pk []byte) (Row, error)
	ScanRows(ctx context.Context, db, table string, opts ScanOptions) (RowIterator, error)
	ScanIndex(ctx context.Context, db, table, indexName string, opts IndexScanOptions) (RowIterator, error)
	GetTableSchema(ctx context.Context, db, table string) (*TableSchema, error)
	ListTables(ctx context.Context, db string) ([]string, error)
	ListDatabases(ctx context.Context) ([]string, error)
}

// RowWriter provides write access to relational data.
type RowWriter interface {
	InsertRow(ctx context.Context, db, table string, pk []byte, row Row) error
	UpdateRow(ctx context.Context, db, table string, pk []byte, row Row) error
	DeleteRow(ctx context.Context, db, table string, pk []byte) error
	CreateTable(ctx context.Context, db string, schema *TableSchema) error
	DropTable(ctx context.Context, db, table string) error
	CreateDatabase(ctx context.Context, db string) error
	DropDatabase(ctx context.Context, db string) error
	CreateIndex(ctx context.Context, db, table, indexName string, columns []string, unique bool) error
	DropIndex(ctx context.Context, db, table, indexName string) error
}

// RowStore combines read and write access with lifecycle management.
type RowStore interface {
	RowReader
	RowWriter
	Close() error
}
