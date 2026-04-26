package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestDeepCopyTable(t *testing.T) {
	t.Run("nil returns nil", func(t *testing.T) {
		assert.Nil(t, deepCopyTable(nil))
	})

	t.Run("shallow fields copied", func(t *testing.T) {
		orig := &dbstore.Table{
			Name:      "test",
			ARN:       "arn:test",
			Status:    dbstore.TableStatusActive,
			ItemCount: 42,
		}
		cp := deepCopyTable(orig)
		assert.Equal(t, orig.Name, cp.Name)
		assert.Equal(t, orig.ARN, cp.ARN)
		assert.Equal(t, orig.Status, cp.Status)
		assert.Equal(t, orig.ItemCount, cp.ItemCount)
	})

	t.Run("KeySchema deep copy isolates mutations (Bug D-1)", func(t *testing.T) {
		orig := &dbstore.Table{
			Name: "test",
			KeySchema: []*dbstore.KeySchemaElement{
				{AttributeName: "id", KeyType: dbstore.KeyTypeHash},
			},
		}
		cp := deepCopyTable(orig)

		cp.KeySchema[0].AttributeName = "modified"
		assert.Equal(t, "id", orig.KeySchema[0].AttributeName,
			"mutating copy should not affect original")
	})

	t.Run("AttributeDefinitions deep copy isolates mutations", func(t *testing.T) {
		orig := &dbstore.Table{
			Name: "test",
			AttributeDefinitions: []*dbstore.AttributeDefinition{
				{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS},
			},
		}
		cp := deepCopyTable(orig)

		cp.AttributeDefinitions[0].AttributeName = "modified"
		assert.Equal(t, "id", orig.AttributeDefinitions[0].AttributeName)
	})

	t.Run("ProvisionedThroughput deep copy", func(t *testing.T) {
		orig := &dbstore.Table{
			Name: "test",
			ProvisionedThroughput: &dbstore.ProvisionedThroughput{
				ReadCapacityUnits:  10,
				WriteCapacityUnits: 5,
			},
		}
		cp := deepCopyTable(orig)

		cp.ProvisionedThroughput.ReadCapacityUnits = 999
		assert.Equal(t, int64(10), orig.ProvisionedThroughput.ReadCapacityUnits)
	})

	t.Run("GSI deep copy with Projection NonKeyAttributes", func(t *testing.T) {
		orig := &dbstore.Table{
			Name: "test",
			GlobalSecondaryIndexes: []*dbstore.GlobalSecondaryIndex{
				{
					IndexName: "gsi1",
					KeySchema: []*dbstore.KeySchemaElement{
						{AttributeName: "gsi_pk", KeyType: dbstore.KeyTypeHash},
					},
					Projection: &dbstore.Projection{
						ProjectionType:   "ALL",
						NonKeyAttributes: []string{"attr1", "attr2"},
					},
				},
			},
		}
		cp := deepCopyTable(orig)

		cp.GlobalSecondaryIndexes[0].Projection.NonKeyAttributes[0] = "modified"
		assert.Equal(t, "attr1", orig.GlobalSecondaryIndexes[0].Projection.NonKeyAttributes[0])

		cp.GlobalSecondaryIndexes[0].KeySchema[0].AttributeName = "mod"
		assert.Equal(t, "gsi_pk", orig.GlobalSecondaryIndexes[0].KeySchema[0].AttributeName)
	})

	t.Run("LSI deep copy", func(t *testing.T) {
		orig := &dbstore.Table{
			Name: "test",
			LocalSecondaryIndexes: []*dbstore.LocalSecondaryIndex{
				{
					IndexName: "lsi1",
					KeySchema: []*dbstore.KeySchemaElement{
						{AttributeName: "lsi_sk", KeyType: dbstore.KeyTypeRange},
					},
				},
			},
		}
		cp := deepCopyTable(orig)

		cp.LocalSecondaryIndexes[0].KeySchema[0].AttributeName = "mod"
		assert.Equal(t, "lsi_sk", orig.LocalSecondaryIndexes[0].KeySchema[0].AttributeName)
	})

	t.Run("StreamSpecification deep copy", func(t *testing.T) {
		orig := &dbstore.Table{
			Name: "test",
			StreamSpecification: &dbstore.StreamSpecification{
				StreamEnabled:  true,
				StreamViewType: dbstore.StreamViewTypeNewAndOldImages,
			},
		}
		cp := deepCopyTable(orig)

		cp.StreamSpecification.StreamEnabled = false
		assert.True(t, orig.StreamSpecification.StreamEnabled)
	})

	t.Run("TTL deep copy", func(t *testing.T) {
		orig := &dbstore.Table{
			Name: "test",
			TimeToLive: &dbstore.TimeToLiveSpecification{
				Enabled:       true,
				AttributeName: "expires",
			},
		}
		cp := deepCopyTable(orig)

		cp.TimeToLive.Enabled = false
		assert.True(t, orig.TimeToLive.Enabled)
	})
}
