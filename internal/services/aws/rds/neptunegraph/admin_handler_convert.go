package neptunegraph

// This file is the sole admin handler file that imports the store package.
// It contains store-accessing wrapper methods used by admin_handler.go,
// plus pure proto conversion helpers (toPb* functions) that translate
// store types to proto types for response marshalling.

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"vorpalstacks/internal/common/defaults"

	"vorpalstacks/internal/utils/ptrutil"
	"vorpalstacks/internal/utils/timeutils"

	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"

	pb "vorpalstacks/internal/pb/aws/neptunegraph"
)

// getStore returns the per-region NeptuneGraph store for the given header.
func (h *AdminHandler) getStore(header http.Header) (*ngstore.NeptuneGraphStore, error) {
	region := defaults.GetRegionFromHeader(header)
	return h.service.GetStoreForRegion(region)
}

// ---------------------------------------------------------------------------

func (h *AdminHandler) getGraphPb(header http.Header, id string) (*pb.GetGraphOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	graph, err := h.service.getGraphByIDCore(store, id)
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
	result, err := h.service.listGraphsCore(store, &ListGraphsInput{})
	if err != nil {
		return nil, err
	}
	summaries := make([]*pb.GraphSummary, 0, len(result.Graphs))
	for _, g := range result.Graphs {
		summaries = append(summaries, graphSummaryToPb(g))
	}
	return summaries, nil
}

func (h *AdminHandler) getGraphSnapshotPb(header http.Header, id string) (*pb.GetGraphSnapshotOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	snapshot, err := h.service.getGraphSnapshotCore(store, &GetGraphSnapshotInput{SnapshotIdentifier: id})
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
	result, err := h.service.listGraphSnapshotsCore(store, &ListGraphSnapshotsInput{GraphIdentifier: graphID})
	if err != nil {
		return nil, err
	}
	summaries := make([]*pb.GraphSnapshotSummary, 0, len(result.Snapshots))
	for _, s := range result.Snapshots {
		summaries = append(summaries, snapshotSummaryToPb(s))
	}
	return summaries, nil
}

func (h *AdminHandler) getPrivateGraphEndpointPb(header http.Header, graphID, vpcID string) (*pb.GetPrivateGraphEndpointOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	ep, err := h.service.getPrivateGraphEndpointCore(store, &GetPrivateGraphEndpointInput{GraphIdentifier: graphID, VpcId: vpcID})
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
	endpoints, err := h.service.listPrivateGraphEndpointsCore(store, &ListPrivateGraphEndpointsInput{GraphIdentifier: graphID})
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
	return h.service.listTagsForResourceCore(store, &ListTagsForResourceInput{ResourceArn: arn})
}

func (h *AdminHandler) tagResourcePb(header http.Header, arn string, tags map[string]string) error {
	if len(tags) == 0 {
		return nil
	}
	store, err := h.getStore(header)
	if err != nil {
		return err
	}
	_, err = h.service.tagResourceCore(store, &TagResourceInput{ResourceArn: arn, Tags: tags})
	return err
}

func (h *AdminHandler) untagResourcePb(header http.Header, arn string, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	store, err := h.getStore(header)
	if err != nil {
		return err
	}
	_, err = h.service.untagResourceCore(store, &UntagResourceInput{ResourceArn: arn, TagKeys: keys})
	return err
}

func (h *AdminHandler) getImportTaskPb(header http.Header, id string) (*pb.GetImportTaskOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	task, err := h.service.getImportTaskCore(store, id)
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
	result, err := h.service.listImportTasksCore(store, &ListImportTasksInput{})
	if err != nil {
		return nil, err
	}
	summaries := make([]*pb.ImportTaskSummary, 0, len(result.Tasks))
	for _, t := range result.Tasks {
		summaries = append(summaries, importTaskSummaryToPb(t))
	}
	return summaries, nil
}

