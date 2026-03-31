package gremlinparser

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	tokens []Token
	pos    int
}

func Parse(input string) (*Query, error) {
	tokens, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	p := &Parser{tokens: tokens, pos: 0}
	return p.parseQuery()
}

func (p *Parser) cur() Token {
	if p.pos >= len(p.tokens) {
		return Token{Kind: tokEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peek(offset int) Token {
	idx := p.pos + offset
	if idx >= len(p.tokens) {
		return Token{Kind: tokEOF}
	}
	return p.tokens[idx]
}

func (p *Parser) advance() Token {
	t := p.cur()
	if t.Kind != tokEOF {
		p.pos++
	}
	return t
}

func (p *Parser) expect(kind TokenKind) (Token, error) {
	t := p.cur()
	if t.Kind != kind {
		return t, fmt.Errorf("gremlin: expected token kind %d, got %d (%q) at position %d", kind, t.Kind, t.Text, t.Pos)
	}
	return p.advance(), nil
}

func (p *Parser) match(kind TokenKind) (Token, bool) {
	if p.cur().Kind == kind {
		return p.advance(), true
	}
	return Token{}, false
}

func (p *Parser) matchIdent(text string) bool {
	t := p.cur()
	if t.Kind == tokIdentifier && t.Text == text {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) isIdent(text string) bool {
	t := p.cur()
	return t.Kind == tokIdentifier && t.Text == text
}

func (p *Parser) parseQuery() (*Query, error) {
	q := &Query{}

	for p.cur().Kind != tokEOF {
		if tok, ok := p.match(tokSemicolon); ok {
			_ = tok
			continue
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			q.Statements = append(q.Statements, *stmt)
		}
	}

	return q, nil
}

func (p *Parser) parseStatement() (*Statement, error) {
	if p.cur().Kind == tokEOF {
		return nil, nil
	}

	ident := p.cur()
	if ident.Kind == tokIdentifier && ident.Text == "g" {
		return p.parseTraversalStatement()
	}

	return nil, fmt.Errorf("gremlin: expected traversal starting with 'g' or '__', got %q at position %d", ident.Text, ident.Pos)
}

func (p *Parser) parseTraversalStatement() (*Statement, error) {
	g := p.advance()
	if g.Text != "g" {
		return nil, fmt.Errorf("gremlin: expected 'g', got %q", g.Text)
	}

	if _, ok := p.match(tokDot); !ok {
		return nil, fmt.Errorf("gremlin: expected '.' after 'g' at position %d", g.Pos)
	}

	trav, err := p.parseChainedTraversal()
	if err != nil {
		return nil, err
	}

	if p.cur().Kind == tokDot {
		p.advance()
		term, err := p.parseTerminal()
		if err != nil {
			return nil, err
		}
		trav.Terminal = term
	}

	return &Statement{Traversal: trav}, nil
}

func (p *Parser) parseChainedTraversal() (*Traversal, error) {
	source := ""
	if p.cur().Kind == tokUnderscoreUnderscore {
		p.advance()
		if _, ok := p.match(tokDot); !ok {
			return nil, fmt.Errorf("gremlin: expected '.' after '__'")
		}
		source = "__"
	}

	trav := &Traversal{Source: source}

	for {
		if p.cur().Kind != tokIdentifier {
			break
		}

		stepName := p.cur().Text
		if isTerminalMethod(stepName) {
			break
		}

		if isModulatorStep(stepName) && len(trav.Steps) > 0 {
			p.advance()

			if _, ok := p.match(tokDot); ok {
				lastStep := &trav.Steps[len(trav.Steps)-1]
				lastStep.Modulators = append(lastStep.Modulators, Step{Name: stepName, Args: nil})
				continue
			}

			if _, ok := p.match(tokLParen); !ok {
				return nil, fmt.Errorf("gremlin: expected '(' or '.' after modulator step %q", stepName)
			}

			args, err := p.parseArguments()
			if err != nil {
				return nil, err
			}

			if _, err := p.expect(tokRParen); err != nil {
				return nil, fmt.Errorf("gremlin: expected ')' after arguments for modulator step %q", stepName)
			}

			lastStep := &trav.Steps[len(trav.Steps)-1]
			lastStep.Modulators = append(lastStep.Modulators, Step{Name: stepName, Args: args})

			if p.cur().Kind == tokDot && p.peek(1).Kind == tokIdentifier && isTerminalMethod(p.peek(1).Text) {
				break
			}

			if _, ok := p.match(tokDot); !ok {
				break
			}
			continue
		}

		p.advance()

		if _, ok := p.match(tokDot); ok {
			trav.Steps = append(trav.Steps, Step{Name: stepName, Args: nil})
			continue
		}

		if _, ok := p.match(tokLParen); !ok {
			return nil, fmt.Errorf("gremlin: expected '(' or '.' after step %q", stepName)
		}

		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(tokRParen); err != nil {
			return nil, fmt.Errorf("gremlin: expected ')' after arguments for step %q", stepName)
		}

		trav.Steps = append(trav.Steps, Step{Name: stepName, Args: args})

		if p.cur().Kind == tokDot && p.peek(1).Kind == tokIdentifier && isTerminalMethod(p.peek(1).Text) {
			break
		}

		if _, ok := p.match(tokDot); !ok {
			break
		}
	}

	return trav, nil
}

func (p *Parser) parseArguments() ([]Argument, error) {
	if p.cur().Kind == tokRParen {
		return nil, nil
	}

	var args []Argument
	for {
		arg, err := p.parseArgument()
		if err != nil {
			return nil, err
		}
		args = append(args, *arg)

		if _, ok := p.match(tokComma); !ok {
			break
		}
	}
	return args, nil
}

func (p *Parser) parseArgument() (*Argument, error) {
	t := p.cur()

	switch t.Kind {
	case tokString:
		p.advance()
		return &Argument{Kind: ArgString, Str: t.Text}, nil

	case tokInteger:
		p.advance()
		val, err := parseInt(t.Text)
		if err != nil {
			return nil, err
		}
		return &Argument{Kind: ArgInt, Int: val}, nil

	case tokFloat:
		p.advance()
		val, err := strconv.ParseFloat(strings.TrimRightFunc(t.Text, func(r rune) bool {
			return r == 'f' || r == 'F' || r == 'd' || r == 'D' || r == 'm' || r == 'M'
		}), 64)
		if err != nil {
			return nil, fmt.Errorf("gremlin: invalid float %q: %w", t.Text, err)
		}
		return &Argument{Kind: ArgFloat, Float: val}, nil

	case tokTrue:
		p.advance()
		return &Argument{Kind: ArgBool, Bool: true}, nil

	case tokFalse:
		p.advance()
		return &Argument{Kind: ArgBool, Bool: false}, nil

	case tokNull:
		p.advance()
		return &Argument{Kind: ArgNull}, nil

	case tokLBracket:
		return p.parseListOrSetOrRange()

	case tokLBrace:
		return p.parseSetLiteral()

	case tokUnderscoreUnderscore:
		return p.parseNestedTraversal()

	case tokIdentifier:
		return p.parseIdentArgument()

	default:
		return nil, fmt.Errorf("gremlin: unexpected token %q (%d) in argument at position %d", t.Text, t.Kind, t.Pos)
	}
}

func (p *Parser) parseIdentArgument() (*Argument, error) {
	t := p.cur()

	if t.Text == "P" && p.peek(1).Kind == tokDot {
		return p.parsePredicate()
	}

	if t.Text == "TextP" && p.peek(1).Kind == tokDot {
		return p.parseTextPredicate()
	}

	if isEnumPrefix(t.Text) && p.peek(1).Kind == tokDot {
		return p.parseEnum()
	}

	if isBareEnum(t.Text) && !(isBarePredicate(t.Text) && p.peek(1).Kind == tokLParen) {
		p.advance()
		return &Argument{Kind: ArgEnum, Enum: &EnumRef{Value: t.Text}}, nil
	}

	if isBarePredicate(t.Text) && p.peek(1).Kind == tokLParen {
		return p.parseBarePredicate()
	}

	if (p.peek(1).Kind == tokLParen || p.peek(1).Kind == tokDot) && !isTerminalMethod(t.Text) {
		return p.parseBareTraversal()
	}

	p.advance()
	return &Argument{Kind: ArgString, Str: t.Text}, nil
}

func (p *Parser) parseBareTraversal() (*Argument, error) {
	if p.cur().Kind == tokIdentifier && p.cur().Text == "g" && p.peek(1).Kind == tokDot {
		p.advance()
		p.advance()
	}
	trav, err := p.parseChainedTraversal()
	if err != nil {
		return nil, err
	}
	return &Argument{Kind: ArgNestedTraversal, Trav: trav}, nil
}

func (p *Parser) parseBarePredicate() (*Argument, error) {
	method := p.cur()
	p.advance()

	predType := "P"
	if isTextPredicate(method.Text) {
		predType = "TextP"
	}

	if _, err := p.expect(tokLParen); err != nil {
		return nil, fmt.Errorf("gremlin: expected '(' after bare predicate %s", method.Text)
	}

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return nil, fmt.Errorf("gremlin: expected ')' in %s()", method.Text)
	}

	pred := &Predicate{Type: predType, Method: method.Text, Args: args}
	return &Argument{Kind: ArgPredicate, Pred: pred}, nil
}

func (p *Parser) parsePredicate() (*Argument, error) {
	p.advance()
	if _, err := p.expect(tokDot); err != nil {
		return nil, fmt.Errorf("gremlin: expected '.' after 'P'")
	}
	method := p.cur()
	if method.Kind != tokIdentifier {
		return nil, fmt.Errorf("gremlin: expected predicate method after 'P.', got %q", method.Text)
	}
	p.advance()

	if _, ok := p.match(tokLParen); !ok {
		return nil, fmt.Errorf("gremlin: expected '(' after P.%s", method.Text)
	}

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return nil, fmt.Errorf("gremlin: expected ')' in P.%s()", method.Text)
	}

	pred := &Predicate{Type: "P", Method: method.Text, Args: args}

	for {
		if p.isIdent("and") && p.peek(1).Kind == tokLParen {
			p.advance()
			if _, err := p.expect(tokLParen); err != nil {
				return nil, fmt.Errorf("gremlin: expected '(' after .and")
			}
			innerArg, err := p.parseArgument()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(tokRParen); err != nil {
				return nil, fmt.Errorf("gremlin: expected ')' in .and()")
			}
			var innerPred *Predicate
			if innerArg.Kind == ArgPredicate {
				innerPred = innerArg.Pred
			} else {
				innerPred = &Predicate{Type: "P", Method: "eq", Args: []Argument{*innerArg}}
			}
			pred = &Predicate{Type: pred.Type, Method: "and", Args: []Argument{{Kind: ArgPredicate, Pred: pred}, {Kind: ArgPredicate, Pred: innerPred}}}
			continue
		}
		if p.isIdent("or") && p.peek(1).Kind == tokLParen {
			p.advance()
			if _, err := p.expect(tokLParen); err != nil {
				return nil, fmt.Errorf("gremlin: expected '(' after .or")
			}
			innerArg, err := p.parseArgument()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(tokRParen); err != nil {
				return nil, fmt.Errorf("gremlin: expected ')' in .or()")
			}
			var innerPred *Predicate
			if innerArg.Kind == ArgPredicate {
				innerPred = innerArg.Pred
			} else {
				innerPred = &Predicate{Type: "P", Method: "eq", Args: []Argument{*innerArg}}
			}
			pred = &Predicate{Type: pred.Type, Method: "or", Args: []Argument{{Kind: ArgPredicate, Pred: pred}, {Kind: ArgPredicate, Pred: innerPred}}}
			continue
		}
		if p.isIdent("negate") && p.peek(1).Kind == tokLParen {
			p.advance()
			if _, err := p.expect(tokLParen); err != nil {
				return nil, fmt.Errorf("gremlin: expected '(' after .negate")
			}
			if _, err := p.expect(tokRParen); err != nil {
				return nil, fmt.Errorf("gremlin: expected ')' in .negate()")
			}
			pred.Negate = true
			continue
		}
		break
	}

	return &Argument{Kind: ArgPredicate, Pred: pred}, nil
}

func (p *Parser) parseTextPredicate() (*Argument, error) {
	p.advance()
	if _, err := p.expect(tokDot); err != nil {
		return nil, fmt.Errorf("gremlin: expected '.' after 'TextP'")
	}
	method := p.cur()
	if method.Kind != tokIdentifier {
		return nil, fmt.Errorf("gremlin: expected text predicate method after 'TextP.', got %q", method.Text)
	}
	p.advance()

	if _, ok := p.match(tokLParen); !ok {
		return nil, fmt.Errorf("gremlin: expected '(' after TextP.%s", method.Text)
	}

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return nil, fmt.Errorf("gremlin: expected ')' in TextP.%s()", method.Text)
	}

	pred := &Predicate{Type: "TextP", Method: method.Text, Args: args}
	return &Argument{Kind: ArgPredicate, Pred: pred}, nil
}

func (p *Parser) parseEnum() (*Argument, error) {
	category := p.cur().Text
	p.advance()
	if _, err := p.expect(tokDot); err != nil {
		return nil, fmt.Errorf("gremlin: expected '.' after enum category %q", category)
	}
	value := p.cur()
	if value.Kind != tokIdentifier {
		return nil, fmt.Errorf("gremlin: expected enum value after '%s.', got %q", category, value.Text)
	}
	p.advance()

	if _, ok := p.match(tokLParen); ok {
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return nil, fmt.Errorf("gremlin: expected ')' in %s.%s()", category, value.Text)
		}
		if len(args) == 1 {
			return &Argument{Kind: ArgEnum, Enum: &EnumRef{Category: category, Value: value.Text}, List: []Argument{args[0]}}, nil
		}
		return &Argument{Kind: ArgEnum, Enum: &EnumRef{Category: category, Value: value.Text}, List: args}, nil
	}

	return &Argument{Kind: ArgEnum, Enum: &EnumRef{Category: category, Value: value.Text}}, nil
}

