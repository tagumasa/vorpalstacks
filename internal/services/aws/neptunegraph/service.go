package neptunegraph

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/graphengine"
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
		graphDir := filepath.Join(s.dataPath, graphDataDirPrefix, g.Id)
		db, err := graphengine.Open(graphDir, s.engineOptions())
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
	if replicaCount < 0 || replicaCount > maxReplicaCount {
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

	graphDir := filepath.Join(s.dataPath, graphDataDirPrefix, graphID)
	db, err := graphengine.Open(graphDir, s.engineOptions())
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

	if !skipSnapshot {
		snapshotID := generateID("gs-")
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

	graphDir := filepath.Join(s.dataPath, graphDataDirPrefix, graphID)
	os.RemoveAll(graphDir)

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

	graphDir := filepath.Join(s.dataPath, graphDataDirPrefix, graphID)
	db, err := graphengine.Open(graphDir, s.engineOptions())
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

	snapshot, err := store.GetSnapshot(snapshotID)
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

	graphDir := filepath.Join(s.dataPath, graphDataDirPrefix, graphID)

	srcDir := filepath.Join(s.dataPath, graphDataDirPrefix, snapshot.SourceGraphId)
	if _, err := os.Stat(srcDir); err == nil {
		if err := copyGraphData(srcDir, graphDir); err != nil {
			logs.Warn("Failed to copy graph data during restore from snapshot",
				logs.String("graphId", graphID),
				logs.String("srcDir", srcDir),
				logs.Err(err))
		}
	}

	db, err := graphengine.Open(graphDir, s.engineOptions())
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

func copyGraphData(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	srcGraph := filepath.Join(src, "graph")
	dstGraph := filepath.Join(dst, "graph")
	if _, err := os.Stat(srcGraph); err == nil {
		if err := os.MkdirAll(dstGraph, 0755); err != nil {
			return err
		}
		if err := filepath.Walk(srcGraph, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			rel, _ := filepath.Rel(srcGraph, path)
			destPath := filepath.Join(dstGraph, rel)
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			if err := os.WriteFile(destPath, data, 0644); err != nil {
				logs.Warn("failed to copy graph data file", logs.String("path", destPath), logs.Err(err))
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

// CreateGraphSnapshot creates a point-in-time snapshot of the specified graph.
func (s *NeptuneGraphService) CreateGraphSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	_, err = store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	snapshotName := request.GetStringParam(req.Parameters, "snapshotName")
	if snapshotName == "" || strings.HasPrefix(snapshotName, "gs-") || !snapshotNameRegex.MatchString(snapshotName) || len(snapshotName) > 63 {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "snapshotName")
	}

	region := reqCtx.GetRegion()
	snapshotID := generateID("gs-")
	now := time.Now().UTC()

	snapshot := &ngstore.GraphSnapshot{
		Id:                 snapshotID,
		Name:               snapshotName,
		Arn:                s.arnBuilder.NeptuneGraph().Snapshot(snapshotID),
		Status:             "AVAILABLE",
		SourceGraphId:      graphID,
		SnapshotCreateTime: &now,
		AccountID:          s.accountID,
		Region:             region,
	}

	if err := store.CreateSnapshot(snapshot); err != nil {
		if ngstore.IsAlreadyExists(err) {
			return nil, newConflictException("CONCURRENT_MODIFICATION")
		}
		return nil, err
	}

	return snapshotToResponse(snapshot), nil
}

// GetGraphSnapshot retrieves a snapshot by its identifier.
func (s *NeptuneGraphService) GetGraphSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshotID := request.GetStringParam(req.Parameters, "snapshotIdentifier")
	if snapshotID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "snapshotIdentifier")
	}

	snapshot, err := store.GetSnapshot(snapshotID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("snapshot", snapshotID)
		}
		return nil, err
	}

	if snapshot.Status == "DELETING" {
		return nil, newResourceNotFoundException("snapshot", snapshotID)
	}

	return snapshotToResponse(snapshot), nil
}

// ListGraphSnapshots returns a paginated list of graph snapshots, optionally filtered by graph.
func (s *NeptuneGraphService) ListGraphSnapshots(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := storecommon.ListOptions{
		MaxItems: clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		Marker:   request.GetStringParam(req.Parameters, "nextToken"),
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	snapshots, nextToken, truncated, err := store.ListSnapshots(opts, graphID)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(snapshots))
	for _, sn := range snapshots {
		items = append(items, snapshotToResponse(sn))
	}

	result := map[string]interface{}{
		"graphSnapshots": items,
	}
	if truncated {
		result["nextToken"] = nextToken
	}
	return result, nil
}

// DeleteGraphSnapshot removes a graph snapshot by its identifier.
func (s *NeptuneGraphService) DeleteGraphSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshotID := request.GetStringParam(req.Parameters, "snapshotIdentifier")
	if snapshotID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "snapshotIdentifier")
	}

	snapshot, err := store.GetSnapshot(snapshotID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("snapshot", snapshotID)
		}
		return nil, err
	}

	response := snapshotToResponse(snapshot)

	if err := store.DeleteSnapshot(snapshotID); err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("snapshot", snapshotID)
		}
		return nil, err
	}

	return response, nil
}

// CreatePrivateGraphEndpoint creates a VPC-private endpoint for accessing the specified graph.
func (s *NeptuneGraphService) CreatePrivateGraphEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	_, err = store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	vpcID := request.GetStringParam(req.Parameters, "vpcId")

	if vpcID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "vpcId")
	}

	var subnetIds []string
	if v, ok := req.Parameters["subnetIds"]; ok {
		if arr, ok := v.([]interface{}); ok {
			for _, item := range arr {
				if str, ok := item.(string); ok {
					subnetIds = append(subnetIds, str)
				}
			}
		}
	}

	if len(subnetIds) == 0 {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "subnetIds must not be empty")
	}

	if s.eventBus == nil {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "EC2 service not available for VPC/Subnet validation")
	}
	ec2 := s.eventBus.EC2Invoker()
	if ec2 == nil {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "EC2 service not available for VPC/Subnet validation")
	}
	for _, subnetId := range subnetIds {
		subnetVpcId, _, err := ec2.LookupSubnet(ctx, reqCtx.GetRegion(), subnetId)
		if err != nil {
			return nil, newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("subnet %s not found: %v", subnetId, err))
		}
		if subnetVpcId != vpcID {
			return nil, newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("subnet %s belongs to VPC %s, not %s", subnetId, subnetVpcId, vpcID))
		}
	}

	ep := &ngstore.PrivateGraphEndpoint{
		GraphId:   graphID,
		VpcId:     vpcID,
		SubnetIds: subnetIds,
		Status:    "AVAILABLE",
		AccountID: s.accountID,
		Region:    reqCtx.GetRegion(),
	}

	if err := store.CreateEndpoint(ep); err != nil {
		if ngstore.IsAlreadyExists(err) {
			return nil, newConflictException("CONCURRENT_MODIFICATION")
		}
		return nil, err
	}

	return endpointToResponse(ep), nil
}

// GetPrivateGraphEndpoint retrieves a private graph endpoint by graph and VPC identifiers.
func (s *NeptuneGraphService) GetPrivateGraphEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	vpcID := request.GetStringParam(req.Parameters, "vpcId")

	ep, err := store.GetEndpoint(graphID, vpcID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("endpoint", vpcID)
		}
		return nil, err
	}

	return endpointToResponse(ep), nil
}

