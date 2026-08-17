package cloudwatchlogs

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Row-shaping pipeline commands: projection, filter, sort, dedup, filldown,
// expand, and unnest.

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

func jsonUnmarshalString(s string, out interface{}) error {
	return json.Unmarshal([]byte(strings.TrimSpace(s)), out)
}
