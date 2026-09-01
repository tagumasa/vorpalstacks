package neptunedata

// This file holds the admin console's Core-backed wrapper methods and pure
// proto conversion helpers. It performs no direct store access: every data
// path routes through the service's Core functions.

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"vorpalstacks/internal/common/defaults"

	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/neptunedata"
)

// ---------------------------------------------------------------------------
// Store-accessing wrapper methods — each returns proto-ready types so that
// admin_handler.go needs zero store imports.
// ---------------------------------------------------------------------------

// queryStatusPb retrieves a single query by ID and returns a proto response
// suitable for both Gremlin and OpenCypher query status endpoints.
func (h *AdminHandler) queryStatusPb(header http.Header, queryId string) (*pb.GetOpenCypherQueryStatusOutput, error) {
	result, err := h.service.getQueryStatusCore(&GetQueryStatusInput{
		QueryId: queryId,
		Region:  defaults.GetRegionFromHeader(header),
	})
	if err != nil {
		return nil, err
	}

	stats, _ := result["queryEvalStats"].(map[string]interface{})
	elapsedStr := strconv.FormatInt(stats["elapsed"].(int64), 10)

	return &pb.GetOpenCypherQueryStatusOutput{
		Queryid:     result["queryId"].(string),
		Querystring: result["queryString"].(string),
		Queryevalstats: &pb.QueryEvalStats{
			Elapsed: elapsedStr,
		},
	}, nil
}

// queryListPb lists queries of the given type and returns proto-ready
// results. The list entries carry the identifiers only; per-query elapsed
// time is available on the per-query status endpoint.
func (h *AdminHandler) queryListPb(header http.Header, queryType string) ([]*pb.GremlinQueryStatus, int, int, error) {
	result, err := h.service.listQueriesCore(&ListQueriesInput{
		QueryType:      queryType,
		IncludeWaiting: true,
		Region:         defaults.GetRegionFromHeader(header),
	})
	if err != nil {
		return nil, 0, 0, err
	}

	entries, _ := result["queries"].([]interface{})
	pbQueries := make([]*pb.GremlinQueryStatus, 0, len(entries))
	for _, e := range entries {
		entry, _ := e.(map[string]interface{})
		pbQueries = append(pbQueries, &pb.GremlinQueryStatus{
			Queryid:     entry["queryId"].(string),
			Querystring: entry["queryString"].(string),
		})
	}
	accepted := int(result["acceptedQueryCount"].(int32)) + int(result["runningQueryCount"].(int32))
	running := int(result["runningQueryCount"].(int32))
	return pbQueries, accepted, running, nil
}

// loaderJobStatusPb retrieves a single loader job status via the shared
// Core and renders the console payload.
func (h *AdminHandler) loaderJobStatusPb(header http.Header, loadId string) (*pb.GetLoaderJobStatusOutput, error) {
	result, err := h.service.getLoaderJobStatusCore(&GetLoaderJobStatusInput{
		LoadId: loadId,
		Region: defaults.GetRegionFromHeader(header),
	})
	if err != nil {
		return nil, err
	}
	status, _ := result["status"].(string)
	return &pb.GetLoaderJobStatusOutput{
		Status:  status,
		Payload: fmt.Sprintf(`{"loadId":"%s","status":"%s"}`, loadId, status),
	}, nil
}

// loaderJobListPb lists loader job IDs via the shared Core.
func (h *AdminHandler) loaderJobListPb(header http.Header) (*pb.ListLoaderJobsOutput, error) {
	result, err := h.service.listLoaderJobsCore(&ListLoaderJobsInput{
		IncludeQueuedLoads: true,
		Region:             defaults.GetRegionFromHeader(header),
	})
	if err != nil {
		return nil, err
	}
	payload, _ := result["payload"].(map[string]interface{})
	rawIds, _ := payload["loadIds"].([]string)
	loadIds := make([]string, 0, len(rawIds))
	loadIds = append(loadIds, rawIds...)

	return &pb.ListLoaderJobsOutput{
		Status: "200 OK",
		Payload: &pb.LoaderIdResult{
			Loadids: loadIds,
		},
	}, nil
}

// propertygraphStatisticsPb computes statistics and returns a proto response
// for the admin console. This method accesses the service state directly
// (no store needed) but is placed here for architectural consistency.
func (h *AdminHandler) propertygraphStatisticsPb(header http.Header) (*pb.GetPropertygraphStatisticsOutput, error) {
	s := h.service
	region := defaults.GetRegionFromHeader(header)

	s.mu.RLock()
	statsDisabled := s.statsDisabled
	autoCompute := s.autoComputeEnabled
	s.mu.RUnlock()

	if statsDisabled {
		return &pb.GetPropertygraphStatisticsOutput{
			Status: "200 OK",
			Payload: &pb.Statistics{
				Active:      "false",
				Autocompute: fmt.Sprintf("%t", autoCompute),
				Note:        "Statistics auto-compute is disabled. Call ManagePropertygraphStatistics with mode 'refresh' or 'enableAutoCompute' to generate statistics.",
			},
		}, nil
	}

	s.refreshStatisticsForRegion(region)
	st := s.getStats(region)
	nodeCount, _, labelCounts, relCounts, _, _ := st.snapshot()
	sigCount := int64(len(labelCounts))
	predCount := int64(len(relCounts))

	stats := &pb.Statistics{
		Active:       "true",
		Autocompute:  fmt.Sprintf("%t", autoCompute),
		Date:         time.Now().UTC().Format(timeutils.ISO8601UTCFormat),
		Note:         "Automatically computed",
		Statisticsid: "auto-statistics",
		Signatureinfo: &pb.StatisticsSummary{
			Signaturecount: strconv.FormatInt(sigCount, 10),
			Instancecount:  strconv.FormatInt(nodeCount, 10),
			Predicatecount: strconv.FormatInt(predCount, 10),
		},
	}

	return &pb.GetPropertygraphStatisticsOutput{
		Status:  "200 OK",
		Payload: stats,
	}, nil
}

// propertygraphSummaryPb computes a graph summary and returns a proto response
// for the admin console.
func (h *AdminHandler) propertygraphSummaryPb(header http.Header) (*pb.GetPropertygraphSummaryOutput, error) {
	s := h.service
	region := defaults.GetRegionFromHeader(header)

	s.mu.RLock()
	statsDisabled := s.statsDisabled
	s.mu.RUnlock()

	if !statsDisabled {
		s.refreshStatisticsForRegion(region)
	}
	st := s.getStats(region)
	nodeCount, edgeCount, _, _, _, _ := st.snapshot()

	summaryMap := &pb.PropertygraphSummaryValueMap{
		Graphsummary: &pb.PropertygraphSummary{
			Numnodes: strconv.FormatInt(nodeCount, 10),
			Numedges: strconv.FormatInt(edgeCount, 10),
		},
	}

	return &pb.GetPropertygraphSummaryOutput{
		Statuscode: "200 OK",
		Payload:    summaryMap,
	}, nil
}
