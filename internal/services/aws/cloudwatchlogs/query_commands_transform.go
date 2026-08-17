package cloudwatchlogs

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Row-transforming pipeline commands: parse extraction, accum,
// autoregress, addtotals, outlier, and sessionize.

// --- parse (glob and regular expression extraction) ---

type parseExtractCommand struct {
	source  exprNode
	pattern string
	re      *regexp.Regexp
	names   []string
	headTk  token
}

func (c *parseExtractCommand) name() string { return "parse" }

func parseParseCommand(args []token, head token) (command, error) {
	if len(args) < 2 {
		return nil, newQueryCompileError("Syntax error: parse requires a field and a pattern", head.start, head.end)
	}
	c := &parseExtractCommand{headTk: head}
	p := &exprParser{toks: args}
	src, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	c.source = src
	pt, ok := p.next()
	if !ok {
		return nil, newQueryCompileError("Syntax error: parse requires a pattern", head.start, head.end)
	}
	switch pt.kind {
	case tokString:
		c.pattern = pt.text
	case tokRegex:
		re, err := regexp.Compile(pt.text)
		if err != nil {
			return nil, newQueryCompileError(fmt.Sprintf("Invalid regular expression: %v", err), pt.start, pt.end)
		}
		c.re = re
	case tokOp:
		if pt.text != "/" {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", pt.raw), pt.start, pt.end)
		}
		var raw strings.Builder
		closing, ok := p.next()
		for ok && !(closing.kind == tokOp && closing.text == "/") {
			if raw.Len() > 0 {
				raw.WriteByte(' ')
			}
			raw.WriteString(closing.raw)
			closing, ok = p.next()
		}
		if !ok {
			return nil, newQueryCompileError("Unterminated regular expression", pt.start, pt.end)
		}
		re, err := regexp.Compile(raw.String())
		if err != nil {
			return nil, newQueryCompileError(fmt.Sprintf("Invalid regular expression: %v", err), pt.start, closing.end)
		}
		c.re = re
	default:
		return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", pt.raw), pt.start, pt.end)
	}
	if !p.acceptKeyword("as") {
		return nil, newQueryCompileError("Syntax error: parse requires as", head.start, head.end)
	}
	for {
		n, ok := p.next()
		if !ok || (n.kind != tokIdent && n.kind != tokBacktickIdent) {
			return nil, newQueryCompileError("Syntax error: parse requires field names after as", head.start, head.end)
		}
		c.names = append(c.names, n.text)
		t, ok := p.next()
		if !ok {
			break
		}
		if t.kind != tokComma {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
	}
	if c.re == nil {
		// Glob mode: each * captures one field.
		wildcards := strings.Count(c.pattern, "*")
		if wildcards != len(c.names) {
			return nil, newQueryCompileError(
				fmt.Sprintf("parse pattern has %d wildcards but %d field names", wildcards, len(c.names)), head.start, head.end)
		}
	}
	return c, nil
}

func (c *parseExtractCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	for i := range rows {
		row := rows[i]
		src := asString(c.source.eval(&row, ctx))
		if c.re != nil {
			m := c.re.FindStringSubmatch(src)
			if m == nil {
				continue
			}
			for j, n := range c.names {
				if j+1 < len(m) {
					rows[i].set(n, m[j+1])
				}
			}
			continue
		}
		matches := globCapture(c.pattern, src)
		for j, n := range c.names {
			if j < len(matches) {
				rows[i].set(n, matches[j])
			}
		}
	}
	return rows
}

// globCapture matches text against a glob pattern and returns the values
// captured by the wildcards.
func globCapture(pattern, text string) []string {
	parts := strings.Split(pattern, "*")
	if len(parts) < 2 {
		return nil
	}
	var re strings.Builder
	for i, p := range parts {
		if i > 0 {
			re.WriteString("(.*?)")
		}
		re.WriteString(regexp.QuoteMeta(p))
	}
	compiled, err := regexp.Compile("^" + re.String() + "$")
	if err != nil {
		return nil
	}
	m := compiled.FindStringSubmatch(text)
	if m == nil {
		return nil
	}
	return m[1:]
}

// --- accum ---

type accumCommand struct {
	field exprNode
	out   string
	head  token
}

func (c *accumCommand) name() string { return "accum" }

func parseAccumCommand(args []token, head token) (command, error) {
	if len(args) == 0 {
		return nil, newQueryCompileError("Syntax error: accum requires a field", head.start, head.end)
	}
	p := &exprParser{toks: args}
	e, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	c := &accumCommand{field: e, head: head}
	if fn, ok := e.(*fieldNode); ok {
		c.out = fn.segs[0].name
	} else {
		c.out = "accum"
	}
	if p.acceptKeyword("as") {
		alias, ok := p.next()
		if !ok || (alias.kind != tokIdent && alias.kind != tokBacktickIdent) {
			return nil, newQueryCompileError("Syntax error: as requires a name", head.start, head.end)
		}
		c.out = alias.text
	}
	return c, nil
}

func (c *accumCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	total := 0.0
	for i := range rows {
		row := rows[i]
		if f, ok := asNumber(c.field.eval(&row, ctx)); ok {
			total += f
			rows[i].set(c.out, formatNumber(total))
		} else {
			rows[i].set(c.out, "")
		}
	}
	return rows
}

// --- autoregress ---

type autoregressCommand struct {
	field string
	alias string
	pFrom int
	pTo   int
	head  token
}

func (c *autoregressCommand) name() string { return "autoregress" }

func parseAutoregressCommand(args []token, head token) (command, error) {
	if len(args) == 0 {
		return nil, newQueryCompileError("Syntax error: autoregress requires a field", head.start, head.end)
	}
	c := &autoregressCommand{head: head, pFrom: 1, pTo: 1}
	c.alias = args[0].text
	c.field = args[0].text
	i := 1
	for i < len(args) {
		t := args[i]
		if t.kind == tokIdent && strings.EqualFold(t.text, "as") {
			if i+1 >= len(args) {
				return nil, newQueryCompileError("Syntax error: as requires a name", head.start, head.end)
			}
			c.alias = args[i+1].text
			i += 2
			continue
		}
		if t.kind == tokIdent && strings.EqualFold(t.text, "p") && i+1 < len(args) && args[i+1].kind == tokOp && args[i+1].text == "=" {
			// The lag depth lexes as number, optional minus, number.
			if i+2 >= len(args) || args[i+2].kind != tokNumber {
				return nil, newQueryCompileError("Syntax error: expected lag depth", t.start, t.end)
			}
			spec := args[i+2].text
			consumed := 3
			if i+4 < len(args) && args[i+3].kind == tokOp && args[i+3].text == "-" && args[i+4].kind == tokNumber {
				spec = args[i+2].text + "-" + args[i+4].text
				consumed = 5
			}
			from, to, err := parseLagSpecRaw(spec, args[i+2])
			if err != nil {
				return nil, err
			}
			c.pFrom, c.pTo = from, to
			i += consumed
			continue
		}
		return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
	}
	return c, nil
}

func parseLagSpecRaw(s string, t token) (int, int, error) {
	if idx := strings.Index(s, "-"); idx > 0 {
		from, err1 := strconv.Atoi(s[:idx])
		to, err2 := strconv.Atoi(s[idx+1:])
		if err1 != nil || err2 != nil || from <= 0 || to < from {
			return 0, 0, newQueryCompileError(fmt.Sprintf("Invalid lag range '%s'", s), t.start, t.end)
		}
		return from, to, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0, 0, newQueryCompileError(fmt.Sprintf("Invalid lag depth '%s'", s), t.start, t.end)
	}
	return n, n, nil
}

func (c *autoregressCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	for lag := c.pFrom; lag <= c.pTo; lag++ {
		out := c.alias + "_p" + strconv.Itoa(lag)
		for i := range rows {
			if i-lag >= 0 {
				rows[i].set(out, readField(rows[i-lag], c.field, ctx))
			} else {
				rows[i].set(out, "")
			}
		}
	}
	return rows
}

// --- addtotals ---

type addTotalsCommand struct {
	fieldName string
	row       bool
	col       bool
	fields    []string
	head      token
}

func (c *addTotalsCommand) name() string { return "addtotals" }

func parseAddTotalsCommand(args []token, head token) (command, error) {
	c := &addTotalsCommand{head: head, fieldName: "Total", row: true, col: false}
	for i := 0; i < len(args); i++ {
		t := args[i]
		if t.kind == tokComma {
			continue
		}
		if t.kind != tokIdent && t.kind != tokBacktickIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
		lower := strings.ToLower(t.text)
		switch {
		case strings.HasPrefix(lower, "fieldname") && i+2 < len(args) && args[i+1].kind == tokOp && args[i+1].text == "=":
			c.fieldName = args[i+2].text
			i += 2
		case lower == "row" && i+2 < len(args) && args[i+1].kind == tokOp && args[i+1].text == "=":
			c.row = strings.EqualFold(args[i+2].text, "true")
			i += 2
		case lower == "col" && i+2 < len(args) && args[i+1].kind == tokOp && args[i+1].text == "=":
			c.col = strings.EqualFold(args[i+2].text, "true")
			i += 2
		default:
			c.fields = append(c.fields, t.text)
		}
	}
	return c, nil
}

func (c *addTotalsCommand) targetFields(row queryResultRow) []string {
	if len(c.fields) > 0 {
		return c.fields
	}
	var numeric []string
	for k, v := range row.fields {
		if _, ok := asNumber(v); ok {
			numeric = append(numeric, k)
		}
	}
	sort.Strings(numeric)
	return numeric
}

func (c *addTotalsCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	if len(rows) == 0 {
		return rows
	}
	targets := c.targetFields(rows[0])
	if c.row {
		for i := range rows {
			sum := 0.0
			any := false
			for _, f := range targets {
				if v, ok := asNumber(rows[i].fields[f]); ok {
					sum += v
					any = true
				}
			}
			if any {
				rows[i].set(c.fieldName, formatNumber(sum))
			}
		}
	}
	if c.col {
		totals := queryResultRow{}
		for _, f := range targets {
			sum := 0.0
			any := false
			for i := range rows {
				if v, ok := asNumber(rows[i].fields[f]); ok {
					sum += v
					any = true
				}
			}
			if any {
				totals.set(f, formatNumber(sum))
			}
		}
		rows = append(rows, totals)
	}
	return rows
}

// --- outlier ---

type outlierCommand struct {
	action   string // remove or transform
	param    float64
	useLower bool
	mark     bool
	fields   []string
	head     token
}

func (c *outlierCommand) name() string { return "outlier" }

func parseOutlierCommand(args []token, head token) (command, error) {
	c := &outlierCommand{head: head, action: "transform", param: 2.5}
	for i := 0; i < len(args); i++ {
		t := args[i]
		if t.kind != tokIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
		lower := strings.ToLower(t.text)
		isKV := i+2 < len(args) && args[i+1].kind == tokOp && args[i+1].text == "="
		switch {
		case lower == "action" && isKV:
			c.action = strings.ToLower(args[i+2].text)
			i += 2
		case lower == "param" && isKV:
			f, err := strconv.ParseFloat(args[i+2].text, 64)
			if err != nil {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid param '%s'", args[i+2].text), args[i+2].start, args[i+2].end)
			}
			c.param = f
			i += 2
		case lower == "uselower" && isKV:
			c.useLower = strings.EqualFold(args[i+2].text, "true")
			i += 2
		case lower == "mark" && isKV:
			c.mark = strings.EqualFold(args[i+2].text, "true")
			i += 2
		default:
			c.fields = append(c.fields, t.text)
		}
	}
	if c.action != "remove" && c.action != "transform" {
		return nil, newQueryCompileError(fmt.Sprintf("Invalid outlier action '%s'", c.action), head.start, head.end)
	}
	if len(c.fields) == 0 {
		return nil, newQueryCompileError("Syntax error: outlier requires at least one field", head.start, head.end)
	}
	return c, nil
}

func (c *outlierCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	bounds := make(map[string][2]float64, len(c.fields))
	for _, f := range c.fields {
		var vals []float64
		for i := range rows {
			if v, ok := asNumber(readField(rows[i], f, ctx)); ok {
				vals = append(vals, v)
			}
		}
		if len(vals) == 0 {
			continue
		}
		q1 := percentile(vals, 25)
		q3 := percentile(vals, 75)
		iqr := q3 - q1
		bounds[f] = [2]float64{q1 - c.param*iqr, q3 + c.param*iqr}
	}
	out := make([]queryResultRow, 0, len(rows))
	for i := range rows {
		row := &rows[i]
		isOutlier := false
		for _, f := range c.fields {
			b, ok := bounds[f]
			if !ok {
				continue
			}
			v, ok := asNumber(readField(*row, f, ctx))
			if !ok {
				continue
			}
			high := v > b[1]
			low := c.useLower && v < b[0]
			if c.mark {
				row.set(f+"_outlier", strconv.FormatBool(high || low))
			}
			if high || low {
				isOutlier = true
				if c.action == "transform" {
					if high {
						row.set(f, formatNumber(b[1]))
					} else {
						row.set(f, formatNumber(b[0]))
					}
					isOutlier = false
				}
			}
		}
		if c.action == "remove" && isOutlier {
			continue
		}
		out = append(out, *row)
	}
	return out
}

// percentile computes the nearest-rank percentile of a value slice. The
// slice is sorted in place.
func percentile(vals []float64, p float64) float64 {
	sort.Float64s(vals)
	if len(vals) == 1 {
		return vals[0]
	}
	rank := int(math.Ceil(p / 100 * float64(len(vals))))
	if rank < 1 {
		rank = 1
	}
	if rank > len(vals) {
		rank = len(vals)
	}
	return vals[rank-1]
}

// --- sessionize ---

type sessionizeCommand struct {
	fields  []string
	maxSpan int64
	out     string
	head    token
}

func (c *sessionizeCommand) name() string { return "sessionize" }

func parseSessionizeCommand(args []token, head token) (command, error) {
	c := &sessionizeCommand{head: head, maxSpan: 30 * 60 * 1000, out: "session_id"}
	for i := 0; i < len(args); i++ {
		t := args[i]
		if t.kind != tokIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
		lower := strings.ToLower(t.text)
		switch {
		case lower == "maxspan" && i+1 < len(args):
			text, span, ok := periodArgAt(args, i+1)
			if !ok {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid maxspan '%s'", args[i+1].text), args[i+1].start, args[i+1].end)
			}
			ms, ok2 := parsePeriodMillis(text)
			if !ok2 {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid maxspan '%s'", text), args[i+1].start, args[i+1].end)
			}
			c.maxSpan = ms
			i += span
		case lower == "as":
			if i+1 >= len(args) {
				return nil, newQueryCompileError("Syntax error: as requires a name", head.start, head.end)
			}
			c.out = args[i+1].text
			i++
		default:
			c.fields = append(c.fields, t.text)
		}
	}
	if len(c.fields) == 0 {
		return nil, newQueryCompileError("Syntax error: sessionize requires identity fields", head.start, head.end)
	}
	return c, nil
}

func (c *sessionizeCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	// Sessions are formed per identity value combination: events ordered by
	// @timestamp start a new session when the inactivity gap exceeds maxspan.
	type sessionState struct {
		lastTS int64
		id     int
	}
	sessions := make(map[string]*sessionState)
	counter := 0
	idx := make([]int, len(rows))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(a, b int) bool {
		va, _ := asNumber(rows[idx[a]].fields["@timestamp"])
		vb, _ := asNumber(rows[idx[b]].fields["@timestamp"])
		return va < vb
	})
	for _, i := range idx {
		var key strings.Builder
		for _, f := range c.fields {
			key.WriteString(readField(rows[i], f, ctx))
			key.WriteByte(0)
		}
		ts, _ := asNumber(rows[i].fields["@timestamp"])
		st, ok := sessions[key.String()]
		if !ok || int64(ts)-st.lastTS > c.maxSpan {
			counter++
			st = &sessionState{id: counter}
			sessions[key.String()] = st
		}
		st.lastTS = int64(ts)
		rows[i].set(c.out, strconv.Itoa(st.id))
	}
	return rows
}
