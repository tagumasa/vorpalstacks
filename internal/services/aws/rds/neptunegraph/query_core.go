package neptunegraph

// Query-plane Core functions: the single validation and persistence path for
// the query record and graph summary operations, plus the ExecuteQuery
// execution flow shared with cross-service callers.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
	"vorpalstacks/pkg/cypherparser"
)

// GetQueryInput carries the wire-parsed GetQuery request.
type GetQueryInput struct {
	GraphIdentifier string
	QueryID         string
}

// ListQueriesInput carries the wire-parsed ListQueries request with the
// already-clamped page size and the upper-cased state filter.
type ListQueriesInput struct {
	GraphIdentifier string
	MaxResults      int
	State           string
}

// CancelQueryInput carries the wire-parsed CancelQuery request.
type CancelQueryInput struct {
	GraphIdentifier string
	QueryID         string
}

// GetGraphSummaryInput carries the wire-parsed GetGraphSummary request with
// the upper-cased mode.
type GetGraphSummaryInput struct {
	GraphIdentifier string
	Mode            string
}

// GetGraphSummaryResult carries the computed graph statistics and the
// statistics computation timestamp.
type GetGraphSummaryResult struct {
	Summary   *ngstore.GraphDataSummary
	StatsTime time.Time
}

// getQueryCore retrieves a stored query record by identifier.
func (s *NeptuneGraphService) getQueryCore(store *ngstore.NeptuneGraphStore, in *GetQueryInput) (*ngstore.QueryRecord, error) {
	queryID := in.QueryID
	if queryID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "queryId")
	}

	graphID := in.GraphIdentifier
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

	return query, nil
}

// listQueriesCore validates the state filter and returns query records for
// the graph, either the raw page or the state-filtered set.
func (s *NeptuneGraphService) listQueriesCore(store *ngstore.NeptuneGraphStore, in *ListQueriesInput) ([]*ngstore.QueryRecord, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	maxResults := in.MaxResults

	stateFilter := in.State
	if err := validateQueryStateInput(stateFilter); err != nil {
		return nil, err
	}
	if stateFilter != "" && stateFilter != "ALL" {
		allQueries, err := store.ListQueries(graphID, 0)
		if err != nil {
			return nil, err
		}
		queries := make([]*ngstore.QueryRecord, 0, maxResults)
		for _, q := range allQueries {
			if q.State != stateFilter {
				continue
			}
			queries = append(queries, q)
			if len(queries) >= maxResults {
				break
			}
		}
		return queries, nil
	}

	return store.ListQueries(graphID, maxResults)
}

// cancelQueryCore cancels a running query by transitioning it to CANCELLING
// state. Per Smithy QueryState model, only RUNNING/WAITING/CANCELLING are
// valid states. Terminal queries are deleted (not stored with a terminal
// state), so a query that still exists is either RUNNING or CANCELLING.
func (s *NeptuneGraphService) cancelQueryCore(store *ngstore.NeptuneGraphStore, in *CancelQueryInput) (*ngstore.QueryRecord, error) {
	queryID := in.QueryID
	if queryID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "queryId")
	}

	graphID := in.GraphIdentifier
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

	// If already CANCELLING, return idempotently.
	if query.State == "CANCELLING" || query.State == "WAITING" {
		return query, nil
	}

	query.State = "CANCELLING"
	if err := store.UpdateQuery(query); err != nil {
		logs.Warn("failed to cancel query", logs.String("queryId", queryID), logs.Err(err))
		return nil, err
	}

	return query, nil
}

// getGraphSummaryCore resolves the graph, acquires its engine and computes
// the structural statistics, honouring the DETAILED mode property walk.
func (s *NeptuneGraphService) getGraphSummaryCore(store *ngstore.NeptuneGraphStore, in *GetGraphSummaryInput) (*GetGraphSummaryResult, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	graph, err := s.resolveGraphIdentifier(store, graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}
	graphID = graph.Id

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
		NumNodes:      proto.Int64(stats.NodeCount),
		NumEdges:      proto.Int64(stats.EdgeCount),
		NumNodeLabels: proto.Int64(int64(len(stats.LabelCounts))),
		NumEdgeLabels: proto.Int64(int64(len(stats.RelCounts))),
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

	// DETAILED mode: compute property statistics by iterating all nodes
	// and edges. This is O(n) but GetGraphSummary is not a hot path.
	mode := in.Mode
	if err := validateGraphSummaryMode(mode); err != nil {
		return nil, err
	}
	if mode == "DETAILED" {
		populateDetailedStats(entry.db, summary)
	}

	return &GetGraphSummaryResult{Summary: summary, StatsTime: now}, nil
}

