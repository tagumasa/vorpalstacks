package lambda

import (
	"context"
	"fmt"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/lambda"
	lambdaconnect "vorpalstacks/internal/pb/aws/lambda/lambdaconnect"
	storecommon "vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// AdminHandler provides Lambda service administration functionality.
// It implements the LambdaServiceHandler interface for gRPC-Web communication.
// It delegates to the shared LambdaService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
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

// getStoreFromHeader extracts the region from request headers and returns a FunctionStore.
func (h *AdminHandler) getStoreFromHeader(header http.Header) (*lambdastore.FunctionStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	return h.service.GetFunctionStoreForRegion(region), nil
}

// ListFunctions lists the Lambda functions in the region.
// It returns all functions with their configurations including runtime, memory size, and timeout.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListFunctions(ctx context.Context, req *connect.Request[pb.ListFunctionsRequest]) (*connect.Response[pb.ListFunctionsResponse], error) {
	functionStore, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	maxItems := int(req.Msg.GetMaxitems())
	if maxItems <= 0 {
		maxItems = 50
	}

	opts := storecommon.ListOptions{
		Marker:   req.Msg.Marker,
		MaxItems: maxItems,
	}
	result, err := functionStore.List(opts)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	functions := make([]*pb.FunctionConfiguration, len(result.Items))
	for i, f := range result.Items {
		functions[i] = functionToProto(f)
	}

	resp := &pb.ListFunctionsResponse{Functions: functions}
	if result.NextMarker != "" {
		resp.Nextmarker = result.NextMarker
	}
	return connect.NewResponse(resp), nil
}

// CreateFunction creates a new Lambda function via the admin console.
// It accepts essential parameters: FunctionName (required), Runtime (required),
// Handler (required), Role (required), plus optional MemorySize, Timeout, Description.
func (h *AdminHandler) CreateFunction(ctx context.Context, req *connect.Request[pb.CreateFunctionRequest]) (*connect.Response[pb.FunctionConfiguration], error) {
	if req.Msg.Functionname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("functionName is required"))
	}
	if req.Msg.Handler == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("handler is required"))
	}
	if req.Msg.Role == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("role is required"))
	}

	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	memorySize := req.Msg.GetMemorysize()
	if memorySize == 0 {
		memorySize = 128
	}
	timeout := req.Msg.GetTimeout()
	if timeout == 0 {
		timeout = 3
	}

	fn, err := store.Create(&lambdastore.Function{
		FunctionName: req.Msg.Functionname,
		Runtime:      lambdastore.Runtime(req.Msg.Runtime.String()),
		Role:         req.Msg.Role,
		Handler:      req.Msg.Handler,
		MemorySize:   memorySize,
		Timeout:      timeout,
		Description:  req.Msg.Description,
		PackageType:  req.Msg.Packagetype.String(),
	})
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Apply tags if provided — mirrors the HTTP API's CreateFunction
	// which accepts Tags alongside the function configuration.
	if len(req.Msg.Tags) > 0 {
		tags := make(map[string]string, len(req.Msg.Tags))
		for k, v := range req.Msg.Tags {
			tags[k] = v
		}
		if err := store.TagStore.Tag(fn.FunctionName, tags); err != nil {
			return nil, svcerrors.StoreErrorToGRPC(err)
		}
	}

	resp := functionToProto(fn)
	return connect.NewResponse(resp), nil
}

// DeleteFunction deletes a Lambda function by name via the admin console.
func (h *AdminHandler) DeleteFunction(ctx context.Context, req *connect.Request[pb.DeleteFunctionRequest]) (*connect.Response[pb.DeleteFunctionResponse], error) {
	if req.Msg.Functionname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("functionName is required"))
	}

	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := store.Delete(req.Msg.Functionname); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteFunctionResponse{Statuscode: proto.Int32(204)}), nil
}

// safeRuntime converts a store Runtime to a proto Runtime, falling back to
// nodejs16x when the value is not present in the proto enum map.
func safeRuntime(v lambdastore.Runtime) pb.Runtime {
	if val, ok := pb.Runtime_value[string(v)]; ok {
		return pb.Runtime(val)
	}
	return pb.Runtime_RUNTIME_NODEJS16X
}

// safeState converts a store State to a proto State, defaulting to Active.
func safeState(v lambdastore.State) pb.State {
	if val, ok := pb.State_value[string(v)]; ok {
		return pb.State(val)
	}
	return pb.State_STATE_ACTIVE
}

// safeStateReasonCode maps a state reason code string to the proto enum,
// falling back to Idle when the value is unrecognised.
func safeStateReasonCode(v string) pb.StateReasonCode {
	if val, ok := pb.StateReasonCode_value[v]; ok {
		return pb.StateReasonCode(val)
	}
	return pb.StateReasonCode_STATE_REASON_CODE_IDLE
}

// safePackageType maps a package type string to the proto enum, falling back
// to Zip when the value is unrecognised.
func safePackageType(v string) pb.PackageType {
	if val, ok := pb.PackageType_value[v]; ok {
		return pb.PackageType(val)
	}
	return pb.PackageType_PACKAGE_TYPE_ZIP
}

// functionToProto converts a store Function to its proto representation,
// using safe enum mappers to avoid silent zero-value fallbacks.
func functionToProto(f *lambdastore.Function) *pb.FunctionConfiguration {
	pbFn := &pb.FunctionConfiguration{
		Functionname:    f.FunctionName,
		Functionarn:     f.FunctionArn,
		Runtime:         safeRuntime(f.Runtime),
		Role:            f.Role,
		Handler:         f.Handler,
		Codesize:        proto.Int64(f.CodeSize),
		Codesha256:      f.CodeSha256,
		Description:     f.Description,
		Timeout:         proto.Int32(f.Timeout),
		Memorysize:      proto.Int32(f.MemorySize),
		Lastmodified:    f.LastModified.Format(timeutils.ISO8601UTCFormat),
		Revisionid:      f.RevisionId,
		State:           safeState(f.State),
		Statereason:     f.StateReason,
		Statereasoncode: safeStateReasonCode(f.StateReasonCode),
		Packagetype:     safePackageType(f.PackageType),
	}
	if f.EphemeralStorage != nil {
		pbFn.Ephemeralstorage = &pb.EphemeralStorage{Size: f.EphemeralStorage.Size}
	}
	return pbFn
}

// NewConnectHandler creates a gRPC-Web connect handler for the Lambda admin console.
func NewConnectHandler(svc *LambdaService) (string, http.Handler) {
	return lambdaconnect.NewLambdaServiceHandler(NewAdminHandler(svc))
}
