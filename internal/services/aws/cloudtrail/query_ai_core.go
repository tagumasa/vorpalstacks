package cloudtrail

import (
	"fmt"
	"sort"
	"strings"

	"vorpalstacks/internal/common/pagination"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// GenerateQueryInput carries the natural-language prompt and the target
// event data stores for GenerateQuery. EventDataStores holds the wire values
// (IDs or ARNs); the Core normalises them to storage keys.
type GenerateQueryInput struct {
	EventDataStores []string
	Prompt          string
}

// SearchSampleQueriesInput carries the search phrase and paging members
// for SearchSampleQueries.
type SearchSampleQueriesInput struct {
	SearchPhrase string
	MaxResults   int
	NextToken    string
}

// generateQueryCore is the single entry point for GenerateQuery: it verifies
// the EventDataStores list and Prompt, validates that the first event data
// store exists, and maps the prompt onto a CloudTrail Lake SQL statement via
// keyword-based template matching.
func (s *CloudTrailService) generateQueryCore(store cloudtrailstore.CloudTrailStoreInterface, in GenerateQueryInput) (map[string]interface{}, error) {
	// The model bounds EventDataStoreList to exactly one entry and marks it
	// and Prompt required; Prompt carries the length 3-500 and printable
	// pattern constraints.
	if len(in.EventDataStores) == 0 {
		return nil, newInvalidParameterException(
			"EventDataStores is required")
	}
	if len(in.EventDataStores) > 1 {
		return nil, newInvalidParameterException(
			"EventDataStores must contain exactly one event data store")
	}
	if len(in.Prompt) < cloudtrailstore.MinPromptLength || len(in.Prompt) > cloudtrailstore.MaxPromptLength || !promptPattern.MatchString(in.Prompt) {
		return nil, newInvalidParameterException(
			fmt.Sprintf("Prompt must be between %d and %d printable characters",
				cloudtrailstore.MinPromptLength, cloudtrailstore.MaxPromptLength))
	}

	edsID := cloudtrailstore.ExtractEventDataStoreID(in.EventDataStores[0])

	// Validate EDS exists.
	_, err := s.resolveEventDataStore(store, in.EventDataStores[0])
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// Generate SQL via keyword-based template matching.
	sql := generateQueryFromPrompt(in.Prompt, edsID)
	alias := generateAlias(in.Prompt)

	return map[string]interface{}{
		"QueryStatement":               sql,
		"QueryAlias":                   alias,
		"EventDataStoreOwnerAccountId": store.GetAccountID(),
	}, nil
}

// searchSampleQueriesCore is the single entry point for SearchSampleQueries:
// it validates the required search phrase and the MaxResults bound, scores
// the curated sample list against the phrase, orders the matches by
// descending relevance, and pages them with the shared positional
// paginator — the service's one offset-token policy.
func (s *CloudTrailService) searchSampleQueriesCore(in SearchSampleQueriesInput) (map[string]interface{}, error) {
	// SearchPhrase is model-required with length 2-1000 and the printable
	// pattern; MaxResults carries the model range 1-50.
	if len(in.SearchPhrase) < cloudtrailstore.MinSearchPhraseLength || len(in.SearchPhrase) > cloudtrailstore.MaxSearchPhraseLength || !promptPattern.MatchString(in.SearchPhrase) {
		return nil, newInvalidParameterException(
			fmt.Sprintf("SearchPhrase must be between %d and %d printable characters",
				cloudtrailstore.MinSearchPhraseLength, cloudtrailstore.MaxSearchPhraseLength))
	}
	if in.MaxResults < 0 || in.MaxResults > cloudtrailstore.MaxSampleQueriesResults {
		return nil, newInvalidParameterException(
			fmt.Sprintf("MaxResults must be between 1 and %d", cloudtrailstore.MaxSampleQueriesResults))
	}
	maxResults := in.MaxResults
	if maxResults == 0 {
		maxResults = cloudtrailstore.DefaultSampleQueriesResults
	}

	phrase := strings.ToLower(in.SearchPhrase)

	// Filter sample queries by search phrase. computeRelevance returns 1
	// for all queries when the phrase is empty, so no special-casing needed.
	matched := make([]map[string]interface{}, 0, len(cloudTrailSampleQueries))
	for _, sq := range cloudTrailSampleQueries {
		relevance := computeRelevance(phrase, sq)
		if relevance > 0 {
			matched = append(matched, map[string]interface{}{
				"Name":        sq.Name,
				"Description": sq.Description,
				"SQL":         sq.SQL,
				"Relevance":   relevance,
			})
		}
	}

	// Sort by relevance descending (highest relevance first).
	sort.Slice(matched, func(i, j int) bool {
		ri, _ := matched[i]["Relevance"].(int)
		rj, _ := matched[j]["Relevance"].(int)
		return ri > rj
	})

	paged := pagination.PaginateSliceByPosition(matched, in.NextToken, maxResults)
	resp := map[string]interface{}{
		"SearchResults": paged.Items,
	}
	if paged.IsTruncated {
		resp["NextToken"] = paged.NextMarker
	}
	return resp, nil
}

// generateQueryFromPrompt maps a natural language prompt to a CloudTrail Lake
// SQL query using keyword matching against the sample query templates.
func generateQueryFromPrompt(prompt, edsID string) string {
	for _, sq := range cloudTrailSampleQueries {
		name := strings.ToLower(sq.Name)
		desc := strings.ToLower(sq.Description)
		// Check if any keyword from the sample name or description appears in the prompt.
		keywords := strings.Fields(name + " " + desc)
		for _, kw := range keywords {
			if len(kw) < 4 {
				continue
			}
			if strings.Contains(prompt, kw) {
				return formatSampleSQL(sq.SQL, edsID)
			}
		}
	}

	// Fallback: default management events query.
	return formatSampleSQL(cloudTrailSampleQueries[0].SQL, edsID)
}

// generateAlias creates a short alias from the prompt.
func generateAlias(prompt string) string {
	words := strings.Fields(prompt)
	if len(words) > 3 {
		words = words[:3]
	}
	return strings.Join(words, "_")
}

// formatSampleSQL fills in the EDS ID placeholder in a sample SQL template.
func formatSampleSQL(template, edsID string) string {
	return strings.ReplaceAll(template, "%s", edsID)
}
