// Package athena provides Athena data store implementations for vorpalstacks.
package athena

import (
	"time"
)

// WorkGroupState represents the state of an Athena work group.
type WorkGroupState string

const (
	// WorkGroupStateEnabled indicates that the work group is enabled.
	WorkGroupStateEnabled WorkGroupState = "ENABLED"
	// WorkGroupStateDisabled indicates that the work group is disabled.
	WorkGroupStateDisabled WorkGroupState = "DISABLED"
)

// QueryExecutionState represents the state of a query execution.
type QueryExecutionState string

const (
	// QueryExecutionStateQueued indicates that the query is queued.
	QueryExecutionStateQueued QueryExecutionState = "QUEUED"
	// QueryExecutionStateRunning indicates that the query is running.
	QueryExecutionStateRunning QueryExecutionState = "RUNNING"
	// QueryExecutionStateSucceeded indicates that the query succeeded.
	QueryExecutionStateSucceeded QueryExecutionState = "SUCCEEDED"
	// QueryExecutionStateFailed indicates that the query failed.
	QueryExecutionStateFailed QueryExecutionState = "FAILED"
	// QueryExecutionStateCancelled indicates that the query was cancelled.
	QueryExecutionStateCancelled QueryExecutionState = "CANCELLED"
)

// StatementType represents the type of SQL statement.
type StatementType string

const (
	// StatementTypeDDL represents a DDL (Data Definition Language) statement.
	StatementTypeDDL StatementType = "DDL"
	// StatementTypeDML represents a DML (Data Manipulation Language) statement.
	StatementTypeDML StatementType = "DML"
	// StatementTypeUtility represents a utility statement.
	StatementTypeUtility StatementType = "UTILITY"
	// StatementTypeUnknown represents an unknown statement type.
	StatementTypeUnknown StatementType = "UNKNOWN"
)

// WorkGroupConfiguration holds the configuration for an Athena work group.
type WorkGroupConfiguration struct {
	ResultConfiguration                    *ResultConfiguration                    `json:"resultConfiguration,omitempty"`
	EnforceWorkGroupConfiguration          bool                                    `json:"enforceWorkGroupConfiguration,omitempty"`
	PublishCloudWatchMetricsEnabled        bool                                    `json:"publishCloudWatchMetricsEnabled,omitempty"`
	BytesScannedCutoffPerQuery             int64                                   `json:"bytesScannedCutoffPerQuery,omitempty"`
	RequesterPaysEnabled                   bool                                    `json:"requesterPaysEnabled,omitempty"`
	EngineVersion                          *EngineVersion                          `json:"engineVersion,omitempty"`
	AdditionalConfiguration                string                                  `json:"additionalConfiguration,omitempty"`
	ExecutionRole                          string                                  `json:"executionRole,omitempty"`
	CustomerContentEncryptionConfiguration *CustomerContentEncryptionConfiguration `json:"customerContentEncryptionConfiguration,omitempty"`
	EnableMinimumEncryptionConfiguration   bool                                    `json:"enableMinimumEncryptionConfiguration,omitempty"`
}

// ResultConfiguration holds the configuration for query results.
type ResultConfiguration struct {
	OutputLocation          string                   `json:"outputLocation,omitempty"`
	EncryptionConfiguration *EncryptionConfiguration `json:"encryptionConfiguration,omitempty"`
	ExpectedBucketOwner     string                   `json:"expectedBucketOwner,omitempty"`
	ACLConfiguration        *ACLConfiguration        `json:"aclConfiguration,omitempty"`
}

// EncryptionConfiguration holds the encryption configuration for query results.
type EncryptionConfiguration struct {
	EncryptionOption string `json:"encryptionOption,omitempty"`
	KmsKey           string `json:"kmsKey,omitempty"`
}

// ACLConfiguration holds the ACL configuration for query results.
type ACLConfiguration struct {
	S3ACLOption string `json:"s3AclOption,omitempty"`
}

// EngineVersion represents the Athena engine version.
type EngineVersion struct {
	SelectedEngineVersion  string `json:"selectedEngineVersion,omitempty"`
	EffectiveEngineVersion string `json:"effectiveEngineVersion,omitempty"`
}

// CustomerContentEncryptionConfiguration holds the customer content encryption configuration.
type CustomerContentEncryptionConfiguration struct {
	KmsKey string `json:"kmsKey,omitempty"`
}

// WorkGroup represents an Athena work group.
type WorkGroup struct {
	Name          string                  `json:"name"`
	Description   string                  `json:"description,omitempty"`
	Configuration *WorkGroupConfiguration `json:"configuration,omitempty"`
	State         WorkGroupState          `json:"state"`
	CreatedTime   time.Time               `json:"createdTime"`
}

// QueryExecution represents an Athena query execution.
type QueryExecution struct {
	QueryExecutionId      string                    `json:"queryExecutionId"`
	Query                 string                    `json:"query"`
	StatementType         StatementType             `json:"statementType,omitempty"`
	WorkGroup             string                    `json:"workGroup,omitempty"`
	QueryExecutionContext *QueryExecutionContext    `json:"queryExecutionContext,omitempty"`
	ResultConfiguration   *ResultConfiguration      `json:"resultConfiguration,omitempty"`
	Status                *QueryExecutionStatus     `json:"status"`
	Statistics            *QueryExecutionStatistics `json:"statistics,omitempty"`
	SubstatementType      string                    `json:"substatementType,omitempty"`
}

// QueryExecutionContext holds the context for a query execution.
type QueryExecutionContext struct {
	Database string `json:"database,omitempty"`
	Catalog  string `json:"catalog,omitempty"`
}

