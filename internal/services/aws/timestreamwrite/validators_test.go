package timestreamwrite

import (
	"strings"
	"testing"

	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// These tests pin the enum and S3 shape validators to the Smithy model so
// the validators cannot silently lose their call sites again: the record
// enum contract is exercised through validateRecordEnums (the WriteRecords
// post-merge entry point) and the S3 shape validators through direct
// boundary cases.

func TestValidateRecordEnums(t *testing.T) {
	validRecord := tsstore.Record{
		MeasureName:      "cpu",
		MeasureValue:     "1.5",
		MeasureValueType: tsstore.MeasureValueTypeDouble,
		Time:             "1700000000000",
		TimeUnit:         tsstore.TimeUnitMilliseconds,
		Dimensions: []tsstore.Dimension{
			{Name: "host", Value: "a1", DimensionValueType: tsstore.DimensionValueTypeVarchar},
		},
	}

	t.Run("valid record passes", func(t *testing.T) {
		if err := validateRecordEnums([]tsstore.Record{validRecord}); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("invalid MeasureValueType rejected", func(t *testing.T) {
		rec := validRecord
		rec.MeasureValueType = tsstore.MeasureValueType("FOO")
		if err := validateRecordEnums([]tsstore.Record{rec}); err == nil {
			t.Fatal("expected error for invalid MeasureValueType")
		}
	})

	t.Run("invalid TimeUnit rejected", func(t *testing.T) {
		rec := validRecord
		rec.TimeUnit = tsstore.TimeUnit("CENTURIES")
		if err := validateRecordEnums([]tsstore.Record{rec}); err == nil {
			t.Fatal("expected error for invalid TimeUnit")
		}
	})

	t.Run("invalid DimensionValueType rejected", func(t *testing.T) {
		rec := validRecord
		rec.Dimensions = []tsstore.Dimension{{Name: "host", Value: "a1", DimensionValueType: tsstore.DimensionValueType("INT")}}
		if err := validateRecordEnums([]tsstore.Record{rec}); err == nil {
			t.Fatal("expected error for invalid DimensionValueType")
		}
	})

	t.Run("MeasureValues without Type rejected", func(t *testing.T) {
		rec := validRecord
		rec.MeasureValueType = tsstore.MeasureValueTypeMulti
		rec.MeasureValues = []tsstore.MeasureValue{{Name: "m1", Value: "1"}}
		if err := validateRecordEnums([]tsstore.Record{rec}); err == nil {
			t.Fatal("expected error for missing MeasureValue.Type")
		}
	})

	t.Run("MeasureValues with invalid Type rejected", func(t *testing.T) {
		rec := validRecord
		rec.MeasureValueType = tsstore.MeasureValueTypeMulti
		rec.MeasureValues = []tsstore.MeasureValue{{Name: "m1", Value: "1", Type: tsstore.MeasureValueType("DECIMAL")}}
		if err := validateRecordEnums([]tsstore.Record{rec}); err == nil {
			t.Fatal("expected error for invalid MeasureValue.Type")
		}
	})

	t.Run("MULTI without MeasureValues rejected", func(t *testing.T) {
		rec := validRecord
		rec.MeasureValueType = tsstore.MeasureValueTypeMulti
		if err := validateRecordEnums([]tsstore.Record{rec}); err == nil {
			t.Fatal("expected error for MULTI without MeasureValues")
		}
	})

	t.Run("MeasureValues with scalar type rejected", func(t *testing.T) {
		rec := validRecord
		rec.MeasureValues = []tsstore.MeasureValue{{Name: "m1", Value: "1", Type: tsstore.MeasureValueTypeDouble}}
		if err := validateRecordEnums([]tsstore.Record{rec}); err == nil {
			t.Fatal("expected error for MeasureValues with scalar MeasureValueType")
		}
	})

	t.Run("MULTI with typed MeasureValues passes", func(t *testing.T) {
		rec := validRecord
		rec.MeasureValueType = tsstore.MeasureValueTypeMulti
		rec.MeasureValues = []tsstore.MeasureValue{{Name: "m1", Value: "1", Type: tsstore.MeasureValueTypeDouble}}
		if err := validateRecordEnums([]tsstore.Record{rec}); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})
}

func TestValidateS3ObjectKeyPrefix(t *testing.T) {
	patternValid := func(n int) string {
		return "a" + strings.Repeat("b", n-1)
	}
	cases := []struct {
		name   string
		prefix string
		want   bool
	}{
		{"short prefix valid", "data/", true},
		{"896 chars valid (old upper bound)", patternValid(896), true},
		{"928 chars valid (Smithy upper bound)", patternValid(928), true},
		{"929 chars rejected", patternValid(929), false},
		{"space rejected by pattern", "pre fix/", false},
		{"backslash rejected by pattern", "pre\\fix", false},
		{"empty rejected", "", false},
	}
	for _, tc := range cases {
		if got := validateS3ObjectKeyPrefix(tc.prefix); got != tc.want {
			t.Errorf("%s: validateS3ObjectKeyPrefix(len=%d) = %v, want %v", tc.name, len(tc.prefix), got, tc.want)
		}
	}
}

func TestValidateS3ObjectKey(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"data.csv", true},
		{"ab", true},
		{"a", false}, // pattern requires at least 2 chars
		{"path/to/file.csv", true},
		{"file name.csv", false},            // space not in pattern
		{"a\\b.csv", false},                 // backslash not in the model pattern
		{"", false},                         // below min length 1
		{string(make([]byte, 1025)), false}, // above max length 1024
	}
	for _, tc := range cases {
		if got := validateS3ObjectKey(tc.key); got != tc.want {
			t.Errorf("validateS3ObjectKey(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

// TestValidateKmsKeyIdAndClientTokenUnicodeLengths pins that the
// S3Configuration KmsKeyId shape (StringValue2048) @length(1,2048) and the
// ClientToken parameter are counted in Unicode characters; neither shape
// carries a pattern.
func TestValidateKmsKeyIdAndClientTokenUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if !validateKmsKeyId(strings.Repeat(cjk, 2048)) {
		t.Error("2048-character CJK KmsKeyId rejected")
	}
	if validateKmsKeyId(strings.Repeat(cjk, 2049)) {
		t.Error("2049-character CJK KmsKeyId accepted")
	}
	if validateKmsKeyId("") {
		t.Error("empty KmsKeyId accepted")
	}
	if !validateClientToken(strings.Repeat(cjk, 64)) {
		t.Error("64-character CJK client token rejected")
	}
	if validateClientToken(strings.Repeat(cjk, 65)) {
		t.Error("65-character CJK client token accepted")
	}
}
