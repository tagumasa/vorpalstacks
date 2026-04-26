package graphengine

import (
	"fmt"
	"testing"
)

func TestOutEdgesByLabel(t *testing.T) {
	db := openTestDB(t)

	alice, _ := db.AddNode(nil, Props{"name": "Alice"})
	bob, _ := db.AddNode(nil, Props{"name": "Bob"})
	carol, _ := db.AddNode(nil, Props{"name": "Carol"})

	db.AddEdge(alice, bob, "knows", nil)
	db.AddEdge(alice, carol, "likes", nil)

	edges, err := db.OutEdgesByLabel(alice, "knows")
	if err != nil {
		t.Fatalf("OutEdgesByLabel failed: %v", err)
	}
	if len(edges) != 1 {
		t.Fatalf("expected 1 knows edge, got %d", len(edges))
	}
	if edges[0].To != bob {
		t.Fatalf("expected edge to Bob, got %d", edges[0].To)
	}

	edges, err = db.OutEdgesByLabel(alice, "likes")
	if err != nil {
		t.Fatalf("OutEdgesByLabel failed: %v", err)
	}
	if len(edges) != 1 {
		t.Fatalf("expected 1 likes edge, got %d", len(edges))
	}

	edges, err = db.OutEdgesByLabel(alice, "nonexistent")
	if err != nil {
		t.Fatalf("OutEdgesByLabel failed: %v", err)
	}
	if len(edges) != 0 {
		t.Fatalf("expected 0 edges for nonexistent label, got %d", len(edges))
	}
}

func TestInEdgesByLabel(t *testing.T) {
	db := openTestDB(t)

	alice, _ := db.AddNode(nil, nil)
	bob, _ := db.AddNode(nil, nil)
	carol, _ := db.AddNode(nil, nil)

	db.AddEdge(alice, carol, "knows", nil)
	db.AddEdge(bob, carol, "likes", nil)

	edges, err := db.InEdgesByLabel(carol, "knows")
	if err != nil {
		t.Fatalf("InEdgesByLabel failed: %v", err)
	}
	if len(edges) != 1 {
		t.Fatalf("expected 1 knows edge, got %d", len(edges))
	}
	if edges[0].From != alice {
		t.Fatalf("expected edge from Alice, got %d", edges[0].From)
	}
}

func TestGetAdjacentNodesByLabel(t *testing.T) {
	db := openTestDB(t)

	alice, _ := db.AddNode(nil, Props{"name": "Alice"})
	bob, _ := db.AddNode(nil, Props{"name": "Bob"})
	carol, _ := db.AddNode(nil, Props{"name": "Carol"})

	db.AddEdge(alice, bob, "knows", nil)
	db.AddEdge(alice, carol, "likes", nil)

	nodes, err := db.GetAdjacentNodesByLabel(alice, Outgoing, "knows")
	if err != nil {
		t.Fatalf("GetAdjacentNodesByLabel failed: %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 adjacent node, got %d", len(nodes))
	}
	if nodes[0].GetString("name") != "Bob" {
		t.Fatalf("expected Bob, got %s", nodes[0].GetString("name"))
	}
}

func TestGetEdgesWithLabelFilter(t *testing.T) {
	db := openTestDB(t)

	alice, _ := db.AddNode(nil, nil)
	bob, _ := db.AddNode(nil, nil)
	carol, _ := db.AddNode(nil, nil)

	db.AddEdge(alice, bob, "knows", nil)
	db.AddEdge(alice, carol, "likes", nil)

	edges, err := db.GetEdges(alice, Outgoing, "knows")
	if err != nil {
		t.Fatalf("GetEdges with filter failed: %v", err)
	}
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
	if edges[0].To != bob {
		t.Fatal("wrong edge target")
	}

	edges, err = db.GetEdges(alice, Outgoing, "")
	if err != nil {
		t.Fatalf("GetEdges without filter failed: %v", err)
	}
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges without filter, got %d", len(edges))
	}
}

