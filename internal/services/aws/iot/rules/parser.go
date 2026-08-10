package rules

import (
	"fmt"
	"strings"
	"unicode"
)

type SelectExpr struct {
	Fields    []SelectField
	FromTopic string
	Where     Expr
}

type SelectField struct {
	Expression Expr
	Alias      string
}

type Expr interface {
	exprNode()
}

type (
	StringLiteral struct{ Value string }
	NumberLiteral struct{ Value float64 }
	BoolLiteral   struct{ Value bool }
	NullLiteral   struct{}
	Identifier    struct{ Name string }
	JsonPath      struct{ Path []string }
	BinaryExpr    struct {
		Op          string
		Left, Right Expr
	}
	UnaryExpr struct {
		Op      string
		Operand Expr
	}
	FunctionCall struct {
		Name string
		Args []Expr
	}
	InExpr struct {
		Expr   Expr
		Not    bool
		Values []Expr
	}
	StarExpr struct{}
)

func (*StringLiteral) exprNode() {}
func (*NumberLiteral) exprNode() {}
func (*BoolLiteral) exprNode()   {}
func (*NullLiteral) exprNode()   {}
func (*Identifier) exprNode()    {}
func (*JsonPath) exprNode()      {}
func (*BinaryExpr) exprNode()    {}
func (*UnaryExpr) exprNode()     {}
func (*FunctionCall) exprNode()  {}
func (*InExpr) exprNode()        {}
func (*StarExpr) exprNode()      {}

type TokenKind int

const (
	TokenEOF TokenKind = iota
	TokenSelect
	TokenFrom
	TokenWhere
	TokenAnd
	TokenOr
	TokenNot
	TokenIdent
	TokenString
	TokenNumber
	TokenComma
	TokenDot
	TokenStar
	TokenLParen
	TokenRParen
	TokenLBracket
	TokenRBracket
	TokenEQ
	TokenNEQ
	TokenLT
	TokenGT
	TokenLTE
	TokenGTE
	TokenLike
	TokenIn
	TokenIs
	TokenNull
	TokenPlus
	TokenMinus
	TokenSlash
	TokenPercent
)

type Token struct {
	Kind  TokenKind
	Value string
	Pos   int
}

func (t Token) String() string {
	switch t.Kind {
	case TokenEOF:
		return "EOF"
	case TokenIdent, TokenString, TokenNumber:
		return fmt.Sprintf("%s(%q)", tokenName(t.Kind), t.Value)
	default:
		return tokenName(t.Kind)
	}
}

func tokenName(k TokenKind) string {
	switch k {
	case TokenEOF:
		return "EOF"
	case TokenSelect:
		return "SELECT"
	case TokenFrom:
		return "FROM"
	case TokenWhere:
		return "WHERE"
	case TokenAnd:
		return "AND"
	case TokenOr:
		return "OR"
	case TokenNot:
		return "NOT"
	case TokenIdent:
		return "IDENT"
	case TokenString:
		return "STRING"
	case TokenNumber:
		return "NUMBER"
	case TokenComma:
		return "COMMA"
	case TokenDot:
		return "DOT"
	case TokenStar:
		return "STAR"
	case TokenLParen:
		return "LPAREN"
	case TokenRParen:
		return "RPAREN"
	case TokenLBracket:
		return "LBRACKET"
	case TokenRBracket:
		return "RBRACKET"
	case TokenEQ:
		return "="
	case TokenNEQ:
		return "<>"
	case TokenLT:
		return "<"
	case TokenGT:
		return ">"
	case TokenLTE:
		return "<="
	case TokenGTE:
		return ">="
	case TokenLike:
		return "LIKE"
	case TokenIn:
		return "IN"
	case TokenIs:
		return "IS"
	case TokenNull:
		return "NULL"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenSlash:
		return "/"
	case TokenPercent:
		return "%"
	default:
		return "UNKNOWN"
	}
}

