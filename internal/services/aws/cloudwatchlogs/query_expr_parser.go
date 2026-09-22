package cloudwatchlogs

import (
	"fmt"
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
	case t.kind == tokOp && (t.text == "=" || t.text == "!=" ||
		t.text == "<" || t.text == "<=" || t.text == ">" || t.text == ">="):
		p.pos++
		r, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}
		return &binOpNode{op: t.text, l: l, r: r}, nil
	case t.kind == tokOp && t.text == "=~":
		p.pos++
		return p.parseRegexMatchTail(l)
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
