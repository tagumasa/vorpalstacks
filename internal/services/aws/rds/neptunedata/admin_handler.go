// Package neptunedata implements AWS Neptune Data API operations including
// property graph queries (OpenCypher, Gremlin), bulk loading, statistics,
// and streaming.
//
// NeptuneData is a sub-service of Neptune and shares the Neptune store
// (internal/store/aws/rds/neptune) for cluster metadata access.
package neptunedata

import (
	"context"
	"fmt"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/neptunedata"
	neptunedataconnect "vorpalstacks/internal/pb/aws/neptunedata/neptunedataconnect"
)

// AdminHandler implements the Neptune Data API gRPC-Web admin console handler.
// It shares the same NeptuneDataService instance as the HTTP dispatcher,
// providing read access to query states, loader jobs, and graph statistics.
type AdminHandler struct {
	neptunedataconnect.UnimplementedNeptunedataServiceHandler
	service *NeptuneDataService
}

var _ neptunedataconnect.NeptunedataServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Neptune Data API admin console handler backed
// by the given NeptuneDataService instance.
func NewAdminHandler(svc *NeptuneDataService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// errUnsupportedAdmin returns a gRPC error wrapping the AWS-compatible
// UnsupportedOperationException (HTTP 400), matching the HTTP API behaviour
// for unsupported operations such as SPARQL, ML, and streaming.
func errUnsupportedAdmin(category string) error {
	return svcerrors.AWSErrorToGRPC(unsupported(fmt.Sprintf("%s operations are not supported by vorpalstacks", category)))
}

// GetEngineStatus returns the health status and engine version information
// for the Neptune-compatible graph engine.
func (h *AdminHandler) GetEngineStatus(ctx context.Context, req *connect.Request[pbcommon.Empty]) (*connect.Response[pb.GetEngineStatusOutput], error) {
	s := h.service
	s.mu.RLock()
	startTime := s.startTime.UTC().Format(timeutils.ISO8601UTCFormat)
	s.mu.RUnlock()

	return connect.NewResponse(&pb.GetEngineStatusOutput{
		Status:    "healthy",
		Starttime: startTime,
		Gremlin: &pb.QueryLanguageVersion{
			Version: "3.7.x",
		},
		Opencypher: &pb.QueryLanguageVersion{
			Version: "2023-08-01",
		},
		Sparql: &pb.QueryLanguageVersion{
			Version: "1.1",
		},
		Labmode: map[string]string{},
		Settings: map[string]string{
			"neptune lab mode": "DISABLED",
		},
		Role: "writer",
	}), nil
}

// GetGremlinQueryStatus returns the status and timing of a Gremlin query.
func (h *AdminHandler) GetGremlinQueryStatus(ctx context.Context, req *connect.Request[pb.GetGremlinQueryStatusInput]) (*connect.Response[pb.GetGremlinQueryStatusOutput], error) {
	out, err := h.queryStatusPb(req.Header(), req.Msg.Queryid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.GetGremlinQueryStatusOutput{
		Queryid:        out.Queryid,
		Querystring:    out.Querystring,
		Queryevalstats: out.Queryevalstats,
	}), nil
}

// GetOpenCypherQueryStatus returns the status and timing of an openCypher query.
func (h *AdminHandler) GetOpenCypherQueryStatus(ctx context.Context, req *connect.Request[pb.GetOpenCypherQueryStatusInput]) (*connect.Response[pb.GetOpenCypherQueryStatusOutput], error) {
	out, err := h.queryStatusPb(req.Header(), req.Msg.Queryid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(out), nil
}

// ListGremlinQueries returns the status of all accepted and running Gremlin queries.
func (h *AdminHandler) ListGremlinQueries(ctx context.Context, req *connect.Request[pb.ListGremlinQueriesInput]) (*connect.Response[pb.ListGremlinQueriesOutput], error) {
	queries, accepted, running, err := h.queryListPb(req.Header(), "gremlin")
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.ListGremlinQueriesOutput{
		Queries:            queries,
		Acceptedquerycount: fmt.Sprintf("%d", accepted),
		Runningquerycount:  fmt.Sprintf("%d", running),
	}), nil
}

// ListOpenCypherQueries returns the status of all accepted and running openCypher queries.
func (h *AdminHandler) ListOpenCypherQueries(ctx context.Context, req *connect.Request[pb.ListOpenCypherQueriesInput]) (*connect.Response[pb.ListOpenCypherQueriesOutput], error) {
	queries, accepted, running, err := h.queryListPb(req.Header(), "opencypher")
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.ListOpenCypherQueriesOutput{
		Queries:            queries,
		Acceptedquerycount: fmt.Sprintf("%d", accepted),
		Runningquerycount:  fmt.Sprintf("%d", running),
	}), nil
}

