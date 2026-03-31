package neptunedata

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/pb/aws/common"
	neptunedatapb "vorpalstacks/internal/pb/aws/neptunedata"
	neptunedataconnect "vorpalstacks/internal/pb/aws/neptunedata/neptunedataconnect"
)

type AdminHandler struct {
	neptunedataconnect.UnimplementedNeptunedataServiceHandler
}

var _ neptunedataconnect.NeptunedataServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func nop(name string) error {
	return connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for %s", name))
}

func (h *AdminHandler) GetEngineStatus(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[neptunedatapb.GetEngineStatusOutput], error) {
	return nil, nop("GetEngineStatus")
}

func (h *AdminHandler) ExecuteGremlinQuery(ctx context.Context, req *connect.Request[neptunedatapb.ExecuteGremlinQueryInput]) (*connect.Response[neptunedatapb.ExecuteGremlinQueryOutput], error) {
	return nil, nop("ExecuteGremlinQuery")
}

func (h *AdminHandler) ExecuteOpenCypherQuery(ctx context.Context, req *connect.Request[neptunedatapb.ExecuteOpenCypherQueryInput]) (*connect.Response[neptunedatapb.ExecuteOpenCypherQueryOutput], error) {
	return nil, nop("ExecuteOpenCypherQuery")
}

func (h *AdminHandler) CancelGremlinQuery(ctx context.Context, req *connect.Request[neptunedatapb.CancelGremlinQueryInput]) (*connect.Response[neptunedatapb.CancelGremlinQueryOutput], error) {
	return nil, nop("CancelGremlinQuery")
}

func (h *AdminHandler) GetGremlinQueryStatus(ctx context.Context, req *connect.Request[neptunedatapb.GetGremlinQueryStatusInput]) (*connect.Response[neptunedatapb.GetGremlinQueryStatusOutput], error) {
	return nil, nop("GetGremlinQueryStatus")
}

func (h *AdminHandler) ListGremlinQueries(ctx context.Context, req *connect.Request[neptunedatapb.ListGremlinQueriesInput]) (*connect.Response[neptunedatapb.ListGremlinQueriesOutput], error) {
	return nil, nop("ListGremlinQueries")
}

func (h *AdminHandler) ExecuteGremlinExplainQuery(ctx context.Context, req *connect.Request[neptunedatapb.ExecuteGremlinExplainQueryInput]) (*connect.Response[neptunedatapb.ExecuteGremlinExplainQueryOutput], error) {
	return nil, nop("ExecuteGremlinExplainQuery")
}

func (h *AdminHandler) ExecuteGremlinProfileQuery(ctx context.Context, req *connect.Request[neptunedatapb.ExecuteGremlinProfileQueryInput]) (*connect.Response[neptunedatapb.ExecuteGremlinProfileQueryOutput], error) {
	return nil, nop("ExecuteGremlinProfileQuery")
}

func (h *AdminHandler) ExecuteOpenCypherExplainQuery(ctx context.Context, req *connect.Request[neptunedatapb.ExecuteOpenCypherExplainQueryInput]) (*connect.Response[neptunedatapb.ExecuteOpenCypherExplainQueryOutput], error) {
	return nil, nop("ExecuteOpenCypherExplainQuery")
}

func (h *AdminHandler) GetOpenCypherQueryStatus(ctx context.Context, req *connect.Request[neptunedatapb.GetOpenCypherQueryStatusInput]) (*connect.Response[neptunedatapb.GetOpenCypherQueryStatusOutput], error) {
	return nil, nop("GetOpenCypherQueryStatus")
}

func (h *AdminHandler) ListOpenCypherQueries(ctx context.Context, req *connect.Request[neptunedatapb.ListOpenCypherQueriesInput]) (*connect.Response[neptunedatapb.ListOpenCypherQueriesOutput], error) {
	return nil, nop("ListOpenCypherQueries")
}

func (h *AdminHandler) CancelOpenCypherQuery(ctx context.Context, req *connect.Request[neptunedatapb.CancelOpenCypherQueryInput]) (*connect.Response[neptunedatapb.CancelOpenCypherQueryOutput], error) {
	return nil, nop("CancelOpenCypherQuery")
}

func (h *AdminHandler) ExecuteFastReset(ctx context.Context, req *connect.Request[neptunedatapb.ExecuteFastResetInput]) (*connect.Response[neptunedatapb.ExecuteFastResetOutput], error) {
	return nil, nop("ExecuteFastReset")
}

