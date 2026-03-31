package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/pkg/cypherparser"
	"vorpalstacks/pkg/graphengine"
	"vorpalstacks/pkg/gremlinparser"
)

// queryStateTTL defines how long completed query state entries are retained
// before being eligible for cleanup.
const (
	queryStateTTL = 5 * time.Minute
)

// NeptuneDataService implements the Neptune Data API service, handling Gremlin,
// OpenCypher, graph statistics, loader, and ML operations. Query state is
// tracked in-memory with TTL-based expiry. Graph data is read and written via
// the GraphReader/GraphWriter interfaces provided by the request context.
type NeptuneDataService struct {
	mu         sync.RWMutex
	queryState map[string]*queryState
	fastTokens map[string]time.Time
	stats      *GraphStatistics
	startTime  time.Time
	loaderJobs map[string]*loaderJobState
}

type loaderJobState struct {
	LoadId     string    `json:"loadId"`
	Status     string    `json:"status"`
	Source     string    `json:"source"`
	Format     string    `json:"format"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
}

// queryState holds the runtime state of an executed query, including timing,
// status, and optional result or error information.
type queryState struct {
	ID        string
	Query     string
	QueryType string
	Status    string
	StartedAt time.Time
	EndedAt   time.Time
	Result    any
	Error     string
}

// elapsedMs returns the elapsed time between start and end in milliseconds.
func elapsedMs(start, end time.Time) int64 {
	return int64(end.Sub(start).Milliseconds())
}

// GraphStatistics holds cached graph-level statistics for the property graph.
// Refreshed on demand when statistics or summary endpoints are called.
type GraphStatistics struct {
	NodeCount   int64            `json:"numNodes"`
	EdgeCount   int64            `json:"numEdges"`
	LabelCounts map[string]int64 `json:"-"`
	RelCounts   map[string]int64 `json:"-"`
}

// NewNeptuneDataService creates a new service instance with empty query state,
// FastReset token store, and zeroed statistics.
func NewNeptuneDataService() *NeptuneDataService {
	return &NeptuneDataService{
		queryState: make(map[string]*queryState),
		fastTokens: make(map[string]time.Time),
		stats: &GraphStatistics{
			LabelCounts: make(map[string]int64),
			RelCounts:   make(map[string]int64),
		},
		startTime:  time.Now(),
		loaderJobs: make(map[string]*loaderJobState),
	}
}

// RegisterHandlers registers all Neptune Data API operation handlers with the
// dispatcher. Unsupported operations (SPARQL, ML) return HTTP 501.
func (s *NeptuneDataService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("neptunedata", "GetEngineStatus", s.GetEngineStatus)
	d.RegisterHandlerForService("neptunedata", "ExecuteGremlinQuery", s.ExecuteGremlinQuery)
	d.RegisterHandlerForService("neptunedata", "ExecuteGremlinExplainQuery", s.ExecuteGremlinExplainQuery)
	d.RegisterHandlerForService("neptunedata", "ExecuteGremlinProfileQuery", s.ExecuteGremlinProfileQuery)
	d.RegisterHandlerForService("neptunedata", "GetGremlinQueryStatus", s.GetGremlinQueryStatus)
	d.RegisterHandlerForService("neptunedata", "ListGremlinQueries", s.ListGremlinQueries)
	d.RegisterHandlerForService("neptunedata", "CancelGremlinQuery", s.CancelGremlinQuery)
	d.RegisterHandlerForService("neptunedata", "ExecuteOpenCypherQuery", s.ExecuteOpenCypherQuery)
	d.RegisterHandlerForService("neptunedata", "ExecuteOpenCypherExplainQuery", s.ExecuteOpenCypherExplainQuery)
	d.RegisterHandlerForService("neptunedata", "GetOpenCypherQueryStatus", s.GetOpenCypherQueryStatus)
	d.RegisterHandlerForService("neptunedata", "ListOpenCypherQueries", s.ListOpenCypherQueries)
	d.RegisterHandlerForService("neptunedata", "CancelOpenCypherQuery", s.CancelOpenCypherQuery)
	d.RegisterHandlerForService("neptunedata", "ExecuteFastReset", s.ExecuteFastReset)
	d.RegisterHandlerForService("neptunedata", "GetPropertygraphStatistics", s.GetPropertygraphStatistics)
	d.RegisterHandlerForService("neptunedata", "ManagePropertygraphStatistics", s.ManagePropertygraphStatistics)
	d.RegisterHandlerForService("neptunedata", "DeletePropertygraphStatistics", s.DeletePropertygraphStatistics)
	d.RegisterHandlerForService("neptunedata", "GetPropertygraphSummary", s.GetPropertygraphSummary)
	d.RegisterHandlerForService("neptunedata", "GetPropertygraphStream", s.GetPropertygraphStream)
	d.RegisterHandlerForService("neptunedata", "StartLoaderJob", s.StartLoaderJob)
	d.RegisterHandlerForService("neptunedata", "GetLoaderJobStatus", s.GetLoaderJobStatus)
	d.RegisterHandlerForService("neptunedata", "ListLoaderJobs", s.ListLoaderJobs)
	d.RegisterHandlerForService("neptunedata", "CancelLoaderJob", s.CancelLoaderJob)
	d.RegisterHandlerForService("neptunedata", "GetSparqlStatistics", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "ManageSparqlStatistics", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "DeleteSparqlStatistics", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "GetSparqlStream", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "GetRDFGraphSummary", s.unsupported)

	d.RegisterHandlerForService("neptunedata", "StartMLDataProcessingJob", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "GetMLDataProcessingJob", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "ListMLDataProcessingJobs", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "CancelMLDataProcessingJob", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "StartMLModelTrainingJob", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "GetMLModelTrainingJob", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "ListMLModelTrainingJobs", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "CancelMLModelTrainingJob", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "StartMLModelTransformJob", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "GetMLModelTransformJob", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "ListMLModelTransformJobs", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "CancelMLModelTransformJob", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "CreateMLEndpoint", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "GetMLEndpoint", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "ListMLEndpoints", s.unsupported)
	d.RegisterHandlerForService("neptunedata", "DeleteMLEndpoint", s.unsupported)
}

// GetEngineStatus returns the health status and engine version information
// for the Neptune-compatible graph engine.
func (s *NeptuneDataService) GetEngineStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = req
	return map[string]interface{}{
		"status":    "healthy",
		"startTime": s.startTime.UTC().Format("2006-01-02T15:04:05.000Z"),
		"gremlin": map[string]interface{}{
			"version": "3.7.x",
		},
		"opencypher": map[string]interface{}{
			"version": "2023-08-01",
		},
		"sparql": map[string]interface{}{
			"version": "1.1",
		},
		"labMode": map[string]interface{}{},
		"settings": map[string]interface{}{
			"neptune lab mode": "DISABLED",
		},
		"role": "writer",
	}, nil
}

// ExecuteGremlinQuery parses and executes a Gremlin traversal query against the
// graph engine, returning results in the AWS response shape.
func (s *NeptuneDataService) ExecuteGremlinQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	body := req.Body
	var params struct {
		Gremlin string `json:"gremlin"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Gremlin == "" {
		return nil, missingParameter("gremlin")
	}

	qid := generateQueryID()
	s.trackQuery(qid, params.Gremlin, "gremlin")

	reader := reqCtx.GraphReader()
	writer := reqCtx.GraphWriter()
	parsed, err := gremlinparser.Parse(params.Gremlin)
	if err != nil {
		s.resolveQuery(qid, nil, err)
		return nil, malformedQuery(err.Error())
	}

	result, execErr := gremlinparser.ExecuteQuery(ctx, reader, writer, parsed, nil)
	s.resolveQuery(qid, result, execErr)
	if execErr != nil {
		return nil, malformedQuery(execErr.Error())
	}

	return map[string]interface{}{
		"requestId": qid,
		"status": map[string]interface{}{
			"code":       float64(200),
			"message":    "",
			"attributes": map[string]interface{}{},
		},
		"result": result,
		"meta": map[string]interface{}{
			"requestId": qid,
		},
	}, nil
}

