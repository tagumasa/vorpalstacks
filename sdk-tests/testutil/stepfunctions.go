package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"vorpalstacks-sdk-tests/config"
)

type sfnTestContext struct {
	client      *sfn.Client
	iamClient   *iam.Client
	ctx         context.Context
	runner      *TestRunner
	roleARN     string
	roleName    string
	trustPolicy string
	definition  string
}

func newSFNTestContext(r *TestRunner) (*sfnTestContext, func(), error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("load config: %w", err)
	}

	tc := &sfnTestContext{
		client:      sfn.NewFromConfig(cfg),
		iamClient:   iam.NewFromConfig(cfg),
		ctx:         context.Background(),
		runner:      r,
		trustPolicy: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"states.amazonaws.com"},"Action":"sts:AssumeRole"}]}`,
		definition:  `{"Comment":"A Hello World example","StartAt":"HelloWorld","States":{"HelloWorld":{"Type":"Pass","End":true}}}`,
	}

	tc.roleName = fmt.Sprintf("TestSfnRole-%d", time.Now().UnixNano())
	tc.roleARN = fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, tc.roleName)

	if err := IAMCreateRole(tc.iamClient, tc.roleName, tc.trustPolicy); err != nil {
		return nil, nil, fmt.Errorf("create IAM role: %w", err)
	}

	cleanup := func() {
		IAMDeleteRoleCtx(tc.ctx, tc.iamClient, tc.roleName)
	}

	return tc, cleanup, nil
}

func (tc *sfnTestContext) createTestSM(name string) (string, error) {
	resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(name),
		Definition: aws.String(tc.definition),
		RoleArn:    aws.String(tc.roleARN),
	})
	if err != nil {
		return "", err
	}
	return *resp.StateMachineArn, nil
}

func (tc *sfnTestContext) createPassSM(name, comment string) (string, error) {
	def := map[string]interface{}{
		"Comment": comment,
		"StartAt": "Pass",
		"States": map[string]interface{}{
			"Pass": map[string]interface{}{
				"Type":   "Pass",
				"Result": map[string]string{"hello": "world"},
				"End":    true,
			},
		},
	}
	defJSON, err := json.Marshal(def)
	if err != nil {
		return "", err
	}
	resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(name),
		Definition: aws.String(string(defJSON)),
		RoleArn:    aws.String(tc.roleARN),
	})
	if err != nil {
		return "", err
	}
	return *resp.StateMachineArn, nil
}

func (tc *sfnTestContext) createRoleForSM(name string) (string, string, func()) {
	roleName := fmt.Sprintf("%s-%d", name, time.Now().UnixNano())
	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.runner.AccountID(), roleName)
	IAMCreateRole(tc.iamClient, roleName, tc.trustPolicy)
	cleanup := func() { IAMDeleteRole(tc.iamClient, roleName) }
	return roleName, roleARN, cleanup
}

// rawJSONCall issues a request over the raw JSON-1.0 protocol and returns
// the HTTP status together with the decoded response body. It is used for
// operations the SDK client resolves onto sync-prefixed endpoints the
// local resolver cannot name, and for wire members absent from the typed
// SDK inputs.
func (tc *sfnTestContext) rawJSONCall(target string, payload map[string]interface{}) (int, map[string]interface{}, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequestWithContext(tc.ctx, "POST", tc.runner.endpoint, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")
	req.Header.Set("X-Amz-Target", target)
	resp, err := testHTTPClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	if derr := json.NewDecoder(resp.Body).Decode(&result); derr != nil {
		return resp.StatusCode, result, derr
	}
	return resp.StatusCode, result, nil
}

// rawTestState invokes the TestState operation over the raw JSON-1.0
// protocol: the SDK client resolves TestState onto a sync-prefixed
// endpoint the local resolver cannot name.
func (tc *sfnTestContext) rawTestState(payload map[string]interface{}) (map[string]interface{}, error) {
	status, result, err := tc.rawJSONCall("AWSStepFunctions.TestState", payload)
	if err != nil {
		return result, err
	}
	if status != http.StatusOK {
		return result, fmt.Errorf("status %d, want 200: %v", status, result)
	}
	return result, nil
}

// awaitTerminal polls DescribeExecution until the execution leaves the
// RUNNING state and returns the last description.
func (tc *sfnTestContext) awaitTerminal(executionArn string, interval time.Duration, attempts int) (*sfn.DescribeExecutionOutput, error) {
	var desc *sfn.DescribeExecutionOutput
	for i := 0; i < attempts; i++ {
		out, err := tc.client.DescribeExecution(tc.ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionArn),
		})
		if err != nil {
			return nil, err
		}
		desc = out
		if desc.Status != types.ExecutionStatusRunning {
			return desc, nil
		}
		time.Sleep(interval)
	}
	return desc, fmt.Errorf("execution did not finish")
}

// startExecution starts an execution with the given input; the name is
// optional and omitted when empty.
func (tc *sfnTestContext) startExecution(smArn, execName, input string) (string, error) {
	in := &sfn.StartExecutionInput{
		StateMachineArn: aws.String(smArn),
		Input:           aws.String(input),
	}
	if execName != "" {
		in.Name = aws.String(execName)
	}
	resp, err := tc.client.StartExecution(tc.ctx, in)
	if err != nil {
		return "", err
	}
	return aws.ToString(resp.ExecutionArn), nil
}

// runWithInput starts an execution and waits for a terminal state. On
// success it returns the execution ARN and the output payload; any other
// terminal state is an error carrying the execution Error and Cause.
func (tc *sfnTestContext) runWithInput(smArn, execName, input string) (string, string, error) {
	execArn, err := tc.startExecution(smArn, execName, input)
	if err != nil {
		return "", "", err
	}
	desc, err := tc.awaitTerminal(execArn, 200*time.Millisecond, 60)
	if err != nil {
		return execArn, "", err
	}
	if desc.Status != types.ExecutionStatusSucceeded {
		return execArn, "", fmt.Errorf("execution ended %s: %s %s",
			desc.Status, aws.ToString(desc.Error), aws.ToString(desc.Cause))
	}
	return execArn, aws.ToString(desc.Output), nil
}

// createSingleStateSM creates a standard state machine whose workflow is
// the single given state.
func (tc *sfnTestContext) createSingleStateSM(name string, state map[string]interface{}) (string, error) {
	def, err := json.Marshal(map[string]interface{}{
		"StartAt": "S",
		"States":  map[string]interface{}{"S": state},
	})
	if err != nil {
		return "", err
	}
	resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(name),
		Definition: aws.String(string(def)),
		RoleArn:    aws.String(tc.roleARN),
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(resp.StateMachineArn), nil
}

// firstMapRunFor polls ListMapRuns until the parent execution lists a Map
// Run and returns its ARN.
func (tc *sfnTestContext) firstMapRunFor(executionArn string) (string, error) {
	for i := 0; i < 20; i++ {
		time.Sleep(100 * time.Millisecond)
		page, err := tc.client.ListMapRuns(tc.ctx, &sfn.ListMapRunsInput{
			ExecutionArn: aws.String(executionArn),
		})
		if err != nil {
			return "", err
		}
		for _, mr := range page.MapRuns {
			if mr.MapRunArn != nil && *mr.MapRunArn != "" {
				return *mr.MapRunArn, nil
			}
		}
	}
	return "", fmt.Errorf("no map runs found after polling")
}

// createRoleBackedSM creates a dedicated IAM role and a standard state
// machine over the definition. The returned cleanup deletes both.
func (tc *sfnTestContext) createRoleBackedSM(namePrefix, definition string) (string, func(), error) {
	_, roleARN, roleCleanup := tc.createRoleForSM(namePrefix)
	resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(fmt.Sprintf("%s-%d", namePrefix, time.Now().UnixNano())),
		Definition: aws.String(definition),
		RoleArn:    aws.String(roleARN),
	})
	if err != nil {
		roleCleanup()
		return "", func() {}, err
	}
	arn := aws.ToString(resp.StateMachineArn)
	return arn, func() {
		tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(arn)})
		roleCleanup()
	}, nil
}

func (r *TestRunner) RunStepFunctionsTests() []TestResult {
	var results []TestResult

	tc, cleanup, err := newSFNTestContext(r)
	if err != nil {
		return append(results, TestResult{
			Service:  "stepfunctions",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to initialise test context: %v", err),
		})
	}
	defer cleanup()

	results = append(results, r.runSFNStateMachineTests(tc)...)
	results = append(results, r.runSFNExecutionTests(tc)...)
	results = append(results, r.runSFNActivityTests(tc)...)
	results = append(results, r.runSFNVersionTests(tc)...)
	results = append(results, r.runSFNAliasTests(tc)...)
	results = append(results, r.runSFNTagTests(tc)...)
	results = append(results, r.runSFNAdvancedTests(tc)...)
	results = append(results, r.runSFNEdgeTests(tc)...)
	results = append(results, r.runSFNWaitTests(tc)...)
	results = append(results, r.runSFNContractTests(tc)...)
	results = append(results, r.runSFNItemReaderTests(tc)...)
	results = append(results, r.runSFNMapFeatureTests(tc)...)
	results = append(results, r.runSFNIntrinsicTests(tc)...)

	return results
}
