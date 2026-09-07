package dynamodb

import (
	"math"
	"strconv"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// newVectorFixture creates a hash-only table (id S) whose single vector index
// "vec" indexes the attribute "embedding" with the given distance function
// and dimension count.
func newVectorFixture(t *testing.T, distanceFn string, dims int64) *DynamoDBStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(CreateTableParams{
		Name:                 "VecTbl",
		KeySchema:            []*KeySchemaElement{{AttributeName: "id", KeyType: KeyTypeHash}},
		AttributeDefinitions: []*AttributeDefinition{{AttributeName: "id", AttributeType: ScalarAttributeTypeS}},
		BillingMode:          BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := store.Tables().Update("VecTbl", func(table *Table) error {
		table.VectorIndexes = []*VectorIndex{{
			IndexName:           "vec",
			VectorAttributeName: "embedding",
			Dimensions:          dims,
			DistanceFunction:    distanceFn,
			Projection:          &Projection{ProjectionType: "ALL"},
			IndexStatus:         IndexStatusActive,
		}}
		return nil
	}); err != nil {
		t.Fatalf("add vector index: %v", err)
	}
	return store
}

// numListAttr builds an L-of-N attribute value — the wire shape of a vector.
func numListAttr(vals ...float64) *AttributeValue {
	l := make([]*AttributeValue, len(vals))
	for i, v := range vals {
		s := strconv.FormatFloat(v, 'g', -1, 64)
		l[i] = &AttributeValue{N: &s}
	}
	return &AttributeValue{L: l}
}

// putVectorFixtureItem stores one item whose embedding attribute is vec (nil
// omits the attribute) and maintains its index entries through the same
// transactional write path the service uses.
func putVectorFixtureItem(t *testing.T, store *DynamoDBStore, id string, vec []float64, extra map[string]*AttributeValue) {
	t.Helper()
	key := map[string]*AttributeValue{"id": strAttr(id)}
	attrs := map[string]*AttributeValue{}
	for k, v := range extra {
		attrs[k] = v
	}
	if vec != nil {
		attrs["embedding"] = numListAttr(vec...)
	}
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		if err := txn.PutItem("VecTbl", key, attrs); err != nil {
			return err
		}
		return txn.PutIndexEntries("VecTbl", &Item{TableName: "VecTbl", Key: key, Attributes: attrs})
	}); err != nil {
		t.Fatalf("put fixture item %s: %v", id, err)
	}
}

func vectorTopKIDs(t *testing.T, store *DynamoDBStore, query []float64, k int, filter func(*Item) bool) ([]string, []float64) {
	t.Helper()
	var ids []string
	var scores []float64
	if err := store.View(t.Context(), func(txn *DynamoDBTxn) error {
		hits, err := txn.VectorTopK("VecTbl", "vec", query, k, filter)
		if err != nil {
			return err
		}
		for _, h := range hits {
			ids = append(ids, *h.Item.Key["id"].S)
			scores = append(scores, h.Score)
		}
		return nil
	}); err != nil {
		t.Fatalf("vector top-k: %v", err)
	}
	return ids, scores
}

func scoreNear(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s: got %g, want %g", name, got, want)
	}
}

// Vector elements are 32-bit IEEE-754 values (SearchVectors API reference):
// extraction rounds every element to its float32 value, so decimal spellings
// that differ only below float32 precision are the same element and cannot
// be distinguished by score.
func TestVectorElementsQuantisedToFloat32(t *testing.T) {
	vi := &VectorIndex{VectorAttributeName: "embedding", Dimensions: 1}
	quantised := float64(float32(0.1))
	if quantised == 0.1 {
		t.Fatal("precondition: 0.1 must not be float32-exact for this test")
	}
	for _, spelling := range []string{"0.1", "0.100000001"} {
		n := spelling
		item := &Item{Attributes: map[string]*AttributeValue{
			"embedding": {L: []*AttributeValue{{N: &n}}},
		}}
		vec := extractVector(vi, item)
		if len(vec) != 1 || vec[0] != quantised {
			t.Fatalf("spelling %q extracted %v, want the float32 value %g", spelling, vec, quantised)
		}
	}

	// Sub-float32 element differences score identically end to end.
	store := newVectorFixture(t, "EUCLIDEAN", 1)
	putVectorFixtureItem(t, store, "a", []float64{0.1}, nil)
	putVectorFixtureItem(t, store, "b", []float64{0.100000001}, nil)
	ids, scores := vectorTopKIDs(t, store, []float64{0.5}, 2, nil)
	if len(ids) != 2 {
		t.Fatalf("top-k = %v, want both items", ids)
	}
	if scores[0] != scores[1] {
		t.Fatalf("sub-float32 difference distinguished: %g vs %g", scores[0], scores[1])
	}
	scoreNear(t, "quantised euclidean score", scores[0], math.Abs(float64(float32(0.5))-quantised))
}

