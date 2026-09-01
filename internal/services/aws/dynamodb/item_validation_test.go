package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestIsValidDynamoDBNumber(t *testing.T) {
	valid := []string{
		"0",
		"-1.5",
		"+3",
		"1e10",
		"9.9999999999999999999999999999999999999E+125",
		"1E-130",
		"0.000000000000000000000000000000000000000001",
		"10000000000000000000000000000000000000000",
	}
	for _, value := range valid {
		assert.True(t, isValidDynamoDBNumber(value), "expected %q to be valid", value)
	}

	invalid := []string{
		"",
		"abc",
		"1/3",     // big.Rat fraction form, not part of the grammar
		"1e400",   // magnitude above 9.99...E+125
		"-1E-131", // magnitude below 1E-130
		"1.23456789012345678901234567890123456789", // 39 significant digits
		"1.2.3",
		"--1",
	}
	for _, value := range invalid {
		assert.False(t, isValidDynamoDBNumber(value), "expected %q to be invalid", value)
	}
}

func TestCountSignificantDigits(t *testing.T) {
	assert.Equal(t, 0, countSignificantDigits("0"))
	assert.Equal(t, 0, countSignificantDigits("0.000"))
	assert.Equal(t, 1, countSignificantDigits("1000"))
	assert.Equal(t, 3, countSignificantDigits("1.23"))
	assert.Equal(t, 1, countSignificantDigits("0.0001"))
	assert.Equal(t, 2, countSignificantDigits("-2.5e10"))
}

func TestSetDuplicateDetection(t *testing.T) {
	assert.True(t, hasDuplicateString([]string{"a", "b", "a"}))
	assert.False(t, hasDuplicateString([]string{"a", "b", ""}))
	assert.True(t, hasDuplicateNumber([]string{"1", "2", "1.0"}))
	assert.False(t, hasDuplicateNumber([]string{"1", "2", "3"}))
	assert.True(t, hasDuplicateBinary([][]byte{[]byte("x"), []byte("x")}))
	assert.False(t, hasDuplicateBinary([][]byte{[]byte("x"), []byte("")}))
}

func TestResolveNameStrict(t *testing.T) {
	resolved, err := resolveNameStrict("#n", map[string]string{"#n": "name"})
	require.NoError(t, err)
	assert.Equal(t, "name", resolved)

	resolved, err = resolveNameStrict("plain", nil)
	require.NoError(t, err)
	assert.Equal(t, "plain", resolved)

	_, err = resolveNameStrict("#missing", map[string]string{"#n": "name"})
	assert.Error(t, err)

	_, err = resolveNameStrict("#missing", nil)
	assert.Error(t, err)
}

func TestUpdateExpressionStrictness(t *testing.T) {
	values := map[string]*dbstore.AttributeValue{
		":v": dbstore.NumberValue("2"),
		":w": dbstore.NumberValue("3"),
	}

	attrs := map[string]*dbstore.AttributeValue{
		"src": dbstore.StringValue("copied"),
	}

	// The multiply operator is not part of SET arithmetic.
	_, err := applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET a = :v * :w", nil, values)
	assert.Error(t, err)

	// An undefined value placeholder must not silently skip the action.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET a = :missing", nil, values)
	assert.Error(t, err)

	// An undefined name placeholder must not pass through as a literal.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET #missing = :v", nil, values)
	assert.Error(t, err)

	// SET a = src copies the existing attribute value.
	updated, err := applyUpdateExpressionWithTracking(attrs, "SET a = src", nil, nil)
	require.NoError(t, err)
	assert.Contains(t, updated, "a")
	assert.Equal(t, "copied", *attrs["a"].S)

	// DELETE with a mismatched operand type is rejected.
	ssAttrs := map[string]*dbstore.AttributeValue{
		"tags": dbstore.StringSet([]string{"a", "b"}),
	}
	_, err = applyUpdateExpressionWithTracking(ssAttrs, "DELETE tags :nums", nil, map[string]*dbstore.AttributeValue{
		":nums": dbstore.NumberSet([]string{"1"}),
	})
	assert.Error(t, err)

	// DELETE of matching elements still works.
	updated, err = applyUpdateExpressionWithTracking(ssAttrs, "DELETE tags :del", nil, map[string]*dbstore.AttributeValue{
		":del": dbstore.StringSet([]string{"a"}),
	})
	require.NoError(t, err)
	assert.Contains(t, updated, "tags")
	assert.Equal(t, []string{"b"}, ssAttrs["tags"].SS)

	// An undefined name placeholder in REMOVE must not pass through as a
	// literal attribute name.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "REMOVE #missing", nil, nil)
	assert.Error(t, err)

	// REMOVE with a defined name placeholder removes the attribute.
	rmAttrs := map[string]*dbstore.AttributeValue{
		"gsik": dbstore.StringValue("g"),
	}
	updated, err = applyUpdateExpressionWithTracking(rmAttrs, "REMOVE #name", map[string]string{"#name": "gsik"}, nil)
	require.NoError(t, err)
	assert.Contains(t, updated, "gsik")
	assert.NotContains(t, rmAttrs, "gsik")
}

// TestAdminItemCoresRejectEmptyTableName pins the empty-table rejection
// in the admin item cores: an omitted TableName is a client error
// (ValidationException) reported before any store or table lookup, so
// the console does not surface a table-not-found for it.
func TestAdminItemCoresRejectEmptyTableName(t *testing.T) {
	svc := &DynamoDBService{}

	if _, err := svc.adminGetItem(nil, "us-east-1", "", nil); err != ErrInvalidParameter {
		t.Errorf("adminGetItem: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.adminScan("us-east-1", "", 10, nil); err != ErrInvalidParameter {
		t.Errorf("adminScan: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.adminPutItem(nil, "us-east-1", "", nil); err != ErrInvalidParameter {
		t.Errorf("adminPutItem: expected ErrInvalidParameter, got %v", err)
	}
	if err := svc.adminDeleteItem(nil, "us-east-1", "", nil); err != ErrInvalidParameter {
		t.Errorf("adminDeleteItem: expected ErrInvalidParameter, got %v", err)
	}
}
