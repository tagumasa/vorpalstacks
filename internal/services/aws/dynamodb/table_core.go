package dynamodb

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"

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
	VectorIndexes             []*dbstore.VectorIndex
	StreamSpecification       *dbstore.StreamSpecification
	Tags                      []tagutil.Tag
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

// validateAndGetTableWithErr behaves like validateAndGetTable but lets the
// caller pick the not-found error sentinel. Operations whose Smithy model
// declares TableNotFoundException (rather than the general
// ResourceNotFoundException) must pass ErrTableNotFoundException here so the
// client receives the individual error code.
func (s *DynamoDBService) validateAndGetTableWithErr(reqCtx *request.RequestContext, params map[string]interface{}, notFoundErr *APIError) (*dbstore.Table, error) {
	tableName := request.GetStringParam(params, "TableName")
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := store.Tables().Get(tableName)
	if err != nil {
		if dbstore.IsTableNotFound(err) || commonstore.IsNotFound(err) {
			return nil, notFoundErr
		}
		return nil, err
	}
	return table, nil
}

// createTableCore is the single entry point for table creation shared by the
// HTTP API and the admin gRPC handler. It applies defaults, performs all
// validation, creates the table in the store, applies post-create field
// updates, tags the table, and returns the fully-created table.
func (s *DynamoDBService) createTableCore(store dbstore.DynamoDBStoreInterface, in CreateTableInput) (*dbstore.Table, error) {
	// 1. Table name validation (length 3-255, allowed characters).
	if !validateResourceName(in.TableName) {
		return nil, ErrInvalidParameter
	}

	// 2. Key schema must be non-empty and well-formed.
	if len(in.KeySchema) == 0 {
		return nil, ErrInvalidParameter
	}
	if !validateKeySchema(in.KeySchema) {
		return nil, ErrInvalidParameter
	}

	// 3. Attribute definitions must cover every key attribute.
	if !validateAttributeDefinitions(in.KeySchema, in.AttributeDefinitions) {
		return nil, ErrInvalidParameter
	}

	// 4. Billing mode default + consistency check.
	if in.BillingMode == "" {
		in.BillingMode = dbstore.BillingModeProvisioned
	}
	if !validateBillingModeValue(in.BillingMode) {
		return nil, ErrInvalidParameter
	}
	if !validateBillingModeConsistency(in.BillingMode, in.ProvisionedThroughput) {
		return nil, ErrInvalidParameter
	}
	// Provisioned throughput values must fall inside the documented quota
	// range (DescribeLimits reports the same constants).
	if in.ProvisionedThroughput != nil && !validateProvisionedThroughputValues(in.ProvisionedThroughput) {
		return nil, ErrInvalidParameter
	}
	// A GSI that carries throughput settings must carry valid ones
	// (PositiveLongObject minimum 1, same as the table-level check).
	for _, gsi := range in.GlobalSecondaryIndexes {
		if gsi.ProvisionedThroughput != nil && !validateProvisionedThroughputValues(gsi.ProvisionedThroughput) {
			return nil, ErrInvalidParameter
		}
	}

	// 5. Cross-index validation: every key attribute across table + GSIs +
	//    LSIs must appear in AttributeDefinitions; LSI partition key must
	//    match the table partition key; index names must be unique.
	if !validateAllKeyAttributesInDefs(in.KeySchema, in.GlobalSecondaryIndexes, in.LocalSecondaryIndexes, in.AttributeDefinitions) {
		return nil, ErrInvalidParameter
	}
	if !validateLSIPartitionKey(in.KeySchema, in.LocalSecondaryIndexes) {
		return nil, ErrInvalidParameter
	}
	if err := validateIndexNameUniqueness(in.GlobalSecondaryIndexes, in.LocalSecondaryIndexes, in.VectorIndexes); err != nil {
		return nil, err
	}
	if err := validateVectorAttributeDimensions(in.VectorIndexes); err != nil {
		return nil, err
	}
	if len(in.VectorIndexes) > dbstore.VectorIndexesPerTable {
		return nil, ErrInvalidParameter
	}
	if in.TableClass != "" && !validateTableClassValue(in.TableClass) {
		return nil, ErrInvalidParameter
	}
	if err := validateStreamSpecification(in.StreamSpecification); err != nil {
		return nil, err
	}

	// 6. Persist.
	table, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                      in.TableName,
		KeySchema:                 in.KeySchema,
		AttributeDefinitions:      in.AttributeDefinitions,
		BillingMode:               in.BillingMode,
		ProvisionedThroughput:     in.ProvisionedThroughput,
		GlobalSecondaryIndexes:    in.GlobalSecondaryIndexes,
		LocalSecondaryIndexes:     in.LocalSecondaryIndexes,
		StreamSpecification:       in.StreamSpecification,
		Tags:                      in.Tags,
		DeletionProtectionEnabled: in.DeletionProtectionEnabled,
	})
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
	if len(in.VectorIndexes) > 0 {
		for _, vi := range in.VectorIndexes {
			vi.IndexArn = table.ARN + "/index/" + vi.IndexName
		}
		table.VectorIndexes = in.VectorIndexes
		needsPersist = true
	}
	if needsPersist {
		if err := store.Tables().Put(table); err != nil {
			return nil, err
		}
	}

	// 8. Tags.
	if len(in.Tags) > 0 {
		if err := store.Tables().Tags().Tag(in.TableName, tagutil.ToMap(in.Tags)); err != nil {
			return nil, err
		}
	}

	return table, nil
}

