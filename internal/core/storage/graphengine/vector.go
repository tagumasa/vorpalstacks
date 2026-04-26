package graphengine

// vector.go provides in-memory vector embedding storage, distance functions,
// brute-force topK search, and vertexFilter evaluation for Neptune Analytics
// vector search algorithms.
//
// Pebble key schema extension:
//
//	vec/{nodeID:u64} → [dimension:u32][float64×dimension] (big-endian)
//
// Embeddings are loaded into an in-memory cache on demand and kept in sync
// with Pebble on upsert/remove. The cache is the "in-memory vector index"
// that Neptune Analytics uses.

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"
)

var vecPrefix = []byte("vec/")

// DistanceMetric enumerates the supported distance functions for vector search.
type DistanceMetric string

// Distance metric constants.
const (
	L2Squared        DistanceMetric = "L2Squared"
	L2               DistanceMetric = "L2"
	CosineSimilarity DistanceMetric = "CosineSimilarity"
	CosineDistance   DistanceMetric = "CosineDistance"
	DotProduct       DistanceMetric = "DotProduct"
)

// ParseDistanceMetric converts a case-insensitive string to a DistanceMetric.
// Returns L2Squared (the default) for unrecognised values.
func ParseDistanceMetric(s string) DistanceMetric {
	switch s {
	case "L2Squared", "l2squared":
		return L2Squared
	case "L2", "l2":
		return L2
	case "CosineSimilarity", "cosinesimilarity":
		return CosineSimilarity
	case "CosineDistance", "cosinedistance":
		return CosineDistance
	case "DotProduct", "dotproduct":
		return DotProduct
	default:
		return L2Squared
	}
}

// ComputeDistance returns the distance between two vectors using the given metric.
// Returns an error if the vectors have different lengths.
func ComputeDistance(a, b []float64, metric DistanceMetric) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("graphengine: vector dimension mismatch: %d vs %d", len(a), len(b))
	}
	switch metric {
	case L2Squared:
		return l2Squared(a, b), nil
	case L2:
		return math.Sqrt(l2Squared(a, b)), nil
	case CosineSimilarity:
		return cosineSimilarity(a, b), nil
	case CosineDistance:
		return 1.0 - cosineSimilarity(a, b), nil
	case DotProduct:
		return dotProduct(a, b), nil
	default:
		return l2Squared(a, b), nil
	}
}

func l2Squared(a, b []float64) float64 {
	var sum float64
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return sum
}

func cosineSimilarity(a, b []float64) float64 {
	var dotA, dotB, dotAB float64
	for i := range a {
		dotA += a[i] * a[i]
		dotB += b[i] * b[i]
		dotAB += a[i] * b[i]
	}
	norm := math.Sqrt(dotA) * math.Sqrt(dotB)
	if norm == 0 {
		return 0
	}
	return dotAB / norm
}

func dotProduct(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// TopKResult represents a single result from a topK search.
type TopKResult struct {
	NodeID NodeID
	Node   *Node
	Score  float64
}

// EmbeddingStore manages in-memory vector embeddings backed by Pebble persistence.
type EmbeddingStore struct {
	db    pebbleDB
	cache map[NodeID][]float64
	mu    sync.RWMutex
}

// pebbleDB is a minimal interface for the vector store to avoid depending
// directly on the Pebble import. The concrete *DB type satisfies this.
type pebbleDB interface {
	get(key []byte) ([]byte, error)
	set(key, value []byte) error
	delete(key []byte) error
	deleteRange(prefix []byte) error
	iteratePrefix(prefix []byte, fn func(key, value []byte) error) error
}

// NewEmbeddingStore creates an embedding store backed by the given Pebble DB.
func NewEmbeddingStore(db *DB) *EmbeddingStore {
	es := &EmbeddingStore{
		db:    &pebbleDBAdapter{db: db},
		cache: make(map[NodeID][]float64),
	}
	es.loadCache()
	return es
}

// pebbleDBAdapter wraps *DB to satisfy the pebbleDB interface.
type pebbleDBAdapter struct {
	db *DB
}

func (a *pebbleDBAdapter) get(key []byte) ([]byte, error) {
	return a.db.getRaw(key)
}

func (a *pebbleDBAdapter) set(key, value []byte) error {
	return a.db.setRaw(key, value)
}

func (a *pebbleDBAdapter) delete(key []byte) error {
	return a.db.deleteRaw(key)
}

func (a *pebbleDBAdapter) deleteRange(prefix []byte) error {
	return a.db.deleteRangeRaw(prefix)
}

func (a *pebbleDBAdapter) iteratePrefix(prefix []byte, fn func(key, value []byte) error) error {
	return a.db.iteratePrefixRaw(prefix, fn)
}

// loadCache reads all vec/ entries from Pebble into the in-memory cache.
func (es *EmbeddingStore) loadCache() {
	es.mu.Lock()
	defer es.mu.Unlock()

	es.cache = make(map[NodeID][]float64)
	if err := es.db.iteratePrefix(vecPrefix, func(key, value []byte) error {
		nodeID := decodeNodeID(key[len(vecPrefix):])
		emb, err := decodeEmbedding(value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "graphengine: skipping corrupted embedding for node %d: %v\n", nodeID, err)
			return nil
		}
		es.cache[nodeID] = emb
		return nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "graphengine: error loading embedding cache: %v\n", err)
	}
}

