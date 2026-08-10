package timestreamwrite

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/aws/types"

	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/timestreamwrite"
	timestreamwriteconnect "vorpalstacks/internal/pb/aws/timestreamwrite/timestreamwriteconnect"
)

// AdminHandler implements the Timestream Write gRPC-Web admin console handler.
// It is a thin adapter: proto request → transport-agnostic Input struct →
// service-layer Core function → proto response conversion. Store packages
// are never imported directly (AGENTS.md #29 compliance).
type AdminHandler struct {
	timestreamwriteconnect.UnimplementedTimestreamWriteServiceHandler
	service *TimestreamWriteService
}

var _ timestreamwriteconnect.TimestreamWriteServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Timestream Write admin console handler backed
// by the given service instance.
func NewAdminHandler(svc *TimestreamWriteService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListDatabases returns a paginated list of Timestream databases in the
// requested region.
func (h *AdminHandler) ListDatabases(ctx context.Context, req *connect.Request[pb.ListDatabasesRequest]) (*connect.Response[pb.ListDatabasesResponse], error) {
	stores, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.listDatabasesCore(ctx, stores, ListDatabasesInput{
		NextToken: req.Msg.Nexttoken,
		MaxItems:  int(req.Msg.GetMaxresults()),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	databases := make([]*pb.Database, 0, len(result.Databases))
	for _, db := range result.Databases {
		databases = append(databases, toPbDatabase(&db))
	}

	return connect.NewResponse(&pb.ListDatabasesResponse{
		Databases: databases,
		Nexttoken: result.NextToken,
	}), nil
}

// ListTables returns a paginated list of Timestream tables in the specified
// database.
func (h *AdminHandler) ListTables(ctx context.Context, req *connect.Request[pb.ListTablesRequest]) (*connect.Response[pb.ListTablesResponse], error) {
	stores, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.listTablesCore(ctx, stores, ListTablesInput{
		DatabaseName: req.Msg.Databasename,
		NextToken:    req.Msg.Nexttoken,
		MaxItems:     int(req.Msg.GetMaxresults()),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	tables := make([]*pb.Table, 0, len(result.Tables))
	for _, t := range result.Tables {
		tables = append(tables, toPbTable(&t))
	}

	return connect.NewResponse(&pb.ListTablesResponse{
		Tables:    tables,
		Nexttoken: result.NextToken,
	}), nil
}

// CreateDatabase creates a new Timestream database via the admin console.
func (h *AdminHandler) CreateDatabase(ctx context.Context, req *connect.Request[pb.CreateDatabaseRequest]) (*connect.Response[pb.CreateDatabaseResponse], error) {
	stores, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	input := CreateDatabaseInput{
		DatabaseName: req.Msg.GetDatabasename(),
		KmsKeyId:     req.Msg.GetKmskeyid(),
	}

	if len(req.Msg.Tags) > 0 {
		input.TagsProvided = true
		tags := make([]types.Tag, 0, len(req.Msg.Tags))
		for _, t := range req.Msg.Tags {
			tags = append(tags, types.Tag{Key: t.Key, Value: t.Value})
		}
		input.Tags = tags
	}

	result, err := h.service.createDatabaseCore(ctx, stores, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateDatabaseResponse{
		Database: toPbDatabase(result),
	}), nil
}

// DeleteDatabase deletes a Timestream database via the admin console.
func (h *AdminHandler) DeleteDatabase(ctx context.Context, req *connect.Request[pb.DeleteDatabaseRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.deleteDatabaseCore(ctx, stores, DeleteDatabaseInput{
		DatabaseName: req.Msg.GetDatabasename(),
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Timestream Write admin console.
func NewConnectHandler(svc *TimestreamWriteService) (string, http.Handler) {
	return timestreamwriteconnect.NewTimestreamWriteServiceHandler(NewAdminHandler(svc))
}
