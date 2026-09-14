package testutil

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) runSFNExecutionTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	execSMName := fmt.Sprintf("ExecSM-%d", time.Now().UnixNano())
	_, execRoleARN, execRoleCleanup := tc.createRoleForSM("ExecRole")
	defer execRoleCleanup()

	passDef := `{"Comment":"pass","StartAt":"B","States":{"B":{"Type":"Pass","Result":"hello","End":true}}}`
	var execSMARN string
	results = append(results, r.RunTest("stepfunctions", "StartExecution", func() error {
		resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(execSMName),
			Definition: aws.String(passDef),
			RoleArn:    aws.String(execRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create SM: %v", err)
		}
		execSMARN = *resp.StateMachineArn

		input := map[string]string{"message": "test"}
		inputJSON, _ := json.Marshal(input)
		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(execSMARN),
			Input:           aws.String(string(inputJSON)),
		})
		if err != nil {
			return err
		}
		if startResp.ExecutionArn == nil {
			return fmt.Errorf("execution ARN is nil")
		}
		if startResp.StartDate.IsZero() {
			return fmt.Errorf("start date is zero")
		}
		return nil
	}))

	var executionARN string
	results = append(results, r.RunTest("stepfunctions", "ListExecutions", func() error {
		resp, err := tc.client.ListExecutions(tc.ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(execSMARN),
		})
		if err != nil {
			return err
		}
		if resp.Executions == nil {
			return fmt.Errorf("executions list is nil")
		}
		for _, ex := range resp.Executions {
			if ex.ExecutionArn != nil {
				executionARN = *ex.ExecutionArn
				break
			}
		}
		if executionARN == "" {
			return fmt.Errorf("no execution found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeExecution", func() error {
		if executionARN == "" {
			return fmt.Errorf("no execution ARN available")
		}
		resp, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionARN),
		})
		if err != nil {
			return err
		}
		if resp.Status == "" {
			return fmt.Errorf("execution status is empty")
		}
		if resp.StateMachineArn == nil || *resp.StateMachineArn != execSMARN {
			return fmt.Errorf("state machine ARN mismatch in execution")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "GetExecutionHistory", func() error {
		if executionARN == "" {
			return fmt.Errorf("no execution ARN available")
		}
		resp, err := tc.client.GetExecutionHistory(tc.ctx, &sfn.GetExecutionHistoryInput{
			ExecutionArn: aws.String(executionARN),
		})
		if err != nil {
			return err
		}
		if resp.Events == nil || len(resp.Events) == 0 {
			return fmt.Errorf("events list is nil or empty")
		}
		// Events are numbered sequentially, starting at one (HistoryEvent
		// API reference); a repeated ID means the server overwrote the
		// earlier event instead of appending.
		for i, ev := range resp.Events {
			if ev.Id != int64(i+1) {
				return fmt.Errorf("event %d id = %d, want %d", i, ev.Id, i+1)
			}
		}
		// State events carry the model's generic detail members; a nil
		// StateEnteredEventDetails means the server serialised unmodelled
		// per-state keys the SDK silently drops.
		var sawStateEntered, sawStartDetails bool
		for _, ev := range resp.Events {
			if ev.Type == types.HistoryEventTypePassStateEntered {
				sawStateEntered = true
				d := ev.StateEnteredEventDetails
				if d == nil {
					return fmt.Errorf("PassStateEntered carries no stateEnteredEventDetails")
				}
				if d.Name == nil || *d.Name != "B" {
					return fmt.Errorf("stateEnteredEventDetails name = %v, want B", d.Name)
				}
				if d.Input == nil || *d.Input != `{"message":"test"}` {
					return fmt.Errorf("stateEnteredEventDetails input = %v, want the execution input", d.Input)
				}
			}
			if ev.Type == types.HistoryEventTypeExecutionStarted {
				sawStartDetails = true
				d := ev.ExecutionStartedEventDetails
				if d == nil {
					return fmt.Errorf("ExecutionStarted carries no executionStartedEventDetails")
				}
				if d.RoleArn == nil || *d.RoleArn == "" {
					return fmt.Errorf("executionStartedEventDetails roleArn is empty")
				}
				if d.InputDetails == nil {
					return fmt.Errorf("executionStartedEventDetails inputDetails is nil")
				}
			}
		}
		if !sawStateEntered {
			return fmt.Errorf("no PassStateEntered event in the history")
		}
		if !sawStartDetails {
			return fmt.Errorf("no ExecutionStarted event in the history")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DescribeStateMachineForExecution", func() error {
		if executionARN == "" {
			return fmt.Errorf("no execution ARN available")
		}
		descResp, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionARN),
		})
		if err != nil {
			return fmt.Errorf("describe execution: %v", err)
		}
		if descResp.StateMachineArn == nil {
			return fmt.Errorf("execution has no state machine ARN")
		}
		if descResp.Status != types.ExecutionStatusSucceeded && descResp.Status != types.ExecutionStatusRunning {
			return fmt.Errorf("execution not in suitable state: %s", descResp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Execution_PassStateOutput", func() error {
		if execSMARN == "" {
			return fmt.Errorf("state machine ARN not available")
		}

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(execSMARN),
			Input:           aws.String(`{"value":42}`),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		descResp, err := tc.awaitTerminal(aws.ToString(startResp.ExecutionArn), 500*time.Millisecond, 10)
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if descResp.Status == types.ExecutionStatusSucceeded {
			if descResp.Output == nil {
				return fmt.Errorf("execution output is nil")
			}
			if *descResp.Output != `"hello"` {
				return fmt.Errorf("expected output %q, got %q", `"hello"`, *descResp.Output)
			}
			return nil
		}
		return fmt.Errorf("execution finished with status %s", descResp.Status)
	}))

	results = append(results, r.RunTest("stepfunctions", "StopExecution", func() error {
		longDef := `{"Comment":"long","StartAt":"Wait","States":{"Wait":{"Type":"Wait","Seconds":300,"End":true}}}`
		longSMName := fmt.Sprintf("LongSM-%d", time.Now().UnixNano())
		_, longRoleARN, longRoleCleanup := tc.createRoleForSM("LongRole")
		defer longRoleCleanup()

		longResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(longSMName),
			Definition: aws.String(longDef),
			RoleArn:    aws.String(longRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		longSMARN := *longResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(longSMARN)})

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(longSMARN),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		stopResp, err := tc.client.StopExecution(tc.ctx, &sfn.StopExecutionInput{
			ExecutionArn: aws.String(*startResp.ExecutionArn),
			Error:        aws.String("TestError"),
			Cause:        aws.String("Integration test stop"),
		})
		if err != nil {
			return err
		}
		if stopResp.StopDate.IsZero() {
			return fmt.Errorf("stop date is zero")
		}

		// The caller's error and cause are the execution's reported pair —
		// on the record and on the ExecutionAborted event — rather than
		// generic strings overwriting them through the cancellation race.
		stopped, err := tc.awaitTerminal(*startResp.ExecutionArn, 500*time.Millisecond, 10)
		if err != nil {
			return fmt.Errorf("await terminal: %v", err)
		}
		if stopped.Status != types.ExecutionStatusAborted {
			return fmt.Errorf("stopped execution status = %s", stopped.Status)
		}
		if stopped.Error == nil || *stopped.Error != "TestError" {
			return fmt.Errorf("stopped execution error = %v, want TestError", stopped.Error)
		}
		if stopped.Cause == nil || *stopped.Cause != "Integration test stop" {
			return fmt.Errorf("stopped execution cause = %v, want the caller's cause", stopped.Cause)
		}
		histResp, err := tc.client.GetExecutionHistory(tc.ctx, &sfn.GetExecutionHistoryInput{
			ExecutionArn: startResp.ExecutionArn,
		})
		if err != nil {
			return fmt.Errorf("get history: %v", err)
		}
		for _, ev := range histResp.Events {
			if ev.Type == types.HistoryEventTypeExecutionAborted {
				if ev.ExecutionAbortedEventDetails == nil {
					return fmt.Errorf("ExecutionAborted carries no details")
				}
				d := ev.ExecutionAbortedEventDetails
				if d.Error == nil || *d.Error != "TestError" {
					return fmt.Errorf("ExecutionAborted error = %v, want TestError", d.Error)
				}
				if d.Cause == nil || *d.Cause != "Integration test stop" {
					return fmt.Errorf("ExecutionAborted cause = %v, want the caller's cause", d.Cause)
				}
				return nil
			}
		}
		return fmt.Errorf("no ExecutionAborted event in the stopped execution's history")
	}))

	results = append(results, r.RunTest("stepfunctions", "StartSyncExecution", func() error {
		syncSMName := fmt.Sprintf("SyncSM-%d", time.Now().UnixNano())
		_, syncRoleARN, syncRoleCleanup := tc.createRoleForSM("SyncRole")
		defer syncRoleCleanup()

		// StartSyncExecution is not available for STANDARD workflows; the
		// synchronous start requires an EXPRESS state machine.
		syncResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(syncSMName),
			Definition: aws.String(passDef),
			RoleArn:    aws.String(syncRoleARN),
			Type:       types.StateMachineTypeExpress,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		syncSMARN := *syncResp.StateMachineArn
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(syncSMARN)})

		status, result, err := tc.rawJSONCall("AWSStepFunctions.StartSyncExecution", map[string]interface{}{
			"stateMachineArn": syncSMARN,
			"input":           `{"sync":true}`,
		})
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status %d: %v", status, result)
		}
		if got, _ := result["status"].(string); got != "SUCCEEDED" {
			return fmt.Errorf("expected SUCCEEDED, got %s", got)
		}
		if _, ok := result["executionArn"]; !ok {
			return fmt.Errorf("executionArn missing from response")
		}
		return nil
	}))

	tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{
		StateMachineArn: aws.String(execSMARN),
	})

	results = append(results, r.RunTest("stepfunctions", "Execution_CallbackTaskToken", func() error {
		scfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		sqsClient := sqs.NewFromConfig(scfg)

		queueName := fmt.Sprintf("sfn-callback-%d", time.Now().UnixNano())
		qResp, err := sqsClient.CreateQueue(tc.ctx, &sqs.CreateQueueInput{QueueName: aws.String(queueName)})
		if err != nil {
			return fmt.Errorf("create queue: %w", err)
		}
		defer sqsClient.DeleteQueue(tc.ctx, &sqs.DeleteQueueInput{QueueUrl: qResp.QueueUrl})

		def := fmt.Sprintf(`{"StartAt":"WaitApproval","States":{"WaitApproval":{"Type":"Task","Resource":"arn:aws:states:::sqs:sendMessage.waitForTaskToken","TimeoutSeconds":30,"Parameters":{"QueueUrl":%q,"MessageBody":{"msg":"approve","TaskToken.$":"$$.Task.Token"}},"End":true}}}`, *qResp.QueueUrl)
		smName := fmt.Sprintf("SfnCallback-%d", time.Now().UnixNano())
		smResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(smName),
			Definition: aws.String(def),
			RoleArn:    aws.String(execRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create SM: %w", err)
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smResp.StateMachineArn})

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: smResp.StateMachineArn,
			Input:           aws.String(`{}`),
		})
		if err != nil {
			return err
		}

		// Consume the queue like the external approver would; the message
		// body carries the exact task token the task attempt registered.
		var token string
		deadline := time.Now().Add(15 * time.Second)
		for token == "" && time.Now().Before(deadline) {
			msgs, rerr := sqsClient.ReceiveMessage(tc.ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            qResp.QueueUrl,
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     1,
			})
			if rerr != nil {
				return fmt.Errorf("receive: %w", rerr)
			}
			for _, m := range msgs.Messages {
				var body struct {
					Msg       string `json:"msg"`
					TaskToken string `json:"TaskToken"`
				}
				if jerr := json.Unmarshal([]byte(*m.Body), &body); jerr == nil && body.TaskToken != "" {
					if body.Msg != "approve" {
						return fmt.Errorf("message body lost its parameters: %s", *m.Body)
					}
					token = body.TaskToken
				}
				_, _ = sqsClient.DeleteMessage(tc.ctx, &sqs.DeleteMessageInput{
					QueueUrl:      qResp.QueueUrl,
					ReceiptHandle: m.ReceiptHandle,
				})
			}
		}
		if token == "" {
			return fmt.Errorf("no token-carrying message arrived")
		}

		if _, err := tc.client.SendTaskSuccess(tc.ctx, &sfn.SendTaskSuccessInput{
			TaskToken: aws.String(token),
			Output:    aws.String(`{"approved":true}`),
		}); err != nil {
			return fmt.Errorf("SendTaskSuccess: %w", err)
		}

		desc, err := tc.awaitTerminal(*startResp.ExecutionArn, 500*time.Millisecond, 20)
		if err != nil {
			return fmt.Errorf("await terminal: %w", err)
		}
		if desc.Status != types.ExecutionStatusSucceeded {
			return fmt.Errorf("callback execution status = %s, error %v cause %v", desc.Status, desc.Error, desc.Cause)
		}
		if aws.ToString(desc.Output) != `{"approved":true}` {
			return fmt.Errorf("execution output must be the SendTaskSuccess output, got %s", aws.ToString(desc.Output))
		}

		hist, err := tc.client.GetExecutionHistory(tc.ctx, &sfn.GetExecutionHistoryInput{
			ExecutionArn: startResp.ExecutionArn,
		})
		if err != nil {
			return err
		}
		for _, ev := range hist.Events {
			if ev.Type != types.HistoryEventTypeTaskSubmitted {
				continue
			}
			d := ev.TaskSubmittedEventDetails
			if d == nil {
				return fmt.Errorf("TaskSubmitted carries no details")
			}
			if aws.ToString(d.Resource) != "arn:aws:states:::sqs:sendMessage.waitForTaskToken" {
				return fmt.Errorf("TaskSubmitted resource = %s, want the pattern-suffixed resource", aws.ToString(d.Resource))
			}
			if aws.ToString(d.ResourceType) != "sqs" {
				return fmt.Errorf("TaskSubmitted resourceType = %s, want sqs", aws.ToString(d.ResourceType))
			}
			if aws.ToString(d.Output) == "" {
				return fmt.Errorf("TaskSubmitted output must carry the submit response")
			}
			return nil
		}
		return fmt.Errorf("no TaskSubmitted event in the callback execution's history")
	}))

	results = append(results, r.RunTest("stepfunctions", "Execution_OptimisedLambdaInvoke", func() error {
		lcfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		ltc := &lambdaTestContext{
			r:      r,
			ctx:    tc.ctx,
			client: lambda.NewFromConfig(lcfg),
			cwl:    cloudwatchlogs.NewFromConfig(lcfg),
			iam:    tc.iamClient,
			ts:     fmt.Sprintf("%d", time.Now().UnixNano()),
		}
		fnName, cleanupFn, err := ltc.setupFunction("SfnEcho", "exports.handler = async (event) => { return { echoed: event }; };")
		if err != nil {
			return fmt.Errorf("create function: %w", err)
		}
		defer cleanupFn()
		fnOut, err := ltc.client.GetFunction(tc.ctx, &lambda.GetFunctionInput{FunctionName: aws.String(fnName)})
		if err != nil {
			return fmt.Errorf("get function ARN: %w", err)
		}
		fnARN := aws.ToString(fnOut.Configuration.FunctionArn)

		def := fmt.Sprintf(`{"StartAt":"Invoke","States":{"Invoke":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","Parameters":{"FunctionName":%q,"Payload":{"n":7}},"End":true}}}`, fnARN)
		smName := fmt.Sprintf("SfnLambdaInvoke-%d", time.Now().UnixNano())
		smResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(smName),
			Definition: aws.String(def),
			RoleArn:    aws.String(execRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create SM: %w", err)
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smResp.StateMachineArn})

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: smResp.StateMachineArn,
			Input:           aws.String(`{}`),
		})
		if err != nil {
			return err
		}
		desc, err := tc.awaitTerminal(*startResp.ExecutionArn, 500*time.Millisecond, 20)
		if err != nil {
			return fmt.Errorf("await terminal: %w", err)
		}
		if desc.Status != types.ExecutionStatusSucceeded {
			return fmt.Errorf("invoke execution status = %s, error %v cause %v", desc.Status, desc.Error, desc.Cause)
		}
		var wrapper struct {
			Payload struct {
				Echoed struct {
					N int `json:"n"`
				} `json:"echoed"`
			} `json:"Payload"`
			StatusCode      int    `json:"StatusCode"`
			ExecutedVersion string `json:"ExecutedVersion"`
		}
		if err := json.Unmarshal([]byte(aws.ToString(desc.Output)), &wrapper); err != nil {
			return fmt.Errorf("output is not the invoke response wrapper: %v (%s)", err, aws.ToString(desc.Output))
		}
		if wrapper.Payload.Echoed.N != 7 {
			return fmt.Errorf("wrapper Payload lost the function result: %s", aws.ToString(desc.Output))
		}
		if wrapper.StatusCode != 200 || wrapper.ExecutedVersion != "$LATEST" {
			return fmt.Errorf("wrapper members wrong: %s", aws.ToString(desc.Output))
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Execution_TaskCredentialsRoleAssumption", func() error {
		lcfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		ltc := &lambdaTestContext{
			r:      r,
			ctx:    tc.ctx,
			client: lambda.NewFromConfig(lcfg),
			cwl:    cloudwatchlogs.NewFromConfig(lcfg),
			iam:    tc.iamClient,
			ts:     fmt.Sprintf("%d", time.Now().UnixNano()),
		}
		fnName, cleanupFn, err := ltc.setupFunction("SfnCredEcho", "exports.handler = async (event) => { return { ok: true }; };")
		if err != nil {
			return fmt.Errorf("create function: %w", err)
		}
		defer cleanupFn()
		fnOut, err := ltc.client.GetFunction(tc.ctx, &lambda.GetFunctionInput{FunctionName: aws.String(fnName)})
		if err != nil {
			return fmt.Errorf("get function ARN: %w", err)
		}
		fnARN := aws.ToString(fnOut.Configuration.FunctionArn)

		// The permitted task role trusts the States service principal and
		// carries an inline policy allowing the invocation; the denied role
		// trusts the principal but grants nothing.
		allowedRole := fmt.Sprintf("SfnTaskCredAllowed-%d", time.Now().UnixNano())
		if err := IAMCreateRole(tc.iamClient, allowedRole, sfnTrustPolicy); err != nil {
			return fmt.Errorf("create allowed task role: %w", err)
		}
		defer func() {
			tc.iamClient.DeleteRolePolicy(tc.ctx, &iam.DeleteRolePolicyInput{RoleName: aws.String(allowedRole), PolicyName: aws.String("invoke")})
			IAMDeleteRole(tc.iamClient, allowedRole)
		}()
		allowedARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.AccountID(), allowedRole)
		allowedDoc := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"lambda:InvokeFunction","Resource":%q}]}`, fnARN)
		if _, err := tc.iamClient.PutRolePolicy(tc.ctx, &iam.PutRolePolicyInput{
			RoleName:       aws.String(allowedRole),
			PolicyName:     aws.String("invoke"),
			PolicyDocument: aws.String(allowedDoc),
		}); err != nil {
			return fmt.Errorf("put allowed role policy: %w", err)
		}

		deniedRole := fmt.Sprintf("SfnTaskCredDenied-%d", time.Now().UnixNano())
		if err := IAMCreateRole(tc.iamClient, deniedRole, sfnTrustPolicy); err != nil {
			return fmt.Errorf("create denied task role: %w", err)
		}
		defer IAMDeleteRole(tc.iamClient, deniedRole)
		deniedARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.AccountID(), deniedRole)

		runCredentialed := func(roleArn string) (*sfn.DescribeExecutionOutput, error) {
			def := fmt.Sprintf(`{"StartAt":"Invoke","States":{"Invoke":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","Credentials":{"RoleArn":%q},"Parameters":{"FunctionName":%q,"Payload":{"n":1}},"End":true}}}`, roleArn, fnARN)
			smName := fmt.Sprintf("SfnTaskCred-%d", time.Now().UnixNano())
			smResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
				Name:       aws.String(smName),
				Definition: aws.String(def),
				RoleArn:    aws.String(execRoleARN),
			})
			if err != nil {
				return nil, fmt.Errorf("create SM: %w", err)
			}
			defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smResp.StateMachineArn})
			startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
				StateMachineArn: smResp.StateMachineArn,
				Input:           aws.String(`{}`),
			})
			if err != nil {
				return nil, err
			}
			return tc.awaitTerminal(*startResp.ExecutionArn, 500*time.Millisecond, 20)
		}

		desc, err := runCredentialed(allowedARN)
		if err != nil {
			return fmt.Errorf("await terminal (allowed): %w", err)
		}
		if desc.Status != types.ExecutionStatusSucceeded {
			return fmt.Errorf("permitted credentials must invoke: status %s, error %v cause %v", desc.Status, desc.Error, desc.Cause)
		}

		desc, err = runCredentialed(deniedARN)
		if err != nil {
			return fmt.Errorf("await terminal (denied): %w", err)
		}
		if desc.Status != types.ExecutionStatusFailed {
			return fmt.Errorf("credentials without the invoke permission must fail: status %s, output %s", desc.Status, aws.ToString(desc.Output))
		}
		if aws.ToString(desc.Error) != "States.Permissions" {
			return fmt.Errorf("denied credentials error = %v, want States.Permissions", desc.Error)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Execution_NestedStartExecutionSync", func() error {
		childName := fmt.Sprintf("SfnNestedChild-%d", time.Now().UnixNano())
		childARN, err := tc.createPassSM(childName, "nested child")
		if err != nil {
			return fmt.Errorf("create child SM: %w", err)
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(childARN)})

		def := fmt.Sprintf(`{"StartAt":"StartChild","States":{"StartChild":{"Type":"Task","Resource":"arn:aws:states:::states:startExecution.sync","Parameters":{"StateMachineArn":%q,"Name":"nested-child-exec","Input":{"seed":1}},"End":true}}}`, childARN)
		parentName := fmt.Sprintf("SfnNestedParent-%d", time.Now().UnixNano())
		smResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(parentName),
			Definition: aws.String(def),
			RoleArn:    aws.String(execRoleARN),
		})
		if err != nil {
			return fmt.Errorf("create parent SM: %w", err)
		}
		defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: smResp.StateMachineArn})

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: smResp.StateMachineArn,
			Input:           aws.String(`{}`),
		})
		if err != nil {
			return err
		}
		desc, err := tc.awaitTerminal(*startResp.ExecutionArn, 500*time.Millisecond, 30)
		if err != nil {
			return fmt.Errorf("await terminal: %w", err)
		}
		if desc.Status != types.ExecutionStatusSucceeded {
			return fmt.Errorf("parent status = %s, error %v cause %v", desc.Status, desc.Error, desc.Cause)
		}
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(aws.ToString(desc.Output)), &result); err != nil {
			return fmt.Errorf("sync result is not JSON: %v (%s)", err, aws.ToString(desc.Output))
		}
		if got, _ := result["Output"].(string); got != `{"hello":"world"}` {
			return fmt.Errorf(".sync Output must be the child output string, got %v", result["Output"])
		}
		if got, _ := result["Status"].(string); got != "SUCCEEDED" {
			return fmt.Errorf(".sync Status = %v, want SUCCEEDED", result["Status"])
		}

		childExecs, err := tc.client.ListExecutions(tc.ctx, &sfn.ListExecutionsInput{
			StateMachineArn: aws.String(childARN),
		})
		if err != nil {
			return err
		}
		if len(childExecs.Executions) == 0 {
			return fmt.Errorf("no child execution was created")
		}
		if childExecs.Executions[0].Status != types.ExecutionStatusSucceeded {
			return fmt.Errorf("child execution status = %s", childExecs.Executions[0].Status)
		}
		return nil
	}))

	return results
}