// Get retrieves the embedding for a node. Returns nil if not found.
func (es *EmbeddingStore) Get(nodeID NodeID) ([]float64, error) {
	es.mu.RLock()
	emb, ok := es.cache[nodeID]
	es.mu.RUnlock()
	if ok {
		return emb, nil
	}
	return nil, nil
}

// Upsert stores an embedding for a node, validating dimension if declaredDim > 0.
func (es *EmbeddingStore) Upsert(nodeID NodeID, embedding []float64, declaredDim int32) error {
	if declaredDim > 0 && int32(len(embedding)) != declaredDim {
		return fmt.Errorf("graphengine: embedding dimension %d does not match declared dimension %d", len(embedding), declaredDim)
	}
	for _, v := range embedding {
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return fmt.Errorf("graphengine: embedding contains invalid float value (INF or NaN)")
		}
	}

	key := make([]byte, len(vecPrefix)+8)
	copy(key, vecPrefix)
	binary.BigEndian.PutUint64(key[len(vecPrefix):], uint64(nodeID))

	value := encodeEmbedding(embedding)

	if err := es.db.set(key, value); err != nil {
		return fmt.Errorf("graphengine: failed to persist embedding: %w", err)
	}

	es.mu.Lock()
	copied := make([]float64, len(embedding))
	copy(copied, embedding)
	es.cache[nodeID] = copied
	es.mu.Unlock()

	return nil
}

// Remove deletes the embedding for a node. Returns false if no embedding existed.
func (es *EmbeddingStore) Remove(nodeID NodeID) (bool, error) {
	key := make([]byte, len(vecPrefix)+8)
	copy(key, vecPrefix)
	binary.BigEndian.PutUint64(key[len(vecPrefix):], uint64(nodeID))

	if err := es.db.delete(key); err != nil {
		return false, fmt.Errorf("graphengine: failed to delete embedding: %w", err)
	}

	es.mu.Lock()
	_, existed := es.cache[nodeID]
	if existed {
		delete(es.cache, nodeID)
	}
	es.mu.Unlock()

	return existed, nil
}

// Clear removes all embeddings from both Pebble and the cache.
func (es *EmbeddingStore) Clear() error {
	if err := es.db.deleteRange(vecPrefix); err != nil {
		return fmt.Errorf("graphengine: failed to clear embeddings: %w", err)
	}
	es.mu.Lock()
	es.cache = make(map[NodeID][]float64)
	es.mu.Unlock()
	return nil
}

// Size returns the number of embeddings in the cache.
func (es *EmbeddingStore) Size() int {
	es.mu.RLock()
	defer es.mu.RUnlock()
	return len(es.cache)
}

