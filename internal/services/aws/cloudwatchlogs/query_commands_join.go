package cloudwatchlogs

import (
	"fmt"
	"strconv"
	"strings"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The secondary-source commands: join (combining the pipeline's rows
// with a second log group's rows on a shared key) and appendcols
// (appending a subquery's columns onto the current rows).

// --- join ---

// joinCommand correlates the current rows with a secondary log group. Only
// equality joins are supported and only one join per query, per the
// documented limitations.
type joinCommand struct {
	joinType   string // inner or left
	leftAlias  string
	rightAlias string
	leftKey    string
	rightKey   string
	source     *sourceCommand
	headTk     token
}

func (c *joinCommand) name() string { return "join" }

func parseJoinCommand(args []token, head token) (command, error) {
	c := &joinCommand{joinType: "inner", leftAlias: "left", rightAlias: "right", headTk: head}
	// The right-side source is the parenthesised trailing group.
	depth := 0
	openIdx := -1
	for i, t := range args {
		switch t.kind {
		case tokLParen:
			if depth == 0 {
				openIdx = i
			}
			depth++
		case tokRParen:
			depth--
			if depth == 0 {
				if openIdx < 0 || i+1 < len(args) {
					return nil, newQueryCompileError("Syntax error in join source", head.start, head.end)
				}
				src, err := parseSourceCommand(args[openIdx+2:i], head)
				if err != nil {
					return nil, err
				}
				sc, ok := src.(*sourceCommand)
				if !ok {
					return nil, newQueryCompileError("Syntax error: join source must be a SOURCE selection", head.start, head.end)
				}
				c.source = sc
				args = args[:openIdx]
			}
		}
	}
	i := 0
	for i < len(args) {
		t := args[i]
		if t.kind != tokIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
		lower := strings.ToLower(t.text)
		// The "=" guard covers its value too: a trailing `join type =`
		// with nothing after the equals must reject as a syntax error,
		// never read past the slice end.
		isKV := i+2 < len(args) && args[i+1].kind == tokOp && args[i+1].text == "="
		switch {
		case lower == "type" && isKV:
			c.joinType = strings.ToLower(args[i+2].text)
			if c.joinType != "inner" && c.joinType != "left" {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid join type '%s'", c.joinType), args[i+2].start, args[i+2].end)
			}
			i += 3
		case lower == "left" && isKV:
			c.leftAlias = args[i+2].text
			i += 3
		case lower == "right" && isKV:
			c.rightAlias = args[i+2].text
			i += 3
		case strings.EqualFold(t.text, "where") && i+3 < len(args):
			lk, ok1 := parseDottedKey(args[i+1])
			if !ok1 || args[i+2].kind != tokOp || args[i+2].text != "=" {
				return nil, newQueryCompileError("Syntax error: join where requires field=field", t.start, t.end)
			}
			rk, ok2 := parseDottedKey(args[i+3])
			if !ok2 {
				return nil, newQueryCompileError("Syntax error: join where requires field=field", t.start, t.end)
			}
			c.leftKey, c.rightKey = lk, rk
			i += 4
		default:
			return nil, newQueryCompileError(
				fmt.Sprintf("Syntax error at '%s': unexpected token in join", t.raw), t.start, t.end)
		}
	}
	if c.leftKey == "" || c.rightKey == "" {
		return nil, newQueryCompileError("Syntax error: join requires a where key condition", head.start, head.end)
	}
	if c.source == nil {
		return nil, newQueryCompileError(
			"Syntax error: join requires a parenthesised SOURCE selection for the right side", head.start, head.end)
	}
	return c, nil
}

// parseDottedKey parses alias.field references in join conditions.
func parseDottedKey(t token) (string, bool) {
	if t.kind != tokIdent || !strings.Contains(t.text, ".") {
		return "", false
	}
	return t.text, true
}

func (c *joinCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	groups, err := c.source.resolve(ctx)
	if err != nil {
		ctx.sourceError = err
		return rows
	}
	rightEvents, err := ctx.fetchEvents(groups, ctx.startTime, ctx.endTime)
	if err != nil {
		ctx.sourceError = err
		return rows
	}
	// The right side is evaluated over the secondary source's raw events.
	rightRows := buildRows(rightEvents, ctx.accountID)
	leftField := aliasField(c.leftKey)
	rightField := aliasField(c.rightKey)

	index := make(map[string][]queryResultRow)
	for i := range rightRows {
		key := rightRows[i].fields[rightField]
		if key == "" {
			continue
		}
		index[key] = append(index[key], rightRows[i])
		// "The number of unique key values in the secondary data source
		// is limited to 50,000" — a right side beyond the cap fails the
		// query rather than silently joining a truncated index.
		if len(index) > logsstore.MaxJoinKeyValues {
			ctx.sourceError = fmt.Errorf(
				"the number of unique key values in the secondary data source is limited to %d",
				logsstore.MaxJoinKeyValues)
			return rows
		}
	}
	var out []queryResultRow
	for i := range rows {
		key := rows[i].fields[leftField]
		matches := index[key]
		if len(matches) == 0 {
			if c.joinType == "left" {
				out = append(out, mergedJoinRow(rows[i], c.leftAlias, queryResultRow{}, c.rightAlias))
			}
			continue
		}
		for _, m := range matches {
			out = append(out, mergedJoinRow(rows[i], c.leftAlias, m, c.rightAlias))
		}
	}
	if out == nil {
		out = []queryResultRow{}
	}
	return out
}

// aliasField strips the alias prefix from alias.field references.
func aliasField(key string) string {
	dot := strings.Index(key, ".")
	if dot < 0 {
		return key
	}
	return key[dot+1:]
}

// mergedJoinRow combines a left and a right row, prefixing both sides'
// fields with their aliases so later commands can reference
// alias.field names. Left columns precede right columns.
func mergedJoinRow(left queryResultRow, leftAlias string, right queryResultRow, rightAlias string) queryResultRow {
	merged := queryResultRow{}
	for _, k := range left.ordered() {
		merged.set(leftAlias+"."+k, left.fields[k])
	}
	for _, k := range right.ordered() {
		merged.set(rightAlias+"."+k, right.fields[k])
	}
	return merged
}

// --- appendcols ---

// appendcolsCommand appends a sub-query's columns to the current results by
// positional row matching.
type appendcolsCommand struct {
	override bool
	max      int
	toks     []token
	headTk   token
}

func (c *appendcolsCommand) name() string { return "appendcols" }

func parseAppendColsCommand(args []token, head token) (command, error) {
	c := &appendcolsCommand{headTk: head, max: 10000}
	// The sub-query is the trailing parenthesised group; flags precede it.
	depth := 0
	openIdx := -1
	closeIdx := -1
	for i := 0; i < len(args); i++ {
		switch args[i].kind {
		case tokLParen:
			if depth == 0 {
				openIdx = i
			}
			depth++
		case tokRParen:
			depth--
			if depth == 0 {
				closeIdx = i
			}
		}
		if closeIdx >= 0 {
			break
		}
	}
	if openIdx < 0 || closeIdx != len(args)-1 {
		return nil, newQueryCompileError("Syntax error: appendcols requires a trailing subquery", head.start, head.end)
	}
	c.toks = args[openIdx+1 : closeIdx]
	if len(c.toks) == 0 {
		return nil, newQueryCompileError("Syntax error: appendcols subquery is empty", head.start, head.end)
	}
	flags := args[:openIdx]
	i := 0
	for i < len(flags) {
		t := flags[i]
		if t.kind != tokIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
		if i+2 >= len(flags) || flags[i+1].kind != tokOp || flags[i+1].text != "=" {
			return nil, newQueryCompileError("Syntax error: appendcols flags use key=value", t.start, t.end)
		}
		val := flags[i+2]
		switch strings.ToLower(t.text) {
		case "override":
			c.override = strings.EqualFold(val.text, "true")
		case "max":
			n, err := strconv.Atoi(val.text)
			if err != nil || n < 1 || n > logsstore.MaxQueryLanguageLimit {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid max '%s'", val.text), val.start, val.end)
			}
			c.max = n
		default:
			return nil, newQueryCompileError(
				fmt.Sprintf("Syntax error at '%s': unknown appendcols flag", t.raw), t.start, t.end)
		}
		i += 3
	}
	return c, nil
}

func (c *appendcolsCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	subRows, err := ctx.runSubquery(c.toks)
	if err != nil {
		// A failed subquery side fails the query — appending nothing
		// silently would report a result set the query never computed.
		ctx.sourceError = err
		return rows
	}
	n := len(rows)
	if n > c.max {
		n = c.max
	}
	for i := 0; i < n; i++ {
		if i >= len(subRows) {
			break
		}
		for _, k := range subRows[i].ordered() {
			if _, exists := rows[i].fields[k]; exists && !c.override {
				continue
			}
			rows[i].set(k, subRows[i].fields[k])
		}
	}
	return rows
}