// ---------------------------------------------------------------------------
// Core functions for DeleteTable, DescribeTable, ListTables
// ---------------------------------------------------------------------------

// deleteTableCore is the single entry point for table deletion shared by the
// HTTP API and the admin gRPC handler. It validates the table name, checks
// for deletion protection, performs the cascade delete, and returns the
// archived table description. The cascade transaction runs under the
// table's record lock so a concurrent metric flush cannot write the record
// back after the cascade deletes it (the cascade queues no metric deltas,
// so the post-commit flush never re-takes the lock).
func (s *DynamoDBService) deleteTableCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName string) (*dbstore.Table, error) {
	if !validateResourceName(tableName) {
		return nil, ErrInvalidParameter
	}

	var deletedTable *dbstore.Table

	err := store.Tables().WithTableLock(tableName, func() error {
		return store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
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
	})
	if err != nil {
		return nil, err
	}

	deletedTable.Status = dbstore.TableStatusDeleting
	return deletedTable, nil
}

// replicasForTable renders the Replicas member a TableDescription reports:
// every region of the table's global-table replication group, or nil for a
// standalone table. Replica tables are created synchronously on this
// platform, so each region renders as a live ACTIVE replica unless the
// record says otherwise.
func (s *DynamoDBService) replicasForTable(store dbstore.DynamoDBStoreInterface, tableName string) []interface{} {
	gt, err := store.GlobalTables().Get(tableName)
	if err != nil || gt == nil || len(gt.ReplicationGroup) == 0 {
		return nil
	}
	replicas := make([]interface{}, 0, len(gt.ReplicationGroup))
	for _, r := range gt.ReplicationGroup {
		status := r.ReplicaStatus
		if status == "" {
			status = "ACTIVE"
		}
		replicas = append(replicas, map[string]interface{}{
			"RegionName":    r.RegionName,
			"ReplicaStatus": status,
		})
	}
	return replicas
}

