package cloudwatchlogs

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// pathSeg is one segment of a field path: either a map key or a list index.
// Structure access follows the documented dot and bracket operators, e.g.
// json_message.users[1].action.
type pathSeg struct {
	name  string
	index int
	isIdx bool
}

// exprNode is a parsed expression. Commands parse their argument expressions
// once and evaluate the AST per row.
type exprNode interface {
	eval(row *queryResultRow, ctx *execContext) interface{}
}

type literalNode struct {
	val interface{}
}

type fieldNode struct {
	segs []pathSeg
}

type binOpNode struct {
	op   string
	l, r exprNode
}

type unaryNode struct {
	op string
	x  exprNode
}

type inNode struct {
	x    exprNode
	list []exprNode
	neg  bool
}

type likeNode struct {
	x       exprNode
	pattern string
	re      *regexp.Regexp
	neg     bool
}

type funcNode struct {
	name string
	args []exprNode
}

// exprParser is a recursive-descent parser over a token slice.
type exprParser struct {
	toks []token
	pos  int
}

func parseExprTokens(toks []token) (exprNode, error) {
	p := &exprParser{toks: toks}
	n, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if p.pos != len(p.toks) {
		t := p.toks[p.pos]
		return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
	}
	return n, nil
}

func (p *exprParser) peek() (token, bool) {
	if p.pos < len(p.toks) {
		return p.toks[p.pos], true
	}
	return token{}, false
}

func (p *exprParser) peekAt(n int) (token, bool) {
	if p.pos+n < len(p.toks) {
		return p.toks[p.pos+n], true
	}
	return token{}, false
}

func (p *exprParser) next() (token, bool) {
	t, ok := p.peek()
	if ok {
		p.pos++
	}
	return t, ok
}

func (p *exprParser) isKeyword(kw string) bool {
	t, ok := p.peek()
	return ok && t.kind == tokIdent && strings.EqualFold(t.text, kw)
}

func (p *exprParser) acceptKeyword(kw string) bool {
	if p.isKeyword(kw) {
		p.pos++
		return true
	}
	return false
}

func (p *exprParser) parseOr() (exprNode, error) {
	l, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.acceptKeyword("or") {
		r, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		l = &binOpNode{op: "or", l: l, r: r}
	}
	return l, nil
}

func (p *exprParser) parseAnd() (exprNode, error) {
	l, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.acceptKeyword("and") {
		r, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		l = &binOpNode{op: "and", l: l, r: r}
	}
	return l, nil
}

func (p *exprParser) parseNot() (exprNode, error) {
	if p.acceptKeyword("not") {
		x, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &unaryNode{op: "not", x: x}, nil
	}
	return p.parseComparison()
}

func (p *exprParser) parseComparison() (exprNode, error) {
	l, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}
	t, ok := p.peek()
	if !ok {
		return l, nil
	}
	switch {
	case t.kind == tokOp && (t.text == "=" || t.text == "==" || t.text == "!=" ||
		t.text == "<" || t.text == "<=" || t.text == ">" || t.text == ">="):
		p.pos++
		r, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}
		op := t.text
		if op == "==" {
			op = "="
		}
		return &binOpNode{op: op, l: l, r: r}, nil
	case t.kind == tokIdent && strings.EqualFold(t.text, "like"):
		p.pos++
		return p.parseLikeTail(l, false)
	case t.kind == tokIdent && strings.EqualFold(t.text, "not"):
		// not like / not in
		p.pos++
		if t2, ok := p.peek(); ok && t2.kind == tokIdent && strings.EqualFold(t2.text, "like") {
			p.pos++
			return p.parseLikeTail(l, true)
		}
		if t2, ok := p.peek(); ok && t2.kind == tokIdent && strings.EqualFold(t2.text, "in") {
			p.pos++
			return p.parseInTail(l, true)
		}
		return l, nil
	case t.kind == tokIdent && strings.EqualFold(t.text, "in"):
		p.pos++
		return p.parseInTail(l, false)
	}
	return l, nil
}

// parseLikeTail parses the pattern after the like keyword. A pattern enclosed
// in forward slashes is a regular expression; a quoted pattern is a wildcard
// glob.
func (p *exprParser) parseLikeTail(x exprNode, neg bool) (exprNode, error) {
	t, ok := p.next()
	if !ok {
		// Report at the end of the last token rather than the zero-value
		// token, so the offset is inside the query string.
		last := p.toks[len(p.toks)-1]
		return nil, newQueryCompileError("Syntax error: like requires a pattern", last.end, last.end)
	}
	switch t.kind {
	case tokString:
		return &likeNode{x: x, pattern: t.text, neg: neg}, nil
	case tokRegex:
		re, err := regexp.Compile(t.text)
		if err != nil {
			return nil, newQueryCompileError(fmt.Sprintf("Invalid regular expression: %v", err), t.start, t.end)
		}
		return &likeNode{x: x, pattern: t.text, re: re, neg: neg}, nil
	case tokOp:
		if t.text == "/" {
			startTok := t
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
				return nil, newQueryCompileError("Unterminated regular expression", startTok.start, startTok.start+1)
			}
			re, err := regexp.Compile(raw.String())
			if err != nil {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid regular expression: %v", err), startTok.start, closing.end)
			}
			return &likeNode{x: x, pattern: raw.String(), re: re, neg: neg}, nil
		}
	}
	return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
}

