// Package rdsdata implements the AWS RDS Data API service, providing SQL
// execution against vmysql instances via the go-mysql-server sqle.Engine.
package rdsdata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
)

// AWS Data API transaction timeouts (see
// https://docs.aws.amazon.com/rdsdataservice/latest/APIReference/API_BeginTransaction.html):
//
//   - Idle timeout: a transaction times out if no calls use its
//     transactionId for 3 minutes. Any ExecuteStatement / CommitTransaction /
//     RollbackTransaction carrying the id resets the idle timer.
//   - Maximum lifetime: a transaction is terminated and rolled back
//     automatically after 24 hours regardless of activity.
const (
	transactionIdleTTL = 3 * time.Minute
	transactionMaxTTL  = 24 * time.Hour
)

// staleEntry holds a transaction session and its deadlines.
//
// Fields:
//   - engine    - sqle.Engine to execute subsequent statements against
//   - created   - absolute creation time (used for the 24 h hard cap)
//   - lastSeen  - wall time of the most recent operation referencing this id
//   - database  - database the transaction was started with (wins over the
//     per-call Database parameter, matching AWS semantics)
//   - schema    - captured for completeness; AWS spec marks schema as
//     currently unsupported
//   - sqlCtx    - the live sql.Context carrying the open transaction state
type staleEntry struct {
	engine   *sqle.Engine
	created  time.Time
	lastSeen time.Time
	database string
	schema   string
	sqlCtx   *sql.Context
	// bgWg tracks outstanding ContinueAfterTimeout statements for
	// CommitTransaction / RollbackTransaction, which must wait for
	// in-flight background statements before issuing COMMIT / ROLLBACK.
	bgWg sync.WaitGroup
	// bgCount is the atomic mirror of bgWg's internal counter. It lets
	// purgeExpired check whether background statements are running
	// without blocking (atomic read instead of WaitGroup.Wait). This
	// eliminates the need for purgeExpired to wait on bgWg, which could
	// block the 30-second cleanup ticker for up to maxBgStatementTime.
	bgCount atomic.Int32
}

// isExpired reports whether the entry has crossed either AWS-defined deadline.
func (e *staleEntry) isExpired(now time.Time) bool {
	if now.Sub(e.lastSeen) > transactionIdleTTL {
		return true
	}
	if now.Sub(e.created) > transactionMaxTTL {
		return true
	}
	return false
}

// touch resets the idle timer without affecting the absolute creation time.
func (e *staleEntry) touch() {
	e.lastSeen = time.Now()
}

// RDSDataService implements the RDS Data API HTTP operations.
type RDSDataService struct {
	mu             sync.RWMutex
	bus            eventbus.Bus
	vmysqlProvider VmysqlProvider
	rdsStore       RDSStoreProvider
	secretResolver SecretResolver
	transactions   map[string]*staleEntry
	cancelCleanup  context.CancelFunc

	// bgMu protects bgCancelFuncs. bgWg tracks all outstanding
	// ContinueAfterTimeout background statements (both transactional and
	// non-transactional). Shutdown cancels every context and then waits
	// on bgWg.
	bgMu          sync.Mutex
	bgCancelFuncs []context.CancelFunc
	bgWg          sync.WaitGroup

	// secretResolverOnce ensures the "secretResolver not configured"
	// warning is emitted at most once per service lifetime instead of
	// on every validateCredentials call, preventing log pollution in
	// dev/test configurations that do not wire a resolver.
	secretResolverOnce sync.Once
}

// VmysqlProvider returns engine information for a given instance ID.
type VmysqlProvider interface {
	GetEngine(instanceID string) *sqle.Engine
	NewContext(instanceID string, database string) *sql.Context
}

// RDSStoreProvider resolves resourceArn to instance identifiers.
type RDSStoreProvider interface {
	GetInstance(identifier string) (instanceName string, engine string, status string, err error)
	GetCluster(identifier string) (clusterName string, instances []string, err error)
}