func TestForEachNodeWithLimit(t *testing.T) {
	db := openTestDB(t)

	for i := 0; i < 10; i++ {
		db.AddNode(nil, Props{"i": i})
	}

	count := 0
	err := db.ForEachNodeN(func(n *Node) error {
		count++
		return nil
	}, 3)
	if err != nil {
		t.Fatalf("ForEachNodeN with limit failed: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 nodes with limit=3, got %d", count)
	}

	count = 0
	err = db.ForEachNodeN(func(n *Node) error {
		count++
		return nil
	}, 0)
	if err != nil {
		t.Fatalf("ForEachNodeN with limit=0 failed: %v", err)
	}
	if count != 10 {
		t.Fatalf("expected 10 nodes with limit=0, got %d", count)
	}
}

func TestFindByPropertyInterface(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, Props{"age": 30})
	b, _ := db.AddNode(nil, Props{"age": 30})
	c, _ := db.AddNode(nil, Props{"age": 25})

	ids, err := db.FindByProperty("age", 30)
	if err != nil {
		t.Fatalf("FindByProperty failed: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 nodes with age=30, got %d (IDs: %v)", len(ids), ids)
	}

	ids, err = db.FindByProperty("age", 25)
	if err != nil {
		t.Fatalf("FindByProperty failed: %v", err)
	}
	if len(ids) != 1 {
		t.Fatalf("expected 1 node with age=25, got %d", len(ids))
	}

	ids, err = db.FindByProperty("age", 99)
	if err != nil {
		t.Fatalf("FindByProperty failed: %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected 0 nodes with age=99, got %d", len(ids))
	}

	_ = a
	_ = b
	_ = c
}

func TestUpdateEdge(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	eid, _ := db.AddEdge(a, b, "knows", Props{"since": 2020})

	err := db.UpdateEdge(eid, Props{"since": 2024, "weight": 1.5})
	if err != nil {
		t.Fatalf("UpdateEdge failed: %v", err)
	}

	edge, _ := db.GetEdge(eid)
	if edge.Props["since"].(float64) != 2024 {
		t.Fatalf("expected since=2024, got %v", edge.Props["since"])
	}
	if edge.Props["weight"].(float64) != 1.5 {
		t.Fatalf("expected weight=1.5, got %v", edge.Props["weight"])
	}
}

func TestUpdateEdgeNotFound(t *testing.T) {
	db := openTestDB(t)

	err := db.UpdateEdge(999, Props{"key": "val"})
	if err == nil {
		t.Fatal("expected error for non-existent edge")
	}
}

func TestRemoveLabel(t *testing.T) {
	db := openTestDB(t)

	id, _ := db.AddNode([]string{"Person", "Employee"}, Props{"name": "Alice"})

	err := db.RemoveLabel(id, "Employee")
	if err != nil {
		t.Fatalf("RemoveLabel failed: %v", err)
	}

	node, _ := db.GetNode(id)
	if len(node.Labels) != 1 || node.Labels[0] != "Person" {
		t.Fatalf("expected [Person], got %v", node.Labels)
	}

	ids, _ := db.FindByLabel("Employee")
	if len(ids) != 0 {
		t.Fatalf("expected 0 Employee nodes after remove, got %d", len(ids))
	}
}

func TestRemoveProperty(t *testing.T) {
	db := openTestDB(t)

	id, _ := db.AddNode(nil, Props{"name": "Alice", "age": 30})

	err := db.RemoveProperty(id, "age")
	if err != nil {
		t.Fatalf("RemoveProperty failed: %v", err)
	}

	node, _ := db.GetNode(id)
	if node.GetString("name") != "Alice" {
		t.Fatalf("name should be preserved, got %s", node.GetString("name"))
	}
	if _, ok := node.Props["age"]; ok {
		t.Fatal("age should be removed")
	}
}

func TestRemovePropertyNotFound(t *testing.T) {
	db := openTestDB(t)

	id, _ := db.AddNode(nil, Props{"name": "Alice"})

	err := db.RemoveProperty(id, "nonexistent")
	if err != nil {
		t.Fatalf("RemoveProperty on missing key should be no-op: %v", err)
	}
}

