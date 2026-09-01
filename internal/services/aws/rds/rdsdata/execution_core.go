// execution_core.go holds the Core execution path for the RDS Data API
// operations: input validation, engine resolution, credential validation,
// transaction-session management, and SQL execution. The HTTP handlers in
// service.go are thin adapters that parse the wire request and delegate
// here; CommitTransaction and RollbackTransaction share the transaction
// Cores in transaction_core.go.
package rdsdata

import (
	"context"
	"fmt"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
)

// executeStatementInTxCore runs a single SQL statement inside an existing
// transaction. Unlike executeStatementCore the transaction id is required:
// an empty id must never silently degrade into an autocommit execution.
// The explicit emptiness check exists because the length validator treats
// an empty transactionId as valid (its AWS-spec minimum is zero), and the
// old invoker path rejected the empty id through the failed transaction
// lookup — a TransactionNotFoundException, restored here.
func (s *RDSDataService) executeStatementInTxCore(ctx context.Context, input *ExecuteStatementInput) (interface{}, error) {
	if input.TransactionID == "" {
		return nil, transactionNotFound("transaction id is required")
	}
	return s.executeStatementCore(ctx, input)
}

// executeStatementCore runs a single SQL statement and returns results.
// transactionId is optional on this path (an omitted id runs the statement
// in autocommit mode), matching the ExecuteStatement/BatchExecuteStatement
// wire contract.
func (s *RDSDataService) executeStatementCore(ctx context.Context, input *ExecuteStatementInput) (interface{}, error) {
	if err := validateCommon(input.ResourceArn, input.SecretArn, input.Database, input.Schema); err != nil {
		return nil, err
	}
	if err := validateSQL(input.Sql); err != nil {
		return nil, err
	}
	if err := validateTransactionID(input.TransactionID, true); err != nil {
		return nil, err
	}

	engine, instanceID, err := s.resolveEngine(input.ResourceArn)
	if err != nil {
		return nil, err
	}
	if err := s.validateCredentials(ctx, input.SecretArn); err != nil {
		return nil, err
	}

	// Validate and substitute parameters BEFORE acquiring s.mu. This
	// has two benefits:
	//  1. If validation fails, we return early without touching the
	//     transaction map or bg counters — no leak of bgCount/bgWg.
	//  2. The s.mu critical section stays short (no regex work under
	//     the lock).
	sqlStr := input.Sql
	if len(input.Parameters) > 0 {
		if err := validateParameters(input.Parameters); err != nil {
			return nil, err
		}
		sqlStr = substituteParameters(sqlStr, input.Parameters)
	}

	database := input.Database
	var txCtx *sql.Context
	var txEntry *staleEntry
	if input.TransactionID != "" {
		now := time.Now()
		s.mu.Lock()
		entry, ok := s.transactions[input.TransactionID]
		if ok {
			if entry.isExpired(now) {
				delete(s.transactions, input.TransactionID)
				entry = nil
				ok = false
			} else {
				entry.touch()
				// If ContinueAfterTimeout will be used, pre-acquire the
				// bg counters under s.mu so purgeExpired cannot delete
				// the entry between Unlock and the goroutine start in
				// executeWithContinueAfterTimeout. The matching decrement
				// is in the goroutine's defer.
				if input.ContinueAfterTimeout {
					entry.bgCount.Add(1)
					entry.bgWg.Add(1)
				}
			}
		}
		s.mu.Unlock()
		if !ok {
			return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
		}
		if entry.database != "" {
			database = entry.database
		}
		txCtx = entry.sqlCtx
		txEntry = entry
		// Serialise concurrent statements on the same transaction.
		// go-mysql-server's *sql.Context is not thread-safe.
		txEntry.execMu.Lock()
	}

	sqlCtx := txCtx
	if sqlCtx == nil {
		sqlCtx = newSQLContext(database)
	}

	if input.ContinueAfterTimeout {
		// execMu will be unlocked by the background goroutine's defer
		// so the lock covers the entire background execution window.
		return s.executeWithContinueAfterTimeout(engine, sqlCtx, txEntry, instanceID, sqlStr,
			input.IncludeResultMetadata, input.FormatRecordsAs, input.ResultSetOptions, input.Sql)
	}

	// Non-ContinueAfterTimeout: release execMu after execution.
	defer func() {
		if txEntry != nil {
			txEntry.execMu.Unlock()
		}
	}()

	result, err := executeSQLOpts(engine, sqlCtx, sqlStr, input.IncludeResultMetadata, input.FormatRecordsAs, input.ResultSetOptions, instanceID)
	if err != nil {
		return nil, mapSQLError(err)
	}

	logs.Info("rdsdata: ExecuteStatement", logs.String("instance", instanceID), logs.String("sql", truncateSQL(input.Sql)))
	return result, nil
}

