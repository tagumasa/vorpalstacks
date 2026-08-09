package neptunegraph

// This file is the sole admin handler file that imports the store package.
// It contains store-accessing wrapper methods used by admin_handler.go,
// plus pure proto conversion helpers (toPb* functions) that translate
// store types to proto types for response marshalling.

import (
	"net/http"
	"strconv"
	"time"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/utils/ptrutil"
	"vorpalstacks/internal/utils/timeutils"

	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"

	pb "vorpalstacks/internal/pb/aws/neptunegraph"
)

// getStore returns the per-region NeptuneGraph store for the given header.
func (h *AdminHandler) getStore(header http.Header) (*ngstore.NeptuneGraphStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	return h.service.GetStoreForRegion(region)
}

// ---------------------------------------------------------------------------
// Store-accessing wrapper methods — each returns proto types so that
// admin_handler.go needs zero store imports.
// ---------------------------------------------------------------------------

func (h *AdminHandler) getGraphPb(header http.Header, id string) (*pb.GetGraphOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	graph, err := store.GetGraph(id)
	if err != nil {
		return nil, err
	}
	return graphToPb(graph), nil
}

func (h *AdminHandler) listGraphsPb(header http.Header) ([]*pb.GraphSummary, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	graphs, _, _, err := store.ListGraphs(storecommon.ListOptions{})
	if err != nil {
		return nil, err
	}
	summaries := make([]*pb.GraphSummary, 0, len(graphs))
	for _, g := range graphs {
		summaries = append(summaries, graphSummaryToPb(g))
	}
	return summaries, nil
}

func (h *AdminHandler) getGraphSnapshotPb(header http.Header, id string) (*pb.GetGraphSnapshotOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	snapshot, err := store.GetSnapshot(id)
	if err != nil {
		return nil, err
	}
	return snapshotToPb(snapshot), nil
}

func (h *AdminHandler) listGraphSnapshotsPb(header http.Header, graphID string) ([]*pb.GraphSnapshotSummary, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	snapshots, _, _, err := store.ListSnapshots(storecommon.ListOptions{}, graphID)
	if err != nil {
		return nil, err
	}
	summaries := make([]*pb.GraphSnapshotSummary, 0, len(snapshots))
	for _, s := range snapshots {
		summaries = append(summaries, snapshotSummaryToPb(s))
	}
	return summaries, nil
}

func (h *AdminHandler) getPrivateGraphEndpointPb(header http.Header, graphID, vpcID string) (*pb.GetPrivateGraphEndpointOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	ep, err := store.GetEndpoint(graphID, vpcID)
	if err != nil {
		return nil, err
	}
	return endpointToPb(ep), nil
}

func (h *AdminHandler) listPrivateGraphEndpointsPb(header http.Header, graphID string) ([]*pb.PrivateGraphEndpointSummary, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	endpoints, err := store.ListEndpoints(graphID)
	if err != nil {
		return nil, err
	}
	summaries := make([]*pb.PrivateGraphEndpointSummary, 0, len(endpoints))
	for _, ep := range endpoints {
		summaries = append(summaries, endpointSummaryToPb(ep))
	}
	return summaries, nil
}

func (h *AdminHandler) listTagsPb(header http.Header, arn string) (map[string]string, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	return store.GetTags(arn)
}

func (h *AdminHandler) tagResourcePb(header http.Header, arn string, tags map[string]string) error {
	if len(tags) == 0 {
		return nil
	}
	store, err := h.getStore(header)
	if err != nil {
		return err
	}
	return store.AddTags(arn, tags)
}

func (h *AdminHandler) untagResourcePb(header http.Header, arn string, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	store, err := h.getStore(header)
	if err != nil {
		return err
	}
	return store.RemoveTags(arn, keys)
}

func (h *AdminHandler) getImportTaskPb(header http.Header, id string) (*pb.GetImportTaskOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	task, err := store.GetImportTask(id)
	if err != nil {
		return nil, err
	}
	return importTaskToPb(task), nil
}

func (h *AdminHandler) listImportTasksPb(header http.Header) ([]*pb.ImportTaskSummary, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	tasks, _, _, err := store.ListImportTasks(storecommon.ListOptions{})
	if err != nil {
		return nil, err
	}
	summaries := make([]*pb.ImportTaskSummary, 0, len(tasks))
	for _, t := range tasks {
		summaries = append(summaries, importTaskSummaryToPb(t))
	}
	return summaries, nil
}

func (h *AdminHandler) getExportTaskPb(header http.Header, id string) (*pb.GetExportTaskOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	task, err := store.GetExportTask(id)
	if err != nil {
		return nil, err
	}
	return exportTaskToPb(task), nil
}

func (h *AdminHandler) listExportTasksPb(header http.Header, graphID string) ([]*pb.ExportTaskSummary, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	tasks, _, _, err := store.ListExportTasks(storecommon.ListOptions{}, graphID)
	if err != nil {
		return nil, err
	}
	summaries := make([]*pb.ExportTaskSummary, 0, len(tasks))
	for _, t := range tasks {
		summaries = append(summaries, exportTaskSummaryToPb(t))
	}
	return summaries, nil
}

// ---------------------------------------------------------------------------
// Proto conversion helpers (store types → proto types)
// ---------------------------------------------------------------------------

func graphToPb(g *ngstore.Graph) *pb.GetGraphOutput {
	return &pb.GetGraphOutput{
		Id:                 g.Id,
		Name:               g.Name,
		Arn:                g.Arn,
		Status:             graphStatusToPb(g.Status),
		Statusreason:       g.StatusReason,
		Provisionedmemory:  g.ProvisionedMemory,
		Replicacount:       g.ReplicaCount,
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
		Provisionedmemory:  g.ProvisionedMemory,
		Replicacount:       g.ReplicaCount,
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
	return t.UTC().Format(timeutils.ISO8601UTCFormat)
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
		Errordetails:         ptrutil.DerefOrZero(d.ErrorDetails),
		Status:               ptrutil.DerefOrZero(d.Status),
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
			Outputtype:         ptrutil.DerefOrZero(v.OutputType),
			Sourcepropertyname: ptrutil.DerefOrZero(v.SourcePropertyName),
		}
	}
	return pbElem
}
