package cloudwatchlogs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Recursive-descent parser for query expressions: tokens in, exprNode out.

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

// periodArgAt folds the token at toks[i] with a directly following period
// unit word (5m lexes as the two tokens 5 and m) into one period literal,
// returning the joined text and how many tokens it spans. ok is false when
// toks[i] is not a number.
func periodArgAt(toks []token, i int) (text string, span int, ok bool) {
	if i >= len(toks) || toks[i].kind != tokNumber {
		return "", 0, false
	}
	if i+1 < len(toks) && toks[i+1].kind == tokIdent && isPeriodUnitWord(toks[i+1].text) {
		return toks[i].raw + toks[i+1].text, 2, true
	}
	return toks[i].text, 1, true
}

// parseOffsetDuration consumes the duration after the offset keyword,
// folding a number and unit word written without separation (1h) into one
// period literal.
func (p *exprParser) parseOffsetDuration(head token) (int64, error) {
	d, ok := p.next()
	if !ok {
		return 0, newQueryCompileError("Syntax error: offset requires a duration", head.start, head.end)
	}
	text := d.text
	if u, ok := p.peek(); ok && u.kind == tokIdent && isPeriodUnitWord(u.text) {
		p.pos++
		text = d.raw + u.text
	}
	ms, ok2 := parsePeriodMillis(text)
	if !ok2 {
		return 0, newQueryCompileError(fmt.Sprintf("Invalid offset '%s'", text), d.start, d.end)
	}
	return ms, nil
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
