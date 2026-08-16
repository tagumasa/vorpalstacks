package cloudwatchlogs

import (
	"encoding/base64"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"
)

// queryResultRow represents a single row in the query results. columns
// records the order in which field names were first written so that output
// field order and "first field" selection stay deterministic.
type queryResultRow struct {
	fields  map[string]string
	columns []string
}

// set stores a field value, appending new field names to the column order.
func (r *queryResultRow) set(k, v string) {
	if r.fields == nil {
		r.fields = make(map[string]string)
	}
	if _, ok := r.fields[k]; !ok {
		r.columns = append(r.columns, k)
	}
	r.fields[k] = v
}

// ordered returns the row's field names in insertion order. Fields written
// directly to the map without set are appended in sorted order as a
// fallback so the result is always deterministic.
func (r *queryResultRow) ordered() []string {
	out := make([]string, 0, len(r.fields))
	seen := make(map[string]bool, len(r.fields))
	for _, k := range r.columns {
		if _, ok := r.fields[k]; !ok || seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, k)
	}
	var rest []string
	for k := range r.fields {
		if !seen[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	return append(out, rest...)
}

// cloneRow deep-copies a row including its column order.
func cloneRow(src queryResultRow) queryResultRow {
	fields := make(map[string]string, len(src.fields))
	for k, v := range src.fields {
		fields[k] = v
	}
	columns := make([]string, len(src.columns))
	copy(columns, src.columns)
	return queryResultRow{fields: fields, columns: columns}
}

// queryResultField represents a field-value pair in the output.
type queryResultField struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

// queryStats holds statistics about the query execution.
type queryStats struct {
	recordsScanned int64
	recordsMatched int64
	bytesScanned   int64
}

// logEventWithContext carries a log event with its group/stream context.
type logEventWithContext struct {
	timestamp     int64
	message       string
	ingestionTime int64
	logGroup      string
	logStream     string
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

	subqueryCache map[string][]interface{}
	lookupCache   map[string]*parsedLookupTable
	sourceError   error
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
func (ctx *execContext) runSubquery(toks []token) ([]queryResultRow, error) {
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
	if len(cmds) > 0 {
		if _, ok := cmds[0].cmd.(*sourceCommand); ok {
			// The SOURCE command refetches events into the child context.
		} else {
			child.effectiveGroups = ctx.effectiveGroups
			child.events = ctx.events
		}
	}
	rows := buildRows(child.events, child.accountID)
	for _, c := range cmds {
		rows = c.cmd.apply(child, rows)
		child.preceding = append(child.preceding, c.cmd)
	}
	if child.sourceError != nil {
		return nil, child.sourceError
	}
	return rows, nil
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
func validateQueryPipeline(queryString string) error {
	toks, err := lexQuery(queryString)
	if err != nil {
		return err
	}
	cmds, err := compilePipeline(toks)
	if err != nil {
		return err
	}
	if len(cmds) == 0 {
		return nil
	}
	return validateCommandOrder(cmds)
}

// validateCommandOrder enforces the documented command placement rules.
func validateCommandOrder(cmds []compiledCommand) error {
	lastStats := -1
	var statsHeads []token
	var joinHeads []token
	patternIdx := -1
	sortIdx := -1
	dedupIdx := -1
	for i, c := range cmds {
		switch c.cmd.name() {
		case "stats":
			statsHeads = append(statsHeads, c.head)
			lastStats = i
		case "join":
			joinHeads = append(joinHeads, c.head)
		case "pattern":
			if patternIdx < 0 {
				patternIdx = i
			}
		case "sort":
			sortIdx = i
		case "dedup":
			if dedupIdx < 0 {
				dedupIdx = i
			}
		}
	}
	if len(statsHeads) > 10 {
		// The error points at the first command beyond the limit.
		return newQueryCompileError("A query can have a maximum of 10 stats commands",
			statsHeads[10].start, statsHeads[10].end)
	}
	if len(joinHeads) > 1 {
		// The error points at the second join, which is the violation.
		return newQueryCompileError("Only one join command is supported per query",
			joinHeads[1].start, joinHeads[1].end)
	}
	if patternIdx >= 0 && sortIdx >= 0 && sortIdx < patternIdx {
		return newQueryCompileError(
			"A query is not valid if it includes a pattern command after a sort command", cmds[patternIdx].head.start, cmds[patternIdx].head.end)
	}
	// sort and limit must appear after the last stats command.
	for i, c := range cmds {
		n := c.cmd.name()
		if (n == "sort" || n == "limit") && lastStats >= 0 && i < lastStats {
			return newQueryCompileError(
				"If you use a sort or limit command, it must appear after the last stats command", c.head.start, c.head.end)
		}
	}
	// Only limit may follow dedup.
	if dedupIdx >= 0 {
		for i := dedupIdx + 1; i < len(cmds); i++ {
			if cmds[i].cmd.name() != "limit" {
				return newQueryCompileError(
					"The only query command that you can use after the dedup command is limit", cmds[i].head.start, cmds[i].head.end)
			}
		}
	}
	// SOURCE is only valid as the first command.
	for i, c := range cmds {
		if c.cmd.name() == "SOURCE" && i > 0 {
			return newQueryCompileError("SOURCE is only valid as the first command of a query", c.head.start, c.head.end)
		}
	}
	return nil
}

// buildRows converts raw log events into the initial row set with the
// discoverable fields. JSON log messages contribute their top-level keys as
// discovered fields, matching the automatic field discovery of Logs
// Insights; nested structures are kept as canonical JSON strings.
// maxDiscoveredJSONFields is the documented ceiling on the number of fields
// Logs Insights extracts from a JSON log event; further fields are ignored
// and must be extracted with parse.
const maxDiscoveredJSONFields = 200

func buildRows(events []logEventWithContext, accountID string) []queryResultRow {
	rows := make([]queryResultRow, 0, len(events))
	for _, evt := range events {
		row := queryResultRow{}
		row.set("@timestamp", formatTimestampField(evt.timestamp))
		row.set("@message", evt.message)
		row.set("@logStream", evt.logStream)
		// @log identifies the event's log group as account:group-name.
		row.set("@log", accountID+":"+evt.logGroup)
		row.set("@ingestionTime", formatTimestampField(evt.ingestionTime))
		// @ptr addresses the event for GetLogRecord; aggregate rows built
		// by later commands never carry it.
		row.set("@ptr", eventPointer(evt.logGroup, evt.logStream, evt.timestamp, evt.message))
		discoverJSONFields(evt.message, &row)
		rows = append(rows, row)
	}
	return rows
}

// eventPointer encodes the logGroup|logStream|timestamp|message pointer
// that GetLogRecord accepts, in base64.
func eventPointer(group, stream string, ts int64, msg string) string {
	return base64.StdEncoding.EncodeToString([]byte(group + "|" + stream + "|" + strconv.FormatInt(ts, 10) + "|" + msg))
}

// discoverJSONFields merges the top-level keys of a JSON object message
// into the row fields. Discoverable @-fields are never overwritten; field
// names starting with @ are displayed with an additional @ prefix, per the
// documented discovery behaviour.
func discoverJSONFields(message string, row *queryResultRow) {
	trimmed := strings.TrimSpace(message)
	if !strings.HasPrefix(trimmed, "{") {
		return
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &decoded); err != nil {
		return
	}
	keys := make([]string, 0, len(decoded))
	for k := range decoded {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	discovered := 0
	for _, k := range keys {
		name := k
		if strings.HasPrefix(k, "@") {
			name = "@" + k
		}
		if _, exists := row.fields[name]; exists {
			continue
		}
		if discovered >= maxDiscoveredJSONFields {
			return
		}
		row.set(name, storeValue(decoded[k]))
		discovered++
	}
}

// formatTimestampField renders epoch milliseconds; numeric formatting keeps
// the value numeric for later comparisons.
func formatTimestampField(ms int64) string {
	return formatNumber(float64(ms))
}

// executeQuery runs a CloudWatch Logs Insights query against the given
// events. Commands follow the documented Logs Insights QL grammar.
func executeQuery(queryString string, events []logEventWithContext) ([]queryResultRow, queryStats) {
	ctx := &execContext{events: events, subqueryCache: map[string][]interface{}{}}
	rows, _ := executeQueryContext(ctx, queryString)
	stats := queryStats{
		recordsScanned: int64(len(events)),
	}
	for _, e := range events {
		stats.bytesScanned += int64(len(e.message))
	}
	stats.recordsMatched = int64(len(rows))
	return rows, stats
}

// executeQueryContext compiles and runs the query within an execution
// context. It returns the result rows and any compile error.
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

	rows := buildRows(ctx.events, ctx.accountID)
	for _, c := range cmds {
		rows = c.cmd.apply(ctx, rows)
		ctx.preceding = append(ctx.preceding, c.cmd)
	}
	if ctx.sourceError != nil {
		// SOURCE resolution and refetch failures fail the query rather
		// than silently returning the default groups' rows.
		return nil, ctx.sourceError
	}
	return rows, nil
}
