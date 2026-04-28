// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
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

// ListContributorInsights lists the contributor insights summaries for tables.
func (s *DynamoDBService) ListContributorInsights(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableName := request.GetStringParam(req.Parameters, "TableName")
	maxResults := pagination.GetMaxItems(req.Parameters, 100, "MaxResults")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var tables []string
	if tableName != "" {
		if _, err := s.validateAndGetTable(reqCtx, req.Parameters); err != nil {
			return nil, err
		}
		tables = []string{tableName}
	} else {
		tableList, _, err := store.Tables().List(nextToken, maxResults+1)
		if err != nil {
			return nil, err
		}
		for _, t := range tableList {
			tables = append(tables, t.Name)
		}
	}

	summaries := make([]map[string]interface{}, 0)
	for _, tn := range tables {
		t, err := store.Tables().Get(tn)
		if err != nil {
			continue
		}
		status := "DISABLED"
		if t.ContributorInsightsEnabled {
			status = "ENABLED"
		}
		summaries = append(summaries, map[string]interface{}{
			"TableName":                 tn,
			"ContributorInsightsStatus": status,
		})
	}

	if len(summaries) > maxResults {
		summaries = summaries[:maxResults]
		return pagination.BuildListResponse("ContributorInsightsSummaries", summaries, summaries[len(summaries)-1]["TableName"].(string)), nil
	}

	return pagination.BuildListResponse("ContributorInsightsSummaries", summaries, ""), nil
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
