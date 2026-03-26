package timestreamwrite

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/ingest.timestream"
	ingesttimestreamconnect "vorpalstacks/internal/pb/aws/ingest.timestream/ingest_timestreamconnect"
	"vorpalstacks/internal/store/aws/common"
	timestreamstore "vorpalstacks/internal/store/aws/timestream"
)

// AdminHandler provides Timestream Write service administration functionality.
// It implements the TimestreamWriteServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	ingesttimestreamconnect.UnimplementedTimestreamWriteServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
	dataPath       string
}

var _ ingesttimestreamconnect.TimestreamWriteServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Timestream Write AdminHandler.
// It initialises the handler with the provided storage manager, account ID, and data path.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId, dataPath string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
		dataPath:       dataPath,
	}
}

// getStore retrieves the Timestream Write store for the request.
// It extracts the region from the request header and creates a new Store instance.
func (h *AdminHandler) getStore(req *connect.Request[pb.ListDatabasesRequest]) (*timestreamstore.Store, error) {
	region := req.Header().Get("X-Aws-Region")
	if region == "" {
		region = "us-east-1"
	}
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return timestreamstore.NewStore(regionStorage, h.accountId, region), nil
}

// getTableStore retrieves the Timestream Write table store for the request.
// It extracts the region from the request header and creates a new TableStore instance.
func (h *AdminHandler) getTableStore(req *connect.Request[pb.ListTablesRequest]) (*timestreamstore.TableStore, error) {
	region := req.Header().Get("X-Aws-Region")
	if region == "" {
		region = "us-east-1"
	}
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	dbStore := timestreamstore.NewStore(regionStorage, h.accountId, region)
	return timestreamstore.NewTableStore(regionStorage, dbStore, h.accountId, region), nil
}

// ListDatabases lists databases in Timestream Write.
// It supports pagination via the NextToken parameter.
func (h *AdminHandler) ListDatabases(ctx context.Context, req *connect.Request[pb.ListDatabasesRequest]) (*connect.Response[pb.ListDatabasesResponse], error) {
	store, err := h.getStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	limit := int(req.Msg.Maxresults)
	if limit <= 0 {
		limit = 100
	}

	opts := common.ListOptions{
		MaxItems: limit,
		Marker:   req.Msg.Nexttoken,
	}

	result, err := store.ListDatabases(opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var databases []*pb.Database
	for _, db := range result.Items {
		databases = append(databases, &pb.Database{
			Arn:             db.ARN,
			Databasename:    db.DatabaseName,
			Tablecount:      db.TableCount,
			Kmskeyid:        db.KmsKeyId,
			Creationtime:    db.CreationTime.Format("2006-01-02T15:04:05.000Z"),
			Lastupdatedtime: db.LastUpdatedTime.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	return connect.NewResponse(&pb.ListDatabasesResponse{
		Databases: databases,
		Nexttoken: result.NextMarker,
	}), nil
}

// ListTables lists tables in the specified Timestream Write database.
// It supports pagination via the NextToken parameter.
func (h *AdminHandler) ListTables(ctx context.Context, req *connect.Request[pb.ListTablesRequest]) (*connect.Response[pb.ListTablesResponse], error) {
	store, err := h.getTableStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	limit := int(req.Msg.Maxresults)
	if limit <= 0 {
		limit = 100
	}

	opts := common.ListOptions{
		MaxItems: limit,
		Marker:   req.Msg.Nexttoken,
	}

	result, err := store.ListTables(req.Msg.Databasename, opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var tables []*pb.Table
	for _, t := range result.Items {
		tables = append(tables, &pb.Table{
			Arn:             t.ARN,
			Tablename:       t.TableName,
			Databasename:    t.DatabaseName,
			Creationtime:    t.CreationTime.Format("2006-01-02T15:04:05.000Z"),
			Lastupdatedtime: t.LastUpdatedTime.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	return connect.NewResponse(&pb.ListTablesResponse{
		Tables:    tables,
		Nexttoken: result.NextMarker,
	}), nil
}

// CreateDatabase creates a new database in Timestream Write.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDatabase(ctx context.Context, req *connect.Request[pb.CreateDatabaseRequest]) (*connect.Response[pb.CreateDatabaseResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateTable creates a new table in the specified Timestream Write database.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateTable(ctx context.Context, req *connect.Request[pb.CreateTableRequest]) (*connect.Response[pb.CreateTableResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteDatabase deletes the specified database and all its tables.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDatabase(ctx context.Context, req *connect.Request[pb.DeleteDatabaseRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteTable deletes the specified table and its data.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteTable(ctx context.Context, req *connect.Request[pb.DeleteTableRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeDatabase returns detailed information about the specified database.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeDatabase(ctx context.Context, req *connect.Request[pb.DescribeDatabaseRequest]) (*connect.Response[pb.DescribeDatabaseResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeEndpoints returns the service endpoints for Timestream Write.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeEndpoints(ctx context.Context, req *connect.Request[pb.DescribeEndpointsRequest]) (*connect.Response[pb.DescribeEndpointsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeTable returns detailed information about the specified table.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeTable(ctx context.Context, req *connect.Request[pb.DescribeTableRequest]) (*connect.Response[pb.DescribeTableResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListTagsForResource lists tags for a Timestream Write database or table.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceRequest]) (*connect.Response[pb.ListTagsForResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagResource adds tags to a Timestream Write database or table.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[pb.TagResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagResource removes tags from a Timestream Write database or table.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[pb.UntagResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateDatabase updates the specified database.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateDatabase(ctx context.Context, req *connect.Request[pb.UpdateDatabaseRequest]) (*connect.Response[pb.UpdateDatabaseResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateTable updates the specified table's settings.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateTable(ctx context.Context, req *connect.Request[pb.UpdateTableRequest]) (*connect.Response[pb.UpdateTableResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// WriteRecords writes time-series data to the specified table.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) WriteRecords(ctx context.Context, req *connect.Request[pb.WriteRecordsRequest]) (*connect.Response[pb.WriteRecordsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
