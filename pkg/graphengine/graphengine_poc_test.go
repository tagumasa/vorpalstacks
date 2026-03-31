package graphengine

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	pebblev2 "github.com/cockroachdb/pebble/v2"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_test_%d_%d", time.Now().UnixNano(), rand.Intn(10000)))
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func openTestDB(t *testing.T) *DB {
	t.Helper()
	dir := tempDir(t)
	db, err := Open(dir, DefaultOptions())
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestAddAndGetNode(t *testing.T) {
	db := openTestDB(t)

	id, err := db.AddNode([]string{"Person"}, Props{"name": "Alice", "age": 30})
	if err != nil {
		t.Fatalf("AddNode failed: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero node ID")
	}

	node, err := db.GetNode(id)
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}
	if node.ID != id {
		t.Fatalf("expected ID %d, got %d", id, node.ID)
	}
	if len(node.Labels) != 1 || node.Labels[0] != "Person" {
		t.Fatalf("expected label Person, got %v", node.Labels)
	}
	if node.GetString("name") != "Alice" {
		t.Fatalf("expected name Alice, got %s", node.GetString("name"))
	}
	if node.GetFloat("age") != 30 {
		t.Fatalf("expected age 30, got %f", node.GetFloat("age"))
	}
}

func TestNodeNotFound(t *testing.T) {
	db := openTestDB(t)

	_, err := db.GetNode(999)
	if err == nil {
		t.Fatal("expected error for non-existent node")
	}
}

func TestUpdateNode(t *testing.T) {
	db := openTestDB(t)

	id, _ := db.AddNode([]string{"Person"}, Props{"name": "Alice"})

	err := db.UpdateNode(id, Props{"age": 31})
	if err != nil {
		t.Fatalf("UpdateNode failed: %v", err)
	}

	node, _ := db.GetNode(id)
	if node.GetString("name") != "Alice" {
		t.Fatalf("original prop lost: expected Alice, got %s", node.GetString("name"))
	}
	if node.GetFloat("age") != 31 {
		t.Fatalf("expected age 31, got %f", node.GetFloat("age"))
	}
}

func TestDeleteNode(t *testing.T) {
	db := openTestDB(t)

	alice, _ := db.AddNode([]string{"Person"}, Props{"name": "Alice"})
	bob, _ := db.AddNode([]string{"Person"}, Props{"name": "Bob"})
	db.AddEdge(alice, bob, "knows", nil)

	err := db.DeleteNode(alice)
	if err != nil {
		t.Fatalf("DeleteNode failed: %v", err)
	}

	exists, _ := db.NodeExists(alice)
	if exists {
		t.Fatal("node should be deleted")
	}

	stats := db.Stats()
	if stats.NodeCount != 1 {
		t.Fatalf("expected 1 node, got %d", stats.NodeCount)
	}
	if stats.EdgeCount != 0 {
		t.Fatalf("expected 0 edges, got %d", stats.EdgeCount)
	}
}

func TestAddAndGetEdge(t *testing.T) {
	db := openTestDB(t)

	alice, _ := db.AddNode([]string{"Person"}, Props{"name": "Alice"})
	bob, _ := db.AddNode([]string{"Person"}, Props{"name": "Bob"})

	edgeID, err := db.AddEdge(alice, bob, "knows", Props{"since": 2024})
	if err != nil {
		t.Fatalf("AddEdge failed: %v", err)
	}
	if edgeID == 0 {
		t.Fatal("expected non-zero edge ID")
	}

	edge, err := db.GetEdge(edgeID)
	if err != nil {
		t.Fatalf("GetEdge failed: %v", err)
	}
	if edge.From != alice || edge.To != bob {
		t.Fatalf("expected edge %d->%d, got %d->%d", alice, bob, edge.From, edge.To)
	}
	if edge.Label != "knows" {
		t.Fatalf("expected label knows, got %s", edge.Label)
	}
}

