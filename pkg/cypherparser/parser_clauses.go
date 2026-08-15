// Cypher clause parsing: RETURN, SET, DELETE, REMOVE and CALL.

package cypherparser

import (
	"fmt"
	"strings"
)

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

	if p.is(tokColon) {
		p.advance()
		labelTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return SetItem{}, fmt.Errorf("cypher parser: SET expects label name after ':'")
		}
		return SetItem{Variable: varTok.Text, SetLabel: true, Label: labelTok.Text}, nil
	}

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
