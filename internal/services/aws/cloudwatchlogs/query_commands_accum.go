package cloudwatchlogs

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// The row-window value commands: accum, autoregress, and addtotals —
// each derives per-row values from the neighbouring rows or the column
// totals of the pipeline's current row set.

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
