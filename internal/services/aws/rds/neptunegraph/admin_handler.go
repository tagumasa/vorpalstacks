package neptunegraph

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"

	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/neptunegraph"
	neptunegraphconnect "vorpalstacks/internal/pb/aws/neptunegraph/neptunegraphconnect"
)

// AdminHandler provides gRPC-based admin console handlers for NeptuneGraph resources.
type AdminHandler struct {
	neptunegraphconnect.UnimplementedNeptuneGraphServiceHandler
	service   *NeptuneGraphService
	accountId string
}

var _ neptunegraphconnect.NeptuneGraphServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates an AdminHandler backed by the given service and account identifier.
func NewAdminHandler(svc *NeptuneGraphService, accountId string) *AdminHandler {
	return &AdminHandler{service: svc, accountId: accountId}
}

// ExecuteQuery handles the admin console ExecuteQuery request. Not implemented via admin console.
func (h *AdminHandler) ExecuteQuery(ctx context.Context, req *connect.Request[pb.ExecuteQueryInput]) (*connect.Response[pb.ExecuteQueryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) CancelQuery(ctx context.Context, req *connect.Request[pb.CancelQueryInput]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) GetQuery(ctx context.Context, req *connect.Request[pb.GetQueryInput]) (*connect.Response[pb.GetQueryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) ListQueries(ctx context.Context, req *connect.Request[pb.ListQueriesInput]) (*connect.Response[pb.ListQueriesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) GetGraphSummary(ctx context.Context, req *connect.Request[pb.GetGraphSummaryInput]) (*connect.Response[pb.GetGraphSummaryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) CreateGraph(ctx context.Context, req *connect.Request[pb.CreateGraphInput]) (*connect.Response[pb.CreateGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) DeleteGraph(ctx context.Context, req *connect.Request[pb.DeleteGraphInput]) (*connect.Response[pb.DeleteGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) GetGraph(ctx context.Context, req *connect.Request[pb.GetGraphInput]) (*connect.Response[pb.GetGraphOutput], error) {
	output, err := h.getGraphPb(req.Header(), req.Msg.Graphidentifier)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(output), nil
}

func (h *AdminHandler) ListGraphs(ctx context.Context, req *connect.Request[pb.ListGraphsInput]) (*connect.Response[pb.ListGraphsOutput], error) {
	summaries, err := h.listGraphsPb(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.ListGraphsOutput{Graphs: summaries}), nil
}

func (h *AdminHandler) UpdateGraph(ctx context.Context, req *connect.Request[pb.UpdateGraphInput]) (*connect.Response[pb.UpdateGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) StartGraph(ctx context.Context, req *connect.Request[pb.StartGraphInput]) (*connect.Response[pb.StartGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) StopGraph(ctx context.Context, req *connect.Request[pb.StopGraphInput]) (*connect.Response[pb.StopGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) ResetGraph(ctx context.Context, req *connect.Request[pb.ResetGraphInput]) (*connect.Response[pb.ResetGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) RestoreGraphFromSnapshot(ctx context.Context, req *connect.Request[pb.RestoreGraphFromSnapshotInput]) (*connect.Response[pb.RestoreGraphFromSnapshotOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) CreateGraphSnapshot(ctx context.Context, req *connect.Request[pb.CreateGraphSnapshotInput]) (*connect.Response[pb.CreateGraphSnapshotOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) DeleteGraphSnapshot(ctx context.Context, req *connect.Request[pb.DeleteGraphSnapshotInput]) (*connect.Response[pb.DeleteGraphSnapshotOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) GetGraphSnapshot(ctx context.Context, req *connect.Request[pb.GetGraphSnapshotInput]) (*connect.Response[pb.GetGraphSnapshotOutput], error) {
	output, err := h.getGraphSnapshotPb(req.Header(), req.Msg.Snapshotidentifier)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(output), nil
}

func (h *AdminHandler) ListGraphSnapshots(ctx context.Context, req *connect.Request[pb.ListGraphSnapshotsInput]) (*connect.Response[pb.ListGraphSnapshotsOutput], error) {
	summaries, err := h.listGraphSnapshotsPb(req.Header(), req.Msg.Graphidentifier)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.ListGraphSnapshotsOutput{Graphsnapshots: summaries}), nil
}

func (h *AdminHandler) CreatePrivateGraphEndpoint(ctx context.Context, req *connect.Request[pb.CreatePrivateGraphEndpointInput]) (*connect.Response[pb.CreatePrivateGraphEndpointOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) DeletePrivateGraphEndpoint(ctx context.Context, req *connect.Request[pb.DeletePrivateGraphEndpointInput]) (*connect.Response[pb.DeletePrivateGraphEndpointOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) GetPrivateGraphEndpoint(ctx context.Context, req *connect.Request[pb.GetPrivateGraphEndpointInput]) (*connect.Response[pb.GetPrivateGraphEndpointOutput], error) {
	output, err := h.getPrivateGraphEndpointPb(req.Header(), req.Msg.Graphidentifier, req.Msg.Vpcid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(output), nil
}

func (h *AdminHandler) ListPrivateGraphEndpoints(ctx context.Context, req *connect.Request[pb.ListPrivateGraphEndpointsInput]) (*connect.Response[pb.ListPrivateGraphEndpointsOutput], error) {
	summaries, err := h.listPrivateGraphEndpointsPb(req.Header(), req.Msg.Graphidentifier)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.ListPrivateGraphEndpointsOutput{Privategraphendpoints: summaries}), nil
}

func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceInput]) (*connect.Response[pb.ListTagsForResourceOutput], error) {
	tags, err := h.listTagsPb(req.Header(), req.Msg.Resourcearn)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.ListTagsForResourceOutput{Tags: tags}), nil
}

func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceInput]) (*connect.Response[pb.TagResourceOutput], error) {
	if err := h.tagResourcePb(req.Header(), req.Msg.Resourcearn, req.Msg.Tags); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.TagResourceOutput{}), nil
}

func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceInput]) (*connect.Response[pb.UntagResourceOutput], error) {
	if err := h.untagResourcePb(req.Header(), req.Msg.Resourcearn, req.Msg.Tagkeys); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.UntagResourceOutput{}), nil
}

func (h *AdminHandler) CreateGraphUsingImportTask(ctx context.Context, req *connect.Request[pb.CreateGraphUsingImportTaskInput]) (*connect.Response[pb.CreateGraphUsingImportTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) GetImportTask(ctx context.Context, req *connect.Request[pb.GetImportTaskInput]) (*connect.Response[pb.GetImportTaskOutput], error) {
	output, err := h.getImportTaskPb(req.Header(), req.Msg.Taskidentifier)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(output), nil
}

func (h *AdminHandler) ListImportTasks(ctx context.Context, req *connect.Request[pb.ListImportTasksInput]) (*connect.Response[pb.ListImportTasksOutput], error) {
	summaries, err := h.listImportTasksPb(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.ListImportTasksOutput{Tasks: summaries}), nil
}

func (h *AdminHandler) CancelImportTask(ctx context.Context, req *connect.Request[pb.CancelImportTaskInput]) (*connect.Response[pb.CancelImportTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) StartImportTask(ctx context.Context, req *connect.Request[pb.StartImportTaskInput]) (*connect.Response[pb.StartImportTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) StartExportTask(ctx context.Context, req *connect.Request[pb.StartExportTaskInput]) (*connect.Response[pb.StartExportTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func (h *AdminHandler) GetExportTask(ctx context.Context, req *connect.Request[pb.GetExportTaskInput]) (*connect.Response[pb.GetExportTaskOutput], error) {
	output, err := h.getExportTaskPb(req.Header(), req.Msg.Taskidentifier)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(output), nil
}

func (h *AdminHandler) ListExportTasks(ctx context.Context, req *connect.Request[pb.ListExportTasksInput]) (*connect.Response[pb.ListExportTasksOutput], error) {
	summaries, err := h.listExportTasksPb(req.Header(), req.Msg.Graphidentifier)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.ListExportTasksOutput{Tasks: summaries}), nil
}

func (h *AdminHandler) CancelExportTask(ctx context.Context, req *connect.Request[pb.CancelExportTaskInput]) (*connect.Response[pb.CancelExportTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func errNotImplemented() error {
	return fmt.Errorf("not implemented via admin console")
}

// NewConnectHandler creates a gRPC-Web connect handler for the NeptuneGraph admin console.
func NewConnectHandler(svc *NeptuneGraphService, accountID string) (string, http.Handler) {
	return neptunegraphconnect.NewNeptuneGraphServiceHandler(NewAdminHandler(svc, accountID))
}
