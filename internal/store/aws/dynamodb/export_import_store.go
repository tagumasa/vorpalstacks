// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

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
func (s *ExportStore) Create(tableArn, tableId, exportFormat string) (*ExportDescription, error) {
	now := time.Now().UTC()
	exportArn := s.arnBuilder.Export(tableArn, now.Format("20060102150405"))

	export := &ExportDescription{
		ExportArn:    exportArn,
		ExportStatus: "IN_PROGRESS",
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
func (s *ExportStore) List(tableArn, marker string, maxItems int) ([]*ExportDescription, string, error) {
	if maxItems <= 0 {
		maxItems = 100
	}
	opts := common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	}

	filter := func(e *pb.ExportDescription) bool {
		if tableArn == "" {
			return true
		}
		return e.TableArn == tableArn
	}

	result, err := common.ListProto[*pb.ExportDescription](s.BaseStore, opts, func() *pb.ExportDescription { return &pb.ExportDescription{} }, filter)
	if err != nil {
		return nil, "", err
	}

	exports := make([]*ExportDescription, len(result.Items))
	for i, e := range result.Items {
		exports[i] = ProtoToExportDescription(e)
	}

	nextToken := ""
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return exports, nextToken, nil
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
	importArn := s.arnBuilder.Import(tableArn, now.Format("20060102150405"))

	imp := &ImportTableDescription{
		ImportArn:    importArn,
		ImportStatus: "IN_PROGRESS",
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
func (s *ImportStore) List(tableArn, marker string, maxItems int) ([]*ImportTableDescription, string, error) {
	if maxItems <= 0 {
		maxItems = 100
	}
	opts := common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	}

	filter := func(i *pb.ImportTableDescription) bool {
		if tableArn == "" {
			return true
		}
		return i.TableArn == tableArn
	}

	result, err := common.ListProto[*pb.ImportTableDescription](s.BaseStore, opts, func() *pb.ImportTableDescription { return &pb.ImportTableDescription{} }, filter)
	if err != nil {
		return nil, "", err
	}

	imports := make([]*ImportTableDescription, len(result.Items))
	for i, imp := range result.Items {
		imports[i] = ProtoToImportTableDescription(imp)
	}

	nextToken := ""
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return imports, nextToken, nil
}