func (p *exprParser) parseInTail(x exprNode, neg bool) (exprNode, error) {
	t, ok := p.next()
	if !ok {
		last := p.toks[len(p.toks)-1]
		return nil, newQueryCompileError("Syntax error: in requires a list", last.end, last.end)
	}
	if t.kind == tokLParen {
		// Literal lists use brackets; a parenthesised operand after in is
		// always a nested query. Invalid interiors fail compilation with a
		// positioned syntax error.
		inner, err := p.collectParenInterior(t)
		if err != nil {
			return nil, err
		}
		sub := &subqueryNode{toks: inner, startTk: t}
		return &inNode{x: x, list: []exprNode{sub}, neg: neg}, nil
	}
	if t.kind != tokLBracket {
		return nil, newQueryCompileError("Syntax error: in requires a list", t.start, t.end)
	}
	var list []exprNode
	if t2, ok := p.peek(); ok && t2.kind == tokRBracket {
		p.pos++
		return &inNode{x: x, list: list, neg: neg}, nil
	}
	for {
		e, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		list = append(list, e)
		t2, ok := p.next()
		if !ok {
			last := p.toks[len(p.toks)-1]
			return nil, newQueryCompileError("Syntax error: unterminated list", last.end, last.end)
		}
		if t2.kind == tokRBracket {
			break
		}
		if t2.kind != tokComma {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t2.raw), t2.start, t2.end)
		}
	}
	return &inNode{x: x, list: list, neg: neg}, nil
}

// collectParenInterior consumes tokens up to the parenthesis matching the
// already-consumed opening token and returns the interior.
func (p *exprParser) collectParenInterior(open token) ([]token, error) {
	depth := 1
	var inner []token
	for {
		t, ok := p.next()
		if !ok {
			return nil, newQueryCompileError("Syntax error: missing ')'", open.start, open.end)
		}
		switch t.kind {
		case tokLParen, tokLBracket:
			depth++
		case tokRParen, tokRBracket:
			depth--
			if depth == 0 {
				return inner, nil
			}
		}
		inner = append(inner, t)
	}
}

func (p *exprParser) parseAdditive() (exprNode, error) {
	l, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}
	for {
		t, ok := p.peek()
		if !ok || t.kind != tokOp || (t.text != "+" && t.text != "-") {
			return l, nil
		}
		p.pos++
		r, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		l = &binOpNode{op: t.text, l: l, r: r}
	}
}

func (p *exprParser) parseMultiplicative() (exprNode, error) {
	l, err := p.parsePower()
	if err != nil {
		return nil, err
	}
	for {
		t, ok := p.peek()
		if !ok || t.kind != tokOp || (t.text != "*" && t.text != "/" && t.text != "%") {
			return l, nil
		}
		p.pos++
		r, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		l = &binOpNode{op: t.text, l: l, r: r}
	}
}

func (p *exprParser) parsePower() (exprNode, error) {
	l, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	t, ok := p.peek()
	if ok && t.kind == tokOp && t.text == "^" {
		p.pos++
		r, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		return &binOpNode{op: "^", l: l, r: r}, nil
	}
	return l, nil
}

func (p *exprParser) parseUnary() (exprNode, error) {
	t, ok := p.peek()
	if ok && t.kind == tokOp && t.text == "-" {
		p.pos++
		x, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &unaryNode{op: "-", x: x}, nil
	}
	return p.parsePrimary()
}

func (p *exprParser) parsePrimary() (exprNode, error) {
	t, ok := p.next()
	if !ok {
		return nil, newQueryCompileError("Syntax error: unexpected end of expression", 0, 0)
	}
	switch t.kind {
	case tokNumber:
		n, err := strconv.ParseFloat(t.text, 64)
		if err != nil {
			return nil, newQueryCompileError(fmt.Sprintf("Invalid number '%s'", t.text), t.start, t.end)
		}
		return &literalNode{val: n}, nil
	case tokString:
		return &literalNode{val: t.text}, nil
	case tokOp:
		// count(*) uses the bare star as its argument.
		if t.text == "*" {
			return &literalNode{val: "*"}, nil
		}
	case tokLParen:
		e, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		t2, ok := p.next()
		if !ok || t2.kind != tokRParen {
			if !ok {
				return nil, newQueryCompileError("Syntax error: missing ')'", t.start, t.end)
			}
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s': expected ')'", t2.raw), t2.start, t2.end)
		}
		return e, nil
	case tokLBracket:
		// Array literal, used by in [...] elements such as strings.
		var items []exprNode
		if t2, ok := p.peek(); ok && t2.kind == tokRBracket {
			p.pos++
			return &arrayLiteralNode{items: items}, nil
		}
		for {
			e, err := p.parseOr()
			if err != nil {
				return nil, err
			}
			items = append(items, e)
			t2, ok := p.next()
			if !ok {
				return nil, newQueryCompileError("Syntax error: unterminated list", t.start, t.end)
			}
			if t2.kind == tokRBracket {
				break
			}
			if t2.kind != tokComma {
				return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t2.raw), t2.start, t2.end)
			}
		}
		return &arrayLiteralNode{items: items}, nil
	case tokIdent, tokBacktickIdent:
		// Function call or field path.
		if t.kind == tokIdent {
			if t2, ok := p.peek(); ok && t2.kind == tokLParen && !strings.Contains(t.text, ".") {
				fn, err := p.parseFuncCall(t)
				if err != nil {
					return nil, err
				}
				return p.parsePostfix(fn)
			}
		}
		fld, err := p.parseFieldPath(t)
		if err != nil {
			return nil, err
		}
		return p.parsePostfix(fld)
	}
	return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
}

// parsePostfix wraps a primary expression with .field and [index] access
// chains, e.g. jsonParse(@message).users[1].action.
func (p *exprParser) parsePostfix(e exprNode) (exprNode, error) {
	for {
		t, ok := p.peek()
		if !ok {
			return e, nil
		}
		if t.kind == tokLBracket {
			p.pos++
			idx, ok := p.next()
			if !ok || idx.kind != tokNumber {
				return nil, newQueryCompileError("Syntax error: bracket access requires an index", t.start, t.end)
			}
			n, err := strconv.Atoi(idx.text)
			if err != nil {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid index '%s'", idx.text), idx.start, idx.end)
			}
			t2, ok := p.next()
			if !ok || t2.kind != tokRBracket {
				return nil, newQueryCompileError("Syntax error: missing ']'", t.start, t.end)
			}
			e = &accessNode{base: e, segs: []pathSeg{{index: n, isIdx: true}}}
			continue
		}
		if t.kind == tokOp && t.text == "." {
			p.pos++
			name, ok := p.next()
			if !ok || (name.kind != tokIdent && name.kind != tokBacktickIdent) {
				return nil, newQueryCompileError("Syntax error: expected a field name after '.'", t.start, t.end)
			}
			e = &accessNode{base: e, segs: []pathSeg{{name: name.text}}}
			continue
		}
		if t.kind == tokBacktickIdent {
			p.pos++
			e = &accessNode{base: e, segs: []pathSeg{{name: t.text}}}
			continue
		}
		return e, nil
	}
}