// ListPrivateGraphEndpoints returns all private endpoints for the specified graph.
func (s *NeptuneGraphService) ListPrivateGraphEndpoints(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	endpoints, err := store.ListEndpoints(graphID)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(endpoints))
	for _, ep := range endpoints {
		items = append(items, endpointToResponse(ep))
	}

	return map[string]interface{}{
		"privateGraphEndpoints": items,
	}, nil
}

// DeletePrivateGraphEndpoint removes a private graph endpoint identified by graph and VPC identifiers.
func (s *NeptuneGraphService) DeletePrivateGraphEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	vpcID := request.GetStringParam(req.Parameters, "vpcId")

	ep, err := store.GetEndpoint(graphID, vpcID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("endpoint", vpcID)
		}
		return nil, err
	}

	if err := store.DeleteEndpoint(graphID, vpcID); err != nil {
		logs.Warn("failed to delete endpoint", logs.String("graphId", graphID), logs.String("vpcId", vpcID), logs.Err(err))
	}

	return endpointToResponse(ep), nil
}

// ListTagsForResource returns all tags associated with the specified resource ARN.
func (s *NeptuneGraphService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceArn := request.GetStringParam(req.Parameters, "resourceArn")
	if resourceArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "resourceArn")
	}

	tags, err := store.GetTags(resourceArn)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": tags,
	}, nil
}

// TagResource adds or updates tags on the specified resource, validating key and value constraints.
func (s *NeptuneGraphService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceArn := request.GetStringParam(req.Parameters, "resourceArn")
	if resourceArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "resourceArn")
	}

	tags := parseTagsFromParams(req.Parameters)
	if len(tags) > maxTags {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "too many tags")
	}
	for k, v := range tags {
		if len(k) < 1 || len(k) > 128 || strings.HasPrefix(k, "aws:") || !tagKeyRegex.MatchString(k) {
			return nil, newValidationException("ILLEGAL_ARGUMENT", "invalid tag key")
		}
		if len(v) > 256 {
			return nil, newValidationException("ILLEGAL_ARGUMENT", "tag value too long")
		}
	}

	if err := store.AddTags(resourceArn, tags); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": tags,
	}, nil
}

// UntagResource removes the specified tag keys from a resource's tag set.
func (s *NeptuneGraphService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceArn := request.GetStringParam(req.Parameters, "resourceArn")
	if resourceArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "resourceArn")
	}

	tagKeys := tags.ParseTagKeysAsSlice(req.Parameters, "tagKeys")

	if err := store.RemoveTags(resourceArn, tagKeys); err != nil {
		return nil, err
	}

	tags, _ := store.GetTags(resourceArn)
	return map[string]interface{}{
		"tags": tags,
	}, nil
}

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

// CancelQuery cancels a running query. Not yet implemented.
func (s *NeptuneGraphService) CancelQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return nil, newHTTPError(501, "NotImplementedException", "CancelQuery is not yet implemented")
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

// CreateGraphUsingImportTask creates a new graph and initiates a bulk import task from the specified source.
func (s *NeptuneGraphService) CreateGraphUsingImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphName := request.GetStringParam(req.Parameters, "graphName")
	if graphName == "" || strings.HasPrefix(graphName, "g-") || !graphNameRegex.MatchString(graphName) || len(graphName) > 63 {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphName")
	}

	roleArn := request.GetStringParam(req.Parameters, "roleArn")
	if roleArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "roleArn")
	}

	source := request.GetStringParam(req.Parameters, "source")
	if source == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "source")
	}

	region := reqCtx.GetRegion()
	graphID := generateID("g-")
	taskID := generateID("t-")
	now := time.Now().UTC()

	graph := &ngstore.Graph{
		Id:                 graphID,
		Name:               graphName,
		Arn:                s.arnBuilder.NeptuneGraph().Graph(graphID),
		Status:             "IMPORTING",
		ProvisionedMemory:  int32Ptr(128),
		ReplicaCount:       int32Ptr(1),
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

	graphDir := filepath.Join(s.dataPath, graphDataDirPrefix, graphID)
	db, err := graphengine.Open(graphDir, s.engineOptions())
	if err != nil {
		logs.Warn("failed to open graph engine", logs.String("graphId", graphID), logs.Err(err))
	} else {
		s.enginesMu.Lock()
		s.activeEngines[graphID] = &engineEntry{db: db}
		s.enginesMu.Unlock()
	}

	task := &ngstore.ImportTask{
		TaskId:             taskID,
		GraphId:            graphID,
		Status:             "INITIALIZING",
		Source:             source,
		RoleArn:            roleArn,
		GraphName:          graphName,
		DeletionProtection: request.GetBoolParam(req.Parameters, "deletionProtection"),
		KmsKeyIdentifier:   request.GetStringParam(req.Parameters, "kmsKeyIdentifier"),
		PublicConnectivity: request.GetBoolParam(req.Parameters, "publicConnectivity"),
		Format:             request.GetStringParam(req.Parameters, "format"),
		ParquetType:        request.GetStringParam(req.Parameters, "parquetType"),
		BlankNodeHandling:  request.GetStringParam(req.Parameters, "blankNodeHandling"),
		FailOnError:        request.GetBoolParam(req.Parameters, "failOnError"),
		StartTime:          &now,
	}

	if request.HasParam(req.Parameters, "replicaCount") {
		rc := request.GetIntParam(req.Parameters, "replicaCount")
		task.ReplicaCount = int32Ptr(int32(rc))
	}
	if request.HasParam(req.Parameters, "minProvisionedMemory") {
		task.MinProvisionedMemory = int32Ptr(int32(request.GetIntParam(req.Parameters, "minProvisionedMemory")))
	}
	if request.HasParam(req.Parameters, "maxProvisionedMemory") {
		task.MaxProvisionedMemory = int32Ptr(int32(request.GetIntParam(req.Parameters, "maxProvisionedMemory")))
	}
	if request.HasParam(req.Parameters, "importOptions") {
		task.ImportOptions = parseImportOptions(req.Parameters)
	}
	if graph.VectorSearchConfiguration != nil {
		task.VectorSearchConfiguration = graph.VectorSearchConfiguration
	}

	if err := store.CreateImportTask(task); err != nil {
		return nil, err
	}

	s.taskWg.Add(1)
	go s.advanceImportTask(store, taskID, graphID)

	return importTaskToResponse(task), nil
}

// GetImportTask retrieves an import task by its identifier.
func (s *NeptuneGraphService) GetImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	taskID := request.GetStringParam(req.Parameters, "taskIdentifier")
	if taskID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "taskIdentifier")
	}

	task, err := store.GetImportTask(taskID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("import task", taskID)
		}
		return nil, err
	}

	return importTaskToResponse(task), nil
}

// ListImportTasks returns a paginated list of import task summaries.
func (s *NeptuneGraphService) ListImportTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := storecommon.ListOptions{
		MaxItems: clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		Marker:   request.GetStringParam(req.Parameters, "nextToken"),
	}

	tasks, nextToken, truncated, err := store.ListImportTasks(opts)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, importTaskSummaryToResponse(t))
	}

	result := map[string]interface{}{
		"tasks": items,
	}
	if truncated {
		result["nextToken"] = nextToken
	}
	return result, nil
}

