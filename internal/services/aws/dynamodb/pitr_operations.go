// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// DescribeContinuousBackups returns the continuous backup settings for a table.
func (s *DynamoDBService) DescribeContinuousBackups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeContinuousBackupsCore(ctx, reqCtx, describeContinuousBackupsInput{
		Parameters: req.Parameters,
	})
}

// UpdateContinuousBackups enables or disables continuous backup for a table.
func (s *DynamoDBService) UpdateContinuousBackups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateContinuousBackupsCore(ctx, reqCtx, updateContinuousBackupsInput{
		Parameters: req.Parameters,
	})
}
