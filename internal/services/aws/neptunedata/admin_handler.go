package neptunedata

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"connectrpc.com/connect"

	"vorpalstacks/internal/pb/aws/common"
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

// GetEngineStatus returns the health status and engine version information
// for the Neptune-compatible graph engine.
func (h *AdminHandler) GetEngineStatus(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[pb.GetEngineStatusOutput], error) {
	s := h.service
	s.mu.RLock()
	startTime := s.startTime.UTC().Format("2006-01-02T15:04:05.000Z")
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
	s := h.service
	queryId := req.Msg.Queryid

	qr, err := s.store.GetQuery(queryId)
	if err != nil || qr == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("query not found: %s", queryId))
	}

	output := &pb.GetGremlinQueryStatusOutput{
		Queryid:     qr.GetQueryId(),
		Querystring: qr.GetQueryType(),
	}

	if qr.EndTime != nil && qr.StartTime != nil {
		elapsed := qr.EndTime.AsTime().Sub(qr.StartTime.AsTime()).Milliseconds()
		output.Queryevalstats = &pb.QueryEvalStats{
			Elapsed: fmt.Sprintf("%d", elapsed),
		}
	}

	return connect.NewResponse(output), nil
}

// GetOpenCypherQueryStatus returns the status and timing of an openCypher query.
func (h *AdminHandler) GetOpenCypherQueryStatus(ctx context.Context, req *connect.Request[pb.GetOpenCypherQueryStatusInput]) (*connect.Response[pb.GetOpenCypherQueryStatusOutput], error) {
	s := h.service
	queryId := req.Msg.Queryid

	qr, err := s.store.GetQuery(queryId)
	if err != nil || qr == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("query not found: %s", queryId))
	}

	output := &pb.GetOpenCypherQueryStatusOutput{
		Queryid:     qr.GetQueryId(),
		Querystring: qr.GetQueryType(),
	}

	if qr.EndTime != nil && qr.StartTime != nil {
		elapsed := qr.EndTime.AsTime().Sub(qr.StartTime.AsTime()).Milliseconds()
		output.Queryevalstats = &pb.QueryEvalStats{
			Elapsed: fmt.Sprintf("%d", elapsed),
		}
	}

	return connect.NewResponse(output), nil
}

// ListGremlinQueries returns the status of all accepted and running Gremlin queries.
func (h *AdminHandler) ListGremlinQueries(ctx context.Context, req *connect.Request[pb.ListGremlinQueriesInput]) (*connect.Response[pb.ListGremlinQueriesOutput], error) {
	s := h.service

	queries, err := s.store.ListQueries()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbQueries := make([]*pb.GremlinQueryStatus, 0)
	accepted := 0
	running := 0
	for _, qr := range queries {
		if qr.GetQueryType() != "gremlin" {
			continue
		}
		accepted++
		if qr.GetStatus() == "running" {
			running++
		}
		status := &pb.GremlinQueryStatus{
			Queryid:     qr.GetQueryId(),
			Querystring: qr.GetQueryType(),
		}
		if qr.EndTime != nil && qr.StartTime != nil {
			elapsed := qr.EndTime.AsTime().Sub(qr.StartTime.AsTime()).Milliseconds()
			status.Queryevalstats = &pb.QueryEvalStats{
				Elapsed: fmt.Sprintf("%d", elapsed),
			}
		}
		pbQueries = append(pbQueries, status)
	}

	return connect.NewResponse(&pb.ListGremlinQueriesOutput{
		Queries:            pbQueries,
		Acceptedquerycount: strconv.Itoa(accepted),
		Runningquerycount:  strconv.Itoa(running),
	}), nil
}

