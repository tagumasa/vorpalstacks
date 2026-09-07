package dynamodb

import (
	"errors"
	"fmt"
	"testing"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// vectorIndexParams builds a CreateTable VectorIndexes element in wire form
// (numbers arrive from JSON as float64).
func vectorIndexParams(indexName, attrName string, dims float64, distanceFn string) map[string]interface{} {
	return map[string]interface{}{
		"IndexName":        indexName,
		"VectorAttribute":  map[string]interface{}{"AttributeName": attrName},
		"Dimensions":       dims,
		"DistanceFunction": distanceFn,
		"Projection":       map[string]interface{}{"ProjectionType": "ALL"},
	}
}

func TestParseVectorIndexesRoundTrip(t *testing.T) {
	params := map[string]interface{}{
		"VectorIndexes": []interface{}{
			func() map[string]interface{} {
				vi := vectorIndexParams("vec-idx", "embedding", float64(3), "COSINE")
				vi["SearchSchema"] = []interface{}{
					map[string]interface{}{"AttributeName": "category", "SearchSchemaElementType": "HASH"},
					map[string]interface{}{"AttributeName": "year", "SearchSchemaElementType": "INLINE_FILTER"},
				}
				return vi
			}(),
		},
	}
	idx, err := parseVectorIndexes(params)
	if err != nil {
		t.Fatalf("parseVectorIndexes: %v", err)
	}
	if len(idx) != 1 {
		t.Fatalf("expected 1 index, got %d", len(idx))
	}
	got := idx[0]
	if got.IndexName != "vec-idx" || got.VectorAttributeName != "embedding" {
		t.Errorf("names mismatch: %+v", got)
	}
	if got.Dimensions != 3 {
		t.Errorf("Dimensions: got %d, want 3", got.Dimensions)
	}
	if got.DistanceFunction != "COSINE" {
		t.Errorf("DistanceFunction: got %s", got.DistanceFunction)
	}
	if got.IndexStatus != dbstore.IndexStatusActive {
		t.Errorf("IndexStatus: got %s, want ACTIVE", got.IndexStatus)
	}
	if got.Projection == nil || got.Projection.ProjectionType != "ALL" {
		t.Errorf("Projection: %+v", got.Projection)
	}
	if len(got.SearchSchema) != 2 {
		t.Fatalf("SearchSchema: got %d elements", len(got.SearchSchema))
	}
	if got.SearchSchema[0].AttributeName != "category" || got.SearchSchema[0].SearchSchemaElementType != "HASH" {
		t.Errorf("schema[0]: %+v", got.SearchSchema[0])
	}
	if got.SearchSchema[1].SearchSchemaElementType != "INLINE_FILTER" {
		t.Errorf("schema[1]: %+v", got.SearchSchema[1])
	}

	// Absent member yields nil (no vector indexes).
	none, err := parseVectorIndexes(map[string]interface{}{})
	if err != nil || none != nil {
		t.Errorf("absent VectorIndexes: got (%v, %v)", none, err)
	}
}

func TestParseVectorIndexRejections(t *testing.T) {
	base := func() map[string]interface{} {
		return vectorIndexParams("vec-idx", "embedding", float64(3), "COSINE")
	}
	cases := map[string]func(map[string]interface{}){
		"short index name":    func(m map[string]interface{}) { m["IndexName"] = "vi" },
		"missing index name":  func(m map[string]interface{}) { m["IndexName"] = "" },
		"missing vector attr": func(m map[string]interface{}) { delete(m, "VectorAttribute") },
		"empty vector attr":   func(m map[string]interface{}) { m["VectorAttribute"] = map[string]interface{}{"AttributeName": ""} },
		"long vector attr": func(m map[string]interface{}) {
			m["VectorAttribute"] = map[string]interface{}{"AttributeName": attrNameOfLength(256)}
		},
		"missing dimensions":   func(m map[string]interface{}) { delete(m, "Dimensions") },
		"zero dimensions":      func(m map[string]interface{}) { m["Dimensions"] = float64(0) },
		"oversized dimensions": func(m map[string]interface{}) { m["Dimensions"] = float64(4097) },
		"bad distance fn":      func(m map[string]interface{}) { m["DistanceFunction"] = "MANHATTAN" },
		"missing distance fn":  func(m map[string]interface{}) { m["DistanceFunction"] = "" },
		"missing projection":   func(m map[string]interface{}) { delete(m, "Projection") },
		"empty search schema":  func(m map[string]interface{}) { m["SearchSchema"] = []interface{}{} },
		"bad schema elem type": func(m map[string]interface{}) {
			m["SearchSchema"] = []interface{}{map[string]interface{}{"AttributeName": "a", "SearchSchemaElementType": "RANGE"}}
		},
		"schema elem missing attr": func(m map[string]interface{}) {
			m["SearchSchema"] = []interface{}{map[string]interface{}{"SearchSchemaElementType": "HASH"}}
		},
		"duplicate schema attr": func(m map[string]interface{}) {
			m["SearchSchema"] = []interface{}{
				map[string]interface{}{"AttributeName": "a", "SearchSchemaElementType": "HASH"},
				map[string]interface{}{"AttributeName": "a", "SearchSchemaElementType": "INLINE_FILTER"},
			}
		},
		"two hash elements": func(m map[string]interface{}) {
			m["SearchSchema"] = []interface{}{
				map[string]interface{}{"AttributeName": "a", "SearchSchemaElementType": "HASH"},
				map[string]interface{}{"AttributeName": "b", "SearchSchemaElementType": "HASH"},
			}
		},
		"too many inline filters": func(m map[string]interface{}) {
			elements := make([]interface{}, 0, dbstore.VectorInlineFiltersMax+2)
			elements = append(elements, map[string]interface{}{"AttributeName": "a", "SearchSchemaElementType": "HASH"})
			for i := 0; i < dbstore.VectorInlineFiltersMax+1; i++ {
				elements = append(elements, map[string]interface{}{
					"AttributeName": fmt.Sprintf("f%02d", i), "SearchSchemaElementType": "INLINE_FILTER",
				})
			}
			m["SearchSchema"] = elements
		},
	}
	for name, mutate := range cases {
		m := base()
		mutate(m)
		if _, err := parseVectorIndex(m); err == nil {
			t.Errorf("%s: expected rejection", name)
		}
	}
}

func attrNameOfLength(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}

// TestCreateTableCoreRejectsExcessVectorIndexes pins the per-table quota:
// a CreateTable request may declare at most the documented number of vector
// indexes.
func TestCreateTableCoreRejectsExcessVectorIndexes(t *testing.T) {
	svc := &DynamoDBService{}
	indexes := make([]*dbstore.VectorIndex, 0, dbstore.VectorIndexesPerTable+1)
	for i := 0; i <= dbstore.VectorIndexesPerTable; i++ {
		indexes = append(indexes, &dbstore.VectorIndex{
			IndexName:           fmt.Sprintf("vi-%d", i),
			VectorAttributeName: "embedding",
			Dimensions:          2,
			DistanceFunction:    "COSINE",
			Projection:          &dbstore.Projection{ProjectionType: "ALL"},
		})
	}
	in := CreateTableInput{
		TableName:            "too-many-vi",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "pk", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
		VectorIndexes:        indexes,
	}
	if _, err := svc.createTableCore(nil, in); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("%d vector indexes: expected ErrInvalidParameter, got %v", len(indexes), err)
	}
}

