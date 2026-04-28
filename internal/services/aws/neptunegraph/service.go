package neptunegraph

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/graphengine"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
	"vorpalstacks/internal/utils/aws/arn"
)

const (
	graphNamePattern     = `^[a-z][a-z0-9]*(-[a-z0-9]+)*$`
	snapshotNamePattern  = `^[a-z][a-z0-9]*(-[a-z0-9]+)*$`
	graphDataDirPrefix   = "neptunegraph/graphs"
	minProvisionedMemory = 16
	maxProvisionedMemory = 24576
	maxReplicaCount      = 2
	maxTags              = 50
	tagKeyPattern        = `^[a-zA-Z+\-=._:/]+$`
)

var (
	graphNameRegex    = regexp.MustCompile(graphNamePattern)
	snapshotNameRegex = regexp.MustCompile(snapshotNamePattern)
	tagKeyRegex       = regexp.MustCompile(tagKeyPattern)
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
}

type engineEntry struct {
	db      *graphengine.DB
	mu      sync.RWMutex
	stopped bool
	wg      sync.WaitGroup
}

// NewNeptuneGraphService creates a new NeptuneGraphService for the given account, region, and data path.
func NewNeptuneGraphService(accountID, region, dataPath string) *NeptuneGraphService {
	return &NeptuneGraphService{
		accountID:     accountID,
		region:        region,
		dataPath:      dataPath,
		activeEngines: make(map[string]*engineEntry),
		arnBuilder:    arn.NewARNBuilder(accountID, region),
	}
}

// Close shuts down all active graph engines and releases associated resources.
func (s *NeptuneGraphService) Close() {
	s.enginesMu.Lock()
	for id, entry := range s.activeEngines {
		entry.stopped = true
		entry.wg.Wait()
		entry.db.Close()
		delete(s.activeEngines, id)
	}
	s.enginesMu.Unlock()
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
	return s.GetStoreForRegion(region)
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

	graphName := request.GetStringParam(req.Parameters, "graphName")
	if graphName == "" || strings.HasPrefix(graphName, "g-") || !graphNameRegex.MatchString(graphName) || len(graphName) > 63 {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphName")
	}

	mem := request.GetIntParam(req.Parameters, "provisionedMemory")
	if mem < minProvisionedMemory || mem > maxProvisionedMemory {
		mem = 128
	}

	replicaCount := request.GetIntParam(req.Parameters, "replicaCount")
	if !request.HasParam(req.Parameters, "replicaCount") {
		replicaCount = 1
	} else if replicaCount < 0 || replicaCount > maxReplicaCount {
		replicaCount = 1
	}

	region := reqCtx.GetRegion()
	graphID := generateID("g-")
	now := time.Now().UTC()

	graph := &ngstore.Graph{
		Id:                 graphID,
		Name:               graphName,
		Arn:                s.arnBuilder.NeptuneGraph().Graph(graphID),
		Status:             "CREATING",
		ProvisionedMemory:  int32Ptr(int32(mem)),
		ReplicaCount:       int32Ptr(int32(replicaCount)),
		DeletionProtection: request.GetBoolParam(req.Parameters, "deletionProtection"),
		PublicConnectivity: request.GetBoolParam(req.Parameters, "publicConnectivity"),
		KmsKeyIdentifier:   request.GetStringParam(req.Parameters, "kmsKeyIdentifier"),
		BuildNumber:        "1.0.20250313",
		CreateTime:         &now,
		AccountID:          s.accountID,
		Region:             region,
	}

	if vsc := parseVectorSearchConfig(req.Parameters); vsc != nil {
		graph.VectorSearchConfiguration = vsc
	}

	if err := store.CreateGraph(graph); err != nil {
		if ngstore.IsAlreadyExists(err) {
			return nil, newConflictException("CONCURRENT_MODIFICATION")
		}
		return nil, err
	}

	bucket, err := s.graphBucket(graphID)
	if err != nil {
		return nil, newInternalServerException(err)
	}
	db, err := graphengine.New(bucket, s.engineOptions())
	if err != nil {
		logs.Warn("failed to open graph engine", logs.String("graphId", graphID), logs.Err(err))
		graph.Status = "FAILED"
		graph.StatusReason = fmt.Sprintf("failed to open graph engine: %v", err)
		if updateErr := store.UpdateGraph(graph); updateErr != nil {
			logs.Warn("failed to update graph status to FAILED", logs.String("graphId", graphID), logs.Err(updateErr))
		}
		return graphToResponse(graph, true), nil
	}
	s.enginesMu.Lock()
	s.activeEngines[graphID] = &engineEntry{db: db}
	s.enginesMu.Unlock()

	graph.Status = "AVAILABLE"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("failed to update graph status", logs.Err(err))
	}

	if tags := parseTagsFromParams(req.Parameters); len(tags) > 0 {
		if err := store.AddTags(graph.Arn, tags); err != nil {
			logs.Warn("failed to store tags for graph", logs.String("graphId", graphID), logs.Err(err))
		}
	}

	return graphToResponse(graph, true), nil
}

