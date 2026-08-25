package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

// runSFNWaitTests pins the Wait-state timestamp semantics: a past
// timestamp completes immediately, a future timestamp is honoured, an
// invalid Timestamp literal is rejected at creation, an invalid
// TimestampPath value fails the execution with States.Runtime, and a
// JSONata timestamp expression is honoured.
func (r *TestRunner) runSFNWaitTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	waitSM := func(name, definition string) (string, func()) {
		arn, cleanup, err := tc.createRoleBackedSM(name, definition)
		if err != nil {
			return "", func() {}
		}
		return arn, cleanup
	}

	awaitExecution := func(executionARN string, atLeast time.Duration) (types.ExecutionStatus, string, error) {
		started := time.Now()
		desc, err := tc.awaitTerminal(executionARN, 300*time.Millisecond, 200)
		if err != nil {
			return "", "", err
		}
		if atLeast > 0 && time.Since(started) < atLeast {
			return desc.Status, aws.ToString(desc.Error), fmt.Errorf("execution finished after %v, before the %v wait elapsed", time.Since(started).Round(time.Millisecond), atLeast)
		}
		return desc.Status, aws.ToString(desc.Error), nil
	}

	results = append(results, r.RunTest("stepfunctions", "Wait_PastTimestamp_CompletesImmediately", func() error {
		def := `{"StartAt":"W","States":{"W":{"Type":"Wait","Timestamp":"2024-03-14T01:59:00Z","Next":"P"},"P":{"Type":"Pass","Result":"done","End":true}}}`
		arn, cleanup := waitSM("WaitPast", def)
		defer cleanup()

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{StateMachineArn: aws.String(arn)})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}
		status, execErr, err := awaitExecution(*startResp.ExecutionArn, 0)
		if err != nil {
			return err
		}
		if status != types.ExecutionStatusSucceeded {
			return fmt.Errorf("status = %s, error = %q, want SUCCEEDED", status, execErr)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Wait_FutureTimestamp_IsHonoured", func() error {
		target := time.Now().UTC().Add(3 * time.Second).Format("2006-01-02T15:04:05Z")
		def := fmt.Sprintf(`{"StartAt":"W","States":{"W":{"Type":"Wait","Timestamp":%q,"Next":"P"},"P":{"Type":"Pass","Result":"done","End":true}}}`, target)
		arn, cleanup := waitSM("WaitFuture", def)
		defer cleanup()

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{StateMachineArn: aws.String(arn)})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}
		status, execErr, err := awaitExecution(*startResp.ExecutionArn, 2*time.Second)
		if err != nil {
			return err
		}
		if status != types.ExecutionStatusSucceeded {
			return fmt.Errorf("status = %s, error = %q, want SUCCEEDED", status, execErr)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "CreateStateMachine_InvalidWaitTimestamp_Rejected", func() error {
		def := `{"StartAt":"W","States":{"W":{"Type":"Wait","Timestamp":"2024-03-14 01:59:00","Next":"P"},"P":{"Type":"Pass","End":true}}}`
		_, roleARN, roleCleanup := tc.createRoleForSM("WaitInvalidRole")
		defer roleCleanup()
		_, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
			Name:       aws.String(fmt.Sprintf("WaitInvalid-%d", time.Now().UnixNano())),
			Definition: aws.String(def),
			RoleArn:    aws.String(roleARN),
		})
		if err == nil {
			return fmt.Errorf("definition with an offset-less Timestamp accepted")
		}
		return expectAWSErrorCode(err, "InvalidDefinition")
	}))

	results = append(results, r.RunTest("stepfunctions", "Wait_TimestampPath_InvalidValue_FailsExecution", func() error {
		def := `{"StartAt":"W","States":{"W":{"Type":"Wait","TimestampPath":"$.expiry","Next":"P"},"P":{"Type":"Pass","End":true}}}`
		arn, cleanup := waitSM("WaitPath", def)
		defer cleanup()

		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(arn),
			Input:           aws.String(`{"expiry":"not-a-timestamp"}`),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}
		status, execErr, err := awaitExecution(*startResp.ExecutionArn, 0)
		if err != nil {
			return err
		}
		if status != types.ExecutionStatusFailed {
			return fmt.Errorf("status = %s, want FAILED", status)
		}
		if execErr != "States.Runtime" {
			return fmt.Errorf("execution error = %q, want States.Runtime", execErr)
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "Wait_JSONata_TimestampExpression_IsHonoured", func() error {
		def := `{"QueryLanguage":"JSONata","StartAt":"W","States":{"W":{"Type":"Wait","Timestamp":"{% $states.input.expiry %}","Next":"P"},"P":{"Type":"Pass","QueryLanguage":"JSONata","End":true}}}`
		arn, cleanup := waitSM("WaitJSONata", def)
		defer cleanup()

		target := time.Now().UTC().Add(3 * time.Second).Format("2006-01-02T15:04:05Z")
		startResp, err := tc.client.StartExecution(tc.ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(arn),
			Input:           aws.String(fmt.Sprintf(`{"expiry":%q}`, target)),
		})
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}
		status, execErr, err := awaitExecution(*startResp.ExecutionArn, 2*time.Second)
		if err != nil {
			return err
		}
		if status != types.ExecutionStatusSucceeded {
			return fmt.Errorf("status = %s, error = %q, want SUCCEEDED", status, execErr)
		}
		return nil
	}))

	return results
}
