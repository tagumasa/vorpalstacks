package neptunegraph

// Graph resource Core functions: the single validation and persistence path
// for the graph CRUD/lifecycle operations on both protocol planes.

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/graphengine"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
)

// CreateGraphInput carries the wire-parsed CreateGraph request. VectorSearchConfig
// holds the raw vectorSearchConfiguration wire value; the Core parses and
// validates it so validation lives on the service layer.
type CreateGraphInput struct {
	GraphName          string
	ProvisionedMemory  int
	HasReplicaCount    bool
	ReplicaCount       int
	KmsKeyIdentifier   string
	DeletionProtection bool
	PublicConnectivity bool
	VectorSearchConfig interface{}
	Tags               map[string]string
	Region             string
}

// GetGraphInput carries the wire-parsed GetGraph request.
type GetGraphInput struct {
	GraphIdentifier string
}

// ListGraphsInput carries the wire-parsed ListGraphs request with the
// already-clamped page size.
type ListGraphsInput struct {
	MaxItems int
	Marker   string
}

// ListGraphsResult carries the graph records for one list page.
type ListGraphsResult struct {
	Graphs    []*ngstore.Graph
	NextToken string
	Truncated bool
}

// UpdateGraphInput carries the wire-parsed UpdateGraph request. The Has*
// members preserve the wire presence of the optional update members.
type UpdateGraphInput struct {
	GraphIdentifier       string
	HasProvisionedMemory  bool
	ProvisionedMemory     int
	HasDeletionProtection bool
	DeletionProtection    bool
	HasPublicConnectivity bool
	PublicConnectivity    bool
}

// DeleteGraphInput carries the wire-parsed DeleteGraph request.
type DeleteGraphInput struct {
	GraphIdentifier string
	HasSkipSnapshot bool
	SkipSnapshot    bool
	Region          string
}

// StartGraphInput carries the wire-parsed StartGraph request.
type StartGraphInput struct {
	GraphIdentifier string
}

// StopGraphInput carries the wire-parsed StopGraph request.
type StopGraphInput struct {
	GraphIdentifier string
}

// ResetGraphInput carries the wire-parsed ResetGraph request.
type ResetGraphInput struct {
	GraphIdentifier string
	HasSkipSnapshot bool
	SkipSnapshot    bool
}

// RestoreGraphFromSnapshotInput carries the wire-parsed
// RestoreGraphFromSnapshot request.
type RestoreGraphFromSnapshotInput struct {
	SnapshotIdentifier   string
	GraphName            string
	DeletionProtection   bool
	PublicConnectivity   bool
	HasProvisionedMemory bool
	ProvisionedMemory    int
	HasReplicaCount      bool
	ReplicaCount         int
	Tags                 map[string]string
	Region               string
}

// resolveGraphIdentifier resolves a graphIdentifier parameter (which may be
// either a graph ID like "g-xxxxxxxxxx" or a graph name) to a Graph.
// AWS NeptuneGraph accepts both forms for most operations.
func (s *NeptuneGraphService) resolveGraphIdentifier(store *ngstore.NeptuneGraphStore, identifier string) (*ngstore.Graph, error) {
	graph, err := store.GetGraph(identifier)
	if err == nil {
		return graph, nil
	}
	if !ngstore.IsNotFound(err) {
		return nil, err
	}

	graphs, _, _, listErr := store.ListGraphs(storecommon.ListOptions{})
	if listErr != nil {
		return nil, listErr
	}
	for _, g := range graphs {
		if g.Name == identifier {
			return g, nil
		}
	}
	return nil, ngstore.ErrGraphNotFound
}