// describeTableCore is the single entry point for table lookup shared by the
// HTTP API and the admin gRPC handler. It validates the table name and maps
// store errors to DynamoDB API errors.
func (s *DynamoDBService) describeTableCore(store dbstore.DynamoDBStoreInterface, tableName string) (*dbstore.Table, error) {
	if !validateResourceName(tableName) {
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
// HTTP API and the admin gRPC handler. It validates the Limit range (the
// model's ListTablesInputLimit, 1-100) and the marker's table-name shape
// before delegating to the store; an absent limit arrives as the documented
// maximum from both planes.
func (s *DynamoDBService) listTablesCore(store dbstore.DynamoDBStoreInterface, marker string, limit int) ([]*dbstore.Table, string, error) {
	if !validateListTablesLimit(limit) {
		return nil, "", ErrInvalidParameter
	}
	if marker != "" && !validateResourceName(marker) {
		return nil, "", ErrInvalidParameter
	}
	return store.Tables().List(marker, limit)
}

// ---------------------------------------------------------------------------
// Transport-agnostic DTO for UpdateTable
// ---------------------------------------------------------------------------

// UpdateTableInput carries every field that UpdateTable can modify, in a
// format independent of the wire protocol. Fields with zero values are
// treated as "no change". Both the HTTP API handler and the admin gRPC
// handler build this struct and delegate to updateTableCore.
type UpdateTableInput struct {
	TableName             string
	BillingMode           string                         // "" = no change
	ProvisionedThroughput *dbstore.ProvisionedThroughput // nil = no change
	AttributeDefinitions  []*dbstore.AttributeDefinition // nil = no change
	GSIUpdates            []interface{}                  // nil = no change
	VectorIndexUpdates    []interface{}                  // nil = no change
	StreamSpecification   *dbstore.StreamSpecification   // nil = no change
	SSESpecification      *dbstore.SSEDescription        // nil = no change
	DeletionProtectionSet bool                           // whether DeletionProtectionEnabled was provided
	DeletionProtection    bool                           // value (only if DeletionProtectionSet is true)
	TableClass            string                         // "" = no change
}

// updateTableCore is the single entry point for table updates shared by the
// HTTP API and the admin gRPC handler. It validates all update parameters,
// applies the changes to the table, persists them, and backfills any newly
// created GSIs. The record update runs through the store's locked
// read-modify-write path, so concurrent counter flushes and settings writes
// cannot be lost; the GSI backfill runs after the lock is released because
// its chunk transactions take the same lock for their own metric flushes.
func (s *DynamoDBService) updateTableCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, in UpdateTableInput) (*dbstore.Table, error) {
	if !validateResourceName(in.TableName) {
		return nil, ErrInvalidParameter
	}
	if in.BillingMode != "" && !validateBillingModeValue(dbstore.BillingMode(in.BillingMode)) {
		return nil, ErrInvalidParameter
	}
	if in.TableClass != "" && !validateTableClassValue(in.TableClass) {
		return nil, ErrInvalidParameter
	}
	if err := validateStreamSpecification(in.StreamSpecification); err != nil {
		return nil, err
	}
	if in.ProvisionedThroughput != nil && !validateProvisionedThroughputValues(in.ProvisionedThroughput) {
		return nil, ErrInvalidParameter
	}

	existingGSINames := make(map[string]bool)
	deletedGSINames := []string{}
	vectorCreatedNames := []string{}
	vectorDeletedNames := []string{}
	table, err := store.Tables().Update(in.TableName, func(table *dbstore.Table) error {
		if table.Status != dbstore.TableStatusActive {
			return ErrTableNotActive
		}

		if in.BillingMode != "" {
			table.BillingMode = dbstore.BillingMode(in.BillingMode)
			if table.BillingMode == dbstore.BillingModePayPerRequest {
				// Switching to on-demand discards the provisioned capacity:
				// an on-demand table keeps no throughput settings, and the
				// consistency check below then rejects a request that tries
				// to carry ProvisionedThroughput into the new mode.
				table.ProvisionedThroughput = nil
			}
		}

		if in.ProvisionedThroughput != nil {
			table.ProvisionedThroughput = in.ProvisionedThroughput
		}

		if !validateBillingModeConsistency(table.BillingMode, table.ProvisionedThroughput) {
			return ErrInvalidParameter
		}

		if len(in.AttributeDefinitions) > 0 {
			if !validateAttributeDefinitions(table.KeySchema, in.AttributeDefinitions) {
				return ErrInvalidParameter
			}
			table.AttributeDefinitions = mergeAttributeDefinitions(table.AttributeDefinitions, in.AttributeDefinitions)
		}

		for _, g := range table.GlobalSecondaryIndexes {
			existingGSINames[g.IndexName] = true
		}

		if len(in.GSIUpdates) > 0 {
			updatedGSIs, deleted, err := applyGSIUpdates(table.ARN, table.GlobalSecondaryIndexes, in.GSIUpdates)
			if err != nil {
				return err
			}
			table.GlobalSecondaryIndexes = updatedGSIs
			deletedGSINames = deleted
		}

		if len(in.VectorIndexUpdates) > 0 {
			updatedVector, created, deleted, err := applyVectorIndexUpdates(table.ARN, table.VectorIndexes, in.VectorIndexUpdates)
			if err != nil {
				return err
			}
			table.VectorIndexes = updatedVector
			vectorCreatedNames = created
			vectorDeletedNames = deleted
		}

		if !validateAllKeyAttributesInDefs(table.KeySchema, table.GlobalSecondaryIndexes, table.LocalSecondaryIndexes, table.AttributeDefinitions) {
			return ErrInvalidParameter
		}
		if err := validateIndexNameUniqueness(table.GlobalSecondaryIndexes, table.LocalSecondaryIndexes, table.VectorIndexes); err != nil {
			return err
		}
		if err := validateVectorAttributeDimensions(table.VectorIndexes); err != nil {
			return err
		}
		if len(table.VectorIndexes) > dbstore.VectorIndexesPerTable {
			return ErrInvalidParameter
		}

		if in.StreamSpecification != nil {
			table.StreamSpecification = in.StreamSpecification
			if in.StreamSpecification.StreamEnabled {
				now := time.Now().UTC()
				table.StreamArn = table.ARN + "/stream/" + now.Format("2006-01-02T15:04:05.000")
				table.LatestStreamLabel = now.Format("2006-01-02T15:04:05.000")
			} else {
				table.StreamArn = ""
				table.LatestStreamLabel = ""
			}
		}

		if in.SSESpecification != nil {
			table.SSEDescription = in.SSESpecification
		}

		if in.DeletionProtectionSet {
			table.DeletionProtectionEnabled = in.DeletionProtection
		}

		if in.TableClass != "" {
			table.TableClass = in.TableClass
		}

		table.LastUpdatedDateTime = time.Now().UTC()
		return nil
	})
	if err != nil {
		if dbstore.IsTableNotFound(err) || commonstore.IsNotFound(err) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	// A deleted index's entries must not outlive the index: a re-created
	// same-name index would inherit stale entries that resolve to items no
	// longer holding the index key attributes. The cleanup runs in its own
	// transaction after the locked table update — TableStore.Update exposes
	// no transaction to the mutate callback, and the storage layer is
	// read-committed, so folding it into the table write would not close
	// the stale-reader window either. Item writes between the two commits
	// write no entries for the deleted index: the table record they read no
	// longer lists it. Both deletion sweeps — this one and the vector twin
	// below — run ahead of the backfills: the metadata deletions have
	// already committed and a re-sent request carries no deletion to
	// re-run, so a backfill failure must not be able to strand the sweeps.
	if len(deletedGSINames) > 0 {
		err := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for _, name := range deletedGSINames {
				if err := txn.DeleteIndexEntriesForIndex(table.Name, name); err != nil {
					return fmt.Errorf("delete index entries for %s: %w", name, err)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// A deleted vector index's entries must not outlive it — the same
	// lifecycle the GSI sweep above applies, through the vector twin
	// helper, and under the same ordering rule: ahead of every backfill.
	if len(vectorDeletedNames) > 0 {
		err := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for _, name := range vectorDeletedNames {
				if err := txn.DeleteVectorEntriesForIndex(table.Name, name); err != nil {
					return fmt.Errorf("delete vector entries for %s: %w", name, err)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	for _, g := range table.GlobalSecondaryIndexes {
		if existingGSINames[g.IndexName] {
			continue
		}
		if err := s.backfillGSI(ctx, store, table.Name, g.IndexName); err != nil {
			return nil, err
		}
	}

	// A newly created vector index is populated from the existing items —
	// the vector twin of the GSI backfill above, through the same helper.
	for _, name := range vectorCreatedNames {
		if err := s.backfillVectorIndex(ctx, store, table.Name, name); err != nil {
			return nil, err
		}
	}

	return table, nil
}
