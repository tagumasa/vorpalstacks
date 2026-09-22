package cloudwatchlogs

import (
	"fmt"
	"time"
)

// The query execution engine's core: the execution context (the
// per-query state the pipeline reads), the command pipeline's compile
// and apply loop, and the compile-and-run entry the operations call.

type queryStats struct {
	recordsScanned int64
	recordsMatched int64
	bytesScanned   int64
}

// logEventWithContext carries a log event with its group/stream context.
// transformed marks the message as a transformed copy (the query plane's
// read of the ingestion-time transformation output) — the transformer's
// @transformationError marker then surfaces as the single-@ system field
// the documented query filters on.
type logEventWithContext struct {
	timestamp     int64
	message       string
	ingestionTime int64
	logGroup      string
	logStream     string
	transformed   bool
}

// sourceGroupInfo is the log group metadata SOURCE selection needs.
type sourceGroupInfo struct {
	Name  string
	Class string
	Tags  map[string]string
}

// binFillInfo records the bin grouping of the last stats command so that
// fillmissing can synthesise rows for empty time bins.
type binFillInfo struct {
	dur    int64
	minBin int64
	maxBin int64
}

// execContext carries everything a pipeline execution needs beyond the
// rows: the query window, the underlying events, log group resolution for
// SOURCE/join/subqueries, and per-run scratch state.
type execContext struct {
	startTime int64
	endTime   int64
	accountID string

	events          []logEventWithContext
	effectiveGroups []string
	defaultGroups   []string

	fetchEvents    func(groups []string, start, end int64) ([]logEventWithContext, error)
	listLogGroups  func() ([]sourceGroupInfo, error)
	getLookupTable func(name string) (*parsedLookupTable, error)

	preceding  []command
	sorted     bool
	statsCount int
	lastBins   *binFillInfo
	// currentBinDur is the bin duration of the stats command currently
	// emitting, which the time-series functions scale their window by.
	currentBinDur int64
	// recordsMatched counts the rows the event-plane retained — the
	// documented QueryStatistics basis ("The number of log events that
	// matched the query string"), recorded before aggregation collapses
	// rows or limit truncates them. Subqueries execute on child contexts,
	// so the outer query's count is never clobbered.
	recordsMatched int64

	subqueryCache map[string][]interface{}
	lookupCache   map[string]*parsedLookupTable
	sourceError   error

	// deadline bounds the pipeline's execution — the sixty-minute outer
	// runtime ("Queries time out after 60 minutes of runtime") or the
	// thirty-second inner limit ("Inner query execution is limited to
	// 30 seconds"); zero means unbounded. timedOut marks the breach: the
	// outer lifecycle stamps the Timeout status from it, while the inner
	// limit surfaces as a query failure carrying deadlineMsg's wording.
	deadline    time.Time
	deadlineMsg string
	timedOut    bool

	// startedMs stamps the query processing's start — the fixed point
	// now() reports ("Returns the time that the query processing was
	// started, in epoch seconds"), so the value does not drift as a
	// long-running pipeline evaluates. Zero means unstamped (a test
	// context); now() then falls back to the evaluation clock.
	startedMs int64

	// subqueryDepth counts nested query executions: the outer query runs
	// at zero and each runSubquery call executes one level deeper.
	subqueryDepth int
}

func (ctx *execContext) now() int64 {
	return time.Now().UnixMilli()
}

// runPrecedingOnWindow re-runs the commands preceding diff/logcompare over
// a different time window so that the comparison covers the same analysis.
func (ctx *execContext) runPrecedingOnWindow(start, end int64) ([]queryResultRow, error) {
	groups := ctx.effectiveGroups
	if groups == nil {
		groups = ctx.defaultGroups
	}
	events, err := ctx.fetchEvents(groups, start, end)
	if err != nil {
		return nil, err
	}
	rows := buildRows(events, ctx.accountID)
	for _, cmd := range ctx.preceding {
		rows = cmd.apply(ctx, rows)
	}
	return rows, nil
}

// runSubquery executes a nested query over its own SOURCE selection, or the
// enclosing query's log groups when the subquery has no SOURCE.
// subqueryExecutionLimit is the documented inner bound: "Inner query
// execution is limited to 30 seconds."
const subqueryExecutionLimit = 30 * time.Second

// queryDeadlineNow is the clock the deadline sites read — the seam a test
// uses to advance an execution past its bound without waiting out the
// sixty-minute outer runtime.
var queryDeadlineNow = time.Now

