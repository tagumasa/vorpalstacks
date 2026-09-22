package cloudwatchlogs

import (
	"fmt"
	"strconv"
	"strings"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
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
	// "This function cannot be combined with by or other aggregation
	// functions" — topk shares no stats command with another output.
	if c.hasTopk && len(c.aggs) > 1 {
		return nil, newQueryCompileError("Syntax error: topk cannot be combined with other aggregation functions", head.start, head.end)
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
		} else {
			spec.alias = exprText(e)
		}
		if p.acceptKeyword("as") {
			alias, ok := p.next()
			if !ok || (alias.kind != tokIdent && alias.kind != tokBacktickIdent) {
				return newQueryCompileError("Syntax error: as requires a name", head.start, head.end)
			}
			spec.alias = alias.text
		}
		if fn, ok := e.(*funcNode); ok && strings.EqualFold(fn.name, "topk") {
			c.hasTopk = true
			// "Valid values for k range from 1 to 10000" — the literal k
			// rejects at compile; a computed k bounds again at
			// evaluation.
			if lit, lok := fn.args[0].(*literalNode); lok {
				if k, kok := asNumber(lit.val); !kok || k < 1 || k > logsstore.MaxTopkK {
					return newQueryCompileError(fmt.Sprintf(
						"Invalid topk k value: valid values range from 1 to %d", logsstore.MaxTopkK),
						head.start, head.end)
				}
			}
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
				// The bin group key carries the emitted column's name:
				// the results and the stats visibility set both address
				// the bucket column as @bin, so the key must equal the
				// name a subsequent command reads — a key of bin(<ms>)
				// could never match the emitted @bin.
				c.groups = append(c.groups, groupKeySpec{name: "@bin", bin: spec})
				goto next
			}
		}
		{
			e, err := p.parseOr()
			if err != nil {
				return err
			}
			g := groupKeySpec{expr: e}
			g.name = exprText(e)
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
