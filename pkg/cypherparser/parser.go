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

	"vorpalstacks/pkg/graphengine"
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

// ---------------------------------------------------------------------------
// Pattern parsing
// ---------------------------------------------------------------------------
// Based on goraphdb parsePattern/parseNodePattern/parseRelPattern, extended
// with edge properties support.

func (p *parser) parsePattern() (Pattern, error) {
	pat := Pattern{}
	node, err := p.parseNodePattern()
	if err != nil {
		return pat, err
	}
	pat.Nodes = append(pat.Nodes, node)

	for p.is(tokDash) || p.is(tokLArrow) {
		rel, err := p.parseRelPattern()
		if err != nil {
			return pat, err
		}
		pat.Rels = append(pat.Rels, rel)

		node, err = p.parseNodePattern()
		if err != nil {
			return pat, err
		}
		pat.Nodes = append(pat.Nodes, node)
	}

	for p.is(tokComma) {
		p.advance()
		node, err := p.parseNodePattern()
		if err != nil {
			return pat, err
		}
		pat.Nodes = append(pat.Nodes, node)

		for p.is(tokDash) || p.is(tokLArrow) {
			rel, err := p.parseRelPattern()
			if err != nil {
				return pat, err
			}
			pat.Rels = append(pat.Rels, rel)

			node, err = p.parseNodePattern()
			if err != nil {
				return pat, err
			}
			pat.Nodes = append(pat.Nodes, node)
		}
	}
	return pat, nil
}

func (p *parser) parseNodePattern() (NodePattern, error) {
	np := NodePattern{}
	if _, err := p.expect(tokLParen); err != nil {
		return np, err
	}

	if p.is(tokIdent) {
		np.Variable = p.advance().Text
	}

	for p.is(tokColon) {
		p.advance()
		if !p.is(tokIdent) {
			return np, fmt.Errorf("cypher parser: expected label name after ':' at position %d", p.cur().Pos)
		}
		np.Labels = append(np.Labels, p.advance().Text)
	}

	if p.is(tokLBrace) {
		p.advance()
		props, err := p.parsePropMapExpr()
		if err != nil {
			return np, err
		}
		np.Props = props
	}

	if _, err := p.expect(tokRParen); err != nil {
		return np, err
	}
	return np, nil
}

// parsePropMapExpr parses a property map where values can be any expression.
// Our addition — goraphdb only supported literal values in property maps.
func (p *parser) parsePropMapExpr() (map[string]any, error) {
	props := make(map[string]any)
	for !p.is(tokRBrace) && !p.is(tokEOF) {
		keyTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return nil, fmt.Errorf("cypher parser: expected property key, got %s", tokenKindName(p.cur().Kind))
		}
		if _, err := p.expect(tokColon); err != nil {
			return nil, err
		}

		// Check for $param reference
		if p.is(tokParam) {
			paramTok := p.advance()
			props[keyTok.Text] = paramRef(paramTok.Text)
			continue
		}

		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		// If the expression is a simple literal, store the value directly.
		// Otherwise store the expression AST for later evaluation.
		if val.Kind == ExprLiteral {
			props[keyTok.Text] = val.LitValue
		} else {
			props[keyTok.Text] = val
		}

		p.match(tokComma)
	}

	if _, err := p.expect(tokRBrace); err != nil {
		return nil, err
	}
	return props, nil
}

