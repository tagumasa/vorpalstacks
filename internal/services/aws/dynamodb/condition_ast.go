package dynamodb

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// The condition grammar, one parse away from the wire string.
//
// DynamoDB answers a ConditionExpression or FilterExpression that does not
// parse, references an undefined substitution, or uses an operator outside
// the grammar with a ValidationException before any item is read — a value
// placeholder or attribute alias used in an expression must be defined
// (developer guide, Expression attribute values / names pages, cached).
// This file is that contract for the condition plane: the expression is
// lexed, parsed to a typed tree, and validated once per request; per-item
// evaluation then walks the tree and cannot fail. Absence semantics — an
// operand path that addresses nothing, or an operand of a type the
// operator does not order — make the condition false, never an error:
// absence belongs to evaluation, malformation to validation.
//
// Grammar (operators and functions page, cached):
//
//	condition  := orExpr
//	orExpr     := andExpr ( OR andExpr )*
//	andExpr    := notExpr ( AND notExpr )*
//	notExpr    := NOT notExpr | primary
//	primary    := '(' orExpr ')' | operand ( comparator operand
//	                                            | BETWEEN operand AND operand
//	                                            | IN '(' operand (',' operand)* ')' )
//	          |  attribute_exists '(' path ')'
//	          |  attribute_not_exists '(' path ')'
//	          |  attribute_type '(' path ',' typeOperand ')'
//	          |  begins_with '(' path ',' operand ')'
//	          |  contains '(' path ',' operand ')'
//	operand    := path | size '(' path ')' | placeholder
//	comparator := '=' | '<>' | '<' | '<=' | '>' | '>='
//
// Keywords (AND, OR, NOT, BETWEEN, IN) match case-insensitively; function
// names are case-sensitive (documented). An identifier that case-folds to a
// keyword is still a path when it stands in operand position — the grammar
// recognises keywords only where a keyword is expected, so an attribute
// named "or" remains addressable. Parentheses group conditions, never
// operands. String literals are not operands of this grammar; the one
// tolerated exception is attribute_type's type argument, kept for the
// quoted-literal form the platform has always accepted.

// conditionValidationError is the ValidationException every parse or
// substitution failure answers with. The error class is the documented
// contract; the prose is the platform's own, in AWS's message style.
func conditionValidationError(format string, args ...interface{}) error {
	return NewAPIError("com.amazon.coral.validate#ValidationException",
		fmt.Sprintf(format, args...), http.StatusBadRequest)
}

// ---------------------------------------------------------------------------
// Lexer
// ---------------------------------------------------------------------------

type condTokenKind int

const (
	condTokenPath        condTokenKind = iota // attribute name, document path, keyword or function name
	condTokenPlaceholder                      // :name substitution
	condTokenString                           // 'literal'
	condTokenOperator                         // = <> < <= > >=
	condTokenLParen
	condTokenRParen
	condTokenComma
	condTokenEOF
)

type condToken struct {
	kind condTokenKind
	text string
}

// isCondPathByte reports whether a raw byte continues a path or placeholder
// token. Every delimiter the grammar gives meaning to — whitespace, the
// comparison operators, parentheses, commas and quotes — is ASCII, so the
// complement set accumulates byte-wise and multi-byte UTF-8 runes pass
// through unchanged. A colon only starts a placeholder at token start;
// inside a token it is an ordinary path character.
func isCondPathByte(ch byte) bool {
	switch ch {
	case ' ', '\t', '\n', '\r',
		'=', '<', '>',
		'(', ')', ',', '\'':
		return false
	}
	return true
}

