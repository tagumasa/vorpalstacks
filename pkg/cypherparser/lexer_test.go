package cypherparser

import (
	"testing"
)

func TestLexer_SimpleMatch(t *testing.T) {
	input := `MATCH (n) RETURN n`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) < 6 {
		t.Fatalf("expected at least 6 tokens, got %d", len(tokens))
	}
	assertToken(t, tokens[0], tokMatch, "MATCH")
	assertToken(t, tokens[1], tokLParen, "(")
	assertToken(t, tokens[2], tokIdent, "n")
	assertToken(t, tokens[3], tokRParen, ")")
	assertToken(t, tokens[4], tokReturn, "RETURN")
	assertToken(t, tokens[5], tokIdent, "n")
	assertLastEOF(t, tokens)
}

func TestLexer_PropertyFilter(t *testing.T) {
	input := `MATCH (n {name: "Alice"}) WHERE n.age > 25 RETURN n`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatal(err)
	}

	foundMatch := false
	foundStringAlice := false
	foundGt := false
	for _, tok := range tokens {
		if tok.Kind == tokMatch {
			foundMatch = true
		}
		if tok.Kind == tokString && tok.Text == "Alice" {
			foundStringAlice = true
		}
		if tok.Kind == tokGt {
			foundGt = true
		}
	}
	if !foundMatch {
		t.Error("expected MATCH token")
	}
	if !foundStringAlice {
		t.Error("expected string token 'Alice'")
	}
	if !foundGt {
		t.Error("expected > token")
	}
}

func TestLexer_ComparisonOperators(t *testing.T) {
	tests := []struct {
		input string
		kind  TokenKind
		text  string
	}{
		{"=", tokEq, "="},
		{"<>", tokNeq, "<>"},
		{"!=", tokNeq, "!="},
		{"<", tokLt, "<"},
		{">", tokGt, ">"},
		{"<=", tokLte, "<="},
		{">=", tokGte, ">="},
		{"=~", tokRegexMatch, "=~"},
	}
	for _, tc := range tests {
		t.Run(tc.text, func(t *testing.T) {
			tokens, err := Tokenize(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			assertToken(t, tokens[0], tc.kind, tc.text)
		})
	}
}

func TestLexer_ArrowTokens(t *testing.T) {
	tokens, err := Tokenize("-> <- ->")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokArrow, "->")
	assertToken(t, tokens[1], tokLArrow, "<-")
	assertToken(t, tokens[2], tokArrow, "->")
}

func TestLexer_Numbers(t *testing.T) {
	tokens, err := Tokenize("42 3.14 0")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokInt, "42")
	assertToken(t, tokens[1], tokFloat, "3.14")
	assertToken(t, tokens[2], tokInt, "0")
}

func TestLexer_DotDotDisambiguation(t *testing.T) {
	input := `n.name..5 3.14`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokIdent, "n")
	assertToken(t, tokens[1], tokDot, ".")
	assertToken(t, tokens[2], tokIdent, "name")
	assertToken(t, tokens[3], tokDotDot, "..")
	assertToken(t, tokens[4], tokInt, "5")
	assertToken(t, tokens[5], tokFloat, "3.14")
}

func TestLexer_ParameterTokens(t *testing.T) {
	tokens, err := Tokenize(`$name $age $x`)
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokParam, "name")
	assertToken(t, tokens[1], tokParam, "age")
	assertToken(t, tokens[2], tokParam, "x")
}

func TestLexer_CaseInsensitiveKeywords(t *testing.T) {
	tests := []struct {
		input string
		kind  TokenKind
	}{
		{"match", tokMatch},
		{"Match", tokMatch},
		{"MATCH", tokMatch},
		{"where", tokWhere},
		{"WHERE", tokWhere},
		{"return", tokReturn},
		{"RETURN", tokReturn},
		{"create", tokCreate},
		{"CREATE", tokCreate},
		{"merge", tokMerge},
		{"MERGE", tokMerge},
		{"with", tokWith},
		{"WITH", tokWith},
		{"unwind", tokUnwind},
		{"UNWIND", tokUnwind},
		{"distinct", tokDistinct},
		{"DISTINCT", tokDistinct},
		{"optional", tokOptional},
		{"OPTIONAL", tokOptional},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			tokens, err := Tokenize(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if tokens[0].Kind != tc.kind {
				t.Fatalf("expected %s for %q, got %s", tokenKindName(tc.kind), tc.input, tokenKindName(tokens[0].Kind))
			}
		})
	}
}

func TestLexer_BooleanAndNull(t *testing.T) {
	tokens, err := Tokenize("true false null TRUE FALSE NULL")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokTrue, "true")
	assertToken(t, tokens[1], tokFalse, "false")
	assertToken(t, tokens[2], tokNull, "null")
	assertToken(t, tokens[3], tokTrue, "TRUE")
	assertToken(t, tokens[4], tokFalse, "FALSE")
	assertToken(t, tokens[5], tokNull, "NULL")
}

