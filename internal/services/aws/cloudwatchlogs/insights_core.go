package cloudwatchlogs

import (
	"fmt"
	"time"
)

// StartQueryInput holds validated parameters for StartQuery.
type StartQueryInput struct {
	StartTime           int64
	EndTime             int64
	QueryString         string
	QueryLanguage       string
	LogGroupNames       []string
	LogGroupIdentifiers []string
	Limit               int64
	Region              string
}

// startQueryCore validates input, creates the query state, and launches
// the async execution goroutine.
func (s *LogsService) startQueryCore(input *StartQueryInput) (string, error) {
	if input.QueryString == "" || input.StartTime <= 0 || input.EndTime <= 0 {
		return "", ErrMissingParameter
	}
	if len(input.LogGroupNames) == 0 && len(input.LogGroupIdentifiers) == 0 {
		return "", ErrMissingParameter
	}
	if err := validateQueryPipeline(input.QueryString); err != nil {
		return "", err
	}

	allGroups := input.LogGroupNames
	allGroups = append(allGroups, input.LogGroupIdentifiers...)

	queryId := fmt.Sprintf("query-%d", time.Now().UnixNano())

	qs := &queryState{
		queryId:             queryId,
		logGroupNames:       allGroups,
		logGroupIdentifiers: input.LogGroupIdentifiers,
		startTime:           input.StartTime,
		endTime:             input.EndTime,
		queryString:         input.QueryString,
		queryLanguage:       input.QueryLanguage,
		status:              "Running",
		createdAt:           time.Now(),
	}
	s.queries.Store(queryId, qs)

	limit := input.Limit
	if limit <= 0 {
		limit = 10000
	}
	go s.executeQuery(input.Region, queryId, input.QueryString, allGroups, input.StartTime, input.EndTime, limit)

	return queryId, nil
}

// stopQueryCore validates input and cancels a running query.
func (s *LogsService) stopQueryCore(queryId string) error {
	if queryId == "" {
		return ErrMissingParameter
	}

	val, ok := s.queries.Load(queryId)
	if !ok {
		return NewLogsError("ResourceNotFoundException",
			fmt.Sprintf("Query %s not found", queryId), 400)
	}

	qs := val.(*queryState)
	qs.status = "Cancelled"
	return nil
}

// DescribeQueriesInput holds parameters for DescribeQueries.
type DescribeQueriesInput struct {
	StatusFilter string
	LogGroupName string
	NextToken    string
	MaxResults   int32
}

// describeQueriesCore validates input and lists queries with filtering.
func (s *LogsService) describeQueriesCore(input *DescribeQueriesInput) ([]*queryState, string, error) {
	maxResults := input.MaxResults
	if maxResults <= 0 {
		maxResults = 50
	}

	var allQueries []*queryState
	s.queries.Range(func(key, value interface{}) bool {
		qs := value.(*queryState)
		if input.StatusFilter != "" && qs.status != input.StatusFilter {
			return true
		}
		if input.LogGroupName != "" {
			found := false
			for _, n := range qs.logGroupNames {
				if n == input.LogGroupName {
					found = true
					break
				}
			}
			if !found {
				return true
			}
		}
		allQueries = append(allQueries, qs)
		return true
	})

	nextToken := input.NextToken
	if nextToken == "" {
		nextToken = "0"
	}
	result := paginateQuerySlice(allQueries, nextToken, int(maxResults))
	return result.Items, result.NextToken, nil
}

// GetQueryResultsInput holds parameters for GetQueryResults.
type GetQueryResultsInput struct {
	QueryId   string
	NextToken string
	MaxItems  int32
}

// GetQueryResultsResult holds the query results page.
type GetQueryResultsResult struct {
	QueryId       string
	Status        string
	QueryLanguage string
	Results       []queryResultRow
	Stats         queryStats
	NextToken     string
}

// getQueryResultsCore validates input and retrieves paginated query results.
func (s *LogsService) getQueryResultsCore(input *GetQueryResultsInput) (*GetQueryResultsResult, error) {
	if input.QueryId == "" {
		return nil, ErrMissingParameter
	}

	val, ok := s.queries.Load(input.QueryId)
	if !ok {
		return nil, NewLogsError("ResourceNotFoundException",
			fmt.Sprintf("Query %s not found", input.QueryId), 400)
	}

	qs := val.(*queryState)

	limit := int(input.MaxItems)
	if limit <= 0 {
		limit = 10000
	}

	offset := 0
	if input.NextToken != "" {
		n, err := parseInt(input.NextToken)
		if err != nil {
			return nil, NewLogsError("InvalidParameterException",
				"Invalid nextToken", 400)
		}
		offset = n
	}

	endIdx := offset + limit
	if endIdx > len(qs.results) {
		endIdx = len(qs.results)
	}
	if offset > len(qs.results) {
		offset = len(qs.results)
	}

	result := &GetQueryResultsResult{
		QueryId:       input.QueryId,
		Status:        qs.status,
		QueryLanguage: qs.queryLanguage,
		Results:       qs.results[offset:endIdx],
		Stats:         qs.stats,
	}

	if endIdx < len(qs.results) {
		result.NextToken = fmt.Sprintf("%d", endIdx)
	}

	return result, nil
}

type queryPageResult struct {
	Items     []*queryState
	NextToken string
}

func paginateQuerySlice(items []*queryState, token string, limit int) *queryPageResult {
	offset := 0
	if n, err := parseInt(token); err == nil {
		offset = n
	}
	if offset >= len(items) {
		return &queryPageResult{Items: []*queryState{}, NextToken: ""}
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	page := items[offset:end]
	nextToken := ""
	if end < len(items) {
		nextToken = fmt.Sprintf("%d", end)
	}
	return &queryPageResult{Items: page, NextToken: nextToken}
}
