package timestreamwrite

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/timestreamwrite"
	timestreamwriteconnect "vorpalstacks/internal/pb/aws/timestreamwrite/timestreamwriteconnect"
	storecommon "vorpalstacks/internal/store/aws/common"
	timestreamstore "vorpalstacks/internal/store/aws/timestream"
)

// AdminHandler implements the Timestream Write gRPC-Web admin console handler.
// It exposes list operations for databases and tables for the Flutter
// management UI.
type AdminHandler struct {
	timestreamwriteconnect.UnimplementedTimestreamWriteServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
	dataPath       string
}

var _ timestreamwriteconnect.TimestreamWriteServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Timestream Write admin console handler.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId, dataPath string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
		dataPath:       dataPath,
	}
}

func (h *AdminHandler) getStore(req *connect.Request[pb.ListDatabasesRequest]) (*timestreamstore.Store, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return timestreamstore.NewStore(regionStorage, h.accountId, region), nil
}

func (h *AdminHandler) getTableStore(req *connect.Request[pb.ListTablesRequest]) (*timestreamstore.TableStore, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	dbStore := timestreamstore.NewStore(regionStorage, h.accountId, region)
	return timestreamstore.NewTableStore(regionStorage, dbStore, h.accountId, region), nil
}

// ListDatabases returns a paginated list of Timestream databases in the
// requested region.
func (h *AdminHandler) ListDatabases(ctx context.Context, req *connect.Request[pb.ListDatabasesRequest]) (*connect.Response[pb.ListDatabasesResponse], error) {
	store, err := h.getStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	limit := int(req.Msg.Maxresults)
	if limit <= 0 {
		limit = 100
	}

	opts := storecommon.ListOptions{
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

// ListTables returns a paginated list of Timestream tables in the specified
// database.
func (h *AdminHandler) ListTables(ctx context.Context, req *connect.Request[pb.ListTablesRequest]) (*connect.Response[pb.ListTablesResponse], error) {
	store, err := h.getTableStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	limit := int(req.Msg.Maxresults)
	if limit <= 0 {
		limit = 100
	}

	opts := storecommon.ListOptions{
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

func NewConnectHandler(sm *storage.RegionStorageManager, accountID, dataPath string) (string, http.Handler) {
	return timestreamwriteconnect.NewTimestreamWriteServiceHandler(NewAdminHandler(sm, accountID, dataPath))
}