// TopK performs brute-force topK search using the specified distance metric.
// If metric is empty, L2Squared is used. If filterFn is non-nil, only nodes
// passing the filter are considered.
func (es *EmbeddingStore) TopK(query []float64, k int, metric DistanceMetric, filterFn func(NodeID) bool) ([]TopKResult, error) {
	if k <= 0 {
		k = 10
	}
	if metric == "" {
		metric = L2Squared
	}

	es.mu.RLock()
	type scored struct {
		id    NodeID
		score float64
	}

	candidates := make([]scored, 0, len(es.cache))
	for id, emb := range es.cache {
		if filterFn != nil && !filterFn(id) {
			continue
		}
		if len(emb) != len(query) {
			continue
		}
		dist, err := ComputeDistance(query, emb, metric)
		if err != nil {
			continue
		}
		candidates = append(candidates, scored{id: id, score: dist})
	}
	es.mu.RUnlock()

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score < candidates[j].score
	})

	if k > len(candidates) {
		k = len(candidates)
	}

	results := make([]TopKResult, k)
	for i := 0; i < k; i++ {
		results[i] = TopKResult{
			NodeID: candidates[i].id,
			Score:  candidates[i].score,
		}
	}
	return results, nil
}

// AllNodeIDs returns all node IDs that have embeddings.
func (es *EmbeddingStore) AllNodeIDs() []NodeID {
	es.mu.RLock()
	defer es.mu.RUnlock()
	ids := make([]NodeID, 0, len(es.cache))
	for id := range es.cache {
		ids = append(ids, id)
	}
	return ids
}

// encodeEmbedding serialises a float64 slice to binary:
// [dimension:u32][float64×dimension] (big-endian).
func encodeEmbedding(emb []float64) []byte {
	dim := uint32(len(emb))
	buf := make([]byte, 4+dim*8)
	binary.BigEndian.PutUint32(buf[0:4], dim)
	for i, v := range emb {
		binary.BigEndian.PutUint64(buf[4+i*8:], math.Float64bits(v))
	}
	return buf
}

// decodeEmbedding deserialises a binary-encoded embedding.
func decodeEmbedding(data []byte) ([]float64, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("graphengine: embedding data too short")
	}
	dim := binary.BigEndian.Uint32(data[0:4])
	expected := 4 + dim*8
	if uint32(len(data)) < expected {
		return nil, fmt.Errorf("graphengine: embedding data truncated: expected %d bytes, got %d", expected, len(data))
	}
	emb := make([]float64, dim)
	for i := uint32(0); i < dim; i++ {
		emb[i] = math.Float64frombits(binary.BigEndian.Uint64(data[4+i*8 : 4+i*8+8]))
	}
	return emb, nil
}

// --- vertexFilter evaluator ---

// VertexFilter represents a parsed vertex filter expression tree.
type VertexFilter interface {
	Evaluate(node *Node) bool
}

type filterAndAll struct {
	children []VertexFilter
}

// Evaluate implements the VertexFilter interface.
func (f *filterAndAll) Evaluate(node *Node) bool {
	for _, c := range f.children {
		if !c.Evaluate(node) {
			return false
		}
	}
	return true
}

type filterOrAll struct {
	children []VertexFilter
}

// Evaluate implements the VertexFilter interface.
func (f *filterOrAll) Evaluate(node *Node) bool {
	for _, c := range f.children {
		if c.Evaluate(node) {
			return true
		}
	}
	return false
}

type filterEquals struct {
	property string
	value    interface{}
}

// Evaluate implements the VertexFilter interface.
func (f *filterEquals) Evaluate(node *Node) bool {
	if f.property == "~label" {
		for _, l := range node.Labels {
			if l == fmt.Sprint(f.value) {
				return true
			}
		}
		return false
	}
	v, ok := node.Props[f.property]
	if !ok {
		return false
	}
	return compareValues(v, f.value) == 0
}

type filterNotEquals struct {
	property string
	value    interface{}
}

// Evaluate implements the VertexFilter interface.
func (f *filterNotEquals) Evaluate(node *Node) bool {
	if f.property == "~label" {
		for _, l := range node.Labels {
			if l == fmt.Sprint(f.value) {
				return false
			}
		}
		return true
	}
	v, ok := node.Props[f.property]
	if !ok {
		return true
	}
	return compareValues(v, f.value) != 0
}

