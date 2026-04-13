package policy

import (
	"testing"
)

func TestPolicyError(t *testing.T) {
	err := &PolicyError{
		Code:    "TestCode",
		Message: "TestMessage",
	}

	if err.Error() != "TestCode: TestMessage" {
		t.Errorf("expected 'TestCode: TestMessage', got '%s'", err.Error())
	}
}

func TestNewPolicyError(t *testing.T) {
	err := NewPolicyError("Code", "Message")
	if err.Code != "Code" {
		t.Errorf("expected Code, got %s", err.Code)
	}
	if err.Message != "Message" {
		t.Errorf("expected Message, got %s", err.Message)
	}
}

func TestIsAccessDenied(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "AccessDenied error",
			err:      ErrAccessDenied,
			expected: true,
		},
		{
			name:     "Other error",
			err:      NewPolicyError("OtherCode", "Message"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "non-PolicyError",
			err:      &customError{"message"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAccessDenied(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func TestErrVariables(t *testing.T) {
	if ErrInvalidPolicyDocument.Code != "InvalidPolicyDocument" {
		t.Errorf("expected InvalidPolicyDocument, got %s", ErrInvalidPolicyDocument.Code)
	}
	if ErrMalformedPolicy.Code != "MalformedPolicy" {
		t.Errorf("expected MalformedPolicy, got %s", ErrMalformedPolicy.Code)
	}
	if ErrInvalidEffect.Code != "InvalidEffect" {
		t.Errorf("expected InvalidEffect, got %s", ErrInvalidEffect.Code)
	}
	if ErrMissingAction.Code != "MissingAction" {
		t.Errorf("expected MissingAction, got %s", ErrMissingAction.Code)
	}
	if ErrMissingResource.Code != "MissingResource" {
		t.Errorf("expected MissingResource, got %s", ErrMissingResource.Code)
	}
	if ErrAccessDenied.Code != "AccessDenied" {
		t.Errorf("expected AccessDenied, got %s", ErrAccessDenied.Code)
	}
}
