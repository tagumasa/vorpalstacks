package rdsdata

import (
	"context"
	"fmt"
	"time"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/google/uuid"
)

// ExecuteStatementForInvoker is the EventBus invoker entry point for ExecuteStatement.
func (s *RDSDataService) ExecuteStatementForInvoker(ctx context.Context, resourceArn, secretArn, database, schema, sqlStr string, includeResultMetadata bool, formatRecordsAs string) (interface{}, error) {
	if err := validateCommon(resourceArn, secretArn, database, schema); err != nil {
		return nil, err
	}
	if err := validateSQL(sqlStr); err != nil {
		return nil, err
	}
	engine, _, err := s.resolveEngine(resourceArn)
	if err != nil {
		return nil, err
	}
	if err := s.validateCredentials(ctx, secretArn); err != nil {
		return nil, err
	}

	return executeSQL(engine, newSQLContext(database), sqlStr, includeResultMetadata, formatRecordsAs)
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
	engine, _, err := s.resolveEngine(resourceArn)
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

	sqlCtx := entry.sqlCtx
	if sqlCtx == nil {
		sqlCtx = newSQLContext(entry.database)
	}
	return executeSQL(engine, sqlCtx, sqlStr, includeResultMetadata, formatRecordsAs)
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

	if _, err := executeSQL(engine, sqlCtx, "START TRANSACTION", false, ""); err != nil {
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

	s.mu.Lock()
	entry, ok := s.transactions[transactionId]
	if ok {
		delete(s.transactions, transactionId)
	}
	s.mu.Unlock()

	if !ok {
		return transactionNotFound(transactionId)
	}

	commitCtx := entry.sqlCtx
	if commitCtx == nil {
		commitCtx = newSQLContext(entry.database)
	}
	_, err := executeSQL(entry.engine, commitCtx, "COMMIT", false, "")
	return err
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

	s.mu.Lock()
	entry, ok := s.transactions[transactionId]
	if ok {
		delete(s.transactions, transactionId)
	}
	s.mu.Unlock()

	if !ok {
		return transactionNotFound(transactionId)
	}

	rollbackCtx := entry.sqlCtx
	if rollbackCtx == nil {
		rollbackCtx = newSQLContext(entry.database)
	}
	_, err := executeSQL(entry.engine, rollbackCtx, "ROLLBACK", false, "")
	return err
}
