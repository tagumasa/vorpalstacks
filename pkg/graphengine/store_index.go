package graphengine

import (
	"bytes"
	"fmt"

	"github.com/cockroachdb/pebble/v2"
)

// FindByLabel returns all node IDs that carry the given label, using the
// label index. The scan is bounded by the label prefix and terminated early
// when keys no longer match.
func (d *DB) FindByLabel(label string) ([]NodeID, error) {
	var ids []NodeID

	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: idxLabelPrefixKey(label),
	})
	if err != nil {
		return nil, fmt.Errorf("graphengine: failed to create iterator for FindByLabel: %w", err)
	}

	prefix := idxLabelPrefixKey(label)
	for iter.First(); iter.Valid(); iter.Next() {
		if !hasPrefix(iter.Key(), prefix) {
			break
		}
		id := decodeNodeID(iter.Key()[len(prefix):])
		ids = append(ids, id)
	}
	if err := iter.Error(); err != nil {
		iter.Close()
		return nil, fmt.Errorf("graphengine: FindByLabel iteration error: %w", err)
	}
	iter.Close()

	return ids, nil
}

// ForEachNode iterates over every node in the graph, calling fn for each.
// Iteration is stopped early if fn returns a non-nil error.
func (d *DB) ForEachNode(fn func(*Node) error) error {
	return d.forEachNodeN(fn, 0)
}

// ForEachNodeN iterates over nodes, stopping after at most limit nodes
// have been visited. A limit of 0 means no limit.
func (d *DB) ForEachNodeN(fn func(*Node) error, limit int) error {
	return d.forEachNodeN(fn, limit)
}

func (d *DB) forEachNodeN(fn func(*Node) error, maxNodes int) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: nodePrefix,
		UpperBound: nodeEndKey(),
	})
	if err != nil {
		return fmt.Errorf("graphengine: failed to create iterator for forEachNodeN: %w", err)
	}
	defer iter.Close()

	count := 0
	for iter.First(); iter.Valid(); iter.Next() {
		if maxNodes > 0 && count >= maxNodes {
			break
		}

		id := decodeNodeID(iter.Key()[len(nodePrefix):])

		labels, props, err := decodeNodeData(iter.Value())
		if err != nil {
			return err
		}

		node := &Node{
			ID:     id,
			Labels: labels,
			Props:  props,
		}

		if err := fn(node); err != nil {
			return err
		}
		count++
	}

	if err := iter.Error(); err != nil {
		return fmt.Errorf("graphengine: node iteration error: %w", err)
	}

	return nil
}

// CreateIndex is a placeholder for future property index creation. Currently
// a no-op; property indexes are populated at write time by AddNode.
func (d *DB) CreateIndex(propName string) error {
	return nil
}

// UpdateEdge merges properties into an existing edge.
func (d *DB) UpdateEdge(id EdgeID, props Props) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	val, closer, err := d.db.Get(edgeKey(id))
	if err != nil {
		if err == pebble.ErrNotFound {
			return fmt.Errorf("graphengine: edge %d not found", id)
		}
		return fmt.Errorf("graphengine: failed to get edge %d: %w", id, err)
	}
	closer.Close()

	edge, err := decodeEdgeData(id, val)
	if err != nil {
		return fmt.Errorf("graphengine: failed to decode edge %d: %w", id, err)
	}

	if edge.Props == nil {
		edge.Props = make(Props)
	}
	for k, v := range props {
		edge.Props[k] = v
	}

	newData, err := encodeEdgeData(edge)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode updated edge %d: %w", id, err)
	}
	if err := d.db.Set(edgeKey(id), newData, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to update edge %d: %w", id, err)
	}

	return nil
}

