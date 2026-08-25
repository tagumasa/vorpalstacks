package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"

	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) runSFNAdvancedTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return []TestResult{{Service: "stepfunctions", TestName: "AdvancedSetup", Status: "FAIL", Error: fmt.Sprintf("load config: %v", err)}}
	}
	ddbClient := dynamodb.NewFromConfig(cfg)

	results = append(results, r.RunTest("stepfunctions", "TestState_Pass", func() error {
		def := `{"StartAt":"TestPass","States":{"TestPass":{"Type":"Pass","Result":{"hello":"world"},"End":true}}}`
		result, err := tc.rawTestState(map[string]interface{}{
			"definition": def,
			"stateName":  "TestPass",
			"input":      `{}`,
		})
		if err != nil {
			return err
		}
		if status, _ := result["status"].(string); status != "SUCCEEDED" {
			return fmt.Errorf("expected SUCCEEDED, got %s", status)
		}
		if _, ok := result["output"]; !ok {
			return fmt.Errorf("output is missing")
		}
		return nil
	}))

	mapSMName := fmt.Sprintf("MapSM-%d", time.Now().UnixNano())
	_, mapRoleARN, mapRoleCleanup := tc.createRoleForSM("MapRole")
	defer mapRoleCleanup()

	mapDef := `{"Comment":"map test","StartAt":"Map","States":{"Map":{"Type":"Map","Iterator":{"StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}},"End":true}}}`
	var mapSMARN string
	mapResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(mapSMName),
		Definition: aws.String(mapDef),
		RoleArn:    aws.String(mapRoleARN),
	})
	if err != nil {
		return append(results, TestResult{Service: "stepfunctions", TestName: "MapRunSetup", Status: "FAIL", Error: fmt.Sprintf("create SM: %v", err)})
	}
	mapSMARN = *mapResp.StateMachineArn
	defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(mapSMARN)})

	results = append(results, r.RunTest("stepfunctions", "ListMapRuns", func() error {
		execArn, err := tc.startExecution(mapSMARN, "", `[1,2,3]`)
		if err != nil {
			return fmt.Errorf("start execution: %v", err)
		}
		if _, err := tc.firstMapRunFor(execArn); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeMapRun", func() error {
		execArn, err := tc.startExecution(mapSMARN, "", `[1,2,3]`)
		if err != nil {
			return fmt.Errorf("start execution: %v", err)
		}
		mapRunARN, err := tc.firstMapRunFor(execArn)
		if err != nil {
			return err
		}

		descResp, err := tc.client.DescribeMapRun(tc.ctx, &sfn.DescribeMapRunInput{
			MapRunArn: aws.String(mapRunARN),
		})
		if err != nil {
			return err
		}
		if descResp.MapRunArn == nil || *descResp.MapRunArn == "" {
			return fmt.Errorf("map run ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "UpdateMapRun", func() error {
		execArn, err := tc.startExecution(mapSMARN, "", `[1,2,3]`)
		if err != nil {
			return fmt.Errorf("start execution: %v", err)
		}
		mapRunARN, err := tc.firstMapRunFor(execArn)
		if err != nil {
			return err
		}

		_, err = tc.client.UpdateMapRun(tc.ctx, &sfn.UpdateMapRunInput{
			MapRunArn:      aws.String(mapRunARN),
			MaxConcurrency: aws.Int32(2),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "RedriveExecution", func() error {
		failDef := `{"Comment":"fail then redrive","StartAt":"Fail","States":{"Fail":{"Type":"Fail","Error":"TestError","Cause":"Redrive test"}}}`
		failSMARN, cleanup, err := tc.createRoleBackedSM("RedriveSM", failDef)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer cleanup()

		execArn, err := tc.startExecution(failSMARN, "", `{}`)
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}
		desc, err := tc.awaitTerminal(execArn, 100*time.Millisecond, 30)
		if err != nil {
			return fmt.Errorf("execution did not reach a terminal state before redrive: %v", err)
		}
		if desc.Status == sfntypes.ExecutionStatusSucceeded || desc.Status == sfntypes.ExecutionStatusAborted {
			return fmt.Errorf("expected a failed execution before redrive, got %s", desc.Status)
		}

		redriveResp, err := tc.client.RedriveExecution(tc.ctx, &sfn.RedriveExecutionInput{
			ExecutionArn: aws.String(execArn),
		})
		if err != nil {
			return err
		}
		if redriveResp == nil {
			return fmt.Errorf("redrive response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "RedriveExecution_MultiStateResume", func() error {
		multiDef := `{"Comment":"pass then fail","StartAt":"FirstPass","States":{"FirstPass":{"Type":"Pass","Result":{"step":1},"Next":"ThenFail"},"ThenFail":{"Type":"Fail","Error":"TestError","Cause":"Second state fails"}}}`
		multiSMARN, cleanup, err := tc.createRoleBackedSM("MultiRedriveSM", multiDef)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer cleanup()

		originalArn, err := tc.startExecution(multiSMARN, "", `{}`)
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		desc, err := tc.awaitTerminal(originalArn, 100*time.Millisecond, 30)
		if err != nil {
			return fmt.Errorf("execution did not reach a terminal state before redrive: %v", err)
		}
		if desc.Status != sfntypes.ExecutionStatusFailed {
			return fmt.Errorf("expected FAILED before redrive, got %s", desc.Status)
		}

		redriveResp, err := tc.client.RedriveExecution(tc.ctx, &sfn.RedriveExecutionInput{
			ExecutionArn: aws.String(originalArn),
		})
		if err != nil {
			return fmt.Errorf("redrive: %v", err)
		}
		if redriveResp.RedriveDate == nil {
			return fmt.Errorf("redriveDate is nil in response")
		}

		descAfter, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(originalArn),
		})
		if err != nil {
			return fmt.Errorf("describe after redrive: %v", err)
		}
		if descAfter.RedriveCount == nil || *descAfter.RedriveCount < 1 {
			return fmt.Errorf("expected RedriveCount >= 1, got %v", descAfter.RedriveCount)
		}

		histResp, err := tc.client.GetExecutionHistory(tc.ctx, &sfn.GetExecutionHistoryInput{
			ExecutionArn: aws.String(originalArn),
		})
		if err != nil {
			return fmt.Errorf("get history: %v", err)
		}
		hasRedrivenEvent := false
		for _, evt := range histResp.Events {
			if string(evt.Type) == "ExecutionRedriven" {
				hasRedrivenEvent = true
				break
			}
		}
		if !hasRedrivenEvent {
			return fmt.Errorf("ExecutionRedrived event not found in history")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "RedriveExecution_MapStateResume", func() error {
		mapFailDef := `{"Comment":"map then fail","StartAt":"Map","States":{"Map":{"Type":"Map","ItemsPath":"$","Iterator":{"StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}},"Next":"Fail"},"Fail":{"Type":"Fail","Error":"PostMapFail","Cause":"Fails after map succeeds"}}}`
		mapFailSMARN, cleanup, err := tc.createRoleBackedSM("MapRedriveSM", mapFailDef)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer cleanup()

		originalArn, err := tc.startExecution(mapFailSMARN, "", `[1,2,3]`)
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		desc, err := tc.awaitTerminal(originalArn, 100*time.Millisecond, 30)
		if err != nil {
			return fmt.Errorf("execution did not reach a terminal state before redrive: %v", err)
		}
		if desc.Status != sfntypes.ExecutionStatusFailed {
			return fmt.Errorf("expected FAILED before redrive, got %s", desc.Status)
		}

		_, err = tc.client.RedriveExecution(tc.ctx, &sfn.RedriveExecutionInput{
			ExecutionArn: aws.String(originalArn),
		})
		if err != nil {
			return fmt.Errorf("redrive: %v", err)
		}

		descAfter, err := tc.awaitTerminal(originalArn, 100*time.Millisecond, 30)
		if err != nil {
			return fmt.Errorf("execution did not reach a terminal state after redrive: %v", err)
		}
		if descAfter.Status != sfntypes.ExecutionStatusFailed {
			return fmt.Errorf("expected FAILED after redrive, got %s", descAfter.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "RedriveExecution_ParallelStateResume", func() error {
		parFailDef := `{"Comment":"parallel then fail","StartAt":"Parallel","States":{"Parallel":{"Type":"Parallel","Branches":[{"StartAt":"Pass1","States":{"Pass1":{"Type":"Pass","End":true}}},{"StartAt":"Pass2","States":{"Pass2":{"Type":"Pass","End":true}}}],"Next":"Fail"},"Fail":{"Type":"Fail","Error":"PostParallelFail","Cause":"Fails after parallel succeeds"}}}`
		parFailSMARN, cleanup, err := tc.createRoleBackedSM("ParRedriveSM", parFailDef)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer cleanup()

		originalArn, err := tc.startExecution(parFailSMARN, "", `{}`)
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		desc, err := tc.awaitTerminal(originalArn, 100*time.Millisecond, 30)
		if err != nil {
			return fmt.Errorf("execution did not reach a terminal state before redrive: %v", err)
		}
		if desc.Status != sfntypes.ExecutionStatusFailed {
			return fmt.Errorf("expected FAILED before redrive, got %s", desc.Status)
		}

		_, err = tc.client.RedriveExecution(tc.ctx, &sfn.RedriveExecutionInput{
			ExecutionArn: aws.String(originalArn),
		})
		if err != nil {
			return fmt.Errorf("redrive: %v", err)
		}

		descAfter, err := tc.awaitTerminal(originalArn, 100*time.Millisecond, 30)
		if err != nil {
			return fmt.Errorf("execution did not reach a terminal state after redrive: %v", err)
		}
		if descAfter.Status != sfntypes.ExecutionStatusFailed {
			return fmt.Errorf("expected FAILED after redrive, got %s", descAfter.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DynamoDBTask_UpdateItem_SetAndRemove", func() error {
		tableName := fmt.Sprintf("SFNUpdateItem-%d", time.Now().UnixNano())
		_, err := ddbClient.CreateTable(tc.ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
			},
			KeySchema: []dynamodbtypes.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
			},
			BillingMode: dynamodbtypes.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer ddbClient.DeleteTable(tc.ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

		_, err = ddbClient.PutItem(tc.ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]dynamodbtypes.AttributeValue{
				"pk":  &dynamodbtypes.AttributeValueMemberS{Value: "item1"},
				"val": &dynamodbtypes.AttributeValueMemberS{Value: "original"},
				"num": &dynamodbtypes.AttributeValueMemberN{Value: "10"},
			},
		})
		if err != nil {
			return fmt.Errorf("put item: %v", err)
		}

		updateDef := fmt.Sprintf(`{"StartAt":"Update","States":{"Update":{"Type":"Task","Resource":"arn:aws:states:::dynamodb:updateItem","Parameters":{"TableName":"%s","Key":{"pk":{"S":"item1"}},"UpdateExpression":"SET val=:v REMOVE num","ExpressionAttributeValues":{":v":{"S":"updated"}}},"End":true}}}`, tableName)
		updateSMARN, smCleanup, cerr := tc.createRoleBackedSM("SFNUpdSM", updateDef)
		if cerr != nil {
			return fmt.Errorf("create SM: %v", cerr)
		}
		defer smCleanup()

		execArn, serr := tc.startExecution(updateSMARN, "", `{}`)
		if serr != nil {
			return fmt.Errorf("start exec: %v", serr)
		}
		desc, werr := tc.awaitTerminal(execArn, 100*time.Millisecond, 30)
		if werr != nil {
			return werr
		}
		if desc.Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", desc.Status)
		}

		getResp, err := ddbClient.GetItem(tc.ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]dynamodbtypes.AttributeValue{"pk": &dynamodbtypes.AttributeValueMemberS{Value: "item1"}},
		})
		if err != nil {
			return fmt.Errorf("get item: %v", err)
		}

		valAttr, hasVal := getResp.Item["val"]
		if !hasVal {
			return fmt.Errorf("expected 'val' attribute to exist after SET")
		}
		if s, ok := valAttr.(*dynamodbtypes.AttributeValueMemberS); !ok || s.Value != "updated" {
			return fmt.Errorf("expected val='updated', got %v", valAttr)
		}

		if _, hasNum := getResp.Item["num"]; hasNum {
			return fmt.Errorf("expected 'num' attribute to be removed by REMOVE, but it still exists")
		}

		if _, hasPK := getResp.Item["pk"]; !hasPK {
			return fmt.Errorf("expected 'pk' key attribute to be preserved")
		}

		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DynamoDBTask_UpdateItem_MultiWhitespace", func() error {
		tableName := fmt.Sprintf("SFNMultiWS-%d", time.Now().UnixNano())
		_, err := ddbClient.CreateTable(tc.ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
			},
			KeySchema: []dynamodbtypes.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
			},
			BillingMode: dynamodbtypes.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer ddbClient.DeleteTable(tc.ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

		_, err = ddbClient.PutItem(tc.ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]dynamodbtypes.AttributeValue{
				"pk":  &dynamodbtypes.AttributeValueMemberS{Value: "ws1"},
				"old": &dynamodbtypes.AttributeValueMemberS{Value: "tobedeleted"},
			},
		})
		if err != nil {
			return fmt.Errorf("put item: %v", err)
		}

		updateExpr := `SET\t\tnew=:v  REMOVE\t#old`
		namedOld := "old"
		updateDef := fmt.Sprintf(`{"StartAt":"Update","States":{"Update":{"Type":"Task","Resource":"arn:aws:states:::dynamodb:updateItem","Parameters":{"TableName":"%s","Key":{"pk":{"S":"ws1"}},"UpdateExpression":"%s","ExpressionAttributeValues":{":v":{"S":"set-via-whitespace"}},"ExpressionAttributeNames":{"#old":"%s"}},"End":true}}}`, tableName, updateExpr, namedOld)

		updateSMARN, smCleanup, cerr := tc.createRoleBackedSM("SFNMultiWSSM", updateDef)
		if cerr != nil {
			return fmt.Errorf("create SM: %v", cerr)
		}
		defer smCleanup()

		execArn, serr := tc.startExecution(updateSMARN, "", `{}`)
		if serr != nil {
			return fmt.Errorf("start exec: %v", serr)
		}
		desc, werr := tc.awaitTerminal(execArn, 100*time.Millisecond, 30)
		if werr != nil {
			return werr
		}
		if desc.Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", desc.Status)
		}

		getResp, err := ddbClient.GetItem(tc.ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]dynamodbtypes.AttributeValue{"pk": &dynamodbtypes.AttributeValueMemberS{Value: "ws1"}},
		})
		if err != nil {
			return fmt.Errorf("get item: %v", err)
		}

		if newAttr, hasNew := getResp.Item["new"]; !hasNew {
			return fmt.Errorf("expected 'new' attribute after SET with multi-whitespace")
		} else if s, ok := newAttr.(*dynamodbtypes.AttributeValueMemberS); !ok || s.Value != "set-via-whitespace" {
			return fmt.Errorf("expected new='set-via-whitespace', got %v", newAttr)
		}

		if _, hasOld := getResp.Item["old"]; hasOld {
			return fmt.Errorf("expected 'old' attribute removed via #old placeholder, but it still exists")
		}

		return nil
	}))

	return results
}
