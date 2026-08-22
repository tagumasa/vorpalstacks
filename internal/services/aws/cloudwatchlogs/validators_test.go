package cloudwatchlogs

import (
	"strings"
	"testing"
)

// TestValidateLogStreamNameUnicodeLengths pins that LogStreamName follows the
// Smithy @length(1, 512) trait counted in Unicode characters; the shape's
// ^[^:*]*$ pattern admits multibyte, so rune-legal names must not be
// rejected on byte length.
func TestValidateLogStreamNameUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateLogStreamName(strings.Repeat(cjk, 512)); err != nil {
		t.Errorf("512-character CJK log stream name rejected: %v", err)
	}
	if err := validateLogStreamName(strings.Repeat(cjk, 513)); err == nil {
		t.Error("513-character CJK log stream name accepted")
	}
}

// TestValidateFilterPatternUnicodeLengths pins that FilterPattern follows the
// Smithy @length(0, 1024) trait counted in Unicode characters (no pattern).
func TestValidateFilterPatternUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateFilterPattern(strings.Repeat(cjk, 1024)); err != nil {
		t.Errorf("1024-character CJK filter pattern rejected: %v", err)
	}
	if err := validateFilterPattern(strings.Repeat(cjk, 1025)); err == nil {
		t.Error("1025-character CJK filter pattern accepted")
	}
}

// TestValidatePolicyDocumentUnicodeLengths pins that the PolicyDocument
// @length(1, 51200) trait is counted in Unicode characters (no pattern).
func TestValidatePolicyDocumentUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validatePolicyDocument(strings.Repeat(cjk, 51200)); err != nil {
		t.Errorf("51200-character CJK policy document rejected: %v", err)
	}
	if err := validatePolicyDocument(strings.Repeat(cjk, 51201)); err == nil {
		t.Error("51201-character CJK policy document accepted")
	}
}

// TestValidateQueryDefinitionNameUnicodeLengths pins that
// QueryDefinitionName follows the Smithy @length(1, 255) trait counted in
// Unicode characters (no pattern).
func TestValidateQueryDefinitionNameUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateQueryDefinitionName(strings.Repeat(cjk, 255)); err != nil {
		t.Errorf("255-character CJK query definition name rejected: %v", err)
	}
	if err := validateQueryDefinitionName(strings.Repeat(cjk, 256)); err == nil {
		t.Error("256-character CJK query definition name accepted")
	}
}
