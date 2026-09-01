package neptunegraph

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/request"

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

// ExecuteQuery runs a query against the specified graph's engine via the
// shared execution Core. The console request is rendered onto the same
// query-body JSON contract the HTTP plane sends.
func (h *AdminHandler) ExecuteQuery(ctx context.Context, req *connect.Request[pb.ExecuteQueryInput]) (*connect.Response[pb.ExecuteQueryOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	body := map[string]interface{}{
		"query":    req.Msg.Querystring,
		"language": "OPEN_CYPHER",
	}
	if len(req.Msg.Parameters) > 0 {
		params := make(map[string]interface{}, len(req.Msg.Parameters))
		for k, v := range req.Msg.Parameters {
			params[k] = v
		}
		body["parameters"] = params
	}
	if mode := explainModeInputPbToString(req.Msg.Explainmode); mode != "" {
		body["explain"] = mode
	}
	if pc := planCacheInputPbToString(req.Msg.Plancache); pc != "" {
		body["planCache"] = pc
	}
	if req.Msg.Querytimeoutmilliseconds != "" {
		ms, err := strconv.Atoi(req.Msg.Querytimeoutmilliseconds)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("queryTimeoutMilliseconds must be an integer"))
		}
		body["queryTimeoutMilliseconds"] = ms
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	result, err := h.service.executeQueryCore(ctx, store, nil, &request.ParsedRequest{
		Body:       raw,
		Parameters: map[string]interface{}{"graphidentifier": req.Msg.Graphidentifier},
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	payload, err := json.Marshal(result)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&pb.ExecuteQueryOutput{Payload: payload}), nil
}

func (h *AdminHandler) CancelQuery(ctx context.Context, req *connect.Request[pb.CancelQueryInput]) (*connect.Response[pbcommon.Empty], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if _, err := h.service.cancelQueryCore(store, &CancelQueryInput{
		GraphIdentifier: req.Msg.Graphidentifier,
		QueryID:         req.Msg.Queryid,
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

func (h *AdminHandler) GetQuery(ctx context.Context, req *connect.Request[pb.GetQueryInput]) (*connect.Response[pb.GetQueryOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	query, err := h.service.getQueryCore(store, &GetQueryInput{
		QueryID:         req.Msg.Queryid,
		GraphIdentifier: req.Msg.Graphidentifier,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(queryToGetPb(query)), nil
}

func (h *AdminHandler) ListQueries(ctx context.Context, req *connect.Request[pb.ListQueriesInput]) (*connect.Response[pb.ListQueriesOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	maxResults := 0
	if req.Msg.Maxresults != "" {
		if n, err := strconv.Atoi(req.Msg.Maxresults); err == nil {
			maxResults = n
		}
	}
	queries, err := h.service.listQueriesCore(store, &ListQueriesInput{
		GraphIdentifier: req.Msg.Graphidentifier,
		MaxResults:      maxResults,
		State:           queryStateInputPbToString(req.Msg.State),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	summaries := make([]*pb.QuerySummary, 0, len(queries))
	for _, q := range queries {
		summaries = append(summaries, queryToSummaryPb(q))
	}
	return connect.NewResponse(&pb.ListQueriesOutput{Queries: summaries}), nil
}

func (h *AdminHandler) GetGraphSummary(ctx context.Context, req *connect.Request[pb.GetGraphSummaryInput]) (*connect.Response[pb.GetGraphSummaryOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.getGraphSummaryCore(store, &GetGraphSummaryInput{
		GraphIdentifier: req.Msg.Graphidentifier,
		Mode:            graphSummaryModePbToString(req.Msg.Mode),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(graphSummaryOutputPb(result)), nil
}

func (h *AdminHandler) CreateGraph(ctx context.Context, req *connect.Request[pb.CreateGraphInput]) (*connect.Response[pb.CreateGraphOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	in := &CreateGraphInput{
		GraphName:          req.Msg.Graphname,
		ProvisionedMemory:  int(req.Msg.Provisionedmemory),
		KmsKeyIdentifier:   req.Msg.Kmskeyidentifier,
		DeletionProtection: strToBool(req.Msg.Deletionprotection),
		PublicConnectivity: strToBool(req.Msg.Publicconnectivity),
		Tags:               req.Msg.Tags,
		Region:             defaults.GetRegionFromHeader(req.Header()),
	}
	if req.Msg.Replicacount != nil {
		in.HasReplicaCount = true
		in.ReplicaCount = int(*req.Msg.Replicacount)
	}
	if req.Msg.Vectorsearchconfiguration != nil {
		in.VectorSearchConfig = map[string]interface{}{"dimension": req.Msg.Vectorsearchconfiguration.GetDimension()}
	}
	graph, err := h.service.createGraphCore(ctx, store, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(graphToCreateGraphPb(graph)), nil
}

func (h *AdminHandler) DeleteGraph(ctx context.Context, req *connect.Request[pb.DeleteGraphInput]) (*connect.Response[pb.DeleteGraphOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	graph, err := h.service.deleteGraphCore(store, &DeleteGraphInput{
		GraphIdentifier: req.Msg.Graphidentifier,
		HasSkipSnapshot: req.Msg.Skipsnapshot != "",
		SkipSnapshot:    strToBool(req.Msg.Skipsnapshot),
		Region:          defaults.GetRegionFromHeader(req.Header()),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(graphToDeleteGraphPb(graph)), nil
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
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	in := &UpdateGraphInput{GraphIdentifier: req.Msg.Graphidentifier}
	if req.Msg.Provisionedmemory != nil {
		in.HasProvisionedMemory = true
		in.ProvisionedMemory = int(*req.Msg.Provisionedmemory)
	}
	if req.Msg.Deletionprotection != "" {
		in.HasDeletionProtection = true
		in.DeletionProtection = strToBool(req.Msg.Deletionprotection)
	}
	if req.Msg.Publicconnectivity != "" {
		in.HasPublicConnectivity = true
		in.PublicConnectivity = strToBool(req.Msg.Publicconnectivity)
	}
	graph, err := h.service.updateGraphCore(store, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(graphToUpdateGraphPb(graph)), nil
}

func (h *AdminHandler) StartGraph(ctx context.Context, req *connect.Request[pb.StartGraphInput]) (*connect.Response[pb.StartGraphOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	graph, err := h.service.startGraphCore(store, &StartGraphInput{
		GraphIdentifier: req.Msg.Graphidentifier,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(graphToStartGraphPb(graph)), nil
}

func (h *AdminHandler) StopGraph(ctx context.Context, req *connect.Request[pb.StopGraphInput]) (*connect.Response[pb.StopGraphOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	graph, err := h.service.stopGraphCore(store, &StopGraphInput{
		GraphIdentifier: req.Msg.Graphidentifier,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(graphToStopGraphPb(graph)), nil
}

func (h *AdminHandler) ResetGraph(ctx context.Context, req *connect.Request[pb.ResetGraphInput]) (*connect.Response[pb.ResetGraphOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	graph, err := h.service.resetGraphCore(store, &ResetGraphInput{
		GraphIdentifier: req.Msg.Graphidentifier,
		HasSkipSnapshot: req.Msg.Skipsnapshot != "",
		SkipSnapshot:    strToBool(req.Msg.Skipsnapshot),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(graphToResetGraphPb(graph)), nil
}

func (h *AdminHandler) RestoreGraphFromSnapshot(ctx context.Context, req *connect.Request[pb.RestoreGraphFromSnapshotInput]) (*connect.Response[pb.RestoreGraphFromSnapshotOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	in := &RestoreGraphFromSnapshotInput{
		SnapshotIdentifier: req.Msg.Snapshotidentifier,
		GraphName:          req.Msg.Graphname,
		DeletionProtection: strToBool(req.Msg.Deletionprotection),
		PublicConnectivity: strToBool(req.Msg.Publicconnectivity),
		Tags:               req.Msg.Tags,
		Region:             defaults.GetRegionFromHeader(req.Header()),
	}
	if req.Msg.Provisionedmemory != nil {
		in.HasProvisionedMemory = true
		in.ProvisionedMemory = int(*req.Msg.Provisionedmemory)
	}
	if req.Msg.Replicacount != nil {
		in.HasReplicaCount = true
		in.ReplicaCount = int(*req.Msg.Replicacount)
	}
	graph, err := h.service.restoreGraphFromSnapshotCore(store, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(graphToRestoreGraphFromSnapshotPb(graph)), nil
}

func (h *AdminHandler) CreateGraphSnapshot(ctx context.Context, req *connect.Request[pb.CreateGraphSnapshotInput]) (*connect.Response[pb.CreateGraphSnapshotOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	snapshot, err := h.service.createGraphSnapshotCore(store, &CreateGraphSnapshotInput{
		GraphIdentifier: req.Msg.Graphidentifier,
		SnapshotName:    req.Msg.Snapshotname,
		Region:          defaults.GetRegionFromHeader(req.Header()),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(snapshotToCreateGraphSnapshotPb(snapshot)), nil
}

func (h *AdminHandler) DeleteGraphSnapshot(ctx context.Context, req *connect.Request[pb.DeleteGraphSnapshotInput]) (*connect.Response[pb.DeleteGraphSnapshotOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	snapshot, err := h.service.deleteGraphSnapshotCore(store, &DeleteGraphSnapshotInput{
		SnapshotIdentifier: req.Msg.Snapshotidentifier,
		Region:             defaults.GetRegionFromHeader(req.Header()),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(snapshotToDeleteGraphSnapshotPb(snapshot)), nil
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
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	ep, err := h.service.createPrivateGraphEndpointCore(ctx, store, &CreatePrivateGraphEndpointInput{
		GraphIdentifier: req.Msg.Graphidentifier,
		VpcId:           req.Msg.Vpcid,
		SubnetIds:       req.Msg.Subnetids,
		Region:          defaults.GetRegionFromHeader(req.Header()),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(endpointToCreatePrivateGraphEndpointPb(ep)), nil
}

func (h *AdminHandler) DeletePrivateGraphEndpoint(ctx context.Context, req *connect.Request[pb.DeletePrivateGraphEndpointInput]) (*connect.Response[pb.DeletePrivateGraphEndpointOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	ep, err := h.service.deletePrivateGraphEndpointCore(store, &DeletePrivateGraphEndpointInput{
		GraphIdentifier: req.Msg.Graphidentifier,
		VpcId:           req.Msg.Vpcid,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(endpointToDeletePrivateGraphEndpointPb(ep)), nil
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
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	in := &CreateGraphUsingImportTaskInput{
		GraphName:          req.Msg.Graphname,
		RoleArn:            req.Msg.Rolearn,
		Source:             req.Msg.Source,
		Format:             formatPbToString(req.Msg.Format),
		ParquetType:        parquetTypePbToString(req.Msg.Parquettype),
		BlankNodeHandling:  blankNodeHandlingPbToString(req.Msg.Blanknodehandling),
		KmsKeyIdentifier:   req.Msg.Kmskeyidentifier,
		DeletionProtection: strToBool(req.Msg.Deletionprotection),
		PublicConnectivity: strToBool(req.Msg.Publicconnectivity),
		FailOnError:        strToBool(req.Msg.Failonerror),
		Tags:               req.Msg.Tags,
		Region:             defaults.GetRegionFromHeader(req.Header()),
	}
	if req.Msg.Replicacount != nil {
		in.HasReplicaCount = true
		in.ReplicaCount = int(*req.Msg.Replicacount)
	}
	if req.Msg.Minprovisionedmemory != nil {
		in.HasMinProvisionedMemory = true
		in.MinProvisionedMemory = int(*req.Msg.Minprovisionedmemory)
	}
	if req.Msg.Maxprovisionedmemory != nil {
		in.HasMaxProvisionedMemory = true
		in.MaxProvisionedMemory = int(*req.Msg.Maxprovisionedmemory)
	}
	if req.Msg.Vectorsearchconfiguration != nil {
		in.VectorSearchConfig = map[string]interface{}{"dimension": req.Msg.Vectorsearchconfiguration.GetDimension()}
	}
	if req.Msg.Importoptions != nil {
		in.HasImportOptions = true
		in.ImportOptions = importOptionsFromPb(req.Msg.Importoptions)
	}
	task, err := h.service.createGraphUsingImportTaskCore(ctx, store, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(importTaskToCreateGraphUsingImportTaskPb(task)), nil
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
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	task, err := h.service.cancelImportTaskCore(store, &CancelImportTaskInput{
		TaskIdentifier: req.Msg.Taskidentifier,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(importTaskToCancelImportTaskPb(task)), nil
}

func (h *AdminHandler) StartImportTask(ctx context.Context, req *connect.Request[pb.StartImportTaskInput]) (*connect.Response[pb.StartImportTaskOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	in := &StartImportTaskInput{
		GraphIdentifier:   req.Msg.Graphidentifier,
		RoleArn:           req.Msg.Rolearn,
		Source:            req.Msg.Source,
		Format:            formatPbToString(req.Msg.Format),
		ParquetType:       parquetTypePbToString(req.Msg.Parquettype),
		BlankNodeHandling: blankNodeHandlingPbToString(req.Msg.Blanknodehandling),
		FailOnError:       strToBool(req.Msg.Failonerror),
	}
	if req.Msg.Importoptions != nil {
		in.HasImportOptions = true
		in.ImportOptions = importOptionsFromPb(req.Msg.Importoptions)
	}
	task, err := h.service.startImportTaskCore(store, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(importTaskToStartImportTaskPb(task)), nil
}

func (h *AdminHandler) StartExportTask(ctx context.Context, req *connect.Request[pb.StartExportTaskInput]) (*connect.Response[pb.StartExportTaskOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	in := &StartExportTaskInput{
		GraphIdentifier:  req.Msg.Graphidentifier,
		Format:           exportFormatPbToString(req.Msg.Format),
		ParquetType:      parquetTypePbToString(req.Msg.Parquettype),
		KmsKeyIdentifier: req.Msg.Kmskeyidentifier,
		RoleArn:          req.Msg.Rolearn,
		Destination:      req.Msg.Destination,
	}
	if req.Msg.Exportfilter != nil {
		in.HasExportFilter = true
		in.ExportFilter = exportFilterFromPb(req.Msg.Exportfilter)
	}
	task, err := h.service.startExportTaskCore(store, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(exportTaskToStartExportTaskPb(task)), nil
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
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	task, err := h.service.cancelExportTaskCore(store, &CancelExportTaskInput{
		TaskIdentifier: req.Msg.Taskidentifier,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(exportTaskToCancelExportTaskPb(task)), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the NeptuneGraph admin console.
func NewConnectHandler(svc *NeptuneGraphService, accountID string) (string, http.Handler) {
	return neptunegraphconnect.NewNeptuneGraphServiceHandler(NewAdminHandler(svc, accountID))
}
