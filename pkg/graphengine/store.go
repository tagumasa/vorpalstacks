package graphengine

// store.go provides the Pebble-backed graph storage layer. Nodes, edges, and
// adjacency lists are stored as typed binary keys/values within a single
// Pebble instance. Metadata counters (next-node-ID, next-edge-ID, counts)
// are maintained as atomic in-memory values and persisted on Close.
//
// Pebble key schema:
//
//	node/{id:u64}              → binary node data
//	edge/{id:u64}              → binary edge data
//	adj_out/{nodeID:u64}{edgeID:u64} → binary [targetID:u64][label:bytes]
//	adj_in/{nodeID:u64}{edgeID:u64}  → binary [sourceID:u64][label:bytes]
//	idx_label/{label}\x00{id}  → nil
//	idx_prop/{prop}\x00{value}\x00{id} → nil
//	idx_etype/{label}\x00{edgeID} → nil
//	uc/{label}\x00{prop}\x00{value} → nodeID (u64)
//	meta/next_node_id          → u64
//	meta/next_edge_id          → u64
//	meta/node_count            → u64
//	meta/edge_count            → u64

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"

	"github.com/cockroachdb/pebble/v2"
)

var (
	idxEtypePrefix = []byte("idx_etype/")
	ucPrefix       = []byte("uc/")
	ErrNotFound    = errors.New("graphengine: not found")
)

var (
	keyNextNodeID = []byte("meta/next_node_id")
	keyNextEdgeID = []byte("meta/next_edge_id")
	keyNodeCount  = []byte("meta/node_count")
	keyEdgeCount  = []byte("meta/edge_count")

	nodePrefix     = []byte("node/")
	edgePrefix     = []byte("edge/")
	adjOutPrefix   = []byte("adj_out/")
	adjInPrefix    = []byte("adj_in/")
	idxLabelPrefix = []byte("idx_label/")
	idxPropPrefix  = []byte("idx_prop/")
)

// Options configures the graph database. SyncMode and Logger are reserved for
// future use when configurable durability guarantees are introduced.
type Options struct {
	CacheSize int64
	SyncMode  SyncMode
	Logger    interface{ Printf(string, ...interface{}) }
}

type SyncMode int

const (
	SyncNormal SyncMode = iota
	SyncNone
)

func DefaultOptions() Options {
	return Options{
		CacheSize: 64 << 20,
		SyncMode:  SyncNone,
	}
}

// DB is the top-level graph database handle. All operations are safe for
// concurrent use. Close must be called exactly once to release the Pebble
// instance and the associated cache.
type DB struct {
	db         *pebble.DB
	dir        string
	cache      *pebble.Cache
	nextNodeID atomic.Uint64
	nextEdgeID atomic.Uint64
	nodeCount  atomic.Uint64
	edgeCount  atomic.Uint64
	closed     atomic.Bool
}

// Open creates or opens a graph database at the given directory. A Pebble
// cache is allocated with the configured size and its reference is tracked so
// that Close can release it cleanly.
func Open(dir string, opts Options) (*DB, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("graphengine: failed to create directory %s: %w", dir, err)
	}

	cache := pebble.NewCache(opts.CacheSize)

	pebbleOpts := &pebble.Options{
		Cache:        cache,
		MemTableSize: 4 << 20,
	}

	pdb, err := pebble.Open(filepath.Join(dir, "graph"), pebbleOpts)
	if err != nil {
		cache.Unref()
		return nil, fmt.Errorf("graphengine: failed to open pebble at %s: %w", dir, err)
	}

	d := &DB{
		db:    pdb,
		dir:   dir,
		cache: cache,
	}

	if err := d.loadCounters(); err != nil {
		pdb.Close()
		cache.Unref()
		return nil, err
	}

	return d, nil
}

