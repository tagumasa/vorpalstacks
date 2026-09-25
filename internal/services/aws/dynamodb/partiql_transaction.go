package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ExecuteTransaction executes a TransactWriteItems operation for PartiQL statements.
func (s *DynamoDBService) ExecuteTransaction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statements, ok := req.Parameters["TransactStatements"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	return s.executeTransactionCore(ctx, reqCtx, executeTransactionInput{
		TransactStatements: statements,
		Parameters:         req.Parameters,
	})
}

// BatchExecuteStatement executes multiple PartiQL statements in a batch.
func (s *DynamoDBService) BatchExecuteStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statements, ok := req.Parameters["Statements"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	return s.batchExecuteStatementCore(ctx, reqCtx, batchExecuteStatementInput{
		Statements: statements,
		Parameters: req.Parameters,
	})
}
