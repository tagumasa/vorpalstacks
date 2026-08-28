package neptunegraph

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/serviceports"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/graphengine"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
	"vorpalstacks/internal/utils/aws/arn"
)

const (
	graphDataDirPrefix  = "neptunegraph/graphs"
	taskCleanupInterval = 24 * time.Hour
	taskTTL             = 24 * time.Hour
	taskTTLTestMode     = 30 * time.Minute
	taskCleanupTestMode = 5 * time.Minute

	// neptuneGraphBuildNumber is the build identifier returned in graph
	// responses. Neptune Analytics (neptunegraph) is a distinct AWS service
	// from Neptune Database and maintains its own release pipeline; its
	// build number has no relationship to the Neptune DB engine version
	// (1.x.y.z). The value below mirrors the AWS Neptune Analytics build
	// identifier format (major.minor.YYYYMMDD).
	neptuneGraphBuildNumber = "1.0.20250313"
)

// NeptuneGraphService implements NeptuneGraph API operations, managing graphs, snapshots, endpoints, and tasks.
type NeptuneGraphService struct {
	accountID      string
	region         string
	dataPath       string
	storageManager *storage.RegionStorageManager
	stores         sync.Map
	activeEngines  map[string]*engineEntry
	enginesMu      sync.RWMutex
	taskWg         sync.WaitGroup
	graphCache     *graphengine.Cache
	arnBuilder     *arn.ARNBuilder
	eventBus       *eventbus.EventBus
	regionCleanups sync.Map
	testMode       bool
	planCache      *queryPlanCache
}

type engineEntry struct {
	db      *graphengine.DB
	mu      sync.RWMutex
	stopped bool
	wg      sync.WaitGroup
}

// NewNeptuneGraphService creates a new NeptuneGraphService instance.
func NewNeptuneGraphService(accountID, region, dataPath string) *NeptuneGraphService {
	return &NeptuneGraphService{
		accountID:     accountID,
		region:        region,
		dataPath:      dataPath,
		activeEngines: make(map[string]*engineEntry),
		arnBuilder:    arn.NewARNBuilder(accountID, region),
		testMode:      os.Getenv("TEST_MODE") == "true",
		planCache:     newQueryPlanCache(planCacheCapacity, planCacheTTLSeconds*time.Second),
	}
}

// Close shuts down all active graph engines and releases associated resources.
func (s *NeptuneGraphService) Close() {
	s.enginesMu.Lock()
	entries := make(map[string]*engineEntry, len(s.activeEngines))
	for id, entry := range s.activeEngines {
		entry.stopped = true
		entries[id] = entry
	}
	s.activeEngines = make(map[string]*engineEntry)
	s.enginesMu.Unlock()

	for _, entry := range entries {
		entry.wg.Wait()
		entry.db.Close()
	}
	s.taskWg.Wait()
}

// SetEventBus injects the shared event bus for cross-service invocations
// such as EC2 subnet lookups.
func (s *NeptuneGraphService) SetEventBus(bus *eventbus.EventBus) {
	s.eventBus = bus
}

// ExecuteQueryOnGraph executes a query against a specific graph engine,
// intended for cross-service callers (e.g. AppSync resolvers via the event bus).
func (s *NeptuneGraphService) ExecuteQueryOnGraph(ctx context.Context, graphID string, query string, language string, parameters map[string]interface{}) (interface{}, error) {
	s.enginesMu.Lock()
	entry, ok := s.activeEngines[graphID]
	if !ok || entry.stopped {
		s.enginesMu.Unlock()
		return nil, fmt.Errorf("graph %q is not available", graphID)
	}
	entry.wg.Add(1)
	s.enginesMu.Unlock()
	defer entry.wg.Done()

	entry.mu.RLock()
	defer entry.mu.RUnlock()

	store, err := s.GetStoreForRegion(s.region)
	if err != nil {
		return nil, err
	}

	return executeCypherQuery(ctx, s, nil, &request.ParsedRequest{
		Body:       jsonQueryBody(query, language, parameters),
		Parameters: map[string]interface{}{"graphidentifier": graphID},
	}, graphID, entry, store)
}

func jsonQueryBody(query, language string, parameters map[string]interface{}) json.RawMessage {
	body := map[string]interface{}{
		"query":    query,
		"language": language,
	}
	if len(parameters) > 0 {
		body["parameters"] = parameters
	}
	b, _ := json.Marshal(body)
	return b
}

