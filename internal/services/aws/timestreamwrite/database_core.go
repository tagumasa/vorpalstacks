package timestreamwrite

import (
	"context"
	"time"

	types "vorpalstacks/internal/common/tags"
	storecommon "vorpalstacks/internal/store/aws/common"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Output DTOs — Database
// ---------------------------------------------------------------------------

// CreateDatabaseInput carries every field that CreateDatabase needs, in a
// format independent of the wire protocol (HTTP Query/JSON vs gRPC-Web).
type CreateDatabaseInput struct {
	DatabaseName string
	KmsKeyId     string
	Tags         []types.Tag
	TagsProvided bool
}

// UpdateDatabaseInput carries the fields for UpdateDatabase.
type UpdateDatabaseInput struct {
	DatabaseName string
	KmsKeyId     string
}

// DescribeDatabaseInput carries the fields for DescribeDatabase.
type DescribeDatabaseInput struct {
	DatabaseName string
}

// DeleteDatabaseInput carries the fields for DeleteDatabase.
type DeleteDatabaseInput struct {
	DatabaseName string
}

// ListDatabasesInput carries the fields for ListDatabases.
type ListDatabasesInput struct {
	NextToken string
	MaxItems  int
}

// DatabaseResult is the transport-agnostic database representation returned
// by Core functions. It stores time.Time values so each transport layer can
// format them appropriately (epoch float64 for HTTP, ISO8601 for gRPC-Web).
type DatabaseResult struct {
	ARN             string
	DatabaseName    string
	TableCount      int64
	KmsKeyId        string
	CreationTime    time.Time
	LastUpdatedTime time.Time
	Tags            []types.Tag
}

// ListDatabasesResult is the paginated result of listing databases.
type ListDatabasesResult struct {
	Databases []DatabaseResult
	NextToken string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createDatabaseCore is the single entry point for database creation shared by
// the HTTP API and the admin gRPC handler. It performs all Smithy-conformant
// validation and persists to the store.
func (s *TimestreamWriteService) createDatabaseCore(ctx context.Context, stores *tsWriteStores, in CreateDatabaseInput) (*DatabaseResult, error) {
	if in.DatabaseName == "" {
		return nil, ErrValidationException
	}
	if !isValidTimestreamName(in.DatabaseName) {
		return nil, ErrValidationException
	}

	kmsKeyID := in.KmsKeyId
	if kmsKeyID != "" && !validateKmsKeyId(kmsKeyID) {
		return nil, ErrValidationException
	}

	db, err := stores.store.CreateDatabase(in.DatabaseName, kmsKeyID)
	if err != nil {
		if err == tsstore.ErrDatabaseAlreadyExists {
			return nil, ErrConflictException
		}
		return nil, ErrInternalServer
	}

	if in.TagsProvided && len(in.Tags) > 0 {
		if err := stores.store.TagFromSlice(db.ARN, in.Tags); err != nil {
			return nil, ErrInternalServer
		}
	}

	tags, _ := stores.store.ListAsSlice(db.ARN)

	return databaseToResult(db, tags), nil
}

// describeDatabaseCore is the single entry point for DescribeDatabase.
func (s *TimestreamWriteService) describeDatabaseCore(ctx context.Context, stores *tsWriteStores, in DescribeDatabaseInput) (*DatabaseResult, error) {
	if in.DatabaseName == "" {
		return nil, ErrValidationException
	}

	db, err := stores.store.GetDatabase(in.DatabaseName)
	if err != nil {
		if err == tsstore.ErrDatabaseNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	tags, _ := stores.store.ListAsSlice(db.ARN)

	return databaseToResult(db, tags), nil
}

// listDatabasesCore is the single entry point for ListDatabases.
func (s *TimestreamWriteService) listDatabasesCore(ctx context.Context, stores *tsWriteStores, in ListDatabasesInput) (*ListDatabasesResult, error) {
	maxResults := in.MaxItems
	if maxResults <= 0 {
		maxResults = maxListDatabasesResults
	}
	if maxResults > maxListDatabasesResults {
		maxResults = maxListDatabasesResults
	}

	opts := storecommon.ListOptions{MaxItems: maxResults}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}

	result, err := stores.store.ListDatabases(opts)
	if err != nil {
		return nil, ErrInternalServer
	}

	databases := make([]DatabaseResult, 0, len(result.Items))
	for _, db := range result.Items {
		tags, _ := stores.store.ListAsSlice(db.ARN)
		databases = append(databases, *databaseToResult(db, tags))
	}

	return &ListDatabasesResult{
		Databases: databases,
		NextToken: result.NextMarker,
	}, nil
}

// updateDatabaseCore is the single entry point for UpdateDatabase.
func (s *TimestreamWriteService) updateDatabaseCore(ctx context.Context, stores *tsWriteStores, in UpdateDatabaseInput) (*DatabaseResult, error) {
	if in.DatabaseName == "" {
		return nil, ErrValidationException
	}
	if in.KmsKeyId == "" {
		return nil, ErrValidationException
	}
	if !validateKmsKeyId(in.KmsKeyId) {
		return nil, ErrValidationException
	}

	db, err := stores.store.UpdateDatabase(in.DatabaseName, in.KmsKeyId)
	if err != nil {
		if err == tsstore.ErrDatabaseNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	tags, _ := stores.store.ListAsSlice(db.ARN)

	return databaseToResult(db, tags), nil
}

// deleteDatabaseCore is the single entry point for DeleteDatabase.
func (s *TimestreamWriteService) deleteDatabaseCore(ctx context.Context, stores *tsWriteStores, in DeleteDatabaseInput) error {
	if in.DatabaseName == "" {
		return ErrValidationException
	}

	err := stores.store.DeleteDatabase(in.DatabaseName)
	if err != nil {
		if err == tsstore.ErrDatabaseNotFound {
			return ErrResourceNotFound
		}
		if err == tsstore.ErrDatabaseNotEmpty {
			return ErrValidationException
		}
		return ErrInternalServer
	}

	stores.recordStore.DeleteDatabaseChunks(in.DatabaseName)

	return nil
}

// databaseToResult converts a store-layer Database to a transport-agnostic
// DatabaseResult, including tags.
func databaseToResult(db *tsstore.Database, tags []types.Tag) *DatabaseResult {
	return &DatabaseResult{
		ARN:             db.ARN,
		DatabaseName:    db.DatabaseName,
		TableCount:      db.TableCount,
		KmsKeyId:        db.KmsKeyId,
		CreationTime:    db.CreationTime,
		LastUpdatedTime: db.LastUpdatedTime,
		Tags:            tags,
	}
}
