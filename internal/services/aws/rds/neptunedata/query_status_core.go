package neptunedata

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
)

// GetQueryStatusInput carries the parsed query-status wire payload.
type GetQueryStatusInput struct {
	QueryId string
	Region  string
}

// getQueryStatusCore returns the status and evaluation statistics of a query
// identified by queryId. Shared by both Gremlin and OpenCypher query status
// handlers.
func (s *NeptuneDataService) getQueryStatusCore(in *GetQueryStatusInput) (map[string]interface{}, error) {
	if in.QueryId == "" {
		return nil, missingParameter("queryId")
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	qr, err := store.GetQuery(in.QueryId)
	if err != nil || qr == nil {
		return nil, badRequest(fmt.Sprintf("query not found: %s", in.QueryId))
	}

	var elapsed int64
	if qr.StartTime != nil {
		if qr.EndTime != nil {
			elapsed = qr.EndTime.AsTime().Sub(qr.StartTime.AsTime()).Milliseconds()
		} else {
			elapsed = time.Since(qr.StartTime.AsTime()).Milliseconds()
		}
	}

	return map[string]interface{}{
		"queryId":     qr.GetQueryId(),
		"queryString": qr.GetQueryString(),
		"queryEvalStats": map[string]interface{}{
			"cancelled":  qr.GetStatus() == "cancelled",
			"elapsed":    elapsed,
			"waited":     0,
			"subqueries": []interface{}{},
		},
	}, nil
}

// ListQueriesInput carries the parsed list-queries wire payload.
type ListQueriesInput struct {
	QueryType      string
	IncludeWaiting bool
	Region         string
}

// listQueriesCore returns all submitted queries of the given type, optionally
// including those in a waiting state. Shared by both Gremlin and OpenCypher
// list queries handlers.
func (s *NeptuneDataService) listQueriesCore(in *ListQueriesInput) (map[string]interface{}, error) {
	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	queries, err := store.ListQueries()
	if err != nil {
		return nil, err
	}

	var result []interface{}
	var acceptedCount, runningCount int32

	for _, qr := range queries {
		if qr.GetQueryType() != in.QueryType {
			continue
		}
		st := qr.GetStatus()
		if st == "complete" || st == "failed" || st == "cancelled" {
			continue
		}
		if st == "waiting" && !in.IncludeWaiting {
			continue
		}
		entry := map[string]interface{}{
			"queryId":     qr.GetQueryId(),
			"queryString": qr.GetQueryString(),
		}
		if st == "running" {
			runningCount++
		} else {
			acceptedCount++
		}
		result = append(result, entry)
	}
	if result == nil {
		result = []interface{}{}
	}

	return map[string]interface{}{
		"queries":            result,
		"acceptedQueryCount": acceptedCount,
		"runningQueryCount":  runningCount,
	}, nil
}

// CancelQueryInput carries the parsed cancel-query wire payload. If Silent is
// true, an empty body is returned instead of the standard cancellation
// confirmation.
type CancelQueryInput struct {
	QueryId        string
	Silent         bool
	IncludePayload bool
	Region         string
}

// cancelQueryCore cancels a running query identified by queryId. Shared by
// both Gremlin and OpenCypher cancel handlers.
func (s *NeptuneDataService) cancelQueryCore(in *CancelQueryInput) (map[string]interface{}, error) {
	if in.QueryId == "" {
		return nil, missingParameter("queryId")
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	qr, err := store.GetQuery(in.QueryId)
	if err != nil || qr == nil {
		return nil, badRequest(fmt.Sprintf("query not found: %s", in.QueryId))
	}
	switch qr.GetStatus() {
	case "complete", "failed", "cancelled":
		return nil, badRequest(fmt.Sprintf("cannot cancel query in terminal state: %s", qr.GetStatus()))
	}
	qr.Status = "cancelled"
	qr.EndTime = timestamppb.Now()
	if err := store.UpdateQuery(qr); err != nil {
		logs.Warn("failed to persist query cancellation", logs.String("queryId", in.QueryId), logs.Err(err))
	}

	if in.Silent {
		return map[string]interface{}{}, nil
	}

	resp := map[string]interface{}{
		"status": "200 OK",
	}
	if in.IncludePayload {
		resp["payload"] = true
	}
	return resp, nil
}

// trackQuery records the start of a query execution in Pebble storage.
func (s *NeptuneDataService) trackQuery(store *neptunestore.NeptuneStore, id, query, queryType string) {
	if queryType != "gremlin" && queryType != "opencypher" && queryType != "sparql" {
		queryType = "unknown"
	}
	qr := &pb.QueryState{
		QueryId:     id,
		QueryType:   queryType,
		QueryString: query,
		Status:      "running",
		StartTime:   timestamppb.Now(),
	}
	if err := store.CreateQuery(qr); err != nil {
		logs.Warn("failed to track query", logs.String("queryId", id), logs.Err(err))
	}
}

// resolveQuery records the completion of a query in Pebble storage.
func (s *NeptuneDataService) resolveQuery(store *neptunestore.NeptuneStore, id string, result any, err error) {
	qr, storeErr := store.GetQuery(id)
	if storeErr != nil || qr == nil {
		return
	}
	qr.EndTime = timestamppb.Now()

	if err != nil {
		qr.Status = "failed"
		qr.Error = err.Error()
	} else {
		qr.Status = "complete"
	}
	if updateErr := store.UpdateQuery(qr); updateErr != nil {
		logs.Warn("failed to resolve query", logs.String("queryId", id), logs.Err(updateErr))
	}
}
