package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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
// end: a one-time schedule targeting the default event bus with a non-object
// Input fails the PutEvents-style delivery, and with MaximumRetryAttempts=0
// the engine routes the raw Input to the dead-letter queue. A handler that
// answered the drop with an empty result would record the delivery as
// successful and the queue would stay empty.
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
	busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/default", ic.region)
	fireAt := time.Now().UTC().Add(20 * time.Second).Format("2006-01-02T15:04:05")

	_, err = ic.scheduler.CreateSchedule(ic.ctx, &scheduler.CreateScheduleInput{
		Name:                       aws.String(scheduleName),
		GroupName:                  aws.String(groupName),
		ScheduleExpression:         aws.String("at(" + fireAt + ")"),
		ScheduleExpressionTimezone: aws.String("UTC"),
		Target: &schedulertypes.Target{
			Arn:     aws.String(busARN),
			RoleArn: aws.String(intRoleARN(roleName, ic.accountID)),
			// A JSON scalar is not a PutEvents payload: the delivery must
			// fail and the dead-letter route must carry the event.
			Input: aws.String(`"not-an-object"`),
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
			if m.Body != nil && strings.Contains(*m.Body, "not-an-object") {
				return nil
			}
		}
		return fmt.Errorf("dead-letter queue carries no message for the failed delivery (got %d messages)", len(msgs))
	})
}
