package athena

// WorkGroupStoreInterface defines operations for managing Athena work groups.
type WorkGroupStoreInterface interface {
	CreateWorkGroup(wg *WorkGroup) error
	GetWorkGroup(name string) (*WorkGroup, error)
	UpdateWorkGroup(wg *WorkGroup) error
	DeleteWorkGroup(name string) error
	ListWorkGroups() ([]*WorkGroup, error)
	GetARN(name string) string
}

// NamedQueryStoreInterface defines operations for managing named queries.
type NamedQueryStoreInterface interface {
	CreateNamedQuery(nq *NamedQuery) error
	GetNamedQuery(id string) (*NamedQuery, error)
	DeleteNamedQuery(id string) error
	ListNamedQueries(workGroup string) ([]*NamedQuery, error)
}

// PreparedStatementStoreInterface defines operations for managing prepared statements.
type PreparedStatementStoreInterface interface {
	CreatePreparedStatement(ps *PreparedStatement) error
	GetPreparedStatement(workGroup, name string) (*PreparedStatement, error)
	UpdatePreparedStatement(ps *PreparedStatement) error
	DeletePreparedStatement(workGroup, name string) error
	ListPreparedStatements(workGroup string) ([]*PreparedStatement, error)
}

// QueryExecutionStoreInterface defines operations for managing query executions.
type QueryExecutionStoreInterface interface {
	CreateQueryExecution(qe *QueryExecution) error
	GetQueryExecution(id string) (*QueryExecution, error)
	UpdateQueryExecution(qe *QueryExecution) error
	DeleteQueryExecution(id string) error
	ListQueryExecutionIDs(workGroup string, maxResults int) ([]string, error)
}

// ResultStoreInterface defines operations for storing query results.
type ResultStoreInterface interface {
	StoreResult(queryExecutionId string, result *QueryResult) error
	GetResult(queryExecutionId string) (*QueryResult, error)
	DeleteResult(queryExecutionId string) error
}

// DataCatalogStoreInterface defines operations for managing data catalogues.
type DataCatalogStoreInterface interface {
	CreateDataCatalog(dc *DataCatalog) error
	GetDataCatalog(name string) (*DataCatalog, error)
	UpdateDataCatalog(dc *DataCatalog) error
	DeleteDataCatalog(name string) error
	ListDataCatalogs() ([]*DataCatalog, error)
}

// DatabaseStoreInterface defines operations for managing databases.
type DatabaseStoreInterface interface {
	CreateDatabase(catalog string, db *Database) error
	GetDatabase(catalog, name string) (*Database, error)
	ListDatabases(catalog string) ([]*Database, error)
	DeleteDatabase(catalog, name string) error
}

// TableStoreInterface defines operations for managing tables.
type TableStoreInterface interface {
	CreateTable(catalog, database string, table *TableMetadata) error
	GetTable(catalog, database, name string) (*TableMetadata, error)
	ListTables(catalog, database string) ([]*TableMetadata, error)
	DeleteTable(catalog, database, name string) error
	UpdateTable(catalog, database string, table *TableMetadata) error
}

// TableDataStoreInterface defines operations for storing table data.
type TableDataStoreInterface interface {
	StoreTableData(catalog, database, table string, storedTable *StoredTable) error
	GetTableData(catalog, database, table string) (*StoredTable, error)
	DeleteTableData(catalog, database, table string) error
}

// AthenaStoresInterface defines access to all Athena stores.
type AthenaStoresInterface interface {
	WorkGroupStore() WorkGroupStoreInterface
	NamedQueryStore() NamedQueryStoreInterface
	PreparedStatementStore() PreparedStatementStoreInterface
	QueryExecutionStore() QueryExecutionStoreInterface
	ResultStore() ResultStoreInterface
	DataCatalogStore() DataCatalogStoreInterface
	DatabaseStore() DatabaseStoreInterface
	TableStore() TableStoreInterface
	TableDataStore() TableDataStoreInterface
	WorkGroupStoreRaw() *WorkGroupStore
	NamedQueryStoreRaw() *NamedQueryStore
	PreparedStatementStoreRaw() *PreparedStatementStore
	QueryExecutionStoreRaw() *QueryExecutionStore
	ResultStoreRaw() *ResultStore
	DataCatalogStoreRaw() *DataCatalogStore
	DatabaseStoreRaw() *DatabaseStore
	TableStoreRaw() *TableStore
	TableDataStoreRaw() *TableDataStore
}

// AthenaStore provides access to all Athena stores.
type AthenaStore struct {
	workGroupStore         *WorkGroupStore
	namedQueryStore        *NamedQueryStore
	preparedStatementStore *PreparedStatementStore
	queryExecutionStore    *QueryExecutionStore
	resultStore            *ResultStore
	dataCatalogStore       *DataCatalogStore
	databaseStore          *DatabaseStore
	tableStore             *TableStore
	tableDataStore         *TableDataStore
}