func (h *AdminHandler) GetPropertygraphStatistics(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[neptunedatapb.GetPropertygraphStatisticsOutput], error) {
	return nil, nop("GetPropertygraphStatistics")
}

func (h *AdminHandler) ManagePropertygraphStatistics(ctx context.Context, req *connect.Request[neptunedatapb.ManagePropertygraphStatisticsInput]) (*connect.Response[neptunedatapb.ManagePropertygraphStatisticsOutput], error) {
	return nil, nop("ManagePropertygraphStatistics")
}

func (h *AdminHandler) DeletePropertygraphStatistics(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[neptunedatapb.DeletePropertygraphStatisticsOutput], error) {
	return nil, nop("DeletePropertygraphStatistics")
}

func (h *AdminHandler) GetPropertygraphSummary(ctx context.Context, req *connect.Request[neptunedatapb.GetPropertygraphSummaryInput]) (*connect.Response[neptunedatapb.GetPropertygraphSummaryOutput], error) {
	return nil, nop("GetPropertygraphSummary")
}

func (h *AdminHandler) GetPropertygraphStream(ctx context.Context, req *connect.Request[neptunedatapb.GetPropertygraphStreamInput]) (*connect.Response[neptunedatapb.GetPropertygraphStreamOutput], error) {
	return nil, nop("GetPropertygraphStream")
}

func (h *AdminHandler) StartLoaderJob(ctx context.Context, req *connect.Request[neptunedatapb.StartLoaderJobInput]) (*connect.Response[neptunedatapb.StartLoaderJobOutput], error) {
	return nil, nop("StartLoaderJob")
}

func (h *AdminHandler) GetLoaderJobStatus(ctx context.Context, req *connect.Request[neptunedatapb.GetLoaderJobStatusInput]) (*connect.Response[neptunedatapb.GetLoaderJobStatusOutput], error) {
	return nil, nop("GetLoaderJobStatus")
}

func (h *AdminHandler) ListLoaderJobs(ctx context.Context, req *connect.Request[neptunedatapb.ListLoaderJobsInput]) (*connect.Response[neptunedatapb.ListLoaderJobsOutput], error) {
	return nil, nop("ListLoaderJobs")
}

func (h *AdminHandler) CancelLoaderJob(ctx context.Context, req *connect.Request[neptunedatapb.CancelLoaderJobInput]) (*connect.Response[neptunedatapb.CancelLoaderJobOutput], error) {
	return nil, nop("CancelLoaderJob")
}

func (h *AdminHandler) GetSparqlStatistics(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[neptunedatapb.GetSparqlStatisticsOutput], error) {
	return nil, nop("GetSparqlStatistics")
}

func (h *AdminHandler) ManageSparqlStatistics(ctx context.Context, req *connect.Request[neptunedatapb.ManageSparqlStatisticsInput]) (*connect.Response[neptunedatapb.ManageSparqlStatisticsOutput], error) {
	return nil, nop("ManageSparqlStatistics")
}

func (h *AdminHandler) DeleteSparqlStatistics(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[neptunedatapb.DeleteSparqlStatisticsOutput], error) {
	return nil, nop("DeleteSparqlStatistics")
}

func (h *AdminHandler) GetSparqlStream(ctx context.Context, req *connect.Request[neptunedatapb.GetSparqlStreamInput]) (*connect.Response[neptunedatapb.GetSparqlStreamOutput], error) {
	return nil, nop("GetSparqlStream")
}

func (h *AdminHandler) GetRDFGraphSummary(ctx context.Context, req *connect.Request[neptunedatapb.GetRDFGraphSummaryInput]) (*connect.Response[neptunedatapb.GetRDFGraphSummaryOutput], error) {
	return nil, nop("GetRDFGraphSummary")
}

func (h *AdminHandler) StartMLDataProcessingJob(ctx context.Context, req *connect.Request[neptunedatapb.StartMLDataProcessingJobInput]) (*connect.Response[neptunedatapb.StartMLDataProcessingJobOutput], error) {
	return nil, nop("StartMLDataProcessingJob")
}

func (h *AdminHandler) GetMLDataProcessingJob(ctx context.Context, req *connect.Request[neptunedatapb.GetMLDataProcessingJobInput]) (*connect.Response[neptunedatapb.GetMLDataProcessingJobOutput], error) {
	return nil, nop("GetMLDataProcessingJob")
}