type filterComparison struct {
	property string
	value    float64
	op       string
}

// Evaluate implements the VertexFilter interface.
func (f *filterComparison) Evaluate(node *Node) bool {
	v, ok := node.Props[f.property]
	if !ok {
		return false
	}
	nv := toFloat64(v)
	switch f.op {
	case "greaterThan":
		return nv > f.value
	case "greaterThanOrEquals":
		return nv >= f.value
	case "lessThan":
		return nv < f.value
	case "lessThanOrEquals":
		return nv <= f.value
	}
	return false
}

type filterIn struct {
	property string
	values   []interface{}
}

// Evaluate implements the VertexFilter interface.
func (f *filterIn) Evaluate(node *Node) bool {
	v, ok := node.Props[f.property]
	if !ok {
		return false
	}
	for _, val := range f.values {
		if compareValues(v, val) == 0 {
			return true
		}
	}
	return false
}

type filterNotIn struct {
	property string
	values   []interface{}
}

// Evaluate implements the VertexFilter interface.
func (f *filterNotIn) Evaluate(node *Node) bool {
	v, ok := node.Props[f.property]
	if !ok {
		return true
	}
	for _, val := range f.values {
		if compareValues(v, val) == 0 {
			return false
		}
	}
	return true
}

type filterStartsWith struct {
	property string
	value    string
}

// Evaluate implements the VertexFilter interface.
func (f *filterStartsWith) Evaluate(node *Node) bool {
	v, ok := node.Props[f.property]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	return len(s) >= len(f.value) && s[:len(f.value)] == f.value
}

type filterStringContains struct {
	property string
	value    string
}

// Evaluate implements the VertexFilter interface.
func (f *filterStringContains) Evaluate(node *Node) bool {
	v, ok := node.Props[f.property]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	for i := 0; i <= len(s)-len(f.value); i++ {
		if s[i:i+len(f.value)] == f.value {
			return true
		}
	}
	return false
}

// ParseVertexFilter parses a JSON vertex filter expression string.
func ParseVertexFilter(jsonStr string) (VertexFilter, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return nil, fmt.Errorf("graphengine: invalid vertexFilter JSON: %w", err)
	}

	for op, val := range raw {
		switch op {
		case "andAll":
			return parseAndAll(val)
		case "orAll":
			return parseOrAll(val)
		case "equals":
			return parseEquals(val)
		case "notEquals":
			return parseNotEquals(val)
		case "greaterThan":
			return parseComparison(val, "greaterThan")
		case "greaterThanOrEquals":
			return parseComparison(val, "greaterThanOrEquals")
		case "lessThan":
			return parseComparison(val, "lessThan")
		case "lessThanOrEquals":
			return parseComparison(val, "lessThanOrEquals")
		case "in":
			return parseIn(val)
		case "notIn":
			return parseNotIn(val)
		case "startsWith":
			return parseStartsWith(val)
		case "stringContains":
			return parseStringContains(val)
		default:
			return nil, fmt.Errorf("graphengine: unknown vertexFilter operator: %s", op)
		}
	}

	return nil, fmt.Errorf("graphengine: empty vertexFilter expression")
}

func parseAndAll(raw json.RawMessage) (VertexFilter, error) {
	var children []json.RawMessage
	if err := json.Unmarshal(raw, &children); err != nil {
		return nil, fmt.Errorf("graphengine: invalid andAll: %w", err)
	}
	filters := make([]VertexFilter, 0, len(children))
	for _, c := range children {
		f, err := ParseVertexFilter(string(c))
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return &filterAndAll{children: filters}, nil
}

func parseOrAll(raw json.RawMessage) (VertexFilter, error) {
	var children []json.RawMessage
	if err := json.Unmarshal(raw, &children); err != nil {
		return nil, fmt.Errorf("graphengine: invalid orAll: %w", err)
	}
	filters := make([]VertexFilter, 0, len(children))
	for _, c := range children {
		f, err := ParseVertexFilter(string(c))
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return &filterOrAll{children: filters}, nil
}

func parseEquals(raw json.RawMessage) (VertexFilter, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("graphengine: invalid equals: %w", err)
	}
	prop, ok := m["property"].(string)
	if !ok {
		return nil, fmt.Errorf("graphengine: equals missing property")
	}
	val, ok := m["value"]
	if !ok {
		return nil, fmt.Errorf("graphengine: equals missing value")
	}
	return &filterEquals{property: prop, value: val}, nil
}

func parseNotEquals(raw json.RawMessage) (VertexFilter, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("graphengine: invalid notEquals: %w", err)
	}
	prop, ok := m["property"].(string)
	if !ok {
		return nil, fmt.Errorf("graphengine: notEquals missing property")
	}
	val, ok := m["value"]
	if !ok {
		return nil, fmt.Errorf("graphengine: notEquals missing value")
	}
	return &filterNotEquals{property: prop, value: val}, nil
}

