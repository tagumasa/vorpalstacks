package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
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

// contributorInsightsSummary renders one table's summary: the status its
// record carries (DISABLED unless insights are enabled) and the mode when
// one is recorded.
func contributorInsightsSummary(t *dbstore.Table) map[string]interface{} {
	status := "DISABLED"
	if t.ContributorInsightsEnabled {
		status = "ENABLED"
	}
	summary := map[string]interface{}{
		"TableName":                 t.Name,
		"ContributorInsightsStatus": status,
	}
	if t.ContributorInsightsMode != "" {
		summary["ContributorInsightsMode"] = t.ContributorInsightsMode
	}
	return summary
}

// listContributorInsightsCore returns the contributor insights summaries,
// optionally scoped to one table. The summaries describe resources carrying
// contributor-insights state: the unfiltered walk keeps the tables whose
// record an UpdateContributorInsights call has stamped (the update
// timestamp), reading the state from the listed records themselves — the
// scoped call answers the named table's own summary the way
// DescribeContributorInsights does, DISABLED included for a table insights
// were never configured on.
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

	summaries := make([]map[string]interface{}, 0)
	if tableName != "" {
		table, err := s.validateAndGetTable(reqCtx, in.Parameters)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, contributorInsightsSummary(table))
		return pagination.BuildListResponse("ContributorInsightsSummaries", summaries, ""), nil
	}

	// The unfiltered walk pages through the table records, keeping those
	// carrying insights state. A result page fills slower than the record
	// pages it consumes (records without state are skipped), so the walk
	// continues across record pages until the result page fills or the
	// tables run out; the continuation marker is the last table name
	// consumed, examined or not, so the next call resumes after it and a
	// final marker may answer an empty page when only state-less records
	// followed it.
	remaining := maxResults
	marker := nextToken
	for remaining > 0 {
		page, next, err := store.Tables().List(marker, remaining+1)
		if err != nil {
			return nil, err
		}
		if len(page) == 0 {
			break
		}
		filled := false
		for _, t := range page {
			marker = t.Name
			if t.ContributorInsightsUpdatedAt.IsZero() {
				continue
			}
			summaries = append(summaries, contributorInsightsSummary(t))
			remaining--
			if remaining == 0 {
				filled = true
				break
			}
		}
		if filled {
			return pagination.BuildListResponse("ContributorInsightsSummaries", summaries, marker), nil
		}
		if next == "" {
			break
		}
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
	if err := store.Tables().SetContributorInsights(tableName, enabled, dbstore.ContributorInsightsMode(mode)); err != nil {
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

// describeContributorInsightsCore owns the DescribeContributorInsights read:
// the optional index name is validated against the resource-name rules, and
// the status, rules, and mode come from the resolved table.
func (s *DynamoDBService) describeContributorInsightsCore(table *dbstore.Table, params map[string]interface{}) (map[string]interface{}, error) {
	indexName := request.GetStringParam(params, "IndexName")
	if indexName != "" {
		if !validateResourceName(indexName) {
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
