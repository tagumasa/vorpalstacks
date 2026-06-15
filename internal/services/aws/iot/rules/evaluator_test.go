package rules

import "testing"

func TestMatchLikeUnderscore(t *testing.T) {
	cases := []struct {
		text    string
		pattern string
		want    bool
	}{
		{"dev1c", "dev_c", true},
		{"dev12c", "dev_c", false},
		{"devc", "dev_c", false},
		{"a_c", "a_c", true},
	}
	for _, c := range cases {
		got := matchLike(c.text, c.pattern)
		if got != c.want {
			t.Errorf("matchLike(%q, %q) = %v, want %v", c.text, c.pattern, got, c.want)
		}
	}
}

func TestMatchLikePercent(t *testing.T) {
	cases := []struct {
		text    string
		pattern string
		want    bool
	}{
		{"device123", "dev%", true},
		{"my_device", "%device", true},
		{"my_device_42", "%dev%", true},
		{"hello", "hello", true},
		{"world", "hello", false},
	}
	for _, c := range cases {
		got := matchLike(c.text, c.pattern)
		if got != c.want {
			t.Errorf("matchLike(%q, %q) = %v, want %v", c.text, c.pattern, got, c.want)
		}
	}
}

func TestMatchLikeInterleaved(t *testing.T) {
	cases := []struct {
		text    string
		pattern string
		want    bool
	}{
		{"abcde", "a%c_e", true},
		{"axbcye", "a%c_e", true},
		{"abcde", "a%c_d", false},
	}
	for _, c := range cases {
		got := matchLike(c.text, c.pattern)
		if got != c.want {
			t.Errorf("matchLike(%q, %q) = %v, want %v", c.text, c.pattern, got, c.want)
		}
	}
}

func TestSubstringOneIndexed(t *testing.T) {
	_, result, err := EvaluateSQL(
		"SELECT substring('hello', 1) as s FROM 't'",
		map[string]interface{}{}, "t", "c",
	)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	s, ok := result["s"].(string)
	if !ok || s != "hello" {
		t.Errorf("substring('hello', 1) = %v, want \"hello\"", result["s"])
	}
}

func TestSubstringWithLength(t *testing.T) {
	_, result, err := EvaluateSQL(
		"SELECT substring('hello', 1, 2) as s FROM 't'",
		map[string]interface{}{}, "t", "c",
	)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	s, ok := result["s"].(string)
	if !ok || s != "he" {
		t.Errorf("substring('hello', 1, 2) = %v, want \"he\"", result["s"])
	}
}

func TestSubstringMidPosition(t *testing.T) {
	_, result, err := EvaluateSQL(
		"SELECT substring('hello', 3) as s FROM 't'",
		map[string]interface{}{}, "t", "c",
	)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	s, ok := result["s"].(string)
	if !ok || s != "llo" {
		t.Errorf("substring('hello', 3) = %v, want \"llo\"", result["s"])
	}
}

func TestSubstringBeyondLength(t *testing.T) {
	_, result, err := EvaluateSQL(
		"SELECT substring('hello', 6) as s FROM 't'",
		map[string]interface{}{}, "t", "c",
	)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	s, ok := result["s"].(string)
	if !ok || s != "" {
		t.Errorf("substring('hello', 6) = %v, want \"\"", result["s"])
	}
}

func TestCastNumeric(t *testing.T) {
	_, result, err := EvaluateSQL(
		"SELECT cast(x AS NUMERIC) as v FROM 't'",
		map[string]interface{}{"x": 123}, "t", "c",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := result["v"].(float64)
	if v != 123.0 {
		t.Errorf("cast(123 AS NUMERIC) = %v, want 123.0", v)
	}
}

func TestCastString(t *testing.T) {
	_, result, err := EvaluateSQL(
		"SELECT cast(x AS STRING) as v FROM 't'",
		map[string]interface{}{"x": 123}, "t", "c",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := result["v"].(string)
	if v != "123" {
		t.Errorf("cast(123 AS STRING) = %v, want '123'", v)
	}
}

func TestCastBoolean(t *testing.T) {
	_, result, err := EvaluateSQL(
		"SELECT cast(x AS BOOLEAN) as v FROM 't'",
		map[string]interface{}{"x": "true"}, "t", "c",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := result["v"].(bool)
	if !v {
		t.Errorf("cast('true' AS BOOLEAN) = %v, want true", v)
	}
}

func TestCastInt(t *testing.T) {
	_, result, err := EvaluateSQL(
		"SELECT cast(x AS INT) as v FROM 't'",
		map[string]interface{}{"x": 42.7}, "t", "c",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := result["v"].(int64)
	if v != 42 {
		t.Errorf("cast(42.7 AS INT) = %v, want 42", v)
	}
}
