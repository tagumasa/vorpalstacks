package dynamodb

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/dynamodb"
	dynamodbconnect "vorpalstacks/internal/pb/aws/dynamodb/dynamodbconnect"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// AdminHandler provides DynamoDB service administration functionality.
type AdminHandler struct {
	dynamodbconnect.UnimplementedDynamoDBServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ dynamodbconnect.DynamoDBServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new DynamoDB AdminHandler.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

// getTableStore retrieves the DynamoDB table store for the request.
// It extracts the region from the request header and creates a new TableStore instance.
func (h *AdminHandler) getTableStore(req *connect.Request[pb.ListTablesInput]) (*dynamodbstore.TableStore, error) {
	region := req.Header().Get("X-Aws-Region")
	if region == "" {
		region = "us-east-1"
	}
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return dynamodbstore.NewTableStore(regionStorage, h.accountId, region), nil
}

// ListTables lists tables in DynamoDB.
func (h *AdminHandler) ListTables(ctx context.Context, req *connect.Request[pb.ListTablesInput]) (*connect.Response[pb.ListTablesOutput], error) {
	tableStore, err := h.getTableStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	marker := ""
	limit := 100
	if req.Msg.Limit > 0 {
		limit = int(req.Msg.Limit)
	}

	tables, nextMarker, err := tableStore.List(marker, limit)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var tableNames []string
	for _, t := range tables {
		tableNames = append(tableNames, t.Name)
	}

	return connect.NewResponse(&pb.ListTablesOutput{
		Tablenames:             tableNames,
		Lastevaluatedtablename: nextMarker,
	}), nil
}

// BatchExecuteStatement executes multiple SQL statements in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) BatchExecuteStatement(ctx context.Context, req *connect.Request[pb.BatchExecuteStatementInput]) (*connect.Response[pb.BatchExecuteStatementOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// BatchGetItem retrieves multiple items from one or more tables in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) BatchGetItem(ctx context.Context, req *connect.Request[pb.BatchGetItemInput]) (*connect.Response[pb.BatchGetItemOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// BatchWriteItem puts or deletes multiple items in one or more tables in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) BatchWriteItem(ctx context.Context, req *connect.Request[pb.BatchWriteItemInput]) (*connect.Response[pb.BatchWriteItemOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateTable creates a new table with the specified schema and settings.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateTable(ctx context.Context, req *connect.Request[pb.CreateTableInput]) (*connect.Response[pb.CreateTableOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteTable deletes the specified table and all its items.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteTable(ctx context.Context, req *connect.Request[pb.DeleteTableInput]) (*connect.Response[pb.DeleteTableOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeTable returns information about the specified table.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeTable(ctx context.Context, req *connect.Request[pb.DescribeTableInput]) (*connect.Response[pb.DescribeTableOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetItem retrieves a single item from the specified table by its primary key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetItem(ctx context.Context, req *connect.Request[pb.GetItemInput]) (*connect.Response[pb.GetItemOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutItem creates a new item or replaces an existing item in the specified table.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutItem(ctx context.Context, req *connect.Request[pb.PutItemInput]) (*connect.Response[pb.PutItemOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// Query retrieves items from a table based on the partition key and optional sort key conditions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) Query(ctx context.Context, req *connect.Request[pb.QueryInput]) (*connect.Response[pb.QueryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// Scan retrieves all items from the specified table or a secondary index.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) Scan(ctx context.Context, req *connect.Request[pb.ScanInput]) (*connect.Response[pb.ScanOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateItem edits an existing item's attributes or adds a new item if it does not exist.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateItem(ctx context.Context, req *connect.Request[pb.UpdateItemInput]) (*connect.Response[pb.UpdateItemOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateTable modifies the specified table's settings and schema.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateTable(ctx context.Context, req *connect.Request[pb.UpdateTableInput]) (*connect.Response[pb.UpdateTableOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteItem deletes an item from the specified table by its primary key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteItem(ctx context.Context, req *connect.Request[pb.DeleteItemInput]) (*connect.Response[pb.DeleteItemOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
