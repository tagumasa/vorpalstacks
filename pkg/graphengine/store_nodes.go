package graphengine

import (
	"fmt"

	"github.com/cockroachdb/pebble/v2"
)

func (d *DB) allocNodeID() NodeID {
	return NodeID(d.nextNodeID.Add(1))
}

// AddNode creates a new node with the given labels and properties, allocates
// a monotonically increasing ID, and updates the label index. If the Pebble
// write fails the ID is consumed (a gap in the sequence); this is acceptable
// as IDs carry no semantic meaning. Unique constraints are checked and
// enforced before the write.
func (d *DB) AddNode(labels []string, props Props) (NodeID, error) {
	if d.closed.Load() {
		return 0, fmt.Errorf("graphengine: database is closed")
	}

	if err := d.enforceUniqueConstraints(labels, props, 0); err != nil {
		return 0, err
	}

	id := d.allocNodeID()

	batch := d.db.NewBatch()
	defer batch.Close()

	data, err := encodeNodeData(labels, props)
	if err != nil {
		return 0, err
	}
	_ = batch.Set(nodeKey(id), data, nil)

	for _, label := range labels {
		_ = batch.Set(idxLabelKey(label, id), nil, nil)
	}

	for k, v := range props {
		_ = batch.Set(idxPropKey(k, propIndexValue(v), id), nil, nil)
	}

	d.writeUniqueConstraintEntries(batch, labels, props, id)

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return 0, fmt.Errorf("graphengine: failed to add node: %w", err)
	}

	d.nodeCount.Add(1)
	return id, nil
}

// AddNodeBatch creates multiple nodes in a single Pebble batch. IDs are
// allocated atomically before the batch is applied; on failure the allocated
// IDs are consumed but not stored, leaving a harmless gap.
func (d *DB) AddNodeBatch(items []struct {
	Labels []string
	Props  Props
}) ([]NodeID, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	ids := make([]NodeID, len(items))
	batch := d.db.NewBatch()
	defer batch.Close()

	seenUC := make(map[string]struct{})

	for i, item := range items {
		ucProps := d.getConstrainedProps(item.Labels)
		labelSet := make(map[string]bool, len(item.Labels))
		for _, l := range item.Labels {
			labelSet[l] = true
		}
		for prop, val := range item.Props {
			if cl, ok := ucProps[prop]; ok {
				for _, label := range cl {
					if !labelSet[label] {
						continue
					}
					ucKeyStr := string(ucKey(label, prop, val))
					if _, dup := seenUC[ucKeyStr]; dup {
						return nil, fmt.Errorf("graphengine: unique constraint violation for %s.%s=%v (duplicate within batch)", label, prop, val)
					}
					seenUC[ucKeyStr] = struct{}{}
				}
			}
		}

		if err := d.enforceUniqueConstraints(item.Labels, item.Props, 0); err != nil {
			return nil, err
		}

		id := d.allocNodeID()
		ids[i] = id

		data, err := encodeNodeData(item.Labels, item.Props)
		if err != nil {
			return nil, err
		}
		_ = batch.Set(nodeKey(id), data, nil)

		for _, label := range item.Labels {
			_ = batch.Set(idxLabelKey(label, id), nil, nil)
		}

		for k, v := range item.Props {
			_ = batch.Set(idxPropKey(k, propIndexValue(v), id), nil, nil)
		}

		d.writeUniqueConstraintEntries(batch, item.Labels, item.Props, id)
	}

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return nil, fmt.Errorf("graphengine: batch add failed: %w", err)
	}

	d.nodeCount.Add(uint64(len(items)))
	return ids, nil
}

// GetNode retrieves a node by its ID. Returns an error if the node does not exist.
func (d *DB) GetNode(id NodeID) (*Node, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	val, closer, err := d.db.Get(nodeKey(id))
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, fmt.Errorf("graphengine: node %d not found", id)
		}
		return nil, fmt.Errorf("graphengine: failed to get node %d: %w", id, err)
	}
	defer closer.Close()

	labels, props, err := decodeNodeData(val)
	if err != nil {
		return nil, fmt.Errorf("graphengine: failed to decode node %d: %w", id, err)
	}

	return &Node{
		ID:     id,
		Labels: labels,
		Props:  props,
	}, nil
}