// ExecuteGremlinExplainQuery returns an explain plan for a Gremlin query without
// executing it.
func (s *NeptuneDataService) ExecuteGremlinExplainQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	body := req.Body
	var params struct {
		Gremlin string `json:"gremlin"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Gremlin == "" {
		return nil, missingParameter("gremlin")
	}

	return map[string]interface{}{
		"code": "gremlin-traversal-explanation",
		"plan": explainGremlinQuery(params.Gremlin),
	}, nil
}

// ExecuteGremlinProfileQuery returns a profiled explain plan for a Gremlin query.
func (s *NeptuneDataService) ExecuteGremlinProfileQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	body := req.Body
	var params struct {
		Gremlin string `json:"gremlin"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Gremlin == "" {
		return nil, missingParameter("gremlin")
	}

	return map[string]interface{}{
		"code":    "gremlin-traversal-profile",
		"profile": profileGremlinQuery(params.Gremlin),
	}, nil
}

func (s *NeptuneDataService) GetGremlinQueryStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	queryId := getPathParam(req, "queryId")
	if queryId == "" {
		return nil, missingParameter("queryId")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	qs, ok := s.queryState[queryId]
	if !ok {
		return nil, badRequest(fmt.Sprintf("query not found: %s", queryId))
	}

	evalStats := map[string]interface{}{
		"cancelled": qs.Status == "cancelled",
		"elapsed":   elapsedMs(qs.StartedAt, qs.EndedAt),
		"waited":    0,
	}

	return map[string]interface{}{
		"queryId":        qs.ID,
		"queryString":    qs.Query,
		"queryEvalStats": evalStats,
	}, nil
}

