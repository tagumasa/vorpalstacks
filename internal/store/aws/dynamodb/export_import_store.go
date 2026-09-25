// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"crypto/rand"
	"fmt"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// tableScopedId mints the distinguishing segment of an export, import, or
// backup ARN: the creation time at millisecond granularity followed by
// a random 8-hex-digit suffix, the shape the documented ARN examples of
// all three families show (export 01234567890123-a1b2c3d4, backup
// table/Music/backup/01489602797149-73d8d5bc). The timestamp alone cannot
// tell two creations apart inside one millisecond; the random suffix
// is what keeps concurrent same-table creations from sharing the ARN
// their records are keyed by — the identity mechanism the documented
// concurrency quotas rely on, with no one-in-progress guard on the API to
// mask a collision. For backups the same mechanism is what lets
// same-named generations of one table name coexist, each addressed by its
// own ARN.
func tableScopedId(now time.Time) (string, error) {
	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return "", fmt.Errorf("generate table-scoped resource id suffix: %w", err)
	}
	return fmt.Sprintf("%014d-%x", now.UnixMilli(), suffix[:]), nil
}

func exportBucketName(region string) string {
	return "dynamodb_exports-" + region
}

func importBucketName(region string) string {
	return "dynamodb_imports-" + region
}

// ExportStore manages DynamoDB table exports to S3.
type ExportStore struct {
	*common.BaseStore
	arnBuilder *svcarn.DynamoDBBuilder
}

// NewExportStore creates a new export store for DynamoDB exports.
func NewExportStore(store storage.BasicStorage, accountId, region string) *ExportStore {
	return &ExportStore{
		BaseStore:  common.NewBaseStore(store.Bucket(exportBucketName(region)), "dynamodb_exports"),
		arnBuilder: svcarn.NewARNBuilder(accountId, region).DynamoDB(),
	}
}

// Get retrieves an export description by its ARN.
func (s *ExportStore) Get(exportArn string) (*ExportDescription, error) {
	var export pb.ExportDescription
	if err := s.BaseStore.GetProto(exportArn, &export); err != nil {
		return nil, err
	}
	return ProtoToExportDescription(&export), nil
}

// Create initiates a new export of a DynamoDB table to S3.
func (s *ExportStore) Create(tableArn, tableId string, exportFormat ExportFormat) (*ExportDescription, error) {
	now := time.Now().UTC()
	exportId, err := tableScopedId(now)
	if err != nil {
		return nil, err
	}
	exportArn := s.arnBuilder.Export(tableArn, exportId)

	export := &ExportDescription{
		ExportArn:    exportArn,
		ExportStatus: ExportStatusInProgress,
		StartTime:    now,
		TableArn:     tableArn,
		TableId:      tableId,
		ExportFormat: exportFormat,
	}

	if err := s.BaseStore.PutProto(exportArn, ExportDescriptionToProto(export)); err != nil {
		return nil, err
	}

	return export, nil
}

// Put updates an existing export description.
func (s *ExportStore) Put(export *ExportDescription) error {
	return s.BaseStore.PutProto(export.ExportArn, ExportDescriptionToProto(export))
}

// List returns exports, optionally filtered by table ARN, with pagination.
// A non-positive maxItems reads as the pagination layer's default (the
// normalisation the shared list path applies to every family).
func (s *ExportStore) List(tableArn, marker string, maxItems int) ([]*ExportDescription, string, error) {
	filter := func(e *pb.ExportDescription) bool {
		if tableArn == "" {
			return true
		}
		return e.TableArn == tableArn
	}
	return listProtoConverted(s.BaseStore, marker, maxItems, func() *pb.ExportDescription { return &pb.ExportDescription{} }, ProtoToExportDescription, filter)
}

// ImportStore manages DynamoDB table imports from S3.
type ImportStore struct {
	*common.BaseStore
	arnBuilder *svcarn.DynamoDBBuilder
}

// NewImportStore creates a new import store for DynamoDB imports.
func NewImportStore(store storage.BasicStorage, accountId, region string) *ImportStore {
	return &ImportStore{
		BaseStore:  common.NewBaseStore(store.Bucket(importBucketName(region)), "dynamodb_imports"),
		arnBuilder: svcarn.NewARNBuilder(accountId, region).DynamoDB(),
	}
}

// Get retrieves an import description by its ARN.
func (s *ImportStore) Get(importArn string) (*ImportTableDescription, error) {
	var imp pb.ImportTableDescription
	if err := s.BaseStore.GetProto(importArn, &imp); err != nil {
		return nil, err
	}
	return ProtoToImportTableDescription(&imp), nil
}

// Create initiates a new import of a DynamoDB table from S3.
func (s *ImportStore) Create(tableArn, tableId string) (*ImportTableDescription, error) {
	now := time.Now().UTC()
	importId, err := tableScopedId(now)
	if err != nil {
		return nil, err
	}
	importArn := s.arnBuilder.Import(tableArn, importId)

	imp := &ImportTableDescription{
		ImportArn:    importArn,
		ImportStatus: ImportStatusInProgress,
		TableArn:     tableArn,
		TableId:      tableId,
		StartTime:    now,
	}

	if err := s.BaseStore.PutProto(importArn, ImportTableDescriptionToProto(imp)); err != nil {
		return nil, err
	}

	return imp, nil
}

// Put updates an existing import description.
func (s *ImportStore) Put(imp *ImportTableDescription) error {
	return s.BaseStore.PutProto(imp.ImportArn, ImportTableDescriptionToProto(imp))
}

// List returns imports, optionally filtered by table ARN, with pagination.
// A non-positive maxItems reads as the pagination layer's default (the
// normalisation the shared list path applies to every family).
func (s *ImportStore) List(tableArn, marker string, maxItems int) ([]*ImportTableDescription, string, error) {
	filter := func(i *pb.ImportTableDescription) bool {
		if tableArn == "" {
			return true
		}
		return i.TableArn == tableArn
	}
	return listProtoConverted(s.BaseStore, marker, maxItems, func() *pb.ImportTableDescription { return &pb.ImportTableDescription{} }, ProtoToImportTableDescription, filter)
}
