package timestreamwrite

import (
	"context"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/ingest.timestream"
	ingesttimestreamconnect "vorpalstacks/internal/pb/aws/ingest.timestream/ingest_timestreamconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	storecommon "vorpalstacks/internal/store/aws/common"
	timestreamstore "vorpalstacks/internal/store/aws/timestream"
)

type AdminHandler struct {
	ingesttimestreamconnect.UnimplementedTimestreamWriteServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
	dataPath       string
}

var _ ingesttimestreamconnect.TimestreamWriteServiceHandler = (*AdminHandler)(nil)

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
