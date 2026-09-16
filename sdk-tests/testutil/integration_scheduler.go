package testutil

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	ebtypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulertypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runSchedulerToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-sched-lambda-%s", ts)
	roleName := fmt.Sprintf("integ-sched-role-%s", ts)
	lambdaRoleName := fmt.Sprintf("integ-sched-lambda-fn-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-lambda-%s", ts)
	groupName := fmt.Sprintf("integ-sched-group-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	IAMCreateRole(ic.iam, lambdaRoleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, lambdaRoleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	ic.createLambda(fnName, lambdaRoleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:%s", ic.region, fnName)

	_, err := ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:               aws.String(scheduleName),
		GroupName:          aws.String(groupName),
		ScheduleExpression: aws.String("rate(1 minute)"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(fnARN),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Lambda", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(scheduleName),
			GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_Lambda", schedulerPollTimeout, func() error {
		return ic.verifyLambdaInvoked(fnName)
	})
}

func (r *TestRunner) runSchedulerToSQS(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sched-sqs-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-sqs-%s", ts)
	groupName := fmt.Sprintf("integ-sched-sqs-group-%s", ts)
	queueName := fmt.Sprintf("integ-sched-sqs-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_SQS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	queueARN := fmt.Sprintf("arn:aws:sqs:%s:000000000000:%s", ic.region, queueName)

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:               aws.String(scheduleName),
		GroupName:          aws.String(groupName),
		ScheduleExpression: aws.String("rate(1 minute)"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(queueARN),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_SQS", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(scheduleName),
			GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_SQS", schedulerPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected message from scheduler in queue, got 0")
		}
		return nil
	})
}

func (r *TestRunner) runSchedulerToSNS(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sched-sns-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-sns-%s", ts)
	groupName := fmt.Sprintf("integ-sched-sns-group-%s", ts)
	topicName := fmt.Sprintf("integ-sched-sns-t-%s", ts)
	queueName := fmt.Sprintf("integ-sched-sns-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_SNS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_SNS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:%s:000000000000:%s", ic.region, queueName)),
	})

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:               aws.String(scheduleName),
		GroupName:          aws.String(groupName),
		ScheduleExpression: aws.String("rate(1 minute)"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(topicARN),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_SNS", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(scheduleName),
			GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_SNS", schedulerPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("expected message from scheduler (via SNS) in queue, got 0")
		}
		return nil
	})
}

// runSchedulerToKinesis pins the templated Kinesis delivery contract end to
// end: a schedule targeting a Kinesis stream delivers its Input, and an SDK
// consumer reads the record back through GetRecords — the client base64-
// decodes the Data member of the response, so a deliverer that stored the
// raw JSON payload fails the read-back with a client-side decode error.
func (r *TestRunner) runSchedulerToKinesis(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sched-kin-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-kin-%s", ts)
	groupName := fmt.Sprintf("integ-sched-kin-group-%s", ts)
	streamName := fmt.Sprintf("integ-sched-kin-stream-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	if err := ic.createKinesisStream(streamName); err != nil {
		return r.RunTest(integSvc, "Scheduler_Kinesis", func() error { return fmt.Errorf("create stream: %w", err) })
	}
	defer ic.deleteStream(streamName)

	if err := ic.pollStreamActive(streamName, schedulerPollTimeout); err != nil {
		return r.RunTest(integSvc, "Scheduler_Kinesis", func() error { return fmt.Errorf("stream not active: %w", err) })
	}

	streamARN := fmt.Sprintf("arn:aws:kinesis:%s:000000000000:stream/%s", ic.region, streamName)

	_, err := ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String(r.universalFireAt()),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(streamARN),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			Input:   aws.String(`{"kinesisPin":"templated"}`),
			KinesisParameters: &schedulertypes.KinesisParameters{
				PartitionKey: aws.String("templated"),
			},
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Kinesis", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName), GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_Kinesis", schedulerPollTimeout, func() error {
		streamDesc, err := ic.describeStream(streamName)
		if err != nil {
			return fmt.Errorf("describe stream: %w", err)
		}
		if len(streamDesc.StreamDescription.Shards) == 0 {
			return fmt.Errorf("no shards in stream")
		}
		iter, err := ic.createIteratorFromHorizon(streamName, *streamDesc.StreamDescription.Shards[0].ShardId)
		if err != nil {
			return fmt.Errorf("create iterator: %w", err)
		}
		records, err := ic.getRecords(iter)
		if err != nil {
			return fmt.Errorf("get records: %w", err)
		}
		if len(records.Records) == 0 {
			return fmt.Errorf("expected a record from the templated delivery, got 0")
		}
		if data := string(records.Records[0].Data); !strings.Contains(data, `"kinesisPin":"templated"`) {
			return fmt.Errorf("record data does not carry the decoded schedule input, got: %s", data)
		}
		return nil
	})
}