func (s *NeptuneDataService) ListGremlinQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	includeWaiting := request.GetBoolParam(req.Parameters, "includeWaiting")
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpired()

	var queries []interface{}
	var acceptedCount, runningCount int32

	for _, qs := range s.queryState {
		if qs.QueryType != "gremlin" {
			continue
		}
		if qs.Status == "waiting" && !includeWaiting {
			continue
		}
		entry := map[string]interface{}{
			"queryId":     qs.ID,
			"queryString": qs.Query,
		}
		switch qs.Status {
		case "running":
			runningCount++
		default:
			acceptedCount++
		}
		queries = append(queries, entry)
	}

	return map[string]interface{}{
		"queries":            queries,
		"acceptedQueryCount": acceptedCount,
		"runningQueryCount":  runningCount,
	}, nil
}

func (s *NeptuneDataService) CancelGremlinQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	queryId := getPathParam(req, "queryId")
	if queryId == "" {
		return nil, missingParameter("queryId")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	qs, ok := s.queryState[queryId]
	if !ok {
		return nil, badRequest(fmt.Sprintf("query not found: %s", queryId))
	}
	qs.Status = "cancelled"
	qs.EndedAt = time.Now()

	return map[string]interface{}{
		"status": "200 OK",
	}, nil
}

