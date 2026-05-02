package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (r *TestRunner) runAlarmToSNS(ic *integClients, ts string) TestResult {
	topicName := fmt.Sprintf("integ-alarm-sns-%s", ts)
	queueName := fmt.Sprintf("integ-alarm-sns-q-%s", ts)
	alarmName := fmt.Sprintf("integ-alarm-sns-%s", ts)

	topicARN, err := ic.createTopic(topicName)
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_SNS", func() error { return fmt.Errorf("create topic: %w", err) })
	}
	defer ic.deleteTopic(topicARN)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_SNS", func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	ic.sns.Subscribe(ic.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:%s", queueName)),
	})

	_, err = ic.cw.PutMetricAlarm(ic.ctx, &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		MetricName:         aws.String("CPUUtilization"),
		Namespace:          aws.String("AWS/EC2"),
		Statistic:          cloudwatchtypes.StatisticAverage,
		Period:             aws.Int32(1),
		EvaluationPeriods:  aws.Int32(1),
		Threshold:          aws.Float64(0),
		ComparisonOperator: cloudwatchtypes.ComparisonOperatorGreaterThanThreshold,
		AlarmActions:       []string{topicARN},
	})
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_SNS", func() error { return fmt.Errorf("put alarm: %w", err) })
	}
	defer func() {
		ic.cw.DeleteAlarms(ic.ctx, &cloudwatch.DeleteAlarmsInput{AlarmNames: []string{alarmName}})
	}()

	ic.cw.PutMetricData(ic.ctx, &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("AWS/EC2"),
		MetricData: []cloudwatchtypes.MetricDatum{{
			MetricName: aws.String("CPUUtilization"),
			Value:      aws.Float64(100),
		}},
	})

	return r.pollVerify("CWAlarm_SNS", defaultPollTimeout, func() error {
		alarmResp, err := ic.cw.DescribeAlarms(ic.ctx, &cloudwatch.DescribeAlarmsInput{AlarmNames: []string{alarmName}})
		if err != nil {
			return err
		}
		if len(alarmResp.MetricAlarms) == 0 || alarmResp.MetricAlarms[0].StateValue != cloudwatchtypes.StateValueAlarm {
			return fmt.Errorf("expected alarm state ALARM")
		}
		return ic.verifyMessageContains(queueURL, alarmName)
	})
}

func (r *TestRunner) runAlarmToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-alarm-lambda-%s", ts)
	roleName := fmt.Sprintf("integ-alarm-lambda-role-%s", ts)
	alarmName := fmt.Sprintf("integ-alarm-lambda-%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.createLambda(fnName, roleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fnName)

	_, err := ic.cw.PutMetricAlarm(ic.ctx, &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		MetricName:         aws.String("MemoryUtilization"),
		Namespace:          aws.String("AWS/EC2"),
		Statistic:          cloudwatchtypes.StatisticAverage,
		Period:             aws.Int32(1),
		EvaluationPeriods:  aws.Int32(1),
		Threshold:          aws.Float64(0),
		ComparisonOperator: cloudwatchtypes.ComparisonOperatorGreaterThanThreshold,
		AlarmActions:       []string{fnARN},
	})
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_Lambda", func() error { return fmt.Errorf("put alarm: %w", err) })
	}
	defer func() {
		ic.cw.DeleteAlarms(ic.ctx, &cloudwatch.DeleteAlarmsInput{AlarmNames: []string{alarmName}})
	}()

	ic.cw.PutMetricData(ic.ctx, &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("AWS/EC2"),
		MetricData: []cloudwatchtypes.MetricDatum{{
			MetricName: aws.String("MemoryUtilization"),
			Value:      aws.Float64(100),
		}},
	})

	return r.pollVerify("CWAlarm_Lambda", defaultPollTimeout, func() error {
		alarmResp, err := ic.cw.DescribeAlarms(ic.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe alarms: %w", err)
		}
		if len(alarmResp.MetricAlarms) == 0 {
			return fmt.Errorf("alarm %s not found", alarmName)
		}
		if alarmResp.MetricAlarms[0].StateValue != cloudwatchtypes.StateValueAlarm {
			return fmt.Errorf("expected alarm state ALARM, got %s", alarmResp.MetricAlarms[0].StateValue)
		}
		return ic.verifyLambdaInvoked(fnName)
	})
}

func (r *TestRunner) runAlarmToStepFunctions(ic *integClients, ts string) TestResult {
	roleName := fmt.Sprintf("integ-alarm-sfn-role-%s", ts)
	alarmName := fmt.Sprintf("integ-alarm-sfn-%s", ts)
	smName := fmt.Sprintf("integ-alarm-sfn-sm-%s", ts)

	IAMCreateRole(ic.iam, roleName, sfnTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	_, err := ic.sfn.CreateStateMachine(ic.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(smName),
		RoleArn:    aws.String(intRoleARN(roleName)),
		Definition: aws.String(`{"StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}}`),
	})
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_StepFunctions", func() error { return fmt.Errorf("create state machine: %w", err) })
	}
	defer func() {
		ic.sfn.DeleteStateMachine(ic.ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", smName)),
		})
	}()

	smARN := fmt.Sprintf("arn:aws:states:us-east-1:000000000000:stateMachine:%s", smName)

	_, err = ic.cw.PutMetricAlarm(ic.ctx, &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		MetricName:         aws.String("MemoryUtilization"),
		Namespace:          aws.String("AWS/EC2"),
		Statistic:          cloudwatchtypes.StatisticAverage,
		Period:             aws.Int32(1),
		EvaluationPeriods:  aws.Int32(1),
		Threshold:          aws.Float64(0),
		ComparisonOperator: cloudwatchtypes.ComparisonOperatorGreaterThanThreshold,
		AlarmActions:       []string{smARN},
	})
	if err != nil {
		return r.RunTest(integSvc, "CWAlarm_StepFunctions", func() error { return fmt.Errorf("put alarm: %w", err) })
	}
	defer func() {
		ic.cw.DeleteAlarms(ic.ctx, &cloudwatch.DeleteAlarmsInput{AlarmNames: []string{alarmName}})
	}()

	ic.cw.PutMetricData(ic.ctx, &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("AWS/EC2"),
		MetricData: []cloudwatchtypes.MetricDatum{{
			MetricName: aws.String("DiskSpaceUtilization"),
			Value:      aws.Float64(100),
		}},
	})

	return r.pollVerify("CWAlarm_StepFunctions", defaultPollTimeout, func() error {
		resp, err := ic.sfn.ListExecutions(ic.ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(smARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Executions) == 0 {
			return fmt.Errorf("expected at least 1 execution from alarm, got 0")
		}
		if resp.Executions[0].Status != sfntypes.ExecutionStatusSucceeded {
			return fmt.Errorf("expected SUCCEEDED, got %s", resp.Executions[0].Status)
		}
		return nil
	})
}
