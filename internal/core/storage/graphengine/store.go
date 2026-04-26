package graphengine

import (
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
	idxMetaPrefix  = []byte("idx_meta/")
)

type Cache struct {
	cache *pebble.Cache
}

func NewSharedCache(size int64) *Cache {
	return &Cache{cache: pebble.NewCache(size)}
}

func (c *Cache) Release() {
	c.cache.Unref()
}

type Options struct {
	CacheSize   int64
	SyncMode    SyncMode
	Logger      interface{ Printf(string, ...interface{}) }
	SharedCache *Cache
}

type SyncMode int

const (
	SyncNormal SyncMode = iota
	SyncNone
)

const DefaultCacheSize int64 = 8 << 20

func DefaultOptions() Options {
	return Options{
		CacheSize: DefaultCacheSize,
		SyncMode:  SyncNone,
	}
}

type DB struct {
	backend     kvBackend
	dir         string
	cache       *pebble.Cache
	ownsCache   bool
	nextNodeID  atomic.Uint64
	nextEdgeID  atomic.Uint64
	nodeCount   atomic.Uint64
	edgeCount   atomic.Uint64
	closed      atomic.Bool
	Embeddings  *EmbeddingStore
}

// Open creates a graph engine backed by an independent Pebble database.
// Deprecated: use New(bucket, opts) for shared storage via RegionStorageManager.
func Open(dir string, opts Options) (*DB, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("graphengine: failed to create directory %s: %w", dir, err)
	}

	var pebbleCache *pebble.Cache
	ownsCache := opts.SharedCache == nil
	if ownsCache {
		pebbleCache = pebble.NewCache(opts.CacheSize)
	} else {
		opts.SharedCache.cache.Ref()
		pebbleCache = opts.SharedCache.cache
	}

	pebbleOpts := &pebble.Options{
		Cache:        pebbleCache,
		MemTableSize: 4 << 20,
	}

	pdb, err := pebble.Open(filepath.Join(dir, "graph"), pebbleOpts)
	if err != nil {
		if ownsCache {
			pebbleCache.Unref()
		} else {
			opts.SharedCache.cache.Unref()
		}
		return nil, fmt.Errorf("graphengine: failed to open pebble at %s: %w", dir, err)
	}

	backend := &pebbleBackend{db: pdb, cache: pebbleCache}

	d := &DB{
		backend:   backend,
		dir:       dir,
		cache:     pebbleCache,
		ownsCache: ownsCache,
	}

	if err := d.loadCounters(); err != nil {
		pdb.Close()
		if ownsCache {
			pebbleCache.Unref()
		} else {
			opts.SharedCache.cache.Unref()
		}
		return nil, err
	}

	d.Embeddings = newEmbeddingStore(d)
	return d, nil
}

func (d *DB) Close() error {
	if d.closed.Swap(true) {
		return nil
	}
	var err error
	if err = d.persistCounters(); err != nil {
		d.backend.close()
		if d.cache != nil {
			d.cache.Unref()
		}
		return err
	}
	err = d.backend.close()
	if d.cache != nil {
		d.cache.Unref()
	}
	return err
}

func (d *DB) Dir() string {
	return d.dir
}

func (d *DB) Stats() *GraphStats {
	labelCounts, _ := d.GetLabelCounts()
	relCounts, _ := d.GetRelCounts()
	return &GraphStats{
		NodeCount:     int64(d.nodeCount.Load()),
		EdgeCount:     int64(d.edgeCount.Load()),
		LabelCounts:   labelCounts,
		RelCounts:     relCounts,
		DiskSizeBytes: d.backend.diskSize(),
	}
}

func (d *DB) CountNodes() int64 {
	return int64(d.nodeCount.Load())
}

func (d *DB) CountEdges() int64 {
	return int64(d.edgeCount.Load())
}