// parseVectorSearchConfigValue parses the raw vectorSearchConfiguration wire
// value, validating the dimension per the Smithy VectorSearchDimension shape.
// Absent or non-object values are ignored, matching the wire behaviour.
func parseVectorSearchConfigValue(v interface{}) (*ngstore.VectorSearchConfig, error) {
	if v == nil {
		return nil, nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, nil
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
	if err := validateVectorSearchDimension(dim); err != nil {
		return nil, err
	}
	return &ngstore.VectorSearchConfig{Dimension: int32(dim)}, nil
}

// createGraphCore validates the request, persists the graph record, opens its
// query engine and applies create-time tags.
func (s *NeptuneGraphService) createGraphCore(ctx context.Context, store *ngstore.NeptuneGraphStore, in *CreateGraphInput) (*ngstore.Graph, error) {
	graphName := in.GraphName
	if err := validateGraphName(graphName); err != nil {
		return nil, err
	}

	mem := in.ProvisionedMemory
	if err := validateProvisionedMemory(mem, true); err != nil {
		return nil, err
	}

	replicaCount := 1
	if in.HasReplicaCount {
		replicaCount = in.ReplicaCount
		if err := validateReplicaCount(replicaCount); err != nil {
			return nil, err
		}
	}

	kmsKey := in.KmsKeyIdentifier
	if err := validateKmsKeyArn(kmsKey, false); err != nil {
		return nil, err
	}

	graphID := generateID("g-")
	now := time.Now().UTC()

	graph := &ngstore.Graph{
		Id:                 graphID,
		Name:               graphName,
		Arn:                s.arnBuilder.NeptuneGraph().Graph(graphID),
		Status:             "CREATING",
		ProvisionedMemory:  proto.Int32(int32(mem)),
		ReplicaCount:       proto.Int32(int32(replicaCount)),
		DeletionProtection: in.DeletionProtection,
		PublicConnectivity: in.PublicConnectivity,
		KmsKeyIdentifier:   in.KmsKeyIdentifier,
		BuildNumber:        neptuneGraphBuildNumber,
		CreateTime:         &now,
		AccountID:          s.accountID,
		Region:             in.Region,
	}

	if vsc, err := parseVectorSearchConfigValue(in.VectorSearchConfig); err != nil {
		return nil, err
	} else if vsc != nil {
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
		logs.Error("failed to open graph engine", logs.String("graphId", graphID), logs.Err(err))
		graph.Status = "FAILED"
		graph.StatusReason = fmt.Sprintf("failed to open graph engine: %v", err)
		if updateErr := store.UpdateGraph(graph); updateErr != nil {
			logs.Error("failed to update graph status to FAILED", logs.String("graphId", graphID), logs.Err(updateErr))
		}
		return nil, newInternalServerException(err)
	}
	s.enginesMu.Lock()
	s.activeEngines[graphID] = &engineEntry{db: db}
	s.enginesMu.Unlock()

	graph.Status = "AVAILABLE"
	graph.Endpoint = s.graphEndpoint(graphID)
	if err := store.UpdateGraph(graph); err != nil {
		logs.Error("failed to update graph status to AVAILABLE", logs.String("graphId", graphID), logs.Err(err))
		s.enginesMu.Lock()
		delete(s.activeEngines, graphID)
		s.enginesMu.Unlock()
		db.Close()
		graph.Status = "FAILED"
		graph.StatusReason = "failed to persist graph status"
		graph.Endpoint = ""
		if updateErr := store.UpdateGraph(graph); updateErr != nil {
			logs.Error("failed to update graph status to FAILED after AVAILABLE failure", logs.String("graphId", graphID), logs.Err(updateErr))
		}
		return nil, newInternalServerException(fmt.Errorf("failed to update graph status: %w", err))
	}

	if tags := in.Tags; len(tags) > 0 {
		if err := store.AddTags(graph.Arn, tags); err != nil {
			logs.Warn("failed to store tags for graph", logs.String("graphId", graphID), logs.Err(err))
		}
	}

	return graph, nil
}

// getGraphCore resolves a graph by identifier (ID or name), mapping a missing
// graph to the documented ResourceNotFoundException.
func (s *NeptuneGraphService) getGraphCore(store *ngstore.NeptuneGraphStore, in *GetGraphInput) (*ngstore.Graph, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := s.resolveGraphIdentifier(store, graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	return graph, nil
}

// getGraphByIDCore reads a graph by raw identifier without the name fallback,
// returning the store error unchanged for the admin plane.
func (s *NeptuneGraphService) getGraphByIDCore(store *ngstore.NeptuneGraphStore, graphID string) (*ngstore.Graph, error) {
	return store.GetGraph(graphID)
}

// listGraphsCore reads one page of graph records, returning the store error
// unchanged so both planes keep their own error mapping.
func (s *NeptuneGraphService) listGraphsCore(store *ngstore.NeptuneGraphStore, in *ListGraphsInput) (*ListGraphsResult, error) {
	graphs, nextToken, truncated, err := store.ListGraphs(storecommon.ListOptions{
		MaxItems: in.MaxItems,
		Marker:   in.Marker,
	})
	if err != nil {
		return nil, err
	}
	return &ListGraphsResult{Graphs: graphs, NextToken: nextToken, Truncated: truncated}, nil
}

// updateGraphCore applies presence-based configuration updates to an
// AVAILABLE graph.
func (s *NeptuneGraphService) updateGraphCore(store *ngstore.NeptuneGraphStore, in *UpdateGraphInput) (*ngstore.Graph, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := s.resolveGraphIdentifier(store, graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	if graph.Status != "AVAILABLE" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graph is not in AVAILABLE state")
	}

	if in.HasProvisionedMemory {
		mem := in.ProvisionedMemory
		if mem < minProvisionedMemory || mem > maxProvisionedMemory {
			return nil, newValidationException("CONSTRAINT_VIOLATION", "provisionedMemory")
		}
		graph.ProvisionedMemory = proto.Int32(int32(mem))
	}
	if in.HasDeletionProtection {
		graph.DeletionProtection = in.DeletionProtection
	}
	if in.HasPublicConnectivity {
		graph.PublicConnectivity = in.PublicConnectivity
	}

	if err := store.UpdateGraph(graph); err != nil {
		return nil, err
	}

	return graph, nil
}

// deleteGraphCore removes a graph, optionally creating a final snapshot, and
// cleans up tags, engine and storage bucket.
func (s *NeptuneGraphService) deleteGraphCore(store *ngstore.NeptuneGraphStore, in *DeleteGraphInput) (*ngstore.Graph, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	if !in.HasSkipSnapshot {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "skipSnapshot is required")
	}

	graph, err := s.resolveGraphIdentifier(store, graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}
	graphID = graph.Id

	if graph.DeletionProtection {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graph has deletion protection enabled")
	}

	skipSnapshot := in.SkipSnapshot

	var snapshotID string
	if !skipSnapshot && graph.Status != "DELETING" {
		snapshotID = generateID("gs-")
		now := time.Now().UTC()
		snapshot := &ngstore.GraphSnapshot{
			Id:                 snapshotID,
			Name:               graph.Name + "-auto-snapshot",
			Arn:                s.arnBuilder.NeptuneGraph().Snapshot(snapshotID),
			Status:             "CREATING",
			SourceGraphId:      graphID,
			SnapshotCreateTime: &now,
			AccountID:          s.accountID,
			Region:             graph.Region,
		}
		if err := store.CreateSnapshot(snapshot); err != nil {
			logs.Warn("failed to create auto snapshot during graph deletion", logs.String("graphId", graphID), logs.Err(err))
			snapshotID = ""
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
	var engineEntry *engineEntry
	if e, ok := s.activeEngines[graphID]; ok {
		e.stopped = true
		engineEntry = e
		delete(s.activeEngines, graphID)
	}
	s.enginesMu.Unlock()
	s.planCache.purgeByGraph(graphID)

	if engineEntry != nil {
		engineEntry.wg.Wait()
		engineEntry.db.Close()
	}

	if snapshotID != "" {
		srcBkt, srcErr := s.graphBucket(graphID)
		dstBkt, dstErr := s.graphBucket("snapshot:" + snapshotID)
		copyOK := false
		if srcErr == nil && dstErr == nil {
			if err := copyGraphBucket(srcBkt, dstBkt); err != nil {
				logs.Warn("failed to copy graph data to snapshot bucket during deletion",
					logs.String("graphId", graphID), logs.String("snapshotId", snapshotID), logs.Err(err))
			} else {
				copyOK = true
			}
		}
		finalStatus := "AVAILABLE"
		if !copyOK {
			finalStatus = "FAILED"
		}
		if snap, getErr := store.GetSnapshot(snapshotID); getErr == nil {
			snap.Status = finalStatus
			if updateErr := store.UpdateSnapshot(snap); updateErr != nil {
				logs.Warn("failed to update snapshot status after copy",
					logs.String("snapshotId", snapshotID), logs.Err(updateErr))
			}
		}
	}

	if rs, err := s.storageManager.GetStorage(in.Region); err == nil {
		rs.DeleteBucket("neptunegraph:graph:" + graphID)
	}

	if err := store.DeleteGraph(graphID); err != nil {
		logs.Error("failed to delete graph from store", logs.String("graphId", graphID), logs.Err(err))
		return nil, newInternalServerException(fmt.Errorf("failed to delete graph: %w", err))
	}

	return graph, nil
}

// startGraphCore transitions a STOPPED graph to AVAILABLE by reopening its
// query engine.
func (s *NeptuneGraphService) startGraphCore(store *ngstore.NeptuneGraphStore, in *StartGraphInput) (*ngstore.Graph, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := s.resolveGraphIdentifier(store, graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}
	graphID = graph.Id

	if graph.Status != "STOPPED" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graph is not in STOPPED state")
	}

	graph.Status = "STARTING"
	if err := store.UpdateGraph(graph); err != nil {
		graph.Status = "STOPPED"
		return nil, newInternalServerException(fmt.Errorf("failed to update graph status to STARTING: %w", err))
	}

	bucket, err := s.graphBucket(graphID)
	if err != nil {
		graph.Status = "STOPPED"
		if updateErr := store.UpdateGraph(graph); updateErr != nil {
			logs.Error("failed to restore graph status to STOPPED after graphBucket failure",
				logs.String("graphId", graphID), logs.Err(updateErr))
		}
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
	graph.Endpoint = s.graphEndpoint(graphID)
	if err := store.UpdateGraph(graph); err != nil {
		logs.Error("failed to update graph status to AVAILABLE", logs.String("graphId", graphID), logs.Err(err))
		s.enginesMu.Lock()
		delete(s.activeEngines, graphID)
		s.enginesMu.Unlock()
		db.Close()
		graph.Status = "STOPPED"
		graph.Endpoint = ""
		if updateErr := store.UpdateGraph(graph); updateErr != nil {
			logs.Error("failed to restore graph status to STOPPED", logs.String("graphId", graphID), logs.Err(updateErr))
		}
		return nil, newInternalServerException(fmt.Errorf("failed to update graph status to AVAILABLE: %w", err))
	}

	return graph, nil
}

// stopGraphCore gracefully shuts down a graph's query engine and transitions
// it to STOPPED state.
func (s *NeptuneGraphService) stopGraphCore(store *ngstore.NeptuneGraphStore, in *StopGraphInput) (*ngstore.Graph, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := s.resolveGraphIdentifier(store, graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}
	graphID = graph.Id

	if graph.Status != "AVAILABLE" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graph is not in AVAILABLE state")
	}

	graph.Status = "STOPPING"
	if err := store.UpdateGraph(graph); err != nil {
		graph.Status = "AVAILABLE"
		return nil, newInternalServerException(fmt.Errorf("failed to update graph status to STOPPING: %w", err))
	}

	s.enginesMu.Lock()
	var entry *engineEntry
	if e, ok := s.activeEngines[graphID]; ok {
		e.stopped = true
		entry = e
		delete(s.activeEngines, graphID)
	}
	s.enginesMu.Unlock()

	if entry != nil {
		entry.wg.Wait()
		entry.db.Close()
	}

	graph.Status = "STOPPED"
	graph.Endpoint = ""
	if err := store.UpdateGraph(graph); err != nil {
		logs.Error("failed to update graph status to STOPPED", logs.String("graphId", graphID), logs.Err(err))
		return nil, newInternalServerException(fmt.Errorf("failed to update graph status to STOPPED: %w", err))
	}

	return graph, nil
}

// resetGraphCore clears all data from an AVAILABLE graph's engine while
// keeping the graph resource intact.
func (s *NeptuneGraphService) resetGraphCore(store *ngstore.NeptuneGraphStore, in *ResetGraphInput) (*ngstore.Graph, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	if !in.HasSkipSnapshot {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "skipSnapshot is required")
	}

	graph, err := s.resolveGraphIdentifier(store, graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}
	graphID = graph.Id

	if graph.Status != "AVAILABLE" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graph is not in AVAILABLE state")
	}

	graph.Status = "RESETTING"
	if err := store.UpdateGraph(graph); err != nil {
		graph.Status = "AVAILABLE"
		return nil, newInternalServerException(fmt.Errorf("failed to update graph status to RESETTING: %w", err))
	}

	// When the request does not skip it, a final snapshot of the graph data
	// is taken before the engine is cleared, mirroring DeleteGraph.
	var resetSnapshotID string
	if !in.SkipSnapshot {
		resetSnapshotID = generateID("gs-")
		now := time.Now().UTC()
		snapshot := &ngstore.GraphSnapshot{
			Id:                 resetSnapshotID,
			Name:               graph.Name + "-auto-snapshot",
			Arn:                s.arnBuilder.NeptuneGraph().Snapshot(resetSnapshotID),
			Status:             "CREATING",
			SourceGraphId:      graphID,
			SnapshotCreateTime: &now,
			AccountID:          s.accountID,
			Region:             graph.Region,
		}
		if err := store.CreateSnapshot(snapshot); err != nil {
			logs.Warn("failed to create final snapshot during graph reset", logs.String("graphId", graphID), logs.Err(err))
			resetSnapshotID = ""
		}
	}

	if resetSnapshotID != "" {
		srcBkt, srcErr := s.graphBucket(graphID)
		dstBkt, dstErr := s.graphBucket("snapshot:" + resetSnapshotID)
		copyOK := false
		if srcErr == nil && dstErr == nil {
			if err := copyGraphBucket(srcBkt, dstBkt); err != nil {
				logs.Warn("failed to copy graph data to final snapshot bucket during reset",
					logs.String("graphId", graphID), logs.String("snapshotId", resetSnapshotID), logs.Err(err))
			} else {
				copyOK = true
			}
		}
		finalStatus := "AVAILABLE"
		if !copyOK {
			finalStatus = "FAILED"
		}
		if snap, getErr := store.GetSnapshot(resetSnapshotID); getErr == nil {
			snap.Status = finalStatus
			if updateErr := store.UpdateSnapshot(snap); updateErr != nil {
				logs.Warn("failed to update final snapshot status after copy",
					logs.String("snapshotId", resetSnapshotID), logs.Err(updateErr))
			}
		}
	}

	s.enginesMu.RLock()
	entry, ok := s.activeEngines[graphID]
	s.enginesMu.RUnlock()

	if ok {
		entry.mu.Lock()
		if err := entry.db.Clear(); err != nil {
			entry.mu.Unlock()

			graph.Status = "FAILED"
			graph.StatusReason = err.Error()
			if storeErr := store.UpdateGraph(graph); storeErr != nil {
				logs.Warn("Failed to update graph status to FAILED after Clear error",
					logs.String("graphId", graphID),
					logs.Err(storeErr))
			}
			return nil, newInternalServerException(err)
		}
		entry.mu.Unlock()
	}

	graph.Status = "AVAILABLE"
	if err := store.UpdateGraph(graph); err != nil {
		logs.Error("failed to update graph status to AVAILABLE after reset", logs.String("graphId", graphID), logs.Err(err))
		return nil, newInternalServerException(fmt.Errorf("failed to update graph status to AVAILABLE: %w", err))
	}

	return graph, nil
}

