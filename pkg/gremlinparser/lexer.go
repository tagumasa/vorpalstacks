package gremlinparser

import (
	"fmt"
)

type TokenKind int

const (
	tokEOF TokenKind = iota
	tokDot
	tokComma
	tokLParen
	tokRParen
	tokLBracket
	tokRBracket
	tokLBrace
	tokRBrace
	tokColon
	tokSemicolon
	tokUnderscoreUnderscore
	tokString
	tokInteger
	tokFloat
	tokIdentifier
	tokTrue
	tokFalse
	tokNull
)

type Token struct {
	Kind TokenKind
	Text string
	Pos  int
}

type Lexer struct {
	input  []byte
	pos    int
	tokens []Token
}

func Tokenize(input string) ([]Token, error) {
	l := &Lexer{input: []byte(input), pos: 0}
	for l.pos < len(l.input) {
		l.skipWhitespaceAndComments()
		if l.pos >= len(l.input) {
			break
		}
		if err := l.nextToken(); err != nil {
			return nil, err
		}
	}
	l.tokens = append(l.tokens, Token{Kind: tokEOF, Pos: l.pos})
	return l.tokens, nil
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.pos++
		} else if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		} else if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			commentStart := l.pos
			l.pos += 2
			found := false
			for l.pos < len(l.input) {
				if l.input[l.pos] == '*' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
					l.pos += 2
					found = true
					break
				}
				l.pos++
			}
			if !found {
				panic(fmt.Sprintf("gremlin: unterminated block comment starting at position %d", commentStart))
			}
		} else {
			break
		}
	}
}

func (l *Lexer) nextToken() error {
	ch := l.input[l.pos]
	start := l.pos

	switch ch {
	case '.':
		if l.pos+1 < len(l.input) && l.input[l.pos+1] == '.' {
			l.pos += 2
			l.tokens = append(l.tokens, Token{Kind: tokIdentifier, Text: "..", Pos: start})
			return nil
		}
		l.pos++
		l.tokens = append(l.tokens, Token{Kind: tokDot, Text: ".", Pos: start})
		return nil
	case ',':
		l.pos++
		l.tokens = append(l.tokens, Token{Kind: tokComma, Text: ",", Pos: start})
		return nil
	case '(':
		l.pos++
		l.tokens = append(l.tokens, Token{Kind: tokLParen, Text: "(", Pos: start})
		return nil
	case ')':
		l.pos++
		l.tokens = append(l.tokens, Token{Kind: tokRParen, Text: ")", Pos: start})
		return nil
	case '[':
		l.pos++
		l.tokens = append(l.tokens, Token{Kind: tokLBracket, Text: "[", Pos: start})
		return nil
	case ']':
		l.pos++
		l.tokens = append(l.tokens, Token{Kind: tokRBracket, Text: "]", Pos: start})
		return nil
	case '{':
		l.pos++
		l.tokens = append(l.tokens, Token{Kind: tokLBrace, Text: "{", Pos: start})
		return nil
	case '}':
		l.pos++
		l.tokens = append(l.tokens, Token{Kind: tokRBrace, Text: "}", Pos: start})
		return nil
	case ':':
		l.pos++
		l.tokens = append(l.tokens, Token{Kind: tokColon, Text: ":", Pos: start})
		return nil
	case ';':
		l.pos++
		l.tokens = append(l.tokens, Token{Kind: tokSemicolon, Text: ";", Pos: start})
		return nil
	case '\'':
		return l.scanString('\'')
	case '"':
		return l.scanString('"')
	}

	if ch == '_' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '_' {
		l.pos += 2
		l.tokens = append(l.tokens, Token{Kind: tokUnderscoreUnderscore, Pos: start})
		return nil
	}

	if ch >= '0' && ch <= '9' || (ch == '-' && l.pos+1 < len(l.input) && (l.input[l.pos+1] >= '0' && l.input[l.pos+1] <= '9') && !isValueToken(l.prevTokenKind())) {
		return l.scanNumber()
	}

	if isIdentStart(ch) {
		return l.scanIdentOrKeyword()
	}

	return fmt.Errorf("gremlin: unexpected character %q at position %d", ch, l.pos)
}

