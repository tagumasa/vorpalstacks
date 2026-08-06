package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) runSFNAdvancedTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("stepfunctions", "TestState_Pass", func() error {
		def := `{"StartAt":"TestPass","States":{"TestPass":{"Type":"Pass","Result":{"hello":"world"},"End":true}}}`
		body := map[string]interface{}{
			"definition": def,
			"stateName":  "TestPass",
			"input":      `{}`,
		}
		bodyBytes, _ := json.Marshal(body)
		testReq, _ := http.NewRequestWithContext(tc.ctx, "POST", r.endpoint, bytes.NewReader(bodyBytes))
		testReq.Header.Set("Content-Type", "application/x-amz-json-1.0")
		testReq.Header.Set("X-Amz-Target", "AWSStepFunctions.TestState")

		testResp, err := testHTTPClient.Do(testReq)
		if err != nil {
			return err
		}
		defer testResp.Body.Close()
		if testResp.StatusCode != 200 {
			var errBody map[string]interface{}
			json.NewDecoder(testResp.Body).Decode(&errBody)
			return fmt.Errorf("status %d: %v", testResp.StatusCode, errBody)
		}
		var result map[string]interface{}
		json.NewDecoder(testResp.Body).Decode(&result)
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
		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(mapSMARN),
			Input:           aws.String(`[1,2,3]`),
		})
		if err != nil {
			return fmt.Errorf("start execution: %v", err)
		}

		var listResp *sfn.ListMapRunsOutput
		for i := 0; i < 20; i++ {
			time.Sleep(100 * time.Millisecond)
			listResp, err = tc.client.ListMapRuns(tc.ctx, &sfn.ListMapRunsInput{
				ExecutionArn: aws.String(*execResp.ExecutionArn),
			})
			if err != nil {
				return err
			}
			if len(listResp.MapRuns) > 0 {
				break
			}
		}
		if listResp == nil || len(listResp.MapRuns) == 0 {
			return fmt.Errorf("no map runs found after polling")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeMapRun", func() error {
		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(mapSMARN),
			Input:           aws.String(`[1,2,3]`),
		})
		if err != nil {
			return fmt.Errorf("start execution: %v", err)
		}

		var mapRunARN string
		for i := 0; i < 20; i++ {
			time.Sleep(100 * time.Millisecond)
			listResp, err := tc.client.ListMapRuns(tc.ctx, &sfn.ListMapRunsInput{
				ExecutionArn: aws.String(*execResp.ExecutionArn),
			})
			if err != nil {
				return err
			}
			if len(listResp.MapRuns) > 0 && listResp.MapRuns[0].MapRunArn != nil {
				mapRunARN = *listResp.MapRuns[0].MapRunArn
				break
			}
		}
		if mapRunARN == "" {
			return fmt.Errorf("no map runs found after polling")
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
		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(mapSMARN),
			Input:           aws.String(`[1,2,3]`),
		})
		if err != nil {
			return fmt.Errorf("start execution: %v", err)
		}

		var mapRunARN string
		for i := 0; i < 20; i++ {
			time.Sleep(100 * time.Millisecond)
			listResp, err := tc.client.ListMapRuns(tc.ctx, &sfn.ListMapRunsInput{
				ExecutionArn: aws.String(*execResp.ExecutionArn),
			})
			if err != nil {
				return err
			}
			if len(listResp.MapRuns) > 0 && listResp.MapRuns[0].MapRunArn != nil {
				mapRunARN = *listResp.MapRuns[0].MapRunArn
				break
			}
		}
		if mapRunARN == "" {
			return fmt.Errorf("no map runs found after polling")
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
		failSMName := fmt.Sprintf("RedriveSM-%d", time.Now().UnixNano())
		_, failRoleARN, failRoleCleanup := tc.createRoleForSM("RedriveRole")
		defer failRoleCleanup()

		failResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(failSMName),
			Definition: aws.String(failDef),
			RoleArn:    aws.String(failRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		failSMARN := *failResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(failSMARN)})

		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(failSMARN),
			Input:           aws.String(`{}`),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		var execStatus string
		time.Sleep(200 * time.Millisecond)
		for i := 0; i < 30; i++ {
			desc, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
				ExecutionArn: execResp.ExecutionArn,
			})
			if err != nil {
				return fmt.Errorf("describe: %v", err)
			}
			execStatus = string(desc.Status)
			if execStatus != "RUNNING" {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if execStatus == "RUNNING" || execStatus == "" {
			return fmt.Errorf("execution status %q after polling, cannot redrive", execStatus)
		}

		redriveResp, err := tc.client.RedriveExecution(tc.ctx, &sfn.RedriveExecutionInput{
			ExecutionArn: aws.String(*execResp.ExecutionArn),
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
		multiSMName := fmt.Sprintf("MultiRedriveSM-%d", time.Now().UnixNano())
		_, multiRoleARN, multiRoleCleanup := tc.createRoleForSM("MultiRedriveRole")
		defer multiRoleCleanup()

		multiResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(multiSMName),
			Definition: aws.String(multiDef),
			RoleArn:    aws.String(multiRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		multiSMARN := *multiResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(multiSMARN)})

		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(multiSMARN),
			Input:           aws.String(`{}`),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}
		originalArn := *execResp.ExecutionArn

		var execStatus string
		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond)
			desc, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
				ExecutionArn: aws.String(originalArn),
			})
			if err != nil {
				return fmt.Errorf("describe: %v", err)
			}
			execStatus = string(desc.Status)
			if execStatus != "RUNNING" {
				break
			}
		}
		if execStatus != "FAILED" {
			return fmt.Errorf("expected FAILED before redrive, got %s", execStatus)
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
		mapFailSMName := fmt.Sprintf("MapRedriveSM-%d", time.Now().UnixNano())
		_, mapFailRoleARN, mapFailRoleCleanup := tc.createRoleForSM("MapRedriveRole")
		defer mapFailRoleCleanup()

		mapFailResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(mapFailSMName),
			Definition: aws.String(mapFailDef),
			RoleArn:    aws.String(mapFailRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		mapFailSMARN := *mapFailResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(mapFailSMARN)})

		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(mapFailSMARN),
			Input:           aws.String(`[1,2,3]`),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}
		originalArn := *execResp.ExecutionArn

		var execStatus string
		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond)
			desc, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
				ExecutionArn: aws.String(originalArn),
			})
			if err != nil {
				return fmt.Errorf("describe: %v", err)
			}
			execStatus = string(desc.Status)
			if execStatus != "RUNNING" {
				break
			}
		}
		if execStatus != "FAILED" {
			return fmt.Errorf("expected FAILED before redrive, got %s", execStatus)
		}

		_, err = tc.client.RedriveExecution(tc.ctx, &sfn.RedriveExecutionInput{
			ExecutionArn: aws.String(originalArn),
		})
		if err != nil {
			return fmt.Errorf("redrive: %v", err)
		}

		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond)
			desc, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
				ExecutionArn: aws.String(originalArn),
			})
			if err != nil {
				return fmt.Errorf("describe after redrive: %v", err)
			}
			execStatus = string(desc.Status)
			if execStatus != "RUNNING" {
				break
			}
		}
		if execStatus != "FAILED" {
			return fmt.Errorf("expected FAILED after redrive, got %s", execStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "RedriveExecution_ParallelStateResume", func() error {
		parFailDef := `{"Comment":"parallel then fail","StartAt":"Parallel","States":{"Parallel":{"Type":"Parallel","Branches":[{"StartAt":"Pass1","States":{"Pass1":{"Type":"Pass","End":true}}},{"StartAt":"Pass2","States":{"Pass2":{"Type":"Pass","End":true}}}],"Next":"Fail"},"Fail":{"Type":"Fail","Error":"PostParallelFail","Cause":"Fails after parallel succeeds"}}}`
		parFailSMName := fmt.Sprintf("ParRedriveSM-%d", time.Now().UnixNano())
		_, parFailRoleARN, parFailRoleCleanup := tc.createRoleForSM("ParRedriveRole")
		defer parFailRoleCleanup()

		parFailResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(parFailSMName),
			Definition: aws.String(parFailDef),
			RoleArn:    aws.String(parFailRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		parFailSMARN := *parFailResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(parFailSMARN)})

		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(parFailSMARN),
			Input:           aws.String(`{}`),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}
		originalArn := *execResp.ExecutionArn

		var execStatus string
		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond)
			desc, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
				ExecutionArn: aws.String(originalArn),
			})
			if err != nil {
				return fmt.Errorf("describe: %v", err)
			}
			execStatus = string(desc.Status)
			if execStatus != "RUNNING" {
				break
			}
		}
		if execStatus != "FAILED" {
			return fmt.Errorf("expected FAILED before redrive, got %s", execStatus)
		}

		_, err = tc.client.RedriveExecution(tc.ctx, &sfn.RedriveExecutionInput{
			ExecutionArn: aws.String(originalArn),
		})
		if err != nil {
			return fmt.Errorf("redrive: %v", err)
		}

		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond)
			desc, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
				ExecutionArn: aws.String(originalArn),
			})
			if err != nil {
				return fmt.Errorf("describe after redrive: %v", err)
			}
			execStatus = string(desc.Status)
			if execStatus != "RUNNING" {
				break
			}
		}
		if execStatus != "FAILED" {
			return fmt.Errorf("expected FAILED after redrive, got %s", execStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DynamoDBTask_UpdateItem_SetAndRemove", func() error {
		cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		ddbClient := dynamodb.NewFromConfig(cfg)

		tableName := fmt.Sprintf("SFNUpdateItem-%d", time.Now().UnixNano())
		_, err = ddbClient.CreateTable(tc.ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer ddbClient.DeleteTable(tc.ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

		_, err = ddbClient.PutItem(tc.ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"pk":  &types.AttributeValueMemberS{Value: "item1"},
				"val": &types.AttributeValueMemberS{Value: "original"},
				"num": &types.AttributeValueMemberN{Value: "10"},
			},
		})
		if err != nil {
			return fmt.Errorf("put item: %v", err)
		}

		updateDef := fmt.Sprintf(`{"StartAt":"Update","States":{"Update":{"Type":"Task","Resource":"arn:aws:states:::dynamodb:updateItem","Parameters":{"TableName":"%s","Key":{"pk":{"S":"item1"}},"UpdateExpression":"SET val=:v REMOVE num","ExpressionAttributeValues":{":v":{"S":"updated"}}},"End":true}}}`, tableName)
		updateSMName := fmt.Sprintf("SFNUpdSM-%d", time.Now().UnixNano())
		_, updateRoleARN, updateRoleCleanup := tc.createRoleForSM("SFNUpdRole")
		defer updateRoleCleanup()

		smResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(updateSMName),
			Definition: aws.String(updateDef),
			RoleArn:    aws.String(updateRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create SM: %v", err)
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smResp.StateMachineArn})

		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: smResp.StateMachineArn,
			Input:           aws.String(`{}`),
		})
		if err != nil {
			return fmt.Errorf("start exec: %v", err)
		}

		var status string
		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond)
			desc, _ := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
				ExecutionArn: execResp.ExecutionArn,
			})
			status = string(desc.Status)
			if status != "RUNNING" {
				break
			}
		}
		if status != "SUCCEEDED" {
			return fmt.Errorf("expected SUCCEEDED, got %s", status)
		}

		getResp, err := ddbClient.GetItem(tc.ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: "item1"}},
		})
		if err != nil {
			return fmt.Errorf("get item: %v", err)
		}

		valAttr, hasVal := getResp.Item["val"]
		if !hasVal {
			return fmt.Errorf("expected 'val' attribute to exist after SET")
		}
		if s, ok := valAttr.(*types.AttributeValueMemberS); !ok || s.Value != "updated" {
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
		cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		ddbClient := dynamodb.NewFromConfig(cfg)

		tableName := fmt.Sprintf("SFNMultiWS-%d", time.Now().UnixNano())
		_, err = ddbClient.CreateTable(tc.ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer ddbClient.DeleteTable(tc.ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

		_, err = ddbClient.PutItem(tc.ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"pk":  &types.AttributeValueMemberS{Value: "ws1"},
				"old": &types.AttributeValueMemberS{Value: "tobedeleted"},
			},
		})
		if err != nil {
			return fmt.Errorf("put item: %v", err)
		}

		updateExpr := `SET\t\tnew=:v  REMOVE\t#old`
		namedOld := "old"
		updateDef := fmt.Sprintf(`{"StartAt":"Update","States":{"Update":{"Type":"Task","Resource":"arn:aws:states:::dynamodb:updateItem","Parameters":{"TableName":"%s","Key":{"pk":{"S":"ws1"}},"UpdateExpression":"%s","ExpressionAttributeValues":{":v":{"S":"set-via-whitespace"}},"ExpressionAttributeNames":{"#old":"%s"}},"End":true}}}`, tableName, updateExpr, namedOld)

		updateSMName := fmt.Sprintf("SFNMultiWSSM-%d", time.Now().UnixNano())
		_, updateRoleARN, updateRoleCleanup := tc.createRoleForSM("SFNMultiWSRole")
		defer updateRoleCleanup()

		smResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(updateSMName),
			Definition: aws.String(updateDef),
			RoleArn:    aws.String(updateRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create SM: %v", err)
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smResp.StateMachineArn})

		execResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: smResp.StateMachineArn,
			Input:           aws.String(`{}`),
		})
		if err != nil {
			return fmt.Errorf("start exec: %v", err)
		}

		var status string
		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond)
			desc, _ := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
				ExecutionArn: execResp.ExecutionArn,
			})
			status = string(desc.Status)
			if status != "RUNNING" {
				break
			}
		}
		if status != "SUCCEEDED" {
			return fmt.Errorf("expected SUCCEEDED, got %s", status)
		}

		getResp, err := ddbClient.GetItem(tc.ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: "ws1"}},
		})
		if err != nil {
			return fmt.Errorf("get item: %v", err)
		}

		if newAttr, hasNew := getResp.Item["new"]; !hasNew {
			return fmt.Errorf("expected 'new' attribute after SET with multi-whitespace")
		} else if s, ok := newAttr.(*types.AttributeValueMemberS); !ok || s.Value != "set-via-whitespace" {
			return fmt.Errorf("expected new='set-via-whitespace', got %v", newAttr)
		}

		if _, hasOld := getResp.Item["old"]; hasOld {
			return fmt.Errorf("expected 'old' attribute removed via #old placeholder, but it still exists")
		}

		return nil
	}))

	return results
}