// lexCondition splits the expression into grammar tokens.
func lexCondition(expr string) ([]condToken, error) {
	var tokens []condToken
	i := 0
	for i < len(expr) {
		ch := expr[i]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
			i++
		case ch == '(':
			tokens = append(tokens, condToken{condTokenLParen, "("})
			i++
		case ch == ')':
			tokens = append(tokens, condToken{condTokenRParen, ")"})
			i++
		case ch == ',':
			tokens = append(tokens, condToken{condTokenComma, ","})
			i++
		case ch == ':' && i+1 < len(expr) && isCondPathByte(expr[i+1]):
			start := i
			i++
			for i < len(expr) && isCondPathByte(expr[i]) {
				i++
			}
			tokens = append(tokens, condToken{condTokenPlaceholder, expr[start:i]})
		case ch == '=':
			tokens = append(tokens, condToken{condTokenOperator, "="})
			i++
		case ch == '<':
			if i+1 < len(expr) && (expr[i+1] == '=' || expr[i+1] == '>') {
				tokens = append(tokens, condToken{condTokenOperator, expr[i : i+2]})
				i += 2
			} else {
				tokens = append(tokens, condToken{condTokenOperator, "<"})
				i++
			}
		case ch == '>':
			if i+1 < len(expr) && expr[i+1] == '=' {
				tokens = append(tokens, condToken{condTokenOperator, ">="})
				i += 2
			} else {
				tokens = append(tokens, condToken{condTokenOperator, ">"})
				i++
			}
		case ch == '\'':
			// A literal runs to the closing quote; the grammar has no
			// escape form. Unterminated, it is a syntax error here rather
			// than a token that silently swallows the rest of the
			// expression.
			j := i + 1
			for j < len(expr) && expr[j] != '\'' {
				j++
			}
			if j == len(expr) {
				return nil, conditionValidationError("Invalid ConditionExpression: Syntax error; unterminated literal near %q", expr[i:])
			}
			tokens = append(tokens, condToken{condTokenString, expr[i : j+1]})
			i = j + 1
		default:
			start := i
			for i < len(expr) && isCondPathByte(expr[i]) {
				i++
			}
			tokens = append(tokens, condToken{condTokenPath, expr[start:i]})
		}
	}
	tokens = append(tokens, condToken{condTokenEOF, ""})
	return tokens, nil
}

// ---------------------------------------------------------------------------
// Tree
// ---------------------------------------------------------------------------

// condNode is one node of the parsed condition tree. Evaluation cannot
// fail: every malformation was rejected at parse time.
type condNode interface {
	eval(item *dbstore.Item) bool
}

type condOrNode struct{ operands []condNode }

func (n *condOrNode) eval(item *dbstore.Item) bool {
	for _, operand := range n.operands {
		if operand.eval(item) {
			return true
		}
	}
	return false
}

type condAndNode struct{ operands []condNode }

func (n *condAndNode) eval(item *dbstore.Item) bool {
	for _, operand := range n.operands {
		if !operand.eval(item) {
			return false
		}
	}
	return true
}

type condNotNode struct{ operand condNode }

func (n *condNotNode) eval(item *dbstore.Item) bool {
	return !n.operand.eval(item)
}

// condOperand is one operand position of the grammar: a document path, the
// size of one, or a substitution value. value reports nil when the operand
// addresses nothing (absence) — the condition forms above turn that into
// false.
type condOperand interface {
	value(item *dbstore.Item) *dbstore.AttributeValue
}

// operandPath resolves a parsed document path against the item.
type operandPath struct{ parts []docPathPart }

func (o operandPath) value(item *dbstore.Item) *dbstore.AttributeValue {
	return getDocPathValue(item.Attributes, o.parts)
}

// operandSize is the size() form: the value is the operand's length as a
// number attribute. A path that addresses nothing, or addresses a value
// with no length (a number, a boolean, null), has no size to compare and
// reports nil.
type operandSize struct{ parts []docPathPart }

func (o operandSize) value(item *dbstore.Item) *dbstore.AttributeValue {
	attr := getDocPathValue(item.Attributes, o.parts)
	if attr == nil {
		return nil
	}
	size, ok := computeAttributeSize(attr)
	if !ok {
		return nil
	}
	sizeStr := strconv.Itoa(size)
	return &dbstore.AttributeValue{N: &sizeStr}
}

// operandValue is a substitution placeholder, resolved to its concrete
// value at parse time — an undefined placeholder is a validation error,
// never a silently false condition.
type operandValue struct{ v *dbstore.AttributeValue }

func (o operandValue) value(*dbstore.Item) *dbstore.AttributeValue { return o.v }

type condCompareNode struct {
	op       string
	lhs, rhs condOperand
}

func (n *condCompareNode) eval(item *dbstore.Item) bool {
	lhs := n.lhs.value(item)
	if lhs == nil {
		return false
	}
	rhs := n.rhs.value(item)
	if rhs == nil {
		return false
	}
	return compareAttributeValues(lhs, n.op, rhs)
}

type condBetweenNode struct{ operand, low, high condOperand }

