package cloudtrail

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
)

// sampleQuery represents a predefined CloudTrail Lake sample query.
type sampleQuery struct {
	Name        string
	Description string
	SQL         string
}

// cloudTrailSampleQueries is a curated list of common CloudTrail Lake queries
// used by SearchSampleQueries and as templates for GenerateQuery.
var cloudTrailSampleQueries = []sampleQuery{
	{
		Name:        "AllManagementEvents",
		Description: "Returns all management events in the event data store.",
		SQL:         "SELECT eventId, eventTime, eventName, eventSource, userIdentity FROM %s",
	},
	{
		Name:        "S3Events",
		Description: "Returns all Amazon S3 API call events.",
		SQL:         "SELECT eventId, eventTime, eventName, resources FROM %s WHERE eventSource = 's3.amazonaws.com'",
	},
	{
		Name:        "IAMEvents",
		Description: "Returns all IAM API call events, useful for auditing access changes.",
		SQL:         "SELECT eventId, eventTime, eventName, userIdentity FROM %s WHERE eventSource = 'iam.amazonaws.com'",
	},
	{
		Name:        "EC2Events",
		Description: "Returns all Amazon EC2 API call events.",
		SQL:         "SELECT eventId, eventTime, eventName, resources FROM %s WHERE eventSource = 'ec2.amazonaws.com'",
	},
	{
		Name:        "ConsoleLogins",
		Description: "Returns all console login events, useful for detecting unauthorised access.",
		SQL:         "SELECT eventId, eventTime, eventName, userIdentity, sourceIpAddress FROM %s WHERE eventName = 'ConsoleLogin'",
	},
	{
		Name:        "RootAccountActivity",
		Description: "Returns all events performed by the root account.",
		SQL:         "SELECT eventId, eventTime, eventName, userIdentity FROM %s WHERE userIdentity LIKE '%root%'",
	},
	{
		Name:        "FailedEvents",
		Description: "Returns all events that resulted in an error.",
		SQL:         "SELECT eventId, eventTime, eventName, errorMessage FROM %s WHERE errorMessage != ''",
	},
	{
		Name:        "ResourceCreation",
		Description: "Returns all resource creation events (Create* API calls).",
		SQL:         "SELECT eventId, eventTime, eventName, resources FROM %s WHERE eventName LIKE 'Create%'",
	},
	{
		Name:        "ResourceDeletion",
		Description: "Returns all resource deletion events (Delete* API calls).",
		SQL:         "SELECT eventId, eventTime, eventName, resources FROM %s WHERE eventName LIKE 'Delete%'",
	},
	{
		Name:        "LambdaInvocations",
		Description: "Returns all AWS Lambda invocation and management events.",
		SQL:         "SELECT eventId, eventTime, eventName FROM %s WHERE eventSource = 'lambda.amazonaws.com'",
	},
}

// GenerateQuery generates a CloudTrail Lake SQL query from a natural language
// prompt using keyword-based template matching.
func (s *CloudTrailService) GenerateQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// Parse EventDataStores (list of EDS ARNs or IDs); the Core normalises
	// each entry to the storage key.
	edsListRaw := req.Parameters["EventDataStores"]
	var edsIDs []string
	if arr, ok := edsListRaw.([]interface{}); ok {
		for _, item := range arr {
			if str, ok := item.(string); ok {
				edsIDs = append(edsIDs, str)
			}
		}
	}

	return s.generateQueryCore(store, GenerateQueryInput{
		EventDataStores: edsIDs,
		Prompt:          strings.ToLower(request.GetStringParam(req.Parameters, "Prompt")),
	})
}

// SearchSampleQueries returns sample CloudTrail Lake queries matching the
// search phrase.
func (s *CloudTrailService) SearchSampleQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.searchSampleQueriesCore(SearchSampleQueriesInput{
		SearchPhrase: request.GetStringParam(req.Parameters, "SearchPhrase"),
		MaxResults:   request.GetIntParam(req.Parameters, "MaxResults"),
		NextToken:    request.GetStringParam(req.Parameters, "NextToken"),
	})
}

// computeRelevance scores a sample query against the search phrase.
// Higher scores indicate better matches.
func computeRelevance(searchPhrase string, sq sampleQuery) int {
	if searchPhrase == "" {
		return 1
	}
	name := strings.ToLower(sq.Name)
	desc := strings.ToLower(sq.Description)

	score := 0
	if strings.Contains(name, searchPhrase) {
		score += 10
	}
	if strings.Contains(desc, searchPhrase) {
		score += 5
	}
	// Word-level matching.
	words := strings.Fields(searchPhrase)
	for _, word := range words {
		if len(word) < 3 {
			continue
		}
		if strings.Contains(name, word) {
			score += 3
		}
		if strings.Contains(desc, word) {
			score += 2
		}
	}
	return score
}
