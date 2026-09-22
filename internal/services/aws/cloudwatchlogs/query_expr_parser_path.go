package cloudwatchlogs

import (
	"fmt"
	"strconv"
	"strings"
)

// The call-and-path half of the recursive-descent expression parser:
// function-call arity checking, the period-literal forms the datetime
// and bin vocabulary takes, and field-path segmentation.

func checkFunctionArity(name token, args []exprNode) error {
	arity, ok := queryFunctionArities[strings.ToLower(name.text)]
	if !ok {
		return newQueryCompileError(
			fmt.Sprintf("Unknown function: %s", name.text), name.start, name.end)
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

// splitPathSegs splits a dotted field path into structure segments,
// preserving empty-free order. A numeric dot segment addresses an array
// position: the documented discoverable-field form names the first item of
// requestParameters.instancesSet.items as
// requestParameters.instancesSet.items.0.instanceId — the number refers to
// the position of the field's values.
func splitPathSegs(text string) []pathSeg {
	var segs []pathSeg
	for _, part := range strings.Split(text, ".") {
		if part == "" {
			continue
		}
		if n, err := strconv.Atoi(part); err == nil {
			segs = append(segs, pathSeg{index: n, isIdx: true})
			continue
		}
		segs = append(segs, pathSeg{name: part})
	}
	return segs
}

// parseFieldPath parses a field reference with dot and bracket access. The
// first token may already contain dot-separated segments; bracket indexing
// and post-bracket dot continuation extend the path.
func (p *exprParser) parseFieldPath(first token) (exprNode, error) {
	node := &fieldNode{}
	appendSegs := func(text string) {
		node.segs = append(node.segs, splitPathSegs(text)...)
	}
	if first.kind == tokBacktickIdent {
		node.segs = append(node.segs, pathSeg{name: first.text})
	} else if first.kind == tokIdent && first.text == "@" {
		// The system-field form @`name`: "place the @ symbol outside the
		// backticks and enclose only the field name portion" — e.g.
		// filter @`aws.tag.aws:cloudformation:stack-name` = "my-stack".
		if bt, ok := p.peek(); ok && bt.kind == tokBacktickIdent {
			p.pos++
			node.segs = append(node.segs, pathSeg{name: "@" + bt.text})
		} else {
			node.segs = append(node.segs, pathSeg{name: "@"})
		}
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

// collectSlashRegex consumes the tokens of a slash-delimited regular
// expression literal up to (not including) the closing slash, appending
// each token's space-joined raw form to raw: one lexical rule shared by
// the like-expression's pattern and the parse command's regex argument.
func (p *exprParser) collectSlashRegex(raw *strings.Builder) (token, bool) {
	closing, ok := p.next()
	for ok && !(closing.kind == tokOp && closing.text == "/") {
		if raw.Len() > 0 {
			raw.WriteByte(' ')
		}
		raw.WriteString(closing.raw)
		closing, ok = p.next()
	}
	return closing, ok
}