// ListOpenCypherQueries returns the status of all accepted and running openCypher queries.
func (h *AdminHandler) ListOpenCypherQueries(ctx context.Context, req *connect.Request[pb.ListOpenCypherQueriesInput]) (*connect.Response[pb.ListOpenCypherQueriesOutput], error) {
	s := h.service

	queries, err := s.store.ListQueries()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbQueries := make([]*pb.GremlinQueryStatus, 0)
	accepted := 0
	running := 0
	for _, qr := range queries {
		if qr.GetQueryType() != "opencypher" {
			continue
		}
		accepted++
		if qr.GetStatus() == "running" {
			running++
		}
		status := &pb.GremlinQueryStatus{
			Queryid:     qr.GetQueryId(),
			Querystring: qr.GetQueryType(),
		}
		if qr.EndTime != nil && qr.StartTime != nil {
			elapsed := qr.EndTime.AsTime().Sub(qr.StartTime.AsTime()).Milliseconds()
			status.Queryevalstats = &pb.QueryEvalStats{
				Elapsed: fmt.Sprintf("%d", elapsed),
			}
		}
		pbQueries = append(pbQueries, status)
	}

	return connect.NewResponse(&pb.ListOpenCypherQueriesOutput{
		Queries:            pbQueries,
		Acceptedquerycount: strconv.Itoa(accepted),
		Runningquerycount:  strconv.Itoa(running),
	}), nil
}

// GetLoaderJobStatus returns the status of a bulk loader job.
func (h *AdminHandler) GetLoaderJobStatus(ctx context.Context, req *connect.Request[pb.GetLoaderJobStatusInput]) (*connect.Response[pb.GetLoaderJobStatusOutput], error) {
	s := h.service
	loadId := req.Msg.Loadid

	job, err := s.store.GetLoaderJob(loadId)
	if err != nil || job == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("loader job not found: %s", loadId))
	}

	return connect.NewResponse(&pb.GetLoaderJobStatusOutput{
		Status:  job.GetStatus(),
		Payload: fmt.Sprintf(`{"loadId":"%s","status":"%s"}`, job.GetLoadId(), job.GetStatus()),
	}), nil
}

// ListLoaderJobs returns the IDs of all known bulk loader jobs.
func (h *AdminHandler) ListLoaderJobs(ctx context.Context, req *connect.Request[pb.ListLoaderJobsInput]) (*connect.Response[pb.ListLoaderJobsOutput], error) {
	s := h.service

	jobs, err := s.store.ListLoaderJobs()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	loadIds := make([]string, 0, len(jobs))
	for _, job := range jobs {
		loadIds = append(loadIds, job.GetLoadId())
	}

	return connect.NewResponse(&pb.ListLoaderJobsOutput{
		Status: "200 OK",
		Payload: &pb.LoaderIdResult{
			Loadids: loadIds,
		},
	}), nil
}

// GetPropertygraphStatistics returns node and edge counts for the property graph.
func (h *AdminHandler) GetPropertygraphStatistics(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[pb.GetPropertygraphStatisticsOutput], error) {
	s := h.service

	s.mu.RLock()
	sigCount := int64(0)
	predCount := int64(0)
	for _, cnt := range s.stats.LabelCounts {
		sigCount += cnt
	}
	for _, cnt := range s.stats.RelCounts {
		predCount += cnt
	}
	stats := &pb.Statistics{
		Active:       "true",
		Autocompute:  "true",
		Date:         time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		Note:         "Automatically computed",
		Statisticsid: "auto-statistics",
		Signatureinfo: &pb.StatisticsSummary{
			Signaturecount: fmt.Sprintf("%d", sigCount),
			Instancecount:  fmt.Sprintf("%d", s.stats.NodeCount),
			Predicatecount: fmt.Sprintf("%d", predCount),
		},
	}
	s.mu.RUnlock()

	return connect.NewResponse(&pb.GetPropertygraphStatisticsOutput{
		Status:  "200 OK",
		Payload: stats,
	}), nil
}