func (p *parser) parseRelPattern() (RelPattern, error) {
	rp := RelPattern{
		Dir:     graphengine.Outgoing,
		MinHops: 1,
		MaxHops: 1,
	}

	leftArrow := false
	if p.is(tokLArrow) {
		leftArrow = true
		p.advance()
	} else if p.is(tokDash) {
		p.advance()
	} else {
		return rp, fmt.Errorf("cypher parser: expected '-' or '<-' at position %d", p.cur().Pos)
	}

	if _, err := p.expect(tokLBracket); err != nil {
		return rp, err
	}

	if p.is(tokIdent) {
		rp.Variable = p.advance().Text
	}

	if p.is(tokColon) {
		p.advance()
		labelTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return rp, fmt.Errorf("cypher parser: expected relationship label after ':'")
		}
		rp.Label = labelTok.Text
	}

	// Edge properties (our addition — goraphdb did not support this).
	if p.is(tokLBrace) {
		p.advance()
		props, err := p.parsePropMapExpr()
		if err != nil {
			return rp, err
		}
		rp.Props = props
	}

	if p.is(tokStar) {
		p.advance()
		rp.VarLength = true
		rp.MinHops = 1
		rp.MaxHops = -1

		if p.is(tokInt) {
			n, _ := strconv.Atoi(p.advance().Text)
			rp.MinHops = n
			rp.MaxHops = n
		}

		if p.is(tokDotDot) {
			p.advance()
			rp.MaxHops = -1
			if p.is(tokInt) {
				n, _ := strconv.Atoi(p.advance().Text)
				rp.MaxHops = n
			}
		}
	}

	if _, err := p.expect(tokRBracket); err != nil {
		return rp, err
	}

	if leftArrow {
		if _, err := p.expect(tokDash); err != nil {
			return rp, fmt.Errorf("cypher parser: expected '-' to close '<-[...]-' pattern")
		}
		rp.Dir = graphengine.Incoming
	} else {
		if p.is(tokArrow) {
			p.advance()
			rp.Dir = graphengine.Outgoing
		} else if p.is(tokDash) {
			p.advance()
			rp.Dir = graphengine.Both
		} else {
			return rp, fmt.Errorf("cypher parser: expected '->' or '-' after relationship pattern at position %d", p.cur().Pos)
		}
	}

	return rp, nil
}

// ---------------------------------------------------------------------------
// RETURN clause
// ---------------------------------------------------------------------------
// Based on goraphdb parseReturnClause, extended with DISTINCT support.

func (p *parser) parseReturnClause() (ReturnClause, error) {
	rc := ReturnClause{}
	if _, err := p.expect(tokReturn); err != nil {
		return rc, err
	}
	if p.match(tokDistinct) {
		rc.Distinct = true
	}
	for {
		item, err := p.parseReturnItem()
		if err != nil {
			return rc, err
		}
		rc.Items = append(rc.Items, item)
		if !p.match(tokComma) {
			break
		}
	}
	return rc, nil
}

func (p *parser) parseReturnItem() (ReturnItem, error) {
	if p.is(tokStar) {
		p.advance()
		return ReturnItem{Expr: Expression{Kind: ExprLiteral, LitValue: "*"}}, nil
	}
	expr, err := p.parseExpr()
	if err != nil {
		return ReturnItem{}, err
	}
	ri := ReturnItem{Expr: expr}
	if p.is(tokAs) {
		p.advance()
		aliasTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return ri, fmt.Errorf("cypher parser: expected alias after AS")
		}
		ri.Alias = aliasTok.Text
	}
	return ri, nil
}

// ---------------------------------------------------------------------------
// ORDER BY
// ---------------------------------------------------------------------------
// Based on goraphdb parseOrderItems.

func (p *parser) parseOrderItems() ([]OrderItem, error) {
	var items []OrderItem
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		item := OrderItem{Expr: expr}
		if p.is(tokDesc) {
			p.advance()
			item.Desc = true
		} else {
			p.match(tokAsc)
		}
		items = append(items, item)
		if !p.match(tokComma) {
			break
		}
	}
	return items, nil
}

// ---------------------------------------------------------------------------
// SET clause
// ---------------------------------------------------------------------------
// Based on goraphdb parseSetItems/parseSetItem, extended to accept expressions
// as values (not just primary/literal).

func (p *parser) parseSetItems() ([]SetItem, error) {
	var items []SetItem
	for {
		item, err := p.parseSetItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if !p.match(tokComma) {
			break
		}
	}
	return items, nil
}

