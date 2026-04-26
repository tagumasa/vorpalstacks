package graphengine

// traversal.go implements graph traversal algorithms: breadth-first search,
// shortest path (unweighted and Dijkstra-weighted), connected components, and
// topological sort. All traversals use the lightweight adjacency lists
// (loadProps=false) unless the caller explicitly requires full edge data.

import (
	"container/heap"
	"container/list"
	"encoding/json"
	"fmt"
	"reflect"
)

// BFS performs a breadth-first traversal from startID up to maxDepth levels.
// The visitor callback is invoked for each discovered node; returning false
// stops the traversal early. For Both direction, outgoing edges resolve their
// To node and incoming edges resolve their From node as the next hop.
func (d *DB) BFS(startID NodeID, maxDepth int, direction Direction, edgeFilter EdgeFilter, visitor Visitor) error {
	if d.closed.Load() {
		return fmt.Errorf("graphengine: database is closed")
	}

	startNode, err := d.GetNode(startID)
	if err != nil {
		return err
	}

	visited := make(map[NodeID]bool)
	visited[startID] = true

	pathTo := make(map[NodeID][]*Edge)
	pathTo[startID] = nil

	if !visitor(&TraversalResult{Node: startNode, Depth: 0, Path: nil}) {
		return nil
	}

	type queueEntry struct {
		nodeID NodeID
		depth  int
	}

	queue := list.New()
	queue.PushBack(queueEntry{nodeID: startID, depth: 0})

	for queue.Len() > 0 {
		front := queue.Front()
		queue.Remove(front)
		entry := front.Value.(queueEntry)

		if maxDepth >= 0 && entry.depth >= maxDepth {
			continue
		}

		edges, err := d.getEdgesForNode(entry.nodeID, direction, false, "")
		if err != nil {
			return err
		}

		for _, edge := range edges {
			if edgeFilter != nil && !edgeFilter(edge) {
				continue
			}

			targetID := edge.To
			if direction == Incoming {
				targetID = edge.From
			} else if direction == Both {
				targetID = edge.To
				if edge.To == entry.nodeID {
					targetID = edge.From
				}
			}

			if visited[targetID] {
				continue
			}
			visited[targetID] = true

			parentPath := pathTo[entry.nodeID]
			edgePath := make([]*Edge, len(parentPath)+1)
			copy(edgePath, parentPath)
			edgePath[len(parentPath)] = edge
			pathTo[targetID] = edgePath

			targetNode, err := d.GetNode(targetID)
			if err != nil {
				continue
			}

			newDepth := entry.depth + 1
			result := &TraversalResult{
				Node:  targetNode,
				Depth: newDepth,
				Path:  edgePath,
			}

			if !visitor(result) {
				return nil
			}

			queue.PushBack(queueEntry{nodeID: targetID, depth: newDepth})
		}
	}

	return nil
}

// BFSCollect is a convenience wrapper that collects all BFS results into a
// slice.
func (d *DB) BFSCollect(startID NodeID, maxDepth int, direction Direction) ([]*TraversalResult, error) {
	var results []*TraversalResult

	err := d.BFS(startID, maxDepth, direction, nil, func(r *TraversalResult) bool {
		results = append(results, r)
		return true
	})

	return results, err
}

// ShortestPath returns the shortest path between two nodes by hop count.
// Only outgoing edges are traversed. Returns nil, nil if no path exists.
func (d *DB) ShortestPath(fromID, toID NodeID) (*PathResult, error) {
	return d.ShortestPathLabeled(fromID, toID, "")
}

