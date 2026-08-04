package rdsdata

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
)

// ExecuteStatementForInvoker is the EventBus invoker entry point for ExecuteStatement.
func (s *RDSDataService) ExecuteStatementForInvoker(ctx context.Context, resourceArn, secretArn, database, schema, sqlStr string, includeResultMetadata bool, formatRecordsAs string) (interface{}, error) {
	if err := validateCommon(resourceArn, secretArn, database, schema); err != nil {
		return nil, err
	}
	if err := validateSQL(sqlStr); err != nil {
		return nil, err
	}
	engine, instanceID, err := s.resolveEngine(resourceArn)
	if err != nil {
		return nil, err
	}
	if err := s.validateCredentials(ctx, secretArn); err != nil {
		return nil, err
	}

	return executeSQL(engine, newSQLContext(database), sqlStr, includeResultMetadata, formatRecordsAs, instanceID)
}

// ExecuteStatementInTxForInvoker runs a SQL statement within an existing transaction.
func (s *RDSDataService) ExecuteStatementInTxForInvoker(ctx context.Context, resourceArn, secretArn, transactionId, sqlStr string, includeResultMetadata bool, formatRecordsAs string) (interface{}, error) {
	if err := validateCommon(resourceArn, secretArn, "", ""); err != nil {
		return nil, err
	}
	if err := validateSQL(sqlStr); err != nil {
		return nil, err
	}
	if err := validateTransactionID(transactionId, false); err != nil {
		return nil, err
	}
	engine, instanceID, err := s.resolveEngine(resourceArn)
	if err != nil {
		return nil, err
	}
	if err := s.validateCredentials(ctx, secretArn); err != nil {
		return nil, err
	}

	now := time.Now()
	s.mu.Lock()
	entry, ok := s.transactions[transactionId]
	if ok {
		if entry.isExpired(now) {
			delete(s.transactions, transactionId)
			entry = nil
			ok = false
		} else {
			entry.touch()
		}
	}
	s.mu.Unlock()
	if !ok {
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", transactionId))
	}

	// Serialise concurrent statements on the same transaction.
	entry.execMu.Lock()
	defer entry.execMu.Unlock()

	sqlCtx := entry.sqlCtx
	if sqlCtx == nil {
		sqlCtx = newSQLContext(entry.database)
	}
	return executeSQL(engine, sqlCtx, sqlStr, includeResultMetadata, formatRecordsAs, instanceID)
}

// BeginTransactionForInvoker is the EventBus invoker entry point for BeginTransaction.
func (s *RDSDataService) BeginTransactionForInvoker(ctx context.Context, resourceArn, secretArn, database, schema string) (string, error) {
	if err := validateCommon(resourceArn, secretArn, database, schema); err != nil {
		return "", err
	}
	engine, instanceID, err := s.resolveEngine(resourceArn)
	if err != nil {
		return "", err
	}
	if err := s.validateCredentials(ctx, secretArn); err != nil {
		return "", err
	}

	var sqlCtx *sql.Context
	if s.vmysqlProvider != nil {
		sqlCtx = s.vmysqlProvider.NewContext(instanceID, database)
	}
	if sqlCtx == nil {
		sqlCtx = newSQLContext(database)
	}

	if _, err := executeSQL(engine, sqlCtx, "START TRANSACTION", false, "", instanceID); err != nil {
		return "", mapSQLError(err)
	}

	txID := uuid.New().String()
	now := time.Now()
	s.mu.Lock()
	s.transactions[txID] = &staleEntry{
		engine:   engine,
		created:  now,
		lastSeen: now,
		database: database,
		schema:   schema,
		sqlCtx:   sqlCtx,
	}
	s.mu.Unlock()

	return txID, nil
}

// CommitTransactionForInvoker is the EventBus invoker entry point for CommitTransaction.
func (s *RDSDataService) CommitTransactionForInvoker(ctx context.Context, resourceArn, secretArn, transactionId string) error {
	if err := validateCommon(resourceArn, secretArn, "", ""); err != nil {
		return err
	}
	if err := validateTransactionID(transactionId, false); err != nil {
		return err
	}
	if err := s.validateCredentials(ctx, secretArn); err != nil {
		return err
	}

	return s.commitTransactionCore(transactionId)
}

// RollbackTransactionForInvoker is the EventBus invoker entry point for RollbackTransaction.
func (s *RDSDataService) RollbackTransactionForInvoker(ctx context.Context, resourceArn, secretArn, transactionId string) error {
	if err := validateCommon(resourceArn, secretArn, "", ""); err != nil {
		return err
	}
	if err := validateTransactionID(transactionId, false); err != nil {
		return err
	}
	if err := s.validateCredentials(ctx, secretArn); err != nil {
		return err
	}

	return s.rollbackTransactionCore(transactionId)
}

// BatchExecuteStatementForInvoker is the EventBus invoker entry point for
// BatchExecuteStatement. It runs the same SQL once per parameter set,
// collecting generated fields for each execution.
func (s *RDSDataService) BatchExecuteStatementForInvoker(ctx context.Context, resourceArn, secretArn, database, schema, sqlStr string, parameterSets [][]SqlParameter) (interface{}, error) {
	if err := validateCommon(resourceArn, secretArn, database, schema); err != nil {
		return nil, err
	}
	if err := validateSQL(sqlStr); err != nil {
		return nil, err
	}
	engine, instanceID, err := s.resolveEngine(resourceArn)
	if err != nil {
		return nil, err
	}
	if err := s.validateCredentials(ctx, secretArn); err != nil {
		return nil, err
	}

	if len(parameterSets) == 0 {
		return nil, invalidParam("parameterSets is required and must contain at least one entry")
	}

	var results []UpdateResult
	for _, params := range parameterSets {
		if err := validateParameters(params); err != nil {
			return nil, err
		}
		sqlWithParams := substituteParameters(sqlStr, params)
		res, err := executeSQL(engine, newSQLContext(database), sqlWithParams, false, "", instanceID)
		if err != nil {
			return nil, mapSQLError(err)
		}
		results = append(results, UpdateResult{GeneratedFields: res.GeneratedFields})
	}

	logs.Info("rdsdata: BatchExecuteStatementForInvoker", logs.String("instance", instanceID), logs.Int("paramSets", len(parameterSets)))
	return &BatchExecuteStatementResponse{UpdateResults: results}, nil
}

// SqlParameterFromInterface converts a raw map (as received from AppSync
// VTL templates via the EventBus) to a SqlParameter. The map is expected
// to have "name" (string) and "value" (map with one of longValue,
// stringValue, doubleValue, booleanValue, blobValue, isNull) and an
// optional "typeHint" (string).
func SqlParameterFromInterface(raw interface{}) (SqlParameter, error) {
	data, err := json.Marshal(raw)
	if err != nil {
		return SqlParameter{}, fmt.Errorf("failed to marshal parameter: %w", err)
	}
	var p SqlParameter
	if err := json.Unmarshal(data, &p); err != nil {
		return SqlParameter{}, fmt.Errorf("failed to unmarshal parameter: %w", err)
	}
	return p, nil
}