// accessNode accesses structure attributes on the result of an arbitrary
// expression, e.g. jsonParse(@message).status.
type accessNode struct {
	base exprNode
	segs []pathSeg
}

func (n *accessNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	return descendPath(n.base.eval(row, ctx), n.segs)
}

type arrayLiteralNode struct {
	items []exprNode
}

func (n *arrayLiteralNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	out := make([]interface{}, 0, len(n.items))
	for _, e := range n.items {
		out = append(out, e.eval(row, ctx))
	}
	return out
}

// queryFunctionArities records the accepted argument counts for every
// documented function, so calls with a wrong number of arguments are
// rejected at compile time with a positioned error instead of panicking or
// silently misbehaving at evaluation time. The pair is {min, max} with a
// max of -1 meaning unbounded.
var queryFunctionArities = map[string][2]int{
	// General functions.
	"ispresent": {1, 1}, "ispresentornull": {1, 1}, "coalesce": {1, -1},
	"case": {2, -1}, "if": {3, 3}, "isnumeric": {1, 1},
	"messagesize": {1, 1}, "querystarttime": {0, 0}, "queryendtime": {0, 0}, "querytimerange": {0, 0},
	// Hashing functions.
	"md5": {1, 1}, "sha256": {1, 1},
	// Numeric functions.
	"abs": {1, 1}, "ceil": {1, 1}, "floor": {1, 1}, "greatest": {1, -1}, "least": {1, -1},
	"log": {1, 1}, "round": {1, 2}, "sqrt": {1, 1}, "haversine": {4, 4},
	"tonumber": {1, 1}, "toint": {1, 1}, "tolong": {1, 1}, "todouble": {1, 1},
	// String functions.
	"isempty": {1, 1}, "isblank": {1, 1}, "concat": {1, -1},
	"ltrim": {1, 2}, "rtrim": {1, 2}, "trim": {1, 2}, "strlen": {1, 1},
	"toupper": {1, 1}, "tolower": {1, 1}, "substr": {2, 3},
	"replace": {3, 3}, "regexreplace": {3, 3}, "strcontains": {2, 3},
	"startswith": {2, 2}, "endswith": {2, 2}, "urlencode": {1, 1}, "urldecode": {1, 1},
	"base64encode": {1, 1}, "base64decode": {1, 1}, "split": {2, 2},
	"hextoascii": {1, 1}, "hextodec": {1, 1}, "dectohex": {1, 1},
	// IP address functions.
	"isvalidip": {1, 1}, "isvalidipv4": {1, 1}, "isvalidipv6": {1, 1},
	"isipinsubnet": {2, 2}, "isipv4insubnet": {2, 2}, "isipv6insubnet": {2, 2},
	"ipv4tonumber": {1, 1}, "isprivateip": {1, 1}, "ispublicip": {1, 1}, "isreservedip": {1, 1},
	// JSON functions.
	"jsonparse": {1, 1}, "jsonstringify": {1, 1}, "jsonarraysize": {1, 1}, "jsonarraycontains": {2, 2},
	// Datetime functions.
	"datefloor": {2, 2}, "dateceil": {2, 2}, "frommillis": {1, 1}, "tomillis": {1, 1},
	"now": {0, 0}, "parsedate": {2, 3}, "formatdate": {2, 3}, "strftime": {2, 3},
	// Aggregation functions (stats context).
	"avg": {1, 1}, "count": {0, 1}, "countdistinct": {1, 1}, "max": {1, 1}, "min": {1, 1},
	"pct": {2, 2}, "stddev": {1, 1}, "sum": {1, 1}, "values": {1, 1}, "collect_values": {1, 1},
	"variance": {1, 1}, "topk": {2, 2}, "earliest": {1, 1}, "latest": {1, 1},
	"sortsfirst": {1, 1}, "sortslast": {1, 1},
	// Time-series functions (stats context).
	"rate": {2, 2}, "countovertime": {1, 1}, "sumovertime": {1, 1}, "histogram": {2, 2},
}

// checkFunctionArity rejects calls whose argument count falls outside the
// documented range. Functions absent from the table are not constrained.
func checkFunctionArity(name token, args []exprNode) error {
	arity, ok := queryFunctionArities[strings.ToLower(name.text)]
	if !ok {
		return nil
	}
	n := len(args)
	if n < arity[0] || (arity[1] >= 0 && n > arity[1]) {
		if arity[1] < 0 {
			return newQueryCompileError(
				fmt.Sprintf("Function %s requires at least %d argument(s)", name.text, arity[0]), name.start, name.end)
		}
		if arity[0] == arity[1] {
			return newQueryCompileError(
				fmt.Sprintf("Function %s requires exactly %d argument(s)", name.text, arity[0]), name.start, name.end)
		}
		return newQueryCompileError(
			fmt.Sprintf("Function %s requires between %d and %d arguments", name.text, arity[0], arity[1]), name.start, name.end)
	}
	// case accepts an unbounded argument count but at most ten
	// condition/value branches, with an optional trailing default.
	if strings.ToLower(name.text) == "case" {
		branches := n / 2
		if n%2 == 1 {
			branches = (n - 1) / 2
		}
		if branches > 10 {
			return newQueryCompileError(
				"Function case supports up to 10 branches", name.start, name.end)
		}
	}
	return nil
}

