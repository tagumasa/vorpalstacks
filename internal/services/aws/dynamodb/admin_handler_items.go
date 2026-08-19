package dynamodb

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"
	pb "vorpalstacks/internal/pb/aws/dynamodb"
)

// GetItem retrieves a single DynamoDB item by primary key.
func (h *AdminHandler) GetItem(ctx context.Context, req *connect.Request[pb.GetItemInput]) (*connect.Response[pb.GetItemOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}
	if len(req.Msg.GetKey()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("key is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	attrs, err := h.service.adminGetItem(ctx, region, req.Msg.GetTablename(), req.Msg.GetKey())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.GetItemOutput{
		Item: attrs,
	}), nil
}

// Scan returns all items in a DynamoDB table with optional pagination.
func (h *AdminHandler) Scan(ctx context.Context, req *connect.Request[pb.ScanInput]) (*connect.Response[pb.ScanOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	result, err := h.service.adminScan(region, req.Msg.GetTablename(), req.Msg.GetLimit(), req.Msg.GetExclusivestartkey())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	output := &pb.ScanOutput{
		Items:        result.Items,
		Count:        proto.Int32(result.Count),
		Scannedcount: proto.Int32(result.Count),
	}
	if result.LastEvaluatedKey != nil {
		output.Lastevaluatedkey = result.LastEvaluatedKey
	}

	return connect.NewResponse(output), nil
}

// PutItem inserts or replaces a DynamoDB item.
func (h *AdminHandler) PutItem(ctx context.Context, req *connect.Request[pb.PutItemInput]) (*connect.Response[pb.PutItemOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}
	if len(req.Msg.GetItem()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("item is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	attrs, err := h.service.adminPutItem(ctx, region, req.Msg.GetTablename(), req.Msg.GetItem())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.PutItemOutput{
		Attributes: attrs,
	}), nil
}

// DeleteItem removes a DynamoDB item by primary key.
func (h *AdminHandler) DeleteItem(ctx context.Context, req *connect.Request[pb.DeleteItemInput]) (*connect.Response[pb.DeleteItemOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}
	if len(req.Msg.GetKey()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("key is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	if err := h.service.adminDeleteItem(ctx, region, req.Msg.GetTablename(), req.Msg.GetKey()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteItemOutput{}), nil
}
