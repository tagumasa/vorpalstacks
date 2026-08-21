package tags

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestCheckTagsStandardLimits(t *testing.T) {
	limits := StandardLimits()

	valid := [][]Tag{
		nil,
		{},
		{{Key: "env", Value: "prod"}},
		{{Key: strings128(), Value: strings256()}},
	}
	for _, tags := range valid {
		if v, _ := CheckTags(tags, limits); v != OK {
			t.Errorf("CheckTags(%v) = %v, want OK", tags, v)
		}
	}

	invalid := []struct {
		tags []Tag
		want Violation
	}{
		{tags: tooManyTags(51), want: TooManyTags},
		{tags: []Tag{{Key: "", Value: "v"}}, want: TagKeyTooShort},
		{tags: []Tag{{Key: strings129(), Value: "v"}}, want: TagKeyTooLong},
		{tags: []Tag{{Key: "env", Value: strings257()}}, want: TagValueTooLong},
		{tags: []Tag{{Key: "aws:reserved", Value: "v"}}, want: ReservedTagKey},
		{tags: []Tag{{Key: "AWS:reserved", Value: "v"}}, want: ReservedTagKey},
	}
	for _, tc := range invalid {
		if v, _ := CheckTags(tc.tags, limits); v != tc.want {
			t.Errorf("CheckTags(key=%q) = %v, want %v", firstKey(tc.tags), v, tc.want)
		}
	}
}

func TestCheckTagsProfileVariants(t *testing.T) {
	// A service with no count limit, a 200-key cap and no reservation.
	limits := TagLimits{MaxKeyLength: 200}
	tags := tooManyTags(60)
	if v, _ := CheckTags(tags, limits); v != OK {
		t.Errorf("CheckTags with MaxCount=0 = %v, want OK", v)
	}
	if v, _ := CheckTags([]Tag{{Key: "aws:anything", Value: "v"}}, limits); v != OK {
		t.Errorf("CheckTags without ReservedPrefix = %v, want OK", v)
	}

	// A service with a case-sensitive reservation.
	limits = TagLimits{ReservedPrefix: "aws:", ReservedCaseSensitive: true}
	if v, _ := CheckTags([]Tag{{Key: "AWS:reserved", Value: "v"}}, limits); v != OK {
		t.Errorf("case-sensitive reservation rejected mixed case: %v", v)
	}
	if v, _ := CheckTags([]Tag{{Key: "aws:reserved", Value: "v"}}, limits); v != ReservedTagKey {
		t.Errorf("case-sensitive reservation missed lowercase: %v", v)
	}
}

func TestCheckStringTagsMatchesCheckTags(t *testing.T) {
	limits := StandardLimits()
	mapForm := map[string]string{"aws:reserved": "v"}
	if v, _ := CheckStringTags(mapForm, limits); v != ReservedTagKey {
		t.Errorf("CheckStringTags = %v, want ReservedTagKey", v)
	}
	if v, _ := CheckStringTags(map[string]string{"env": "prod"}, limits); v != OK {
		t.Errorf("CheckStringTags = %v, want OK", v)
	}
	if v, _ := CheckStringTags(make(map[string]string, 0), limits); v != OK {
		t.Errorf("CheckStringTags empty = %v, want OK", v)
	}
}

func TestValidateTagsWrapperSentinels(t *testing.T) {
	if err := ValidateTags([]Tag{{Key: "env", Value: "prod"}}); !errors.Is(err, nil) && err != nil {
		t.Errorf("ValidateTags valid = %v, want nil", err)
	}
	if err := ValidateTags(tooManyTags(51)); !errors.Is(err, ErrTooManyTags) {
		t.Errorf("ValidateTags too many = %v, want ErrTooManyTags", err)
	}
	if err := ValidateTags([]Tag{{Key: "aws:x", Value: "v"}}); !errors.Is(err, ErrReservedTagKey) {
		t.Errorf("ValidateTags reserved = %v, want ErrReservedTagKey", err)
	}
}

func tooManyTags(n int) []Tag {
	tags := make([]Tag, 0, n)
	for i := 0; i < n; i++ {
		tags = append(tags, Tag{Key: "k" + strconv.Itoa(i), Value: "v"})
	}
	return tags
}

func firstKey(tags []Tag) string {
	if len(tags) == 0 {
		return ""
	}
	return tags[0].Key
}

func strings128() string { return repeat('k', 128) }
func strings129() string { return repeat('k', 129) }
func strings256() string { return repeat('v', 256) }
func strings257() string { return repeat('v', 257) }