// GetGraph retrieves a graph by its identifier.
func (s *NeptuneGraphService) GetGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	return graphToResponse(graph, false), nil
}

// ListGraphs returns a paginated list of graph summaries.
func (s *NeptuneGraphService) ListGraphs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := storecommon.ListOptions{
		MaxItems: clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		Marker:   request.GetStringParam(req.Parameters, "nextToken"),
	}

	graphs, nextToken, truncated, err := store.ListGraphs(opts)
	if err != nil {
		return nil, err
	}

	summaries := make([]interface{}, 0, len(graphs))
	for _, g := range graphs {
		summaries = append(summaries, graphSummaryToResponse(g))
	}

	result := map[string]interface{}{
		"graphs": summaries,
	}
	if truncated {
		result["nextToken"] = nextToken
	}
	return result, nil
}

// UpdateGraph modifies configuration of an existing graph that is in AVAILABLE state.
func (s *NeptuneGraphService) UpdateGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	if graph.Status != "AVAILABLE" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graph is not in AVAILABLE state")
	}

	graph.Status = "UPDATING"

	if request.HasParam(req.Parameters, "provisionedMemory") {
		mem := request.GetIntParam(req.Parameters, "provisionedMemory")
		if mem < minProvisionedMemory || mem > maxProvisionedMemory {
			return nil, newValidationException("ILLEGAL_ARGUMENT", "provisionedMemory")
		}
		graph.ProvisionedMemory = int32Ptr(int32(mem))
	}
	if request.HasParam(req.Parameters, "deletionProtection") {
		graph.DeletionProtection = request.GetBoolParam(req.Parameters, "deletionProtection")
	}
	if request.HasParam(req.Parameters, "publicConnectivity") {
		graph.PublicConnectivity = request.GetBoolParam(req.Parameters, "publicConnectivity")
	}

	if err := store.UpdateGraph(graph); err != nil {
		return nil, err
	}

	graph.Status = "AVAILABLE"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("failed to update graph status to AVAILABLE", logs.String("graphId", graph.Id), logs.Err(err))
	}

	return graphToResponse(graph, false), nil
}

// DeleteGraph removes a graph, optionally creating an automatic snapshot, and cleans up all associated resources.
func (s *NeptuneGraphService) DeleteGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	if graph.DeletionProtection {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graph has deletion protection enabled")
	}

	skipSnapshot := request.GetBoolParam(req.Parameters, "skipSnapshot")

	var snapshotID string
	if !skipSnapshot {
		snapshotID = generateID("gs-")
		now := time.Now().UTC()
		snapshot := &ngstore.GraphSnapshot{
			Id:                 snapshotID,
			Name:               graph.Name + "-auto-snapshot",
			Arn:                s.arnBuilder.NeptuneGraph().Snapshot(snapshotID),
			Status:             "AVAILABLE",
			SourceGraphId:      graphID,
			SnapshotCreateTime: &now,
			AccountID:          s.accountID,
			Region:             graph.Region,
		}
		if err := store.CreateSnapshot(snapshot); err != nil {
			logs.Warn("failed to create auto snapshot during graph deletion", logs.String("graphId", graphID), logs.Err(err))
		}
	}

	graph.Status = "DELETING"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("failed to update graph status to DELETING", logs.String("graphId", graphID), logs.Err(err))
	}

	if graph.Arn != "" {
		existingTags, _ := store.GetTags(graph.Arn)
		if len(existingTags) > 0 {
			keys := make([]string, 0, len(existingTags))
			for k := range existingTags {
				keys = append(keys, k)
			}
			if err := store.RemoveTags(graph.Arn, keys); err != nil {
				logs.Warn("failed to remove tags during graph deletion", logs.String("arn", graph.Arn), logs.Err(err))
			}
		}
	}

	s.enginesMu.Lock()
	if entry, ok := s.activeEngines[graphID]; ok {
		entry.stopped = true
		entry.wg.Wait()
		entry.db.Close()
		delete(s.activeEngines, graphID)
	}
	s.enginesMu.Unlock()

	if snapshotID != "" {
		srcBkt, srcErr := s.graphBucket(graphID)
		dstBkt, dstErr := s.graphBucket("snapshot:" + snapshotID)
		if srcErr == nil && dstErr == nil {
			if err := copyGraphBucket(srcBkt, dstBkt); err != nil {
				logs.Warn("failed to copy graph data to snapshot bucket during deletion",
					logs.String("graphId", graphID), logs.String("snapshotId", snapshotID), logs.Err(err))
			}
		}
	}

	if rs, err := s.storageManager.GetStorage(s.region); err == nil {
		rs.DeleteBucket("neptunegraph:graph:" + graphID)
	}

	if err := store.DeleteGraph(graphID); err != nil {
		logs.Warn("failed to delete graph", logs.String("graphId", graphID), logs.Err(err))
	}

	return graphToResponse(graph, false), nil
}

