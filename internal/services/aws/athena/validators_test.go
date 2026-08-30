package athena

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
)

// TestValidateTagsReportsOffenderAccurately pins that tag violation
// messages quote the offender found by the shared checker: the reported
// length is the Unicode character count of the actual offending key or
// value, never an in-range count of a rune-legal multibyte key, and the
// choice does not depend on map iteration order.
func TestValidateTagsReportsOffenderAccurately(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	// The 129-character ASCII key is the offender; the 100-character CJK
	// key (300 bytes) is rune-legal and must never be blamed for it.
	tags := map[string]string{
		strings.Repeat("a", 129): "v",
		strings.Repeat(cjk, 100): "v",
	}
	want := "TagKey length must be between 1 and 128 (got 129)"
	for i := 0; i < 100; i++ {
		err := validateTags(tags)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("validateTags = %v, want message containing %q", err, want)
		}
	}

	// Value lengths are counted in Unicode characters, not bytes.
	err := validateTags(map[string]string{"env": strings.Repeat(cjk, 257)})
	want = "TagValue length must be between 0 and 256 (got 257)"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("validateTags 257-character CJK value = %v, want message containing %q", err, want)
	}
	if err := validateTags(map[string]string{"env": strings.Repeat(cjk, 256)}); err != nil {
		t.Errorf("validateTags 256-character CJK value = %v, want nil", err)
	}
}

// TestValidateQueryStringSizeUnicodeLengths pins that the QueryString
// @length(1, 262144) trait is counted in Unicode characters: the shape
// carries no pattern, so a query whose text is rune-legal must not be
// rejected on byte length.
func TestValidateQueryStringSizeUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateQueryStringSize(strings.Repeat(cjk, 262144)); err != nil {
		t.Errorf("262144-character CJK query rejected: %v", err)
	}
	if err := validateQueryStringSize(strings.Repeat(cjk, 262145)); err == nil {
		t.Error("262145-character CJK query accepted")
	}
	if err := validateQueryStringSize(""); err == nil {
		t.Error("empty query accepted")
	}
}

// TestParameterViolationsEmitInvalidRequestException pins that every shared
// parameter validator emits InvalidRequestException with HTTP 400: the Athena
// model defines no InvalidParameterException shape and every operation's
// errors list declares InvalidRequestException, so a constraint violation
// must ride the operation-declared error.
func TestParameterViolationsEmitInvalidRequestException(t *testing.T) {
	tooManyTags := make(map[string]string, 51)
	for i := 0; i < 51; i++ {
		tooManyTags[fmt.Sprintf("key%d", i)] = "v"
	}

	cases := []struct {
		name string
		err  error
	}{
		{"validateStringLength", validateStringLength("Name", "", 1, 128)},
		{"validateWorkGroupName", validateWorkGroupName("bad name!")},
		{"validateStatementName", validateStatementName("9starts-with-digit")},
		{"validateCapacityReservationName", validateCapacityReservationName("bad!")},
		{"validateTags too many", validateTags(tooManyTags)},
		{"validateTags reserved key", validateTags(map[string]string{"aws:reserved": "v"})},
		{"validateClientRequestToken", validateClientRequestToken("short")},
		{"validateExecutionRole", validateExecutionRole("not-a-valid-arn-but-long-enough-to-pass-length")},
		{"validateBytesScannedCutoff", validateBytesScannedCutoff(1)},
		{"validateDataCatalogType", validateDataCatalogType("BOGUS")},
		{"validateWorkGroupState", validateWorkGroupState("BOGUS")},
		{"validateMaxResults", func() error {
			_, err := validateMaxResults(map[string]interface{}{"MaxResults": 51}, 50, 1, 50)
			return err
		}()},
		{"resolveMaxResults", func() error { _, err := resolveMaxResults(51, true, 50, 1, 50); return err }()},
	}
	for _, tc := range cases {
		if tc.err == nil {
			t.Errorf("%s: expected a violation error", tc.name)
			continue
		}
		var awsErr *awserrors.AWSError
		if !errors.As(tc.err, &awsErr) {
			t.Errorf("%s: error is not *awserrors.AWSError: %v", tc.name, tc.err)
			continue
		}
		if awsErr.Code != "InvalidRequestException" {
			t.Errorf("%s: code = %q, want InvalidRequestException", tc.name, awsErr.Code)
		}
		if awsErr.HTTPStatus != http.StatusBadRequest {
			t.Errorf("%s: HTTP status = %d, want 400", tc.name, awsErr.HTTPStatus)
		}
	}
}
