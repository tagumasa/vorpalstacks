package timestreamwrite

import (
	"context"
	"fmt"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/timestreamwrite"
	timestreamwriteconnect "vorpalstacks/internal/pb/aws/timestreamwrite/timestreamwriteconnect"
	storecommon "vorpalstacks/internal/store/aws/common"
	timestreamstore "vorpalstacks/internal/store/aws/timestream"
)

// AdminHandler implements the Timestream Write gRPC-Web admin console handler.
// It exposes list operations for databases and tables for the Flutter
// management UI.
// It delegates to the shared TimestreamWriteService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
type AdminHandler struct {
	timestreamwriteconnect.UnimplementedTimestreamWriteServiceHandler
	service *TimestreamWriteService
}

var _ timestreamwriteconnect.TimestreamWriteServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Timestream Write admin console handler backed
// by the given service instance.
func NewAdminHandler(svc *TimestreamWriteService) *AdminHandler {
	return &AdminHandler{
		service: svc,
	}
}

func (h *AdminHandler) getStoreFromHeader(header http.Header) (*timestreamstore.Store, error) {
	region := svccommon.GetRegionFromHeader(header)
	return h.service.GetDatabaseStoreForRegion(region)
}

func (h *AdminHandler) getTableStoreFromHeader(header http.Header) (*timestreamstore.TableStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	return h.service.GetTableStoreForRegion(region)
}

// ListDatabases returns a paginated list of Timestream databases in the
// requested region.
func (h *AdminHandler) ListDatabases(ctx context.Context, req *connect.Request[pb.ListDatabasesRequest]) (*connect.Response[pb.ListDatabasesResponse], error) {
	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var databases []*pb.Database
	for _, db := range result.Items {
		databases = append(databases, &pb.Database{
			Arn:             db.ARN,
			Databasename:    db.DatabaseName,
			Tablecount:      db.TableCount,
			Kmskeyid:        db.KmsKeyId,
			Creationtime:    db.CreationTime.Format(timeutils.ISO8601UTCFormat),
			Lastupdatedtime: db.LastUpdatedTime.Format(timeutils.ISO8601UTCFormat),
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
	store, err := h.getTableStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var tables []*pb.Table
	for _, t := range result.Items {
		tables = append(tables, &pb.Table{
			Arn:             t.ARN,
			Tablename:       t.TableName,
			Databasename:    t.DatabaseName,
			Creationtime:    t.CreationTime.Format(timeutils.ISO8601UTCFormat),
			Lastupdatedtime: t.LastUpdatedTime.Format(timeutils.ISO8601UTCFormat),
		})
	}

	return connect.NewResponse(&pb.ListTablesResponse{
		Tables:    tables,
		Nexttoken: result.NextMarker,
	}), nil
}

// CreateDatabase creates a new Timestream database via the admin console.
func (h *AdminHandler) CreateDatabase(ctx context.Context, req *connect.Request[pb.CreateDatabaseRequest]) (*connect.Response[pb.CreateDatabaseResponse], error) {
	if req.Msg.GetDatabasename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("DatabaseName is required"))
	}

	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	db, err := store.CreateDatabase(req.Msg.GetDatabasename(), req.Msg.GetKmskeyid())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateDatabaseResponse{
		Database: &pb.Database{
			Arn:             db.ARN,
			Databasename:    db.DatabaseName,
			Kmskeyid:        db.KmsKeyId,
			Creationtime:    db.CreationTime.Format(timeutils.ISO8601UTCFormat),
			Lastupdatedtime: db.LastUpdatedTime.Format(timeutils.ISO8601UTCFormat),
		},
	}), nil
}

// DeleteDatabase deletes a Timestream database via the admin console.
func (h *AdminHandler) DeleteDatabase(ctx context.Context, req *connect.Request[pb.DeleteDatabaseRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.GetDatabasename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("DatabaseName is required"))
	}

	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := store.DeleteDatabase(req.Msg.GetDatabasename()); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Timestream Write admin console.
func NewConnectHandler(svc *TimestreamWriteService) (string, http.Handler) {
	return timestreamwriteconnect.NewTimestreamWriteServiceHandler(NewAdminHandler(svc))
}
