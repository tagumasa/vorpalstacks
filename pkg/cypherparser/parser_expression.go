// Cypher RETURN, SET, DELETE, REMOVE clause parsing and expression parsing.
//
// Based on goraphdb/cypher_parser.go parseExpr with precedence climbing.
// Extended with IS [NOT] NULL, IN, STARTS WITH, ENDS WITH, CONTAINS, EXISTS
// predicates, CASE/WHEN/THEN/ELSE expressions, list/map literals, arithmetic
// operators, aggregation function parsing, and CREATE/MERGE/DDL query parsing.

package cypherparser

import (
	"fmt"
	"strconv"
	"strings"
)

// Expression parsing with proper precedence
// ---------------------------------------------------------------------------
// Based on goraphdb parseExpr/parseOrExpr/parseAndExpr/parseNotExpr/parseComparison,
// extended with:
//   - Addition/multiplication precedence levels
//   - IS [NOT] NULL, IN, STARTS WITH, ENDS WITH, CONTAINS as postfix operators
//   - CASE expression parsing
//   - EXISTS { pattern } parsing
//   - List/map literal parsing

func (p *parser) parseExpr() (Expression, error) {
	return p.parseOrExpr()
}

func (p *parser) parseOrExpr() (Expression, error) {
	left, err := p.parseAndExpr()
	if err != nil {
		return Expression{}, err
	}
	if !p.is(tokOr) {
		return left, nil
	}
	operands := []Expression{left}
	for p.match(tokOr) {
		right, err := p.parseAndExpr()
		if err != nil {
			return Expression{}, err
		}
		operands = append(operands, right)
	}
	return orExpr(operands...), nil
}

func (p *parser) parseAndExpr() (Expression, error) {
	left, err := p.parseNotExpr()
	if err != nil {
		return Expression{}, err
	}
	if !p.is(tokAnd) {
		return left, nil
	}
	operands := []Expression{left}
	for p.match(tokAnd) {
		right, err := p.parseNotExpr()
		if err != nil {
			return Expression{}, err
		}
		operands = append(operands, right)
	}
	return andExpr(operands...), nil
}

func (p *parser) parseNotExpr() (Expression, error) {
	if p.match(tokNot) {
		inner, err := p.parseNotExpr()
		if err != nil {
			return Expression{}, err
		}
		return notExpr(inner), nil
	}
	return p.parseComparison()
}

func (p *parser) parseComparison() (Expression, error) {
	left, err := p.parseAddition()
	if err != nil {
		return Expression{}, err
	}

	switch {
	case p.is(tokEq):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return compExpr(left, OpEq, right), nil

	case p.is(tokNeq):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return compExpr(left, OpNeq, right), nil

	case p.is(tokLt):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return compExpr(left, OpLt, right), nil

	case p.is(tokGt):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return compExpr(left, OpGt, right), nil

	case p.is(tokLte):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return compExpr(left, OpLte, right), nil

	case p.is(tokGte):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return compExpr(left, OpGte, right), nil

	case p.is(tokRegexMatch):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return regexMatchExpr(left, right), nil

	case p.is(tokIn):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return inExpr(left, right), nil

	case p.is(tokStartsWith):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return startsWithExpr(left, right), nil

	case p.is(tokEndsWith):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return endsWithExpr(left, right), nil

	case p.is(tokContains):
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return Expression{}, err
		}
		return containsExpr(left, right), nil

	case p.is(tokIs):
		p.advance()
		if p.match(tokNot) {
			if _, err := p.expect(tokNull); err != nil {
				return Expression{}, fmt.Errorf("cypher parser: expected NULL after IS NOT")
			}
			return isNotNullExpr(left), nil
		}
		if _, err := p.expect(tokNull); err != nil {
			return Expression{}, fmt.Errorf("cypher parser: expected NULL after IS")
		}
		return isNullExpr(left), nil
	}

	return left, nil
}

func (p *parser) parseAddition() (Expression, error) {
	left, err := p.parseMultiplication()
	if err != nil {
		return Expression{}, err
	}
	for {
		if p.is(tokPlus) {
			p.advance()
			right, err := p.parseMultiplication()
			if err != nil {
				return Expression{}, err
			}
			left = addExpr(left, right)
		} else if p.is(tokDash) {
			p.advance()
			right, err := p.parseMultiplication()
			if err != nil {
				return Expression{}, err
			}
			left = subExpr(left, right)
		} else {
			break
		}
	}
	return left, nil
}

