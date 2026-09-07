package dynamodb

import (
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

	// The returned hash value is the storage partition prefix, so it must be
	// the store's key encoding — the same rendering the item and index keys
	// are built with. Any other rendering (raw number text, base64 of the
	// binary) matches no stored key of that type.
	t.Run("binary key uses the store key encoding", func(t *testing.T) {
		hashKey, _, _ := extractPrimaryKeyCondition(table, "id = :v", nil, values)
		expected := dbstore.EncodeKeyValue(&dbstore.AttributeValue{B: rawBytes})
		assert.NotEmpty(t, expected)
		assert.Equal(t, expected, hashKey,
			"binary key must render as the store key encoding, not raw string()")
	})

	t.Run("string key uses the store key encoding", func(t *testing.T) {
		strValues := map[string]*dbstore.AttributeValue{
			":v": {S: ptrStr("my-key")},
		}
		hashKey, _, _ := extractPrimaryKeyCondition(table, "id = :v", nil, strValues)
		assert.Equal(t, dbstore.EncodeKeyValue(&dbstore.AttributeValue{S: ptrStr("my-key")}), hashKey)
	})

	t.Run("numeric key uses the store key encoding", func(t *testing.T) {
		numValues := map[string]*dbstore.AttributeValue{
			":v": {N: ptrStr("42")},
		}
		hashKey, _, _ := extractPrimaryKeyCondition(table, "id = :v", nil, numValues)
		expected := dbstore.EncodeKeyValue(&dbstore.AttributeValue{N: ptrStr("42")})
		assert.NotEqual(t, "42", expected,
			"a raw number string cannot match the encoded storage prefix")
		assert.Equal(t, expected, hashKey)
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
		hashKey, _, sortCond := extractPrimaryKeyCondition(table, "pk = :pk AND sk = :sk", nil, values)
		assert.Equal(t, dbstore.EncodeKeyValue(values[":pk"]), hashKey)
		assert.NotNil(t, sortCond)
		assert.Equal(t, "=", sortCond.op)
	})
}