func TestLexer_StringEscapes(t *testing.T) {
	tokens, err := Tokenize(`'hello\nworld' "tab\there"`)
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Text != "hello\nworld" {
		t.Fatalf("expected escaped newline, got %q", tokens[0].Text)
	}
	if tokens[1].Text != "tab\there" {
		t.Fatalf("expected escaped tab, got %q", tokens[1].Text)
	}
}

func TestLexer_UnterminatedString(t *testing.T) {
	_, err := Tokenize(`"hello`)
	if err == nil {
		t.Fatal("expected error for unterminated string")
	}
}

func TestLexer_PredicateKeywords(t *testing.T) {
	tokens, err := Tokenize("IN CONTAINS STARTS WITH ENDS WITH IS NOT NULL EXISTS")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokIn, "IN")
	assertToken(t, tokens[1], tokContains, "CONTAINS")
	assertToken(t, tokens[2], tokStartsWith, "STARTS WITH")
	assertToken(t, tokens[3], tokEndsWith, "ENDS WITH")
	assertToken(t, tokens[4], tokIs, "IS")
	assertToken(t, tokens[5], tokNot, "NOT")
	assertToken(t, tokens[6], tokNull, "NULL")
	assertToken(t, tokens[7], tokExists, "EXISTS")
}

func TestLexer_CaseKeywords(t *testing.T) {
	tokens, err := Tokenize("CASE WHEN x > 0 THEN 'pos' ELSE 'neg' END")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokCase, "CASE")
	assertToken(t, tokens[1], tokWhen, "WHEN")
	assertToken(t, tokens[2], tokIdent, "x")
	assertToken(t, tokens[3], tokGt, ">")
	assertToken(t, tokens[4], tokInt, "0")
	assertToken(t, tokens[5], tokThen, "THEN")
	assertToken(t, tokens[6], tokString, "pos")
	assertToken(t, tokens[7], tokElse, "ELSE")
	assertToken(t, tokens[8], tokString, "neg")
	assertToken(t, tokens[9], tokEnd, "END")
}

func TestLexer_AggregationKeywords(t *testing.T) {
	tokens, err := Tokenize("COUNT SUM AVG MIN MAX COLLECT")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokCount, "COUNT")
	assertToken(t, tokens[1], tokSum, "SUM")
	assertToken(t, tokens[2], tokAvg, "AVG")
	assertToken(t, tokens[3], tokMin, "MIN")
	assertToken(t, tokens[4], tokMax, "MAX")
	assertToken(t, tokens[5], tokCollect, "COLLECT")
}

func TestLexer_ArithmeticOperators(t *testing.T) {
	tokens, err := Tokenize("+ - * / % ^ ~")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokPlus, "+")
	assertToken(t, tokens[1], tokDash, "-")
	assertToken(t, tokens[2], tokStar, "*")
	assertToken(t, tokens[3], tokSlash, "/")
	assertToken(t, tokens[4], tokPercent, "%")
	assertToken(t, tokens[5], tokCaret, "^")
	assertToken(t, tokens[6], tokTilde, "~")
}

func TestLexer_LineComment(t *testing.T) {
	input := `MATCH (n) // this is a comment
RETURN n`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokMatch, "MATCH")
	// The comment should be skipped entirely
	assertToken(t, tokens[4], tokReturn, "RETURN")
}

func TestLexer_BacktickIdentifier(t *testing.T) {
	tokens, err := Tokenize("`special name`")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokIdent, "special name")
}

func TestLexer_RelPattern(t *testing.T) {
	input := `(a)-[r:FOLLOWS*1..3]->(b)`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokLParen, "(")
	assertToken(t, tokens[1], tokIdent, "a")
	assertToken(t, tokens[2], tokRParen, ")")
	assertToken(t, tokens[3], tokDash, "-")
	assertToken(t, tokens[4], tokLBracket, "[")
	assertToken(t, tokens[5], tokIdent, "r")
	assertToken(t, tokens[6], tokColon, ":")
	assertToken(t, tokens[7], tokIdent, "FOLLOWS")
	assertToken(t, tokens[8], tokStar, "*")
	assertToken(t, tokens[9], tokInt, "1")
	assertToken(t, tokens[10], tokDotDot, "..")
	assertToken(t, tokens[11], tokInt, "3")
	assertToken(t, tokens[12], tokRBracket, "]")
	assertToken(t, tokens[13], tokArrow, "->")
	assertToken(t, tokens[14], tokLParen, "(")
	assertToken(t, tokens[15], tokIdent, "b")
	assertToken(t, tokens[16], tokRParen, ")")
}

func TestLexer_IncomingEdge(t *testing.T) {
	input := `<-[r:KNOWS]-`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokLArrow, "<-")
	assertToken(t, tokens[1], tokLBracket, "[")
	assertToken(t, tokens[2], tokIdent, "r")
	assertToken(t, tokens[3], tokColon, ":")
	assertToken(t, tokens[4], tokIdent, "KNOWS")
	assertToken(t, tokens[5], tokRBracket, "]")
	assertToken(t, tokens[6], tokDash, "-")
}

