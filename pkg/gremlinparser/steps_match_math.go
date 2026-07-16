// Match and math step executors with expression evaluation.

package gremlinparser

import (
	"fmt"
	mathpkg "math"
	"strconv"
	"strings"

	"vorpalstacks/internal/core/storage/graphengine"
)

// execMatch implements the match() step: evaluates multiple traversal patterns
// (cross-product join) and returns traversers whose label bindings satisfy all patterns.
//
//	g.V().match(
//	    __.as('a').out('knows').as('b'),
//	    __.as('b').out('created').as('c')
//	)
//
// Each pattern argument is an ArgNestedTraversal. For each input traverser, the
// algorithm evaluates patterns sequentially, propagating as() bindings so that
// subsequent patterns can reference labels bound by earlier patterns.
func execMatch(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	patterns := make([]*Traversal, 0, len(step.Args))
	for _, arg := range step.Args {
		if arg.Kind == ArgNestedTraversal {
			patterns = append(patterns, arg.Trav)
		}
	}
	if len(patterns) == 0 {
		return traversers, nil
	}

	var result []*Traverser
	for _, t := range traversers {
		candidates := []*Traverser{t}
		for _, pattern := range patterns {
			var next []*Traverser
			for _, cand := range candidates {
				matched, err := executeNestedTraversal(ec, cand, pattern)
				if err != nil {
					continue
				}
				for _, m := range matched {
					merged := cand.clone()
					for k, v := range m.Tags {
						merged.Tags[k] = v
					}
					next = append(next, merged)
				}
			}
			candidates = next
			if len(candidates) == 0 {
				break
			}
		}
		result = append(result, candidates...)
	}
	return result, nil
}

// execMath implements the math() step: evaluates an arithmetic expression on each traverser.
//
//	g.V().values('age').math('(_ * 2) + 1')
//	g.V().has('age').math('sin(_ * PI / 180)')
//
// The underscore '_' refers to the current traverser element. Supported operators:
// +, -, *, /, %, ^. Parentheses for grouping. Functions: sin, cos, tan, sqrt, abs,
// log, log10, ceil, floor, round, min, max, pow, PI, E.
func execMath(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	expr := argString(step.Args[0])
	if expr == "" {
		return traversers, nil
	}

	result := make([]*Traverser, 0, len(traversers))
	for _, t := range traversers {
		val, err := evalMathExpr(expr, t.Element)
		if err != nil {
			return nil, fmt.Errorf("gremlin: math() error for %v: %w", t.Element, err)
		}
		nt := t.clone()
		nt.Element = val
		nt.Path = append(nt.Path, val)
		result = append(result, nt)
	}
	return result, nil
}

type mathParser struct {
	expr          string
	pos           int
	underscoreVal float64
}

func evalMathExpr(expr string, element any) (float64, error) {
	uv, ok := toFloat64(element)
	if !ok {
		return 0, fmt.Errorf("math() requires numeric element, got %T", element)
	}
	p := &mathParser{expr: expr, pos: 0, underscoreVal: uv}
	p.skipSpaces()
	val, err := p.parseExpr()
	if err != nil {
		return 0, err
	}
	p.skipSpaces()
	if p.pos < len(p.expr) {
		return 0, fmt.Errorf("unexpected character at position %d: %q", p.pos, string(p.expr[p.pos]))
	}
	return val, nil
}

func (p *mathParser) skipSpaces() {
	for p.pos < len(p.expr) && p.expr[p.pos] == ' ' {
		p.pos++
	}
}

func (p *mathParser) peek() byte {
	if p.pos >= len(p.expr) {
		return 0
	}
	return p.expr[p.pos]
}

func (p *mathParser) consume(ch byte) bool {
	if p.peek() == ch {
		p.pos++
		return true
	}
	return false
}

func (p *mathParser) parseExpr() (float64, error) {
	return p.parseAddSub()
}

func (p *mathParser) parseAddSub() (float64, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return 0, err
	}
	for {
		p.skipSpaces()
		if p.consume('+') {
			p.skipSpaces()
			right, err := p.parseMulDiv()
			if err != nil {
				return 0, err
			}
			left += right
		} else if p.consume('-') {
			p.skipSpaces()
			right, err := p.parseMulDiv()
			if err != nil {
				return 0, err
			}
			left -= right
		} else {
			break
		}
	}
	return left, nil
}

