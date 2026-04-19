// Package cypherparser implements an openCypher query parser and executor.
// It targets the subset required by the Neptune Data API (ExecuteOpenCypherQuery).
//
// Architecture:
//
//	cypher string → [Lexer] → []Token → [Parser] → AST → [Executor] → CypherResult
//
// The lexer and parser are hand-written recursive descent. The executor operates
// against the graphengine.GraphStore interface, decoupled from Pebble storage.
//
// Design references:
//   - goraphdb/cypher_lexer.go: inspired the hand-written batch tokenizer approach
//     and the token kind enum layout. Token types, scanString, scanNumber, and
//     scanIdentOrKeyword patterns were used as a starting point but significantly
//     extended (multi-word keywords, backtick identifiers, line comments, arithmetic
//     operators, != syntax, =~ regex match, aggregation function tokens).
//   - openCypher specification (opencypher.org): defines the full grammar.

package cypherparser

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenKind identifies the type of a lexer token.
// Inspired by goraphdb TokenKind enum, extended with WITH, UNWIND, DISTINCT,
// IN, CONTAINS, STARTS WITH, ENDS WITH, CASE/WHEN/THEN/ELSE/END, IS, EXISTS,
// aggregation functions (COUNT/SUM/AVG/MIN/MAX/COLLECT), arithmetic operators,
// != and =~ comparison operators, and backtick identifier support.
type TokenKind int

const (
	tokEOF TokenKind = iota
	tokError

	// Literals
	tokIdent  // unquoted or backtick-quoted identifier
	tokString // 'hello' or "hello"
	tokInt    // 42
	tokFloat  // 3.14

	// Query keywords (case-insensitive)
	tokMatch
	tokWhere
	tokReturn
	tokOrder
	tokBy
	tokLimit
	tokAnd
	tokOr
	tokNot
	tokAs
	tokAsc
	tokDesc
	tokTrue
	tokFalse
	tokNull
	tokType
	tokExplain
	tokProfile
	tokOptional
	tokCreate
	tokMerge
	tokSet
	tokDelete
	tokRemove
	tokOn
	tokSkip
	tokDetach
	tokShow
	tokIndex
	tokConstraint
	tokIndexes
	tokConstraints
	tokCall
	tokYield

	// Pipeline keywords (our addition beyond goraphdb)
	tokWith
	tokUnwind
	tokDistinct

	// Predicate keywords (our addition beyond goraphdb)
	tokIn
	tokContains
	tokStartsWith // "STARTS WITH" (multi-word keyword)
	tokEndsWith   // "ENDS WITH" (multi-word keyword)

	// CASE expression keywords (our addition beyond goraphdb)
	tokCase
	tokWhen
	tokThen
	tokElse
	tokEnd
	tokIs

	// Built-in function names treated as keywords for expression parsing
	tokExists
	tokCount
	tokSum
	tokAvg
	tokMin
	tokMax
	tokCollect
	tokProperties
	tokLabels

	// Comparison operators
	tokEq         // =
	tokNeq        // <> or !=
	tokLt         // <
	tokGt         // >
	tokLte        // <=
	tokGte        // >=
	tokRegexMatch // =~ (our addition beyond goraphdb)

	// Punctuation and operators
	tokLParen
	tokRParen
	tokLBracket
	tokRBracket
	tokLBrace
	tokRBrace
	tokColon
	tokDot
	tokDotDot
	tokComma
	tokStar
	tokPlus    // arithmetic (our addition)
	tokPlusEq  // += (SET n += {map})
	tokSlash   // division (our addition)
	tokPercent // modulo (our addition)
	tokCaret   // exponentiation (our addition)
	tokTilde   // string concatenation (our addition)
	tokDash    // arrow prefix in patterns: ->
	tokArrow   // ->
	tokLArrow  // <- (incoming edge)
	tokParam   // $paramName
	tokSemicolon
	tokPipe
)

// Token is a single lexer token with its kind, raw text, and byte offset.
// Based on goraphdb Token struct.
type Token struct {
	Kind TokenKind
	Text string
	Pos  int
}

// String returns a human-readable representation of the token.
func (t Token) String() string {
	return fmt.Sprintf("%s(%q)@%d", tokenKindName(t.Kind), t.Text, t.Pos)
}