// QueryExecutionStatus represents the status of a query execution.
type QueryExecutionStatus struct {
	State              QueryExecutionState `json:"state"`
	StateChangeReason  string              `json:"stateChangeReason,omitempty"`
	SubmissionDateTime time.Time           `json:"submissionDateTime"`
	CompletionDateTime time.Time           `json:"completionDateTime,omitempty"`
	AthenaError        *AthenaError        `json:"athenaError,omitempty"`
}

// AthenaError represents an error from Athena.
type AthenaError struct {
	ErrorCategory     int    `json:"errorCategory,omitempty"`
	ErrorType         string `json:"errorType,omitempty"`
	Retryable         bool   `json:"retryable,omitempty"`
	ErrorMessage      string `json:"errorMessage,omitempty"`
	SyntaxErrorRow    int    `json:"syntaxErrorRow,omitempty"`
	SyntaxErrorColumn int    `json:"syntaxErrorColumn,omitempty"`
}

// QueryExecutionStatistics holds statistics about a query execution.
type QueryExecutionStatistics struct {
	EngineExecutionTimeInMillis   int64                   `json:"engineExecutionTimeInMillis,omitempty"`
	DataScannedInBytes            int64                   `json:"dataScannedInBytes,omitempty"`
	DataManifestLocation          string                  `json:"dataManifestLocation,omitempty"`
	TotalExecutionTimeInMillis    int64                   `json:"totalExecutionTimeInMillis,omitempty"`
	QueryQueueTimeInMillis        int64                   `json:"queryQueueTimeInMillis,omitempty"`
	QueryPlanningTimeInMillis     int64                   `json:"queryPlanningTimeInMillis,omitempty"`
	ServiceProcessingTimeInMillis int64                   `json:"serviceProcessingTimeInMillis,omitempty"`
	ResultReuseInformation        *ResultReuseInformation `json:"resultReuseInformation,omitempty"`
}

// ResultReuseInformation holds information about result reuse.
type ResultReuseInformation struct {
	ReusedPreviousResult bool `json:"reusedPreviousResult,omitempty"`
}

// QueryResult represents the result of a query execution.
type QueryResult struct {
	QueryExecutionId string     `json:"queryExecutionId"`
	ResultSet        *ResultSet `json:"resultSet"`
	NextToken        string     `json:"nextToken,omitempty"`
	UpdateCount      int64      `json:"updateCount,omitempty"`
}

// ResultSet represents a set of query results.
type ResultSet struct {
	Rows              []Row              `json:"rows"`
	ResultSetMetadata *ResultSetMetadata `json:"resultSetMetadata"`
}

// Row represents a row in a query result.
type Row struct {
	Data []Datum `json:"data"`
}

// Datum represents a single data value in a query result.
type Datum struct {
	VarCharValue string `json:"varCharValue,omitempty"`
}

// ResultSetMetadata holds metadata about a result set.
type ResultSetMetadata struct {
	ColumnInfo []ColumnInfo `json:"columnInfo"`
}

// ColumnInfo represents information about a column in a result set.
type ColumnInfo struct {
	Label         string `json:"label,omitempty"`
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	Precision     int32  `json:"precision,omitempty"`
	Scale         int32  `json:"scale,omitempty"`
	Nullable      string `json:"nullable,omitempty"`
	CaseSensitive bool   `json:"caseSensitive,omitempty"`
	SchemaName    string `json:"schemaName,omitempty"`
	TableName     string `json:"tableName,omitempty"`
	CatalogName   string `json:"catalogName,omitempty"`
}

// NamedQuery represents a saved Athena query.
type NamedQuery struct {
	NamedQueryId string    `json:"namedQueryId"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Database     string    `json:"database"`
	QueryString  string    `json:"queryString"`
	WorkGroup    string    `json:"workGroup"`
	CreatedTime  time.Time `json:"createdTime"`
}

// PreparedStatement represents a prepared Athena statement.
type PreparedStatement struct {
	StatementName    string    `json:"statementName"`
	WorkGroupName    string    `json:"workGroupName"`
	QueryStatement   string    `json:"queryStatement"`
	Description      string    `json:"description,omitempty"`
	LastModifiedTime time.Time `json:"lastModifiedTime"`
	CreatedTime      time.Time `json:"createdTime"`
}

// DataCatalog represents an Athena data catalog.
type DataCatalog struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"type"`
	Parameters  map[string]string `json:"parameters,omitempty"`
}

// Database represents an Athena database.
type Database struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Parameters  map[string]string `json:"parameters,omitempty"`
}

// TableMetadata represents metadata about an Athena table.
type TableMetadata struct {
	Name          string            `json:"name"`
	DatabaseName  string            `json:"databaseName,omitempty"`
	Description   string            `json:"description,omitempty"`
	TableType     string            `json:"tableType,omitempty"`
	Columns       []Column          `json:"columns,omitempty"`
	PartitionKeys []Column          `json:"partitionKeys,omitempty"`
	Parameters    map[string]string `json:"parameters,omitempty"`
}

// Column represents a column in an Athena table.
type Column struct {
	Name    string `json:"name"`
	Type    string `json:"type,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// StoredRow represents a stored row in an Athena table.
type StoredRow struct {
	Values map[string]interface{} `json:"values"`
}

// StoredTable represents a stored Athena table.
type StoredTable struct {
	DatabaseName string       `json:"databaseName"`
	TableName    string       `json:"tableName"`
	Columns      []Column     `json:"columns"`
	Rows         []*StoredRow `json:"rows"`
}