func (l *Lexer) scanString(quote byte) error {
	start := l.pos
	l.pos++
	var buf []byte
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == quote {
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == quote {
				buf = append(buf, quote)
				l.pos += 2
			} else {
				l.pos++
				l.tokens = append(l.tokens, Token{Kind: tokString, Text: string(buf), Pos: start})
				return nil
			}
		} else if ch == '\\' {
			l.pos++
			if l.pos >= len(l.input) {
				return fmt.Errorf("gremlin: unterminated string escape at position %d", start)
			}
			esc := l.input[l.pos]
			switch esc {
			case 'n':
				buf = append(buf, '\n')
			case 't':
				buf = append(buf, '\t')
			case 'r':
				buf = append(buf, '\r')
			case '\\':
				buf = append(buf, '\\')
			case '\'':
				buf = append(buf, '\'')
			case '"':
				buf = append(buf, '"')
			case 'u':
				if l.pos+4 >= len(l.input) {
					return fmt.Errorf("gremlin: incomplete unicode escape at position %d", l.pos)
				}
				hex := l.input[l.pos+1 : l.pos+5]
				var r rune
				for _, c := range hex {
					r = r<<4 | rune(hexDigit(c))
				}
				if r >= 0xD800 && r <= 0xDFFF {
					return fmt.Errorf("gremlin: invalid unicode surrogate at position %d", l.pos)
				}
				buf = append(buf, []byte(string(r))...)
				l.pos += 4
			default:
				buf = append(buf, esc)
			}
			l.pos++
		} else {
			buf = append(buf, ch)
			l.pos++
		}
	}
	return fmt.Errorf("gremlin: unterminated string starting at position %d", start)
}

func (l *Lexer) scanNumber() error {
	start := l.pos
	isFloat := false

	if l.input[l.pos] == '-' {
		l.pos++
	}

	if l.pos < len(l.input) && l.input[l.pos] == '0' && l.pos+1 < len(l.input) {
		next := l.input[l.pos+1]
		if next == 'x' || next == 'X' {
			l.pos += 2
			for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
				l.pos++
			}
			l.tokens = append(l.tokens, Token{Kind: tokInteger, Text: string(l.input[start:l.pos]), Pos: start})
			return nil
		}
	}

	for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
		l.pos++
	}

	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		if l.pos+1 < len(l.input) && l.input[l.pos+1] == '.' {
			l.tokens = append(l.tokens, Token{Kind: tokInteger, Text: string(l.input[start:l.pos]), Pos: start})
			return nil
		}
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	}

	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		isFloat = true
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	}

	if l.pos < len(l.input) {
		suffix := l.input[l.pos]
		if suffix == 'f' || suffix == 'F' || suffix == 'd' || suffix == 'D' || suffix == 'm' || suffix == 'M' {
			isFloat = true
			l.pos++
		} else if suffix == 'i' || suffix == 'I' || suffix == 's' || suffix == 'S' || suffix == 'l' || suffix == 'L' || suffix == 'n' || suffix == 'N' {
			l.pos++
		}
	}

	kind := tokInteger
	if isFloat {
		kind = tokFloat
	}
	l.tokens = append(l.tokens, Token{Kind: kind, Text: string(l.input[start:l.pos]), Pos: start})
	return nil
}

func (l *Lexer) scanIdentOrKeyword() error {
	start := l.pos
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.pos++
	}
	text := string(l.input[start:l.pos])
	kind := tokIdentifier
	switch text {
	case "true":
		kind = tokTrue
	case "false":
		kind = tokFalse
	case "null":
		kind = tokNull
	}
	l.tokens = append(l.tokens, Token{Kind: kind, Text: text, Pos: start})
	return nil
}

func isIdentStart(ch byte) bool {
	return ch == '_' || ch == '$' || ch == '~' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isIdentPart(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func hexDigit(ch byte) byte {
	if ch >= '0' && ch <= '9' {
		return ch - '0'
	}
	if ch >= 'a' && ch <= 'f' {
		return ch - 'a' + 10
	}
	return ch - 'A' + 10
}

func (l *Lexer) prevTokenKind() TokenKind {
	if len(l.tokens) > 0 {
		return l.tokens[len(l.tokens)-1].Kind
	}
	return tokEOF
}

func isValueToken(kind TokenKind) bool {
	switch kind {
	case tokInteger, tokFloat, tokString, tokTrue, tokFalse, tokNull, tokRParen, tokRBracket, tokIdentifier:
		return true
	}
	return false
}
