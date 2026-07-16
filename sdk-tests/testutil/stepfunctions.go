package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
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
	tc.roleARN = fmt.Sprintf("arn:aws:iam::000000000000:role/%s", tc.roleName)

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
	defJSON, _ := json.Marshal(def)
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
	roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)
	IAMCreateRole(tc.iamClient, roleName, tc.trustPolicy)
	cleanup := func() { IAMDeleteRole(tc.iamClient, roleName) }
	return roleName, roleARN, cleanup
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

	return results
}
