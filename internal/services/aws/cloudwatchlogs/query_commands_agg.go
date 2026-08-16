package cloudwatchlogs

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func jsonUnmarshalString(s string, out interface{}) error {
	return json.Unmarshal([]byte(strings.TrimSpace(s)), out)
}

// --- stats ---

// binSpec describes a bin() grouping: bucket duration plus the optional
// offset that shifts the bin boundaries.
type binSpec struct {
	dur      int64
	offset   int64
	hasOffst bool
}

type groupKeySpec struct {
	expr exprNode
	name string
	bin  *binSpec
}

type statsCommand struct {
	aggs    []aggSpec
	groups  []groupKeySpec
	headTk  token
	hasTopk bool
}

func (c *statsCommand) name() string { return "stats" }

// aggSpec is one aggregation output: an expression tree whose leaves are
// aggregation function calls and literals, combined with arithmetic.
type aggSpec struct {
	expr  exprNode
	alias string
}

func parseStatsCommand(args []token, head token) (command, error) {
	if len(args) == 0 {
		return nil, newQueryCompileError("Syntax error: stats requires aggregations", head.start, head.end)
	}
	// Split the aggregation list from the optional by-clause at the
	// top-level "by" keyword.
	byIdx := -1
	depth := 0
	for i, t := range args {
		switch t.kind {
		case tokLParen, tokLBracket:
			depth++
		case tokRParen, tokRBracket:
			depth--
		case tokIdent:
			if depth == 0 && strings.EqualFold(t.text, "by") {
				byIdx = i
			}
		}
		if byIdx >= 0 {
			break
		}
	}
	aggToks := args
	var byToks []token
	if byIdx >= 0 {
		aggToks = args[:byIdx]
		byToks = args[byIdx+1:]
	}
	c := &statsCommand{headTk: head}
	if err := parseAggList(aggToks, c, head); err != nil {
		return nil, err
	}
	if c.hasTopk && len(byToks) > 0 {
		return nil, newQueryCompileError("Syntax error: topk cannot be combined with by", head.start, head.end)
	}
	if len(byToks) > 0 {
		if err := parseGroupList(byToks, c, head); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func parseAggList(toks []token, c *statsCommand, head token) error {
	p := &exprParser{toks: toks}
	for {
		e, err := p.parseOr()
		if err != nil {
			return err
		}
		spec := aggSpec{expr: e}
		if fn, ok := e.(*funcNode); ok {
			spec.alias = fn.name + "(" + exprDisplay(fn.args) + ")"
		} else if lit, ok := e.(*literalNode); ok {
			spec.alias = asString(lit.val)
		}
		if p.acceptKeyword("as") {
			alias, ok := p.next()
			if !ok || (alias.kind != tokIdent && alias.kind != tokBacktickIdent) {
				return newQueryCompileError("Syntax error: as requires a name", head.start, head.end)
			}
			spec.alias = alias.text
		}
		if strings.HasPrefix(spec.alias, "topk") {
			c.hasTopk = true
		}
		c.aggs = append(c.aggs, spec)
		t, ok := p.next()
		if !ok {
			return nil
		}
		if t.kind != tokComma {
			return newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
	}
}

func exprDisplay(args []exprNode) string {
	parts := make([]string, 0, len(args))
	for _, a := range args {
		parts = append(parts, exprText(a))
	}
	return strings.Join(parts, ", ")
}

func exprText(e exprNode) string {
	switch n := e.(type) {
	case *literalNode:
		return asString(n.val)
	case *fieldNode:
		var b strings.Builder
		for i, s := range n.segs {
			if i > 0 {
				b.WriteByte('.')
			}
			if s.isIdx {
				b.WriteString("[")
				b.WriteString(strconv.Itoa(s.index))
				b.WriteString("]")
			} else {
				b.WriteString(s.name)
			}
		}
		return b.String()
	case *funcNode:
		return n.name + "(" + exprDisplay(n.args) + ")"
	case *binOpNode:
		return exprText(n.l) + " " + n.op + " " + exprText(n.r)
	}
	return ""
}

func parseGroupList(toks []token, c *statsCommand, head token) error {
	p := &exprParser{toks: toks}
	for {
		// bin(period) and offset are handled specially; datefloor/dateceil
		// are ordinary expressions over timestamps.
		if t, ok := p.peek(); ok && t.kind == tokIdent && strings.EqualFold(t.text, "bin") {
			if t2, ok2 := p.peekAt(1); ok2 && t2.kind == tokLParen {
				spec, err := parseBinCall(p)
				if err != nil {
					return err
				}
				c.groups = append(c.groups, groupKeySpec{name: "bin(" + strconv.FormatInt(spec.dur, 10) + ")", bin: spec})
				goto next
			}
		}
		{
			e, err := p.parseOr()
			if err != nil {
				return err
			}
			g := groupKeySpec{expr: e}
			if fn, ok := e.(*fieldNode); ok {
				g.name = fn.segs[0].name
			} else {
				g.name = exprText(e)
			}
			// stats ... by bin(5m) offset 1h
			if p.acceptKeyword("offset") {
				d, ok := p.next()
				if !ok {
					return newQueryCompileError("Syntax error: offset requires a duration", head.start, head.end)
				}
				ms, ok2 := parsePeriodMillis(d.text)
				if !ok2 {
					return newQueryCompileError(fmt.Sprintf("Invalid offset '%s'", d.text), d.start, d.end)
				}
				bin := &binSpec{dur: 0, offset: ms, hasOffst: true}
				if last := len(c.groups) - 1; last >= 0 && c.groups[last].bin != nil {
					c.groups[last].bin.offset = ms
					c.groups[last].bin.hasOffst = true
				} else {
					c.groups = append(c.groups, groupKeySpec{name: "offset", bin: bin})
				}
			}
			c.groups = append(c.groups, g)
		}
	next:
		t, ok := p.next()
		if !ok {
			return nil
		}
		if t.kind != tokComma {
			return newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
	}
}

// parseBinCall parses bin(period) with the documented unit caps.
func parseBinCall(p *exprParser) (*binSpec, error) {
	p.pos += 2 // bin (
	t, ok := p.next()
	if !ok {
		return nil, newQueryCompileError("Syntax error: bin requires a period", 0, 0)
	}
	period := t.text
	if t2, ok := p.peek(); ok && t2.kind == tokIdent && !strings.EqualFold(t2.text, "as") {
		// Units written as a separate word, e.g. bin(5 m).
		p.pos++
		period = t.text + t2.text
	}
	ms, ok := parsePeriodMillis(period)
	if !ok {
		return nil, newQueryCompileError(fmt.Sprintf("Invalid bin period '%s'", period), t.start, t.end)
	}
	unit := unitMillisOf(period)
	rp, ok := p.next()
	if !ok || rp.kind != tokRParen {
		return nil, newQueryCompileError("Syntax error: missing ')'", t.start, t.end)
	}
	return &binSpec{dur: capPeriodMillis(ms, unit)}, nil
}

// unitMillisOf returns the magnitude of one unit of the period for capping.
func unitMillisOf(period string) int64 {
	numPart := strings.TrimLeft(period, "0123456789.")
	if numPart == "" {
		return 1
	}
	if u, ok := parsePeriodMillis(numPart); ok && u > 0 {
		return u
	}
	return 1
}

// evalAggExpr evaluates an aggregation expression over a group of rows.
// Arithmetic over aggregation calls is supported per the documented grammar.
func evalAggExpr(e exprNode, rows []queryResultRow, ctx *execContext) interface{} {
	switch n := e.(type) {
	case *funcNode:
		return computeAggregation(strings.ToLower(n.name), n.args, rows, ctx)
	case *binOpNode:
		l := evalAggExpr(n.l, rows, ctx)
		r := evalAggExpr(n.r, rows, ctx)
		ln, lok := asNumber(l)
		rn, rok := asNumber(r)
		if !lok || !rok {
			return nil
		}
		return applyArith(n.op, ln, rn)
	case *literalNode:
		return n.val
	}
	return nil
}

// computeAggregation evaluates one aggregation function over the rows of a
// group.
func computeAggregation(name string, args []exprNode, rows []queryResultRow, ctx *execContext) interface{} {
	argValues := func() []interface{} {
		out := make([]interface{}, len(rows))
		for i := range rows {
			row := rows[i]
			if len(args) == 0 {
				out[i] = nil
				continue
			}
			out[i] = args[0].eval(&row, ctx)
		}
		return out
	}
	numeric := func() []float64 {
		var out []float64
		for _, v := range argValues() {
			if f, ok := asNumber(v); ok {
				out = append(out, f)
			}
		}
		return out
	}
	present := func() []interface{} {
		var out []interface{}
		for _, v := range argValues() {
			if v != nil && asString(v) != "" {
				out = append(out, v)
			}
		}
		return out
	}

	switch name {
	case "count":
		if len(args) == 0 {
			return float64(len(rows))
		}
		return float64(len(present()))
	case "avg":
		vals := numeric()
		if len(vals) == 0 {
			return nil
		}
		sum := 0.0
		for _, v := range vals {
			sum += v
		}
		return sum / float64(len(vals))
	case "sum":
		vals := numeric()
		if len(vals) == 0 {
			return float64(0)
		}
		sum := 0.0
		for _, v := range vals {
			sum += v
		}
		return sum
	case "min":
		vals := argValues()
		if len(vals) == 0 {
			return nil
		}
		best := vals[0]
		for _, v := range vals[1:] {
			if compareForSort(v, best) < 0 {
				best = v
			}
		}
		return wrapTimestampResult(args, best)
	case "max":
		vals := argValues()
		if len(vals) == 0 {
			return nil
		}
		best := vals[0]
		for _, v := range vals[1:] {
			if compareForSort(v, best) > 0 {
				best = v
			}
		}
		return wrapTimestampResult(args, best)
	case "countdistinct":
		seen := make(map[string]bool)
		for _, v := range present() {
			seen[asString(v)] = true
		}
		return float64(len(seen))
	case "stddev":
		vals := numeric()
		if len(vals) == 0 {
			return nil
		}
		mean := 0.0
		for _, v := range vals {
			mean += v
		}
		mean /= float64(len(vals))
		variance := 0.0
		for _, v := range vals {
			variance += (v - mean) * (v - mean)
		}
		variance /= float64(len(vals))
		return math.Sqrt(variance)
	case "variance":
		vals := numeric()
		if len(vals) == 0 {
			return nil
		}
		mean := 0.0
		for _, v := range vals {
			mean += v
		}
		mean /= float64(len(vals))
		variance := 0.0
		for _, v := range vals {
			variance += (v - mean) * (v - mean)
		}
		variance /= float64(len(vals))
		return variance
	case "pct":
		vals := numeric()
		if len(vals) == 0 || len(args) < 2 {
			return nil
		}
		p, ok := asNumber(args[1].eval(&rows[0], ctx))
		if !ok {
			return nil
		}
		return wrapTimestampResult(args, percentile(vals, p))
	case "values", "collect_values":
		seen := make(map[string]bool)
		var out []interface{}
		for _, v := range present() {
			s := asString(v)
			if !seen[s] {
				seen[s] = true
				out = append(out, s)
			}
		}
		if out == nil {
			out = []interface{}{}
		}
		return out
	case "topk":
		if len(args) < 2 || len(rows) == 0 {
			return nil
		}
		k, ok := asNumber(args[0].eval(&rows[0], ctx))
		if !ok || k < 1 {
			return nil
		}
		counts := make(map[string]int)
		for _, v := range present() {
			counts[asString(v)]++
		}
		type kv struct {
			k string
			v int
		}
		var list []kv
		for k, v := range counts {
			list = append(list, kv{k, v})
		}
		sort.Slice(list, func(i, j int) bool {
			if list[i].v != list[j].v {
				return list[i].v > list[j].v
			}
			return list[i].k < list[j].k
		})
		limit := int(k)
		if limit > len(list) {
			limit = len(list)
		}
		var out []interface{}
		for i := 0; i < limit; i++ {
			out = append(out, list[i].k)
		}
		return out

	// Time-series functions. The per-time-bin windowing comes from the
	// stats command's by bin() grouping; each aggregation here evaluates
	// within one bin (or the whole query window when ungrouped).
	case "countovertime":
		return float64(len(present()))
	case "sumovertime":
		sum := 0.0
		for _, v := range numeric() {
			sum += v
		}
		return sum
	case "rate":
		if len(args) < 2 || len(rows) == 0 {
			return nil
		}
		interval, ok := parsePeriodMillis(asString(args[1].eval(&rows[0], ctx)))
		if !ok || interval <= 0 {
			return nil
		}
		window := ctx.currentBinDur
		if window <= 0 {
			window = ctx.endTime - ctx.startTime
		}
		sum := 0.0
		for _, v := range numeric() {
			sum += v
		}
		return sum / (float64(window) / float64(interval))
	case "histogram":
		vals := numeric()
		if len(vals) == 0 || len(args) < 2 || len(rows) == 0 {
			return nil
		}
		buckets, ok := asNumber(args[1].eval(&rows[0], ctx))
		if !ok || buckets < 1 {
			return nil
		}
		n := int(buckets)
		min, max := vals[0], vals[0]
		for _, v := range vals[1:] {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
		width := (max - min) / float64(n)
		if width <= 0 {
			width = 1
		}
		counts := make(map[string]interface{}, n)
		for _, v := range vals {
			idx := int((v - min) / width)
			if idx >= n {
				idx = n - 1
			}
			if idx < 0 {
				idx = 0
			}
			key := formatNumber(min+width*float64(idx)) + "-" + formatNumber(min+width*float64(idx+1))
			if c, ok := counts[key].(int); ok {
				counts[key] = c + 1
			} else {
				counts[key] = 1
			}
		}
		return counts
	case "earliest", "latest":
		bestRow := -1
		var bestTS float64
		for i := range rows {
			row := rows[i]
			ts, ok := asNumber(row.fields["@timestamp"])
			if !ok {
				continue
			}
			if bestRow < 0 ||
				(name == "earliest" && ts < bestTS) ||
				(name == "latest" && ts > bestTS) {
				bestRow = i
				bestTS = ts
			}
		}
		if bestRow < 0 || len(args) == 0 {
			return nil
		}
		row := rows[bestRow]
		return wrapTimestampResult(args, args[0].eval(&row, ctx))
	case "sortsfirst", "sortslast":
		vals := present()
		if len(vals) == 0 {
			return nil
		}
		best := vals[0]
		for _, v := range vals[1:] {
			cmp := compareForSort(v, best)
			if (name == "sortsfirst" && cmp < 0) || (name == "sortslast" && cmp > 0) {
				best = v
			}
		}
		return wrapTimestampResult(args, best)
	}
	return nil
}

// isTimestampExpr reports whether the expression yields a timestamp-typed
// value: a direct event timestamp field, a datetime function whose
// documented result type is Timestamp, or the passthrough aggregation
// functions (documented result type LogField) applied to one of those.
func isTimestampExpr(e exprNode) bool {
	switch n := e.(type) {
	case *fieldNode:
		if len(n.segs) == 1 {
			switch strings.ToLower(n.segs[0].name) {
			case "@timestamp", "@ingestiontime":
				return true
			}
		}
	case *funcNode:
		switch strings.ToLower(n.name) {
		case "frommillis", "datefloor", "dateceil":
			return true
		case "earliest", "latest", "min", "max", "pct", "sortsfirst", "sortslast":
			return len(n.args) > 0 && isTimestampExpr(n.args[0])
		}
	}
	return false
}

// wrapTimestampResult renders a passthrough aggregation result as a
// timestamp when its argument is timestamp-typed, so the value keeps the
// documented Timestamp rendering in result rows while remaining numeric
// inside further expressions.
func wrapTimestampResult(args []exprNode, v interface{}) interface{} {
	if _, ok := v.(timestampValue); ok {
		return v
	}
	if len(args) == 0 || !isTimestampExpr(args[0]) {
		return v
	}
	if ms, ok := asNumber(v); ok {
		return timestampValue(int64(ms))
	}
	return v
}

func (c *statsCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	ctx.statsCount++
	rowsByGroup := make(map[string][]queryResultRow)
	partsByGroup := make(map[string][]string)
	var order []string
	for i := range rows {
		row := rows[i]
		var parts []string
		for _, g := range c.groups {
			if g.bin != nil {
				parts = append(parts, binKeyValue(row, g.bin))
				continue
			}
			parts = append(parts, asString(g.expr.eval(&row, ctx)))
		}
		key := strings.Join(parts, "\x00")
		if _, ok := rowsByGroup[key]; !ok {
			rowsByGroup[key] = []queryResultRow{}
			partsByGroup[key] = parts
			order = append(order, key)
		}
		rowsByGroup[key] = append(rowsByGroup[key], row)
	}

	emit := func(rowsIn []queryResultRow, key string) queryResultRow {
		row := queryResultRow{}
		for gi, gk := range c.groups {
			if gk.bin != nil {
				// The bin grouping is emitted as a timestamp-valued column
				// so that fillmissing and downstream sort can address it.
				row.set("@bin", formatResultTimestamp(partsByGroup[key][gi]))
				continue
			}
			row.set(gk.name, partsByGroup[key][gi])
		}
		for _, a := range c.aggs {
			row.set(a.alias, storeValue(evalAggExpr(a.expr, rowsIn, ctx)))
		}
		return row
	}

	binIdx := -1
	var binDur int64
	for gi, g := range c.groups {
		if g.bin != nil {
			binIdx = gi
			binDur = g.bin.dur
			break
		}
	}
	// The time-series functions scale by the bin window while this command
	// emits; without a bin grouping the query window applies.
	ctx.currentBinDur = binDur
	defer func() { ctx.currentBinDur = 0 }()

	if len(c.groups) == 0 {
		ctx.lastBins = nil
		return []queryResultRow{emit(rows, "")}
	}

	out := make([]queryResultRow, 0, len(order))
	for _, key := range order {
		out = append(out, emit(rowsByGroup[key], key))
	}
	if binIdx >= 0 && len(order) > 0 {
		minBin, maxBin := binRange(partsByGroup, order, binIdx)
		ctx.lastBins = &binFillInfo{
			dur:    binDur,
			minBin: minBin,
			maxBin: maxBin,
		}
	} else {
		ctx.lastBins = nil
	}
	return out
}

func binKeyValue(row queryResultRow, b *binSpec) string {
	ts, ok := asNumber(row.fields["@timestamp"])
	if !ok {
		return ""
	}
	v := int64(ts)
	if b.hasOffst {
		v -= b.offset
	}
	v = (v / b.dur) * b.dur
	if b.hasOffst {
		v += b.offset
	}
	return strconv.FormatInt(v, 10)
}

func binRange(partsByGroup map[string][]string, order []string, binIdx int) (int64, int64) {
	var min, max int64
	first := true
	for _, key := range order {
		v, ok := asNumber(partsByGroup[key][binIdx])
		if !ok {
			continue
		}
		if first || int64(v) < min {
			min = int64(v)
		}
		if first || int64(v) > max {
			max = int64(v)
		}
		first = false
	}
	return min, max
}

// --- fillmissing ---

type fillmissingCommand struct {
	fills  []fillConst
	headTk token
}

type fillConst struct {
	value string
	field string
}

func (c *fillmissingCommand) name() string { return "fillmissing" }

func parseFillmissingCommand(args []token, head token) (command, error) {
	c := &fillmissingCommand{headTk: head}
	i := 0
	for i < len(args) {
		t := args[i]
		if !(t.kind == tokIdent && strings.EqualFold(t.text, "with")) {
			return nil, newQueryCompileError(
				fmt.Sprintf("Syntax error at '%s': fillmissing expects with <value> for <field>", t.raw), t.start, t.end)
		}
		if i+4 > len(args) {
			return nil, newQueryCompileError("Syntax error: incomplete with clause", head.start, head.end)
		}
		val := args[i+1]
		if val.kind != tokIdent && val.kind != tokNumber && val.kind != tokString {
			return nil, newQueryCompileError("Syntax error: with requires a value", val.start, val.end)
		}
		if args[i+2].kind != tokIdent || !strings.EqualFold(args[i+2].text, "for") {
			return nil, newQueryCompileError("Syntax error: expected for", args[i+2].start, args[i+2].end)
		}
		fld := args[i+3]
		if fld.kind != tokIdent && fld.kind != tokBacktickIdent {
			return nil, newQueryCompileError("Syntax error: for requires a field", fld.start, fld.end)
		}
		c.fills = append(c.fills, fillConst{value: val.text, field: fld.text})
		i += 4
		if i < len(args) && args[i].kind == tokComma {
			i++
		}
	}
	return c, nil
}

func (c *fillmissingCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	info := ctx.lastBins
	if info == nil || info.dur <= 0 {
		return rows
	}
	existing := make(map[string]bool, len(rows))
	for _, r := range rows {
		existing[r.fields["@bin"]] = true
	}
	var out []queryResultRow
	out = append(out, rows...)
	for bin := info.minBin; bin <= info.maxBin; bin += info.dur {
		key := formatResultTimestamp(strconv.FormatInt(bin, 10))
		if existing[key] {
			continue
		}
		synth := queryResultRow{}
		synth.set("@bin", key)
		for _, f := range c.fills {
			synth.set(f.field, f.value)
		}
		out = append(out, synth)
	}
	// Bins are stored in their timestamp rendering; the fixed-width layout
	// orders chronologically, with the numeric form kept as a fallback for
	// defensive sorting of unformatted values.
	binOrder := func(v string) int64 {
		if ms, ok := parseResultTimestamp(v); ok {
			return ms
		}
		ms, _ := asNumber(v)
		return int64(ms)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return binOrder(out[i].fields["@bin"]) < binOrder(out[j].fields["@bin"])
	})
	return out
}

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
	// Comparison queries analyse patterns: the current rows are clustered,
	// the comparison window is re-queried with the same preceding pipeline,
	// and the per-pattern count differences are reported.
	curPatterns := patternCounts(rows)

	ps, pe := c.comparisonWindow(ctx.startTime, ctx.endTime)
	prevRows, err := ctx.runPrecedingOnWindow(ps, pe)
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
	if len(args) != 2 || args[0].kind != tokIdent || !strings.EqualFold(args[0].text, "timeshift") {
		return nil, newQueryCompileError("Syntax error: logcompare requires timeshift <duration>", head.start, head.end)
	}
	ms, ok := parsePeriodMillis(args[1].text)
	if !ok {
		return nil, newQueryCompileError(fmt.Sprintf("Invalid timeshift '%s'", args[1].text), args[1].start, args[1].end)
	}
	return &logCompareCommand{timeshift: ms, headTk: head}, nil
}

func (c *logCompareCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	curPatterns := patternCounts(rows)
	ps, pe := ctx.startTime-c.timeshift, ctx.endTime-c.timeshift
	prevRows, err := ctx.runPrecedingOnWindow(ps, pe)
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
