package metricmath

import (
	"fmt"
	"strconv"
)

// Parse parses a metric math expression string and returns the root AST
// node. The grammar follows standard arithmetic operator precedence:
//
//	expr   := term (('+' | '-') term)*
//	term   := factor (('*' | '/') factor)*
//	power  := unary ('^' power)?
//	unary  := '-' unary | atom
//	atom   := number | ident | ident '(' args ')' | '(' expr ')'
//	args   := expr (',' expr)*
func Parse(input string) (Expr, error) {
	tokens, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	p := &parser{tokens: tokens}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.peek().Type != TokEOF {
		return nil, fmt.Errorf("metricmath: unexpected token %q at position %d", p.peek().Value, p.peek().Pos)
	}
	return expr, nil
}

type parser struct {
	tokens []Token
	pos    int
}

func (p *parser) peek() Token { return p.tokens[p.pos] }

func (p *parser) consume() Token {
	t := p.tokens[p.pos]
	p.pos++
	return t
}

func (p *parser) parseExpr() (Expr, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == TokPlus || p.peek().Type == TokMinus {
		op := p.consume()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &binaryExpr{op: op, left: left, right: right}
	}
	return left, nil
}

func (p *parser) parseTerm() (Expr, error) {
	left, err := p.parsePower()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == TokStar || p.peek().Type == TokSlash {
		op := p.consume()
		right, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		left = &binaryExpr{op: op, left: left, right: right}
	}
	return left, nil
}

func (p *parser) parsePower() (Expr, error) {
	base, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	if p.peek().Type == TokCaret {
		op := p.consume()
		exp, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		return &binaryExpr{op: op, left: base, right: exp}, nil
	}
	return base, nil
}

func (p *parser) parseUnary() (Expr, error) {
	if p.peek().Type == TokMinus {
		p.consume()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &unaryExpr{operand: operand}, nil
	}
	return p.parseAtom()
}

func (p *parser) parseAtom() (Expr, error) {
	tok := p.peek()
	switch tok.Type {
	case TokNumber:
		p.consume()
		val, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("metricmath: invalid number %q at position %d", tok.Value, tok.Pos)
		}
		return &numberExpr{value: val}, nil

	case TokIdent:
		p.consume()
		if p.peek().Type == TokLParen {
			p.consume()
			var args []Expr
			if p.peek().Type != TokRParen {
				for {
					arg, err := p.parseExpr()
					if err != nil {
						return nil, err
					}
					args = append(args, arg)
					if p.peek().Type == TokComma {
						p.consume()
						continue
					}
					break
				}
			}
			if p.peek().Type != TokRParen {
				return nil, fmt.Errorf("metricmath: expected ')' at position %d", p.peek().Pos)
			}
			p.consume()
			return &funcCallExpr{name: tok.Value, args: args}, nil
		}
		return &varRefExpr{name: tok.Value}, nil

	case TokLParen:
		p.consume()
		inner, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.peek().Type != TokRParen {
			return nil, fmt.Errorf("metricmath: expected ')' at position %d", p.peek().Pos)
		}
		p.consume()
		return inner, nil

	default:
		return nil, fmt.Errorf("metricmath: unexpected token %q at position %d", tok.Value, tok.Pos)
	}
}
