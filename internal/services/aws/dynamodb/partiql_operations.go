package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ExecuteStatement executes a PartiQL statement.
func (s *DynamoDBService) ExecuteStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.executeStatementCore(ctx, reqCtx, executeStatementInput{
		Parameters: req.Parameters,
	})
}
