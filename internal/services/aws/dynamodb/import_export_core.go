package dynamodb

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Import / Export Core — single validation + persistence path for table
// import and export operations.
//
// These methods encapsulate import/export lifecycle logic, including S3
// invocations for data transfer. Both the HTTP API handlers
// (import_export_operations.go) and any future admin handler delegate to
// these methods.
// ---------------------------------------------------------------------------

// describeExportCore returns an export by ARN.
func (s *DynamoDBService) describeExportCore(store dbstore.DynamoDBStoreInterface, exportArn string) (*dbstore.ExportDescription, error) {
	export, err := store.Exports().Get(exportArn)
	if err != nil {
		return nil, ErrExportNotFound
	}
	return export, nil
}

// listExportsCore returns a paginated list of exports filtered by table ARN.
func (s *DynamoDBService) listExportsCore(store dbstore.DynamoDBStoreInterface, tableArn string, nextToken string, maxResults int) ([]*dbstore.ExportDescription, string, error) {
	return store.Exports().List(tableArn, nextToken, maxResults)
}

// describeImportCore returns an import by ARN.
func (s *DynamoDBService) describeImportCore(store dbstore.DynamoDBStoreInterface, importArn string) (*dbstore.ImportTableDescription, error) {
	imp, err := store.Imports().Get(importArn)
	if err != nil {
		return nil, ErrImportNotFound
	}
	return imp, nil
}

// listImportsCore returns a paginated list of imports filtered by table ARN.
func (s *DynamoDBService) listImportsCore(store dbstore.DynamoDBStoreInterface, tableArn, nextToken string, pageSize int) ([]*dbstore.ImportTableDescription, string, error) {
	return store.Imports().List(tableArn, nextToken, pageSize)
}

// ExportTableCoreInput is the service-layer DTO for ExportTableToPointInTime.
type ExportTableCoreInput struct {
	TableArn       string
	TableName      string
	ExportFormat   string
	S3Bucket       string
	S3Prefix       string
	ClientToken    string
	Region         string
	TableItemCount int64 // used when S3 invoker is nil
}

// exportTableCore creates an export record, scans the table, uploads the
// data to S3 (if the S3 invoker is available), and returns the export in
// its final state (COMPLETED or FAILED).
func (s *DynamoDBService) exportTableCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, in ExportTableCoreInput) (*dbstore.ExportDescription, error) {
	export, err := store.Exports().Create(in.TableArn, in.TableName, in.ExportFormat)
	if err != nil {
		return nil, err
	}

	itemCount := int64(0)
	billedSizeBytes := int64(0)

	s3 := s.s3invoker()
	exportFailed := false
	var failureCode, failureMessage string
	if s3 != nil {
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)

		err = store.Items().Scan(in.TableName, func(item *dbstore.Item) error {
			itemCount++
			itemJSON := buildDynamoDBJSONItem(item)
			line, mErr := json.Marshal(itemJSON)
			if mErr != nil {
				return mErr
			}
			if _, wErr := writer.Write(line); wErr != nil {
				return wErr
			}
			if _, wErr := writer.Write([]byte("\n")); wErr != nil {
				return wErr
			}
			return nil
		})
		if err != nil {
			logs.Error("Failed to scan items for export",
				logs.Err(err),
				logs.String("tableName", in.TableName))
			exportFailed = true
			failureCode = "InternalFailure"
			failureMessage = fmt.Sprintf("failed to scan items: %v", err)
		}

		if flushErr := writer.Flush(); flushErr != nil {
			logs.Error("Failed to flush export buffer",
				logs.Err(flushErr),
				logs.String("tableName", in.TableName))
			exportFailed = true
			failureCode = "InternalFailure"
			failureMessage = fmt.Sprintf("failed to flush buffer: %v", flushErr)
		}

		billedSizeBytes = int64(buf.Len())
		dataFileID := fmt.Sprintf("%d", time.Now().UnixNano())
		objectKey := fmt.Sprintf("%s/AWSDynamoDB/%s/data/%s.json",
			in.S3Prefix, export.ExportArn, dataFileID)
		if putErr := s3.PutObject(ctx, in.Region, in.S3Bucket, objectKey, buf.Bytes(), "application/octet-stream"); putErr != nil {
			logs.Error("Failed to write export to S3",
				logs.Err(putErr),
				logs.String("bucket", in.S3Bucket),
				logs.String("key", objectKey))
			exportFailed = true
			failureCode = "S3AccessDenied"
			failureMessage = fmt.Sprintf("failed to write to S3: %v", putErr)
		}
	} else {
		itemCount = in.TableItemCount
	}

	if exportFailed {
		export.ExportStatus = "FAILED"
		export.FailureCode = failureCode
		export.FailureMessage = failureMessage
	} else {
		export.ExportStatus = "COMPLETED"
	}
	export.S3Bucket = in.S3Bucket
	export.S3Prefix = in.S3Prefix
	export.EndTime = time.Now()
	export.ItemCount = itemCount
	export.BilledSizeBytes = billedSizeBytes
	if err := store.Exports().Put(export); err != nil {
		return nil, err
	}
	return export, nil
}

