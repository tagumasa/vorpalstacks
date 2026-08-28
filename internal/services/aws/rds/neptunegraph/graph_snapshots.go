package neptunegraph

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateGraphSnapshot creates a point-in-time snapshot of the specified graph.
func (s *NeptuneGraphService) CreateGraphSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &CreateGraphSnapshotInput{
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
		SnapshotName:    request.GetStringParam(req.Parameters, "snapshotName"),
		Region:          reqCtx.GetRegion(),
	}

	snapshot, err := s.createGraphSnapshotCore(store, in)
	if err != nil {
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

	in := &GetGraphSnapshotInput{
		SnapshotIdentifier: request.GetStringParam(req.Parameters, "snapshotIdentifier"),
	}

	snapshot, err := s.getGraphSnapshotCore(store, in)
	if err != nil {
		return nil, err
	}
	return snapshotToResponse(snapshot), nil
}

// ListGraphSnapshots returns a paginated list of graph snapshots, optionally filtered by graph.
func (s *NeptuneGraphService) ListGraphSnapshots(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &ListGraphSnapshotsInput{
		MaxItems:        clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		Marker:          request.GetStringParam(req.Parameters, "nextToken"),
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
	}

	res, err := s.listGraphSnapshotsCore(store, in)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(res.Snapshots))
	for _, sn := range res.Snapshots {
		items = append(items, snapshotToResponse(sn))
	}

	result := map[string]interface{}{
		"graphSnapshots": items,
	}
	if res.Truncated {
		result["nextToken"] = res.NextToken
	}
	return result, nil
}

// DeleteGraphSnapshot removes a graph snapshot by its identifier.
func (s *NeptuneGraphService) DeleteGraphSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &DeleteGraphSnapshotInput{
		SnapshotIdentifier: request.GetStringParam(req.Parameters, "snapshotIdentifier"),
		Region:             reqCtx.GetRegion(),
	}

	snapshot, err := s.deleteGraphSnapshotCore(store, in)
	if err != nil {
		return nil, err
	}
	return snapshotToResponse(snapshot), nil
}