type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	if l.pos >= len(l.input) {
		return Token{Kind: TokenEOF, Pos: l.pos}
	}

	ch := l.input[l.pos]

	if ch == '\'' {
		return l.readString()
	}
	if ch == '"' {
		return l.readDoubleString()
	}
	if unicode.IsDigit(ch) || (ch == '.' && l.pos+1 < len(l.input) && unicode.IsDigit(l.input[l.pos+1])) {
		return l.readNumber()
	}
	if ch == '*' {
		l.pos++
		return Token{Kind: TokenStar, Value: "*", Pos: l.pos - 1}
	}
	if ch == '.' {
		l.pos++
		return Token{Kind: TokenDot, Value: ".", Pos: l.pos - 1}
	}
	if ch == ',' {
		l.pos++
		return Token{Kind: TokenComma, Value: ",", Pos: l.pos - 1}
	}
	if ch == '(' {
		l.pos++
		return Token{Kind: TokenLParen, Value: "(", Pos: l.pos - 1}
	}
	if ch == ')' {
		l.pos++
		return Token{Kind: TokenRParen, Value: ")", Pos: l.pos - 1}
	}
	if ch == '[' {
		l.pos++
		return Token{Kind: TokenLBracket, Value: "[", Pos: l.pos - 1}
	}
	if ch == ']' {
		l.pos++
		return Token{Kind: TokenRBracket, Value: "]", Pos: l.pos - 1}
	}
	if ch == '+' {
		l.pos++
		return Token{Kind: TokenPlus, Value: "+", Pos: l.pos - 1}
	}
	if ch == '-' {
		l.pos++
		return Token{Kind: TokenMinus, Value: "-", Pos: l.pos - 1}
	}
	if ch == '/' {
		l.pos++
		return Token{Kind: TokenSlash, Value: "/", Pos: l.pos - 1}
	}
	if ch == '%' {
		l.pos++
		return Token{Kind: TokenPercent, Value: "%", Pos: l.pos - 1}
	}
	if ch == '=' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
		l.pos += 2
		return Token{Kind: TokenEQ, Value: "==", Pos: l.pos - 2}
	}
	if ch == '!' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
		l.pos += 2
		return Token{Kind: TokenNEQ, Value: "!=", Pos: l.pos - 2}
	}
	if ch == '<' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
		l.pos += 2
		return Token{Kind: TokenLTE, Value: "<=", Pos: l.pos - 2}
	}
	if ch == '<' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '>' {
		l.pos += 2
		return Token{Kind: TokenNEQ, Value: "<>", Pos: l.pos - 2}
	}
	if ch == '<' {
		l.pos++
		return Token{Kind: TokenLT, Value: "<", Pos: l.pos - 1}
	}
	if ch == '>' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
		l.pos += 2
		return Token{Kind: TokenGTE, Value: ">=", Pos: l.pos - 2}
	}
	if ch == '>' {
		l.pos++
		return Token{Kind: TokenGT, Value: ">", Pos: l.pos - 1}
	}
	if ch == '=' {
		l.pos++
		return Token{Kind: TokenEQ, Value: "=", Pos: l.pos - 1}
	}
	if isIdentStart(ch) {
		return l.readIdent()
	}

	l.pos++
	return Token{Kind: TokenEOF, Pos: l.pos}
}

func isIdentStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isIdentPart(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '-'
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
		l.pos++
	}
}

func (l *Lexer) readIdent() Token {
	start := l.pos
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.pos++
	}
	val := string(l.input[start:l.pos])
	upper := strings.ToUpper(val)
	kind := TokenIdent
	switch upper {
	case "SELECT":
		kind = TokenSelect
	case "FROM":
		kind = TokenFrom
	case "WHERE":
		kind = TokenWhere
	case "AND":
		kind = TokenAnd
	case "OR":
		kind = TokenOr
	case "NOT":
		kind = TokenNot
	case "LIKE":
		kind = TokenLike
	case "IN":
		kind = TokenIn
	case "IS":
		kind = TokenIs
	case "NULL":
		kind = TokenNull
	}
	return Token{Kind: kind, Value: val, Pos: start}
}

func (l *Lexer) readString() Token {
	l.pos++
	start := l.pos
	for l.pos < len(l.input) {
		if l.input[l.pos] == '\'' {
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
				l.pos += 2
				continue
			}
			break
		}
		l.pos++
	}
	val := strings.ReplaceAll(string(l.input[start:l.pos]), "''", "'")
	l.pos++
	return Token{Kind: TokenString, Value: val, Pos: start}
}

