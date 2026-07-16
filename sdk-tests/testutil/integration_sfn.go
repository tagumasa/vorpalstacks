package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgetypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runSFNTaskLambda(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sfn-lambda-role-%s", ts)
	lambdaRoleName := fmt.Sprintf("integ-sfn-lambda-fn-role-%s", ts)
	smName := fmt.Sprintf("integ-sfn-lambda-%s", ts)
	fnName := fmt.Sprintf("integ-sfn-lambda-fn-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	IAMCreateRole(ic.iam, lambdaRoleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, lambdaRoleName)

	ic.createLambda(fnName, lambdaRoleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fnName)

	definition := fmt.Sprintf(`{
		"StartAt":"InvokeLambda",
		"States":{
			"InvokeLambda":{
				"Type":"Task",
				"Resource":"%s",
				"End":true
			}
		}
	}`, fnARN)

	createResp, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(definition),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_Lambda", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	smARN := createResp.StateMachineArn
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smARN})
	}()

	execResp, err := ic.sfn.StartExecution(ic.ctx, &sfn.StartExecutionInput{
		StateMachineArn: smARN,
		Input:           aws.String(`{"test":"sfn-task-lambda"}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_Lambda", func() error { return fmt.Errorf("start execution: %w", err) })
	}

	return r.pollVerify("SFN_Task_Lambda", defaultPollTimeout, func() error {
		descResp, err := ic.sfn.DescribeExecution(ic.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: execResp.ExecutionArn,
		})
		if err != nil {
			return err
		}
		if descResp.Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", descResp.Status)
		}
		return nil
	})
}

func (r *TestRunner) runSFNTaskSQS(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sfn-sqs-role-%s", ts)
	smName := fmt.Sprintf("integ-sfn-sqs-%s", ts)
	queueName := fmt.Sprintf("integ-sfn-sqs-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SQS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	definition := `{
		"StartAt":"SendMsg",
		"States":{
			"SendMsg":{
				"Type":"Task",
				"Resource":"arn:aws:states:::sqs:sendMessage",
				"Parameters":{
					"QueueUrl":"QUEUE_URL_PLACEHOLDER",
					"MessageBody":{"test":"sfn-to-sqs"}
				},
				"End":true
			}
		}
	}`
	definition = strings.Replace(definition, "QUEUE_URL_PLACEHOLDER", queueURL, 1)

	createResp, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(definition),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SQS", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	smARN := createResp.StateMachineArn
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smARN})
	}()

	startResp, err := ic.sfn.StartExecution(ic.ctx, &sfn.StartExecutionInput{
		StateMachineArn: smARN,
		Input:           aws.String(`{}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SQS", func() error { return fmt.Errorf("start execution: %w", err) })
	}

	return r.pollVerify("SFN_Task_SQS", defaultPollTimeout, func() error {
		descResp, err := ic.sfn.DescribeExecution(ic.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: startResp.ExecutionArn,
		})
		if err != nil {
			return err
		}
		if descResp.Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", descResp.Status)
		}
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return fmt.Errorf("receive messages: %w", err)
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected SQS message from SFN Task, got 0")
		}
		return nil
	})
}

