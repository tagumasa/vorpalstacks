package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// vectorNumber builds the DynamoDB list-of-numbers wire shape a vector
// attribute carries inside an item.
func vectorNumber(vals ...float64) types.AttributeValue {
	elements := make([]types.AttributeValue, len(vals))
	for i, v := range vals {
		elements[i] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%g", v)}
	}
	return &types.AttributeValueMemberL{Value: elements}
}

// searchVectorElems builds the SearchVector request member: the vector as a
// flat list of number AttributeValues, one element per dimension.
func searchVectorElems(vals ...float64) []types.AttributeValue {
	elements := make([]types.AttributeValue, len(vals))
	for i, v := range vals {
		elements[i] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%g", v)}
	}
	return elements
}

// createVectorTestTable creates a throwaway table with one COSINE vector
// index over the 2-dimensional "embedding" attribute, partitioned by the
// HASH attribute "category" with "year" as an INLINE_FILTER attribute.
func createVectorTestTable(ctx context.Context, client *dynamodb.Client, name string) (func(), error) {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(name),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		},
		BillingMode: types.BillingModePayPerRequest,
		VectorIndexes: []types.VectorIndex{{
			IndexName:        aws.String("vec"),
			VectorAttribute:  &types.VectorAttributeDefinition{AttributeName: aws.String("embedding")},
			Dimensions:       aws.Int64(2),
			DistanceFunction: types.VectorDistanceFunctionCosine,
			Projection:       &types.Projection{ProjectionType: types.ProjectionTypeAll},
			SearchSchema: []types.SearchSchemaElement{
				{AttributeName: aws.String("category"), SearchSchemaElementType: types.SearchSchemaElementTypeHash},
				{AttributeName: aws.String("year"), SearchSchemaElementType: types.SearchSchemaElementTypeInlineFilter},
			},
		}},
	})
	if err != nil {
		return func() {}, fmt.Errorf("create vector table %s: %w", name, err)
	}
	return func() {
		_, _ = client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})
	}, nil
}

func putVectorItem(ctx context.Context, client *dynamodb.Client, table, id string, vec []float64, category string, year int32) error {
	item := map[string]types.AttributeValue{
		"id":       &types.AttributeValueMemberS{Value: id},
		"category": &types.AttributeValueMemberS{Value: category},
		"year":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", year)},
	}
	if vec != nil {
		item["embedding"] = vectorNumber(vec...)
	}
	_, err := client.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(table), Item: item})
	return err
}