// RemoveLabel removes a label from a node.
func (d *DB) RemoveLabel(id NodeID, label string) error {
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

	labels, props, err := decodeNodeData(val)
	if err != nil {
		return fmt.Errorf("graphengine: failed to decode node %d: %w", id, err)
	}

	newLabels := make([]string, 0, len(labels))
	for _, l := range labels {
		if l != label {
			newLabels = append(newLabels, l)
		}
	}

	batch := d.db.NewBatch()
	defer batch.Close()

	nodeData, err := encodeNodeData(newLabels, props)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode node after removing label: %w", err)
	}
	_ = batch.Set(nodeKey(id), nodeData, nil)
	_ = batch.Delete(idxLabelKey(label, id), nil)

	_ = d.forEachUniqueConstraint(func(ucLabel, prop string) error {
		if ucLabel == label {
			if val, ok := props[prop]; ok {
				_ = batch.Delete(ucKey(label, prop, val), nil)
			}
		}
		return nil
	})

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to remove label from node %d: %w", id, err)
	}

	return nil
}

// AddLabel adds a label to a node, enforcing unique constraints.
func (d *DB) AddLabel(id NodeID, label string) error {
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

	labels, props, err := decodeNodeData(val)
	if err != nil {
		return fmt.Errorf("graphengine: failed to decode node %d: %w", id, err)
	}

	for _, l := range labels {
		if l == label {
			return nil
		}
	}
	labels = append(labels, label)

	newLabels := append([]string{}, labels...)
	ucProps := d.getConstrainedProps(newLabels)

	batch := d.db.NewBatch()
	defer batch.Close()

	for prop, constraintLabels := range ucProps {
		for _, cl := range constraintLabels {
			if cl == label {
				if val, ok := props[prop]; ok {
					existingID, existingCloser, err := d.db.Get(ucKey(label, prop, val))
					if err == nil && existingID != nil {
						existingCloser.Close()
						existingNodeID := decodeNodeID(existingID)
						if existingNodeID != id {
							return fmt.Errorf("graphengine: unique constraint violation for %s.%s on node %d", label, prop, existingNodeID)
						}
					} else if err == nil {
						existingCloser.Close()
					}
					_ = batch.Set(ucKey(label, prop, val), encodeNodeID(id), nil)
				}
			}
		}
	}

	nodeData, err := encodeNodeData(labels, props)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode node after adding label: %w", err)
	}
	_ = batch.Set(nodeKey(id), nodeData, nil)
	_ = batch.Set(idxLabelKey(label, id), nil, nil)

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to add label to node %d: %w", id, err)
	}

	return nil
}

// RemoveProperty removes a property from a node.
func (d *DB) RemoveProperty(id NodeID, key string) error {
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

	labels, props, err := decodeNodeData(val)
	if err != nil {
		return fmt.Errorf("graphengine: failed to decode node %d: %w", id, err)
	}

	if _, exists := props[key]; !exists {
		return nil
	}

	oldValue := props[key]
	delete(props, key)

	ucProps := d.getConstrainedProps(labels)

	batch := d.db.NewBatch()
	defer batch.Close()

	nodeData, err := encodeNodeData(labels, props)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode node after removing property: %w", err)
	}
	_ = batch.Set(nodeKey(id), nodeData, nil)
	_ = batch.Delete(idxPropKey(key, propIndexValue(oldValue), id), nil)

	if constraintLabels, isConstrained := ucProps[key]; isConstrained {
		for _, cl := range constraintLabels {
			_ = batch.Delete(ucKey(cl, key, oldValue), nil)
		}
	}

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to remove property from node %d: %w", id, err)
	}

	return nil
}

// RemoveEdgeProperty removes a property from an edge.
func (d *DB) RemoveEdgeProperty(id EdgeID, key string) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	val, closer, err := d.db.Get(edgeKey(id))
	if err != nil {
		if err == pebble.ErrNotFound {
			return fmt.Errorf("graphengine: edge %d not found", id)
		}
		return fmt.Errorf("graphengine: failed to get edge %d: %w", id, err)
	}
	closer.Close()

	edge, err := decodeEdgeData(id, val)
	if err != nil {
		return fmt.Errorf("graphengine: failed to decode edge %d: %w", id, err)
	}

	if _, exists := edge.Props[key]; !exists {
		return nil
	}

	delete(edge.Props, key)

	edgeData, err := encodeEdgeData(edge)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode edge after removing property: %w", err)
	}
	if err := d.db.Set(edgeKey(id), edgeData, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to remove property from edge %d: %w", id, err)
	}

	return nil
}