func (n *condBetweenNode) eval(item *dbstore.Item) bool {
	v := n.operand.value(item)
	if v == nil {
		return false
	}
	low := n.low.value(item)
	high := n.high.value(item)
	if low == nil || high == nil {
		return false
	}
	return compareAttributeValues(v, ">=", low) && compareAttributeValues(v, "<=", high)
}

type condInNode struct {
	operand    condOperand
	candidates []condOperand
}

func (n *condInNode) eval(item *dbstore.Item) bool {
	v := n.operand.value(item)
	if v == nil {
		return false
	}
	for _, candidate := range n.candidates {
		cv := candidate.value(item)
		if cv != nil && attributeValuesEqual(v, cv) {
			return true
		}
	}
	return false
}

// condLegacyInNode is the legacy-parameter IN. The documented legacy
// semantics are per-element membership — a scalar candidate matches when
// it is a member of a set attribute ("If any elements of the input set
// are present in the item attribute, the expression evaluates to true")
// — while a non-set attribute keeps the whole-value equality-any
// reading. The expression IN (condInNode) compares whole values only.
type condLegacyInNode struct {
	operand    condOperand
	candidates []condOperand
}

func (n *condLegacyInNode) eval(item *dbstore.Item) bool {
	v := n.operand.value(item)
	if v == nil {
		return false
	}
	for _, candidate := range n.candidates {
		cv := candidate.value(item)
		if cv != nil && legacyInMatches(v, cv) {
			return true
		}
	}
	return false
}

// legacyInMatches answers one legacy-IN candidate against an attribute: a
// set attribute matches a same-kind scalar candidate that equals any of
// its members; every other attribute keeps whole-value equality.
func legacyInMatches(attr, candidate *dbstore.AttributeValue) bool {
	if setContainsMember(attr, candidate) {
		return true
	}
	return attributeValuesEqual(attr, candidate)
}

// condFuncKind enumerates the condition functions beyond the comparator
// forms.
type condFuncKind int

const (
	condFuncExists condFuncKind = iota
	condFuncNotExists
	condFuncBegins
	condFuncContains
)

// condFuncNode is attribute_exists, attribute_not_exists, begins_with or
// contains. The first argument is a document path; the second (begins_with
// and contains only) is a full operand.
type condFuncNode struct {
	kind condFuncKind
	path []docPathPart
	arg  condOperand
}

func (n *condFuncNode) eval(item *dbstore.Item) bool {
	switch n.kind {
	case condFuncExists:
		return getDocPathValue(item.Attributes, n.path) != nil
	case condFuncNotExists:
		return getDocPathValue(item.Attributes, n.path) == nil
	}
	attr := getDocPathValue(item.Attributes, n.path)
	if attr == nil {
		return false
	}
	arg := n.arg.value(item)
	if arg == nil {
		return false
	}
	switch n.kind {
	case condFuncBegins:
		return beginsWithMatches(attr, arg)
	case condFuncContains:
		return containsMatches(attr, arg)
	}
	return false
}

// condAttrTypeNode is attribute_type. The type argument resolves at parse
// time — a placeholder carrying its string, or the tolerated quoted
// literal — so evaluation compares the addressed value's type descriptor
// against a fixed string.
type condAttrTypeNode struct {
	path      []docPathPart
	typeValue string
}

func (n *condAttrTypeNode) eval(item *dbstore.Item) bool {
	attr := getDocPathValue(item.Attributes, n.path)
	if attr == nil {
		return false
	}
	return getAttributeTypeName(attr) == n.typeValue
}

// ---------------------------------------------------------------------------
// Parser
// ---------------------------------------------------------------------------

type condParser struct {
	tokens []condToken
	pos    int
	names  map[string]string
	values map[string]*dbstore.AttributeValue
}

func (p *condParser) peek() condToken { return p.tokens[p.pos] }

func (p *condParser) next() condToken {
	tok := p.tokens[p.pos]
	if tok.kind != condTokenEOF {
		p.pos++
	}
	return tok
}

// keyword reports whether the upcoming token is the given keyword,
// case-insensitively, without consuming it.
func (p *condParser) keyword(word string) bool {
	tok := p.peek()
	return tok.kind == condTokenPath && strings.EqualFold(tok.text, word)
}

