package testutil

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *TestRunner) dynamoDBTypeTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "PutItem_GetItem_AllScalarTypes", func() error {
		allTypesTable := fmt.Sprintf("AllTypes-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, allTypesTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		binaryData := []byte("\x00\x01\x02\xff\xfe")
		putItem := map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: "alltypes1"},
			"str_val":  &types.AttributeValueMemberS{Value: "hello"},
			"num_val":  &types.AttributeValueMemberN{Value: "3.14"},
			"bin_val":  &types.AttributeValueMemberB{Value: binaryData},
			"bool_val": &types.AttributeValueMemberBOOL{Value: true},
			"null_val": &types.AttributeValueMemberNULL{Value: true},
		}
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(allTypesTable),
			Item:      putItem,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(allTypesTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "alltypes1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Item == nil {
			return fmt.Errorf("item is nil")
		}

		s, ok := resp.Item["str_val"].(*types.AttributeValueMemberS)
		if !ok || s.Value != "hello" {
			return fmt.Errorf("str_val mismatch: got %v", resp.Item["str_val"])
		}
		n, ok := resp.Item["num_val"].(*types.AttributeValueMemberN)
		if !ok || n.Value != "3.14" {
			return fmt.Errorf("num_val mismatch: got %v", resp.Item["num_val"])
		}
		b, ok := resp.Item["bin_val"].(*types.AttributeValueMemberB)
		if !ok || !bytes.Equal(b.Value, binaryData) {
			return fmt.Errorf("bin_val mismatch: got %v", resp.Item["bin_val"])
		}
		bo, ok := resp.Item["bool_val"].(*types.AttributeValueMemberBOOL)
		if !ok || bo.Value != true {
			return fmt.Errorf("bool_val mismatch: got %v", resp.Item["bool_val"])
		}
		nu, ok := resp.Item["null_val"].(*types.AttributeValueMemberNULL)
		if !ok || nu.Value != true {
			return fmt.Errorf("null_val mismatch: got %v", resp.Item["null_val"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_GetItem_SetTypes", func() error {
		setTable := fmt.Sprintf("SetTypes-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, setTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		nsItem := map[string]types.AttributeValue{
			"id":   &types.AttributeValueMemberS{Value: "setitem1"},
			"nums": &types.AttributeValueMemberNS{Value: []string{"1", "2.5", "-10"}},
			"strs": &types.AttributeValueMemberSS{Value: []string{"alpha", "beta"}},
			"bins": &types.AttributeValueMemberBS{Value: [][]byte{{0xCA, 0xFE}, {0xDE, 0xAD}}},
			"map_v": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"inner_str": &types.AttributeValueMemberS{Value: "deep"},
				"inner_num": &types.AttributeValueMemberN{Value: "42"},
			}},
			"list_v": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberN{Value: "1"},
				&types.AttributeValueMemberS{Value: "two"},
				&types.AttributeValueMemberBOOL{Value: false},
			}},
		}
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(setTable),
			Item:      nsItem,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(setTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "setitem1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Item == nil {
			return fmt.Errorf("item is nil")
		}

		ns, ok := resp.Item["nums"].(*types.AttributeValueMemberNS)
		if !ok {
			return fmt.Errorf("nums: expected NS, got %T", resp.Item["nums"])
		}
		nsSet := make(map[string]bool)
		for _, v := range ns.Value {
			nsSet[v] = true
		}
		for _, expected := range []string{"1", "2.5", "-10"} {
			if !nsSet[expected] {
				return fmt.Errorf("nums missing %q: got %v", expected, ns.Value)
			}
		}

		ss, ok := resp.Item["strs"].(*types.AttributeValueMemberSS)
		if !ok {
			return fmt.Errorf("strs: expected SS, got %T", resp.Item["strs"])
		}
		ssSet := make(map[string]bool)
		for _, v := range ss.Value {
			ssSet[v] = true
		}
		for _, expected := range []string{"alpha", "beta"} {
			if !ssSet[expected] {
				return fmt.Errorf("strs missing %q: got %v", expected, ss.Value)
			}
		}

		bs, ok := resp.Item["bins"].(*types.AttributeValueMemberBS)
		if !ok {
			return fmt.Errorf("bins: expected BS, got %T", resp.Item["bins"])
		}
		if len(bs.Value) != 2 {
			return fmt.Errorf("bins: expected 2 elements, got %d", len(bs.Value))
		}

		m, ok := resp.Item["map_v"].(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("map_v: expected M, got %T", resp.Item["map_v"])
		}
		innerS, ok := m.Value["inner_str"].(*types.AttributeValueMemberS)
		if !ok || innerS.Value != "deep" {
			return fmt.Errorf("map_v.inner_str mismatch: got %v", m.Value["inner_str"])
		}
		innerN, ok := m.Value["inner_num"].(*types.AttributeValueMemberN)
		if !ok || innerN.Value != "42" {
			return fmt.Errorf("map_v.inner_num mismatch: got %v", m.Value["inner_num"])
		}

		l, ok := resp.Item["list_v"].(*types.AttributeValueMemberL)
		if !ok || len(l.Value) != 3 {
			return fmt.Errorf("list_v: expected 3 elements, got %d", len(l.Value))
		}
		ln, ok := l.Value[0].(*types.AttributeValueMemberN)
		if !ok || ln.Value != "1" {
			return fmt.Errorf("list_v[0] mismatch: got %v", l.Value[0])
		}
		ls, ok := l.Value[1].(*types.AttributeValueMemberS)
		if !ok || ls.Value != "two" {
			return fmt.Errorf("list_v[1] mismatch: got %v", l.Value[1])
		}
		lb, ok := l.Value[2].(*types.AttributeValueMemberBOOL)
		if !ok || lb.Value != false {
			return fmt.Errorf("list_v[2] mismatch: got %v", l.Value[2])
		}
		return nil
	}))

	return results
}
