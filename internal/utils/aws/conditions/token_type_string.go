// Package conditions provides AWS-specific utility functions for vorpalstacks.
package conditions

// TokenType represents the type of a token in DynamoDB expressions.
type TokenType int

// Token types for parsing DynamoDB condition expressions.
const (
	// TokenEOF is the end of file token.
	TokenEOF TokenType = iota
	// TokenIdentifier is an identifier token.
	TokenIdentifier
	// TokenOperator is an operator token.
	TokenOperator
	// TokenFunction is a function token.
	TokenFunction
	// TokenStringLiteral is a string literal token.
	TokenStringLiteral
	// TokenNumberLiteral is a number literal token.
	TokenNumberLiteral
	// TokenBooleanLiteral is a boolean literal token.
	TokenBooleanLiteral
	// TokenNullLiteral is a null literal token.
	TokenNullLiteral
	// TokenLeftParen is a left parenthesis token.
	TokenLeftParen
	// TokenRightParen is a right parenthesis token.
	TokenRightParen
	// TokenComma is a comma token.
	TokenComma
	// TokenAnd is an AND token.
	TokenAnd
	// TokenOr is an OR token.
	TokenOr
	// TokenNot is a NOT token.
	TokenNot
	// TokenBetween is a BETWEEN token.
	TokenBetween
	// TokenIn is an IN token.
	TokenIn
	// TokenPathSeparator is a path separator token.
	TokenPathSeparator
	// TokenColon is a colon token.
	TokenColon
)

// TokenTypeString returns a string representation of a token type.
//
// Parameters:
//   - t: The token type
//
// Returns:
//   - string: The string representation of the token type
func TokenTypeString(t TokenType) string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenOperator:
		return "OPERATOR"
	case TokenFunction:
		return "FUNCTION"
	case TokenStringLiteral:
		return "STRING_LITERAL"
	case TokenNumberLiteral:
		return "NUMBER_LITERAL"
	case TokenBooleanLiteral:
		return "BOOLEAN_LITERAL"
	case TokenNullLiteral:
		return "NULL_LITERAL"
	case TokenLeftParen:
		return "LEFT_PAREN"
	case TokenRightParen:
		return "RIGHT_PAREN"
	case TokenComma:
		return "COMMA"
	case TokenAnd:
		return "AND"
	case TokenOr:
		return "OR"
	case TokenNot:
		return "NOT"
	case TokenBetween:
		return "BETWEEN"
	case TokenIn:
		return "IN"
	case TokenPathSeparator:
		return "PATH_SEPARATOR"
	case TokenColon:
		return "COLON"
	default:
		return "UNKNOWN"
	}
}
