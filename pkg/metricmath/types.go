// Package metricmath implements the CloudWatch metric math expression
// engine. It tokenises, parses, and evaluates arithmetic expressions
// over time-series data points, supporting binary operators, unary
// operators, function calls, and variable references.
package metricmath

import "time"

// DataPoint represents a single timestamped value in a time series.
type DataPoint struct {
	Timestamp time.Time
	Value     float64
}

// Token represents a lexical token in a metric math expression.
type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

// TokenType identifies the kind of a token.
type TokenType int

const (
	TokNumber TokenType = iota
	TokIdent
	TokPlus
	TokMinus
	TokStar
	TokSlash
	TokCaret
	TokLParen
	TokRParen
	TokComma
	TokEOF
)