func (p *Parser) parseNestedTraversal() (*Argument, error) {
	p.advance()
	if _, err := p.expect(tokDot); err != nil {
		return nil, fmt.Errorf("gremlin: expected '.' after '__'")
	}

	trav, err := p.parseChainedTraversal()
	if err != nil {
		return nil, err
	}
	trav.Source = "__"

	return &Argument{Kind: ArgNestedTraversal, Trav: trav}, nil
}

func (p *Parser) parseListOrSetOrRange() (*Argument, error) {
	p.advance()

	if p.cur().Kind == tokColon {
		p.advance()
		if _, err := p.expect(tokRBracket); err != nil {
			return nil, fmt.Errorf("gremlin: expected ']' for empty map literal")
		}
		return &Argument{Kind: ArgMap, Map: nil}, nil
	}

	if p.cur().Kind == tokRBracket {
		p.advance()
		return &Argument{Kind: ArgList, List: nil}, nil
	}

	first, err := p.parseArgument()
	if err != nil {
		return nil, err
	}

	if p.cur().Kind == tokColon {
		p.advance()
		second, err := p.parseArgument()
		if err != nil {
			return nil, err
		}

		entries := []MapEntry{{Key: *first, Value: *second}}
		for {
			if _, ok := p.match(tokComma); !ok {
				break
			}
			k, err := p.parseArgument()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(tokColon); err != nil {
				return nil, fmt.Errorf("gremlin: expected ':' in map entry")
			}
			v, err := p.parseArgument()
			if err != nil {
				return nil, err
			}
			entries = append(entries, MapEntry{Key: *k, Value: *v})
		}

		if _, err := p.expect(tokRBracket); err != nil {
			return nil, fmt.Errorf("gremlin: expected ']' to close map literal")
		}
		return &Argument{Kind: ArgMap, Map: entries}, nil
	}

	items := []Argument{*first}
	for {
		if _, ok := p.match(tokComma); !ok {
			break
		}
		item, err := p.parseArgument()
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	if _, err := p.expect(tokRBracket); err != nil {
		return nil, fmt.Errorf("gremlin: expected ']' to close list literal")
	}
	return &Argument{Kind: ArgList, List: items}, nil
}

func (p *Parser) parseSetLiteral() (*Argument, error) {
	p.advance()

	if p.cur().Kind == tokRBrace {
		p.advance()
		return &Argument{Kind: ArgList, List: nil}, nil
	}

	items := []Argument{}
	for {
		item, err := p.parseArgument()
		if err != nil {
			return nil, err
		}
		items = append(items, *item)

		if _, ok := p.match(tokComma); !ok {
			break
		}
	}

	if _, err := p.expect(tokRBrace); err != nil {
		return nil, fmt.Errorf("gremlin: expected '}' to close set literal")
	}
	return &Argument{Kind: ArgList, List: items}, nil
}

func (p *Parser) parseTerminal() (*TerminalStep, error) {
	t := p.cur()
	if t.Kind != tokIdentifier {
		return nil, fmt.Errorf("gremlin: expected terminal method, got %q at position %d", t.Text, t.Pos)
	}
	p.advance()

	if _, ok := p.match(tokLParen); !ok {
		return &TerminalStep{Name: t.Text}, nil
	}

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokRParen); err != nil {
		return nil, fmt.Errorf("gremlin: expected ')' in terminal %q", t.Text)
	}

	return &TerminalStep{Name: t.Text, Args: args}, nil
}

