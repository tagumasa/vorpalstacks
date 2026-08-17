package cloudwatchlogs

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Pattern-analysis commands: countFrequent, pattern, diff, logcompare,
// relevantfields, and anomaly, with the shared pattern masking machinery.

// --- countFrequent ---

type countFrequentCommand struct {
	fields []string
	head   token
}

func (c *countFrequentCommand) name() string { return "countFrequent" }

func parseCountFrequentCommand(args []token, head token) (command, error) {
	if len(args) == 0 {
		return nil, newQueryCompileError("Syntax error: countFrequent requires fields", head.start, head.end)
	}
	c := &countFrequentCommand{head: head}
	for i := 0; i < len(args); i++ {
		if args[i].kind != tokIdent && args[i].kind != tokBacktickIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", args[i].raw), args[i].start, args[i].end)
		}
		c.fields = append(c.fields, args[i].text)
		if i+1 < len(args) && args[i+1].kind == tokComma {
			i++
		}
	}
	return c, nil
}

func (c *countFrequentCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	counts := make(map[string]int)
	parts := make(map[string][]string)
	var order []string
	for i := range rows {
		vals := make([]string, len(c.fields))
		for fi, f := range c.fields {
			vals[fi] = readField(rows[i], f, ctx)
		}
		key := strings.Join(vals, "\x00")
		if _, ok := counts[key]; !ok {
			order = append(order, key)
			parts[key] = vals
		}
		counts[key]++
	}
	sort.Slice(order, func(i, j int) bool {
		if counts[order[i]] != counts[order[j]] {
			return counts[order[i]] > counts[order[j]]
		}
		return order[i] < order[j]
	})
	out := make([]queryResultRow, 0, len(order))
	for _, key := range order {
		counted := queryResultRow{}
		for fi, f := range c.fields {
			counted.set(f, parts[key][fi])
		}
		counted.set("_approxcount", strconv.Itoa(counts[key]))
		out = append(out, counted)
	}
	return out
}

// --- pattern ---

type patternCommand struct {
	expr   exprNode
	headTk token
}

func (c *patternCommand) name() string { return "pattern" }

func parsePatternCommand(args []token, head token) (command, error) {
	if len(args) == 0 {
		return nil, newQueryCompileError("Syntax error: pattern requires a field", head.start, head.end)
	}
	e, err := parseExprTokens(args)
	if err != nil {
		return nil, err
	}
	return &patternCommand{expr: e, headTk: head}, nil
}

func (c *patternCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	type patGroup struct {
		count    int
		samples  []string
		pattern  string
		order    int
		severity string
	}
	groups := make(map[string]*patGroup)
	var order []string
	total := 0
	for i := range rows {
		row := rows[i]
		v := asString(c.expr.eval(&row, ctx))
		if v == "" {
			continue
		}
		pat := maskPattern(v)
		total++
		g, ok := groups[pat]
		if !ok {
			g = &patGroup{pattern: pat, order: len(order)}
			groups[pat] = g
			order = append(order, pat)
		}
		g.count++
		if len(g.samples) < 10 {
			g.samples = append(g.samples, v)
		}
	}
	out := make([]queryResultRow, 0, len(order))
	for _, pat := range order {
		g := groups[pat]
		patternRow := queryResultRow{}
		patternRow.set("@pattern", g.pattern)
		patternRow.set("@sampleCount", strconv.Itoa(g.count))
		patternRow.set("@ratio", strconv.FormatFloat(float64(g.count)/float64(total), 'f', 2, 64))
		patternRow.set("@severityLabel", patternSeverity(g.samples))
		out = append(out, patternRow)
	}
	sort.SliceStable(out, func(i, j int) bool {
		ci, _ := asNumber(out[i].fields["@sampleCount"])
		cj, _ := asNumber(out[j].fields["@sampleCount"])
		return ci > cj
	})
	return out
}

// patternSeverity classifies a pattern's severity from its samples: ERROR
// when the word Error appears, WARN when Warn appears without Error, INFO
// otherwise.
func patternSeverity(samples []string) string {
	hasErr, hasWarn := false, false
	for _, s := range samples {
		if strings.Contains(s, "Error") || strings.Contains(s, "ERROR") {
			hasErr = true
		}
		if strings.Contains(s, "Warn") || strings.Contains(s, "WARN") {
			hasWarn = true
		}
	}
	switch {
	case hasErr:
		return "ERROR"
	case hasWarn:
		return "WARN"
	}
	return "INFO"
}

