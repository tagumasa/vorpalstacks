package graphengine

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (d *DB) FindByLabel(label string) ([]NodeID, error) {
	var ids []NodeID
	prefix := idxLabelPrefixKey(label)

	it, err := d.backend.newIter(prefix, nil)
	if err != nil {
		return nil, fmt.Errorf("graphengine: failed to create iterator for FindByLabel: %w", err)
	}

	for it.first(); it.valid(); it.next() {
		if !hasPrefix(it.key(), prefix) {
			break
		}
		id := decodeNodeID(it.key()[len(prefix):])
		ids = append(ids, id)
	}
	if err := it.err(); err != nil {
		it.close()
		return nil, fmt.Errorf("graphengine: FindByLabel iteration error: %w", err)
	}
	it.close()

	return ids, nil
}

func (d *DB) ForEachNode(fn func(*Node) error) error {
	return d.forEachNodeN(fn, 0)
}

func (d *DB) ForEachNodeN(fn func(*Node) error, limit int) error {
	return d.forEachNodeN(fn, limit)
}

func (d *DB) forEachNodeN(fn func(*Node) error, maxNodes int) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	it, err := d.backend.newIter(nodePrefix, nodeEndKey())
	if err != nil {
		return fmt.Errorf("graphengine: failed to create iterator for forEachNodeN: %w", err)
	}
	defer it.close()

	count := 0
	for it.first(); it.valid(); it.next() {
		if maxNodes > 0 && count >= maxNodes {
			break
		}

		id := decodeNodeID(it.key()[len(nodePrefix):])

		labels, props, err := decodeNodeData(it.value())
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

	if err := it.err(); err != nil {
		return fmt.Errorf("graphengine: node iteration error: %w", err)
	}

	return nil
}

// IndexInfo describes a registered property index on a label.
type IndexInfo struct {
	Label    string
	Property string
}

// idxMetaKey returns the storage key for an index metadata entry.
func idxMetaKey(label, prop string) []byte {
	key := append(append([]byte(nil), idxMetaPrefix...), []byte(label)...)
	key = append(key, 0x00)
	return append(key, []byte(prop)...)
}

// CreateIndex registers a property index for the given label and property.
// Property data is already auto-indexed via idx_prop/ keys; this method
// records the index in metadata so that SHOW INDEXES can list it.
func (d *DB) CreateIndex(label, prop string) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	key := idxMetaKey(label, prop)
	val, err := d.backend.get(key)
	if err != nil {
		return fmt.Errorf("graphengine: failed to check index: %w", err)
	}
	if val != nil {
		return fmt.Errorf("graphengine: index already exists for %s.%s", label, prop)
	}

	meta := map[string]string{"label": label, "property": prop}
	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("graphengine: failed to marshal index meta for %s.%s: %w", label, prop, err)
	}
	if err := d.backend.set(key, data); err != nil {
		return fmt.Errorf("graphengine: failed to create index for %s.%s: %w", label, prop, err)
	}

	return nil
}

// ShowIndexes returns all registered property indexes.
func (d *DB) ShowIndexes() ([]IndexInfo, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	var indexes []IndexInfo
	it, err := d.backend.newIter(idxMetaPrefix, nil)
	if err != nil {
		return nil, fmt.Errorf("graphengine: failed to create iterator for ShowIndexes: %w", err)
	}
	defer it.close()

	for it.first(); it.valid(); it.next() {
		if !hasPrefix(it.key(), idxMetaPrefix) {
			break
		}
		var meta map[string]string
		if err := json.Unmarshal(it.value(), &meta); err == nil {
			indexes = append(indexes, IndexInfo{Label: meta["label"], Property: meta["property"]})
		}
	}

	if err := it.err(); err != nil {
		return nil, fmt.Errorf("graphengine: ShowIndexes iteration error: %w", err)
	}

	return indexes, nil
}

// DropIndex removes a registered property index.
func (d *DB) DropIndex(label, prop string) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	key := idxMetaKey(label, prop)
	val, err := d.backend.get(key)
	if err != nil {
		return fmt.Errorf("graphengine: failed to check index: %w", err)
	}
	if val == nil {
		return fmt.Errorf("graphengine: index not found for %s.%s", label, prop)
	}

	d.backend.delete(key)

	return nil
}

