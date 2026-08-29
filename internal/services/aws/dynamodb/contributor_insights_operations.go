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

	indexName := request.GetStringParam(req.Parameters, "IndexName")
	if indexName != "" {
		if !validateIndexName(indexName) {
			return nil, ErrInvalidParameter
		}
	}

	status := "DISABLED"
	if table.ContributorInsightsEnabled {
		status = "ENABLED"
	}

	result := map[string]interface{}{
		"TableName":                 table.Name,
		"ContributorInsightsStatus": status,
	}
	if ruleNames := ContributorInsightsRuleNames(table); len(ruleNames) > 0 {
		result["ContributorInsightsRuleList"] = ruleNames
	}
	if !table.ContributorInsightsUpdatedAt.IsZero() {
		result["LastUpdateDateTime"] = table.ContributorInsightsUpdatedAt.Unix()
	}
	if table.ContributorInsightsMode != "" {
		result["ContributorInsightsMode"] = table.ContributorInsightsMode
	}
	if indexName != "" {
		result["IndexName"] = indexName
	}
	return result, nil
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
