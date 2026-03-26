// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
)

// DescribeContributorInsights returns the contributor insights status for a table.
func (s *DynamoDBService) DescribeContributorInsights(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	indexName := request.GetStringParam(req.Parameters, "IndexName")

	status := "DISABLED"
	if table.ContributorInsightsEnabled {
		status = "ENABLED"
	}

	result := map[string]interface{}{
		"TableName":                 table.Name,
		"ContributorInsightsStatus": status,
		"LastUpdateDateTime":        time.Now().Unix(),
	}
	if indexName != "" {
		result["IndexName"] = indexName
	}
	return result, nil
}

// ListContributorInsights lists the contributor insights summaries for a table.
func (s *DynamoDBService) ListContributorInsights(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableName := request.GetStringParam(req.Parameters, "TableName")

	if tableName != "" {
		if _, err := s.validateAndGetTable(reqCtx, req.Parameters); err != nil {
			return nil, err
		}
	}

	var summaries []map[string]interface{}
	if tableName != "" {
		summaries = append(summaries, map[string]interface{}{
			"TableName":                 tableName,
			"ContributorInsightsStatus": "ENABLED",
		})
	}

	return map[string]interface{}{
		"ContributorInsightsSummaries": summaries,
	}, nil
}

// UpdateContributorInsights enables or disables contributor insights for a table.
func (s *DynamoDBService) UpdateContributorInsights(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	enabled := true
	if e, ok := req.Parameters["ContributorInsightsAction"].(string); ok {
		if e == "DISABLE" {
			enabled = false
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Tables().SetContributorInsights(tableName, enabled); err != nil {
		return nil, err
	}

	status := "ENABLED"
	if !enabled {
		status = "DISABLED"
	}

	return map[string]interface{}{
		"TableName":                 tableName,
		"ContributorInsightsStatus": status,
	}, nil
}