func (r *TestRunner) runSchedulerToStepFunctions(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sched-sfn-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-sfn-%s", ts)
	groupName := fmt.Sprintf("integ-sched-sfn-group-%s", ts)
	smName := fmt.Sprintf("integ-sched-sfn-sm-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	IAMCreateRole(ic.iam, fmt.Sprintf("%s-sfn", roleName), sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, fmt.Sprintf("%s-sfn", roleName))

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	_, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(fmt.Sprintf("%s-sfn", roleName), ic.accountID)),
		Definition: aws.String(`{"StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_StepFunctions", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:%s:000000000000:stateMachine:%s", ic.region, smName)),
		})
	}()

	smARN := fmt.Sprintf("arn:aws:states:%s:000000000000:stateMachine:%s", ic.region, smName)

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:               aws.String(scheduleName),
		GroupName:          aws.String(groupName),
		ScheduleExpression: aws.String("rate(1 minute)"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(smARN),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_StepFunctions", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(scheduleName),
			GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_StepFunctions", schedulerPollTimeout, func() error {
		resp, err := ic.sfn.ListExecutions(ic.ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Executions) == 0 {
			return fmt.Errorf("expected at least 1 execution from scheduler, got 0")
		}
		if resp.Executions[0].Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", resp.Executions[0].Status)
		}
		return nil
	})
}

// runSchedulerToEventBridgeDLQ pins the delivery-failure contract end to
// end: a one-time schedule targeting an event bus that does not exist fails
// the PutEvents-style delivery (the ingress plane reports a missing bus as
// a delivery failure), and with MaximumRetryAttempts=0 the engine routes
// the raw Input to the dead-letter queue. A handler that answered the drop
// with an empty result would record the delivery as successful and the
// queue would stay empty. The failure driver is the missing bus, not a
// malformed payload: Target.Input on an EventBridge target must be a valid
// JSON object, so a malformed payload is rejected at creation and can no
// longer reach the delivery path.
func (r *TestRunner) runSchedulerToEventBridgeDLQ(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sched-eb-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-eb-dlq-%s", ts)
	groupName := fmt.Sprintf("integ-sched-eb-group-%s", ts)
	queueName := fmt.Sprintf("integ-sched-eb-dlq-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_EventBridge_DLQ", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	queueARN := fmt.Sprintf("arn:aws:sqs:%s:000000000000:%s", ic.region, queueName)
	// A bus that was never created: the schedule is creatable (target
	// validation is structural), but every delivery to it fails.
	busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/missing-%s", ic.region, ts)
	fireAt := time.Now().UTC().Add(20 * time.Second).Format("2006-01-02T15:04:05")

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String("at(" + fireAt + ")"),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(busARN),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			// The object Input rides to the dead-letter queue verbatim
			// when the delivery to the missing bus exhausts its retry
			// policy.
			Input: aws.String(`{"dlqPin":true}`),
			EventBridgeParameters: &schedulertypes.EventBridgeParameters{
				// The 'aws.' prefix is reserved and rejected at creation;
				// a customer source keeps the schedule creatable.
				Source:     aws.String("vorpal.test"),
				DetailType: aws.String("dlq-pin"),
			},
			RetryPolicy: &schedulertypes.RetryPolicy{
				MaximumRetryAttempts:     aws.Int32(0),
				MaximumEventAgeInSeconds: aws.Int32(60),
			},
			DeadLetterConfig: &schedulertypes.DeadLetterConfig{Arn: aws.String(queueARN)},
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_EventBridge_DLQ", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(scheduleName),
			GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_EventBridge_DLQ", schedulerPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		for _, m := range msgs {
			if m.Body != nil && strings.Contains(*m.Body, "dlqPin") {
				return nil
			}
		}
		return fmt.Errorf("dead-letter queue carries no message for the failed delivery (got %d messages)", len(msgs))
	})
}

// runSchedulerToEventBridgeCustomBus pins the named-bus delivery contract
// end to end: a schedule targeting a CUSTOM event bus delivers through
// that bus — the rule bound to the bus fires and its SQS target receives
// the event. A delivery that fell back to the default bus leaves the
// queue empty: the rule exists only on the custom bus.
func (r *TestRunner) runSchedulerToEventBridgeCustomBus(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-sched-ebcb-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-sched-ebcb-%s", ts)
	groupName := fmt.Sprintf("integ-sched-ebcb-group-%s", ts)
	queueName := fmt.Sprintf("integ-sched-ebcb-q-%s", ts)
	busName := fmt.Sprintf("integ-sched-ebcb-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-sched-ebcb-rule-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_EventBridge_CustomBus", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{
		Name:         aws.String(ruleName),
		EventBusName: aws.String(busName),
		EventPattern: aws.String(`{"source":["vorpal.integration.test"]}`),
	})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	queueARN := fmt.Sprintf("arn:aws:sqs:%s:000000000000:%s", ic.region, queueName)
	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule:         aws.String(ruleName),
		EventBusName: aws.String(busName),
		Targets:      []ebtypes.Target{{Id: aws.String("t1"), Arn: aws.String(queueARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule:         aws.String(ruleName),
			EventBusName: aws.String(busName),
			Ids:          []string{"t1"},
		})
	}()

	busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", ic.region, busName)
	fireAt := time.Now().UTC().Add(20 * time.Second).Format("2006-01-02T15:04:05")

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String("at(" + fireAt + ")"),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(busARN),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			Input:   aws.String(`{"customBus":true}`),
			EventBridgeParameters: &schedulertypes.EventBridgeParameters{
				Source:     aws.String("vorpal.integration.test"),
				DetailType: aws.String("custom-bus-pin"),
			},
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_EventBridge_CustomBus", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(scheduleName),
			GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_EventBridge_CustomBus", schedulerPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		for _, m := range msgs {
			if m.Body != nil && strings.Contains(*m.Body, `"customBus":true`) {
				return nil
			}
		}
		return fmt.Errorf("queue carries no event delivered through the custom bus (got %d messages)", len(msgs))
	})
}

// Universal-target integration pins: a schedule whose Target.Arn carries
// the aws-sdk form (arn:aws:scheduler:::aws-sdk:{service}:{action})
// delivers through the family invokers with Target.Input as the target
// API's request JSON, and an ARN naming an unimplemented service fails at
// delivery into the dead-letter queue with the raw Input riding along.

func (r *TestRunner) universalFireAt() string {
	return "at(" + time.Now().UTC().Add(20*time.Second).Format("2006-01-02T15:04:05") + ")"
}

func (r *TestRunner) runSchedulerUniversalToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-univ-lambda-%s", ts)
	roleName := fmt.Sprintf("integ-univ-lambda-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-univ-lambda-%s", ts)
	groupName := fmt.Sprintf("integ-univ-lambda-group-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	if _, err := ic.createLambda(fnName, roleName); err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_Lambda", func() error { return fmt.Errorf("create lambda: %w", err) })
	}
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:%s", ic.region, fnName)

	_, err := ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String(r.universalFireAt()),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			Arn:     aws.String("arn:aws:scheduler:::aws-sdk:lambda:invoke"),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			Input:   aws.String(fmt.Sprintf(`{"FunctionName":%q,"InvocationType":"Event","Payload":"{\"universal\":true}"}`, fnARN)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_Lambda", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName), GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_Universal_Lambda", schedulerPollTimeout, func() error {
		return ic.verifyLambdaInvoked(fnName)
	})
}

func (r *TestRunner) runSchedulerUniversalToSQS(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-univ-sqs-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-univ-sqs-%s", ts)
	groupName := fmt.Sprintf("integ-univ-sqs-group-%s", ts)
	queueName := fmt.Sprintf("integ-univ-sqs-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_SQS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String(r.universalFireAt()),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			Arn:     aws.String("arn:aws:scheduler:::aws-sdk:sqs:sendMessage"),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			Input:   aws.String(fmt.Sprintf(`{"QueueUrl":%q,"MessageBody":"universal-sqs-pin"}`, queueURL)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_SQS", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName), GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_Universal_SQS", schedulerPollTimeout, func() error {
		return ic.verifyMessageContains(queueURL, "universal-sqs-pin")
	})
}

func (r *TestRunner) runSchedulerUniversalToSNS(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-univ-sns-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-univ-sns-%s", ts)
	groupName := fmt.Sprintf("integ-univ-sns-group-%s", ts)
	topicName := fmt.Sprintf("integ-univ-sns-t-%s", ts)
	queueName := fmt.Sprintf("integ-univ-sns-q-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_SNS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_SNS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:%s:000000000000:%s", ic.region, queueName)),
	})

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String(r.universalFireAt()),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			Arn:     aws.String("arn:aws:scheduler:::aws-sdk:sns:publish"),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			Input:   aws.String(fmt.Sprintf(`{"TopicArn":%q,"Message":"universal-sns-pin","Subject":"univ"}`, topicARN)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_SNS", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName), GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_Universal_SNS", schedulerPollTimeout, func() error {
		return ic.verifyMessageContains(queueURL, "universal-sns-pin")
	})
}

func (r *TestRunner) runSchedulerUniversalToKinesis(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-univ-kin-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-univ-kin-%s", ts)
	groupName := fmt.Sprintf("integ-univ-kin-group-%s", ts)
	streamName := fmt.Sprintf("integ-univ-kin-stream-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	if err := ic.createKinesisStream(streamName); err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_Kinesis", func() error { return fmt.Errorf("create stream: %w", err) })
	}
	defer ic.deleteStream(streamName)

	if err := ic.pollStreamActive(streamName, schedulerPollTimeout); err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_Kinesis", func() error { return fmt.Errorf("stream not active: %w", err) })
	}

	streamARN := fmt.Sprintf("arn:aws:kinesis:%s:000000000000:stream/%s", ic.region, streamName)

	_, err := ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String(r.universalFireAt()),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			Arn:     aws.String("arn:aws:scheduler:::aws-sdk:kinesis:putRecord"),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			Input: aws.String(fmt.Sprintf(`{"StreamARN":%q,"Data":%q,"PartitionKey":"univ"}`,
				streamARN, base64.StdEncoding.EncodeToString([]byte("universal-kinesis-pin")))),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_Kinesis", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName), GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_Universal_Kinesis", schedulerPollTimeout, func() error {
		streamDesc, err := ic.describeStream(streamName)
		if err != nil {
			return fmt.Errorf("describe stream: %w", err)
		}
		if len(streamDesc.StreamDescription.Shards) == 0 {
			return fmt.Errorf("no shards in stream")
		}
		iter, err := ic.createIteratorFromHorizon(streamName, *streamDesc.StreamDescription.Shards[0].ShardId)
		if err != nil {
			return fmt.Errorf("create iterator: %w", err)
		}
		records, err := ic.getRecords(iter)
		if err != nil {
			return fmt.Errorf("get records: %w", err)
		}
		if len(records.Records) == 0 {
			return fmt.Errorf("expected a record from the universal putRecord, got 0")
		}
		if data := string(records.Records[0].Data); !strings.Contains(data, "universal-kinesis-pin") {
			return fmt.Errorf("record data does not contain the decoded pin, got: %s", data)
		}
		return nil
	})
}

func (r *TestRunner) runSchedulerUniversalToStepFunctions(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-univ-sfn-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-univ-sfn-%s", ts)
	groupName := fmt.Sprintf("integ-univ-sfn-group-%s", ts)
	smName := fmt.Sprintf("integ-univ-sfn-sm-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	IAMCreateRole(ic.iam, fmt.Sprintf("%s-sfn", roleName), sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, fmt.Sprintf("%s-sfn", roleName))

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	if _, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(fmt.Sprintf("%s-sfn", roleName), ic.accountID)),
		Definition: aws.String(`{"StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}}`),
	}); err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_StepFunctions", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:%s:000000000000:stateMachine:%s", ic.region, smName)),
		})
	}()

	smARN := fmt.Sprintf("arn:aws:states:%s:000000000000:stateMachine:%s", ic.region, smName)

	_, err := ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String(r.universalFireAt()),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			Arn:     aws.String("arn:aws:scheduler:::aws-sdk:sfn:startExecution"),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			Input:   aws.String(fmt.Sprintf(`{"StateMachineArn":%q,"Input":"{\"universal\":true}"}`, smARN)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_StepFunctions", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName), GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_Universal_StepFunctions", schedulerPollTimeout, func() error {
		resp, err := ic.sfn.ListExecutions(ic.ctx, &sfn.ListExecutionsInput{StateMachineArn: aws.String(smARN)})
		if err != nil {
			return err
		}
		if len(resp.Executions) == 0 {
			return fmt.Errorf("expected at least 1 execution from the universal target, got 0")
		}
		if resp.Executions[0].Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", resp.Executions[0].Status)
		}
		return nil
	})
}