// CancelImportTask cancels an in-progress import task, transitioning it to CANCELLED state.
func (s *NeptuneGraphService) CancelImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	taskID := request.GetStringParam(req.Parameters, "taskIdentifier")
	if taskID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "taskIdentifier")
	}

	task, err := store.GetImportTask(taskID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("import task", taskID)
		}
		return nil, err
	}

	if task.Status == "SUCCEEDED" || task.Status == "FAILED" || task.Status == "CANCELLED" {
		return importTaskSummaryToResponse(task), nil
	}

	task.Status = "CANCELLING"
	if err := store.UpdateImportTask(task); err != nil {
		logs.Warn("failed to update import task status to CANCELLING", logs.String("taskId", taskID), logs.Err(err))
	}

	task.Status = "CANCELLED"
	task.StatusReason = "Cancelled by user"
	if err := store.UpdateImportTask(task); err != nil {
		logs.Warn("failed to update import task status to CANCELLED", logs.String("taskId", taskID), logs.Err(err))
	}

	return importTaskSummaryToResponse(task), nil
}

// StartImportTask initiates a bulk import task on an existing graph from the specified source.
func (s *NeptuneGraphService) StartImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	_, err = store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	roleArn := request.GetStringParam(req.Parameters, "roleArn")
	if roleArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "roleArn")
	}

	source := request.GetStringParam(req.Parameters, "source")
	if source == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "source")
	}

	taskID := generateID("t-")
	now := time.Now().UTC()

	task := &ngstore.ImportTask{
		TaskId:            taskID,
		GraphId:           graphID,
		Status:            "INITIALIZING",
		Source:            source,
		RoleArn:           roleArn,
		Format:            request.GetStringParam(req.Parameters, "format"),
		ParquetType:       request.GetStringParam(req.Parameters, "parquetType"),
		BlankNodeHandling: request.GetStringParam(req.Parameters, "blankNodeHandling"),
		FailOnError:       request.GetBoolParam(req.Parameters, "failOnError"),
		StartTime:         &now,
	}

	if request.HasParam(req.Parameters, "importOptions") {
		task.ImportOptions = parseImportOptions(req.Parameters)
	}

	if err := store.CreateImportTask(task); err != nil {
		return nil, err
	}

	s.taskWg.Add(1)
	go s.advanceImportTask(store, taskID, graphID)

	return importTaskToResponse(task), nil
}

// StartExportTask initiates a bulk export of graph data to the specified S3 destination.
func (s *NeptuneGraphService) StartExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	_, err = store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	destination := request.GetStringParam(req.Parameters, "destination")
	if destination == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "destination")
	}

	format := request.GetStringParam(req.Parameters, "format")
	if format == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "format")
	}

	kmsKeyID := request.GetStringParam(req.Parameters, "kmsKeyIdentifier")
	if kmsKeyID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "kmsKeyIdentifier")
	}

	roleArn := request.GetStringParam(req.Parameters, "roleArn")
	if roleArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "roleArn")
	}

	taskID := generateID("t-")
	now := time.Now().UTC()

	task := &ngstore.ExportTask{
		TaskId:           taskID,
		GraphId:          graphID,
		Status:           "INITIALIZING",
		Format:           format,
		ParquetType:      request.GetStringParam(req.Parameters, "parquetType"),
		Destination:      destination,
		RoleArn:          roleArn,
		KmsKeyIdentifier: kmsKeyID,
		StartTime:        &now,
	}

	if request.HasParam(req.Parameters, "exportFilter") {
		task.ExportFilter = parseExportFilter(req.Parameters)
	}

	if err := store.CreateExportTask(task); err != nil {
		return nil, err
	}

	s.taskWg.Add(1)
	go s.advanceExportTask(store, taskID)

	return exportTaskToResponse(task), nil
}

// GetExportTask retrieves an export task by its identifier.
func (s *NeptuneGraphService) GetExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	taskID := request.GetStringParam(req.Parameters, "taskIdentifier")
	if taskID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "taskIdentifier")
	}

	task, err := store.GetExportTask(taskID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("export task", taskID)
		}
		return nil, err
	}

	return exportTaskToResponse(task), nil
}

// ListExportTasks returns a paginated list of export task summaries, optionally filtered by graph.
func (s *NeptuneGraphService) ListExportTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := storecommon.ListOptions{
		MaxItems: clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		Marker:   request.GetStringParam(req.Parameters, "nextToken"),
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	tasks, nextToken, truncated, err := store.ListExportTasks(opts, graphID)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, exportTaskSummaryToResponse(t))
	}

	result := map[string]interface{}{
		"tasks": items,
	}
	if truncated {
		result["nextToken"] = nextToken
	}
	return result, nil
}

// CancelExportTask cancels an in-progress export task, transitioning it to CANCELLED state.
func (s *NeptuneGraphService) CancelExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	taskID := request.GetStringParam(req.Parameters, "taskIdentifier")
	if taskID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "taskIdentifier")
	}

	task, err := store.GetExportTask(taskID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("export task", taskID)
		}
		return nil, err
	}

	if task.Status == "SUCCEEDED" || task.Status == "FAILED" || task.Status == "CANCELLED" {
		return exportTaskSummaryToResponse(task), nil
	}

	task.Status = "CANCELLING"
	if err := store.UpdateExportTask(task); err != nil {
		logs.Warn("failed to update export task status to CANCELLING", logs.String("taskId", taskID), logs.Err(err))
	}

	task.Status = "CANCELLED"
	task.StatusReason = "Cancelled by user"
	if err := store.UpdateExportTask(task); err != nil {
		logs.Warn("failed to update export task status to CANCELLED", logs.String("taskId", taskID), logs.Err(err))
	}

	return exportTaskSummaryToResponse(task), nil
}