func (h *AdminHandler) getExportTaskPb(header http.Header, id string) (*pb.GetExportTaskOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}
	task, err := h.service.getExportTaskCore(store, id)
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
	result, err := h.service.listExportTasksCore(store, &ListExportTasksInput{GraphIdentifier: graphID})
	if err != nil {
		return nil, err
	}
	summaries := make([]*pb.ExportTaskSummary, 0, len(result.Tasks))
	for _, t := range result.Tasks {
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

// queryToGetPb renders a query record for the GetQuery admin response.
func queryToGetPb(q *ngstore.QueryRecord) *pb.GetQueryOutput {
	return &pb.GetQueryOutput{
		Id:          q.Id,
		Querystring: q.QueryString,
		Elapsed:     strconv.Itoa(int(q.Elapsed)),
		Waited:      strconv.Itoa(int(q.Waited)),
	}
}

// queryToSummaryPb renders a query record for the ListQueries summaries.
func queryToSummaryPb(q *ngstore.QueryRecord) *pb.QuerySummary {
	return &pb.QuerySummary{
		Id:          q.Id,
		Querystring: q.QueryString,
		Elapsed:     strconv.Itoa(int(q.Elapsed)),
		Waited:      strconv.Itoa(int(q.Waited)),
		State:       queryStateToPb(q.State),
	}
}

// queryStateToPb maps a stored query state onto the console enum.
func queryStateToPb(s string) pb.QueryState {
	switch s {
	case "WAITING":
		return pb.QueryState_QUERY_STATE_WAITING
	case "CANCELLING":
		return pb.QueryState_QUERY_STATE_CANCELLING
	default:
		return pb.QueryState_QUERY_STATE_RUNNING
	}
}

// queryStateInputPbToString maps the console state filter onto the Core's
// wire-style state string.
func queryStateInputPbToString(s pb.QueryStateInput) string {
	switch s {
	case pb.QueryStateInput_QUERY_STATE_INPUT_RUNNING:
		return "RUNNING"
	case pb.QueryStateInput_QUERY_STATE_INPUT_WAITING:
		return "WAITING"
	case pb.QueryStateInput_QUERY_STATE_INPUT_CANCELLING:
		return "CANCELLING"
	default:
		return "ALL"
	}
}

// graphSummaryModePbToString maps the console summary mode onto the Core's
// wire-style mode string.
func graphSummaryModePbToString(m pb.GraphSummaryMode) string {
	if m == pb.GraphSummaryMode_GRAPH_SUMMARY_MODE_BASIC {
		return "BASIC"
	}
	return "DETAILED"
}

// longValuedMapListToPb converts an int-valued property map list entry.
func longValuedMapListToPb(m map[string]int64) *pb.LongValuedMapListEntry {
	value := make(map[string]string, len(m))
	for k, v := range m {
		value[k] = strconv.FormatInt(v, 10)
	}
	return &pb.LongValuedMapListEntry{Value: value}
}

// graphSummaryOutputPb renders the computed graph statistics.
func graphSummaryOutputPb(result *GetGraphSummaryResult) *pb.GetGraphSummaryOutput {
	s := result.Summary
	out := &pb.GetGraphSummaryOutput{
		Laststatisticscomputationtime: result.StatsTime.UTC().Format(timeutils.ISO8601UTCFormat),
		Version:                       "v1",
	}
	if s == nil {
		return out
	}
	out.Graphsummary = &pb.GraphDataSummary{
		Numnodes:                int64ToStr(s.NumNodes),
		Numedges:                int64ToStr(s.NumEdges),
		Numnodelabels:           int64ToStr(s.NumNodeLabels),
		Numedgelabels:           int64ToStr(s.NumEdgeLabels),
		Numnodeproperties:       int64ToStr(s.NumNodeProperties),
		Numedgeproperties:       int64ToStr(s.NumEdgeProperties),
		Totalnodepropertyvalues: int64ToStr(s.TotalNodePropertyValues),
		Totaledgepropertyvalues: int64ToStr(s.TotalEdgePropertyValues),
		Nodelabels:              s.NodeLabels,
		Edgelabels:              s.EdgeLabels,
	}
	for _, m := range s.NodeProperties {
		out.Graphsummary.Nodeproperties = append(out.Graphsummary.Nodeproperties, longValuedMapListToPb(m))
	}
	for _, m := range s.EdgeProperties {
		out.Graphsummary.Edgeproperties = append(out.Graphsummary.Edgeproperties, longValuedMapListToPb(m))
	}
	for _, ns := range s.NodeStructures {
		out.Graphsummary.Nodestructures = append(out.Graphsummary.Nodestructures, &pb.NodeStructure{
			Count:                      int64ToStr(ns.Count),
			Distinctoutgoingedgelabels: ns.DistinctOutgoingEdgeLabels,
			Nodeproperties:             ns.NodeProperties,
		})
	}
	for _, es := range s.EdgeStructures {
		out.Graphsummary.Edgestructures = append(out.Graphsummary.Edgestructures, &pb.EdgeStructure{
			Count:          int64ToStr(es.Count),
			Edgeproperties: es.EdgeProperties,
		})
	}
	return out
}

// strToBool parses the console's string-encoded boolean members.
func strToBool(s string) bool {
	return strings.EqualFold(s, "true")
}

// vectorSearchConfigToPb renders the stored vector-search configuration.
func vectorSearchConfigToPb(v *ngstore.VectorSearchConfig) *pb.VectorSearchConfiguration {
	if v == nil {
		return nil
	}
	return &pb.VectorSearchConfiguration{Dimension: v.Dimension}
}

func graphFieldsPb(g *ngstore.Graph) (id, name, arn, statusReason, deletionProtection, publicConnectivity, endpoint, kmsKeyIdentifier, buildNumber, createTime, sourceSnapshotId string, status pb.GraphStatus, provisionedMemory, replicaCount *int32, vectorSearch *pb.VectorSearchConfiguration) {
	return g.Id, g.Name, g.Arn, g.StatusReason, boolToStr(g.DeletionProtection), boolToStr(g.PublicConnectivity), g.Endpoint, g.KmsKeyIdentifier, g.BuildNumber, timePtrToStr(g.CreateTime), g.SourceSnapshotId, graphStatusToPb(g.Status), g.ProvisionedMemory, g.ReplicaCount, vectorSearchConfigToPb(g.VectorSearchConfiguration)
}

func graphToCreateGraphPb(g *ngstore.Graph) *pb.CreateGraphOutput {
	id, name, arn, sreason, delprot, pubconn, ep, kms, build, ct, srcsnap, status, pm, rc, vsc := graphFieldsPb(g)
	return &pb.CreateGraphOutput{Id: id, Name: name, Arn: arn, Status: status, Statusreason: sreason, Provisionedmemory: pm, Replicacount: rc, Deletionprotection: delprot, Publicconnectivity: pubconn, Endpoint: ep, Kmskeyidentifier: kms, Buildnumber: build, Createtime: ct, Sourcesnapshotid: srcsnap, Vectorsearchconfiguration: vsc}
}

func graphToUpdateGraphPb(g *ngstore.Graph) *pb.UpdateGraphOutput {
	id, name, arn, sreason, delprot, pubconn, ep, kms, build, ct, srcsnap, status, pm, rc, vsc := graphFieldsPb(g)
	return &pb.UpdateGraphOutput{Id: id, Name: name, Arn: arn, Status: status, Statusreason: sreason, Provisionedmemory: pm, Replicacount: rc, Deletionprotection: delprot, Publicconnectivity: pubconn, Endpoint: ep, Kmskeyidentifier: kms, Buildnumber: build, Createtime: ct, Sourcesnapshotid: srcsnap, Vectorsearchconfiguration: vsc}
}

func graphToDeleteGraphPb(g *ngstore.Graph) *pb.DeleteGraphOutput {
	id, name, arn, sreason, delprot, pubconn, ep, kms, build, ct, srcsnap, status, pm, rc, vsc := graphFieldsPb(g)
	return &pb.DeleteGraphOutput{Id: id, Name: name, Arn: arn, Status: status, Statusreason: sreason, Provisionedmemory: pm, Replicacount: rc, Deletionprotection: delprot, Publicconnectivity: pubconn, Endpoint: ep, Kmskeyidentifier: kms, Buildnumber: build, Createtime: ct, Sourcesnapshotid: srcsnap, Vectorsearchconfiguration: vsc}
}

func graphToResetGraphPb(g *ngstore.Graph) *pb.ResetGraphOutput {
	id, name, arn, sreason, delprot, pubconn, ep, kms, build, ct, srcsnap, status, pm, rc, vsc := graphFieldsPb(g)
	return &pb.ResetGraphOutput{Id: id, Name: name, Arn: arn, Status: status, Statusreason: sreason, Provisionedmemory: pm, Replicacount: rc, Deletionprotection: delprot, Publicconnectivity: pubconn, Endpoint: ep, Kmskeyidentifier: kms, Buildnumber: build, Createtime: ct, Sourcesnapshotid: srcsnap, Vectorsearchconfiguration: vsc}
}

func graphToStartGraphPb(g *ngstore.Graph) *pb.StartGraphOutput {
	id, name, arn, sreason, delprot, pubconn, ep, kms, build, ct, srcsnap, status, pm, rc, vsc := graphFieldsPb(g)
	return &pb.StartGraphOutput{Id: id, Name: name, Arn: arn, Status: status, Statusreason: sreason, Provisionedmemory: pm, Replicacount: rc, Deletionprotection: delprot, Publicconnectivity: pubconn, Endpoint: ep, Kmskeyidentifier: kms, Buildnumber: build, Createtime: ct, Sourcesnapshotid: srcsnap, Vectorsearchconfiguration: vsc}
}

func graphToStopGraphPb(g *ngstore.Graph) *pb.StopGraphOutput {
	id, name, arn, sreason, delprot, pubconn, ep, kms, build, ct, srcsnap, status, pm, rc, vsc := graphFieldsPb(g)
	return &pb.StopGraphOutput{Id: id, Name: name, Arn: arn, Status: status, Statusreason: sreason, Provisionedmemory: pm, Replicacount: rc, Deletionprotection: delprot, Publicconnectivity: pubconn, Endpoint: ep, Kmskeyidentifier: kms, Buildnumber: build, Createtime: ct, Sourcesnapshotid: srcsnap, Vectorsearchconfiguration: vsc}
}

func graphToRestoreGraphFromSnapshotPb(g *ngstore.Graph) *pb.RestoreGraphFromSnapshotOutput {
	id, name, arn, sreason, delprot, pubconn, ep, kms, build, ct, srcsnap, status, pm, rc, vsc := graphFieldsPb(g)
	return &pb.RestoreGraphFromSnapshotOutput{Id: id, Name: name, Arn: arn, Status: status, Statusreason: sreason, Provisionedmemory: pm, Replicacount: rc, Deletionprotection: delprot, Publicconnectivity: pubconn, Endpoint: ep, Kmskeyidentifier: kms, Buildnumber: build, Createtime: ct, Sourcesnapshotid: srcsnap, Vectorsearchconfiguration: vsc}
}

// formatPbToString maps the console format enum onto the Core's wire string.
func formatPbToString(f pb.Format) string {
	switch f {
	case pb.Format_FORMAT_OPEN_CYPHER:
		return "OPEN_CYPHER"
	case pb.Format_FORMAT_PARQUET:
		return "PARQUET"
	case pb.Format_FORMAT_CSV:
		return "CSV"
	default:
		return "NTRIPLES"
	}
}

// parquetTypePbToString maps the console parquet-type enum.
func parquetTypePbToString(p pb.ParquetType) string {
	if p == pb.ParquetType_PARQUET_TYPE_COLUMNAR {
		return "COLUMNAR"
	}
	return ""
}

// blankNodeHandlingPbToString maps the console blank-node-handling enum.
func blankNodeHandlingPbToString(b pb.BlankNodeHandling) string {
	if b == pb.BlankNodeHandling_BLANK_NODE_HANDLING_CONVERT_TO_IRI {
		return "CONVERT_TO_IRI"
	}
	return "PERMUTE"
}

// exportFormatPbToString maps the console export-format enum.
func exportFormatPbToString(f pb.ExportFormat) string {
	if f == pb.ExportFormat_EXPORT_FORMAT_CSV {
		return "CSV"
	}
	return "PARQUET"
}

// explainModeInputPbToString maps the console explain-mode request enum onto
// the Core's wire string. STATIC is the proto3 zero value, so an unset
// member is indistinguishable from an explicit STATIC; it maps to the empty
// string (the Core default) so ordinary console queries do not come back as
// explain plans. DETAILED is expressible explicitly.
func explainModeInputPbToString(m pb.ExplainMode) string {
	if m == pb.ExplainMode_EXPLAIN_MODE_DETAILS {
		return "DETAILS"
	}
	return ""
}

// planCacheInputPbToString maps the console plan-cache request enum. AUTO is
// the proto3 zero value; the empty string is behaviourally identical to AUTO
// in the execution Core, so it maps through unchanged.
func planCacheInputPbToString(p pb.PlanCacheType) string {
	switch p {
	case pb.PlanCacheType_PLAN_CACHE_TYPE_DISABLED:
		return "DISABLED"
	case pb.PlanCacheType_PLAN_CACHE_TYPE_ENABLED:
		return "ENABLED"
	default:
		return ""
	}
}

// importOptionsFromPb converts the console import options onto the store type.
func importOptionsFromPb(in *pb.ImportOptions) *ngstore.ImportOptions {
	if in == nil || in.Neptune == nil {
		return nil
	}
	n := in.Neptune
	opts := &ngstore.ImportOptions{Neptune: &ngstore.NeptuneImportOptions{
		S3ExportPath:     n.S3Exportpath,
		S3ExportKmsKeyId: n.S3Exportkmskeyid,
	}}
	if v := strToBool(n.Preservedefaultvertexlabels); v {
		opts.Neptune.PreserveDefaultVertexLabels = &v
	}
	if v := strToBool(n.Preserveedgeids); v {
		opts.Neptune.PreserveEdgeIds = &v
	}
	return opts
}

// exportFilterFromPb converts the console export filter onto the store type.
func exportFilterFromPb(in *pb.ExportFilter) *ngstore.ExportFilter {
	if in == nil {
		return nil
	}
	out := &ngstore.ExportFilter{
		EdgeFilter:   make(map[string]ngstore.ExportFilterElement),
		VertexFilter: make(map[string]ngstore.ExportFilterElement),
	}
	for k, elem := range in.Edgefilter {
		out.EdgeFilter[k] = exportFilterElementFromPb(elem)
	}
	for k, elem := range in.Vertexfilter {
		out.VertexFilter[k] = exportFilterElementFromPb(elem)
	}
	return out
}

// exportFilterElementFromPb converts one filter element.
func exportFilterElementFromPb(elem *pb.ExportFilterElement) ngstore.ExportFilterElement {
	out := ngstore.ExportFilterElement{Properties: make(map[string]ngstore.ExportFilterPropertyAttributes)}
	for k, attrs := range elem.GetProperties() {
		mvh := "PICK_FIRST"
		if attrs.GetMultivaluehandling() == pb.MultiValueHandlingType_MULTI_VALUE_HANDLING_TYPE_TO_LIST {
			mvh = "TO_LIST"
		}
		outType := attrs.GetOutputtype()
		srcProp := attrs.GetSourcepropertyname()
		out.Properties[k] = ngstore.ExportFilterPropertyAttributes{
			MultiValueHandling: mvh,
			OutputType:         &outType,
			SourcePropertyName: &srcProp,
		}
	}
	return out
}

func snapshotFieldsPb(s *ngstore.GraphSnapshot) (id, name, arn, sourceGraphID, kms, createTime string, status pb.SnapshotStatus) {
	return s.Id, s.Name, s.Arn, s.SourceGraphId, s.KmsKeyIdentifier, timePtrToStr(s.SnapshotCreateTime), snapshotStatusToPb(s.Status)
}

func snapshotToCreateGraphSnapshotPb(s *ngstore.GraphSnapshot) *pb.CreateGraphSnapshotOutput {
	id, name, arn, src, kms, ct, status := snapshotFieldsPb(s)
	return &pb.CreateGraphSnapshotOutput{Id: id, Name: name, Arn: arn, Sourcegraphid: src, Status: status, Kmskeyidentifier: kms, Snapshotcreatetime: ct}
}

func snapshotToDeleteGraphSnapshotPb(s *ngstore.GraphSnapshot) *pb.DeleteGraphSnapshotOutput {
	id, name, arn, src, kms, ct, status := snapshotFieldsPb(s)
	return &pb.DeleteGraphSnapshotOutput{Id: id, Name: name, Arn: arn, Sourcegraphid: src, Status: status, Kmskeyidentifier: kms, Snapshotcreatetime: ct}
}

func endpointFieldsPb(ep *ngstore.PrivateGraphEndpoint) (vpcID, vpcEndpointID string, status pb.PrivateGraphEndpointStatus, subnetIDs []string) {
	return ep.VpcId, ep.VpcEndpointId, endpointStatusToPb(ep.Status), ep.SubnetIds
}

func endpointToCreatePrivateGraphEndpointPb(ep *ngstore.PrivateGraphEndpoint) *pb.CreatePrivateGraphEndpointOutput {
	vpcID, vpcEp, status, subnets := endpointFieldsPb(ep)
	return &pb.CreatePrivateGraphEndpointOutput{Vpcid: vpcID, Vpcendpointid: vpcEp, Status: status, Subnetids: subnets}
}

func endpointToDeletePrivateGraphEndpointPb(ep *ngstore.PrivateGraphEndpoint) *pb.DeletePrivateGraphEndpointOutput {
	vpcID, vpcEp, status, subnets := endpointFieldsPb(ep)
	return &pb.DeletePrivateGraphEndpointOutput{Vpcid: vpcID, Vpcendpointid: vpcEp, Status: status, Subnetids: subnets}
}

func importTaskFieldsPb(t *ngstore.ImportTask) (taskID, graphID, source string, format pb.Format, roleArn string, parquetType pb.ParquetType, status pb.ImportTaskStatus, statusReason string) {
	return t.TaskId, t.GraphId, t.Source, formatToPb(t.Format), t.RoleArn, parquetTypeToPb(t.ParquetType), importTaskStatusToPb(t.Status), t.StatusReason
}

func importTaskToCreateGraphUsingImportTaskPb(t *ngstore.ImportTask) *pb.CreateGraphUsingImportTaskOutput {
	taskID, graphID, source, format, roleArn, pt, status, _ := importTaskFieldsPb(t)
	return &pb.CreateGraphUsingImportTaskOutput{Taskid: taskID, Graphid: graphID, Source: source, Format: format, Rolearn: roleArn, Parquettype: pt, Status: status}
}

func importTaskToStartImportTaskPb(t *ngstore.ImportTask) *pb.StartImportTaskOutput {
	taskID, graphID, source, format, roleArn, pt, status, _ := importTaskFieldsPb(t)
	return &pb.StartImportTaskOutput{Taskid: taskID, Graphid: graphID, Source: source, Format: format, Rolearn: roleArn, Parquettype: pt, Status: status}
}

func importTaskToCancelImportTaskPb(t *ngstore.ImportTask) *pb.CancelImportTaskOutput {
	taskID, graphID, source, format, roleArn, pt, status, _ := importTaskFieldsPb(t)
	return &pb.CancelImportTaskOutput{Taskid: taskID, Graphid: graphID, Source: source, Format: format, Rolearn: roleArn, Parquettype: pt, Status: status}
}

func exportTaskFieldsPb(t *ngstore.ExportTask) (taskID, graphID, destination string, format pb.ExportFormat, roleArn string, parquetType pb.ParquetType, kms, statusReason string, status pb.ExportTaskStatus, filter *pb.ExportFilter) {
	return t.TaskId, t.GraphId, t.Destination, exportFormatToPb(t.Format), t.RoleArn, parquetTypeToPb(t.ParquetType), t.KmsKeyIdentifier, t.StatusReason, exportTaskStatusToPb(t.Status), exportFilterToPb(t.ExportFilter)
}

func exportTaskToStartExportTaskPb(t *ngstore.ExportTask) *pb.StartExportTaskOutput {
	taskID, graphID, dest, format, roleArn, pt, kms, sreason, status, filter := exportTaskFieldsPb(t)
	return &pb.StartExportTaskOutput{Taskid: taskID, Graphid: graphID, Destination: dest, Format: format, Rolearn: roleArn, Parquettype: pt, Kmskeyidentifier: kms, Status: status, Statusreason: sreason, Exportfilter: filter}
}

func exportTaskToCancelExportTaskPb(t *ngstore.ExportTask) *pb.CancelExportTaskOutput {
	taskID, graphID, dest, format, roleArn, pt, kms, sreason, status, _ := exportTaskFieldsPb(t)
	return &pb.CancelExportTaskOutput{Taskid: taskID, Graphid: graphID, Destination: dest, Format: format, Rolearn: roleArn, Parquettype: pt, Kmskeyidentifier: kms, Status: status, Statusreason: sreason}
}