func (l *Lexer) readDoubleString() Token {
	l.pos++
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		if l.input[l.pos] == '\\' && l.pos+1 < len(l.input) {
			l.pos++
		}
		l.pos++
	}
	val := string(l.input[start:l.pos])
	if l.pos < len(l.input) {
		l.pos++
	}
	return Token{Kind: TokenString, Value: val, Pos: start}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	hasDot := false
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '.' && !hasDot {
			hasDot = true
			l.pos++
			continue
		}
		if !unicode.IsDigit(ch) {
			break
		}
		l.pos++
	}
	return Token{Kind: TokenNumber, Value: string(l.input[start:l.pos]), Pos: start}
}

type Parser struct {
	lexer *Lexer
	peek  Token
}

func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.peek = p.lexer.NextToken()
	return p
}

func (p *Parser) next() Token {
	tok := p.peek
	p.peek = p.lexer.NextToken()
	return tok
}

func (p *Parser) expect(kind TokenKind) (Token, error) {
	tok := p.next()
	if tok.Kind != kind {
		return tok, fmt.Errorf("parser: expected %s, got %s at pos %d", tokenName(kind), tok, tok.Pos)
	}
	return tok, nil
}

func (p *Parser) Parse() (result *SelectExpr, err error) {
	result, err = p.parseSelect()
	return
}

func (p *Parser) parseSelect() (*SelectExpr, error) {
	if _, err := p.expect(TokenSelect); err != nil {
		return nil, err
	}
	fields, err := p.parseSelectFields()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(TokenFrom); err != nil {
		return nil, err
	}
	topicTok, err := p.expect(TokenString)
	if err != nil {
		return nil, err
	}
	fromTopic := topicTok.Value

	var where Expr
	if p.peek.Kind == TokenWhere {
		p.next()
		where, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	return &SelectExpr{Fields: fields, FromTopic: fromTopic, Where: where}, nil
}

func (p *Parser) parseSelectFields() ([]SelectField, error) {
	var fields []SelectField
	if p.peek.Kind == TokenStar {
		p.next()
		return append(fields, SelectField{Expression: &StarExpr{}}), nil
	}
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		alias := ""
		if p.peek.Kind == TokenIdent && strings.ToUpper(p.peek.Value) == "AS" {
			p.next()
			aliasTok, err := p.expect(TokenIdent)
			if err != nil {
				return nil, err
			}
			alias = aliasTok.Value
		}
		fields = append(fields, SelectField{Expression: expr, Alias: alias})
		if p.peek.Kind == TokenComma {
			p.next()
		} else {
			break
		}
	}
	return fields, nil
}

func (p *Parser) parseExpr() (Expr, error) {
	return p.parseOr()
}

