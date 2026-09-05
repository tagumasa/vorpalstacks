package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"vorpalstacks-sdk-tests/config"
)

type schedTestContext struct {
	runner    *TestRunner
	client    *scheduler.Client
	iamClient *iam.Client
	kmsClient *kms.Client
	ctx       context.Context
	region    string
	accountID string
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
		kmsClient: kms.NewFromConfig(cfg),
		ctx:       ctx,
		region:    r.region,
		accountID: r.accountID,
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
	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.accountID, roleName)
	return roleName, roleARN
}

func (tc *schedTestContext) deleteIAMRole(roleName string) {
	IAMDeleteRole(tc.iamClient, roleName)
}

func (tc *schedTestContext) defaultTarget(roleARN string) *types.Target {
	return &types.Target{
		Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:TestFunction", tc.region, tc.accountID)),
		RoleArn: aws.String(roleARN),
		Input:   aws.String(`{"message":"test message"}`),
	}
}

func (tc *schedTestContext) lambdaARN() string {
	return fmt.Sprintf("arn:aws:lambda:%s:%s:function:TestFunction", tc.region, tc.accountID)
}

// createSchedule issues a plain CreateSchedule: the flexible time window is
// always OFF and no input-level members beyond name, expression and target
// are set. Tests exercising optional members or the operation's own
// semantics build their inputs inline.
func (tc *schedTestContext) createSchedule(name, expr string, target *types.Target) (*scheduler.CreateScheduleOutput, error) {
	resp, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
		Name:               aws.String(name),
		ScheduleExpression: aws.String(expr),
		Target:             target,
		FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		return nil, fmt.Errorf("CreateSchedule %s: %w", name, err)
	}
	return resp, nil
}

// cleanupSchedule deletes a schedule by name and ignores the error: it backs
// per-test defers, where the schedule may never have been created.
func (tc *schedTestContext) cleanupSchedule(name string) {
	tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(name)})
}

// getSchedule issues a plain GetSchedule addressed by name only.
func (tc *schedTestContext) getSchedule(name string) (*scheduler.GetScheduleOutput, error) {
	return tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{Name: aws.String(name)})
}

// createScheduleGroup issues a plain CreateScheduleGroup without tags.
func (tc *schedTestContext) createScheduleGroup(name string) (*scheduler.CreateScheduleGroupOutput, error) {
	resp, err := tc.client.CreateScheduleGroup(tc.ctx, &scheduler.CreateScheduleGroupInput{Name: aws.String(name)})
	if err != nil {
		return nil, fmt.Errorf("CreateScheduleGroup %s: %w", name, err)
	}
	return resp, nil
}

// cleanupScheduleGroup deletes a schedule group by name and ignores the
// error: it backs per-test defers.
func (tc *schedTestContext) cleanupScheduleGroup(name string) {
	tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(name)})
}

func (tc *schedTestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func (tc *schedTestContext) scheduleARN(scheduleName string) string {
	return fmt.Sprintf("arn:aws:scheduler:%s:%s:schedule/default/%s", tc.region, tc.accountID, scheduleName)
}