// isPeriodUnitWord reports whether the identifier text is one of the
// documented period unit words, used to join a number and unit written
// without separation (e.g. 5m) into one period argument.
func isPeriodUnitWord(s string) bool {
	switch strings.ToLower(s) {
	case "ms", "msec", "msecs", "millisecond", "milliseconds",
		"s", "sec", "secs", "second", "seconds",
		"m", "min", "mins", "minute", "minutes",
		"h", "hr", "hrs", "hour", "hours",
		"d", "day", "days", "w", "week", "weeks",
		"mo", "mon", "mons", "month", "months",
		"q", "qtr", "qtrs", "quarter", "quarters",
		"y", "yr", "yrs", "year", "years":
		return true
	}
	return false
}

// parsePeriodArg parses one function argument, folding a number followed
// directly by a period unit word into a single period literal node.
func (p *exprParser) parsePeriodArg() (exprNode, error) {
	t, ok := p.peek()
	if !ok || t.kind != tokNumber {
		return p.parseOr()
	}
	if u, ok := p.peekAt(1); ok && u.kind == tokIdent && isPeriodUnitWord(u.text) {
		p.pos += 2
		return &literalNode{val: t.raw + u.text}, nil
	}
	return p.parseOr()
}

func (p *exprParser) parseFuncCall(name token) (exprNode, error) {
	p.pos++ // consume '('
	var args []exprNode
	if t, ok := p.peek(); ok && t.kind == tokRParen {
		p.pos++
		if err := checkFunctionArity(name, args); err != nil {
			return nil, err
		}
		return &funcNode{name: name.text, args: args}, nil
	}
	for {
		e, err := p.parsePeriodArg()
		if err != nil {
			return nil, err
		}
		args = append(args, e)
		t, ok := p.next()
		if !ok {
			return nil, newQueryCompileError("Syntax error: unterminated function call", name.start, name.end)
		}
		if t.kind == tokRParen {
			break
		}
		if t.kind != tokComma {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
	}
	if err := checkFunctionArity(name, args); err != nil {
		return nil, err
	}
	return &funcNode{name: name.text, args: args}, nil
}

// parseFieldPath parses a field reference with dot and bracket access. The
// first token may already contain dot-separated segments; bracket indexing
// and post-bracket dot continuation extend the path.
func (p *exprParser) parseFieldPath(first token) (exprNode, error) {
	node := &fieldNode{}
	appendSegs := func(text string) {
		for _, part := range strings.Split(text, ".") {
			if part != "" {
				node.segs = append(node.segs, pathSeg{name: part})
			}
		}
	}
	if first.kind == tokBacktickIdent {
		node.segs = append(node.segs, pathSeg{name: first.text})
	} else {
		appendSegs(first.text)
	}
	if len(node.segs) == 0 {
		return nil, newQueryCompileError("Syntax error: empty field name", first.start, first.end)
	}
	for {
		t, ok := p.peek()
		if !ok {
			break
		}
		if t.kind == tokLBracket {
			p.pos++
			idx, ok := p.next()
			if !ok || idx.kind != tokNumber {
				return nil, newQueryCompileError("Syntax error: bracket access requires an index", t.start, t.end)
			}
			n, err := strconv.Atoi(idx.text)
			if err != nil {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid index '%s'", idx.text), idx.start, idx.end)
			}
			node.segs = append(node.segs, pathSeg{index: n, isIdx: true})
			t2, ok := p.next()
			if !ok || t2.kind != tokRBracket {
				return nil, newQueryCompileError("Syntax error: missing ']'", t.start, t.end)
			}
			continue
		}
		if t.kind == tokOp && t.text == "." {
			p.pos++
			name, ok := p.next()
			if !ok || (name.kind != tokIdent && name.kind != tokBacktickIdent) {
				return nil, newQueryCompileError("Syntax error: expected a field name after '.'", t.start, t.end)
			}
			if name.kind == tokBacktickIdent {
				node.segs = append(node.segs, pathSeg{name: name.text})
			} else {
				appendSegs(name.text)
			}
			continue
		}
		if t.kind == tokBacktickIdent {
			p.pos++
			node.segs = append(node.segs, pathSeg{name: t.text})
			continue
		}
		break
	}
	return node, nil
}

// --- evaluation ---

func (n *literalNode) eval(row *queryResultRow, ctx *execContext) interface{} { return n.val }

func (n *fieldNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	var cur interface{}
	name := n.segs[0].name
	raw, ok := (*row).fields[name]
	if !ok {
		// Fields without the @ prefix resolve against the @-prefixed
		// discoverable fields, e.g. timestamp matches @timestamp.
		raw, ok = (*row).fields["@"+name]
	}
	if !ok {
		return nil
	}
	cur = raw
	return descendPath(cur, n.segs[1:])
}

// descendPath walks structure segments: maps by key, lists by index.
// String values read from rows are decoded from JSON on demand so that
// jsonParse results round-trip through the string-typed rows.
func descendPath(start interface{}, segs []pathSeg) interface{} {
	cur := start
	for _, seg := range segs {
		if cur == nil {
			return nil
		}
		if s, isStr := cur.(string); isStr {
			trimmed := strings.TrimSpace(s)
			if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
				var decoded interface{}
				if err := json.Unmarshal([]byte(trimmed), &decoded); err == nil {
					cur = decoded
				}
			}
		}
		if seg.isIdx {
			list, ok := cur.([]interface{})
			if !ok || seg.index < 0 || seg.index >= len(list) {
				return nil
			}
			cur = list[seg.index]
			continue
		}
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		v, ok := m[seg.name]
		if !ok {
			return nil
		}
		cur = v
	}
	return cur
}

func (n *unaryNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	v := n.x.eval(row, ctx)
	switch n.op {
	case "not":
		return !truthy(v)
	case "-":
		if f, ok := asNumber(v); ok {
			return -f
		}
		return nil
	}
	return nil
}

func (n *binOpNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	switch n.op {
	case "and":
		return truthy(n.l.eval(row, ctx)) && truthy(n.r.eval(row, ctx))
	case "or":
		return truthy(n.l.eval(row, ctx)) || truthy(n.r.eval(row, ctx))
	}
	l := n.l.eval(row, ctx)
	r := n.r.eval(row, ctx)
	switch n.op {
	case "=", "!=":
		eq := valuesEqual(l, r)
		if n.op == "=" {
			return eq
		}
		return !eq
	case "<", "<=", ">", ">=":
		return compareValues(l, r, n.op)
	case "+", "-", "*", "/", "%", "^":
		// String concatenation for + when either side is a non-numeric
		// string, numeric arithmetic otherwise.
		ln, lok := asNumber(l)
		rn, rok := asNumber(r)
		if !lok || !rok {
			if n.op == "+" {
				return asString(l) + asString(r)
			}
			return nil
		}
		return applyArith(n.op, ln, rn)
	}
	return nil
}