// ForEachEdge iterates over every edge in the graph, calling fn for each.
func (d *DB) ForEachEdge(fn func(*Edge) error) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: edgePrefix,
		UpperBound: edgeEndKey(),
	})
	if err != nil {
		return fmt.Errorf("graphengine: failed to create iterator for ForEachEdge: %w", err)
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		id := EdgeID(decodeUint64(iter.Key()[len(edgePrefix):]))

		edge, err := decodeEdgeData(id, iter.Value())
		if err != nil {
			return err
		}

		if err := fn(edge); err != nil {
			return err
		}
	}

	if err := iter.Error(); err != nil {
		return fmt.Errorf("graphengine: edge iteration error: %w", err)
	}

	return nil
}

func (d *DB) forEachUniqueConstraint(fn func(label, prop string) error) error {
	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: ucPrefix,
		UpperBound: append([]byte(nil), ucPrefix[0]+1),
	})
	if err != nil {
		return fmt.Errorf("graphengine: failed to create iterator for forEachUniqueConstraint: %w", err)
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		if !hasPrefix(iter.Key(), ucPrefix) {
			break
		}

		suffix := iter.Key()[len(ucPrefix):]
		parts := splitNUL(suffix)
		if len(parts) != 2 {
			continue
		}

		label := string(parts[0])
		prop := string(parts[1])
		if err := fn(label, prop); err != nil {
			return err
		}
	}

	if err := iter.Error(); err != nil {
		return fmt.Errorf("graphengine: forEachUniqueConstraint iteration error: %w", err)
	}

	return nil
}

func splitNUL(data []byte) [][]byte {
	var parts [][]byte
	start := 0
	for i, b := range data {
		if b == 0x00 {
			parts = append(parts, data[start:i])
			start = i + 1
		}
	}
	parts = append(parts, data[start:])
	return parts
}

func (d *DB) enforceUniqueConstraints(labels []string, props Props, excludeID NodeID) error {
	labelSet := make(map[string]bool, len(labels))
	for _, l := range labels {
		labelSet[l] = true
	}

	return d.forEachUniqueConstraint(func(label, prop string) error {
		if !labelSet[label] {
			return nil
		}

		val, ok := props[prop]
		if !ok {
			return nil
		}

		prefix := ucKey(label, prop, val)
		iter, err := d.db.NewIter(&pebble.IterOptions{
			LowerBound: prefix,
		})
		if err != nil {
			return fmt.Errorf("graphengine: failed to create iterator for unique constraint check: %w", err)
		}
		defer iter.Close()

		for iter.First(); iter.Valid(); iter.Next() {
			if !hasPrefix(iter.Key(), prefix) {
				break
			}
			existingID := decodeNodeID(iter.Value())
			if existingID != excludeID {
				return fmt.Errorf("graphengine: unique constraint violation for %s.%s=%v (existing node %d)", label, prop, val, existingID)
			}
		}

		if err := iter.Error(); err != nil {
			return fmt.Errorf("graphengine: unique constraint iteration error: %w", err)
		}

		return nil
	})
}

func (d *DB) writeUniqueConstraintEntries(batch *pebble.Batch, labels []string, props Props, id NodeID) {
	labelSet := make(map[string]bool, len(labels))
	for _, l := range labels {
		labelSet[l] = true
	}

	_ = d.forEachUniqueConstraint(func(label, prop string) error {
		if !labelSet[label] {
			return nil
		}

		val, ok := props[prop]
		if !ok {
			return nil
		}

		key := ucKey(label, prop, val)
		_ = batch.Set(key, encodeNodeID(id), nil)
		return nil
	})
}

