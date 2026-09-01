// Package rdsdata implements the AWS RDS Data API service, providing SQL
// execution against vmysql instances via the go-mysql-server sqle.Engine.
package rdsdata

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/sql"

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
	// execMu serialises SQL execution within the same transaction.
	// go-mysql-server's *sql.Context is not thread-safe; without this
	// mutex, two concurrent ExecuteStatement calls sharing the same
	// TransactionID would race on the shared sqlCtx.
	execMu sync.Mutex
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
	return s.executeStatementCore(ctx, &input)
}

// BatchExecuteStatement runs a SQL statement with multiple parameter sets.
func (s *RDSDataService) BatchExecuteStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input BatchExecuteStatementInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
	return s.batchExecuteStatementCore(ctx, &input)
}

// ExecuteSql runs one or more SQL statements (deprecated operation).
func (s *RDSDataService) ExecuteSql(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input ExecuteSqlInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
	return s.executeSqlCore(ctx, &input)
}

// BeginTransaction starts a new transaction and returns a transaction ID.
func (s *RDSDataService) BeginTransaction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input BeginTransactionInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
	return s.beginTransactionCore(ctx, &input)
}

// CommitTransaction commits a transaction.
func (s *RDSDataService) CommitTransaction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input CommitTransactionInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
	if err := s.commitTransactionCore(ctx, &input); err != nil {
		return nil, err
	}
	return &CommitTransactionResponse{TransactionStatus: "COMMIT"}, nil
}

// RollbackTransaction rolls back a transaction.
func (s *RDSDataService) RollbackTransaction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input RollbackTransactionInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
	if err := s.rollbackTransactionCore(ctx, &input); err != nil {
		return nil, err
	}
	return &RollbackTransactionResponse{TransactionStatus: "ROLLBACK"}, nil
}

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
		// Acquire execMu to avoid racing with an in-flight ExecuteStatement
		// that holds execMu but was not yet visible to the bgCount check
		// above (normal statements do not increment bgCount).
		entry.execMu.Lock()
		rollbackCtx := entry.sqlCtx
		if rollbackCtx == nil {
			rollbackCtx = newSQLContext(entry.database)
		}
		if _, err := executeSQL(entry.engine, rollbackCtx, "ROLLBACK", false, "", ""); err != nil {
			logs.Warn("rdsdata: failed to roll back expired transaction",
				logs.Err(err),
				logs.String("database", entry.database))
		}
		entry.execMu.Unlock()
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