func TestOutEdges(t *testing.T) {
	db := openTestDB(t)

	alice, _ := db.AddNode(nil, Props{"name": "Alice"})
	bob, _ := db.AddNode(nil, Props{"name": "Bob"})
	carol, _ := db.AddNode(nil, Props{"name": "Carol"})

	db.AddEdge(alice, bob, "knows", nil)
	db.AddEdge(alice, carol, "knows", nil)
	db.AddEdge(bob, carol, "knows", nil)

	edges, err := db.OutEdges(alice)
	if err != nil {
		t.Fatalf("OutEdges failed: %v", err)
	}
	if len(edges) != 2 {
		t.Fatalf("expected 2 outgoing edges, got %d", len(edges))
	}

	edges, err = db.InEdges(carol)
	if err != nil {
		t.Fatalf("InEdges failed: %v", err)
	}
	if len(edges) != 2 {
		t.Fatalf("expected 2 incoming edges, got %d", len(edges))
	}
}

func TestDeleteEdge(t *testing.T) {
	db := openTestDB(t)

	alice, _ := db.AddNode(nil, nil)
	bob, _ := db.AddNode(nil, nil)
	edgeID, _ := db.AddEdge(alice, bob, "knows", nil)

	err := db.DeleteEdge(edgeID)
	if err != nil {
		t.Fatalf("DeleteEdge failed: %v", err)
	}

	_, err = db.GetEdge(edgeID)
	if err == nil {
		t.Fatal("expected error for deleted edge")
	}

	stats := db.Stats()
	if stats.EdgeCount != 0 {
		t.Fatalf("expected 0 edges, got %d", stats.EdgeCount)
	}
}

func TestFindByLabel(t *testing.T) {
	db := openTestDB(t)

	db.AddNode([]string{"Person"}, Props{"name": "Alice"})
	db.AddNode([]string{"Person"}, Props{"name": "Bob"})
	db.AddNode([]string{"City"}, Props{"name": "London"})

	ids, err := db.FindByLabel("Person")
	if err != nil {
		t.Fatalf("FindByLabel failed: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 Person nodes, got %d", len(ids))
	}

	ids, err = db.FindByLabel("City")
	if err != nil {
		t.Fatalf("FindByLabel failed: %v", err)
	}
	if len(ids) != 1 {
		t.Fatalf("expected 1 City node, got %d", len(ids))
	}
}

func TestForEachNode(t *testing.T) {
	db := openTestDB(t)

	db.AddNode(nil, Props{"name": "Alice"})
	db.AddNode(nil, Props{"name": "Bob"})
	db.AddNode(nil, Props{"name": "Carol"})

	count := 0
	err := db.ForEachNode(func(n *Node) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("ForEachNode failed: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 nodes, got %d", count)
	}
}

func TestGetAdjacentNodes(t *testing.T) {
	db := openTestDB(t)

	alice, _ := db.AddNode(nil, Props{"name": "Alice"})
	bob, _ := db.AddNode(nil, Props{"name": "Bob"})
	carol, _ := db.AddNode(nil, Props{"name": "Carol"})

	db.AddEdge(alice, bob, "knows", nil)
	db.AddEdge(alice, carol, "knows", nil)

	nodes, err := db.GetAdjacentNodes(alice, Outgoing)
	if err != nil {
		t.Fatalf("GetAdjacentNodes failed: %v", err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 adjacent nodes, got %d", len(nodes))
	}
}

func TestBFS(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, Props{"name": "A"})
	b, _ := db.AddNode(nil, Props{"name": "B"})
	c, _ := db.AddNode(nil, Props{"name": "C"})
	d, _ := db.AddNode(nil, Props{"name": "D"})

	db.AddEdge(a, b, "", nil)
	db.AddEdge(b, c, "", nil)
	db.AddEdge(c, d, "", nil)

	results, err := db.BFSCollect(a, 2, Outgoing)
	if err != nil {
		t.Fatalf("BFS failed: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 nodes (A, B, C) at depth 0-2, got %d", len(results))
	}

	names := make(map[string]bool)
	for _, r := range results {
		names[r.Node.GetString("name")] = true
	}
	for _, n := range []string{"A", "B", "C"} {
		if !names[n] {
			t.Fatalf("expected node %s in BFS results", n)
		}
	}
	if names["D"] {
		t.Fatal("D should not be in BFS results at depth 2")
	}
}

