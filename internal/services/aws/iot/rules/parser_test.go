package rules

import "testing"

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