// tokenKindName returns a human-readable name for a token kind.
func tokenKindName(k TokenKind) string {
	switch k {
	case tokEOF:
		return "EOF"
	case tokError:
		return "ERROR"
	case tokIdent:
		return "IDENT"
	case tokString:
		return "STRING"
	case tokInt:
		return "INT"
	case tokFloat:
		return "FLOAT"
	case tokMatch:
		return "MATCH"
	case tokWhere:
		return "WHERE"
	case tokReturn:
		return "RETURN"
	case tokOrder:
		return "ORDER"
	case tokBy:
		return "BY"
	case tokLimit:
		return "LIMIT"
	case tokAnd:
		return "AND"
	case tokOr:
		return "OR"
	case tokNot:
		return "NOT"
	case tokAs:
		return "AS"
	case tokAsc:
		return "ASC"
	case tokDesc:
		return "DESC"
	case tokTrue:
		return "TRUE"
	case tokFalse:
		return "FALSE"
	case tokNull:
		return "NULL"
	case tokType:
		return "TYPE"
	case tokExplain:
		return "EXPLAIN"
	case tokProfile:
		return "PROFILE"
	case tokOptional:
		return "OPTIONAL"
	case tokCreate:
		return "CREATE"
	case tokMerge:
		return "MERGE"
	case tokSet:
		return "SET"
	case tokDelete:
		return "DELETE"
	case tokRemove:
		return "REMOVE"
	case tokOn:
		return "ON"
	case tokSkip:
		return "SKIP"
	case tokDetach:
		return "DETACH"
	case tokShow:
		return "SHOW"
	case tokIndex:
		return "INDEX"
	case tokConstraint:
		return "CONSTRAINT"
	case tokIndexes:
		return "INDEXES"
	case tokConstraints:
		return "CONSTRAINTS"
	case tokCall:
		return "CALL"
	case tokYield:
		return "YIELD"
	case tokWith:
		return "WITH"
	case tokUnwind:
		return "UNWIND"
	case tokDistinct:
		return "DISTINCT"
	case tokIn:
		return "IN"
	case tokContains:
		return "CONTAINS"
	case tokStartsWith:
		return "STARTS_WITH"
	case tokEndsWith:
		return "ENDS_WITH"
	case tokCase:
		return "CASE"
	case tokWhen:
		return "WHEN"
	case tokThen:
		return "THEN"
	case tokElse:
		return "ELSE"
	case tokEnd:
		return "END"
	case tokIs:
		return "IS"
	case tokExists:
		return "EXISTS"
	case tokCount:
		return "COUNT"
	case tokSum:
		return "SUM"
	case tokAvg:
		return "AVG"
	case tokMin:
		return "MIN"
	case tokMax:
		return "MAX"
	case tokCollect:
		return "COLLECT"
	case tokProperties:
		return "PROPERTIES"
	case tokLabels:
		return "LABELS"
	case tokEq:
		return "="
	case tokNeq:
		return "<>"
	case tokLt:
		return "<"
	case tokGt:
		return ">"
	case tokLte:
		return "<="
	case tokGte:
		return ">="
	case tokRegexMatch:
		return "=~"
	case tokLParen:
		return "("
	case tokRParen:
		return ")"
	case tokLBracket:
		return "["
	case tokRBracket:
		return "]"
	case tokLBrace:
		return "{"
	case tokRBrace:
		return "}"
	case tokColon:
		return ":"
	case tokDot:
		return "."
	case tokDotDot:
		return ".."
	case tokComma:
		return ","
	case tokStar:
		return "*"
	case tokPlus:
		return "+"
	case tokPlusEq:
		return "+="
	case tokSlash:
		return "/"
	case tokPercent:
		return "%"
	case tokCaret:
		return "^"
	case tokTilde:
		return "~"
	case tokDash:
		return "-"
	case tokArrow:
		return "->"
	case tokLArrow:
		return "<-"
	case tokParam:
		return "PARAM"
	case tokSemicolon:
		return ";"
	case tokPipe:
		return "|"
	default:
		return "???"
	}
}