// Close persists metadata counters to disk, closes the Pebble instance, and
// releases the cache. Safe to call multiple times; subsequent calls are no-ops.
func (d *DB) Close() error {
	if d.closed.Swap(true) {
		return nil
	}
	var err error
	if err = d.persistCounters(); err != nil {
		d.db.Close()
		d.cache.Unref()
		return err
	}
	err = d.db.Close()
	d.cache.Unref()
	return err
}

func (d *DB) Dir() string {
	return d.dir
}

func (d *DB) Stats() *GraphStats {
	return &GraphStats{
		NodeCount:     int64(d.nodeCount.Load()),
		EdgeCount:     int64(d.edgeCount.Load()),
		LabelCounts:   d.GetLabelCounts(),
		RelCounts:     d.GetRelCounts(),
		DiskSizeBytes: int64(d.db.Metrics().DiskSpaceUsage()),
	}
}

func (d *DB) CountNodes() int64 {
	return int64(d.nodeCount.Load())
}

func (d *DB) CountEdges() int64 {
	return int64(d.edgeCount.Load())
}

func (d *DB) GetLabelCounts() map[string]int64 {
	counts := make(map[string]int64)
	lower := append([]byte(nil), idxLabelPrefix...)
	upper := append([]byte(nil), idxLabelPrefix...)
	upper = append(upper, 0xff)
	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: lower,
		UpperBound: upper,
	})
	if err != nil {
		return counts
	}
	defer iter.Close()
	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		rest := key[len(idxLabelPrefix):]
		idx := 0
		for idx < len(rest) && rest[idx] != 0x00 {
			idx++
		}
		if idx > 0 {
			label := string(rest[:idx])
			counts[label]++
		}
	}
	if err := iter.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: GetLabelCounts iteration error: %v\n", err)
	}
	return counts
}

func (d *DB) GetRelCounts() map[string]int64 {
	counts := make(map[string]int64)
	lower := append([]byte(nil), idxEtypePrefix...)
	upper := append([]byte(nil), idxEtypePrefix...)
	upper = append(upper, 0xff)
	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: lower,
		UpperBound: upper,
	})
	if err != nil {
		return counts
	}
	defer iter.Close()
	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		rest := key[len(idxEtypePrefix):]
		idx := 0
		for idx < len(rest) && rest[idx] != 0x00 {
			idx++
		}
		if idx > 0 {
			label := string(rest[:idx])
			counts[label]++
		}
	}
	if err := iter.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: GetRelCounts iteration error: %v\n", err)
	}
	return counts
}

// loadCounters reads persisted metadata from Pebble. If all counters are zero
// (fresh database or corrupted metadata), selfHealCounters scans all node and
// edge keys to reconstruct the correct values.
func (d *DB) loadCounters() error {
	var loadedNodeID, loadedEdgeID uint64
	var nodeIDLoaded, edgeIDLoaded bool

	val, closer, err := d.db.Get(keyNextNodeID)
	if err != nil && err != pebble.ErrNotFound {
		return err
	}
	if err == nil {
		loadedNodeID = decodeUint64(val)
		d.nextNodeID.Store(loadedNodeID)
		nodeIDLoaded = true
		closer.Close()
	}

	val, closer, err = d.db.Get(keyNextEdgeID)
	if err != nil && err != pebble.ErrNotFound {
		return err
	}
	if err == nil {
		loadedEdgeID = decodeUint64(val)
		d.nextEdgeID.Store(loadedEdgeID)
		edgeIDLoaded = true
		closer.Close()
	}

	val, closer, err = d.db.Get(keyNodeCount)
	if err != nil && err != pebble.ErrNotFound {
		return err
	}
	if err == nil {
		d.nodeCount.Store(decodeUint64(val))
		closer.Close()
	}

	val, closer, err = d.db.Get(keyEdgeCount)
	if err != nil && err != pebble.ErrNotFound {
		return err
	}
	if err == nil {
		d.edgeCount.Store(decodeUint64(val))
		closer.Close()
	}

	if d.nextNodeID.Load() == 0 || d.nextEdgeID.Load() == 0 {
		d.selfHealCounters()
		if nodeIDLoaded && loadedNodeID > d.nextNodeID.Load() {
			d.nextNodeID.Store(loadedNodeID)
		}
		if edgeIDLoaded && loadedEdgeID > d.nextEdgeID.Load() {
			d.nextEdgeID.Store(loadedEdgeID)
		}
	}

	return nil
}

