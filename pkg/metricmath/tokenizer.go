package metricmath

import (
	"fmt"
	"unicode"
)

// Tokenize converts a metric math expression string into a sequence of
// tokens. Identifiers may contain letters, digits, and underscores but
// must start with a letter. Numbers may contain a decimal point and an
// optional leading sign is handled by the parser as a unary operator.
func Tokenize(input string) ([]Token, error) {
	var tokens []Token
	i := 0
	for i < len(input) {
		ch := input[i]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
			i++
		case ch == '+':
			tokens = append(tokens, Token{TokPlus, "+", i})
			i++
		case ch == '-':
			tokens = append(tokens, Token{TokMinus, "-", i})
			i++
		case ch == '*':
			tokens = append(tokens, Token{TokStar, "*", i})
			i++
		case ch == '/':
			tokens = append(tokens, Token{TokSlash, "/", i})
			i++
		case ch == '^':
			tokens = append(tokens, Token{TokCaret, "^", i})
			i++
		case ch == '(':
			tokens = append(tokens, Token{TokLParen, "(", i})
			i++
		case ch == ')':
			tokens = append(tokens, Token{TokRParen, ")", i})
			i++
		case ch == ',':
			tokens = append(tokens, Token{TokComma, ",", i})
			i++
		case unicode.IsLetter(rune(ch)) || ch == '_':
			start := i
			for i < len(input) && (unicode.IsLetter(rune(input[i])) || unicode.IsDigit(rune(input[i])) || input[i] == '_') {
				i++
			}
			tokens = append(tokens, Token{TokIdent, input[start:i], start})
		case unicode.IsDigit(rune(ch)) || ch == '.':
			start := i
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
				i++
			}
			tokens = append(tokens, Token{TokNumber, input[start:i], start})
		default:
			return nil, fmt.Errorf("metricmath: unexpected character %q at position %d", ch, i)
		}
	}
	tokens = append(tokens, Token{TokEOF, "", len(input)})
	return tokens, nil
}
