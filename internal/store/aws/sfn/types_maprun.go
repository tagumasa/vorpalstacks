package sfn

// This file holds the Distributed Map run records and their counters.

// MapRunItemCounts holds the per-status count of items processed by a
// distributed map run. Matches Smithy MapRunItemCounts shape.
type MapRunItemCounts struct {
	Pending               int64 `json:"pending"`
	Running               int64 `json:"running"`
	Succeeded             int64 `json:"succeeded"`
	Failed                int64 `json:"failed"`
	TimedOut              int64 `json:"timedOut"`
	Aborted               int64 `json:"aborted"`
	Total                 int64 `json:"total"`
	ResultsWritten        int64 `json:"resultsWritten"`
	FailuresNotRedrivable int64 `json:"failuresNotRedrivable,omitempty"`
	PendingRedrive        int64 `json:"pendingRedrive,omitempty"`
}

// MapRunExecutionCounts holds the per-status count of child workflow
// executions started by a distributed map run. Matches Smithy
// MapRunExecutionCounts shape.
type MapRunExecutionCounts struct {
	Pending               int64 `json:"pending"`
	Running               int64 `json:"running"`
	Succeeded             int64 `json:"succeeded"`
	Failed                int64 `json:"failed"`
	TimedOut              int64 `json:"timedOut"`
	Aborted               int64 `json:"aborted"`
	Total                 int64 `json:"total"`
	ResultsWritten        int64 `json:"resultsWritten"`
	FailuresNotRedrivable int64 `json:"failuresNotRedrivable,omitempty"`
	PendingRedrive        int64 `json:"pendingRedrive,omitempty"`
}

// MapRun represents a distributed map state execution within a Step Functions
// state machine. Map runs are persisted in Pebble so they survive restarts.
type MapRun struct {
	MapRunArn                  string                `json:"mapRunArn"`
	ExecutionArn               string                `json:"executionArn"`
	StateMachineArn            string                `json:"stateMachineArn"`
	Name                       string                `json:"name"`
	Status                     string                `json:"status"`
	StartDate                  int64                 `json:"startDate"`
	StopDate                   int64                 `json:"stopDate,omitempty"`
	ItemCounts                 MapRunItemCounts      `json:"itemCounts"`
	ExecutionCounts            MapRunExecutionCounts `json:"executionCounts"`
	MaxConcurrency             int64                 `json:"maxConcurrency"`
	ToleratedFailureCount      int64                 `json:"toleratedFailureCount,omitempty"`
	ToleratedFailurePercentage float32               `json:"toleratedFailurePercentage,omitempty"`
	RedriveCount               int64                 `json:"redriveCount,omitempty"`
	RedriveDate                int64                 `json:"redriveDate,omitempty"`
	// Inline marks a record kept purely as the redrive checkpoint of an
	// inline Map state. "When you run a Map state in Distributed mode,
	// Step Functions creates a Map Run resource" — an inline run is engine
	// bookkeeping and stays invisible on the Map Run API surface.
	Inline bool `json:"inline,omitempty"`
	// Attempt records the state-level retry attempt this run belongs to
	// ("When you retry a Map state, it creates a new Map Run"): every run
	// dispatches its own child workflow executions, so the attempt stamps
	// the child-execution names and keeps each run's child set disjoint
	// from its predecessors'. A redrive keeps the run — and therefore the
	// attempt — it reclaims.
	Attempt          int64          `json:"attempt,omitempty"`
	CompletedResults map[int]string `json:"completedResults,omitempty"`
}

// MapRunListResult holds a paginated list of map runs.
type MapRunListResult struct {
	MapRuns   []*MapRun
	NextToken string
}
