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

// TestValidateNamespaceUnicodeLengths pins that Namespace follows the Smithy
// @length(1, 255) trait counted in Unicode characters; the shape's pattern
// only forbids a leading colon, so multibyte namespaces are valid input.
func TestValidateNamespaceUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateNamespace(strings.Repeat(cjk, 255)); err != nil {
		t.Errorf("255-character CJK namespace rejected: %v", err)
	}
	if err := validateNamespace(strings.Repeat(cjk, 256)); err == nil {
		t.Error("256-character CJK namespace accepted")
	}
}

// TestValidateMetricNameUnicodeLengths pins that MetricName follows the
// Smithy @length(1, 255) trait counted in Unicode characters (no pattern).
func TestValidateMetricNameUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateMetricName(strings.Repeat(cjk, 255)); err != nil {
		t.Errorf("255-character CJK metric name rejected: %v", err)
	}
	if err := validateMetricName(strings.Repeat(cjk, 256)); err == nil {
		t.Error("256-character CJK metric name accepted")
	}
}

// TestValidateThresholdMetricIdUnicodeLengths pins that ThresholdMetricId
// follows the Smithy MetricId @length(1, 255) trait counted in Unicode
// characters (no pattern).
func TestValidateThresholdMetricIdUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateThresholdMetricId(strings.Repeat(cjk, 255)); err != nil {
		t.Errorf("255-character CJK threshold metric id rejected: %v", err)
	}
	if err := validateThresholdMetricId(strings.Repeat(cjk, 256)); err == nil {
		t.Error("256-character CJK threshold metric id accepted")
	}
	if err := validateThresholdMetricId(""); err != nil {
		t.Errorf("empty threshold metric id rejected: %v", err)
	}
}
