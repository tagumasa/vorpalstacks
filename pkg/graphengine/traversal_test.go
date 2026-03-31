package graphengine

import (
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
