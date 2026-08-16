package cloudwatchlogs

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// command is one stage of a parsed query pipeline.
type command interface {
	name() string
	apply(ctx *execContext, rows []queryResultRow) []queryResultRow
}

// parsePipelineCommands splits a token stream into pipeline commands at
// top-level pipes. Parenthesised and bracketed groups keep their inner pipes
// so that subqueries and join sources stay intact.
func parsePipelineCommands(toks []token) ([][]token, error) {
	var segs [][]token
	var cur []token
	depth := 0
	for _, t := range toks {
		switch t.kind {
		case tokLParen, tokLBracket, tokLBrace:
			depth++
		case tokRParen, tokRBracket, tokRBrace:
			depth--
		case tokPipe:
			if depth == 0 {
				segs = append(segs, cur)
				cur = nil
				continue
			}
		}
		cur = append(cur, t)
	}
	if len(cur) > 0 || len(segs) > 0 {
		segs = append(segs, cur)
	}
	return segs, nil
}

// parseCommand parses one pipeline segment into a command.
func parseCommand(toks []token) (command, error) {
	if len(toks) == 0 {
		return nil, nil
	}
	head := toks[0]
	args := toks[1:]
	if head.kind == tokIdent && strings.EqualFold(head.text, "SOURCE") {
		return parseSourceCommand(args, head)
	}
	if head.kind != tokIdent && head.kind != tokBacktickIdent {
		return nil, newQueryCompileError(
			fmt.Sprintf("Syntax error at '%s': expected a command name", head.raw), head.start, head.end)
	}
	cmdName := strings.ToLower(head.text)
	switch cmdName {
	case "fields":
		return parseProjectionCommand("fields", args, head)
	case "display":
		return parseProjectionCommand("display", args, head)
	case "filter", "where", "filterindex":
		return parseFilterCommand(cmdName, args, head)
	case "sort":
		return parseSortCommand(args, head)
	case "limit":
		return parseLimitCommand(args, head)
	case "parse":
		return parseParseCommand(args, head)
	case "dedup":
		return parseDedupCommand(args, head)
	case "filldown":
		return parseFilldownCommand(args, head)
	case "accum":
		return parseAccumCommand(args, head)
	case "autoregress":
		return parseAutoregressCommand(args, head)
	case "addtotals":
		return parseAddTotalsCommand(args, head)
	case "expand":
		return parseExpandCommand(args, head)
	case "unnest":
		return parseUnnestCommand(args, head)
	case "outlier":
		return parseOutlierCommand(args, head)
	case "sessionize":
		return parseSessionizeCommand(args, head)
	case "countfrequent", "count_frequent":
		return parseCountFrequentCommand(args, head)
	case "unmask":
		return &unmaskCommand{cmdName: cmdName, headTk: head}, nil
	case "stats":
		return parseStatsCommand(args, head)
	case "pattern":
		return parsePatternCommand(args, head)
	case "diff":
		return parseDiffCommand(args, head)
	case "logcompare":
		return parseLogCompareCommand(args, head)
	case "relevantfields":
		return parseRelevantFieldsCommand(args, head)
	case "anomaly":
		return &anomalyCommand{headTk: head}, nil
	case "fillmissing":
		return parseFillmissingCommand(args, head)
	case "join":
		return parseJoinCommand(args, head)
	case "appendcols":
		return parseAppendColsCommand(args, head)
	case "lookup":
		return parseLookupCommand(args, head)
	case "cidrlookup":
		return parseCidrLookupCommand(args, head)
	}
	return nil, newQueryCompileError(
		fmt.Sprintf("Unknown command '%s' in query string", head.raw), head.start, head.end)
}

// pathNode builds a field accessor for a possibly dotted field name so that
// commands taking bare field names resolve JSON structure paths the same way
// expressions do.
func pathNode(name string) *fieldNode {
	fn := &fieldNode{}
	for _, part := range strings.Split(name, ".") {
		if part != "" {
			fn.segs = append(fn.segs, pathSeg{name: part})
		}
	}
	return fn
}

// readField resolves a field name (with dots) against a row.
func readField(row queryResultRow, name string, ctx *execContext) string {
	return asString(pathNode(name).eval(&row, ctx))
}

// --- projection: fields / display ---

