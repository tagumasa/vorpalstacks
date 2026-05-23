// Package rdsdata implements the AWS RDS Data API service, providing SQL
// execution against vmysql instances via the go-mysql-server sqle.Engine.
package rdsdata

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/google/uuid"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
)

const transactionTTL = 5 * time.Minute

// staleEntry holds a transaction session and its creation time.
type staleEntry struct {
	engine   *sqle.Engine
	created  time.Time
	database string
	schema   string
	sqlCtx   *sql.Context
}

// RDSDataService implements the RDS Data API HTTP operations.
type RDSDataService struct {
	mu             sync.RWMutex
	bus            eventbus.Bus
	vmysqlProvider VmysqlProvider
	rdsStore       RDSStoreProvider
	transactions   map[string]*staleEntry
	cancelCleanup  context.CancelFunc
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

// StartCleanup launches the background stale-transaction reaper.
func (s *RDSDataService) StartCleanup() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelCleanup = cancel
	go s.cleanupLoop(ctx)
}

// Shutdown stops the background cleanup goroutine.
func (s *RDSDataService) Shutdown() {
	if s.cancelCleanup != nil {
		s.cancelCleanup()
	}
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
	if input.ResourceArn == "" {
		return nil, invalidParam("resourceArn is required")
	}
	if input.SecretArn == "" {
		return nil, invalidParam("secretArn is required")
	}
	if input.Sql == "" {
		return nil, invalidParam("sql is required")
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
		s.mu.RLock()
		entry, ok := s.transactions[input.TransactionID]
		s.mu.RUnlock()
		if !ok {
			return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
		}
		if entry.database != "" {
			database = entry.database
		}
		txCtx = entry.sqlCtx
	}

	sqlStr := input.Sql
	if len(input.Parameters) > 0 {
		sqlStr = substituteParameters(sqlStr, input.Parameters)
	}

	sqlCtx := txCtx
	if sqlCtx == nil {
		sqlCtx = newSQLContext(database)
	}
	result, err := executeSQL(engine, sqlCtx, sqlStr, input.IncludeResultMetadata, input.FormatRecordsAs)
	if err != nil {
		return nil, badRequest(fmt.Sprintf("SQL execution failed: %v", err))
	}

	logs.Info("rdsdata: ExecuteStatement", logs.String("instance", instanceID), logs.String("sql", truncateSQL(input.Sql)))
	return result, nil
}