// ShortestPathLabeled returns the shortest path between two nodes, optionally
// restricted to edges with the given label. An empty label matches all edges.
func (d *DB) ShortestPathLabeled(fromID, toID NodeID, label string) (*PathResult, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	if fromID == toID {
		node, err := d.GetNode(fromID)
		if err != nil {
			return nil, err
		}
		return &PathResult{
			Nodes: []*Node{node},
			Edges: nil,
			Cost:  0,
		}, nil
	}

	type bfsEntry struct {
		nodeID NodeID
		path   []*Edge
	}

	visited := make(map[NodeID]bool)
	visited[fromID] = true

	queue := list.New()
	queue.PushBack(bfsEntry{nodeID: fromID, path: nil})

	for queue.Len() > 0 {
		front := queue.Front()
		queue.Remove(front)
		entry := front.Value.(bfsEntry)

		edges, err := d.getEdgesForNode(entry.nodeID, Outgoing, false, label)
		if err != nil {
			return nil, err
		}

		for _, edge := range edges {
			if visited[edge.To] {
				continue
			}
			visited[edge.To] = true

			newPath := make([]*Edge, len(entry.path)+1)
			copy(newPath, entry.path)
			newPath[len(entry.path)] = edge

			if edge.To == toID {
				return d.buildPathResult(fromID, newPath, float64(len(newPath)))
			}

			queue.PushBack(bfsEntry{nodeID: edge.To, path: newPath})
		}
	}

	return nil, nil
}

// ShortestPathWeighted returns the shortest path between two nodes using
// Dijkstra's algorithm. Edge weights are read from the property named
// weightProp; edges lacking that property fall back to defaultWeight. The
// cost in the returned PathResult reflects the total weighted distance.
func (d *DB) ShortestPathWeighted(fromID, toID NodeID, weightProp string, defaultWeight float64) (*PathResult, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	if fromID == toID {
		node, err := d.GetNode(fromID)
		if err != nil {
			return nil, err
		}
		return &PathResult{
			Nodes: []*Node{node},
			Edges: nil,
			Cost:  0,
		}, nil
	}

	dist := make(map[NodeID]float64)
	prevEdge := make(map[NodeID]*Edge)

	dist[fromID] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &pqItem{nodeID: fromID, cost: 0})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*pqItem)
		current := item.nodeID

		if current == toID {
			break
		}

		if item.cost > dist[current] {
			continue
		}

		edges, err := d.getEdgesForNode(current, Outgoing, true, "")
		if err != nil {
			return nil, err
		}

		for _, edge := range edges {
			weight := defaultWeight
			if w, ok := edge.Props[weightProp]; ok {
				switch v := w.(type) {
				case float64:
					weight = v
				case int:
					weight = float64(v)
				case int64:
					weight = float64(v)
				case json.Number:
					f, err := v.Float64()
					if err != nil {
						break
					}
					weight = f
				default:
					rv := reflect.ValueOf(w)
					switch rv.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						weight = float64(rv.Int())
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						weight = float64(rv.Uint())
					case reflect.Float32, reflect.Float64:
						weight = rv.Float()
					}
				}
			}

			if weight < 0 {
				fmt.Printf("[graphengine] warning: negative edge weight %.4f skipped in Dijkstra traversal\n", weight)
				continue
			}

			newDist := dist[current] + weight
			if existing, ok := dist[edge.To]; !ok || newDist < existing {
				dist[edge.To] = newDist
				prevEdge[edge.To] = edge
				heap.Push(pq, &pqItem{nodeID: edge.To, cost: newDist})
			}
		}
	}

	if _, ok := dist[toID]; !ok {
		return nil, nil
	}

	var pathEdges []*Edge
	current := toID
	for current != fromID {
		edge, ok := prevEdge[current]
		if !ok {
			return nil, nil
		}
		pathEdges = append([]*Edge{edge}, pathEdges...)
		current = edge.From
	}

	return d.buildPathResult(fromID, pathEdges, dist[toID])
}

// HasPath reports whether a path exists between two nodes.
func (d *DB) HasPath(fromID, toID NodeID) (bool, error) {
	path, err := d.ShortestPath(fromID, toID)
	if err != nil {
		return false, err
	}
	return path != nil, nil
}