func (d *DB) UpdateEdge(id EdgeID, props Props) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	val, err := d.backend.get(edgeKey(id))
	if err != nil {
		return fmt.Errorf("graphengine: failed to get edge %d: %w", id, err)
	}
	if val == nil {
		return fmt.Errorf("graphengine: edge %d not found", id)
	}

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
	if err := d.backend.set(edgeKey(id), newData); err != nil {
		return fmt.Errorf("graphengine: failed to update edge %d: %w", id, err)
	}

	return nil
}

func (d *DB) RemoveLabel(id NodeID, label string) error {
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

	batch := d.backend.newBatch()
	defer batch.close()

	nodeData, err := encodeNodeData(newLabels, props)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode node after removing label: %w", err)
	}
	batch.put(nodeKey(id), nodeData)
	batch.del(idxLabelKey(label, id))

	if err := d.forEachUniqueConstraint(func(ucLabel, prop string) error {
		if ucLabel == label {
			if val, ok := props[prop]; ok {
				batch.del(ucKey(label, prop, val))
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("graphengine: failed to iterate unique constraints when removing label from node %d: %w", id, err)
	}

	if err := batch.commit(); err != nil {
		return fmt.Errorf("graphengine: failed to remove label from node %d: %w", id, err)
	}

	return nil
}

func (d *DB) AddLabel(id NodeID, label string) error {
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

	batch := d.backend.newBatch()
	defer batch.close()

	for prop, constraintLabels := range ucProps {
		for _, cl := range constraintLabels {
			if cl == label {
				if val, ok := props[prop]; ok {
					existingVal, err := d.backend.get(ucKey(label, prop, val))
					if err != nil {
						return err
					}
					if existingVal != nil {
						existingNodeID := decodeNodeID(existingVal)
						if existingNodeID != id {
							return fmt.Errorf("graphengine: unique constraint violation for %s.%s on node %d", label, prop, existingNodeID)
						}
					}
					batch.put(ucKey(label, prop, val), encodeNodeID(id))
				}
			}
		}
	}

	nodeData, err := encodeNodeData(labels, props)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode node after adding label: %w", err)
	}
	batch.put(nodeKey(id), nodeData)
	batch.put(idxLabelKey(label, id), nil)

	if err := batch.commit(); err != nil {
		return fmt.Errorf("graphengine: failed to add label to node %d: %w", id, err)
	}

	return nil
}

func (d *DB) RemoveProperty(id NodeID, key string) error {
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

	batch := d.backend.newBatch()
	defer batch.close()

	nodeData, err := encodeNodeData(labels, props)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode node after removing property: %w", err)
	}
	batch.put(nodeKey(id), nodeData)
	batch.del(idxPropKey(key, propIndexValue(oldValue), id))

	if constraintLabels, isConstrained := ucProps[key]; isConstrained {
		for _, cl := range constraintLabels {
			batch.del(ucKey(cl, key, oldValue))
		}
	}

	if err := batch.commit(); err != nil {
		return fmt.Errorf("graphengine: failed to remove property from node %d: %w", id, err)
	}

	return nil
}

func (d *DB) RemoveEdgeProperty(id EdgeID, key string) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	val, err := d.backend.get(edgeKey(id))
	if err != nil {
		return fmt.Errorf("graphengine: failed to get edge %d: %w", id, err)
	}
	if val == nil {
		return fmt.Errorf("graphengine: edge %d not found", id)
	}

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
	if err := d.backend.set(edgeKey(id), edgeData); err != nil {
		return fmt.Errorf("graphengine: failed to remove property from edge %d: %w", id, err)
	}

	return nil
}

func (d *DB) ForEachEdge(fn func(*Edge) error) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	it, err := d.backend.newIter(edgePrefix, edgeEndKey())
	if err != nil {
		return fmt.Errorf("graphengine: failed to create iterator for ForEachEdge: %w", err)
	}
	defer it.close()

	for it.first(); it.valid(); it.next() {
		id := EdgeID(decodeUint64(it.key()[len(edgePrefix):]))

		edge, err := decodeEdgeData(id, it.value())
		if err != nil {
			return err
		}

		if err := fn(edge); err != nil {
			return err
		}
	}

	if err := it.err(); err != nil {
		return fmt.Errorf("graphengine: edge iteration error: %w", err)
	}

	return nil
}

