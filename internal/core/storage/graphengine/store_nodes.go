package graphengine

import (
	"fmt"
)

func (d *DB) allocNodeID() NodeID {
	return NodeID(d.nextNodeID.Add(1))
}

func (d *DB) AddNode(labels []string, props Props) (NodeID, error) {
	if d.closed.Load() {
		return 0, fmt.Errorf("graphengine: database is closed")
	}

	if err := d.enforceUniqueConstraints(labels, props, 0); err != nil {
		return 0, err
	}

	id := d.allocNodeID()

	batch := d.backend.newBatch()
	defer batch.close()

	data, err := encodeNodeData(labels, props)
	if err != nil {
		return 0, err
	}
	batch.put(nodeKey(id), data)

	for _, label := range labels {
		batch.put(idxLabelKey(label, id), nil)
	}

	for k, v := range props {
		batch.put(idxPropKey(k, propIndexValue(v), id), nil)
	}

	if err := d.writeUniqueConstraintEntries(batch, labels, props, id); err != nil {
		return 0, fmt.Errorf("graphengine: failed to write unique constraint entries for new node: %w", err)
	}

	if err := batch.commit(); err != nil {
		return 0, fmt.Errorf("graphengine: failed to add node: %w", err)
	}

	d.nodeCount.Add(1)
	return id, nil
}

func (d *DB) AddNodeBatch(items []struct {
	Labels []string
	Props  Props
}) ([]NodeID, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	ids := make([]NodeID, len(items))
	batch := d.backend.newBatch()
	defer batch.close()

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
		batch.put(nodeKey(id), data)

		for _, label := range item.Labels {
			batch.put(idxLabelKey(label, id), nil)
		}

		for k, v := range item.Props {
			batch.put(idxPropKey(k, propIndexValue(v), id), nil)
		}

		if err := d.writeUniqueConstraintEntries(batch, item.Labels, item.Props, id); err != nil {
			return nil, fmt.Errorf("graphengine: failed to write unique constraint entries in batch: %w", err)
		}
	}

	if err := batch.commit(); err != nil {
		return nil, fmt.Errorf("graphengine: batch add failed: %w", err)
	}

	d.nodeCount.Add(uint64(len(items)))
	return ids, nil
}

func (d *DB) GetNode(id NodeID) (*Node, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	val, err := d.backend.get(nodeKey(id))
	if err != nil {
		return nil, fmt.Errorf("graphengine: failed to get node %d: %w", id, err)
	}
	if val == nil {
		return nil, fmt.Errorf("graphengine: node %d not found", id)
	}

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

func (d *DB) UpdateNode(id NodeID, props Props) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	val, err := d.backend.get(nodeKey(id))
	if err != nil {
		return fmt.Errorf("graphengine: failed to get node %d: %w", id, err)
	}
	if val == nil {
		return fmt.Errorf("graphengine: node %d not found", id)
	}

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

	batch := d.backend.newBatch()
	defer batch.close()

	for k, v := range props {
		if oldVal, ok := existingProps[k]; ok && propIndexValue(oldVal) != propIndexValue(v) {
			batch.del(idxPropKey(k, propIndexValue(oldVal), id))
			if constraintLabels, isConstrained := ucProps[k]; isConstrained {
				for _, cl := range constraintLabels {
					batch.del(ucKey(cl, k, oldVal))
					batch.put(ucKey(cl, k, v), encodeNodeID(id))
				}
			}
		} else if _, ok := existingProps[k]; !ok {
			if constraintLabels, isConstrained := ucProps[k]; isConstrained {
				for _, cl := range constraintLabels {
					batch.put(ucKey(cl, k, v), encodeNodeID(id))
				}
			}
		}
		existingProps[k] = v
		batch.put(idxPropKey(k, propIndexValue(v), id), nil)
	}

	data, err := encodeNodeData(labels, existingProps)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode updated node %d: %w", id, err)
	}
	batch.put(nodeKey(id), data)

	if err := batch.commit(); err != nil {
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
	if err := d.forEachUniqueConstraint(func(label, prop string) error {
		if labelSet[label] {
			result[prop] = append(result[prop], label)
		}
		return nil
	}); err != nil {
		return nil
	}
	return result
}

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

	batch := d.backend.newBatch()
	defer batch.close()

	for _, e := range edges {
		batch.del(edgeKey(e.ID))
		batch.del(adjOutKey(e.From, e.ID))
		batch.del(adjInKey(e.To, e.ID))
		if e.Label != "" {
			batch.del(idxEtypeKey(e.Label, e.ID))
		}
	}

	batch.del(nodeKey(id))

	for _, label := range node.Labels {
		batch.del(idxLabelKey(label, id))
	}

	for k, v := range node.Props {
		batch.del(idxPropKey(k, propIndexValue(v), id))
	}

	if err := d.deleteUniqueConstraintEntries(batch, node.Labels, node.Props); err != nil {
		return fmt.Errorf("graphengine: failed to delete unique constraint entries for node %d: %w", id, err)
	}

	if err := batch.commit(); err != nil {
		return fmt.Errorf("graphengine: failed to delete node %d: %w", id, err)
	}

	d.nodeCount.Add(^uint64(0))
	if len(edges) > 0 {
		d.edgeCount.Add(^uint64(uint64(len(edges) - 1)))
	}

	if d.Embeddings != nil {
		if _, err := d.Embeddings.Remove(id); err != nil {
			return fmt.Errorf("graphengine: failed to remove embedding for node %d: %w", id, err)
		}
	}

	return nil
}

func (d *DB) NodeExists(id NodeID) (bool, error) {
	if d.closed.Load() {
		return false, fmt.Errorf("graphengine: database is closed")
	}

	val, err := d.backend.get(nodeKey(id))
	if err != nil {
		return false, err
	}
	return val != nil, nil
}