func (r *TestRunner) runSchedulerUniversalToEventBridge(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-univ-eb-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-univ-eb-%s", ts)
	groupName := fmt.Sprintf("integ-univ-eb-group-%s", ts)
	queueName := fmt.Sprintf("integ-univ-eb-q-%s", ts)
	busName := fmt.Sprintf("integ-univ-eb-bus-%s", ts)
	ruleName := fmt.Sprintf("integ-univ-eb-rule-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_EventBridge", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.eb.CreateEventBus(ic.ctx, &eventbridge.CreateEventBusInput{Name: aws.String(busName)})
	defer func() {
		ic.eb.DeleteEventBus(ic.ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(busName)})
	}()

	ic.eb.PutRule(ic.ctx, &eventbridge.PutRuleInput{
		Name:         aws.String(ruleName),
		EventBusName: aws.String(busName),
		EventPattern: aws.String(`{"source":["vorpal.universal.test"]}`),
	})
	defer func() {
		ic.eb.DeleteRule(ic.ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
	}()

	queueARN := fmt.Sprintf("arn:aws:sqs:%s:000000000000:%s", ic.region, queueName)
	ic.eb.PutTargets(ic.ctx, &eventbridge.PutTargetsInput{
		Rule:         aws.String(ruleName),
		EventBusName: aws.String(busName),
		Targets:      []ebtypes.Target{{Id: aws.String("t1"), Arn: aws.String(queueARN)}},
	})
	defer func() {
		ic.eb.RemoveTargets(ic.ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"t1"},
		})
	}()

	busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", ic.region, busName)

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String(r.universalFireAt()),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			Arn:     aws.String("arn:aws:scheduler:::aws-sdk:eventbridge:putEvents"),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			Input:   aws.String(fmt.Sprintf(`{"Entries":[{"Source":"vorpal.universal.test","DetailType":"universal-eb-pin","Detail":"{\"universal\":true}","EventBusName":%q}]}`, busARN)),
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_EventBridge", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName), GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_Universal_EventBridge", schedulerPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		for _, m := range msgs {
			if m.Body != nil && strings.Contains(*m.Body, "universal-eb-pin") && strings.Contains(*m.Body, `"universal":true`) {
				return nil
			}
		}
		return fmt.Errorf("queue carries no event put through the universal target (got %d messages)", len(msgs))
	})
}

