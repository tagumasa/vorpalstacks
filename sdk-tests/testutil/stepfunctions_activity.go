package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

func (r *TestRunner) runSFNActivityTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	activityName := fmt.Sprintf("TestActivity-%d", time.Now().UnixNano())
	var activityARN string

	results = append(results, r.RunTest("stepfunctions", "CreateActivity", func() error {
		resp, err := tc.client.CreateActivity(tc.ctx, &sfn.CreateActivityInput{
			Name: aws.String(activityName),
		})
		if err != nil {
			return err
		}
		if resp.ActivityArn == nil || *resp.ActivityArn == "" {
			return fmt.Errorf("activity ARN is nil or empty")
		}
		activityARN = *resp.ActivityArn
		if resp.CreationDate.IsZero() {
			return fmt.Errorf("creation date is zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeActivity", func() error {
		resp, err := tc.client.DescribeActivity(tc.ctx, &sfn.DescribeActivityInput{
			ActivityArn: aws.String(activityARN),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != activityName {
			return fmt.Errorf("activity name mismatch: got %v, want %q", resp.Name, activityName)
		}
		if resp.ActivityArn == nil || *resp.ActivityArn != activityARN {
			return fmt.Errorf("activity ARN mismatch")
		}
		if resp.CreationDate.IsZero() {
			return fmt.Errorf("creation date is zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeActivity_IdempotentRead", func() error {
		resp, err := tc.client.DescribeActivity(tc.ctx, &sfn.DescribeActivityInput{
			ActivityArn: aws.String(activityARN),
		})
		if err != nil {
			return err
		}
		if resp.ActivityArn == nil || *resp.ActivityArn != activityARN {
			return fmt.Errorf("activity ARN mismatch on repeated describe")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ListActivities", func() error {
		resp, err := tc.client.ListActivities(tc.ctx, &sfn.ListActivitiesInput{})
		if err != nil {
			return err
		}
		if resp.Activities == nil {
			return fmt.Errorf("activities list is nil")
		}
		found := false
		for _, a := range resp.Activities {
			if a.ActivityArn != nil && *a.ActivityArn == activityARN {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created activity not found in list")
		}
		return nil
	}))

	taskSMName := fmt.Sprintf("TaskSM-%d", time.Now().UnixNano())
	var taskSMARN string
	results = append(results, r.RunTest("stepfunctions", "GetActivityTask_SendTaskSuccess", func() error {
		taskDef := fmt.Sprintf(`{"Comment":"task","StartAt":"Task","States":{"Task":{"Type":"Task","Resource":"arn:aws:states:%s:%s:activity:%s","End":true}}}`, r.region, r.accountID, activityName)
		_, taskRoleARN, taskRoleCleanup := tc.createRoleForSM("TaskRole")
		defer taskRoleCleanup()

		resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(taskSMName),
			Definition: aws.String(taskDef),
			RoleArn:    aws.String(taskRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create task SM: %v", err)
		}
		taskSMARN = *resp.StateMachineArn

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(taskSMARN),
			Input:           aws.String(`{"task":"input"}`),
		})
		if err != nil {
			return fmt.Errorf("start execution: %v", err)
		}
		execARN := *startResp.ExecutionArn

		taskResp, err := tc.client.GetActivityTask(tc.ctx, &sfn.GetActivityTaskInput{
			ActivityArn: aws.String(activityARN),
			WorkerName:  aws.String("test-worker"),
		})
		if err != nil {
			return fmt.Errorf("get task: %v", err)
		}
		if taskResp.TaskToken == nil || *taskResp.TaskToken == "" {
			return fmt.Errorf("task token is nil or empty")
		}

		_, err = tc.client.SendTaskSuccess(tc.ctx, &sfn.SendTaskSuccessInput{
			TaskToken: taskResp.TaskToken,
			Output:    aws.String(`{"result":"ok"}`),
		})
		if err != nil {
			return fmt.Errorf("send success: %v", err)
		}

		// The history event names must be the model's HistoryEventType
		// members: the activity family is ActivityScheduled/ActivitySucceeded
		// (never the ActivityTask* forms), and the task token must not ride
		// in the scheduled event.
		if _, err := tc.awaitTerminal(execARN, 200*time.Millisecond, 50); err != nil {
			return fmt.Errorf("await terminal: %v", err)
		}
		hist, err := tc.client.GetExecutionHistory(tc.ctx, &sfn.GetExecutionHistoryInput{ExecutionArn: aws.String(execARN)})
		if err != nil {
			return fmt.Errorf("get history: %v", err)
		}
		typesSeen := map[string]bool{}
		for _, ev := range hist.Events {
			typesSeen[string(ev.Type)] = true
		}
		if !typesSeen["ActivityScheduled"] || !typesSeen["ActivitySucceeded"] {
			return fmt.Errorf("history lacks the model activity events, got %v", typesSeen)
		}
		for _, forbidden := range []string{"TaskRetried", "ActivityTaskScheduled", "ActivityTaskStarted", "ActivityTaskSucceeded", "ActivityTaskFailed", "ActivityTaskTimedOut"} {
			if typesSeen[forbidden] {
				return fmt.Errorf("history contains non-model event type %q", forbidden)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "SendTaskFailure", func() error {
		_, taskRoleARN2, taskRole2Cleanup := tc.createRoleForSM("TaskRole2")
		defer taskRole2Cleanup()

		taskSM2Name := fmt.Sprintf("TaskSM2-%d", time.Now().UnixNano())
		failDef := fmt.Sprintf(`{"Comment":"fail task","StartAt":"Task","States":{"Task":{"Type":"Task","Resource":"arn:aws:states:%s:%s:activity:%s","End":true}}}`, r.region, r.accountID, activityName)
		failResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(taskSM2Name),
			Definition: aws.String(failDef),
			RoleArn:    aws.String(taskRoleARN2),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		failSMARN := *failResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(failSMARN)})

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(failSMARN),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}
		failExecARN := *startResp.ExecutionArn

		taskResp, err := tc.client.GetActivityTask(tc.ctx, &sfn.GetActivityTaskInput{
			ActivityArn: aws.String(activityARN),
			WorkerName:  aws.String("fail-worker"),
		})
		if err != nil {
			return fmt.Errorf("get task: %v", err)
		}
		if taskResp.TaskToken == nil || *taskResp.TaskToken == "" {
			return fmt.Errorf("no task token available")
		}

		_, err = tc.client.SendTaskFailure(tc.ctx, &sfn.SendTaskFailureInput{
			TaskToken: taskResp.TaskToken,
			Error:     aws.String("TaskError"),
			Cause:     aws.String("Test failure"),
		})
		if err != nil {
			return err
		}

		if _, err := tc.awaitTerminal(failExecARN, 200*time.Millisecond, 50); err != nil {
			return fmt.Errorf("await terminal: %v", err)
		}
		hist, err := tc.client.GetExecutionHistory(tc.ctx, &sfn.GetExecutionHistoryInput{ExecutionArn: aws.String(failExecARN)})
		if err != nil {
			return fmt.Errorf("get history: %v", err)
		}
		seenActivityFailed := false
		for _, ev := range hist.Events {
			if string(ev.Type) == "ActivityFailed" {
				seenActivityFailed = true
			}
			if strings.HasPrefix(string(ev.Type), "ActivityTask") {
				return fmt.Errorf("history contains non-model event type %q", string(ev.Type))
			}
		}
		if !seenActivityFailed {
			return fmt.Errorf("history lacks ActivityFailed for the failed activity task")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "SendTaskHeartbeat", func() error {
		heartbeatSMName := fmt.Sprintf("HBSM-%d", time.Now().UnixNano())
		hbDef := fmt.Sprintf(`{"Comment":"hb","StartAt":"Task","States":{"Task":{"Type":"Task","Resource":"arn:aws:states:%s:%s:activity:%s","HeartbeatSeconds":60,"End":true}}}`, r.region, r.accountID, activityName)
		_, hbRoleARN, hbRoleCleanup := tc.createRoleForSM("HBRole")
		defer hbRoleCleanup()

		hbResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(heartbeatSMName),
			Definition: aws.String(hbDef),
			RoleArn:    aws.String(hbRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		hbSMARN := *hbResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(hbSMARN)})

		_, err = tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(hbSMARN),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		taskResp, err := tc.client.GetActivityTask(tc.ctx, &sfn.GetActivityTaskInput{
			ActivityArn: aws.String(activityARN),
			WorkerName:  aws.String("hb-worker"),
		})
		if err != nil {
			return fmt.Errorf("get task: %v", err)
		}
		if taskResp.TaskToken == nil || *taskResp.TaskToken == "" {
			return fmt.Errorf("no task token available")
		}

		_, err = tc.client.SendTaskHeartbeat(tc.ctx, &sfn.SendTaskHeartbeatInput{
			TaskToken: taskResp.TaskToken,
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteActivity", func() error {
		resp, err := tc.client.DeleteActivity(tc.ctx, &sfn.DeleteActivityInput{
			ActivityArn: aws.String(activityARN),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	if taskSMARN != "" {
		tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{
			StateMachineArn: aws.String(taskSMARN),
		})
	}

	return results
}