func (p *Parser) parseOr() (Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.peek.Kind == TokenOr {
		p.next()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Op: "OR", Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseAnd() (Expr, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.peek.Kind == TokenAnd {
		p.next()
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Op: "AND", Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseNot() (Expr, error) {
	if p.peek.Kind == TokenNot {
		p.next()
		operand, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{Op: "NOT", Operand: operand}, nil
	}
	return p.parseComparison()
}

func (p *Parser) parseComparison() (Expr, error) {
	left, err := p.parseAddSub()
	if err != nil {
		return nil, err
	}
	switch p.peek.Kind {
	case TokenEQ, TokenNEQ, TokenLT, TokenGT, TokenLTE, TokenGTE, TokenLike:
		op := strings.ToUpper(p.next().Value)
		right, err := p.parseAddSub()
		if err != nil {
			return nil, err
		}
		return &BinaryExpr{Op: op, Left: left, Right: right}, nil
	case TokenIs:
		p.next()
		if p.peek.Kind == TokenNot {
			p.next()
			if p.peek.Kind == TokenNull {
				p.next()
				return &BinaryExpr{Op: "IS NOT NULL", Left: left, Right: &NullLiteral{}}, nil
			}
			return nil, fmt.Errorf("parser: expected NULL after IS NOT")
		}
		if p.peek.Kind == TokenNull {
			p.next()
			return &BinaryExpr{Op: "IS NULL", Left: left, Right: &NullLiteral{}}, nil
		}
		return nil, fmt.Errorf("parser: expected NULL or NOT after IS")
	case TokenIn:
		return p.parseInExpr(left, false)
	case TokenNot:
		p.next()
		if p.peek.Kind == TokenIn {
			return p.parseInExpr(left, true)
		}
		return nil, fmt.Errorf("parser: expected IN after NOT")
	}
	return left, nil
}

func (p *Parser) parseInExpr(left Expr, negated bool) (Expr, error) {
	p.next()

	// AWS IoT SQL supports two IN syntaxes:
	//   1. IN (expr, expr, ...) — comma-separated list
	//   2. IN expr              — single expression evaluating to an array
	if p.peek.Kind != TokenLParen {
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &InExpr{Expr: left, Not: negated, Values: []Expr{val}}, nil
	}

	p.next()
	var values []Expr
	for p.peek.Kind != TokenRParen && p.peek.Kind != TokenEOF {
		if len(values) > 0 {
			if p.peek.Kind != TokenComma {
				return nil, fmt.Errorf("parser: expected ',' or ')' in IN list")
			}
			p.next()
		}
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	if p.peek.Kind != TokenRParen {
		return nil, fmt.Errorf("parser: expected ')' to close IN list")
	}
	p.next()
	return &InExpr{Expr: left, Not: negated, Values: values}, nil
}

func (p *Parser) parseAddSub() (Expr, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return nil, err
	}
	for p.peek.Kind == TokenPlus || p.peek.Kind == TokenMinus {
		op := p.next().Value
		right, err := p.parseMulDiv()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseMulDiv() (Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.peek.Kind == TokenStar || p.peek.Kind == TokenSlash || p.peek.Kind == TokenPercent {
		op := p.next().Value
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseUnary() (Expr, error) {
	if p.peek.Kind == TokenMinus {
		p.next()
		operand, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{Op: "-", Operand: operand}, nil
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (Expr, error) {
	tok := p.peek

	switch tok.Kind {
	case TokenString:
		p.next()
		return &StringLiteral{Value: tok.Value}, nil
	case TokenNumber:
		p.next()
		val := 0.0
		fmt.Sscanf(tok.Value, "%f", &val)
		return &NumberLiteral{Value: val}, nil
	case TokenLParen:
		p.next()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(TokenRParen); err != nil {
			return nil, err
		}
		return expr, nil
	case TokenLBracket:
		p.next()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(TokenRBracket); err != nil {
			return nil, err
		}
		return expr, nil
	case TokenIdent:
		p.next()
		name := tok.Value

		if p.peek.Kind == TokenLParen {
			return p.parseFunctionCall(name)
		}
		if p.peek.Kind == TokenDot {
			return p.parseJsonPath(name)
		}
		return &Identifier{Name: name}, nil
	case TokenNull:
		p.next()
		return &NullLiteral{}, nil
	default:
		return nil, fmt.Errorf("parser: unexpected token %s at pos %d", tok, tok.Pos)
	}
}

func (p *Parser) parseFunctionCall(name string) (Expr, error) {
	if _, err := p.expect(TokenLParen); err != nil {
		return nil, err
	}

	if strings.ToUpper(name) == "CAST" {
		return p.parseCastCall()
	}

	var args []Expr
	for p.peek.Kind != TokenRParen {
		if len(args) > 0 {
			if _, err := p.expect(TokenComma); err != nil {
				return nil, err
			}
		}
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	if _, err := p.expect(TokenRParen); err != nil {
		return nil, err
	}
	return &FunctionCall{Name: name, Args: args}, nil
}

func (p *Parser) parseCastCall() (Expr, error) {
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.peek.Kind == TokenComma {
		p.next()
	} else if p.peek.Kind == TokenIdent && strings.ToUpper(p.peek.Value) == "AS" {
		p.next()
	} else {
		return nil, fmt.Errorf("expected AS or ',' in cast(), got %v", p.peek)
	}
	typExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(TokenRParen); err != nil {
		return nil, err
	}
	return &FunctionCall{Name: "CAST", Args: []Expr{val, typExpr}}, nil
}

func (p *Parser) parseJsonPath(prefix string) (Expr, error) {
	parts := []string{prefix}
	for p.peek.Kind == TokenDot {
		p.next()
		partTok, err := p.expect(TokenIdent)
		if err != nil {
			return nil, err
		}
		parts = append(parts, partTok.Value)
	}
	return &JsonPath{Path: parts}, nil
}