// RestoreEngines reopens graph engines for all AVAILABLE graphs after a service restart.
func (s *NeptuneGraphService) RestoreEngines() {
	store, err := s.GetStoreForRegion(s.region)
	if err != nil {
		logs.Warn("failed to get store for engine restore", logs.Err(err))
		return
	}

	graphs, _, _, err := store.ListGraphs(storecommon.ListOptions{})
	if err != nil {
		logs.Warn("failed to list graphs for engine restore", logs.Err(err))
		return
	}

	s.enginesMu.Lock()
	defer s.enginesMu.Unlock()

	for _, g := range graphs {
		if g.Status != "AVAILABLE" {
			continue
		}
		bucket, err := s.graphBucket(g.Id)
		if err != nil {
			logs.Warn("failed to get graph bucket", logs.String("graphId", g.Id), logs.Err(err))
			continue
		}
		db, err := graphengine.New(bucket, s.engineOptions())
		if err != nil {
			logs.Warn("failed to restore graph engine", logs.String("graphId", g.Id), logs.Err(err))
			continue
		}
		s.activeEngines[g.Id] = &engineEntry{db: db}
		logs.Info("restored graph engine", logs.String("graphId", g.Id))
	}
}

// SetStorageManager injects the region storage manager used to back persistent stores.
func (s *NeptuneGraphService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

func (s *NeptuneGraphService) graphBucket(graphID string) (storage.BatchBucket, error) {
	if s.storageManager == nil {
		return nil, fmt.Errorf("storage manager not set")
	}
	rs, err := s.storageManager.GetStorage(s.region)
	if err != nil {
		return nil, err
	}
	bkt := rs.Bucket("neptunegraph:graph:" + graphID)
	bb, ok := bkt.(storage.BatchBucket)
	if !ok {
		return nil, fmt.Errorf("storage bucket does not support batch operations")
	}
	return bb, nil
}

// SetGraphCache injects a shared Pebble block cache for all graph engine
// instances. Must be called before RestoreEngines or any graph creation.
func (s *NeptuneGraphService) SetGraphCache(cache *graphengine.Cache) {
	s.graphCache = cache
}

func (s *NeptuneGraphService) engineOptions() graphengine.Options {
	opts := graphengine.DefaultOptions()
	if s.graphCache != nil {
		opts.SharedCache = s.graphCache
	}
	return opts
}

func (s *NeptuneGraphService) graphEndpoint(graphID string) string {
	return fmt.Sprintf("%s.graph.%s.vorpalstacks.localhost:%d", graphID, s.region, serviceports.HTTP)
}

// GetStoreForRegion returns a lazily-initialised NeptuneGraphStore for the given region.
func (s *NeptuneGraphService) GetStoreForRegion(region string) (*ngstore.NeptuneGraphStore, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (*ngstore.NeptuneGraphStore, error) {
		if s.storageManager == nil {
			return nil, fmt.Errorf("storage manager not set")
		}
		rs, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}
		return ngstore.NewNeptuneGraphStore(rs), nil
	})
}

func (s *NeptuneGraphService) store(reqCtx *request.RequestContext) (*ngstore.NeptuneGraphStore, error) {
	region := reqCtx.GetRegion()
	st, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	lastVal, loaded := s.regionCleanups.LoadOrStore(region, now)
	shouldCleanup := !loaded
	if loaded {
		interval := taskCleanupInterval
		if s.testMode {
			interval = taskCleanupTestMode
		}
		if now.Sub(lastVal.(time.Time)) >= interval {
			shouldCleanup = true
		}
	}
	if shouldCleanup {
		s.regionCleanups.Store(region, now)
		go s.cleanupExpiredTasks(st)
	}

	return st, nil
}