func (s *NeptuneGraphService) advanceImportTask(store *ngstore.NeptuneGraphStore, taskID, graphID string) {
	defer s.taskWg.Done()
	defer func() { resilience.RecoverPanic("NeptuneGraph advanceImportTask") }()

	task, err := store.GetImportTask(taskID)
	if err != nil {
		logs.Warn("failed to get import task", logs.String("taskId", taskID), logs.Err(err))
		return
	}

	source := task.Source
	format := strings.ToLower(task.Format)

	if strings.HasPrefix(strings.ToLower(source), "s3://") {
		err = store.TryAdvanceImportTask(taskID, "INITIALIZING", func(t *ngstore.ImportTask) {
			t.Status = "FAILED"
			t.StatusReason = "S3 sources are not accessible in standalone mode"
			now := time.Now().UTC()
			sinceStart := int64(now.Sub(*t.StartTime).Seconds())
			t.ImportTaskDetails = &ngstore.ImportTaskDetails{
				ProgressPercentage: int32Ptr(0),
				StartTime:          t.StartTime,
				TimeElapsedSeconds: int64Ptr(sinceStart),
				StatementCount:     int64Ptr(0),
				ErrorCount:         int32Ptr(1),
				Status:             stringPtr("FAILED"),
			}
		})
		if err != nil {
			logs.Warn("failed to advance import task to FAILED", logs.String("taskId", taskID), logs.Err(err))
		}
		return
	}

	filePath := strings.TrimPrefix(source, "file://")
	if filePath == source {
		filePath = source
	}

	s.enginesMu.RLock()
	entry, ok := s.activeEngines[graphID]
	s.enginesMu.RUnlock()

	if !ok {
		logs.Warn("engine not found for graph during import", logs.String("graphId", graphID))
		err = store.TryAdvanceImportTask(taskID, "INITIALIZING", func(t *ngstore.ImportTask) {
			t.Status = "FAILED"
			t.StatusReason = "Graph engine not available"
			now := time.Now().UTC()
			sinceStart := int64(now.Sub(*t.StartTime).Seconds())
			t.ImportTaskDetails = &ngstore.ImportTaskDetails{
				ProgressPercentage: int32Ptr(0),
				StartTime:          t.StartTime,
				TimeElapsedSeconds: int64Ptr(sinceStart),
				StatementCount:     int64Ptr(0),
				ErrorCount:         int32Ptr(1),
				Status:             stringPtr("FAILED"),
			}
		})
		if err != nil {
			logs.Warn("failed to advance import task to FAILED", logs.String("taskId", taskID), logs.Err(err))
		}
		return
	}

	err = store.TryAdvanceImportTask(taskID, "INITIALIZING", func(t *ngstore.ImportTask) {
		t.Status = "IN_PROGRESS"
	})
	if err != nil {
		logs.Warn("failed to advance import task to IN_PROGRESS", logs.String("taskId", taskID), logs.Err(err))
		return
	}

	var statementCount, dictionaryCount, errorCount int64

	switch {
	case format == "" || format == "csv":
		statementCount, dictionaryCount, errorCount = s.importCSV(entry.db, filePath)
	case format == "json":
		statementCount, dictionaryCount, errorCount = s.importJSON(entry.db, filePath)
	case format == "ntriples" || format == "nquad":
		statementCount, dictionaryCount, errorCount = s.importRDF(entry.db, filePath)
	default:
		errorCount = 1
		logs.Warn("unsupported import format", logs.String("format", format), logs.String("taskId", taskID))
	}

	finalStatus := "SUCCEEDED"
	statusReason := ""
	if errorCount > 0 && task.FailOnError {
		finalStatus = "FAILED"
		statusReason = fmt.Sprintf("%d errors during import", errorCount)
	}

	now := time.Now().UTC()
	sinceStart := int64(now.Sub(*task.StartTime).Seconds())
	details := &ngstore.ImportTaskDetails{
		ProgressPercentage:   int32Ptr(100),
		StartTime:            task.StartTime,
		TimeElapsedSeconds:   int64Ptr(sinceStart),
		StatementCount:       int64Ptr(statementCount),
		DictionaryEntryCount: int64Ptr(dictionaryCount),
		ErrorCount:           int32Ptr(int32(errorCount)),
		Status:               stringPtr(finalStatus),
	}

	err = store.TryAdvanceImportTask(taskID, "IN_PROGRESS", func(t *ngstore.ImportTask) {
		t.Status = finalStatus
		t.StatusReason = statusReason
		t.ImportTaskDetails = details
	})
	if err != nil {
		logs.Warn("failed to advance import task to final state", logs.String("taskId", taskID), logs.Err(err))
		return
	}

	if finalStatus == "SUCCEEDED" {
		graph, err := store.GetGraph(graphID)
		if err == nil && graph.Status == "IMPORTING" {
			graph.Status = "AVAILABLE"
			if err := store.UpdateGraph(graph); err != nil {
				logs.Warn("Failed to update graph status to AVAILABLE", logs.Err(err))
			}
		}
	}
}