type projectionItem struct {
	expr  exprNode
	alias string
	isRef bool // plain field reference without alias
	name  string
}

type projectionCommand struct {
	cmd    string
	items  []projectionItem
	headTk token
}

func (c *projectionCommand) name() string { return c.cmd }

// fieldPathDisplay renders a field path for use as an output field name.
func fieldPathDisplay(fn *fieldNode) string {
	parts := make([]string, 0, len(fn.segs))
	for _, seg := range fn.segs {
		if seg.isIdx {
			parts = append(parts, "["+strconv.Itoa(seg.index)+"]")
		} else {
			parts = append(parts, seg.name)
		}
	}
	return strings.Join(parts, ".")
}

// parseProjectionCommand parses `fields a, b as x, strlen(@message) as len`.
func parseProjectionCommand(cmd string, args []token, head token) (command, error) {
	c := &projectionCommand{cmd: cmd, headTk: head}
	if len(args) == 0 {
		return nil, newQueryCompileError(fmt.Sprintf("Syntax error: %s requires a field list", cmd), head.start, head.end)
	}
	p := &exprParser{toks: args}
	for {
		e, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		item := projectionItem{expr: e}
		if fn, ok := e.(*fieldNode); ok {
			item.isRef = true
			item.name = fieldPathDisplay(fn)
		}
		if p.acceptKeyword("as") {
			alias, ok := p.next()
			if !ok || (alias.kind != tokIdent && alias.kind != tokBacktickIdent) {
				return nil, newQueryCompileError("Syntax error: as requires a name", head.start, head.end)
			}
			item.alias = alias.text
			item.name = alias.text
		}
		c.items = append(c.items, item)
		t, ok := p.next()
		if !ok {
			break
		}
		if t.kind != tokComma {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
	}
	return c, nil
}

func (c *projectionCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	out := make([]queryResultRow, 0, len(rows))
	for i := range rows {
		row := rows[i]
		projected := queryResultRow{}
		for _, it := range c.items {
			v := it.expr.eval(&row, ctx)
			if v == nil && it.isRef {
				continue // absent fields are omitted from the projection
			}
			projected.set(it.name, storeValue(v))
		}
		// Event rows keep their @ptr through projection so that results
		// remain addressable with GetLogRecord.
		if ptr, ok := row.fields["@ptr"]; ok {
			projected.set("@ptr", ptr)
		}
		out = append(out, projected)
	}
	return out
}

// --- filter / where / filterIndex ---

type filterCommand struct {
	cmd    string
	expr   exprNode
	headTk token
}

func (c *filterCommand) name() string { return c.cmd }

func parseFilterCommand(cmd string, args []token, head token) (command, error) {
	if len(args) == 0 {
		return nil, newQueryCompileError(fmt.Sprintf("Syntax error: %s requires an expression", cmd), head.start, head.end)
	}
	expr, err := parseExprTokens(args)
	if err != nil {
		return nil, err
	}
	return &filterCommand{cmd: cmd, expr: expr, headTk: head}, nil
}

func (c *filterCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	var out []queryResultRow
	for i := range rows {
		row := rows[i]
		if truthy(c.expr.eval(&row, ctx)) {
			out = append(out, row)
		}
	}
	if out == nil {
		out = []queryResultRow{}
	}
	return out
}

// --- sort ---

type sortKey struct {
	expr exprNode
	desc bool
}

type sortCommand struct {
	keys   []sortKey
	headTk token
}

func (c *sortCommand) name() string { return "sort" }

func parseSortCommand(args []token, head token) (command, error) {
	if len(args) == 0 {
		return nil, newQueryCompileError("Syntax error: sort requires a field", head.start, head.end)
	}
	c := &sortCommand{headTk: head}
	p := &exprParser{toks: args}
	for {
		e, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		key := sortKey{expr: e}
		if p.acceptKeyword("desc") {
			key.desc = true
		} else if p.acceptKeyword("asc") {
			key.desc = false
		}
		c.keys = append(c.keys, key)
		t, ok := p.next()
		if !ok {
			break
		}
		if t.kind != tokComma {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
	}
	return c, nil
}

func (c *sortCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	sort.SliceStable(rows, func(i, j int) bool {
		for _, k := range c.keys {
			vi := k.expr.eval(&rows[i], ctx)
			vj := k.expr.eval(&rows[j], ctx)
			if vi == nil && vj == nil {
				continue
			}
			if vi == nil {
				return false
			}
			if vj == nil {
				return true
			}
			cmp := compareForSort(vi, vj)
			if cmp == 0 {
				continue
			}
			if k.desc {
				return cmp > 0
			}
			return cmp < 0
		}
		return false
	})
	ctx.sorted = true
	return rows
}

// compareForSort returns -1/0/1 with numeric comparison when both values are
// numeric and string comparison otherwise.
func compareForSort(a, b interface{}) int {
	an, aok := asNumber(a)
	bn, bok := asNumber(b)
	if aok && bok {
		switch {
		case an < bn:
			return -1
		case an > bn:
			return 1
		}
		return 0
	}
	return strings.Compare(asString(a), asString(b))
}

// --- limit ---

type limitCommand struct {
	n      int
	any    bool
	headTk token
}

func (c *limitCommand) name() string { return "limit" }

func parseLimitCommand(args []token, head token) (command, error) {
	c := &limitCommand{headTk: head}
	p := &exprParser{toks: args}
	if p.acceptKeyword("any") {
		c.any = true
	}
	n, ok := p.next()
	if !ok || n.kind != tokNumber {
		return nil, newQueryCompileError("Syntax error: limit requires a number", head.start, head.end)
	}
	v, err := strconv.Atoi(n.text)
	if err != nil || v <= 0 {
		return nil, newQueryCompileError(fmt.Sprintf("Invalid limit '%s'", n.text), n.start, n.end)
	}
	c.n = v
	if _, ok := p.next(); ok {
		return nil, newQueryCompileError("Syntax error: limit takes a single number", head.start, head.end)
	}
	return c, nil
}

func (c *limitCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	if len(rows) <= c.n {
		return rows
	}
	return rows[:c.n]
}

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
		matches := globCapture(c.pattern, src, len(c.names))
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
func globCapture(pattern, text string, want int) []string {
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

// --- dedup ---

type dedupCommand struct {
	fields []exprNode
	headTk token
}

func (c *dedupCommand) name() string { return "dedup" }

func parseDedupCommand(args []token, head token) (command, error) {
	c := &dedupCommand{headTk: head}
	p := &exprParser{toks: args}
	if len(args) == 0 {
		return c, nil
	}
	for {
		e, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		c.fields = append(c.fields, e)
		t, ok := p.next()
		if !ok {
			break
		}
		if t.kind != tokComma {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
	}
	return c, nil
}

func (c *dedupCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	work := rows
	if !ctx.sorted {
		// Without an explicit sort, duplicates are discarded against the
		// default descending @timestamp order.
		work = make([]queryResultRow, len(rows))
		copy(work, rows)
		sort.SliceStable(work, func(i, j int) bool {
			vi, _ := asNumber(work[i].fields["@timestamp"])
			vj, _ := asNumber(work[j].fields["@timestamp"])
			return vi > vj
		})
	}
	seen := make(map[string]bool)
	out := make([]queryResultRow, 0, len(work))
	for i := range work {
		row := work[i]
		if len(c.fields) == 0 {
			key := rowSignature(row)
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, row)
			continue
		}
		// Rows with a missing value for any dedup field are retained; null
		// values are not considered duplicates.
		missing := false
		var key strings.Builder
		for _, f := range c.fields {
			v := f.eval(&row, ctx)
			if v == nil || asString(v) == "" {
				missing = true
				break
			}
			key.WriteString(asString(v))
			key.WriteByte(0)
		}
		if missing {
			out = append(out, row)
			continue
		}
		if seen[key.String()] {
			continue
		}
		seen[key.String()] = true
		out = append(out, row)
	}
	return out
}

func rowSignature(row queryResultRow) string {
	keys := make([]string, 0, len(row.fields))
	for k := range row.fields {
		// @ptr is unique per event; excluding it keeps the signature about
		// the event content rather than its storage address.
		if k == "@ptr" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(row.fields[k])
		b.WriteByte(0)
	}
	return b.String()
}

// --- filldown ---

type filldownCommand struct {
	patterns []string // empty means all fields
	headTk   token
}

func (c *filldownCommand) name() string { return "filldown" }

func parseFilldownCommand(args []token, head token) (command, error) {
	c := &filldownCommand{headTk: head}
	for _, t := range args {
		if t.kind != tokIdent && t.kind != tokBacktickIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
		c.patterns = append(c.patterns, t.text)
	}
	return c, nil
}

func (c *filldownCommand) matches(field string) bool {
	if len(c.patterns) == 0 {
		return true
	}
	for _, p := range c.patterns {
		if globMatch(p, field) || p == field {
			return true
		}
	}
	return false
}

func (c *filldownCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	last := make(map[string]string)
	for i := range rows {
		for k, v := range rows[i].fields {
			if v != "" && c.matches(k) {
				last[k] = v
			}
		}
		for k, v := range last {
			if rows[i].fields[k] == "" && c.matches(k) {
				rows[i].set(k, v)
			}
		}
	}
	return rows
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

// --- expand / unnest ---

type expandCommand struct {
	field string
	head  token
}

func (c *expandCommand) name() string { return "expand" }

func parseExpandCommand(args []token, head token) (command, error) {
	if len(args) != 1 || (args[0].kind != tokIdent && args[0].kind != tokBacktickIdent) {
		return nil, newQueryCompileError("Syntax error: expand requires a single field", head.start, head.end)
	}
	return &expandCommand{field: args[0].text, head: head}, nil
}

func (c *expandCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	out := make([]queryResultRow, 0, len(rows))
	for i := range rows {
		raw := readField(rows[i], c.field, ctx)
		if raw == "" {
			out = append(out, rows[i])
			continue
		}
		var decoded interface{}
		if err := jsonUnmarshalString(raw, &decoded); err != nil {
			out = append(out, rows[i])
			continue
		}
		list, ok := decoded.([]interface{})
		if !ok {
			out = append(out, rows[i])
			continue
		}
		for _, item := range list {
			clone := cloneRow(rows[i])
			clone.set(c.field, storeValue(item))
			out = append(out, clone)
		}
	}
	return out
}

type unnestCommand struct {
	field  string
	target string
	head   token
}

func (c *unnestCommand) name() string { return "unnest" }

func parseUnnestCommand(args []token, head token) (command, error) {
	if len(args) < 1 || (args[0].kind != tokIdent && args[0].kind != tokBacktickIdent) {
		return nil, newQueryCompileError("Syntax error: unnest requires a field", head.start, head.end)
	}
	c := &unnestCommand{field: args[0].text, head: head, target: args[0].text}
	if len(args) >= 3 && args[1].kind == tokIdent && strings.EqualFold(args[1].text, "into") {
		if args[2].kind != tokIdent && args[2].kind != tokBacktickIdent {
			return nil, newQueryCompileError("Syntax error: into requires a name", head.start, head.end)
		}
		c.target = args[2].text
	}
	return c, nil
}

func (c *unnestCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	out := make([]queryResultRow, 0, len(rows))
	for i := range rows {
		raw := readField(rows[i], c.field, ctx)
		if raw == "" {
			out = append(out, rows[i])
			continue
		}
		var decoded interface{}
		if err := jsonUnmarshalString(raw, &decoded); err != nil {
			out = append(out, unnestRow(rows[i], c.target, raw))
			continue
		}
		list, ok := decoded.([]interface{})
		if !ok {
			// Non-list values are treated as a list with one item.
			out = append(out, unnestRow(rows[i], c.target, raw))
			continue
		}
		for _, item := range list {
			out = append(out, unnestRow(rows[i], c.target, storeValue(item)))
		}
	}
	return out
}

func unnestRow(src queryResultRow, target, value string) queryResultRow {
	clone := cloneRow(src)
	clone.set(target, value)
	return clone
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
			ms, ok := parsePeriodMillis(args[i+1].text)
			if !ok {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid maxspan '%s'", args[i+1].text), args[i+1].start, args[i+1].end)
			}
			c.maxSpan = ms
			i++
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

// --- unmask ---

type unmaskCommand struct {
	cmdName string
	headTk  token
}

func (c *unmaskCommand) name() string { return c.cmdName }

// apply is the identity transform: the platform does not mask log event
// content, so unmask returns the full content exactly as stored.
func (c *unmaskCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	return rows
}
