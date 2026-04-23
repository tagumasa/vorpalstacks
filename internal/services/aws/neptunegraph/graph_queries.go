package neptunegraph

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
)

// ExecuteQuery runs a Cypher query against the specified graph's query engine.
func (s *NeptuneGraphService) ExecuteQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		graphID = request.GetStringParam(req.Parameters, "graphidentifier")
	}
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	s.enginesMu.Lock()
	entry, ok := s.activeEngines[graphID]
	if !ok || entry.stopped {
		s.enginesMu.Unlock()
		return nil, newValidationException("UNSUPPORTED_OPERATION", "graph is not available")
	}
	entry.wg.Add(1)
	s.enginesMu.Unlock()
	defer entry.wg.Done()

	entry.mu.RLock()
	defer entry.mu.RUnlock()

	return executeCypherQuery(ctx, s, reqCtx, req, graphID, entry, store)
}

// GetQuery retrieves the details and results of a previously executed query.
func (s *NeptuneGraphService) GetQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryID := request.GetStringParam(req.Parameters, "queryId")
	if queryID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "queryId")
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		graphID = request.GetStringParam(req.Parameters, "graphidentifier")
	}
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	query, err := store.GetQuery(graphID, queryID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("query", queryID)
		}
		return nil, err
	}

	return queryToResponse(query), nil
}

// ListQueries returns query records for a graph, optionally filtered by state.
func (s *NeptuneGraphService) ListQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		graphID = request.GetStringParam(req.Parameters, "graphidentifier")
	}
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	maxResults := request.GetIntParam(req.Parameters, "maxResults")
	if maxResults < 1 || maxResults > 100 {
		maxResults = 100
	}

	stateFilter := request.GetStringParam(req.Parameters, "state")
	if stateFilter != "" && stateFilter != "ALL" {
		allQueries, err := store.ListQueries(graphID, maxResults*3)
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, maxResults)
		for _, q := range allQueries {
			if q.State != stateFilter {
				continue
			}
			items = append(items, queryToResponse(q))
			if len(items) >= maxResults {
				break
			}
		}
		return map[string]interface{}{
			"queries": items,
		}, nil
	}

	queries, err := store.ListQueries(graphID, maxResults)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(queries))
	for _, q := range queries {
		items = append(items, queryToResponse(q))
	}

	return map[string]interface{}{
		"queries": items,
	}, nil
}

// CancelQuery cancels a running query by transitioning it to CANCELLED state.
// If the query has already reached a terminal state (COMPLETE, FAILED, CANCELLED),
// the current record is returned idempotently.
func (s *NeptuneGraphService) CancelQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryID := request.GetStringParam(req.Parameters, "queryId")
	if queryID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "queryId")
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		graphID = request.GetStringParam(req.Parameters, "graphidentifier")
	}
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	query, err := store.GetQuery(graphID, queryID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("query", queryID)
		}
		return nil, err
	}

	switch query.State {
	case "COMPLETE", "FAILED", "CANCELLED":
		return queryToResponse(query), nil
	}

	query.State = "CANCELLED"
	if err := store.UpdateQuery(query); err != nil {
		logs.Warn("failed to cancel query", logs.String("queryId", queryID), logs.Err(err))
		return nil, err
	}

	return queryToResponse(query), nil
}

// GetGraphSummary returns structural statistics about the specified graph's data.
func (s *NeptuneGraphService) GetGraphSummary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		graphID = request.GetStringParam(req.Parameters, "graphidentifier")
	}
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	s.enginesMu.Lock()
	entry, ok := s.activeEngines[graphID]
	if !ok {
		s.enginesMu.Unlock()
		return nil, newResourceNotFoundException("graph", graphID)
	}
	if entry.stopped {
		s.enginesMu.Unlock()
		return nil, newResourceNotFoundException("graph", graphID)
	}
	entry.wg.Add(1)
	s.enginesMu.Unlock()
	defer entry.wg.Done()

	entry.mu.RLock()
	stats := entry.db.Stats()
	entry.mu.RUnlock()
	now := time.Now().UTC()

	summary := &ngstore.GraphDataSummary{
		NumNodes:      int64Ptr(stats.NodeCount),
		NumEdges:      int64Ptr(stats.EdgeCount),
		NumNodeLabels: int64Ptr(int64(len(stats.LabelCounts))),
		NumEdgeLabels: int64Ptr(int64(len(stats.RelCounts))),
	}

	if len(stats.LabelCounts) > 0 {
		labels := make([]string, 0, len(stats.LabelCounts))
		for label := range stats.LabelCounts {
			labels = append(labels, label)
		}
		summary.NodeLabels = labels
	}

	if len(stats.RelCounts) > 0 {
		labels := make([]string, 0, len(stats.RelCounts))
		for label := range stats.RelCounts {
			labels = append(labels, label)
		}
		summary.EdgeLabels = labels
	}

	return map[string]interface{}{
		"graphSummary":                  graphDataSummaryToResponse(summary),
		"lastStatisticsComputationTime": now.Format("2006-01-02T15:04:05.000Z"),
		"version":                       "1.0",
	}, nil
}