// StartGraph transitions a STOPPED graph to AVAILABLE by reopening its query engine.
func (s *NeptuneGraphService) StartGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	if graph.Status != "STOPPED" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graph is not in STOPPED state")
	}

	graph.Status = "STARTING"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("failed to update graph status to STARTING", logs.String("graphId", graphID), logs.Err(err))
	}

	bucket, err := s.graphBucket(graphID)
	if err != nil {
		return nil, newInternalServerException(err)
	}
	db, err := graphengine.New(bucket, s.engineOptions())
	if err != nil {
		graph.Status = "FAILED"
		graph.StatusReason = err.Error()
		if updateErr := store.UpdateGraph(graph); updateErr != nil {
			logs.Error("failed to update graph status to FAILED after engine open failure", logs.String("graphId", graphID), logs.Err(updateErr))
		}
		return nil, newInternalServerException(err)
	}

	s.enginesMu.Lock()
	s.activeEngines[graphID] = &engineEntry{db: db}
	s.enginesMu.Unlock()

	graph.Status = "AVAILABLE"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("Failed to update graph status to AVAILABLE", logs.Err(err))
	}

	return graphToResponse(graph, false), nil
}

// StopGraph gracefully shuts down a graph's query engine and transitions it to STOPPED state.
func (s *NeptuneGraphService) StopGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	if graph.Status != "AVAILABLE" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graph is not in AVAILABLE state")
	}

	graph.Status = "STOPPING"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("Failed to update graph status to STOPPING", logs.Err(err))
	}

	s.enginesMu.Lock()
	if entry, ok := s.activeEngines[graphID]; ok {
		entry.stopped = true
		entry.wg.Wait()
		entry.db.Close()
		delete(s.activeEngines, graphID)
	}
	s.enginesMu.Unlock()

	graph.Status = "STOPPED"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("Failed to update graph status to STOPPED", logs.Err(err))
	}

	return graphToResponse(graph, false), nil
}

// ResetGraph clears all data from an AVAILABLE graph's engine while keeping the graph resource intact.
func (s *NeptuneGraphService) ResetGraph(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	if graph.Status != "AVAILABLE" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graph is not in AVAILABLE state")
	}

	graph.Status = "RESETTING"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("Failed to update graph status to RESETTING", logs.Err(err))
	}

	s.enginesMu.RLock()
	entry, ok := s.activeEngines[graphID]

	if ok {
		if err := entry.db.Clear(); err != nil {
			s.enginesMu.RUnlock()

			graph.Status = "FAILED"
			graph.StatusReason = err.Error()
			if storeErr := store.UpdateGraph(graph); storeErr != nil {
				logs.Warn("Failed to update graph status to FAILED after Clear error",
					logs.String("graphId", graphID),
					logs.Err(storeErr))
			}
			return nil, newInternalServerException(err)
		}
	}
	s.enginesMu.RUnlock()

	graph.Status = "AVAILABLE"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("Failed to update graph status to AVAILABLE", logs.Err(err))
	}

	return graphToResponse(graph, false), nil
}