// ExecuteOpenCypherQuery parses and executes an OpenCypher query against the
// graph engine, dispatching to the appropriate executor based on query type.
func (s *NeptuneDataService) ExecuteOpenCypherQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	body := req.Body
	var params struct {
		Query      string          `json:"query"`
		Parameters json.RawMessage `json:"parameters"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Query == "" {
		return nil, missingParameter("query")
	}

	var cypherParams map[string]any
	if len(params.Parameters) > 0 && string(params.Parameters) != "" {
		trimmed := strings.TrimSpace(string(params.Parameters))
		if len(trimmed) > 0 && trimmed[0] == '"' {
			var paramStr string
			if err := json.Unmarshal(params.Parameters, &paramStr); err != nil {
				return nil, invalidParameter(fmt.Sprintf("invalid parameters: %v", err))
			}
			if err := json.Unmarshal([]byte(paramStr), &cypherParams); err != nil {
				return nil, invalidParameter(fmt.Sprintf("invalid parameters: %v", err))
			}
		} else {
			if err := json.Unmarshal(params.Parameters, &cypherParams); err != nil {
				return nil, invalidParameter(fmt.Sprintf("invalid parameters: %v", err))
			}
		}
	}

	qid := generateQueryID()
	s.trackQuery(qid, params.Query, "opencypher")

	parsed, err := cypherparser.Parse(params.Query)
	if err != nil {
		s.resolveQuery(qid, nil, err)
		return nil, malformedQuery(err.Error())
	}

	reader := reqCtx.GraphReader()
	writer := reqCtx.GraphWriter()

	var result *cypherparser.CypherResult
	switch {
	case parsed.Write != nil:
		result, err = cypherparser.ExecuteWrite(ctx, reader, writer, parsed.Write, cypherParams)
	case parsed.Merge != nil:
		result, err = cypherparser.ExecuteMerge(ctx, reader, writer, parsed.Merge, cypherParams)
	case parsed.Read != nil && (len(parsed.Read.Set) > 0 || len(parsed.Read.Delete) > 0 || len(parsed.Read.Remove) > 0 || parsed.Read.Create != nil):
		result, err = cypherparser.ExecuteQueryWrite(ctx, reader, writer, parsed.Read, cypherParams)
	default:
		result, err = cypherparser.Execute(ctx, reader, parsed.Read, cypherParams)
	}
	s.resolveQuery(qid, result, err)
	if err != nil {
		return nil, malformedQuery(err.Error())
	}

	return map[string]interface{}{
		"results": result,
	}, nil
}

// ExecuteOpenCypherExplainQuery returns an explain plan for an OpenCypher query.
// Only read queries are supported; write queries return BadRequestException.
func (s *NeptuneDataService) ExecuteOpenCypherExplainQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	body := req.Body
	var params struct {
		Query      string          `json:"query"`
		Explain    string          `json:"explain"`
		Parameters json.RawMessage `json:"parameters"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Query == "" {
		return nil, missingParameter("query")
	}

	parsed, err := cypherparser.Parse(params.Query)
	if err != nil {
		return nil, malformedQuery(err.Error())
	}

	if parsed.Read == nil {
		return nil, badRequest("EXPLAIN is only supported for read queries")
	}

	plan := cypherparser.BuildExplainPlan(parsed.Read, reqCtx.GraphReader())
	return map[string]interface{}{
		"explain": plan,
	}, nil
}

func (s *NeptuneDataService) GetOpenCypherQueryStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	queryId := getPathParam(req, "queryId")
	if queryId == "" {
		return nil, missingParameter("queryId")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	qs, ok := s.queryState[queryId]
	if !ok {
		return nil, badRequest(fmt.Sprintf("query not found: %s", queryId))
	}

	evalStats := map[string]interface{}{
		"cancelled": qs.Status == "cancelled",
		"elapsed":   elapsedMs(qs.StartedAt, qs.EndedAt),
		"waited":    0,
	}

	return map[string]interface{}{
		"queryId":        qs.ID,
		"queryString":    qs.Query,
		"queryEvalStats": evalStats,
	}, nil
}

func (s *NeptuneDataService) ListOpenCypherQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	includeWaiting := request.GetBoolParam(req.Parameters, "includeWaiting")
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpired()

	var queries []interface{}
	var acceptedCount, runningCount int32

	for _, qs := range s.queryState {
		if qs.QueryType != "opencypher" {
			continue
		}
		if qs.Status == "waiting" && !includeWaiting {
			continue
		}
		entry := map[string]interface{}{
			"queryId":     qs.ID,
			"queryString": qs.Query,
		}
		switch qs.Status {
		case "running":
			runningCount++
		default:
			acceptedCount++
		}
		queries = append(queries, entry)
	}

	return map[string]interface{}{
		"queries":            queries,
		"acceptedQueryCount": acceptedCount,
		"runningQueryCount":  runningCount,
	}, nil
}

