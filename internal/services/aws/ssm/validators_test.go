package ssm

import (
	"errors"
	"strings"
	"testing"
)

// TestValidateParameterNameList_Length pins the Smithy ParameterNameList
// contract (required, @length 1-10): an empty or oversized list is rejected
// with ValidationException. A missing Names member is unreachable through
// the AWS SDK (client-side required validation), so the omission shape is
// pinned here as the nil case.
func TestValidateParameterNameList_Length(t *testing.T) {
	if err := validateParameterNameList(nil); !errors.Is(err, ErrValidationException) {
		t.Fatalf("nil err = %v, want ErrValidationException", err)
	}
	if err := validateParameterNameList([]string{}); !errors.Is(err, ErrValidationException) {
		t.Fatalf("empty err = %v, want ErrValidationException", err)
	}
	names := make([]string, 10)
	if err := validateParameterNameList(names); err != nil {
		t.Fatalf("10 names rejected: %v", err)
	}
	if err := validateParameterNameList(append(names, "extra")); !errors.Is(err, ErrValidationException) {
		t.Fatalf("11 names err = %v, want ErrValidationException", err)
	}
}

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
// input is rejected instead of being silently dropped or silently reduced to
// its valid subset (a dropped filter would return unfiltered results), while
// the optional Values member may be omitted (key-existence filter).
func TestFiltersFromList_FailClosed(t *testing.T) {
	// Wire-type violations never deserialise into the modelled awsJson1_1
	// shape and map to SerializationException.
	if _, err := filtersFromList("not-a-list"); !errors.Is(err, ErrSerializationException) {
		t.Fatalf("string err = %v, want ErrSerializationException", err)
	}
	if _, err := filtersFromList([]interface{}{"not-an-object"}); !errors.Is(err, ErrSerializationException) {
		t.Fatalf("non-object entry err = %v, want ErrSerializationException", err)
	}
	nonStringKey := []interface{}{map[string]interface{}{"Key": 5, "Values": []interface{}{"p1"}}}
	if _, err := filtersFromList(nonStringKey); !errors.Is(err, ErrSerializationException) {
		t.Fatalf("non-string Key err = %v, want ErrSerializationException", err)
	}
	// The Values member is optional ("Required: No"): an entry with a Key
	// and no Values member is accepted as a key-existence filter.
	noValues := []interface{}{map[string]interface{}{"Key": "Name"}}
	filters, err := filtersFromList(noValues)
	if err != nil || len(filters) != 1 || len(filters[0].Values) != 0 {
		t.Fatalf("omitted Values err = %v, filters = %+v, want accepted key-existence filter", err, filters)
	}
	// A present-but-empty Values member violates the value-list length
	// minimum of 1 and stays rejected.
	emptyValues := []interface{}{map[string]interface{}{"Key": "Name", "Values": []interface{}{}}}
	if _, err := filtersFromList(emptyValues); !errors.Is(err, ErrInvalidFilterValue) {
		t.Fatalf("empty Values err = %v, want ErrInvalidFilterValue", err)
	}
	// Every Values entry must satisfy the ParameterStringFilterValue length
	// minimum of 1; a mixed list is invalid as a whole, not reduced to its
	// non-empty entries.
	mixedValues := []interface{}{map[string]interface{}{"Key": "Name", "Values": []interface{}{"", "p1"}}}
	if _, err := filtersFromList(mixedValues); !errors.Is(err, ErrInvalidFilterValue) {
		t.Fatalf("mixed empty/valid Values err = %v, want ErrInvalidFilterValue", err)
	}
	allEmpty := []interface{}{map[string]interface{}{"Key": "Name", "Values": []interface{}{""}}}
	if _, err := filtersFromList(allEmpty); !errors.Is(err, ErrInvalidFilterValue) {
		t.Fatalf("single empty Values entry err = %v, want ErrInvalidFilterValue", err)
	}
	nonStringEntry := []interface{}{map[string]interface{}{"Key": "Name", "Values": []interface{}{1}}}
	if _, err := filtersFromList(nonStringEntry); !errors.Is(err, ErrSerializationException) {
		t.Fatalf("non-string Values entry err = %v, want ErrSerializationException", err)
	}
	// The Option member is optional; an explicitly empty one violates the
	// ParameterStringQueryOption length minimum of 1 and a non-string one is
	// a wire-type violation — neither may be treated as an omitted member.
	emptyOption := []interface{}{map[string]interface{}{"Key": "Name", "Option": ""}}
	if _, err := filtersFromList(emptyOption); !errors.Is(err, ErrInvalidFilterOption) {
		t.Fatalf("empty Option err = %v, want ErrInvalidFilterOption", err)
	}
	nonStringOption := []interface{}{map[string]interface{}{"Key": "Name", "Option": 5}}
	if _, err := filtersFromList(nonStringOption); !errors.Is(err, ErrSerializationException) {
		t.Fatalf("non-string Option err = %v, want ErrSerializationException", err)
	}
	nullOption := []interface{}{map[string]interface{}{"Key": "Name", "Option": nil}}
	filters, err = filtersFromList(nullOption)
	if err != nil || len(filters) != 1 || filters[0].Option != "" {
		t.Fatalf("null Option err = %v, filters = %+v, want accepted as omitted", err, filters)
	}
	valid := []interface{}{map[string]interface{}{"Key": "Name", "Values": []interface{}{"p1"}}}
	filters, err = filtersFromList(valid)
	if err != nil || len(filters) != 1 {
		t.Fatalf("valid filters err = %v, len = %d", err, len(filters))
	}
	if _, err := filtersFromList(nil); err != nil {
		t.Fatalf("absent key err = %v, want nil", err)
	}
}

