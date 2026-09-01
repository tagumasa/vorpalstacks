package rdsdata

import (
	"context"
	"encoding/json"
	"fmt"
)

// ExecuteStatementForInvoker is the EventBus invoker entry point for ExecuteStatement.
func (s *RDSDataService) ExecuteStatementForInvoker(ctx context.Context, resourceArn, secretArn, database, schema, sqlStr string, includeResultMetadata bool, formatRecordsAs string) (interface{}, error) {
	return s.executeStatementCore(ctx, &ExecuteStatementInput{
		ResourceArn:           resourceArn,
		SecretArn:             secretArn,
		Database:              database,
		Schema:                schema,
		Sql:                   sqlStr,
		IncludeResultMetadata: includeResultMetadata,
		FormatRecordsAs:       formatRecordsAs,
	})
}

// ExecuteStatementInTxForInvoker runs a SQL statement within an existing transaction.
func (s *RDSDataService) ExecuteStatementInTxForInvoker(ctx context.Context, resourceArn, secretArn, transactionId, sqlStr string, includeResultMetadata bool, formatRecordsAs string) (interface{}, error) {
	return s.executeStatementInTxCore(ctx, &ExecuteStatementInput{
		ResourceArn:           resourceArn,
		SecretArn:             secretArn,
		Sql:                   sqlStr,
		TransactionID:         transactionId,
		IncludeResultMetadata: includeResultMetadata,
		FormatRecordsAs:       formatRecordsAs,
	})
}

// BeginTransactionForInvoker is the EventBus invoker entry point for BeginTransaction.
func (s *RDSDataService) BeginTransactionForInvoker(ctx context.Context, resourceArn, secretArn, database, schema string) (string, error) {
	resp, err := s.beginTransactionCore(ctx, &BeginTransactionInput{
		ResourceArn: resourceArn,
		SecretArn:   secretArn,
		Database:    database,
		Schema:      schema,
	})
	if err != nil {
		return "", err
	}
	return resp.(*BeginTransactionResponse).TransactionID, nil
}

// CommitTransactionForInvoker is the EventBus invoker entry point for CommitTransaction.
func (s *RDSDataService) CommitTransactionForInvoker(ctx context.Context, resourceArn, secretArn, transactionId string) error {
	return s.commitTransactionCore(ctx, &CommitTransactionInput{
		ResourceArn:   resourceArn,
		SecretArn:     secretArn,
		TransactionID: transactionId,
	})
}

// RollbackTransactionForInvoker is the EventBus invoker entry point for RollbackTransaction.
func (s *RDSDataService) RollbackTransactionForInvoker(ctx context.Context, resourceArn, secretArn, transactionId string) error {
	return s.rollbackTransactionCore(ctx, &RollbackTransactionInput{
		ResourceArn:   resourceArn,
		SecretArn:     secretArn,
		TransactionID: transactionId,
	})
}

// BatchExecuteStatementForInvoker is the EventBus invoker entry point for
// BatchExecuteStatement. It runs the same SQL once per parameter set,
// collecting generated fields for each execution.
func (s *RDSDataService) BatchExecuteStatementForInvoker(ctx context.Context, resourceArn, secretArn, database, schema, sqlStr string, parameterSets [][]SqlParameter) (interface{}, error) {
	return s.batchExecuteStatementCore(ctx, &BatchExecuteStatementInput{
		ResourceArn:   resourceArn,
		SecretArn:     secretArn,
		Database:      database,
		Schema:        schema,
		Sql:           sqlStr,
		ParameterSets: parameterSets,
	})
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