// keywords maps uppercase keyword text to token kind.
// Base set from goraphdb keywords map, extended with pipeline/predicate/CASE tokens.
var keywords = map[string]TokenKind{
	"MATCH":       tokMatch,
	"WHERE":       tokWhere,
	"RETURN":      tokReturn,
	"ORDER":       tokOrder,
	"BY":          tokBy,
	"LIMIT":       tokLimit,
	"AND":         tokAnd,
	"OR":          tokOr,
	"NOT":         tokNot,
	"AS":          tokAs,
	"ASC":         tokAsc,
	"DESC":        tokDesc,
	"TRUE":        tokTrue,
	"FALSE":       tokFalse,
	"NULL":        tokNull,
	"TYPE":        tokType,
	"EXPLAIN":     tokExplain,
	"PROFILE":     tokProfile,
	"OPTIONAL":    tokOptional,
	"CREATE":      tokCreate,
	"MERGE":       tokMerge,
	"SET":         tokSet,
	"DELETE":      tokDelete,
	"REMOVE":      tokRemove,
	"ON":          tokOn,
	"SKIP":        tokSkip,
	"DETACH":      tokDetach,
	"SHOW":        tokShow,
	"INDEX":       tokIndex,
	"CONSTRAINT":  tokConstraint,
	"INDEXES":     tokIndexes,
	"CONSTRAINTS": tokConstraints,
	"CALL":        tokCall,
	"YIELD":       tokYield,
	"WITH":        tokWith,
	"UNWIND":      tokUnwind,
	"DISTINCT":    tokDistinct,
	"IN":          tokIn,
	"CONTAINS":    tokContains,
	"STARTS":      tokStartsWith,
	"ENDS":        tokEndsWith,
	"CASE":        tokCase,
	"WHEN":        tokWhen,
	"THEN":        tokThen,
	"ELSE":        tokElse,
	"END":         tokEnd,
	"IS":          tokIs,
	"EXISTS":      tokExists,
	"COUNT":       tokCount,
	"SUM":         tokSum,
	"AVG":         tokAvg,
	"MIN":         tokMin,
	"MAX":         tokMax,
	"COLLECT":     tokCollect,
	"PROPERTIES":  tokProperties,
	"LABELS":      tokLabels,
}

// lexer holds the state for tokenising a Cypher string.
// Based on goraphdb lexer struct with the same scan-once batch approach.
type lexer struct {
	input  string
	pos    int
	tokens []Token
}

// Tokenize converts a Cypher query string into a slice of tokens.
func Tokenize(input string) ([]Token, error) {
	l := &lexer{input: input}
	if err := l.scan(); err != nil {
		return nil, err
	}
	return l.tokens, nil
}

// scan is the main tokenisation loop.
// Structure inspired by goraphdb scan(), but extended with:
//   - != syntax (in addition to <>)
//   - =~ regex match operator
//   - line comment support (//)
//   - backtick-quoted identifiers
//   - arithmetic operators (+, %, ^, ~)
//   - semicolon and pipe tokens
func (l *lexer) scan() error {
	for l.pos < len(l.input) {
		l.skipWhitespace()
		if l.pos >= len(l.input) {
			break
		}

		ch := l.input[l.pos]

		switch {
		case ch == '(':
			l.emit(tokLParen, "(")
		case ch == ')':
			l.emit(tokRParen, ")")
		case ch == '[':
			l.emit(tokLBracket, "[")
		case ch == ']':
			l.emit(tokRBracket, "]")
		case ch == '{':
			l.emit(tokLBrace, "{")
		case ch == '}':
			l.emit(tokRBrace, "}")
		case ch == ':':
			l.emit(tokColon, ":")
		case ch == ',':
			l.emit(tokComma, ",")
		case ch == '*':
			l.emit(tokStar, "*")
		case ch == '+':
			if l.peek(1) == '=' {
				l.emitN(tokPlusEq, "+=", 2)
			} else {
				l.emit(tokPlus, "+")
			}
		case ch == '%':
			l.emit(tokPercent, "%")
		case ch == '^':
			l.emit(tokCaret, "^")
		case ch == ';':
			l.emit(tokSemicolon, ";")
		case ch == '|':
			l.emit(tokPipe, "|")
		case ch == '~':
			l.emit(tokTilde, "~")

		case ch == '.':
			// Disambiguate: ".." is the range operator, "." is property access.
			if l.peek(1) == '.' {
				l.emitN(tokDotDot, "..", 2)
			} else {
				l.emit(tokDot, ".")
			}

		case ch == '-':
			// "->" is the outgoing arrow; "-" is a dash in patterns or arithmetic minus.
			if l.peek(1) == '>' {
				l.emitN(tokArrow, "->", 2)
			} else {
				l.emit(tokDash, "-")
			}

		case ch == '<':
			// "<-" is the incoming arrow prefix, "<=" is LTE, "<>" is NEQ.
			next := l.peek(1)
			switch {
			case next == '-':
				l.emitN(tokLArrow, "<-", 2)
			case next == '=':
				l.emitN(tokLte, "<=", 2)
			case next == '>':
				l.emitN(tokNeq, "<>", 2)
			default:
				l.emit(tokLt, "<")
			}

		case ch == '>':
			if l.peek(1) == '=' {
				l.emitN(tokGte, ">=", 2)
			} else {
				l.emit(tokGt, ">")
			}

		case ch == '=':
			// "=~" is regex match; "=" is equality.
			if l.peek(1) == '~' {
				l.emitN(tokRegexMatch, "=~", 2)
			} else {
				l.emit(tokEq, "=")
			}

		case ch == '!':
			// "!=" is the alternative inequality operator (our addition).
			if l.peek(1) == '=' {
				l.emitN(tokNeq, "!=", 2)
			} else {
				return fmt.Errorf("cypher lexer: unexpected character '!' at position %d (did you mean '!=' or '<>')", l.pos)
			}

		case ch == '/':
			if l.peek(1) == '/' {
				l.skipLineComment()
			} else if l.peek(1) == '*' {
				if err := l.skipBlockComment(); err != nil {
					return err
				}
			} else {
				l.emit(tokSlash, "/")
			}

		case ch == '\'' || ch == '"':
			if err := l.scanString(ch); err != nil {
				return err
			}

		case isDigit(ch):
			l.scanNumber()

		case ch == '$':
			l.scanParam()

		case isIdentStart(ch):
			l.scanIdentOrKeyword()

		case ch == '`':
			// Backtick-quoted identifier (our addition — not in goraphdb).
			if err := l.scanBacktickString(); err != nil {
				return err
			}

		default:
			return fmt.Errorf("cypher lexer: unexpected character %q at position %d", ch, l.pos)
		}
	}

	l.tokens = append(l.tokens, Token{Kind: tokEOF, Text: "", Pos: l.pos})
	return nil
}

