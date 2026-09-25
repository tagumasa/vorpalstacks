package dynamodb

import (
	"testing"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// Binary sets are unordered: equality under = and IN is membership-based,
// the contract the SS and NS branches of attributeValuesEqual already
// apply — not element-by-element position comparison, which would reject
// a reordered rendering of the same set.
func TestBinarySetEqualityIsMembershipBased(t *testing.T) {
	low := []byte{0x01, 0x02}
	high := []byte{0x03, 0x04}
	forward := &dbstore.AttributeValue{BS: [][]byte{low, high}}
	reversed := &dbstore.AttributeValue{BS: [][]byte{high, low}}

	if !attributeValuesEqual(forward, reversed) {
		t.Fatal("reordered binary set: membership equality failed")
	}
	if !compareAttributeValues(forward, "=", reversed) {
		t.Fatal("reordered binary set under '=': must compare equal")
	}
	if compareAttributeValues(forward, "<>", reversed) {
		t.Fatal("reordered binary set under '<>': must not report a difference")
	}

	// A set carrying a different element is unequal in either order.
	changed := &dbstore.AttributeValue{BS: [][]byte{low, []byte{0x03, 0x05}}}
	if attributeValuesEqual(forward, changed) {
		t.Fatal("different binary sets must not compare equal")
	}
	if attributeValuesEqual(reversed, changed) {
		t.Fatal("different binary sets must not compare equal in reverse order either")
	}
}