func isTerminalMethod(name string) bool {
	switch name {
	case "toList", "toSet", "toBulkSet", "next", "iterate", "hasNext", "tryNext", "explain":
		return true
	}
	return false
}

func isModulatorStep(name string) bool {
	switch name {
	case "by", "as", "emit", "until", "times", "option", "with", "from", "to", "without", "scope":
		return true
	}
	return false
}

func isEnumPrefix(name string) bool {
	switch name {
	case "T", "Order", "Scope", "Cardinality", "Pop", "Column", "Direction",
		"Operator", "Merge", "Pick", "DT", "GType",
		"WithOptions", "ConnectedComponent", "ShortestPath":
		return true
	}
	return false
}

func isBareEnum(name string) bool {
	switch name {
	case "id", "label", "key", "value",
		"asc", "desc", "shuffle",
		"local", "global",
		"single", "set", "list",
		"first", "last", "all", "mixed",
		"keys", "values",
		"in", "out", "both",
		"onCreate", "onMatch",
		"any", "none", "unproductive":
		return true
	}
	return false
}

func isBarePredicate(name string) bool {
	switch name {
	case "eq", "neq", "gt", "gte", "lt", "lte",
		"between", "inside", "outside", "within", "without",
		"containing", "notContaining",
		"startingWith", "notStartingWith",
		"endingWith", "notEndingWith",
		"regex", "notRegex":
		return true
	}
	return false
}

func isTextPredicate(name string) bool {
	switch name {
	case "containing", "notContaining",
		"startingWith", "notStartingWith",
		"endingWith", "notEndingWith",
		"regex", "notRegex":
		return true
	}
	return false
}

func parseInt(text string) (int64, error) {
	trimmed := strings.TrimRightFunc(text, func(r rune) bool {
		return r == 'i' || r == 'I' || r == 's' || r == 'S' || r == 'l' || r == 'L'
	})
	val, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}