func (d *DB) deleteUniqueConstraintEntries(batch *pebble.Batch, labels []string, props Props) {
	labelSet := make(map[string]bool, len(labels))
	for _, l := range labels {
		labelSet[l] = true
	}

	_ = d.forEachUniqueConstraint(func(label, prop string) error {
		if !labelSet[label] {
			return nil
		}

		val, ok := props[prop]
		if !ok {
			return nil
		}

		key := ucKey(label, prop, val)
		_ = batch.Delete(key, nil)
		return nil
	})
}

// clearPrefixes lists all Pebble key prefixes in the graph schema. Used by
// Clear to delete each range in a single operation instead of iterating
// individual keys.
var clearPrefixes = [][]byte{
	nodePrefix,
	edgePrefix,
	adjOutPrefix,
	adjInPrefix,
	idxLabelPrefix,
	idxPropPrefix,
	idxEtypePrefix,
	ucPrefix,
	vecPrefix,
	[]byte("meta/"),
}

// Clear removes all data from the graph database.
func (d *DB) Clear() error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	batch := d.db.NewBatch()
	defer batch.Close()

	for _, prefix := range clearPrefixes {
		end := append([]byte(nil), prefix...)
		end = append(end, 0xFF)
		_ = batch.DeleteRange(prefix, end, nil)
	}

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to clear: %w", err)
	}

	d.nodeCount.Store(0)
	d.edgeCount.Store(0)
	d.nextNodeID.Store(1)
	d.nextEdgeID.Store(1)

	if d.Embeddings != nil {
		_ = d.Embeddings.Clear()
	}

	if err := d.persistCounters(); err != nil {
		return fmt.Errorf("graphengine: failed to persist counters after clear: %w", err)
	}

	return nil
}

// CreateUniqueConstraint creates a unique constraint for a label-property pair and backfills the index.
func (d *DB) CreateUniqueConstraint(label, prop string) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	key := append(append([]byte(nil), ucPrefix...), []byte(label)...)
	key = append(key, 0x00)
	key = append(key, []byte(prop)...)

	_, closer, err := d.db.Get(key)
	if err == nil {
		closer.Close()
		return fmt.Errorf("graphengine: unique constraint already exists for %s.%s", label, prop)
	}
	if err != pebble.ErrNotFound {
		return fmt.Errorf("graphengine: failed to check constraint: %w", err)
	}

	if err := d.db.Set(key, nil, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to create unique constraint: %w", err)
	}

	ids, err := d.FindByLabel(label)
	if err != nil {
		return fmt.Errorf("graphengine: failed to backfill unique constraint: %w", err)
	}

	batch := d.db.NewBatch()
	defer batch.Close()

	for _, id := range ids {
		node, err := d.GetNode(id)
		if err != nil {
			continue
		}
		val, ok := node.Props[prop]
		if !ok {
			continue
		}
		ucKey := ucKey(label, prop, val)
		existingVal, closer, err := d.db.Get(ucKey)
		if closer != nil {
			closer.Close()
		}
		if err == nil && existingVal != nil {
			existingNodeID := decodeNodeID(existingVal)
			if existingNodeID != id {
				batch.Close()
				return fmt.Errorf("graphengine: unique constraint violation during backfill for %s.%s=%v (existing node %d)", label, prop, val, existingNodeID)
			}
			continue
		}
		_ = batch.Set(ucKey, encodeNodeID(id), nil)
	}

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to backfill unique constraint index: %w", err)
	}

	return nil
}

// HasUniqueConstraint reports whether a unique constraint exists for the given label and property.
func (d *DB) HasUniqueConstraint(label, prop string) bool {
	key := append(append([]byte(nil), ucPrefix...), []byte(label)...)
	key = append(key, 0x00)
	key = append(key, []byte(prop)...)

	_, closer, err := d.db.Get(key)
	if err != nil {
		return false
	}
	closer.Close()
	return true
}

