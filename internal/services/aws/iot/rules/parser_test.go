package rules

import (
	"fmt"
	"strings"
	"testing"
)

func TestNotEqualOperator(t *testing.T) {
	p := NewParser("SELECT * FROM 'topic' WHERE 1 <> 2")
	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if expr.Where == nil {
		t.Fatal("expected WHERE clause")
	}
	eval := NewEvaluator(map[string]interface{}{}, "topic", "client")
	result, err := eval.Eval(expr.Where)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	if !toBool(result) {
		t.Fatal("1 <> 2 should be true")
	}
}

func TestLessThanOrEqual(t *testing.T) {
	p := NewParser("SELECT * FROM 'topic' WHERE 1 <= 2")
	_, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error for <=: %v", err)
	}
}

func TestLessThan(t *testing.T) {
	p := NewParser("SELECT * FROM 'topic' WHERE 1 < 2")
	_, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error for <: %v", err)
	}
}

func TestNotEqualBang(t *testing.T) {
	p := NewParser("SELECT * FROM 'topic' WHERE 1 != 2")
	_, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error for !=: %v", err)
	}
}

func TestSelectStar(t *testing.T) {
	p := NewParser("SELECT * FROM 'topic'")
	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(expr.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(expr.Fields))
	}
	if _, ok := expr.Fields[0].Expression.(*StarExpr); !ok {
		t.Fatal("expected StarExpr")
	}
}

func TestParseErrorMissingFrom(t *testing.T) {
	p := NewParser("SELECT * 'topic'")
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for missing FROM keyword")
	}
}

func TestParseErrorMissingCloseParen(t *testing.T) {
	p := NewParser("SELECT abs(5 FROM 'topic'")
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for missing closing paren")
	}
}

func TestParseErrorTruncatedCast(t *testing.T) {
	p := NewParser("SELECT CAST(5 FROM 'topic'")
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for truncated CAST call")
	}
}

func TestInOperatorParse(t *testing.T) {
	tests := []struct {
		sql     string
		wantErr bool
	}{
		{"SELECT * FROM 't' WHERE color IN ('red', 'blue')", false},
		{"SELECT * FROM 't' WHERE id IN (1, 2, 3)", false},
		{"SELECT * FROM 't' WHERE id NOT IN (1, 2)", false},
		{"SELECT * FROM 't' WHERE x IN 5", true},
		{"SELECT * FROM 't' WHERE x NOT 5", true},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := NewParser(tt.sql)
			_, err := p.Parse()
			if tt.wantErr && err == nil {
				t.Fatal("expected parse error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}
		})
	}
}

func TestInOperatorEval(t *testing.T) {
	tests := []struct {
		sql     string
		payload map[string]interface{}
		want    bool
	}{
		{"SELECT * FROM 't' WHERE c IN ('red', 'blue')", map[string]interface{}{"c": "red"}, true},
		{"SELECT * FROM 't' WHERE c IN ('red', 'blue')", map[string]interface{}{"c": "green"}, false},
		{"SELECT * FROM 't' WHERE c NOT IN ('red')", map[string]interface{}{"c": "green"}, true},
		{"SELECT * FROM 't' WHERE c NOT IN ('red')", map[string]interface{}{"c": "red"}, false},
		{"SELECT * FROM 't' WHERE n IN (1, 2, 3)", map[string]interface{}{"n": float64(2)}, true},
		{"SELECT * FROM 't' WHERE n IN (1, 2, 3)", map[string]interface{}{"n": float64(9)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.sql+"_"+strings.TrimSpace(fmt.Sprintf("%v", tt.payload)), func(t *testing.T) {
			p := NewParser(tt.sql)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			if expr.Where == nil {
				t.Fatal("expected WHERE")
			}
			eval := NewEvaluator(tt.payload, "t", "c")
			result, err := eval.Eval(expr.Where)
			if err != nil {
				t.Fatalf("eval error: %v", err)
			}
			if toBool(result) != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, toBool(result))
			}
		})
	}
}