// selfHealCounters performs a full scan of node and edge keys to reconstruct
// the next-ID and count metadata. Called only when loadCounters detects that
// all persisted values are zero, indicating a fresh or uncleanly closed
// database.
func (d *DB) selfHealCounters() {
	var maxNodeID, maxEdgeID, nodeCt, edgeCt uint64

	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: nodePrefix,
		UpperBound: nodeEndKey(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: failed to create iterator for selfHealCounters (nodes): %v\n", err)
		return
	}
	for iter.First(); iter.Valid(); iter.Next() {
		nodeCt++
		id := NodeID(decodeUint64(iter.Key()[len(nodePrefix):]))
		if uint64(id) > maxNodeID {
			maxNodeID = uint64(id)
		}
	}
	if err := iter.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: selfHealCounters node iteration error: %v\n", err)
		iter.Close()
		return
	}
	iter.Close()

	iter2, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: edgePrefix,
		UpperBound: edgeEndKey(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: failed to create iterator for selfHealCounters (edges): %v\n", err)
		return
	}
	for iter2.First(); iter2.Valid(); iter2.Next() {
		edgeCt++
		id := EdgeID(decodeUint64(iter2.Key()[len(edgePrefix):]))
		if uint64(id) > maxEdgeID {
			maxEdgeID = uint64(id)
		}
	}
	if err := iter2.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: selfHealCounters edge iteration error: %v\n", err)
		iter2.Close()
		return
	}
	iter2.Close()

	d.nodeCount.Store(nodeCt)
	d.edgeCount.Store(edgeCt)
	if maxNodeID > 0 {
		d.nextNodeID.Store(maxNodeID + 1)
	}
	if maxEdgeID > 0 {
		d.nextEdgeID.Store(maxEdgeID + 1)
	}
}

// persistCounters writes the current metadata counters to Pebble with Sync
// to ensure durability across restarts.
func (d *DB) persistCounters() error {
	batch := d.db.NewBatch()
	defer batch.Close()

	batch.Set(keyNextNodeID, encodeUint64(d.nextNodeID.Load()), nil)
	batch.Set(keyNextEdgeID, encodeUint64(d.nextEdgeID.Load()), nil)
	batch.Set(keyNodeCount, encodeUint64(d.nodeCount.Load()), nil)
	batch.Set(keyEdgeCount, encodeUint64(d.edgeCount.Load()), nil)

	return d.db.Apply(batch, pebble.Sync)
}

// --- Key construction helpers ---
//
// All key helpers allocate new slices (via append to nil) to avoid sharing
// backing arrays between calls, which Pebble requires for correctness.

func nodeKey(id NodeID) []byte {
	return append(append([]byte(nil), nodePrefix...), encodeNodeID(id)...)
}

// nodeEndKey returns the exclusive upper bound for node key iteration.
// Since node keys are "node/{id}", the byte after 'n' ('o') suffices.
func nodeEndKey() []byte {
	return append([]byte(nil), nodePrefix[0]+1)
}

func edgeKey(id EdgeID) []byte {
	return append(append([]byte(nil), edgePrefix...), encodeEdgeID(id)...)
}

// edgeEndKey returns the exclusive upper bound for edge key iteration.
func edgeEndKey() []byte {
	return append([]byte(nil), edgePrefix[0]+1)
}

func adjOutKey(nodeID NodeID, edgeID EdgeID) []byte {
	return append(append([]byte(nil), adjOutPrefix...), encodeAdjKey(nodeID, edgeID)...)
}

