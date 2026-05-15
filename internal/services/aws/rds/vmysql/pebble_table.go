package vmysql

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/dolthub/vitess/go/vt/proto/query"

	"vorpalstacks/internal/core/storage/rdbengine"
)

// sqlTypeToColumnType converts a go-mysql-server sql.Type to rdbengine.ColumnType.
func sqlTypeToColumnType(t sql.Type) rdbengine.ColumnType {
	switch t.Type() {
	case query.Type_INT8:
		return rdbengine.ColumnTypeInt8
	case query.Type_INT16:
		return rdbengine.ColumnTypeInt16
	case query.Type_INT24:
		return rdbengine.ColumnTypeInt32
	case query.Type_INT32:
		return rdbengine.ColumnTypeInt32
	case query.Type_INT64:
		return rdbengine.ColumnTypeInt64
	case query.Type_UINT8:
		return rdbengine.ColumnTypeInt8
	case query.Type_UINT16:
		return rdbengine.ColumnTypeInt16
	case query.Type_UINT24:
		return rdbengine.ColumnTypeInt32
	case query.Type_UINT32:
		return rdbengine.ColumnTypeInt32
	case query.Type_UINT64:
		return rdbengine.ColumnTypeInt64
	case query.Type_FLOAT32:
		return rdbengine.ColumnTypeFloat32
	case query.Type_FLOAT64:
		return rdbengine.ColumnTypeFloat64
	case query.Type_DECIMAL:
		return rdbengine.ColumnTypeDecimal
	case query.Type_CHAR, query.Type_VARCHAR, query.Type_SET, query.Type_ENUM:
		return rdbengine.ColumnTypeString
	case query.Type_TEXT:
		return rdbengine.ColumnTypeText
	case query.Type_BLOB:
		return rdbengine.ColumnTypeBlob
	case query.Type_BINARY, query.Type_VARBINARY:
		return rdbengine.ColumnTypeBytes
	case query.Type_DATE:
		return rdbengine.ColumnTypeDate
	case query.Type_DATETIME, query.Type_TIMESTAMP:
		return rdbengine.ColumnTypeTimestamp
	case query.Type_JSON:
		return rdbengine.ColumnTypeJSON
	case query.Type_BIT:
		return rdbengine.ColumnTypeBool
	case query.Type_NULL_TYPE:
		return rdbengine.ColumnTypeUnknown
	default:
		return rdbengine.ColumnTypeString
	}
}

// columnTypeToSQLType converts an rdbengine.ColumnType back to go-mysql-server sql.Type.
func columnTypeToSQLType(ct rdbengine.ColumnType) sql.Type {
	switch ct {
	case rdbengine.ColumnTypeBool:
		return types.Boolean
	case rdbengine.ColumnTypeInt8:
		return types.Int8
	case rdbengine.ColumnTypeInt16:
		return types.Int16
	case rdbengine.ColumnTypeInt32, rdbengine.ColumnTypeSerial:
		return types.Int32
	case rdbengine.ColumnTypeInt64, rdbengine.ColumnTypeBigSerial:
		return types.Int64
	case rdbengine.ColumnTypeFloat32:
		return types.Float32
	case rdbengine.ColumnTypeFloat64:
		return types.Float64
	case rdbengine.ColumnTypeDecimal:
		return types.MustCreateDecimalType(10, 0)
	case rdbengine.ColumnTypeString:
		return types.MustCreateStringWithDefaults(query.Type_VARCHAR, 255)
	case rdbengine.ColumnTypeText:
		return types.MustCreateStringWithDefaults(query.Type_TEXT, 65535)
	case rdbengine.ColumnTypeBytes, rdbengine.ColumnTypeBlob:
		return types.Blob
	case rdbengine.ColumnTypeDate:
		return types.Date
	case rdbengine.ColumnTypeTimestamp, rdbengine.ColumnTypeTimestampTZ:
		return types.Datetime
	case rdbengine.ColumnTypeJSON:
		return types.JSON
	case rdbengine.ColumnTypeUUID:
		return types.MustCreateStringWithDefaults(query.Type_VARCHAR, 36)
	default:
		return types.MustCreateStringWithDefaults(query.Type_VARCHAR, 255)
	}
}

// defToSQLColumn converts an rdbengine.ColumnDef to go-mysql-server sql.Column.
func defToSQLColumn(d rdbengine.ColumnDef) *sql.Column {
	return &sql.Column{
		Name:          d.Name,
		Type:          columnTypeToSQLType(d.Type),
		Nullable:      d.Nullable,
		PrimaryKey:    d.PrimaryKey,
		AutoIncrement: d.AutoIncr,
	}
}