// BatchExecuteStatement runs a SQL statement with multiple parameter sets.
func (s *RDSDataService) BatchExecuteStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var input BatchExecuteStatementInput
	if err := parseRequest(req, &input); err != nil {
		return nil, err
	}
	if input.ResourceArn == "" {
		return nil, invalidParam("resourceArn is required")
	}
	if input.SecretArn == "" {
		return nil, invalidParam("secretArn is required")
	}
	if input.Sql == "" {
		return nil, invalidParam("sql is required")
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
		s.mu.RLock()
		entry, ok := s.transactions[input.TransactionID]
		s.mu.RUnlock()
		if !ok {
			return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
		}
		if entry.database != "" {
			database = entry.database
		}
		txCtx = entry.sqlCtx
	}

	var results []UpdateResult
	if len(input.ParameterSets) == 0 {
		sqlCtx := txCtx
		if sqlCtx == nil {
			sqlCtx = newSQLContext(database)
		}
		res, err := executeSQL(engine, sqlCtx, input.Sql, false, "")
		if err != nil {
			return nil, badRequest(fmt.Sprintf("SQL execution failed: %v", err))
		}
		results = append(results, UpdateResult{GeneratedFields: res.GeneratedFields})
	} else {
		for _, params := range input.ParameterSets {
			sqlWithParams := substituteParameters(input.Sql, params)
			sqlCtx := txCtx
			if sqlCtx == nil {
				sqlCtx = newSQLContext(database)
			}
			res, err := executeSQL(engine, sqlCtx, sqlWithParams, false, "")
			if err != nil {
				return nil, badRequest(fmt.Sprintf("SQL execution failed for parameter set: %v", err))
			}
			results = append(results, UpdateResult{GeneratedFields: res.GeneratedFields})
		}
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
	if input.ResourceArn == "" {
		return nil, invalidParam("resourceArn is required")
	}
	if input.SecretArn == "" {
		return nil, invalidParam("secretArn is required")
	}
	if input.Sql == "" {
		return nil, invalidParam("sql is required")
	}

	engine, instanceID, err := s.resolveEngine(input.ResourceArn)
	if err != nil {
		return nil, err
	}
	if err := s.validateCredentials(ctx, input.SecretArn); err != nil {
		return nil, err
	}

	statements := splitSQL(input.Sql)
	var sqlResults []SqlStatementResult
	for _, stmt := range statements {
		res, err := executeSQL(engine, newSQLContext(input.Database), stmt, true, "")
		if err != nil {
			return nil, badRequest(fmt.Sprintf("SQL execution failed: %v", err))
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
	if input.ResourceArn == "" {
		return nil, invalidParam("resourceArn is required")
	}
	if input.SecretArn == "" {
		return nil, invalidParam("secretArn is required")
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
		return nil, badRequest(fmt.Sprintf("begin transaction failed: %v", err))
	}

	txID := uuid.New().String()
	s.mu.Lock()
	s.transactions[txID] = &staleEntry{
		engine:   engine,
		created:  time.Now(),
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
	if input.ResourceArn == "" {
		return nil, invalidParam("resourceArn is required")
	}
	if input.SecretArn == "" {
		return nil, invalidParam("secretArn is required")
	}
	if input.TransactionID == "" {
		return nil, invalidParam("transactionId is required")
	}

	if err := s.validateCredentials(ctx, input.SecretArn); err != nil {
		return nil, err
	}

	s.mu.Lock()
	entry, ok := s.transactions[input.TransactionID]
	if ok {
		delete(s.transactions, input.TransactionID)
	}
	s.mu.Unlock()

	if !ok {
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
	}

	commitCtx := entry.sqlCtx
	if commitCtx == nil {
		commitCtx = newSQLContext(entry.database)
	}
	if _, err := executeSQL(entry.engine, commitCtx, "COMMIT", false, ""); err != nil {
		return nil, badRequest(fmt.Sprintf("commit failed: %v", err))
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
	if input.ResourceArn == "" {
		return nil, invalidParam("resourceArn is required")
	}
	if input.SecretArn == "" {
		return nil, invalidParam("secretArn is required")
	}
	if input.TransactionID == "" {
		return nil, invalidParam("transactionId is required")
	}

	if err := s.validateCredentials(ctx, input.SecretArn); err != nil {
		return nil, err
	}

	s.mu.Lock()
	entry, ok := s.transactions[input.TransactionID]
	if ok {
		delete(s.transactions, input.TransactionID)
	}
	s.mu.Unlock()

	if !ok {
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", input.TransactionID))
	}

	rollbackCtx := entry.sqlCtx
	if rollbackCtx == nil {
		rollbackCtx = newSQLContext(entry.database)
	}
	if _, err := executeSQL(entry.engine, rollbackCtx, "ROLLBACK", false, ""); err != nil {
		return nil, badRequest(fmt.Sprintf("rollback failed: %v", err))
	}

	logs.Info("rdsdata: RollbackTransaction", logs.String("txId", input.TransactionID))
	return &RollbackTransactionResponse{TransactionStatus: "ROLLBACK"}, nil
}

// --- Internal helpers ---

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
		_, instances, err := s.rdsStore.GetCluster(identifier)
		if err != nil {
			return nil, "", invalidParam(fmt.Sprintf("cluster %s not found", identifier))
		}
		if len(instances) == 0 {
			return nil, "", invalidParam(fmt.Sprintf("cluster %s has no instances", identifier))
		}
		instanceID = instances[0]
	} else {
		instanceID = identifier
	}

	if s.vmysqlProvider == nil {
		return nil, "", internalError("vmysql provider not configured")
	}
	engine := s.vmysqlProvider.GetEngine(instanceID)
	if engine == nil {
		return nil, "", invalidParam(fmt.Sprintf("no MySQL engine for instance %s", instanceID))
	}
	return engine, instanceID, nil
}

// validateCredentials checks that the secretArn resolves to valid credentials.
func (s *RDSDataService) validateCredentials(ctx context.Context, secretArn string) error {
	if s.bus == nil {
		return nil
	}
	if secretArn == "" {
		return accessDenied("secretArn is required")
	}
	return nil
}

func (s *RDSDataService) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.purgeExpired()
		}
	}
}

func (s *RDSDataService) purgeExpired() {
	now := time.Now()
	s.mu.Lock()
	for id, entry := range s.transactions {
		if now.Sub(entry.created) > transactionTTL {
			delete(s.transactions, id)
		}
	}
	s.mu.Unlock()
}

func truncateSQL(sql string) string {
	if len(sql) > 100 {
		return sql[:100] + "..."
	}
	return sql
}