// GetPropertygraphSummary returns a summary of the property graph metadata.
func (h *AdminHandler) GetPropertygraphSummary(ctx context.Context, req *connect.Request[pb.GetPropertygraphSummaryInput]) (*connect.Response[pb.GetPropertygraphSummaryOutput], error) {
	s := h.service

	s.mu.RLock()
	summaryMap := &pb.PropertygraphSummaryValueMap{
		Graphsummary: &pb.PropertygraphSummary{
			Numnodes: fmt.Sprintf("%d", s.stats.NodeCount),
			Numedges: fmt.Sprintf("%d", s.stats.EdgeCount),
		},
	}
	s.mu.RUnlock()

	return connect.NewResponse(&pb.GetPropertygraphSummaryOutput{
		Statuscode: "200 OK",
		Payload:    summaryMap,
	}), nil
}

// GetMLDataProcessingJob returns the status of an ML data processing job.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) GetMLDataProcessingJob(ctx context.Context, req *connect.Request[pb.GetMLDataProcessingJobInput]) (*connect.Response[pb.GetMLDataProcessingJobOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ML operations return HTTP 501"))
}

// ListMLDataProcessingJobs lists all ML data processing job IDs.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) ListMLDataProcessingJobs(ctx context.Context, req *connect.Request[pb.ListMLDataProcessingJobsInput]) (*connect.Response[pb.ListMLDataProcessingJobsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ML operations return HTTP 501"))
}

// GetMLModelTrainingJob returns the status of an ML model training job.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) GetMLModelTrainingJob(ctx context.Context, req *connect.Request[pb.GetMLModelTrainingJobInput]) (*connect.Response[pb.GetMLModelTrainingJobOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ML operations return HTTP 501"))
}

// ListMLModelTrainingJobs lists all ML model training job IDs.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) ListMLModelTrainingJobs(ctx context.Context, req *connect.Request[pb.ListMLModelTrainingJobsInput]) (*connect.Response[pb.ListMLModelTrainingJobsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ML operations return HTTP 501"))
}

// GetMLModelTransformJob returns the status of an ML model transform job.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) GetMLModelTransformJob(ctx context.Context, req *connect.Request[pb.GetMLModelTransformJobInput]) (*connect.Response[pb.GetMLModelTransformJobOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ML operations return HTTP 501"))
}

// ListMLModelTransformJobs lists all ML model transform job IDs.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) ListMLModelTransformJobs(ctx context.Context, req *connect.Request[pb.ListMLModelTransformJobsInput]) (*connect.Response[pb.ListMLModelTransformJobsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ML operations return HTTP 501"))
}

// GetMLEndpoint returns the status of an ML endpoint.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) GetMLEndpoint(ctx context.Context, req *connect.Request[pb.GetMLEndpointInput]) (*connect.Response[pb.GetMLEndpointOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ML operations return HTTP 501"))
}

// ListMLEndpoints lists all ML endpoint IDs.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) ListMLEndpoints(ctx context.Context, req *connect.Request[pb.ListMLEndpointsInput]) (*connect.Response[pb.ListMLEndpointsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ML operations return HTTP 501"))
}

// GetSparqlStatistics returns SPARQL graph statistics.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) GetSparqlStatistics(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[pb.GetSparqlStatisticsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("SPARQL operations return HTTP 501"))
}

// GetRDFGraphSummary returns an RDF graph summary.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) GetRDFGraphSummary(ctx context.Context, req *connect.Request[pb.GetRDFGraphSummaryInput]) (*connect.Response[pb.GetRDFGraphSummaryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("SPARQL operations return HTTP 501"))
}

// GetSparqlStream returns an RDF change data stream.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) GetSparqlStream(ctx context.Context, req *connect.Request[pb.GetSparqlStreamInput]) (*connect.Response[pb.GetSparqlStreamOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("SPARQL operations return HTTP 501"))
}

// GetPropertygraphStream returns a property graph change data stream.
// Not yet implemented; returns HTTP 501.
func (h *AdminHandler) GetPropertygraphStream(ctx context.Context, req *connect.Request[pb.GetPropertygraphStreamInput]) (*connect.Response[pb.GetPropertygraphStreamOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("streaming operations return HTTP 501"))
}
