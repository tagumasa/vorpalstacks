// Package dynamodb provides DynamoDB service operations for vorpalstacks.
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
	assert.Equal(t, 0, dbstore.CountSignificantDigits("0"))
	assert.Equal(t, 0, dbstore.CountSignificantDigits("0.000"))
	assert.Equal(t, 1, dbstore.CountSignificantDigits("1000"))
	assert.Equal(t, 3, dbstore.CountSignificantDigits("1.23"))
	assert.Equal(t, 1, dbstore.CountSignificantDigits("0.0001"))
	assert.Equal(t, 2, dbstore.CountSignificantDigits("-2.5e10"))
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
