package graphengine

import (
	"fmt"

	"github.com/cockroachdb/pebble/v2"
)

func (d *DB) allocEdgeID() EdgeID {
	return EdgeID(d.nextEdgeID.Add(1))
}

// AddEdge creates a directed edge from one node to another. Both source and
// target nodes must exist; otherwise an error is returned. Three Pebble keys
// are written atomically: the edge data, an outgoing adjacency entry, and an
// incoming adjacency entry.
func (d *DB) AddEdge(from, to NodeID, label string, props Props) (EdgeID, error) {
	if d.closed.Load() {
		return 0, fmt.Errorf("graphengine: database is closed")
	}

	if exists, err := d.NodeExists(from); err != nil {
		return 0, err
	} else if !exists {
		return 0, fmt.Errorf("graphengine: source node %d not found", from)
	}
	if exists, err := d.NodeExists(to); err != nil {
		return 0, err
	} else if !exists {
		return 0, fmt.Errorf("graphengine: target node %d not found", to)
	}

	id := d.allocEdgeID()

	edge := &Edge{
		ID:    id,
		From:  from,
		To:    to,
		Label: label,
		Props: props,
	}

	batch := d.db.NewBatch()
	defer batch.Close()

	propsBytes := marshalProps(props)

	edgeData, err := encodeEdgeDataRaw(edge, propsBytes)
	if err != nil {
		return 0, err
	}
	_ = batch.Set(edgeKey(id), edgeData, nil)
	_ = batch.Set(adjOutKey(from, id), encodeAdjValue(to, label), nil)
	_ = batch.Set(adjInKey(to, id), encodeAdjValue(from, label), nil)
	if label != "" {
		_ = batch.Set(idxEtypeKey(label, id), nil, nil)
	}

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return 0, fmt.Errorf("graphengine: failed to add edge: %w", err)
	}

	d.edgeCount.Add(1)
	return id, nil
}

// AddEdgeBatch adds multiple edges in a single batch. Caller must ensure all
// source and target nodes exist; no existence checks are performed to avoid
// per-edge read overhead within the batch.
func (d *DB) AddEdgeBatch(edges []Edge) ([]EdgeID, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	ids := make([]EdgeID, len(edges))
	batch := d.db.NewBatch()
	defer batch.Close()

	for i := range edges {
		id := d.allocEdgeID()
		ids[i] = id

		edge := edges[i]
		edge.ID = id

		propsBytes := marshalProps(edge.Props)

		edgeData, err := encodeEdgeDataRaw(&edge, propsBytes)
		if err != nil {
			return nil, err
		}
		_ = batch.Set(edgeKey(id), edgeData, nil)
		_ = batch.Set(adjOutKey(edge.From, id), encodeAdjValue(edge.To, edge.Label), nil)
		_ = batch.Set(adjInKey(edge.To, id), encodeAdjValue(edge.From, edge.Label), nil)
		if edge.Label != "" {
			_ = batch.Set(idxEtypeKey(edge.Label, id), nil, nil)
		}
	}

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return nil, fmt.Errorf("graphengine: batch edge add failed: %w", err)
	}

	d.edgeCount.Add(uint64(len(edges)))
	return ids, nil
}

// GetEdge retrieves an edge by its ID. Returns an error if the edge does not exist.
func (d *DB) GetEdge(id EdgeID) (*Edge, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	val, closer, err := d.db.Get(edgeKey(id))
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, fmt.Errorf("graphengine: edge %d not found", id)
		}
		return nil, fmt.Errorf("graphengine: failed to get edge %d: %w", id, err)
	}
	defer closer.Close()

	return decodeEdgeData(id, val)
}

// DeleteEdge removes an edge and its adjacency and index entries.
func (d *DB) DeleteEdge(id EdgeID) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	edge, err := d.GetEdge(id)
	if err != nil {
		return err
	}

	batch := d.db.NewBatch()
	defer batch.Close()

	_ = batch.Delete(edgeKey(id), nil)
	_ = batch.Delete(adjOutKey(edge.From, id), nil)
	_ = batch.Delete(adjInKey(edge.To, id), nil)
	if edge.Label != "" {
		_ = batch.Delete(idxEtypeKey(edge.Label, id), nil)
	}

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to delete edge %d: %w", id, err)
	}

	d.edgeCount.Add(^uint64(0))
	return nil
}