// ImportTableCoreInput is the service-layer DTO for ImportTable. It
// contains all parameters needed to create the target table and import
// data from S3.
type ImportTableCoreInput struct {
	TableName      string
	KeySchema      []*dbstore.KeySchemaElement
	AttributeDefs  []*dbstore.AttributeDefinition
	BillingMode    dbstore.BillingMode
	ProvThroughput *dbstore.ProvisionedThroughput
	GSI            []*dbstore.GlobalSecondaryIndex
	LSI            []*dbstore.LocalSecondaryIndex
	InputFormat    string
	S3Bucket       string
	S3Prefix       string
	S3BucketOwner  string
	ClientToken    string
	CSVDelimiter   string
	CSVHeaderList  []string
}

// importTableCore creates the target table, creates an import record,
// downloads data from S3 (if available), parses and inserts items, and
// returns the import in its final state.
func (s *DynamoDBService) importTableCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, region string, in ImportTableCoreInput) (*dbstore.ImportTableDescription, error) {
	table, err := store.Tables().Create(
		in.TableName, in.KeySchema, in.AttributeDefs, in.BillingMode,
		in.ProvThroughput, in.GSI, in.LSI, nil, nil, false,
	)
	if err != nil {
		return nil, err
	}

	imp, err := store.Imports().Create(table.ARN, in.TableName)
	if err != nil {
		if delErr := store.Tables().Delete(in.TableName); delErr != nil {
			logs.Error("Failed to rollback table creation during import",
				logs.Err(delErr),
				logs.String("tableName", in.TableName))
		}
		return nil, err
	}

	importedCount := int64(0)
	processedSizeBytes := int64(0)
	importFailed := false
	var failureCode, failureMessage string

	s3 := s.s3invoker()
	if s3 != nil {
		keys, listErr := s3.ListObjects(ctx, region, in.S3Bucket, in.S3Prefix, 0)
		if listErr != nil {
			logs.Error("Failed to list S3 objects for import",
				logs.Err(listErr),
				logs.String("bucket", in.S3Bucket),
				logs.String("prefix", in.S3Prefix))
			importFailed = true
			failureCode = "S3AccessDenied"
			failureMessage = fmt.Sprintf("failed to list S3 objects: %v", listErr)
		} else {
			for _, key := range keys {
				data, getErr := s3.GetObject(ctx, region, in.S3Bucket, key, 0)
				if getErr != nil {
					logs.Error("Failed to read S3 object for import",
						logs.Err(getErr),
						logs.String("bucket", in.S3Bucket),
						logs.String("key", key))
					importFailed = true
					failureCode = "S3AccessDenied"
					failureMessage = fmt.Sprintf("failed to read S3 object %s: %v", key, getErr)
					continue
				}
				processedSizeBytes += int64(len(data))
				var count int64
				var parseErr error
				switch in.InputFormat {
				case "CSV":
					count, parseErr = importCSVData(ctx, data, in.TableName, store, in.CSVDelimiter, in.CSVHeaderList)
				case "ION":
					importFailed = true
					failureCode = "UnsupportedFormat"
					failureMessage = "ION format parsing is not yet implemented"
					continue
				default:
					count, parseErr = importDynamoDBJSONData(ctx, data, in.TableName, store)
				}
				if parseErr != nil {
					logs.Error("Failed to parse import data",
						logs.Err(parseErr),
						logs.String("key", key))
					importFailed = true
					failureCode = "DocumentAccessException"
					failureMessage = fmt.Sprintf("failed to parse data in %s: %v", key, parseErr)
					continue
				}
				importedCount += count
			}
		}
	}

	if importFailed {
		imp.ImportStatus = "FAILED"
		imp.FailureCode = failureCode
		imp.FailureMessage = failureMessage
	} else {
		imp.ImportStatus = "COMPLETED"
	}
	imp.TableArn = table.ARN
	imp.InputFormat = in.InputFormat
	imp.S3BucketSource = &dbstore.S3BucketSource{
		S3Bucket:      in.S3Bucket,
		S3Prefix:      in.S3Prefix,
		S3BucketOwner: in.S3BucketOwner,
	}
	imp.ProcessedItemCount = importedCount
	imp.ProcessedSizeBytes = processedSizeBytes
	imp.EndTime = time.Now()
	if err := store.Imports().Put(imp); err != nil {
		if delErr := store.Tables().Delete(in.TableName); delErr != nil {
			logs.Error("Failed to rollback table creation during import Put",
				logs.Err(delErr),
				logs.String("tableName", in.TableName))
		}
		return nil, err
	}
	return imp, nil
}
