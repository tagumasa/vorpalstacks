package testutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *TestRunner) dynamoDBBasicPartiQLTests(ctx context.Context, client *dynamodb.Client, tableName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement (PartiQL)", func() error {
		resp, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("INSERT INTO \"" + tableName + "\" VALUE {'id': 'partiql1', 'name': 'PartiQL Item'}"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("ExecuteStatement response is nil")
		}
		verifyResp, verifyErr := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "partiql1"},
			},
		})
		if verifyErr != nil {
			return fmt.Errorf("PartiQL INSERT verification failed: %w", verifyErr)
		}
		if verifyResp.Item == nil {
			return fmt.Errorf("item not found after PartiQL INSERT")
		}
		name, ok := verifyResp.Item["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "PartiQL Item" {
			return fmt.Errorf("expected name='PartiQL Item', got %v", verifyResp.Item["name"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement (SELECT)", func() error {
		resp, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("SELECT * FROM \"" + tableName + "\" WHERE id = 'partiql1'"),
		})
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			return fmt.Errorf("no items found")
		}
		name, ok := resp.Items[0]["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "PartiQL Item" {
			return fmt.Errorf("expected name='PartiQL Item', got %v", resp.Items[0]["name"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "BatchExecuteStatement", func() error {
		resp, err := client.BatchExecuteStatement(ctx, &dynamodb.BatchExecuteStatementInput{
			Statements: []types.BatchStatementRequest{
				{
					Statement: aws.String("UPDATE \"" + tableName + "\" SET #n = :name WHERE id = 'batch1'"),
					Parameters: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: "Updated via Batch"},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Responses == nil {
			return fmt.Errorf("BatchExecuteStatement Responses is nil")
		}
		if len(resp.Responses) == 0 {
			return fmt.Errorf("expected at least one response in BatchExecuteStatement")
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBPartiQLInsertParamsTest(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	// PartiQL INSERT with parameterised values (?) substitutes placeholders correctly.
	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_InsertWithParams", func() error {
		paramTable := fmt.Sprintf("ParamInsert-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, paramTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		_, err = client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("INSERT INTO \"" + paramTable + "\" VALUE {'id': ?, 'name': ?}"),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "param1"},
				&types.AttributeValueMemberS{Value: "Parametrised Item"},
			},
		})
		if err != nil {
			return fmt.Errorf("parametrised INSERT failed: %v", err)
		}
		getResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(paramTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "param1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get after param insert: %v", err)
		}
		if getResp.Item == nil {
			return fmt.Errorf("item not found after parametrised INSERT")
		}
		name, ok := getResp.Item["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "Parametrised Item" {
			return fmt.Errorf("expected name='Parametrised Item', got %v", getResp.Item["name"])
		}

		// The parameter list must carry exactly one value per placeholder:
		// both the over-count and the under-count direction are rejected.
		if _, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("INSERT INTO \"" + paramTable + "\" VALUE {'id': ?, 'name': ?}"),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "over"},
				&types.AttributeValueMemberS{Value: "count"},
				&types.AttributeValueMemberS{Value: "extra"},
			},
		}); err == nil {
			return fmt.Errorf("more parameters than placeholders must be rejected")
		} else if e := expectAWSErrorCode(err, "ValidationException"); e != nil {
			return e
		}
		if _, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement:  aws.String("INSERT INTO \"" + paramTable + "\" VALUE {'id': ?, 'name': ?}"),
			Parameters: []types.AttributeValue{&types.AttributeValueMemberS{Value: "under"}},
		}); err == nil {
			return fmt.Errorf("fewer parameters than placeholders must be rejected")
		} else if e := expectAWSErrorCode(err, "ValidationException"); e != nil {
			return e
		}

		// Number literals follow the DynamoDB Number contract: an
		// out-of-range magnitude is a ValidationException, never a stored
		// value, and a signed literal stores as a Number.
		if _, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("INSERT INTO \"" + paramTable + "\" VALUE {'id': 'num1', 'n': 1e-500}"),
		}); err == nil {
			return fmt.Errorf("an out-of-range number literal must be rejected")
		} else if e := expectAWSErrorCode(err, "ValidationException"); e != nil {
			return e
		}
		if _, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("INSERT INTO \"" + paramTable + "\" VALUE {'id': 'num2', 'n': 1234567890123456789012345678901234567890}"),
		}); err == nil {
			return fmt.Errorf("a 39-significant-digit literal must be rejected")
		} else if e := expectAWSErrorCode(err, "ValidationException"); e != nil {
			return e
		}
		if _, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("INSERT INTO \"" + paramTable + "\" VALUE {'id': 'num3', 'n': -5}"),
		}); err != nil {
			return fmt.Errorf("signed number literal INSERT failed: %v", err)
		}
		numResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(paramTable),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "num3"}},
		})
		if err != nil {
			return fmt.Errorf("get signed literal: %v", err)
		}
		if n, ok := numResp.Item["n"].(*types.AttributeValueMemberN); !ok || n.Value != "-5" {
			return fmt.Errorf("expected n=-5, got %v", numResp.Item["n"])
		}
		return nil
	}))

	// ADD-clause parameters bind in statement order across the clause and
	// the WHERE: the value parameter drives the ADD and the key parameter
	// drives the equality, with no cross-binding between statement segments.
	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_AddClauseParamsBindInOrder", func() error {
		addTable := fmt.Sprintf("PQAddParam-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, addTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(addTable),
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "add1"},
				"n":  &types.AttributeValueMemberN{Value: "1"},
			},
		}); err != nil {
			return err
		}

		if _, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("UPDATE \"" + addTable + "\" ADD n ? WHERE id = ?"),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberN{Value: "5"},
				&types.AttributeValueMemberS{Value: "add1"},
			},
		}); err != nil {
			return fmt.Errorf("parametrised ADD update failed: %v", err)
		}

		getResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(addTable),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "add1"}},
		})
		if err != nil {
			return fmt.Errorf("get after ADD update: %v", err)
		}
		n, ok := getResp.Item["n"].(*types.AttributeValueMemberN)
		if !ok || n.Value != "6" {
			return fmt.Errorf("expected n=6 after ADD 5, got %v", getResp.Item["n"])
		}
		return nil
	}))

	// ExecuteTransaction binds per-statement parameters on INSERT too: the
	// stored item holds the bound values, not the placeholder text.
	results = append(results, r.RunTest("dynamodb", "ExecuteTransaction_InsertWithParams", func() error {
		txTable := fmt.Sprintf("ParamTxInsert-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, txTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		_, err = client.ExecuteTransaction(ctx, &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{
					Statement: aws.String("INSERT INTO \"" + txTable + "\" VALUE {'id': ?, 'name': ?}"),
					Parameters: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: "txparam1"},
						&types.AttributeValueMemberS{Value: "Transaction Parametrised Item"},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("ExecuteTransaction parametrised INSERT failed: %v", err)
		}
		getResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(txTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "txparam1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get after transactional param insert: %v", err)
		}
		if getResp.Item == nil {
			return fmt.Errorf("item not found after transactional parametrised INSERT")
		}
		name, ok := getResp.Item["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "Transaction Parametrised Item" {
			return fmt.Errorf("expected name='Transaction Parametrised Item', got %v", getResp.Item["name"])
		}
		return nil
	}))

	// Every AttributeValue union member round-trips through parameter
	// binding, including placeholders nested inside objects and lists.
	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_InsertWithParams_AllTypes", func() error {
		atTable := fmt.Sprintf("ParamAllTypes-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, atTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		binary := []byte{0x01, 0x02, 0x03}
		_, err = client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String(`INSERT INTO "` + atTable + `"
				VALUE {'id': ?, 's': ?, 'n': ?, 'b': ?, 'bool': ?, 'null': ?,
					'ss': ?, 'ns': ?, 'bs': ?, 'l': ?, 'm': ?,
					'nested': {'x': ?, 'inner': [?]}}`),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "all1"},
				&types.AttributeValueMemberS{Value: "str"},
				&types.AttributeValueMemberN{Value: "42"},
				&types.AttributeValueMemberB{Value: binary},
				&types.AttributeValueMemberBOOL{Value: true},
				&types.AttributeValueMemberNULL{Value: true},
				&types.AttributeValueMemberSS{Value: []string{"x", "y"}},
				&types.AttributeValueMemberNS{Value: []string{"1", "2"}},
				&types.AttributeValueMemberBS{Value: [][]byte{{0x03}}},
				&types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberN{Value: "7"},
				}},
				&types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"inner": &types.AttributeValueMemberS{Value: "v"},
				}},
				&types.AttributeValueMemberS{Value: "nx"},
				&types.AttributeValueMemberN{Value: "9"},
			},
		})
		if err != nil {
			return fmt.Errorf("all-types parametrised INSERT failed: %v", err)
		}
		getResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(atTable),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "all1"}},
		})
		if err != nil {
			return fmt.Errorf("get after all-types param insert: %v", err)
		}
		if getResp.Item == nil {
			return fmt.Errorf("item not found after all-types parametrised INSERT")
		}
		item := getResp.Item
		if v, ok := item["s"].(*types.AttributeValueMemberS); !ok || v.Value != "str" {
			return fmt.Errorf("expected s='str', got %v", item["s"])
		}
		if v, ok := item["n"].(*types.AttributeValueMemberN); !ok || v.Value != "42" {
			return fmt.Errorf("expected n='42', got %v", item["n"])
		}
		if v, ok := item["b"].(*types.AttributeValueMemberB); !ok || !bytes.Equal(v.Value, binary) {
			return fmt.Errorf("expected b=%v, got %v", binary, item["b"])
		}
		if v, ok := item["bool"].(*types.AttributeValueMemberBOOL); !ok || !v.Value {
			return fmt.Errorf("expected bool=true, got %v", item["bool"])
		}
		if _, ok := item["null"].(*types.AttributeValueMemberNULL); !ok {
			return fmt.Errorf("expected NULL member, got %v", item["null"])
		}
		if v, ok := item["ss"].(*types.AttributeValueMemberSS); !ok || len(v.Value) != 2 {
			return fmt.Errorf("expected ss of 2, got %v", item["ss"])
		}
		if v, ok := item["ns"].(*types.AttributeValueMemberNS); !ok || len(v.Value) != 2 {
			return fmt.Errorf("expected ns of 2, got %v", item["ns"])
		}
		if v, ok := item["bs"].(*types.AttributeValueMemberBS); !ok || len(v.Value) != 1 || v.Value[0][0] != 0x03 {
			return fmt.Errorf("expected bs=[[3]], got %v", item["bs"])
		}
		if v, ok := item["l"].(*types.AttributeValueMemberL); !ok || len(v.Value) != 1 {
			return fmt.Errorf("expected l of 1, got %v", item["l"])
		} else if n, ok := v.Value[0].(*types.AttributeValueMemberN); !ok || n.Value != "7" {
			return fmt.Errorf("expected l[0]='7', got %v", v.Value[0])
		}
		if v, ok := item["m"].(*types.AttributeValueMemberM); !ok {
			return fmt.Errorf("expected m, got %v", item["m"])
		} else if s, ok := v.Value["inner"].(*types.AttributeValueMemberS); !ok || s.Value != "v" {
			return fmt.Errorf("expected m.inner='v', got %v", v.Value["inner"])
		}
		nested, ok := item["nested"].(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("expected nested map, got %v", item["nested"])
		}
		if s, ok := nested.Value["x"].(*types.AttributeValueMemberS); !ok || s.Value != "nx" {
			return fmt.Errorf("expected nested.x='nx', got %v", nested.Value["x"])
		}
		inner, ok := nested.Value["inner"].(*types.AttributeValueMemberL)
		if !ok || len(inner.Value) != 1 {
			return fmt.Errorf("expected nested.inner list of 1, got %v", nested.Value["inner"])
		}
		if n, ok := inner.Value[0].(*types.AttributeValueMemberN); !ok || n.Value != "9" {
			return fmt.Errorf("expected nested.inner[0]='9', got %v", inner.Value[0])
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBPartiQLEdgeCaseTests(ctx context.Context, client *dynamodb.Client, compTableName string) []TestResult {
	var results []TestResult

	// === PARTIQL EDGE CASES ===

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_SelectWhere", func() error {
		resp, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("SELECT * FROM \"" + compTableName + "\" WHERE pk = 'user1' AND sk = 'order2'"),
		})
		if err != nil {
			return err
		}
		if len(resp.Items) != 1 {
			return fmt.Errorf("expected 1 item, got %d", len(resp.Items))
		}
		amt, ok := resp.Items[0]["amount"].(*types.AttributeValueMemberN)
		if !ok || amt.Value != "200" {
			return fmt.Errorf("expected amount=200, got %v", resp.Items[0]["amount"])
		}
		status, ok := resp.Items[0]["status"].(*types.AttributeValueMemberS)
		if !ok || status.Value != "pending" {
			return fmt.Errorf("expected status=pending, got %v", resp.Items[0]["status"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_Update", func() error {
		puTable := fmt.Sprintf("PQUpd-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, puTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(puTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "pu1"},
				"val": &types.AttributeValueMemberN{Value: "10"},
			},
		})

		resp, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("UPDATE \"" + puTable + "\" SET val = 20 WHERE id = 'pu1'"),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}

		getResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(puTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "pu1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		val, ok := getResp.Item["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "20" {
			return fmt.Errorf("expected val=20 after PartiQL UPDATE, got %v", val)
		}
		return nil
	}))

	// Key attributes identify the item; UPDATE clauses that write them are
	// rejected on both the standalone and the transactional engine, and the
	// stored item stays reachable under its original key.
	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_UpdateKeyAttributeRejected", func() error {
		kaTable := fmt.Sprintf("PQKeyAttr-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, kaTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(kaTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "ka1"},
				"val": &types.AttributeValueMemberN{Value: "1"},
			},
		}); err != nil {
			return err
		}

		if _, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("UPDATE \"" + kaTable + "\" SET id = 'ka2' WHERE id = 'ka1'"),
		}); err == nil {
			return fmt.Errorf("SET on the partition key must be rejected")
		} else if e := expectAWSErrorCode(err, "ValidationException"); e != nil {
			return e
		}

		if _, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("UPDATE \"" + kaTable + "\" REMOVE id WHERE id = 'ka1'"),
		}); err == nil {
			return fmt.Errorf("REMOVE on the partition key must be rejected")
		} else if e := expectAWSErrorCode(err, "ValidationException"); e != nil {
			return e
		}

		getResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(kaTable),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "ka1"}},
		})
		if err != nil {
			return fmt.Errorf("get after rejection: %v", err)
		}
		if getResp.Item == nil {
			return fmt.Errorf("item must remain under its original key after rejection")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteTransaction_UpdateKeyAttributeRejected", func() error {
		kaTxTable := fmt.Sprintf("PQKeyAttrTx-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, kaTxTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(kaTxTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "ka1"},
				"val": &types.AttributeValueMemberN{Value: "1"},
			},
		}); err != nil {
			return err
		}

		_, err = client.ExecuteTransaction(ctx, &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{
					Statement: aws.String("UPDATE \"" + kaTxTable + "\" SET id = 'ka2' WHERE id = 'ka1'"),
				},
			},
		})
		if err == nil {
			return fmt.Errorf("transactional SET on the partition key must be rejected")
		}
		return expectAWSErrorCode(err, "ValidationException")
	}))

	// ExecuteTransaction UPDATE/DELETE must target a single item: a WHERE
	// without a partition-key equality, or one matching several items, is
	// rejected instead of rewriting or wiping the matching rows.
	results = append(results, r.RunTest("dynamodb", "ExecuteTransaction_SingleItemConstraints", func() error {
		mmTable := fmt.Sprintf("PQSingle-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, mmTable, withDynamoKeySchema(
			[]types.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS},
			},
			[]types.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
				{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange},
			},
		))
		if err != nil {
			return err
		}
		defer cleanupTable()

		for _, sk := range []string{"s1", "s2"} {
			if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(mmTable),
				Item: map[string]types.AttributeValue{
					"pk":   &types.AttributeValueMemberS{Value: "a"},
					"sk":   &types.AttributeValueMemberS{Value: sk},
					"tags": &types.AttributeValueMemberS{Value: "keep"},
				},
			}); err != nil {
				return err
			}
		}

		if _, err := client.ExecuteTransaction(ctx, &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{Statement: aws.String("UPDATE \"" + mmTable + "\" SET tags = 'rewritten' WHERE tags = 'keep'")},
			},
		}); err == nil {
			return fmt.Errorf("UPDATE without a partition-key equality must be rejected")
		} else if e := expectAWSErrorCode(err, "ValidationException"); e != nil {
			return e
		}

		if _, err := client.ExecuteTransaction(ctx, &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{Statement: aws.String("UPDATE \"" + mmTable + "\" SET tags = 'rewritten' WHERE pk = 'a'")},
			},
		}); err == nil {
			return fmt.Errorf("UPDATE matching two items must be rejected")
		} else if e := expectTransactionCanceled(err, "ValidationError"); e != nil {
			return e
		}

		if _, err := client.ExecuteTransaction(ctx, &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{Statement: aws.String("DELETE FROM \"" + mmTable + "\" WHERE pk = 'a'")},
			},
		}); err == nil {
			return fmt.Errorf("DELETE matching two items must be rejected")
		} else if e := expectTransactionCanceled(err, "ValidationError"); e != nil {
			return e
		}

		// A failing singleton operation cancels the whole transaction: the
		// duplicate-key INSERT reports ConditionalCheckFailed at its position,
		// the unaffected statements report the literal code "None", and
		// nothing commits.
		if _, err := client.ExecuteTransaction(ctx, &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{Statement: aws.String("INSERT INTO \"" + mmTable + "\" VALUE {'pk':'z','sk':'z1','tags':'new'}")},
				{Statement: aws.String("INSERT INTO \"" + mmTable + "\" VALUE {'pk':'a','sk':'s1','tags':'dup'}")},
				{Statement: aws.String("DELETE FROM \"" + mmTable + "\" WHERE pk = 'z' AND sk = 'z1'")},
			},
		}); err == nil {
			return fmt.Errorf("duplicate-key INSERT inside a transaction must cancel it")
		} else {
			var tce *types.TransactionCanceledException
			if !errors.As(err, &tce) {
				return fmt.Errorf("expected TransactionCanceledException, got %v", err)
			}
			if len(tce.CancellationReasons) != 3 {
				return fmt.Errorf("cancellation reasons = %d, want 3", len(tce.CancellationReasons))
			}
			for i, want := range []string{"None", "ConditionalCheckFailed", "None"} {
				if code := aws.ToString(tce.CancellationReasons[i].Code); code != want {
					return fmt.Errorf("cancellation reason %d code = %q, want %q", i, code, want)
				}
			}
		}
		out, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(mmTable),
			Key:       map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: "z"}, "sk": &types.AttributeValueMemberS{Value: "z1"}},
		})
		if err != nil {
			return fmt.Errorf("get after cancelled transaction: %v", err)
		}
		if len(out.Item) != 0 {
			return fmt.Errorf("cancelled transaction must not commit its statements")
		}

		scanResp, err := client.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(mmTable)})
		if err != nil {
			return fmt.Errorf("scan after rejections: %v", err)
		}
		if scanResp.Count != 2 {
			return fmt.Errorf("expected both items to survive the rejected statements, got %d", scanResp.Count)
		}
		for _, item := range scanResp.Items {
			if tags, ok := item["tags"].(*types.AttributeValueMemberS); !ok || tags.Value != "keep" {
				return fmt.Errorf("expected tags='keep' to be untouched, got %v", item["tags"])
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_Delete", func() error {
		pdTable := fmt.Sprintf("PQDel-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, pdTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(pdTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "pd1"},
				"val": &types.AttributeValueMemberN{Value: "99"},
			},
		})

		_, err = client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("DELETE FROM \"" + pdTable + "\" WHERE id = 'pd1'"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		getResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(pdTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "pd1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(getResp.Item) != 0 {
			return fmt.Errorf("item should be deleted after PartiQL DELETE")
		}
		return nil
	}))

	// SELECT with column projection returns only the specified columns, not all attributes.
	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_SelectProjection", func() error {
		projTable := fmt.Sprintf("PQProj-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, projTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(projTable),
			Item: map[string]types.AttributeValue{
				"id":     &types.AttributeValueMemberS{Value: "proj1"},
				"name":   &types.AttributeValueMemberS{Value: "Alice"},
				"age":    &types.AttributeValueMemberN{Value: "30"},
				"active": &types.AttributeValueMemberBOOL{Value: true},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("SELECT name, age FROM \"" + projTable + "\" WHERE id = 'proj1'"),
		})
		if err != nil {
			return fmt.Errorf("select projection: %v", err)
		}
		if len(resp.Items) != 1 {
			return fmt.Errorf("expected 1 item, got %d", len(resp.Items))
		}
		item := resp.Items[0]
		if _, ok := item["name"]; !ok {
			return fmt.Errorf("missing projected column 'name'")
		}
		if _, ok := item["age"]; !ok {
			return fmt.Errorf("missing projected column 'age'")
		}
		if _, ok := item["active"]; ok {
			return fmt.Errorf("unprojected column 'active' should not be present")
		}
		if _, ok := item["id"]; ok {
			return fmt.Errorf("unprojected column 'id' should not be present")
		}
		return nil
	}))

	// A NextToken replayed past the end of a shrunken result set returns an
	// empty final page, not an error.
	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_StaleNextToken", func() error {
		staleTable := fmt.Sprintf("PQStale-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, staleTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		for i := 0; i < 3; i++ {
			if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(staleTable),
				Item: map[string]types.AttributeValue{
					"id": &types.AttributeValueMemberS{Value: fmt.Sprintf("stale%d", i)},
				},
			}); err != nil {
				return fmt.Errorf("seed: %v", err)
			}
		}

		stmt := `SELECT id FROM "` + staleTable + `"`
		page1, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String(stmt),
			Limit:     aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("page 1: %v", err)
		}
		if len(page1.Items) != 1 || page1.NextToken == nil {
			return fmt.Errorf("page 1: %d items, token %v", len(page1.Items), page1.NextToken)
		}

		page2, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String(stmt),
			Limit:     aws.Int32(1),
			NextToken: page1.NextToken,
		})
		if err != nil {
			return fmt.Errorf("page 2: %v", err)
		}
		if len(page2.Items) != 1 {
			return fmt.Errorf("page 2: expected 1 item, got %d", len(page2.Items))
		}

		// Shrink the table below the offset the token carries, then replay it.
		for i := 0; i < 3; i++ {
			if _, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
				TableName: aws.String(staleTable),
				Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: fmt.Sprintf("stale%d", i)}},
			}); err != nil {
				return fmt.Errorf("shrink: %v", err)
			}
		}
		stale, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String(stmt),
			Limit:     aws.Int32(1),
			NextToken: page1.NextToken,
		})
		if err != nil {
			return fmt.Errorf("stale token replay must not fail: %v", err)
		}
		if len(stale.Items) != 0 {
			return fmt.Errorf("stale token replay: expected empty page, got %d items", len(stale.Items))
		}
		if stale.NextToken != nil {
			return fmt.Errorf("stale token replay: unexpected further NextToken %v", stale.NextToken)
		}
		return nil
	}))

	return results
}