// getEdgesForNode retrieves edges for a node in the given direction. When
// loadProps is false, only the peer node ID and label are populated (from the
// adjacency value), avoiding N+1 GetEdge reads during traversals. When true,
// the full edge data is loaded via a separate GetEdge call. When labelFilter
// is non-empty, only edges with a matching label are returned. Self-loop
// edges are deduplicated when dir is Both.
func (d *DB) getEdgesForNode(id NodeID, dir Direction, loadProps bool, labelFilter string) ([]*Edge, error) {
	var edges []*Edge
	seenSelfLoops := make(map[EdgeID]bool)

	if dir == Outgoing || dir == Both {
		iter, err := d.db.NewIter(&pebble.IterOptions{
			LowerBound: adjOutPrefixKey(id),
			UpperBound: adjOutEndKey(id),
		})
		if err != nil {
			return nil, fmt.Errorf("graphengine: failed to create iterator for outgoing edges of node %d: %w", id, err)
		}

		for iter.First(); iter.Valid(); iter.Next() {
			_, edgeID := decodeAdjKey(iter.Key()[len(adjOutPrefix):])
			targetID, label := decodeAdjValue(iter.Value())
			if labelFilter != "" && label != labelFilter {
				continue
			}
			if dir == Both && targetID == id {
				seenSelfLoops[edgeID] = true
			}
			edge := &Edge{ID: edgeID, From: id, To: targetID, Label: label}
			if loadProps {
				full, err := d.GetEdge(edgeID)
				if err != nil {
					return nil, fmt.Errorf("graphengine: failed to load edge %d props: %w", edgeID, err)
				}
				edge = full
			}
			edges = append(edges, edge)
		}
		if err := iter.Error(); err != nil {
			iter.Close()
			return nil, fmt.Errorf("graphengine: outgoing edge iteration error for node %d: %w", id, err)
		}
		iter.Close()
	}

	if dir == Incoming || dir == Both {
		iter, err := d.db.NewIter(&pebble.IterOptions{
			LowerBound: adjInPrefixKey(id),
			UpperBound: adjInEndKey(id),
		})
		if err != nil {
			return nil, fmt.Errorf("graphengine: failed to create iterator for incoming edges of node %d: %w", id, err)
		}

		for iter.First(); iter.Valid(); iter.Next() {
			_, edgeID := decodeAdjKey(iter.Key()[len(adjInPrefix):])
			if dir == Both && seenSelfLoops[edgeID] {
				continue
			}
			sourceID, label := decodeAdjValue(iter.Value())
			if labelFilter != "" && label != labelFilter {
				continue
			}
			edge := &Edge{ID: edgeID, From: sourceID, To: id, Label: label}
			if loadProps {
				full, err := d.GetEdge(edgeID)
				if err != nil {
					return nil, fmt.Errorf("graphengine: failed to load edge %d props: %w", edgeID, err)
				}
				edge = full
			}
			edges = append(edges, edge)
		}
		if err := iter.Error(); err != nil {
			iter.Close()
			return nil, fmt.Errorf("graphengine: incoming edge iteration error for node %d: %w", id, err)
		}
		iter.Close()
	}

	return edges, nil
}

// OutEdges returns all outgoing edges of the given node.
func (d *DB) OutEdges(id NodeID) ([]*Edge, error) {
	return d.getEdgesForNode(id, Outgoing, true, "")
}

// InEdges returns all incoming edges of the given node.
func (d *DB) InEdges(id NodeID) ([]*Edge, error) {
	return d.getEdgesForNode(id, Incoming, true, "")
}

// OutEdgesByLabel returns outgoing edges of the given node filtered by label.
func (d *DB) OutEdgesByLabel(id NodeID, label string) ([]*Edge, error) {
	return d.getEdgesForNode(id, Outgoing, true, label)
}

// InEdgesByLabel returns incoming edges of the given node filtered by label.
func (d *DB) InEdgesByLabel(id NodeID, label string) ([]*Edge, error) {
	return d.getEdgesForNode(id, Incoming, true, label)
}

// GetEdges returns edges of the given node in the specified direction, optionally filtered by label.
func (d *DB) GetEdges(nodeID NodeID, dir Direction, labelFilter string) ([]*Edge, error) {
	return d.getEdgesForNode(nodeID, dir, true, labelFilter)
}

// GetAdjacentNodes returns the deduplicated set of nodes adjacent to the
// given node in the specified direction. For Both direction, outgoing edges
// contribute their To node and incoming edges contribute their From node.
// When labelFilter is non-empty, only edges with a matching label are
// considered.
// GetAdjacentNodes returns the deduplicated set of adjacent nodes in the specified direction.
func (d *DB) GetAdjacentNodes(nodeID NodeID, dir Direction) ([]*Node, error) {
	return d.getAdjacentNodesFiltered(nodeID, dir, "")
}

// GetAdjacentNodesByLabel returns adjacent nodes filtered by edge label.
func (d *DB) GetAdjacentNodesByLabel(nodeID NodeID, dir Direction, label string) ([]*Node, error) {
	return d.getAdjacentNodesFiltered(nodeID, dir, label)
}

func (d *DB) getAdjacentNodesFiltered(nodeID NodeID, dir Direction, labelFilter string) ([]*Node, error) {
	edges, err := d.getEdgesForNode(nodeID, dir, false, labelFilter)
	if err != nil {
		return nil, err
	}

	nodes := make([]*Node, 0, len(edges))
	seen := make(map[NodeID]bool)
	for _, e := range edges {
		targetID := e.To
		if dir == Incoming || (dir == Both && e.To == nodeID) {
			targetID = e.From
		}
		if seen[targetID] {
			continue
		}
		seen[targetID] = true
		n, err := d.GetNode(targetID)
		if err != nil {
			continue
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}
