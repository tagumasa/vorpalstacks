// Cypher pattern parsing: node patterns, relationship patterns, and property maps.
//
// Based on goraphdb/cypher_parser.go parsePattern/parseNodePattern/parseRelPattern.
// Extended with edge properties in RelPattern and expression-valued property maps
// (goraphdb only supported literal values).

package cypherparser

import (
	"fmt"
	"strconv"

	"vorpalstacks/internal/core/storage/graphengine"
)

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

			node, err := p.parseNodePattern()
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

		if p.is(tokParam) {
			paramTok := p.advance()
			props[keyTok.Text] = paramRef(paramTok.Text)
			continue
		}

		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
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