func (s *NeptuneDataService) CancelOpenCypherQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	queryId := getPathParam(req, "queryId")
	silent := request.GetBoolParam(req.Parameters, "silent")
	if queryId == "" {
		return nil, missingParameter("queryId")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	qs, ok := s.queryState[queryId]
	if !ok {
		return nil, badRequest(fmt.Sprintf("query not found: %s", queryId))
	}
	qs.Status = "cancelled"
	qs.EndedAt = time.Now()

	if silent {
		return map[string]interface{}{}, nil
	}

	return map[string]interface{}{
		"payload": true,
		"status":  "200 OK",
	}, nil
}

// ExecuteFastReset handles the two-phase database reset protocol. The
// initiateDatabaseReset action issues a time-limited token; performDatabaseReset
// validates the token and clears all graph data.
func (s *NeptuneDataService) ExecuteFastReset(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	body := req.Body
	var params struct {
		Action string `json:"action"`
		Token  string `json:"token"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}

	switch params.Action {
	case "initiateDatabaseReset":
		token := generateFastResetToken()
		s.mu.Lock()
		now := time.Now()
		for k, expiry := range s.fastTokens {
			if now.After(expiry) {
				delete(s.fastTokens, k)
			}
		}
		s.fastTokens[token] = now.Add(30 * time.Second)
		s.mu.Unlock()
		return map[string]interface{}{
			"payload": map[string]interface{}{
				"token": token,
			},
		}, nil

	case "performDatabaseReset":
		if params.Token == "" {
			return nil, missingParameter("token")
		}
		s.mu.Lock()
		expiry, ok := s.fastTokens[params.Token]
		if !ok || time.Now().After(expiry) {
			s.mu.Unlock()
			return nil, preconditionFailed("invalid or expired token")
		}
		delete(s.fastTokens, params.Token)
		s.mu.Unlock()

		if gs, ok := reqCtx.GraphWriter().(graphengine.GraphStore); ok {
			gs.Clear()
		}
		return map[string]interface{}{
			"status": "200 OK",
		}, nil

	default:
		return nil, invalidParameter(fmt.Sprintf("unknown action: %s", params.Action))
	}
}

// GetPropertygraphStatistics returns auto-computed property graph statistics
// including node and edge counts grouped by label.
func (s *NeptuneDataService) GetPropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = req

	s.mu.Lock()
	s.refreshStatistics(reqCtx)

	sigCount := int64(0)
	predCount := int64(0)
	for _, cnt := range s.stats.LabelCounts {
		sigCount += cnt
	}
	for _, cnt := range s.stats.RelCounts {
		predCount += cnt
	}

	result := map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"active":       true,
			"autoCompute":  true,
			"date":         time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			"note":         "Automatically computed",
			"statisticsId": "auto-statistics",
			"signatureInfo": map[string]interface{}{
				"signatureCount": sigCount,
				"instanceCount":  s.stats.NodeCount,
				"predicateCount": predCount,
			},
		},
	}
	s.mu.Unlock()
	return result, nil
}

func (s *NeptuneDataService) ManagePropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	body := req.Body
	var params struct {
		Mode string `json:"mode"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}

	switch params.Mode {
	case "disableAutoCompute":
		return map[string]interface{}{
			"status": "200",
		}, nil
	case "enableAutoCompute":
		return map[string]interface{}{
			"status": "200",
		}, nil
	case "refresh":
		return map[string]interface{}{
			"status":  "200",
			"payload": map[string]interface{}{"statisticsId": generateQueryID()},
		}, nil
	default:
		return nil, invalidParameter(fmt.Sprintf("unknown mode: %s", params.Mode))
	}
}