func (p *mathParser) parseMulDiv() (float64, error) {
	left, err := p.parsePower()
	if err != nil {
		return 0, err
	}
	for {
		p.skipSpaces()
		if p.consume('*') {
			p.skipSpaces()
			right, err := p.parsePower()
			if err != nil {
				return 0, err
			}
			left *= right
		} else if p.consume('/') {
			p.skipSpaces()
			right, err := p.parsePower()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left /= right
		} else if p.consume('%') {
			p.skipSpaces()
			right, err := p.parsePower()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, fmt.Errorf("modulo by zero")
			}
			left = mathpkg.Mod(left, right)
		} else {
			break
		}
	}
	return left, nil
}

func (p *mathParser) parsePower() (float64, error) {
	base, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	p.skipSpaces()
	if p.consume('^') {
		p.skipSpaces()
		exp, err := p.parsePower()
		if err != nil {
			return 0, err
		}
		base = mathpkg.Pow(base, exp)
	}
	return base, nil
}

func (p *mathParser) parseUnary() (float64, error) {
	p.skipSpaces()
	if p.consume('-') {
		p.skipSpaces()
		val, err := p.parseAtom()
		if err != nil {
			return 0, err
		}
		return -val, nil
	}
	if p.consume('+') {
		p.skipSpaces()
	}
	return p.parseAtom()
}

func (p *mathParser) parseAtom() (float64, error) {
	p.skipSpaces()
	ch := p.peek()

	if ch == '(' {
		p.pos++
		val, err := p.parseExpr()
		if err != nil {
			return 0, err
		}
		p.skipSpaces()
		if !p.consume(')') {
			return 0, fmt.Errorf("expected ')'")
		}
		return val, nil
	}

	if ch == '_' {
		p.pos++
		return p.underscoreVal, nil
	}

	if isDigit(ch) || ch == '.' {
		return p.parseNumber()
	}

	return p.parseFunction()
}

func (p *mathParser) parseNumber() (float64, error) {
	start := p.pos
	for p.pos < len(p.expr) && isDigit(p.expr[p.pos]) {
		p.pos++
	}
	if p.pos < len(p.expr) && p.expr[p.pos] == '.' {
		p.pos++
		for p.pos < len(p.expr) && isDigit(p.expr[p.pos]) {
			p.pos++
		}
	}
	s := p.expr[start:p.pos]
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q", s)
	}
	return val, nil
}

func (p *mathParser) parseFunction() (float64, error) {
	start := p.pos
	for p.pos < len(p.expr) && isAlpha(p.expr[p.pos]) {
		p.pos++
	}
	name := strings.ToUpper(p.expr[start:p.pos])

	p.skipSpaces()
	switch name {
	case "PI":
		return mathpkg.Pi, nil
	case "E":
		return mathpkg.E, nil
	}

	if !p.consume('(') {
		return 0, fmt.Errorf("expected '(' after function %s", name)
	}
	p.skipSpaces()
	val, err := p.parseExpr()
	if err != nil {
		return 0, err
	}
	p.skipSpaces()
	if !p.consume(')') {
		return 0, fmt.Errorf("expected ')' after function argument")
	}

	switch name {
	case "SIN":
		return mathpkg.Sin(val), nil
	case "COS":
		return mathpkg.Cos(val), nil
	case "TAN":
		return mathpkg.Tan(val), nil
	case "ASIN":
		return mathpkg.Asin(val), nil
	case "ACOS":
		return mathpkg.Acos(val), nil
	case "ATAN":
		return mathpkg.Atan(val), nil
	case "SQRT":
		return mathpkg.Sqrt(val), nil
	case "ABS":
		return mathpkg.Abs(val), nil
	case "LOG":
		return mathpkg.Log(val), nil
	case "LOG10":
		return mathpkg.Log10(val), nil
	case "CEIL":
		return mathpkg.Ceil(val), nil
	case "FLOOR":
		return mathpkg.Floor(val), nil
	case "ROUND":
		return mathpkg.Round(val), nil
	case "EXP":
		return mathpkg.Exp(val), nil
	case "SIGN":
		if val < 0 {
			return -1, nil
		}
		if val > 0 {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unknown math function: %s", name)
	}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_'
}

// elementOrValueEqual compares two values for equality, handling node and edge identity comparison.
func elementOrValueEqual(a, b any) bool {
	if a == b {
		return true
	}
	an, aok := a.(*graphengine.Node)
	bn, bok := b.(*graphengine.Node)
	if aok && bok {
		return an.ID == bn.ID
	}
	ae, aok := a.(*graphengine.Edge)
	be, bok := b.(*graphengine.Edge)
	if aok && bok {
		return ae.ID == be.ID
	}
	return false
}