func TestBFSEdgeFilter(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	c, _ := db.AddNode(nil, nil)

	db.AddEdge(a, b, "friend", nil)
	db.AddEdge(a, c, "enemy", nil)

	var results []*TraversalResult
	err := db.BFS(a, 1, Outgoing, func(e *Edge) bool {
		return e.Label == "friend"
	}, func(r *TraversalResult) bool {
		results = append(results, r)
		return true
	})
	if err != nil {
		t.Fatalf("BFS with filter failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 nodes (a + b), got %d", len(results))
	}
}

func TestShortestPath(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, Props{"name": "A"})
	b, _ := db.AddNode(nil, Props{"name": "B"})
	c, _ := db.AddNode(nil, Props{"name": "C"})
	d, _ := db.AddNode(nil, Props{"name": "D"})

	db.AddEdge(a, b, "", nil)
	db.AddEdge(b, c, "", nil)
	db.AddEdge(a, d, "", nil)

	path, err := db.ShortestPath(a, c)
	if err != nil {
		t.Fatalf("ShortestPath failed: %v", err)
	}
	if path == nil {
		t.Fatal("expected path to exist")
	}
	if path.Cost != 2 {
		t.Fatalf("expected cost 2, got %f", path.Cost)
	}
	if len(path.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(path.Nodes))
	}
	if path.Nodes[0].GetString("name") != "A" || path.Nodes[2].GetString("name") != "C" {
		t.Fatal("wrong path endpoints")
	}
}

func TestShortestPathWeighted(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	c, _ := db.AddNode(nil, nil)

	db.AddEdge(a, b, "", Props{"weight": 10.0})
	db.AddEdge(a, c, "", Props{"weight": 1.0})
	db.AddEdge(c, b, "", Props{"weight": 1.0})

	path, err := db.ShortestPathWeighted(a, b, "weight", 1.0)
	if err != nil {
		t.Fatalf("ShortestPathWeighted failed: %v", err)
	}
	if path == nil {
		t.Fatal("expected path to exist")
	}
	if path.Cost != 2.0 {
		t.Fatalf("expected cost 2.0 (a->c->b), got %f", path.Cost)
	}
}

func TestHasPath(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	c, _ := db.AddNode(nil, nil)

	db.AddEdge(a, b, "", nil)

	has, err := db.HasPath(a, b)
	if err != nil || !has {
		t.Fatal("expected path a->b to exist")
	}

	has, err = db.HasPath(a, c)
	if err != nil || has {
		t.Fatal("expected no path a->c")
	}
}

func TestTopologicalSort(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, Props{"name": "A"})
	b, _ := db.AddNode(nil, Props{"name": "B"})
	c, _ := db.AddNode(nil, Props{"name": "C"})

	db.AddEdge(a, b, "", nil)
	db.AddEdge(b, c, "", nil)

	sorted, err := db.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}
	if len(sorted) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(sorted))
	}

	pos := make(map[NodeID]int)
	for i, id := range sorted {
		pos[id] = i
	}
	if pos[a] >= pos[b] || pos[b] >= pos[c] {
		t.Fatal("wrong topological order")
	}
}

func TestConnectedComponents(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	c, _ := db.AddNode(nil, nil)
	d, _ := db.AddNode(nil, nil)

	db.AddEdge(a, b, "", nil)
	db.AddEdge(c, d, "", nil)

	components, err := db.ConnectedComponents()
	if err != nil {
		t.Fatalf("ConnectedComponents failed: %v", err)
	}
	if len(components) != 2 {
		t.Fatalf("expected 2 components, got %d", len(components))
	}
}

func TestStats(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	db.AddEdge(a, b, "", nil)

	stats := db.Stats()
	if stats.NodeCount != 2 {
		t.Fatalf("expected 2 nodes, got %d", stats.NodeCount)
	}
	if stats.EdgeCount != 1 {
		t.Fatalf("expected 1 edge, got %d", stats.EdgeCount)
	}
}

func TestGetAdjacentNodesBoth(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, Props{"name": "A"})
	b, _ := db.AddNode(nil, Props{"name": "B"})
	c, _ := db.AddNode(nil, Props{"name": "C"})

	db.AddEdge(a, b, "knows", nil)
	db.AddEdge(c, a, "knows", nil)

	nodes, err := db.GetAdjacentNodes(a, Both)
	if err != nil {
		t.Fatalf("GetAdjacentNodes(Both) failed: %v", err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 adjacent nodes (B and C), got %d", len(nodes))
	}
	names := make(map[string]bool)
	for _, n := range nodes {
		names[n.GetString("name")] = true
	}
	if !names["B"] || !names["C"] {
		t.Fatalf("expected B and C as adjacent, got %v", names)
	}
}