func (p *parser) parseMultiplication() (Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return Expression{}, err
	}
	for {
		if p.is(tokStar) {
			p.advance()
			right, err := p.parseUnary()
			if err != nil {
				return Expression{}, err
			}
			left = mulExpr(left, right)
		} else if p.is(tokSlash) {
			p.advance()
			right, err := p.parseUnary()
			if err != nil {
				return Expression{}, err
			}
			left = divExpr(left, right)
		} else if p.is(tokPercent) {
			p.advance()
			right, err := p.parseUnary()
			if err != nil {
				return Expression{}, err
			}
			left = modExpr(left, right)
		} else {
			break
		}
	}
	return left, nil
}

func (p *parser) parseUnary() (Expression, error) {
	if p.is(tokDash) {
		p.advance()
		inner, err := p.parsePrimary()
		if err != nil {
			return Expression{}, err
		}
		return subExpr(litExpr(0), inner), nil
	}
	return p.parsePrimary()
}

func (p *parser) parsePrimary() (Expression, error) {
	t := p.cur()

	if t.Kind == tokCase {
		return p.parseCaseExpr()
	}

	if t.Kind == tokExists {
		p.advance()
		if _, err := p.expect(tokLBrace); err != nil {
			return Expression{}, fmt.Errorf("cypher parser: EXISTS requires '{ pattern }'")
		}
		pat, err := p.parsePattern()
		if err != nil {
			return Expression{}, err
		}
		if _, err := p.expect(tokRBrace); err != nil {
			return Expression{}, fmt.Errorf("cypher parser: EXISTS requires closing '}'")
		}
		return existsExpr(&pat), nil
	}

	if t.Kind == tokLBracket {
		return p.parseListLiteral()
	}

	if t.Kind == tokLBrace {
		return p.parseMapLiteralExpr()
	}

	if t.Kind == tokCount || t.Kind == tokSum || t.Kind == tokAvg ||
		t.Kind == tokMin || t.Kind == tokMax || t.Kind == tokCollect {
		return p.parseAggregation(t.Kind)
	}

	if t.Kind == tokLabels {
		p.advance()
		if _, err := p.expect(tokLParen); err != nil {
			return Expression{}, fmt.Errorf("cypher parser: LABELS requires '('")
		}
		arg, err := p.parseExpr()
		if err != nil {
			return Expression{}, err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return Expression{}, fmt.Errorf("cypher parser: LABELS requires ')'")
		}
		return labelsExpr(arg), nil
	}
	if t.Kind == tokProperties {
		p.advance()
		if _, err := p.expect(tokLParen); err != nil {
			return Expression{}, fmt.Errorf("cypher parser: PROPERTIES requires '('")
		}
		arg, err := p.parseExpr()
		if err != nil {
			return Expression{}, err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return Expression{}, fmt.Errorf("cypher parser: PROPERTIES requires ')'")
		}
		return propertiesExpr(arg), nil
	}

	if t.Kind == tokIdent || t.Kind == tokType {
		name := t.Text
		p.advance()

		if p.is(tokDot) {
			p.advance()
			propTok, err := p.expectIdentOrKeyword()
			if err != nil {
				return Expression{}, fmt.Errorf("cypher parser: expected property name after '.'")
			}
			expr := propExpr(name, propTok.Text)
			for p.is(tokDot) {
				p.advance()
				nextProp, err := p.expectIdentOrKeyword()
				if err != nil {
					return Expression{}, fmt.Errorf("cypher parser: expected property name after '.'")
				}
				expr = propExprForExpr(expr, nextProp.Text)
			}
			return expr, nil
		}

		if p.is(tokLParen) {
			p.advance()
			var args []Expression
			if !p.is(tokRParen) {
				arg, err := p.parseExpr()
				if err != nil {
					return Expression{}, err
				}
				args = append(args, arg)
				for p.match(tokComma) {
					arg, err = p.parseExpr()
					if err != nil {
						return Expression{}, err
					}
					args = append(args, arg)
				}
			}
			if _, err := p.expect(tokRParen); err != nil {
				return Expression{}, err
			}
			return funcCallExpr(name, args...), nil
		}

		return varRefExpr(name), nil
	}

	if t.Kind == tokParam {
		p.advance()
		return Expression{Kind: ExprParam, ParamName: t.Text}, nil
	}

	switch t.Kind {
	case tokString:
		p.advance()
		return litExpr(t.Text), nil
	case tokInt:
		p.advance()
		n, _ := strconv.ParseInt(t.Text, 10, 64)
		return litExpr(n), nil
	case tokFloat:
		p.advance()
		f, _ := strconv.ParseFloat(t.Text, 64)
		return litExpr(f), nil
	case tokTrue:
		p.advance()
		return litExpr(true), nil
	case tokFalse:
		p.advance()
		return litExpr(false), nil
	case tokNull:
		p.advance()
		return litExpr(nil), nil
	}

	if t.Kind == tokLParen {
		p.advance()
		inner, err := p.parseExpr()
		if err != nil {
			return Expression{}, err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return Expression{}, fmt.Errorf("cypher parser: expected ')' after expression")
		}
		return inner, nil
	}

	return Expression{}, fmt.Errorf("cypher parser: unexpected token %s at position %d",
		tokenKindName(t.Kind), t.Pos)
}