// UpdateNode merges the given properties into the existing node. Existing
// keys are overwritten; new keys are added. Labels are preserved unchanged.
// Unique constraint index entries are updated for changed property values.
func (d *DB) UpdateNode(id NodeID, props Props) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	val, closer, err := d.db.Get(nodeKey(id))
	if err != nil {
		if err == pebble.ErrNotFound {
			return fmt.Errorf("graphengine: node %d not found", id)
		}
		return fmt.Errorf("graphengine: failed to get node %d: %w", id, err)
	}
	closer.Close()

	labels, existingProps, err := decodeNodeData(val)
	if err != nil {
		return fmt.Errorf("graphengine: failed to decode node %d: %w", id, err)
	}

	if existingProps == nil {
		existingProps = make(Props)
	}

	ucProps := d.getConstrainedProps(labels)

	if err := d.enforceUniqueConstraints(labels, props, id); err != nil {
		return err
	}

	batch := d.db.NewBatch()
	defer batch.Close()

	for k, v := range props {
		if oldVal, ok := existingProps[k]; ok && propIndexValue(oldVal) != propIndexValue(v) {
			_ = batch.Delete(idxPropKey(k, propIndexValue(oldVal), id), nil)
			if constraintLabels, isConstrained := ucProps[k]; isConstrained {
				for _, cl := range constraintLabels {
					_ = batch.Delete(ucKey(cl, k, oldVal), nil)
					_ = batch.Set(ucKey(cl, k, v), encodeNodeID(id), nil)
				}
			}
		} else if _, ok := existingProps[k]; !ok {
			if constraintLabels, isConstrained := ucProps[k]; isConstrained {
				for _, cl := range constraintLabels {
					_ = batch.Set(ucKey(cl, k, v), encodeNodeID(id), nil)
				}
			}
		}
		existingProps[k] = v
		_ = batch.Set(idxPropKey(k, propIndexValue(v), id), nil, nil)
	}

	data, err := encodeNodeData(labels, existingProps)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode updated node %d: %w", id, err)
	}
	_ = batch.Set(nodeKey(id), data, nil)

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to update node %d: %w", id, err)
	}

	return nil
}

func (d *DB) getConstrainedProps(labels []string) map[string][]string {
	result := make(map[string][]string)
	labelSet := make(map[string]bool, len(labels))
	for _, l := range labels {
		labelSet[l] = true
	}
	_ = d.forEachUniqueConstraint(func(label, prop string) error {
		if labelSet[label] {
			result[prop] = append(result[prop], label)
		}
		return nil
	})
	return result
}

// DeleteNode removes a node and all of its incident edges (both outgoing and
// incoming), along with their adjacency entries and label index entries. The
// in-memory counters are decremented accordingly.
func (d *DB) DeleteNode(id NodeID) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	node, err := d.GetNode(id)
	if err != nil {
		return err
	}

	edges, err := d.getEdgesForNode(id, Both, false, "")
	if err != nil {
		return err
	}

	batch := d.db.NewBatch()
	defer batch.Close()

	for _, e := range edges {
		_ = batch.Delete(edgeKey(e.ID), nil)
		_ = batch.Delete(adjOutKey(e.From, e.ID), nil)
		_ = batch.Delete(adjInKey(e.To, e.ID), nil)
		if e.Label != "" {
			_ = batch.Delete(idxEtypeKey(e.Label, e.ID), nil)
		}
	}

	_ = batch.Delete(nodeKey(id), nil)

	for _, label := range node.Labels {
		_ = batch.Delete(idxLabelKey(label, id), nil)
	}

	for k, v := range node.Props {
		_ = batch.Delete(idxPropKey(k, propIndexValue(v), id), nil)
	}

	d.deleteUniqueConstraintEntries(batch, node.Labels, node.Props)

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to delete node %d: %w", id, err)
	}

	d.nodeCount.Add(^uint64(0))
	if len(edges) > 0 {
		d.edgeCount.Add(^uint64(uint64(len(edges) - 1)))
	}

	if d.Embeddings != nil {
		_, _ = d.Embeddings.Remove(id)
	}

	return nil
}

// NodeExists reports whether a node with the given ID exists.
func (d *DB) NodeExists(id NodeID) (bool, error) {
	if d.closed.Load() {
		return false, fmt.Errorf("graphengine: database is closed")
	}

	_, closer, err := d.db.Get(nodeKey(id))
	if err != nil {
		if err == pebble.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	closer.Close()
	return true, nil
}