func TestAddEdgeBatchDoesNotMutateInput(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)

	edges := make([]Edge, 2)
	edges[0] = Edge{From: a, To: b, Label: "knows", Props: Props{"w": 1.0}}
	edges[1] = Edge{From: a, To: b, Label: "likes", Props: Props{"w": 2.0}}

	ids, err := db.AddEdgeBatch(edges)
	if err != nil {
		t.Fatalf("AddEdgeBatch failed: %v", err)
	}

	if edges[0].ID != 0 || edges[1].ID != 0 {
		t.Fatal("AddEdgeBatch should not mutate input slice IDs")
	}

	for _, id := range ids {
		if id == 0 {
			t.Fatal("expected non-zero allocated edge ID")
		}
	}
}

func TestShortestPathWeightedCost(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	c, _ := db.AddNode(nil, nil)

	db.AddEdge(a, b, "", Props{"weight": 5.0})
	db.AddEdge(b, c, "", Props{"weight": 3.0})

	path, err := db.ShortestPathWeighted(a, c, "weight", 1.0)
	if err != nil {
		t.Fatalf("ShortestPathWeighted failed: %v", err)
	}
	if path == nil {
		t.Fatal("expected path to exist")
	}
	if path.Cost != 8.0 {
		t.Fatalf("expected weighted cost 8.0, got %f", path.Cost)
	}
}

func TestCloseAndReopen(t *testing.T) {
	dir := tempDir(t)

	db, err := Open(dir, DefaultOptions())
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	id, _ := db.AddNode([]string{"Person"}, Props{"name": "Alice"})
	db.Close()

	db2, err := Open(dir, DefaultOptions())
	if err != nil {
		t.Fatalf("Reopen failed: %v", err)
	}
	defer db2.Close()

	node, err := db2.GetNode(id)
	if err != nil {
		t.Fatalf("GetNode after reopen failed: %v", err)
	}
	if node.GetString("name") != "Alice" {
		t.Fatalf("expected Alice, got %s", node.GetString("name"))
	}

	stats := db2.Stats()
	if stats.NodeCount != 1 {
		t.Fatalf("expected 1 node after reopen, got %d", stats.NodeCount)
	}
}

// --- Benchmarks ---

func BenchmarkAddNode(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_add_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.AddNode([]string{"Person"}, Props{
			"name": fmt.Sprintf("node_%d", i),
			"age":  i % 100,
		})
	}
}

func BenchmarkAddNodeBatch100(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_batch_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	items := make([]struct {
		Labels []string
		Props  Props
	}, 100)
	for i := range items {
		items[i] = struct {
			Labels []string
			Props  Props
		}{
			Labels: []string{"Person"},
			Props: Props{
				"name": fmt.Sprintf("node_%d", i),
				"age":  i % 100,
			},
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.AddNodeBatch(items)
	}
}

func BenchmarkAddEdge(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_edge_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	a, _ := db.AddNode(nil, nil)
	nodeB, _ := db.AddNode(nil, nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.AddEdge(a, nodeB, "knows", Props{"weight": float64(i)})
	}
}

func BenchmarkAddEdgeBatch100(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_edgebatch_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	a, _ := db.AddNode(nil, nil)
	nodeB, _ := db.AddNode(nil, nil)

	edges := make([]Edge, 100)
	for i := range edges {
		edges[i] = Edge{From: a, To: nodeB, Label: "knows", Props: Props{"weight": float64(i)}}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.AddEdgeBatch(edges)
	}
}

func BenchmarkGetNode(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_get_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	id, _ := db.AddNode([]string{"Person"}, Props{"name": "Alice", "age": 30})

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.GetNode(id)
	}
}

func BenchmarkOutEdges(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_out_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	a, _ := db.AddNode(nil, nil)
	for i := 0; i < 100; i++ {
		target, _ := db.AddNode(nil, nil)
		db.AddEdge(a, target, "knows", nil)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.OutEdges(a)
	}
}