func repeat(c rune, n int) string {
	out := make([]rune, n)
	for i := range out {
		out[i] = c
	}
	return string(out)
}

// TestCheckTagsUnicodeLengths pins that key and value limits are counted
// in Unicode characters, not bytes: AWS describes the limits as
// "128/256 Unicode characters in UTF-8".
func TestCheckTagsUnicodeLengths(t *testing.T) {
	limits := StandardLimits()
	cjk := strings.Repeat("\u65e5", 1) // one CJK character, 3 bytes

	key128 := strings.Repeat(cjk, 128)
	if v, _ := CheckTags([]Tag{{Key: key128, Value: "v"}}, limits); v != OK {
		t.Errorf("128-character CJK key rejected: %v", v)
	}
	key129 := strings.Repeat(cjk, 129)
	if v, _ := CheckTags([]Tag{{Key: key129, Value: "v"}}, limits); v != TagKeyTooLong {
		t.Errorf("129-character CJK key = %v, want TagKeyTooLong", v)
	}
	value256 := strings.Repeat(cjk, 256)
	if v, _ := CheckTags([]Tag{{Key: "env", Value: value256}}, limits); v != OK {
		t.Errorf("256-character CJK value rejected: %v", v)
	}
	value257 := strings.Repeat(cjk, 257)
	if v, _ := CheckTags([]Tag{{Key: "env", Value: value257}}, limits); v != TagValueTooLong {
		t.Errorf("257-character CJK value = %v, want TagValueTooLong", v)
	}
}

// TestCheckTagsReportsKeyViolationsFirst pins the deterministic
// reporting order: when a tag both uses a reserved prefix and carries an
// over-long value, the key-side violation is reported — matching the
// error codes AWS services emit for the aws: reservation.
func TestCheckTagsReportsKeyViolationsFirst(t *testing.T) {
	limits := StandardLimits()
	if v, _ := CheckTags([]Tag{{Key: "aws:reserved", Value: strings257()}}, limits); v != ReservedTagKey {
		t.Errorf("CheckTags reserved+long value = %v, want ReservedTagKey", v)
	}

	// Map iteration must not change the outcome either.
	tags := map[string]string{"aws:reserved": strings257()}
	for i := 0; i < 100; i++ {
		if v, _ := CheckStringTags(tags, limits); v != ReservedTagKey {
			t.Fatalf("CheckStringTags reserved+long value = %v, want ReservedTagKey", v)
		}
	}
}

// TestCheckStringTagsReturnsOffendingKey pins that the checker reports
// WHICH key violated the limits, selected deterministically by sorted key
// order with key-side checks before value-side ones, so consumers format
// error messages from the finding instead of re-deriving it.
func TestCheckStringTagsReturnsOffendingKey(t *testing.T) {
	limits := StandardLimits()
	cjk := "\u65e5" // one CJK character, 3 bytes

	// A rune-legal multibyte key (100 characters, 300 bytes) must never be
	// reported when the actual violation is the 129-character ASCII key.
	tags := map[string]string{
		strings.Repeat("a", 129): "v",
		"zzz":                    "v",
		strings.Repeat(cjk, 100): "v",
	}
	for i := 0; i < 100; i++ {
		v, key := CheckStringTags(tags, limits)
		if v != TagKeyTooLong || key != strings.Repeat("a", 129) {
			t.Fatalf("CheckStringTags = (%v, %q), want (TagKeyTooLong, the 129-character key)", v, key)
		}
	}

	// Value violations report the offending key as well.
	tags = map[string]string{"a": strings.Repeat("v", 257), "b": "ok"}
	for i := 0; i < 100; i++ {
		v, key := CheckStringTags(tags, limits)
		if v != TagValueTooLong || key != "a" {
			t.Fatalf("CheckStringTags = (%v, %q), want (TagValueTooLong, \"a\")", v, key)
		}
	}

	// No violation carries no key.
	if v, key := CheckStringTags(map[string]string{"a": "b"}, limits); v != OK || key != "" {
		t.Errorf("CheckStringTags OK = (%v, %q), want (OK, \"\")", v, key)
	}

	// The slice form reports the first offending tag in list order.
	v, key := CheckTags([]Tag{
		{Key: strings.Repeat(cjk, 100), Value: "v"},
		{Key: strings.Repeat("a", 129), Value: "v"},
	}, limits)
	if v != TagKeyTooLong || key != strings.Repeat("a", 129) {
		t.Errorf("CheckTags = (%v, %q), want (TagKeyTooLong, the 129-character key)", v, key)
	}
}