// executeCypherQuery handles the full ExecuteQuery flow: parse, dispatch,
// track query state, and return results.
func executeCypherQuery(ctx context.Context, s *NeptuneGraphService, reqCtx *request.RequestContext, req *request.ParsedRequest, graphID string, entry *engineEntry, store *ngstore.NeptuneGraphStore) (interface{}, error) {
	var params struct {
		Query                    string          `json:"query"`
		Language                 string          `json:"language"`
		Parameters               json.RawMessage `json:"parameters"`
		PlanCache                string          `json:"planCache"`
		ExplainMode              string          `json:"explain"`
		QueryTimeoutMilliseconds int             `json:"queryTimeoutMilliseconds"`
	}
	if err := json.Unmarshal(req.Body, &params); err != nil {
		return nil, newValidationException("BAD_REQUEST", "invalid request body")
	}
	if params.Query == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "query is required")
	}

	const maxQueryBytes = 1 << 20
	if len(params.Query) > maxQueryBytes {
		return nil, newValidationException("QUERY_TOO_LARGE", fmt.Sprintf("query exceeds %d byte limit", maxQueryBytes))
	}

	lang := strings.ToUpper(params.Language)
	if err := validateQueryLanguage(lang); err != nil {
		return nil, err
	}

	if err := validatePlanCache(strings.ToUpper(params.PlanCache)); err != nil {
		return nil, err
	}

	if params.ExplainMode != "" {
		if err := validateExplainMode(strings.ToUpper(params.ExplainMode)); err != nil {
			return nil, err
		}
	}

	var cypherParams map[string]any
	if len(params.Parameters) > 0 {
		if err := json.Unmarshal(params.Parameters, &cypherParams); err != nil {
			return nil, newValidationException("BAD_REQUEST", fmt.Sprintf("invalid parameters: %v", err))
		}
	}

	queryID := generateID("q-")
	now := time.Now().UTC()

	queryRecord := &ngstore.QueryRecord{
		Id:          queryID,
		QueryString: params.Query,
		Language:    lang,
		State:       "RUNNING",
		GraphId:     graphID,
		StartedAt:   now,
	}
	if err := store.CreateQuery(queryRecord); err != nil {
		logs.Warn("failed to create query record", logs.Err(err))
	}

	finaliseQuery := func() {
		elapsed := int32(time.Since(now).Milliseconds())
		// Update elapsed time before removing the record.
		_ = store.TryAdvanceQuery(graphID, queryID, "RUNNING", func(q *ngstore.QueryRecord) {
			q.Elapsed = elapsed
		})
		// Terminal queries are removed from the store per Smithy
		// QueryState model — RUNNING, WAITING, and CANCELLING are the
		// only defined enum values. Completed/failed queries are not
		// retrievable via GetQuery (returns 404), matching AWS behaviour.
		if err := store.DeleteQuery(graphID, queryID); err != nil && !ngstore.IsNotFound(err) {
			logs.Warn("failed to delete query record after completion", logs.Err(err))
		}
	}

	usePlanCache := params.PlanCache == "" ||
		strings.ToUpper(params.PlanCache) == "ENABLED" ||
		strings.ToUpper(params.PlanCache) == "AUTO"

	var parsed *cypherparser.ParsedCypher
	var err error
	if usePlanCache {
		cacheKey := planCacheKey(graphID, params.Query, cypherParams)
		if cached, ok := s.planCache.get(cacheKey); ok {
			if ast, ok := cached.(*cypherparser.ParsedCypher); ok {
				parsed = ast
			}
		}
		if parsed == nil {
			parsed, err = cypherparser.Parse(params.Query)
			if err != nil {
				finaliseQuery()
				return nil, newValidationException("MALFORMED_QUERY", err.Error())
			}
			s.planCache.put(cacheKey, parsed)
		}
	} else {
		parsed, err = cypherparser.Parse(params.Query)
		if err != nil {
			finaliseQuery()
			return nil, newValidationException("MALFORMED_QUERY", err.Error())
		}
	}

	// Apply query timeout if specified. When the deadline is exceeded,
	// return UnprocessableException with QUERY_TIMEOUT reason.
	if params.QueryTimeoutMilliseconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(params.QueryTimeoutMilliseconds)*time.Millisecond)
		defer cancel()
	}

	var execErr error
	var result interface{}

	switch {
	case parsed.Call != nil:
		result, execErr = executeCallQuery(parsed.Call, entry, store, graphID)
	case parsed.DDL != nil:
		var r *cypherparser.CypherResult
		r, execErr = cypherparser.ExecuteDDL(ctx, entry.db, parsed.DDL)
		if r != nil {
			result = map[string]interface{}{"results": r}
		}
	case parsed.Write != nil:
		var r *cypherparser.CypherResult
		r, execErr = cypherparser.ExecuteWrite(ctx, entry.db, entry.db, parsed.Write, cypherParams)
		if r != nil {
			result = map[string]interface{}{"results": r}
		}
	case parsed.Merge != nil:
		var r *cypherparser.CypherResult
		r, execErr = cypherparser.ExecuteMerge(ctx, entry.db, entry.db, parsed.Merge, cypherParams)
		if r != nil {
			result = map[string]interface{}{"results": r}
		}
	case parsed.Read != nil:
		if len(parsed.Read.Set) > 0 || len(parsed.Read.Delete) > 0 || len(parsed.Read.Remove) > 0 || parsed.Read.Create != nil {
			var r *cypherparser.CypherResult
			r, execErr = cypherparser.ExecuteQueryWrite(ctx, entry.db, entry.db, parsed.Read, cypherParams)
			if r != nil {
				result = map[string]interface{}{"results": r}
			}
		} else {
			var r *cypherparser.CypherResult
			r, execErr = cypherparser.Execute(ctx, entry.db, parsed.Read, cypherParams)
			if r != nil {
				result = map[string]interface{}{"results": r}
			}
		}
	default:
		finaliseQuery()
		return nil, newValidationException("MALFORMED_QUERY", "unsupported query type")
	}

	if execErr != nil {
		finaliseQuery()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, newUnprocessableException("QUERY_TIMEOUT", "query exceeded timeout")
		}
		return nil, newValidationException("MALFORMED_QUERY", execErr.Error())
	}

	// Explain mode: add query plan information to the existing response
	// map. Execution branches already set result = {"results": r}, so we
	// add the "explain" key to the same map instead of re-wrapping.
	// For DDL or other operations that produce no result payload (result
	// is nil), create a minimal response containing only the explain.
	explainMode := strings.ToUpper(params.ExplainMode)
	if explainMode == "DETAILS" || explainMode == "STATIC" {
		if resultMap, ok := result.(map[string]interface{}); ok {
			resultMap["explain"] = buildExplainOutput(parsed, explainMode)
		} else {
			result = map[string]interface{}{
				"explain": buildExplainOutput(parsed, explainMode),
			}
		}
	}

	finaliseQuery()
	return result, nil
}

func executeCallQuery(call *cypherparser.CypherCall, entry *engineEntry, store *ngstore.NeptuneGraphStore, graphID string) (interface{}, error) {
	graph, err := store.GetGraph(graphID)
	if err != nil {
		return nil, newResourceNotFoundException("graph", graphID)
	}

	var dim int32
	if graph.VectorSearchConfiguration != nil {
		dim = graph.VectorSearchConfiguration.Dimension
	}

	pd := newProcedureDispatcher(entry.db, dim)
	result, err := pd.ExecuteCall(call, nil)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "ValidationException") {
			return nil, newValidationException("ILLEGAL_ARGUMENT", strings.TrimPrefix(errStr, "ValidationException: "))
		}
		if strings.Contains(errStr, "ResourceNotFoundException") {
			return nil, newResourceNotFoundException("node", "")
		}
		return nil, newHTTPError(500, "InternalServerException", errStr)
	}

	return map[string]interface{}{"results": result}, nil
}
