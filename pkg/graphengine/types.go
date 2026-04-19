package graphengine

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// NodeID is a unique identifier for a node, allocated as a monotonically
// increasing unsigned 64-bit integer. IDs start at 1; 0 is never allocated.
type NodeID uint64

// EdgeID is a unique identifier for an edge, allocated independently from
// NodeID using the same monotonically increasing scheme.
type EdgeID uint64

// Props represents a property map attached to a node or edge. Values are
// deserialised from JSON and may be any JSON-compatible type (string, number,
// bool, nil, slice, or nested map).
type Props map[string]interface{}

// Direction controls which adjacency lists are traversed when querying edges.
type Direction byte

const (
	Outgoing Direction = 0x01
	Incoming Direction = 0x02
	Both     Direction = 0x03
)

// Node represents a graph node with an optional set of labels and properties.
type Node struct {
	ID     NodeID   `json:"id"`
	Labels []string `json:"labels,omitempty"`
	Props  Props    `json:"props,omitempty"`
}

// GetString returns the property value for the given key as a string.
// Returns an empty string if the key does not exist or the value is not a string.
func (n *Node) GetString(key string) string {
	if v, ok := n.Props[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetFloat returns the property value for the given key as a float64.
// Handles all integer and floating-point JSON types, including json.Number.
// Returns 0 if the key does not exist or the value is not numeric.
func (n *Node) GetFloat(key string) float64 {
	if v, ok := n.Props[key]; ok {
		switch f := v.(type) {
		case float64:
			return f
		case float32:
			return float64(f)
		case int:
			return float64(f)
		case int8:
			return float64(f)
		case int16:
			return float64(f)
		case int32:
			return float64(f)
		case int64:
			return float64(f)
		case uint:
			return float64(f)
		case uint8:
			return float64(f)
		case uint16:
			return float64(f)
		case uint32:
			return float64(f)
		case uint64:
			return float64(f)
		case json.Number:
			val, _ := f.Float64()
			return val
		}
	}
	return 0
}

// Edge represents a directed edge between two nodes with a label and optional
// properties. The From and To fields denote the source and target nodes
// respectively.
type Edge struct {
	ID    EdgeID `json:"id"`
	From  NodeID `json:"from"`
	To    NodeID `json:"to"`
	Label string `json:"label"`
	Props Props  `json:"props,omitempty"`
}

// String returns a human-readable representation of the edge.
func (e *Edge) String() string {
	return fmt.Sprintf("(%d)--%s-->(%d)", e.From, e.Label, e.To)
}

// TraversalResult is yielded to a Visitor during BFS traversal, carrying the
// discovered node, its depth from the start node, and the edge path taken to
// reach it.
type TraversalResult struct {
	Node  *Node
	Depth int
	Path  []*Edge
}

// PathResult represents a path between two nodes, including all intermediate
// nodes, the edges traversed, and the total cost (hop count for unweighted,
// weighted sum for Dijkstra).
type PathResult struct {
	Nodes []*Node
	Edges []*Edge
	Cost  float64
}

// EdgeFilter is a function that returns true if an edge should be included in a traversal.
type EdgeFilter func(*Edge) bool

// Visitor is a callback invoked for each node discovered during a BFS traversal. Returning false stops the traversal early.
type Visitor func(result *TraversalResult) bool

// GraphStats holds high-level statistics about the graph.
type GraphStats struct {
	NodeCount     int64            `json:"node_count"`
	EdgeCount     int64            `json:"edge_count"`
	LabelCounts   map[string]int64 `json:"label_counts,omitempty"`
	RelCounts     map[string]int64 `json:"rel_counts,omitempty"`
	DiskSizeBytes int64            `json:"disk_size_bytes"`
}

// --- Binary encoding helpers (big-endian) ---

func encodeUint64(v uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	return buf
}

func decodeUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func encodeNodeID(id NodeID) []byte {
	return encodeUint64(uint64(id))
}

func decodeNodeID(b []byte) NodeID {
	return NodeID(decodeUint64(b))
}

func encodeEdgeID(id EdgeID) []byte {
	return encodeUint64(uint64(id))
}

func decodeEdgeID(b []byte) EdgeID {
	return EdgeID(decodeUint64(b))
}

func putUint16(buf []byte, v uint16) {
	buf[0] = byte(v >> 8)
	buf[1] = byte(v)
}

func getUint16(b []byte) uint16 {
	return uint16(b[0])<<8 | uint16(b[1])
}

func encodeAdjKey(nodeID NodeID, edgeID EdgeID) []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:8], uint64(nodeID))
	binary.BigEndian.PutUint64(buf[8:], uint64(edgeID))
	return buf
}

func decodeAdjKey(b []byte) (NodeID, EdgeID) {
	nodeID := NodeID(binary.BigEndian.Uint64(b[:8]))
	edgeID := EdgeID(binary.BigEndian.Uint64(b[8:]))
	return nodeID, edgeID
}

var emptyPropsJSON = []byte("{}")

func marshalProps(props Props) []byte {
	if len(props) == 0 {
		return emptyPropsJSON
	}
	b, err := json.Marshal(props)
	if err != nil {
		return emptyPropsJSON
	}
	return b
}

func unmarshalProps(data []byte) (Props, error) {
	var props Props
	if err := json.Unmarshal(data, &props); err != nil {
		return nil, err
	}
	if props == nil {
		props = make(Props)
	}
	return props, nil
}