func (s *NeptuneGraphService) importCSV(db *graphengine.DB, filePath string) (statementCount, dictionaryCount, errorCount int64) {
	f, err := os.Open(filePath)
	if err != nil {
		errorCount++
		logs.Warn("failed to open import file", logs.String("path", filePath), logs.Err(err))
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.LazyQuotes = true

	header, err := r.Read()
	if err != nil {
		errorCount++
		logs.Warn("failed to read CSV header", logs.Err(err))
		return
	}

	headerLower := make([]string, len(header))
	for i, h := range header {
		headerLower[i] = strings.ToLower(strings.TrimSpace(h))
	}

	hasFromTo := false
	for _, h := range headerLower {
		if h == "~from" || h == "~to" {
			hasFromTo = true
			break
		}
	}

	isEdge := hasFromTo

	idIndex := -1
	labelIndex := -1
	fromIndex := -1
	toIndex := -1
	for i, h := range headerLower {
		switch h {
		case "~id":
			idIndex = i
		case "~label":
			labelIndex = i
		case "~from":
			fromIndex = i
		case "~to":
			toIndex = i
		}
	}

	propIndices := make(map[string]int)
	for i, h := range headerLower {
		if strings.HasPrefix(h, "~") {
			continue
		}
		propIndices[h] = i
	}

	var nodeBatch []struct {
		Labels []string
		Props  graphengine.Props
	}
	var edgeBatch []graphengine.Edge
	extToInternal := make(map[string]graphengine.NodeID)

	flushNodeBatch := func() {
		if len(nodeBatch) == 0 {
			return
		}
		ids, err := db.AddNodeBatch(nodeBatch)
		if err != nil {
			errorCount += int64(len(nodeBatch))
			logs.Warn("AddNodeBatch failed", logs.Err(err))
		} else {
			for i, id := range ids {
				if idIndex >= 0 && i < len(nodeBatch) {
					if extID, ok := nodeBatch[i].Props["~id"].(string); ok {
						extToInternal[extID] = id
						delete(nodeBatch[i].Props, "~id")
					}
				}
			}
			dictionaryCount += int64(len(ids))
		}
		nodeBatch = nil
	}

	for {
		record, err := r.Read()
		if err != nil {
			if err.Error() != "EOF" {
				errorCount++
			}
			break
		}

		statementCount++

		if isEdge {
			flushNodeBatch()

			if fromIndex < 0 || toIndex < 0 {
				errorCount++
				continue
			}

			fromExt := strings.TrimSpace(record[fromIndex])
			toExt := strings.TrimSpace(record[toIndex])
			label := ""
			if labelIndex >= 0 && labelIndex < len(record) {
				label = strings.TrimSpace(record[labelIndex])
			}

			fromID, ok := extToInternal[fromExt]
			if !ok {
				errorCount++
				continue
			}
			toID, ok := extToInternal[toExt]
			if !ok {
				errorCount++
				continue
			}

			props := make(graphengine.Props)
			for k, idx := range propIndices {
				if idx < len(record) && record[idx] != "" {
					props[k] = strings.TrimSpace(record[idx])
				}
			}

			edgeBatch = append(edgeBatch, graphengine.Edge{
				From:  fromID,
				To:    toID,
				Label: label,
				Props: props,
			})
		} else {
			var labels []string
			if labelIndex >= 0 && labelIndex < len(record) {
				for _, l := range strings.Split(strings.TrimSpace(record[labelIndex]), ":") {
					if l != "" {
						labels = append(labels, strings.TrimSpace(l))
					}
				}
			}

			props := make(graphengine.Props)
			if idIndex >= 0 && idIndex < len(record) {
				props["~id"] = strings.TrimSpace(record[idIndex])
			}
			for k, idx := range propIndices {
				if idx < len(record) && record[idx] != "" {
					props[k] = strings.TrimSpace(record[idx])
				}
			}

			nodeBatch = append(nodeBatch, struct {
				Labels []string
				Props  graphengine.Props
			}{Labels: labels, Props: props})
		}

		if len(nodeBatch) >= 500 {
			flushNodeBatch()
		}

		if len(edgeBatch) >= 500 {
			_, err := db.AddEdgeBatch(edgeBatch)
			if err != nil {
				errorCount += int64(len(edgeBatch))
				logs.Warn("AddEdgeBatch failed", logs.Err(err))
			} else {
				dictionaryCount += int64(len(edgeBatch))
			}
			edgeBatch = nil
		}
	}

	if len(nodeBatch) > 0 {
		flushNodeBatch()
	}

	if len(edgeBatch) > 0 {
		_, err := db.AddEdgeBatch(edgeBatch)
		if err != nil {
			errorCount += int64(len(edgeBatch))
			logs.Warn("AddEdgeBatch failed", logs.Err(err))
		} else {
			dictionaryCount += int64(len(edgeBatch))
		}
	}

	return
}

func (s *NeptuneGraphService) importJSON(db *graphengine.DB, filePath string) (statementCount, dictionaryCount, errorCount int64) {
	f, err := os.Open(filePath)
	if err != nil {
		errorCount++
		logs.Warn("failed to open import file", logs.String("path", filePath), logs.Err(err))
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	var nodeBatch []struct {
		Labels []string
		Props  graphengine.Props
	}
	var edgeBatch []graphengine.Edge
	extToInternal := make(map[string]graphengine.NodeID)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var doc map[string]interface{}
		if err := json.Unmarshal([]byte(line), &doc); err != nil {
			errorCount++
			continue
		}

		statementCount++

		if fromVal, hasFrom := doc["~from"]; hasFrom {
			fromExt := fmt.Sprintf("%v", fromVal)
			toExt := fmt.Sprintf("%v", doc["~to"])
			label := fmt.Sprintf("%v", doc["~type"])

			fromID, ok := extToInternal[fromExt]
			if !ok {
				errorCount++
				continue
			}
			toID, ok := extToInternal[toExt]
			if !ok {
				errorCount++
				continue
			}

			props := make(graphengine.Props)
			for k, v := range doc {
				if strings.HasPrefix(k, "~") {
					continue
				}
				props[k] = v
			}

			edgeBatch = append(edgeBatch, graphengine.Edge{
				From:  fromID,
				To:    toID,
				Label: label,
				Props: props,
			})
		} else {
			var labels []string
			if lbl, ok := doc["~label"]; ok {
				switch l := lbl.(type) {
				case string:
					for _, part := range strings.Split(l, ":") {
						if part != "" {
							labels = append(labels, part)
						}
					}
				case []interface{}:
					for _, item := range l {
						if s, ok := item.(string); ok {
							labels = append(labels, s)
						}
					}
				}
			}

			props := make(graphengine.Props)
			for k, v := range doc {
				if strings.HasPrefix(k, "~") {
					continue
				}
				props[k] = v
			}

			nodeBatch = append(nodeBatch, struct {
				Labels []string
				Props  graphengine.Props
			}{Labels: labels, Props: props})
		}

		if len(nodeBatch) >= 500 {
			ids, err := db.AddNodeBatch(nodeBatch)
			if err != nil {
				errorCount += int64(len(nodeBatch))
				logs.Warn("AddNodeBatch failed", logs.Err(err))
			} else {
				for i, id := range ids {
					if i < len(nodeBatch) {
						if extID, ok := nodeBatch[i].Props["~id"].(string); ok {
							extToInternal[extID] = id
							delete(nodeBatch[i].Props, "~id")
						}
					}
				}
				dictionaryCount += int64(len(ids))
			}
			nodeBatch = nil
		}

		if len(edgeBatch) >= 500 {
			_, err := db.AddEdgeBatch(edgeBatch)
			if err != nil {
				errorCount += int64(len(edgeBatch))
				logs.Warn("AddEdgeBatch failed", logs.Err(err))
			} else {
				dictionaryCount += int64(len(edgeBatch))
			}
			edgeBatch = nil
		}
	}

	if len(nodeBatch) > 0 {
		ids, err := db.AddNodeBatch(nodeBatch)
		if err != nil {
			errorCount += int64(len(nodeBatch))
			logs.Warn("AddNodeBatch failed", logs.Err(err))
		} else {
			for i, id := range ids {
				if i < len(nodeBatch) {
					if extID, ok := nodeBatch[i].Props["~id"].(string); ok {
						extToInternal[extID] = id
						delete(nodeBatch[i].Props, "~id")
					}
				}
			}
			dictionaryCount += int64(len(ids))
		}
	}

	if len(edgeBatch) > 0 {
		_, err := db.AddEdgeBatch(edgeBatch)
		if err != nil {
			errorCount += int64(len(edgeBatch))
			logs.Warn("AddEdgeBatch failed", logs.Err(err))
		} else {
			dictionaryCount += int64(len(edgeBatch))
		}
	}

	return
}

func (s *NeptuneGraphService) importRDF(db *graphengine.DB, filePath string) (statementCount, dictionaryCount, errorCount int64) {
	f, err := os.Open(filePath)
	if err != nil {
		errorCount++
		logs.Warn("failed to open import file", logs.String("path", filePath), logs.Err(err))
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	extToInternal := make(map[string]graphengine.NodeID)

	var nodeBatch []struct {
		Labels []string
		Props  graphengine.Props
	}
	var edgeBatch []graphengine.Edge

	ensureNode := func(extID string) graphengine.NodeID {
		if id, ok := extToInternal[extID]; ok {
			return id
		}
		nodeID, err := db.AddNode([]string{"Resource"}, graphengine.Props{"uri": extID})
		if err != nil {
			return 0
		}
		extToInternal[extID] = nodeID
		dictionaryCount++
		return nodeID
	}

	flushNodes := func() {
		if len(nodeBatch) == 0 {
			return
		}
		ids, err := db.AddNodeBatch(nodeBatch)
		if err != nil {
			errorCount += int64(len(nodeBatch))
			logs.Warn("AddNodeBatch failed", logs.Err(err))
		} else {
			for i, id := range ids {
				if i < len(nodeBatch) {
					if uri, ok := nodeBatch[i].Props["uri"].(string); ok {
						extToInternal[uri] = id
					}
				}
			}
			dictionaryCount += int64(len(ids))
		}
		nodeBatch = nil
	}

	flushEdges := func() {
		if len(edgeBatch) == 0 {
			return
		}
		_, err := db.AddEdgeBatch(edgeBatch)
		if err != nil {
			errorCount += int64(len(edgeBatch))
			logs.Warn("AddEdgeBatch failed", logs.Err(err))
		} else {
			dictionaryCount += int64(len(edgeBatch))
		}
		edgeBatch = nil
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		subject, predicate, obj, ok := parseNTriplesLine(line)
		if !ok {
			errorCount++
			continue
		}

		statementCount++

		subjID := ensureNode(subject)

		isURI := strings.HasPrefix(obj, "<")
		if isURI {
			objURI := strings.Trim(obj, "<>")
			objID := ensureNode(objURI)

			predLabel := extractLocalName(predicate)

			edgeBatch = append(edgeBatch, graphengine.Edge{
				From:  subjID,
				To:    objID,
				Label: predLabel,
			})
		} else {
			predKey := extractLocalName(predicate)

			node, err := db.GetNode(subjID)
			if err == nil && node != nil {
				updatedProps := make(graphengine.Props)
				for k, v := range node.Props {
					updatedProps[k] = v
				}
				if strings.HasPrefix(obj, "\"") {
					closing := strings.Index(obj[1:], "\"")
					if closing >= 0 {
						obj = obj[1 : closing+1]
					}
				}
				updatedProps[predKey] = obj
				_ = db.UpdateNode(subjID, updatedProps)
			}
		}

		if len(edgeBatch) >= 500 {
			flushEdges()
		}
	}

	flushNodes()
	flushEdges()
	return
}

func (s *NeptuneGraphService) advanceExportTask(store *ngstore.NeptuneGraphStore, taskID string) {
	defer s.taskWg.Done()
	defer func() { resilience.RecoverPanic("NeptuneGraph advanceExportTask") }()

	task, err := store.GetExportTask(taskID)
	if err != nil {
		logs.Warn("failed to get export task", logs.String("taskId", taskID), logs.Err(err))
		return
	}

	err = store.TryAdvanceExportTask(taskID, "INITIALIZING", func(t *ngstore.ExportTask) {
		t.Status = "EXPORTING"
	})
	if err != nil {
		logs.Warn("failed to advance export task to EXPORTING", logs.String("taskId", taskID), logs.Err(err))
		return
	}

	s.enginesMu.RLock()
	entry, ok := s.activeEngines[task.GraphId]
	s.enginesMu.RUnlock()
	if !ok {
		failExportTask(store, taskID, task, "graph engine not found for graphId: "+task.GraphId)
		return
	}

	dest := task.Destination
	format := strings.ToUpper(task.Format)

	if strings.HasPrefix(dest, "s3://") {
		failExportTask(store, taskID, task, "S3 destination is not accessible in standalone mode")
		return
	}

	if !strings.HasPrefix(dest, "file://") {
		failExportTask(store, taskID, task, "unsupported destination scheme: "+dest)
		return
	}

	filePath := strings.TrimPrefix(dest, "file://")
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		failExportTask(store, taskID, task, fmt.Sprintf("failed to create export directory: %v", err))
		return
	}

	var nodeCount, edgeCount int64
	if format == "CSV" || format == "CSV+BINARY" {
		nodeCount, edgeCount, err = exportGraphCSV(entry.db, filePath)
	} else {
		failExportTask(store, taskID, task, "unsupported export format: "+task.Format)
		return
	}

	if err != nil {
		failExportTask(store, taskID, task, fmt.Sprintf("export failed: %v", err))
		return
	}

	task, err = store.GetExportTask(taskID)
	if err != nil || task.Status == "CANCELLED" {
		return
	}

	err = store.TryAdvanceExportTask(taskID, "EXPORTING", func(t *ngstore.ExportTask) {
		t.Status = "SUCCEEDED"
		now := time.Now().UTC()
		sinceStart := int64(now.Sub(*t.StartTime).Seconds())
		details := t.ExportTaskDetails
		if details == nil {
			details = &ngstore.ExportTaskDetails{}
		}
		details.ProgressPercentage = int32Ptr(100)
		details.StartTime = t.StartTime
		details.TimeElapsedSeconds = int64Ptr(sinceStart)
		details.NumVerticesWritten = int64Ptr(nodeCount)
		details.NumEdgesWritten = int64Ptr(edgeCount)
		t.ExportTaskDetails = details
	})
	if err != nil {
		logs.Warn("failed to advance export task to SUCCEEDED", logs.String("taskId", taskID), logs.Err(err))
	}
}

func failExportTask(store *ngstore.NeptuneGraphStore, taskID string, task *ngstore.ExportTask, reason string) {
	logs.Warn("export task failed", logs.String("taskId", taskID), logs.String("reason", reason))
	now := time.Now().UTC()
	sinceStart := int64(0)
	if task.StartTime != nil {
		sinceStart = int64(now.Sub(*task.StartTime).Seconds())
	}
	store.TryAdvanceExportTask(taskID, "EXPORTING", func(t *ngstore.ExportTask) {
		t.Status = "FAILED"
		t.StatusReason = reason
		t.ExportTaskDetails = &ngstore.ExportTaskDetails{
			ProgressPercentage: int32Ptr(0),
			StartTime:          t.StartTime,
			TimeElapsedSeconds: int64Ptr(sinceStart),
		}
	})
}

func exportGraphCSV(db *graphengine.DB, filePath string) (int64, int64, error) {
	nodesFile := filePath
	if !strings.HasSuffix(strings.ToLower(filePath), ".csv") {
		nodesFile = filePath + "_nodes.csv"
	}
	edgesFile := strings.TrimSuffix(nodesFile, "_nodes.csv") + "_edges.csv"
	if nodesFile == filePath {
		edgesFile = filePath + "_edges.csv"
	}

	var nodeCount, edgeCount int64

	nf, err := os.Create(nodesFile)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create nodes file: %w", err)
	}
	defer nf.Close()

	nodeW := csv.NewWriter(nf)
	defer nodeW.Flush()

	ef, err := os.Create(edgesFile)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create edges file: %w", err)
	}
	defer ef.Close()

	edgeW := csv.NewWriter(ef)
	defer edgeW.Flush()

	var propKeys []string

	err = db.ForEachNode(func(node *graphengine.Node) error {
		propKeys = sortedPropKeys(node.Props)
		record := make([]string, 0, 2+len(propKeys))
		record = append(record, fmt.Sprintf("%d", node.ID))
		record = append(record, strings.Join(node.Labels, ";"))
		for _, k := range propKeys {
			record = append(record, formatPropValue(node.Props[k]))
		}
		if err := nodeW.Write(record); err != nil {
			return err
		}
		nodeCount++
		return nil
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to enumerate nodes: %w", err)
	}

	err = db.ForEachEdge(func(edge *graphengine.Edge) error {
		propKeys = sortedPropKeys(edge.Props)
		record := make([]string, 0, 4+len(propKeys))
		record = append(record, fmt.Sprintf("%d", edge.ID))
		record = append(record, edge.Label)
		record = append(record, fmt.Sprintf("%d", edge.From))
		record = append(record, fmt.Sprintf("%d", edge.To))
		for _, k := range propKeys {
			record = append(record, formatPropValue(edge.Props[k]))
		}
		if err := edgeW.Write(record); err != nil {
			return err
		}
		edgeCount++
		return nil
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to enumerate edges: %w", err)
	}

	return nodeCount, edgeCount, nil
}

func sortedPropKeys(props graphengine.Props) []string {
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}

func formatPropValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(b)
	}
}

