// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"fmt"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// DynamoDBService provides DynamoDB operations for managing tables, items, and other resources.
type DynamoDBService struct {
	accountID string
}

// NewDynamoDBService creates a new DynamoDB service instance.
func NewDynamoDBService(accountID string) *DynamoDBService {
	return &DynamoDBService{
		accountID: accountID,
	}
}

func (s *DynamoDBService) store(reqCtx *request.RequestContext) (dynamodbstore.DynamoDBStoreInterface, error) {
	store := reqCtx.GetDynamoDBStore()
	if store != nil {
		return store, nil
	}
	region := reqCtx.GetRegion()
	basicStorage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage: %w", err)
	}
	txnStorage, ok := basicStorage.(storage.TransactionalStorageWith2PC)
	if !ok {
		return nil, fmt.Errorf("storage does not implement TransactionalStorageWith2PC")
	}
	return dynamodbstore.NewDynamoDBStore(txnStorage, s.accountID, region), nil
}

// RegisterHandlers registers all DynamoDB operation handlers with the dispatcher.
func (s *DynamoDBService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("dynamodb", "CreateTable", s.CreateTable)
	d.RegisterHandlerForService("dynamodb", "DeleteTable", s.DeleteTable)
	d.RegisterHandlerForService("dynamodb", "DescribeTable", s.DescribeTable)
	d.RegisterHandlerForService("dynamodb", "ListTables", s.ListTables)
	d.RegisterHandlerForService("dynamodb", "UpdateTable", s.UpdateTable)
	d.RegisterHandlerForService("dynamodb", "PutItem", s.PutItem)
	d.RegisterHandlerForService("dynamodb", "GetItem", s.GetItem)
	d.RegisterHandlerForService("dynamodb", "DeleteItem", s.DeleteItem)
	d.RegisterHandlerForService("dynamodb", "UpdateItem", s.UpdateItem)
	d.RegisterHandlerForService("dynamodb", "Query", s.Query)
	d.RegisterHandlerForService("dynamodb", "Scan", s.Scan)
	d.RegisterHandlerForService("dynamodb", "BatchGetItem", s.BatchGetItem)
	d.RegisterHandlerForService("dynamodb", "BatchWriteItem", s.BatchWriteItem)
	d.RegisterHandlerForService("dynamodb", "TagResource", s.TagResource)
	d.RegisterHandlerForService("dynamodb", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("dynamodb", "ListTagsOfResource", s.ListTagsForResource)

	d.RegisterHandlerForService("dynamodb", "TransactGetItems", s.TransactGetItems)
	d.RegisterHandlerForService("dynamodb", "TransactWriteItems", s.TransactWriteItems)

	d.RegisterHandlerForService("dynamodb", "DescribeTimeToLive", s.DescribeTimeToLive)
	d.RegisterHandlerForService("dynamodb", "UpdateTimeToLive", s.UpdateTimeToLive)

	d.RegisterHandlerForService("dynamodb", "DescribeEndpoints", s.DescribeEndpoints)
	d.RegisterHandlerForService("dynamodb", "DescribeLimits", s.DescribeLimits)

	d.RegisterHandlerForService("dynamodb", "CreateBackup", s.CreateBackup)
	d.RegisterHandlerForService("dynamodb", "DeleteBackup", s.DeleteBackup)
	d.RegisterHandlerForService("dynamodb", "DescribeBackup", s.DescribeBackup)
	d.RegisterHandlerForService("dynamodb", "ListBackups", s.ListBackups)
	d.RegisterHandlerForService("dynamodb", "RestoreTableFromBackup", s.RestoreTableFromBackup)
	d.RegisterHandlerForService("dynamodb", "RestoreTableToPointInTime", s.RestoreTableToPointInTime)

	d.RegisterHandlerForService("dynamodb", "CreateGlobalTable", s.CreateGlobalTable)
	d.RegisterHandlerForService("dynamodb", "DescribeGlobalTable", s.DescribeGlobalTable)
	d.RegisterHandlerForService("dynamodb", "DescribeGlobalTableSettings", s.DescribeGlobalTableSettings)
	d.RegisterHandlerForService("dynamodb", "ListGlobalTables", s.ListGlobalTables)
	d.RegisterHandlerForService("dynamodb", "UpdateGlobalTable", s.UpdateGlobalTable)
	d.RegisterHandlerForService("dynamodb", "UpdateGlobalTableSettings", s.UpdateGlobalTableSettings)

	d.RegisterHandlerForService("dynamodb", "DeleteResourcePolicy", s.DeleteResourcePolicy)
	d.RegisterHandlerForService("dynamodb", "GetResourcePolicy", s.GetResourcePolicy)
	d.RegisterHandlerForService("dynamodb", "PutResourcePolicy", s.PutResourcePolicy)

	d.RegisterHandlerForService("dynamodb", "DescribeContinuousBackups", s.DescribeContinuousBackups)
	d.RegisterHandlerForService("dynamodb", "UpdateContinuousBackups", s.UpdateContinuousBackups)

	d.RegisterHandlerForService("dynamodb", "DescribeKinesisStreamingDestination", s.DescribeKinesisStreamingDestination)
	d.RegisterHandlerForService("dynamodb", "EnableKinesisStreamingDestination", s.EnableKinesisStreamingDestination)
	d.RegisterHandlerForService("dynamodb", "DisableKinesisStreamingDestination", s.DisableKinesisStreamingDestination)
	d.RegisterHandlerForService("dynamodb", "UpdateKinesisStreamingDestination", s.UpdateKinesisStreamingDestination)

	d.RegisterHandlerForService("dynamodb", "ExportTableToPointInTime", s.ExportTableToPointInTime)
	d.RegisterHandlerForService("dynamodb", "DescribeExport", s.DescribeExport)
	d.RegisterHandlerForService("dynamodb", "ListExports", s.ListExports)
	d.RegisterHandlerForService("dynamodb", "ImportTable", s.ImportTable)
	d.RegisterHandlerForService("dynamodb", "DescribeImport", s.DescribeImport)
	d.RegisterHandlerForService("dynamodb", "ListImports", s.ListImports)

	d.RegisterHandlerForService("dynamodb", "DescribeContributorInsights", s.DescribeContributorInsights)
	d.RegisterHandlerForService("dynamodb", "ListContributorInsights", s.ListContributorInsights)
	d.RegisterHandlerForService("dynamodb", "UpdateContributorInsights", s.UpdateContributorInsights)

	d.RegisterHandlerForService("dynamodb", "DescribeTableReplicaAutoScaling", s.DescribeTableReplicaAutoScaling)
	d.RegisterHandlerForService("dynamodb", "UpdateTableReplicaAutoScaling", s.UpdateTableReplicaAutoScaling)

	d.RegisterHandlerForService("dynamodb", "ExecuteStatement", s.ExecuteStatement)
	d.RegisterHandlerForService("dynamodb", "ExecuteTransaction", s.ExecuteTransaction)
	d.RegisterHandlerForService("dynamodb", "BatchExecuteStatement", s.BatchExecuteStatement)
}