var (
	reISOTime   = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}([T ]\d{2}:\d{2}(:\d{2})?)?`)
	reTimeOnly  = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}`)
	reEpoch     = regexp.MustCompile(`^\d{10,13}(\.\d+)?$`)
	reUUID      = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	reIPAddr    = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	reHex       = regexp.MustCompile(`^[0-9a-fA-F]{6,}$`)
	reNumberTok = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	reIdentTok  = regexp.MustCompile(`^[0-9A-Za-z][0-9A-Za-z_-]{2,}$`)
	reHasDigit  = regexp.MustCompile(`\d`)
)

// maskPattern converts a log message into its pattern by replacing dynamic
// tokens with typed placeholders such as <Time-1> or <ID-2>. The number is
// the ordinal position of the token within the pattern.
func maskPattern(msg string) string {
	fields := strings.Fields(msg)
	counts := make(map[string]int)
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		kind := tokenKindOf(f)
		if kind == "" {
			out = append(out, f)
			continue
		}
		counts[kind]++
		out = append(out, fmt.Sprintf("<%s-%d>", kind, counts[kind]))
	}
	return strings.Join(out, " ")
}

// tokenKindOf classifies a dynamic token: timestamps, identifiers, numbers,
// IP addresses, and hexadecimal identifiers; unknown dynamic tokens fall
// back to Token.
func tokenKindOf(f string) string {
	trimmed := strings.Trim(f, ",;:()[]{}\"'")
	switch {
	case trimmed == "":
		return ""
	case reISOTime.MatchString(trimmed) || reTimeOnly.MatchString(trimmed) || reEpoch.MatchString(trimmed):
		return "Time"
	case reUUID.MatchString(trimmed):
		return "ID"
	case reIPAddr.MatchString(trimmed):
		return "IP"
	case reNumberTok.MatchString(trimmed):
		return "Num"
	case reHex.MatchString(trimmed) && len(trimmed)%2 == 0:
		return "ID"
	}
	if reIdentTok.MatchString(trimmed) {
		// Alphanumeric identifiers containing digits are dynamic tokens;
		// plain dictionary-like words that recur verbatim stay literal so
		// that stable structure words keep the pattern readable.
		if reHasDigit.MatchString(trimmed) {
			return "ID"
		}
		return ""
	}
	if len(trimmed) <= 2 {
		return ""
	}
	return "Token"
}

// --- diff ---

type diffCommand struct {
	mode   string // "", previousDay, previousWeek, previousMonth
	headTk token
}

func (c *diffCommand) name() string { return "diff" }

func parseDiffCommand(args []token, head token) (command, error) {
	c := &diffCommand{headTk: head}
	if len(args) > 0 {
		if args[0].kind != tokIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", args[0].raw), args[0].start, args[0].end)
		}
		switch strings.ToLower(args[0].text) {
		case "previousday", "previousweek", "previousmonth":
			c.mode = strings.ToLower(args[0].text)
		default:
			return nil, newQueryCompileError(
				fmt.Sprintf("Invalid diff modifier '%s'", args[0].text), args[0].start, args[0].end)
		}
	}
	return c, nil
}

// comparisonWindow returns the comparison window for the diff mode. The
// window always has the same length as the query window.
func (c *diffCommand) comparisonWindow(start, end int64) (int64, int64) {
	span := end - start
	switch c.mode {
	case "previousday":
		return start - 24*3600*1000, end - 24*3600*1000
	case "previousweek":
		return start - 7*24*3600*1000, end - 7*24*3600*1000
	case "previousmonth":
		return start - 30*24*3600*1000, end - 30*24*3600*1000
	}
	return start - span, start
}

func (c *diffCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	ps, pe := c.comparisonWindow(ctx.startTime, ctx.endTime)
	return comparePatternWindows(ctx, rows, ps, pe)
}

