package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
	"vorpalstacks/pkg/graphengine"
)

// queryStateTTL defines how long completed query state entries are retained
// before being eligible for cleanup.
const (
	queryStateTTL = 5 * time.Minute
)

// NeptuneDataService implements the Neptune Data API service, handling Gremlin,
// OpenCypher, graph statistics, loader, and ML operations. Query state and
// loader jobs are persisted in Pebble storage via the shared Neptune store.
// Graph data is read and written via the GraphReader/GraphWriter interfaces
// provided by the request context.
type NeptuneDataService struct {
	store      *neptunestore.NeptuneStore
	mu         sync.RWMutex
	fastTokens map[string]time.Time
	stats      *GraphStatistics
	startTime  time.Time
}

// elapsedMs returns the elapsed time between two Unix-millisecond timestamps.
func elapsedMs(startMs, endMs int64) int64 {
	return endMs - startMs
}

// GraphStatistics holds cached graph-level statistics for the property graph.
// Refreshed on demand when statistics or summary endpoints are called.
type GraphStatistics struct {
	NodeCount   int64            `json:"numNodes"`
	EdgeCount   int64            `json:"numEdges"`
	LabelCounts map[string]int64 `json:"-"`
	RelCounts   map[string]int64 `json:"-"`
}

// NewNeptuneDataService creates a new service instance backed by the shared
// Neptune Pebble store for query states and loader jobs.
func NewNeptuneDataService(store *neptunestore.NeptuneStore) *NeptuneDataService {
	return &NeptuneDataService{
		store:      store,
		fastTokens: make(map[string]time.Time),
		stats: &GraphStatistics{
			LabelCounts: make(map[string]int64),
			RelCounts:   make(map[string]int64),
		},
		startTime: time.Now(),
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

// trackQuery records the start of a query execution in Pebble storage.
func (s *NeptuneDataService) trackQuery(id, query, queryType string) {
	qr := &pb.QueryState{
		QueryId:   id,
		QueryType: queryType,
		Status:    "running",
		StartTime: timestamppb.Now(),
	}
	_ = s.store.CreateQuery(qr)
}

// resolveQuery records the completion of a query in Pebble storage.
func (s *NeptuneDataService) resolveQuery(id string, result any, err error) {
	qr, storeErr := s.store.GetQuery(id)
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
	_ = s.store.UpdateQuery(qr)
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
// protocol. Uses atomic operations for thread safety.
var tokenCounter int64

func generateFastResetToken() string {
	id := atomic.AddInt64(&tokenCounter, 1)
	return fmt.Sprintf("frt-%d-%d", time.Now().UnixNano(), id)
}
