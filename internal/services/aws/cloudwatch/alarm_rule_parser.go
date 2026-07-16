package cloudwatch

import (
	"fmt"
	"strings"
)

// alarmRuleNode is the interface implemented by all nodes in the parsed
// AlarmRule expression tree. Composite alarms reference other alarms via
// ALARM("name") terms combined with boolean operators.
type alarmRuleNode interface {
	// evaluate returns true if this sub-expression is satisfied given
	// the provided set of alarm names that are currently in ALARM state.
	evaluate(alarmStates map[string]bool) bool

	// childAlarmNames returns all alarm names referenced by this node.
	childAlarmNames() []string
}

// alarmRefNode references a single alarm by name. It evaluates to true
// when the referenced alarm is in ALARM state.
type alarmRefNode struct {
	name string
}

func (n *alarmRefNode) evaluate(alarmStates map[string]bool) bool {
	return alarmStates[n.name]
}

func (n *alarmRefNode) childAlarmNames() []string {
	return []string{n.name}
}

// trueNode always evaluates to true (the TRUE literal in AlarmRule syntax).
type trueNode struct{}

func (n *trueNode) evaluate(_ map[string]bool) bool { return true }
func (n *trueNode) childAlarmNames() []string       { return nil }

// falseNode always evaluates to false (the FALSE literal in AlarmRule syntax).
type falseNode struct{}

func (n *falseNode) evaluate(_ map[string]bool) bool { return false }
func (n *falseNode) childAlarmNames() []string       { return nil }

// andNode represents a logical AND of its children. All children must
// evaluate to true.
type andNode struct {
	children []alarmRuleNode
}

func (n *andNode) evaluate(alarmStates map[string]bool) bool {
	for _, c := range n.children {
		if !c.evaluate(alarmStates) {
			return false
		}
	}
	return true
}

func (n *andNode) childAlarmNames() []string {
	var names []string
	for _, c := range n.children {
		names = append(names, c.childAlarmNames()...)
	}
	return names
}

// orNode represents a logical OR of its children. At least one child must
// evaluate to true.
type orNode struct {
	children []alarmRuleNode
}

func (n *orNode) evaluate(alarmStates map[string]bool) bool {
	for _, c := range n.children {
		if c.evaluate(alarmStates) {
			return true
		}
	}
	return false
}

func (n *orNode) childAlarmNames() []string {
	var names []string
	for _, c := range n.children {
		names = append(names, c.childAlarmNames()...)
	}
	return names
}

// notNode represents a logical NOT of its single child.
type notNode struct {
	child alarmRuleNode
}

func (n *notNode) evaluate(alarmStates map[string]bool) bool {
	return !n.child.evaluate(alarmStates)
}

func (n *notNode) childAlarmNames() []string {
	return n.child.childAlarmNames()
}

// --- Tokenizer ---

type ruleTokenType int

const (
	tokAlarm ruleTokenType = iota
	tokAnd
	tokOr
	tokNot
	tokTrue
	tokFalse
	tokLParen
	tokRParen
	tokString
	tokEOF
)

type ruleToken struct {
	typ ruleTokenType
	val string
	pos int
}

