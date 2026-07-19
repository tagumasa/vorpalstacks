package rdsdata

// ExecuteStatementInput represents the ExecuteStatement request.
type ExecuteStatementInput struct {
	ResourceArn string `json:"resourceArn"`
	SecretArn   string `json:"secretArn"`
	Sql         string `json:"sql"`
	Database    string `json:"database,omitempty"`
	// Schema is accepted for AWS SDK compatibility but, per the AWS spec,
	// "isn't currently supported." It is captured for round-trip fidelity
	// and silently ignored at execution time.
	Schema                string `json:"schema,omitempty"`
	TransactionID         string `json:"transactionId,omitempty"`
	IncludeResultMetadata bool   `json:"includeResultMetadata,omitempty"`
	FormatRecordsAs       string `json:"formatRecordsAs,omitempty"`
	// ContinueAfterTimeout controls whether the statement continues
	// running after the Data API's 45-second call timeout (AWS Aurora
	// User Guide, "Controlling Data API timeout behavior").
	//
	// When true:
	//   - The client receives StatementTimeoutException at 45 s.
	//   - The statement continues in a background goroutine on a
	//     detached context for up to maxBgStatementTime (1 hour; see
	//     engine_bridge.go). DDL changes take effect even though the
	//     client already received the timeout error.
	//   - For transactional calls (transactionId set), CommitTransaction
	//     / RollbackTransaction wait for the background statement —
	//     bounded at 45 s — before issuing COMMIT / ROLLBACK. If the
	//     wait itself times out, the transaction stays open and the
	//     client can retry.
	//
	// When false (default): the statement is aborted at the 45-second
	// timeout and StatementTimeoutException is returned.
	//
	// AWS spec: "A value that indicates whether to continue running the
	// statement after the call times out."
	ContinueAfterTimeout bool              `json:"continueAfterTimeout,omitempty"`
	Parameters           []SqlParameter    `json:"parameters,omitempty"`
	ResultSetOptions     *ResultSetOptions `json:"resultSetOptions,omitempty"`
}

// BatchExecuteStatementInput represents the BatchExecuteStatement request.
type BatchExecuteStatementInput struct {
	ResourceArn   string           `json:"resourceArn"`
	SecretArn     string           `json:"secretArn"`
	Sql           string           `json:"sql"`
	Database      string           `json:"database,omitempty"`
	Schema        string           `json:"schema,omitempty"`
	TransactionID string           `json:"transactionId,omitempty"`
	ParameterSets [][]SqlParameter `json:"parameterSets,omitempty"`
}

// ExecuteSqlInput represents the deprecated ExecuteSql request.
// The AWS SDK uses different field names than ExecuteStatement:
// awsSecretStoreArn (not secretArn), dbClusterOrInstanceArn (not resourceArn),
// sqlStatements (not sql).
type ExecuteSqlInput struct {
	ResourceArn string `json:"dbClusterOrInstanceArn,omitempty"`
	SecretArn   string `json:"awsSecretStoreArn,omitempty"`
	Sql         string `json:"sqlStatements,omitempty"`
	Database    string `json:"database,omitempty"`
	Schema      string `json:"schema,omitempty"`
}

// BeginTransactionInput represents the BeginTransaction request.
type BeginTransactionInput struct {
	ResourceArn string `json:"resourceArn"`
	SecretArn   string `json:"secretArn"`
	Database    string `json:"database,omitempty"`
	// Schema is accepted for AWS SDK compatibility but, per the AWS spec,
	// "isn't currently supported." Captured for round-trip fidelity.
	Schema string `json:"schema,omitempty"`
}

// CommitTransactionInput represents the CommitTransaction request.
type CommitTransactionInput struct {
	ResourceArn   string `json:"resourceArn"`
	SecretArn     string `json:"secretArn"`
	TransactionID string `json:"transactionId"`
}

// RollbackTransactionInput represents the RollbackTransaction request.
type RollbackTransactionInput struct {
	ResourceArn   string `json:"resourceArn"`
	SecretArn     string `json:"secretArn"`
	TransactionID string `json:"transactionId"`
}

// SqlParameter represents a named SQL parameter.
type SqlParameter struct {
	Name     string `json:"name,omitempty"`
	Value    *Field `json:"value,omitempty"`
	TypeHint string `json:"typeHint,omitempty"`
}

// Field represents a typed value in RDS Data API responses.
type Field struct {
	IsNull       *bool       `json:"isNull,omitempty"`
	StringValue  *string     `json:"stringValue,omitempty"`
	LongValue    *int64      `json:"longValue,omitempty"`
	DoubleValue  *float64    `json:"doubleValue,omitempty"`
	BooleanValue *bool       `json:"booleanValue,omitempty"`
	BlobValue    []byte      `json:"blobValue,omitempty"`
	ArrayValue   interface{} `json:"arrayValue,omitempty"`
}

