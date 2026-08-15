// Cypher DDL parsing: CREATE/DROP INDEX and CONSTRAINT, SHOW INDEXES/
// CONSTRAINTS.

package cypherparser

import (
	"fmt"
)

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

func (p *parser) parseDDLDrop() (*ParsedCypher, error) {
	if _, err := p.expect(tokDrop); err != nil {
		return nil, err
	}
	if p.is(tokIndex) {
		p.advance()
		if p.is(tokOn) {
			p.advance()
		}
		labelTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return nil, fmt.Errorf("cypher parser: expected label name after DROP INDEX")
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
			return nil, fmt.Errorf("cypher parser: unexpected token after DROP INDEX")
		}
		return &ParsedCypher{DDL: &DDLStatement{
			Kind:     "DROP_INDEX",
			Label:    labelTok.Text,
			Property: propTok.Text,
		}}, nil
	}
	if p.is(tokConstraint) {
		p.advance()
		p.match(tokOn)
		labelTok, err := p.expectIdentOrKeyword()
		if err != nil {
			return nil, fmt.Errorf("cypher parser: expected label name after DROP CONSTRAINT")
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
			return nil, fmt.Errorf("cypher parser: unexpected token after DROP CONSTRAINT")
		}
		return &ParsedCypher{DDL: &DDLStatement{
			Kind:     "DROP_CONSTRAINT",
			Label:    labelTok.Text,
			Property: propTok.Text,
		}}, nil
	}
	return nil, fmt.Errorf("cypher parser: expected INDEX or CONSTRAINT after DROP")
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
