package gremlinparser

import (
	"testing"
)

func TestLexer_Basic(t *testing.T) {
	tokens, err := Tokenize("g.V().out('knows')")
	if err != nil {
		t.Fatalf("tokenize error: %v", err)
	}
	expected := []struct {
		kind TokenKind
		text string
	}{
		{tokIdentifier, "g"},
		{tokDot, "."},
		{tokIdentifier, "V"},
		{tokLParen, "("},
		{tokRParen, ")"},
		{tokDot, "."},
		{tokIdentifier, "out"},
		{tokLParen, "("},
		{tokString, "knows"},
		{tokRParen, ")"},
		{tokEOF, ""},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %v", len(expected), len(tokens), tokens)
	}
	for i, exp := range expected {
		if tokens[i].Kind != exp.kind {
			t.Errorf("token %d: expected kind %d, got %d", i, exp.kind, tokens[i].Kind)
		}
		if tokens[i].Text != exp.text {
			t.Errorf("token %d: expected text %q, got %q", i, exp.text, tokens[i].Text)
		}
	}
}

func TestLexer_Numbers(t *testing.T) {
	tests := []struct {
		input    string
		expected []struct {
			kind TokenKind
			text string
		}
	}{
		{"42", []struct {
			kind TokenKind
			text string
		}{{tokInteger, "42"}}},
		{"3.14", []struct {
			kind TokenKind
			text string
		}{{tokFloat, "3.14"}}},
		{"-1", []struct {
			kind TokenKind
			text string
		}{{tokInteger, "-1"}}},
		{"0xFF", []struct {
			kind TokenKind
			text string
		}{{tokInteger, "0xFF"}}},
		{"1.0e-10", []struct {
			kind TokenKind
			text string
		}{{tokFloat, "1.0e-10"}}},
		{"32i", []struct {
			kind TokenKind
			text string
		}{{tokInteger, "32i"}}},
		{"3.14f", []struct {
			kind TokenKind
			text string
		}{{tokFloat, "3.14f"}}},
	}
	for _, tt := range tests {
		tokens, err := Tokenize(tt.input)
		if err != nil {
			t.Errorf("tokenize(%q) error: %v", tt.input, err)
			continue
		}
		if len(tokens) != len(tt.expected)+1 {
			t.Errorf("tokenize(%q): expected %d tokens, got %d", tt.input, len(tt.expected)+1, len(tokens))
			continue
		}
		for i, exp := range tt.expected {
			if tokens[i].Kind != exp.kind || tokens[i].Text != exp.text {
				t.Errorf("tokenize(%q)[%d]: expected (%d, %q), got (%d, %q)", tt.input, i, exp.kind, exp.text, tokens[i].Kind, tokens[i].Text)
			}
		}
	}
}

func TestLexer_Strings(t *testing.T) {
	tests := []struct {
		input string
		text  string
	}{
		{"'hello'", "hello"},
		{`"hello"`, "hello"},
		{`'it\'s'`, "it's"},
		{`"say \"hi\""`, `say "hi"`},
		{"''", ""},
	}
	for _, tt := range tests {
		tokens, err := Tokenize(tt.input)
		if err != nil {
			t.Errorf("tokenize(%q) error: %v", tt.input, err)
			continue
		}
		if tokens[0].Kind != tokString || tokens[0].Text != tt.text {
			t.Errorf("tokenize(%q): expected string %q, got (%d, %q)", tt.input, tt.text, tokens[0].Kind, tokens[0].Text)
		}
	}
}

func TestLexer_Comments(t *testing.T) {
	tokens, err := Tokenize("g.V() // comment\n.out()")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	idents := []string{}
	for _, tok := range tokens {
		if tok.Kind == tokIdentifier {
			idents = append(idents, tok.Text)
		}
	}
	if len(idents) != 3 || idents[0] != "g" || idents[1] != "V" || idents[2] != "out" {
		t.Errorf("expected [g, V, out], got %v", idents)
	}
}

func TestLexer_Booleans(t *testing.T) {
	tokens, err := Tokenize("true false null")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if tokens[0].Kind != tokTrue || tokens[1].Kind != tokFalse || tokens[2].Kind != tokNull {
		t.Errorf("expected true/false/null, got %v %v %v", tokens[0].Kind, tokens[1].Kind, tokens[2].Kind)
	}
}

func TestLexer_AnonymousTraversal(t *testing.T) {
	tokens, err := Tokenize("__")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if tokens[0].Kind != tokUnderscoreUnderscore {
		t.Errorf("expected __, got %d", tokens[0].Kind)
	}
}

func TestLexer_MapLiteral(t *testing.T) {
	tokens, err := Tokenize("[name: 'marko', age: 29]")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	kinds := []TokenKind{}
	for _, tok := range tokens {
		if tok.Kind != tokEOF {
			kinds = append(kinds, tok.Kind)
		}
	}
	expected := []TokenKind{tokLBracket, tokIdentifier, tokColon, tokString, tokComma, tokIdentifier, tokColon, tokInteger, tokRBracket}
	if len(kinds) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %v", len(expected), len(kinds), kinds)
	}
	for i, k := range expected {
		if kinds[i] != k {
			t.Errorf("token %d: expected %d, got %d", i, k, kinds[i])
		}
	}
}
