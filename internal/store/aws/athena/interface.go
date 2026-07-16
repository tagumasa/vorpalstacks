package athena

import (
	"time"

	common "vorpalstacks/internal/store/aws/common"
)

// WorkGroupStoreInterface defines operations for managing Athena work groups.
type WorkGroupStoreInterface interface {
	CreateWorkGroup(wg *WorkGroup) error
	GetWorkGroup(name string) (*WorkGroup, error)
	UpdateWorkGroup(wg *WorkGroup) error
	DeleteWorkGroup(name string) error
	ListWorkGroups(opts common.ListOptions) (*common.ListResult[WorkGroup], error)
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
	DeleteExpiredQueryExecutions(olderThan time.Time) (int, []string, error)
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

// CapacityReservationStoreInterface defines operations for managing capacity reservations.
type CapacityReservationStoreInterface interface {
	CreateCapacityReservation(cr *CapacityReservation) (*CapacityReservation, error)
	GetCapacityReservation(name string) (*CapacityReservation, error)
	UpdateCapacityReservation(cr *CapacityReservation) error
	ListCapacityReservations(workGroup string) ([]*CapacityReservation, error)
	DeleteCapacityReservation(name string) error
}

var (
	_ WorkGroupStoreInterface           = (*WorkGroupStore)(nil)
	_ NamedQueryStoreInterface          = (*NamedQueryStore)(nil)
	_ PreparedStatementStoreInterface   = (*PreparedStatementStore)(nil)
	_ QueryExecutionStoreInterface      = (*QueryExecutionStore)(nil)
	_ ResultStoreInterface              = (*ResultStore)(nil)
	_ DataCatalogStoreInterface         = (*DataCatalogStore)(nil)
	_ DatabaseStoreInterface            = (*DatabaseStore)(nil)
	_ TableStoreInterface               = (*TableStore)(nil)
	_ TableDataStoreInterface           = (*TableDataStore)(nil)
	_ CapacityReservationStoreInterface = (*CapacityReservationStore)(nil)
)
