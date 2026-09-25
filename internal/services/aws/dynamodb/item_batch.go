package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// BatchGetItem retrieves multiple items from one or more tables in a single request.
func (s *DynamoDBService) BatchGetItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	requestItems, _ := req.Parameters["RequestItems"].(map[string]interface{})
	return s.batchGetItemCore(ctx, reqCtx, batchGetItemInput{
		RequestItems: requestItems,
		Parameters:   req.Parameters,
	})
}

// BatchWriteItem inserts, updates, or deletes multiple items across one or more tables.
func (s *DynamoDBService) BatchWriteItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	requestItems, _ := req.Parameters["RequestItems"].(map[string]interface{})
	return s.batchWriteItemCore(ctx, reqCtx, batchWriteItemInput{
		RequestItems: requestItems,
		Parameters:   req.Parameters,
	})
}