func (d *DB) forEachUniqueConstraint(fn func(label, prop string) error) error {
	ucEnd := append([]byte(nil), ucPrefix[0]+1)
	it, err := d.backend.newIter(ucPrefix, ucEnd)
	if err != nil {
		return fmt.Errorf("graphengine: failed to create iterator for forEachUniqueConstraint: %w", err)
	}
	defer it.close()

	for it.first(); it.valid(); it.next() {
		if !hasPrefix(it.key(), ucPrefix) {
			break
		}

		suffix := it.key()[len(ucPrefix):]
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

	if err := it.err(); err != nil {
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
		it, err := d.backend.newIter(prefix, nil)
		if err != nil {
			return fmt.Errorf("graphengine: failed to create iterator for unique constraint check: %w", err)
		}
		defer it.close()

		for it.first(); it.valid(); it.next() {
			if !hasPrefix(it.key(), prefix) {
				break
			}
			existingID := decodeNodeID(it.value())
			if existingID != excludeID {
				return fmt.Errorf("graphengine: unique constraint violation for %s.%s=%v (existing node %d)", label, prop, val, existingID)
			}
		}

		if err := it.err(); err != nil {
			return fmt.Errorf("graphengine: unique constraint iteration error: %w", err)
		}

		return nil
	})
}

func (d *DB) writeUniqueConstraintEntries(batch kvBatch, labels []string, props Props, id NodeID) error {
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

		key := ucKey(label, prop, val)
		batch.put(key, encodeNodeID(id))
		return nil
	})
}

func (d *DB) deleteUniqueConstraintEntries(batch kvBatch, labels []string, props Props) error {
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

		key := ucKey(label, prop, val)
		batch.del(key)
		return nil
	})
}

var clearPrefixes = [][]byte{
	nodePrefix,
	edgePrefix,
	adjOutPrefix,
	adjInPrefix,
	idxLabelPrefix,
	idxPropPrefix,
	idxEtypePrefix,
	idxMetaPrefix,
	ucPrefix,
	vecPrefix,
	[]byte("meta/"),
}

func (d *DB) Clear() error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	batch := d.backend.newBatch()
	defer batch.close()

	for _, prefix := range clearPrefixes {
		end := append([]byte(nil), prefix...)
		end = append(end, 0xFF)
		batch.deleteRange(prefix, end)
	}

	if err := batch.commit(); err != nil {
		return fmt.Errorf("graphengine: failed to clear: %w", err)
	}

	d.nodeCount.Store(0)
	d.edgeCount.Store(0)
	d.nextNodeID.Store(1)
	d.nextEdgeID.Store(1)

	if d.Embeddings != nil {
		if err := d.Embeddings.Clear(); err != nil {
			return fmt.Errorf("graphengine: failed to clear embeddings: %w", err)
		}
	}

	if err := d.persistCounters(); err != nil {
		return fmt.Errorf("graphengine: failed to persist counters after clear: %w", err)
	}

	return nil
}

func (d *DB) CreateUniqueConstraint(label, prop string) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	key := append(append([]byte(nil), ucPrefix...), []byte(label)...)
	key = append(key, 0x00)
	key = append(key, []byte(prop)...)

	val, err := d.backend.get(key)
	if err != nil {
		return fmt.Errorf("graphengine: failed to check constraint: %w", err)
	}
	if val != nil {
		return fmt.Errorf("graphengine: unique constraint already exists for %s.%s", label, prop)
	}

	if err := d.backend.set(key, nil); err != nil {
		return fmt.Errorf("graphengine: failed to create unique constraint: %w", err)
	}

	ids, err := d.FindByLabel(label)
	if err != nil {
		return fmt.Errorf("graphengine: failed to backfill unique constraint: %w", err)
	}

	batch := d.backend.newBatch()
	defer batch.close()

	for _, id := range ids {
		node, err := d.GetNode(id)
		if err != nil {
			continue
		}
		val, ok := node.Props[prop]
		if !ok {
			continue
		}
		ucKeyBytes := ucKey(label, prop, val)
		existingVal, err := d.backend.get(ucKeyBytes)
		if err != nil {
			return err
		}
		if existingVal != nil {
			existingNodeID := decodeNodeID(existingVal)
			if existingNodeID != id {
				return fmt.Errorf("graphengine: unique constraint violation during backfill for %s.%s=%v (existing node %d)", label, prop, val, existingNodeID)
			}
			continue
		}
		batch.put(ucKeyBytes, encodeNodeID(id))
	}

	if err := batch.commit(); err != nil {
		return fmt.Errorf("graphengine: failed to backfill unique constraint index: %w", err)
	}

	return nil
}

