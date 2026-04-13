package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
	"vorpalstacks/pkg/graphengine"
)

const (
	queryStateTTL      = 5 * time.Minute
	loaderJobTTL       = 1 * time.Hour
	statsCacheTTL      = 10 * time.Minute
	statsLastAccessTTL = 30 * time.Minute
)

// NeptuneDataService implements the Neptune Data API service, handling Gremlin,
// OpenCypher, graph statistics, loader, and ML operations. Query state and
// loader jobs are persisted in Pebble storage via per-region Neptune stores.
// Graph data is read and written via the GraphReader/GraphWriter interfaces
// provided by the request context.
type NeptuneDataService struct {
	mu                 sync.RWMutex
	fastTokens         sync.Map
	statsMap           sync.Map
	statsDisabled      bool
	autoComputeEnabled bool
	startTime          time.Time
	storageManager     *storage.RegionStorageManager
	stores             sync.Map
	cancelCleanup      context.CancelFunc
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
	LastRefresh time.Time        `json:"-"`
	LastAccess  time.Time        `json:"-"`
}

// NewNeptuneDataService creates a new service instance. Per-region stores are
// created lazily via the RegionStorageManager. A background goroutine is
// started to periodically purge expired query states from Pebble storage.
func NewNeptuneDataService() *NeptuneDataService {
	ctx, cancel := context.WithCancel(context.Background())
	s := &NeptuneDataService{
		autoComputeEnabled: true,
		startTime:          time.Now(),
		cancelCleanup:      cancel,
	}
	go s.cleanupExpiredQueries(ctx)
	return s
}

// Close stops the background query-state cleanup goroutine.
func (s *NeptuneDataService) Close() {
	if s.cancelCleanup != nil {
		s.cancelCleanup()
	}
}

// cleanupExpiredQueries periodically scans all per-region stores and deletes
// terminal query states whose EndTime has exceeded queryStateTTL.
func (s *NeptuneDataService) cleanupExpiredQueries(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.purgeExpiredQueries()
		}
	}
}

// purgeExpiredQueries iterates every region's store and removes expired
// terminal queries ("complete", "failed", "cancelled"). Also purges expired
// loader jobs and evicts stale statistics cache entries.
func (s *NeptuneDataService) purgeExpiredQueries() {
	terminalStates := map[string]bool{
		"complete":  true,
		"failed":    true,
		"cancelled": true,
	}
	terminalLoaderStates := map[string]bool{
		"LOAD_COMPLETED": true,
		"LOAD_FAILED":    true,
		"CANCELLED":      true,
	}

	s.stores.Range(func(_, value any) bool {
		store, ok := value.(*neptunestore.NeptuneStore)
		if !ok {
			return true
		}

		queries, err := store.ListQueries()
		if err != nil {
			logs.Warn("failed to list queries for cleanup", logs.Err(err))
		} else {
			for _, q := range queries {
				if !terminalStates[q.Status] {
					continue
				}
				if q.EndTime == nil {
					continue
				}
				endTime := q.EndTime.AsTime()
				if time.Since(endTime) > queryStateTTL {
					if delErr := store.DeleteQuery(q.QueryId); delErr != nil {
						logs.Warn("failed to delete expired query", logs.String("queryId", q.QueryId), logs.Err(delErr))
					}
				}
			}
		}

		jobs, err := store.ListLoaderJobs()
		if err != nil {
			logs.Warn("failed to list loader jobs for cleanup", logs.Err(err))
		} else {
			for _, j := range jobs {
				if !terminalLoaderStates[j.GetStatus()] {
					continue
				}
				if j.EndTime == nil {
					continue
				}
				endTime := j.EndTime.AsTime()
				if time.Since(endTime) > loaderJobTTL {
					if delErr := store.DeleteLoaderJob(j.GetLoadId()); delErr != nil {
						logs.Warn("failed to delete expired loader job", logs.String("loadId", j.GetLoadId()), logs.Err(delErr))
					}
				}
			}
		}

		return true
	})

	now := time.Now()
	s.mu.RLock()
	s.statsMap.Range(func(key, value any) bool {
		st := value.(*GraphStatistics)
		if now.Sub(st.LastAccess) > statsLastAccessTTL {
			s.statsMap.Delete(key)
		}
		return true
	})
	s.mu.RUnlock()
}

