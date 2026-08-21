package acm

import (
	"strings"
	"testing"

	types "vorpalstacks/internal/common/tags"
)

// TestValidateACMTagsUnicodeLengths pins that tag key and value limits are
// counted in Unicode characters per the Smithy TagKey/TagValue @length
// traits, and that the reported length is the character count, not the
// byte count.
func TestValidateACMTagsUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateACMTags([]types.Tag{{Key: strings.Repeat(cjk, 128), Value: strings.Repeat(cjk, 256)}}); err != nil {
		t.Errorf("128-character key with 256-character value rejected: %v", err)
	}
	err := validateACMTags([]types.Tag{{Key: strings.Repeat(cjk, 129), Value: "v"}})
	if err == nil || !strings.Contains(err.Error(), "got 129") {
		t.Errorf("129-character CJK key = %v, want message containing \"got 129\"", err)
	}
	err = validateACMTags([]types.Tag{{Key: "k", Value: strings.Repeat(cjk, 257)}})
	if err == nil || !strings.Contains(err.Error(), "got 257") {
		t.Errorf("257-character CJK value = %v, want message containing \"got 257\"", err)
	}
}
