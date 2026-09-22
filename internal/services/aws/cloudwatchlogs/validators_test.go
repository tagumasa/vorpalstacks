package cloudwatchlogs

import (
	"strings"
	"testing"
	"unicode/utf8"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// TestValidatorUnicodeLengthLimits pins that the string validators enforce
// their Smithy @length traits in Unicode characters: LogStreamName's
// ^[^:*]*$ pattern admits multibyte and the other three members carry no
// pattern, so rune-legal values must not be rejected on byte length. Each
// row is one member's trait window, pinned at both edges; the subtests
// carry the names the four per-member tests wore before the table.
func TestValidatorUnicodeLengthLimits(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes
	tests := []struct {
		name     string
		validate func(string) error
		limit    int
	}{
		// LogStreamName @length(1, 512).
		{"TestValidateLogStreamNameUnicodeLengths", validateLogStreamName, 512},
		// FilterPattern @length(0, 1024).
		{"TestValidateFilterPatternUnicodeLengths", validateFilterPattern, 1024},
		// PolicyDocument @length(1, 51200).
		{"TestValidatePolicyDocumentUnicodeLengths", validatePolicyDocument, 51200},
		// QueryDefinitionName @length(1, 255).
		{"TestValidateQueryDefinitionNameUnicodeLengths", validateQueryDefinitionName, 255},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.validate(strings.Repeat(cjk, tt.limit)); err != nil {
				t.Errorf("limit-length CJK value rejected: %v", err)
			}
			if err := tt.validate(strings.Repeat(cjk, tt.limit+1)); err == nil {
				t.Error("over-limit CJK value accepted")
			}
		})
	}
}

// A negative limit is a supplied value outside the range trait's
// 1..max window, not an absent member: it rejects rather than silently
// becoming the default page.
func TestResolveListLimitNegativeRejects(t *testing.T) {
	if got, err := validateListLimit(0, 50, 10000); err != nil || got != 50 {
		t.Fatalf("absent limit: got %d, %v; want the default 50", got, err)
	}
	if _, err := validateListLimit(-5, 50, 10000); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("negative limit: %v, want InvalidParameterException", err)
	}
}

// The read operations' name-or-identifier pair is exclusive and
// exhaustive: both members reject, neither rejects.
func TestLogGroupNameOrIdentifierExclusive(t *testing.T) {
	if name, err := logGroupNameOrIdentifier(map[string]interface{}{"logGroupName": "g"}); err != nil || name != "g" {
		t.Fatalf("name alone: %q, %v", name, err)
	}
	if name, err := logGroupNameOrIdentifier(map[string]interface{}{"logGroupIdentifier": "arn:aws:logs:us-east-1:000000000000:log-group:x"}); err != nil || name != "x" {
		t.Fatalf("identifier alone: %q, %v", name, err)
	}
	if _, err := logGroupNameOrIdentifier(map[string]interface{}{
		"logGroupName": "g", "logGroupIdentifier": "arn:aws:logs:us-east-1:000000000000:log-group:x",
	}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("both members: %v, want InvalidParameterException", err)
	}
	if _, err := logGroupNameOrIdentifier(map[string]interface{}{}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("neither member: %v, want InvalidParameterException", err)
	}
}

// The transformed-message digest separates its input fields by length:
// '|' is legal in both stream names and messages, so separator-only
// framing let two different (stream, message) pairs collide.
func TestTransformedMessageDigestNoSeparatorCollision(t *testing.T) {
	if logsstore.TransformedMessageDigest(5, "a|b", "c") == logsstore.TransformedMessageDigest(5, "a", "b|c") {
		t.Fatal("separator-framed digest collided on distinct (stream, message) pairs")
	}
}

// A long object key with multi-byte runes truncates to a valid UTF-8
// stream name — mid-rune slicing would mint an invalid one.
func TestImportStreamNameRuneSafeTruncation(t *testing.T) {
	long := strings.Repeat("オ", logsstore.MaxLogStreamNameLength)
	name := importStreamName(long)
	if len(name) > logsstore.MaxLogStreamNameLength {
		t.Fatalf("truncated name length %d exceeds the limit", len(name))
	}
	if !utf8.ValidString(name) {
		t.Fatal("truncated stream name is not valid UTF-8")
	}
}

// The bearer auth-scheme matches case-insensitively (RFC 7235).
func TestHasBearerSchemeCaseInsensitive(t *testing.T) {
	for _, auth := range []string{"Bearer tok", "bearer tok", "BEARER tok"} {
		if !HasBearerScheme(auth) {
			t.Fatalf("%q must match the bearer scheme", auth)
		}
	}
	if HasBearerScheme("AWS4-HMAC-SHA256 Credential=x") || HasBearerScheme("") {
		t.Fatal("non-bearer schemes must not match")
	}
}