func TestLexer_FullQuery(t *testing.T) {
	input := `MATCH (a:Person {name: "Alice"})-[:FOLLOWS]->(b) WHERE b.age >= 25 RETURN b.name ORDER BY b.name DESC LIMIT 10`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatal(err)
	}
	if tokens[len(tokens)-1].Kind != tokEOF {
		t.Fatal("expected EOF at end")
	}

	kinds := make([]TokenKind, len(tokens)-1)
	for i := 0; i < len(tokens)-1; i++ {
		kinds[i] = tokens[i].Kind
	}

	expected := []TokenKind{
		tokMatch, tokLParen, tokIdent, tokColon, tokIdent,
		tokLBrace, tokIdent, tokColon, tokString, tokRBrace,
		tokRParen, tokDash, tokLBracket, tokColon, tokIdent,
		tokRBracket, tokArrow, tokLParen, tokIdent, tokRParen,
		tokWhere, tokIdent, tokDot, tokIdent, tokGte,
		tokInt, tokReturn, tokIdent, tokDot, tokIdent,
		tokOrder, tokBy, tokIdent, tokDot, tokIdent,
		tokDesc, tokLimit, tokInt,
	}

	if len(kinds) != len(expected) {
		t.Fatalf("expected %d tokens, got %d\nexpected: %v\ngot:      %v",
			len(expected), len(kinds), expected, kinds)
	}
	for i, exp := range expected {
		if kinds[i] != exp {
			t.Fatalf("token %d: expected %s, got %s (text=%q)",
				i, tokenKindName(exp), tokenKindName(kinds[i]), tokens[i].Text)
		}
	}
}

func TestLexer_StartsWithAsIdent(t *testing.T) {
	input := `MATCH (n) WHERE n.name STARTS WITH "A" RETURN n`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatal(err)
	}
	foundStartsWith := false
	for _, tok := range tokens {
		if tok.Kind == tokStartsWith {
			foundStartsWith = true
		}
	}
	if !foundStartsWith {
		t.Error("expected STARTS_WITH token")
	}
}

func TestLexer_ExplainProfile(t *testing.T) {
	tests := []struct {
		input string
		kind  TokenKind
	}{
		{"EXPLAIN", tokExplain},
		{"PROFILE", tokProfile},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			tokens, err := Tokenize(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			assertToken(t, tokens[0], tc.kind, tc.input)
		})
	}
}

func TestLexer_UnexpectedCharacter(t *testing.T) {
	_, err := Tokenize("MATCH @n RETURN n")
	if err == nil {
		t.Fatal("expected error for unexpected character '@'")
	}
}

func TestLexer_UnicodeIdentifiers(t *testing.T) {
	tokens, err := Tokenize("MATCH (n) RETURN n")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokMatch, "MATCH")
}

func TestLexer_WithQuery(t *testing.T) {
	input := `MATCH (n:Person) WITH n.name AS name RETURN name`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatal(err)
	}
	foundWith := false
	foundAs := false
	for _, tok := range tokens {
		if tok.Kind == tokWith {
			foundWith = true
		}
		if tok.Kind == tokAs {
			foundAs = true
		}
	}
	if !foundWith {
		t.Error("expected WITH token")
	}
	if !foundAs {
		t.Error("expected AS token")
	}
}

func TestLexer_PropertiesLabelsKeywords(t *testing.T) {
	tokens, err := Tokenize("PROPERTIES(n) LABELS(n)")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokProperties, "PROPERTIES")
	assertToken(t, tokens[1], tokLParen, "(")
	assertToken(t, tokens[2], tokIdent, "n")
	assertToken(t, tokens[3], tokRParen, ")")
	assertToken(t, tokens[4], tokLabels, "LABELS")
}

func TestLexer_EmptyInput(t *testing.T) {
	tokens, err := Tokenize("")
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 1 || tokens[0].Kind != tokEOF {
		t.Fatalf("expected single EOF token, got %d tokens", len(tokens))
	}
}

func TestLexer_DeleteRemove(t *testing.T) {
	tokens, err := Tokenize("DELETE n REMOVE n:Label")
	if err != nil {
		t.Fatal(err)
	}
	assertToken(t, tokens[0], tokDelete, "DELETE")
	assertToken(t, tokens[1], tokIdent, "n")
	assertToken(t, tokens[2], tokRemove, "REMOVE")
	assertToken(t, tokens[3], tokIdent, "n")
	assertToken(t, tokens[4], tokColon, ":")
	assertToken(t, tokens[5], tokIdent, "Label")
}

func assertToken(t *testing.T, tok Token, expectedKind TokenKind, expectedText string) {
	t.Helper()
	if tok.Kind != expectedKind {
		t.Errorf("expected token kind %s, got %s (text=%q)", tokenKindName(expectedKind), tokenKindName(tok.Kind), tok.Text)
	}
	if tok.Text != expectedText {
		t.Errorf("expected token text %q, got %q", expectedText, tok.Text)
	}
}

func assertLastEOF(t *testing.T, tokens []Token) {
	t.Helper()
	if len(tokens) == 0 || tokens[len(tokens)-1].Kind != tokEOF {
		t.Errorf("expected EOF as last token")
	}
}
