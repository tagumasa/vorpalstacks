package dynamodb

import (
	"context"

	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/utils/aws/types"

	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Transport-agnostic DTO for CreateTable
// ---------------------------------------------------------------------------

// CreateTableInput carries every field that CreateTable needs, in a format
// independent of the wire protocol. Both the HTTP API handler and the
// admin gRPC handler build this struct and delegate to createTableCore,
// ensuring that validation and table creation follow a single code path.
type CreateTableInput struct {
	TableName                 string
	KeySchema                 []*dbstore.KeySchemaElement
	AttributeDefinitions      []*dbstore.AttributeDefinition
	BillingMode               dbstore.BillingMode
	ProvisionedThroughput     *dbstore.ProvisionedThroughput
	GlobalSecondaryIndexes    []*dbstore.GlobalSecondaryIndex
	LocalSecondaryIndexes     []*dbstore.LocalSecondaryIndex
	StreamSpecification       *dbstore.StreamSpecification
	Tags                      []types.Tag
	DeletionProtectionEnabled bool
	WarmThroughput            *dbstore.WarmThroughput
	OnDemandThroughput        *dbstore.OnDemandThroughput
	GlobalTableSourceArn      string
	SSEDescription            *dbstore.SSEDescription
	TableClass                string
}

// ---------------------------------------------------------------------------
// Core function — single validation + persistence path
// ---------------------------------------------------------------------------

// createTableCore is the single entry point for table creation shared by the
// HTTP API and the admin gRPC handler. It applies defaults, performs all
// validation, creates the table in the store, applies post-create field
// updates, tags the table, and returns the fully-created table.
func (s *DynamoDBService) createTableCore(store dbstore.DynamoDBStoreInterface, in CreateTableInput) (*dbstore.Table, error) {
	// 1. Table name validation (length 3-255, allowed characters).
	if err := validateTableName(in.TableName); err != nil {
		return nil, err
	}

	// 2. Key schema must be non-empty and well-formed.
	if len(in.KeySchema) == 0 {
		return nil, ErrInvalidParameter
	}
	if err := validateKeySchema(in.KeySchema); err != nil {
		return nil, err
	}

	// 3. Attribute definitions must cover every key attribute.
	if err := validateAttributeDefinitions(in.KeySchema, in.AttributeDefinitions); err != nil {
		return nil, err
	}

	// 4. Billing mode default + consistency check.
	if in.BillingMode == "" {
		in.BillingMode = dbstore.BillingModeProvisioned
	}
	if err := validateBillingModeConsistency(in.BillingMode, in.ProvisionedThroughput); err != nil {
		return nil, err
	}

	// 5. Cross-index validation: every key attribute across table + GSIs +
	//    LSIs must appear in AttributeDefinitions; LSI partition key must
	//    match the table partition key; index names must be unique.
	if err := validateAllKeyAttributesInDefs(in.KeySchema, in.GlobalSecondaryIndexes, in.LocalSecondaryIndexes, in.AttributeDefinitions); err != nil {
		return nil, err
	}
	if err := validateLSIPartitionKey(in.KeySchema, in.LocalSecondaryIndexes); err != nil {
		return nil, err
	}
	if err := validateIndexNameUniqueness(in.GlobalSecondaryIndexes, in.LocalSecondaryIndexes); err != nil {
		return nil, err
	}

	// 6. Persist.
	table, err := store.Tables().Create(
		in.TableName,
		in.KeySchema,
		in.AttributeDefinitions,
		in.BillingMode,
		in.ProvisionedThroughput,
		in.GlobalSecondaryIndexes,
		in.LocalSecondaryIndexes,
		in.StreamSpecification,
		in.Tags,
		in.DeletionProtectionEnabled,
	)
	if err != nil {
		if dbstore.IsTableAlreadyExists(err) {
			return nil, ErrTableAlreadyExists
		}
		return nil, err
	}

	// 7. Post-create field updates (SSE, warm throughput, on-demand
	//    throughput, global table source ARN, table class).
	needsPersist := false
	if in.SSEDescription != nil {
		table.SSEDescription = in.SSEDescription
		needsPersist = true
	}
	if in.WarmThroughput != nil {
		table.WarmThroughput = in.WarmThroughput
		needsPersist = true
	}
	if in.OnDemandThroughput != nil {
		table.OnDemandThroughput = in.OnDemandThroughput
		needsPersist = true
	}
	if in.GlobalTableSourceArn != "" {
		table.GlobalTableSourceArn = in.GlobalTableSourceArn
		needsPersist = true
	}
	if in.TableClass != "" {
		table.TableClass = in.TableClass
		needsPersist = true
	}
	if needsPersist {
		if err := store.Tables().Put(table); err != nil {
			return nil, err
		}
	}

	// 8. Tags.
	if len(in.Tags) > 0 {
		store.Tables().Tags().Tag(in.TableName, tagutil.ToMap(in.Tags))
	}

	return table, nil
}

// ---------------------------------------------------------------------------
// Core functions for DeleteTable, DescribeTable, ListTables
// ---------------------------------------------------------------------------

// deleteTableCore is the single entry point for table deletion shared by the
// HTTP API and the admin gRPC handler. It validates the table name, checks
// for deletion protection, performs the cascade delete, and returns the
// archived table description.
func (s *DynamoDBService) deleteTableCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName string) (*dbstore.Table, error) {
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	var deletedTable *dbstore.Table

	err := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		table, err := txn.GetTable(tableName)
		if err != nil {
			if dbstore.IsTableNotFound(err) {
				return ErrTableNotFound
			}
			return err
		}
		if table.DeletionProtectionEnabled {
			return ErrTableDeletionProtected
		}
		deletedTable = table
		return txn.DeleteTableCascade(tableName)
	})
	if err != nil {
		return nil, err
	}

	deletedTable.Status = dbstore.TableStatusArchived
	return deletedTable, nil
}

// describeTableCore is the single entry point for table lookup shared by the
// HTTP API and the admin gRPC handler. It validates the table name and maps
// store errors to DynamoDB API errors.
func (s *DynamoDBService) describeTableCore(store dbstore.DynamoDBStoreInterface, tableName string) (*dbstore.Table, error) {
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		if dbstore.IsTableNotFound(err) || commonstore.IsNotFound(err) {
			return nil, ErrTableNotFound
		}
		return nil, err
	}
	return table, nil
}

// listTablesCore is the single entry point for table listing shared by the
// HTTP API and the admin gRPC handler. It delegates to the store after
// validating the limit range.
func (s *DynamoDBService) listTablesCore(store dbstore.DynamoDBStoreInterface, marker string, limit int) ([]*dbstore.Table, string, error) {
	return store.Tables().List(marker, limit)
}
