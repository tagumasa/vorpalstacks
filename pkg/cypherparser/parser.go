// Cypher recursive descent parser.
//
// Based on goraphdb/cypher_parser.go parser structure (parseQuery, parsePattern,
// parseNodePattern, parseRelPattern, parseExpr with precedence climbing), extended with:
//   - WITH / UNWIND / DISTINCT clause parsing
//   - REMOVE clause parsing
//   - Multiple MATCH clause support
//   - Edge properties in RelPattern (parseRelPattern extended)
//   - IS [NOT] NULL, IN, STARTS WITH, ENDS WITH, CONTAINS, EXISTS predicates
//   - CASE/WHEN/THEN/ELSE/END expression parsing
//   - List literal [...] and map literal {key: val, ...} parsing
//   - Arithmetic operators (+, -, *, /, %) with proper precedence
//   - Aggregation function parsing (COUNT, SUM, AVG, MIN, MAX, COLLECT)
//   - Expression values in property maps (not just literals)
//   - LABELS() and PROPERTIES() function parsing
//
// Grammar (extended beyond goraphdb):
//
//	Query       → [EXPLAIN|PROFILE]
//	              [WITH [DISTINCT] items]
//	              [MATCH Pattern [WHERE Expr]]*
//	              [OPTIONAL MATCH Pattern [WHERE Expr]]
//	              [UNWIND expr AS var]
//	              [SET items]
//	              [DELETE vars]
//	              [REMOVE items]
//	              RETURN [DISTINCT] Items
//	              [ORDER BY Items] [SKIP n] [LIMIT n]
//
//	Expr        → OrExpr
//	OrExpr      → AndExpr ( OR AndExpr )*
//	AndExpr     → NotExpr ( AND NotExpr )*
//	NotExpr     → [NOT] Comparison
//	Comparison  → Addition [ ( '=' | '<>' | '<' | '>' | '<=' | '>=' | '=~' | IN | STARTS WITH | ENDS WITH | CONTAINS | IS [NOT] NULL ) Addition ]
//	Addition    → Multiplication ( ('+' | '-') Multiplication )*
//	Multiplication → Unary ( ('*' | '/' | '%') Unary )*
//	Unary       → ['-'] Primary
//	Primary     → Literal | VarRef | PropAccess | FuncCall | CaseExpr
//	              | ExistsExpr | ListLiteral | '(' Expr ')' | ParamRef

package cypherparser

import (
	"fmt"
	"strconv"
	"strings"
)

type parser struct {
	tokens []Token
	pos    int
}

// ParsedCypher holds the result of parsing a Cypher query string. Exactly one of Read, Write, Merge, Call, or DDL is non-nil.
type ParsedCypher struct {
	Read  *CypherQuery
	Write *CypherWrite
	Merge *CypherMerge
	Call  *CypherCall
	DDL   *DDLStatement
}

// Parse tokenises and parses a Cypher query string.
func Parse(input string) (*ParsedCypher, error) {
	tokens, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	p := &parser{tokens: tokens}

	if p.is(tokCreate) {
		savedPos := p.pos
		p.advance()
		if p.is(tokIndex) || p.is(tokConstraint) {
			p.pos = savedPos
			return p.parseDDLCreate()
		}
		p.pos = savedPos
		w, err := p.parseCreateQuery()
		if err != nil {
			return nil, err
		}
		return &ParsedCypher{Write: w}, nil
	}

	if p.is(tokShow) {
		return p.parseShow()
	}

	if p.is(tokMerge) {
		m, err := p.parseMergeQuery()
		if err != nil {
			return nil, err
		}
		return &ParsedCypher{Merge: m}, nil
	}

	if p.is(tokCall) {
		c, err := p.parseCallStatement()
		if err != nil {
			return nil, err
		}
		return &ParsedCypher{Call: c}, nil
	}

	q, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	return &ParsedCypher{Read: q}, nil
}

func (p *parser) cur() Token {
	if p.pos >= len(p.tokens) {
		return Token{Kind: tokEOF}
	}
	return p.tokens[p.pos]
}

func (p *parser) advance() Token {
	t := p.cur()
	p.pos++
	return t
}

func (p *parser) expect(kind TokenKind) (Token, error) {
	t := p.cur()
	if t.Kind != kind {
		return t, fmt.Errorf("cypher parser: expected %s but got %s at position %d",
			tokenKindName(kind), tokenKindName(t.Kind), t.Pos)
	}
	p.pos++
	return t, nil
}