// executeWithContinueAfterTimeout implements the AWS Data API
// ContinueAfterTimeout contract: the statement runs on a detached context
// that outlives the 45-second advisory deadline. The client receives
// StatementTimeoutException at 45 s; the goroutine continues until
// completion (or maxBgStatementTime, whichever is first).
//
// Transactional calls: bgWg.Add(1) and bgCount.Add(1) were already
// performed under s.mu in executeStatementCore (before this function is
// called) to close the race window with purgeExpired. The goroutine's
// defer handles the matching Done / Add(-1).
//
// Non-transactional calls: the service-level bgWg is incremented here
// (no purgeExpired interaction) and Shutdown waits for completion.
func (s *RDSDataService) executeWithContinueAfterTimeout(
	engine *sqle.Engine,
	sqlCtx *sql.Context,
	txEntry *staleEntry,
	instanceID, sqlStr string,
	includeMetadata bool,
	formatRecordsAs string,
	resultSetOpts *ResultSetOptions,
	originalSQL string,
) (interface{}, error) {
	bgCtx, bgCancel := context.WithTimeout(context.Background(), maxBgStatementTime)
	bgSqlCtx := sqlCtx.WithContext(bgCtx)

	type execResult struct {
		resp *ExecuteStatementResponse
		err  error
	}
	resultCh := make(chan execResult, 1)

	// Pick the correct WaitGroup. For transactional calls, bgWg was
	// already incremented under s.mu in executeStatementCore; we only
	// decrement in the goroutine's defer. For non-transactional calls,
	// increment the service-level bgWg here.
	wg := &s.bgWg
	if txEntry != nil {
		wg = &txEntry.bgWg
	} else {
		wg.Add(1)
	}

	go func() {
		defer wg.Done()
		defer bgCancel()
		if txEntry != nil {
			defer txEntry.execMu.Unlock()
		}
		if txEntry != nil {
			defer txEntry.bgCount.Add(-1)
		}
		defer func() {
			if re := recover(); re != nil {
				logs.Error("rdsdata: panic in background ContinueAfterTimeout statement",
					logs.Any("panic", re),
					logs.String("instance", instanceID),
					logs.String("sql", truncateSQL(originalSQL)))
				resultCh <- execResult{nil, fmt.Errorf("internal panic: %v", re)}
			}
		}()
		resp, err := executeSQLOpts(engine, bgSqlCtx, sqlStr, includeMetadata, formatRecordsAs, resultSetOpts, instanceID)
		if err != nil {
			logs.Warn("rdsdata: background ContinueAfterTimeout statement failed",
				logs.Err(err),
				logs.String("instance", instanceID),
				logs.String("sql", truncateSQL(originalSQL)))
		}
		resultCh <- execResult{resp, err}
	}()

	timer := time.NewTimer(defaultStatementTimeout)
	defer timer.Stop()
	select {
	case r := <-resultCh:
		if r.err != nil {
			return nil, mapSQLError(r.err)
		}
		logs.Info("rdsdata: ExecuteStatement (ContinueAfterTimeout, completed within deadline)",
			logs.String("instance", instanceID), logs.String("sql", truncateSQL(originalSQL)))
		return r.resp, nil
	case <-timer.C:
		// StatementTimeoutException per AWS spec; goroutine continues.
		s.trackBgCancel(bgCancel)
		logs.Info("rdsdata: ExecuteStatement (ContinueAfterTimeout, timed out, running in background)",
			logs.String("instance", instanceID), logs.String("sql", truncateSQL(originalSQL)))
		return nil, statementTimeout(fmt.Sprintf(
			"statement exceeded %v timeout; execution continues in the background",
			defaultStatementTimeout), instanceID)
	}
}

