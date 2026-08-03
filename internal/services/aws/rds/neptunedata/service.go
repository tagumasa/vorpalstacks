package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"

	"google.golang.org/protobuf/types/known/timestamppb"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/graphengine"
	"vorpalstacks/internal/eventbus"
	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	"vorpalstacks/internal/server/listener"
	"vorpalstacks/internal/server/portalloc"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	storecommon "vorpalstacks/internal/store/aws/common"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	"vorpalstacks/internal/utils/timeutils"
)

var (
	_ rdssvc.Engine    = (*NeptuneDataService)(nil)
	_ rdssvc.GetPorter = (*NeptuneDataService)(nil)
)

const (
	queryStateTTL      = 5 * time.Minute
	loaderJobTTL       = 1 * time.Hour
	statsCacheTTL      = 10 * time.Minute
	statsLastAccessTTL = 30 * time.Minute
	maxConcurrentLoads = 64
)

// NeptuneDataService implements the Neptune Data API service, handling Gremlin,
// OpenCypher, graph statistics, loader, and ML operations. Each Neptune DB
// cluster gets its own isolated graph engine backed by a dedicated Pebble
// bucket. Query state and loader jobs are persisted in Pebble storage via
// per-region Neptune stores.
type NeptuneDataService struct {
	mu                 sync.RWMutex
	fastTokens         sync.Map
	statsMap           sync.Map
	statsDisabled      bool
	autoComputeEnabled bool
	startTime          time.Time
	storageManager     *storage.RegionStorageManager
	region             string
	stores             sync.Map
	loaderWg           sync.WaitGroup
	cancelCleanup      context.CancelFunc
	loaderCancelChs    map[string]chan struct{}
	loaderMu           sync.Mutex
	s3Invoker          eventbus.S3Invoker
	portAllocator      *portalloc.Allocator
	listenerManager    *listener.Manager
	dispatcherHandler  func() http.Handler
	activeEngines      map[string]*clusterEngineEntry
	enginesMu          sync.RWMutex
	graphCache         *graphengine.Cache
	dispatcherRunning  atomic.Bool
}

// clusterEngineEntry holds the graph engine instance for a single DB cluster.
type clusterEngineEntry struct {
	db     *graphengine.DB
	region string
}