func (p *parser) parseSetItem() (SetItem, error) {
	varTok, err := p.expectIdentOrKeyword()
	if err != nil {
		return SetItem{}, fmt.Errorf("cypher parser: SET expects a variable name, got %s at position %d",
			tokenKindName(p.cur().Kind), p.cur().Pos)
	}

	// SET var:Label
	if p.is(tokColon) {
		p.advance()
		labelTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return SetItem{}, fmt.Errorf("cypher parser: SET expects label name after ':'")
		}
		return SetItem{Variable: varTok.Text, SetLabel: true, Label: labelTok.Text}, nil
	}

	// SET var.prop = expr
	if p.is(tokDot) {
		p.advance()
		propTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return SetItem{}, fmt.Errorf("cypher parser: SET expects a property name after '.', got %s at position %d",
				tokenKindName(p.cur().Kind), p.cur().Pos)
		}

		if _, err := p.expect(tokEq); err != nil {
			return SetItem{}, fmt.Errorf("cypher parser: SET expects '=' after %s.%s at position %d",
				varTok.Text, propTok.Text, p.cur().Pos)
		}

		val, err := p.parseExpr()
		if err != nil {
			return SetItem{}, fmt.Errorf("cypher parser: SET value: %w", err)
		}

		return SetItem{
			Variable: varTok.Text,
			Property: propTok.Text,
			Value:    val,
		}, nil
	}

	// SET var += expr  or  SET var = expr
	if p.is(tokPlusEq) {
		p.advance()
		val, err := p.parseExpr()
		if err != nil {
			return SetItem{}, fmt.Errorf("cypher parser: SET value: %w", err)
		}
		return SetItem{
			Variable: varTok.Text,
			Value:    val,
			Merge:    true,
		}, nil
	}

	if _, err := p.expect(tokEq); err != nil {
		return SetItem{}, fmt.Errorf("cypher parser: SET expects '=' or '.property =' or '+=' after variable at position %d", p.cur().Pos)
	}

	val, err := p.parseExpr()
	if err != nil {
		return SetItem{}, fmt.Errorf("cypher parser: SET value: %w", err)
	}

	return SetItem{
		Variable:   varTok.Text,
		Value:      val,
		ReplaceAll: true,
	}, nil
}

// ---------------------------------------------------------------------------
// DELETE clause
// ---------------------------------------------------------------------------

func (p *parser) parseDeleteVars() ([]string, error) {
	var vars []string
	for {
		tok, err := p.expectIdentOrKeyword()
		if err != nil {
			return nil, fmt.Errorf("cypher parser: DELETE expects a variable name, got %s at position %d",
				tokenKindName(p.cur().Kind), p.cur().Pos)
		}
		vars = append(vars, tok.Text)
		if !p.match(tokComma) {
			break
		}
	}
	return vars, nil
}

// ---------------------------------------------------------------------------
// REMOVE clause
// ---------------------------------------------------------------------------

func (p *parser) parseRemoveItems() ([]RemoveItem, error) {
	var items []RemoveItem
	for {
		item, err := p.parseRemoveItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if !p.match(tokComma) {
			break
		}
	}
	return items, nil
}

func (p *parser) parseRemoveItem() (RemoveItem, error) {
	varTok, err := p.expectIdentOrKeyword()
	if err != nil {
		return RemoveItem{}, fmt.Errorf("cypher parser: REMOVE expects a variable name")
	}

	if p.is(tokColon) {
		p.advance()
		labelTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return RemoveItem{}, fmt.Errorf("cypher parser: REMOVE expects label name after ':'")
		}
		return RemoveItem{Variable: varTok.Text, Kind: RemoveLabel, Label: labelTok.Text}, nil
	}

	if _, err := p.expect(tokDot); err != nil {
		return RemoveItem{}, fmt.Errorf("cypher parser: REMOVE expects ':' (label) or '.' (property)")
	}
	propTok, err := p.expectIdentOrKeyword()
	if err != nil {
		return RemoveItem{}, fmt.Errorf("cypher parser: REMOVE expects property name after '.'")
	}
	return RemoveItem{Variable: varTok.Text, Kind: RemoveProperty, Property: propTok.Text}, nil
}

// ---------------------------------------------------------------------------
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

	// Check for comparison/postfix operators
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
			// IS NOT NULL
			if _, err := p.expect(tokNull); err != nil {
				return Expression{}, fmt.Errorf("cypher parser: expected NULL after IS NOT")
			}
			return isNotNullExpr(left), nil
		}
		// IS NULL
		if _, err := p.expect(tokNull); err != nil {
			return Expression{}, fmt.Errorf("cypher parser: expected NULL after IS")
		}
		return isNullExpr(left), nil
	}

	return left, nil
}

// parseAddition handles + and - (our addition — not in goraphdb).
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

// parseMultiplication handles *, /, % (our addition — not in goraphdb).
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

// parseUnary handles unary minus (our addition — not in goraphdb).
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