func TestRemoveEdgeProperty(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	eid, _ := db.AddEdge(a, b, "knows", Props{"since": 2020, "weight": 1.0})

	err := db.RemoveEdgeProperty(eid, "since")
	if err != nil {
		t.Fatalf("RemoveEdgeProperty failed: %v", err)
	}

	edge, _ := db.GetEdge(eid)
	if _, ok := edge.Props["since"]; ok {
		t.Fatal("since should be removed")
	}
	if edge.Props["weight"] != 1.0 {
		t.Fatalf("weight should be preserved, got %v", edge.Props["weight"])
	}
}

func TestForEachEdge(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	c, _ := db.AddNode(nil, nil)

	db.AddEdge(a, b, "knows", nil)
	db.AddEdge(b, c, "knows", nil)

	count := 0
	err := db.ForEachEdge(func(e *Edge) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("ForEachEdge failed: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 edges, got %d", count)
	}
}

func TestClear(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, Props{"name": "Alice"})
	b, _ := db.AddNode(nil, Props{"name": "Bob"})
	db.AddEdge(a, b, "knows", nil)

	err := db.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	stats := db.Stats()
	if stats.NodeCount != 0 || stats.EdgeCount != 0 {
		t.Fatalf("expected empty stats after Clear, got nodes=%d edges=%d", stats.NodeCount, stats.EdgeCount)
	}

	id, _ := db.AddNode(nil, nil)
	if id == 0 {
		t.Fatal("expected non-zero ID after Clear")
	}
}

func TestCreateUniqueConstraint(t *testing.T) {
	db := openTestDB(t)

	err := db.CreateUniqueConstraint("Person", "email")
	if err != nil {
		t.Fatalf("CreateUniqueConstraint failed: %v", err)
	}

	if !db.HasUniqueConstraint("Person", "email") {
		t.Fatal("expected constraint to exist")
	}

	err = db.CreateUniqueConstraint("Person", "email")
	if err == nil {
		t.Fatal("expected error for duplicate constraint")
	}
}

func TestFindByUniqueConstraint(t *testing.T) {
	db := openTestDB(t)

	db.CreateUniqueConstraint("Person", "email")

	aliceID, _ := db.AddNode([]string{"Person"}, Props{"email": "alice@example.com"})
	bobID, _ := db.AddNode([]string{"Person"}, Props{"email": "bob@example.com"})

	found, err := db.FindByUniqueConstraint("Person", "email", "alice@example.com")
	if err != nil {
		t.Fatalf("FindByUniqueConstraint failed: %v", err)
	}
	if found != aliceID {
		t.Fatalf("expected Alice's ID %d, got %d", aliceID, found)
	}

	found, err = db.FindByUniqueConstraint("Person", "email", "bob@example.com")
	if err != nil {
		t.Fatalf("FindByUniqueConstraint failed: %v", err)
	}
	if found != bobID {
		t.Fatalf("expected Bob's ID %d, got %d", bobID, found)
	}

	found, err = db.FindByUniqueConstraint("Person", "email", "nonexistent@example.com")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing value, got: %v", err)
	}
	if found != 0 {
		t.Fatalf("expected 0 for missing value, got %d", found)
	}
}

func TestUniqueConstraintViolation(t *testing.T) {
	db := openTestDB(t)

	db.CreateUniqueConstraint("Person", "email")

	_, err := db.AddNode([]string{"Person"}, Props{"email": "alice@example.com"})
	if err != nil {
		t.Fatalf("first AddNode should succeed: %v", err)
	}

	_, err = db.AddNode([]string{"Person"}, Props{"email": "alice@example.com"})
	if err == nil {
		t.Fatal("expected unique constraint violation")
	}

	_, err = db.AddNode([]string{"Person"}, Props{"email": "bob@example.com"})
	if err != nil {
		t.Fatalf("different email should succeed: %v", err)
	}
}

func TestUniqueConstraintCleanupOnDelete(t *testing.T) {
	db := openTestDB(t)

	db.CreateUniqueConstraint("Person", "email")

	aliceID, _ := db.AddNode([]string{"Person"}, Props{"email": "alice@example.com"})
	db.DeleteNode(aliceID)

	_, err := db.AddNode([]string{"Person"}, Props{"email": "alice@example.com"})
	if err != nil {
		t.Fatalf("should allow reuse after delete: %v", err)
	}
}