// expectKeyword consumes the given keyword or fails.
func (p *condParser) expectKeyword(word string) error {
	if !p.keyword(word) {
		return p.syntaxErrorf("expected %s", strings.ToUpper(word))
	}
	p.next()
	return nil
}

func (p *condParser) expectKind(kind condTokenKind, what string) (condToken, error) {
	tok := p.peek()
	if tok.kind != kind {
		return tok, p.syntaxErrorf("expected %s", what)
	}
	return p.next(), nil
}

func (p *condParser) syntaxErrorf(format string, args ...interface{}) error {
	tok := p.peek()
	shown := tok.text
	if tok.kind == condTokenEOF {
		shown = "<end of expression>"
	}
	return conditionValidationError("Invalid ConditionExpression: Syntax error; token: %s, near: %s: "+fmt.Sprintf(format, args...),
		tokenKindName(tok.kind), shown)
}

func tokenKindName(kind condTokenKind) string {
	switch kind {
	case condTokenPath:
		return "identifier"
	case condTokenPlaceholder:
		return "placeholder"
	case condTokenString:
		return "literal"
	case condTokenOperator:
		return "operator"
	case condTokenLParen:
		return "'('"
	case condTokenRParen:
		return "')'"
	case condTokenComma:
		return "','"
	case condTokenEOF:
		return "end"
	}
	return "token"
}

// parse reads one complete condition expression.
func (p *condParser) parse() (condNode, error) {
	node, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if p.peek().kind != condTokenEOF {
		return nil, p.syntaxErrorf("unexpected trailing input")
	}
	return node, nil
}

func (p *condParser) parseOr() (condNode, error) {
	first, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	if !p.keyword("OR") {
		return first, nil
	}
	operands := []condNode{first}
	for p.keyword("OR") {
		p.next()
		next, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		operands = append(operands, next)
	}
	return &condOrNode{operands: operands}, nil
}

func (p *condParser) parseAnd() (condNode, error) {
	first, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	if !p.keyword("AND") {
		return first, nil
	}
	operands := []condNode{first}
	for p.keyword("AND") {
		p.next()
		next, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		operands = append(operands, next)
	}
	return &condAndNode{operands: operands}, nil
}

// parseNot binds NOT to the immediately following condition: "NOT a AND b"
// is (NOT a) AND b, matching the precedence the grammar's grouping gives.
func (p *condParser) parseNot() (condNode, error) {
	if p.keyword("NOT") {
		p.next()
		inner, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &condNotNode{operand: inner}, nil
	}
	return p.parsePrimary()
}

// condPathFunctionNames are the functions whose whole argument list is a
// single document path; size is not among them — it is an operand form.
func condPathFunctionName(text string) (condFuncKind, bool) {
	switch text {
	case "attribute_exists":
		return condFuncExists, true
	case "attribute_not_exists":
		return condFuncNotExists, true
	}
	return 0, false
}

func (p *condParser) parsePrimary() (condNode, error) {
	tok := p.peek()

	if tok.kind == condTokenLParen {
		p.next()
		inner, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expectKind(condTokenRParen, "')'"); err != nil {
			return nil, err
		}
		return inner, nil
	}

	if tok.kind == condTokenPath {
		switch {
		case tok.text == "attribute_type":
			return p.parseAttributeType()
		case tok.text == "begins_with":
			return p.parseTwoArgFunction(condFuncBegins)
		case tok.text == "contains":
			return p.parseTwoArgFunction(condFuncContains)
		}
		if kind, ok := condPathFunctionName(tok.text); ok {
			p.next()
			if _, err := p.expectKind(condTokenLParen, "'('"); err != nil {
				return nil, err
			}
			path, err := p.parseFunctionPath()
			if err != nil {
				return nil, err
			}
			if _, err := p.expectKind(condTokenRParen, "')'"); err != nil {
				return nil, err
			}
			return &condFuncNode{kind: kind, path: path}, nil
		}
	}

	// Operand-first form: comparison, BETWEEN or IN.
	operand, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	tok = p.peek()
	switch {
	case tok.kind == condTokenOperator:
		op := p.next().text
		rhs, err := p.parseOperand()
		if err != nil {
			return nil, err
		}
		return &condCompareNode{op: op, lhs: operand, rhs: rhs}, nil
	case p.keyword("BETWEEN"):
		p.next()
		low, err := p.parseOperand()
		if err != nil {
			return nil, err
		}
		if err := p.expectKeyword("AND"); err != nil {
			return nil, err
		}
		high, err := p.parseOperand()
		if err != nil {
			return nil, err
		}
		return &condBetweenNode{operand: operand, low: low, high: high}, nil
	case p.keyword("IN"):
		p.next()
		if _, err := p.expectKind(condTokenLParen, "'('"); err != nil {
			return nil, err
		}
		var candidates []condOperand
		for {
			candidate, err := p.parseOperand()
			if err != nil {
				return nil, err
			}
			candidates = append(candidates, candidate)
			if p.peek().kind == condTokenComma {
				p.next()
				continue
			}
			break
		}
		if _, err := p.expectKind(condTokenRParen, "')'"); err != nil {
			return nil, err
		}
		return &condInNode{operand: operand, candidates: candidates}, nil
	}
	return nil, p.syntaxErrorf("expected a comparator, BETWEEN or IN after the operand")
}