// SecretResolver validates the secretArn supplied to a Data API call by
// retrieving the referenced Secrets Manager secret and returning its
// parsed credential payload. Implementations must distinguish three
// failure modes that map to distinct AWS exceptions:
//
//   - ErrSecretNotFound  → SecretsErrorException   (HTTP 400)
//   - ErrSecretInvalid   → InvalidSecretException  (HTTP 400)
//   - ErrSecretsService  → SecretsErrorException   (HTTP 400, transient)
//
// A nil return indicates the secret exists, is decryptable, and contains
// at minimum a username and password field.
type SecretResolver interface {
	ResolveSecret(ctx context.Context, secretArn string) (*SecretCredential, error)
}

// NewRDSDataService creates a new RDS Data API service.
func NewRDSDataService() *RDSDataService {
	return &RDSDataService{
		transactions: make(map[string]*staleEntry),
	}
}

// SetEventBus injects the event bus for Secrets Manager access.
func (s *RDSDataService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
}

// SetVmysqlProvider injects the vmysql engine provider.
func (s *RDSDataService) SetVmysqlProvider(p VmysqlProvider) {
	s.vmysqlProvider = p
}

// SetRDSStoreProvider injects the RDS store for ARN resolution.
func (s *RDSDataService) SetRDSStoreProvider(p RDSStoreProvider) {
	s.rdsStore = p
}

// SetSecretResolver injects the secret-lookup implementation used by
// validateCredentials. When nil, validateCredentials falls back to
// structural validation (length / non-empty) only. Production deployments
// must inject a real resolver; nil-resolution is a developer convenience
// for unit tests and the boot-time vmysql test-instance.
func (s *RDSDataService) SetSecretResolver(r SecretResolver) {
	s.secretResolver = r
}

// StartCleanup launches the background stale-transaction reaper.
func (s *RDSDataService) StartCleanup() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelCleanup = cancel
	go s.cleanupLoop(ctx)
}

// Shutdown stops the background cleanup goroutine and waits for all
// outstanding ContinueAfterTimeout statements to finish (or be cancelled).
func (s *RDSDataService) Shutdown() {
	if s.cancelCleanup != nil {
		s.cancelCleanup()
	}
	// Cancel every outstanding background statement context so the
	// goroutines exit promptly instead of running for up to
	// maxBgStatementTime.
	s.bgMu.Lock()
	for _, cancel := range s.bgCancelFuncs {
		cancel()
	}
	s.bgCancelFuncs = nil
	s.bgMu.Unlock()
	s.bgWg.Wait()
}

// trackBgCancel records a CancelFunc for an outstanding ContinueAfterTimeout
// statement so that Shutdown can cancel it on graceful exit.
func (s *RDSDataService) trackBgCancel(cancel context.CancelFunc) {
	s.bgMu.Lock()
	s.bgCancelFuncs = append(s.bgCancelFuncs, cancel)
	s.bgMu.Unlock()
}

// RegisterHandlers registers all 6 RDS Data API operation handlers.
func (s *RDSDataService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("rdsdata", "ExecuteStatement", s.ExecuteStatement)
	d.RegisterHandlerForService("rdsdata", "BatchExecuteStatement", s.BatchExecuteStatement)
	d.RegisterHandlerForService("rdsdata", "ExecuteSql", s.ExecuteSql)
	d.RegisterHandlerForService("rdsdata", "BeginTransaction", s.BeginTransaction)
	d.RegisterHandlerForService("rdsdata", "CommitTransaction", s.CommitTransaction)
	d.RegisterHandlerForService("rdsdata", "RollbackTransaction", s.RollbackTransaction)
}

// ExecuteStatement runs a single SQL statement and returns results.
func (s *RDSDataService) ExecuteStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input ExecuteStatementInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
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
				// executeWithContinueAfterTimeout (L-10 race fix).
				// The matching decrement is in the goroutine's defer.
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
	}

	sqlCtx := txCtx
	if sqlCtx == nil {
		sqlCtx = newSQLContext(database)
	}

	if input.ContinueAfterTimeout {
		return s.executeWithContinueAfterTimeout(engine, sqlCtx, txEntry, instanceID, sqlStr,
			input.IncludeResultMetadata, input.FormatRecordsAs, input.ResultSetOptions, input.Sql)
	}

	result, err := executeSQLOpts(engine, sqlCtx, sqlStr, input.IncludeResultMetadata, input.FormatRecordsAs, input.ResultSetOptions)
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
// performed under s.mu in ExecuteStatement (before this function is
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
	// already incremented under s.mu in ExecuteStatement; we only
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
			defer txEntry.bgCount.Add(-1)
		}
		resp, err := executeSQLOpts(engine, bgSqlCtx, sqlStr, includeMetadata, formatRecordsAs, resultSetOpts)
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
			defaultStatementTimeout))
	}
}