// parsePrimary is the innermost expression parser.
// Based on goraphdb parsePrimary, extended with CASE, EXISTS, list/map literals,
// aggregation functions, LABELS(), PROPERTIES().
func (p *parser) parsePrimary() (Expression, error) {
	t := p.cur()

	// CASE expression (our addition)
	if t.Kind == tokCase {
		return p.parseCaseExpr()
	}

	// EXISTS { pattern } (our addition)
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

	// List literal [...] (our addition)
	if t.Kind == tokLBracket {
		return p.parseListLiteral()
	}

	// Map literal {key: expr, ...} (our addition)
	if t.Kind == tokLBrace {
		return p.parseMapLiteralExpr()
	}

	// Aggregation functions (our addition)
	if t.Kind == tokCount || t.Kind == tokSum || t.Kind == tokAvg ||
		t.Kind == tokMin || t.Kind == tokMax || t.Kind == tokCollect {
		return p.parseAggregation(t.Kind)
	}

	// LABELS() and PROPERTIES() (our addition)
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

	// Identifier-led expressions
	if t.Kind == tokIdent || t.Kind == tokType {
		name := t.Text
		p.advance()

		// ident '.' ident → property access (supports chaining: n.address.city)
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

		// ident '(' args ')' → function call
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

	// Parameter reference: $paramName
	if t.Kind == tokParam {
		p.advance()
		return Expression{Kind: ExprParam, ParamName: t.Text}, nil
	}

	// Literal values
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

	// Parenthesised expression
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

// parseCaseExpr handles CASE ... WHEN ... THEN ... ELSE ... END.
// Our addition — not in goraphdb.
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

// parseListLiteral handles [...] list expressions.
// Our addition — not in goraphdb.
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

// parseMapLiteralExpr handles {key: expr, ...} map literal expressions.
// Our addition — not in goraphdb.
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

// parseAggregation handles COUNT, SUM, AVG, MIN, MAX, COLLECT.
// Our addition — not in goraphdb.
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

	// COUNT(*) has no argument
	if kind == tokCount && p.is(tokStar) {
		p.advance()
		if _, err := p.expect(tokRParen); err != nil {
			return Expression{}, err
		}
		return aggExpr(fn, nil), nil
	}

	// DISTINCT inside aggregation (e.g., COUNT(DISTINCT n))
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

// ---------------------------------------------------------------------------
// CREATE query
// ---------------------------------------------------------------------------
// Based on goraphdb parseCreateQuery, extended with DISTINCT on RETURN.

func (p *parser) parseCreateQuery() (*CypherWrite, error) {
	w := &CypherWrite{}
	if _, err := p.expect(tokCreate); err != nil {
		return nil, err
	}
	for {
		pat, err := p.parseCreatePattern()
		if err != nil {
			return nil, err
		}
		w.Creates = append(w.Creates, pat)
		if !p.match(tokComma) {
			break
		}
	}
	if p.is(tokReturn) {
		ret, err := p.parseReturnClause()
		if err != nil {
			return nil, err
		}
		w.Return = &ret

		if p.is(tokOrder) {
			p.advance()
			if _, err := p.expect(tokBy); err != nil {
				return nil, err
			}
			items, err := p.parseOrderItems()
			if err != nil {
				return nil, err
			}
			w.OrderBy = items
		}

		if p.is(tokSkip) {
			p.advance()
			tok, err := p.expectOneOf(tokInt, tokParam)
			if err != nil {
				return nil, fmt.Errorf("cypher parser: SKIP requires an integer")
			}
			if tok.Kind == tokParam {
				return nil, fmt.Errorf("cypher parser: SKIP parameter not supported in CREATE")
			}
			n, _ := strconv.Atoi(tok.Text)
			w.Skip = intPtr(n)
		}

		if p.is(tokLimit) {
			p.advance()
			tok, err := p.expectOneOf(tokInt, tokParam)
			if err != nil {
				return nil, fmt.Errorf("cypher parser: LIMIT requires an integer")
			}
			if tok.Kind == tokParam {
				return nil, fmt.Errorf("cypher parser: LIMIT parameter not supported in CREATE")
			}
			n, _ := strconv.Atoi(tok.Text)
			w.Limit = intPtr(n)
		}
	}
	if !p.is(tokEOF) {
		return nil, fmt.Errorf("cypher parser: unexpected token %s at position %d after CREATE",
			tokenKindName(p.cur().Kind), p.cur().Pos)
	}
	return w, nil
}

func (p *parser) parseCreatePattern() (CreatePattern, error) {
	cp := CreatePattern{}
	node, err := p.parseNodePattern()
	if err != nil {
		return cp, err
	}
	cp.Nodes = append(cp.Nodes, node)
	for p.is(tokDash) || p.is(tokLArrow) {
		rel, err := p.parseRelPattern()
		if err != nil {
			return cp, err
		}
		cp.Rels = append(cp.Rels, rel)
		node, err = p.parseNodePattern()
		if err != nil {
			return cp, err
		}
		cp.Nodes = append(cp.Nodes, node)
	}
	return cp, nil
}

// ---------------------------------------------------------------------------
// MERGE query
// ---------------------------------------------------------------------------
// Based on goraphdb parseMergeQuery.

func (p *parser) parseMergeQuery() (*CypherMerge, error) {
	m := &CypherMerge{}
	if _, err := p.expect(tokMerge); err != nil {
		return nil, err
	}
	node, err := p.parseNodePattern()
	if err != nil {
		return nil, fmt.Errorf("cypher parser: MERGE requires a node pattern: %w", err)
	}
	if len(node.Labels) == 0 {
		return nil, fmt.Errorf("cypher parser: MERGE node must have at least one label at position %d", p.cur().Pos)
	}
	m.Pattern = MergePattern(node)

	for p.is(tokOn) {
		p.advance()
		tok := p.cur()
		switch tok.Kind {
		case tokCreate:
			p.advance()
			if _, err := p.expect(tokSet); err != nil {
				return nil, fmt.Errorf("cypher parser: expected SET after ON CREATE at position %d", tok.Pos)
			}
			items, err := p.parseSetItems()
			if err != nil {
				return nil, err
			}
			m.OnCreateSet = items
		case tokMatch:
			p.advance()
			if _, err := p.expect(tokSet); err != nil {
				return nil, fmt.Errorf("cypher parser: expected SET after ON MATCH at position %d", tok.Pos)
			}
			items, err := p.parseSetItems()
			if err != nil {
				return nil, err
			}
			m.OnMatchSet = items
		default:
			return nil, fmt.Errorf("cypher parser: expected CREATE or MATCH after ON at position %d, got %s",
				tok.Pos, tokenKindName(tok.Kind))
		}
	}

	if p.is(tokReturn) {
		ret, err := p.parseReturnClause()
		if err != nil {
			return nil, err
		}
		m.Return = &ret

		if p.is(tokOrder) {
			p.advance()
			if _, err := p.expect(tokBy); err != nil {
				return nil, err
			}
			items, err := p.parseOrderItems()
			if err != nil {
				return nil, err
			}
			m.OrderBy = items
		}

		if p.is(tokSkip) {
			p.advance()
			tok, err := p.expectOneOf(tokInt, tokParam)
			if err != nil {
				return nil, fmt.Errorf("cypher parser: SKIP requires an integer")
			}
			if tok.Kind == tokParam {
				return nil, fmt.Errorf("cypher parser: SKIP parameter not supported in MERGE")
			}
			n, _ := strconv.Atoi(tok.Text)
			m.Skip = intPtr(n)
		}

		if p.is(tokLimit) {
			p.advance()
			tok, err := p.expectOneOf(tokInt, tokParam)
			if err != nil {
				return nil, fmt.Errorf("cypher parser: LIMIT requires an integer")
			}
			if tok.Kind == tokParam {
				return nil, fmt.Errorf("cypher parser: LIMIT parameter not supported in MERGE")
			}
			n, _ := strconv.Atoi(tok.Text)
			m.Limit = intPtr(n)
		}
	}

	if !p.is(tokEOF) {
		return nil, fmt.Errorf("cypher parser: unexpected token %s at position %d after MERGE",
			tokenKindName(p.cur().Kind), p.cur().Pos)
	}
	return m, nil
}

// ---------------------------------------------------------------------------
// DDL: CREATE INDEX, CREATE CONSTRAINT, SHOW INDEXES, SHOW CONSTRAINTS
// ---------------------------------------------------------------------------

func (p *parser) parseDDLCreate() (*ParsedCypher, error) {
	if _, err := p.expect(tokCreate); err != nil {
		return nil, err
	}
	if p.is(tokIndex) {
		p.advance()
		if p.is(tokOn) {
			p.advance()
		}
		labelTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return nil, fmt.Errorf("cypher parser: expected label name after CREATE INDEX")
		}
		var propTok Token
		if p.match(tokLParen) {
			propTok, err = p.expectIdentOrKeyword()
			if err != nil {
				return nil, fmt.Errorf("cypher parser: expected property name")
			}
			if _, expectErr := p.expect(tokRParen); expectErr != nil {
				return nil, fmt.Errorf("cypher parser: expected ')' after property name")
			}
		}
		if !p.is(tokEOF) {
			return nil, fmt.Errorf("cypher parser: unexpected token after CREATE INDEX")
		}
		return &ParsedCypher{DDL: &DDLStatement{
			Kind:     "INDEX",
			Label:    labelTok.Text,
			Property: propTok.Text,
		}}, nil
	}
	if p.is(tokConstraint) {
		p.advance()
		p.match(tokOn)
		labelTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return nil, fmt.Errorf("cypher parser: expected label name after CREATE CONSTRAINT")
		}
		if _, expectErr := p.expect(tokLParen); expectErr != nil {
			return nil, fmt.Errorf("cypher parser: expected '(' after constraint label")
		}
		propTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return nil, fmt.Errorf("cypher parser: expected property name")
		}
		if _, expectErr := p.expect(tokRParen); expectErr != nil {
			return nil, fmt.Errorf("cypher parser: expected ')' after constraint property")
		}
		if !p.is(tokEOF) {
			return nil, fmt.Errorf("cypher parser: unexpected token after CREATE CONSTRAINT")
		}
		return &ParsedCypher{DDL: &DDLStatement{
			Kind:     "CONSTRAINT",
			Label:    labelTok.Text,
			Property: propTok.Text,
		}}, nil
	}
	return nil, fmt.Errorf("cypher parser: expected INDEX or CONSTRAINT after CREATE")
}