// parseTwoArgFunction parses begins_with/contains: a path first argument
// and a full operand second argument. The path and the operand must be
// distinct — an operand addressing the same value as the path is the
// rejected same-operand form.
func (p *condParser) parseTwoArgFunction(kind condFuncKind) (condNode, error) {
	p.next()
	if _, err := p.expectKind(condTokenLParen, "'('"); err != nil {
		return nil, err
	}
	path, err := p.parseFunctionPath()
	if err != nil {
		return nil, err
	}
	if _, err := p.expectKind(condTokenComma, "','"); err != nil {
		return nil, err
	}
	arg, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	if pathArg, ok := arg.(*operandPath); ok && docPathPartsEqual(path, pathArg.parts) {
		name := "contains"
		if kind == condFuncBegins {
			name = "begins_with"
		}
		return nil, conditionValidationError("The path and the operand of %s must be distinct", name)
	}
	if _, err := p.expectKind(condTokenRParen, "')'"); err != nil {
		return nil, err
	}
	return &condFuncNode{kind: kind, path: path, arg: arg}, nil
}

// parseAttributeType parses attribute_type(path, type): the type argument
// is a placeholder carrying the type string, or the quoted literal form
// the platform has always tolerated.
func (p *condParser) parseAttributeType() (condNode, error) {
	p.next()
	if _, err := p.expectKind(condTokenLParen, "'('"); err != nil {
		return nil, err
	}
	path, err := p.parseFunctionPath()
	if err != nil {
		return nil, err
	}
	if _, err := p.expectKind(condTokenComma, "','"); err != nil {
		return nil, err
	}
	var typeValue string
	switch tok := p.next(); tok.kind {
	case condTokenString:
		typeValue = strings.Trim(tok.text, "'")
	case condTokenPlaceholder:
		v, ok := p.values[tok.text]
		if !ok {
			return nil, conditionValidationError(
				"An expression attribute value used in the expression is not defined; attribute value: %s", tok.text)
		}
		// The type operand is a String naming a type; any other defined
		// value is a validation error, never an always-false condition.
		if v == nil || v.S == nil {
			return nil, conditionValidationError(
				"The attribute_type operand must be a String; attribute value: %s", tok.text)
		}
		typeValue = *v.S
	default:
		return nil, p.syntaxErrorf("attribute_type requires a type operand")
	}
	if _, err := p.expectKind(condTokenRParen, "')'"); err != nil {
		return nil, err
	}
	return &condAttrTypeNode{path: path, typeValue: typeValue}, nil
}

// parseFunctionPath reads the path argument of a condition function and
// resolves its alias segments.
func (p *condParser) parseFunctionPath() ([]docPathPart, error) {
	tok := p.peek()
	if tok.kind != condTokenPath {
		return nil, p.syntaxErrorf("expected a document path")
	}
	if tok.text == "size" && p.tokens[p.pos+1].kind == condTokenLParen {
		return nil, p.syntaxErrorf("size is an operand, not the path argument of a function")
	}
	p.next()
	return p.resolvePath(tok.text)
}

