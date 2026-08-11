package rules

import "testing"

func TestOpMod_ZeroDivision(t *testing.T) {
	result, err := opMod(float64(10), float64(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(unknownValue); !ok {
		t.Fatalf("expected unknownValue for modulo by zero, got %T", result)
	}
}

func TestOpMod_Normal(t *testing.T) {
	result, err := opMod(float64(10), float64(3))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != float64(1) {
		t.Errorf("expected 10 mod 3 = 1, got %v", result)
	}
}

func TestOpDiv_ZeroDivision(t *testing.T) {
	result, err := opDiv(float64(10), float64(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(unknownValue); !ok {
		t.Fatalf("expected unknownValue for division by zero, got %T", result)
	}
}

func TestOpDiv_Normal(t *testing.T) {
	result, err := opDiv(float64(10), float64(3))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != float64(10)/float64(3) {
		t.Errorf("expected 10/3, got %v", result)
	}
}