// NewAthenaStore creates a new AthenaStore with the given stores.
func NewAthenaStore(
	workGroupStore *WorkGroupStore,
	namedQueryStore *NamedQueryStore,
	preparedStatementStore *PreparedStatementStore,
	queryExecutionStore *QueryExecutionStore,
	resultStore *ResultStore,
	dataCatalogStore *DataCatalogStore,
	databaseStore *DatabaseStore,
	tableStore *TableStore,
	tableDataStore *TableDataStore,
) *AthenaStore {
	return &AthenaStore{
		workGroupStore:         workGroupStore,
		namedQueryStore:        namedQueryStore,
		preparedStatementStore: preparedStatementStore,
		queryExecutionStore:    queryExecutionStore,
		resultStore:            resultStore,
		dataCatalogStore:       dataCatalogStore,
		databaseStore:          databaseStore,
		tableStore:             tableStore,
		tableDataStore:         tableDataStore,
	}
}

// WorkGroupStore returns the work group store.
func (s *AthenaStore) WorkGroupStore() WorkGroupStoreInterface { return s.workGroupStore }

// NamedQueryStore returns the named query store.
func (s *AthenaStore) NamedQueryStore() NamedQueryStoreInterface { return s.namedQueryStore }

// PreparedStatementStore returns the prepared statement store.
func (s *AthenaStore) PreparedStatementStore() PreparedStatementStoreInterface {
	return s.preparedStatementStore
}

// QueryExecutionStore returns the query execution store.
func (s *AthenaStore) QueryExecutionStore() QueryExecutionStoreInterface {
	return s.queryExecutionStore
}

// ResultStore returns the result store.
func (s *AthenaStore) ResultStore() ResultStoreInterface { return s.resultStore }

// DataCatalogStore returns the data catalog store.
func (s *AthenaStore) DataCatalogStore() DataCatalogStoreInterface { return s.dataCatalogStore }

// DatabaseStore returns the database store.
func (s *AthenaStore) DatabaseStore() DatabaseStoreInterface { return s.databaseStore }

// TableStore returns the table store.
func (s *AthenaStore) TableStore() TableStoreInterface { return s.tableStore }

// TableDataStore returns the table data store.
func (s *AthenaStore) TableDataStore() TableDataStoreInterface { return s.tableDataStore }

// WorkGroupStoreRaw returns the raw work group store.
func (s *AthenaStore) WorkGroupStoreRaw() *WorkGroupStore { return s.workGroupStore }

// NamedQueryStoreRaw returns the raw named query store.
func (s *AthenaStore) NamedQueryStoreRaw() *NamedQueryStore { return s.namedQueryStore }

// PreparedStatementStoreRaw returns the raw prepared statement store.
func (s *AthenaStore) PreparedStatementStoreRaw() *PreparedStatementStore {
	return s.preparedStatementStore
}

// QueryExecutionStoreRaw returns the raw query execution store.
func (s *AthenaStore) QueryExecutionStoreRaw() *QueryExecutionStore { return s.queryExecutionStore }

// ResultStoreRaw returns the raw result store.
func (s *AthenaStore) ResultStoreRaw() *ResultStore { return s.resultStore }

// DataCatalogStoreRaw returns the raw data catalog store.
func (s *AthenaStore) DataCatalogStoreRaw() *DataCatalogStore { return s.dataCatalogStore }

// DatabaseStoreRaw returns the raw database store.
func (s *AthenaStore) DatabaseStoreRaw() *DatabaseStore { return s.databaseStore }

// TableStoreRaw returns the raw table store.
func (s *AthenaStore) TableStoreRaw() *TableStore { return s.tableStore }

// TableDataStoreRaw returns the raw table data store.
func (s *AthenaStore) TableDataStoreRaw() *TableDataStore { return s.tableDataStore }

var (
	_ WorkGroupStoreInterface         = (*WorkGroupStore)(nil)
	_ NamedQueryStoreInterface        = (*NamedQueryStore)(nil)
	_ PreparedStatementStoreInterface = (*PreparedStatementStore)(nil)
	_ QueryExecutionStoreInterface    = (*QueryExecutionStore)(nil)
	_ ResultStoreInterface            = (*ResultStore)(nil)
	_ DataCatalogStoreInterface       = (*DataCatalogStore)(nil)
	_ DatabaseStoreInterface          = (*DatabaseStore)(nil)
	_ TableStoreInterface             = (*TableStore)(nil)
	_ TableDataStoreInterface         = (*TableDataStore)(nil)
	_ AthenaStoresInterface           = (*AthenaStore)(nil)
)