// SetStorageManager injects the region storage manager for per-region store
// caching and admin console access.
func (s *NeptuneDataService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// GetStoreForRegion returns the cached Neptune store for the given region,
// creating one if not already cached.
func (s *NeptuneDataService) GetStoreForRegion(region string) (*neptunestore.NeptuneStore, error) {
	val, err := storecommon.GetOrCreateStoreE(&s.stores, region, func() (*neptunestore.NeptuneStore, error) {
		if s.storageManager == nil {
			return nil, fmt.Errorf("storage manager not set")
		}
		rs, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}
		return neptunestore.NewNeptuneStore(rs), nil
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

// store resolves the NeptuneStore for the current request context, using the
// per-region cache.
func (s *NeptuneDataService) store(reqCtx *request.RequestContext) (*neptunestore.NeptuneStore, error) {
	region := reqCtx.GetRegion()
	return s.GetStoreForRegion(region)
}

// RegisterHandlers registers all Neptune Data API operation handlers with the
// dispatcher. Unsupported operations (SPARQL, ML) return HTTP 501.
func (s *NeptuneDataService) RegisterHandlers(d handler.Registrar) {
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
		now := time.Now()
		s.fastTokens.Store(token, now.Add(30*time.Second))
		s.fastTokens.Range(func(key, value any) bool {
			if now.After(value.(time.Time)) {
				s.fastTokens.Delete(key)
			}
			return true
		})
		return map[string]interface{}{
			"payload": map[string]interface{}{
				"token": token,
			},
		}, nil

	case "performDatabaseReset":
		if params.Token == "" {
			return nil, missingParameter("token")
		}
		val, ok := s.fastTokens.Load(params.Token)
		if !ok || time.Now().After(val.(time.Time)) {
			return nil, preconditionFailed("invalid or expired token")
		}
		s.fastTokens.Delete(params.Token)

		if gs, ok := reqCtx.GraphWriter().(graphengine.GraphStore); ok {
			gs.Clear()
		}
		region := reqCtx.GetRegion()
		s.mu.Lock()
		s.statsDisabled = false
		s.autoComputeEnabled = true
		s.mu.Unlock()
		s.statsMap.Store(region, &GraphStatistics{
			LabelCounts: make(map[string]int64),
			RelCounts:   make(map[string]int64),
		})
		if resetStore, err := s.GetStoreForRegion(region); err == nil {
			queries, _ := resetStore.ListQueries()
			for _, q := range queries {
				resetStore.DeleteQuery(q.GetQueryId())
			}
			jobs, _ := resetStore.ListLoaderJobs()
			for _, j := range jobs {
				resetStore.DeleteLoaderJob(j.GetLoadId())
			}
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

func (s *NeptuneDataService) getStats(region string) *GraphStatistics {
	val, _ := s.statsMap.LoadOrStore(region, &GraphStatistics{
		LabelCounts: make(map[string]int64),
		RelCounts:   make(map[string]int64),
		LastAccess:  time.Now(),
		LastRefresh: time.Now(),
	})
	st := val.(*GraphStatistics)
	st.LastAccess = time.Now()
	return st
}

func (s *NeptuneDataService) refreshStatistics(reqCtx *request.RequestContext) {
	if reqCtx == nil {
		return
	}
	reader := reqCtx.GraphReader()
	if reader == nil {
		return
	}
	region := reqCtx.GetRegion()
	stats := s.getStats(region)
	stats.NodeCount = reader.CountNodes()
	stats.EdgeCount = reader.CountEdges()
	stats.LabelCounts, _ = reader.GetLabelCounts()
	stats.RelCounts, _ = reader.GetRelCounts()
	stats.LastRefresh = time.Now()
}

var queryCounter int64

func generateQueryID() string {
	id := atomic.AddInt64(&queryCounter, 1)
	return fmt.Sprintf("query-%d-%d", time.Now().UnixMilli(), id)
}

var tokenCounter int64

func generateFastResetToken() string {
	id := atomic.AddInt64(&tokenCounter, 1)
	return fmt.Sprintf("frt-%d-%d", time.Now().UnixNano(), id)
}
