package neptunedata

// This file is the sole admin handler file that imports the store package.
// It contains store-accessing wrapper methods used by admin_handler.go,
// plus pure proto conversion helpers that translate store types to proto
// types for response marshalling.

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"vorpalstacks/internal/common/defaults"

	"google.golang.org/protobuf/types/known/timestamppb"

	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/neptunedata"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
)

// getStore returns the per-region Neptune store for the given header.
func (h *AdminHandler) getStore(header http.Header) (*neptunestore.NeptuneStore, error) {
	region := defaults.GetRegionFromHeader(header)
	return h.service.GetStoreForRegion(region)
}

// computeElapsedMillis returns the elapsed time in milliseconds for a query.
// For completed queries this is EndTime - StartTime. For running queries
// (EndTime nil but StartTime set) it is time.Since(StartTime). Returns 0
// when StartTime is nil.
//
// The parameters must be concrete *timestamppb.Timestamp pointers (not
// interface-typed) so that a nil EndTime is correctly detected. A nil
// pointer wrapped in an interface would pass an != nil check, causing
// AsTime to return the epoch (1970) and producing a huge negative value.
func computeElapsedMillis(startTime, endTime *timestamppb.Timestamp) int64 {
	if startTime == nil {
		return 0
	}
	if endTime != nil {
		return endTime.AsTime().Sub(startTime.AsTime()).Milliseconds()
	}
	return time.Since(startTime.AsTime()).Milliseconds()
}

// ---------------------------------------------------------------------------
// Store-accessing wrapper methods — each returns proto-ready types so that
// admin_handler.go needs zero store imports.
// ---------------------------------------------------------------------------

// queryStatusPb retrieves a single query by ID and returns a proto response
// suitable for both Gremlin and OpenCypher query status endpoints.
func (h *AdminHandler) queryStatusPb(header http.Header, queryId string) (*pb.GetOpenCypherQueryStatusOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}

	qr, err := store.GetQuery(queryId)
	if err != nil {
		return nil, err
	}
	if qr == nil {
		return nil, queryNotFound(queryId)
	}

	var elapsedStr string
	if qr.StartTime != nil {
		elapsed := computeElapsedMillis(qr.StartTime, qr.EndTime)
		elapsedStr = strconv.FormatInt(elapsed, 10)
	}

	return &pb.GetOpenCypherQueryStatusOutput{
		Queryid:     qr.GetQueryId(),
		Querystring: qr.GetQueryString(),
		Queryevalstats: &pb.QueryEvalStats{
			Elapsed: elapsedStr,
		},
	}, nil
}

// queryListPb lists queries of the given type and returns proto-ready results.
func (h *AdminHandler) queryListPb(header http.Header, queryType string) ([]*pb.GremlinQueryStatus, int, int, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, 0, 0, err
	}

	queries, err := store.ListQueries()
	if err != nil {
		return nil, 0, 0, err
	}

	pbQueries := make([]*pb.GremlinQueryStatus, 0)
	accepted := 0
	running := 0
	for _, qr := range queries {
		if qr.GetQueryType() != queryType {
			continue
		}
		st := qr.GetStatus()
		if st == "complete" || st == "failed" || st == "cancelled" {
			continue
		}
		accepted++
		if st == "running" {
			running++
		}
		status := &pb.GremlinQueryStatus{
			Queryid:     qr.GetQueryId(),
			Querystring: qr.GetQueryString(),
		}
		if qr.StartTime != nil {
			elapsed := computeElapsedMillis(qr.StartTime, qr.EndTime)
			status.Queryevalstats = &pb.QueryEvalStats{
				Elapsed: strconv.FormatInt(elapsed, 10),
			}
		}
		pbQueries = append(pbQueries, status)
	}

	return pbQueries, accepted, running, nil
}

// loaderJobStatusPb retrieves a single loader job by ID and returns a proto
// response for the admin console.
func (h *AdminHandler) loaderJobStatusPb(header http.Header, loadId string) (*pb.GetLoaderJobStatusOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}

	job, err := store.GetLoaderJob(loadId)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, bulkLoadNotFound(loadId)
	}

	return &pb.GetLoaderJobStatusOutput{
		Status:  job.GetStatus(),
		Payload: fmt.Sprintf(`{"loadId":"%s","status":"%s"}`, job.GetLoadId(), job.GetStatus()),
	}, nil
}

// loaderJobListPb lists all loader job IDs and returns a proto response.
func (h *AdminHandler) loaderJobListPb(header http.Header) (*pb.ListLoaderJobsOutput, error) {
	store, err := h.getStore(header)
	if err != nil {
		return nil, err
	}

	jobs, err := store.ListLoaderJobs()
	if err != nil {
		return nil, err
	}

	loadIds := make([]string, 0, len(jobs))
	for _, job := range jobs {
		loadIds = append(loadIds, job.GetLoadId())
	}

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