// BatchExecuteStatement runs a SQL statement with multiple parameter sets.
func (s *RDSDataService) BatchExecuteStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input BatchExecuteStatementInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
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
		res, err := executeSQL(engine, sqlCtx, sqlWithParams, false, "")
		if err != nil {
			return nil, mapSQLError(err)
		}
		results = append(results, UpdateResult{GeneratedFields: res.GeneratedFields})
	}

	logs.Info("rdsdata: BatchExecuteStatement", logs.String("instance", instanceID), logs.Int("paramSets", len(input.ParameterSets)))
	return &BatchExecuteStatementResponse{UpdateResults: results}, nil
}

// ExecuteSql runs one or more SQL statements (deprecated operation).
func (s *RDSDataService) ExecuteSql(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input ExecuteSqlInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
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
		res, err := executeSQL(engine, sqlCtx, stmt, true, "")
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

// BeginTransaction starts a new transaction and returns a transaction ID.
func (s *RDSDataService) BeginTransaction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input BeginTransactionInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
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

	if _, err := executeSQL(engine, sqlCtx, "START TRANSACTION", false, ""); err != nil {
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

// CommitTransaction commits a transaction.
func (s *RDSDataService) CommitTransaction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input CommitTransactionInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
	if err := validateCommon(input.ResourceArn, input.SecretArn, "", ""); err != nil {
		return nil, err
	}
	if err := validateTransactionID(input.TransactionID, false); err != nil {
		return nil, err
	}

	if err := s.validateCredentials(ctx, input.SecretArn); err != nil {
		return nil, err
	}

	s.mu.Lock()
	entry, ok := s.transactions[input.TransactionID]
	s.mu.Unlock()

	if !ok {
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
	}

	// Pre-check expiry: if the transaction has already exceeded its
	// idle or max-life deadline, purgeExpired should have already rolled
	// it back. But the purge ticker has a 30 s gap, so an expired-but-
	// not-yet-purged transaction may still be in the map. Reject it
	// explicitly rather than allowing a stale COMMIT to succeed.
	if entry.isExpired(time.Now()) {
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
	}

	// Wait for outstanding ContinueAfterTimeout background statements
	// before committing. The wait is bounded by the Data API's 45-second
	// call timeout (AWS spec: "the Data API cancels an operation and
	// returns a timeout error if the operation doesn't finish processing
	// within 45 seconds"). If the background statement has not finished,
	// the transaction stays in the map and the client can retry
	// CommitTransaction. Deleting from the map only after the wait
	// succeeds ensures the transaction is retryable on timeout.
	if !waitForBg(&entry.bgWg, defaultStatementTimeout) {
		return nil, statementTimeout(fmt.Sprintf(
			"CommitTransaction timed out waiting for a background statement; retry %s",
			input.TransactionID))
	}

	// Background statements are done (or none were running). Atomically
	// remove from the map so a concurrent retry cannot double-commit.
	s.mu.Lock()
	if _, stillThere := s.transactions[input.TransactionID]; !stillThere {
		s.mu.Unlock()
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
	}
	delete(s.transactions, input.TransactionID)
	s.mu.Unlock()

	commitCtx := entry.sqlCtx
	if commitCtx == nil {
		commitCtx = newSQLContext(entry.database)
	}
	if _, err := executeSQL(entry.engine, commitCtx, "COMMIT", false, ""); err != nil {
		return nil, mapSQLError(err)
	}

	logs.Info("rdsdata: CommitTransaction", logs.String("txId", input.TransactionID))
	return &CommitTransactionResponse{TransactionStatus: "COMMIT"}, nil
}

// RollbackTransaction rolls back a transaction.
func (s *RDSDataService) RollbackTransaction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input RollbackTransactionInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
	if err := validateCommon(input.ResourceArn, input.SecretArn, "", ""); err != nil {
		return nil, err
	}
	if err := validateTransactionID(input.TransactionID, false); err != nil {
		return nil, err
	}

	if err := s.validateCredentials(ctx, input.SecretArn); err != nil {
		return nil, err
	}

	s.mu.Lock()
	entry, ok := s.transactions[input.TransactionID]
	s.mu.Unlock()

	if !ok {
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
	}

	// Pre-check expiry (see CommitTransaction for rationale).
	if entry.isExpired(time.Now()) {
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
	}

	// Bounded wait for background statements (see CommitTransaction).
	if !waitForBg(&entry.bgWg, defaultStatementTimeout) {
		return nil, statementTimeout(fmt.Sprintf(
			"RollbackTransaction timed out waiting for a background statement; retry %s",
			input.TransactionID))
	}

	// Atomically remove from the map (see CommitTransaction).
	s.mu.Lock()
	if _, stillThere := s.transactions[input.TransactionID]; !stillThere {
		s.mu.Unlock()
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
	}
	delete(s.transactions, input.TransactionID)
	s.mu.Unlock()

	rollbackCtx := entry.sqlCtx
	if rollbackCtx == nil {
		rollbackCtx = newSQLContext(entry.database)
	}
	if _, err := executeSQL(entry.engine, rollbackCtx, "ROLLBACK", false, ""); err != nil {
		return nil, mapSQLError(err)
	}

	logs.Info("rdsdata: RollbackTransaction", logs.String("txId", input.TransactionID))
	return &RollbackTransactionResponse{TransactionStatus: "ROLLBACK"}, nil
}

