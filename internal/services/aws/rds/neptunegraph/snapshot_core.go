package neptunegraph

// Snapshot resource Core functions: the single validation and persistence
// path for the graph snapshot operations.

import (
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
)

// CreateGraphSnapshotInput carries the wire-parsed CreateGraphSnapshot request.
type CreateGraphSnapshotInput struct {
	GraphIdentifier string
	SnapshotName    string
	Region          string
}

// GetGraphSnapshotInput carries the wire-parsed GetGraphSnapshot request.
type GetGraphSnapshotInput struct {
	SnapshotIdentifier string
}

// ListGraphSnapshotsInput carries the wire-parsed ListGraphSnapshots request
// with the already-clamped page size.
type ListGraphSnapshotsInput struct {
	MaxItems        int
	Marker          string
	GraphIdentifier string
}

// ListGraphSnapshotsResult carries the snapshot records for one list page.
type ListGraphSnapshotsResult struct {
	Snapshots []*ngstore.GraphSnapshot
	NextToken string
	Truncated bool
}

// DeleteGraphSnapshotInput carries the wire-parsed DeleteGraphSnapshot request.
type DeleteGraphSnapshotInput struct {
	SnapshotIdentifier string
	Region             string
}

// createGraphSnapshotCore validates the request and creates a point-in-time
// snapshot of the specified AVAILABLE graph.
func (s *NeptuneGraphService) createGraphSnapshotCore(store *ngstore.NeptuneGraphStore, in *CreateGraphSnapshotInput) (*ngstore.GraphSnapshot, error) {
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

	snapshotName := in.SnapshotName
	if snapshotName == "" || strings.HasPrefix(snapshotName, "gs-") || !snapshotNameRegex.MatchString(snapshotName) || len(snapshotName) > 63 {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "snapshotName")
	}

	snapshotID := generateID("gs-")
	now := time.Now().UTC()

	snapshot := &ngstore.GraphSnapshot{
		Id:                 snapshotID,
		Name:               snapshotName,
		Arn:                s.arnBuilder.NeptuneGraph().Snapshot(snapshotID),
		Status:             "CREATING",
		SourceGraphId:      graph.Id,
		SnapshotCreateTime: &now,
		AccountID:          s.accountID,
		Region:             in.Region,
	}

	if err := store.CreateSnapshot(snapshot); err != nil {
		if ngstore.IsAlreadyExists(err) {
			return nil, newConflictException("CONCURRENT_MODIFICATION")
		}
		return nil, err
	}

	srcBkt, srcErr := s.graphBucket(graph.Id)
	dstBkt, dstErr := s.graphBucket("snapshot:" + snapshotID)
	copyOK := false
	if srcErr == nil && dstErr == nil {
		if err := copyGraphBucket(srcBkt, dstBkt); err != nil {
			logs.Warn("failed to copy graph data to snapshot bucket",
				logs.String("graphId", graph.Id), logs.String("snapshotId", snapshotID), logs.Err(err))
		} else {
			copyOK = true
		}
	}

	snapshot.Status = "AVAILABLE"
	if !copyOK {
		snapshot.Status = "FAILED"
	}
	if err := store.UpdateSnapshot(snapshot); err != nil {
		logs.Warn("failed to update snapshot status after copy",
			logs.String("snapshotId", snapshotID), logs.Err(err))
	}

	return snapshot, nil
}

// getGraphSnapshotCore retrieves a snapshot by identifier, hiding snapshots
// that are mid-deletion.
func (s *NeptuneGraphService) getGraphSnapshotCore(store *ngstore.NeptuneGraphStore, in *GetGraphSnapshotInput) (*ngstore.GraphSnapshot, error) {
	snapshotID := in.SnapshotIdentifier
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

	return snapshot, nil
}

// listGraphSnapshotsCore reads one page of snapshot records, optionally
// filtered by graph.
func (s *NeptuneGraphService) listGraphSnapshotsCore(store *ngstore.NeptuneGraphStore, in *ListGraphSnapshotsInput) (*ListGraphSnapshotsResult, error) {
	snapshots, nextToken, truncated, err := store.ListSnapshots(storecommon.ListOptions{
		MaxItems: in.MaxItems,
		Marker:   in.Marker,
	}, in.GraphIdentifier)
	if err != nil {
		return nil, err
	}
	return &ListGraphSnapshotsResult{Snapshots: snapshots, NextToken: nextToken, Truncated: truncated}, nil
}

// deleteGraphSnapshotCore removes a snapshot record and its storage bucket.
func (s *NeptuneGraphService) deleteGraphSnapshotCore(store *ngstore.NeptuneGraphStore, in *DeleteGraphSnapshotInput) (*ngstore.GraphSnapshot, error) {
	snapshotID := in.SnapshotIdentifier
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

	if err := store.DeleteSnapshot(snapshotID); err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("snapshot", snapshotID)
		}
		return nil, err
	}

	if rs, err := s.storageManager.GetStorage(in.Region); err == nil {
		rs.DeleteBucket("neptunegraph:graph:snapshot:" + snapshotID)
	}

	return snapshot, nil
}
