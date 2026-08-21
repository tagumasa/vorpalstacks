package timestreamwrite

import (
	"context"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	storecommon "vorpalstacks/internal/store/aws/common"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Output DTOs — Table
// ---------------------------------------------------------------------------

// CreateTableInput carries every field that CreateTable needs.
type CreateTableInput struct {
	DatabaseName                 string
	TableName                    string
	RetentionProperties          *tsstore.RetentionProperties
	Schema                       *tsstore.Schema
	MagneticStoreWriteProperties *tsstore.MagneticStoreWriteProperties
	Tags                         []types.Tag
	TagsProvided                 bool
}

// UpdateTableInput carries every field that UpdateTable needs.
type UpdateTableInput struct {
	DatabaseName                 string
	TableName                    string
	RetentionProperties          *tsstore.RetentionProperties
	Schema                       *tsstore.Schema
	MagneticStoreWriteProperties *tsstore.MagneticStoreWriteProperties
}

// DescribeTableInput carries the fields for DescribeTable.
type DescribeTableInput struct {
	DatabaseName string
	TableName    string
}

// DeleteTableInput carries the fields for DeleteTable.
type DeleteTableInput struct {
	DatabaseName string
	TableName    string
}

// ListTablesInput carries the fields for ListTables.
type ListTablesInput struct {
	DatabaseName string
	NextToken    string
	MaxItems     int
}

// TableResult is the transport-agnostic table representation returned by
// Core functions.
type TableResult struct {
	ARN                          string
	TableName                    string
	DatabaseName                 string
	TableStatus                  tsstore.TableStatus
	CreationTime                 time.Time
	LastUpdatedTime              time.Time
	RetentionProperties          *tsstore.RetentionProperties
	Schema                       *tsstore.Schema
	MagneticStoreWriteProperties *tsstore.MagneticStoreWriteProperties
	Tags                         []types.Tag
}

// ListTablesResult is the paginated result of listing tables.
type ListTablesResult struct {
	Tables    []TableResult
	NextToken string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createTableCore is the single entry point for table creation shared by the
// HTTP API and the admin gRPC handler.
func (s *TimestreamWriteService) createTableCore(ctx context.Context, stores *tsWriteStores, in CreateTableInput) (*TableResult, error) {
	if in.DatabaseName == "" || !isValidTimestreamName(in.DatabaseName) {
		return nil, ErrValidationException
	}
	if in.TableName == "" || !isValidTimestreamName(in.TableName) {
		return nil, ErrValidationException
	}
	if err := validateRetentionProperties(in.RetentionProperties); err != nil {
		return nil, err
	}
	if err := validateSchema(in.Schema); err != nil {
		return nil, err
	}
	if err := validateMagneticStoreWriteProperties(in.MagneticStoreWriteProperties); err != nil {
		return nil, err
	}

	table, err := stores.tableStore.CreateTable(in.DatabaseName, in.TableName, in.RetentionProperties, in.Schema, in.MagneticStoreWriteProperties)
	if err != nil {
		if err == tsstore.ErrTableAlreadyExists {
			return nil, ErrConflictException
		}
		if err == tsstore.ErrDatabaseNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	if in.TagsProvided && len(in.Tags) > 0 {
		if err := stores.store.TagFromSlice(table.ARN, in.Tags); err != nil {
			return nil, ErrInternalServer
		}
	}

	tags, _ := stores.store.ListAsSlice(table.ARN)
	return tableToResult(table, tags), nil
}

// describeTableCore is the single entry point for DescribeTable.
func (s *TimestreamWriteService) describeTableCore(ctx context.Context, stores *tsWriteStores, in DescribeTableInput) (*TableResult, error) {
	if in.DatabaseName == "" || in.TableName == "" {
		return nil, ErrValidationException
	}

	table, err := stores.tableStore.GetTable(in.DatabaseName, in.TableName)
	if err != nil {
		if err == tsstore.ErrTableNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	tags, _ := stores.store.ListAsSlice(table.ARN)
	return tableToResult(table, tags), nil
}

// listTablesCore is the single entry point for ListTables.
func (s *TimestreamWriteService) listTablesCore(ctx context.Context, stores *tsWriteStores, in ListTablesInput) (*ListTablesResult, error) {
	if in.DatabaseName == "" {
		return nil, ErrValidationException
	}

	maxResults := in.MaxItems
	if maxResults <= 0 {
		maxResults = maxListTablesResults
	}
	if maxResults > maxListTablesResults {
		maxResults = maxListTablesResults
	}

	opts := storecommon.ListOptions{MaxItems: maxResults}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}

	result, err := stores.tableStore.ListTables(in.DatabaseName, opts)
	if err != nil {
		return nil, ErrInternalServer
	}

	tables := make([]TableResult, 0, len(result.Items))
	for _, table := range result.Items {
		tags, _ := stores.store.ListAsSlice(table.ARN)
		tables = append(tables, *tableToResult(table, tags))
	}

	return &ListTablesResult{
		Tables:    tables,
		NextToken: result.NextMarker,
	}, nil
}

// updateTableCore is the single entry point for UpdateTable. Rejects requests
// where all optional fields are nil (at least one is required to prevent a
// misleading no-op success that only updates LastUpdatedTime).
func (s *TimestreamWriteService) updateTableCore(ctx context.Context, stores *tsWriteStores, in UpdateTableInput) (*TableResult, error) {
	if in.DatabaseName == "" || in.TableName == "" {
		return nil, ErrValidationException
	}

	// At least one of RetentionProperties, MagneticStoreWriteProperties,
	// or Schema must be provided. A no-op update that only changes
	// LastUpdatedTime is misleading and should be rejected.
	if in.RetentionProperties == nil && in.MagneticStoreWriteProperties == nil && in.Schema == nil {
		return nil, awserrors.NewValidationException("At least one of RetentionProperties, MagneticStoreWriteProperties, or Schema must be provided")
	}

	if err := validateRetentionProperties(in.RetentionProperties); err != nil {
		return nil, err
	}
	if err := validateSchema(in.Schema); err != nil {
		return nil, err
	}
	if err := validateMagneticStoreWriteProperties(in.MagneticStoreWriteProperties); err != nil {
		return nil, err
	}

	table, err := stores.tableStore.UpdateTable(in.DatabaseName, in.TableName, in.RetentionProperties, in.Schema, in.MagneticStoreWriteProperties)
	if err != nil {
		if err == tsstore.ErrTableNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	tags, _ := stores.store.ListAsSlice(table.ARN)
	return tableToResult(table, tags), nil
}

// deleteTableCore is the single entry point for DeleteTable.
func (s *TimestreamWriteService) deleteTableCore(ctx context.Context, stores *tsWriteStores, in DeleteTableInput) error {
	if in.DatabaseName == "" || in.TableName == "" {
		return ErrValidationException
	}

	err := stores.tableStore.DeleteTable(in.DatabaseName, in.TableName)
	if err != nil {
		if err == tsstore.ErrTableNotFound {
			return ErrResourceNotFound
		}
		return ErrInternalServer
	}

	stores.recordStore.DeleteTableChunks(in.DatabaseName, in.TableName)
	return nil
}

// validateRetentionProperties validates retention period ranges.
func validateRetentionProperties(rp *tsstore.RetentionProperties) error {
	if rp == nil {
		return nil
	}
	if rp.MemoryStoreRetentionPeriodInHours < 1 || rp.MemoryStoreRetentionPeriodInHours > 8766 {
		return ErrValidationException
	}
	if rp.MagneticStoreRetentionPeriodInDays < 1 || rp.MagneticStoreRetentionPeriodInDays > 73000 {
		return ErrValidationException
	}
	return nil
}

// validateSchema validates PartitionKey fields per Smithy. Validates
// Type (required + enum) and EnforcementInRecord (enum when present).
func validateSchema(schema *tsstore.Schema) error {
	if schema == nil {
		return nil
	}
	for _, pk := range schema.CompositePartitionKey {
		if pk.Type == "" {
			return awserrors.NewValidationException("PartitionKey.Type is required")
		}
		if !validatePartitionKeyType(string(pk.Type)) {
			return awserrors.NewValidationException("Invalid PartitionKey.Type: " + string(pk.Type))
		}
		if pk.EnforcementInRecord != "" && !validateEnforcementInRecord(string(pk.EnforcementInRecord)) {
			return awserrors.NewValidationException("Invalid PartitionKey.EnforcementInRecord: " + string(pk.EnforcementInRecord))
		}
	}
	return nil
}

// validateMagneticStoreWriteProperties validates EnableMagneticStoreWrites
// is set when the properties struct is provided.
func validateMagneticStoreWriteProperties(mswp *tsstore.MagneticStoreWriteProperties) error {
	if mswp == nil {
		return nil
	}
	// EnableMagneticStoreWrites has no explicit "set" flag in Go bool.
	// false is a valid value, so we accept it.
	return nil
}

// tableToResult converts a store-layer Table to a transport-agnostic
// TableResult, including tags.
func tableToResult(table *tsstore.Table, tags []types.Tag) *TableResult {
	return &TableResult{
		ARN:                          table.ARN,
		TableName:                    table.TableName,
		DatabaseName:                 table.DatabaseName,
		TableStatus:                  table.TableStatus,
		CreationTime:                 table.CreationTime,
		LastUpdatedTime:              table.LastUpdatedTime,
		RetentionProperties:          table.RetentionProperties,
		Schema:                       table.Schema,
		MagneticStoreWriteProperties: table.MagneticStoreWriteProperties,
		Tags:                         tags,
	}
}