func (l *lexer) emit(kind TokenKind, text string) {
	l.tokens = append(l.tokens, Token{Kind: kind, Text: text, Pos: l.pos})
	l.pos++
}

func (l *lexer) emitN(kind TokenKind, text string, n int) {
	l.tokens = append(l.tokens, Token{Kind: kind, Text: text, Pos: l.pos})
	l.pos += n
}

func (l *lexer) peek(offset int) byte {
	idx := l.pos + offset
	if idx >= len(l.input) {
		return 0
	}
	return l.input[idx]
}

func (l *lexer) skipWhitespace() {
	for l.pos < len(l.input) && isWhitespace(l.input[l.pos]) {
		l.pos++
	}
}

// skipLineComment skips from "//" to end of line. Our addition.
func (l *lexer) skipLineComment() {
	l.pos += 2
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.pos++
	}
}

// skipBlockComment skips from "/*" to "*/". Nested block comments are not supported.
func (l *lexer) skipBlockComment() error {
	l.pos += 2
	for l.pos+1 < len(l.input) {
		if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
			l.pos += 2
			return nil
		}
		l.pos++
	}
	return fmt.Errorf("cypher lexer: unterminated block comment")
}

// scanString scans a single- or double-quoted string literal.
// Based on goraphdb scanString, extended with \r and \0 escape sequences.
func (l *lexer) scanString(quote byte) error {
	start := l.pos
	l.pos++
	var b strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\\' && l.pos+1 < len(l.input) {
			l.pos++
			switch l.input[l.pos] {
			case '\\':
				b.WriteByte('\\')
			case '\'':
				b.WriteByte('\'')
			case '"':
				b.WriteByte('"')
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case '0':
				b.WriteByte(0)
			default:
				b.WriteByte('\\')
				b.WriteByte(l.input[l.pos])
			}
			l.pos++
			continue
		}
		if ch == quote {
			l.pos++
			l.tokens = append(l.tokens, Token{Kind: tokString, Text: b.String(), Pos: start})
			return nil
		}
		b.WriteByte(ch)
		l.pos++
	}
	return fmt.Errorf("cypher lexer: unterminated string starting at position %d", start)
}

// scanBacktickString scans a backtick-quoted identifier (our addition).
// Cypher allows `backtick quoted` identifiers for names containing special chars.
func (l *lexer) scanBacktickString() error {
	start := l.pos
	l.pos++
	var b strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '`' {
			l.pos++
			l.tokens = append(l.tokens, Token{Kind: tokIdent, Text: b.String(), Pos: start})
			return nil
		}
		if ch == '\\' && l.pos+1 < len(l.input) {
			l.pos++
			b.WriteByte(l.input[l.pos])
			l.pos++
			continue
		}
		b.WriteByte(ch)
		l.pos++
	}
	return fmt.Errorf("cypher lexer: unterminated backtick string starting at position %d", start)
}