// --- Node binary encoding ---
//
// Layout: [numLabels:u16] [labelLen:u16]... [label:bytes]... [jsonProps:bytes]
// Labels are stored with individual length prefixes to avoid a separator byte,
// allowing any label content including NUL. Properties are JSON-serialised as
// the trailing payload; when absent, an empty object "{}" is stored so that
// the decode side can distinguish zero labels from zero labels + empty props.

func encodeNodeDataRaw(labels []string, propsBytes []byte) ([]byte, error) {
	if len(labels) > 65535 {
		return nil, fmt.Errorf("graphengine: too many labels (%d, max 65535)", len(labels))
	}
	for i, l := range labels {
		if len(l) > 65535 {
			return nil, fmt.Errorf("graphengine: label %d too long (%d bytes, max 65535)", i, len(l))
		}
	}
	numLabels := uint16(len(labels))
	headerSize := int(2 + numLabels*2)
	labelDataSize := 0
	for _, l := range labels {
		labelDataSize += len(l)
	}
	totalSize := headerSize + labelDataSize + len(propsBytes)
	buf := make([]byte, totalSize)
	putUint16(buf[0:2], numLabels)
	offset := 2
	for _, label := range labels {
		putUint16(buf[offset:offset+2], uint16(len(label)))
		offset += 2
	}
	for _, label := range labels {
		copy(buf[offset:], label)
		offset += len(label)
	}
	copy(buf[offset:], propsBytes)
	return buf, nil
}

func encodeNodeData(labels []string, props Props) ([]byte, error) {
	return encodeNodeDataRaw(labels, marshalProps(props))
}

func decodeNodeData(data []byte) ([]string, Props, error) {
	if len(data) < 2 {
		return nil, nil, fmt.Errorf("graphengine: node data too short")
	}
	numLabels := int(getUint16(data[0:2]))
	offset := 2
	labelLens := make([]int, numLabels)
	for i := 0; i < numLabels; i++ {
		if offset+2 > len(data) {
			return nil, nil, fmt.Errorf("graphengine: node data truncated at label length")
		}
		labelLens[i] = int(getUint16(data[offset : offset+2]))
		offset += 2
	}
	labels := make([]string, numLabels)
	for i := 0; i < numLabels; i++ {
		if offset+labelLens[i] > len(data) {
			return nil, nil, fmt.Errorf("graphengine: node data truncated at label %d", i)
		}
		labels[i] = string(data[offset : offset+labelLens[i]])
		offset += labelLens[i]
	}
	var props Props
	if offset < len(data) {
		var err error
		props, err = unmarshalProps(data[offset:])
		if err != nil {
			return nil, nil, err
		}
	}
	return labels, props, nil
}

// --- Edge binary encoding ---
//
// Layout: [fromID:u64] [toID:u64] [labelLen:u16] [label:bytes] [jsonProps:bytes]
// Properties are serialised with json.Marshal (which uses sync.Pool internally
// for the encoder, making it competitive with hand-rolled binary for small maps).

func encodeEdgeDataRaw(e *Edge, propsBytes []byte) ([]byte, error) {
	if len(e.Label) > 65535 {
		return nil, fmt.Errorf("graphengine: edge label too long (%d bytes, max 65535)", len(e.Label))
	}
	labelLen := uint16(len(e.Label))
	headerSize := 8 + 8 + 2
	totalSize := headerSize + len(e.Label) + len(propsBytes)
	buf := make([]byte, totalSize)
	binary.BigEndian.PutUint64(buf[0:8], uint64(e.From))
	binary.BigEndian.PutUint64(buf[8:16], uint64(e.To))
	putUint16(buf[16:18], labelLen)
	copy(buf[18:], e.Label)
	copy(buf[18+len(e.Label):], propsBytes)
	return buf, nil
}

func encodeEdgeData(e *Edge) ([]byte, error) {
	return encodeEdgeDataRaw(e, marshalProps(e.Props))
}

func decodeEdgeData(edgeID EdgeID, data []byte) (*Edge, error) {
	if len(data) < 18 {
		return nil, fmt.Errorf("graphengine: edge data too short")
	}
	from := NodeID(binary.BigEndian.Uint64(data[0:8]))
	to := NodeID(binary.BigEndian.Uint64(data[8:16]))
	labelLen := int(getUint16(data[16:18]))
	if len(data) < 18+labelLen {
		return nil, fmt.Errorf("graphengine: edge data truncated at label")
	}
	label := string(data[18 : 18+labelLen])
	var props Props
	if 18+labelLen < len(data) {
		var err error
		props, err = unmarshalProps(data[18+labelLen:])
		if err != nil {
			return nil, err
		}
	}
	return &Edge{
		ID:    edgeID,
		From:  from,
		To:    to,
		Label: label,
		Props: props,
	}, nil
}

// --- Adjacency list encoding ---
//
// Adjacency values are lightweight: only the peer node ID and edge label are
// stored. Full edge properties are loaded lazily via GetEdge when loadProps
// is true, avoiding N+1 read amplification during BFS and topological sort.

func encodeAdjValue(targetID NodeID, label string) []byte {
	buf := make([]byte, 8+len(label))
	binary.BigEndian.PutUint64(buf[:8], uint64(targetID))
	copy(buf[8:], label)
	return buf
}

func decodeAdjValue(data []byte) (NodeID, string) {
	if len(data) < 8 {
		return 0, ""
	}
	targetID := NodeID(binary.BigEndian.Uint64(data[:8]))
	label := string(data[8:])
	return targetID, label
}