func (r *TestRunner) runSFNTaskSNS(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sfn-sns-role-%s", ts)
	smName := fmt.Sprintf("integ-sfn-sns-%s", ts)
	topicName := fmt.Sprintf("integ-sfn-sns-t-%s", ts)
	queueName := fmt.Sprintf("integ-sfn-sns-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SNS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SNS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)),
	})

	definition := fmt.Sprintf(`{
		"StartAt":"Publish",
		"States":{
			"Publish":{
				"Type":"Task",
				"Resource":"arn:aws:states:::sns:publish",
				"Parameters":{
					"TopicArn":"%s",
					"Message":{"test":"sfn-to-sns"}
				},
				"End":true
			}
		}
	}`, topicARN)

	createResp, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(definition),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SNS", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	smARN := createResp.StateMachineArn
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smARN})
	}()

	startResp, err := ic.sfn.StartExecution(ic.ctx, &sfn.StartExecutionInput{
		StateMachineArn: smARN,
		Input:           aws.String(`{}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_SNS", func() error { return fmt.Errorf("start execution: %w", err) })
	}

	return r.pollVerify("SFN_Task_SNS", defaultPollTimeout, func() error {
		descResp, err := ic.sfn.DescribeExecution(ic.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: startResp.ExecutionArn,
		})
		if err != nil {
			return err
		}
		if descResp.Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", descResp.Status)
		}
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return fmt.Errorf("receive messages: %w", err)
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected SNS message from SFN Task in queue, got 0")
		}
		return nil
	})
}

func (r *TestRunner) runSFNTaskEventBridge(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sfn-eb-role-%s", ts)
	smName := fmt.Sprintf("integ-sfn-eb-%s", ts)
	busName := fmt.Sprintf("integ-sfn-eb-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-sfn-eb-rule-%s", ts)
	queueName := fmt.Sprintf("integ-sfn-eb-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_EventBridge", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	queueARN := fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{
		Name:         aws.String(ruleName),
		EventBusName: aws.String(busName),
	})
	defer ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})

	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule:         aws.String(ruleName),
		EventBusName: aws.String(busName),
		Targets:      []eventbridgetypes.Target{{Id: aws.String("t1"), Arn: aws.String(queueARN)}},
	})
	defer ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"}})

	definition := fmt.Sprintf(`{
		"StartAt":"PutEvent",
		"States":{
			"PutEvent":{
				"Type":"Task",
				"Resource":"arn:aws:states:::events:putEvents",
				"Parameters":{
					"Entries":[{
						"EventBusName":"%s",
						"Source":"com.integration.sfn",
						"DetailType":"SFNEventBridgeTest",
						"Detail":{"test":"sfn-to-eventbridge"}
					}]
				},
				"End":true
			}
		}
	}`, busName)

	createResp, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(definition),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_EventBridge", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	smARN := createResp.StateMachineArn
	defer ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smARN})

	startResp, err := ic.sfn.StartExecution(ic.ctx, &sfn.StartExecutionInput{
		StateMachineArn: smARN,
		Input:           aws.String(`{}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_EventBridge", func() error { return fmt.Errorf("start execution: %w", err) })
	}

	return r.pollVerify("SFN_Task_EventBridge", defaultPollTimeout, func() error {
		descResp, err := ic.sfn.DescribeExecution(ic.ctx, &sfn.DescribeExecutionInput{ExecutionArn: startResp.ExecutionArn})
		if err != nil {
			return err
		}
		if descResp.Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", descResp.Status)
		}
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected EventBridge event from SFN in queue, got 0")
		}
		return nil
	})
}

func (r *TestRunner) runSFNTaskDynamoDB(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sfn-ddb-role-%s", ts)
	smName := fmt.Sprintf("integ-sfn-ddb-%s", ts)
	tableName := fmt.Sprintf("integ-sfn-ddb-%s", strings.ToLower(ts))

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	_, err := ic.dynamodb.CreateTable(ic.ctx, &dynamodb.CreateTableInput{
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
		return r.RunTest(integSvc, "SFN_Task_DynamoDB", func() error { return fmt.Errorf("create table: %w", err) })
	}
	defer ic.dynamodb.DeleteTable(ic.ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

	putDef := fmt.Sprintf(`{
		"StartAt":"PutItem",
		"States":{
			"PutItem":{
				"Type":"Task",
				"Resource":"arn:aws:states:::dynamodb:putItem",
				"Parameters":{
					"TableName":"%s",
					"Item":{
						"pk":{"S":"test-key"},
						"data":{"S":"sfn-dynamodb-integration"}
					}
				},
				"End":true
			}
		}
	}`, tableName)

	createResp, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName + "-put"),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(putDef),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_DynamoDB", func() error { return fmt.Errorf("create put SM: %w", err) })
	}
	putSMARN := createResp.StateMachineArn
	defer ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: putSMARN})

	putExec, err := ic.sfn.StartExecution(ic.ctx, &sfn.StartExecutionInput{
		StateMachineArn: putSMARN,
		Input:           aws.String(`{}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "SFN_Task_DynamoDB", func() error { return fmt.Errorf("start put execution: %w", err) })
	}

	return r.pollVerify("SFN_Task_DynamoDB", defaultPollTimeout, func() error {
		descResp, err := ic.sfn.DescribeExecution(ic.ctx, &sfn.DescribeExecutionInput{ExecutionArn: putExec.ExecutionArn})
		if err != nil {
			return err
		}
		if descResp.Status != sfntypes.ExecutionStatusSucceeded {
			if descResp.Status == sfntypes.ExecutionStatusFailed {
				return fmt.Errorf("put execution failed: %s", aws.ToString(descResp.Error))
			}
			return fmt.Errorf("expected SUCCEEDED, got %s", descResp.Status)
		}

		getResp, err := ic.dynamodb.GetItem(ic.ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]dynamodbtypes.AttributeValue{"pk": &dynamodbtypes.AttributeValueMemberS{Value: "test-key"}},
		})
		if err != nil {
			return fmt.Errorf("GetItem verify: %w", err)
		}
		if len(getResp.Item) == 0 {
			return fmt.Errorf("item not found in DynamoDB after putItem task")
		}
		if pk, ok := getResp.Item["pk"].(*dynamodbtypes.AttributeValueMemberS); !ok || pk.Value != "test-key" {
			return fmt.Errorf("pk mismatch: expected test-key")
		}
		if data, ok := getResp.Item["data"].(*dynamodbtypes.AttributeValueMemberS); !ok || data.Value != "sfn-dynamodb-integration" {
			return fmt.Errorf("data mismatch: expected sfn-dynamodb-integration")
		}
		return nil
	})
}
