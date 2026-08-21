package ssm

import (
	"errors"
	"strings"
	"testing"
)

// TestValidateDataType_Integration verifies that the AWS-documented
// aws:ssm:integration DataType value is accepted.
func TestValidateDataType_Integration(t *testing.T) {
	for _, dt := range []string{"text", "aws:ec2:image", "aws:ssm:integration"} {
		if err := validateDataType(dt); err != nil {
			t.Fatalf("validateDataType(%q) = %v, want nil", dt, err)
		}
	}
	if err := validateDataType("unknown:type"); !errors.Is(err, ErrInvalidParameterValue) {
		t.Fatalf("validateDataType(unknown) = %v, want ErrInvalidParameterValue", err)
	}
}

// TestValidateKeyID_Length verifies the Smithy ParameterKeyId 256-character
// length cap.
func TestValidateKeyID_Length(t *testing.T) {
	if err := validateKeyID(strings.Repeat("k", 256)); err != nil {
		t.Fatalf("256-char key id rejected: %v", err)
	}
	if err := validateKeyID(strings.Repeat("k", 257)); !errors.Is(err, ErrInvalidParameterValue) {
		t.Fatalf("257-char key id err = %v, want ErrInvalidParameterValue", err)
	}
}

// TestValidateAllowedPattern_Length verifies the Smithy AllowedPattern
// 1024-character length cap.
func TestValidateAllowedPattern_Length(t *testing.T) {
	if err := validateAllowedPattern(strings.Repeat("a", 1024)); err != nil {
		t.Fatalf("1024-char pattern rejected: %v", err)
	}
	if err := validateAllowedPattern(strings.Repeat("a", 1025)); !errors.Is(err, ErrInvalidParameterValue) {
		t.Fatalf("1025-char pattern err = %v, want ErrInvalidParameterValue", err)
	}
}

// TestNormalisePutParameter_RequiredValue verifies that an empty Value is
// rejected (Smithy marks Value as required).
func TestNormalisePutParameter_RequiredValue(t *testing.T) {
	_, err := normalisePutParameter(ParameterPutFields{Name: "p", Value: ""})
	if !errors.Is(err, ErrInvalidParameterValue) {
		t.Fatalf("err = %v, want ErrInvalidParameterValue", err)
	}
}

// TestNormalisePutParameter_DescriptionLength verifies the Smithy
// ParameterDescription 1024-character length cap.
func TestNormalisePutParameter_DescriptionLength(t *testing.T) {
	_, err := normalisePutParameter(ParameterPutFields{
		Name:        "p",
		Value:       "v",
		Description: strings.Repeat("d", 1025),
	})
	if !errors.Is(err, ErrInvalidParameterValue) {
		t.Fatalf("err = %v, want ErrInvalidParameterValue", err)
	}
}

// TestNormalisePutParameter_DescriptionUnicodeLength pins that the
// ParameterDescription cap is counted in Unicode characters per the
// Smithy @length trait: 1024 CJK characters (3072 bytes) are in range.
func TestNormalisePutParameter_DescriptionUnicodeLength(t *testing.T) {
	cjk := strings.Repeat("\u65e5", 1024)
	if _, err := normalisePutParameter(ParameterPutFields{
		Name:        "p",
		Value:       "v",
		Description: cjk,
	}); err != nil {
		t.Fatalf("1024-character CJK description rejected: %v", err)
	}
}

// TestFiltersFromList_FailClosed verifies that malformed ParameterFilters
// input is rejected instead of being silently dropped (a dropped filter
// would return unfiltered results).
func TestFiltersFromList_FailClosed(t *testing.T) {
	if _, err := filtersFromList("not-a-list"); !errors.Is(err, ErrInvalidFilterValue) {
		t.Fatalf("string err = %v, want ErrInvalidFilterValue", err)
	}
	if _, err := filtersFromList([]interface{}{"not-an-object"}); !errors.Is(err, ErrInvalidFilterValue) {
		t.Fatalf("non-object entry err = %v, want ErrInvalidFilterValue", err)
	}
	noValues := []interface{}{map[string]interface{}{"Key": "Name"}}
	if _, err := filtersFromList(noValues); !errors.Is(err, ErrInvalidFilterValue) {
		t.Fatalf("missing Values err = %v, want ErrInvalidFilterValue", err)
	}
	valid := []interface{}{map[string]interface{}{"Key": "Name", "Values": []interface{}{"p1"}}}
	filters, err := filtersFromList(valid)
	if err != nil || len(filters) != 1 {
		t.Fatalf("valid filters err = %v, len = %d", err, len(filters))
	}
	if _, err := filtersFromList(nil); err != nil {
		t.Fatalf("absent key err = %v, want nil", err)
	}
}

func TestFiltersFromQueryParams_FailClosed(t *testing.T) {
	// An entry with a Key but no Values must be rejected identically to
	// the JSON body path, and must not silently drop subsequent entries.
	params := map[string]interface{}{
		"ParameterFilters.member.1.Key":             "Name",
		"ParameterFilters.member.2.Key":             "Type",
		"ParameterFilters.member.2.Values.member.1": "String",
	}
	if _, err := filtersFromQueryParams(params, "ParameterFilters"); !errors.Is(err, ErrInvalidFilterValue) {
		t.Fatalf("missing Values err = %v, want ErrInvalidFilterValue", err)
	}

	valid := map[string]interface{}{
		"ParameterFilters.member.1.Key":             "Name",
		"ParameterFilters.member.1.Values.member.1": "p1",
		"ParameterFilters.member.1.Values.member.2": "p2",
		"ParameterFilters.member.2.Key":             "Type",
		"ParameterFilters.member.2.Option":          "Equals",
		"ParameterFilters.member.2.Values.member.1": "String",
	}
	filters, err := filtersFromQueryParams(valid, "ParameterFilters")
	if err != nil || len(filters) != 2 {
		t.Fatalf("valid filters err = %v, len = %d", err, len(filters))
	}
	if len(filters[0].Values) != 2 || filters[1].Option != "Equals" {
		t.Fatalf("unexpected filters: %+v", filters)
	}
}
