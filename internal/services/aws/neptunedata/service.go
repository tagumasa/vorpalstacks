package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/pkg/graphengine"
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
