package athena

import (
	"context"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/athena"
	athenaconnect "vorpalstacks/internal/pb/aws/athena/athenaconnect"
)

// AdminHandler implements the Athena admin console gRPC-Web handler.
// It delegates to core functions in workgroup_core.go so that validation,
// cascade cleanup, and error handling are shared with the HTTP API path.
type AdminHandler struct {
	athenaconnect.UnimplementedAthenaServiceHandler
	service *AthenaService
}

var _ athenaconnect.AthenaServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(svc *AthenaService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// CreateWorkGroup creates a new Athena work group via the admin console.
func (h *AdminHandler) CreateWorkGroup(ctx context.Context, req *connect.Request[pb.CreateWorkGroupInput]) (*connect.Response[pb.CreateWorkGroupOutput], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := createWorkGroupCore(stores, protoToCreateInput(req.Msg)); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateWorkGroupOutput{}), nil
}

// DeleteWorkGroup deletes an Athena work group via the admin console.
// Delegates to deleteWorkGroupCore for cascade cleanup of dependent resources.
func (h *AdminHandler) DeleteWorkGroup(ctx context.Context, req *connect.Request[pb.DeleteWorkGroupInput]) (*connect.Response[pb.DeleteWorkGroupOutput], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := deleteWorkGroupCore(stores, req.Msg.Workgroup); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteWorkGroupOutput{}), nil
}

// ListWorkGroups retrieves Athena work groups with pagination support.
func (h *AdminHandler) ListWorkGroups(ctx context.Context, req *connect.Request[pb.ListWorkGroupsInput]) (*connect.Response[pb.ListWorkGroupsOutput], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	maxResults := clampMaxResults(int(req.Msg.GetMaxresults()), athenaMaxWorkGroupsResults, athenaMaxWorkGroupsResults)

	result, err := listWorkGroupsCore(stores, maxResults, req.Msg.Nexttoken)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.ListWorkGroupsOutput{
		Workgroups: toPbWorkGroupSummaries(result.Items),
		Nexttoken:  result.NextMarker,
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Athena admin console.
func NewConnectHandler(svc *AthenaService) (string, http.Handler) {
	return athenaconnect.NewAthenaServiceHandler(NewAdminHandler(svc))
}
