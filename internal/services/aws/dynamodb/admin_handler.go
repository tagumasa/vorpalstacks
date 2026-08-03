package dynamodb

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/dynamodb"
	dynamodbconnect "vorpalstacks/internal/pb/aws/dynamodb/dynamodbconnect"
)

// AdminHandler implements the gRPC admin console handlers for DynamoDB.
// It delegates exclusively to admin-facing service methods, keeping the
// handler free of any store-layer dependency — following the same pattern
// established by ACM's admin handler.
type AdminHandler struct {
	dynamodbconnect.UnimplementedDynamoDBServiceHandler
	service *DynamoDBService
}

var _ dynamodbconnect.DynamoDBServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new DynamoDB admin handler backed by the given
// service instance.
func NewAdminHandler(svc *DynamoDBService) *AdminHandler {
	return &AdminHandler{
		service: svc,
	}
}

// ListTables returns all DynamoDB table names for the admin console.
func (h *AdminHandler) ListTables(ctx context.Context, req *connect.Request[pb.ListTablesInput]) (*connect.Response[pb.ListTablesOutput], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	limit := req.Msg.GetLimit()

	names, nextMarker, err := h.service.adminListTables(region, req.Msg.GetExclusivestarttablename(), limit)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.ListTablesOutput{
		Tablenames:             names,
		Lastevaluatedtablename: nextMarker,
	}), nil
}

// DescribeTable returns detailed metadata for a single DynamoDB table.
func (h *AdminHandler) DescribeTable(ctx context.Context, req *connect.Request[pb.DescribeTableInput]) (*connect.Response[pb.DescribeTableOutput], error) {
	region := svccommon.GetRegionFromHeader(req.Header())

	desc, err := h.service.adminDescribeTable(region, req.Msg.GetTablename())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DescribeTableOutput{
		Table: desc,
	}), nil
}

// CreateTable creates a new DynamoDB table from the admin console request.
func (h *AdminHandler) CreateTable(ctx context.Context, req *connect.Request[pb.CreateTableInput]) (*connect.Response[pb.CreateTableOutput], error) {
	region := svccommon.GetRegionFromHeader(req.Header())

	desc, err := h.service.adminCreateTable(region, req.Msg)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateTableOutput{
		Tabledescription: desc,
	}), nil
}

// DeleteTable removes a DynamoDB table via the admin console.
func (h *AdminHandler) DeleteTable(ctx context.Context, req *connect.Request[pb.DeleteTableInput]) (*connect.Response[pb.DeleteTableOutput], error) {
	region := svccommon.GetRegionFromHeader(req.Header())

	desc, err := h.service.adminDeleteTable(ctx, region, req.Msg.GetTablename())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteTableOutput{
		Tabledescription: desc,
	}), nil
}

// NewConnectHandler returns the connect RPC path and handler for DynamoDB admin.
func NewConnectHandler(svc *DynamoDBService) (string, http.Handler) {
	return dynamodbconnect.NewDynamoDBServiceHandler(NewAdminHandler(svc))
}