func graphToResponse(g *ngstore.Graph, isCreate bool) map[string]interface{} {
	r := map[string]interface{}{
		"id":                 g.Id,
		"name":               g.Name,
		"arn":                g.Arn,
		"status":             g.Status,
		"endpoint":           g.Endpoint,
		"deletionProtection": g.DeletionProtection,
		"publicConnectivity": g.PublicConnectivity,
		"buildNumber":        g.BuildNumber,
		"sourceSnapshotId":   g.SourceSnapshotId,
		"accountId":          g.AccountID,
		"region":             g.Region,
	}
	if g.ProvisionedMemory != nil {
		r["provisionedMemory"] = *g.ProvisionedMemory
	}
	if g.ReplicaCount != nil {
		r["replicaCount"] = *g.ReplicaCount
	}
	if g.VectorSearchConfiguration != nil {
		r["vectorSearchConfiguration"] = map[string]interface{}{
			"dimension": g.VectorSearchConfiguration.Dimension,
		}
	}
	if g.KmsKeyIdentifier != "" {
		r["kmsKeyIdentifier"] = g.KmsKeyIdentifier
	}
	if g.StatusReason != "" {
		r["statusReason"] = g.StatusReason
	}
	if g.CreateTime != nil {
		r["createTime"] = g.CreateTime.Unix()
	}
	return r
}

func graphSummaryToResponse(g *ngstore.Graph) map[string]interface{} {
	r := map[string]interface{}{
		"id":                 g.Id,
		"name":               g.Name,
		"arn":                g.Arn,
		"status":             g.Status,
		"endpoint":           g.Endpoint,
		"deletionProtection": g.DeletionProtection,
		"publicConnectivity": g.PublicConnectivity,
	}
	if g.ProvisionedMemory != nil {
		r["provisionedMemory"] = *g.ProvisionedMemory
	}
	if g.ReplicaCount != nil {
		r["replicaCount"] = *g.ReplicaCount
	}
	if g.KmsKeyIdentifier != "" {
		r["kmsKeyIdentifier"] = g.KmsKeyIdentifier
	}
	return r
}

