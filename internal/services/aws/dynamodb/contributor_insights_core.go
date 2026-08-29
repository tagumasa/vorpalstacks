package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
)

// ---------------------------------------------------------------------------
// Contributor Insights Core — single validation + persistence path for the
// contributor insights list and update operations.
//
// Both the HTTP API handlers (contributor_insights_operations.go) and any
// future admin handler delegate to these methods to ensure identical
// behaviour.
// ---------------------------------------------------------------------------

// listContributorInsightsInput carries the raw wire parameters for
// ListContributorInsights.
type listContributorInsightsInput struct {
	Parameters map[string]interface{}
}

// listContributorInsightsCore returns the contributor insights summaries,
// optionally scoped to one table.
func (s *DynamoDBService) listContributorInsightsCore(ctx context.Context, reqCtx *request.RequestContext, in listContributorInsightsInput) (interface{}, error) {
	tableName := request.GetStringParam(in.Parameters, "TableName")
	maxResults := listContributorMaxLimit
	if _, ok := in.Parameters["MaxResults"]; ok {
		v := request.GetIntParam(in.Parameters, "MaxResults")
		if !validateListContributorInsightsLimit(v) {
			return nil, ErrInvalidParameter
		}
		maxResults = v
	}
	nextToken := pagination.GetMarker(in.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var tables []string
	if tableName != "" {
		if _, err := s.validateAndGetTable(reqCtx, in.Parameters); err != nil {
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

// updateContributorInsightsInput carries the raw wire parameters for
// UpdateContributorInsights.
type updateContributorInsightsInput struct {
	Parameters map[string]interface{}
}

// updateContributorInsightsCore enables or disables contributor insights
// for the named table.
func (s *DynamoDBService) updateContributorInsightsCore(ctx context.Context, reqCtx *request.RequestContext, in updateContributorInsightsInput) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, in.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	action, ok := in.Parameters["ContributorInsightsAction"].(string)
	if !ok || (action != "ENABLE" && action != "DISABLE") {
		return nil, ErrInvalidParameter
	}
	enabled := action == "ENABLE"

	mode := request.GetStringParam(in.Parameters, "ContributorInsightsMode")
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
