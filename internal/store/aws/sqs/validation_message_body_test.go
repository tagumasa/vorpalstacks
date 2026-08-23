package sqs

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateMessageBody(t *testing.T) {
	cases := []struct {
		name string
		body string
		want error
	}{
		{"plain text", "hello world", nil},
		{"tab newline carriage return", "a\tb\nc\rd", nil},
		{"json", `{"key": "value"}`, nil},
		{"nul byte", "a\x00b", ErrInvalidMessageContents},
		{"control byte", "\x01", ErrInvalidMessageContents},
		{"0x1f delimiter", "a\x1fb", ErrInvalidMessageContents},
		{"unicode supplement", strings.Repeat("😀", 4), nil},
	}
	for _, tc := range cases {
		got := validateMessageBody(tc.body)
		if !errors.Is(got, tc.want) {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestValidateFifoIdentifier(t *testing.T) {
	if err := validateFifoIdentifier(""); err != nil {
		t.Errorf("empty identifier must be allowed: %v", err)
	}
	if err := validateFifoIdentifier("group-1.2_3!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"); err != nil {
		t.Errorf("documented punctuation set rejected: %v", err)
	}
	if err := validateFifoIdentifier(strings.Repeat("a", MaxFifoIdLength)); err != nil {
		t.Errorf("128-character identifier rejected: %v", err)
	}
	if err := validateFifoIdentifier(strings.Repeat("a", MaxFifoIdLength+1)); err == nil {
		t.Errorf("129-character identifier accepted")
	}
	if err := validateFifoIdentifier("has space"); err == nil {
		t.Errorf("identifier with space accepted")
	}
	if err := validateFifoIdentifier("non-ascii-é"); err == nil {
		t.Errorf("identifier with non-ASCII character accepted")
	}
}