func (h *AdminHandler) ListMLDataProcessingJobs(ctx context.Context, req *connect.Request[neptunedatapb.ListMLDataProcessingJobsInput]) (*connect.Response[neptunedatapb.ListMLDataProcessingJobsOutput], error) {
	return nil, nop("ListMLDataProcessingJobs")
}

func (h *AdminHandler) CancelMLDataProcessingJob(ctx context.Context, req *connect.Request[neptunedatapb.CancelMLDataProcessingJobInput]) (*connect.Response[neptunedatapb.CancelMLDataProcessingJobOutput], error) {
	return nil, nop("CancelMLDataProcessingJob")
}

func (h *AdminHandler) StartMLModelTrainingJob(ctx context.Context, req *connect.Request[neptunedatapb.StartMLModelTrainingJobInput]) (*connect.Response[neptunedatapb.StartMLModelTrainingJobOutput], error) {
	return nil, nop("StartMLModelTrainingJob")
}

func (h *AdminHandler) GetMLModelTrainingJob(ctx context.Context, req *connect.Request[neptunedatapb.GetMLModelTrainingJobInput]) (*connect.Response[neptunedatapb.GetMLModelTrainingJobOutput], error) {
	return nil, nop("GetMLModelTrainingJob")
}

func (h *AdminHandler) ListMLModelTrainingJobs(ctx context.Context, req *connect.Request[neptunedatapb.ListMLModelTrainingJobsInput]) (*connect.Response[neptunedatapb.ListMLModelTrainingJobsOutput], error) {
	return nil, nop("ListMLModelTrainingJobs")
}

func (h *AdminHandler) CancelMLModelTrainingJob(ctx context.Context, req *connect.Request[neptunedatapb.CancelMLModelTrainingJobInput]) (*connect.Response[neptunedatapb.CancelMLModelTrainingJobOutput], error) {
	return nil, nop("CancelMLModelTrainingJob")
}

func (h *AdminHandler) StartMLModelTransformJob(ctx context.Context, req *connect.Request[neptunedatapb.StartMLModelTransformJobInput]) (*connect.Response[neptunedatapb.StartMLModelTransformJobOutput], error) {
	return nil, nop("StartMLModelTransformJob")
}

func (h *AdminHandler) GetMLModelTransformJob(ctx context.Context, req *connect.Request[neptunedatapb.GetMLModelTransformJobInput]) (*connect.Response[neptunedatapb.GetMLModelTransformJobOutput], error) {
	return nil, nop("GetMLModelTransformJob")
}

func (h *AdminHandler) ListMLModelTransformJobs(ctx context.Context, req *connect.Request[neptunedatapb.ListMLModelTransformJobsInput]) (*connect.Response[neptunedatapb.ListMLModelTransformJobsOutput], error) {
	return nil, nop("ListMLModelTransformJobs")
}

func (h *AdminHandler) CancelMLModelTransformJob(ctx context.Context, req *connect.Request[neptunedatapb.CancelMLModelTransformJobInput]) (*connect.Response[neptunedatapb.CancelMLModelTransformJobOutput], error) {
	return nil, nop("CancelMLModelTransformJob")
}

func (h *AdminHandler) CreateMLEndpoint(ctx context.Context, req *connect.Request[neptunedatapb.CreateMLEndpointInput]) (*connect.Response[neptunedatapb.CreateMLEndpointOutput], error) {
	return nil, nop("CreateMLEndpoint")
}

func (h *AdminHandler) GetMLEndpoint(ctx context.Context, req *connect.Request[neptunedatapb.GetMLEndpointInput]) (*connect.Response[neptunedatapb.GetMLEndpointOutput], error) {
	return nil, nop("GetMLEndpoint")
}

func (h *AdminHandler) ListMLEndpoints(ctx context.Context, req *connect.Request[neptunedatapb.ListMLEndpointsInput]) (*connect.Response[neptunedatapb.ListMLEndpointsOutput], error) {
	return nil, nop("ListMLEndpoints")
}

func (h *AdminHandler) DeleteMLEndpoint(ctx context.Context, req *connect.Request[neptunedatapb.DeleteMLEndpointInput]) (*connect.Response[neptunedatapb.DeleteMLEndpointOutput], error) {
	return nil, nop("DeleteMLEndpoint")
}