func TestScanEdgesByType(t *testing.T) {
	db := openTestDB(t)

	alice, _ := db.AddNode(nil, Props{"name": "Alice"})
	bob, _ := db.AddNode(nil, Props{"name": "Bob"})
	carol, _ := db.AddNode(nil, Props{"name": "Carol"})

	db.AddEdge(alice, bob, "knows", nil)
	db.AddEdge(bob, carol, "knows", nil)
	db.AddEdge(alice, carol, "likes", nil)

	var found []string
	err := db.ScanEdgesByType("knows", func(edge *Edge, src, dst *Node) error {
		found = append(found, src.GetString("name")+"->"+dst.GetString("name"))
		return nil
	})
	if err != nil {
		t.Fatalf("ScanEdgesByType failed: %v", err)
	}
	if len(found) != 2 {
		t.Fatalf("expected 2 knows edges, got %d", len(found))
	}
}

func TestGraphStoreInterface(t *testing.T) {
	db := openTestDB(t)

	var _ GraphStore = db

	var _ GraphReader = db
	var _ GraphWriter = db
}

func TestEdgeTypeIndexCleanup(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	eid, _ := db.AddEdge(a, b, "knows", nil)

	var before int
	db.ScanEdgesByType("knows", func(_ *Edge, _, _ *Node) error {
		before++
		return nil
	})
	if before != 1 {
		t.Fatalf("expected 1 edge before delete, got %d", before)
	}

	db.DeleteEdge(eid)

	var after int
	db.ScanEdgesByType("knows", func(_ *Edge, _, _ *Node) error {
		after++
		return nil
	})
	if after != 0 {
		t.Fatalf("expected 0 edges after delete, got %d", after)
	}
}

func TestAddNodeBatch(t *testing.T) {
	db := openTestDB(t)

	items := []struct {
		Labels []string
		Props  Props
	}{
		{Labels: []string{"Person"}, Props: Props{"name": "Alice", "age": 30}},
		{Labels: []string{"Person"}, Props: Props{"name": "Bob", "age": 25}},
		{Labels: []string{"City"}, Props: Props{"name": "London"}},
	}

	ids, err := db.AddNodeBatch(items)
	if err != nil {
		t.Fatalf("AddNodeBatch failed: %v", err)
	}
	if len(ids) != 3 {
		t.Fatalf("expected 3 IDs, got %d", len(ids))
	}

	for i, id := range ids {
		if id == 0 {
			t.Fatalf("expected non-zero ID at index %d", i)
		}
		if id == ids[0] && i > 0 {
			t.Fatalf("expected unique IDs, duplicate at index %d", i)
		}
	}

	alice := ids[0]
	node, err := db.GetNode(alice)
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}
	if len(node.Labels) != 1 || node.Labels[0] != "Person" {
		t.Fatalf("expected label Person, got %v", node.Labels)
	}
	if node.GetString("name") != "Alice" {
		t.Fatalf("expected name Alice, got %s", node.GetString("name"))
	}

	stats := db.Stats()
	if stats.NodeCount != 3 {
		t.Fatalf("expected 3 nodes, got %d", stats.NodeCount)
	}

	labelIDs, err := db.FindByLabel("Person")
	if err != nil {
		t.Fatalf("FindByLabel failed: %v", err)
	}
	if len(labelIDs) != 2 {
		t.Fatalf("expected 2 Person nodes, got %d", len(labelIDs))
	}
}

func TestAddNodeBatchEmpty(t *testing.T) {
	db := openTestDB(t)

	ids, err := db.AddNodeBatch(nil)
	if err != nil {
		t.Fatalf("AddNodeBatch(nil) failed: %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected 0 IDs for empty batch, got %d", len(ids))
	}

	stats := db.Stats()
	if stats.NodeCount != 0 {
		t.Fatalf("expected 0 nodes, got %d", stats.NodeCount)
	}
}

