// Cypher CREATE and MERGE query parsing.

package cypherparser

import (
	"fmt"
	"strconv"
)

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