func (d *DB) GetLabelCounts() (map[string]int64, error) {
	counts := make(map[string]int64)
	lower := append([]byte(nil), idxLabelPrefix...)
	upper := append([]byte(nil), idxLabelPrefix...)
	upper = append(upper, 0xff)

	it, err := d.backend.newIter(lower, upper)
	if err != nil {
		return counts, fmt.Errorf("graphengine: GetLabelCounts iterator creation failed: %w", err)
	}
	defer it.close()

	for it.first(); it.valid(); it.next() {
		key := it.key()
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
	if err := it.err(); err != nil {
		return nil, fmt.Errorf("graphengine: GetLabelCounts iteration error: %w", err)
	}
	return counts, nil
}

func (d *DB) GetRelCounts() (map[string]int64, error) {
	counts := make(map[string]int64)
	lower := append([]byte(nil), idxEtypePrefix...)
	upper := append([]byte(nil), idxEtypePrefix...)
	upper = append(upper, 0xff)

	it, err := d.backend.newIter(lower, upper)
	if err != nil {
		return counts, fmt.Errorf("graphengine: GetRelCounts iterator creation failed: %w", err)
	}
	defer it.close()

	for it.first(); it.valid(); it.next() {
		key := it.key()
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
	if err := it.err(); err != nil {
		return nil, fmt.Errorf("graphengine: GetRelCounts iteration error: %w", err)
	}
	return counts, nil
}

func (d *DB) loadCounters() error {
	var loadedNodeID, loadedEdgeID uint64
	var nodeIDLoaded, edgeIDLoaded bool

	val, err := d.backend.get(keyNextNodeID)
	if err != nil {
		return err
	}
	if val != nil {
		loadedNodeID = decodeUint64(val)
		d.nextNodeID.Store(loadedNodeID)
		nodeIDLoaded = true
	}

	val, err = d.backend.get(keyNextEdgeID)
	if err != nil {
		return err
	}
	if val != nil {
		loadedEdgeID = decodeUint64(val)
		d.nextEdgeID.Store(loadedEdgeID)
		edgeIDLoaded = true
	}

	val, err = d.backend.get(keyNodeCount)
	if err != nil {
		return err
	}
	if val != nil {
		d.nodeCount.Store(decodeUint64(val))
	}

	val, err = d.backend.get(keyEdgeCount)
	if err != nil {
		return err
	}
	if val != nil {
		d.edgeCount.Store(decodeUint64(val))
	}

	d.selfHealCounters()
	if nodeIDLoaded && loadedNodeID > d.nextNodeID.Load() {
		d.nextNodeID.Store(loadedNodeID)
	}
	if edgeIDLoaded && loadedEdgeID > d.nextEdgeID.Load() {
		d.nextEdgeID.Store(loadedEdgeID)
	}

	return nil
}

func (d *DB) selfHealCounters() {
	var maxNodeID, maxEdgeID, nodeCt, edgeCt uint64

	it, err := d.backend.newIter(nodePrefix, nodeEndKey())
	if err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: failed to create iterator for selfHealCounters (nodes): %v\n", err)
		return
	}
	for it.first(); it.valid(); it.next() {
		nodeCt++
		id := NodeID(decodeUint64(it.key()[len(nodePrefix):]))
		if uint64(id) > maxNodeID {
			maxNodeID = uint64(id)
		}
	}
	if err := it.err(); err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: selfHealCounters node iteration error: %v\n", err)
		it.close()
		return
	}
	it.close()

	it2, err := d.backend.newIter(edgePrefix, edgeEndKey())
	if err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: failed to create iterator for selfHealCounters (edges): %v\n", err)
		return
	}
	for it2.first(); it2.valid(); it2.next() {
		edgeCt++
		id := EdgeID(decodeUint64(it2.key()[len(edgePrefix):]))
		if uint64(id) > maxEdgeID {
			maxEdgeID = uint64(id)
		}
	}
	if err := it2.err(); err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: selfHealCounters edge iteration error: %v\n", err)
		it2.close()
		return
	}
	it2.close()

	d.nodeCount.Store(nodeCt)
	d.edgeCount.Store(edgeCt)
	if maxNodeID > 0 {
		d.nextNodeID.Store(maxNodeID + 1)
	}
	if maxEdgeID > 0 {
		d.nextEdgeID.Store(maxEdgeID + 1)
	}
}

func (d *DB) persistCounters() error {
	batch := d.backend.newBatch()
	defer batch.close()

	batch.put(keyNextNodeID, encodeUint64(d.nextNodeID.Load()))
	batch.put(keyNextEdgeID, encodeUint64(d.nextEdgeID.Load()))
	batch.put(keyNodeCount, encodeUint64(d.nodeCount.Load()))
	batch.put(keyEdgeCount, encodeUint64(d.edgeCount.Load()))

	return batch.commitSync()
}

func nodeKey(id NodeID) []byte {
	return append(append([]byte(nil), nodePrefix...), encodeNodeID(id)...)
}

func nodeEndKey() []byte {
	return append(append([]byte(nil), nodePrefix...), 0xFF)
}

func edgeKey(id EdgeID) []byte {
	return append(append([]byte(nil), edgePrefix...), encodeEdgeID(id)...)
}

func edgeEndKey() []byte {
	return append(append([]byte(nil), edgePrefix...), 0xFF)
}

func adjOutKey(nodeID NodeID, edgeID EdgeID) []byte {
	return append(append([]byte(nil), adjOutPrefix...), encodeAdjKey(nodeID, edgeID)...)
}

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