func TestVectorTopKCOSINEOrdering(t *testing.T) {
	store := newVectorFixture(t, "COSINE", 2)
	putVectorFixtureItem(t, store, "identical", []float64{1, 0}, nil)
	putVectorFixtureItem(t, store, "orthogonal", []float64{0, 1}, nil)
	putVectorFixtureItem(t, store, "opposite", []float64{-1, 0}, nil)

	ids, scores := vectorTopKIDs(t, store, []float64{1, 0}, 3, nil)
	if len(ids) != 3 || ids[0] != "identical" || ids[1] != "orthogonal" || ids[2] != "opposite" {
		t.Fatalf("cosine order = %v, want identical,orthogonal,opposite", ids)
	}
	scoreNear(t, "identical score", scores[0], 0)
	scoreNear(t, "orthogonal score", scores[1], 1)
	scoreNear(t, "opposite score", scores[2], 2)

	// TopK bounds the result set to the k most similar items.
	ids, _ = vectorTopKIDs(t, store, []float64{1, 0}, 2, nil)
	if len(ids) != 2 || ids[0] != "identical" || ids[1] != "orthogonal" {
		t.Fatalf("cosine k=2 = %v", ids)
	}
}

func TestVectorTopKEUCLIDEANOrdering(t *testing.T) {
	store := newVectorFixture(t, "EUCLIDEAN", 2)
	putVectorFixtureItem(t, store, "origin", []float64{0, 0}, nil)
	putVectorFixtureItem(t, store, "unit", []float64{1, 0}, nil)
	putVectorFixtureItem(t, store, "threefour", []float64{3, 4}, nil)

	ids, scores := vectorTopKIDs(t, store, []float64{0, 0}, 3, nil)
	if len(ids) != 3 || ids[0] != "origin" || ids[1] != "unit" || ids[2] != "threefour" {
		t.Fatalf("euclidean order = %v, want origin,unit,threefour", ids)
	}
	scoreNear(t, "origin score", scores[0], 0)
	scoreNear(t, "unit score", scores[1], 1)
	// The score is the Euclidean distance (the square root), not its square.
	scoreNear(t, "threefour score", scores[2], 5)
}

func TestVectorTopKDotProductDescending(t *testing.T) {
	store := newVectorFixture(t, "DOT_PRODUCT", 2)
	putVectorFixtureItem(t, store, "a", []float64{1, 0}, nil)
	putVectorFixtureItem(t, store, "b", []float64{0, 1}, nil)
	putVectorFixtureItem(t, store, "c", []float64{-1, 0}, nil)

	ids, scores := vectorTopKIDs(t, store, []float64{1, 2}, 3, nil)
	if len(ids) != 3 || ids[0] != "b" || ids[1] != "a" || ids[2] != "c" {
		t.Fatalf("dot-product order = %v, want b,a,c (highest similarity first)", ids)
	}
	scoreNear(t, "b score", scores[0], 2)
	scoreNear(t, "a score", scores[1], 1)
	scoreNear(t, "c score", scores[2], -1)
}