// batchExecuteStatementCore runs a SQL statement with multiple parameter sets.
func (s *RDSDataService) batchExecuteStatementCore(ctx context.Context, input *BatchExecuteStatementInput) (interface{}, error) {
	if err := validateCommon(input.ResourceArn, input.SecretArn, input.Database, input.Schema); err != nil {
		return nil, err
	}
	if err := validateSQL(input.Sql); err != nil {
		return nil, err
	}
	if err := validateTransactionID(input.TransactionID, true); err != nil {
		return nil, err
	}

	engine, instanceID, err := s.resolveEngine(input.ResourceArn)
	if err != nil {
		return nil, err
	}
	if err := s.validateCredentials(ctx, input.SecretArn); err != nil {
		return nil, err
	}

	database := input.Database
	var txCtx *sql.Context
	if input.TransactionID != "" {
		now := time.Now()
		s.mu.Lock()
		entry, ok := s.transactions[input.TransactionID]
		if ok {
			if entry.isExpired(now) {
				delete(s.transactions, input.TransactionID)
				entry = nil
				ok = false
			} else {
				entry.touch()
			}
		}
		s.mu.Unlock()
		if !ok {
			return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
		}
		if entry.database != "" {
			// AWS spec: when a transaction is bound to a database at
			// BeginTransaction, every subsequent ExecuteStatement /
			// BatchExecuteStatement in that transaction targets the
			// same database regardless of the per-call Database value.
			// Setting Database on the per-call input is allowed but
			// silently overridden.
			database = entry.database
		}
		txCtx = entry.sqlCtx
		// Serialise concurrent statements on the same transaction.
		entry.execMu.Lock()
		defer entry.execMu.Unlock()
	}

	// AWS Data API spec: BatchExecuteStatement exists to run the same SQL
	// multiple times with different parameter sets. An empty parameterSets
	// is a usage error rather than a request to execute once; surface it
	// explicitly so callers do not silently receive a single UpdateResult
	// when they expected N.
	if len(input.ParameterSets) == 0 {
		return nil, invalidParam("parameterSets is required and must contain at least one entry")
	}

	var results []UpdateResult
	for _, params := range input.ParameterSets {
		if err := validateParameters(params); err != nil {
			return nil, err
		}
		sqlWithParams := substituteParameters(input.Sql, params)
		sqlCtx := txCtx
		if sqlCtx == nil {
			sqlCtx = newSQLContext(database)
		}
		res, err := executeSQL(engine, sqlCtx, sqlWithParams, false, "", instanceID)
		if err != nil {
			return nil, mapSQLError(err)
		}
		results = append(results, UpdateResult{GeneratedFields: res.GeneratedFields})
	}

	logs.Info("rdsdata: BatchExecuteStatement", logs.String("instance", instanceID), logs.Int("paramSets", len(input.ParameterSets)))
	return &BatchExecuteStatementResponse{UpdateResults: results}, nil
}

// executeSqlCore runs one or more SQL statements (deprecated operation).
func (s *RDSDataService) executeSqlCore(ctx context.Context, input *ExecuteSqlInput) (interface{}, error) {
	if err := validateCommon(input.ResourceArn, input.SecretArn, input.Database, input.Schema); err != nil {
		return nil, err
	}
	if err := validateSQL(input.Sql); err != nil {
		return nil, err
	}

	engine, instanceID, err := s.resolveEngine(input.ResourceArn)
	if err != nil {
		return nil, err
	}
	if err := s.validateCredentials(ctx, input.SecretArn); err != nil {
		return nil, err
	}

	statements := splitSQL(input.Sql)
	// AWS spec for the deprecated ExecuteSql operation: statements run
	// sequentially within a single session so that session state (user
	// variables, temporary tables, prepared statements) carries across.
	// Creating a fresh sql.Context per statement would silently drop
	// that state and produce different results from AWS.
	sqlCtx := newSQLContext(input.Database)
	var sqlResults []SqlStatementResult
	for _, stmt := range statements {
		res, err := executeSQL(engine, sqlCtx, stmt, true, "", instanceID)
		if err != nil {
			return nil, mapSQLError(err)
		}
		sr := SqlStatementResult{
			NumberOfRecordsUpdated: res.NumberOfRecordsUpdated,
		}
		if len(res.Records) > 0 || len(res.ColumnMetadata) > 0 {
			records := make([]Record, len(res.Records))
			for i, row := range res.Records {
				vals := make([]Value, len(row))
				for j, f := range row {
					vals[j] = fieldToValue(f)
				}
				records[i] = Record{Values: vals}
			}
			sr.ResultFrame = &ResultFrame{
				ResultSetMetadata: &ResultSetMetadata{
					ColumnCount:    int64(len(res.ColumnMetadata)),
					ColumnMetadata: res.ColumnMetadata,
				},
				Records: records,
			}
		}
		sqlResults = append(sqlResults, sr)
	}

	logs.Info("rdsdata: ExecuteSql", logs.String("instance", instanceID), logs.Int("statements", len(statements)))
	return &ExecuteSqlResponse{SqlStatementResults: sqlResults}, nil
}