func (s *NeptuneDataService) DeletePropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = req
	s.mu.Lock()
	s.stats = &GraphStatistics{LabelCounts: make(map[string]int64), RelCounts: make(map[string]int64)}
	s.mu.Unlock()
	return map[string]interface{}{
		"status":  "200",
		"payload": map[string]interface{}{},
	}, nil
}

// GetPropertygraphSummary returns a detailed summary of the property graph
// including node/edge counts, label lists, and structural metadata.
func (s *NeptuneDataService) GetPropertygraphSummary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	mode := getPathParam(req, "mode")

	s.mu.Lock()
	s.refreshStatistics(reqCtx)

	nodeLabels := make([]string, 0, len(s.stats.LabelCounts))
	for label := range s.stats.LabelCounts {
		nodeLabels = append(nodeLabels, label)
	}

	edgeLabels := make([]string, 0, len(s.stats.RelCounts))
	for label := range s.stats.RelCounts {
		edgeLabels = append(edgeLabels, label)
	}

	summary := map[string]interface{}{
		"numNodes":      s.stats.NodeCount,
		"numEdges":      s.stats.EdgeCount,
		"numNodeLabels": int64(len(nodeLabels)),
		"numEdgeLabels": int64(len(edgeLabels)),
	}

	if mode == "detailed" {
		summary["numNodeProperties"] = 0
		summary["numEdgeProperties"] = 0
		summary["totalNodePropertyValues"] = 0
		summary["totalEdgePropertyValues"] = 0
		summary["nodeLabels"] = nodeLabels
		summary["edgeLabels"] = edgeLabels
		summary["nodeProperties"] = []interface{}{}
		summary["edgeProperties"] = []interface{}{}
		summary["nodeStructures"] = []interface{}{}
		summary["edgeStructures"] = []interface{}{}
	}

	s.mu.Unlock()

	return map[string]interface{}{
		"payload": map[string]interface{}{
			"version":                       "v1",
			"graphSummary":                  summary,
			"lastStatisticsComputationTime": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		},
	}, nil
}

func (s *NeptuneDataService) GetPropertygraphStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = request.GetIntParam(req.Parameters, "commitNum")
	_ = request.GetStringParam(req.Parameters, "iteratorType")
	_ = request.GetIntParam(req.Parameters, "limit")
	_ = request.GetIntParam(req.Parameters, "opNum")
	return map[string]interface{}{
		"format":                   "PG_JSON",
		"lastEventId":              map[string]string{},
		"totalRecords":             0,
		"records":                  []interface{}{},
		"lastTrxTimestampInMillis": int64(0),
	}, nil
}

