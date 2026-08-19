// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

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
	tableName := request.GetStringParam(req.Parameters, "TableName")
	maxResults := listContributorMaxLimit
	if _, ok := req.Parameters["MaxResults"]; ok {
		v := request.GetIntParam(req.Parameters, "MaxResults")
		if !validateListContributorInsightsLimit(v) {
			return nil, ErrInvalidParameter
		}
		maxResults = v
	}
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
		summary := map[string]interface{}{
			"TableName":                 tn,
			"ContributorInsightsStatus": status,
		}
		if t.ContributorInsightsMode != "" {
			summary["ContributorInsightsMode"] = t.ContributorInsightsMode
		}
		summaries = append(summaries, summary)
	}

	if len(summaries) > maxResults && maxResults > 0 {
		summaries = summaries[:maxResults]
		lastTableName, _ := summaries[len(summaries)-1]["TableName"].(string)
		return pagination.BuildListResponse("ContributorInsightsSummaries", summaries, lastTableName), nil
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

	action, ok := req.Parameters["ContributorInsightsAction"].(string)
	if !ok || (action != "ENABLE" && action != "DISABLE") {
		return nil, ErrInvalidParameter
	}
	enabled := action == "ENABLE"

	mode := request.GetStringParam(req.Parameters, "ContributorInsightsMode")
	if !validateContributorInsightsMode(mode) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Tables().SetContributorInsights(tableName, enabled, mode); err != nil {
		return nil, err
	}

	// Return the transition state when the requested value differs from
	// the existing value; otherwise return the steady state.
	alreadyEnabled := table.ContributorInsightsEnabled
	status := "ENABLED"
	if !enabled {
		status = "DISABLED"
	}
	if enabled != alreadyEnabled {
		if enabled {
			status = "ENABLING"
		} else {
			status = "DISABLING"
		}
	}

	return map[string]interface{}{
		"TableName":                 tableName,
		"ContributorInsightsStatus": status,
	}, nil
}
