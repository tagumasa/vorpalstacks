package eventbridge

import (
	"strings"
	"testing"
)

// TestValidateDetailTypeUnicodeLengths pins that the PutEvents DetailType
// maximum of 128 characters is counted in Unicode characters (member
// documentation: "maximum of 128 characters").
func TestValidateDetailTypeUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if !validateDetailType(strings.Repeat(cjk, 128)) {
		t.Error("128-character CJK detail type rejected")
	}
	if validateDetailType(strings.Repeat(cjk, 129)) {
		t.Error("129-character CJK detail type accepted")
	}
	if validateDetailType("") {
		t.Error("empty detail type accepted")
	}
}

// TestValidateSourceUnicodeLengths pins that the PutEvents Source bound
// ("Length constraints: minimum length of 1, maximum length of 256" per the
// API reference) is counted in Unicode characters, matching DetailType.
func TestValidateSourceUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if !validateSource(strings.Repeat(cjk, 256)) {
		t.Error("256-character CJK source rejected")
	}
	if validateSource(strings.Repeat(cjk, 257)) {
		t.Error("257-character CJK source accepted")
	}
	if validateSource("") {
		t.Error("empty source accepted")
	}
}

// TestValidateDescriptionUnicodeLengths pins that the shared description
// @length(0,512) trait is counted in Unicode characters; the description
// shapes carry no pattern, so multibyte text is valid input.
func TestValidateDescriptionUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if !validateDescription(strings.Repeat(cjk, 512)) {
		t.Error("512-character CJK description rejected")
	}
	if validateDescription(strings.Repeat(cjk, 513)) {
		t.Error("513-character CJK description accepted")
	}
	if !validateDescription("") {
		t.Error("empty description rejected")
	}
}

// TestValidateEventPatternLengthUnicodeLengths pins that the EventPattern
// @length(0,4096) trait is counted in Unicode characters; the shape carries
// no pattern, so JSON patterns containing multibyte values are rune-legal.
func TestValidateEventPatternLengthUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if !validateEventPatternLength(strings.Repeat(cjk, 4096)) {
		t.Error("4096-character CJK event pattern rejected")
	}
	if validateEventPatternLength(strings.Repeat(cjk, 4097)) {
		t.Error("4097-character CJK event pattern accepted")
	}
	if !validateEventPatternLength("") {
		t.Error("empty event pattern rejected")
	}
}

// TestValidateEventBusPolicySize pins the documented 10 KB (10240-byte)
// ceiling on event bus resource policies: "The permission policy on the
// event bus cannot exceed 10 KB in size" (PutPermission API Reference).
func TestValidateEventBusPolicySize(t *testing.T) {
	if err := validateEventBusPolicySize(strings.Repeat("a", 10000)); err != nil {
		t.Errorf("10000-byte policy rejected: %v", err)
	}
	if err := validateEventBusPolicySize(strings.Repeat("a", 10240)); err != nil {
		t.Errorf("10240-byte policy rejected: %v", err)
	}
	if err := validateEventBusPolicySize(strings.Repeat("a", 10241)); err == nil {
		t.Error("10241-byte policy accepted")
	}
}