// beginTransactionCore starts a new transaction and returns a transaction ID.
func (s *RDSDataService) beginTransactionCore(ctx context.Context, input *BeginTransactionInput) (interface{}, error) {
	if err := validateCommon(input.ResourceArn, input.SecretArn, input.Database, input.Schema); err != nil {
		return nil, err
	}

	engine, instanceID, err := s.resolveEngine(input.ResourceArn)
	if err != nil {
		return nil, err
	}
	if err := s.validateCredentials(ctx, input.SecretArn); err != nil {
		return nil, err
	}

	var sqlCtx *sql.Context
	if s.vmysqlProvider != nil {
		sqlCtx = s.vmysqlProvider.NewContext(instanceID, input.Database)
	}
	if sqlCtx == nil {
		sqlCtx = newSQLContext(input.Database)
	}

	if _, err := executeSQL(engine, sqlCtx, "START TRANSACTION", false, "", instanceID); err != nil {
		return nil, mapSQLError(err)
	}

	txID := uuid.New().String()
	now := time.Now()
	s.mu.Lock()
	s.transactions[txID] = &staleEntry{
		engine:   engine,
		created:  now,
		lastSeen: now,
		database: input.Database,
		schema:   input.Schema,
		sqlCtx:   sqlCtx,
	}
	s.mu.Unlock()

	logs.Info("rdsdata: BeginTransaction", logs.String("txId", txID))
	return &BeginTransactionResponse{TransactionID: txID}, nil
}

// resolveEngine resolves a resource ARN to the sqle.Engine for that instance.
//
// AWS Data API error mapping (see API_Reference > CommitTransaction > Errors):
//
//   - malformed ARN, unknown resource type  -> InvalidParameterException (400)
//   - resourceArn not found                  -> NotFoundException (404)
//   - cluster has no DB instance             -> DatabaseNotFoundException (404)
//   - writer instance unavailable            -> DatabaseUnavailableException (504)
//   - store / engine not configured          -> InternalServerErrorException (500)
func (s *RDSDataService) resolveEngine(resourceArn string) (*sqle.Engine, string, error) {
	identifier, resourceType := parseArn(resourceArn)
	if identifier == "" {
		return nil, "", invalidParam(fmt.Sprintf("invalid resourceArn: %s", resourceArn))
	}

	var instanceID string
	if resourceType == "cluster" {
		if s.rdsStore == nil {
			return nil, "", internalError("RDS store not configured")
		}
		clusterName, instances, err := s.rdsStore.GetCluster(identifier)
		if err != nil {
			// Store distinguishes not-found from internal errors; treat any
			// other failure as service-side. AWS exposes NotFoundException
			// for missing resourceArn specifically.
			return nil, "", notFoundError(fmt.Sprintf("cluster %s not found: %v", identifier, err))
		}
		if clusterName == "" && len(instances) == 0 {
			return nil, "", notFoundError(fmt.Sprintf("cluster %s not found", identifier))
		}
		if len(instances) == 0 {
			// 'The DB cluster doesn't have a DB instance.' (AWS spec)
			return nil, "", databaseNotFound(fmt.Sprintf("cluster %s has no DB instance", identifier))
		}
		instanceID = instances[0]
	} else if resourceType == "db" {
		instanceID = identifier
	} else {
		return nil, "", invalidParam(fmt.Sprintf("unsupported resource type in ARN: %s", resourceArn))
	}

	if s.vmysqlProvider == nil {
		return nil, "", internalError("vmysql provider not configured")
	}
	engine := s.vmysqlProvider.GetEngine(instanceID)
	if engine == nil {
		// AWS surfaces this condition as DatabaseUnavailableException when
		// the writer exists in metadata but is not currently serving. For
		// our purposes 'no engine' is treated as unavailable rather than
		// missing, since the ARN may resolve to a cluster member that has
		// not had its engine started yet.
		return nil, "", databaseUnavailable(fmt.Sprintf("no MySQL engine for instance %s", instanceID))
	}
	return engine, instanceID, nil
}

// validateCredentials ensures secretArn refers to a real, decryptable
// Secrets Manager secret that carries database credentials. The resolver
// distinguishes between 'not found / decrypt failed' (SecretsErrorException)
// and 'found but malformed' (InvalidSecretException), matching the AWS
// Data API contract.
//
// When no SecretResolver is wired (unit-test mode), the function falls back
// to structural checks only and logs a warning so the gap is visible in
// production deployments that forget to call SetSecretResolver.
func (s *RDSDataService) validateCredentials(ctx context.Context, secretArn string) error {
	if s.secretResolver == nil {
		s.secretResolverOnce.Do(func() {
			logs.Warn("rdsdata: secretResolver not configured; accepting secretArn without Secrets Manager lookup")
		})
		return nil
	}
	cred, err := s.secretResolver.ResolveSecret(ctx, secretArn)
	if err != nil {
		return err
	}
	if cred == nil {
		return invalidSecret(fmt.Sprintf("secret %s is empty", secretArn))
	}
	if cred.Username == "" || cred.Password == "" {
		return invalidSecret(fmt.Sprintf("secret %s is missing username or password", secretArn))
	}
	return nil
}
