package gremlinparser

import (
	"fmt"
	"testing"

	"vorpalstacks/pkg/graphengine"
)

func TestExec_V_All(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Bob"})

	result := execQuery(t, db, "g.V().count().toList()")
	val := toInt64(result)
	if val != 2 {
		t.Fatalf("expected 2, got %d", val)
	}
}

func TestExec_V_ByID(t *testing.T) {
	db := openTestDB(t)
	alice, _ := db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})

	result := execQuery(t, db, "g.V('"+fmt.Sprintf("%d", alice)+"').values('name').toList()")
	names := toStringSlice(result)
	if len(names) != 1 || names[0] != "Alice" {
		t.Fatalf("expected [Alice], got %v", names)
	}
}

func TestExec_Out(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Alice').out('knows').values('name').toList()")
	names := toStringSlice(result)
	if !containsAll(names, "Bob", "Charlie") {
		t.Fatalf("expected Bob and Charlie, got %v", names)
	}
}

func TestExec_In(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Bob').in('knows').values('name').toList()")
	names := toStringSlice(result)
	if len(names) != 1 || names[0] != "Alice" {
		t.Fatalf("expected [Alice], got %v", names)
	}
}

func TestExec_Both(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Bob').both().values('name').toList()")
	names := toStringSlice(result)
	if !containsAll(names, "Alice", "Dave") {
		t.Fatalf("expected Alice and Dave, got %v", names)
	}
}

func TestExec_OutE_InV(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Alice').outE('knows').inV().values('name').toList()")
	names := toStringSlice(result)
	if !containsAll(names, "Bob", "Charlie") {
		t.Fatalf("expected Bob and Charlie, got %v", names)
	}
}

func TestExec_HasString(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Alice').values('age').toList()")
	ages := toInt64Slice(result)
	if len(ages) != 1 || ages[0] != 30 {
		t.Fatalf("expected [30], got %v", ages)
	}
}

func TestExec_HasLabel(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})
	db.AddNode([]string{"animal"}, graphengine.Props{"name": "Rex"})

	result := execQuery(t, db, "g.V().hasLabel('person').values('name').toList()")
	names := toStringSlice(result)
	if len(names) != 1 || names[0] != "Alice" {
		t.Fatalf("expected [Alice], got %v", names)
	}
}

func TestExec_HasPredicate(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('age',gt(28)).values('name').toList()")
	names := toStringSlice(result)
	if !containsAll(names, "Alice", "Charlie") {
		t.Fatalf("expected Alice and Charlie (age>28), got %v", names)
	}
}

func TestExec_HasBetween(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('age',between(25,30)).values('name').toList()")
	names := toStringSlice(result)
	if !containsAll(names, "Alice", "Bob", "Dave") {
		t.Fatalf("expected Alice, Bob, Dave (25<=age<=30), got %v", names)
	}
}

func TestExec_HasWithin(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name',within('Alice','Eve')).values('name').toList()")
	names := toStringSlice(result)
	if !containsAll(names, "Alice", "Eve") {
		t.Fatalf("expected Alice and Eve, got %v", names)
	}
}

func TestExec_Values(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice", "age": 30})

	result := execQuery(t, db, "g.V().values('name','age').toList()")
	if len(toStringSlice(result)) != 2 {
		t.Fatalf("expected 2 values, got %v", toStringSlice(result))
	}
}

func TestExec_Id(t *testing.T) {
	db := openTestDB(t)
	alice, _ := db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})

	result := execQuery(t, db, "g.V().has('name','Alice').id().toList()")
	ids := toStringSlice(result)
	if len(ids) != 1 || ids[0] != fmt.Sprintf("%d", alice) {
		t.Fatalf("expected %d, got %v", alice, ids)
	}
}

func TestExec_Label(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})

	result := execQuery(t, db, "g.V().label().toList()")
	labels := toStringSlice(result)
	if len(labels) != 1 || labels[0] != "person" {
		t.Fatalf("expected [person], got %v", labels)
	}
}

func TestExec_Limit(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().limit(2).count().toList()")
	val := toInt64(result)
	if val != 2 {
		t.Fatalf("expected 2, got %d", val)
	}
}

func TestExec_Skip(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().skip(3).count().toList()")
	val := toInt64(result)
	if val != 2 {
		t.Fatalf("expected 2, got %d", val)
	}
}

func TestExec_Range(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().range(1,3).count().toList()")
	val := toInt64(result)
	if val != 2 {
		t.Fatalf("expected 2, got %d", val)
	}
}

func TestExec_Dedup(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})

	result := execQuery(t, db, "g.V().values('name').dedup().count().toList()")
	val := toInt64(result)
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
}

func TestExec_Count(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Bob"})

	result := execQuery(t, db, "g.V().count().toList()")
	val := toInt64(result)
	if val != 2 {
		t.Fatalf("expected 2, got %d", val)
	}
}

func TestExec_AddV(t *testing.T) {
	db := openTestDB(t)

	result := execQueryWrite(t, db, "g.addV('person').property('name','Alice').id().toList()")
	if len(toStringSlice(result)) != 1 {
		t.Fatalf("expected 1 ID, got %v", result)
	}

	count := execQueryWrite(t, db, "g.V().count().toList()")
	if toInt64(count) != 1 {
		t.Fatalf("expected 1 vertex, got %d", toInt64(count))
	}
}

func TestExec_Drop(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})

	execQueryWrite(t, db, "g.V().has('name','Alice').drop()")

	count := execQuery(t, db, "g.V().count().toList()")
	if toInt64(count) != 0 {
		t.Fatalf("expected 0 vertices after drop, got %d", toInt64(count))
	}
}

