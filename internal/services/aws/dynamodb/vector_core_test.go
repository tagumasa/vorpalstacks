package dynamodb

import (
	"strconv"
	"testing"

	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// newVectorCoreFixture opens a real store with a hash-only table whose
// vector index "vec" (COSINE over the 2-dimensional "embedding" attribute)
// partitions on the HASH attribute "category" and inline-filters on the
// numeric attribute "year".
func newVectorCoreFixture(t *testing.T) *dbstore.DynamoDBStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	store := dbstore.NewDynamoDBStore(st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "VecTbl",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := store.Tables().Update("VecTbl", func(table *dbstore.Table) error {
		table.VectorIndexes = []*dbstore.VectorIndex{{
			IndexName:           "vec",
			VectorAttributeName: "embedding",
			Dimensions:          2,
			DistanceFunction:    "COSINE",
			Projection:          &dbstore.Projection{ProjectionType: "ALL"},
			SearchSchema: []*dbstore.SearchSchemaElement{
				{AttributeName: "category", SearchSchemaElementType: "HASH"},
				{AttributeName: "year", SearchSchemaElementType: "INLINE_FILTER"},
			},
			IndexStatus: dbstore.IndexStatusActive,
		}}
		return nil
	}); err != nil {
		t.Fatalf("attach vector index: %v", err)
	}
	return store
}

func putVecCoreItem(t *testing.T, store *dbstore.DynamoDBStore, id string, vec []float64, category string, year int64, extra map[string]*dbstore.AttributeValue) {
	t.Helper()
	key := map[string]*dbstore.AttributeValue{"id": {S: strPtr(id)}}
	attrs := map[string]*dbstore.AttributeValue{
		"category": {S: strPtr(category)},
		"year":     {N: strPtr(strconv.FormatInt(year, 10))},
	}
	for k, v := range extra {
		attrs[k] = v
	}
	if vec != nil {
		l := make([]*dbstore.AttributeValue, len(vec))
		for i, v := range vec {
			l[i] = &dbstore.AttributeValue{N: strPtr(strconv.FormatFloat(v, 'g', -1, 64))}
		}
		attrs["embedding"] = &dbstore.AttributeValue{L: l}
	}
	if err := store.Update(t.Context(), func(txn *dbstore.DynamoDBTxn) error {
		if err := txn.PutItem("VecTbl", key, attrs); err != nil {
			return err
		}
		return txn.PutIndexEntries("VecTbl", &dbstore.Item{TableName: "VecTbl", Key: key, Attributes: attrs})
	}); err != nil {
		t.Fatalf("put item %s: %v", id, err)
	}
}

func strPtr(s string) *string { return &s }

func vecSearchParams(table, index string, vec []float64, topK float64) map[string]interface{} {
	sv := make([]interface{}, len(vec))
	for i, v := range vec {
		sv[i] = map[string]interface{}{"N": strconv.FormatFloat(v, 'g', -1, 64)}
	}
	return map[string]interface{}{
		"TableName":    table,
		"IndexName":    index,
		"SearchVector": sv,
		"TopK":         topK,
	}
}

// addCategoryCondition equips a search with the partition-key equality
// clause the fixture index requires: its HASH element "category" must appear
// in the SearchConditionExpression.
func addCategoryCondition(params map[string]interface{}, category string) {
	params["SearchConditionExpression"] = "category = :cat"
	params["ExpressionAttributeValues"] = map[string]interface{}{":cat": map[string]interface{}{"S": category}}
}

func TestSearchVectorsCoreOrdersAndScores(t *testing.T) {
	store := newVectorCoreFixture(t)
	putVecCoreItem(t, store, "same", []float64{1, 0}, "a", 2020, nil)
	putVecCoreItem(t, store, "orth", []float64{0, 1}, "a", 2021, nil)
	putVecCoreItem(t, store, "opp", []float64{-1, 0}, "a", 2022, nil)

	svc := &DynamoDBService{}
	params := vecSearchParams("VecTbl", "vec", []float64{1, 0}, 3)
	addCategoryCondition(params, "a")
	result, err := svc.searchVectorsCore(t.Context(), store, params)
	if err != nil {
		t.Fatalf("searchVectorsCore: %v", err)
	}
	if len(result.Items) != 3 {
		t.Fatalf("results = %d, want 3", len(result.Items))
	}
	if got := result.Items[0]["id"].(map[string]interface{})["S"]; got != "same" {
		t.Errorf("first result = %v, want same", got)
	}
	if result.Scores[0] != 0 || result.Scores[1] != 1 || result.Scores[2] != 2 {
		t.Errorf("scores = %v, want [0 1 2]", result.Scores)
	}
	if result.VectorSearchRequestBytes <= 0 {
		t.Errorf("VectorSearchRequestBytes = %d, want > 0", result.VectorSearchRequestBytes)
	}
}

func TestSearchVectorsCoreValidation(t *testing.T) {
	store := newVectorCoreFixture(t)
	putVecCoreItem(t, store, "k1", []float64{1, 0}, "a", 2020, map[string]*dbstore.AttributeValue{"title": {S: strPtr("A Title")}})
	svc := &DynamoDBService{}

	cases := map[string]map[string]interface{}{
		"empty table name":     vecSearchParams("", "vec", []float64{1, 0}, 1),
		"empty index name":     vecSearchParams("VecTbl", "", []float64{1, 0}, 1),
		"missing searchVector": {"TableName": "VecTbl", "IndexName": "vec", "TopK": float64(1)},
		"topK zero":            vecSearchParams("VecTbl", "vec", []float64{1, 0}, 0),
		"topK over max":        vecSearchParams("VecTbl", "vec", []float64{1, 0}, 101),
		"dimension mismatch":   vecSearchParams("VecTbl", "vec", []float64{1, 0, 0}, 1),
		"non-number element": {"TableName": "VecTbl", "IndexName": "vec", "TopK": float64(1),
			"SearchVector": []interface{}{map[string]interface{}{"S": "x"}}},
	}
	for name, params := range cases {
		if _, err := svc.searchVectorsCore(t.Context(), store, params); err == nil {
			t.Errorf("%s: expected rejection", name)
		}
	}

	// Missing table and missing index are ResourceNotFoundException, not
	// validation errors.
	if _, err := svc.searchVectorsCore(t.Context(), store, vecSearchParams("Absent", "vec", []float64{1, 0}, 1)); err != ErrResourceNotFound {
		t.Errorf("absent table: got %v, want ErrResourceNotFound", err)
	}
	if _, err := svc.searchVectorsCore(t.Context(), store, vecSearchParams("VecTbl", "absent", []float64{1, 0}, 1)); err != ErrResourceNotFound {
		t.Errorf("absent index: got %v, want ErrResourceNotFound", err)
	}
}

func TestSearchVectorsCoreCondition(t *testing.T) {
	store := newVectorCoreFixture(t)
	putVecCoreItem(t, store, "near-a", []float64{1, 0}, "a", 2020, nil)
	putVecCoreItem(t, store, "far-a", []float64{0, 1}, "a", 2024, nil)
	putVecCoreItem(t, store, "near-b", []float64{0.9, 0.1}, "b", 2020, nil)
	svc := &DynamoDBService{}

	hashEq := vecSearchParams("VecTbl", "vec", []float64{1, 0}, 3)
	addCategoryCondition(hashEq, "a")
	result, err := svc.searchVectorsCore(t.Context(), store, hashEq)
	if err != nil {
		t.Fatalf("hash equality condition: %v", err)
	}
	if len(result.Items) != 2 || result.Items[0]["id"].(map[string]interface{})["S"] != "near-a" {
		t.Fatalf("hash-equality results = %v", result.Items)
	}

	// An INLINE_FILTER attribute joins the condition with the equality
	// operator — comparison, range, and set-membership operators are not
	// part of the contract for either attribute role.
	inlineEq := vecSearchParams("VecTbl", "vec", []float64{1, 0}, 3)
	inlineEq["SearchConditionExpression"] = "category = :cat AND year = :y"
	inlineEq["ExpressionAttributeValues"] = map[string]interface{}{
		":cat": map[string]interface{}{"S": "a"},
		":y":   map[string]interface{}{"N": "2024"},
	}
	result, err = svc.searchVectorsCore(t.Context(), store, inlineEq)
	if err != nil {
		t.Fatalf("inline-filter equality condition: %v", err)
	}
	if len(result.Items) != 1 || result.Items[0]["id"].(map[string]interface{})["S"] != "far-a" {
		t.Fatalf("inline-filter results = %v", result.Items)
	}

	rangeCond := vecSearchParams("VecTbl", "vec", []float64{1, 0}, 3)
	rangeCond["SearchConditionExpression"] = "category = :cat AND year BETWEEN :lo AND :hi"
	rangeCond["ExpressionAttributeValues"] = map[string]interface{}{
		":cat": map[string]interface{}{"S": "a"},
		":lo":  map[string]interface{}{"N": "2021"},
		":hi":  map[string]interface{}{"N": "2030"},
	}
	if _, err := svc.searchVectorsCore(t.Context(), store, rangeCond); err == nil {
		t.Errorf("range operator on INLINE_FILTER attribute: expected rejection")
	}

	ltCond := vecSearchParams("VecTbl", "vec", []float64{1, 0}, 3)
	ltCond["SearchConditionExpression"] = "category = :cat AND year < :y"
	ltCond["ExpressionAttributeValues"] = map[string]interface{}{
		":cat": map[string]interface{}{"S": "a"},
		":y":   map[string]interface{}{"N": "2030"},
	}
	if _, err := svc.searchVectorsCore(t.Context(), store, ltCond); err == nil {
		t.Errorf("comparison operator on INLINE_FILTER attribute: expected rejection")
	}

	// The HASH element's value must be provided: without the partition-key
	// clause — or without any expression at all — the search is rejected.
	noHash := vecSearchParams("VecTbl", "vec", []float64{1, 0}, 3)
	noHash["SearchConditionExpression"] = "year = :y"
	noHash["ExpressionAttributeValues"] = map[string]interface{}{":y": map[string]interface{}{"N": "2020"}}
	if _, err := svc.searchVectorsCore(t.Context(), store, noHash); err == nil {
		t.Errorf("condition without the HASH element's value: expected rejection")
	}
	if _, err := svc.searchVectorsCore(t.Context(), store, vecSearchParams("VecTbl", "vec", []float64{1, 0}, 3)); err == nil {
		t.Errorf("empty condition on a partitioned index: expected rejection")
	}

	// A comparison on the HASH attribute is rejected, as is an attribute
	// outside the search schema.
	hashLt := vecSearchParams("VecTbl", "vec", []float64{1, 0}, 3)
	hashLt["SearchConditionExpression"] = "category < :cat"
	hashLt["ExpressionAttributeValues"] = map[string]interface{}{":cat": map[string]interface{}{"S": "z"}}
	if _, err := svc.searchVectorsCore(t.Context(), store, hashLt); err == nil {
		t.Errorf("comparison on HASH attribute: expected rejection")
	}

	outside := vecSearchParams("VecTbl", "vec", []float64{1, 0}, 3)
	outside["SearchConditionExpression"] = "title = :t"
	outside["ExpressionAttributeValues"] = map[string]interface{}{":t": map[string]interface{}{"S": "x"}}
	if _, err := svc.searchVectorsCore(t.Context(), store, outside); err == nil {
		t.Errorf("attribute outside the search schema: expected rejection")
	}
}

// TestSearchVectorsCoreUnpartitionedIndex pins the counterpart of the HASH
// requirement: an index whose search schema declares no partition key
// searches the whole index and needs no condition expression.
func TestSearchVectorsCoreUnpartitionedIndex(t *testing.T) {
	store := newVectorCoreFixture(t)
	if _, err := store.Tables().Update("VecTbl", func(table *dbstore.Table) error {
		table.VectorIndexes[0].SearchSchema = []*dbstore.SearchSchemaElement{
			{AttributeName: "year", SearchSchemaElementType: "INLINE_FILTER"},
		}
		return nil
	}); err != nil {
		t.Fatalf("drop partition key: %v", err)
	}
	putVecCoreItem(t, store, "k1", []float64{1, 0}, "a", 2020, nil)
	putVecCoreItem(t, store, "k2", []float64{0, 1}, "b", 2021, nil)

	svc := &DynamoDBService{}
	result, err := svc.searchVectorsCore(t.Context(), store, vecSearchParams("VecTbl", "vec", []float64{1, 0}, 2))
	if err != nil {
		t.Fatalf("unconditional search on unpartitioned index: %v", err)
	}
	if len(result.Items) != 2 || result.Items[0]["id"].(map[string]interface{})["S"] != "k1" {
		t.Fatalf("unpartitioned results = %v", result.Items)
	}
}

// TestSearchVectorsCoreTableArn pins the TableArn form of the TableName
// member: a DynamoDB table ARN addresses the same table as its name, while
// a foreign-service ARN or a non-table resource is a validation error.
func TestSearchVectorsCoreTableArn(t *testing.T) {
	store := newVectorCoreFixture(t)
	putVecCoreItem(t, store, "k1", []float64{1, 0}, "a", 2020, nil)
	svc := &DynamoDBService{}

	params := vecSearchParams("arn:aws:dynamodb:us-east-1:123456789012:table/VecTbl", "vec", []float64{1, 0}, 1)
	addCategoryCondition(params, "a")
	result, err := svc.searchVectorsCore(t.Context(), store, params)
	if err != nil {
		t.Fatalf("table ARN form: %v", err)
	}
	if len(result.Items) != 1 || result.Items[0]["id"].(map[string]interface{})["S"] != "k1" {
		t.Fatalf("table ARN results = %v", result.Items)
	}

	for name, table := range map[string]string{
		"foreign service":       "arn:aws:s3:::bucket",
		"non-table resource":    "arn:aws:dynamodb:us-east-1:123456789012:backup/VecTbl",
		"stream resource":       "arn:aws:dynamodb:us-east-1:123456789012:table/VecTbl/stream/2026-09-07T00:00:00.000",
		"malformed arn":         "arn:aws:dynamodb",
		"plain name still best": "not a table",
	} {
		bad := vecSearchParams(table, "vec", []float64{1, 0}, 1)
		addCategoryCondition(bad, "a")
		if _, err := svc.searchVectorsCore(t.Context(), store, bad); err == nil {
			t.Errorf("%s: expected rejection", name)
		}
	}
}

func TestSearchVectorsCoreProjection(t *testing.T) {
	// Rebuild the fixture with an INCLUDE projection so the retrievable set
	// is bounded: keys, the vector attribute, "title", and the schema
	// elements.
	store := newVectorCoreFixture(t)
	if _, err := store.Tables().Update("VecTbl", func(table *dbstore.Table) error {
		table.VectorIndexes[0].Projection = &dbstore.Projection{
			ProjectionType:   "INCLUDE",
			NonKeyAttributes: []string{"title"},
		}
		return nil
	}); err != nil {
		t.Fatalf("narrow projection: %v", err)
	}
	putVecCoreItem(t, store, "k1", []float64{1, 0}, "a", 2020, map[string]*dbstore.AttributeValue{"title": {S: strPtr("A Title")}})
	svc := &DynamoDBService{}

	ok := vecSearchParams("VecTbl", "vec", []float64{1, 0}, 1)
	addCategoryCondition(ok, "a")
	ok["ProjectionExpression"] = "title, embedding"
	result, err := svc.searchVectorsCore(t.Context(), store, ok)
	if err != nil {
		t.Fatalf("projected search: %v", err)
	}
	item := result.Items[0]
	for _, attr := range []string{"id", "title", "embedding"} {
		if _, present := item[attr]; !present {
			t.Errorf("projected result missing %q: %v", attr, item)
		}
	}
	if _, present := item["category"]; present {
		t.Errorf("projected result must not carry attributes outside the projection: %v", item)
	}

	outside := vecSearchParams("VecTbl", "vec", []float64{1, 0}, 1)
	outside["ProjectionExpression"] = "body"
	if _, err := svc.searchVectorsCore(t.Context(), store, outside); err == nil {
		t.Errorf("requesting an unprojected attribute: expected rejection")
	}
}
