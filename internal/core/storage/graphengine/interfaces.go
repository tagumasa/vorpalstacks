// Package graphengine provides an embedded property graph database backed by
// Pebble key-value storage, with support for nodes, edges, labels, property
// indexing, Cypher/Gremlin query traversal, and batch import/export.
package graphengine

// GraphReader provides read-only access to graph data, decoupled from the
// concrete Pebble-backed storage. Query executors (Cypher, Gremlin) depend
// on this interface rather than *DB directly, enabling in-memory mock
// implementations for testing.
type GraphReader interface {
	GetNode(id NodeID) (*Node, error)
	NodeExists(id NodeID) (bool, error)
	FindByLabel(label string) ([]NodeID, error)
	FindByProperty(prop string, value interface{}) ([]NodeID, error)
	HasUniqueConstraint(label, prop string) bool
	FindByUniqueConstraint(label, prop string, value interface{}) (NodeID, error)
	ForEachNode(fn func(*Node) error) error
	ForEachEdge(fn func(*Edge) error) error
	GetEdges(nodeID NodeID, dir Direction, labelFilter string) ([]*Edge, error)
	GetEdge(id EdgeID) (*Edge, error)
	ScanEdgesByType(label string, fn func(edge *Edge, src, dst *Node) error) error
	BFS(startID NodeID, maxDepth int, dir Direction, filter EdgeFilter, visit Visitor) error
	Stats() *GraphStats
	CountNodes() int64
	CountEdges() int64
	GetLabelCounts() (map[string]int64, error)
	GetRelCounts() (map[string]int64, error)
}

// GraphWriter provides write access to graph data, decoupled from the
// concrete Pebble-backed storage.
type GraphWriter interface {
	AddNode(labels []string, props Props) (NodeID, error)
	AddNodeBatch(items []struct {
		Labels []string
		Props  Props
	}) ([]NodeID, error)
	AddEdge(from, to NodeID, label string, props Props) (EdgeID, error)
	AddEdgeBatch(edges []Edge) ([]EdgeID, error)
	UpdateNode(id NodeID, props Props) error
	UpdateEdge(id EdgeID, props Props) error
	DeleteNode(id NodeID) error
	DeleteEdge(id EdgeID) error
	RemoveLabel(id NodeID, label string) error
	AddLabel(id NodeID, label string) error
	RemoveProperty(id NodeID, key string) error
	RemoveEdgeProperty(id EdgeID, key string) error
}

// GraphStore combines read and write access with lifecycle management.
// The concrete *DB type satisfies this interface implicitly.
type GraphStore interface {
	GraphReader
	GraphWriter
	Clear() error
	Close() error
	Dir() string
}

// GraphDDL provides data definition language operations for managing
// indexes and constraints. The concrete *DB type satisfies this interface.
type GraphDDL interface {
	CreateIndex(label, prop string) error
	ShowIndexes() ([]IndexInfo, error)
	DropIndex(label, prop string) error
	CreateUniqueConstraint(label, prop string) error
	ShowConstraints() ([]ConstraintInfo, error)
	DropUniqueConstraint(label, prop string) error
	HasUniqueConstraint(label, prop string) bool
}
