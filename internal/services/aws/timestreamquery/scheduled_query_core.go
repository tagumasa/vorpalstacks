package timestreamquery

import (
	"strconv"
	"time"

	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// ---------------------------------------------------------------------------
// DTO types — service-layer representations that decouple the admin handler
// from store types. The admin handler only sees these types and never imports
// store packages directly.
// ---------------------------------------------------------------------------

// ScheduledQuerySummary is the projection of a ScheduledQuery used by the
// ListScheduledQueries admin console operation.
type ScheduledQuerySummary struct {
	ARN                      string
	Name                     string
	CreationTime             time.Time
	State                    string
	LastRunStatus            string
	PreviousRunTime          time.Time
	NextRunTime              time.Time
	ErrorReportConfiguration *tsstore.ErrorReportConfiguration
	TargetConfiguration      *tsstore.TargetConfiguration
}

// ListScheduledQueriesInput is the DTO input for the ListScheduledQueries
// admin console operation.
type ListScheduledQueriesInput struct {
	Region     string
	MaxResults int
	NextToken  int
}

// ListScheduledQueriesResult is the DTO result for the ListScheduledQueries
// admin console operation.
type ListScheduledQueriesResult struct {
	Summaries []*ScheduledQuerySummary
	NextToken string
}

// ---------------------------------------------------------------------------
// Core methods — the admin handler delegates to these, passing the store
// group obtained via admin_handler_convert.go.
// ---------------------------------------------------------------------------

// ListScheduledQueriesCore returns a paginated list of scheduled query
// summaries for the admin console. It operates on the provided store group
// and returns service-layer DTOs.
func (s *TimestreamQueryService) ListScheduledQueriesCore(stores *tsQueryStores, input ListScheduledQueriesInput) (*ListScheduledQueriesResult, error) {
	queries, err := stores.scheduledQueryStore.ListScheduledQueries()
	if err != nil {
		return nil, ErrInternalServer
	}

	maxResults := input.MaxResults
	if maxResults <= 0 {
		maxResults = maxListScheduledQueries
	}
	if maxResults > maxListScheduledQueries {
		maxResults = maxListScheduledQueries
	}

	offset := input.NextToken
	if offset < 0 {
		offset = 0
	}

	var summaries []*ScheduledQuerySummary
	for i, sq := range queries {
		if i < offset {
			continue
		}
		if len(summaries) >= maxResults {
			break
		}
		summary := &ScheduledQuerySummary{
			ARN:                      sq.ARN,
			Name:                     sq.Name,
			CreationTime:             sq.CreationTime,
			State:                    string(sq.ScheduledQueryStatus),
			LastRunStatus:            sq.LastRunStatus,
			PreviousRunTime:          sq.PreviousRunTime,
			NextRunTime:              sq.NextRunTime,
			ErrorReportConfiguration: sq.ErrorReportConfiguration,
			TargetConfiguration:      sq.TargetConfiguration,
		}
		summaries = append(summaries, summary)
	}

	result := &ListScheduledQueriesResult{
		Summaries: summaries,
	}

	if offset+maxResults < len(queries) {
		result.NextToken = strconv.Itoa(offset + maxResults)
	}

	return result, nil
}