func applyArith(op string, l, r float64) interface{} {
	switch op {
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		if r == 0 {
			return nil
		}
		return l / r
	case "%":
		if int64(r) == 0 {
			return nil
		}
		return float64(int64(l) % int64(r))
	case "^":
		return math.Pow(l, r)
	}
	return nil
}

func (n *inNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	x := n.x.eval(row, ctx)
	found := false
	for _, e := range n.list {
		if sub, ok := e.(*subqueryNode); ok {
			// Subquery membership: the nested query materialises a set of
			// values that are compared against the outer field.
			for _, v := range sub.evalSet(ctx) {
				if valuesEqual(x, v) {
					found = true
					break
				}
			}
			if found {
				break
			}
			continue
		}
		if valuesEqual(x, e.eval(row, ctx)) {
			found = true
			break
		}
	}
	if n.neg {
		return !found
	}
	return found
}

func (n *likeNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	v := asString(n.x.eval(row, ctx))
	var matched bool
	if n.re != nil {
		matched = n.re.MatchString(v)
	} else {
		matched = globMatch(n.pattern, v)
	}
	if n.neg {
		return !matched
	}
	return matched
}

func (n *funcNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	return callQueryFunction(strings.ToLower(n.name), n.args, row, ctx)
}

// --- value helpers ---

// timestampValue is a timestamp-typed result value in epoch milliseconds.
// The documented result type of fromMillis, datefloor, dateceil and bin is
// Timestamp: asNumber keeps the numeric interpretation working inside
// expressions while asString renders the timestamp form results present.
type timestampValue int64

func truthy(v interface{}) bool {
	b, ok := v.(bool)
	return ok && b
}

// boolToNum renders a boolean as the documented 1/0 number result.
func boolToNum(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func asNumber(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case timestampValue:
		return float64(x), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(x), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	case bool:
		return 0, false
	}
	return 0, false
}

// formatNumber renders a numeric value the way query results present it.
func formatNumber(f float64) string {
	if f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// asString converts any value to its result-string representation. Maps and
// lists render as canonical JSON, matching how structure values round-trip
// through the string-typed result rows.
func asString(v interface{}) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case float64:
		return formatNumber(x)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case timestampValue:
		return time.UnixMilli(int64(x)).UTC().Format(resultTimestampLayout)
	case []interface{}, map[string]interface{}:
		b, err := json.Marshal(x)
		if err != nil {
			return ""
		}
		return string(b)
	}
	return fmt.Sprintf("%v", v)
}

// storeValue converts an evaluated value to the string stored in a row,
// keeping numbers numeric-looking so later numeric comparisons still work.
func storeValue(v interface{}) string {
	return asString(v)
}

func valuesEqual(l, r interface{}) bool {
	if l == nil || r == nil {
		return false
	}
	if ln, lok := asNumber(l); lok {
		if rn, rok := asNumber(r); rok {
			return ln == rn
		}
	}
	if lb, lok := l.(bool); lok {
		if rb, rok := r.(bool); rok {
			return lb == rb
		}
		return false
	}
	if _, rok := r.(bool); rok {
		return false
	}
	ls, lok := l.(string)
	rs, rok2 := r.(string)
	if lok && rok2 {
		return ls == rs
	}
	return asString(l) == asString(r)
}

func compareValues(l, r interface{}, op string) interface{} {
	if l == nil || r == nil {
		return false
	}
	ln, lok := asNumber(l)
	rn, rok := asNumber(r)
	var cmp int
	if lok && rok {
		switch {
		case ln < rn:
			cmp = -1
		case ln > rn:
			cmp = 1
		}
	} else {
		ls, rs := asString(l), asString(r)
		cmp = strings.Compare(ls, rs)
	}
	switch op {
	case "<":
		return cmp < 0
	case "<=":
		return cmp <= 0
	case ">":
		return cmp > 0
	case ">=":
		return cmp >= 0
	}
	return false
}

// globMatch implements the wildcard matching used by like with quoted
// patterns: * matches any run of characters.
func globMatch(pattern, text string) bool {
	if !strings.Contains(pattern, "*") {
		return pattern == text
	}
	re := regexp.QuoteMeta(pattern)
	re = strings.ReplaceAll(re, `\*`, ".*")
	matched, _ := regexp.MatchString("^"+re+"$", text)
	return matched
}

// --- function library ---
// The functions below implement the documented Logs Insights QL function
// set: general, numeric, string, IP address, JSON, and datetime functions.