func (d *DB) HasUniqueConstraint(label, prop string) bool {
	key := append(append([]byte(nil), ucPrefix...), []byte(label)...)
	key = append(key, 0x00)
	key = append(key, []byte(prop)...)

	val, err := d.backend.get(key)
	if err != nil {
		return false
	}
	return val != nil
}

// ConstraintInfo describes a registered unique constraint.
type ConstraintInfo struct {
	Label    string
	Property string
}

// ShowConstraints returns all registered unique constraints.
func (d *DB) ShowConstraints() ([]ConstraintInfo, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	var constraints []ConstraintInfo
	err := d.forEachUniqueConstraint(func(label, prop string) error {
		constraints = append(constraints, ConstraintInfo{Label: label, Property: prop})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return constraints, nil
}

// DropUniqueConstraint removes a unique constraint and its associated index entries.
func (d *DB) DropUniqueConstraint(label, prop string) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	key := append(append([]byte(nil), ucPrefix...), []byte(label)...)
	key = append(key, 0x00)
	key = append(key, []byte(prop)...)

	val, err := d.backend.get(key)
	if err != nil {
		return fmt.Errorf("graphengine: failed to check constraint: %w", err)
	}
	if val == nil {
		return fmt.Errorf("graphengine: unique constraint not found for %s.%s", label, prop)
	}

	d.backend.delete(key)

	return nil
}

func (d *DB) FindByUniqueConstraint(label, prop string, value interface{}) (NodeID, error) {
	prefix := ucKey(label, prop, value)

	var ids []NodeID
	it, err := d.backend.newIter(prefix, nil)
	if err != nil {
		return 0, fmt.Errorf("graphengine: failed to create iterator for FindByUniqueConstraint: %w", err)
	}
	defer it.close()

	for it.first(); it.valid(); it.next() {
		if !hasPrefix(it.key(), prefix) {
			break
		}
		ids = append(ids, decodeNodeID(it.value()))
	}

	if err := it.err(); err != nil {
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

func (d *DB) ScanEdgesByType(label string, fn func(edge *Edge, src, dst *Node) error) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	prefix := idxEtypePrefixKey(label)

	it, err := d.backend.newIter(prefix, nil)
	if err != nil {
		return fmt.Errorf("graphengine: failed to create iterator for ScanEdgesByType: %w", err)
	}
	defer it.close()

	for it.first(); it.valid(); it.next() {
		if !hasPrefix(it.key(), prefix) {
			break
		}

		edgeID := decodeEdgeID(it.key()[len(prefix):])

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

	if err := it.err(); err != nil {
		return fmt.Errorf("graphengine: ScanEdgesByType iteration error: %w", err)
	}

	return nil
}

func (d *DB) FindByProperty(propName string, value interface{}) ([]NodeID, error) {
	valueStr := propIndexValue(value)
	prefix := idxPropPrefixKey(propName, valueStr)

	var ids []NodeID
	it, err := d.backend.newIter(prefix, nil)
	if err != nil {
		return nil, fmt.Errorf("graphengine: failed to create iterator for FindByProperty: %w", err)
	}
	defer it.close()

	for it.first(); it.valid(); it.next() {
		if !hasPrefix(it.key(), prefix) {
			break
		}
		suffix := it.key()[len(prefix):]
		if len(suffix) >= 8 {
			id := decodeNodeID(suffix)
			ids = append(ids, id)
		}
	}

	if err := it.err(); err != nil {
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
	return d.backend.get(key)
}

func (d *DB) setRaw(key, value []byte) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}
	return d.backend.set(key, value)
}

func (d *DB) deleteRaw(key []byte) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}
	return d.backend.delete(key)
}

func (d *DB) deleteRangeRaw(prefix []byte) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}
	end := append([]byte(nil), prefix...)
	end = append(end, 0xFF)
	return d.backend.deleteRange(prefix, end)
}

func (d *DB) iteratePrefixRaw(prefix []byte, fn func(key, value []byte) error) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}
	end := append([]byte(nil), prefix...)
	end = append(end, 0xFF)
	it, err := d.backend.newIter(prefix, end)
	if err != nil {
		return err
	}
	defer it.close()
	for it.first(); it.valid(); it.next() {
		key := make([]byte, len(it.key()))
		copy(key, it.key())
		val := make([]byte, len(it.value()))
		copy(val, it.value())
		if err := fn(key, val); err != nil {
			return err
		}
	}
	return it.err()
}
