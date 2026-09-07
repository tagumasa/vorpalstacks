// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// DescribeContributorInsights returns the contributor insights status for a table.
func (s *DynamoDBService) DescribeContributorInsights(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	return s.describeContributorInsightsCore(table, req.Parameters)
}

// ListContributorInsights lists the contributor insights summaries for tables.
func (s *DynamoDBService) ListContributorInsights(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listContributorInsightsCore(ctx, reqCtx, listContributorInsightsInput{
		Parameters: req.Parameters,
	})
}

// UpdateContributorInsights enables or disables contributor insights for a table.
func (s *DynamoDBService) UpdateContributorInsights(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateContributorInsightsCore(ctx, reqCtx, updateContributorInsightsInput{
		Parameters: req.Parameters,
	})
}