// AWS Data API parameter length bounds (see
// https://docs.aws.amazon.com/rdsdataservice/latest/APIReference/).
const (
	resourceArnMinLen   = 11
	resourceArnMaxLen   = 100
	secretArnMinLen     = 11
	secretArnMaxLen     = 100
	sqlMaxLen           = 65536
	transactionIDMaxLen = 192
	databaseMaxLen      = 64
	schemaMaxLen        = 64
)

// waitForBg waits up to timeout for the WaitGroup to reach zero. Returns
// true if the wait completed, false if the timeout expired. The helper
// goroutine exits when the WaitGroup reaches zero (bounded by
// maxBgStatementTime for ContinueAfterTimeout statements), so there is
// no permanent goroutine leak.
func waitForBg(wg *sync.WaitGroup, timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// validateLength enforces AWS-spec string length constraints and returns a
// typed InvalidParameterException on violation. The min check is skipped
// when allowEmpty=true (AWS marks database / schema / transactionId as
// optional with min length 0).
func validateLength(field, value string, min, max int, allowEmpty bool) error {
	n := len(value)
	if allowEmpty && n == 0 {
		return nil
	}
	if n < min {
		return invalidParam(fmt.Sprintf("%s length %d is shorter than minimum %d", field, n, min))
	}
	if n > max {
		return invalidParam(fmt.Sprintf("%s length %d exceeds maximum %d", field, n, max))
	}
	return nil
}

// validateResourceArn enforces the AWS-spec resourceArn constraints.
func validateResourceArn(arn string) error {
	return validateLength("resourceArn", arn, resourceArnMinLen, resourceArnMaxLen, false)
}

// validateSecretArn enforces the AWS-spec secretArn constraints.
func validateSecretArn(arn string) error {
	return validateLength("secretArn", arn, secretArnMinLen, secretArnMaxLen, false)
}

// validateSQL enforces the AWS-spec sql constraints (sql is required).
func validateSQL(sqlStr string) error {
	return validateLength("sql", sqlStr, 1, sqlMaxLen, false)
}

// validateTransactionID enforces the AWS-spec transactionId constraints.
// transactionId is optional for ExecuteStatement / BatchExecuteStatement but
// required for CommitTransaction / RollbackTransaction — callers that need
// 'required' semantics pass allowEmpty=false.
func validateTransactionID(id string, allowEmpty bool) error {
	return validateLength("transactionId", id, 0, transactionIDMaxLen, allowEmpty)
}

// validateDatabase enforces the AWS-spec database constraints.
func validateDatabase(db string) error {
	return validateLength("database", db, 0, databaseMaxLen, true)
}

// validateSchema enforces the AWS-spec schema constraints.
func validateSchema(schema string) error {
	return validateLength("schema", schema, 0, schemaMaxLen, true)
}

// validateCommon enforces the resourceArn + secretArn + database + schema
// constraints shared by every Data API operation.
func validateCommon(resourceArn, secretArn, database, schema string) error {
	if err := validateResourceArn(resourceArn); err != nil {
		return err
	}
	if err := validateSecretArn(secretArn); err != nil {
		return err
	}
	if err := validateDatabase(database); err != nil {
		return err
	}
	if err := validateSchema(schema); err != nil {
		return err
	}
	return nil
}

// mapSQLError preserves AWS-spec exception types emitted by executeSQL
// (StatementTimeoutException, DatabaseErrorException) and wraps any other
// engine-emitted error in DatabaseErrorException so callers see the
// correct exception name rather than a generic BadRequestException.
func mapSQLError(err error) error {
	if err == nil {
		return nil
	}
	var awsErr *awserrors.AWSError
	if errors.As(err, &awsErr) {
		return awsErr
	}
	return databaseError(err.Error())
}

func parseRequest(req *request.ParsedRequest, v interface{}) error {
	if req == nil || len(req.Body) == 0 {
		return invalidParam("request body is empty")
	}
	if err := json.Unmarshal(req.Body, v); err != nil {
		return badRequest(fmt.Sprintf("failed to parse request: %v", err))
	}
	return nil
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

func (s *RDSDataService) cleanupLoop(ctx context.Context) {
	// Poll at half the idle TTL so an expired transaction is reaped within
	// ~90 seconds of becoming idle (well under the 3-minute AWS deadline).
	// A 30-second tick gives us roughly two sweeps per idle window.
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			func() {
				defer func() {
					if re := recover(); re != nil {
						logs.Error("rdsdata cleanup panic recovered", logs.Any("panic", re))
					}
				}()
				s.purgeExpired()
			}()
		}
	}
}