// schemaToSQLSchema converts an rdbengine.TableSchema to go-mysql-server sql.Schema.
func schemaToSQLSchema(s *rdbengine.TableSchema) sql.Schema {
	cols := make(sql.Schema, len(s.Columns))
	for i, c := range s.Columns {
		cols[i] = defToSQLColumn(c)
	}
	return cols
}

// pebbleTable implements sql.Table backed by rdbengine, with full
// insert/update/delete support.
type pebbleTable struct {
	dbName string
	name   string
	schema *rdbengine.TableSchema
	store  *rdbengine.Store
	sqlSch sql.Schema
}

func newPebbleTable(dbName string, schema *rdbengine.TableSchema, store *rdbengine.Store) (*pebbleTable, error) {
	sqlSch := schemaToSQLSchema(schema)
	for _, col := range sqlSch {
		col.Source = schema.Name
		col.DatabaseSource = dbName
	}
	return &pebbleTable{
		dbName: dbName,
		name:   schema.Name,
		schema: schema,
		store:  store,
		sqlSch: sqlSch,
	}, nil
}

func (t *pebbleTable) Name() string               { return t.name }
func (t *pebbleTable) String() string             { return t.name }
func (t *pebbleTable) Schema() sql.Schema         { return t.sqlSch }
func (t *pebbleTable) Collation() sql.CollationID { return sql.Collation_Default }

type singlePartition struct{}

func (singlePartition) Key() []byte { return []byte("single") }

func (t *pebbleTable) Partitions(ctx *sql.Context) (sql.PartitionIter, error) {
	return sql.PartitionsToPartitionIter(singlePartition{}), nil
}
func (t *pebbleTable) PartitionRows(ctx *sql.Context, p sql.Partition) (sql.RowIter, error) {
	iter, err := t.store.ScanRows(ctx, t.dbName, t.name, rdbengine.ScanOptions{})
	if err != nil {
		return nil, err
	}
	return &rowIterWrapper{iter: iter, sch: t.sqlSch}, nil
}

// --- InsertableTable ---

func (t *pebbleTable) Inserter(ctx *sql.Context) sql.RowInserter {
	return &pebbleInserter{table: t, ctx: ctx}
}

type pebbleInserter struct {
	table *pebbleTable
	ctx   *sql.Context
}

func (i *pebbleInserter) StatementBegin(ctx *sql.Context)                  {}
func (i *pebbleInserter) DiscardChanges(ctx *sql.Context, err error) error { return nil }
func (i *pebbleInserter) StatementComplete(ctx *sql.Context) error         { return nil }
func (i *pebbleInserter) Insert(ctx *sql.Context, row sql.Row) error {
	pk := encodeSQLPK(i.table.schema, row)
	engineRow := sqlRowToEngineRow(row, i.table.sqlSch, i.table.schema)
	return i.table.store.InsertRow(ctx, i.table.dbName, i.table.name, pk, engineRow)
}
func (i *pebbleInserter) Close(ctx *sql.Context) error { return nil }

// --- DeletableTable ---

func (t *pebbleTable) Deleter(ctx *sql.Context) sql.RowDeleter {
	return &pebbleDeleter{table: t, ctx: ctx}
}

type pebbleDeleter struct {
	table *pebbleTable
	ctx   *sql.Context
}

func (d *pebbleDeleter) StatementBegin(ctx *sql.Context)                  {}
func (d *pebbleDeleter) DiscardChanges(ctx *sql.Context, err error) error { return nil }
func (d *pebbleDeleter) StatementComplete(ctx *sql.Context) error         { return nil }
func (d *pebbleDeleter) Delete(ctx *sql.Context, row sql.Row) error {
	pk := encodeSQLPK(d.table.schema, row)
	return d.table.store.DeleteRow(ctx, d.table.dbName, d.table.name, pk)
}
func (d *pebbleDeleter) Close(ctx *sql.Context) error { return nil }

// --- UpdatableTable ---

func (t *pebbleTable) Updater(ctx *sql.Context) sql.RowUpdater {
	return &pebbleUpdater{table: t, ctx: ctx}
}

type pebbleUpdater struct {
	table *pebbleTable
	ctx   *sql.Context
}