func callQueryFunction(name string, args []exprNode, row *queryResultRow, ctx *execContext) interface{} {
	vals := make([]interface{}, len(args))
	for i, a := range args {
		vals[i] = a.eval(row, ctx)
	}
	arg := func(i int) interface{} {
		if i < len(vals) {
			return vals[i]
		}
		return nil
	}
	num := func(i int) (float64, bool) { return asNumber(arg(i)) }
	str := func(i int) string { return asString(arg(i)) }

	switch name {
	// General functions.
	case "ispresent":
		return arg(0) != nil
	case "ispresentornull":
		return arg(0) != nil
	case "coalesce":
		for _, v := range vals {
			if v != nil && asString(v) != "" {
				return v
			}
		}
		return nil
	case "case":
		// case(cond1, val1, ..., [default]): pairs then optional default.
		i := 0
		for i+1 < len(vals) {
			if truthy(vals[i]) {
				return vals[i+1]
			}
			i += 2
		}
		if len(vals)%2 == 1 {
			return vals[len(vals)-1]
		}
		return nil
	case "if":
		// The condition is validated as a three-argument call at parse time;
		// arg bounds defensively so a nil context cannot panic here.
		if truthy(arg(0)) {
			return arg(1)
		}
		return arg(2)
	case "isnumeric":
		_, ok := num(0)
		return ok
	case "messagesize":
		return float64(len(str(0)))
	case "querystarttime":
		if ctx != nil {
			return float64(ctx.startTime)
		}
		return nil
	case "queryendtime":
		if ctx != nil {
			return float64(ctx.endTime)
		}
		return nil
	case "querytimerange":
		if ctx != nil {
			return float64(ctx.endTime - ctx.startTime)
		}
		return nil

	// Numeric functions.
	case "abs":
		if f, ok := num(0); ok {
			return math.Abs(f)
		}
		return nil
	case "ceil":
		if f, ok := num(0); ok {
			return math.Ceil(f)
		}
		return nil
	case "floor":
		if f, ok := num(0); ok {
			return math.Floor(f)
		}
		return nil
	case "greatest":
		best, ok := num(0)
		if !ok {
			return nil
		}
		for i := 1; i < len(vals); i++ {
			if f, ok := num(i); ok && f > best {
				best = f
			}
		}
		return best
	case "least":
		best, ok := num(0)
		if !ok {
			return nil
		}
		for i := 1; i < len(vals); i++ {
			if f, ok := num(i); ok && f < best {
				best = f
			}
		}
		return best
	case "log":
		if f, ok := num(0); ok && f > 0 {
			return math.Log(f)
		}
		return nil
	case "round":
		if f, ok := num(0); ok {
			d := 0.0
			if len(vals) > 1 {
				if d2, ok := num(1); ok {
					d = d2
				}
			}
			pow := math.Pow(10, d)
			return math.Round(f*pow) / pow
		}
		return nil
	case "sqrt":
		if f, ok := num(0); ok && f >= 0 {
			return math.Sqrt(f)
		}
		return nil
	case "haversine":
		lat1, ok1 := num(0)
		lon1, ok2 := num(1)
		lat2, ok3 := num(2)
		lon2, ok4 := num(3)
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return nil
		}
		const toRad = math.Pi / 180
		const earthKm = 6371.0
		dLat := (lat2 - lat1) * toRad
		dLon := (lon2 - lon1) * toRad
		a := math.Sin(dLat/2)*math.Sin(dLat/2) +
			math.Cos(lat1*toRad)*math.Cos(lat2*toRad)*math.Sin(dLon/2)*math.Sin(dLon/2)
		c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
		return earthKm * c
	case "tonumber":
		if f, ok := num(0); ok {
			return f
		}
		return nil
	case "toint":
		if f, ok := num(0); ok {
			return float64(int32(f))
		}
		return nil
	case "tolong":
		if f, ok := num(0); ok {
			return float64(int64(f))
		}
		return nil
	case "todouble":
		if f, ok := num(0); ok {
			return f
		}
		return nil

	// String functions documented to return Number yield 1/0.
	case "isempty":
		return boolToNum(arg(0) == nil || str(0) == "")
	case "isblank":
		return boolToNum(arg(0) == nil || strings.TrimSpace(str(0)) == "")
	case "concat":
		var b strings.Builder
		for i := range vals {
			b.WriteString(str(i))
		}
		return b.String()
	case "ltrim":
		if len(vals) > 1 {
			return strings.TrimLeft(str(0), str(1))
		}
		return strings.TrimLeft(str(0), " \t\r\n")
	case "rtrim":
		if len(vals) > 1 {
			return strings.TrimRight(str(0), str(1))
		}
		return strings.TrimRight(str(0), " \t\r\n")
	case "trim":
		if len(vals) > 1 {
			return strings.Trim(str(0), str(1))
		}
		return strings.Trim(str(0), " \t\r\n")
	case "strlen":
		return float64(len([]rune(str(0))))
	case "toupper":
		return strings.ToUpper(str(0))
	case "tolower":
		return strings.ToLower(str(0))
	case "substr":
		s := []rune(str(0))
		start, ok := num(1)
		if !ok {
			return nil
		}
		si := int(start)
		if si < 0 {
			si = len(s) + si
		}
		if si < 0 || si > len(s) {
			return ""
		}
		if len(vals) > 2 {
			l, ok := num(2)
			if !ok {
				return nil
			}
			ei := si + int(l)
			if ei > len(s) {
				ei = len(s)
			}
			if ei < si {
				return ""
			}
			return string(s[si:ei])
		}
		return string(s[si:])
	case "replace":
		return strings.ReplaceAll(str(0), str(1), str(2))
	case "regexreplace":
		re, err := regexp.Compile(str(1))
		if err != nil {
			return nil
		}
		return re.ReplaceAllString(str(0), str(2))
	case "strcontains":
		if len(vals) > 2 && truthy(vals[2]) {
			return boolToNum(strings.Contains(strings.ToLower(str(0)), strings.ToLower(str(1))))
		}
		return boolToNum(strings.Contains(str(0), str(1)))
	case "startswith":
		return boolToNum(strings.HasPrefix(str(0), str(1)))
	case "endswith":
		return boolToNum(strings.HasSuffix(str(0), str(1)))
	case "urlencode":
		return url.QueryEscape(str(0))
	case "urldecode":
		if dec, err := url.QueryUnescape(str(0)); err == nil {
			return dec
		}
		return nil
	case "base64encode":
		return base64.StdEncoding.EncodeToString([]byte(str(0)))
	case "base64decode":
		if dec, err := base64.StdEncoding.DecodeString(str(0)); err == nil {
			return string(dec)
		}
		return nil
	case "split":
		var parts []interface{}
		for _, p := range strings.Split(str(0), str(1)) {
			parts = append(parts, p)
		}
		return parts
	case "hextoascii":
		if b, err := hexDecodeString(str(0)); err == nil {
			return string(b)
		}
		return nil
	case "hextodec":
		ui, err := strconv.ParseUint(strings.TrimPrefix(strings.ToLower(str(0)), "0x"), 16, 64)
		if err != nil {
			return nil
		}
		return float64(ui)
	case "dectohex":
		f, ok := num(0)
		if !ok {
			return nil
		}
		n := int64(f)
		prefix := "0x"
		if n < 0 {
			prefix = "-0x"
			n = -n
		}
		return prefix + strconv.FormatUint(uint64(n), 16)

	// IP address functions.
	case "isvalidip":
		return net.ParseIP(str(0)) != nil
	case "isvalidipv4":
		ip := net.ParseIP(str(0))
		return ip != nil && ip.To4() != nil && strings.Count(str(0), ".") == 3
	case "isvalidipv6":
		ip := net.ParseIP(str(0))
		return ip != nil && strings.Contains(str(0), ":")
	case "isipinsubnet":
		return ipInSubnet(str(0), str(1), false)
	case "isipv4insubnet":
		return ipInSubnet(str(0), str(1), true)
	case "isipv6insubnet":
		ip := net.ParseIP(str(0))
		_, cidr, err := net.ParseCIDR(str(1))
		return err == nil && ip != nil && strings.Contains(str(1), ":") && cidr.Contains(ip)
	case "ipv4tonumber":
		ip := net.ParseIP(str(0))
		if ip == nil || ip.To4() == nil {
			return nil
		}
		v := uint32(0)
		for _, b := range ip.To4() {
			v = v<<8 | uint32(b)
		}
		return float64(v)
	case "isprivateip":
		ip := net.ParseIP(str(0))
		if ip == nil || ip.To4() == nil {
			return false
		}
		for _, cidr := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
			_, nw, _ := net.ParseCIDR(cidr)
			if nw.Contains(ip) {
				return true
			}
		}
		return false
	case "ispublicip":
		ip := net.ParseIP(str(0))
		if ip == nil {
			return false
		}
		private := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.0/8", "169.254.0.0/16"}
		reserved := []string{"0.0.0.0/8", "192.0.2.0/24", "198.51.100.0/24", "203.0.113.0/24", "224.0.0.0/4", "240.0.0.0/4"}
		for _, cidr := range append(private, reserved...) {
			_, nw, _ := net.ParseCIDR(cidr)
			if nw.Contains(ip) {
				return false
			}
		}
		return true
	case "isreservedip":
		ip := net.ParseIP(str(0))
		if ip == nil {
			return false
		}
		reserved := []string{"0.0.0.0/8", "10.0.0.0/8", "127.0.0.0/8", "169.254.0.0/16", "172.16.0.0/12",
			"192.0.2.0/24", "192.168.0.0/16", "198.51.100.0/24", "203.0.113.0/24", "224.0.0.0/4", "240.0.0.0/4"}
		for _, cidr := range reserved {
			_, nw, _ := net.ParseCIDR(cidr)
			if nw.Contains(ip) {
				return true
			}
		}
		return false

	// JSON functions.
	case "jsonparse":
		s := strings.TrimSpace(str(0))
		if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
			var decoded interface{}
			if err := json.Unmarshal([]byte(s), &decoded); err == nil {
				return decoded
			}
		}
		return nil
	case "jsonstringify":
		return asString(arg(0))
	case "jsonarraysize":
		var decoded interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(str(0))), &decoded); err != nil {
			return float64(0)
		}
		if list, ok := decoded.([]interface{}); ok {
			return float64(len(list))
		}
		return float64(0)
	case "jsonarraycontains":
		var decoded interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(str(0))), &decoded); err != nil {
			return false
		}
		list, ok := decoded.([]interface{})
		if !ok {
			return false
		}
		for _, item := range list {
			if valuesEqual(item, arg(1)) {
				return true
			}
		}
		return false

	// Datetime functions.
	case "datefloor":
		return dateRound(str(0), str(1), false)
	case "dateceil":
		return dateRound(str(0), str(1), true)
	case "frommillis":
		if f, ok := num(0); ok {
			return timestampValue(int64(f))
		}
		return nil
	case "tomillis":
		if f, ok := num(0); ok {
			return f
		}
		// Timestamps stored as result-formatted or ISO strings are parsed
		// back to milliseconds.
		if ms, ok := parseResultTimestamp(str(0)); ok {
			return float64(ms)
		}
		if t, err := time.Parse(time.RFC3339Nano, str(0)); err == nil {
			return float64(t.UnixMilli())
		}
		return nil
	case "now":
		if ctx != nil {
			return float64(ctx.now() / 1000)
		}
		return float64(time.Now().UnixMilli() / 1000)
	case "parsedate":
		layout := javaTimeLayout(str(1))
		tz := ""
		if len(vals) > 2 {
			tz = str(2)
		}
		loc := time.UTC
		if tz != "" {
			if l, err := time.LoadLocation(tz); err == nil {
				loc = l
			}
		}
		if t, err := time.ParseInLocation(layout, str(0), loc); err == nil {
			return float64(t.UnixMilli())
		}
		return nil
	case "formatdate", "strftime":
		layout := strftimeLayout(str(1))
		tz := ""
		if len(vals) > 2 {
			tz = str(2)
		}
		loc := time.UTC
		if tz != "" {
			if l, err := time.LoadLocation(tz); err == nil {
				loc = l
			}
		}
		ms, ok := num(0)
		if !ok {
			if parsed, ok2 := parseResultTimestamp(str(0)); ok2 {
				return time.UnixMilli(parsed).In(loc).Format(layout)
			}
			if t, err := time.Parse(time.RFC3339Nano, str(0)); err == nil {
				return t.In(loc).Format(layout)
			}
			return nil
		}
		return time.UnixMilli(int64(ms)).In(loc).Format(layout)

	// Hashing functions.
	case "md5":
		sum := md5.Sum([]byte(str(0)))
		return hex.EncodeToString(sum[:])
	case "sha256":
		sum := sha256.Sum256([]byte(str(0)))
		return hex.EncodeToString(sum[:])
	}
	return nil
}

