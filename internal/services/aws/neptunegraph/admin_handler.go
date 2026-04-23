// Package neptunegraph implements AWS NeptuneGraph API operations for managing
// graph resources, snapshots, import/export tasks, and query execution.
package neptunegraph

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/neptunegraph"
	neptunegraphconnect "vorpalstacks/internal/pb/aws/neptunegraph/neptunegraphconnect"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
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

func (h *AdminHandler) getStore(header http.Header) (*ngstore.NeptuneGraphStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	store, err := h.service.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	return store, nil
}

// ExecuteQuery handles the admin console ExecuteQuery request. Not implemented via admin console.
func (h *AdminHandler) ExecuteQuery(ctx context.Context, req *connect.Request[pb.ExecuteQueryInput]) (*connect.Response[pb.ExecuteQueryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// CancelQuery handles the admin console CancelQuery request. Not implemented via admin console.
func (h *AdminHandler) CancelQuery(ctx context.Context, req *connect.Request[pb.CancelQueryInput]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// GetQuery handles the admin console GetQuery request. Not implemented via admin console.
func (h *AdminHandler) GetQuery(ctx context.Context, req *connect.Request[pb.GetQueryInput]) (*connect.Response[pb.GetQueryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// ListQueries handles the admin console ListQueries request. Not implemented via admin console.
func (h *AdminHandler) ListQueries(ctx context.Context, req *connect.Request[pb.ListQueriesInput]) (*connect.Response[pb.ListQueriesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// GetGraphSummary handles the admin console GetGraphSummary request. Not implemented via admin console.
func (h *AdminHandler) GetGraphSummary(ctx context.Context, req *connect.Request[pb.GetGraphSummaryInput]) (*connect.Response[pb.GetGraphSummaryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// CreateGraph handles the admin console CreateGraph request. Not implemented via admin console.
func (h *AdminHandler) CreateGraph(ctx context.Context, req *connect.Request[pb.CreateGraphInput]) (*connect.Response[pb.CreateGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// DeleteGraph handles the admin console DeleteGraph request. Not implemented via admin console.
func (h *AdminHandler) DeleteGraph(ctx context.Context, req *connect.Request[pb.DeleteGraphInput]) (*connect.Response[pb.DeleteGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// GetGraph retrieves a graph by its identifier via the admin console.
func (h *AdminHandler) GetGraph(ctx context.Context, req *connect.Request[pb.GetGraphInput]) (*connect.Response[pb.GetGraphOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	graph, err := store.GetGraph(req.Msg.Graphidentifier)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return connect.NewResponse(graphToPb(graph)), nil
}

// ListGraphs returns all graph summaries via the admin console.
func (h *AdminHandler) ListGraphs(ctx context.Context, req *connect.Request[pb.ListGraphsInput]) (*connect.Response[pb.ListGraphsOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	graphs, _, _, err := store.ListGraphs(storecommon.ListOptions{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	pbGraphs := make([]*pb.GraphSummary, 0, len(graphs))
	for _, g := range graphs {
		pbGraphs = append(pbGraphs, graphSummaryToPb(g))
	}
	return connect.NewResponse(&pb.ListGraphsOutput{Graphs: pbGraphs}), nil
}

// UpdateGraph handles the admin console UpdateGraph request. Not implemented via admin console.
func (h *AdminHandler) UpdateGraph(ctx context.Context, req *connect.Request[pb.UpdateGraphInput]) (*connect.Response[pb.UpdateGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// StartGraph handles the admin console StartGraph request. Not implemented via admin console.
func (h *AdminHandler) StartGraph(ctx context.Context, req *connect.Request[pb.StartGraphInput]) (*connect.Response[pb.StartGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// StopGraph handles the admin console StopGraph request. Not implemented via admin console.
func (h *AdminHandler) StopGraph(ctx context.Context, req *connect.Request[pb.StopGraphInput]) (*connect.Response[pb.StopGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// ResetGraph handles the admin console ResetGraph request. Not implemented via admin console.
func (h *AdminHandler) ResetGraph(ctx context.Context, req *connect.Request[pb.ResetGraphInput]) (*connect.Response[pb.ResetGraphOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// RestoreGraphFromSnapshot handles the admin console RestoreGraphFromSnapshot request. Not implemented via admin console.
func (h *AdminHandler) RestoreGraphFromSnapshot(ctx context.Context, req *connect.Request[pb.RestoreGraphFromSnapshotInput]) (*connect.Response[pb.RestoreGraphFromSnapshotOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// CreateGraphSnapshot handles the admin console CreateGraphSnapshot request. Not implemented via admin console.
func (h *AdminHandler) CreateGraphSnapshot(ctx context.Context, req *connect.Request[pb.CreateGraphSnapshotInput]) (*connect.Response[pb.CreateGraphSnapshotOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// DeleteGraphSnapshot handles the admin console DeleteGraphSnapshot request. Not implemented via admin console.
func (h *AdminHandler) DeleteGraphSnapshot(ctx context.Context, req *connect.Request[pb.DeleteGraphSnapshotInput]) (*connect.Response[pb.DeleteGraphSnapshotOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// GetGraphSnapshot retrieves a graph snapshot by its identifier via the admin console.
func (h *AdminHandler) GetGraphSnapshot(ctx context.Context, req *connect.Request[pb.GetGraphSnapshotInput]) (*connect.Response[pb.GetGraphSnapshotOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	snapshot, err := store.GetSnapshot(req.Msg.Snapshotidentifier)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return connect.NewResponse(snapshotToPb(snapshot)), nil
}

// ListGraphSnapshots returns graph snapshots for the specified graph via the admin console.
func (h *AdminHandler) ListGraphSnapshots(ctx context.Context, req *connect.Request[pb.ListGraphSnapshotsInput]) (*connect.Response[pb.ListGraphSnapshotsOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	snapshots, _, _, err := store.ListSnapshots(storecommon.ListOptions{}, req.Msg.Graphidentifier)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	pbSnapshots := make([]*pb.GraphSnapshotSummary, 0, len(snapshots))
	for _, s := range snapshots {
		pbSnapshots = append(pbSnapshots, snapshotSummaryToPb(s))
	}
	return connect.NewResponse(&pb.ListGraphSnapshotsOutput{Graphsnapshots: pbSnapshots}), nil
}

// CreatePrivateGraphEndpoint handles the admin console CreatePrivateGraphEndpoint request. Not implemented via admin console.
func (h *AdminHandler) CreatePrivateGraphEndpoint(ctx context.Context, req *connect.Request[pb.CreatePrivateGraphEndpointInput]) (*connect.Response[pb.CreatePrivateGraphEndpointOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// DeletePrivateGraphEndpoint handles the admin console DeletePrivateGraphEndpoint request. Not implemented via admin console.
func (h *AdminHandler) DeletePrivateGraphEndpoint(ctx context.Context, req *connect.Request[pb.DeletePrivateGraphEndpointInput]) (*connect.Response[pb.DeletePrivateGraphEndpointOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// GetPrivateGraphEndpoint retrieves a private graph endpoint via the admin console.
func (h *AdminHandler) GetPrivateGraphEndpoint(ctx context.Context, req *connect.Request[pb.GetPrivateGraphEndpointInput]) (*connect.Response[pb.GetPrivateGraphEndpointOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	ep, err := store.GetEndpoint(req.Msg.Graphidentifier, req.Msg.Vpcid)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return connect.NewResponse(endpointToPb(ep)), nil
}

// ListPrivateGraphEndpoints returns all private endpoints for a graph via the admin console.
func (h *AdminHandler) ListPrivateGraphEndpoints(ctx context.Context, req *connect.Request[pb.ListPrivateGraphEndpointsInput]) (*connect.Response[pb.ListPrivateGraphEndpointsOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	endpoints, err := store.ListEndpoints(req.Msg.Graphidentifier)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	pbEndpoints := make([]*pb.PrivateGraphEndpointSummary, 0, len(endpoints))
	for _, ep := range endpoints {
		pbEndpoints = append(pbEndpoints, endpointSummaryToPb(ep))
	}
	return connect.NewResponse(&pb.ListPrivateGraphEndpointsOutput{Privategraphendpoints: pbEndpoints}), nil
}

// ListTagsForResource returns tags for a resource via the admin console.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceInput]) (*connect.Response[pb.ListTagsForResourceOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	tags, err := store.GetTags(req.Msg.Resourcearn)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&pb.ListTagsForResourceOutput{Tags: tags}), nil
}

// TagResource adds tags to a resource via the admin console.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceInput]) (*connect.Response[pb.TagResourceOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if len(req.Msg.Tags) > 0 {
		if err := store.AddTags(req.Msg.Resourcearn, req.Msg.Tags); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	return connect.NewResponse(&pb.TagResourceOutput{}), nil
}

// UntagResource removes tag keys from a resource via the admin console.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceInput]) (*connect.Response[pb.UntagResourceOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if len(req.Msg.Tagkeys) > 0 {
		if err := store.RemoveTags(req.Msg.Resourcearn, req.Msg.Tagkeys); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	return connect.NewResponse(&pb.UntagResourceOutput{}), nil
}

// CreateGraphUsingImportTask handles the admin console CreateGraphUsingImportTask request. Not implemented via admin console.
func (h *AdminHandler) CreateGraphUsingImportTask(ctx context.Context, req *connect.Request[pb.CreateGraphUsingImportTaskInput]) (*connect.Response[pb.CreateGraphUsingImportTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// GetImportTask retrieves an import task by its identifier via the admin console.
func (h *AdminHandler) GetImportTask(ctx context.Context, req *connect.Request[pb.GetImportTaskInput]) (*connect.Response[pb.GetImportTaskOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	task, err := store.GetImportTask(req.Msg.Taskidentifier)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return connect.NewResponse(importTaskToPb(task)), nil
}

// ListImportTasks returns all import task summaries via the admin console.
func (h *AdminHandler) ListImportTasks(ctx context.Context, req *connect.Request[pb.ListImportTasksInput]) (*connect.Response[pb.ListImportTasksOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	tasks, _, _, err := store.ListImportTasks(storecommon.ListOptions{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	pbTasks := make([]*pb.ImportTaskSummary, 0, len(tasks))
	for _, t := range tasks {
		pbTasks = append(pbTasks, importTaskSummaryToPb(t))
	}
	return connect.NewResponse(&pb.ListImportTasksOutput{Tasks: pbTasks}), nil
}

// CancelImportTask handles the admin console CancelImportTask request. Not implemented via admin console.
func (h *AdminHandler) CancelImportTask(ctx context.Context, req *connect.Request[pb.CancelImportTaskInput]) (*connect.Response[pb.CancelImportTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// StartImportTask handles the admin console StartImportTask request. Not implemented via admin console.
func (h *AdminHandler) StartImportTask(ctx context.Context, req *connect.Request[pb.StartImportTaskInput]) (*connect.Response[pb.StartImportTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// StartExportTask handles the admin console StartExportTask request. Not implemented via admin console.
func (h *AdminHandler) StartExportTask(ctx context.Context, req *connect.Request[pb.StartExportTaskInput]) (*connect.Response[pb.StartExportTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

// GetExportTask retrieves an export task by its identifier via the admin console.
func (h *AdminHandler) GetExportTask(ctx context.Context, req *connect.Request[pb.GetExportTaskInput]) (*connect.Response[pb.GetExportTaskOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	task, err := store.GetExportTask(req.Msg.Taskidentifier)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return connect.NewResponse(exportTaskToPb(task)), nil
}

// ListExportTasks returns export task summaries via the admin console, optionally filtered by graph.
func (h *AdminHandler) ListExportTasks(ctx context.Context, req *connect.Request[pb.ListExportTasksInput]) (*connect.Response[pb.ListExportTasksOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	tasks, _, _, err := store.ListExportTasks(storecommon.ListOptions{}, req.Msg.Graphidentifier)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	pbTasks := make([]*pb.ExportTaskSummary, 0, len(tasks))
	for _, t := range tasks {
		pbTasks = append(pbTasks, exportTaskSummaryToPb(t))
	}
	return connect.NewResponse(&pb.ListExportTasksOutput{Tasks: pbTasks}), nil
}

// CancelExportTask handles the admin console CancelExportTask request. Not implemented via admin console.
func (h *AdminHandler) CancelExportTask(ctx context.Context, req *connect.Request[pb.CancelExportTaskInput]) (*connect.Response[pb.CancelExportTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errNotImplemented())
}

func errNotImplemented() error {
	return connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented via admin console"))
}

func graphToPb(g *ngstore.Graph) *pb.GetGraphOutput {
	return &pb.GetGraphOutput{
		Id:                 g.Id,
		Name:               g.Name,
		Arn:                g.Arn,
		Status:             graphStatusToPb(g.Status),
		Statusreason:       g.StatusReason,
		Provisionedmemory:  derefInt32(g.ProvisionedMemory),
		Replicacount:       derefInt32(g.ReplicaCount),
		Deletionprotection: boolToStr(g.DeletionProtection),
		Publicconnectivity: boolToStr(g.PublicConnectivity),
		Endpoint:           g.Endpoint,
		Kmskeyidentifier:   g.KmsKeyIdentifier,
		Buildnumber:        g.BuildNumber,
		Createtime:         timePtrToStr(g.CreateTime),
		Sourcesnapshotid:   g.SourceSnapshotId,
	}
}

func graphSummaryToPb(g *ngstore.Graph) *pb.GraphSummary {
	return &pb.GraphSummary{
		Id:                 g.Id,
		Name:               g.Name,
		Arn:                g.Arn,
		Status:             graphStatusToPb(g.Status),
		Provisionedmemory:  derefInt32(g.ProvisionedMemory),
		Replicacount:       derefInt32(g.ReplicaCount),
		Deletionprotection: boolToStr(g.DeletionProtection),
		Publicconnectivity: boolToStr(g.PublicConnectivity),
		Endpoint:           g.Endpoint,
		Kmskeyidentifier:   g.KmsKeyIdentifier,
	}
}

func snapshotToPb(s *ngstore.GraphSnapshot) *pb.GetGraphSnapshotOutput {
	return &pb.GetGraphSnapshotOutput{
		Id:                 s.Id,
		Name:               s.Name,
		Arn:                s.Arn,
		Sourcegraphid:      s.SourceGraphId,
		Status:             snapshotStatusToPb(s.Status),
		Kmskeyidentifier:   s.KmsKeyIdentifier,
		Snapshotcreatetime: timePtrToStr(s.SnapshotCreateTime),
	}
}

func snapshotSummaryToPb(s *ngstore.GraphSnapshot) *pb.GraphSnapshotSummary {
	return &pb.GraphSnapshotSummary{
		Id:                 s.Id,
		Name:               s.Name,
		Arn:                s.Arn,
		Sourcegraphid:      s.SourceGraphId,
		Status:             snapshotStatusToPb(s.Status),
		Kmskeyidentifier:   s.KmsKeyIdentifier,
		Snapshotcreatetime: timePtrToStr(s.SnapshotCreateTime),
	}
}

func endpointToPb(ep *ngstore.PrivateGraphEndpoint) *pb.GetPrivateGraphEndpointOutput {
	return &pb.GetPrivateGraphEndpointOutput{
		Vpcid:         ep.VpcId,
		Vpcendpointid: ep.VpcEndpointId,
		Status:        endpointStatusToPb(ep.Status),
		Subnetids:     ep.SubnetIds,
	}
}

func endpointSummaryToPb(ep *ngstore.PrivateGraphEndpoint) *pb.PrivateGraphEndpointSummary {
	return &pb.PrivateGraphEndpointSummary{
		Vpcid:         ep.VpcId,
		Vpcendpointid: ep.VpcEndpointId,
		Status:        endpointStatusToPb(ep.Status),
		Subnetids:     ep.SubnetIds,
	}
}

func importTaskToPb(t *ngstore.ImportTask) *pb.GetImportTaskOutput {
	return &pb.GetImportTaskOutput{
		Taskid:            t.TaskId,
		Graphid:           t.GraphId,
		Source:            t.Source,
		Format:            formatToPb(t.Format),
		Rolearn:           t.RoleArn,
		Parquettype:       parquetTypeToPb(t.ParquetType),
		Status:            importTaskStatusToPb(t.Status),
		Statusreason:      t.StatusReason,
		Attemptnumber:     int32ToStr(t.AttemptNumber),
		Importoptions:     importOptionsToPb(t.ImportOptions),
		Importtaskdetails: importTaskDetailsToPb(t.ImportTaskDetails),
	}
}

func importTaskSummaryToPb(t *ngstore.ImportTask) *pb.ImportTaskSummary {
	return &pb.ImportTaskSummary{
		Taskid:      t.TaskId,
		Graphid:     t.GraphId,
		Source:      t.Source,
		Format:      formatToPb(t.Format),
		Rolearn:     t.RoleArn,
		Parquettype: parquetTypeToPb(t.ParquetType),
		Status:      importTaskStatusToPb(t.Status),
	}
}

func exportTaskToPb(t *ngstore.ExportTask) *pb.GetExportTaskOutput {
	return &pb.GetExportTaskOutput{
		Taskid:           t.TaskId,
		Graphid:          t.GraphId,
		Destination:      t.Destination,
		Format:           exportFormatToPb(t.Format),
		Rolearn:          t.RoleArn,
		Parquettype:      parquetTypeToPb(t.ParquetType),
		Kmskeyidentifier: t.KmsKeyIdentifier,
		Status:           exportTaskStatusToPb(t.Status),
		Statusreason:     t.StatusReason,
		Exportfilter:     exportFilterToPb(t.ExportFilter),
	}
}

func exportTaskSummaryToPb(t *ngstore.ExportTask) *pb.ExportTaskSummary {
	return &pb.ExportTaskSummary{
		Taskid:      t.TaskId,
		Graphid:     t.GraphId,
		Destination: t.Destination,
		Format:      exportFormatToPb(t.Format),
		Rolearn:     t.RoleArn,
		Parquettype: parquetTypeToPb(t.ParquetType),
		Status:      exportTaskStatusToPb(t.Status),
	}
}

func derefInt32(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func timePtrToStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}

func int32ToStr(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func int32PtrToStr(v *int32) string {
	if v == nil {
		return ""
	}
	return strconv.FormatInt(int64(*v), 10)
}

func int64ToStr(v *int64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatInt(*v, 10)
}

func strPtrToStr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func graphStatusToPb(s string) pb.GraphStatus {
	switch s {
	case "AVAILABLE":
		return pb.GraphStatus_GRAPH_STATUS_AVAILABLE
	case "CREATING":
		return pb.GraphStatus_GRAPH_STATUS_CREATING
	case "DELETING":
		return pb.GraphStatus_GRAPH_STATUS_DELETING
	case "UPDATING":
		return pb.GraphStatus_GRAPH_STATUS_UPDATING
	case "STARTING":
		return pb.GraphStatus_GRAPH_STATUS_STARTING
	case "STOPPING":
		return pb.GraphStatus_GRAPH_STATUS_STOPPING
	case "STOPPED":
		return pb.GraphStatus_GRAPH_STATUS_STOPPED
	case "RESETTING":
		return pb.GraphStatus_GRAPH_STATUS_RESETTING
	case "SNAPSHOTTING":
		return pb.GraphStatus_GRAPH_STATUS_SNAPSHOTTING
	case "IMPORTING":
		return pb.GraphStatus_GRAPH_STATUS_IMPORTING
	case "FAILED":
		return pb.GraphStatus_GRAPH_STATUS_FAILED
	default:
		return pb.GraphStatus_GRAPH_STATUS_AVAILABLE
	}
}

func snapshotStatusToPb(s string) pb.SnapshotStatus {
	switch s {
	case "AVAILABLE":
		return pb.SnapshotStatus_SNAPSHOT_STATUS_AVAILABLE
	case "CREATING":
		return pb.SnapshotStatus_SNAPSHOT_STATUS_CREATING
	case "DELETING":
		return pb.SnapshotStatus_SNAPSHOT_STATUS_DELETING
	case "FAILED":
		return pb.SnapshotStatus_SNAPSHOT_STATUS_FAILED
	default:
		return pb.SnapshotStatus_SNAPSHOT_STATUS_AVAILABLE
	}
}

func endpointStatusToPb(s string) pb.PrivateGraphEndpointStatus {
	switch s {
	case "AVAILABLE":
		return pb.PrivateGraphEndpointStatus_PRIVATE_GRAPH_ENDPOINT_STATUS_AVAILABLE
	case "CREATING":
		return pb.PrivateGraphEndpointStatus_PRIVATE_GRAPH_ENDPOINT_STATUS_CREATING
	case "DELETING":
		return pb.PrivateGraphEndpointStatus_PRIVATE_GRAPH_ENDPOINT_STATUS_DELETING
	case "FAILED":
		return pb.PrivateGraphEndpointStatus_PRIVATE_GRAPH_ENDPOINT_STATUS_FAILED
	default:
		return pb.PrivateGraphEndpointStatus_PRIVATE_GRAPH_ENDPOINT_STATUS_AVAILABLE
	}
}

func formatToPb(s string) pb.Format {
	switch s {
	case "NTRIPLES":
		return pb.Format_FORMAT_NTRIPLES
	case "OPEN_CYPHER":
		return pb.Format_FORMAT_OPEN_CYPHER
	case "PARQUET":
		return pb.Format_FORMAT_PARQUET
	case "CSV":
		return pb.Format_FORMAT_CSV
	default:
		return pb.Format_FORMAT_NTRIPLES
	}
}

func exportFormatToPb(s string) pb.ExportFormat {
	switch s {
	case "PARQUET":
		return pb.ExportFormat_EXPORT_FORMAT_PARQUET
	case "CSV":
		return pb.ExportFormat_EXPORT_FORMAT_CSV
	default:
		return pb.ExportFormat_EXPORT_FORMAT_PARQUET
	}
}

func parquetTypeToPb(s string) pb.ParquetType {
	if s == "COLUMNAR" {
		return pb.ParquetType_PARQUET_TYPE_COLUMNAR
	}
	return pb.ParquetType_PARQUET_TYPE_COLUMNAR
}

func importTaskStatusToPb(s string) pb.ImportTaskStatus {
	switch s {
	case "INITIALIZING":
		return pb.ImportTaskStatus_IMPORT_TASK_STATUS_INITIALIZING
	case "ANALYZING_DATA":
		return pb.ImportTaskStatus_IMPORT_TASK_STATUS_ANALYZING_DATA
	case "IMPORTING":
		return pb.ImportTaskStatus_IMPORT_TASK_STATUS_IMPORTING
	case "REPROVISIONING":
		return pb.ImportTaskStatus_IMPORT_TASK_STATUS_REPROVISIONING
	case "ROLLING_BACK":
		return pb.ImportTaskStatus_IMPORT_TASK_STATUS_ROLLING_BACK
	case "SUCCEEDED":
		return pb.ImportTaskStatus_IMPORT_TASK_STATUS_SUCCEEDED
	case "FAILED":
		return pb.ImportTaskStatus_IMPORT_TASK_STATUS_FAILED
	case "CANCELLED":
		return pb.ImportTaskStatus_IMPORT_TASK_STATUS_CANCELLED
	case "CANCELLING":
		return pb.ImportTaskStatus_IMPORT_TASK_STATUS_CANCELLING
	default:
		return pb.ImportTaskStatus_IMPORT_TASK_STATUS_INITIALIZING
	}
}

func exportTaskStatusToPb(s string) pb.ExportTaskStatus {
	switch s {
	case "INITIALIZING":
		return pb.ExportTaskStatus_EXPORT_TASK_STATUS_INITIALIZING
	case "EXPORTING":
		return pb.ExportTaskStatus_EXPORT_TASK_STATUS_EXPORTING
	case "SUCCEEDED":
		return pb.ExportTaskStatus_EXPORT_TASK_STATUS_SUCCEEDED
	case "FAILED":
		return pb.ExportTaskStatus_EXPORT_TASK_STATUS_FAILED
	case "CANCELLED":
		return pb.ExportTaskStatus_EXPORT_TASK_STATUS_CANCELLED
	case "CANCELLING":
		return pb.ExportTaskStatus_EXPORT_TASK_STATUS_CANCELLING
	default:
		return pb.ExportTaskStatus_EXPORT_TASK_STATUS_INITIALIZING
	}
}

func importOptionsToPb(opts *ngstore.ImportOptions) *pb.ImportOptions {
	if opts == nil || opts.Neptune == nil {
		return nil
	}
	n := opts.Neptune
	pbNeptune := &pb.NeptuneImportOptions{
		S3Exportpath:     n.S3ExportPath,
		S3Exportkmskeyid: n.S3ExportKmsKeyId,
	}
	if n.PreserveDefaultVertexLabels != nil {
		pbNeptune.Preservedefaultvertexlabels = boolToStr(*n.PreserveDefaultVertexLabels)
	}
	if n.PreserveEdgeIds != nil {
		pbNeptune.Preserveedgeids = boolToStr(*n.PreserveEdgeIds)
	}
	return &pb.ImportOptions{Neptune: pbNeptune}
}

func importTaskDetailsToPb(d *ngstore.ImportTaskDetails) *pb.ImportTaskDetails {
	if d == nil {
		return nil
	}
	return &pb.ImportTaskDetails{
		Progresspercentage:   int32PtrToStr(d.ProgressPercentage),
		Starttime:            timePtrToStr(d.StartTime),
		Timeelapsedseconds:   int64ToStr(d.TimeElapsedSeconds),
		Statementcount:       int64ToStr(d.StatementCount),
		Dictionaryentrycount: int64ToStr(d.DictionaryEntryCount),
		Errorcount:           int32PtrToStr(d.ErrorCount),
		Errordetails:         strPtrToStr(d.ErrorDetails),
		Status:               strPtrToStr(d.Status),
	}
}

func exportFilterToPb(f *ngstore.ExportFilter) *pb.ExportFilter {
	if f == nil {
		return nil
	}
	pbFilter := &pb.ExportFilter{
		Edgefilter:   make(map[string]*pb.ExportFilterElement),
		Vertexfilter: make(map[string]*pb.ExportFilterElement),
	}
	for k, v := range f.EdgeFilter {
		pbFilter.Edgefilter[k] = exportFilterElementToPb(v)
	}
	for k, v := range f.VertexFilter {
		pbFilter.Vertexfilter[k] = exportFilterElementToPb(v)
	}
	return pbFilter
}

func exportFilterElementToPb(e ngstore.ExportFilterElement) *pb.ExportFilterElement {
	pbElem := &pb.ExportFilterElement{
		Properties: make(map[string]*pb.ExportFilterPropertyAttributes),
	}
	for k, v := range e.Properties {
		var mvh pb.MultiValueHandlingType
		switch v.MultiValueHandling {
		case "MULTI_VALUE_HANDLING_TYPE_TO_LIST":
			mvh = pb.MultiValueHandlingType_MULTI_VALUE_HANDLING_TYPE_TO_LIST
		default:
			mvh = pb.MultiValueHandlingType_MULTI_VALUE_HANDLING_TYPE_PICK_FIRST
		}
		pbElem.Properties[k] = &pb.ExportFilterPropertyAttributes{
			Multivaluehandling: mvh,
			Outputtype:         strPtrToStr(v.OutputType),
			Sourcepropertyname: strPtrToStr(v.SourcePropertyName),
		}
	}
	return pbElem
}

// NewConnectHandler creates a gRPC-Web connect handler for the NeptuneGraph admin console.
func NewConnectHandler(svc *NeptuneGraphService, accountID string) (string, http.Handler) {
	return neptunegraphconnect.NewNeptuneGraphServiceHandler(NewAdminHandler(svc, accountID))
}