// RegisterHandlers registers all NeptuneGraph operation handlers with the dispatcher.
func (s *NeptuneGraphService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("neptunegraph", "CreateGraph", s.CreateGraph)
	d.RegisterHandlerForService("neptunegraph", "GetGraph", s.GetGraph)
	d.RegisterHandlerForService("neptunegraph", "ListGraphs", s.ListGraphs)
	d.RegisterHandlerForService("neptunegraph", "UpdateGraph", s.UpdateGraph)
	d.RegisterHandlerForService("neptunegraph", "DeleteGraph", s.DeleteGraph)
	d.RegisterHandlerForService("neptunegraph", "StartGraph", s.StartGraph)
	d.RegisterHandlerForService("neptunegraph", "StopGraph", s.StopGraph)
	d.RegisterHandlerForService("neptunegraph", "ResetGraph", s.ResetGraph)
	d.RegisterHandlerForService("neptunegraph", "RestoreGraphFromSnapshot", s.RestoreGraphFromSnapshot)
	d.RegisterHandlerForService("neptunegraph", "CreateGraphSnapshot", s.CreateGraphSnapshot)
	d.RegisterHandlerForService("neptunegraph", "GetGraphSnapshot", s.GetGraphSnapshot)
	d.RegisterHandlerForService("neptunegraph", "ListGraphSnapshots", s.ListGraphSnapshots)
	d.RegisterHandlerForService("neptunegraph", "DeleteGraphSnapshot", s.DeleteGraphSnapshot)
	d.RegisterHandlerForService("neptunegraph", "CreatePrivateGraphEndpoint", s.CreatePrivateGraphEndpoint)
	d.RegisterHandlerForService("neptunegraph", "GetPrivateGraphEndpoint", s.GetPrivateGraphEndpoint)
	d.RegisterHandlerForService("neptunegraph", "ListPrivateGraphEndpoints", s.ListPrivateGraphEndpoints)
	d.RegisterHandlerForService("neptunegraph", "DeletePrivateGraphEndpoint", s.DeletePrivateGraphEndpoint)
	d.RegisterHandlerForService("neptunegraph", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("neptunegraph", "TagResource", s.TagResource)
	d.RegisterHandlerForService("neptunegraph", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("neptunegraph", "CreateGraphUsingImportTask", s.CreateGraphUsingImportTask)
	d.RegisterHandlerForService("neptunegraph", "GetImportTask", s.GetImportTask)
	d.RegisterHandlerForService("neptunegraph", "ListImportTasks", s.ListImportTasks)
	d.RegisterHandlerForService("neptunegraph", "CancelImportTask", s.CancelImportTask)
	d.RegisterHandlerForService("neptunegraph", "StartImportTask", s.StartImportTask)
	d.RegisterHandlerForService("neptunegraph", "StartExportTask", s.StartExportTask)
	d.RegisterHandlerForService("neptunegraph", "GetExportTask", s.GetExportTask)
	d.RegisterHandlerForService("neptunegraph", "ListExportTasks", s.ListExportTasks)
	d.RegisterHandlerForService("neptunegraph", "CancelExportTask", s.CancelExportTask)
	d.RegisterHandlerForService("neptunegraph", "ExecuteQuery", s.ExecuteQuery)
	d.RegisterHandlerForService("neptunegraph", "GetQuery", s.GetQuery)
	d.RegisterHandlerForService("neptunegraph", "ListQueries", s.ListQueries)
	d.RegisterHandlerForService("neptunegraph", "CancelQuery", s.CancelQuery)
	d.RegisterHandlerForService("neptunegraph", "GetGraphSummary", s.GetGraphSummary)
}

func generateID(prefix string) string {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand.Read failed: %v", err))
	}
	return prefix + hex.EncodeToString(b)
}

// CreateGraph creates a new NeptuneGraph graph resource and initialises its query engine.
func (s *NeptuneGraphService) CreateGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &CreateGraphInput{
		GraphName:          request.GetStringParam(req.Parameters, "graphName"),
		ProvisionedMemory:  request.GetIntParam(req.Parameters, "provisionedMemory"),
		HasReplicaCount:    request.HasParam(req.Parameters, "replicaCount"),
		ReplicaCount:       request.GetIntParam(req.Parameters, "replicaCount"),
		KmsKeyIdentifier:   request.GetStringParam(req.Parameters, "kmsKeyIdentifier"),
		DeletionProtection: request.GetBoolParam(req.Parameters, "deletionProtection"),
		PublicConnectivity: request.GetBoolParam(req.Parameters, "publicConnectivity"),
		VectorSearchConfig: req.Parameters["vectorSearchConfiguration"],
		Tags:               parseTagsFromParams(req.Parameters),
		Region:             reqCtx.GetRegion(),
	}

	graph, err := s.createGraphCore(ctx, store, in)
	if err != nil {
		return nil, err
	}
	return graphToResponse(graph), nil
}

// GetGraph retrieves a graph by its identifier.
func (s *NeptuneGraphService) GetGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &GetGraphInput{
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
	}

	graph, err := s.getGraphCore(store, in)
	if err != nil {
		return nil, err
	}
	return graphToResponse(graph), nil
}