// restoreGraphFromSnapshotCore creates a new graph from an existing snapshot,
// copying the source graph data.
func (s *NeptuneGraphService) restoreGraphFromSnapshotCore(store *ngstore.NeptuneGraphStore, in *RestoreGraphFromSnapshotInput) (*ngstore.Graph, error) {
	snapshotID := in.SnapshotIdentifier
	if snapshotID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "snapshotIdentifier")
	}

	graphName := in.GraphName
	if err := validateGraphName(graphName); err != nil {
		return nil, err
	}

	_, err := store.GetSnapshot(snapshotID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("snapshot", snapshotID)
		}
		return nil, err
	}

	graphID := generateID("g-")
	now := time.Now().UTC()

	graph := &ngstore.Graph{
		Id:                 graphID,
		Name:               graphName,
		Arn:                s.arnBuilder.NeptuneGraph().Graph(graphID),
		Status:             "CREATING",
		ProvisionedMemory:  proto.Int32(128),
		ReplicaCount:       proto.Int32(1),
		DeletionProtection: in.DeletionProtection,
		PublicConnectivity: in.PublicConnectivity,
		BuildNumber:        neptuneGraphBuildNumber,
		SourceSnapshotId:   snapshotID,
		CreateTime:         &now,
		AccountID:          s.accountID,
		Region:             in.Region,
	}

	if in.HasProvisionedMemory {
		mem := in.ProvisionedMemory
		if err := validateProvisionedMemory(mem, false); err != nil {
			return nil, err
		}
		if mem > 0 {
			graph.ProvisionedMemory = proto.Int32(int32(mem))
		}
	}

	if in.HasReplicaCount {
		rc := in.ReplicaCount
		if err := validateReplicaCount(rc); err != nil {
			return nil, err
		}
		graph.ReplicaCount = proto.Int32(int32(rc))
	}

	if err := store.CreateGraph(graph); err != nil {
		if ngstore.IsAlreadyExists(err) {
			return nil, newConflictException("CONCURRENT_MODIFICATION")
		}
		return nil, err
	}

	if tags := in.Tags; len(tags) > 0 {
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
		graph.Status = "FAILED"
		graph.StatusReason = fmt.Sprintf("failed to open graph engine: %v", err)
		if updateErr := store.UpdateGraph(graph); updateErr != nil {
			logs.Error("failed to update graph status to FAILED after restore", logs.String("graphId", graphID), logs.Err(updateErr))
		}
		return nil, newInternalServerException(err)
	}

	s.enginesMu.Lock()
	s.activeEngines[graphID] = &engineEntry{db: db}
	s.enginesMu.Unlock()

	graph.Status = "AVAILABLE"
	graph.Endpoint = s.graphEndpoint(graphID)
	if err := store.UpdateGraph(graph); err != nil {
		logs.Error("failed to update restored graph status to AVAILABLE", logs.String("graphId", graphID), logs.Err(err))
		s.enginesMu.Lock()
		delete(s.activeEngines, graphID)
		s.enginesMu.Unlock()
		db.Close()
		graph.Status = "FAILED"
		graph.StatusReason = "failed to persist graph status"
		graph.Endpoint = ""
		if updateErr := store.UpdateGraph(graph); updateErr != nil {
			logs.Error("failed to update graph status to FAILED after restore AVAILABLE failure", logs.String("graphId", graphID), logs.Err(updateErr))
		}
		return nil, newInternalServerException(fmt.Errorf("failed to update restored graph status: %w", err))
	}

	return graph, nil
}

func copyGraphBucket(src, dst storage.BatchBucket) error {
	return src.ForEach(func(k, v []byte) error {
		return dst.Put(k, v)
	})
}
