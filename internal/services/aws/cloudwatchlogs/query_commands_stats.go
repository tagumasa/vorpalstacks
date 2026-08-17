package cloudwatchlogs

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// The stats command: aggregation parsing, grouping with bin(), evaluation,
// and the fillmissing bin completion that consumes its output.

// --- stats ---

// binSpec describes a bin() grouping: bucket duration plus the optional
// offset that shifts the bin boundaries.
type binSpec struct {
	dur       int64
	offset    int64
	hasOffset bool
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
				// stats ... by bin(5m) offset 1h
				if p.acceptKeyword("offset") {
					ms, err := p.parseOffsetDuration(head)
					if err != nil {
						return err
					}
					spec.offset = ms
					spec.hasOffset = true
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
			// stats ... by bin(5m), otherField offset 1h — the offset
			// applies to the bin grouping earlier in the clause.
			if p.acceptKeyword("offset") {
				last := len(c.groups) - 1
				if last < 0 || c.groups[last].bin == nil {
					return newQueryCompileError("Syntax error: offset requires a bin() grouping", head.start, head.end)
				}
				ms, err := p.parseOffsetDuration(head)
				if err != nil {
					return err
				}
				c.groups[last].bin.offset = ms
				c.groups[last].bin.hasOffset = true
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
	if b.hasOffset {
		v -= b.offset
	}
	v = (v / b.dur) * b.dur
	if b.hasOffset {
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
