package cloudwatchlogs

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// tokenKind classifies lexical tokens of a CloudWatch Logs Insights query.
type tokenKind int

const (
	tokIdent tokenKind = iota
	tokNumber
	tokString
	tokBacktickIdent
	tokOp
	tokPipe
	tokLParen
	tokRParen
	tokLBracket
	tokRBracket
	tokLBrace
	tokRBrace
	tokComma
	tokColon
	tokSemicolon
	tokRegex
)

// token is a single lexical token with its character offsets in the original
// query string, so that compile errors can report exact positions.
type token struct {
	kind  tokenKind
	text  string // decoded text for strings and backtick identifiers
	raw   string // raw source text
	start int
	end   int
}

// lexer converts a query string into a token stream. Comments introduced by
// # run to the end of the line and are skipped. Token offsets are counted
// in characters (runes), matching the AWS error location fields; the lexer
// scans bytes and converts positions through a precomputed index.
type lexer struct {
	src     string
	pos     int
	toks    []token
	runeIdx []int
}

func lexQuery(src string) ([]token, error) {
	l := &lexer{src: src, runeIdx: computeRuneIndex(src)}
	for {
		l.skipSpaceAndComments()
		if l.pos >= len(l.src) {
			return l.toks, nil
		}
		start := l.pos
		c := l.src[l.pos]
		switch {
		case c == '"' || c == '\'':
			if err := l.lexString(); err != nil {
				return nil, err
			}
		case c == '`':
			l.lexBacktick()
		case isDigit(c) || (c == '.' && l.pos+1 < len(l.src) && isDigit(l.src[l.pos+1])):
			l.lexNumber()
		case isIdentStart(c):
			l.lexIdent()
		default:
			if !l.lexOperator() {
				r, _ := utf8.DecodeRuneInString(l.src[start:])
				return nil, newQueryCompileError(fmt.Sprintf("Unexpected character '%c'", r),
					l.runeOffset(start), l.runeOffset(start)+1)
			}
		}
	}
}

// computeRuneIndex maps each byte position of src to its rune position so
// that token offsets can be reported in characters, as the AWS
// startCharOffset/endCharOffset fields expect.
func computeRuneIndex(src string) []int {
	idx := make([]int, len(src)+1)
	r := 0
	for i := 0; i < len(src); {
		idx[i] = r
		_, size := utf8.DecodeRuneInString(src[i:])
		i += size
		r++
	}
	idx[len(src)] = r
	return idx
}

// runeOffset converts a byte position into a character position.
func (l *lexer) runeOffset(bytePos int) int {
	if bytePos < 0 {
		return 0
	}
	if bytePos >= len(l.runeIdx) {
		return l.runeIdx[len(l.runeIdx)-1]
	}
	return l.runeIdx[bytePos]
}

func (l *lexer) skipSpaceAndComments() {
	for l.pos < len(l.src) {
		c := l.src[l.pos]
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			l.pos++
			continue
		}
		if c == '#' {
			for l.pos < len(l.src) && l.src[l.pos] != '\n' {
				l.pos++
			}
			continue
		}
		return
	}
}

func (l *lexer) push(kind tokenKind, raw, text string, start int) {
	l.toks = append(l.toks, token{kind: kind, text: text, raw: raw,
		start: l.runeOffset(start), end: l.runeOffset(l.pos)})
}

func (l *lexer) lexString() error {
	start := l.pos
	quote := l.src[l.pos]
	l.pos++ // opening quote
	var b strings.Builder
	for l.pos < len(l.src) {
		c := l.src[l.pos]
		if c == '\\' && l.pos+1 < len(l.src) {
			l.pos++
			b.WriteByte(l.src[l.pos])
			l.pos++
			continue
		}
		if c == quote {
			l.pos++
			l.push(tokString, l.src[start:l.pos], b.String(), start)
			return nil
		}
		b.WriteByte(c)
		l.pos++
	}
	return newQueryCompileError("Unterminated string literal", l.runeOffset(start), l.runeOffset(l.pos))
}

func (l *lexer) lexBacktick() {
	start := l.pos
	l.pos++ // opening backtick
	var b strings.Builder
	for l.pos < len(l.src) {
		c := l.src[l.pos]
		if c == '`' {
			l.pos++
			break
		}
		b.WriteByte(c)
		l.pos++
	}
	l.push(tokBacktickIdent, l.src[start:l.pos], b.String(), start)
}

func (l *lexer) lexNumber() {
	start := l.pos
	for l.pos < len(l.src) && (isDigit(l.src[l.pos]) || l.src[l.pos] == '.') {
		l.pos++
	}
	// Scientific notation.
	if l.pos < len(l.src) && (l.src[l.pos] == 'e' || l.src[l.pos] == 'E') {
		save := l.pos
		l.pos++
		if l.pos < len(l.src) && (l.src[l.pos] == '+' || l.src[l.pos] == '-') {
			l.pos++
		}
		if l.pos < len(l.src) && isDigit(l.src[l.pos]) {
			for l.pos < len(l.src) && isDigit(l.src[l.pos]) {
				l.pos++
			}
		} else {
			l.pos = save
		}
	}
	l.push(tokNumber, l.src[start:l.pos], l.src[start:l.pos], start)
}

func (l *lexer) lexIdent() {
	start := l.pos
	for l.pos < len(l.src) && isIdentPart(l.src[l.pos]) {
		l.pos++
	}
	l.push(tokIdent, l.src[start:l.pos], l.src[start:l.pos], start)
}

// exprKeywords are the reserved words after which a forward slash starts a
// regular expression literal rather than division.
var exprKeywords = map[string]bool{
	"like": true, "and": true, "or": true, "not": true, "in": true,
	"as": true, "by": true, "asc": true, "desc": true, "where": true,
}

