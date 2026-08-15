package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *TestRunner) dynamoDBUpdateExpressionTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "UpdateItem_CreateNonExistent", func() error {
		uaTable := fmt.Sprintf("UACreate-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(uaTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(uaTable)})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(uaTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "new1"},
			},
			UpdateExpression: aws.String("SET #n = :name"),
			ExpressionAttributeNames: map[string]string{
				"#n": "name",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":name": &types.AttributeValueMemberS{Value: "CreatedViaUpdate"},
			},
			ReturnValues: types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp.Attributes == nil {
			return fmt.Errorf("expected attributes")
		}
		idVal, ok := resp.Attributes["id"].(*types.AttributeValueMemberS)
		if !ok || idVal.Value != "new1" {
			return fmt.Errorf("expected id=new1, got %v", idVal)
		}
		nameVal, ok := resp.Attributes["name"].(*types.AttributeValueMemberS)
		if !ok || nameVal.Value != "CreatedViaUpdate" {
			return fmt.Errorf("expected name=CreatedViaUpdate, got %v", nameVal)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_IfNotExists", func() error {
		ineTable := fmt.Sprintf("INETable-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(ineTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(ineTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(ineTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "ine1"},
				"val": &types.AttributeValueMemberN{Value: "10"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(ineTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "ine1"},
			},
			UpdateExpression:          aws.String("SET #v = if_not_exists(#v, :zero) + :inc"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":zero": &types.AttributeValueMemberN{Value: "0"}, ":inc": &types.AttributeValueMemberN{Value: "5"}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp.Attributes == nil {
			return fmt.Errorf("expected attributes")
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "15" {
			return fmt.Errorf("expected val=15 (10+5), got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_IfNotExists_NoExisting", func() error {
		ineTable := fmt.Sprintf("INENE-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(ineTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(ineTable)})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(ineTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "inene1"},
			},
			UpdateExpression:          aws.String("SET #v = if_not_exists(#v, :zero) + :inc"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":zero": &types.AttributeValueMemberN{Value: "0"}, ":inc": &types.AttributeValueMemberN{Value: "5"}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "5" {
			return fmt.Errorf("expected val=5 (0+5), got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_Arithmetic", func() error {
		arithTable := fmt.Sprintf("Arith-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(arithTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(arithTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(arithTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "a1"},
				"val": &types.AttributeValueMemberN{Value: "100"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(arithTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "a1"},
			},
			UpdateExpression:          aws.String("SET #v = #v - :dec"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":dec": &types.AttributeValueMemberN{Value: "30"}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "70" {
			return fmt.Errorf("expected val=70 (100-30), got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_Remove", func() error {
		rmTable := fmt.Sprintf("RmTable-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(rmTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(rmTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(rmTable),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "rm1"},
				"name":  &types.AttributeValueMemberS{Value: "Alice"},
				"email": &types.AttributeValueMemberS{Value: "alice@test.com"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(rmTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "rm1"},
			},
			UpdateExpression: aws.String("REMOVE email"),
			ReturnValues:     types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if _, ok := resp.Attributes["email"]; ok {
			return fmt.Errorf("expected 'email' to be removed")
		}
		if _, ok := resp.Attributes["name"]; !ok {
			return fmt.Errorf("expected 'name' to remain")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_AddNumber", func() error {
		addTable := fmt.Sprintf("AddN-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(addTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(addTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(addTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "an1"},
				"val": &types.AttributeValueMemberN{Value: "10"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(addTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "an1"},
			},
			UpdateExpression:          aws.String("ADD #v :inc"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":inc": &types.AttributeValueMemberN{Value: "5"}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "15" {
			return fmt.Errorf("expected val=15, got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_AddStringSet", func() error {
		ssTable := fmt.Sprintf("AddSS-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(ssTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(ssTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(ssTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "ss1"},
				"tags": &types.AttributeValueMemberSS{Value: []string{"a", "b"}},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(ssTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "ss1"},
			},
			UpdateExpression:          aws.String("ADD #t :newTags"),
			ExpressionAttributeNames:  map[string]string{"#t": "tags"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":newTags": &types.AttributeValueMemberSS{Value: []string{"b", "c"}}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		tags, ok := resp.Attributes["tags"].(*types.AttributeValueMemberSS)
		if !ok {
			return fmt.Errorf("expected SS type for tags")
		}
		if len(tags.Value) != 3 {
			return fmt.Errorf("expected 3 tags (a,b,c), got %d", len(tags.Value))
		}
		tagMap := make(map[string]bool)
		for _, t := range tags.Value {
			tagMap[t] = true
		}
		for _, exp := range []string{"a", "b", "c"} {
			if !tagMap[exp] {
				return fmt.Errorf("expected tag %q in set", exp)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_DeleteStringSet", func() error {
		dsTable := fmt.Sprintf("DelSS-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(dsTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(dsTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(dsTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "ds1"},
				"tags": &types.AttributeValueMemberSS{Value: []string{"a", "b", "c"}},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(dsTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "ds1"},
			},
			UpdateExpression:          aws.String("DELETE #t :remove"),
			ExpressionAttributeNames:  map[string]string{"#t": "tags"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":remove": &types.AttributeValueMemberSS{Value: []string{"a", "c"}}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		tags, ok := resp.Attributes["tags"].(*types.AttributeValueMemberSS)
		if !ok {
			return fmt.Errorf("expected SS type for tags")
		}
		if len(tags.Value) != 1 || tags.Value[0] != "b" {
			return fmt.Errorf("expected tags=[b], got %v", tags.Value)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_UpdatedOld", func() error {
		uoTable := fmt.Sprintf("UO-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(uoTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(uoTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(uoTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "uo1"},
				"val":  &types.AttributeValueMemberN{Value: "10"},
				"name": &types.AttributeValueMemberS{Value: "Old"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(uoTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "uo1"},
			},
			UpdateExpression:          aws.String("SET #v = :new"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":new": &types.AttributeValueMemberN{Value: "20"}},
			ReturnValues:              types.ReturnValueUpdatedOld,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp.Attributes == nil {
			return fmt.Errorf("expected updated old attributes")
		}
		if _, ok := resp.Attributes["val"]; !ok {
			return fmt.Errorf("expected 'val' in UPDATED_OLD response")
		}
		if _, ok := resp.Attributes["name"]; ok {
			return fmt.Errorf("did not expect unchanged 'name' in UPDATED_OLD response")
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "10" {
			return fmt.Errorf("expected val=10 in UPDATED_OLD, got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_UpdatedNew", func() error {
		unTable := fmt.Sprintf("UN-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(unTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(unTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(unTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "un1"},
				"val":  &types.AttributeValueMemberN{Value: "10"},
				"name": &types.AttributeValueMemberS{Value: "Old"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(unTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "un1"},
			},
			UpdateExpression:          aws.String("SET #v = :new"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":new": &types.AttributeValueMemberN{Value: "20"}},
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp.Attributes == nil {
			return fmt.Errorf("expected updated new attributes")
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "20" {
			return fmt.Errorf("expected val=20 in UPDATED_NEW, got %v", val)
		}
		if _, ok := resp.Attributes["name"]; ok {
			return fmt.Errorf("did not expect unchanged 'name' in UPDATED_NEW response")
		}
		return nil
	}))

	// UpdateItem list_append appends elements to an existing list attribute.
	results = append(results, r.RunTest("dynamodb", "UpdateItem_ListAppendNested", func() error {
		listTable := fmt.Sprintf("ListAppend-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(listTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(listTable)})

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(listTable),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "la1"},
				"items": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: "a"}}},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(listTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "la1"},
			},
			UpdateExpression: aws.String("SET items = list_append(items, :new)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":new": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: "b"}}},
			},
			ReturnValues: types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("list_append update failed: %v", err)
		}
		items, ok := resp.Attributes["items"].(*types.AttributeValueMemberL)
		if !ok {
			return fmt.Errorf("expected list for items, got %T", resp.Attributes["items"])
		}
		if len(items.Value) != 2 {
			return fmt.Errorf("expected 2 items after list_append, got %d", len(items.Value))
		}
		first, ok := items.Value[0].(*types.AttributeValueMemberS)
		if !ok || first.Value != "a" {
			return fmt.Errorf("expected first element 'a', got %v", items.Value[0])
		}
		second, ok := items.Value[1].(*types.AttributeValueMemberS)
		if !ok || second.Value != "b" {
			return fmt.Errorf("expected second element 'b', got %v", items.Value[1])
		}
		return nil
	}))

	// list_append on a non-List attribute (SS, NS, BS, etc.) must return ValidationException.
	results = append(results, r.RunTest("dynamodb", "UpdateItem_ListAppend_TypeMismatch", func() error {
		laTable := fmt.Sprintf("ListAppendErr-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(laTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(laTable)})

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(laTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "ss1"},
				"tags": &types.AttributeValueMemberSS{Value: []string{"a", "b"}},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(laTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "ss1"},
			},
			UpdateExpression: aws.String("SET tags = list_append(tags, :extra)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":extra": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: "c"}}},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for list_append on SS, got nil")
		}
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

// === CONDITION EXPRESSION TESTS ===
