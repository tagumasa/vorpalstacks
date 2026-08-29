// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateGlobalTable creates a new global table.
func (s *DynamoDBService) CreateGlobalTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.createGlobalTableCore(ctx, reqCtx, createGlobalTableInput{
		Parameters: req.Parameters,
	})
}

// DescribeGlobalTable returns information about a global table.
func (s *DynamoDBService) DescribeGlobalTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeGlobalTableCore(ctx, reqCtx, describeGlobalTableInput{
		Parameters: req.Parameters,
	})
}

// DescribeGlobalTableSettings returns the settings of a global table.
func (s *DynamoDBService) DescribeGlobalTableSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeGlobalTableSettingsCore(ctx, reqCtx, describeGlobalTableSettingsInput{
		Parameters: req.Parameters,
	})
}

// ListGlobalTables lists the global tables for a given account.
func (s *DynamoDBService) ListGlobalTables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listGlobalTablesCore(ctx, reqCtx, listGlobalTablesInput{
		Parameters: req.Parameters,
	})
}

// UpdateGlobalTable updates a global table.
func (s *DynamoDBService) UpdateGlobalTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateGlobalTableCore(ctx, reqCtx, updateGlobalTableInput{
		Parameters: req.Parameters,
	})
}

// UpdateGlobalTableSettings updates the settings of a global table.
func (s *DynamoDBService) UpdateGlobalTableSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateGlobalTableSettingsCore(ctx, reqCtx, updateGlobalTableSettingsInput{
		Parameters: req.Parameters,
	})
}