// comparePatternWindows runs the comparison side shared by diff and
// logcompare. Comparison queries analyse patterns: the current rows are
// clustered, the comparison window is re-queried with the same preceding
// pipeline, and the per-pattern count differences are reported.
func comparePatternWindows(ctx *execContext, rows []queryResultRow, winStart, winEnd int64) []queryResultRow {
	curPatterns := patternCounts(rows)

	prevRows, err := ctx.runPrecedingOnWindow(winStart, winEnd)
	if err != nil {
		prevRows = nil
	}
	prevPatterns := patternCounts(prevRows)

	all := make(map[string]bool)
	for k := range curPatterns {
		all[k] = true
	}
	for k := range prevPatterns {
		all[k] = true
	}
	keys := make([]string, 0, len(all))
	for k := range all {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		di := absInt(curPatterns[keys[i]] - prevPatterns[keys[i]])
		dj := absInt(curPatterns[keys[j]] - prevPatterns[keys[j]])
		if di != dj {
			return di > dj
		}
		return keys[i] < keys[j]
	})
	out := make([]queryResultRow, 0, len(keys))
	for _, k := range keys {
		cur, prev := curPatterns[k], prevPatterns[k]
		diffRow := queryResultRow{}
		diffRow.set("@pattern", k)
		diffRow.set("@sampleCount", strconv.Itoa(cur))
		diffRow.set("@diffEventCount", strconv.Itoa(cur-prev))
		diffRow.set("@differenceDescription", diffDescription(cur, prev))
		diffRow.set("@severityLabel", patternSeverity([]string{k}))
		out = append(out, diffRow)
	}
	return out
}

func diffDescription(cur, prev int) string {
	switch {
	case prev == 0 && cur > 0:
		return "new"
	case cur == 0 && prev > 0:
		return "missing"
	case cur > prev:
		return fmt.Sprintf("%d more events", cur-prev)
	case cur < prev:
		return fmt.Sprintf("%d fewer events", prev-cur)
	}
	return "no change"
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// patternCounts clusters rows by @message pattern and counts occurrences.
func patternCounts(rows []queryResultRow) map[string]int {
	out := make(map[string]int)
	for i := range rows {
		msg := rows[i].fields["@message"]
		if msg == "" {
			continue
		}
		out[maskPattern(msg)]++
	}
	return out
}

// --- logcompare ---

type logCompareCommand struct {
	timeshift int64
	headTk    token
}

func (c *logCompareCommand) name() string { return "logcompare" }

func parseLogCompareCommand(args []token, head token) (command, error) {
	if len(args) < 2 || args[0].kind != tokIdent || !strings.EqualFold(args[0].text, "timeshift") {
		return nil, newQueryCompileError("Syntax error: logcompare requires timeshift <duration>", head.start, head.end)
	}
	text, span, ok := periodArgAt(args, 1)
	if !ok || 1+span != len(args) {
		return nil, newQueryCompileError("Syntax error: logcompare requires timeshift <duration>", head.start, head.end)
	}
	ms, ok2 := parsePeriodMillis(text)
	if !ok2 {
		return nil, newQueryCompileError(fmt.Sprintf("Invalid timeshift '%s'", text), args[1].start, args[1].end)
	}
	return &logCompareCommand{timeshift: ms, headTk: head}, nil
}

func (c *logCompareCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	return comparePatternWindows(ctx, rows, ctx.startTime-c.timeshift, ctx.endTime-c.timeshift)
}

// --- relevantfields ---

type relevantFieldsCommand struct {
	fields []string
	where  exprNode
	head   token
}

func (c *relevantFieldsCommand) name() string { return "relevantfields" }

func parseRelevantFieldsCommand(args []token, head token) (command, error) {
	c := &relevantFieldsCommand{head: head}
	whereIdx := -1
	for i, t := range args {
		if t.kind == tokIdent && strings.EqualFold(t.text, "where") {
			whereIdx = i
			break
		}
	}
	if whereIdx < 0 {
		return nil, newQueryCompileError("Syntax error: relevantfields requires where", head.start, head.end)
	}
	for i := 0; i < whereIdx; i++ {
		t := args[i]
		if t.kind == tokComma {
			continue
		}
		if t.kind != tokIdent && t.kind != tokBacktickIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
		c.fields = append(c.fields, t.text)
	}
	expr, err := parseExprTokens(args[whereIdx+1:])
	if err != nil {
		return nil, err
	}
	c.where = expr
	return c, nil
}

func (c *relevantFieldsCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	var cond, base []queryResultRow
	for i := range rows {
		row := rows[i]
		if truthy(c.where.eval(&row, ctx)) {
			cond = append(cond, row)
		} else {
			base = append(base, row)
		}
	}
	if len(cond) == 0 || len(base) == 0 {
		return []queryResultRow{}
	}
	fields := c.fields
	if len(fields) == 0 {
		seen := make(map[string]bool)
		for i := range rows {
			for k := range rows[i].fields {
				if !strings.HasPrefix(k, "@") || k == "@message" {
					seen[k] = true
				}
			}
		}
		for k := range seen {
			fields = append(fields, k)
		}
		sort.Strings(fields)
	}
	type result struct {
		field       string
		score       float64
		contributor string
		condMedian  string
		baseMedian  string
	}
	var results []result
	for _, f := range fields {
		// Categorical frequency shift between the conditional subset and
		// the baseline; numeric fields additionally compare medians.
		condCounts := fieldValueCounts(cond, f)
		baseCounts := fieldValueCounts(base, f)
		bestShift := 0.0
		bestVal := ""
		for v, cc := range condCounts {
			condRate := float64(cc) / float64(len(cond))
			baseRate := float64(baseCounts[v]) / float64(len(base))
			shift := condRate - baseRate
			if shift > bestShift {
				bestShift = shift
				bestVal = v
			}
		}
		r := result{field: f, score: bestShift, contributor: bestVal}
		if cv, ok := numericMedian(cond, f); ok {
			r.condMedian = formatNumber(cv)
		}
		if bv, ok := numericMedian(base, f); ok {
			r.baseMedian = formatNumber(bv)
		}
		results = append(results, r)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].score > results[j].score })
	out := make([]queryResultRow, 0, len(results))
	for _, r := range results {
		relevanceRow := queryResultRow{}
		relevanceRow.set("@fieldName", r.field)
		relevanceRow.set("@relevanceScore", strconv.FormatFloat(r.score, 'f', 3, 64))
		relevanceRow.set("@topRelevanceContributors", r.contributor)
		relevanceRow.set("@conditionalMedian", r.condMedian)
		relevanceRow.set("@baselineMedian", r.baseMedian)
		out = append(out, relevanceRow)
	}
	return out
}