// runSchedulerUniversalAcceptThenFail pins the missing-substrate posture:
// a universal ARN naming a service this platform does not implement is
// creatable (the AWS contract) and fails at delivery, routing the raw
// Input to the dead-letter queue exactly as the templated stubs do.
func (r *TestRunner) runSchedulerUniversalAcceptThenFail(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-univ-af-role-%s", ts)
	scheduleName := fmt.Sprintf("integ-univ-af-%s", ts)
	groupName := fmt.Sprintf("integ-univ-af-group-%s", ts)
	queueName := fmt.Sprintf("integ-univ-af-dlq-%s", ts)

	IAMCreateRole(ic.iam, roleName, schedulerTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.scheduler.CreateScheduleGroup(ic.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(groupName)})
	defer func() {
		ic.scheduler.DeleteScheduleGroup(ic.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})
	}()

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_AcceptFail", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	queueARN := fmt.Sprintf("arn:aws:sqs:%s:000000000000:%s", ic.region, queueName)

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String(r.universalFireAt()),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			// AWS Batch is not implemented on this platform: the schedule
			// is creatable, the delivery fails naming the missing
			// substrate, and the dead-letter queue receives the Input.
			Arn:     aws.String("arn:aws:scheduler:::aws-sdk:batch:submitJob"),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			Input:   aws.String(`{"JobName":"universal-accept-fail-pin"}`),
			RetryPolicy: &schedulertypes.RetryPolicy{
				MaximumRetryAttempts:     aws.Int32(0),
				MaximumEventAgeInSeconds: aws.Int32(60),
			},
			DeadLetterConfig: &schedulertypes.DeadLetterConfig{Arn: aws.String(queueARN)},
		},
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{Mode: schedulertypes.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return r.RunTest(integSvc, "Scheduler_Universal_AcceptFail", func() error { return fmt.Errorf("create schedule: %w", err) })
	}
	defer func() {
		ic.scheduler.DeleteSchedule(ic.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName), GroupName: aws.String(groupName),
		})
	}()

	return r.pollVerify("Scheduler_Universal_AcceptFail", schedulerPollTimeout, func() error {
		msgs, err := ic.receiveMessages(queueURL, 5, 3)
		if err != nil {
			return err
		}
		for _, m := range msgs {
			if m.Body != nil && strings.Contains(*m.Body, "universal-accept-fail-pin") {
				return nil
			}
		}
		return fmt.Errorf("dead-letter queue carries no message for the undeliverable universal target (got %d messages)", len(msgs))
	})
}