func TestVectorEntriesMaintainedOnWritePath(t *testing.T) {
	store := newVectorFixture(t, "COSINE", 2)

	// A valid vector is indexed; a dimension-mismatched or missing attribute
	// clears the entry on overwrite; deleting the item removes it.
	putVectorFixtureItem(t, store, "k1", []float64{1, 0}, nil)
	if ids, _ := vectorTopKIDs(t, store, []float64{1, 0}, 10, nil); len(ids) != 1 {
		t.Fatalf("after valid put: top-k = %v, want [k1]", ids)
	}

	putVectorFixtureItem(t, store, "k1", []float64{1, 0, 0}, nil) // wrong dimensions
	if ids, _ := vectorTopKIDs(t, store, []float64{1, 0}, 10, nil); len(ids) != 0 {
		t.Fatalf("after wrong-dimension overwrite: top-k = %v, want none", ids)
	}

	putVectorFixtureItem(t, store, "k1", []float64{0, 1}, nil)
	if ids, _ := vectorTopKIDs(t, store, []float64{0, 1}, 10, nil); len(ids) != 1 || ids[0] != "k1" {
		t.Fatalf("after re-put with valid vector: top-k = %v, want [k1]", ids)
	}

	putVectorFixtureItem(t, store, "k1", nil, nil) // attribute removed
	if ids, _ := vectorTopKIDs(t, store, []float64{0, 1}, 10, nil); len(ids) != 0 {
		t.Fatalf("after attribute removal: top-k = %v, want none", ids)
	}

	putVectorFixtureItem(t, store, "k1", []float64{1, 1}, nil)
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		if err := txn.DeleteItem("VecTbl", map[string]*AttributeValue{"id": strAttr("k1")}); err != nil {
			return err
		}
		return txn.DeleteIndexEntries("VecTbl", &Item{
			TableName: "VecTbl",
			Key:       map[string]*AttributeValue{"id": strAttr("k1")},
		})
	}); err != nil {
		t.Fatalf("delete item: %v", err)
	}
	if ids, _ := vectorTopKIDs(t, store, []float64{1, 1}, 10, nil); len(ids) != 0 {
		t.Fatalf("after item delete: top-k = %v, want none", ids)
	}
}

func TestVectorTopKFilterPrunes(t *testing.T) {
	store := newVectorFixture(t, "COSINE", 2)
	putVectorFixtureItem(t, store, "keep", []float64{1, 0}, map[string]*AttributeValue{"category": strAttr("a")})
	putVectorFixtureItem(t, store, "drop", []float64{0, 1}, map[string]*AttributeValue{"category": strAttr("b")})

	filter := func(item *Item) bool {
		attr := item.GetAttribute("category")
		return attr != nil && attr.S != nil && *attr.S == "a"
	}
	ids, _ := vectorTopKIDs(t, store, []float64{1, 0}, 2, filter)
	if len(ids) != 1 || ids[0] != "keep" {
		t.Fatalf("filtered top-k = %v, want [keep]", ids)
	}
}

func TestDeleteVectorEntriesForIndexSweeps(t *testing.T) {
	store := newVectorFixture(t, "COSINE", 2)
	putVectorFixtureItem(t, store, "k1", []float64{1, 0}, nil)
	putVectorFixtureItem(t, store, "k2", []float64{0, 1}, nil)
	putVectorFixtureItem(t, store, "k3", []float64{1, 1}, nil)

	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.DeleteVectorEntriesForIndex("VecTbl", "vec")
	}); err != nil {
		t.Fatalf("sweep: %v", err)
	}
	if ids, _ := vectorTopKIDs(t, store, []float64{1, 0}, 10, nil); len(ids) != 0 {
		t.Fatalf("after sweep: top-k = %v, want none", ids)
	}
}

func TestVectorBackfillForIndex(t *testing.T) {
	store := newVectorFixture(t, "COSINE", 2)
	// Items exist before the index is attached — the backfill twin must
	// populate their entries afterwards.
	putVectorFixtureItem(t, store, "k1", []float64{1, 0}, nil)
	putVectorFixtureItem(t, store, "k2", []float64{0, 1}, nil)
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.DeleteVectorEntriesForIndex("VecTbl", "vec")
	}); err != nil {
		t.Fatalf("clear entries: %v", err)
	}

	if err := store.Items().Scan("VecTbl", func(item *Item) error {
		return store.Update(t.Context(), func(txn *DynamoDBTxn) error {
			return txn.PutVectorEntriesForIndex("VecTbl", "vec", item)
		})
	}); err != nil {
		t.Fatalf("backfill: %v", err)
	}
	ids, _ := vectorTopKIDs(t, store, []float64{1, 0}, 10, nil)
	if len(ids) != 2 || ids[0] != "k1" || ids[1] != "k2" {
		t.Fatalf("after backfill: top-k = %v, want [k1 k2]", ids)
	}
}

func TestVectorTopKRejectsUnknownIndex(t *testing.T) {
	store := newVectorFixture(t, "COSINE", 2)
	if err := store.View(t.Context(), func(txn *DynamoDBTxn) error {
		_, err := txn.VectorTopK("VecTbl", "absent", []float64{1, 0}, 1, nil)
		return err
	}); err == nil {
		t.Fatalf("unknown index: expected error")
	}
}
