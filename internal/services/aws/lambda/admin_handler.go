package lambda

import (
	"context"
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

// getFunctionStore returns the Lambda function store for the specified region.
// It extracts the region from the request headers and creates a new function store.
func (h *AdminHandler) getFunctionStore(req *connect.Request[pb.ListFunctionsRequest]) (*lambdastore.FunctionStore, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
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
	functionStore, err := h.getFunctionStore(req)
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

func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return lambdaconnect.NewLambdaServiceHandler(NewAdminHandler(sm, accountID))
}