func TestFiltersFromQueryParams_FailClosed(t *testing.T) {
	// The Values member is optional: an entry with a Key and no Values is a
	// key-existence filter and must not be dropped nor stop later entries
	// from being parsed.
	params := map[string]interface{}{
		"ParameterFilters.member.1.Key":             "Name",
		"ParameterFilters.member.2.Key":             "Type",
		"ParameterFilters.member.2.Values.member.1": "String",
	}
	filters, err := filtersFromQueryParams(params, "ParameterFilters")
	if err != nil || len(filters) != 2 {
		t.Fatalf("key-only entry err = %v, len = %d, want both entries accepted", err, len(filters))
	}
	if len(filters[0].Values) != 0 {
		t.Fatalf("key-only entry collected values: %+v", filters[0])
	}

	valid := map[string]interface{}{
		"ParameterFilters.member.1.Key":             "Name",
		"ParameterFilters.member.1.Values.member.1": "p1",
		"ParameterFilters.member.1.Values.member.2": "p2",
		"ParameterFilters.member.2.Key":             "Type",
		"ParameterFilters.member.2.Option":          "Equals",
		"ParameterFilters.member.2.Values.member.1": "String",
	}
	filters, err = filtersFromQueryParams(valid, "ParameterFilters")
	if err != nil || len(filters) != 2 {
		t.Fatalf("valid filters err = %v, len = %d", err, len(filters))
	}
	if len(filters[0].Values) != 2 || filters[1].Option != "Equals" {
		t.Fatalf("unexpected filters: %+v", filters)
	}
}

// TestValidateAllowedPatternUnicodeLength pins that AllowedPattern
// @length(0,1024) is counted in Unicode characters; the shape carries no
// pattern trait of its own, and regex sources may contain multibyte
// character classes.
func TestValidateAllowedPatternUnicodeLength(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes; a valid regex literal

	if err := validateAllowedPattern(strings.Repeat(cjk, 1024)); err != nil {
		t.Errorf("1024-character CJK allowed pattern rejected: %v", err)
	}
	if err := validateAllowedPattern(strings.Repeat(cjk, 1025)); err == nil {
		t.Error("1025-character CJK allowed pattern accepted")
	}
	if err := validateAllowedPattern(""); err != nil {
		t.Errorf("empty allowed pattern rejected: %v", err)
	}
}

// TestValidateLabelsUnicodeLength pins that ParameterLabel @length(1,100)
// is counted in Unicode characters (the shape carries no pattern).
func TestValidateLabelsUnicodeLength(t *testing.T) {
	cjk := "\u65e5"

	if err := validateLabels([]string{strings.Repeat(cjk, 100)}); err != nil {
		t.Errorf("100-character CJK label rejected: %v", err)
	}
	if err := validateLabels([]string{strings.Repeat(cjk, 101)}); err == nil {
		t.Error("101-character CJK label accepted")
	}
}