// ListGraphs returns a paginated list of graph summaries.
func (s *NeptuneGraphService) ListGraphs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &ListGraphsInput{
		MaxItems: clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		Marker:   request.GetStringParam(req.Parameters, "nextToken"),
	}

	res, err := s.listGraphsCore(store, in)
	if err != nil {
		return nil, err
	}

	summaries := make([]interface{}, 0, len(res.Graphs))
	for _, g := range res.Graphs {
		summaries = append(summaries, graphSummaryToResponse(g))
	}

	result := map[string]interface{}{
		"graphs": summaries,
	}
	if res.Truncated {
		result["nextToken"] = res.NextToken
	}
	return result, nil
}

// UpdateGraph modifies configuration of an existing graph that is in AVAILABLE state.
func (s *NeptuneGraphService) UpdateGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &UpdateGraphInput{
		GraphIdentifier:       request.GetStringParam(req.Parameters, "graphIdentifier"),
		HasProvisionedMemory:  request.HasParam(req.Parameters, "provisionedMemory"),
		ProvisionedMemory:     request.GetIntParam(req.Parameters, "provisionedMemory"),
		HasDeletionProtection: request.HasParam(req.Parameters, "deletionProtection"),
		DeletionProtection:    request.GetBoolParam(req.Parameters, "deletionProtection"),
		HasPublicConnectivity: request.HasParam(req.Parameters, "publicConnectivity"),
		PublicConnectivity:    request.GetBoolParam(req.Parameters, "publicConnectivity"),
	}

	graph, err := s.updateGraphCore(store, in)
	if err != nil {
		return nil, err
	}
	return graphToResponse(graph), nil
}

// DeleteGraph removes a graph, optionally creating an automatic snapshot, and cleans up all associated resources.
func (s *NeptuneGraphService) DeleteGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &DeleteGraphInput{
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
		HasSkipSnapshot: request.HasParam(req.Parameters, "skipSnapshot"),
		SkipSnapshot:    request.GetBoolParam(req.Parameters, "skipSnapshot"),
		Region:          reqCtx.GetRegion(),
	}

	graph, err := s.deleteGraphCore(store, in)
	if err != nil {
		return nil, err
	}
	return graphToResponse(graph), nil
}

// StartGraph transitions a STOPPED graph to AVAILABLE by reopening its query engine.
func (s *NeptuneGraphService) StartGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &StartGraphInput{
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
	}

	graph, err := s.startGraphCore(store, in)
	if err != nil {
		return nil, err
	}
	return graphToResponse(graph), nil
}

// StopGraph gracefully shuts down a graph's query engine and transitions it to STOPPED state.
func (s *NeptuneGraphService) StopGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &StopGraphInput{
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
	}

	graph, err := s.stopGraphCore(store, in)
	if err != nil {
		return nil, err
	}
	return graphToResponse(graph), nil
}

// ResetGraph clears all data from an AVAILABLE graph's engine while keeping the graph resource intact.
func (s *NeptuneGraphService) ResetGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &ResetGraphInput{
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
		HasSkipSnapshot: request.HasParam(req.Parameters, "skipSnapshot"),
		SkipSnapshot:    request.GetBoolParam(req.Parameters, "skipSnapshot"),
	}

	graph, err := s.resetGraphCore(store, in)
	if err != nil {
		return nil, err
	}
	return graphToResponse(graph), nil
}

// RestoreGraphFromSnapshot creates a new graph from an existing snapshot, copying source graph data.
func (s *NeptuneGraphService) RestoreGraphFromSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &RestoreGraphFromSnapshotInput{
		SnapshotIdentifier:   request.GetStringParam(req.Parameters, "snapshotIdentifier"),
		GraphName:            request.GetStringParam(req.Parameters, "graphName"),
		DeletionProtection:   request.GetBoolParam(req.Parameters, "deletionProtection"),
		PublicConnectivity:   request.GetBoolParam(req.Parameters, "publicConnectivity"),
		HasProvisionedMemory: request.HasParam(req.Parameters, "provisionedMemory"),
		ProvisionedMemory:    request.GetIntParam(req.Parameters, "provisionedMemory"),
		HasReplicaCount:      request.HasParam(req.Parameters, "replicaCount"),
		ReplicaCount:         request.GetIntParam(req.Parameters, "replicaCount"),
		Tags:                 parseTagsFromParams(req.Parameters),
		Region:               reqCtx.GetRegion(),
	}

	graph, err := s.restoreGraphFromSnapshotCore(store, in)
	if err != nil {
		return nil, err
	}
	return graphToResponse(graph), nil
}
