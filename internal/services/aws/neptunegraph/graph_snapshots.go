package neptunegraph

import (
	"context"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
)

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
