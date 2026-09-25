// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// TestSecondaryIndexQuotas pins the per-table quotas the quotas page
// documents and nothing enforced: at most 20 global secondary indexes per
// table, at most 5 local secondary indexes per table, at most 100
// user-specified projected attribute names combined across a table's
// secondary indexes, and the NonKeyAttributes list's minimum length of 1
// (an explicitly empty INCLUDE projection is invalid).
func TestSecondaryIndexQuotas(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	nonKey := func(n int) []interface{} {
		names := make([]interface{}, n)
		for i := range names {
			names[i] = fmt.Sprintf("attr%02d", i)
		}
		return names
	}
	gsi := func(i, nonKeyCount int) map[string]interface{} {
		projection := map[string]interface{}{"ProjectionType": "ALL"}
		if nonKeyCount > 0 {
			projection = map[string]interface{}{
				"ProjectionType":   "INCLUDE",
				"NonKeyAttributes": nonKey(nonKeyCount),
			}
		}
		return map[string]interface{}{
			"IndexName":  fmt.Sprintf("gsi-%02d", i),
			"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": fmt.Sprintf("g%02d", i), "KeyType": "HASH"}},
			"Projection": projection,
		}
	}

	create := func(name string, gsiCount, gsiNonKey int, lsiCount int, emptyNonKey bool, lastGsiNonKey ...int) error {
		attrDefs := []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}}
		params := map[string]interface{}{
			"TableName":   name,
			"KeySchema":   []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"BillingMode": "PAY_PER_REQUEST",
		}
		gsis := make([]interface{}, 0, gsiCount)
		for i := 0; i < gsiCount; i++ {
			attrDefs = append(attrDefs, map[string]interface{}{"AttributeName": fmt.Sprintf("g%02d", i), "AttributeType": "S"})
			if emptyNonKey && i == 0 {
				// The minimum-length case: an INCLUDE projection whose
				// NonKeyAttributes list is explicitly empty.
				gsis = append(gsis, map[string]interface{}{
					"IndexName":  fmt.Sprintf("gsi-%02d", i),
					"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": fmt.Sprintf("g%02d", i), "KeyType": "HASH"}},
					"Projection": map[string]interface{}{"ProjectionType": "INCLUDE", "NonKeyAttributes": []interface{}{}},
				})
				continue
			}
			count := gsiNonKey
			if i == gsiCount-1 && len(lastGsiNonKey) > 0 {
				count = lastGsiNonKey[0]
			}
			gsis = append(gsis, gsi(i, count))
		}
		if gsiCount > 0 {
			params["GlobalSecondaryIndexes"] = gsis
		}
		lsis := make([]interface{}, 0, lsiCount)
		for i := 0; i < lsiCount; i++ {
			attrDefs = append(attrDefs, map[string]interface{}{"AttributeName": fmt.Sprintf("l%d", i), "AttributeType": "S"})
			lsis = append(lsis, map[string]interface{}{
				"IndexName":  fmt.Sprintf("lsi-%d", i),
				"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}, map[string]interface{}{"AttributeName": fmt.Sprintf("l%d", i), "KeyType": "RANGE"}},
				"Projection": map[string]interface{}{"ProjectionType": "ALL"},
			})
		}
		if lsiCount > 0 {
			params["LocalSecondaryIndexes"] = lsis
		}
		params["AttributeDefinitions"] = attrDefs
		_, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: params})
		return err
	}

	if err := create("QuotaGsi20Table", dbstore.GlobalSecondaryIndexesPerTable, 0, 0, false); err != nil {
		t.Fatalf("create with %d GSIs: %v", dbstore.GlobalSecondaryIndexesPerTable, err)
	}
	if err := create("QuotaGsi21Table", dbstore.GlobalSecondaryIndexesPerTable+1, 0, 0, false); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("create with %d GSIs: expected ErrInvalidParameter, got %v", dbstore.GlobalSecondaryIndexesPerTable+1, err)
	}

	if err := create("QuotaLsi5Table", 0, 0, dbstore.LocalSecondaryIndexesPerTable, false); err != nil {
		t.Fatalf("create with %d LSIs: %v", dbstore.LocalSecondaryIndexesPerTable, err)
	}
	if err := create("QuotaLsi6Table", 0, 0, dbstore.LocalSecondaryIndexesPerTable+1, false); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("create with %d LSIs: expected ErrInvalidParameter, got %v", dbstore.LocalSecondaryIndexesPerTable+1, err)
	}

	// Five indexes carrying the per-index maximum of 20 names sit exactly
	// at the combined quota; a sixth index's single name crosses it.
	fiveAtMax := dbstore.ProjectedAttributesPerTable / nonKeyAttrListMaxLen
	if err := create("QuotaProj100Table", fiveAtMax, nonKeyAttrListMaxLen, 0, false); err != nil {
		t.Fatalf("create with %d projected names: %v", dbstore.ProjectedAttributesPerTable, err)
	}
	if err := create("QuotaProj101Table", fiveAtMax+1, nonKeyAttrListMaxLen, 0, false, 1); err == nil {
		t.Fatal("create with 101 projected names: expected rejection")
	} else if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("create with 101 projected names: expected ErrInvalidParameter, got %v", err)
	}

	if err := create("QuotaEmptyListTable", 1, 0, 0, true); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("create with an empty NonKeyAttributes list: expected ErrInvalidParameter, got %v", err)
	}

	// The GSI-count bound holds on the update plane too: a table at 20
	// indexes refuses the 21st through UpdateTable.
	if _, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "QuotaGsi20Table",
		"AttributeDefinitions": []interface{}{
			map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"},
			map[string]interface{}{"AttributeName": "gx", "AttributeType": "S"},
		},
		"GlobalSecondaryIndexUpdates": []interface{}{map[string]interface{}{
			"Create": map[string]interface{}{
				"IndexName":  "gsi-extra",
				"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": "gx", "KeyType": "HASH"}},
				"Projection": map[string]interface{}{"ProjectionType": "ALL"},
			},
		}},
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("update adding the 21st GSI: expected ErrInvalidParameter, got %v", err)
	}
}
