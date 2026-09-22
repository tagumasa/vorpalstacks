package cloudwatchlogs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// The primary-and-tail half of the recursive-descent expression parser:
// the operand productions (literals, field paths, function calls,
// parenthesised expressions) and the operator tails that bind to a
// parsed left-hand side (like, regular-expression match, in).

// parseLikeTail parses the pattern after the like keyword. A pattern enclosed
// in forward slashes is a regular expression; a quoted pattern is a literal
// substring ("To match a substring with like and not like, enclose the
// substring that you want to match in single or double quotation marks" —
// the vocabulary carries no wildcard form).
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
		return &likeNode{x: x, substr: t.text, neg: neg}, nil
	case tokRegex:
		re, err := regexp.Compile(t.text)
		if err != nil {
			return nil, newQueryCompileError(fmt.Sprintf("Invalid regular expression: %v", err), t.start, t.end)
		}
		return &likeNode{x: x, re: re, neg: neg}, nil
	case tokOp:
		if t.text == "/" {
			startTok := t
			var raw strings.Builder
			closing, ok := p.collectSlashRegex(&raw)
			if !ok {
				return nil, newQueryCompileError("Unterminated regular expression", startTok.start, startTok.start+1)
			}
			re, err := regexp.Compile(raw.String())
			if err != nil {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid regular expression: %v", err), startTok.start, closing.end)
			}
			return &likeNode{x: x, re: re, neg: neg}, nil
		}
	}
	return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
}

// parseRegexMatchTail parses the right-hand side of the =~ operator: a
// regular expression literal — "You can use the regular expression
// operator =~ to match substrings", with the documented example
// filter f1 =~ /Exception/ (query syntax, filter commands). The match is
// the like engine's unanchored regex mode, so the pattern matches
// substrings exactly as documented.
func (p *exprParser) parseRegexMatchTail(x exprNode) (exprNode, error) {
	t, ok := p.next()
	if !ok {
		last := p.toks[len(p.toks)-1]
		return nil, newQueryCompileError("Syntax error: =~ requires a regular expression", last.end, last.end)
	}
	switch t.kind {
	case tokRegex:
		re, err := regexp.Compile(t.text)
		if err != nil {
			return nil, newQueryCompileError(fmt.Sprintf("Invalid regular expression: %v", err), t.start, t.end)
		}
		return &likeNode{x: x, re: re}, nil
	case tokOp:
		if t.text == "/" {
			startTok := t
			var raw strings.Builder
			closing, ok := p.collectSlashRegex(&raw)
			if !ok {
				return nil, newQueryCompileError("Unterminated regular expression", startTok.start, startTok.start+1)
			}
			re, err := regexp.Compile(raw.String())
			if err != nil {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid regular expression: %v", err), startTok.start, closing.end)
			}
			return &likeNode{x: x, re: re}, nil
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
		// "Nested subqueries are not supported." — a parenthesised
		// operand after in is always a subquery (literal lists use
		// brackets), so an in-followed-by-parenthesis inside the
		// interior is the nesting signature at token level.
		for i := 0; i+1 < len(inner); i++ {
			if inner[i].kind == tokIdent && strings.EqualFold(inner[i].text, "in") &&
				inner[i+1].kind == tokLParen {
				return nil, newQueryCompileError("Nested subqueries are not supported",
					inner[i+1].start, inner[i+1].end)
			}
		}
		sub := &subqueryNode{toks: inner, startTk: t}
		return &inNode{x: x, list: []exprNode{sub}, neg: neg}, nil
	}
	if t.kind != tokLBracket {
		// The documented array form: "To check for elements in an array,
		// put the array after in" — an operand evaluating to a list. The
		// bracket literal and the parenthesised subquery forms are handled
		// above; the operand parses at additive precedence so boolean
		// operators keep binding the enclosing expression.
		p.pos--
		rhs, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}
		return &inNode{x: x, operand: rhs, neg: neg}, nil
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
	"ispresent": {1, 1}, "coalesce": {1, -1},
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

// checkFunctionArity rejects calls whose argument count falls outside
// the documented range and names outside the accepted vocabulary: the
// table is the single inventory of legal function names — the general
// dispatcher's cases and the stats aggregations both sit inside it — so
// an absent name is an unknown function and rejects at compile instead
// of silently evaluating to null.