func TestAddNodeBatchDuplicateWithinBatch(t *testing.T) {
	db := openTestDB(t)

	db.CreateUniqueConstraint("Person", "email")

	items := []struct {
		Labels []string
		Props  Props
	}{
		{Labels: []string{"Person"}, Props: Props{"email": "alice@example.com"}},
		{Labels: []string{"Person"}, Props: Props{"email": "alice@example.com"}},
	}

	_, err := db.AddNodeBatch(items)
	if err == nil {
		t.Fatal("expected unique constraint violation within batch")
	}
}

func TestAddNodeBatchClosedDB(t *testing.T) {
	dir := tempDir(t)

	db, err := Open(dir, DefaultOptions())
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	db.Close()

	items := []struct {
		Labels []string
		Props  Props
	}{
		{Labels: []string{"Person"}, Props: Props{"name": "Alice"}},
	}

	_, err = db.AddNodeBatch(items)
	if err == nil {
		t.Fatal("expected error on closed DB")
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	db := openTestDB(t)

	const numGoroutines = 10
	const numOps = 50

	done := make(chan error, numGoroutines*2)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			for j := 0; j < numOps; j++ {
				_, err := db.AddNode([]string{"Worker"}, Props{"worker": idx, "op": j})
				if err != nil {
					done <- fmt.Errorf("AddNode failed: %v", err)
					return
				}
			}
			done <- nil
		}(i)

		go func(idx int) {
			for j := 0; j < numOps; j++ {
				_, err := db.FindByLabel("Worker")
				if err != nil {
					done <- fmt.Errorf("FindByLabel failed: %v", err)
					return
				}
				_ = db.Stats()
			}
			done <- nil
		}(i)
	}

	for i := 0; i < numGoroutines*2; i++ {
		if err := <-done; err != nil {
			t.Fatalf("concurrent operation failed: %v", err)
		}
	}

	stats := db.Stats()
	expectedNodes := int64(numGoroutines * numOps)
	if stats.NodeCount != expectedNodes {
		t.Fatalf("expected %d nodes after concurrent writes, got %d", expectedNodes, stats.NodeCount)
	}
}

func TestConcurrentEdgeAddAndTraversal(t *testing.T) {
	db := openTestDB(t)

	nodes := make([]NodeID, 5)
	for i := range nodes {
		id, err := db.AddNode(nil, Props{"idx": i})
		if err != nil {
			t.Fatalf("AddNode failed: %v", err)
		}
		nodes[i] = id
	}

	const numGoroutines = 10
	done := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			for j := 0; j < 20; j++ {
				from := nodes[idx%len(nodes)]
				to := nodes[(idx+j+1)%len(nodes)]
				_, err := db.AddEdge(from, to, "rel", Props{"g": idx})
				if err != nil {
					done <- fmt.Errorf("AddEdge failed: %v", err)
					return
				}
			}
			done <- nil
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		if err := <-done; err != nil {
			t.Fatalf("concurrent edge add failed: %v", err)
		}
	}

	for _, id := range nodes {
		edges, err := db.OutEdges(id)
		if err != nil {
			t.Fatalf("OutEdges failed for node %d: %v", id, err)
		}
		if len(edges) == 0 {
			t.Fatalf("expected at least 1 outgoing edge for node %d", id)
		}
	}
}

func TestGetLabelCountsReturnsError(t *testing.T) {
	db := openTestDB(t)

	db.CreateUniqueConstraint("Person", "email")
	db.AddNode([]string{"Person"}, Props{"email": "alice@example.com"})
	db.AddNode([]string{"Person"}, Props{"email": "bob@example.com"})

	counts, err := db.GetLabelCounts()
	if err != nil {
		t.Fatalf("GetLabelCounts should succeed: %v", err)
	}
	if counts["Person"] != 2 {
		t.Fatalf("expected 2 Person nodes, got %d", counts["Person"])
	}
}

func TestGetRelCountsReturnsError(t *testing.T) {
	db := openTestDB(t)

	a, _ := db.AddNode(nil, nil)
	b, _ := db.AddNode(nil, nil)
	db.AddEdge(a, b, "knows", nil)

	counts, err := db.GetRelCounts()
	if err != nil {
		t.Fatalf("GetRelCounts should succeed: %v", err)
	}
	if counts["knows"] != 1 {
		t.Fatalf("expected 1 knows edge, got %d", counts["knows"])
	}
}