func (r *TestRunner) dynamoDBVectorTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "VectorIndexes_SearchVectors", func() error {
		table := fmt.Sprintf("VecTbl-%d", time.Now().UnixNano())
		cleanup, err := createVectorTestTable(ctx, client, table)
		if err != nil {
			return err
		}
		defer cleanup()

		// DescribeTable reports the vector index as ACTIVE.
		desc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(table)})
		if err != nil {
			return err
		}
		if len(desc.Table.VectorIndexes) != 1 {
			return fmt.Errorf("DescribeTable VectorIndexes = %d, want 1", len(desc.Table.VectorIndexes))
		}
		vd := desc.Table.VectorIndexes[0]
		if vd.IndexName == nil || *vd.IndexName != "vec" {
			return fmt.Errorf("VectorIndexDescription IndexName = %v", vd.IndexName)
		}
		if vd.IndexStatus != types.IndexStatusActive {
			return fmt.Errorf("VectorIndexDescription IndexStatus = %s, want ACTIVE", vd.IndexStatus)
		}
		if vd.IndexArn == nil || *vd.IndexArn == "" {
			return fmt.Errorf("VectorIndexDescription IndexArn is empty")
		}

		if err := putVectorItem(ctx, client, table, "same", []float64{1, 0}, "a", 2020); err != nil {
			return fmt.Errorf("put same: %w", err)
		}
		if err := putVectorItem(ctx, client, table, "orth", []float64{0, 1}, "a", 2024); err != nil {
			return fmt.Errorf("put orth: %w", err)
		}
		if err := putVectorItem(ctx, client, table, "near-b", []float64{0.9, 0.1}, "b", 2020); err != nil {
			return fmt.Errorf("put near-b: %w", err)
		}
		// A dimension-mismatched vector attribute is written but not
		// indexed.
		if err := putVectorItem(ctx, client, table, "bad-dim", []float64{1, 0, 0}, "a", 2020); err != nil {
			return fmt.Errorf("put bad-dim: %w", err)
		}

		// A write that indexes a vector reports per-index write bytes in
		// ConsumedCapacity; the dimension-mismatched item reports none.
		putResp, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName:              aws.String(table),
			ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
			Item: map[string]types.AttributeValue{
				"id":        &types.AttributeValueMemberS{Value: "cap"},
				"category":  &types.AttributeValueMemberS{Value: "b"},
				"year":      &types.AttributeValueMemberN{Value: "2020"},
				"embedding": vectorNumber(0.5, 0.5),
			},
		})
		if err != nil {
			return fmt.Errorf("put with capacity: %w", err)
		}
		if putResp.ConsumedCapacity == nil {
			return fmt.Errorf("PutItem ConsumedCapacity is nil")
		}
		vecCap, ok := putResp.ConsumedCapacity.VectorIndexes["vec"]
		if !ok {
			return fmt.Errorf("ConsumedCapacity.VectorIndexes has no vec entry: %+v", putResp.ConsumedCapacity.VectorIndexes)
		}
		if vecCap.VectorWriteRequestBytes == nil || *vecCap.VectorWriteRequestBytes <= 0 {
			return fmt.Errorf("VectorWriteRequestBytes = %v, want > 0", vecCap.VectorWriteRequestBytes)
		}

		// The index is partitioned by the HASH attribute "category": the
		// partition-key value must be provided in the condition expression,
		// so an unconditional search is a ValidationException.
		_, err = client.SearchVectors(ctx, &dynamodb.SearchVectorsInput{
			TableName:    aws.String(table),
			IndexName:    aws.String("vec"),
			SearchVector: searchVectorElems(1, 0),
			TopK:         aws.Int32(3),
		})
		if e := expectAWSErrorCode(err, "ValidationException"); e != nil {
			return fmt.Errorf("unconditional search on partitioned index: %v", e)
		}

		// Equality is the only operator the condition expression supports —
		// a comparison on the INLINE_FILTER attribute is rejected.
		_, err = client.SearchVectors(ctx, &dynamodb.SearchVectorsInput{
			TableName:                 aws.String(table),
			IndexName:                 aws.String("vec"),
			SearchVector:              searchVectorElems(1, 0),
			TopK:                      aws.Int32(3),
			SearchConditionExpression: aws.String("category = :cat AND year > :y"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":cat": &types.AttributeValueMemberS{Value: "a"},
				":y":   &types.AttributeValueMemberN{Value: "2020"},
			},
		})
		if e := expectAWSErrorCode(err, "ValidationException"); e != nil {
			return fmt.Errorf("comparison operator in search condition: %v", e)
		}

		searchCategoryA := func(tableRef string) error {
			resp, err := client.SearchVectors(ctx, &dynamodb.SearchVectorsInput{
				TableName:                 aws.String(tableRef),
				IndexName:                 aws.String("vec"),
				SearchVector:              searchVectorElems(1, 0),
				TopK:                      aws.Int32(3),
				SearchConditionExpression: aws.String("category = :cat"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":cat": &types.AttributeValueMemberS{Value: "a"},
				},
			})
			if err != nil {
				return err
			}
			if len(resp.SearchResults) != 2 {
				return fmt.Errorf("SearchResults = %d, want 2 (near-b excluded by the partition, bad-dim never indexed)", len(resp.SearchResults))
			}
			first := resp.SearchResults[0].Item["id"].(*types.AttributeValueMemberS)
			if first.Value != "same" {
				return fmt.Errorf("first result = %s, want same", first.Value)
			}
			// COSINE scores are distances: identical 0, then ascending.
			if resp.SearchResults[0].Score > 1e-9 {
				return fmt.Errorf("first score = %g, want 0", resp.SearchResults[0].Score)
			}
			if resp.SearchResults[1].Score < resp.SearchResults[0].Score {
				return fmt.Errorf("scores not ascending: %v", resp.SearchResults)
			}
			return nil
		}

		if err := searchCategoryA(table); err != nil {
			return fmt.Errorf("conditional search: %w", err)
		}

		// The TableName member also accepts the table's ARN.
		if desc.Table.TableArn == nil || *desc.Table.TableArn == "" {
			return fmt.Errorf("DescribeTable TableArn is empty")
		}
		if err := searchCategoryA(*desc.Table.TableArn); err != nil {
			return fmt.Errorf("search by table ARN: %w", err)
		}

		// Vector elements are 32-bit IEEE-754 values: a query element
		// rounds to its float32 value, so the decimal "0.1" and the exact
		// float32 decimal are the same element and the identical-vector
		// item scores exactly 0; two item vectors differing only below
		// float32 precision score identically.
		if err := putVectorItem(ctx, client, table, "f32", []float64{0.100000001490116119384765625, 0.2}, "b", 2020); err != nil {
			return fmt.Errorf("put f32: %w", err)
		}
		if err := putVectorItem(ctx, client, table, "sub32a", []float64{0.3, 0.3}, "c", 2020); err != nil {
			return fmt.Errorf("put sub32a: %w", err)
		}
		if err := putVectorItem(ctx, client, table, "sub32b", []float64{0.300000001, 0.300000001}, "c", 2020); err != nil {
			return fmt.Errorf("put sub32b: %w", err)
		}

		searchCategory := func(category string) (*dynamodb.SearchVectorsOutput, error) {
			return client.SearchVectors(ctx, &dynamodb.SearchVectorsInput{
				TableName:                 aws.String(table),
				IndexName:                 aws.String("vec"),
				SearchVector:              searchVectorElems(0.1, 0.2),
				TopK:                      aws.Int32(3),
				SearchConditionExpression: aws.String("category = :cat"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":cat": &types.AttributeValueMemberS{Value: category},
				},
			})
		}

		respB, err := searchCategory("b")
		if err != nil {
			return fmt.Errorf("category-b float32 search: %w", err)
		}
		var f32Score *float64
		for _, r := range respB.SearchResults {
			if id, ok := r.Item["id"].(*types.AttributeValueMemberS); ok && id.Value == "f32" {
				s := r.Score
				f32Score = &s
			}
		}
		if f32Score == nil {
			return fmt.Errorf("category-b search results carry no f32 item: %+v", respB.SearchResults)
		}
		// The tolerance is the same one the identical-vector assertion above
		// uses: the cosine arithmetic over identical vectors leaves a
		// machine-epsilon residue rather than an exact zero.
		if *f32Score > 1e-9 {
			return fmt.Errorf("f32 score = %g, want ~0 (0.1 and its float32 decimal are the same element)", *f32Score)
		}

		respC, err := searchCategory("c")
		if err != nil {
			return fmt.Errorf("category-c sub-float32 search: %w", err)
		}
		if len(respC.SearchResults) != 2 {
			return fmt.Errorf("category-c results = %d, want 2", len(respC.SearchResults))
		}
		if respC.SearchResults[0].Score != respC.SearchResults[1].Score {
			return fmt.Errorf("sub-float32 vectors distinguished: %g vs %g",
				respC.SearchResults[0].Score, respC.SearchResults[1].Score)
		}

		// A condition naming only the INLINE_FILTER attribute still misses
		// the partition key and is rejected.
		_, err = client.SearchVectors(ctx, &dynamodb.SearchVectorsInput{
			TableName:                 aws.String(table),
			IndexName:                 aws.String("vec"),
			SearchVector:              searchVectorElems(1, 0),
			TopK:                      aws.Int32(3),
			SearchConditionExpression: aws.String("year = :y"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":y": &types.AttributeValueMemberN{Value: "2020"},
			},
		})
		if e := expectAWSErrorCode(err, "ValidationException"); e != nil {
			return fmt.Errorf("condition without the partition key: %v", e)
		}

		// Searching an absent index is ResourceNotFoundException.
		_, err = client.SearchVectors(ctx, &dynamodb.SearchVectorsInput{
			TableName:    aws.String(table),
			IndexName:    aws.String("absent"),
			SearchVector: searchVectorElems(1, 0),
			TopK:         aws.Int32(1),
		})
		if rnf := expectResourceNotFound(err); rnf != nil {
			return fmt.Errorf("absent index: %v", rnf)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "VectorIndexes_UpdateTableLifecycle", func() error {
		table := fmt.Sprintf("VecUpd-%d", time.Now().UnixNano())
		cleanup, err := createDynamoTestTable(ctx, client, table, withDynamoHashKey("id"))
		if err != nil {
			return err
		}
		defer cleanup()

		// Items written before the index exists.
		for _, id := range []string{"k1", "k2"} {
			if err := putVectorItem(ctx, client, table, id, []float64{1, 0}, "a", 2020); err != nil {
				return fmt.Errorf("put %s: %w", id, err)
			}
		}

		// Attaching the index via UpdateTable backfills existing items.
		newIndex := types.CreateVectorIndexAction{
			IndexName:        aws.String("late-vec"),
			VectorAttribute:  &types.VectorAttributeDefinition{AttributeName: aws.String("embedding")},
			Dimensions:       aws.Int64(2),
			DistanceFunction: types.VectorDistanceFunctionEuclidean,
			Projection:       &types.Projection{ProjectionType: types.ProjectionTypeAll},
		}
		if _, err := client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
			TableName:          aws.String(table),
			VectorIndexUpdates: []types.VectorIndexUpdate{{Create: &newIndex}},
		}); err != nil {
			return fmt.Errorf("update table create vector index: %w", err)
		}

		checkBackfill := func() error {
			var lastErr error
			for attempt := 0; attempt < 10; attempt++ {
				resp, err := client.SearchVectors(ctx, &dynamodb.SearchVectorsInput{
					TableName:    aws.String(table),
					IndexName:    aws.String("late-vec"),
					SearchVector: searchVectorElems(1, 0),
					TopK:         aws.Int32(10),
				})
				if err == nil && len(resp.SearchResults) == 2 {
					return nil
				}
				lastErr = err
				time.Sleep(100 * time.Millisecond)
			}
			if lastErr != nil {
				return fmt.Errorf("backfilled search: %w", lastErr)
			}
			return fmt.Errorf("backfilled search did not observe both items in time")
		}
		if err := checkBackfill(); err != nil {
			return err
		}

		// Deleting the index removes its data: further searches are
		// ResourceNotFoundException.
		if _, err := client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
			TableName: aws.String(table),
			VectorIndexUpdates: []types.VectorIndexUpdate{
				{Delete: &types.DeleteVectorIndexAction{IndexName: aws.String("late-vec")}},
			},
		}); err != nil {
			return fmt.Errorf("update table delete vector index: %w", err)
		}
		_, err = client.SearchVectors(ctx, &dynamodb.SearchVectorsInput{
			TableName:    aws.String(table),
			IndexName:    aws.String("late-vec"),
			SearchVector: searchVectorElems(1, 0),
			TopK:         aws.Int32(1),
		})
		if rnf := expectResourceNotFound(err); rnf != nil {
			return fmt.Errorf("search after index delete: %v", rnf)
		}

		// Re-creating the same name repopulates from the surviving items,
		// proving the deleted index left no stale entries behind.
		if _, err := client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
			TableName:          aws.String(table),
			VectorIndexUpdates: []types.VectorIndexUpdate{{Create: &newIndex}},
		}); err != nil {
			return fmt.Errorf("re-create vector index: %w", err)
		}
		return checkBackfill()
	}))

	return results
}
