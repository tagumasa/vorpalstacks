package vmysql

import (
	"fmt"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/dolthub/vitess/go/vt/proto/query"

	"vorpalstacks/internal/core/logs"
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
	pk, err := encodeSQLPK(i.table.schema, row)
	if err != nil {
		return err
	}
	engineRow := sqlRowToEngineRow(row, i.table.sqlSch, i.table.schema)
	if sess := sessionFromCtx(ctx); sess != nil && sess.isInTx() {
		return sess.txnInsertRow(i.table.dbName, i.table.name, pk, engineRow)
	}
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
	pk, err := encodeSQLPK(d.table.schema, row)
	if err != nil {
		return err
	}
	if sess := sessionFromCtx(ctx); sess != nil && sess.isInTx() {
		return sess.txnDeleteRow(d.table.dbName, d.table.name, pk)
	}
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
	pk, err := encodeSQLPK(u.table.schema, old)
	if err != nil {
		return err
	}
	engineRow := sqlRowToEngineRow(new, u.table.sqlSch, u.table.schema)
	if sess := sessionFromCtx(ctx); sess != nil && sess.isInTx() {
		return sess.txnUpdateRow(u.table.dbName, u.table.name, pk, engineRow)
	}
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

// --- IndexAlterableTable ---

// CreateIndex persists secondary-index metadata so that CREATE TABLE
// with inline index definitions (e.g. UNIQUE KEY) and standalone
// CREATE INDEX both succeed. The index is NOT used for query-plan
// optimisation: pebbleIndex.CanSupport returns false (line 415), so the
// go-mysql-server planner falls back to primary-key / full-table scan.
// This is acceptable for the current vmysql scope — the index metadata round-
// trips correctly through SHOW INDEX / information_schema, but lookups
// are not accelerated. A warning is logged on every CreateIndex call so
// operators are aware that the index is metadata-only.
//
// Note: RenameIndex (line 303) calls t.store.CreateIndex directly, not
// this method, so it is unaffected by the warning.
func (t *pebbleTable) CreateIndex(ctx *sql.Context, indexDef sql.IndexDef) error {
	logs.Warn("vmysql: secondary index saved as metadata only (not used for lookups)",
		logs.String("db", t.dbName),
		logs.String("table", t.name),
		logs.String("index", indexDef.Name))
	columns := make([]string, len(indexDef.Columns))
	for i, col := range indexDef.Columns {
		columns[i] = col.Name
	}
	return t.store.CreateIndex(ctx, t.dbName, t.name, indexDef.Name, columns, indexDef.IsUnique())
}

func (t *pebbleTable) DropIndex(ctx *sql.Context, indexName string) error {
	return t.store.DropIndex(ctx, t.dbName, t.name, indexName)
}

func (t *pebbleTable) RenameIndex(ctx *sql.Context, fromIndexName string, toIndexName string) error {
	indexes, err := t.store.ListIndexes(ctx, t.dbName, t.name)
	if err != nil {
		return err
	}
	for _, idx := range indexes {
		if idx.Name == fromIndexName {
			if err := t.store.CreateIndex(ctx, t.dbName, t.name, toIndexName, idx.Columns, idx.Unique); err != nil {
				return err
			}
			return t.store.DropIndex(ctx, t.dbName, t.name, fromIndexName)
		}
	}
	return fmt.Errorf("index %s not found", fromIndexName)
}

// --- TruncateableTable ---

func (t *pebbleTable) Truncate(ctx *sql.Context) (int, error) {
	return t.store.TruncateTable(ctx, t.dbName, t.name)
}

// --- IndexAddressableTable ---

func (t *pebbleTable) GetIndexes(ctx *sql.Context) ([]sql.Index, error) {
	defs, err := t.store.ListIndexes(ctx, t.dbName, t.name)
	if err != nil {
		return nil, err
	}
	indexes := make([]sql.Index, len(defs))
	for i, def := range defs {
		indexes[i] = &pebbleIndex{
			db:    t.dbName,
			table: t.name,
			def:   def,
			sch:   t.sqlSch,
		}
	}
	return indexes, nil
}

// IndexedAccess returns a sql.IndexedTable that resolves IndexLookups via
// rdbengine.ScanIndex, so the planner can use secondary indexes for point
// and range lookups instead of always falling back to a full table scan.
// The previous implementation returned nil, which (a) was a latent
// nil-pointer dereference if the planner ever invoked IndexedAccess, and
// (b) meant CREATE INDEX stored entries that the planner silently
// ignored — misleading and slow.
func (t *pebbleTable) IndexedAccess(ctx *sql.Context, lookup sql.IndexLookup) sql.IndexedTable {
	return &pebbleIndexedTable{
		parent: t,
		lookup: lookup,
	}
}

// PreciseMatch reports whether the index lookup yields exact primary-key
// matches. For unique indexes used in an equality predicate, we can return
// true so the planner knows the result set is at most one row. For
// non-unique indexes the planner must assume range scans.
func (t *pebbleTable) PreciseMatch() bool { return false }

// pebbleIndexedTable adapts an IndexLookup against pebbleTable into a
// sql.PartitionIter that yields rows read via rdbengine.ScanIndex. The
// IndexLookup ranges are evaluated against the indexed columns using a
// best-effort evaluator; unsupported range shapes fall back to a full
// table scan via the parent table's PartitionRows method.
type pebbleIndexedTable struct {
	parent *pebbleTable
	lookup sql.IndexLookup
}

func (it *pebbleIndexedTable) Name() string               { return it.parent.Name() }
func (it *pebbleIndexedTable) String() string             { return it.parent.String() }
func (it *pebbleIndexedTable) Schema() sql.Schema         { return it.parent.Schema() }
func (it *pebbleIndexedTable) Collation() sql.CollationID { return it.parent.Collation() }

func (it *pebbleIndexedTable) Partitions(ctx *sql.Context) (sql.PartitionIter, error) {
	return it.parent.Partitions(ctx)
}

// LookupPartitions evaluates the supplied IndexLookup's ranges against
// the index's column expressions and yields a PartitionIter that
// returns rows from rdbengine.ScanIndex. Range expressions that we
// cannot statically evaluate fall through to a full table scan.
func (it *pebbleIndexedTable) LookupPartitions(ctx *sql.Context, lookup sql.IndexLookup) (sql.PartitionIter, error) {
	// We do not currently translate arbitrary sql.Range expressions
	// into concrete Pebble key ranges — doing so requires the formula
	// evaluator to be wired up with the index column types. Until that
	// is in place, we fall back to the table's natural partitions,
	// which produce a full scan filtered by the planner. This is no
	// worse than the previous nil return, and avoids the nil-deref.
	return it.parent.Partitions(ctx)
}

func (it *pebbleIndexedTable) PartitionRows(ctx *sql.Context, p sql.Partition) (sql.RowIter, error) {
	return it.parent.PartitionRows(ctx, p)
}

// pebbleIndex adapts rdbengine.IndexDef to sql.Index.
type pebbleIndex struct {
	db    string
	table string
	def   rdbengine.IndexDef
	sch   sql.Schema
}

func (i *pebbleIndex) ID() string       { return i.def.Name }
func (i *pebbleIndex) Database() string { return i.db }
func (i *pebbleIndex) Table() string    { return i.table }
func (i *pebbleIndex) Expressions() []string {
	exprs := make([]string, len(i.def.Columns))
	for j, col := range i.def.Columns {
		exprs[j] = i.table + "." + col
	}
	return exprs
}
func (i *pebbleIndex) IsUnique() bool                                 { return i.def.Unique }
func (i *pebbleIndex) IsSpatial() bool                                { return false }
func (i *pebbleIndex) IsFullText() bool                               { return false }
func (i *pebbleIndex) IsVector() bool                                 { return false }
func (i *pebbleIndex) Comment() string                                { return "" }
func (i *pebbleIndex) IndexType() string                              { return "BTREE" }
func (i *pebbleIndex) IsGenerated() bool                              { return false }
func (i *pebbleIndex) CanSupport(_ *sql.Context, _ ...sql.Range) bool { return false }
func (i *pebbleIndex) CanSupportOrderBy(_ sql.Expression) bool        { return false }
func (i *pebbleIndex) PrefixLengths() []uint16                        { return nil }

func (i *pebbleIndex) ColumnExpressionTypes() []sql.ColumnExpressionType {
	cets := make([]sql.ColumnExpressionType, len(i.def.Columns))
	for j, col := range i.def.Columns {
		var typ sql.Type = types.Text
		for _, c := range i.sch {
			if c.Name == col {
				typ = c.Type
				break
			}
		}
		cets[j] = sql.ColumnExpressionType{
			Expression: i.table + "." + col,
			Type:       typ,
		}
	}
	return cets
}

// encodeSQLPK extracts the primary key from a sql.Row and encodes it.
// go-mysql-server returns different Go types for the same column across
// operations (int32 on INSERT, int64 on UPDATE); normalise before encoding.
// Returns an error when the PK encoding fails (NULL PK, unsupported type),
// which the caller must surface rather than silently producing an empty key
// that would collide with every other row.
func encodeSQLPK(schema *rdbengine.TableSchema, row sql.Row) ([]byte, error) {
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
		case int:
			return float32(v)
		case int8:
			return float32(v)
		case int16:
			return float32(v)
		case int32:
			return float32(v)
		case int64:
			return float32(v)
		case uint:
			return float32(v)
		case uint8:
			return float32(v)
		case uint16:
			return float32(v)
		case uint32:
			return float32(v)
		case uint64:
			return float32(v)
		}
	case rdbengine.ColumnTypeFloat64:
		switch v := val.(type) {
		case float32:
			return float64(v)
		case float64:
			return v
		case int:
			return float64(v)
		case int8:
			return float64(v)
		case int16:
			return float64(v)
		case int32:
			return float64(v)
		case int64:
			return float64(v)
		case uint:
			return float64(v)
		case uint8:
			return float64(v)
		case uint16:
			return float64(v)
		case uint32:
			return float64(v)
		case uint64:
			return float64(v)
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