// lexOperator consumes one operator or punctuation character. Multi-character
// operators are matched longest-first.
func (l *lexer) lexOperator() bool {
	rest := l.src[l.pos:]
	// The equality operator is "=" alone — the query syntax documents
	// no "==" alias, so a doubled equals lexes as two operators and the
	// parser reports the syntax error AWS reports.
	multi := []string{"!=", ">=", "<=", "=~"}
	for _, op := range multi {
		if strings.HasPrefix(rest, op) {
			start := l.pos
			l.pos += len(op)
			l.push(tokOp, op, op, start)
			return true
		}
	}
	start := l.pos
	c := l.src[l.pos]
	if c == '/' {
		// A forward slash is division when the previous token can end an
		// expression (value or closing bracket) and is not a reserved word;
		// otherwise it opens a regular expression literal. A dotted
		// identifier is a field reference and ends a value like any other
		// identifier — treating it otherwise breaks division on nested
		// fields (bytes.received / 2), core query syntax.
		if n := len(l.toks); n > 0 {
			prev := l.toks[n-1]
			endsValue := prev.kind == tokNumber || prev.kind == tokString ||
				prev.kind == tokRParen || prev.kind == tokRBracket ||
				prev.kind == tokBacktickIdent ||
				(prev.kind == tokIdent && !exprKeywords[strings.ToLower(prev.text)])
			if !endsValue || l.slashOpensRegexLiteral() {
				return l.lexRegex()
			}
		} else {
			return l.lexRegex()
		}
	}
	l.pos++
	switch c {
	case '|':
		l.push(tokPipe, "|", "|", start)
	case '(':
		l.push(tokLParen, "(", "(", start)
	case ')':
		l.push(tokRParen, ")", ")", start)
	case '[':
		l.push(tokLBracket, "[", "[", start)
	case ']':
		l.push(tokRBracket, "]", "]", start)
	case '{':
		l.push(tokLBrace, "{", "{", start)
	case '}':
		l.push(tokRBrace, "}", "}", start)
	case ',':
		l.push(tokComma, ",", ",", start)
	case ':':
		l.push(tokColon, ":", ":", start)
	case ';':
		l.push(tokSemicolon, ";", ";", start)
	case '=', '<', '>', '+', '-', '*', '%', '^', '!', '.', '/':
		l.push(tokOp, string(c), string(c), start)
	default:
		l.pos = start
		return false
	}
	return true
}

// lexRegex scans a /.../ regular expression literal. Backslash-escaped
// slashes stay inside the pattern.
func (l *lexer) lexRegex() bool {
	start := l.pos
	l.pos++ // opening slash
	var b strings.Builder
	for l.pos < len(l.src) {
		c := l.src[l.pos]
		if c == '\\' && l.pos+1 < len(l.src) {
			b.WriteByte(c)
			b.WriteByte(l.src[l.pos+1])
			l.pos += 2
			continue
		}
		if c == '/' {
			l.pos++
			l.push(tokRegex, l.src[start:l.pos], b.String(), start)
			return true
		}
		b.WriteByte(c)
		l.pos++
	}
	l.pos = start
	return false
}

// slashOpensRegexLiteral reports whether the forward slash at l.pos opens
// a regular expression literal even though the preceding token could end a
// value, where a slash otherwise means division. This is the parse
// command's regex mode, whose documented form places the literal directly
// after the source field — parse {{fieldName}} /{{regex}}/ — and the
// closing delimiter is followed by the end of the query, a pipe to the
// next command, a comment, the as clause of the numbered-group form, or
// the multi keyword. A division's right-hand side is a value, so none of
// those continuations follows a division slash and the look-ahead
// separates the two readings. The residual ambiguity — a division whose
// scan happens to close on a later parse regex whose body begins with the
// word as or multi — resolves as a regex literal and fails to compile;
// it cannot silently mis-evaluate.
func (l *lexer) slashOpensRegexLiteral() bool {
	i := l.pos + 1 // past the opening slash
	for i < len(l.src) {
		c := l.src[i]
		if c == '\\' && i+1 < len(l.src) {
			// An escaped character keeps an escaped slash inside the body.
			i += 2
			continue
		}
		if c != '/' {
			i++
			continue
		}
		rest := strings.TrimLeft(l.src[i+1:], " \t\r\n")
		if rest == "" || rest[0] == '|' || rest[0] == '#' {
			return true
		}
		word := identWordLen(rest)
		return strings.EqualFold(rest[:word], "as") || strings.EqualFold(rest[:word], "multi")
	}
	return false
}

// identWordLen returns the length of the leading identifier-shaped run of
// s, so a look-ahead can compare whole keywords (as, multi) against a
// word boundary instead of a prefix.
func identWordLen(s string) int {
	i := 0
	for i < len(s) && isIdentPart(s[i]) {
		i++
	}
	return i
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isIdentStart(c byte) bool {
	return c == '_' || c == '@' || c == '$' ||
		(c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentPart(c byte) bool {
	return isIdentStart(c) || isDigit(c) || c == '.'
}

// newQueryCompileError builds the MalformedQueryException carrying the
// QueryCompileError structure with character offsets, matching the AWS error
// shape returned when a query fails to compile.
func newQueryCompileError(message string, start, end int) error {
	err := NewLogsError("MalformedQueryException", message, 400)
	err.RawFields = map[string]interface{}{
		"queryCompileError": map[string]interface{}{
			"message": message,
			"location": map[string]interface{}{
				"startCharOffset": start,
				"endCharOffset":   end,
			},
		},
	}
	return err
}