// scanNumber scans an integer or float literal.
// Based on goraphdb scanNumber with the same ".." disambiguation logic.
func (l *lexer) scanNumber() {
	start := l.pos
	isFloat := false
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		// A single dot followed by a digit is a decimal point.
		// ".." is the range operator, so we don't consume it.
		if ch == '.' && !isFloat && l.peek(1) != '.' {
			if l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1]) {
				isFloat = true
				l.pos++
				continue
			}
			break
		}
		if !isDigit(ch) {
			break
		}
		l.pos++
	}
	// Scientific notation: e/E followed by optional +/- and digits
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		nextPos := l.pos + 1
		if nextPos < len(l.input) && (l.input[nextPos] == '+' || l.input[nextPos] == '-') {
			nextPos++
		}
		if nextPos < len(l.input) && isDigit(l.input[nextPos]) {
			isFloat = true
			l.pos = nextPos
			for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
				l.pos++
			}
		}
	}
	text := l.input[start:l.pos]
	if isFloat {
		l.tokens = append(l.tokens, Token{Kind: tokFloat, Text: text, Pos: start})
	} else {
		l.tokens = append(l.tokens, Token{Kind: tokInt, Text: text, Pos: start})
	}
}

// scanIdentOrKeyword scans an identifier and promotes it to a keyword if it matches.
// Based on goraphdb scanIdentOrKeyword, extended with multi-word keyword handling
// for "STARTS WITH" and "ENDS WITH" (our addition).
func (l *lexer) scanIdentOrKeyword() {
	start := l.pos
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.pos++
	}
	text := l.input[start:l.pos]
	upper := strings.ToUpper(text)

	if kind, ok := keywords[upper]; ok {
		// Multi-word keywords: STARTS WITH, ENDS WITH
		// These need lookahead to distinguish from plain identifiers "starts"/"ends".
		if kind == tokStartsWith || kind == tokEndsWith {
			l.scanMultiWordKeyword(start, text, upper)
			return
		}
		l.tokens = append(l.tokens, Token{Kind: kind, Text: text, Pos: start})
		return
	}

	l.tokens = append(l.tokens, Token{Kind: tokIdent, Text: text, Pos: start})
}

// scanMultiWordKeyword handles compound keywords like "STARTS WITH" and "ENDS WITH".
// Our addition — not present in goraphdb.
// If the next word completes the compound keyword, a single token is emitted.
// Otherwise the first word is emitted as a plain identifier and the position
// is restored so the second word can be scanned independently.
func (l *lexer) scanMultiWordKeyword(start int, text, upper string) {
	savedPos := l.pos
	l.skipWhitespace()
	next := l.pos
	if next >= len(l.input) {
		l.tokens = append(l.tokens, Token{Kind: tokIdent, Text: text, Pos: start})
		return
	}
	wordStart := next
	for next < len(l.input) && isIdentPart(l.input[next]) {
		next++
	}
	secondWord := l.input[wordStart:next]
	upperSecond := strings.ToUpper(secondWord)

	if upper == "STARTS" && upperSecond == "WITH" {
		l.pos = next
		l.tokens = append(l.tokens, Token{Kind: tokStartsWith, Text: text + " " + secondWord, Pos: start})
		return
	}
	if upper == "ENDS" && upperSecond == "WITH" {
		l.pos = next
		l.tokens = append(l.tokens, Token{Kind: tokEndsWith, Text: text + " " + secondWord, Pos: start})
		return
	}

	// Not a compound keyword — emit "starts"/"ends" as identifier,
	// restore pos so the second word (e.g. "WITH") is scanned normally.
	l.pos = savedPos
	l.tokens = append(l.tokens, Token{Kind: tokIdent, Text: text, Pos: start})
}

// scanParam scans a parameter reference: $name.
// Based on goraphdb scanParam. Stores just the name without the '$' prefix.
func (l *lexer) scanParam() {
	start := l.pos
	l.pos++
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.pos++
	}
	text := l.input[start+1 : l.pos]
	if len(text) == 0 {
		return
	}
	l.tokens = append(l.tokens, Token{Kind: tokParam, Text: text, Pos: start})
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isIdentStart(ch byte) bool {
	return ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || unicode.IsLetter(rune(ch))
}

func isIdentPart(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}