func fieldValueCounts(rows []queryResultRow, field string) map[string]int {
	out := make(map[string]int)
	for i := range rows {
		if v, ok := rows[i].fields[field]; ok && v != "" {
			out[v]++
		}
	}
	return out
}

func numericMedian(rows []queryResultRow, field string) (float64, bool) {
	var vals []float64
	for i := range rows {
		if v, ok := asNumber(rows[i].fields[field]); ok {
			vals = append(vals, v)
		}
	}
	if len(vals) == 0 {
		return 0, false
	}
	return percentile(vals, 50), true
}

// --- anomaly ---

// anomalyCommand flags unusual patterns in rows produced by pattern. The
// documented output contract preserves the input fields and appends anomaly
// detection results; the marker field names follow the documented anomaly
// concepts (score, type, flag). Frequency and new-pattern anomalies are the
// types computable from pattern output; the token and numerical token types
// require the raw token values that pattern output no longer carries.
type anomalyCommand struct {
	headTk token
}

func (c *anomalyCommand) name() string { return "anomaly" }

func (c *anomalyCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	var counts []float64
	for i := range rows {
		if v, ok := asNumber(rows[i].fields["@sampleCount"]); ok {
			counts = append(counts, v)
		}
	}
	if len(counts) < 2 {
		return rows
	}
	mean := 0.0
	for _, v := range counts {
		mean += v
	}
	mean /= float64(len(counts))
	variance := 0.0
	for _, v := range counts {
		variance += (v - mean) * (v - mean)
	}
	variance /= float64(len(counts))
	stddev := math.Sqrt(variance)
	for i := range rows {
		v, ok := asNumber(rows[i].fields["@sampleCount"])
		if !ok {
			continue
		}
		z := 0.0
		if stddev > 0 {
			z = (v - mean) / stddev
		}
		if z < 0 {
			z = -z
		}
		score := z / 3.0
		if score > 1 {
			score = 1
		}
		anomalous := z >= 2.0
		anomalyType := ""
		switch {
		case anomalous && v > mean:
			anomalyType = "PATTERN_FREQUENCY"
		case v <= 1 && len(rows) > 3:
			anomalyType = "NEW_PATTERN"
			anomalous = true
		}
		rows[i].set("@anomalyScore", strconv.FormatFloat(score, 'f', 3, 64))
		rows[i].set("@anomalyType", anomalyType)
		rows[i].set("@isAnomaly", strconv.FormatBool(anomalous))
	}
	return rows
}