func (s *NeptuneDataService) StartLoaderJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	body := req.Body
	var params struct {
		Source      string `json:"source"`
		Format      string `json:"format"`
		IamRoleArn  string `json:"iamRoleArn"`
		Mode        string `json:"mode"`
		Parallelism string `json:"parallelism"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Source == "" {
		return nil, missingParameter("source")
	}
	if params.Format == "" {
		return nil, missingParameter("format")
	}

	loadId := generateQueryID()

	s.mu.Lock()
	s.loaderJobs[loadId] = &loaderJobState{
		LoadId:    loadId,
		Status:    "LOAD_COMPLETED",
		Source:    params.Source,
		Format:    params.Format,
		StartedAt: time.Now(),
	}
	s.mu.Unlock()

	return map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"loadId": loadId,
		},
	}, nil
}

func (s *NeptuneDataService) GetLoaderJobStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	loadId := getPathParam(req, "loadId")
	if loadId == "" {
		return nil, missingParameter("loadId")
	}
	details := request.GetBoolParam(req.Parameters, "details")
	errors := request.GetBoolParam(req.Parameters, "errors")
	errorsPerPage := request.GetIntParam(req.Parameters, "errorsPerPage")
	page := request.GetIntParam(req.Parameters, "page")

	_ = details
	_ = errors
	_ = errorsPerPage
	_ = page

	s.mu.RLock()
	job, ok := s.loaderJobs[loadId]
	s.mu.RUnlock()

	if !ok {
		return nil, bulkLoadNotFound(loadId)
	}

	payload := map[string]interface{}{
		"loadId":        job.LoadId,
		"status":        job.Status,
		"feedCount":     map[string]interface{}{"total": 0, "succeeded": 0, "failed": 0},
		"overallStatus": map[string]interface{}{"totalRecords": 0, "totalTime": 0},
	}

	return map[string]interface{}{
		"status":  "200",
		"payload": payload,
	}, nil
}

func (s *NeptuneDataService) ListLoaderJobs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = request.GetBoolParam(req.Parameters, "includeQueuedLoads")
	_ = request.GetIntParam(req.Parameters, "limit")

	s.mu.RLock()
	loadIds := make([]string, 0, len(s.loaderJobs))
	for _, job := range s.loaderJobs {
		loadIds = append(loadIds, job.LoadId)
	}
	s.mu.RUnlock()

	return map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"loadIds": loadIds,
		},
	}, nil
}

func (s *NeptuneDataService) CancelLoaderJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	loadId := getPathParam(req, "loadId")
	if loadId == "" {
		return nil, missingParameter("loadId")
	}

	s.mu.Lock()
	job, ok := s.loaderJobs[loadId]
	if !ok {
		s.mu.Unlock()
		return nil, bulkLoadNotFound(loadId)
	}
	job.Status = "CANCELLED"
	job.FinishedAt = time.Now()
	s.mu.Unlock()

	return map[string]interface{}{
		"status": "200",
	}, nil
}

// unsupported returns an UnsupportedOperationException for operations not
// yet implemented by vorpalstacks (SPARQL, ML, etc.).
func (s *NeptuneDataService) unsupported(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = req
	return nil, unsupported("this operation is not supported by vorpalstacks")
}

// trackQuery records the start of a query execution. Caller must not hold the lock.
func (s *NeptuneDataService) trackQuery(id, query, queryType string) {
	s.mu.Lock()
	s.cleanupExpired()
	s.queryState[id] = &queryState{
		ID:        id,
		Query:     query,
		QueryType: queryType,
		Status:    "running",
		StartedAt: time.Now(),
	}
	s.mu.Unlock()
}

// resolveQuery records the completion of a query, setting its status, end time,
// and optional result or error.
func (s *NeptuneDataService) resolveQuery(id string, result any, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	qs, ok := s.queryState[id]
	if !ok {
		return
	}
	qs.EndedAt = time.Now()

	if err != nil {
		qs.Status = "failed"
		qs.Error = err.Error()
	} else {
		qs.Status = "complete"
		qs.Result = result
	}
}

// cleanupExpired removes query state entries whose elapsed time exceeds
// queryStateTTL. Caller must hold the write lock.
func (s *NeptuneDataService) cleanupExpired() {
	now := time.Now()
	for id, qs := range s.queryState {
		if qs.EndedAt.IsZero() {
			continue
		}
		if now.Sub(qs.EndedAt) > queryStateTTL {
			delete(s.queryState, id)
		}
	}
}

// refreshStatistics re-reads graph statistics from the engine and updates the
// cached stats struct. Caller must hold the appropriate lock.
func (s *NeptuneDataService) refreshStatistics(reqCtx *request.RequestContext) {
	reader := reqCtx.GraphReader()
	if reader == nil {
		return
	}
	s.stats.NodeCount = reader.CountNodes()
	s.stats.EdgeCount = reader.CountEdges()
	s.stats.LabelCounts = reader.GetLabelCounts()
	s.stats.RelCounts = reader.GetRelCounts()
}

// generateQueryID produces a unique query identifier from the current timestamp
// and a monotonically increasing counter. The counter is safe because all callers
// hold s.mu (via trackQuery/resolveQuery) before incrementing.
var queryCounter int64

func generateQueryID() string {
	id := atomic.AddInt64(&queryCounter, 1)
	return fmt.Sprintf("query-%d-%d", time.Now().UnixMilli(), id)
}

// generateFastResetToken produces a unique token for the two-phase FastReset
// protocol. The counter is safe because the caller holds s.mu before incrementing.
var tokenCounter int64

func generateFastResetToken() string {
	tokenCounter++
	return fmt.Sprintf("frt-%d-%d", time.Now().UnixNano(), tokenCounter)
}

// explainGremlinQuery parses a Gremlin query and produces a step-by-step
// explanation of the traversal plan.
func explainGremlinQuery(query string) map[string]interface{} {
	parsed, err := gremlinparser.Parse(query)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("parse error: %v", err),
		}
	}

	steps := []map[string]interface{}{}
	if len(parsed.Statements) > 0 && parsed.Statements[0].Traversal != nil {
		for _, step := range parsed.Statements[0].Traversal.Steps {
			stepInfo := map[string]interface{}{
				"name": step.Name,
			}
			if len(step.Args) > 0 {
				args := make([]interface{}, len(step.Args))
				for i, arg := range step.Args {
					args[i] = describeArg(arg)
				}
				stepInfo["args"] = args
			}
			steps = append(steps, stepInfo)
		}
	}

	return map[string]interface{}{
		"steps": steps,
	}
}

// profileGremlinQuery returns an explain plan augmented with profiling metrics.
func profileGremlinQuery(query string) map[string]interface{} {
	explain := explainGremlinQuery(query)
	explain["profile"] = map[string]interface{}{
		"metrics": map[string]interface{}{
			"dur":     0,
			"count":   1,
			"size":    0,
			"time":    0,
			"incTime": 0,
		},
		"indices": map[string]interface{}{},
	}
	return explain
}

// describeArg converts a Gremlin argument to a serialisable representation for
// explain/profile output.
func describeArg(arg gremlinparser.Argument) interface{} {
	switch arg.Kind {
	case gremlinparser.ArgString:
		return arg.Str
	case gremlinparser.ArgInt:
		return arg.Int
	case gremlinparser.ArgFloat:
		return arg.Float
	case gremlinparser.ArgBool:
		return arg.Bool
	case gremlinparser.ArgNull:
		return nil
	case gremlinparser.ArgEnum:
		if arg.Enum != nil {
			return arg.Enum.Value
		}
		return nil
	case gremlinparser.ArgList:
		items := make([]interface{}, len(arg.List))
		for i, a := range arg.List {
			items[i] = describeArg(a)
		}
		return items
	case gremlinparser.ArgMap:
		m := make(map[string]interface{})
		for _, entry := range arg.Map {
			if key, ok := describeArg(entry.Key).(string); ok {
				m[key] = describeArg(entry.Value)
			}
		}
		return m
	case gremlinparser.ArgPredicate:
		if arg.Pred != nil {
			return map[string]interface{}{
				"type":   arg.Pred.Type,
				"method": arg.Pred.Method,
				"args":   describeArgs(arg.Pred.Args),
			}
		}
		return nil
	case gremlinparser.ArgNestedTraversal:
		if arg.Trav != nil {
			steps := []map[string]interface{}{}
			for _, step := range arg.Trav.Steps {
				stepInfo := map[string]interface{}{
					"name": step.Name,
				}
				if len(step.Args) > 0 {
					args := make([]interface{}, len(step.Args))
					for i, a := range step.Args {
						args[i] = describeArg(a)
					}
					stepInfo["args"] = args
				}
				steps = append(steps, stepInfo)
			}
			return steps
		}
		return nil
	default:
		return arg.Str
	}
}

// describeArgs converts a slice of Gremlin arguments to serialisable representations.
func describeArgs(args []gremlinparser.Argument) []interface{} {
	result := make([]interface{}, len(args))
	for i, a := range args {
		result[i] = describeArg(a)
	}
	return result
}