// GraphStatistics holds cached graph-level statistics for the property graph.
// Refreshed on demand when statistics or summary endpoints are called.
type GraphStatistics struct {
	mu          sync.Mutex
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
func NewNeptuneDataService(allocator *portalloc.Allocator) *NeptuneDataService {
	ctx, cancel := context.WithCancel(context.Background())
	s := &NeptuneDataService{
		autoComputeEnabled: true,
		startTime:          time.Now(),
		cancelCleanup:      cancel,
		portAllocator:      allocator,
		activeEngines:      make(map[string]*clusterEngineEntry),
		loaderCancelChs:    make(map[string]chan struct{}),
	}
	go s.cleanupExpiredQueries(ctx)
	return s
}

// Close stops the background query-state cleanup goroutine and closes all
// per-cluster graph engines.
func (s *NeptuneDataService) Shutdown() {
	if s.cancelCleanup != nil {
		s.cancelCleanup()
	}
	s.loaderWg.Wait()
	s.enginesMu.Lock()
	engines := make(map[string]*clusterEngineEntry, len(s.activeEngines))
	for id, entry := range s.activeEngines {
		engines[id] = entry
	}
	s.activeEngines = make(map[string]*clusterEngineEntry)
	s.enginesMu.Unlock()

	for _, entry := range engines {
		entry.db.Close()
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
			func() {
				defer func() {
					if re := recover(); re != nil {
						logs.Error("neptunedata cleanup panic recovered", logs.Any("panic", re))
					}
				}()
				s.purgeExpiredQueries()
				s.purgeExpiredFastTokens()
			}()
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
		"LOAD_COMPLETED":              true,
		"LOAD_FAILED":                 true,
		"CANCELLED":                   true,
		"LOAD_CANCELLED_DUE_TO_ERROR": true,
		"LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED": true,
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
	s.statsMap.Range(func(key, value any) bool {
		st, ok := value.(*GraphStatistics)
		if !ok {
			return true
		}
		st.mu.Lock()
		expired := now.Sub(st.LastAccess) > statsLastAccessTTL
		st.mu.Unlock()
		if expired {
			s.statsMap.Delete(key)
		}
		return true
	})
}

func (s *NeptuneDataService) purgeExpiredFastTokens() {
	now := time.Now()
	s.fastTokens.Range(func(key, value any) bool {
		t, ok := value.(time.Time)
		if !ok {
			return true
		}
		if now.After(t) {
			s.fastTokens.Delete(key)
		}
		return true
	})
}

// SetStorageManager injects the region storage manager for per-region store
// caching and admin console access.
func (s *NeptuneDataService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
	if regions := sm.GetActiveRegions(); len(regions) > 0 {
		s.region = regions[0]
	}
}

// SetGraphCache injects a shared graph node/edge cache used by all per-cluster
// engines. Must be called before Open or RestoreEngines.
func (s *NeptuneDataService) SetGraphCache(cache *graphengine.Cache) {
	s.graphCache = cache
}

// SetListenerManager injects the listener manager for dynamic per-cluster
// listener registration. Must be called before Open or RestoreEngines.
func (s *NeptuneDataService) SetListenerManager(lm *listener.Manager) {
	s.listenerManager = lm
}

// SetDispatcherHandler sets a lazy resolver for the main HTTP dispatcher.
// The handler is resolved at request time so listeners can be registered
// before the main server has started.
func (s *NeptuneDataService) SetDispatcherHandler(fn func() http.Handler) {
	s.dispatcherHandler = fn
}

// SetS3Invoker injects the S3 reader for bulk loader jobs that load data
// from S3 sources.
func (s *NeptuneDataService) SetS3Invoker(invoker eventbus.S3Invoker) {
	s.s3Invoker = invoker
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

// clusterBucket returns a dedicated Pebble bucket for the given cluster in
// the specified region, following the same bucketBackend pattern used by
// NeptuneGraph.
func (s *NeptuneDataService) clusterBucket(region, clusterID string) (storage.BatchBucket, error) {
	if s.storageManager == nil {
		return nil, fmt.Errorf("storage manager not set")
	}
	rs, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	bkt := rs.Bucket("neptunedata:cluster:" + clusterID)
	bb, ok := bkt.(storage.BatchBucket)
	if !ok {
		return nil, fmt.Errorf("storage bucket does not support batch operations")
	}
	return bb, nil
}

// Open creates and opens a new isolated graph engine for the
// given cluster in the specified region. The engine is stored in the
// activeEngines map. Returns the dynamically allocated port number.
func (s *NeptuneDataService) Open(region, clusterID string) (int, error) {
	bucket, err := s.clusterBucket(region, clusterID)
	if err != nil {
		return 0, fmt.Errorf("neptunedata: failed to get cluster bucket: %w", err)
	}

	opts := graphengine.DefaultOptions()
	if s.graphCache != nil {
		opts.SharedCache = s.graphCache
	}
	db, err := graphengine.New(bucket, opts)
	if err != nil {
		return 0, fmt.Errorf("neptunedata: failed to open cluster engine: %w", err)
	}

	port, err := s.portAllocator.Get("neptunedata", clusterID)
	if err != nil {
		db.Close()
		return 0, fmt.Errorf("neptunedata: failed to allocate port: %w", err)
	}

	s.enginesMu.Lock()
	entry := &clusterEngineEntry{db: db, region: region}
	s.activeEngines[clusterID] = entry
	s.enginesMu.Unlock()

	if s.listenerManager != nil && s.dispatcherHandler != nil {
		s.registerClusterListener(clusterID, entry)
	}

	logs.Info("opened cluster engine", logs.String("cluster", clusterID), logs.Int("port", port))
	return port, nil
}

// Close closes the graph engine for the given cluster and
// releases its dynamically allocated port.
func (s *NeptuneDataService) Close(clusterID string) error {
	s.enginesMu.Lock()
	entry, ok := s.activeEngines[clusterID]
	if ok {
		delete(s.activeEngines, clusterID)
	}
	s.enginesMu.Unlock()

	if !ok {
		return nil
	}

	entry.db.Close()

	if s.listenerManager != nil {
		s.listenerManager.Unregister("neptune_cluster_" + clusterID)
	}

	if err := s.portAllocator.Release("neptunedata", clusterID); err != nil {
		logs.Warn("failed to release cluster port", logs.String("cluster", clusterID), logs.Err(err))
	}

	logs.Info("closed cluster engine", logs.String("cluster", clusterID))
	return nil
}

// GetClusterEngine returns the graph engine for the specified cluster, or nil if
// the cluster has no active engine.
func (s *NeptuneDataService) GetClusterEngine(clusterID string) *graphengine.DB {
	s.enginesMu.RLock()
	entry, ok := s.activeEngines[clusterID]
	s.enginesMu.RUnlock()
	if !ok {
		return nil
	}
	return entry.db
}

// GetPort returns the allocated listener port for a cluster.
func (s *NeptuneDataService) GetPort(clusterID string) (int, error) {
	s.enginesMu.RLock()
	_, exists := s.activeEngines[clusterID]
	s.enginesMu.RUnlock()
	if !exists {
		return 0, fmt.Errorf("neptunedata: cluster %s has no active engine", clusterID)
	}
	return s.portAllocator.Get("neptunedata", clusterID)
}

func (s *NeptuneDataService) EngineType() string { return "neptune" }

// RestoreEngines reopens graph engines for all existing clusters after a
// service restart. Reads cluster IDs from the Neptune store.
func (s *NeptuneDataService) RestoreEngines() {
	if s.storageManager == nil {
		return
	}
	store, err := s.GetStoreForRegion(s.region)
	if err != nil {
		logs.Warn("failed to get store for engine restore", logs.Err(err))
		return
	}

	clusters, err := store.ListClusters()
	if err != nil {
		logs.Warn("failed to list clusters for engine restore", logs.Err(err))
		return
	}

	for _, c := range clusters {
		if c.Status != "available" {
			continue
		}
		if _, err := s.Open(s.region, c.DBClusterIdentifier); err != nil {
			logs.Warn("failed to restore cluster engine", logs.String("cluster", c.DBClusterIdentifier), logs.Err(err))
		}
	}
}

// DataPlaneHandler returns an HTTP handler for the specified cluster's data
// plane. The handler detects WebSocket upgrade requests on /gremlin and
// forwards them to the GremlinWSServer for TinkerPop protocol handling. All
// other requests are injected with the cluster's graph engine context and
// forwarded to the main HTTP dispatcher.
func (s *NeptuneDataService) DataPlaneHandler(clusterID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Route WebSocket upgrades on /gremlin to the TinkerPop handler
		if r.URL.Path == "/gremlin" && websocket.IsWebSocketUpgrade(r) {
			wsServer := NewGremlinWSServer(s)
			wsServer.ServeHTTP(w, r, clusterID)
			return
		}

		db := s.GetClusterEngine(clusterID)
		if db == nil {
			http.Error(w, fmt.Sprintf("cluster %s engine not available", clusterID), http.StatusServiceUnavailable)
			return
		}
		if s.dispatcherHandler == nil {
			http.Error(w, "dispatcher not available", http.StatusServiceUnavailable)
			return
		}
		handler := s.dispatcherHandler()
		if handler == nil {
			http.Error(w, "dispatcher not available", http.StatusServiceUnavailable)
			return
		}
		ctx := request.WithGraphDBOverride(r.Context(), db)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RegisterClusterListeners registers HTTP listeners for all active cluster
// engines. Must be called after the main HTTP server has started (so that
// MainHandler is available) and after RestoreEngines has restored existing
// engines. Called from App.Run() or equivalent.
func (s *NeptuneDataService) RegisterClusterListeners() {
	if s.listenerManager == nil {
		return
	}
	s.enginesMu.RLock()
	defer s.enginesMu.RUnlock()
	for clusterID, entry := range s.activeEngines {
		s.registerClusterListener(clusterID, entry)
	}
}

// RegisterClusterListener registers an HTTP listener for a single cluster.
// Called during CreateDBCluster after the server is running.
func (s *NeptuneDataService) RegisterClusterListener(clusterID string) {
	if s.listenerManager == nil {
		return
	}
	s.enginesMu.RLock()
	entry, ok := s.activeEngines[clusterID]
	s.enginesMu.RUnlock()
	if !ok {
		return
	}
	s.registerClusterListener(clusterID, entry)
}

func (s *NeptuneDataService) registerClusterListener(clusterID string, entry *clusterEngineEntry) {
	listenerName := "neptune_cluster_" + clusterID
	if s.listenerManager.IsRunning(listenerName) {
		return
	}
	port, err := s.portAllocator.Get("neptunedata", clusterID)
	if err != nil {
		logs.Warn("failed to get port for cluster listener", logs.String("cluster", clusterID), logs.Err(err))
		return
	}
	s.listenerManager.Register(listener.ListenerConfig{
		Name:        listenerName,
		DefaultPort: port,
		Handler:     s.DataPlaneHandler(clusterID),
	})
	logs.Info("registered cluster listener", logs.String("cluster", clusterID), logs.Int("port", port))
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
		"startTime": s.startTime.UTC().Format(timeutils.ISO8601UTCFormat),
		// M16: dbEngineVersion — the engine version string used by clients
		// for feature detection. Neptune 1.x series.
		"dbEngineVersion": "1.3.3.0",
		"role":            "writer",
		// M16: dfeQueryEngine — DFE (Distributed Forwarding Engine) state.
		// "Disabled" indicates TinkerPop-only execution (no DFE).
		"dfeQueryEngine": "Disabled",
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
		// M16: rollingBackTrxCount — number of currently rolling-back
		// transactions. 0 when healthy.
		"rollingBackTrxCount": 0,
		// M16: rollingBackTrxEarliestStartTime — empty when no rollbacks.
		"rollingBackTrxEarliestStartTime": "",
		// M16: features — engine feature flags (e.g. Streams, ML).
		"features": map[string]interface{}{
			"streams": map[string]interface{}{
				"property_graph": map[string]interface{}{
					"results": map[string]interface{}{
						"status": "enabled",
					},
				},
				"sparql": map[string]interface{}{
					"results": map[string]interface{}{
						"status": "disabled",
					},
				},
			},
		},
		"settings": map[string]interface{}{
			"neptune lab mode": "DISABLED",
		},
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
			if t, ok := value.(time.Time); ok {
				if now.After(t) {
					s.fastTokens.Delete(key)
				}
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
		if !ok {
			return nil, preconditionFailed("invalid or expired token")
		}
		expiry, typeOk := val.(time.Time)
		if !typeOk || time.Now().After(expiry) {
			s.fastTokens.Delete(params.Token)
			return nil, preconditionFailed("invalid or expired token")
		}
		s.fastTokens.Delete(params.Token)

		if gs, ok := reqCtx.GraphWriter().(graphengine.GraphStore); ok {
			if err := gs.Clear(); err != nil {
				logs.Warn("failed to clear graph store during fast reset", logs.Err(err))
			}
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
				if err := resetStore.DeleteQuery(q.GetQueryId()); err != nil {
					logs.Warn("failed to delete query during reset", logs.String("queryId", q.GetQueryId()), logs.Err(err))
				}
			}
			jobs, _ := resetStore.ListLoaderJobs()
			for _, j := range jobs {
				if err := resetStore.DeleteLoaderJob(j.GetLoadId()); err != nil {
					logs.Warn("failed to delete loader job during reset", logs.String("loadId", j.GetLoadId()), logs.Err(err))
				}
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

// getQueryStatus returns the status and evaluation statistics of a query
// identified by queryId. Shared by both Gremlin and OpenCypher query status
// handlers.
func (s *NeptuneDataService) getQueryStatus(reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryId := getPathParam(req, "queryId")
	if queryId == "" {
		return nil, missingParameter("queryId")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	qr, err := store.GetQuery(queryId)
	if err != nil || qr == nil {
		return nil, badRequest(fmt.Sprintf("query not found: %s", queryId))
	}

	var elapsed int64
	if qr.EndTime != nil && qr.StartTime != nil {
		elapsed = qr.EndTime.AsTime().Sub(qr.StartTime.AsTime()).Milliseconds()
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

// listQueries returns all submitted queries of the given type, optionally
// including those in a waiting state. Shared by both Gremlin and OpenCypher
// list queries handlers.
func (s *NeptuneDataService) listQueries(reqCtx *request.RequestContext, req *request.ParsedRequest, queryType string) (interface{}, error) {
	includeWaiting := request.GetBoolParam(req.Parameters, "includeWaiting")

	store, err := s.store(reqCtx)
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
		if qr.GetQueryType() != queryType {
			continue
		}
		st := qr.GetStatus()
		if st == "complete" || st == "failed" || st == "cancelled" {
			continue
		}
		if st == "waiting" && !includeWaiting {
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

// cancelQuery cancels a running query identified by queryId. Shared by both
// Gremlin and OpenCypher cancel handlers. If silent is true, an empty body is
// returned instead of the standard cancellation confirmation.
func (s *NeptuneDataService) cancelQuery(reqCtx *request.RequestContext, req *request.ParsedRequest, silent bool, includePayload bool) (interface{}, error) {
	queryId := getPathParam(req, "queryId")
	if queryId == "" {
		return nil, missingParameter("queryId")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	qr, err := store.GetQuery(queryId)
	if err != nil || qr == nil {
		return nil, badRequest(fmt.Sprintf("query not found: %s", queryId))
	}
	switch qr.GetStatus() {
	case "complete", "failed", "cancelled":
		return nil, badRequest(fmt.Sprintf("cannot cancel query in terminal state: %s", qr.GetStatus()))
	}
	qr.Status = "cancelled"
	qr.EndTime = timestamppb.Now()
	if err := store.UpdateQuery(qr); err != nil {
		logs.Warn("failed to persist query cancellation", logs.String("queryId", queryId), logs.Err(err))
	}

	if silent {
		return map[string]interface{}{}, nil
	}

	resp := map[string]interface{}{
		"status": "200 OK",
	}
	if includePayload {
		resp["payload"] = true
	}
	return resp, nil
}

func (s *NeptuneDataService) getStats(region string) *GraphStatistics {
	val, _ := s.statsMap.LoadOrStore(region, &GraphStatistics{
		LabelCounts: make(map[string]int64),
		RelCounts:   make(map[string]int64),
		LastAccess:  time.Now(),
		LastRefresh: time.Now(),
	})
	st := val.(*GraphStatistics)
	st.mu.Lock()
	st.LastAccess = time.Now()
	st.mu.Unlock()
	return st
}

func (s *NeptuneDataService) refreshStatistics(reqCtx *request.RequestContext) {
	var reader graphengine.GraphReader
	var region string

	if reqCtx != nil {
		if r := reqCtx.GraphReader(); r != nil {
			reader = r
		}
		region = reqCtx.GetRegion()
	}

	refreshStatisticsWithReader(reader, region, s)
}

// refreshStatisticsForRegion refreshes stats for a given region using the
// first available cluster engine. Used by the admin handler which has no
// per-request graph context.
func (s *NeptuneDataService) refreshStatisticsForRegion(region string) {
	s.enginesMu.RLock()
	var reader graphengine.GraphReader
	for _, entry := range s.activeEngines {
		reader = entry.db
		break
	}
	s.enginesMu.RUnlock()
	refreshStatisticsWithReader(reader, region, s)
}

func refreshStatisticsWithReader(reader graphengine.GraphReader, region string, s *NeptuneDataService) {
	if reader == nil {
		// M21: Log when no graph reader is available so silent failures
		// are diagnosable.
		logs.Warn("refreshStatistics: no graph reader available",
			logs.String("region", region))
		return
	}

	stats := s.getStats(region)
	nodeCount := reader.CountNodes()
	edgeCount := reader.CountEdges()
	labelCounts, _ := reader.GetLabelCounts()
	relCounts, _ := reader.GetRelCounts()
	stats.mu.Lock()
	stats.NodeCount = nodeCount
	stats.EdgeCount = edgeCount
	stats.LabelCounts = labelCounts
	stats.RelCounts = relCounts
	stats.LastRefresh = time.Now()
	stats.mu.Unlock()
}

func (st *GraphStatistics) snapshot() (nodeCount, edgeCount int64, labelCounts, relCounts map[string]int64) {
	st.mu.Lock()
	nodeCount = st.NodeCount
	edgeCount = st.EdgeCount
	labelCounts = make(map[string]int64, len(st.LabelCounts))
	for k, v := range st.LabelCounts {
		labelCounts[k] = v
	}
	relCounts = make(map[string]int64, len(st.RelCounts))
	for k, v := range st.RelCounts {
		relCounts[k] = v
	}
	st.mu.Unlock()
	return
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

// startLoaderDispatcher launches the loader job dispatcher if it is not
// already running. The dispatcher scans for LOAD_QUEUED jobs whose
// dependencies are satisfied and starts them, respecting the maxConcurrentLoads
// limit. It exits when no more queued jobs remain.
func (s *NeptuneDataService) startLoaderDispatcher() {
	if s.dispatcherRunning.Swap(true) {
		return
	}
	go s.loaderDispatchLoop()
}

// loaderDispatchLoop repeatedly scans for dispatchable queued loader jobs
// until none remain. Uses a short poll interval to keep latency low without
// busy-waiting.
func (s *NeptuneDataService) loaderDispatchLoop() {
	defer s.dispatcherRunning.Store(false)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		dispatched, err := s.dispatchQueuedJobs()
		if err != nil {
			logs.Warn("loader dispatcher error", logs.Err(err))
		}
		if dispatched == 0 {
			jobs, err := s.countQueuedJobs()
			if err != nil || jobs == 0 {
				return
			}
		}
	}
}

// countQueuedJobs counts LOAD_QUEUED jobs across all regional stores.
func (s *NeptuneDataService) countQueuedJobs() (int, error) {
	count := 0
	s.stores.Range(func(_, val any) bool {
		store, ok := val.(*neptunestore.NeptuneStore)
		if !ok {
			return true
		}
		jobs, err := store.ListLoaderJobs()
		if err != nil {
			return true
		}
		for _, job := range jobs {
			if job.GetStatus() == "LOAD_QUEUED" {
				count++
			}
		}
		return true
	})
	return count, nil
}

// dispatchQueuedJobs scans all stores for LOAD_QUEUED jobs whose dependencies
// are satisfied and launches them, respecting the concurrency limit.
// Returns the number of jobs dispatched in this pass.
func (s *NeptuneDataService) dispatchQueuedJobs() (int, error) {
	dispatched := 0

	s.stores.Range(func(regionVal, val any) bool {
		region, rok := regionVal.(string)
		store, ok := val.(*neptunestore.NeptuneStore)
		if !ok || !rok {
			return true
		}

		jobs, err := store.ListLoaderJobs()
		if err != nil {
			return true
		}

		running := 0
		for _, job := range jobs {
			if job.GetStatus() == "LOAD_IN_PROGRESS" {
				running++
			}
		}

		for _, job := range jobs {
			if job.GetStatus() != "LOAD_QUEUED" {
				continue
			}
			if running >= maxConcurrentLoads {
				break
			}

			if !s.dependenciesSatisfied(store, job) {
				continue
			}

			clusterDB := s.GetClusterEngineForRegion(region)

			// Re-read the job from the store to detect a concurrent
			// CancelLoaderJob before transitioning to LOAD_IN_PROGRESS.
			// Without this, the dispatcher's snapshot can overwrite a
			// CANCELLED status set between ListLoaderJobs and
			// UpdateLoaderJob.
			current, err := store.GetLoaderJob(job.GetLoadId())
			if err != nil || current == nil || current.GetStatus() != "LOAD_QUEUED" {
				continue
			}

			current.Status = "LOAD_IN_PROGRESS"
			_ = store.UpdateLoaderJob(current)

			s.launchLoaderJob(region, current, clusterDB)
			running++
			dispatched++
		}
		return true
	})

	return dispatched, nil
}

// dependenciesSatisfied checks whether all dependency load IDs are in a
// LOAD_COMPLETED state. If any dependency failed, the job is marked
// LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED.
func (s *NeptuneDataService) dependenciesSatisfied(store *neptunestore.NeptuneStore, job *pb.LoaderJob) bool {
	for _, depID := range job.GetDependencies() {
		dep, err := store.GetLoaderJob(depID)
		if err != nil || dep == nil {
			job.Status = "LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED"
			job.EndTime = timestamppb.Now()
			job.ErrorLog = fmt.Sprintf("dependency %s not found", depID)
			_ = store.UpdateLoaderJob(job)
			return false
		}
		depStatus := dep.GetStatus()
		if depStatus == "LOAD_FAILED" || depStatus == "CANCELLED" || depStatus == "LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED" {
			job.Status = "LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED"
			job.EndTime = timestamppb.Now()
			job.ErrorLog = fmt.Sprintf("dependency %s is in failed state: %s", depID, depStatus)
			_ = store.UpdateLoaderJob(job)
			return false
		}
		if depStatus != "LOAD_COMPLETED" {
			return false
		}
	}
	return true
}

// GetClusterEngineForRegion returns the graph engine DB for a cluster in the
// given region. Used by the loader dispatcher which operates outside of a
// per-request context. If multiple clusters exist in the same region, the
// first match is returned.
func (s *NeptuneDataService) GetClusterEngineForRegion(region string) *graphengine.DB {
	s.enginesMu.RLock()
	defer s.enginesMu.RUnlock()
	for _, entry := range s.activeEngines {
		if entry.region == region {
			return entry.db
		}
	}
	return nil
}