// ConnectedComponents returns all weakly connected components in the graph.
// Each component is a list of node IDs. Isolated nodes (no edges) are
// returned as single-element components.
func (d *DB) ConnectedComponents() ([][]NodeID, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	visited := make(map[NodeID]bool)
	var components [][]NodeID

	err := d.ForEachNode(func(node *Node) error {
		if visited[node.ID] {
			return nil
		}

		var component []NodeID
		queue := list.New()
		queue.PushBack(node.ID)
		visited[node.ID] = true

		for queue.Len() > 0 {
			front := queue.Front()
			queue.Remove(front)
			currentID := front.Value.(NodeID)
			component = append(component, currentID)

			edges, err := d.getEdgesForNode(currentID, Both, false, "")
			if err != nil {
				return err
			}

			for _, edge := range edges {
				targetID := edge.To
				if targetID == currentID {
					targetID = edge.From
				}
				if !visited[targetID] {
					visited[targetID] = true
					queue.PushBack(targetID)
				}
			}
		}

		components = append(components, component)
		return nil
	})

	return components, err
}

// TopologicalSort returns nodes in topological order using Kahn's algorithm.
// Only outgoing edges are considered (directed acyclic graph). Returns an
// error if the graph contains cycles.
func (d *DB) TopologicalSort() ([]NodeID, error) {
	if d.closed.Load() {
		return nil, fmt.Errorf("graphengine: database is closed")
	}

	inDegree := make(map[NodeID]int)
	adjacency := make(map[NodeID][]NodeID)

	err := d.ForEachNode(func(node *Node) error {
		if _, exists := inDegree[node.ID]; !exists {
			inDegree[node.ID] = 0
		}

		edges, err := d.getEdgesForNode(node.ID, Outgoing, false, "")
		if err != nil {
			return err
		}
		for _, e := range edges {
			inDegree[e.To]++
			adjacency[node.ID] = append(adjacency[node.ID], e.To)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	totalNodes := len(inDegree)

	queue := list.New()
	for id, deg := range inDegree {
		if deg == 0 {
			queue.PushBack(id)
		}
	}

	var sorted []NodeID
	for queue.Len() > 0 {
		front := queue.Front()
		queue.Remove(front)
		nodeID := front.Value.(NodeID)
		sorted = append(sorted, nodeID)

		for _, to := range adjacency[nodeID] {
			inDegree[to]--
			if inDegree[to] == 0 {
				queue.PushBack(to)
			}
		}
	}

	if len(sorted) != totalNodes {
		return nil, fmt.Errorf("graphengine: graph contains cycles, topological sort not possible")
	}

	return sorted, nil
}

// buildPathResult materialises the full PathResult by resolving all node IDs
// in the path to Node structs. The cost parameter allows callers to supply
// either hop count or weighted distance.
func (d *DB) buildPathResult(fromID NodeID, pathEdges []*Edge, cost float64) (*PathResult, error) {
	result := &PathResult{
		Edges: pathEdges,
		Cost:  cost,
	}

	nodeIDs := make([]NodeID, 0, len(pathEdges)+1)
	nodeIDs = append(nodeIDs, fromID)
	for _, e := range pathEdges {
		nodeIDs = append(nodeIDs, e.To)
	}

	result.Nodes = make([]*Node, 0, len(nodeIDs))
	for _, id := range nodeIDs {
		node, err := d.GetNode(id)
		if err != nil {
			return nil, err
		}
		result.Nodes = append(result.Nodes, node)
	}

	return result, nil
}

// priorityQueue implements container/heap for Dijkstra's algorithm.
// Each item tracks its node ID, accumulated cost, and current heap index
// for O(log n) decrease-key via re-insertion.

type pqItem struct {
	nodeID NodeID
	cost   float64
	index  int
}

type priorityQueue []*pqItem

// Len implements heap.Interface.
func (pq priorityQueue) Len() int { return len(pq) }

// Less implements heap.Interface.
func (pq priorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }

// Swap implements heap.Interface.
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push implements heap.Interface.
func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqItem)
	item.index = n
	*pq = append(*pq, item)
}

// Pop implements heap.Interface.
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[:n-1]
	return item
}