func BenchmarkBFS10KNodes(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_bfs_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	nodes := make([]NodeID, 10000)
	for i := range nodes {
		nodes[i], _ = db.AddNode(nil, Props{"id": i})
	}
	for i := 1; i < len(nodes); i++ {
		db.AddEdge(nodes[i-1], nodes[i], "next", nil)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.BFSCollect(nodes[0], 3, Outgoing)
	}
}

func BenchmarkShortestPath10KNodes(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_sp_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	nodes := make([]NodeID, 10000)
	for i := range nodes {
		nodes[i], _ = db.AddNode(nil, nil)
	}
	for i := 1; i < len(nodes); i++ {
		db.AddEdge(nodes[i-1], nodes[i], "", nil)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.ShortestPath(nodes[0], nodes[9999])
	}
}

func BenchmarkSocialGraphTraversal(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_social_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	numNodes := 10000
	numEdges := 50000

	nodes := make([]NodeID, numNodes)
	for i := range nodes {
		nodes[i], _ = db.AddNode([]string{"Person"}, Props{"name": fmt.Sprintf("user_%d", i)})
	}

	edges := make([]Edge, numEdges)
	for i := range edges {
		from := nodes[rand.Intn(numNodes)]
		to := nodes[rand.Intn(numNodes)]
		if from == to {
			to = nodes[(int(from)+1)%numNodes]
		}
		edges[i] = Edge{From: from, To: to, Label: "follows"}
	}
	db.AddEdgeBatch(edges)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.BFSCollect(nodes[0], 5, Outgoing)
	}
}

func BenchmarkPebbleBatchRaw(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_raw_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		batch := db.db.NewBatch()
		for j := 0; j < 300; j++ {
			key := make([]byte, 20)
			val := make([]byte, 50)
			binary.BigEndian.PutUint64(key[:8], uint64(j))
			batch.Set(key, val, nil)
		}
		batch.Apply(batch, pebblev2.NoSync)
		batch.Close()
	}
}

func BenchmarkMsgpackMarshal(b *testing.B) {
	props := Props{"name": "node_0", "age": 30}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		marshalProps(props)
	}
}

func BenchmarkBinaryEncoding(b *testing.B) {
	edge := Edge{ID: 1, From: 1, To: 2, Label: "knows", Props: Props{"weight": 1.0}}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = encodeEdgeData(&edge)
		encodeAdjValue(edge.To, edge.Label)
		encodeAdjValue(edge.From, edge.Label)
		edgeKey(edge.ID)
		adjOutKey(edge.From, edge.ID)
		adjInKey(edge.To, edge.ID)
	}
}

func BenchmarkJSONStructEncoding(b *testing.B) {
	edge := Edge{ID: 1, From: 1, To: 2, Label: "knows", Props: Props{"weight": 1.0}}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		json.Marshal(&edge)
		encodeAdjValue(edge.To, edge.Label)
		encodeAdjValue(edge.From, edge.Label)
		edgeKey(edge.ID)
		adjOutKey(edge.From, edge.ID)
		adjInKey(edge.To, edge.ID)
	}
}

func BenchmarkPOCStyleEdgeBatch(b *testing.B) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("graphengine_bench_poc_%d", time.Now().UnixNano()))
	defer os.RemoveAll(dir)

	db, _ := Open(dir, DefaultOptions())
	defer db.Close()

	a, _ := db.AddNode(nil, nil)
	nodeB, _ := db.AddNode(nil, nil)

	edges := make([]Edge, 100)
	for i := range edges {
		edges[i] = Edge{From: a, To: nodeB, Label: "knows", Props: Props{"weight": float64(i)}}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		batch := db.db.NewBatch()
		for j := range edges {
			edges[j].ID = EdgeID(j + 1)
			edgeJSON, _ := json.Marshal(&edges[j])
			batch.Set(edgeKey(edges[j].ID), edgeJSON, nil)
			batch.Set(adjOutKey(edges[j].From, edges[j].ID), encodeAdjValue(edges[j].To, edges[j].Label), nil)
			batch.Set(adjInKey(edges[j].To, edges[j].ID), encodeAdjValue(edges[j].From, edges[j].Label), nil)
		}
		batch.Apply(batch, pebblev2.NoSync)
		batch.Close()
	}
}