func snapshotToResponse(s *ngstore.GraphSnapshot) map[string]interface{} {
	r := map[string]interface{}{
		"id":            s.Id,
		"name":          s.Name,
		"arn":           s.Arn,
		"status":        s.Status,
		"sourceGraphId": s.SourceGraphId,
		"accountId":     s.AccountID,
		"region":        s.Region,
	}
	if s.KmsKeyIdentifier != "" {
		r["kmsKeyIdentifier"] = s.KmsKeyIdentifier
	}
	if s.SnapshotCreateTime != nil {
		r["snapshotCreateTime"] = s.SnapshotCreateTime.Unix()
	}
	return r
}

func endpointToResponse(ep *ngstore.PrivateGraphEndpoint) map[string]interface{} {
	r := map[string]interface{}{
		"status":    ep.Status,
		"subnetIds": ep.SubnetIds,
		"accountId": ep.AccountID,
		"region":    ep.Region,
	}
	if ep.VpcId != "" {
		r["vpcId"] = ep.VpcId
	}
	if ep.VpcEndpointId != "" {
		r["vpcEndpointId"] = ep.VpcEndpointId
	}
	return r
}

func queryToResponse(q *ngstore.QueryRecord) map[string]interface{} {
	r := map[string]interface{}{
		"id":          q.Id,
		"queryString": q.QueryString,
		"state":       q.State,
	}
	if q.Elapsed > 0 {
		r["elapsed"] = q.Elapsed
	}
	if q.Waited > 0 {
		r["waited"] = q.Waited
	}
	return r
}

func graphDataSummaryToResponse(s *ngstore.GraphDataSummary) map[string]interface{} {
	r := map[string]interface{}{
		"numNodes":      s.NumNodes,
		"numEdges":      s.NumEdges,
		"numNodeLabels": s.NumNodeLabels,
		"numEdgeLabels": s.NumEdgeLabels,
		"nodeLabels":    s.NodeLabels,
		"edgeLabels":    s.EdgeLabels,
	}
	if s.NumNodeProperties != nil {
		r["numNodeProperties"] = *s.NumNodeProperties
	}
	if s.NumEdgeProperties != nil {
		r["numEdgeProperties"] = *s.NumEdgeProperties
	}
	if s.TotalNodePropertyValues != nil {
		r["totalNodePropertyValues"] = *s.TotalNodePropertyValues
	}
	if s.TotalEdgePropertyValues != nil {
		r["totalEdgePropertyValues"] = *s.TotalEdgePropertyValues
	}
	if len(s.NodeStructures) > 0 {
		ns := make([]interface{}, 0, len(s.NodeStructures))
		for _, n := range s.NodeStructures {
			ns = append(ns, nodeStructureToResponse(&n))
		}
		r["nodeStructures"] = ns
	}
	if len(s.EdgeStructures) > 0 {
		es := make([]interface{}, 0, len(s.EdgeStructures))
		for _, e := range s.EdgeStructures {
			es = append(es, edgeStructureToResponse(&e))
		}
		r["edgeStructures"] = es
	}
	return r
}

func nodeStructureToResponse(n *ngstore.NodeStructure) map[string]interface{} {
	r := map[string]interface{}{
		"nodeProperties": n.NodeProperties,
	}
	if n.Count != nil {
		r["count"] = *n.Count
	}
	if len(n.DistinctOutgoingEdgeLabels) > 0 {
		r["distinctOutgoingEdgeLabels"] = n.DistinctOutgoingEdgeLabels
	}
	return r
}

func edgeStructureToResponse(e *ngstore.EdgeStructure) map[string]interface{} {
	r := map[string]interface{}{
		"edgeProperties": e.EdgeProperties,
	}
	if e.Count != nil {
		r["count"] = *e.Count
	}
	return r
}

func importTaskToResponse(t *ngstore.ImportTask) map[string]interface{} {
	r := map[string]interface{}{
		"taskId":  t.TaskId,
		"graphId": t.GraphId,
		"status":  t.Status,
		"source":  t.Source,
		"roleArn": t.RoleArn,
	}
	if t.Format != "" {
		r["format"] = t.Format
	}
	if t.ParquetType != "" {
		r["parquetType"] = t.ParquetType
	}
	if t.StatusReason != "" {
		r["statusReason"] = t.StatusReason
	}
	if t.ImportTaskDetails != nil {
		r["importTaskDetails"] = importTaskDetailsToResponse(t.ImportTaskDetails)
	}
	if t.ImportOptions != nil {
		r["importOptions"] = importOptionsToResponse(t.ImportOptions)
	}
	return r
}

func importTaskSummaryToResponse(t *ngstore.ImportTask) map[string]interface{} {
	r := map[string]interface{}{
		"taskId":  t.TaskId,
		"graphId": t.GraphId,
		"status":  t.Status,
		"source":  t.Source,
		"roleArn": t.RoleArn,
	}
	if t.Format != "" {
		r["format"] = t.Format
	}
	if t.ParquetType != "" {
		r["parquetType"] = t.ParquetType
	}
	return r
}

func importTaskDetailsToResponse(d *ngstore.ImportTaskDetails) map[string]interface{} {
	r := map[string]interface{}{}
	if d.ProgressPercentage != nil {
		r["progressPercentage"] = *d.ProgressPercentage
	}
	if d.StartTime != nil {
		r["startTime"] = d.StartTime.Unix()
	}
	if d.TimeElapsedSeconds != nil {
		r["timeElapsedSeconds"] = *d.TimeElapsedSeconds
	}
	if d.StatementCount != nil {
		r["statementCount"] = *d.StatementCount
	}
	if d.DictionaryEntryCount != nil {
		r["dictionaryEntryCount"] = *d.DictionaryEntryCount
	}
	if d.ErrorCount != nil {
		r["errorCount"] = *d.ErrorCount
	}
	if d.ErrorDetails != nil {
		r["errorDetails"] = *d.ErrorDetails
	}
	if d.Status != nil {
		r["status"] = *d.Status
	}
	return r
}

func importOptionsToResponse(o *ngstore.ImportOptions) map[string]interface{} {
	if o == nil || o.Neptune == nil {
		return nil
	}
	r := map[string]interface{}{
		"s3ExportPath":     o.Neptune.S3ExportPath,
		"s3ExportKmsKeyId": o.Neptune.S3ExportKmsKeyId,
	}
	if o.Neptune.PreserveDefaultVertexLabels != nil {
		r["preserveDefaultVertexLabels"] = *o.Neptune.PreserveDefaultVertexLabels
	}
	if o.Neptune.PreserveEdgeIds != nil {
		r["preserveEdgeIds"] = *o.Neptune.PreserveEdgeIds
	}
	return map[string]interface{}{"neptune": r}
}

func exportTaskToResponse(t *ngstore.ExportTask) map[string]interface{} {
	r := map[string]interface{}{
		"taskId":           t.TaskId,
		"graphId":          t.GraphId,
		"status":           t.Status,
		"format":           t.Format,
		"destination":      t.Destination,
		"roleArn":          t.RoleArn,
		"kmsKeyIdentifier": t.KmsKeyIdentifier,
	}
	if t.ParquetType != "" {
		r["parquetType"] = t.ParquetType
	}
	if t.StatusReason != "" {
		r["statusReason"] = t.StatusReason
	}
	if t.ExportFilter != nil {
		r["exportFilter"] = exportFilterToResponse(t.ExportFilter)
	}
	if t.ExportTaskDetails != nil {
		r["exportTaskDetails"] = exportTaskDetailsToResponse(t.ExportTaskDetails)
	}
	return r
}

