package cloudwatch

import (
	"strings"
	"testing"
)

// TestValidateAlarmDescriptionUnicodeLengths pins that alarm descriptions
// follow the Smithy AlarmDescription @length(0, 1024) trait counted in
// Unicode characters; the shape carries no pattern, so multibyte text is
// valid input and must not be rejected on byte length.
func TestValidateAlarmDescriptionUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateAlarmDescription(strings.Repeat(cjk, 1024)); err != nil {
		t.Errorf("1024-character CJK alarm description rejected: %v", err)
	}
	if err := validateAlarmDescription(strings.Repeat(cjk, 1025)); err == nil {
		t.Error("1025-character CJK alarm description accepted")
	}
	if err := validateAlarmDescription(""); err != nil {
		t.Errorf("empty alarm description rejected: %v", err)
	}
}

// TestValidateStateReasonUnicodeLengths pins that alarm state reasons
// follow the Smithy StateReason @length(0, 1023) trait counted in Unicode
// characters.
func TestValidateStateReasonUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateStateReason(strings.Repeat(cjk, 1023)); err != nil {
		t.Errorf("1023-character CJK state reason rejected: %v", err)
	}
	if err := validateStateReason(strings.Repeat(cjk, 1024)); err == nil {
		t.Error("1024-character CJK state reason accepted")
	}
}