// ColumnMetadata describes a single column in the result set.
type ColumnMetadata struct {
	ArrayBaseColumnType int32  `json:"arrayBaseColumnType,omitempty"`
	IsAutoIncrement     bool   `json:"isAutoIncrement,omitempty"`
	IsCaseSensitive     bool   `json:"isCaseSensitive,omitempty"`
	IsCurrency          bool   `json:"isCurrency,omitempty"`
	IsSigned            bool   `json:"isSigned,omitempty"`
	Label               string `json:"label,omitempty"`
	Name                string `json:"name,omitempty"`
	Nullable            int32  `json:"nullable,omitempty"`
	Precision           int32  `json:"precision,omitempty"`
	Scale               int32  `json:"scale,omitempty"`
	SchemaName          string `json:"schemaName,omitempty"`
	TableName           string `json:"tableName,omitempty"`
	Type                int32  `json:"type,omitempty"`
	TypeName            string `json:"typeName,omitempty"`
}

// UpdateResult represents the result of an INSERT/UPDATE/DELETE.
type UpdateResult struct {
	GeneratedFields []Field `json:"generatedFields,omitempty"`
}

// ExecuteStatementResponse represents the ExecuteStatement response.
type ExecuteStatementResponse struct {
	Records                [][]Field        `json:"records,omitempty"`
	ColumnMetadata         []ColumnMetadata `json:"columnMetadata,omitempty"`
	NumberOfRecordsUpdated int64            `json:"numberOfRecordsUpdated,omitempty"`
	GeneratedFields        []Field          `json:"generatedFields,omitempty"`
	FormattedRecords       string           `json:"formattedRecords,omitempty"`
}

// BatchExecuteStatementResponse represents the BatchExecuteStatement response.
type BatchExecuteStatementResponse struct {
	UpdateResults []UpdateResult `json:"updateResults,omitempty"`
}

// ExecuteSqlResponse represents the deprecated ExecuteSql response.
type ExecuteSqlResponse struct {
	SqlStatementResults []SqlStatementResult `json:"sqlStatementResults,omitempty"`
}

// SqlStatementResult represents a single SQL result for the deprecated ExecuteSql.
type SqlStatementResult struct {
	NumberOfRecordsUpdated int64        `json:"numberOfRecordsUpdated,omitempty"`
	ResultFrame            *ResultFrame `json:"resultFrame,omitempty"`
}

// ResultFrame represents a result set frame.
type ResultFrame struct {
	ResultSetMetadata *ResultSetMetadata `json:"resultSetMetadata,omitempty"`
	Records           []Record           `json:"records,omitempty"`
}

// ResultSetMetadata contains column count and column metadata.
type ResultSetMetadata struct {
	ColumnCount    int64            `json:"columnCount,omitempty"`
	ColumnMetadata []ColumnMetadata `json:"columnMetadata,omitempty"`
}

// Record represents a single row for the deprecated ExecuteSql.
type Record struct {
	Values []Value `json:"values,omitempty"`
}

// Value represents a column value for the deprecated ExecuteSql.
type Value struct {
	IsNull      *bool    `json:"isNull,omitempty"`
	StringValue *string  `json:"stringValue,omitempty"`
	LongValue   *int64   `json:"longValue,omitempty"`
	DoubleValue *float64 `json:"doubleValue,omitempty"`
	IntValue    *int32   `json:"intValue,omitempty"`
	BigIntValue *int64   `json:"bigIntValue,omitempty"`
	BitValue    *bool    `json:"bitValue,omitempty"`
	RealValue   *float32 `json:"realValue,omitempty"`
	BlobValue   []byte   `json:"blobValue,omitempty"`
}

// BeginTransactionResponse represents the BeginTransaction response.
type BeginTransactionResponse struct {
	TransactionID string `json:"transactionId"`
}

// CommitTransactionResponse represents the CommitTransaction response.
type CommitTransactionResponse struct {
	TransactionStatus string `json:"transactionStatus"`
}

// RollbackTransactionResponse represents the RollbackTransaction response.
type RollbackTransactionResponse struct {
	TransactionStatus string `json:"transactionStatus"`
}

// ResultSetOptions controls how result values are formatted.
type ResultSetOptions struct {
	DecimalReturnType string `json:"decimalReturnType,omitempty"`
	LongReturnType    string `json:"longReturnType,omitempty"`
}

// SecretCredential holds parsed credentials from a Secrets Manager secret.
type SecretCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Engine   string `json:"engine,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
}