// GetLoaderJobStatus returns the status of a bulk loader job.
func (h *AdminHandler) GetLoaderJobStatus(ctx context.Context, req *connect.Request[pb.GetLoaderJobStatusInput]) (*connect.Response[pb.GetLoaderJobStatusOutput], error) {
	out, err := h.loaderJobStatusPb(req.Header(), req.Msg.Loadid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(out), nil
}

// ListLoaderJobs returns the IDs of all known bulk loader jobs.
func (h *AdminHandler) ListLoaderJobs(ctx context.Context, req *connect.Request[pb.ListLoaderJobsInput]) (*connect.Response[pb.ListLoaderJobsOutput], error) {
	out, err := h.loaderJobListPb(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(out), nil
}

// GetPropertygraphStatistics returns node and edge counts for the property graph.
func (h *AdminHandler) GetPropertygraphStatistics(ctx context.Context, req *connect.Request[pbcommon.Empty]) (*connect.Response[pb.GetPropertygraphStatisticsOutput], error) {
	out, err := h.propertygraphStatisticsPb(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(out), nil
}

// GetPropertygraphSummary returns a summary of the property graph metadata.
func (h *AdminHandler) GetPropertygraphSummary(ctx context.Context, req *connect.Request[pb.GetPropertygraphSummaryInput]) (*connect.Response[pb.GetPropertygraphSummaryOutput], error) {
	out, err := h.propertygraphSummaryPb(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(out), nil
}

// GetMLDataProcessingJob returns the status of an ML data processing job.
// ML operations are not supported by vorpalstacks.
func (h *AdminHandler) GetMLDataProcessingJob(ctx context.Context, req *connect.Request[pb.GetMLDataProcessingJobInput]) (*connect.Response[pb.GetMLDataProcessingJobOutput], error) {
	return nil, errUnsupportedAdmin("ML data processing")
}

// ListMLDataProcessingJobs lists all ML data processing job IDs.
// ML operations are not supported by vorpalstacks.
func (h *AdminHandler) ListMLDataProcessingJobs(ctx context.Context, req *connect.Request[pb.ListMLDataProcessingJobsInput]) (*connect.Response[pb.ListMLDataProcessingJobsOutput], error) {
	return nil, errUnsupportedAdmin("ML data processing")
}

// GetMLModelTrainingJob returns the status of an ML model training job.
// ML operations are not supported by vorpalstacks.
func (h *AdminHandler) GetMLModelTrainingJob(ctx context.Context, req *connect.Request[pb.GetMLModelTrainingJobInput]) (*connect.Response[pb.GetMLModelTrainingJobOutput], error) {
	return nil, errUnsupportedAdmin("ML model training")
}

// ListMLModelTrainingJobs lists all ML model training job IDs.
// ML operations are not supported by vorpalstacks.
func (h *AdminHandler) ListMLModelTrainingJobs(ctx context.Context, req *connect.Request[pb.ListMLModelTrainingJobsInput]) (*connect.Response[pb.ListMLModelTrainingJobsOutput], error) {
	return nil, errUnsupportedAdmin("ML model training")
}

// GetMLModelTransformJob returns the status of an ML model transform job.
// ML operations are not supported by vorpalstacks.
func (h *AdminHandler) GetMLModelTransformJob(ctx context.Context, req *connect.Request[pb.GetMLModelTransformJobInput]) (*connect.Response[pb.GetMLModelTransformJobOutput], error) {
	return nil, errUnsupportedAdmin("ML model transform")
}

// ListMLModelTransformJobs lists all ML model transform job IDs.
// ML operations are not supported by vorpalstacks.
func (h *AdminHandler) ListMLModelTransformJobs(ctx context.Context, req *connect.Request[pb.ListMLModelTransformJobsInput]) (*connect.Response[pb.ListMLModelTransformJobsOutput], error) {
	return nil, errUnsupportedAdmin("ML model transform")
}

// GetMLEndpoint returns the status of an ML endpoint.
// ML operations are not supported by vorpalstacks.
func (h *AdminHandler) GetMLEndpoint(ctx context.Context, req *connect.Request[pb.GetMLEndpointInput]) (*connect.Response[pb.GetMLEndpointOutput], error) {
	return nil, errUnsupportedAdmin("ML endpoint")
}

// ListMLEndpoints lists all ML endpoint IDs.
// ML operations are not supported by vorpalstacks.
func (h *AdminHandler) ListMLEndpoints(ctx context.Context, req *connect.Request[pb.ListMLEndpointsInput]) (*connect.Response[pb.ListMLEndpointsOutput], error) {
	return nil, errUnsupportedAdmin("ML endpoint")
}

// GetSparqlStatistics returns SPARQL graph statistics.
// SPARQL operations are not supported by vorpalstacks.
func (h *AdminHandler) GetSparqlStatistics(ctx context.Context, req *connect.Request[pbcommon.Empty]) (*connect.Response[pb.GetSparqlStatisticsOutput], error) {
	return nil, errUnsupportedAdmin("SPARQL")
}

// GetRDFGraphSummary returns an RDF graph summary.
// SPARQL operations are not supported by vorpalstacks.
func (h *AdminHandler) GetRDFGraphSummary(ctx context.Context, req *connect.Request[pb.GetRDFGraphSummaryInput]) (*connect.Response[pb.GetRDFGraphSummaryOutput], error) {
	return nil, errUnsupportedAdmin("SPARQL")
}

// GetSparqlStream returns an RDF change data stream.
// SPARQL operations are not supported by vorpalstacks.
func (h *AdminHandler) GetSparqlStream(ctx context.Context, req *connect.Request[pb.GetSparqlStreamInput]) (*connect.Response[pb.GetSparqlStreamOutput], error) {
	return nil, errUnsupportedAdmin("SPARQL")
}

// GetPropertygraphStream returns a property graph change data stream.
// Streaming operations are not supported via the admin console; use the
// HTTP API endpoint directly.
func (h *AdminHandler) GetPropertygraphStream(ctx context.Context, req *connect.Request[pb.GetPropertygraphStreamInput]) (*connect.Response[pb.GetPropertygraphStreamOutput], error) {
	return nil, errUnsupportedAdmin("streaming")
}

// NewConnectHandler creates a gRPC-Web connect handler for the NeptuneData admin console.
func NewConnectHandler(svc *NeptuneDataService) (string, http.Handler) {
	return neptunedataconnect.NewNeptunedataServiceHandler(NewAdminHandler(svc))
}