// FindByUniqueConstraint returns the node ID matching a unique constraint, or ErrNotFound.
func (d *DB) FindByUniqueConstraint(label, prop string, value interface{}) (NodeID, error) {
	prefix := ucKey(label, prop, value)

	var ids []NodeID
	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
	})
	if err != nil {
		return 0, fmt.Errorf("graphengine: failed to create iterator for FindByUniqueConstraint: %w", err)
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		if !hasPrefix(iter.Key(), prefix) {
			break
		}
		ids = append(ids, decodeNodeID(iter.Value()))
	}

	if err := iter.Error(); err != nil {
		return 0, fmt.Errorf("graphengine: FindByUniqueConstraint iteration error: %w", err)
	}

	if len(ids) == 0 {
		return 0, ErrNotFound
	}
	if len(ids) > 1 {
		return 0, fmt.Errorf("graphengine: unique constraint violation for %s.%s=%v", label, prop, value)
	}
	return ids[0], nil
}

// ScanEdgesByType iterates all edges of the given type, calling fn for each
// with the edge and its source/destination nodes. Edges or nodes that are
// concurrently deleted between the index scan and the individual Get calls
// are silently skipped (optimistic concurrency). The iteration stops early
// if fn returns a non-nil error.
func (d *DB) ScanEdgesByType(label string, fn func(edge *Edge, src, dst *Node) error) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	prefix := idxEtypePrefixKey(label)

	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
	})
	if err != nil {
		return fmt.Errorf("graphengine: failed to create iterator for ScanEdgesByType: %w", err)
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		if !hasPrefix(iter.Key(), prefix) {
			break
		}

		edgeID := decodeEdgeID(iter.Key()[len(prefix):])

		edge, err := d.GetEdge(edgeID)
		if err != nil {
			continue
		}

		src, err := d.GetNode(edge.From)
		if err != nil {
			continue
		}

		dst, err := d.GetNode(edge.To)
		if err != nil {
			continue
		}

		if err := fn(edge, src, dst); err != nil {
			return err
		}
	}

	if err := iter.Error(); err != nil {
		return fmt.Errorf("graphengine: ScanEdgesByType iteration error: %w", err)
	}

	return nil
}

// FindByProperty returns all node IDs matching the given property name and
// value, using the property index.
// FindByProperty returns all node IDs matching the given property name and value.
func (d *DB) FindByProperty(propName string, value interface{}) ([]NodeID, error) {
	valueStr := propIndexValue(value)
	prefix := idxPropPrefixKey(propName, valueStr)

	var ids []NodeID
	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
	})
	if err != nil {
		return nil, fmt.Errorf("graphengine: failed to create iterator for FindByProperty: %w", err)
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		if !hasPrefix(iter.Key(), prefix) {
			break
		}
		suffix := iter.Key()[len(prefix):]
		if len(suffix) >= 8 {
			id := decodeNodeID(suffix)
			ids = append(ids, id)
		}
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("graphengine: FindByProperty iteration error: %w", err)
	}

	return ids, nil
}

func hasPrefix(key, prefix []byte) bool {
	return bytes.HasPrefix(key, prefix)
}

func (d *DB) getRaw(key []byte) ([]byte, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}
	val, closer, err := d.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	defer closer.Close()
	result := make([]byte, len(val))
	copy(result, val)
	return result, nil
}

func (d *DB) setRaw(key, value []byte) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}
	return d.db.Set(key, value, pebble.NoSync)
}

func (d *DB) deleteRaw(key []byte) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}
	return d.db.Delete(key, pebble.NoSync)
}

func (d *DB) deleteRangeRaw(prefix []byte) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}
	end := append([]byte(nil), prefix...)
	end = append(end, 0xFF)
	return d.db.DeleteRange(prefix, end, pebble.NoSync)
}

func (d *DB) iteratePrefixRaw(prefix []byte, fn func(key, value []byte) error) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}
	end := append([]byte(nil), prefix...)
	end = append(end, 0xFF)
	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: end,
	})
	if err != nil {
		return err
	}
	defer iter.Close()
	for iter.First(); iter.Valid(); iter.Next() {
		key := make([]byte, len(iter.Key()))
		copy(key, iter.Key())
		val := make([]byte, len(iter.Value()))
		copy(val, iter.Value())
		if err := fn(key, val); err != nil {
			return err
		}
	}
	return iter.Error()
}