func TestExec_Property(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})

	execQueryWrite(t, db, "g.V().has('name','Alice').property('age',30)")

	result := execQuery(t, db, "g.V().has('name','Alice').values('age').toList()")
	ages := toInt64Slice(result)
	if len(ages) != 1 || ages[0] != 30 {
		t.Fatalf("expected [30], got %v", ages)
	}
}

func TestExec_AsSelect(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Alice').as('a').out('knows').as('b').select('a','b').by('name').toList()")
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestExec_GroupCount(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().hasLabel('person').groupCount().by('name').toList()")
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestExec_Union(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Alice').union(__.out('knows'),__.in('knows')).values('name').dedup().toList()")
	names := toStringSlice(result)
	if !containsAll(names, "Bob", "Charlie") {
		t.Fatalf("expected Bob and Charlie, got %v", names)
	}
}

func TestExec_Coalesce(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Eve').coalesce(__.in('knows'),__.in('follows')).values('name').toList()")
	names := toStringSlice(result)
	if len(names) != 1 || names[0] != "Charlie" {
		t.Fatalf("expected [Charlie], got %v", names)
	}
}

func TestExec_Path(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Alice').out('knows').path().toList()")
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestExec_Fold(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().hasLabel('person').values('age').fold().count().toList()")
	val := toInt64(result)
	if val != 1 {
		t.Fatalf("expected 1 (single folded list), got %d", val)
	}
}

func TestExec_Unfold(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().hasLabel('person').values('age').fold().unfold().count().toList()")
	val := toInt64(result)
	if val != 5 {
		t.Fatalf("expected 5 (unfolded), got %d", val)
	}
}

func TestExec_Constant(t *testing.T) {
	db := openTestDB(t)

	result := execQuery(t, db, "g.V().constant('hello').toList()")
	vals := toStringSlice(result)
	if len(vals) != 0 {
		t.Fatalf("expected 0 (no vertices), got %v", vals)
	}
}

func TestExec_Inject(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice"})

	result := execQuery(t, db, "g.inject('x','y').fold().toList()")
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestExec_MultiHop(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Alice').out().out().values('name').toList()")
	names := toStringSlice(result)
	if !containsAll(names, "Dave", "Eve") {
		t.Fatalf("expected Dave and Eve, got %v", names)
	}
}

func TestExec_HasNot(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice", "age": 30})
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Bob"})

	result := execQuery(t, db, "g.V().hasLabel('person').hasNot('age').values('name').toList()")
	names := toStringSlice(result)
	if len(names) != 1 || names[0] != "Bob" {
		t.Fatalf("expected [Bob], got %v", names)
	}
}

func TestExec_WhereNested(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().as('a').where(__.out('knows').has('age',gt(28))).values('name').toList()")
	names := toStringSlice(result)
	if len(names) != 1 || names[0] != "Alice" {
		t.Fatalf("expected [Alice] (only person who knows someone age>28), got %v", names)
	}
}

func TestExec_Filter(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().filter(__.out('knows').has('age',gt(28))).values('name').toList()")
	names := toStringSlice(result)
	if len(names) != 1 || names[0] != "Alice" {
		t.Fatalf("expected [Alice] (only person who knows someone age>28), got %v", names)
	}
}

func TestExec_ValueMap(t *testing.T) {
	db := openTestDB(t)
	db.AddNode([]string{"person"}, graphengine.Props{"name": "Alice", "age": 30})

	result := execQuery(t, db, "g.V().valueMap().toList()")
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestExec_MinMax(t *testing.T) {
	db := setupSocialGraph(t)

	minResult := execQuery(t, db, "g.V().values('age').min().toList()")
	maxResult := execQuery(t, db, "g.V().values('age').max().toList()")
	if minResult == nil || maxResult == nil {
		t.Fatal("expected results")
	}
	minVal := toInt64(minResult)
	maxVal := toInt64(maxResult)
	if minVal != 22 || maxVal != 35 {
		t.Fatalf("expected min=22 max=35, got min=%d max=%d", minVal, maxVal)
	}
}

func TestExec_Mean(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().values('age').mean().toList()")
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestExec_Sum(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().values('age').sum().toList()")
	val := toInt64(result)
	if val != 140 {
		t.Fatalf("expected 140, got %d", val)
	}
}

func TestExec_Tail(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().values('age').order().tail(2).toList()")
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestExec_OtherV(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Bob').inE('knows').otherV().values('name').toList()")
	names := toStringSlice(result)
	if len(names) != 1 || names[0] != "Alice" {
		t.Fatalf("expected [Alice], got %v", names)
	}
}

func TestExec_EdgeTraversal(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.E().hasLabel('knows').count().toList()")
	val := toInt64(result)
	if val != 4 {
		t.Fatalf("expected 4 knows edges, got %d", val)
	}
}

func TestExec_TextPredicate(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name',TextP.startingWith('Al')).values('name').toList()")
	names := toStringSlice(result)
	if len(names) != 1 || names[0] != "Alice" {
		t.Fatalf("expected [Alice], got %v", names)
	}
}

func TestExec_NestedTraversalArgs(t *testing.T) {
	db := setupSocialGraph(t)

	result := execQuery(t, db, "g.V().has('name','Alice').out('knows').has('age',gt(28)).values('name').toList()")
	names := toStringSlice(result)
	if len(names) != 1 || names[0] != "Charlie" {
		t.Fatalf("expected [Charlie], got %v", names)
	}
}
