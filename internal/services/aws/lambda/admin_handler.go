package lambda

import (
	"context"
	"fmt"
	"net/http"

	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	pb "vorpalstacks/internal/pb/aws/lambda"
	lambdaconnect "vorpalstacks/internal/pb/aws/lambda/lambdaconnect"
)

// AdminHandler provides Lambda service administration functionality.
// It implements the LambdaServiceHandler interface for gRPC-Web communication.
// All operations delegate to service-layer core functions, ensuring that
// validation and persistence follow the same single code path as the HTTP
// API handlers.
type AdminHandler struct {
	lambdaconnect.UnimplementedLambdaServiceHandler
	service *LambdaService
}

var _ lambdaconnect.LambdaServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Lambda AdminHandler backed by the given
// service instance, ensuring the same per-region cached stores are used as
// the HTTP API handlers.
func NewAdminHandler(svc *LambdaService) *AdminHandler {
	return &AdminHandler{
		service: svc,
	}
}

// ListFunctions lists the Lambda functions in the region.
func (h *AdminHandler) ListFunctions(ctx context.Context, req *connect.Request[pb.ListFunctionsRequest]) (*connect.Response[pb.ListFunctionsResponse], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	items, nextMarker, err := h.service.listFunctionsCore(stores, &ListFunctionsInput{
		Marker:   req.Msg.Marker,
		MaxItems: int(req.Msg.GetMaxitems()),
	})
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	functions := make([]*pb.FunctionConfiguration, len(items))
	for i, f := range items {
		functions[i] = functionToProto(f)
	}

	resp := &pb.ListFunctionsResponse{Functions: functions}
	if nextMarker != "" {
		resp.Nextmarker = nextMarker
	}
	return connect.NewResponse(resp), nil
}

// CreateFunction creates a new Lambda function via the admin console.
// It delegates to createFunctionCore so that all field validation
// (FunctionName, Runtime, Handler, Role, etc.) runs identically to the
// HTTP API path.
func (h *AdminHandler) CreateFunction(ctx context.Context, req *connect.Request[pb.CreateFunctionRequest]) (*connect.Response[pb.FunctionConfiguration], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	in := &CreateFunctionInput{
		FunctionName: req.Msg.Functionname,
		Runtime:      protoToStoreRuntime(req.Msg.Runtime),
		Role:         req.Msg.Role,
		Handler:      req.Msg.Handler,
		Description:  req.Msg.Description,
		PackageType:  protoToPackageType(req.Msg.Packagetype),
		MemorySize:   req.Msg.GetMemorysize(),
		Timeout:      req.Msg.GetTimeout(),
		Region:       svccommon.GetRegionFromHeader(req.Header()),
	}

	if len(req.Msg.Tags) > 0 {
		in.Tags = make(map[string]string, len(req.Msg.Tags))
		for k, v := range req.Msg.Tags {
			in.Tags[k] = v
		}
	}

	created, _, err := h.service.createFunctionCore(stores, in)
	if err != nil {
		if lambdaErr, ok := err.(*LambdaError); ok {
			return nil, svcerrors.AWSErrorToGRPC(lambdaErr.AWSError)
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(functionToProto(created)), nil
}

// DeleteFunction deletes a Lambda function by name via the admin console.
// It delegates to deleteFunctionCore so that event-source-mapping existence
// checks and container cleanup run identically to the HTTP API path.
func (h *AdminHandler) DeleteFunction(ctx context.Context, req *connect.Request[pb.DeleteFunctionRequest]) (*connect.Response[pb.DeleteFunctionResponse], error) {
	if req.Msg.Functionname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("functionName is required"))
	}

	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := h.service.deleteFunctionCore(ctx, stores, &DeleteFunctionInput{
		FunctionName: req.Msg.Functionname,
	}); err != nil {
		if lambdaErr, ok := err.(*LambdaError); ok {
			return nil, svcerrors.AWSErrorToGRPC(lambdaErr.AWSError)
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteFunctionResponse{Statuscode: proto.Int32(204)}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Lambda admin console.
func NewConnectHandler(svc *LambdaService) (string, http.Handler) {
	return lambdaconnect.NewLambdaServiceHandler(NewAdminHandler(svc))
}
