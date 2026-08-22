package kinesis

import (
	"strings"
	"testing"
)

// TestValidatePartitionKeyUnicodeLengths pins that PartitionKey follows the
// Smithy @length(1, 256) trait counted in Unicode characters; the shape
// carries no pattern, so multibyte partition keys are valid input and must
// not be rejected on byte length.
func TestValidatePartitionKeyUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if !validatePartitionKey(strings.Repeat(cjk, 256)) {
		t.Error("256-character CJK partition key rejected")
	}
	if validatePartitionKey(strings.Repeat(cjk, 257)) {
		t.Error("257-character CJK partition key accepted")
	}
	if validatePartitionKey("") {
		t.Error("empty partition key accepted")
	}
}