// RestoreGraphFromSnapshot creates a new graph from an existing snapshot, copying source graph data.
func (s *NeptuneGraphService) RestoreGraphFromSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshotID := request.GetStringParam(req.Parameters, "snapshotIdentifier")
	if snapshotID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "snapshotIdentifier")
	}

	graphName := request.GetStringParam(req.Parameters, "graphName")
	if graphName == "" || strings.HasPrefix(graphName, "g-") || !graphNameRegex.MatchString(graphName) || len(graphName) > 63 {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphName")
	}

	_, err = store.GetSnapshot(snapshotID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("snapshot", snapshotID)
		}
		return nil, err
	}

	region := reqCtx.GetRegion()
	graphID := generateID("g-")
	now := time.Now().UTC()

	graph := &ngstore.Graph{
		Id:                 graphID,
		Name:               graphName,
		Arn:                s.arnBuilder.NeptuneGraph().Graph(graphID),
		Status:             "CREATING",
		ProvisionedMemory:  int32Ptr(128),
		ReplicaCount:       int32Ptr(1),
		DeletionProtection: request.GetBoolParam(req.Parameters, "deletionProtection"),
		PublicConnectivity: request.GetBoolParam(req.Parameters, "publicConnectivity"),
		BuildNumber:        "1.0.20250313",
		SourceSnapshotId:   snapshotID,
		CreateTime:         &now,
		AccountID:          s.accountID,
		Region:             region,
	}

	if request.HasParam(req.Parameters, "provisionedMemory") {
		mem := request.GetIntParam(req.Parameters, "provisionedMemory")
		if mem >= minProvisionedMemory && mem <= maxProvisionedMemory {
			graph.ProvisionedMemory = int32Ptr(int32(mem))
		}
	}

	if request.HasParam(req.Parameters, "replicaCount") {
		rc := request.GetIntParam(req.Parameters, "replicaCount")
		if rc >= 0 && rc <= maxReplicaCount {
			graph.ReplicaCount = int32Ptr(int32(rc))
		}
	}

	if err := store.CreateGraph(graph); err != nil {
		if ngstore.IsAlreadyExists(err) {
			return nil, newConflictException("CONCURRENT_MODIFICATION")
		}
		return nil, err
	}

	if tags := parseTagsFromParams(req.Parameters); len(tags) > 0 {
		if err := store.AddTags(graph.Arn, tags); err != nil {
			logs.Warn("failed to store tags for restored graph", logs.String("graphId", graphID), logs.Err(err))
		}
	}

	srcBucket, srcErr := s.graphBucket("snapshot:" + snapshotID)
	dstBucket, dstErr := s.graphBucket(graphID)
	if srcErr == nil && dstErr == nil {
		if err := copyGraphBucket(srcBucket, dstBucket); err != nil {
			logs.Warn("Failed to copy graph data during restore from snapshot",
				logs.String("graphId", graphID),
				logs.String("snapshotId", snapshotID),
				logs.Err(err))
		}
	}

	db, err := graphengine.New(dstBucket, s.engineOptions())
	if err != nil {
		logs.Warn("failed to open graph engine for restored graph", logs.String("graphId", graphID), logs.Err(err))
	} else {
		s.enginesMu.Lock()
		s.activeEngines[graphID] = &engineEntry{db: db}
		s.enginesMu.Unlock()
	}

	graph.Status = "AVAILABLE"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("Failed to update restored graph status", logs.Err(err))
	}

	return graphToResponse(graph, true), nil
}

func copyGraphBucket(src, dst storage.BatchBucket) error {
	return src.ForEach(func(k, v []byte) error {
		return dst.Put(k, v)
	})
}

func parseVectorSearchConfig(params map[string]interface{}) *ngstore.VectorSearchConfig {
	v, ok := params["vectorSearchConfiguration"]
	if !ok || v == nil {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	dim := 0
	if d, ok := m["dimension"]; ok {
		switch v := d.(type) {
		case float64:
			dim = int(v)
		case int:
			dim = v
		}
	}
	if dim < 1 || dim > 65535 {
		return nil
	}
	return &ngstore.VectorSearchConfig{Dimension: int32(dim)}
}

func clampMaxResults(v int) int {
	if v < 1 {
		return 100
	}
	if v > 100 {
		return 100
	}
	return v
}

func int32Ptr(v int32) *int32    { return &v }
func int64Ptr(v int64) *int64    { return &v }
func stringPtr(v string) *string { return &v }