func (p *parser) expectOneOf(kinds ...TokenKind) (Token, error) {
	t := p.cur()
	for _, k := range kinds {
		if t.Kind == k {
			p.pos++
			return t, nil
		}
	}
	names := make([]string, len(kinds))
	for i, k := range kinds {
		names[i] = tokenKindName(k)
	}
	return t, fmt.Errorf("cypher parser: expected one of %s but got %s at position %d",
		strings.Join(names, "|"), tokenKindName(t.Kind), t.Pos)
}

func (p *parser) expectIdentOrKeyword() (Token, error) {
	t := p.cur()
	if t.Kind == tokIdent {
		p.pos++
		return t, nil
	}
	if t.Kind >= tokMatch && t.Kind <= tokLabels {
		p.pos++
		return t, nil
	}
	return t, fmt.Errorf("cypher parser: expected identifier but got %s at position %d",
		tokenKindName(t.Kind), t.Pos)
}

func (p *parser) is(kind TokenKind) bool {
	return p.cur().Kind == kind
}

func (p *parser) match(kind TokenKind) bool {
	if p.is(kind) {
		p.pos++
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Query parsing
// ---------------------------------------------------------------------------
// Based on goraphdb parseQuery, extended with WITH, UNWIND, multiple MATCH,
// REMOVE, DISTINCT, and flexible clause ordering support.

func (p *parser) parseQuery() (*CypherQuery, error) {
	q := &CypherQuery{}

	if p.is(tokExplain) {
		p.advance()
		q.Explain = ExplainOnly
	} else if p.is(tokProfile) {
		p.advance()
		q.Explain = ExplainProfile
	}

	for !p.is(tokEOF) {
		switch {
		case p.is(tokWith):
			wc, err := p.parseWithClause()
			if err != nil {
				return nil, err
			}
			q.With = &wc
			if p.is(tokWhere) {
				p.advance()
				expr, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				q.With.Where = &expr
			}
			if p.is(tokOrder) {
				p.advance()
				if _, err := p.expect(tokBy); err != nil {
					return nil, err
				}
				items, err := p.parseOrderItems()
				if err != nil {
					return nil, err
				}
				q.With.OrderBy = items
			}
			if p.is(tokSkip) {
				p.advance()
				tok, err := p.expectOneOf(tokInt, tokParam)
				if err != nil {
					return nil, fmt.Errorf("cypher parser: SKIP requires an integer")
				}
				if tok.Kind == tokParam {
					q.SkipParam = tok.Text
				} else {
					n, _ := strconv.Atoi(tok.Text)
					q.With.Skip = intPtr(n)
				}
			}
			if p.is(tokLimit) {
				p.advance()
				tok, err := p.expectOneOf(tokInt, tokParam)
				if err != nil {
					return nil, fmt.Errorf("cypher parser: LIMIT requires an integer")
				}
				if tok.Kind == tokParam {
					q.LimitParam = tok.Text
				} else {
					n, _ := strconv.Atoi(tok.Text)
					q.With.Limit = intPtr(n)
				}
			}

		case p.is(tokMatch):
			p.advance()
			pat, err := p.parsePattern()
			if err != nil {
				return nil, err
			}
			mc := MatchClause{Pattern: pat}

			var where *Expression
			if p.is(tokWhere) {
				p.advance()
				expr, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				where = &expr
			}

			q.Matches = append(q.Matches, mc)
			if len(q.Matches) == 1 {
				q.Match = mc
				q.Where = where
			} else if where != nil && q.Where != nil {
				q.Where = &Expression{
					Kind:     ExprAnd,
					Operands: []Expression{*q.Where, *where},
				}
			} else if where != nil {
				q.Where = where
			}

		case p.is(tokOptional):
			p.advance()
			if _, err := p.expect(tokMatch); err != nil {
				return nil, fmt.Errorf("cypher parser: expected MATCH after OPTIONAL")
			}
			optPat, err := p.parsePattern()
			if err != nil {
				return nil, err
			}
			q.OptionalMatch = &MatchClause{Pattern: optPat}
			if p.is(tokWhere) {
				p.advance()
				expr, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				q.OptionalWhere = &expr
			}

		case p.is(tokUnwind):
			uc, err := p.parseUnwindClause()
			if err != nil {
				return nil, err
			}
			q.Unwind = &uc

		case p.is(tokSet):
			p.advance()
			items, err := p.parseSetItems()
			if err != nil {
				return nil, err
			}
			q.Set = items

		case p.is(tokDetach):
			p.advance()
			if _, err := p.expect(tokDelete); err != nil {
				return nil, err
			}
			vars, err := p.parseDeleteVars()
			if err != nil {
				return nil, err
			}
			q.Delete = vars
			q.DeleteDetach = true

		case p.is(tokDelete):
			p.advance()
			detach := p.match(tokDetach)
			vars, err := p.parseDeleteVars()
			if err != nil {
				return nil, err
			}
			q.Delete = vars
			q.DeleteDetach = detach

		case p.is(tokRemove):
			p.advance()
			items, err := p.parseRemoveItems()
			if err != nil {
				return nil, err
			}
			q.Remove = items

		case p.is(tokReturn):
			ret, err := p.parseReturnClause()
			if err != nil {
				return nil, err
			}
			q.Return = ret

			if p.is(tokOrder) {
				p.advance()
				if _, err := p.expect(tokBy); err != nil {
					return nil, err
				}
				items, err := p.parseOrderItems()
				if err != nil {
					return nil, err
				}
				q.OrderBy = items
			}

			if p.is(tokSkip) {
				p.advance()
				tok, err := p.expectOneOf(tokInt, tokParam)
				if err != nil {
					return nil, fmt.Errorf("cypher parser: SKIP requires an integer")
				}
				if tok.Kind == tokParam {
					q.SkipParam = tok.Text
				} else {
					n, _ := strconv.Atoi(tok.Text)
					q.Skip = intPtr(n)
				}
			}

			if p.is(tokLimit) {
				p.advance()
				tok, err := p.expectOneOf(tokInt, tokParam)
				if err != nil {
					return nil, fmt.Errorf("cypher parser: LIMIT requires an integer")
				}
				if tok.Kind == tokParam {
					q.LimitParam = tok.Text
				} else {
					n, _ := strconv.Atoi(tok.Text)
					q.Limit = intPtr(n)
				}
			}

			if !p.is(tokEOF) {
				return nil, fmt.Errorf("cypher parser: unexpected token %s after RETURN at position %d",
					tokenKindName(p.cur().Kind), p.cur().Pos)
			}

		case p.is(tokCreate):
			w, err := p.parseCreateQuery()
			if err != nil {
				return nil, err
			}
			q.Create = w

		default:
			return nil, fmt.Errorf("cypher parser: unexpected token %s at position %d",
				tokenKindName(p.cur().Kind), p.cur().Pos)
		}
	}

	return q, nil
}

// ---------------------------------------------------------------------------
// WITH clause
// ---------------------------------------------------------------------------

func (p *parser) parseWithClause() (WithClause, error) {
	if _, err := p.expect(tokWith); err != nil {
		return WithClause{}, err
	}
	wc := WithClause{}
	if p.match(tokDistinct) {
		wc.Distinct = true
	}
	for {
		item, err := p.parseWithItem()
		if err != nil {
			return wc, err
		}
		wc.Items = append(wc.Items, item)
		if !p.match(tokComma) {
			break
		}
	}
	return wc, nil
}

func (p *parser) parseWithItem() (WithItem, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return WithItem{}, err
	}
	wi := WithItem{Expr: expr}
	if p.is(tokAs) {
		p.advance()
		aliasTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return wi, fmt.Errorf("cypher parser: expected alias after AS")
		}
		wi.Alias = aliasTok.Text
	}
	return wi, nil
}

// ---------------------------------------------------------------------------
// UNWIND clause
// ---------------------------------------------------------------------------

func (p *parser) parseUnwindClause() (UnwindClause, error) {
	if _, err := p.expect(tokUnwind); err != nil {
		return UnwindClause{}, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return UnwindClause{}, err
	}
	if _, err := p.expect(tokAs); err != nil {
		return UnwindClause{}, fmt.Errorf("cypher parser: expected AS after UNWIND expression")
	}
	varTok, err := p.expectIdentOrKeyword()
	if err != nil {
		return UnwindClause{}, fmt.Errorf("cypher parser: expected variable name after UNWIND AS")
	}
	return UnwindClause{Expr: expr, Var: varTok.Text}, nil
}
