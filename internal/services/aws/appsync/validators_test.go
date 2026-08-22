package appsync

import (
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
)

// TestValidateDescriptionUnicodeLengths pins that AppSync descriptions
// follow the Smithy Description @length(0, 255) trait counted in Unicode
// characters; the shape's "^.*$" pattern admits multibyte text.
func TestValidateDescriptionUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateDescription(strings.Repeat(cjk, 255)); err != nil {
		t.Errorf("255-character CJK description rejected: %v", err)
	}
	if err := validateDescription(strings.Repeat(cjk, 256)); err == nil {
		t.Error("256-character CJK description accepted")
	}
	if err := validateDescription(""); err != nil {
		t.Errorf("empty description rejected: %v", err)
	}
}

// TestValidateMappingTemplateUnicodeLengths pins that MappingTemplate
// @length(1,65536) is counted in Unicode characters; the shape's "^.*$"
// pattern admits multibyte Velocity templates.
func TestValidateMappingTemplateUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateMappingTemplate(strings.Repeat(cjk, 65536)); err != nil {
		t.Errorf("65536-character CJK mapping template rejected: %v", err)
	}
	if err := validateMappingTemplate(strings.Repeat(cjk, 65537)); err == nil {
		t.Error("65537-character CJK mapping template accepted")
	}
	if err := validateMappingTemplate(""); err == nil {
		t.Error("empty mapping template accepted")
	}
}

// TestValidateEnvVarValueUnicodeLengths pins that EnvironmentVariableValue
// @length(0,512) is counted in Unicode characters (no pattern; multibyte
// environment variable values are valid input).
func TestValidateEnvVarValueUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateEnvVarValue(strings.Repeat(cjk, 512)); err != nil {
		t.Errorf("512-character CJK environment variable value rejected: %v", err)
	}
	if err := validateEnvVarValue(strings.Repeat(cjk, 513)); err == nil {
		t.Error("513-character CJK environment variable value accepted")
	}
	if err := validateEnvVarValue(""); err != nil {
		t.Errorf("empty environment variable value rejected: %v", err)
	}
}

// TestParsePaginationOptionsTokenPattern pins the PaginationToken
// @pattern ^[\S]+$: a non-empty token containing whitespace is rejected
// with BadRequestException, while an empty token (first page) and opaque
// non-whitespace tokens pass.
func TestParsePaginationOptionsTokenPattern(t *testing.T) {
	mkReq := func(token string) *request.ParsedRequest {
		return &request.ParsedRequest{
			Parameters: map[string]interface{}{"nextToken": token},
		}
	}
	if _, err := parsePaginationOptions(mkReq("opaque-token_1.2")); err != nil {
		t.Errorf("valid token rejected: %v", err)
	}
	if _, err := parsePaginationOptions(mkReq("")); err != nil {
		t.Errorf("empty token (first page) rejected: %v", err)
	}
	for _, bad := range []string{"has space", "tab\tinside", "line\nbreak"} {
		if _, err := parsePaginationOptions(mkReq(bad)); err == nil {
			t.Errorf("token %q accepted", bad)
		}
	}
}