func (p *parser) parseShow() (*ParsedCypher, error) {
	if _, err := p.expect(tokShow); err != nil {
		return nil, err
	}
	if p.is(tokIndex) || p.is(tokIndexes) {
		p.advance()
		if !p.is(tokEOF) {
			return nil, fmt.Errorf("cypher parser: unexpected token after SHOW INDEX")
		}
		return &ParsedCypher{DDL: &DDLStatement{Kind: "SHOW_INDEXES"}}, nil
	}
	if p.is(tokConstraint) || p.is(tokConstraints) {
		p.advance()
		if !p.is(tokEOF) {
			return nil, fmt.Errorf("cypher parser: unexpected token after SHOW CONSTRAINTS")
		}
		return &ParsedCypher{DDL: &DDLStatement{Kind: "SHOW_CONSTRAINTS"}}, nil
	}
	return nil, fmt.Errorf("cypher parser: expected INDEX(ES) or CONSTRAINT(S) after SHOW")
}

func (p *parser) parseCallStatement() (*CypherCall, error) {
	if _, err := p.expect(tokCall); err != nil {
		return nil, err
	}

	var nameParts []string
	for {
		t, err := p.expectIdentOrKeyword()
		if err != nil {
			return nil, fmt.Errorf("cypher parser: expected procedure name after CALL at position %d", p.cur().Pos)
		}
		nameParts = append(nameParts, t.Text)
		if !p.is(tokDot) {
			break
		}
		p.advance()
	}
	procName := strings.Join(nameParts, ".")

	if _, err := p.expect(tokLParen); err != nil {
		return nil, fmt.Errorf("cypher parser: expected '(' after procedure name at position %d", p.cur().Pos)
	}

	var args []Expression
	if !p.is(tokRParen) {
		for {
			arg, err := p.parseExpr()
			if err != nil {
				return nil, fmt.Errorf("cypher parser: failed to parse CALL argument: %w", err)
			}
			args = append(args, arg)
			if !p.is(tokComma) {
				break
			}
			p.advance()
		}
	}

	if _, err := p.expect(tokRParen); err != nil {
		return nil, fmt.Errorf("cypher parser: expected ')' after CALL arguments at position %d", p.cur().Pos)
	}

	call := &CypherCall{
		Name: procName,
		Args: args,
	}

	if p.is(tokYield) {
		p.advance()
		if p.is(tokStar) {
			p.advance()
			call.YieldItems = []string{"*"}
		} else {
			for {
				t, err := p.expectIdentOrKeyword()
				if err != nil {
					return nil, fmt.Errorf("cypher parser: expected YIELD item name at position %d", p.cur().Pos)
				}
				call.YieldItems = append(call.YieldItems, t.Text)
				if !p.is(tokComma) {
					break
				}
				p.advance()
			}
		}
	}

	return call, nil
}