func parseComparison(raw json.RawMessage, op string) (VertexFilter, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("graphengine: invalid %s: %w", op, err)
	}
	prop, ok := m["property"].(string)
	if !ok {
		return nil, fmt.Errorf("graphengine: %s missing property", op)
	}
	val, ok := m["value"]
	if !ok {
		return nil, fmt.Errorf("graphengine: %s missing value", op)
	}
	return &filterComparison{property: prop, value: toFloat64(val), op: op}, nil
}

func parseIn(raw json.RawMessage) (VertexFilter, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("graphengine: invalid in: %w", err)
	}
	prop, ok := m["property"].(string)
	if !ok {
		return nil, fmt.Errorf("graphengine: in missing property")
	}
	valRaw, ok := m["value"]
	if !ok {
		return nil, fmt.Errorf("graphengine: in missing value")
	}
	arr, ok := valRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("graphengine: in value must be an array")
	}
	return &filterIn{property: prop, values: arr}, nil
}

func parseNotIn(raw json.RawMessage) (VertexFilter, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("graphengine: invalid notIn: %w", err)
	}
	prop, ok := m["property"].(string)
	if !ok {
		return nil, fmt.Errorf("graphengine: notIn missing property")
	}
	valRaw, ok := m["value"]
	if !ok {
		return nil, fmt.Errorf("graphengine: notIn missing value")
	}
	arr, ok := valRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("graphengine: notIn value must be an array")
	}
	return &filterNotIn{property: prop, values: arr}, nil
}

func parseStartsWith(raw json.RawMessage) (VertexFilter, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("graphengine: invalid startsWith: %w", err)
	}
	prop, ok := m["property"].(string)
	if !ok {
		return nil, fmt.Errorf("graphengine: startsWith missing property")
	}
	val, ok := m["value"].(string)
	if !ok {
		return nil, fmt.Errorf("graphengine: startsWith value must be a string")
	}
	return &filterStartsWith{property: prop, value: val}, nil
}

func parseStringContains(raw json.RawMessage) (VertexFilter, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("graphengine: invalid stringContains: %w", err)
	}
	prop, ok := m["property"].(string)
	if !ok {
		return nil, fmt.Errorf("graphengine: stringContains missing property")
	}
	val, ok := m["value"].(string)
	if !ok {
		return nil, fmt.Errorf("graphengine: stringContains value must be a string")
	}
	return &filterStringContains{property: prop, value: val}, nil
}

func compareValues(a, b interface{}) int {
	aStr, aIsStr := a.(string)
	bStr, bIsStr := b.(string)
	if aIsStr && bIsStr {
		if aStr < bStr {
			return -1
		}
		if aStr > bStr {
			return 1
		}
		return 0
	}
	af := toFloat64(a)
	bf := toFloat64(b)
	if af < bf {
		return -1
	}
	if af > bf {
		return 1
	}
	return 0
}

func toFloat64(v interface{}) float64 {
	switch f := v.(type) {
	case float64:
		return f
	case int:
		return float64(f)
	case int64:
		return float64(f)
	case uint64:
		return float64(f)
	case string:
		if val, err := strconv.ParseFloat(f, 64); err == nil {
			return val
		}
		return 0
	case bool:
		if f {
			return 1
		}
		return 0
	default:
		return 0
	}
}