func TestApplyVectorIndexUpdates(t *testing.T) {
	existing := []*dbstore.VectorIndex{
		{IndexName: "old-idx", VectorAttributeName: "embedding", Dimensions: 3, DistanceFunction: "COSINE"},
	}

	created, createdNames, deletedNames, err := applyVectorIndexUpdates("arn:aws:dynamodb:us-east-1:123456789012:table/t", existing, []interface{}{
		map[string]interface{}{"Create": vectorIndexParams("new-idx", "embedding", float64(3), "EUCLIDEAN")},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(created) != 2 || createdNames[0] != "new-idx" || len(deletedNames) != 0 {
		t.Fatalf("create result: idx=%v created=%v deleted=%v", created, createdNames, deletedNames)
	}
	var createdIdx *dbstore.VectorIndex
	for _, vi := range created {
		if vi.IndexName == "new-idx" {
			createdIdx = vi
		}
	}
	if createdIdx == nil || createdIdx.IndexArn != "arn:aws:dynamodb:us-east-1:123456789012:table/t/index/new-idx" {
		t.Errorf("created index ARN: %+v", createdIdx)
	}

	afterDelete, _, deletedNames, err := applyVectorIndexUpdates("arn", existing, []interface{}{
		map[string]interface{}{"Delete": map[string]interface{}{"IndexName": "old-idx"}},
	})
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if len(afterDelete) != 0 || deletedNames[0] != "old-idx" {
		t.Fatalf("delete result: idx=%v deleted=%v", afterDelete, deletedNames)
	}

	if _, _, _, err := applyVectorIndexUpdates("arn", existing, []interface{}{
		map[string]interface{}{"Delete": map[string]interface{}{"IndexName": "old-idx"}},
		map[string]interface{}{"Delete": map[string]interface{}{"IndexName": "other"}},
	}); err == nil {
		t.Errorf("two updates in one request: expected rejection")
	}

	if _, _, _, err := applyVectorIndexUpdates("arn", existing, []interface{}{
		map[string]interface{}{"Create": vectorIndexParams("old-idx", "embedding", float64(3), "COSINE")},
	}); err == nil {
		t.Errorf("create over existing name: expected rejection")
	}

	if _, _, _, err := applyVectorIndexUpdates("arn", existing, []interface{}{
		map[string]interface{}{"Delete": map[string]interface{}{"IndexName": "absent"}},
	}); err == nil {
		t.Errorf("delete of absent index: expected rejection")
	}
}

func TestValidateVectorAttributeDimensions(t *testing.T) {
	agree := []*dbstore.VectorIndex{
		{IndexName: "a", VectorAttributeName: "embedding", Dimensions: 3, DistanceFunction: "COSINE"},
		{IndexName: "b", VectorAttributeName: "embedding", Dimensions: 3, DistanceFunction: "EUCLIDEAN"},
	}
	if err := validateVectorAttributeDimensions(agree); err != nil {
		t.Errorf("agreeing dimensions rejected: %v", err)
	}

	disagree := []*dbstore.VectorIndex{
		{IndexName: "a", VectorAttributeName: "embedding", Dimensions: 3, DistanceFunction: "COSINE"},
		{IndexName: "b", VectorAttributeName: "embedding", Dimensions: 4, DistanceFunction: "EUCLIDEAN"},
	}
	if err := validateVectorAttributeDimensions(disagree); err == nil {
		t.Errorf("disagreeing dimensions accepted")
	}
}

func TestValidateIndexNameUniquenessWithVector(t *testing.T) {
	gsi := []*dbstore.GlobalSecondaryIndex{{IndexName: "shared-name"}}
	vector := []*dbstore.VectorIndex{{IndexName: "shared-name", VectorAttributeName: "embedding", Dimensions: 3, DistanceFunction: "COSINE"}}
	if err := validateIndexNameUniqueness(gsi, nil, vector); err == nil {
		t.Errorf("vector name clashing with a GSI name accepted")
	}
	if err := validateIndexNameUniqueness(nil, nil, vector); err != nil {
		t.Errorf("unique vector name rejected: %v", err)
	}
}

func TestVectorWriteCapacityForItems(t *testing.T) {
	table := &dbstore.Table{
		Name: "VecTbl",
		VectorIndexes: []*dbstore.VectorIndex{
			{IndexName: "vec", VectorAttributeName: "embedding", Dimensions: 2, DistanceFunction: "COSINE"},
		},
	}

	vecAttr := func(vals ...float64) *dbstore.AttributeValue {
		l := make([]*dbstore.AttributeValue, len(vals))
		for i, v := range vals {
			s := fmt.Sprintf("%g", v)
			l[i] = &dbstore.AttributeValue{N: &s}
		}
		return &dbstore.AttributeValue{L: l}
	}

	// An indexable vector reports positive per-index write bytes; two items
	// sum their serialised lengths.
	one := &dbstore.Item{Attributes: map[string]*dbstore.AttributeValue{"embedding": vecAttr(1, 0)}}
	got := vectorWriteCapacityForItems(table, one)
	if got == nil {
		t.Fatalf("expected vector write capacity for an indexable item")
	}
	cap1 := got["vec"].(map[string]interface{})["VectorWriteRequestBytes"].(float64)
	if cap1 <= 0 {
		t.Fatalf("VectorWriteRequestBytes = %v, want > 0", cap1)
	}
	got = vectorWriteCapacityForItems(table, one, one)
	if sum := got["vec"].(map[string]interface{})["VectorWriteRequestBytes"].(float64); sum != 2*cap1 {
		t.Fatalf("two items = %v, want 2×%v", sum, cap1)
	}

	// Dimension-mismatched, missing, and non-list vectors report nothing;
	// a table without vector indexes reports nothing at all.
	if m := vectorWriteCapacityForItems(table, &dbstore.Item{Attributes: map[string]*dbstore.AttributeValue{"embedding": vecAttr(1, 0, 0)}}); m != nil {
		t.Errorf("dimension mismatch reported %v", m)
	}
	if m := vectorWriteCapacityForItems(table, &dbstore.Item{Attributes: map[string]*dbstore.AttributeValue{}}); m != nil {
		t.Errorf("missing attribute reported %v", m)
	}
	if m := vectorWriteCapacityForItems(&dbstore.Table{Name: "Plain"}, one); m != nil {
		t.Errorf("table without vector indexes reported %v", m)
	}
}

func TestBuildVectorIndexesResponse(t *testing.T) {
	vis := []*dbstore.VectorIndex{
		{
			IndexName:           "vec-idx",
			IndexArn:            "arn:aws:dynamodb:us-east-1:123456789012:table/t/index/vec-idx",
			VectorAttributeName: "embedding",
			Dimensions:          3,
			DistanceFunction:    "DOT_PRODUCT",
			Projection:          &dbstore.Projection{ProjectionType: "ALL"},
			SearchSchema: []*dbstore.SearchSchemaElement{
				{AttributeName: "category", SearchSchemaElementType: "HASH"},
			},
			IndexStatus: dbstore.IndexStatusActive,
			Backfilling: true,
			ItemCount:   7,
		},
	}
	resp := buildVectorIndexesResponse(vis)
	if len(resp) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(resp))
	}
	idx := resp[0]
	if idx["IndexName"] != "vec-idx" || idx["Dimensions"].(int64) != 3 || idx["DistanceFunction"] != "DOT_PRODUCT" {
		t.Errorf("scalar fields: %v", idx)
	}
	if idx["IndexStatus"] != "ACTIVE" || idx["Backfilling"] != true {
		t.Errorf("status fields: %v", idx)
	}
	va := idx["VectorAttribute"].(map[string]interface{})
	if va["AttributeName"] != "embedding" {
		t.Errorf("VectorAttribute: %v", va)
	}
	schema := idx["SearchSchema"].([]map[string]interface{})
	if len(schema) != 1 || schema[0]["SearchSchemaElementType"] != "HASH" {
		t.Errorf("SearchSchema: %v", schema)
	}
}