// adjOutPrefixKey / adjOutEndKey bound the outgoing adjacency scan for a
// given node. The end key increments the last byte of the prefix to create
// an exclusive upper bound within the same node-ID space.
func adjOutPrefixKey(nodeID NodeID) []byte {
	return append(append([]byte(nil), adjOutPrefix...), encodeNodeID(nodeID)...)
}

func adjOutEndKey(nodeID NodeID) []byte {
	prefix := adjOutPrefixKey(nodeID)
	end := make([]byte, len(prefix)+8)
	copy(end, prefix)
	for i := len(prefix); i < len(end); i++ {
		end[i] = 0xFF
	}
	return end
}

func adjInKey(nodeID NodeID, edgeID EdgeID) []byte {
	return append(append([]byte(nil), adjInPrefix...), encodeAdjKey(nodeID, edgeID)...)
}

// adjInPrefixKey / adjInEndKey bound the incoming adjacency scan for a
// given node.
func adjInPrefixKey(nodeID NodeID) []byte {
	return append(append([]byte(nil), adjInPrefix...), encodeNodeID(nodeID)...)
}

func adjInEndKey(nodeID NodeID) []byte {
	prefix := adjInPrefixKey(nodeID)
	end := make([]byte, len(prefix)+8)
	copy(end, prefix)
	for i := len(prefix); i < len(end); i++ {
		end[i] = 0xFF
	}
	return end
}

func idxLabelKey(label string, id NodeID) []byte {
	prefix := append(append([]byte(nil), idxLabelPrefix...), []byte(label)...)
	prefix = append(prefix, 0x00)
	return append(prefix, encodeNodeID(id)...)
}

func idxLabelPrefixKey(label string) []byte {
	prefix := append(append([]byte(nil), idxLabelPrefix...), []byte(label)...)
	return append(prefix, 0x00)
}

func idxPropKey(prop, value string, id NodeID) []byte {
	prefix := append(append([]byte(nil), idxPropPrefix...), prop...)
	prefix = append(prefix, 0x00)
	prefix = append(prefix, value...)
	prefix = append(prefix, 0x00)
	return append(prefix, encodeNodeID(id)...)
}

func idxPropPrefixKey(prop, value string) []byte {
	prefix := append(append([]byte(nil), idxPropPrefix...), prop...)
	prefix = append(prefix, 0x00)
	prefix = append(prefix, value...)
	return append(prefix, 0x00)
}

func idxEtypeKey(label string, id EdgeID) []byte {
	prefix := append(append([]byte(nil), idxEtypePrefix...), []byte(label)...)
	prefix = append(prefix, 0x00)
	return append(prefix, encodeEdgeID(id)...)
}

func idxEtypePrefixKey(label string) []byte {
	prefix := append(append([]byte(nil), idxEtypePrefix...), []byte(label)...)
	return append(prefix, 0x00)
}

func ucKey(label, prop string, value interface{}) []byte {
	prefix := append(append([]byte(nil), ucPrefix...), []byte(label)...)
	prefix = append(prefix, 0x00)
	prefix = append(prefix, []byte(prop)...)
	prefix = append(prefix, 0x00)
	prefix = append(prefix, propIndexValue(value)...)
	return prefix
}

func propIndexValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int64:
		return strconv.FormatInt(val, 10)
	case int:
		return strconv.Itoa(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}

func (d *DB) allocNodeID() NodeID {
	return NodeID(d.nextNodeID.Add(1))
}

