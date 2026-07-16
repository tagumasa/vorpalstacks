package dynamodb

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestExtractPrimaryKeyCondition_BinaryKey(t *testing.T) {
	table := &dbstore.Table{
		Name: "test",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "id", KeyType: dbstore.KeyTypeHash},
		},
	}

	rawBytes := []byte{0x00, 0x01, 0x02, 0x03}
	values := map[string]*dbstore.AttributeValue{
		":v": {B: rawBytes},
	}

	t.Run("binary key encoded as base64 (Bug D-4)", func(t *testing.T) {
		hashKey, _ := extractPrimaryKeyCondition(table, "id = :v", nil, values)
		expected := base64.StdEncoding.EncodeToString(rawBytes)
		assert.Equal(t, expected, hashKey,
			"binary key should be base64-encoded, not raw string()")
	})

	t.Run("string key unchanged", func(t *testing.T) {
		strValues := map[string]*dbstore.AttributeValue{
			":v": {S: ptrStr("my-key")},
		}
		hashKey, _ := extractPrimaryKeyCondition(table, "id = :v", nil, strValues)
		assert.Equal(t, "my-key", hashKey)
	})

	t.Run("numeric key", func(t *testing.T) {
		numValues := map[string]*dbstore.AttributeValue{
			":v": {N: ptrStr("42")},
		}
		hashKey, _ := extractPrimaryKeyCondition(table, "id = :v", nil, numValues)
		assert.Equal(t, "42", hashKey)
	})
}

func TestExtractPrimaryKeyCondition_WithSortKey(t *testing.T) {
	table := &dbstore.Table{
		Name: "test",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
	}

	values := map[string]*dbstore.AttributeValue{
		":pk": {S: ptrStr("user1")},
		":sk": {S: ptrStr("order1")},
	}

	t.Run("hash + sort key condition", func(t *testing.T) {
		hashKey, sortCond := extractPrimaryKeyCondition(table, "pk = :pk AND sk = :sk", nil, values)
		assert.Equal(t, "user1", hashKey)
		assert.NotNil(t, sortCond)
		assert.Equal(t, "=", sortCond.op)
	})
}