func hexDecodeString(s string) ([]byte, error) {
	s = strings.TrimPrefix(strings.ToLower(s), "0x")
	return hex.DecodeString(s)
}

// ipInSubnet reports whether the address is within the subnet; when v4Only
// is set the address must be IPv4 and the subnet IPv4 as well.
func ipInSubnet(addr, subnet string, v4Only bool) bool {
	ip := net.ParseIP(addr)
	_, cidr, err := net.ParseCIDR(subnet)
	if err != nil || ip == nil {
		return false
	}
	if v4Only && (ip.To4() == nil || strings.Contains(subnet, ":")) {
		return false
	}
	return cidr.Contains(ip)
}

// dateFloorCeil rounds a timestamp down (floor) or up (ceil) to the period
// boundary and truncates, yielding a timestamp-typed value.
func dateRound(ts, period string, ceil bool) interface{} {
	ms, ok := asNumber(ts)
	if !ok {
		if parsed, ok2 := parseResultTimestamp(ts); ok2 {
			ms = float64(parsed)
		} else if t, err := time.Parse(time.RFC3339Nano, ts); err == nil {
			ms = float64(t.UnixMilli())
		} else {
			return ""
		}
	}
	dur, ok := parsePeriodMillis(period)
	if !ok {
		return ""
	}
	v := int64(ms)
	if ceil {
		v = ((v + dur - 1) / dur) * dur
	} else {
		v = (v / dur) * dur
	}
	return timestampValue(v)
}

