package cloudwatchlogs

import (
	"fmt"
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
