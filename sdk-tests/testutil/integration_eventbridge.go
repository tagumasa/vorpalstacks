package testutil

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runEventBridgeToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-eb-lambda-%s", ts)
	roleName := fmt.Sprintf("integ-eb-lambda-role-%s", ts)
	busName := fmt.Sprintf("integ-eb-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-eb-rule-%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	ic.createLambda(fnName, roleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fnName)
	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule: aws.String(ruleName), EventBusName: aws.String(busName),
		Targets: []types.Target{{Id: aws.String("t1"), Arn: aws.String(fnARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	detail := map[string]string{"test": "eventbridge-lambda"}
	detailJSON, _ := json.Marshal(detail)
	ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.test"),
			DetailType:   aws.String("IntegrationTest"),
			Detail:       aws.String(string(detailJSON)),
		}},
	})

	return r.pollVerify("EventBridge_Lambda", defaultPollTimeout, func() error {
		return ic.verifyLambdaLogContains(fnName, "eventbridge-lambda")
	})
}

func (r *TestRunner) runEventBridgeToStepFunctions(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-eb-sfn-role-%s", ts)
	busName := fmt.Sprintf("integ-eb-sfn-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-eb-sfn-rule-%s", ts)
	smName := fmt.Sprintf("integ-eb-sfn-sm-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	_, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(`{"StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "EventBridge_StepFunctions", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", smName))})
	}()

	smARN := fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", smName)
	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule: aws.String(ruleName), EventBusName: aws.String(busName),
		Targets: []types.Target{{Id: aws.String("t1"), Arn: aws.String(smARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.test"),
			DetailType:   aws.String("SFNTrigger"),
			Detail:       aws.String(`{"test":"eb-to-sfn"}`),
		}},
	})

	return r.pollVerify("EventBridge_StepFunctions", defaultPollTimeout, func() error {
		resp, err := ic.sfn.ListExecutions(ic.ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Executions) == 0 {
			return fmt.Errorf("expected at least 1 execution, got 0")
		}
		if resp.Executions[0].Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", resp.Executions[0].Status)
		}
		return nil
	})
}

func (r *TestRunner) runEventBridgeToSQS(ic *integClients, ts string) TestResult {
	queueName := fmt.Sprintf("integ-eb-sqs-%s", ts)
	busName := fmt.Sprintf("integ-eb-sqs-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-eb-sqs-rule-%s", ts)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "EventBridge_SQS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	queueARN := fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule: aws.String(ruleName), EventBusName: aws.String(busName),
		Targets: []types.Target{{Id: aws.String("t1"), Arn: aws.String(queueARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.test"),
			DetailType:   aws.String("SQSTest"),
			Detail:       aws.String(`{"message":"eb-to-sqs"}`),
		}},
	})

	return r.pollVerify("EventBridge_SQS", defaultPollTimeout, func() error {
		return ic.verifyMessageContains(queueURL, "eb-to-sqs")
	})
}

func (r *TestRunner) runEventBridgeToSNS(ic *integClients, ts string) TestResult {
	topicName := fmt.Sprintf("integ-eb-sns-%s", ts)
	busName := fmt.Sprintf("integ-eb-sns-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-eb-sns-rule-%s", ts)
	queueName := fmt.Sprintf("integ-eb-sns-sqs-%s", ts)

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "EventBridge_SNS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "EventBridge_SNS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)),
	})

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule: aws.String(ruleName), EventBusName: aws.String(busName),
		Targets: []types.Target{{Id: aws.String("t1"), Arn: aws.String(topicARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.test"),
			DetailType:   aws.String("SNSTest"),
			Detail:       aws.String(`{"message":"eb-to-sns"}`),
		}},
	})

	return r.pollVerify("EventBridge_SNS", defaultPollTimeout, func() error {
		return ic.verifyMessageContains(queueURL, "eb-to-sns")
	})
}

func (r *TestRunner) runEventBridgeToKinesis(ic *integClients, ts string) TestResult {
	streamName := fmt.Sprintf("integ-eb-kinesis-%s", ts)
	busName := fmt.Sprintf("integ-eb-kin-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-eb-kin-rule-%s", ts)

	err := ic.createKinesisStream(streamName)
	if err != nil {
		return r.RunTest(integSvc, "EventBridge_Kinesis", func() error { return fmt.Errorf("create stream: %w", err) })
	}
	defer ic.deleteStream(streamName)

	if err := ic.pollStreamActive(streamName, defaultPollTimeout); err != nil {
		return r.RunTest(integSvc, "EventBridge_Kinesis", func() error { return fmt.Errorf("stream not active: %w", err) })
	}

	streamARN := fmt.Sprintf("arn:aws:kinesis:us-east-1:000000000000:stream/%s", streamName)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule: aws.String(ruleName), EventBusName: aws.String(busName),
		Targets: []types.Target{{Id: aws.String("t1"), Arn: aws.String(streamARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	ic.eb.PutEvents(ic.ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String("com.integration.test"),
			DetailType:   aws.String("KinesisTest"),
			Detail:       aws.String(`{"message":"eb-to-kinesis"}`),
		}},
	})

	return r.pollVerify("EventBridge_Kinesis", defaultPollTimeout, func() error {
		streamDesc, err := ic.describeStream(streamName)
		if err != nil {
			return fmt.Errorf("describe stream: %w", err)
		}
		if len(streamDesc.StreamDescription.Shards) == 0 {
			return fmt.Errorf("no shards in stream")
		}
		shardID := streamDesc.StreamDescription.Shards[0].ShardId
		iter, err := ic.createIteratorFromHorizon(streamName, *shardID)
		if err != nil {
			return fmt.Errorf("create iterator: %w", err)
		}
		records, err := ic.getRecords(iter)
		if err != nil {
			return fmt.Errorf("get records: %w", err)
		}
		if len(records.Records) == 0 {
			return fmt.Errorf("expected records in kinesis stream, got 0")
		}
		data := string(records.Records[0].Data)
		if !strings.Contains(data, "eb-to-kinesis") {
			return fmt.Errorf("record data does not contain expected content, got: %s", data)
		}
		return nil
	})
}
