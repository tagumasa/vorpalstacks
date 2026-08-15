package iot

// AWS specification limits for the IoT fleet indexing query surface. These
// exported constants are the single definition of each value: service
// handlers and tests must reference them rather than restate the number.
const (
	// SearchIndexMaxResultsMin is the inclusive lower bound of the
	// SearchIndex maxResults parameter. The Smithy model attaches
	// range(min=1) to SearchQueryMaxResults.
	SearchIndexMaxResultsMin = 1

	// SearchIndexMaxResultsCap is the inclusive upper bound of the
	// SearchIndex maxResults parameter. The Smithy model documents
	// SearchQueryMaxResults as "The maximum number of results to return
	// per page at one time. This maximum number cannot exceed 100."
	SearchIndexMaxResultsCap = 100
)
