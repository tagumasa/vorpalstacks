package lambda

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/lambda"
	lambdaconnect "vorpalstacks/internal/pb/aws/lambda/lambdaconnect"
	storecommon "vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// AdminHandler provides Lambda service administration functionality.
// It implements the LambdaServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	lambdaconnect.UnimplementedLambdaServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ lambdaconnect.LambdaServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Lambda AdminHandler.
// It initialises the handler with the provided storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

// getStoreFromHeader extracts the region from request headers and returns a FunctionStore.
func (h *AdminHandler) getStoreFromHeader(header http.Header) (*lambdastore.FunctionStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return lambdastore.NewFunctionStore(regionStorage, h.accountId, region), nil
}

// ListFunctions lists the Lambda functions in the region.
// It returns all functions with their configurations including runtime, memory size, and timeout.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListFunctions(ctx context.Context, req *connect.Request[pb.ListFunctionsRequest]) (*connect.Response[pb.ListFunctionsResponse], error) {
	functionStore, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	opts := storecommon.ListOptions{
		Marker:   req.Msg.Marker,
		MaxItems: int(req.Msg.Maxitems),
	}
	result, err := functionStore.List(opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	functions := make([]*pb.FunctionConfiguration, len(result.Items))
	for i, f := range result.Items {
		functions[i] = &pb.FunctionConfiguration{
			Functionname:    f.FunctionName,
			Functionarn:     f.FunctionArn,
			Runtime:         pb.Runtime(pb.Runtime_value[string(f.Runtime)]),
			Role:            f.Role,
			Handler:         f.Handler,
			Codesize:        f.CodeSize,
			Codesha256:      f.CodeSha256,
			Description:     f.Description,
			Timeout:         f.Timeout,
			Memorysize:      f.MemorySize,
			Lastmodified:    f.LastModified.Format("2006-01-02T15:04:05.000Z"),
			Revisionid:      f.RevisionId,
			State:           pb.State(pb.State_value[string(f.State)]),
			Statereason:     f.StateReason,
			Statereasoncode: pb.StateReasonCode(pb.StateReasonCode_value[f.StateReasonCode]),
			Packagetype:     pb.PackageType(pb.PackageType_value[f.PackageType]),
		}
		if f.EphemeralStorage != nil {
			functions[i].Ephemeralstorage = &pb.EphemeralStorage{Size: f.EphemeralStorage.Size}
		}
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
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("FunctionName is required"))
	}
	if req.Msg.Role == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("Role is required"))
	}

	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	memorySize := req.Msg.Memorysize
	if memorySize == 0 {
		memorySize = 128
	}
	timeout := req.Msg.Timeout
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
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &pb.FunctionConfiguration{
		Functionname: fn.FunctionName,
		Functionarn:  fn.FunctionArn,
		Runtime:      pb.Runtime(pb.Runtime_value[string(fn.Runtime)]),
		Role:         fn.Role,
		Handler:      fn.Handler,
		Description:  fn.Description,
		Timeout:      fn.Timeout,
		Memorysize:   fn.MemorySize,
		Packagetype:  pb.PackageType(pb.PackageType_value[fn.PackageType]),
	}
	return connect.NewResponse(resp), nil
}

// DeleteFunction deletes a Lambda function by name via the admin console.
func (h *AdminHandler) DeleteFunction(ctx context.Context, req *connect.Request[pb.DeleteFunctionRequest]) (*connect.Response[pb.DeleteFunctionResponse], error) {
	if req.Msg.Functionname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("FunctionName is required"))
	}

	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := store.Delete(req.Msg.Functionname); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.DeleteFunctionResponse{Statuscode: 204}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Lambda admin console.
func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return lambdaconnect.NewLambdaServiceHandler(NewAdminHandler(sm, accountID))
}