// parsePeriodMillis parses a period literal such as 5m, 10s, 1h, 2d, 1mo
// into milliseconds. Month, quarter, and year use calendar approximations.
// Full-word units and their documented abbreviations (with optional plural
// s) are accepted.
func parsePeriodMillis(p string) (int64, bool) {
	s := strings.TrimSpace(strings.ToLower(p))
	if s == "" {
		return 0, false
	}
	// Longest unit names first so that "milliseconds" is not read as "ms"
	// followed by stray characters.
	units := []struct {
		name string
		ms   float64
	}{
		{"milliseconds", 1}, {"millisecond", 1}, {"msec", 1}, {"msecs", 1}, {"ms", 1},
		{"seconds", 1000}, {"second", 1000}, {"secs", 1000}, {"sec", 1000}, {"s", 1000},
		{"minutes", 60 * 1000}, {"minute", 60 * 1000}, {"mins", 60 * 1000}, {"min", 60 * 1000}, {"m", 60 * 1000},
		{"hours", 3600 * 1000}, {"hrs", 3600 * 1000}, {"hr", 3600 * 1000}, {"hour", 3600 * 1000}, {"h", 3600 * 1000},
		{"days", 24 * 3600 * 1000}, {"day", 24 * 3600 * 1000}, {"d", 24 * 3600 * 1000},
		{"weeks", 7 * 24 * 3600 * 1000}, {"week", 7 * 24 * 3600 * 1000}, {"w", 7 * 24 * 3600 * 1000},
		{"months", 30 * 24 * 3600 * 1000}, {"month", 30 * 24 * 3600 * 1000}, {"mons", 30 * 24 * 3600 * 1000}, {"mon", 30 * 24 * 3600 * 1000}, {"mo", 30 * 24 * 3600 * 1000},
		{"quarters", 91 * 24 * 3600 * 1000}, {"quarter", 91 * 24 * 3600 * 1000}, {"qtrs", 91 * 24 * 3600 * 1000}, {"qtr", 91 * 24 * 3600 * 1000}, {"q", 91 * 24 * 3600 * 1000},
		{"years", 365 * 24 * 3600 * 1000}, {"year", 365 * 24 * 3600 * 1000}, {"yrs", 365 * 24 * 3600 * 1000}, {"yr", 365 * 24 * 3600 * 1000}, {"y", 365 * 24 * 3600 * 1000},
	}
	for _, u := range units {
		if strings.HasSuffix(s, u.name) {
			numStr := strings.TrimSpace(strings.TrimSuffix(s, u.name))
			if numStr == "" {
				numStr = "1"
			}
			n, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, false
			}
			return int64(n * u.ms), true
		}
	}
	// Bare numbers are treated as milliseconds.
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return int64(n), true
	}
	return 0, false
}

// capPeriodMillis applies the documented caps for time units: 1000 for
// milliseconds, 60 for seconds and minutes, 24 for hours. The cap depends on
// the unit the period was written in, not on the total duration.
func capPeriodMillis(ms, unitMs int64) int64 {
	var cap int64
	switch unitMs {
	case 1:
		cap = 1000
	case 1000:
		cap = 60 * 1000
	case 60 * 1000:
		cap = 60 * 60 * 1000
	case 3600 * 1000:
		cap = 24 * 3600 * 1000
	default:
		return ms
	}
	if ms > cap {
		return cap
	}
	return ms
}

// javaTimeLayout converts the documented subset of Java DateTimeFormatter
// patterns to Go layouts for parseDate.
func javaTimeLayout(pattern string) string {
	repl := []struct{ java, goLayout string }{
		{"yyyy", "2006"}, {"yy", "06"}, {"MM", "01"}, {"dd", "02"},
		{"HH", "15"}, {"mm", "04"}, {"ss", "05"}, {"SSS", "000"},
		{"a", "PM"}, {"EEE", "Mon"}, {"EEEE", "Monday"}, {"Z", "Z07:00"}, {"X", "Z07:00"},
	}
	out := pattern
	for _, r := range repl {
		out = strings.ReplaceAll(out, r.java, r.goLayout)
	}
	return out
}

// strftimeLayout converts the documented strftime-style specifiers to Go
// layouts for formatDate.
func strftimeLayout(pattern string) string {
	repl := []struct{ spec, goLayout string }{
		{"%Y", "2006"}, {"%y", "06"}, {"%m", "01"}, {"%d", "02"},
		{"%H", "15"}, {"%M", "04"}, {"%S", "05"}, {"%j", "002"},
		{"%p", "PM"}, {"%B", "January"}, {"%b", "Jan"}, {"%Z", "MST"},
	}
	out := pattern
	for _, r := range repl {
		out = strings.ReplaceAll(out, r.spec, r.goLayout)
	}
	return out
}