func (ctx *execContext) runSubquery(toks []token) ([]queryResultRow, error) {
	// "Nested subqueries are not supported." — a subquery executing
	// inside another subquery's context rejects; the filter-in nesting
	// form is additionally rejected at compile.
	if ctx.subqueryDepth > 0 {
		return nil, fmt.Errorf("nested subqueries are not supported")
	}
	cmds, err := compilePipeline(toks)
	if err != nil {
		return nil, err
	}
	child := &execContext{
		startTime:       ctx.startTime,
		endTime:         ctx.endTime,
		accountID:       ctx.accountID,
		defaultGroups:   ctx.defaultGroups,
		effectiveGroups: ctx.defaultGroups,
		events:          ctx.events,
		fetchEvents:     ctx.fetchEvents,
		listLogGroups:   ctx.listLogGroups,
		getLookupTable:  ctx.getLookupTable,
		subqueryCache:   ctx.subqueryCache,
		lookupCache:     ctx.lookupCache,
	}
	child.deadline = queryDeadlineNow().Add(subqueryExecutionLimit)
	child.deadlineMsg = "Inner query execution is limited to 30 seconds"
	child.subqueryDepth = ctx.subqueryDepth + 1
	if len(cmds) > 0 {
		if _, ok := cmds[0].cmd.(*sourceCommand); ok {
			// The SOURCE command refetches events into the child context.
		} else {
			child.effectiveGroups = ctx.effectiveGroups
			child.events = ctx.events
		}
	}
	rows := applyPipelineCommands(child, cmds)
	// The thirty-second inner limit is a query failure, unlike the
	// sixty-minute outer bound: a subquery that breached its deadline
	// fails the enclosing query rather than appending a silently
	// truncated row set.
	if child.timedOut {
		return nil, fmt.Errorf("%s", child.deadlineMsg)
	}
	if child.sourceError != nil {
		return nil, child.sourceError
	}
	return rows, nil
}

// applyPipelineCommands builds the source rows and applies the compiled
// commands in order — the one execution loop the outer query and every
// subquery ride. It records the matched-records count at the boundary
// where rows stop corresponding to source events.
func applyPipelineCommands(ctx *execContext, cmds []compiledCommand) []queryResultRow {
	rows := buildRows(ctx.events, ctx.accountID)
	counted := false
	for _, c := range cmds {
		// A breached deadline marks the context and stops the pipeline
		// without seeding a failure: the outer bound's Timeout status is
		// stamped from timedOut by the lifecycle, and only the inner
		// bound (surfaced by runSubquery) fails the query.
		if !ctx.deadline.IsZero() && queryDeadlineNow().After(ctx.deadline) {
			ctx.timedOut = true
			return rows
		}
		name := c.cmd.name()
		// "These values reflect the full raw results of the query"
		// (QueryStatistics, GetQueryResults): the count is taken as the
		// first aggregating command is about to consume the event rows,
		// or before limit truncates them, so a million matched events
		// aggregating to one row still report the million and a limit
		// never shrinks the count.
		if !counted && (aggregatesRows(name) || name == "limit") {
			ctx.recordsMatched = int64(len(rows))
			counted = true
		}
		rows = c.cmd.apply(ctx, rows)
		ctx.preceding = append(ctx.preceding, c.cmd)
	}
	if !counted {
		ctx.recordsMatched = int64(len(rows))
	}
	return rows
}

// aggregatesRows reports whether a command consumes event rows into an
// aggregate shape, after which the row count no longer counts matched
// events.
func aggregatesRows(name string) bool {
	switch name {
	case "stats", "pattern", "diff", "logcompare", "relevantfields",
		"anomaly", "join", "countFrequent", "estimate":
		return true
	}
	return false
}

// compiledCommand pairs a command with its head token for validation
// reporting.
type compiledCommand struct {
	cmd  command
	head token
}

// compilePipeline lexes, splits, and parses a full query into commands.
func compilePipeline(toks []token) ([]compiledCommand, error) {
	segs, err := parsePipelineCommands(toks)
	if err != nil {
		return nil, err
	}
	var out []compiledCommand
	for _, seg := range segs {
		if len(seg) == 0 {
			continue
		}
		cmd, err := parseCommand(seg)
		if err != nil {
			return nil, err
		}
		if cmd == nil {
			continue
		}
		out = append(out, compiledCommand{cmd: cmd, head: seg[0]})
	}
	return out, nil
}

// validateQueryPipeline compiles the query and checks the documented
// structural rules. AWS rejects queries that fail to compile with
// MalformedQueryException carrying a QueryCompileError with character
// offsets; without this check an unknown command would be silently ignored
// at execution time and the query would report success over unintended
// results.
func executeQueryContext(ctx *execContext, queryString string) ([]queryResultRow, error) {
	toks, err := lexQuery(queryString)
	if err != nil {
		return nil, err
	}
	cmds, err := compilePipeline(toks)
	if err != nil {
		return nil, err
	}
	if err := validateCommandOrder(cmds); err != nil {
		return nil, err
	}
	ctx.effectiveGroups = ctx.defaultGroups

	rows := applyPipelineCommands(ctx, cmds)
	if ctx.sourceError != nil {
		// SOURCE resolution and refetch failures fail the query rather
		// than silently returning the default groups' rows.
		return nil, ctx.sourceError
	}
	return rows, nil
}
