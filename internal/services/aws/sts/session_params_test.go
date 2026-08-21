package sts

import (
	"strings"
	"testing"
)

// TestExtractSessionTagsUnicodeLengths pins that session tag key and value
// limits are counted in Unicode characters, matching the Smithy tagKeyType
// and tagValueType @length traits.
func TestExtractSessionTagsUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	params := map[string]interface{}{
		"Tags.member.1.Key":   strings.Repeat(cjk, 128),
		"Tags.member.1.Value": strings.Repeat(cjk, 256),
	}
	tags, err := extractSessionTags(params)
	if err != nil || len(tags) != 1 {
		t.Errorf("128-character key with 256-character value = %v, %v; want accepted", tags, err)
	}

	params = map[string]interface{}{
		"Tags.member.1.Key":   strings.Repeat(cjk, 129),
		"Tags.member.1.Value": "v",
	}
	if _, err := extractSessionTags(params); err != ErrInvalidSessionTag {
		t.Errorf("129-character CJK key = %v, want ErrInvalidSessionTag", err)
	}

	params = map[string]interface{}{
		"Tags.member.1.Key":   "k",
		"Tags.member.1.Value": strings.Repeat(cjk, 257),
	}
	if _, err := extractSessionTags(params); err != ErrInvalidSessionTag {
		t.Errorf("257-character CJK value = %v, want ErrInvalidSessionTag", err)
	}
}
