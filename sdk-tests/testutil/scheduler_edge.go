package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

func (tc *schedTestContext) runEdgeTests() []TestResult {
	var results []TestResult

	rn, rARN := tc.createIAMRole()
	defer tc.deleteIAMRole(rn)

	results = append(results, tc.runner.RunTest("scheduler", "GetSchedule_NonExistent", func() error {
		_, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{
			Name: aws.String("nonexistent-schedule-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "DeleteSchedule_NonExistent", func() error {
		_, err := tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String("nonexistent-schedule-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_DuplicateName", func() error {
		dupName := tc.uniqueName("DupSched")
		dupRN, dupRARN := tc.createIAMRole()
		defer tc.deleteIAMRole(dupRN)

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(dupName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             tc.defaultTarget(dupRARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(dupName)})

		_, err = tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(dupName),
			ScheduleExpression: aws.String("rate(60 minutes)"),
			Target:             tc.defaultTarget(dupRARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ConflictException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "UpdateSchedule_NonExistent", func() error {
		_, err := tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String("nonexistent-schedule-xyz"),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_InvalidExpression", func() error {
		invName := tc.uniqueName("InvExprSched")
		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(invName),
			ScheduleExpression: aws.String("not-a-valid-expression"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