func (d *DB) allocEdgeID() EdgeID {
	return EdgeID(d.nextEdgeID.Add(1))
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
	batch.Set(nodeKey(id), data, nil)

	for _, label := range labels {
		batch.Set(idxLabelKey(label, id), nil, nil)
	}

	for k, v := range props {
		batch.Set(idxPropKey(k, propIndexValue(v), id), nil, nil)
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

	for i, item := range items {
		if err := d.enforceUniqueConstraints(item.Labels, item.Props, 0); err != nil {
			return nil, err
		}

		id := d.allocNodeID()
		ids[i] = id

		data, err := encodeNodeData(item.Labels, item.Props)
		if err != nil {
			return nil, err
		}
		batch.Set(nodeKey(id), data, nil)

		for _, label := range item.Labels {
			batch.Set(idxLabelKey(label, id), nil, nil)
		}

		for k, v := range item.Props {
			batch.Set(idxPropKey(k, propIndexValue(v), id), nil, nil)
		}

		d.writeUniqueConstraintEntries(batch, item.Labels, item.Props, id)
	}

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return nil, fmt.Errorf("graphengine: batch add failed: %w", err)
	}

	d.nodeCount.Add(uint64(len(items)))
	return ids, nil
}

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

	batch := d.db.NewBatch()
	defer batch.Close()

	for k, v := range props {
		if oldVal, ok := existingProps[k]; ok && propIndexValue(oldVal) != propIndexValue(v) {
			batch.Delete(idxPropKey(k, propIndexValue(oldVal), id), nil)
			if constraintLabels, isConstrained := ucProps[k]; isConstrained {
				for _, cl := range constraintLabels {
					batch.Delete(ucKey(cl, k, oldVal), nil)
					batch.Set(ucKey(cl, k, v), encodeNodeID(id), nil)
				}
			}
		} else if _, ok := existingProps[k]; !ok {
			if constraintLabels, isConstrained := ucProps[k]; isConstrained {
				for _, cl := range constraintLabels {
					batch.Set(ucKey(cl, k, v), encodeNodeID(id), nil)
				}
			}
		}
		existingProps[k] = v
		batch.Set(idxPropKey(k, propIndexValue(v), id), nil, nil)
	}

	data, err := encodeNodeData(labels, existingProps)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode updated node %d: %w", id, err)
	}
	batch.Set(nodeKey(id), data, nil)

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
	d.forEachUniqueConstraint(func(label, prop string) error {
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
		batch.Delete(edgeKey(e.ID), nil)
		batch.Delete(adjOutKey(e.From, e.ID), nil)
		batch.Delete(adjInKey(e.To, e.ID), nil)
		if e.Label != "" {
			batch.Delete(idxEtypeKey(e.Label, e.ID), nil)
		}
	}

	batch.Delete(nodeKey(id), nil)

	for _, label := range node.Labels {
		batch.Delete(idxLabelKey(label, id), nil)
	}

	for k, v := range node.Props {
		batch.Delete(idxPropKey(k, propIndexValue(v), id), nil)
	}

	d.deleteUniqueConstraintEntries(batch, node.Labels, node.Props)

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to delete node %d: %w", id, err)
	}

	d.nodeCount.Add(^uint64(0))
	if len(edges) > 0 {
		d.edgeCount.Add(^uint64(uint64(len(edges) - 1)))
	}
	return nil
}

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
	batch.Set(edgeKey(id), edgeData, nil)
	batch.Set(adjOutKey(from, id), encodeAdjValue(to, label), nil)
	batch.Set(adjInKey(to, id), encodeAdjValue(from, label), nil)
	if label != "" {
		batch.Set(idxEtypeKey(label, id), nil, nil)
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
		batch.Set(edgeKey(id), edgeData, nil)
		batch.Set(adjOutKey(edge.From, id), encodeAdjValue(edge.To, edge.Label), nil)
		batch.Set(adjInKey(edge.To, id), encodeAdjValue(edge.From, edge.Label), nil)
		if edge.Label != "" {
			batch.Set(idxEtypeKey(edge.Label, id), nil, nil)
		}
	}

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return nil, fmt.Errorf("graphengine: batch edge add failed: %w", err)
	}

	d.edgeCount.Add(uint64(len(edges)))
	return ids, nil
}

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

	batch.Delete(edgeKey(id), nil)
	batch.Delete(adjOutKey(edge.From, id), nil)
	batch.Delete(adjInKey(edge.To, id), nil)
	if edge.Label != "" {
		batch.Delete(idxEtypeKey(edge.Label, id), nil)
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
				if full, err := d.GetEdge(edgeID); err == nil {
					edge = full
				}
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
				if full, err := d.GetEdge(edgeID); err == nil {
					edge = full
				}
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

// GetAdjacentNodes returns the deduplicated set of nodes adjacent to the
// given node in the specified direction. For Both direction, outgoing edges
// contribute their To node and incoming edges contribute their From node.
// When labelFilter is non-empty, only edges with a matching label are
// considered.
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
	batch.Set(nodeKey(id), nodeData, nil)
	batch.Delete(idxLabelKey(label, id), nil)

	d.forEachUniqueConstraint(func(ucLabel, prop string) error {
		if ucLabel == label {
			if val, ok := props[prop]; ok {
				batch.Delete(ucKey(label, prop, val), nil)
			}
		}
		return nil
	})

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to remove label from node %d: %w", id, err)
	}

	return nil
}

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
					batch.Set(ucKey(label, prop, val), encodeNodeID(id), nil)
				}
			}
		}
	}

	nodeData, err := encodeNodeData(labels, props)
	if err != nil {
		return fmt.Errorf("graphengine: failed to encode node after adding label: %w", err)
	}
	batch.Set(nodeKey(id), nodeData, nil)
	batch.Set(idxLabelKey(label, id), nil, nil)

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to add label to node %d: %w", id, err)
	}

	return nil
}

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
	batch.Set(nodeKey(id), nodeData, nil)
	batch.Delete(idxPropKey(key, propIndexValue(oldValue), id), nil)

	if constraintLabels, isConstrained := ucProps[key]; isConstrained {
		for _, cl := range constraintLabels {
			batch.Delete(ucKey(cl, key, oldValue), nil)
		}
	}

	if err := d.db.Apply(batch, pebble.NoSync); err != nil {
		return fmt.Errorf("graphengine: failed to remove property from node %d: %w", id, err)
	}

	return nil
}

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

	d.forEachUniqueConstraint(func(label, prop string) error {
		if !labelSet[label] {
			return nil
		}

		val, ok := props[prop]
		if !ok {
			return nil
		}

		key := ucKey(label, prop, val)
		batch.Set(key, encodeNodeID(id), nil)
		return nil
	})
}

