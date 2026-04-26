package graphengine

import (
	"fmt"
)

func (d *DB) allocEdgeID() EdgeID {
	return EdgeID(d.nextEdgeID.Add(1))
}

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

	batch := d.backend.newBatch()
	defer batch.close()

	propsBytes := marshalProps(props)

	edgeData, err := encodeEdgeDataRaw(edge, propsBytes)
	if err != nil {
		return 0, err
	}
	batch.put(edgeKey(id), edgeData)
	batch.put(adjOutKey(from, id), encodeAdjValue(to, label))
	batch.put(adjInKey(to, id), encodeAdjValue(from, label))
	if label != "" {
		batch.put(idxEtypeKey(label, id), nil)
	}

	if err := batch.commit(); err != nil {
		return 0, fmt.Errorf("graphengine: failed to add edge: %w", err)
	}

	d.edgeCount.Add(1)
	return id, nil
}

func (d *DB) AddEdgeBatch(edges []Edge) ([]EdgeID, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	ids := make([]EdgeID, len(edges))
	batch := d.backend.newBatch()
	defer batch.close()

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
		batch.put(edgeKey(id), edgeData)
		batch.put(adjOutKey(edge.From, id), encodeAdjValue(edge.To, edge.Label))
		batch.put(adjInKey(edge.To, id), encodeAdjValue(edge.From, edge.Label))
		if edge.Label != "" {
			batch.put(idxEtypeKey(edge.Label, id), nil)
		}
	}

	if err := batch.commit(); err != nil {
		return nil, fmt.Errorf("graphengine: batch edge add failed: %w", err)
	}

	d.edgeCount.Add(uint64(len(edges)))
	return ids, nil
}

func (d *DB) GetEdge(id EdgeID) (*Edge, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	val, err := d.backend.get(edgeKey(id))
	if err != nil {
		return nil, fmt.Errorf("graphengine: failed to get edge %d: %w", id, err)
	}
	if val == nil {
		return nil, fmt.Errorf("graphengine: edge %d not found", id)
	}

	return decodeEdgeData(id, val)
}

func (d *DB) DeleteEdge(id EdgeID) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	edge, err := d.GetEdge(id)
	if err != nil {
		return err
	}

	batch := d.backend.newBatch()
	defer batch.close()

	batch.del(edgeKey(id))
	batch.del(adjOutKey(edge.From, id))
	batch.del(adjInKey(edge.To, id))
	if edge.Label != "" {
		batch.del(idxEtypeKey(edge.Label, id))
	}

	if err := batch.commit(); err != nil {
		return fmt.Errorf("graphengine: failed to delete edge %d: %w", id, err)
	}

	d.edgeCount.Add(^uint64(0))
	return nil
}

func (d *DB) getEdgesForNode(id NodeID, dir Direction, loadProps bool, labelFilter string) ([]*Edge, error) {
	var edges []*Edge
	seenSelfLoops := make(map[EdgeID]bool)

	if dir == Outgoing || dir == Both {
		it, err := d.backend.newIter(adjOutPrefixKey(id), adjOutEndKey(id))
		if err != nil {
			return nil, fmt.Errorf("graphengine: failed to create iterator for outgoing edges of node %d: %w", id, err)
		}

		for it.first(); it.valid(); it.next() {
			_, edgeID := decodeAdjKey(it.key()[len(adjOutPrefix):])
			targetID, label := decodeAdjValue(it.value())
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
					it.close()
					return nil, fmt.Errorf("graphengine: failed to load edge %d props: %w", edgeID, err)
				}
				edge = full
			}
			edges = append(edges, edge)
		}
		if err := it.err(); err != nil {
			it.close()
			return nil, fmt.Errorf("graphengine: outgoing edge iteration error for node %d: %w", id, err)
		}
		it.close()
	}

	if dir == Incoming || dir == Both {
		it, err := d.backend.newIter(adjInPrefixKey(id), adjInEndKey(id))
		if err != nil {
			return nil, fmt.Errorf("graphengine: failed to create iterator for incoming edges of node %d: %w", id, err)
		}

		for it.first(); it.valid(); it.next() {
			_, edgeID := decodeAdjKey(it.key()[len(adjInPrefix):])
			if dir == Both && seenSelfLoops[edgeID] {
				continue
			}
			sourceID, label := decodeAdjValue(it.value())
			if labelFilter != "" && label != labelFilter {
				continue
			}
			edge := &Edge{ID: edgeID, From: sourceID, To: id, Label: label}
			if loadProps {
				full, err := d.GetEdge(edgeID)
				if err != nil {
					it.close()
					return nil, fmt.Errorf("graphengine: failed to load edge %d props: %w", edgeID, err)
				}
				edge = full
			}
			edges = append(edges, edge)
		}
		if err := it.err(); err != nil {
			it.close()
			return nil, fmt.Errorf("graphengine: incoming edge iteration error for node %d: %w", id, err)
		}
		it.close()
	}

	return edges, nil
}

func (d *DB) OutEdges(id NodeID) ([]*Edge, error) {
	return d.getEdgesForNode(id, Outgoing, true, "")
}

func (d *DB) InEdges(id NodeID) ([]*Edge, error) {
	return d.getEdgesForNode(id, Incoming, true, "")
}

func (d *DB) OutEdgesByLabel(id NodeID, label string) ([]*Edge, error) {
	return d.getEdgesForNode(id, Outgoing, true, label)
}

func (d *DB) InEdgesByLabel(id NodeID, label string) ([]*Edge, error) {
	return d.getEdgesForNode(id, Incoming, true, label)
}

func (d *DB) GetEdges(nodeID NodeID, dir Direction, labelFilter string) ([]*Edge, error) {
	return d.getEdgesForNode(nodeID, dir, true, labelFilter)
}

func (d *DB) GetAdjacentNodes(nodeID NodeID, dir Direction) ([]*Node, error) {
	return d.getAdjacentNodesFiltered(nodeID, dir, "")
}

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
