package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"vorpalstacks-sdk-tests/config"
)

type schedTestContext struct {
	runner    *TestRunner
	client    *scheduler.Client
	iamClient *iam.Client
	ctx       context.Context
	region    string
}

func (r *TestRunner) RunSchedulerTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "scheduler",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	ctx := context.Background()
	tc := &schedTestContext{
		runner:    r,
		client:    scheduler.NewFromConfig(cfg),
		iamClient: iam.NewFromConfig(cfg),
		ctx:       ctx,
		region:    r.region,
	}

	results = append(results, tc.runScheduleTests()...)
	results = append(results, tc.runGroupTests()...)
	results = append(results, tc.runListTests()...)
	results = append(results, tc.runTagTests()...)
	results = append(results, tc.runEdgeTests()...)

	return results
}

func (tc *schedTestContext) createIAMRole() (string, string) {
	roleName := fmt.Sprintf("SchedRole-%d", time.Now().UnixNano())
	trustPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"scheduler.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
	IAMCreateRole(tc.iamClient, roleName, trustPolicy)
	roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)
	return roleName, roleARN
}

func (tc *schedTestContext) deleteIAMRole(roleName string) {
	IAMDeleteRole(tc.iamClient, roleName)
}

func (tc *schedTestContext) defaultTarget(roleARN string) *types.Target {
	return &types.Target{
		Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFunction", tc.region)),
		RoleArn: aws.String(roleARN),
		Input:   aws.String(`{"message":"test message"}`),
	}
}

func (tc *schedTestContext) lambdaARN() string {
	return fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFunction", tc.region)
}

func (tc *schedTestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func (tc *schedTestContext) scheduleARN(scheduleName string) string {
	return fmt.Sprintf("arn:aws:scheduler:%s:000000000000:schedule/default/%s", tc.region, scheduleName)
}