func (u *pebbleUpdater) StatementBegin(ctx *sql.Context)                  {}
func (u *pebbleUpdater) DiscardChanges(ctx *sql.Context, err error) error { return nil }
func (u *pebbleUpdater) StatementComplete(ctx *sql.Context) error         { return nil }
func (u *pebbleUpdater) Update(ctx *sql.Context, old sql.Row, new sql.Row) error {
	pk := encodeSQLPK(u.table.schema, old)
	engineRow := sqlRowToEngineRow(new, u.table.sqlSch, u.table.schema)
	return u.table.store.UpdateRow(ctx, u.table.dbName, u.table.name, pk, engineRow)
}
func (u *pebbleUpdater) Close(ctx *sql.Context) error { return nil }

// --- AutoIncrementTable ---

func (t *pebbleTable) GetNextAutoIncrementValue(ctx *sql.Context, insertVal interface{}) (uint64, error) {
	return t.store.NextAutoIncrement(t.dbName, t.name)
}

func (t *pebbleTable) PeekNextAutoIncrementValue(ctx *sql.Context) (uint64, error) {
	return t.store.GetAutoIncrement(t.dbName, t.name)
}

func (t *pebbleTable) AutoIncrementSetter(ctx *sql.Context) sql.AutoIncrementSetter {
	return &pebbleAutoIncSetter{table: t}
}

type pebbleAutoIncSetter struct {
	table *pebbleTable
}

func (s *pebbleAutoIncSetter) SetAutoIncrementValue(ctx *sql.Context, val uint64) error {
	return s.table.store.SetAutoIncrement(s.table.dbName, s.table.name, val)
}

func (s *pebbleAutoIncSetter) AcquireAutoIncrementLock(ctx *sql.Context) (func(), error) {
	mu := s.table.store.TableLock(s.table.dbName, s.table.name)
	mu.Lock()
	return func() { mu.Unlock() }, nil
}

func (s *pebbleAutoIncSetter) Close(ctx *sql.Context) error { return nil }

// encodeSQLPK extracts the primary key from a sql.Row and encodes it.
// go-mysql-server returns different Go types for the same column across
// operations (int32 on INSERT, int64 on UPDATE); normalise before encoding.
func encodeSQLPK(schema *rdbengine.TableSchema, row sql.Row) []byte {
	engineRow := make(rdbengine.Row)
	for i, col := range schema.Columns {
		if i < len(row) && row[i] != nil {
			engineRow[col.Name] = rdbengine.ColumnValue{
				Type:  col.Type,
				Value: normaliseSQLValue(row[i], col.Type),
			}
		}
	}
	return rdbengine.EncodePK(schema, engineRow)
}

func normaliseSQLValue(val interface{}, ct rdbengine.ColumnType) interface{} {
	switch ct {
	case rdbengine.ColumnTypeInt8, rdbengine.ColumnTypeInt16, rdbengine.ColumnTypeInt32,
		rdbengine.ColumnTypeSerial:
		switch v := val.(type) {
		case int:
			return int32(v)
		case int8:
			return int32(v)
		case int16:
			return int32(v)
		case int32:
			return v
		case int64:
			return int32(v)
		case uint:
			return int32(v)
		case uint8:
			return int32(v)
		case uint16:
			return int32(v)
		case uint32:
			return int32(v)
		case uint64:
			return int32(v)
		}
	case rdbengine.ColumnTypeInt64, rdbengine.ColumnTypeBigSerial:
		switch v := val.(type) {
		case int:
			return int64(v)
		case int8:
			return int64(v)
		case int16:
			return int64(v)
		case int32:
			return int64(v)
		case int64:
			return v
		case uint:
			return int64(v)
		case uint8:
			return int64(v)
		case uint16:
			return int64(v)
		case uint32:
			return int64(v)
		case uint64:
			return int64(v)
		}
	case rdbengine.ColumnTypeFloat32:
		switch v := val.(type) {
		case float32:
			return v
		case float64:
			return float32(v)
		}
	}
	return val
}

// sqlRowToEngineRow converts a sql.Row to rdbengine.Row.
func sqlRowToEngineRow(row sql.Row, sqlSch sql.Schema, schema *rdbengine.TableSchema) rdbengine.Row {
	engineRow := make(rdbengine.Row)
	for i, col := range sqlSch {
		if i >= len(row) {
			continue
		}
		def := schema.ColumnByName(col.Name)
		if def == nil {
			continue
		}
		var val interface{}
		if row[i] != nil {
			val = normaliseSQLValue(row[i], def.Type)
		}
		engineRow[col.Name] = rdbengine.ColumnValue{
			Type:  def.Type,
			Value: val,
		}
	}
	return engineRow
}