func tokenizeAlarmRule(input string) ([]ruleToken, error) {
	var tokens []ruleToken
	i := 0
	for i < len(input) {
		ch := input[i]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
			i++
		case ch == '(':
			tokens = append(tokens, ruleToken{tokLParen, "(", i})
			i++
		case ch == ')':
			tokens = append(tokens, ruleToken{tokRParen, ")", i})
			i++
		case ch == '"':
			start := i + 1
			i++
			var sb strings.Builder
			for i < len(input) && input[i] != '"' {
				sb.WriteByte(input[i])
				i++
			}
			if i >= len(input) {
				return nil, fmt.Errorf("alarm rule: unterminated string at position %d", start-1)
			}
			i++
			tokens = append(tokens, ruleToken{tokString, sb.String(), start})
		default:
			start := i
			for i < len(input) && input[i] != ' ' && input[i] != '\t' && input[i] != '\n' &&
				input[i] != '\r' && input[i] != '(' && input[i] != ')' && input[i] != '"' {
				i++
			}
			word := input[start:i]
			upper := strings.ToUpper(word)
			switch upper {
			case "ALARM":
				tokens = append(tokens, ruleToken{tokAlarm, word, start})
			case "AND":
				tokens = append(tokens, ruleToken{tokAnd, word, start})
			case "OR":
				tokens = append(tokens, ruleToken{tokOr, word, start})
			case "NOT":
				tokens = append(tokens, ruleToken{tokNot, word, start})
			case "TRUE":
				tokens = append(tokens, ruleToken{tokTrue, word, start})
			case "FALSE":
				tokens = append(tokens, ruleToken{tokFalse, word, start})
			default:
				return nil, fmt.Errorf("alarm rule: unexpected token %q at position %d", word, start)
			}
		}
	}
	tokens = append(tokens, ruleToken{tokEOF, "", len(input)})
	return tokens, nil
}

// --- Parser (recursive descent) ---

type ruleParser struct {
	tokens []ruleToken
	pos    int
}

func parseAlarmRule(input string) (alarmRuleNode, error) {
	tokens, err := tokenizeAlarmRule(input)
	if err != nil {
		return nil, err
	}
	p := &ruleParser{tokens: tokens}
	node, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if p.peek().typ != tokEOF {
		return nil, fmt.Errorf("alarm rule: unexpected token %q at position %d", p.peek().val, p.peek().pos)
	}
	return node, nil
}

func (p *ruleParser) peek() ruleToken {
	return p.tokens[p.pos]
}

func (p *ruleParser) consume() ruleToken {
	t := p.tokens[p.pos]
	p.pos++
	return t
}

func (p *ruleParser) parseOr() (alarmRuleNode, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	children := []alarmRuleNode{left}
	for p.peek().typ == tokOr {
		p.consume()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		children = append(children, right)
	}
	if len(children) == 1 {
		return children[0], nil
	}
	return &orNode{children: children}, nil
}

func (p *ruleParser) parseAnd() (alarmRuleNode, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	children := []alarmRuleNode{left}
	for p.peek().typ == tokAnd {
		p.consume()
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		children = append(children, right)
	}
	if len(children) == 1 {
		return children[0], nil
	}
	return &andNode{children: children}, nil
}

func (p *ruleParser) parseNot() (alarmRuleNode, error) {
	if p.peek().typ == tokNot {
		p.consume()
		child, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &notNode{child: child}, nil
	}
	return p.parsePrimary()
}

func (p *ruleParser) parsePrimary() (alarmRuleNode, error) {
	tok := p.peek()
	switch tok.typ {
	case tokAlarm:
		p.consume()
		if p.peek().typ != tokLParen {
			return nil, fmt.Errorf("alarm rule: expected '(' after ALARM at position %d", tok.pos)
		}
		p.consume()
		if p.peek().typ != tokString {
			return nil, fmt.Errorf("alarm rule: expected alarm name string at position %d", p.peek().pos)
		}
		nameTok := p.consume()
		if p.peek().typ != tokRParen {
			return nil, fmt.Errorf("alarm rule: expected ')' after alarm name at position %d", p.peek().pos)
		}
		p.consume()
		return &alarmRefNode{name: nameTok.val}, nil

	case tokTrue:
		p.consume()
		return &trueNode{}, nil

	case tokFalse:
		p.consume()
		return &falseNode{}, nil

	case tokLParen:
		p.consume()
		node, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if p.peek().typ != tokRParen {
			return nil, fmt.Errorf("alarm rule: expected ')' at position %d", p.peek().pos)
		}
		p.consume()
		return node, nil

	default:
		return nil, fmt.Errorf("alarm rule: unexpected token %q at position %d", tok.val, tok.pos)
	}
}