func exportTaskSummaryToResponse(t *ngstore.ExportTask) map[string]interface{} {
	r := map[string]interface{}{
		"taskId":           t.TaskId,
		"graphId":          t.GraphId,
		"status":           t.Status,
		"format":           t.Format,
		"destination":      t.Destination,
		"roleArn":          t.RoleArn,
		"kmsKeyIdentifier": t.KmsKeyIdentifier,
	}
	if t.ParquetType != "" {
		r["parquetType"] = t.ParquetType
	}
	if t.StatusReason != "" {
		r["statusReason"] = t.StatusReason
	}
	return r
}

func exportTaskDetailsToResponse(d *ngstore.ExportTaskDetails) map[string]interface{} {
	r := map[string]interface{}{}
	if d.ProgressPercentage != nil {
		r["progressPercentage"] = *d.ProgressPercentage
	}
	if d.StartTime != nil {
		r["startTime"] = d.StartTime.Unix()
	}
	if d.TimeElapsedSeconds != nil {
		r["timeElapsedSeconds"] = *d.TimeElapsedSeconds
	}
	if d.NumEdgesWritten != nil {
		r["numEdgesWritten"] = *d.NumEdgesWritten
	}
	if d.NumVerticesWritten != nil {
		r["numVerticesWritten"] = *d.NumVerticesWritten
	}
	return r
}

func exportFilterToResponse(f *ngstore.ExportFilter) map[string]interface{} {
	if f == nil {
		return nil
	}
	r := map[string]interface{}{}
	if len(f.EdgeFilter) > 0 {
		ef := make(map[string]interface{})
		for k, v := range f.EdgeFilter {
			ef[k] = exportFilterElementToResponse(&v)
		}
		r["edgeFilter"] = ef
	}
	if len(f.VertexFilter) > 0 {
		vf := make(map[string]interface{})
		for k, v := range f.VertexFilter {
			vf[k] = exportFilterElementToResponse(&v)
		}
		r["vertexFilter"] = vf
	}
	return r
}

func exportFilterElementToResponse(e *ngstore.ExportFilterElement) map[string]interface{} {
	if e == nil {
		return nil
	}
	r := map[string]interface{}{}
	for k, v := range e.Properties {
		props := map[string]interface{}{
			"multiValueHandling": v.MultiValueHandling,
		}
		if v.OutputType != nil {
			props["outputType"] = *v.OutputType
		}
		if v.SourcePropertyName != nil {
			props["sourcePropertyName"] = *v.SourcePropertyName
		}
		r[k] = props
	}
	return r
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

func parseTagsFromParams(params map[string]interface{}) map[string]string {
	v, ok := params["tags"]
	if !ok || v == nil {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	tags := make(map[string]string)
	key, hasKey := m["key"]
	val, hasVal := m["value"]
	if hasKey && hasVal {
		if ks, ok := key.(string); ok {
			if vs, ok := val.(string); ok {
				tags[ks] = vs
				return tags
			}
		}
	}
	keyUpper, hasKeyUpper := m["Key"]
	valUpper, hasValUpper := m["Value"]
	if hasKeyUpper && hasValUpper {
		if ks, ok := keyUpper.(string); ok {
			if vs, ok := valUpper.(string); ok {
				tags[ks] = vs
				return tags
			}
		}
	}
	for k, v := range m {
		if s, ok := v.(string); ok {
			tags[k] = s
		}
	}
	return tags
}

func parseImportOptions(params map[string]interface{}) *ngstore.ImportOptions {
	v, ok := params["importOptions"]
	if !ok || v == nil {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	neptune, ok := m["neptune"]
	if !ok || neptune == nil {
		return nil
	}
	nm, ok := neptune.(map[string]interface{})
	if !ok {
		return nil
	}
	opts := &ngstore.ImportOptions{
		Neptune: &ngstore.NeptuneImportOptions{},
	}
	if s, ok := nm["s3ExportPath"].(string); ok {
		opts.Neptune.S3ExportPath = s
	}
	if s, ok := nm["s3ExportKmsKeyId"].(string); ok {
		opts.Neptune.S3ExportKmsKeyId = s
	}
	if b, ok := nm["preserveDefaultVertexLabels"].(bool); ok {
		opts.Neptune.PreserveDefaultVertexLabels = &b
	}
	if b, ok := nm["preserveEdgeIds"].(bool); ok {
		opts.Neptune.PreserveEdgeIds = &b
	}
	return opts
}

func parseExportFilter(params map[string]interface{}) *ngstore.ExportFilter {
	v, ok := params["exportFilter"]
	if !ok || v == nil {
		return nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var f ngstore.ExportFilter
	if err := json.Unmarshal(data, &f); err != nil {
		return nil
	}
	return &f
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

func parseNTriplesLine(line string) (subject, predicate, object string, ok bool) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return "", "", "", false
	}
	if line[len(line)-1] == '.' {
		line = strings.TrimSpace(line[:len(line)-1])
	}

	var pos int

	subject, pos, ok = parseNTTerm(line, pos)
	if !ok {
		return "", "", "", false
	}
	for pos < len(line) && line[pos] == ' ' {
		pos++
	}

	predicate, pos, ok = parseNTTerm(line, pos)
	if !ok {
		return "", "", "", false
	}
	for pos < len(line) && line[pos] == ' ' {
		pos++
	}

	object, _, ok = parseNTTerm(line, pos)
	if !ok {
		return "", "", "", false
	}

	return subject, predicate, object, true
}

func parseNTTerm(s string, pos int) (string, int, bool) {
	if pos >= len(s) {
		return "", pos, false
	}

	if s[pos] == '<' {
		end := strings.IndexByte(s[pos+1:], '>')
		if end < 0 {
			return "", pos, false
		}
		return s[pos : pos+end+2], pos + end + 2, true
	}

	if s[pos] == '"' {
		i := pos + 1
		for i < len(s) {
			if s[i] == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			if s[i] == '"' {
				break
			}
			i++
		}
		if i >= len(s) {
			return "", pos, false
		}
		end := i + 1
		if end < len(s) && s[end] == '@' {
			j := end + 1
			for j < len(s) && isNTAlphaNumHyphen(s[j]) {
				j++
			}
			return s[pos:j], j, true
		}
		if end+1 < len(s) && s[end] == '^' && end+2 < len(s) && s[end+1] == '^' && end+2 < len(s) && s[end+2] == '<' {
			close := strings.IndexByte(s[end+2:], '>')
			if close >= 0 {
				return s[pos : end+2+close+1], end + 2 + close + 1, true
			}
		}
		return s[pos:end], end, true
	}

	if strings.HasPrefix(s[pos:], "_:") {
		i := pos + 2
		for i < len(s) && s[i] != ' ' && s[i] != '\t' && s[i] != '.' {
			i++
		}
		return s[pos:i], i, true
	}

	return "", pos, false
}

func isNTAlphaNumHyphen(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-'
}

func extractLocalName(uri string) string {
	uri = strings.Trim(uri, "<>")
	if idx := strings.LastIndex(uri, "/"); idx >= 0 {
		return uri[idx+1:]
	}
	if idx := strings.LastIndex(uri, "#"); idx >= 0 {
		return uri[idx+1:]
	}
	return uri
}
