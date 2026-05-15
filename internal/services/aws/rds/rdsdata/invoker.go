package rdsdata

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ExecuteStatementForInvoker is the EventBus invoker entry point for ExecuteStatement.
func (s *RDSDataService) ExecuteStatementForInvoker(ctx context.Context, resourceArn, secretArn, database, schema, sql string, includeResultMetadata bool, formatRecordsAs string) (interface{}, error) {
	engine, _, err := s.resolveEngine(resourceArn)
	if err != nil {
		return nil, err
	}
	if err := s.validateCredentials(ctx, secretArn); err != nil {
		return nil, err
	}

	return executeSQL(engine, sql, database, includeResultMetadata, formatRecordsAs)
}

// BeginTransactionForInvoker is the EventBus invoker entry point for BeginTransaction.
func (s *RDSDataService) BeginTransactionForInvoker(ctx context.Context, resourceArn, secretArn, database, schema string) (string, error) {
	engine, _, err := s.resolveEngine(resourceArn)
	if err != nil {
		return "", err
	}
	if err := s.validateCredentials(ctx, secretArn); err != nil {
		return "", err
	}

	txID := uuid.New().String()
	s.mu.Lock()
	s.transactions[txID] = &staleEntry{
		engine:   engine,
		created:  time.Now(),
		database: database,
		schema:   schema,
	}
	s.mu.Unlock()

	return txID, nil
}

// CommitTransactionForInvoker is the EventBus invoker entry point for CommitTransaction.
func (s *RDSDataService) CommitTransactionForInvoker(ctx context.Context, resourceArn, secretArn, transactionId string) error {
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

	_, err := executeSQL(entry.engine, "COMMIT", entry.database, false, "")
	return err
}

// RollbackTransactionForInvoker is the EventBus invoker entry point for RollbackTransaction.
func (s *RDSDataService) RollbackTransactionForInvoker(ctx context.Context, resourceArn, secretArn, transactionId string) error {
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

	_, err := executeSQL(entry.engine, "ROLLBACK", entry.database, false, "")
	return err
}