// parseOperand reads one operand: a document path, size(path), or a value
// placeholder. A string literal is not an operand of this grammar.
func (p *condParser) parseOperand() (condOperand, error) {
	tok := p.peek()
	switch tok.kind {
	case condTokenPath:
		if tok.text == "size" && p.tokens[p.pos+1].kind == condTokenLParen {
			p.next()
			p.next()
			path, err := p.parseFunctionPath()
			if err != nil {
				return nil, err
			}
			if _, err := p.expectKind(condTokenRParen, "')'"); err != nil {
				return nil, err
			}
			return &operandSize{parts: path}, nil
		}
		p.next()
		parts, err := p.resolvePath(tok.text)
		if err != nil {
			return nil, err
		}
		return &operandPath{parts: parts}, nil
	case condTokenPlaceholder:
		p.next()
		v, ok := p.values[tok.text]
		if !ok {
			return nil, conditionValidationError(
				"An expression attribute value used in the expression is not defined; attribute value: %s", tok.text)
		}
		return &operandValue{v: v}, nil
	}
	return nil, p.syntaxErrorf("expected an operand")
}

// resolvePath parses a document-path token and resolves every alias
// segment through the names map. An alias absent from the map is a
// validation error: the names page requires every alias used in the
// expression to be defined.
func (p *condParser) resolvePath(path string) ([]docPathPart, error) {
	parts, err := resolveDocPathParts(path, p.names)
	if err != nil {
		return nil, conditionValidationError("Invalid ConditionExpression: invalid document path %q", path)
	}
	if len(parts) == 0 {
		return nil, conditionValidationError("Invalid ConditionExpression: empty document path")
	}
	return parts, nil
}

// ---------------------------------------------------------------------------
// Compilation entry
// ---------------------------------------------------------------------------

// compiledCondition is a condition expression parsed and validated once;
// evaluation against any item cannot fail.
type compiledCondition struct{ root condNode }

// matches reports whether the condition holds for the item.
func (c *compiledCondition) matches(item *dbstore.Item) bool {
	if c == nil || c.root == nil {
		return true
	}
	return c.root.eval(item)
}

// compileConditionExpression parses and validates one condition
// expression: grammar, operator set, and substitution resolution are all
// checked here, so the returned tree evaluates without error.
func compileConditionExpression(expr string, names map[string]string, values map[string]*dbstore.AttributeValue) (*compiledCondition, error) {
	tokens, err := lexCondition(expr)
	if err != nil {
		return nil, err
	}
	parser := &condParser{tokens: tokens, names: names, values: values}
	root, err := parser.parse()
	if err != nil {
		return nil, err
	}
	return &compiledCondition{root: root}, nil
}

// ---------------------------------------------------------------------------
// Function value matching (shared by the condition tree)
// ---------------------------------------------------------------------------

// beginsWithMatches applies begins_with's operand typing: a string or
// binary attribute with an operand of the same kind, prefix-tested.
func beginsWithMatches(attr, arg *dbstore.AttributeValue) bool {
	if attr.S != nil && arg.S != nil {
		return strings.HasPrefix(*attr.S, *arg.S)
	}
	if attr.B != nil && arg.B != nil {
		return bytes.HasPrefix(attr.B, arg.B)
	}
	return false
}

// containsMatches applies contains's operand typing: a substring of a
// string, a subsequence of a binary, a member of a string/number/binary
// set, or a deep-equal element of a list.
func containsMatches(attr, arg *dbstore.AttributeValue) bool {
	if attr.S != nil && arg.S != nil {
		return strings.Contains(*attr.S, *arg.S)
	}
	if attr.B != nil && arg.B != nil {
		return bytes.Contains(attr.B, arg.B)
	}
	if setContainsMember(attr, arg) {
		return true
	}
	if attr.L != nil {
		return listContainsValue(attr.L, arg)
	}
	return false
}

// setContainsMember answers whether a scalar operand of the same kind is a
// member of a set attribute (SS/NS/BS); number members compare through the
// normalised number string.
func setContainsMember(attr, arg *dbstore.AttributeValue) bool {
	switch {
	case attr.SS != nil && arg.S != nil:
		for _, s := range attr.SS {
			if s == *arg.S {
				return true
			}
		}
	case attr.NS != nil && arg.N != nil:
		for _, n := range attr.NS {
			if normalizeNumberString(n) == normalizeNumberString(*arg.N) {
				return true
			}
		}
	case attr.BS != nil && arg.B != nil:
		for _, b := range attr.BS {
			if bytesEqual(b, arg.B) {
				return true
			}
		}
	}
	return false
}

// bytesEqual compares two binary operands bytewise.
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