func (p *parser) parseCaseExpr() (Expression, error) {
	if _, err := p.expect(tokCase); err != nil {
		return Expression{}, err
	}

	var operand *Expression
	if !p.is(tokWhen) {
		expr, err := p.parseExpr()
		if err != nil {
			return Expression{}, err
		}
		operand = &expr
	}

	var whens []CaseWhen
	for p.is(tokWhen) {
		p.advance()
		cond, err := p.parseExpr()
		if err != nil {
			return Expression{}, err
		}
		if _, err := p.expect(tokThen); err != nil {
			return Expression{}, fmt.Errorf("cypher parser: expected THEN after WHEN condition")
		}
		result, err := p.parseExpr()
		if err != nil {
			return Expression{}, err
		}
		whens = append(whens, CaseWhen{Condition: &cond, Result: &result})
	}

	var elseExpr *Expression
	if p.match(tokElse) {
		expr, err := p.parseExpr()
		if err != nil {
			return Expression{}, err
		}
		elseExpr = &expr
	}

	if _, err := p.expect(tokEnd); err != nil {
		return Expression{}, fmt.Errorf("cypher parser: expected END to close CASE expression")
	}

	return caseExpr(operand, whens, elseExpr), nil
}

func (p *parser) parseListLiteral() (Expression, error) {
	if _, err := p.expect(tokLBracket); err != nil {
		return Expression{}, err
	}
	var items []Expression
	if !p.is(tokRBracket) {
		item, err := p.parseExpr()
		if err != nil {
			return Expression{}, err
		}
		items = append(items, item)
		for p.match(tokComma) {
			item, err = p.parseExpr()
			if err != nil {
				return Expression{}, err
			}
			items = append(items, item)
		}
	}
	if _, err := p.expect(tokRBracket); err != nil {
		return Expression{}, err
	}
	return listExpr(items...), nil
}

func (p *parser) parseMapLiteralExpr() (Expression, error) {
	if _, err := p.expect(tokLBrace); err != nil {
		return Expression{}, err
	}
	var pairs []MapPair
	for !p.is(tokRBrace) && !p.is(tokEOF) {
		keyTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return Expression{}, fmt.Errorf("cypher parser: expected map key, got %s", tokenKindName(p.cur().Kind))
		}
		if _, err := p.expect(tokColon); err != nil {
			return Expression{}, err
		}
		val, err := p.parseExpr()
		if err != nil {
			return Expression{}, err
		}
		pairs = append(pairs, MapPair{Key: keyTok.Text, Value: val})
		p.match(tokComma)
	}
	if _, err := p.expect(tokRBrace); err != nil {
		return Expression{}, err
	}
	return mapExpr(pairs...), nil
}

func (p *parser) parseAggregation(kind TokenKind) (Expression, error) {
	funcName := p.cur().Text
	p.advance()
	if _, err := p.expect(tokLParen); err != nil {
		return Expression{}, fmt.Errorf("cypher parser: %s requires '('", strings.ToUpper(funcName))
	}

	var fn AggFunc
	switch kind {
	case tokCount:
		fn = AggCount
	case tokSum:
		fn = AggSum
	case tokAvg:
		fn = AggAvg
	case tokMin:
		fn = AggMin
	case tokMax:
		fn = AggMax
	case tokCollect:
		fn = AggCollect
	}

	if kind == tokCount && p.is(tokStar) {
		p.advance()
		if _, err := p.expect(tokRParen); err != nil {
			return Expression{}, err
		}
		return aggExpr(fn, nil), nil
	}

	if p.match(tokDistinct) {
		arg, err := p.parseExpr()
		if err != nil {
			return Expression{}, err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return Expression{}, err
		}
		expr := aggExpr(fn, &arg)
		expr.AggDistinct = true
		return expr, nil
	}

	arg, err := p.parseExpr()
	if err != nil {
		return Expression{}, err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return Expression{}, err
	}
	return aggExpr(fn, &arg), nil
}