func (s *RDSDataService) purgeExpired() {
	now := time.Now()
	var expired []*staleEntry
	s.mu.Lock()
	for id, entry := range s.transactions {
		if !entry.isExpired(now) {
			continue
		}
		// Non-blocking check: if background ContinueAfterTimeout
		// statements are still running, skip this entry and leave it
		// in the map. The next tick will pick it up once bgCount
		// reaches zero. This prevents the cleanup ticker from blocking
		// for up to maxBgStatementTime on a single long-running DDL.
		if entry.bgCount.Load() > 0 {
			continue
		}
		expired = append(expired, entry)
		delete(s.transactions, id)
	}
	s.mu.Unlock()

	// AWS Data API rolls back idle / max-life transactions automatically.
	// We issue an explicit ROLLBACK against each expired session so that
	// engine-side state (open transaction, locks, deferred constraints) is
	// released deterministically rather than relying on GC of the
	// pebbleSession. Errors are best-effort: the entry has already been
	// removed from the map, so a failure here only leaks engine state.
	//
	// No bgWg.Wait is needed here because every entry in the expired
	// slice has bgCount == 0 (checked above under s.mu), which means all
	// background statements have completed.
	for _, entry := range expired {
		rollbackCtx := entry.sqlCtx
		if rollbackCtx == nil {
			rollbackCtx = newSQLContext(entry.database)
		}
		if _, err := executeSQL(entry.engine, rollbackCtx, "ROLLBACK", false, ""); err != nil {
			logs.Warn("rdsdata: failed to roll back expired transaction",
				logs.Err(err),
				logs.String("database", entry.database))
		}
	}
}

// truncateSQL shortens a SQL string for inclusion in request logs. The limit
// is generous (500 chars) so that the logged prefix carries enough context
// to identify the statement (table name, primary operation) without
// disclosing full parameter payloads.
func truncateSQL(sql string) string {
	if len(sql) > 500 {
		return sql[:500] + "..."
	}
	return sql
}
