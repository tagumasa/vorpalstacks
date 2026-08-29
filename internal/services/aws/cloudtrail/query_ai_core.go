package cloudtrail

import (
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
)

// GenerateQueryInput carries the natural-language prompt and the target
// event data store IDs for GenerateQuery.
type GenerateQueryInput struct {
	EventDataStores []string
	Prompt          string
}

// generateQueryCore is the single entry point for GenerateQuery: it verifies
// the EventDataStores list and Prompt, validates that the first event data
// store exists, and maps the prompt onto a CloudTrail Lake SQL statement via
// keyword-based template matching.
func (s *CloudTrailService) generateQueryCore(store StoreInterface, in GenerateQueryInput) (map[string]interface{}, error) {
	edsIDs := in.EventDataStores
	if len(edsIDs) == 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"EventDataStores is required", 400)
	}

	// Validate EDS exists.
	_, err := store.GetEventDataStore(edsIDs[0])
	if err != nil {
		return nil, awserrors.NewAWSError("EventDataStoreNotFoundException",
			"Event data store not found", 404)
	}

	prompt := in.Prompt
	if prompt == "" {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"Prompt is required", 400)
	}

	// Generate SQL via keyword-based template matching.
	sql := generateQueryFromPrompt(prompt, edsIDs[0])
	alias := generateAlias(prompt)

	return map[string]interface{}{
		"QueryStatement":               sql,
		"QueryAlias":                   alias,
		"EventDataStoreOwnerAccountId": store.GetAccountID(),
	}, nil
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