func (d *DB) deleteUniqueConstraintEntries(batch *pebble.Batch, labels []string, props Props) {
	labelSet := make(map[string]bool, len(labels))
	for _, l := range labels {
		labelSet[l] = true
	}

	d.forEachUniqueConstraint(func(label, prop string) error {
		if !labelSet[label] {
			return nil
		}

		val, ok := props[prop]
		if !ok {
			return nil
		}

		key := ucKey(label, prop, val)
		batch.Delete(key, nil)
		return nil
	})
}

func (d *DB) Clear() error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	iter, err := d.db.NewIter(nil)
	if err != nil {
		return fmt.Errorf("graphengine: failed to create iterator for Clear: %w", err)
	}
	defer iter.Close()

	batch := d.db.NewBatch()

	count := 0
	for iter.First(); iter.Valid(); iter.Next() {
		batch.Delete(iter.Key(), nil)
		count++
		if count >= 10000 {
			if err := d.db.Apply(batch, pebble.NoSync); err != nil {
				batch.Close()
				return fmt.Errorf("graphengine: failed to clear: %w", err)
			}
			batch.Close()
			batch = d.db.NewBatch()
			count = 0
		}
	}

	if err := iter.Error(); err != nil {
		batch.Close()
		return fmt.Errorf("graphengine: Clear iteration error: %w", err)
	}

	if count > 0 {
		if err := d.db.Apply(batch, pebble.NoSync); err != nil {
			batch.Close()
			return fmt.Errorf("graphengine: failed to clear: %w", err)
		}
	}
	batch.Close()

	d.nodeCount.Store(0)
	d.edgeCount.Store(0)
	d.nextNodeID.Store(1)
	d.nextEdgeID.Store(1)

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

	return nil
}

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
